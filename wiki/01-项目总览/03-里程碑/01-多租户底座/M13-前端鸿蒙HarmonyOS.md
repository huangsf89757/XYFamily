# M13 — 前端鸿蒙 HarmonyOS App

> XYFamily 鸿蒙 HarmonyOS（ArkTS + ArkUI）原生 App 开发，移动端三阶段之第三阶段。

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
| 关联标签 | 核心文档、里程碑、前端、鸿蒙、HarmonyOS、移动端 |
| 关联目录 | 01-项目总览/里程碑/V1.0.0 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建：新增鸿蒙 HarmonyOS 里程碑，移动端第三阶段 | CatPaw |

---

## 一、里程碑概览

| 项目 | 内容 |
|------|------|
| 里程碑 ID | M13 |
| 所属阶段 | Phase 5: 延后功能（移动端第三阶段） |
| 开始时间 | 2027-01 W1 |
| 结束时间 | 2027-01 W2 |
| 当前状态 | 🔲 待排期 |
| 前置依赖 | M12 — 生产部署 + 上线（系统已稳定运行） |
| 后续里程碑 | M14 — P2 功能开发 |

### 目标

完成鸿蒙 HarmonyOS 原生 App 的核心功能开发，覆盖登录认证、个人中心、组织/团队/小组管理全流程。作为移动端三阶段开发的第三阶段，在 iOS 和 Android 端上线稳定后启动，复用已验证的 API 对接方案与交互设计。

### 移动端三阶段规划

| 阶段 | 里程碑 | 平台 | 时间 | 状态 |
|------|--------|------|------|------|
| 第一阶段 | M9 | iOS | 2026-11 W2-W3 | 待开始 |
| 第二阶段 | M10 | Android | 2026-11 W4 - 2026-12 W1 | 待开始 |
| ✅ 第三阶段 | M13 | 鸿蒙 HarmonyOS | 2027-01 W1-W2 | 待排期 |

---

## 二、任务清单

| # | 任务 | 优先级 | 状态 | 验收标准 | 说明 |
|---|------|--------|------|---------|------|
| 1 | 项目初始化 | P1 | 🔲 | DevEco Studio 编译运行成功 | ArkTS / ArkUI / API 12+ |
| 2 | 统一 API 调用层 | P1 | 🔲 | Token 自动刷新 + 统一错误处理 | @ohos.net.http 封装 |
| 3 | 登录 + 认证流程 | P1 | 🔲 | 登录成功跳转主页 | @ohos.security.huks 存储 Refresh Token |
| 4 | 个人中心 | P1 | 🔲 | 信息查看/修改 + 密码修改 | - |
| 5 | 组织管理 | P1 | 🔲 | CRUD + 成员管理 + 角色分配 | - |
| 6 | 团队/小组管理 | P1 | 🔲 | 层级展示 + CRUD + 成员管理 | - |
| 7 | 邀请管理 | P1 | 🔲 | 查看待确认邀请 + 接受/拒绝 | - |
| 8 | 鸿蒙单元测试 | P2 | 🔲 | 核心逻辑测试通过 | @ohos/hypium |
| 9 | 华为应用市场上架 | P1 | 🔲 | 提交审核材料 | 截图、描述、隐私政策 |

---

## 三、交付产物

| 产物 | 路径 | 说明 |
|------|------|------|
| 鸿蒙源码 | `code/harmony/` | ArkTS / ArkUI |
| 鸿蒙 API 层 | `code/harmony/entry/src/main/ets/network/` | @ohos.net.http 封装 |
| 鸿蒙页面 | `code/harmony/entry/src/main/ets/pages/` | ArkUI 页面 |
| 鸿蒙测试 | `code/harmony/entry/src/test/` | 单元测试 |

---

## 四、技术要点

### 4.1 鸿蒙技术栈

| 类别 | 选型 | 版本 |
|------|------|------|
| 语言 | ArkTS | - |
| UI 框架 | ArkUI (声明式) | - |
| 网络 | @ohos.net.http | API 12+ |
| 本地存储 | @ohos.data.preferences | - |
| 安全存储 | @ohos.security.huks | Refresh Token 存储 |
| 依赖注入 | - | 手动管理 / 状态管理 |

### 4.2 Token 管理

| 平台 | Access Token | Refresh Token |
|------|-------------|---------------|
| 鸿蒙 | 内存（TokenManager） | @ohos.security.huks |

### 4.3 统一 API 层

- **Token 自动刷新**：401 时自动调用 `/auth/refresh`，成功后重试原请求
- **多租户 Header**：自动附加 `X-Organization-ID` / `X-Team-ID` / `X-Group-ID`
- **统一错误码处理**：错误码 → 中文文案映射 → promptAction.showToast 展示

### 4.4 与 iOS/Android 端的对齐

| 维度 | iOS (M9) | Android (M10) | 鸿蒙 (M13) |
|------|----------|---------------|------------|
| API 接口 | 相同 | 相同 | 相同 |
| 交互设计 | 参考基准 | 对齐 iOS | 对齐 iOS/Android，适配 ArkUI |
| Token 策略 | Keychain | EncryptedSharedPreferences | HUKS |
| 错误处理 | Alert | Snackbar | Toast |
| UI 框架 | SwiftUI | Jetpack Compose | ArkUI (声明式) |

---

## 五、页面清单

| # | 页面 | 鸿蒙 (ArkUI) | 说明 |
|---|------|-------------|------|
| 1 | 登录页 |LoginPage | 手机号/邮箱 + 密码 |
| 2 | 主页 | HomePage | 组织列表 |
| 3 | 个人中心 | AccountPage | 信息 + 密码 |
| 4 | 组织详情 | OrgDetailPage | 团队列表 + 成员 |
| 5 | 团队详情 | TeamDetailPage | 小组列表 + 成员 |
| 6 | 小组详情 | GroupDetailPage | 成员管理 |
| 7 | 邀请管理 | InvitationsPage | 待确认邀请 |
| 8 | 成员管理 | MembersPage | 列表 + 角色分配 |

---

## 六、风险评估

| 风险 | 等级 | 应对措施 |
|------|------|---------|
| 鸿蒙生态成熟度 | 中 | 关注 HarmonyOS NEXT 更新，使用稳定 API |
| ArkTS 学习曲线 | 中 | 团队培训 + 参考 iOS/Android 已有实现 |
| 鸿蒙设备适配 | 低 | 使用 ArkUI 自适应布局 |
| 华为应用市场审核 | 低 | 提前准备审核材料 |

---

## 七、关联文档

- [里程碑-多租户底座](./多租户底座.md)
- [M9 — iOS 原生 App](./M9-前端iOS原生App.md)
- [M10 — Android 原生 App](./M10-前端Android原生App.md)
- [M12 — 生产部署 + 上线](./M12-生产部署与上线.md)
- [M14 — P2 功能开发](./M14-P2功能开发.md)
- [前端通用代码规范](../../02-标准与规范/03-分端代码规范/07-前端通用代码规范.md)
- [接口总览](../../../04-接口文档/接口文档.md)
