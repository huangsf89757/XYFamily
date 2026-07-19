# 原型导出说明 - 网页端

> 交付达（Jiao）· 导出交付专家产出
> 导出时间：2025-07-15 · 源文件已通过 5 维质量评审（23/25 PASS）并完成定向精修

---

## 一、交付物清单

| 交付物 | 路径 | 说明 |
|--------|------|------|
| 主交付物（HTML） | `01-网页端/index.html` | 自包含单文件 SPA，双击即可离线运行 |
| 分发打包（ZIP） | `01-网页端/XYFamily-网页端原型.zip` | 内含唯一文件 `index.html`，一键下载/分发 |

> 交付说明本文件：`01-网页端/EXPORT-README.md`

---

## 二、文件信息

- **文件名**：`index.html`
- **行数**：2,669 行
- **体积**：189,179 字节（≈ 185 KB；磁盘占用 188 KB）
- **压缩后（ZIP 内）**：约 39 KB（deflate 79%）
- **技术形态**：HTML + 内联 CSS + 内联 JS，零前端框架，纯原生 SPA（基于 `location.hash` 路由）

---

## 三、打开方式

1. **直接打开（推荐）**：双击 `index.html`，在任意现代浏览器（Chrome / Edge / Safari / Firefox）中即可运行。
2. **本地静态服务**（如需避免个别浏览器对 `file://` 的限制）：
   ```bash
   # 在 01-网页端/ 目录下任选其一
   python3 -m http.server 8080      # 然后访问 http://localhost:8080
   npx serve .
   ```
3. **分发方式**：将 `XYFamily-网页端原型.zip` 解压后，同样双击其中的 `index.html` 即可。

---

## 四、资源内联校验（交付前）

| 校验项 | 结果 |
|--------|------|
| 外部 CSS 文件 | ✅ 无（全部内联 `<style>`） |
| 外部 JS 文件 | ✅ 无（全部内联 `<script>`） |
| 图片 / 图标 | ✅ 全部内联（SVG 直写 + `data:image/svg+xml` mask），无 `.png/.jpg/.gif` 外链 |
| 字体 | ⚠️ 仅 Google Fonts CDN（DM Sans / JetBrains Mono），**已含系统字体栈回退**；离线时自动降级为系统字体，不报错、不影响布局 |
| 其他外链 | ✅ 无（`http(s)://` 仅出现在字体 CDN 与 SVG 命名空间声明） |

**结论**：除字体 CDN 外零外部依赖，离线双击可完整运行。

---

## 五、覆盖范围（路由清单）

**核心业务模块（9 大模块）**

1. **控制台 Dashboard** — `#/`
2. **组织管理 Org** — `#/org`（含 `#/org/members` 成员管理、`#/org/permissions` 角色权限）
3. **团队管理 Team** — `#/team/:id`（含 `:id/members` 团队成员、`:id/groups` 团队小组）
4. **小组管理 Group** — `#/group/:id`（含 `:id/members` 小组成员）
5. **角色权限 / 权限矩阵** — `#/org/permissions`、`#/permissions/matrix`、`#/permissions/overview`
6. **审计日志 Audit** — `#/audit`
7. **个人中心** — `#/profile`（个人信息）、`#/security`（密码与安全）、`#/account`（账号设置）
8. **超级管理员 SuperAdmin** — `#/superadmin`、`#/init`（系统初始化）
9. **成员管理入口** — 收敛于组织管理下的成员视图（邀请 / 移除 / 角色分配）

**全局态（Global States）**

- **认证态**：`#/login`、`#/register`、`#/reset-password`
- **错误态**：`#/403`（无权限）、`#/404`（页面不存在）
- **主题态**：明暗双主题切换（`toggle-theme`）
- **角色态**：角色切换器（`set-role`），覆盖 SuperAdmin / OrgAdmin / SecurityAdmin / AuditAdmin / TeamAdmin 等 9 种角色 Badge

---

## 六、设计对齐

- **设计语言**：Pinguo 设计语言（清爽、克制、信息密度适中）
- **单一主色**：`--brand-500: #007AFF`（全局唯一主色，集中于 `--primary` / `--sidebar-primary`）
- **明暗双主题**：通过 `[data-theme]` + CSS 变量重映射实现（浅色 `--primary:#007aff`，深色 `--primary:#0a84ff/brand-400`）
- **响应式断点**：适配 1920 / 1440 / 1280 三档宽度，侧边栏可折叠、移动端抽屉态
- **角色 Badge 配色一致**：各角色 Badge 采用统一语义色板（如 `security` 紫 `#AF52DE`/`#BF5AF2`），已修复撞色问题

---

## 七、质量结论

- **5 维质量评审**：**23 / 25 PASS**（Anti-Slop 达标，品牌级水准）
- **定向精修已落地**：
  - **P1（已修）**：角色 Badge 撞色 —— 重新分配语义色板，确保 9 种角色 Badge 互不混淆
  - **P2（顺手修）**：圆角令牌化（统一 `--brand-radius-*`）、下拉框去硬编码样式、表格横向滚动兜底、清理死代码

---

## 八、已知简化（演示边界）

- 部分二级列表为带「查看全部」入口的轻量占位，非全量数据。
- 登录错误态通过**演示切换器**展示（如「登录过期」「账号锁定」），非真实校验后端。
- 全流程**无真实后端**：数据均为前端内置 mock，交互（增删改、角色切换、主题切换）可演示但无持久化。

---

## 九、使用提示

- 首次打开建议确认网络可访问 Google Fonts；若离线，字体自动回退为系统字体，布局不受影响。
- 路由基于 URL hash，刷新 / 分享链接可精确定位到具体页面（如 `#/org/members`）。
- 切换角色后部分菜单按权限动态显隐，属预期行为。
