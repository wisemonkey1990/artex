#!/usr/bin/env python3
"""检查静态导出的 Web UI（`web/out`）的 HTML 中是否混入了韩文（谚文）。

本 fork 的核心成果与差异点是"用户可见界面没有韩文"（对应 `out` HTML 谚文
零出现这一属性）。如果只靠人工检查来维护这一属性，贡献者硬编码韩文字符串、
或翻译遗漏导致的回归，会一直悄悄存在，直到下一个人用眼睛发现为止。本脚本
在合并关卡上拦住这类回归。

检查方式
--------
遍历 `web/out` 下所有 `.html` 文件的**原始内容**（不只是渲染后的标记，也包括
内联脚本·React Server Components 的 flight 数据），只要出现一个谚文字符就
以退出码 1 结束。之所以看原始内容，是因为用户可见文本会以 flight 数据
（JSON 形式的字符串）的形式随 HTML 一起下发。只剥标签的话会漏掉这条泄露路径。

检查范围之外（如实说明局限）
----------------------------
这个关卡只看静态导出的产物。因此**MOCK 模式（`NEXT_PUBLIC_MOCK=1`，公开演示
使用的模式）下浏览器运行时绘制的内容**检查不到。`web/src/lib/mock/data.ts`·
`handler.ts` 以及各页面的硬编码样例数据，是在 hydration 之后由 `useEffect`
调 `mockHandle()` 取回并以 React state 渲染的，所以只存在于运行时 JS bundle
里，不在 `out` HTML 中。也就是说经这条路径混入的韩文能绕过这个关卡。

这一面靠代码评审来守：渲染相关的源文件（`web/src` 下的 `.tsx`·客户端数据）
中出现的韩文，必须属于以下情形之一：代码注释、被 `t.has(...)` 挡住实际不会
触发的 i18n 兜底值、后端线上格式的 token（如聊天提及分类标签、`[模型]`
决策来源哨兵值）、仅在 locale 为 `ko` 时才会显示的分支文案，或者为保留性能
而原样镜像的智能体（段 [A]）提示词。不属于以上任何一种、"就是会显示在界面上
的韩文"，需要翻译成中文。

允许·检测范围
--------------
- **汉字允许**。谚文音节块与汉字在 Unicode 中不重叠，只检测谚文音节块，
  汉字会自动放行。
- **只检测谚文（韩文字母）**：谚文音节块（`가`–`힣`，U+AC00–U+D7A3）与
  谚文字母区（U+1100–U+11FF、U+3130–U+318F）。平假名·片假名·全角标点
  不是"谚文"，不在这个关卡的检测范围内（把范围收窄到谚文，是为了避免原文
  保留的代码块或证据引用混进来造成误报）。

需要为上游保留而故意留下的韩文
------------------------------
目前 `out` HTML 中没有任何谚文字符，允许列表为空。将来如果因为保留性能，
必须故意在界面上留一些原文（韩文），把对应字符加进下面的 `ALLOWED_HANGUL`，
并用注释写明原因。这样既能保持关卡的严格性，又能窄窄地放过有意为之的例外。

只使用标准库，不访问网络。
运行：`python3 -I scripts/check-web-hangul.py [out 目录]`
（默认是仓库根目录的 `web/out`。CI 在 `web` 工作目录下用
`python3 -I ../scripts/check-web-hangul.py` 调用。）
"""
import os
import re
import subprocess
import sys

# 仅谚文（韩文字母）区块。汉字·平假名·片假名·全角标点不在其中。
# 范围：谚文音节（U+AC00–U+D7A3）、谚文字母 Jamo（U+1100–U+11FF）、
#       兼容 Jamo（U+3130–U+318F）。
HANGUL = re.compile(r"[가-힣ᄀ-ᇿ㄰-㆏]")

# 界面上故意保留的韩文（例如语言切换器里的原文标签）。目前为空 ——
# 留空意味着任何谚文字符都会被拦住。添加例外时请连同字符和理由一起写明。
# 例：ALLOWED_HANGUL = {"한", "글"}  # 语言切换器的韩语标签（对照上游保留）
ALLOWED_HANGUL: set = set()


def find_out_dir() -> str:
    """确定要检查的 `out` 目录。

    传了参数就用参数指定的路径，否则用 git 仓库根目录下的 `web/out`。
    """
    if len(sys.argv) > 1:
        return os.path.abspath(sys.argv[1])
    root = subprocess.check_output(
        ["git", "rev-parse", "--show-toplevel"], text=True
    ).strip()
    return os.path.join(root, "web", "out")


def main() -> int:
    out_dir = find_out_dir()
    if not os.path.isdir(out_dir):
        print(
            f"错误：待检查的目录不存在 —— {out_dir}\n"
            "请先用 `npm run build:static` 生成静态导出产物后再运行。"
        )
        return 2

    checked = 0
    offenders = []  # (相对路径, 谚文字符数, 预览)
    for dirpath, _dirnames, filenames in os.walk(out_dir):
        for name in filenames:
            if not name.endswith(".html"):
                continue
            checked += 1
            path = os.path.join(dirpath, name)
            with open(path, encoding="utf-8", errors="ignore") as fh:
                text = fh.read()
            found = [ch for ch in HANGUL.findall(text) if ch not in ALLOWED_HANGUL]
            if found:
                rel = os.path.relpath(path, out_dir)
                kinds = "".join(sorted(set(found)))
                offenders.append((rel, len(found), kinds))

    print(f"已检查 {checked} 个 HTML 文件（目标：{out_dir}）")
    if offenders:
        print(f"发现 {len(offenders)} 个文件混入韩文（谚文）—— 用户可见 HTML 中含有韩文：")
        for rel, count, kinds in offenders:
            preview = kinds if len(kinds) <= 30 else kinds[:30] + "…"
            print(f"  out/{rel}：韩文 {count} 个 · 字符 {preview}")
        print(
            "\n请改写为中文，或者如果确实需要为上游保留而故意留下韩文，"
            "请在 scripts/check-web-hangul.py 的 ALLOWED_HANGUL 中添加并写明理由。"
        )
        return 1
    if checked == 0:
        print("警告：没有检查到任何 HTML 文件 —— 请确认构建是否成功。")
        return 2
    print("韩文混入数：0 —— 用户可见 HTML 全部为中文/非韩文。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
