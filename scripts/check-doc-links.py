#!/usr/bin/env python3
"""检查仓库 Markdown 文档中的内部链接、图片路径和标题锚点。仅使用 Python 标准库。"""
import os
import re
import subprocess
import sys
import unicodedata
from collections import defaultdict

# 说明。
MD_LINK = re.compile(r"(?<!\!)\[[^\]]*\]\(([^)\s]+)")
# 说明。
MD_IMAGE = re.compile(r"!\[[^\]]*\]\(([^)\s]+)")
# 说明。
HTML_IMAGE = re.compile(r"<img[^>]*\bsrc=[\"']([^\"']+)[\"']", re.IGNORECASE)

# 说明。
EXTERNAL_PREFIXES = ("http://", "https://", "mailto:", "tel:", "data:")

ATX_HEADING = re.compile(r"^(#{1,6})\s+(.*?)\s*#*\s*$")
# 说明。
MD_INLINE_LINK = re.compile(r"\[([^\]]*)\]\([^)]*\)")


def is_external(target: str) -> bool:
    return target.startswith(EXTERNAL_PREFIXES)


def _is_word_char(ch: str) -> bool:
    """判断字符是否符合 GitHub 锚点 slug 保留的 `\\p{Word}` 规则。

    Ruby 正则 `\\p{Word}` = Letter(L*) + Mark(M*) + Decimal_Number(Nd)
    + Connector_Punctuation(Pc) + Join_Control(U+200C·U+200D)，包括韩文、汉字、
    下划线和表情符号变体选择符(Mn)等会保留。
    """
    if ch == "_":
        return True
    if ch in ("‌", "‍"):
        return True
    category = unicodedata.category(ch)
    if category[0] in ("L", "M"):
        return True
    return category in ("Nd", "Pc")


def slugify(text: str) -> str:
    """将标题文本转换为 GitHub 锚点 slug（不处理重复项后缀）。"""
    text = text.strip().lower()
    out = []
    for ch in text:
        if ch == " ":
            out.append("-")
        elif ch == "-" or _is_word_char(ch):
            out.append(ch)
    return "".join(out)


def heading_slugs(root: str, md: str) -> set:
    """收集文件中所有 ATX 标题对应的 GitHub 锚点 slug。

    slug 重复时，按出现顺序添加 `-1`、`-2` 后缀，与 GitHub 行为一致。
    """
    slugs = set()
    counts = defaultdict(int)
    in_fence = False
    with open(os.path.join(root, md), encoding="utf-8", errors="ignore") as fh:
        for line in fh:
            stripped = line.lstrip()
            if stripped.startswith("```") or stripped.startswith("~~~"):
                in_fence = not in_fence
                continue
            if in_fence:
                continue
            match = ATX_HEADING.match(line.rstrip("\n"))
            if not match:
                continue
            raw = match.group(2)
            raw = MD_INLINE_LINK.sub(r"\1", raw)
            raw = raw.replace("`", "").replace("*", "").replace("_", "")
            base = slugify(raw)
            seen = counts[base]
            counts[base] += 1
            slugs.add(base if seen == 0 else f"{base}-{seen}")
    return slugs


def main() -> int:
    root = subprocess.check_output(
        ["git", "rev-parse", "--show-toplevel"], text=True
    ).strip()
    tracked = subprocess.check_output(
        ["git", "ls-files"], cwd=root, text=True
    ).splitlines()
    md_files = [f for f in tracked if f.endswith(".md")]

    slug_cache: dict = {}

    def slugs_of(md_path: str) -> set:
        if md_path not in slug_cache:
            slug_cache[md_path] = heading_slugs(root, md_path)
        return slug_cache[md_path]

    checked_files = 0
    checked_anchors = 0
    broken_files = []
    broken_anchors = []
    for md in md_files:
        with open(os.path.join(root, md), encoding="utf-8", errors="ignore") as fh:
            lines = fh.readlines()
        base_dir = os.path.dirname(md)
        in_fence = False
        for lineno, line in enumerate(lines, 1):
            stripped = line.lstrip()
            if stripped.startswith("```") or stripped.startswith("~~~"):
                in_fence = not in_fence
                continue
            if in_fence:
                continue
            for pattern in (MD_LINK, MD_IMAGE, HTML_IMAGE):
                for match in pattern.finditer(line):
                    target = match.group(1).strip()
                    if is_external(target):
                        continue

                    if target.startswith("#"):
                        file_part, anchor = "", target[1:]
                    elif "#" in target:
                        file_part, anchor = target.split("#", 1)
                    else:
                        file_part, anchor = target, ""
                    file_part = file_part.split("?", 1)[0]
                    anchor = anchor.strip()

                    # 说明。
                    # 说明。
                    target_md = None
                    if file_part:
                        checked_files += 1
                        resolved = os.path.normpath(
                            os.path.join(root, base_dir, file_part)
                        )
                        if not os.path.exists(resolved):
                            broken_files.append((md, lineno, target))
                            continue
                        if file_part.endswith(".md"):
                            target_md = os.path.normpath(
                                os.path.join(base_dir, file_part)
                            )
                    else:
                        target_md = md

                    # 说明。
                    # 说明。
                    if anchor and target_md is not None and pattern is MD_LINK:
                        checked_anchors += 1
                        if anchor not in slugs_of(target_md):
                            broken_anchors.append((md, lineno, target, target_md))

    print(
        f"已检查 {len(md_files)} 个 Markdown 文件、{checked_files} 个文件引用、"
        f"{checked_anchors} 个锚点引用"
    )
    if broken_files:
        print(f"发现 {len(broken_files)} 个失效文件引用，目标文件不存在：")
        for md, lineno, target in broken_files:
            print(f"  {md}:{lineno} -> {target}")
    if broken_anchors:
        print(f"发现 {len(broken_anchors)} 个失效锚点，目标文档中没有对应标题：")
        for md, lineno, target, target_md in broken_anchors:
            print(f"  {md}:{lineno} -> {target}  (目标：{target_md})")
    if broken_files or broken_anchors:
        return 1
    print("失效引用为 0：所有内部链接、图片和锚点均指向有效目标。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
