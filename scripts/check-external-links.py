#!/usr/bin/env python3
"""检查受跟踪的 Markdown 文档引用的外部链接（http·https）是否仍然有效。

姐妹脚本 `check-doc-links.py` 只看仓库内部链接·锚点，按设计跳过外部 URL
（依赖网络会导致结果不稳定，不适合放进合并关卡）。这个脚本承担那部分
"另外检查"的工作。README 开头·防御指南·detections README 里给访问者和
防御者指引"请访问这里"的外部链接（事件报告渠道 boho.or.kr·privacy.go.kr·
pipc.go.kr·fsec.or.kr，CISA KEV，OWASP·SigmaHQ·Suricata·MITRE·MISP，
原始演示 artex-demo.vercel.app，GitHub 徽章等）一旦变更·迁移·下线就会悄悄
失效，本脚本负责定期·手动地抓出这些问题。

选取待检查 URL 的规则
- 遍历所有受跟踪的 `.md`，但跳过围栏代码块（``` 或 ~~~）和行内代码片段
  （`` `...` ``）里的 URL。那些 URL 是命令示例·配置值·引用的外部资产
  （例如旁敲验证文档里的 `id.redhaze.top`），不是"供读者点击的参考链接"。
- 从 Markdown 链接·图片（`[text](url)`·`![alt](url)`）、自动链接（`<url>`），
  以及正文中裸露的 URL 里收集 http·https 地址。
- 排除预留·占位主机：localhost·私有/回环 IP（127.·10.·192.168.·
  169.254.·172.16~31.·0.0.0.0·::1）、RFC 2606/6761 预留域名
  （example.com/org/net/edu·`*.example.*`·`.test`·`.invalid`·`.local`·`.tld`）、
  没有点号不是 FQDN 的名字（例如 `target`）。

存活判定方式（把 MAINTAINING.md 8.2 的经验落实成代码）
部分国内公共·安全机构网站会拒绝 HEAD 请求·默认 User-Agent，或做多次重定向，
简单检测会把正常链接误判为失效。所以本脚本用**浏览器 User-Agent、GET 请求、
跟随重定向**来检测，把状态分三类：
- OK：最终状态是 2xx·3xx，链接有效。
- RESTRICTED：401·403·405·429。主机本身是活的，只是检测方式被服务器访问
  策略（防爬虫·拒绝该方法·限流）挡住了，不算失效链接，会报告但不计入失败。
- DOWN：404·410·5xx（重试后仍是）·DNS/连接/超时/SSL 错误（重试后仍是）。
  很可能是真的失效了。

网络错误·5xx·429 会在小幅延迟后重试，用来区分瞬时抖动和真正故障。
只使用标准库。本脚本只会对**仓库公开文档已经引用的**链接发 GET 请求来确认
存活，不会扫描或探测任何其他目标。

已知例外（allowlist）
如果 `scripts/external-links-allowlist.txt` 中列出的 URL 检测为 DOWN，会单独
归类为"ALLOWED"，不计入 strict 模式的退出码。这是为了透明记录那些我们不拥有、
无法修复、从上游原文保留文件（例如上游 CHANGELOG.zh.md）继承来的失效链接，
让周期性的 strict 检查不会因为这一个问题永远变红，只在**出现新的失效链接**时
才变红。

运行方式
- 默认（仅报告，退出码 0）：`python3 -I scripts/check-external-links.py`
- 非合并关卡的周期/发布检查（有新的失效链接就变红）：`--strict`（出现任何不在
  allowlist 中的 DOWN 就退出码 1；RESTRICTED·ALLOWED 在 strict 下也不算失败）。
- 不发网络请求，只预览提取到的 URL 集合：`--list`（打印待检查的 URL 及出处，不发请求）。
"""
import argparse
import http.client
import os
import re
import socket
import ssl
import subprocess
import sys
import time
import urllib.error
import urllib.request
from collections import defaultdict
from urllib.parse import urlsplit

