# XYFamily 网页端 · 设计令牌与组件规范（DESIGN-TOKENS）

> **用途**：本文档是原型构建师的**唯一令牌与组件依据**。拿到后无需再读《全局设计规范》《网页端设计规范》，即可写出完全对齐的 CSS 变量与组件。
> **品牌锁定**：Pinguo 设计语言（克制、优雅、苹果风）。单一品牌主色 `#007AFF` + 中性灰阶 + 语义状态色。**组件只消费语义令牌，禁止硬编码色值。**
> **版本**：V1.0.0 · 提炼自 `wiki/.../02-设计规范` 与 `xyfamily-admin-design` 现存原型（login.html / account.html / admin-shell.html）。

---

## 0. 使用铁律（先读）

1. **两层架构**：底层是原始色板（brand / background / text / icon / state-*），上层是语义令牌（`--primary` / `--card` / `--border` …）。**组件与页面只能引用语义令牌**，原始色板仅在 `:root` 映射中使用。
2. **禁止硬编码**：页面里出现 `#007AFF`、`#1D1D1F`、`rgb(...)` 等字面色值均视为缺陷（现存原型已有此类问题，见第 9 节）。
3. **主题切换**：`<head>` 内联脚本在页面解析前设 `class`（`.light` / `.dark`），从 localStorage 读取，防闪烁。所有颜色令牌在 `:root`（light）与 `.dark` 下各定义一份。
4. **新页面一律用 `--brand-radius-sm/md/lg`**；保留 `--radius` 仅为了**不破坏旧页面**，新代码不要再直接用 `--radius`。

---

## 1. 色彩令牌

### 1.1 原始色板（Primitive Palette，仅用于 `:root` 映射）

**品牌色（Brand）** — 锚定 Apple System Blue

| 级别 | 色值 | 用途 |
|------|------|------|
| brand-50 | `#E8F2FF` | 最浅背景着色（浅底蓝） |
| brand-100 | `#CFE5FF` | 浅色背景、hover 态 |
| brand-200 | `#9FCBFF` | 辅助元素 |
| brand-300 | `#66ABFF` | 禁用态主色 |
| **brand-400** | **`#2E8DFF`** | **深色主题主色（dark `--primary`）** |
| **brand-500** | **`#007AFF`** | **主色（light `--primary`）** |
| brand-600 | `#0064D6` | 按压态 |
| brand-700 | `#004FAD` | 深色强调 |
| brand-800 | `#003B82` | 深色背景文字 |
| brand-900 | `#00275A` | 最深色 |

**中性灰阶（Background / Text / Icon）** — 10 级

| 级 | background | text | icon | 典型用途 |
|----|-----------|------|------|---------|
| 50 | `#FFFFFF` | `#F5F5F7` | `#F5F5F7` | 页面/卡片背景（light）；深主题主文字/图标 |
| 100 | `#F7F7FA` | `#E3E3E8` | `#E5E5EA` | 次级背景、强调背景、表头浅灰 |
| 200 | `#F2F2F7` | `#C7C7CC` | `#D1D1D6` | 分组/侧边栏背景（light）；disabled 文字 |
| 300 | `#E5E5EA` | `#AEAEB2` | `#C7C7CC` | 边框、分割线；占位文字 |
| 400 | `#D1D1D6` | `#8E8E93` | `#AEAEB2` | 输入框边框；次级图标、次要说明 |
| 500 | `#AEAEB2` | `#6E6E73` | `#8E8E93` | 次级图标；辅助信息文字 |
| 600 | `#8E8E93` | `#48484A` | `#6E6E73` | 深色文字；图标 |
| 700 | `#3A3A3C` | `#3C3C43` | `#48484A` | 深色卡片；标题文字 |
| 800 | `#1C1C1E` | `#1D1D1F` | `#2C2C2E` | 深色背景；主文字/主图标 |
| 900 | `#000000` | `#000000` | `#1D1D1F` | 纯黑；弹出层文字 |

**语义状态色（State Colors）** — 锚定 Apple 系统色，**仅用于真实语义状态**

| 语义 | Light | Dark | surface（浅底） | foreground |
|------|-------|------|----------------|-----------|
| Success | `#34C759` | `#30D158` | `#E9F9EE` | `#FFFFFF` |
| Error | `#FF3B30` | `#FF453A` | `#FFECEA` | `#FFFFFF` |
| Warning | `#FF9500` | `#FF9F0A` | `#FFF4E6` | `#FFFFFF` |
| Info | `#5AC8FA` | `#64D2FF` | `#E6F7FE` | `#FFFFFF` |

> ⚠ 现存原型**只定义了 state-success / state-error**，缺失 state-warning / state-info。本文档已补全（见第 8 节完整 CSS）。组件里出现的 `--state-warning`、`--state-info` 现已可用。

### 1.2 语义令牌（Semantic Tokens，组件消费层）

Light（`:root`）/ Dark（`.dark`）映射关系（完整 CSS 见第 8 节）：

