#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
XYFamily Wiki 知识图谱构建与维护工具（零依赖、纯本地、只读）。

落地指南核心能力映射：
  1. 原生索引 / 零数据搬运：直接遍历 Vault 内的 .md 文件，解析相对链接与元数据表格，
     不导出到向量库、不修改任何源文件。
  2. Delta 追踪：基于文件内容 sha256 指纹，仅处理新增/修改文档，state.json 记录上次指纹。
  3. cross-linker + 置信度评分：扫描标题/标签/术语重叠，发现"潜在交叉引用"，给出 0~1 置信度，
     不自动写入源文件（保持只读，避免破坏 Obsidian/Markdown 结构）。
  4. 安全隔离：解析"文档密级"与敏感关键词（key/密码/secret/PII 等），生成 visibility 标记；
     safe_graph 自动过滤机密/PII 节点与跨密级边，供对外查询。

产物（输出到 wiki/.kg/）：
  - knowledge_graph.json  全量图谱（含敏感节点，仅内部使用）
  - safe_graph.json       过滤后的安全图谱（剔除机密/PII）
  - GRAPH.md              人类可读报告（结构化，安全视图）
  - graph.html            可交互可视化（含潜在链接，安全视图）
  - state.json            指纹状态，供下次增量 Delta 追踪