# [text](url) 正文链接（非图片）与 ![alt](url) 图片。路径取到空白或右括号之前。
MD_LINK = re.compile(r"(?<!\!)\[[^\]]*\]\(([^)\s]+)")
MD_IMAGE = re.compile(r"!\[[^\]]*\]\(([^)\s]+)")
# <https://...> 自动链接。
AUTOLINK = re.compile(r"<(https?://[^>\s]+)>")
# 正文中裸露的 URL。末尾常见的标点会在后面去掉。
BARE_URL = re.compile(r"https?://[^\s)>\]\"'`]+")
# 行内代码片段 `...`（单反引号）。提取前先替换成空格。
INLINE_CODE = re.compile(r"`[^`]*`")
# 正文 URL 末尾常见的标点符号。
TRAILING_PUNCT = ".,;:!?\"'»)]}>"

# IPv4 私有/回环/链路本地网段（172.16~31. 单独一条规则）。
PRIVATE_IPV4 = re.compile(r"^(127\.|10\.|192\.168\.|169\.254\.|0\.0\.0\.0$)")
PRIVATE_IPV4_172 = re.compile(r"^172\.(1[6-9]|2\d|3[01])\.")

BROWSER_UA = (
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "
    "AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0 Safari/537.36"
)

# 主机本身是活的，只是检测方式被服务器策略挡住的状态码（不算失效）。
RESTRICTED_CODES = {401, 403, 405, 429}

ALLOWLIST_PATH = os.path.join("scripts", "external-links-allowlist.txt")


def load_allowlist(root: str) -> set:
    """读取我们无法修复的已知 DOWN URL 集合（不存在则返回空集合）。

    每行一个 URL，`#` 之后是注释，空行·纯注释行忽略。
    """
    path = os.path.join(root, ALLOWLIST_PATH)
    allow = set()
    if not os.path.exists(path):
        return allow
    with open(path, encoding="utf-8") as fh:
        for line in fh:
            line = line.split("#", 1)[0].strip()
            if line:
                allow.add(line)
    return allow


def strip_code(line: str) -> str:
    """去掉围栏之外一行里的行内代码片段（其中的 URL 不是引用）。"""
    return INLINE_CODE.sub(" ", line)


def is_checkable(url: str) -> bool:
    """过滤掉预留·占位主机，只留下真正值得检查的外部 URL。"""
    parts = urlsplit(url)
    if parts.scheme not in ("http", "https"):
        return False
    host = parts.hostname
    if not host:
        return False
    host = host.lower()
    if host == "localhost" or host.endswith(".localhost"):
        return False
    if host == "::1":
        return False
    if PRIVATE_IPV4.match(host) or PRIVATE_IPV4_172.match(host):
        return False
    if host in ("example.com", "example.org", "example.net", "example.edu"):
        return False
    if host.endswith((".example.com", ".example.org", ".example.net", ".example")):
        return False
    if host.endswith((".test", ".invalid", ".local", ".localdomain", ".tld")):
        return False
    if "." not in host:  # 不是 FQDN 的裸名字（例如 target）
        return False
    return True


def extract(root: str, md_files: list) -> dict:
    """从受跟踪的 .md 中构建 待检查外部 URL → [(文件, 行号), ...] 的字典。"""
    found = defaultdict(list)
    for md in md_files:
        with open(os.path.join(root, md), encoding="utf-8", errors="ignore") as fh:
            lines = fh.readlines()
        in_fence = False
        for lineno, raw in enumerate(lines, 1):
            stripped = raw.lstrip()
            if stripped.startswith("```") or stripped.startswith("~~~"):
                in_fence = not in_fence
                continue
            if in_fence:
                continue
            line = strip_code(raw)
            candidates = []
            for pattern in (MD_LINK, MD_IMAGE, AUTOLINK):
                candidates.extend(m.group(1) for m in pattern.finditer(line))
            for m in BARE_URL.finditer(line):
                candidates.append(m.group(0).rstrip(TRAILING_PUNCT))
            for url in candidates:
                url = url.strip().rstrip(TRAILING_PUNCT)
                if is_checkable(url):
                    where = (md, lineno)
                    if where not in found[url]:
                        found[url].append(where)
    return found