| 令牌 | Light 映射 | Dark 映射 | 用途 |
|------|-----------|-----------|------|
| `--background` | background-50 | background-900 | 页面背景 |
| `--foreground` | text-800 | text-50 | 主文字 |
| `--card` | background-50 | background-800 | 卡片背景 |
| `--card-foreground` | text-800 | text-50 | 卡片文字 |
| `--popover` | background-50 | background-700 | 弹出层背景 |
| `--popover-foreground` | text-900 | text-50 | 弹出层文字 |
| `--primary` | **brand-500** | **brand-400** | 主色 |
| `--primary-foreground` | background-50 | background-900 | 主色上文字 |
| `--secondary` | background-200 | background-800 | 次级背景 |
| `--secondary-foreground` | text-800 | text-50 | 次级文字 |
| `--muted` | background-200 | background-800 | 静默背景 |
| `--muted-foreground` | text-400 | text-400 | 静默文字 |
| `--accent` | background-100 | background-700 | 强调背景 |
| `--accent-foreground` | text-800 | text-50 | 强调文字 |
| `--destructive` | state-error | state-error-dark | 危险操作色 |
| `--destructive-foreground` | state-error-foreground | state-error-foreground | 危险色上文字 |
| `--success` | state-success | state-success-dark | 成功色 |
| `--success-foreground` | state-success-foreground | state-success-foreground | 成功色上文字 |
| `--warning` | state-warning | state-warning-dark | 警告色 |
| `--warning-foreground` | state-warning-foreground | state-warning-foreground | 警告色上文字 |
| `--info` | state-info | state-info-dark | 信息色 |
| `--info-foreground` | state-info-foreground | state-info-foreground | 信息色上文字 |
| `--border` | background-300 | background-700 | 边框 |
| `--input` | background-400 | background-700 | 输入框边框 |
| `--ring` | **brand-500** | **brand-400** | 焦点环 |
| `--icon` | icon-900 | icon-50 | 主图标 |
| `--icon-muted` | icon-500 | icon-500 | 次级图标 |

**侧边栏专用语义令牌**

| 令牌 | Light | Dark | 用途 |
|------|-------|------|------|
| `--sidebar` | background-200 | background-800 | 侧边栏背景 |
| `--sidebar-foreground` | text-800 | text-50 | 侧边栏文字 |
| `--sidebar-primary` | brand-500 | brand-400 | 侧边栏主色（选中/Logo） |
| `--sidebar-primary-foreground` | background-50 | background-900 | 主色上文字 |
| `--sidebar-accent` | background-300 | background-700 | 侧边栏强调 |
| `--sidebar-accent-foreground` | text-800 | text-50 | 强调上文字 |
| `--sidebar-border` | background-300 | background-700 | 侧边栏边框 |
| `--sidebar-ring` | brand-500 | brand-400 | 侧边栏焦点环 |

**图表色（Chart）**：chart-1 success · chart-2 brand · chart-3 warning · chart-4 `#5856D6` · chart-5 `#AF52DE`（dark 下 chart-3 `#FF9F0A`、chart-4 `#5E5CE6`、chart-5 `#BF5AF2`）。

### 1.3 角色 Badge 配色（Role Badge）— 浅底深字

| 角色 | 背景 | 文字 | 说明 |
|------|------|------|------|
| SuperAdmin | `#FFECEA`（state-error-surface） | `#FF3B30`（state-error） | 超级管理员 |
| OrgAdmin | `#E8F2FF`（brand-50） | `#007AFF`（brand-500） | 组织核心管理员 |
| TeamAdmin | `#DBEAFE` | `#2563EB` | 团队管理员（固定蓝，非 brand 阶） |
| GroupAdmin | `#E9F9EE`（state-success-surface） | `#34C759`（state-success） | 小组管理员 |
| Member | `#F2F2F7`（background-200） | `#6E6E73`（text-500） | 普通成员 |
| Public | `#F2F2F7`（background-200） | `#8E8E93`（text-400） | 公开访客 |

令牌化：`--badge-superadmin-bg/fg`、`--badge-orgadmin-bg/fg`、`--badge-teamadmin-bg/fg`、`--badge-groupadmin-bg/fg`、`--badge-member-bg/fg`、`--badge-public-bg/fg`（见第 8 节）。

> TeamAdmin 的 `#2563EB` / `#DBEAFE` 是规范指定的固定值（Tailwind blue-600 / blue-50），**不属于 brand 10 级色阶**，请直接写死这两个 hex（允许的唯一硬编码例外，因它是角色语义而非装饰色）。

---

## 2. 字体令牌

| 令牌 | 值 | 用途 |
|------|----|------|
| `--font-sans` | `"DM Sans", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Noto Sans CJK SC", sans-serif` | 正文、UI |
| `--font-serif` | 同 `--font-sans`（保持克制，标题与正文同族） | 标题 |
| `--font-mono` | `"JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace` | 代码、ID、配置值 |

**字号阶梯（px）**

| 令牌 | px | 用途 |
|------|----|------|
| `--text-xs` | 12 | 辅助说明、标签、表头、Badge |
| `--text-sm` | 14 | 次要正文、表单、表格、导航项 |
| `--text-base` | 16 | 正文默认 |
| `--text-lg` | 18 | 卡片标题、列表项标题 |
| `--text-xl` | 20 | 页面区块标题 |
| `--text-2xl` | 24 | 页面标题 |
| `--text-3xl` | 30 | 大标题 |
| `--text-4xl` | 36 | Hero 标题 |

**字重**：`--font-weight-normal: 400`（正文）/ `medium: 500`（按钮、表单标签、列表项）/ `semibold: 600`（卡片标题、导航选中、页面标题）/ `bold: 700`（页面标题、Hero）。

**行高**：`--leading-tight: 1.2`（标题）/ `--leading-normal: 1.5`（正文）/ `--leading-relaxed: 1.7`（长文本）。

