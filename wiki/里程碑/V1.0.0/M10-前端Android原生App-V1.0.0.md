# M10 — 前端 Android 原生 App

> XYFamily Android（Kotlin + Jetpack Compose）原生 App 开发，移动端三阶段之第二阶段。

---

## 文档信息

| 项目 | 内容 |
|------|------|
| 文档密级 | 内部 |
| 文档版本 | V1.0.0 |
| 编写人 | CatPaw |
| 审核人 | - |
| 生效时间 | 2026-07-12 |
| 废弃时间 | - |
| 关联标签 | 核心文档、里程碑、前端、Android、移动端 |
| 关联目录 | 00-项目总览/00.03-项目里程碑 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建：从原 M9（iOS+Android）拆分为独立 Android 里程碑 | CatPaw |

---

## 一、里程碑概览

| 项目 | 内容 |
|------|------|
| 里程碑 ID | M10 |
| 所属阶段 | Phase 3: 前端开发（移动端第二阶段） |
| 开始时间 | 2026-11 W4 |
| 结束时间 | 2026-12 W1 |
| 当前状态 | 🔲 待开始 |
| 前置依赖 | M9 — iOS 原生 App（复用 API 层设计与交互模式） |
| 后续里程碑 | M11 — 集成测试 + 安全审查 |

### 目标

完成 Android 原生 App 的核心功能开发，覆盖登录认证、个人中心、组织/团队/小组管理全流程。作为移动端三阶段开发的第二阶段，复用 iOS 端验证过的 API 对接方案与交互设计，快速实现 Android 对等版本。

### 移动端三阶段规划

| 阶段 | 里程碑 | 平台 | 时间 |
|------|--------|------|------|
| 第一阶段 | M9 | iOS | 2026-11 W2-W3 |
| ✅ 第二阶段 | M10 | Android | 2026-11 W4 - 2026-12 W1 |
| 第三阶段 | M13 | 鸿蒙 HarmonyOS | 2027-01 W1-W2（上线后） |

---

## 二、任务清单

| # | 任务 | 优先级 | 状态 | 验收标准 | 说明 |
|---|------|--------|------|---------|------|
| 1 | 项目初始化 | P0 | 🔲 | Android Studio 编译运行成功 | Kotlin / Compose / Retrofit / Room / Hilt |
| 2 | 统一 API 调用层 | P0 | 🔲 | Token 自动刷新 + 统一错误处理 | Retrofit + OkHttp 拦截器 |
| 3 | 登录 + 认证流程 | P0 | 🔲 | 登录成功跳转主页 | EncryptedSharedPreferences 存储 |
| 4 | 个人中心 | P0 | 🔲 | 信息查看/修改 + 密码修改 | - |
| 5 | 组织管理 | P0 | 🔲 | CRUD + 成员管理 + 角色分配 | - |
| 6 | 团队/小组管理 | P0 | 🔲 | 层级展示 + CRUD + 成员管理 | - |
| 7 | 邀请管理 | P0 | 🔲 | 查看待确认邀请 + 接受/拒绝 | - |
| 8 | Android 单元测试 | P1 | 🔲 | 核心逻辑测试通过 | JUnit + Espresso |
| 9 | Google Play 上架准备 | P1 | 🔲 | 提交审核材料 | 截图、描述、隐私政策 |

---

## 三、交付产物

| 产物 | 路径 | 说明 |
|------|------|------|
| Android 源码 | `code/android/` | Kotlin / Compose / Retrofit |
| Android API 层 | `code/android/app/src/main/java/com/xyfamily/data/` | Retrofit 接口 |
| Android UI | `code/android/app/src/main/java/com/xyfamily/presentation/` | Compose UI |
| Android 测试 | `code/android/app/src/test/` | 单元测试 |

---

## 四、技术要点

### 4.1 Android 技术栈

| 类别 | 选型 | 版本 |
|------|------|------|
| 语言 | Kotlin | 1.9+ |
| UI 框架 | Jetpack Compose | latest |
| 网络 | Retrofit + OkHttp | - |
| 本地存储 | Room | - |
| 安全存储 | EncryptedSharedPreferences | - |
| 依赖注入 | Hilt (Dagger) | - |

### 4.2 Token 管理

| 平台 | Access Token | Refresh Token |
|------|-------------|---------------|
| Android | 内存（TokenManager） | EncryptedSharedPreferences |

### 4.3 统一 API 层

- **Token 自动刷新**：401 时自动调用 `/auth/refresh`，成功后重试原请求
- **多租户 Header**：自动附加 `X-Organization-ID` / `X-Team-ID` / `X-Group-ID`
- **统一错误码处理**：错误码 → 中文文案映射 → Toast/Snackbar 展示

### 4.4 与 iOS 端的对齐

| 维度 | iOS (M9) | Android (M10) |
|------|----------|---------------|
| API 接口 | 相同（共用后端接口总览） | 相同 |
| 交互设计 | 参考基准 | 对齐 iOS，适配 Material Design |
| Token 策略 | Keychain | EncryptedSharedPreferences |
| 错误处理 | Alert | Snackbar |

---

## 五、页面清单

| # | 页面 | Android (Compose) | 说明 |
|---|------|-------------------|------|
| 1 | 登录页 | LoginScreen | 手机号/邮箱 + 密码 |
| 2 | 主页 | HomeScreen | 组织列表 |
| 3 | 个人中心 | AccountScreen | 信息 + 密码 |
| 4 | 组织详情 | OrgDetailScreen | 团队列表 + 成员 |
| 5 | 团队详情 | TeamDetailScreen | 小组列表 + 成员 |
| 6 | 小组详情 | GroupDetailScreen | 成员管理 |
| 7 | 邀请管理 | InvitationsScreen | 待确认邀请 |
| 8 | 成员管理 | MembersScreen | 列表 + 角色分配 |

---

## 六、风险评估

| 风险 | 等级 | 应对措施 |
|------|------|---------|
| Android 设备碎片化 | 中 | 优先适配主流分辨率，使用 Compose 自适应布局 |
| EncryptedSharedPreferences 兼容 | 低 | 使用 AndroidX Security 最新稳定版 |
| Google Play 审核周期 | 中 | 提前准备审核材料，预留 1 周审核时间 |

---

## 七、关联文档

- [里程碑总览](./项目里程碑-V1.0.0.md)
- [M9 — iOS 原生 App](./M9-前端iOS原生App-V1.0.0.md)
- [M11 — 集成测试 + 安全审查](./M11-集成测试与安全审查-V1.0.0.md)
- [M13 — 鸿蒙 HarmonyOS App](./M13-前端鸿蒙HarmonyOS-V1.0.0.md)
- [安卓 Kotlin 代码规范](../../04-开发规范与编码手册/04.02-分端代码规范/安卓Kotlin代码规范-V1.0.0.md)
- [前端通用代码规范](../../04-开发规范与编码手册/04.02-分端代码规范/前端通用代码规范-V1.0.0.md)
- [接口总览](../../05-接口与模块落地文档/05.01-接口总览/接口总览-V1.0.0.md)