"""

import hashlib
import json
import os
import re
from collections import defaultdict
from datetime import datetime, timezone

WIKI_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SCAN_ROOT = os.path.join(WIKI_ROOT, "wiki")          # 仅扫描 wiki/ 子树
KG_DIR = os.path.join(SCAN_ROOT, ".kg")
STATE_FILE = os.path.join(KG_DIR, "state.json")

# 指南关键词表：识别"机密/PII/密钥"相关内容
SENSITIVE_KEYWORDS = [
    "密钥", "密码", "secret", "token 密钥", "private key", "access key", "ak/sk",
    "数据库账号", "数据库密码", "明文", "客户隐私", "身份证", "手机号", "邮箱",
    "access_token", "refresh_token", "jwt secret", "加密盐", "salt",
]
# 文档密级 -> 可见性级别（数字越大越敏感）
CLASSIFICATION_LEVEL = {"公开": 0, "内部": 1, "机密": 2}

# 业务术语词典（用于 cross-linker 的术语重叠评分）
TERM_DICT = [
    "多租户", "租户隔离", "rbac", "权限", "角色", "鉴权", "认证", "token", "jwt",
    "审计", "组织架构", "团队", "小组", "账号", "会话", "缓存", "中间件", "链路",
    "容灾", "多活", "高可用", "可观测", "合规", "数据脱敏", "等保", "隐私", "安全",
    "部署", "k8s", "docker", "灰度", "故障", "监控", "告警", "数据库", "契约", "接口",
]

LINK_RE = re.compile(r"\[[^\]]+\]\(([^)]+)\)")
WIKILINK_RE = re.compile(r"\[\[([^\]]+)\]\]")
HEADING_RE = re.compile(r"^#{1,6}\s+(.*)$")
# 去掉文件名里的版本后缀（如 -V1.0.0 / -v2.1.0），用于命名漂移匹配
VERSION_STRIP_RE = re.compile(r"[-_ ]?v?\d+\.\d+(\.\d+)?$", re.IGNORECASE)


def strip_version(basename):
    """去掉文件名里的版本后缀并转小写，用于命名漂移匹配。"""
    name, _ = os.path.splitext(basename)
    return VERSION_STRIP_RE.sub("", name).lower()
TABLE_CELL_RE = re.compile(r"\|([^|]+)\|")
FRONT_MATTER_RE = re.compile(r"^---\n(.*?)\n---\n", re.DOTALL)
PROP_RE = re.compile(r"^([\u4e00-\u9fa5A-Za-z0-9 /]+)\s*[:：]\s*(.+)$")


def read_text(path):
    with open(path, "r", encoding="utf-8") as f:
        return f.read()


def sha256(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(65536), b""):
            h.update(chunk)
    return h.hexdigest()


def classify_sensitivity(text, classification):
    """返回 (visibility, pii_flag, reasons)。

    可见性以文档自带的「文档密级」字段为权威（公开/内部/机密）；
    关键词扫描仅作为 PII/密钥泄漏的「复核标记」，不直接提升密级，
    避免把仅提及"密钥/密码"字样的普通文档误判为机密。
    """
    reasons = []
    visibility = CLASSIFICATION_LEVEL.get(classification, 1)
    vis = "public" if visibility == 0 else ("internal" if visibility == 1 else "confidential")
    if vis == "confidential":
        reasons.append(f"文档密级={classification}")
    low = text.lower()
    pii = False
    for kw in SENSITIVE_KEYWORDS:
        if kw.lower() in low:
            pii = True
            reasons.append(f"含敏感词:{kw}")
    return vis, pii, reasons


def parse_metadata(text):
    """优先解析 YAML frontmatter；否则解析"文档信息"表格里的 字段: 值。"""
    meta = {}
    fm = FRONT_MATTER_RE.match(text)
    if fm:
        for line in fm.group(1).splitlines():
            m = PROP_RE.match(line.strip())
            if m:
                meta[m.group(1).strip()] = m.group(2).strip()
    # 文档信息表格（兼容本库格式：| 字段 | 内容 |）
    in_doc_info = False
    for line in text.splitlines():
        if line.strip() == "## 文档信息" or line.strip() == "## 文档信息 ":
            in_doc_info = True
            continue
        if in_doc_info and line.strip().startswith("## "):
            in_doc_info = False
        if in_doc_info and line.strip().startswith("|"):
            cells = [c.strip() for c in line.strip().strip("|").split("|")]
            if len(cells) >= 2 and cells[0] and cells[0] not in ("项目", "------"):
                k = cells[0]
                v = cells[1]
                meta.setdefault(k, v)
    return meta


def norm_target(src_dir, link):
    """把相对链接解析为相对 wiki 根的路径（去锚点/查询，做存在性校验）。"""
    link = link.split("#")[0].split("?")[0].strip()
    if not link or link.startswith(("http://", "https://", "mailto:", "//")):
        return None
    if link.startswith("/"):
        rel = link.lstrip("/")
    else:
        rel = os.path.normpath(os.path.join(src_dir, link))
    return rel.replace(os.sep, "/")


def build_index():
    nodes = {}
    links = []          # 已有链接（已验证 / 断链）
    path_to_id = {}
    dir_index = set()
    for root, dirs, files in os.walk(SCAN_ROOT):
        if ".kg" in root.split(os.sep):
            continue
        rel_root = os.path.relpath(root, SCAN_ROOT).replace(os.sep, "/")
        if rel_root != ".":
            dir_index.add(rel_root.lower())
        for fn in files:
            if not fn.lower().endswith(".md"):
                continue
            full = os.path.join(root, fn)
            rel = os.path.relpath(full, SCAN_ROOT).replace(os.sep, "/")
            nodes[rel] = {
                "path": rel,
                "title": fn,
                "headings": [],
                "sha": sha256(full),
                "meta": {},
                "visibility": "internal",
                "pii_flag": False,
                "sensitive_reasons": [],
                "size": os.path.getsize(full),
            }
            path_to_id[rel.lower()] = rel

    # 版本后缀剥离索引：去掉 -Vx.x.x 后仍匹配的，判为"命名漂移"而非真断链
    base_index = {}          # stripped_basename -> rel（兜底，可能跨目录碰撞）
    path_base_index = {}     # (dir_lower, stripped_basename) -> rel（精确）
    for rel in nodes:
        d = os.path.dirname(rel).lower()
        base = rel.rsplit("/", 1)[-1]
        sv = strip_version(base)
        base_index.setdefault(sv, rel)
        path_base_index.setdefault((d, sv), rel)

    # 解析元数据、标题、链接
    for rel, node in nodes.items():
        text = read_text(os.path.join(SCAN_ROOT, rel))
        meta = parse_metadata(text)
        node["meta"] = meta
        classification = meta.get("文档密级", "内部")
        vis, pii, reasons = classify_sensitivity(text, classification)
        node["visibility"] = vis
        node["pii_flag"] = pii
        node["sensitive_reasons"] = reasons
        node["classification"] = classification
        node["tags"] = [t.strip() for t in meta.get("关联标签", "").replace("、", ",").split(",") if t.strip()]
        # 标题
        for line in text.splitlines():
            h = HEADING_RE.match(line)
            if h:
                node["headings"].append(h.group(1).strip())
            if node["title"] == fn and line.startswith("# "):
                node["title"] = line[2:].strip()
        # 链接
        src_dir = os.path.dirname(rel)
        seen = set()
        for m in LINK_RE.finditer(text):
            raw = m.group(1)
            tgt = norm_target(src_dir, raw)
            if tgt is None or tgt in seen:
                continue
            seen.add(tgt)
            links.append({"source": rel, "raw": raw, "norm": tgt, "kind": "markdown"})
        for m in WIKILINK_RE.finditer(text):
            tgt_name = m.group(1).split("|")[0].split("#")[0].strip()
            links.append({"source": rel, "raw": tgt_name, "norm": tgt_name,
                          "kind": "wikilink"})
    return nodes, links, path_to_id, dir_index, base_index, path_base_index


ASSET_EXT = {".html", ".htm", ".zip", ".png", ".jpg", ".jpeg", ".gif",
              ".svg", ".pdf", ".docx", ".xlsx", ".pptx"}


def compute_edges(nodes, links, path_to_id, dir_index, base_index, path_base_index):
    edges = []
    broken = []
    dir_links = []
    asset_links = []
    for l in links:
        if l["kind"] == "wikilink":
            # 本库实际未使用 wikilink；若未来启用，按名称模糊匹配为潜在建议
            continue
        tgt = l["norm"]
        tl = tgt.lower()
        if tl in path_to_id:
            if l["source"] == path_to_id[tl]:
                continue  # 跳过自环
            edges.append({"source": l["source"], "target": path_to_id[tl],
                          "type": "reference", "confidence": 1.0, "drift": False})
        elif tl in dir_index or tgt.endswith("/"):
            dir_links.append({"source": l["source"], "target": tgt})
        else:
            ext = os.path.splitext(tgt)[1].lower()
            if ext in ASSET_EXT:
                asset_links.append({"source": l["source"], "target": tgt})
                continue
            # 尝试去掉版本后缀匹配（命名漂移）
            base = tgt.rsplit("/", 1)[-1]
            stripped = strip_version(base)
            tdir = os.path.dirname(tgt).lower()
            resolved = path_base_index.get((tdir, stripped)) or base_index.get(stripped)
            if resolved:
                if l["source"] == resolved:
                    continue  # 自环
                edges.append({"source": l["source"], "target": resolved,
                              "type": "reference", "confidence": 1.0, "drift": True})
            else:
                broken.append({"source": l["source"], "target": tgt, "raw": l["raw"]})
    # 去重
    seen = set()
    dedup = []
    for e in edges:
        key = (e["source"], e["target"])
        if key not in seen:
            seen.add(key)
            dedup.append(e)
    return dedup, broken, dir_links, asset_links


def cross_link_suggestions(nodes, edge_set):
    """发现潜在交叉引用（cross-linker），给出置信度评分。"""
    # 标题/标签 token
    def tokens(node):
        toks = set()
        for h in node["headings"]:
            for t in TERM_DICT:
                if t in h.lower():
                    toks.add(t)
        for tag in node["tags"]:
            for t in TERM_DICT:
                if t in tag.lower():
                    toks.add(t)
        return toks

    tok_map = {rel: tokens(n) for rel, n in nodes.items()}
    suggestions = []
    rels = list(nodes.keys())
    for i in range(len(rels)):
        for j in range(i + 1, len(rels)):
            a, b = rels[i], rels[j]
            if (a, b) in edge_set or (b, a) in edge_set:
                continue  # 已有显式引用则不重复建议
            ta, tb = tok_map[a], tok_map[b]
            overlap = ta & tb
            if not overlap:
                continue
            # 置信度：共享术语数 / 较小方术语数，封顶
            denom = min(len(ta), len(tb)) or 1
            conf = round(min(1.0, 0.45 + 0.3 * len(overlap) / denom), 2)
            # 收紧：至少共享 2 个术语且置信度 >= 0.7，避免泛泛关联
            if conf >= 0.7 and len(overlap) >= 2:
                suggestions.append({"source": a, "target": b,
                                     "type": "potential", "confidence": conf,
                                     "shared_terms": sorted(overlap)})
    suggestions.sort(key=lambda s: -s["confidence"])
    return suggestions


def filter_safe(nodes, edges, suggestions):
    """安全隔离：剔除 confidential 节点，剔除跨密级边与含敏感词的潜在建议。"""
    safe_nodes = {r: n for r, n in nodes.items() if n["visibility"] != "confidential"}
    safe_ids = set(safe_nodes.keys())
    safe_edges = [e for e in edges if e["source"] in safe_ids and e["target"] in safe_ids]
    # 潜在链接：两端都在安全集，且不含敏感词
    safe_sug = [s for s in suggestions
                if s["source"] in safe_ids and s["target"] in safe_ids]
    return safe_nodes, safe_edges, safe_sug


def topo_insights(nodes, edges):
    indeg = defaultdict(int)
    outdeg = defaultdict(int)
    for e in edges:
        outdeg[e["source"]] += 1
        indeg[e["target"]] += 1
    orphans = [r for r in nodes if indeg[r] == 0 and outdeg[r] == 0]
    hubs = sorted(((r, indeg[r]) for r in nodes), key=lambda x: -x[1])[:10]
    return orphans, hubs, indeg, outdeg


def write_markdown(nodes, edges, suggestions, broken, orphans, hubs, dir_links, drift_count, asset_links):
    lines = []
    lines.append("# XYFamily Wiki 知识图谱报告\n")
    lines.append(f"> 生成时间: {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M UTC')}  ")
    lines.append(f"> 节点(文档): {len(nodes)}  边(已验证引用): {len(edges)}  ")
    lines.append(f"> 潜在交叉引用: {len(suggestions)}  真断链: {len(broken)}  ")
    lines.append(f"> 命名漂移边(可修复): {drift_count}  目录索引链接: {len(dir_links)}  "
                  f"资源链接: {len(asset_links)}  ")
    if asset_links:
        lines.append("")
        lines.append("**资源链接**（指向 .html/.zip 等非文档资源，不计为断链）：")
        for a in asset_links[:20]:
            lines.append(f"- `{a['source']}` → `{a['target']}`")
        if len(asset_links) > 20:
            lines.append(f"- …共 {len(asset_links)} 条")
    lines.append(f"> 孤立页: {len(orphans)}\n")
    lines.append("---\n")
    lines.append("## 一、目录拓扑（按被引用数 Top10 枢纽页）\n")
    lines.append("| 文档 | 被引用数 | 密级 | 标签 |")
    lines.append("|------|----------|------|------|")
    for r, d in hubs:
        n = nodes[r]
        lines.append(f"| `{r}` | {d} | {n.get('classification','-')} | {','.join(n['tags']) or '-'} |")
    lines.append("\n## 二、孤立页（无任何入链/出链，建议补充引用或合并）\n")
    for r in sorted(orphans):
        lines.append(f"- `{r}` — {nodes[r].get('classification','-')}")
    lines.append("\n## 三、断链（链接目标不存在，需修复）\n")
    if broken:
        for b in broken:
            lines.append(f"- `{b['source']}` → `{b['target']}`")
    else:
        lines.append("- 无")
    lines.append("\n## 四、潜在交叉引用建议（cross-linker，按置信度排序，Top 30）\n")
    lines.append("> 置信度基于标题/标签术语重叠。仅供人工审阅，未自动写入源文件（保持只读）。\n")
    lines.append("| 置信度 | 来源 | 目标 | 共享术语 |")
    lines.append("|--------|------|------|----------|")
    for s in suggestions[:30]:
        lines.append(f"| {s['confidence']} | `{s['source']}` | `{s['target']}` | {','.join(s['shared_terms'])} |")
    lines.append("\n## 五、安全隔离说明\n")
    conf = [r for r, n in nodes.items() if n["visibility"] == "confidential"]
    pii = [r for r, n in nodes.items() if n["pii_flag"] and n["visibility"] != "confidential"]
    lines.append(f"- 机密文档（按「文档密级=机密」，**已排除出安全视图** `safe_graph.json`）：{len(conf)} 篇")
    for r in conf:
        lines.append(f"  - `{r}` — {nodes[r]['sensitive_reasons']}")
    lines.append(f"- 疑似涉密未标注（提及密钥/密码/PII 但密级非机密，**保留在安全视图、建议复核密级**）：{len(pii)} 篇")
    for r in pii[:25]:
        lines.append(f"  - `{r}` — {nodes[r]['sensitive_reasons']}")
    if len(pii) > 25:
        lines.append(f"  - …共 {len(pii)} 篇")
    lines.append("\n---\n*由 wiki_kg.py 生成，零依赖、只读、可增量重建。*")
    return "\n".join(lines)


def write_html(safe_nodes, safe_edges, safe_sug):
    nodes_json = json.dumps(
        [{"id": r, "vis": n["visibility"], "title": n["title"]} for r, n in safe_nodes.items()],
        ensure_ascii=False)
    edges_json = json.dumps(
        [{"s": e["source"], "t": e["target"], "c": e["confidence"]} for e in safe_edges]
        + [{"s": s["source"], "t": s["target"], "c": s["confidence"], "potential": True}
           for s in safe_sug],
        ensure_ascii=False)
    html = """<!doctype html><html lang="zh"><head><meta charset="utf-8">