> 字体经 CDN 加载（DM Sans / JetBrains Mono）。中文回退 PingFang SC / Noto Sans CJK SC。

---

## 3. 间距令牌

**基准**：`--spacing: 0.24rem`（≈3.84px，Tailwind v4 乘数）。精确 px 比例尺（组件直接消费，避免 Tailwind 0.24 基数带来的小数偏差）：

| 令牌 | px | 用途 |
|------|----|------|
| `--space-1` | 4 | 图标与文字间距 |
| `--space-2` | 8 | 紧凑内间距 |
| `--space-3` | 12 | 元素间默认间距 |
| `--space-4` | 16 | 卡片内间距、表单元素 |
| `--space-5` | 20 | 区块间距 |
| `--space-6` | 24 | 卡片间距、表单分组 |
| `--space-8` | 32 | 区块大间距、页面边距 |
| `--space-10` | 40 | 页面边距（大） |
| `--space-12` | 48 | 大区块间距 |
| `--space-16` | 64 | 页面垂直间距 |

**常用组合**：卡片内边距 16–24（`--space-4`~`--space-6`）；页面左右边距 32（`--space-8`）；区块间距 32–48（`--space-8`~`--space-12`）；表单元素内 8–12；图标+文字 4。

---

## 4. 圆角令牌（⚠ 现存缺三档，本文档已补全）

| 令牌 | 值 | 用途 |
|------|----|------|
| `--brand-radius-sm` | **8px** | 按钮、输入框、标签、小卡片、分页按钮、骨架条 |
| `--brand-radius-md` | **12px** | 卡片、下拉菜单、Popover、Alert、Toast、Table 容器 |
| `--brand-radius-lg` | **19.2px** | 模态/对话框、认证页大卡片、大型容器 |
| `--brand-radius-full` | **9999px** | 头像、Pill、Toggle、圆形图标按钮、Badge |
| `--radius` | 1.2rem（19.2px） | **保留兼容旧页面**，新代码勿用 |

> 现存原型 `:root` 仅定义了 `--radius`，但 `login.html`/`account.html` 的 Tailwind `@theme inline` 已引用 `--brand-radius-sm/md/lg`（**未定义 → 半径失效**）。第 8 节 CSS 已补全这四个变量。

---

## 5. 阴影令牌

| 令牌 | Light 值 | 用途 |
|------|---------|------|
| `--shadow-2xs` | `0 1px 2px -1px rgba(0,0,0,0.04)` | 细线 |
| `--shadow-xs` | `0 1px 2px 0 rgba(0,0,0,0.04)` | 微妙 |
| `--shadow-sm` | `0 1px 2px 0 rgba(0,0,0,0.05), 0 1px 3px -1px rgba(0,0,0,0.05)` | **静态卡片** |
| `--shadow` | `0 2px 4px -1px rgba(0,0,0,0.06), 0 1px 2px -1px rgba(0,0,0,0.05)` | 悬浮 |
| `--shadow-md` | `0 4px 8px -2px rgba(0,0,0,0.06), 0 2px 4px -2px rgba(0,0,0,0.05)` | **卡片 hover** |
| `--shadow-lg` | `0 8px 24px -8px rgba(0,0,0,0.08), 0 4px 8px -4px rgba(0,0,0,0.05)` | **下拉·Popover·侧边栏·Toast** |
| `--shadow-xl` | `0 16px 40px -10px rgba(0,0,0,0.10), 0 8px 16px -8px rgba(0,0,0,0.06)` | **认证卡片（Overlay）** |
| `--shadow-2xl` | `0 24px 64px -12px rgba(0,0,0,0.12)` | **模态框** |

Dark 模式阴影透明度整体抬高（深度感）：2xs/xs 0.30 · sm/`--shadow` 0.36~0.40 · md/lg 0.44~0.50 · xl/2xl 0.55~0.60（完整值见第 8 节 `.dark`）。

---

## 6. 布局令牌

| 令牌 | 值 | 用途 |
|------|----|------|
| `--sidebar-w` | **256px** | 侧边栏展开宽度（现存 240px ❌→256） |
| `--sidebar-w-collapsed` | **64px** | 侧边栏收起（图标）宽度 |
| `--topbar-h` | **56px** | 顶部栏高度（h-14） |
| `--content-max-w` | **1408px** | 内容区最大宽度（现存 max-w-1200 ❌→1408） |
| `--content-padding` | **32px** | 内容区左右内边距（px-8） |
| `--brand-logo-size` | **32px** | 侧边栏/顶部栏 Logo 尺寸 |

布局骨架（Desktop First，最小 1024px）：

```
侧边栏 fixed 256/64px  │  内容区 margin-left 随侧边栏  │  max-width 1408 居中  │  padding 32
                      │  顶部栏 sticky top-0 h-56 z-40  │  背景 --background-200
```

侧边栏过渡：`width` + 内容区 `margin-left` 均 `400ms var(--ease-in-out)`。

---

## 7. 组件令牌与规范速查

> 所有组件消费语义令牌。下表给「变体 / 尺寸 / 状态 / 关键样式」速查。Lucide 图标默认 20px，导航 16px，空态 64px。

### 7.1 按钮 Button — 6 变体

