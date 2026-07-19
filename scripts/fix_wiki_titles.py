#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
批量将 wiki 文档的 H1 标题改为「与文件名同名（去掉 NN- 序号前缀）」。

规则：
  1. 取文件名，去掉 .md 后缀；
  2. 去掉开头的纯数字序号前缀（^\\d{1,3}-），保留 M8-/G0-/Phase 等业务标识；
  3. 文件首行若为 '# '，则整行替换为 '# <新标题>'。

排除：
  - README.md           （仓库/板块索引，标题自由）
  - EXPORT-README.md     （Figma 工具生成的原型导出说明，标题固定为「原型导出说明 - XX端」）

用法：
  python3 scripts/fix_wiki_titles.py [wiki根目录，默认当前目录]
"""
import os
import re
import sys

# 排除的特殊文件名
SKIP_BASENAMES = {"README.md", "EXPORT-README.md"}

# 去掉纯数字序号前缀（目录排序用，如 01- / 00- / 12-）
NUM_PREFIX_RE = re.compile(r"^\d{1,3}-")


def title_from_filename(basename: str) -> str:
    name = basename[:-3] if basename.endswith(".md") else basename
    name = NUM_PREFIX_RE.sub("", name)
    return name.strip()


def process_file(path: str) -> bool:
    base = os.path.basename(path)
    if base in SKIP_BASENAMES:
        return False
    with open(path, encoding="utf-8") as f:
        lines = f.readlines()
    if not lines or not lines[0].startswith("# "):
        return False
    new_title = title_from_filename(base)
    new_first = "# " + new_title + "\n"
    if lines[0].rstrip("\n") == new_first.rstrip("\n"):
        return False  # 已经一致，跳过
    lines[0] = new_first
    with open(path, "w", encoding="utf-8") as f:
        f.writelines(lines)
    return True


def main() -> None:
    root = sys.argv[1] if len(sys.argv) > 1 else "."
    changed = []
    for dp, _dn, fnlist in os.walk(root):
        # 跳过隐藏目录（如 .obsidian / .kg 等生成或元数据目录）
        if "/." in dp or dp.startswith("."):
            continue
        for fn in fnlist:
            if not fn.endswith(".md"):
                continue
            p = os.path.join(dp, fn)
            if process_file(p):
                changed.append(os.path.normpath(p))
    print(f"已更新标题的文件数: {len(changed)}")
    for c in changed:
        print("  ", c)


if __name__ == "__main__":
    main()
