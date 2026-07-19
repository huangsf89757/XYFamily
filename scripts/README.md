# 临时脚本 / Scripts

本目录集中存放项目维护用的临时/工具脚本，与 `wiki/` 内容协同工作。
应用代码请见仓库根目录的 `code/`（Go 后端等），不属于此处范畴。

> ⚠️ 注意：所有脚本默认在**项目根目录（XYFamily/）**下运行，路径均基于脚本位置向上一级定位到项目根。

---

## 1. `fix_links.py` — Wiki 链接修复与交叉引用回写

扫描 `wiki/` 下的 Markdown 文件，校验内部链接是否指向真实存在的文件，
并解析正文中形如 `[[文件夹/页面]]` 的 wikilink，将生成的交叉引用卡片回写到各页面。

**依赖**：Python 3（仅标准库，无需安装第三方包）。

**用法**：
```bash
# 仅预览（dry-run），不修改任何文件
python3 scripts/fix_links.py

# 实际修改文件（回写交叉引用、修正链接）
python3 scripts/fix_links.py --apply

# 跳过交叉引用回写，只做链接校验
python3 scripts/fix_links.py --no-backlinks
```

**输出**：
- 控制台报告：✅ 有效链接 / ⚠️ 悬空链接（列出了目标缺失的明细）。
- `--apply` 时生成 `wiki/manual_review_needed.md`，汇总需要人工确认的悬空链接。

---

## 2. `wiki_kg.py` — Wiki 知识图谱构建

扫描 `wiki/` 目录，解析知识图谱特征（节点规模、文件夹层级、悬空链接、
孤立节点、链接数 Top、交叉引用数量等），输出 JSON 到 `graphify-out/wiki_kg.json`
（供 graphify 知识图谱工具消费）。

**依赖**：Python 3（仅标准库）。

**用法**：
```bash
python3 scripts/wiki_kg.py
```

**输出**：`graphify-out/wiki_kg.json`

---

## 3. `server.py` — Wiki 浏览页本地开发服务器

基于 Python 标准库 `http.server` 的静态文件服务器，用于本地预览 `wiki/浏览页/`
下生成的 HTML 浏览页（支持 SPA 回退、`Brotli`/`gzip` 压缩）。

**依赖**：Python 3（仅标准库）。

**用法**：
```bash
python3 scripts/server.py            # 默认 http://localhost:8000
python3 scripts/server.py 9000      # 指定端口
```

然后用浏览器打开 `http://localhost:8000` 即可浏览。

---

## 4. `fix-head.mjs` — 浏览页 HTML 头部修复（一次性脚本）

为 `xyfamily-admin-design/pages/` 下导出的 HTML 浏览页注入统一的
`<head>`（meta、字体、Apple 设计库 CSS 引用、标题映射等）。

**依赖**：Node.js（仅标准库 `fs`/`path`）。

> ⚠️ 该脚本为**一次性/项目特定**脚本，内部路径为硬编码绝对路径（见文件顶部
> `CSS_PATH` 与 `PAGES_DIR`），指向本地设计导出目录。若要在其他环境复用，
> 请先修改这两个常量，或改为读取环境变量/命令行参数。

**用法**：
```bash
node scripts/fix-head.mjs
```

# 说明
- 本目录内的脚本均为**维护/构建辅助工具**，不随应用一起部署。
- 改动脚本后请重新运行语法检查：`python3 -m py_compile scripts/*.py` / `node --check scripts/fix-head.mjs`。