| 变体 | 背景 | 文字 | 边框 |
|------|------|------|------|
| Primary | `--primary` | `--primary-foreground` | 无 |
| Secondary | `--secondary` | `--secondary-foreground` | 无 |
| Outline | 透明 | `--foreground` | 1px `--border` |
| Ghost | 透明 | `--foreground` | 无 |
| Destructive | `--destructive` | `--destructive-foreground` | 无 |
| Link | 透明 | `--primary` | 无 |

| 尺寸 | 高 | 水平 padding | 字号 | 字重 | 圆角 | 图标 |
|------|----|-------------|------|------|------|------|
| sm | 32px | 12px | 13px | 500 | `--brand-radius-sm` | 14px |
| md | 40px | 20px | 14px | 500 | `--brand-radius-sm` | 16px |
| lg | 48px | 24px | 14px | 600 | `--brand-radius-sm` | 16px |

状态：Hover → Primary `--brand-600`、Secondary `--background-300`、Ghost `rgba(0,0,0,.04)`；Active → Primary `--brand-700`；Focus → `box-shadow: var(--focus-ring)`；Disabled → `opacity:.5; cursor:not-allowed`；**Loading** → 文字换 Spinner + “处理中…”，按钮禁用。带图标：左/右图标间距 8px；仅图标按钮需 `aria-label`。

### 7.2 输入框 Input — 4 态

| 尺寸 | 高 | 水平 padding | 字号 | 圆角 |
|------|----|-------------|------|------|
| sm | 36px | 12px | 13px | `--brand-radius-sm` |
| md | 40px | 16px | 14px | `--brand-radius-sm` |
| lg | 48px | 16px | 14px | `--brand-radius-sm` |

> ⚠ 输入框圆角统一用 `--brand-radius-sm`（8px）。（网页端规范 5.2.2 曾写 md/lg=12px，本文档按 lead 指令统一为 8px，见第 9 节偏差 #8。）

四态：Default（`--background-50` 底 / `--background-300` 框）→ Hover（框变 `--background-400`）→ **Focus**（框 `--primary` + `box-shadow: var(--focus-ring)` / `0 0 0 2px color-mix(in srgb,var(--ring) 18%,transparent)`）→ Error（框 `--destructive` + 下方错误文字 13px `--destructive`）→ Disabled（`--background-100` 底，`opacity:.5`）。前后缀图标 16px，前缀 `left-3`、后缀 `right-3`，有前缀 `pl-10`、有后缀 `pr-10`。Select 面板：`--popover` 底 + `--shadow-lg` + `--brand-radius-md`，选项高 36px，hover `--accent`，选中 `--brand-50`/`--brand-500`。Textarea：`min-h-80px`、`p-3`、`--brand-radius-sm`、`resize:vertical`、14px。

### 7.3 表格 Table

| 区域 | 规格 |
|------|------|
| 容器 | `--card` 底、`--brand-radius-md`、1px `--border`、`overflow:hidden`、宽 100% |
| 表头 | 高 40px、`--background-100`/近似浅灰底、12px semibold uppercase、`--text-500`、`px-4`、`white-space:nowrap`、操作列右对齐 |
| 行 | 高 **48px**、下边框 1px `--background-100`、hover 背景 `rgba(0,122,255,.03)`、`px-4`、14px（`--text-800` 主 / `--text-500` 次）、过渡 150ms |
| 操作列 | 右对齐，Ghost 文字按钮，间距 12px，主操作 `--primary`、危险 `--destructive`、13px、hover 下划线 |
| 分页 | 下方两端对齐；总数 14px `--text-500`；页码按钮 36×36、`--brand-radius-sm`、间距 8px；默认白底 1px `--border` 文字 `--text-600`，选中 `--primary` 白字，hover `--background-100`，禁用 `--text-400` |
| 空状态 | 容器内居中，上下 48px，图标 64px `--text-300`，标题 16px semibold `--text-600`，描述 14px `--text-400`，可选 Primary 按钮 |
| 加载 | 骨架屏 5~8 行，行高 48px，条 `--background-200` 底 `--brand-radius-sm`，`animate-pulse` |

### 7.4 对话框 Dialog / Modal

| 尺寸 | 宽 | 用途 |
|------|----|------|
| sm | 400px | 删除确认 |
| md | 520px | 简单表单（编辑名/邀请） |
| lg | 680px | 复杂表单（建组织/角色配置） |
| xl | 880px | 权限矩阵/详情 |

结构：遮罩 `rgba(0,0,0,.4)` + `backdrop-blur-sm` + z-50；框 `--card` 底 + `--brand-radius-lg` + `--shadow-2xl` + `fixed` 居中 + `max-h-85vh` 可滚；头部 `p-6` 标题 18px bold + 右上关闭 X；内容 `p-6`；底部 `p-6` 右对齐（Primary+Secondary）。
关闭：点遮罩（确认类除外）/ 点 X / **ESC**；**焦点陷阱**（Tab 循环）+ 关闭后焦点归还触发元素。
动画：进入 400ms `ease-out` 缩放 0.95→1 + 淡入；退出 250ms `ease-in` 缩放 1→0.95 + 淡出；遮罩 250ms `ease-default` 淡。

### 7.5 Badge / 标签

通用：padding `2.5px 10px`、full 圆角、12px medium、`inline-flex`。角色配色见 §1.3。状态标签：正常/启用=`--state-success` 浅底深字；已禁用=`--state-error`；待审核=`--state-warning`；进行中=`--brand-50`/`--brand-500`。

### 7.6 头像 Avatar

