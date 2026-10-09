#!/usr/bin/env python3
"""检查静态导出的网页是否包含韩文字符。"""
import os
import re
import subprocess
import sys

HANGUL = re.compile(r"[\u1100-\u11ff\u3130-\u318f\uac00-\ud7af]")


def find_out_dir() -> str:
    """确定要检查的 out 目录。"""
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
            found = HANGUL.findall(text)
            if found:
                rel = os.path.relpath(path, out_dir)
                kinds = "".join(sorted(set(found)))
                offenders.append((rel, len(found), kinds))

    print(f"已检查 {checked} 个 HTML 文件（目录：{out_dir}）")
    if offenders:
        print(f"有 {len(offenders)} 个文件包含韩文字符：")
        for rel, count, kinds in offenders:
            preview = kinds if len(kinds) <= 30 else kinds[:30] + "…"
            print(f"  out/{rel}: 包含 {count} 个韩文字符 · 示例 {preview}")
        return 1
    if checked == 0:
        print("警告：没有检查到 HTML 文件，请确认静态构建成功。")
        return 2
    print("韩文字符残留为 0：面向用户的 HTML 中未发现韩文。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