<title>XYFamily Wiki 知识图谱</title>
<script src="https://cdn.jsdelivr.net/npm/vis-network@9/standalone/umd/vis-network.min.js"></script>
<style>body{font-family:sans-serif;margin:0}#net{height:92vh;width:100%}</style></head>
<body><div id="net"></div><script>
const nodes=__NODES__;
const edges=__EDGES__;
const nds=nodes.map(n=>({id:n.id,label:n.id.split('/').pop(),group:n.vis,title:n.title}));
const eds=edges.map(e=>({from:e.s,to:e.t,
  dashes:e.potential||false,value:(e.c||1)*2,
  color:e.potential?'#e0a000':'#4a90d9',
  title:e.potential?('潜在链接 置信度 '+(e.c||'')):'已验证引用'}));
new vis.Network(document.getElementById('net'),{nodes:nds,edges:eds},{
  physics:{barnesHut:{springLength:120}},groups:{
  public:{color:{background:'#a5d6a7'}},internal:{color:{background:'#90caf9'}},
  confidential:{color:{background:'#ef9a9a'}}}});
</script></body></html>"""
    html = html.replace("__NODES__", nodes_json).replace("__EDGES__", edges_json)
    return html


def main():
    os.makedirs(KG_DIR, exist_ok=True)
    prev_state = {}
    if os.path.exists(STATE_FILE):
        try:
            prev_state = json.load(open(STATE_FILE, encoding="utf-8")).get("hashes", {})
        except Exception:
            prev_state = {}

    nodes, links, path_to_id, dir_index, base_index, path_base_index = build_index()
    edges, broken, dir_links, asset_links = compute_edges(nodes, links, path_to_id, dir_index, base_index, path_base_index)
    edge_set = set((e["source"], e["target"]) for e in edges)
    drift_count = sum(1 for e in edges if e.get("drift"))
    suggestions = cross_link_suggestions(nodes, edge_set)
    safe_nodes, safe_edges, safe_sug = filter_safe(nodes, edges, suggestions)
    orphans, hubs, indeg, outdeg = topo_insights(nodes, edges)

    # Delta 追踪统计
    changed, added, unchanged = [], [], []
    for rel, n in nodes.items():
        old = prev_state.get(rel)
        if old is None:
            added.append(rel)
        elif old != n["sha"]:
            changed.append(rel)
        else:
            unchanged.append(rel)

    full = {
        "meta": {"generated": datetime.now(timezone.utc).isoformat(),
                 "node_count": len(nodes), "edge_count": len(edges),
                 "drift_edge_count": drift_count, "dir_link_count": len(dir_links),
                 "asset_link_count": len(asset_links),
                 "potential_count": len(suggestions), "broken_count": len(broken),
                 "orphan_count": len(orphans),
                 "delta": {"added": len(added), "changed": len(changed),
                           "unchanged": len(unchanged)}},
        "nodes": nodes, "edges": edges,
        "potential_links": suggestions, "broken_links": broken,
        "dir_links": dir_links, "asset_links": asset_links,
    }
    safe = {"nodes": {r: {"path": n["path"], "title": n["title"],
                          "classification": n.get("classification"),
                          "tags": n["tags"], "visibility": n["visibility"]}
                      for r, n in safe_nodes.items()},
            "edges": safe_edges, "potential_links": safe_sug}

    with open(os.path.join(KG_DIR, "knowledge_graph.json"), "w", encoding="utf-8") as f:
        json.dump(full, f, ensure_ascii=False, indent=2)
    with open(os.path.join(KG_DIR, "safe_graph.json"), "w", encoding="utf-8") as f:
        json.dump(safe, f, ensure_ascii=False, indent=2)

    md = write_markdown(nodes, edges, suggestions, broken, orphans, hubs, dir_links, drift_count, asset_links)
    with open(os.path.join(KG_DIR, "GRAPH.md"), "w", encoding="utf-8") as f:
        f.write(md)
    html = write_html(safe_nodes, safe_edges, safe_sug)
    with open(os.path.join(KG_DIR, "graph.html"), "w", encoding="utf-8") as f:
        f.write(html)

    state = {"updated": datetime.now(timezone.utc).isoformat(),
             "hashes": {r: n["sha"] for r, n in nodes.items()}}
    with open(STATE_FILE, "w", encoding="utf-8") as f:
        json.dump(state, f, ensure_ascii=False, indent=2)

    print(f"[OK] 节点={len(nodes)} 边={len(edges)} (命名漂移边={drift_count}) "
          f"目录链接={len(dir_links)} 资源链接={len(asset_links)} 潜在={len(suggestions)} "
          f"断链={len(broken)} 孤立={len(orphans)}")
    print(f"[Delta] 新增={len(added)} 变更={len(changed)} 未变={len(unchanged)}")
    print(f"[安全] 机密(按密级)={sum(1 for n in nodes.values() if n['visibility']=='confidential')} "
          f"疑似涉密未标注(pii_flag)={sum(1 for n in nodes.values() if n['pii_flag'] and n['visibility']!='confidential')}")
    print(f"[产物] {KG_DIR}")


if __name__ == "__main__":
    main()