| 尺寸 | 直径 | 字号 | 用途 |
|------|------|------|------|
| xs | 24px | 10px | 表格行内 |
| sm | 32px | 12px | 表格成员列 |
| md | 36px | 14px | 侧边栏用户区 |
| lg | 40px | 16px | 顶部栏用户菜单 |
| xl | 80px | 28px | 资料页/认证 Logo |

圆角 full；文字头像：底 `--brand-50`、字 `--brand-500`、取用户名首字（中文取末字）；图片 `object-fit:cover`；可选 2px `--background-50` 边框用于叠加。

### 7.7 Toast 通知

位置：右上、距顶/右 16px；宽 **360px**；`--brand-radius-md`；`--shadow-lg`；`--card` 底；`p-4`；z-100；多条垂直间距 12px。类型：Success/Error/Warning/Info，左 20px 图标（色对应状态），左边框 4px 同色；标题 14px semibold `--text-800`，描述 13px `--text-500`，右上 X 16px `--text-400`；自动关闭 4000ms（**Error 不自动关**）；底部 2px 进度条；hover 暂停。动画：进入/退出 250ms `ease-out`/`ease-in` 右侧滑入滑出。

### 7.8 Tabs 标签页

`flex` 水平 + 底部 1px `--border`；标签间距 24~32px、`pb-3`；选中 14px semibold `--text-800` + **2px solid `--primary`** 下划线（`-mb-[1px]`）；未选中 14px normal `--text-500`；hover → `--text-800`；切换 150ms 淡入；内容区 `mt-6`。

### 7.9 Toggle 开关

**40×24px**、full 圆角；关 `--background-400`、开 `--primary`（或 `--state-success`）；滑块 20×20 白底 `--shadow-xs`，关 `left-0.5`、开 `translate-x-4`；过渡 150ms；禁用 `opacity:.5`。

### 7.10 Tree 树形

节点高 36px；展开图标 ChevronRight/Down 16px `--text-400`；节点图标 16px `--text-500`；文字 14px `--text-700`；选中 `--brand-500` semibold + `--brand-50` 底；hover `--accent`；**每级缩进 20px**；可选 1px dashed `--background-300` 连接线；半选（Checkbox）显示减号。

### 7.11 图标 Icon（Lucide）

| 令牌 | px | 用途 |
|------|----|------|
| `--icon-xs` | 14 | 行内/标签前缀 |
| `--icon-sm` | 16 | 按钮内、表单后缀、**导航项** |
| `--icon-md` | 20 | 默认尺寸、Toast、**顶部栏功能区** |
| `--icon-lg` | 24 | 导航大图标、列表前缀 |
| `--icon-xl` | 32 | 大按钮 |
| `--icon-empty` | 64 | 空状态插图 |

颜色：默认 `--icon`；次级/禁用 `--icon-muted`；按钮内继承文字色；状态图标用对应状态色。

### 7.12 卡片 Card / Alert / 骨架

- **Card**：`--card` 底、`--brand-radius-md`、1px `--border`、`--shadow-sm`、`p-6`；带标题：标题 16px semibold `--text-800` + `mb-6`，右侧 Ghost 操作；可交互 hover → `--shadow-md` + `cursor:pointer`（250ms）；统计卡：图标容器 48×48 `--brand-radius-md` 浅底，数字 28px bold `--text-900`，标签 14px `--text-500`，`p-5`。
- **Alert**：`--brand-radius-md`、`p-3`、`mb-6`、1px 对应状态色边框、`flex gap-3`；图标 18px 顶部对齐；Success/Error/Warning/Info 用 §1.1 状态色（**Warning 用 `--state-warning` 体系，勿用 amber 硬编码**，见偏差 #6）；可选关闭 X。
- **Skeleton**：`--background-200` 底、`--brand-radius-sm`（文字类 4px）、`animate-pulse` 1.5s。

---

## 8. 完整 CSS 变量块（直接复制）

> 合并「现存令牌 + 本文档补全部分」。把整段放入每个页面的 `<style id="theme-vars">`（或抽为共享 `tokens.css` 由所有页面 `<link>` 引入）。Tailwind `@theme inline` 的 `radius-sm/md/lg` 映射现已有效。