def probe(url: str, timeout: float, retries: int, delay: float):
    """用浏览器 UA 对 URL 发 GET 请求，返回（分类, 详情）。

    分类为 "OK" · "RESTRICTED" · "DOWN"。网络错误·5xx·429 会重试。
    """
    req = urllib.request.Request(
        url,
        method="GET",
        headers={
            "User-Agent": BROWSER_UA,
            "Accept": "*/*",
            "Accept-Language": "ko,en;q=0.8",
        },
    )
    last = ""
    for attempt in range(retries + 1):
        try:
            with urllib.request.urlopen(req, timeout=timeout) as resp:
                code = resp.getcode()
                if code in RESTRICTED_CODES:
                    return "RESTRICTED", f"HTTP {code}"
                return "OK", f"HTTP {code}"
        except urllib.error.HTTPError as exc:
            code = exc.code
            if code in RESTRICTED_CODES:
                return "RESTRICTED", f"HTTP {code}"
            if code >= 500 or code == 408:
                last = f"HTTP {code}"  # 可能是暂时的，重试
            else:
                return "DOWN", f"HTTP {code}"  # 404·410 等是确定的，无需重试
        except (
            urllib.error.URLError,
            http.client.HTTPException,
            ssl.SSLError,
            socket.timeout,
            ConnectionError,
            TimeoutError,
            OSError,
        ) as exc:
            reason = getattr(exc, "reason", exc)
            last = f"{type(exc).__name__}: {reason}"
        if attempt < retries:
            time.sleep(delay * (attempt + 1))
    return "DOWN", last


def main() -> int:
    ap = argparse.ArgumentParser(description="检查受跟踪 Markdown 中外部链接的存活状态")
    ap.add_argument("--strict", action="store_true",
                    help="只要有一个 DOWN 就退出码 1（用于周期/发布检查）")
    ap.add_argument("--list", action="store_true",
                    help="不发网络请求，只打印待检查的 URL 及出处")
    ap.add_argument("--timeout", type=float, default=15.0, help="请求超时（秒）")
    ap.add_argument("--retries", type=int, default=2, help="网络错误·5xx·429 的重试次数")
    ap.add_argument("--delay", type=float, default=0.5, help="请求之间·重试之间的基础延迟（秒）")
    args = ap.parse_args()

    root = subprocess.check_output(
        ["git", "rev-parse", "--show-toplevel"], text=True
    ).strip()
    tracked = subprocess.check_output(["git", "ls-files"], cwd=root, text=True).splitlines()
    md_files = [f for f in tracked if f.endswith(".md")]

    allow = load_allowlist(root)
    found = extract(root, md_files)
    urls = sorted(found)
    print(f"从受跟踪 Markdown {len(md_files)} 个文件中收集到待检查外部 URL {len(urls)} 个")

    if args.list:
        for url in urls:
            where = ", ".join(f"{f}:{ln}" for f, ln in found[url])
            tag = " [allowlist]" if url in allow else ""
            print(f"  {url}{tag}  ({where})")
        return 0

    results = {"OK": [], "RESTRICTED": [], "ALLOWED": [], "DOWN": []}
    for i, url in enumerate(urls):
        if i:
            time.sleep(args.delay)  # 请求之间的小幅延迟（规避沙箱限流）
        verdict, detail = probe(url, args.timeout, args.retries, args.delay)
        if verdict == "DOWN" and url in allow:
            verdict = "ALLOWED"  # 无法修复的已知上游继承 DOWN
        results[verdict].append((url, detail))
        print(f"  [{verdict:10}] {url} — {detail}")

    print(
        f"\n汇总：OK {len(results['OK'])} · RESTRICTED {len(results['RESTRICTED'])} · "
        f"ALLOWED {len(results['ALLOWED'])} · DOWN {len(results['DOWN'])}"
    )
    if results["RESTRICTED"]:
        print("RESTRICTED（主机存活·检测方式被挡住，不是失效链接）：")
        for url, detail in results["RESTRICTED"]:
            print(f"  {url} — {detail}")
    if results["ALLOWED"]:
        print("ALLOWED（allowlist 中记录的已知 DOWN —— 我们无法修复，strict 模式不计入）：")
        for url, detail in results["ALLOWED"]:
            where = ", ".join(f"{f}:{ln}" for f, ln in found[url])
            print(f"  {url} — {detail}  （引用处：{where}）")
    if results["DOWN"]:
        print("DOWN（很可能已失效 —— 需要确认）：")
        for url, detail in results["DOWN"]:
            where = ", ".join(f"{f}:{ln}" for f, ln in found[url])
            print(f"  {url} — {detail}  （引用处：{where}）")

    if args.strict and results["DOWN"]:
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
