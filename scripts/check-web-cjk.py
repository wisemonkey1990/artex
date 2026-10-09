#!/usr/bin/env python3
"""检查静态导出的网页中是否意外包含中文字符。"""
import os
import re
import subprocess
import sys

# 说明。
# 说明。
# 说明。
HAN = re.compile(r"[㐀-鿿豈-﫿\U00020000-\U0002ffff]")

# 说明。
# 说明。
# 说明。
ALLOWED_HAN: set = set()


def find_out_dir() -> str:
    """确定要检查的 `out` 目录。

    如果指定了路径则使用该路径，否则检查 Git 仓库根目录下的 `web/out`。
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
            f"错误：待检查目录不存在 — {out_dir}\n"
            "请先运行 `npm run build:static` 生成静态导出文件。"
        )
        return 2

    checked = 0
    offenders = []
    for dirpath, _dirnames, filenames in os.walk(out_dir):
        for name in filenames:
            if not name.endswith(".html"):
                continue
            checked += 1
            path = os.path.join(dirpath, name)
            with open(path, encoding="utf-8", errors="ignore") as fh:
                text = fh.read()
            found = [ch for ch in HAN.findall(text) if ch not in ALLOWED_HAN]
            if found:
                rel = os.path.relpath(path, out_dir)
                kinds = "".join(sorted(set(found)))
                offenders.append((rel, len(found), kinds))

    print(f"已检查 {checked} 个 HTML 文件（目录：{out_dir}）")
    if offenders:
        print(f"有 {len(offenders)} 个文件包含中文字符：")
        for rel, count, kinds in offenders:
            preview = kinds if len(kinds) <= 30 else kinds[:30] + "…"
            print(f"  out/{rel}: 包含 {count} 个汉字 · 示例 {preview}")
        print(
            "\n请将其翻译为中文；如果是必须保留的上游原文汉字，请在"
            "请在 scripts/check-web-cjk.py 的 ALLOWED_HAN 中添加并说明原因。"
        )
        return 1
    if checked == 0:
        print("警告：没有检查到 HTML 文件，请确认静态构建成功。")
        return 2
    print("中文字符泄漏为 0：面向用户的 HTML 不包含中文字符。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