```css
:root {
  /* ===== PRIMITIVE: brand ===== */
  --brand-50: #e8f2ff;
  --brand-100: #cfe5ff;
  --brand-200: #9fcbff;
  --brand-300: #66abff;
  --brand-400: #2e8dff;
  --brand-500: #007aff; /* @primary */
  --brand-600: #0064d6;
  --brand-700: #004fad;
  --brand-800: #003b82;
  --brand-900: #00275a;

  /* ===== PRIMITIVE: background (10-step neutral) ===== */
  --background-50: #ffffff;
  --background-100: #f7f7fa;
  --background-200: #f2f2f7;
  --background-300: #e5e5ea;
  --background-400: #d1d1d6;
  --background-500: #aeaeb2;
  --background-600: #8e8e93;
  --background-700: #3a3a3c;
  --background-800: #1c1c1e;
  --background-900: #000000;

  /* ===== PRIMITIVE: text ===== */
  --text-50: #f5f5f7;
  --text-100: #e3e3e8;
  --text-200: #c7c7cc;
  --text-300: #aeaeb2;
  --text-400: #8e8e93;
  --text-500: #6e6e73;
  --text-600: #48484a;
  --text-700: #3c3c43;
  --text-800: #1d1d1f;
  --text-900: #000000;

  /* ===== PRIMITIVE: icon ===== */
  --icon-50: #f5f5f7;
  --icon-100: #e5e5ea;
  --icon-200: #d1d1d6;
  --icon-300: #c7c7cc;
  --icon-400: #aeaeb2;
  --icon-500: #8e8e93;
  --icon-600: #6e6e73;
  --icon-700: #48484a;
  --icon-800: #2c2c2e;
  --icon-900: #1d1d1f;

  /* ===== PRIMITIVE: state colors ===== */
  --state-success: #34c759;
  --state-success-dark: #30d158;
  --state-success-surface: #e9f9ee;
  --state-success-foreground: #ffffff;

  --state-error: #ff3b30;
  --state-error-dark: #ff453a;
  --state-error-surface: #ffecea;
  --state-error-foreground: #ffffff;

  /* ⚠ 补全：warning / info（现存缺失） */
  --state-warning: #ff9500;
  --state-warning-dark: #ff9f0a;
  --state-warning-surface: #fff4e6;
  --state-warning-foreground: #ffffff;

  --state-info: #5ac8fa;
  --state-info-dark: #64d2ff;
  --state-info-surface: #e6f7fe;
  --state-info-foreground: #ffffff;

  /* ===== SEMANTIC: surfaces / text ===== */
  --background: var(--background-50);
  --foreground: var(--text-800);
  --card: var(--background-50);
  --card-foreground: var(--text-800);
  --popover: var(--background-50);
  --popover-foreground: var(--text-900);

  --primary: var(--brand-500);
  --primary-foreground: var(--background-50);
  --secondary: var(--background-200);
  --secondary-foreground: var(--text-800);
  --muted: var(--background-200);
  --muted-foreground: var(--text-400);
  --accent: var(--background-100);
  --accent-foreground: var(--text-800);
  --destructive: var(--state-error);
  --destructive-foreground: var(--state-error-foreground);

  --success: var(--state-success);
  --success-foreground: var(--state-success-foreground);
  --warning: var(--state-warning);
  --warning-foreground: var(--state-warning-foreground);
  --info: var(--state-info);
  --info-foreground: var(--state-info-foreground);

  --border: var(--background-300);
  --input: var(--background-400);
  --ring: var(--brand-500);

  --icon: var(--icon-900);
  --icon-muted: var(--icon-500);

  /* ===== SEMANTIC: chart ===== */
  --chart-1: var(--state-success);
  --chart-2: var(--brand-500);
  --chart-3: var(--state-warning);
  --chart-4: #5856d6;
  --chart-5: #af52de;

  /* ===== SEMANTIC: sidebar ===== */
  --sidebar: var(--background-200);
  --sidebar-foreground: var(--text-800);
  --sidebar-primary: var(--brand-500);
  --sidebar-primary-foreground: var(--background-50);
  --sidebar-accent: var(--background-300);
  --sidebar-accent-foreground: var(--text-800);
  --sidebar-border: var(--background-300);
  --sidebar-ring: var(--brand-500);

  /* ===== TYPOGRAPHY ===== */
  --font-sans: "DM Sans", ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "PingFang SC", "Noto Sans CJK SC", sans-serif;
  --font-serif: var(--font-sans);
  --font-mono: "JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  --tracking-normal: 0em;

  --text-xs: 12px;
  --text-sm: 14px;
  --text-base: 16px;
  --text-lg: 18px;
  --text-xl: 20px;
  --text-2xl: 24px;
  --text-3xl: 30px;
  --text-4xl: 36px;

  --font-weight-normal: 400;
  --font-weight-medium: 500;
  --font-weight-semibold: 600;
  --font-weight-bold: 700;

  --leading-tight: 1.2;
  --leading-normal: 1.5;
  --leading-relaxed: 1.7;

  /* ===== RADIUS（⚠ 补全 sm/md/lg/full；保留 --radius 兼容旧页） ===== */
  --brand-radius-sm: 8px;
  --brand-radius-md: 12px;
  --brand-radius-lg: 19.2px;
  --brand-radius-full: 9999px;
  --radius: 1.2rem; /* = 19.2px，仅旧页面 */

  /* ===== SHADOW ===== */
  --shadow-2xs: 0 1px 2px -1px rgba(0, 0, 0, 0.04);
  --shadow-xs: 0 1px 2px 0 rgba(0, 0, 0, 0.04);
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.05), 0 1px 3px -1px rgba(0, 0, 0, 0.05);
  --shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.06), 0 1px 2px -1px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 8px -2px rgba(0, 0, 0, 0.06), 0 2px 4px -2px rgba(0, 0, 0, 0.05);
  --shadow-lg: 0 8px 24px -8px rgba(0, 0, 0, 0.08), 0 4px 8px -4px rgba(0, 0, 0, 0.05);
  --shadow-xl: 0 16px 40px -10px rgba(0, 0, 0, 0.10), 0 8px 16px -8px rgba(0, 0, 0, 0.06);
  --shadow-2xl: 0 24px 64px -12px rgba(0, 0, 0, 0.12);

  /* ===== SPACING（px 精确值；--spacing 为 Tailwind 基数） ===== */
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;
  --space-12: 48px;
  --space-16: 64px;
  --spacing: 0.24rem;

  /* ===== LAYOUT ===== */
  --sidebar-w: 256px;
  --sidebar-w-collapsed: 64px;
  --topbar-h: 56px;
  --content-max-w: 1408px;
  --content-padding: 32px;
  --brand-logo-size: 32px;

  /* ===== MOTION ===== */
  --duration-instant: 0ms;
  --duration-fast: 150ms;
  --duration-normal: 250ms;
  --duration-slow: 400ms;
  --duration-slower: 600ms;
  --ease-default: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-in: cubic-bezier(0.4, 0, 1, 1);
  --ease-out: cubic-bezier(0, 0, 0.2, 1);
  --ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);
  --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);

  /* focus ring（输入框/按钮/可聚焦元素复用） */
  --focus-ring: 0 0 0 2px color-mix(in srgb, var(--ring) 18%, transparent);

  /* ===== ROLE BADGE ===== */
  --badge-superadmin-bg: var(--state-error-surface);
  --badge-superadmin-fg: var(--state-error);
  --badge-orgadmin-bg: var(--brand-50);
  --badge-orgadmin-fg: var(--brand-500);
  --badge-teamadmin-bg: #dbeafe;        /* 固定蓝，非 brand 阶 */
  --badge-teamadmin-fg: #2563eb;        /* 固定蓝，非 brand 阶 */
  --badge-groupadmin-bg: var(--state-success-surface);
  --badge-groupadmin-fg: var(--state-success);
  --badge-member-bg: var(--background-200);
  --badge-member-fg: var(--text-500);
  --badge-public-bg: var(--background-200);
  --badge-public-fg: var(--text-400);

  /* ===== ICON SIZES ===== */
  --icon-xs: 14px;
  --icon-sm: 16px;
  --icon-md: 20px;
  --icon-lg: 24px;
  --icon-xl: 32px;
  --icon-empty: 64px;
}

/* ===== DARK（原始色板继承，仅重映射语义 + 阴影） ===== */
.dark {
  --background: var(--background-900);
  --foreground: var(--text-50);
  --card: var(--background-800);
  --card-foreground: var(--text-50);
  --popover: var(--background-700);
  --popover-foreground: var(--text-50);

  --primary: var(--brand-400);
  --primary-foreground: var(--background-900);
  --secondary: var(--background-800);
  --secondary-foreground: var(--text-50);
  --muted: var(--background-800);
  --muted-foreground: var(--text-400);
  --accent: var(--background-700);
  --accent-foreground: var(--text-50);
  --destructive: var(--state-error-dark);
  --destructive-foreground: var(--state-error-foreground);

  --success: var(--state-success-dark);
  --success-foreground: var(--state-success-foreground);
  --warning: var(--state-warning-dark);
  --warning-foreground: var(--state-warning-foreground);
  --info: var(--state-info-dark);
  --info-foreground: var(--state-info-foreground);

  --border: var(--background-700);
  --input: var(--background-700);
  --ring: var(--brand-400);

  --icon: var(--icon-50);
  --icon-muted: var(--icon-500);

  --chart-1: var(--state-success-dark);
  --chart-2: var(--brand-400);
  --chart-3: #ff9f0a;
  --chart-4: #5e5ce6;
  --chart-5: #bf5af2;

  --sidebar: var(--background-800);
  --sidebar-foreground: var(--text-50);
  --sidebar-primary: var(--brand-400);
  --sidebar-primary-foreground: var(--background-900);
  --sidebar-accent: var(--background-700);
  --sidebar-accent-foreground: var(--text-50);
  --sidebar-border: var(--background-700);
  --sidebar-ring: var(--brand-400);

  --shadow-2xs: 0 1px 2px -1px rgba(0, 0, 0, 0.30);
  --shadow-xs: 0 1px 2px 0 rgba(0, 0, 0, 0.30);
  --shadow-sm: 0 1px 2px 0 rgba(0, 0, 0, 0.36), 0 1px 3px -1px rgba(0, 0, 0, 0.36);
  --shadow: 0 2px 4px -1px rgba(0, 0, 0, 0.40), 0 1px 2px -1px rgba(0, 0, 0, 0.36);
  --shadow-md: 0 4px 8px -2px rgba(0, 0, 0, 0.44), 0 2px 4px -2px rgba(0, 0, 0, 0.36);
  --shadow-lg: 0 8px 24px -8px rgba(0, 0, 0, 0.50), 0 4px 8px -4px rgba(0, 0, 0, 0.40);
  --shadow-xl: 0 16px 40px -10px rgba(0, 0, 0, 0.55), 0 8px 16px -8px rgba(0, 0, 0, 0.44);
  --shadow-2xl: 0 24px 64px -12px rgba(0, 0, 0, 0.60);
}
```

