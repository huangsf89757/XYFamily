#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
XYFamily Wiki 链接修复与交叉引用回写工具

功能：
1. 修复命名漂移 - 自动修复缺失 .md 后缀的链接 + 目录结构重命名
2. 交叉引用回写 - 将知识图谱潜在交叉引用建议写入各文档的"关联文档"章节

用法：
  python3 fix_links.py          # 仅预览
  python3 fix_links.py --apply  # 实际修改文件
"""

import os
import re
import json
import sys
from collections import defaultdict

WIKI_ROOT = os.path.join(os.path.dirname(os.path.abspath(__file__)), "wiki")
KG_FILE = os.path.join(WIKI_ROOT, ".kg", "knowledge_graph.json")

LINK_RE = re.compile(r"\[([^\]]*)\]\(([^)]+)\)")
VERSION_STRIP_RE = re.compile(r"[-_ ]?v?\d+\.\d+(\.\d+)?$", re.IGNORECASE)

# 已知的旧目录→新目录映射（实际重命名历史）
# 注意：norm_target 会将相对链接解析为以 wiki 根为基准的路径
DIR_MAPPINGS = [
    # 旧目录编号/名称                    → 新目录编号/名称
    ("00-项目总览",                       "01-项目总览"),
    ("01-立项与项目管理前置",              "01-项目总览"),
    # 接口文档 05→04，错误码 03→02
    ("01-项目总览/05-接口文档/03-错误码",  "04-接口文档/02-错误码"),
    ("01-项目总览/05-接口文档/05.01-接口总览/接口总览", "04-接口文档/接口文档"),
    ("01-项目总览/05-接口文档",            "04-接口文档"),
    # 技术架构目录重命名
    ("03-技术架构与方案设计/03.04-核心技术专项方案", "03-架构与方案设计"),
    ("03-技术架构与方案设计",              "03-架构与方案设计"),
    ("01-项目总览/03-技术架构与方案设计",   "03-架构与方案设计"),
    # PRD 目录重命名
    ("02-需求与产品设计/02.02-产品PRD",    "02-需求与产品设计/01-产品PRD"),
    ("01-项目总览/02.02-产品PRD",          "02-需求与产品设计/01-产品PRD"),
    ("01-项目总览/02-需求与产品设计",       "02-需求与产品设计"),
    # 开发规范目录重命名
    ("01-项目总览/04-开发规范与编码手册",   "01-项目总览/02-标准与规范/04-开发规范"),
    ("04-开发规范与编码手册",              "01-项目总览/02-标准与规范/04-开发规范"),
    # 接口文档中的 README → 主文档
    ("04-接口文档/README.md",              "04-接口文档/接口文档.md"),
    ("01-项目总览/05-接口文档/README.md",   "04-接口文档/接口文档.md"),
    # 旧里程碑路径
    ("01-项目总览/00-项目总览/00.03-项目里程碑/项目里程碑", "01-项目总览/03-里程碑/里程碑"),
    ("00-项目总览/00.03-项目里程碑/项目里程碑", "01-项目总览/03-里程碑/里程碑"),
    ("01-项目总览/03-里程碑/V1.0.0/项目里程碑", "01-项目总览/03-里程碑/里程碑"),
    # 接口文档目录中错误嵌套的需求PRD路径
    ("04-接口文档/01-标准接口/02-需求与产品设计", "02-需求与产品设计"),
    # 错误码目录下的接口规范
    ("04-接口文档/02-错误码/接口规范",      "04-接口文档/02-错误码/01-多租户底座"),
    # 表单中的PRD需求背景与目标（可能不存在，尝试匹配）
    ("02-需求与产品设计/01-产品PRD/01-多租户底座/01-需求背景与目标", "02-需求与产品设计/01-产品PRD/01-多租户底座/多租户底座"),
    # 原型与UI设计的效果图
    ("02-需求与产品设计/02-原型与UI设计/02-效果图", "02-需求与产品设计/02-原型与UI设计/原型与UI设计"),
    # 开发规范编号变化
    ("01-项目总览/02-标准与规范/04-开发规范/03-接口文档规范", "01-项目总览/02-标准与规范/04-开发规范/02-接口文档规范"),
]

# 链接目标明确不存在的写法（旧结构中的 README 应改为目录主文档）
README_REMAP = {
    # 旧路径中的 README → 新路径的主文档
    "02-需求与产品设计/02-原型与UI设计/README.md": "02-需求与产品设计/02-原型与UI设计/原型与UI设计.md",
}


def build_file_index():
    """构建 wiki 目录下所有 .md 文件的索引"""
    index = {}
    for root, dirs, files in os.walk(WIKI_ROOT):
        if ".kg" in root.split(os.sep):
            continue
        for fn in files:
            if not fn.lower().endswith(".md"):
                continue
            full = os.path.join(root, fn)
            rel = os.path.relpath(full, WIKI_ROOT).replace(os.sep, "/")
            index[rel.lower()] = rel
    return index


def strip_version(basename):
    name, _ = os.path.splitext(basename)
    return VERSION_STRIP_RE.sub("", name).lower()


def build_base_index(index):
    """构建版本剥离后的索引（同目录 + 全局）"""
    path_base = {}
    base = {}
    for rel in index.values():
        d = os.path.dirname(rel).lower()
        bn = rel.rsplit("/", 1)[-1]
        sv = strip_version(bn)
        path_base.setdefault((d, sv), rel)
        base.setdefault(sv, rel)
    return base, path_base


def resolve_norm(source_dir, link_target):
    """将相对链接解析为相对 wiki 根的规范化路径"""
    link_target = link_target.split("#")[0].split("?")[0].strip()
    if not link_target or link_target.startswith(("http://", "https://", "mailto:", "//")):
        return None
    if link_target.startswith("/"):
        return link_target.lstrip("/")
    return os.path.normpath(os.path.join(source_dir, link_target)).replace(os.sep, "/")


def find_correct_path(norm_path, index, path_base, base):
    """尝试找到规范路径的正确目标。返回 (correct_path, fix_type)。"""
    # 1. 精确匹配
    if norm_path.lower() in index:
        return index[norm_path.lower()], "exact"

    # 2. 明确的重映射表
    if norm_path in README_REMAP:
        remapped = README_REMAP[norm_path]
        if remapped.lower() in index:
            return index[remapped.lower()], "readme_remap"

    # 3. 添加 .md 后缀
    with_md = norm_path + ".md"
    if with_md.lower() in index:
        return index[with_md.lower()], "add_md"

    # 4. 目录映射修复
    for old_prefix, new_prefix in DIR_MAPPINGS:
        if norm_path.lower().startswith(old_prefix.lower()):
            new_path = new_prefix + norm_path[len(old_prefix):]
            if new_path.lower() in index:
                return index[new_path.lower()], "dir_map"
            new_with_md = new_path + ".md"
            if new_with_md.lower() in index:
                return index[new_with_md.lower()], "dir_map_md"

    # 5. 同目录下版本剥离匹配（如 File-V1.0.0 → File）
    bn = norm_path.rsplit("/", 1)[-1]
    stripped = strip_version(bn)
    tdir = os.path.dirname(norm_path).lower()
    resolved = path_base.get((tdir, stripped))
    if resolved:
        return resolved, "drift"

    # 6. 全局唯一名称匹配（作为最后手段）
    all_matches = [p for (d, ns), p in path_base.items() if ns == stripped]
    if len(all_matches) == 1:
        return all_matches[0], "drift_name_unique"

    return None, None


def compute_rel_path(source_file, target_file):
    """计算从 source 到 target 的正确相对路径"""
    source_dir = os.path.dirname(source_file)
    if not source_dir or source_dir == ".":
        return "./" + target_file
    sp = source_dir.split("/")
    tp = target_file.split("/")
    i = 0
    while i < len(sp) and i < len(tp) and sp[i] == tp[i]:
        i += 1
    up = len(sp) - i
    down = tp[i:]
    if up == 0:
        return "./" + "/".join(down)
    else:
        return "../" * up + "/".join(down)


def find_all_links_to_fix():
    """扫描所有 .md 文件，分类找出需要修复的链接。"""
    index = build_file_index()
    base, path_base = build_base_index(index)

    auto_fixes = []
    manual_review = []

    for root, dirs, files in os.walk(WIKI_ROOT):
        if ".kg" in root.split(os.sep):
            continue
        for fn in files:
            if not fn.lower().endswith(".md"):
                continue
            full_path = os.path.join(root, fn)
            rel_path = os.path.relpath(full_path, WIKI_ROOT).replace(os.sep, "/")
            source_dir = os.path.dirname(rel_path)

            with open(full_path, "r", encoding="utf-8") as f:
                lines = f.readlines()

            for line_no, line in enumerate(lines):
                for m in LINK_RE.finditer(line):
                    display = m.group(1)
                    raw_link = m.group(2)

                    if raw_link.startswith(("http://", "https://", "mailto:", "//")):
                        continue

                    norm = resolve_norm(source_dir, raw_link)
                    if norm is None:
                        continue

                    result, fix_type = find_correct_path(norm, index, path_base, base)

                    if fix_type is None:
                        # 完全无法解析的断链
                        manual_review.append((rel_path, line_no + 1, raw_link, norm, "broken", display))
                        continue
                    if fix_type == "exact":
                        continue  # 没问题

                    if result:
                        new_raw = compute_rel_path(source_file=rel_path, target_file=result)
                        if new_raw == raw_link:
                            continue

                        entry = (rel_path, line_no + 1, raw_link, new_raw, fix_type, display, norm, result)

                        # 安全自动修复的类型
                        if fix_type in ("add_md", "dir_map", "dir_map_md", "readme_remap", "drift"):
                            auto_fixes.append(entry)
                        else:
                            manual_review.append(entry)

    return auto_fixes, manual_review


def read_file(path):
    with open(path, "r", encoding="utf-8") as f:
        return f.read()


def write_file(path, content):
    with open(path, "w", encoding="utf-8") as f:
        f.write(content)


def apply_link_fixes(fixes):
    """应用链接修复到文件"""
    by_file = defaultdict(list)
    for fix in fixes:
        by_file[fix[0]].append(fix)

    changed = 0
    for rel_path, file_fixes in sorted(by_file.items()):
        full_path = os.path.join(WIKI_ROOT, rel_path)
        content = read_file(full_path)
        lines = content.splitlines(keepends=True)

        file_changed = False
        # 从后往前处理，避免行号偏移
        for fix in sorted(file_fixes, key=lambda x: -x[1]):
            rp, line_no, old_raw, new_raw, fix_type, display, norm, result = fix
            idx = line_no - 1
            if idx < len(lines):
                old_line = lines[idx]
                new_line = old_line.replace(f"]({old_raw})", f"]({new_raw})")
                if old_line != new_line:
                    lines[idx] = new_line
                    file_changed = True
                    print(f"  [{fix_type}] {rel_path}:{line_no} '{old_raw}' -> '{new_raw}'")

        if file_changed:
            write_file(full_path, "".join(lines))
            changed += 1

    return changed


def load_potential_links():
    if not os.path.exists(KG_FILE):
        return []
    with open(KG_FILE, "r", encoding="utf-8") as f:
        return json.load(f).get("potential_links", [])


def apply_cross_references(potential_links):
    """将潜在交叉引用写入各文档的"关联文档"章节"""
    by_source = defaultdict(list)
    for pl in potential_links:
        by_source[pl["source"]].append(pl)

    changed = 0
    skipped = 0

    for source_rel, suggestions in sorted(by_source.items()):
        full_path = os.path.join(WIKI_ROOT, source_rel)
        if not os.path.exists(full_path):
            skipped += 1
            continue

        content = read_file(full_path)

        # 收集已有链接用于去重
        existing_targets = set()
        for m in LINK_RE.finditer(content):
            raw = m.group(2)
            if raw.startswith(("http://", "https://", "mailto:")):
                continue
            existing_targets.add(raw.lower())

        # 过滤已存在链接
        new_suggestions = []
        for s in suggestions:
            tr = compute_rel_path(source_file=source_rel, target_file=s["target"])
            if tr.lower() in existing_targets:
                continue
            new_suggestions.append(s)

        if not new_suggestions:
            skipped += 1
            continue

        # 构建建议条目
        append_lines = []
        for s in sorted(new_suggestions, key=lambda x: -x["confidence"]):
            tr = compute_rel_path(source_file=source_rel, target_file=s["target"])
            display_name = s["target"].rsplit("/", 1)[-1].replace(".md", "")
            terms = "、".join(s.get("shared_terms", []))
            append_lines.append(
                f"- [{display_name}]({tr}) — 共享术语：{terms}（置信度 {s['confidence']}）"
            )

        new_entries = "\n".join(append_lines)
        prefix = "\n> 以下为知识图谱自动推荐的交叉引用，建议人工审阅确认后保留。\n\n"

        # 查找现有的"关联文档"章节
        related_match = re.search(r"^(#{1,6}\s*关联文档)", content, re.MULTILINE)

        if related_match:
            hdr = related_match.group(1)
            hdr_level = len(hdr.split()[0])
            hdr_end = related_match.end()

            # 找下一个同级或更高级标题
            rest = content[hdr_end:]
            next_hdr = re.search(rf"^#{{1,{hdr_level}}}\s", rest, re.MULTILINE)

            if next_hdr:
                insert_pos = hdr_end + next_hdr.start()
            else:
                insert_pos = len(content)

            new_content = (
                content[:insert_pos].rstrip()
                + "\n\n"
                + prefix
                + new_entries
                + "\n"
                + content[insert_pos:]
            )
        else:
            # 文件末尾新增
            new_content = content.rstrip() + "\n\n## 关联文档\n\n" + prefix + new_entries + "\n"

        if new_content != content:
            write_file(full_path, new_content)
            changed += 1
            print(f"  [关联文档] {source_rel} (+{len(new_suggestions)} 条建议)")

    return changed, skipped


def main():
    apply_flag = "--apply" in sys.argv

    print("=" * 60)
    print("XYFamily Wiki 链接修复与交叉引用回写工具")
    print("=" * 60)

    if not apply_flag:
        print("\n⚠️  预览模式，不会修改文件。加 --apply 执行实际修改。\n")

    # 第一步：扫描
    print("🔍 第一步：扫描命名漂移与断链...")
    auto_fixes, manual_review = find_all_links_to_fix()

    type_counts = defaultdict(int)
    for f in auto_fixes:
        type_counts[f[4]] += 1
    manual_counts = defaultdict(int)
    for f in manual_review:
        manual_counts[f[4]] += 1

    print(f"\n可自动修复: {len(auto_fixes)} 处")
    for ft, cnt in sorted(type_counts.items()):
        print(f"  {ft}: {cnt} 处")
    if manual_review:
        print(f"需人工审查: {len(manual_review)} 处")
        for ft, cnt in sorted(manual_counts.items()):
            print(f"  {ft}: {cnt} 处")

    if auto_fixes:
        if apply_flag:
            print("\n🔧 正在自动修复链接...")
            changed = apply_link_fixes(auto_fixes)
            print(f"✅ 已修改 {changed} 个文件")
        else:
            print("\n📋 预览（前 30 条）：")
            for fix in auto_fixes[:30]:
                rp, line_no, old_raw, new_raw, ft, display, norm, result = fix
                print(f"  [{ft}] {rp}:{line_no}")
                print(f"    旧: {old_raw}")
                print(f"    新: {new_raw}")
            if len(auto_fixes) > 30:
                print(f"  ...共 {len(auto_fixes)} 条")

    if manual_review:
        # 写入审查报告
        report_path = os.path.join(
            os.path.dirname(os.path.abspath(__file__)), "manual_review_needed.md"
        )
        with open(report_path, "w", encoding="utf-8") as rf:
            rf.write("# 需人工审查的链接修复项\n\n")
            rf.write(f"共 {len(manual_review)} 条\n\n")
            rf.write("| 类型 | 文件 | 行号 | 旧链接 | 规范路径 |\n")
            rf.write("|------|------|------|--------|----------|\n")
            for fix in sorted(manual_review, key=lambda x: (x[4], x[0], x[1])):
                if len(fix) == 6:
                    rp, line_no, raw_link, info, ft, display = fix
                    rf.write(
                        f"| {ft} | `{rp}` | {line_no} | `{raw_link}` | `{info}` |\n"
                    )
                else:
                    rp, line_no, raw_link, new_raw, ft, display, norm, result = fix
                    rf.write(
                        f"| {ft} | `{rp}` | {line_no} | `{raw_link}` | `{norm}` → `{result}` |\n"
                    )
        print(f"\n📝 需人工审查的项目已写入: {report_path}")

    # 第二步：交叉引用回写
    print("\n🔗 第二步：交叉引用建议回写...")
    potential_links = load_potential_links()
    print(f"知识图谱中共有 {len(potential_links)} 条潜在交叉引用")

    if potential_links:
        if apply_flag:
            changed, skipped = apply_cross_references(potential_links)
            print(f"✅ 已向 {changed} 个文档添加交叉引用建议（跳过 {skipped} 个）")
        else:
            by_source = defaultdict(list)
            for pl in potential_links:
                by_source[pl["source"]].append(pl)
            print(f"\n📋 预览（按源文件分布，前 15 个）：")
            for src, sugs in sorted(by_source.items(), key=lambda x: -len(x[1]))[:15]:
                print(f"  {src}: {len(sugs)} 条建议")

    print("\n" + "=" * 60)
    if not apply_flag:
        print("💡 使用 python3 fix_links.py --apply 执行实际修改")
    else:
        print("🎉 全部完成！请 git diff 检查修改内容。")
    print("=" * 60)


if __name__ == "__main__":
    main()