**Tailwind v4 `@theme inline` 片段**（已有，确认 radius 现已有效）：

```css
@theme inline {
  /* …color 映射同现存… */
  --radius-sm: var(--brand-radius-sm);
  --radius-md: var(--brand-radius-md);
  --radius-lg: var(--brand-radius-lg);
}
```

---

## 8b. 主题切换规范（防闪烁）

在 `<head>` 最前内联脚本，早于样式与 body 解析：

```html
<script>
  (function () {
    try {
      var t = localStorage.getItem('theme');
      if (!t) t = matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
      document.documentElement.classList.add(t); // 加到 <html>
    } catch (e) { document.documentElement.classList.add('light'); }
  })();
</script>
```

- `<html>` 默认 `class="light"`（兜底）；脚本覆盖为实际主题。
- 切换：用户点击 Sun/Moon → 切换 `documentElement.classList` 的 `light`/`dark` → `localStorage.setItem('theme', ...)`。
- 全局生效：所有页面共享同一 `:root`/`.dark` 令牌，无需逐页处理。

---

## 9. 与现存原型的偏差修正清单（原型构建师执行项）

| # | 偏差（现存） | 规范/本文档要求 | 修改点 |
|---|------------|----------------|--------|
| 1 | `admin-shell.html` 侧边栏 `w-[240px]`、内容区 `ml-[240px]` | 展开 **256px** / 收起 64px | 改为 `var(--sidebar-w)` 与 `var(--sidebar-w-collapsed)`；收起态内容区 `margin-left: var(--sidebar-w-collapsed)` |
| 2 | `admin-shell.html` 内容区 `max-w-[1200px]`（且网页端规范 6.2.2 自相矛盾写 1200） | **1408px** | 改为 `max-width: var(--content-max-w)`；以 1408 为准，规范 6.2.2 的 1200 作废 |
| 3 | `:root` 仅定义 `--radius`，但 Tailwind `@theme` 引用 `--brand-radius-sm/md/lg` **未定义**（半径失效） | 补全 sm8/md12/lg19.2/full9999 | 第 8 节 CSS 已补全；新页面统一用语义半径 |
| 4 | 现存只定义 state-success / state-error；组件引用 `--state-warning`/`--state-info` 报错 | 补全 warning/info 原语 + 语义（含 dark） | 第 8 节已补全；组件可直接用 |
| 5 | `admin-shell.html` 侧边栏 Logo 硬编码 `bg-[var(--brand-500)]`（line 9） | 跟随主题：dark 应为 brand-400 | 改为 `bg-[var(--sidebar-primary)]`（或 `--primary`） |
| 6 | `account.html` Alert SVG 硬编码 `stroke="#007aff"`（line 354-356）；网页端 5.10.2 Alert Warning 用 amber `#FBBF24/#D97706/#92400E` 硬编码 | 全部走语义令牌 | Alert 图标改 `currentColor` 或 `var(--brand-500)`；Warning 统一用 `--state-warning` 体系（surface/fg） |
| 7 | 输入框圆角网页端规范 5.2.2 写 md/lg=12px | 本文档统一 **8px**（`--brand-radius-sm`） | 输入框（含 select/textarea）圆角改用 `--brand-radius-sm` |
| 8 | 认证页卡片圆角现存/规范 3.4 写 20px | 大容器用 `--brand-radius-lg`(19.2px) | 认证卡片圆角改 `--brand-radius-lg`（与规范体系一致；20px 视为近似） |
| 9 | 组件类无样式：`account.html` 用 `.sidebar/.card/.btn-danger/.badge` 等但文件内**无对应 CSS**（仅 `theme-vars`） | 组件样式须基于令牌实现 | 建立共享组件样式表（或逐页实现），所有类映射到第 7 节规范，禁止新增硬编码色 |
| 10 | 深色模式焦点环/主色：令牌层已正确（`--primary`=brand-400），但页面若用 `--brand-500` 硬编码则深色下不切换 | 一律用语义令牌 | 全文检索 `--brand-500` / `#007aff` 硬编码并替换为 `--primary`/`--sidebar-primary` |
| 11 | 现存 `html` 默认 `class="light"` 且未内联防闪烁脚本 | 加 `<head>` 内联脚本（8b） | 每页 `<head>` 顶部加入主题脚本，避免深色闪白 |

**执行优先级**：#3（半径失效，阻断新页面）、#4（warning/info 缺失，阻断引用）、#1/#2（布局硬偏差）、#5/#6/#10（硬编码，违反铁律）、#11（闪烁）、#7/#8（圆角统一）、#9（组件 CSS 落地）。

---

## 10. 动效 / 无障碍速查（补充）

- **时长**：instant 0 / fast 150 / normal 250 / slow 400 / slower 600（ms）。
- **缓动**：默认 `cubic-bezier(0.4,0,0.2,1)`；进入 `ease-out`；退出 `ease-in`；状态切换 `ease-in-out`；微交互可 `spring`。
- **焦点环**：所有可交互元素 `:focus-visible { box-shadow: var(--focus-ring); }`；Tab 顺序与视觉一致；模态焦点陷阱 + 归还。
- **对比度**：正文 ≥4.5:1、大字/图标 ≥3:1（WCAG AA）；禁用态可辨识即可。
- **ARIA**：语义标签（`nav/main/section/button`）；动态内容 `aria-live`（Toast、行内错误）；图片 `alt`；装饰图空 `alt`。
- **触控目标**：移动端适配最小 44×44px。

---

## 11. 错误码 → 提示映射（交互行为，构建师实现）

| 错误码 | 行为 | 提示/跳转 |
|--------|------|----------|
| 400 | 行内错误 | 字段级红框 + 下方错误文字（`--destructive`），不弹窗 |
| 401 | 过期/未登录 | 跳登录页（清除会话），可带 `redirect` 回跳 |
| 429 | 限流 | Toast Warning “操作过于频繁，请稍后再试” |
| 500 | 通用服务端错误 | Toast Error “服务异常，请稍后重试”（Error 不自动关闭） |
| 其他 4xx/5xx | 通用 | 对应语义 Toast；字段错误优先行内 |

---

> **交付说明**：本文档即原型构建师的完整依据。第 8 节 CSS 可直接复制；第 7 节为组件速查；第 9 节为对现存原型的修正清单。品牌已锁定 Pinguo，无需选型。
