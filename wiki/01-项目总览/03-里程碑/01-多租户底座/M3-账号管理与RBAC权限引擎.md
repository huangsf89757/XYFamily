# M3 — 账号管理 + RBAC 权限引擎

> XYFamily 账号管理模块与 RBAC 权限引擎开发，包含个人信息、密码修改、注销/恢复、权限中间件、Redis 缓存。

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
| 关联标签 | 核心文档、里程碑、后端、账号、RBAC |
| 关联目录 | 01-项目总览/里程碑/V1.0.0 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建 | CatPaw |

---

## 一、里程碑概览

| 项目 | 内容 |
|------|------|
| 里程碑 ID | M3 |
| 所属阶段 | Phase 2: 后端核心 |
| 开始时间 | 2026-09 W1 |
| 结束时间 | 2026-09 W2 |
| 当前状态 | 🔲 待开始 |
| 前置依赖 | M2 — 核心认证能力（JWT 基础设施） |
| 后续里程碑 | M4 — 组织管理 |

### 目标

完成账号自管理功能（个人信息/密码/注销），并实现 RBAC 权限引擎（权限中间件 + Redis 缓存 + 层级回溯 + 缓存失效），为后续组织/团队/小组模块提供权限校验基础设施。

### V2.0.0 变更说明

> **原计划中 RBAC 位于 M6（审计日志之后），V2.0.0 将其前置至 M3。** 原因：M4 组织管理模块的成员角色分配/降级、成员移除等操作均依赖权限中间件校验，RBAC 必须在 M4 之前完成。

---

## 二、任务清单

### 2.1 账号管理

| # | 任务 | 优先级 | 状态 | 验收标准 | 涉及表 |
|---|------|--------|------|---------|--------|
| 1 | `GET /account/profile` | P0 | 🔲 | 返回个人信息 | accounts |
| 2 | `PUT /account/profile` | P0 | 🔲 | 修改昵称/头像 | accounts |
| 3 | `PUT /account/password` | P0 | 🔲 | 旧密码校验 + 新密码加密 + 当前 Token 黑名单 | accounts, Redis |
| 4 | `POST /account/deactivate` | P0 | 🔲 | 状态 → deactivating，member role 清空，30 天宽限期 | accounts, *_members |
| 5 | `POST /account/undeactivate` | P0 | 🔲 | 状态 deactivating → active（宽限期内） | accounts |
| 6 | 匿名化定时任务 | P1 | 🔲 | 到期账号 PII 替换 + role 清空 | accounts, *_members |

### 2.2 RBAC 权限引擎

| # | 任务 | 优先级 | 状态 | 验收标准 | 说明 |
|---|------|--------|------|---------|------|
| 7 | 权限中间件 | P0 | 🔲 | 请求到达 Handler 前完成权限校验 | 缓存优先 + DB 兜底 |
| 8 | 层级回溯 | P0 | 🔲 | L5→L1 权限继承正确 | SuperAdmin 拥有全部权限 |
| 9 | Redis 缓存 | P0 | 🔲 | Key: `perm:account:{id}:org:{org_id}`，TTL 24h | 减少 DB 查询 |
| 10 | 缓存失效 | P0 | 🔲 | 角色变更/成员变更时主动删除缓存 | 保证一致性 |
| 11 | 权限点校验 | P0 | 🔲 | `permissionMW("perm.code")` 声明式校验 | Gin 中间件 |

---

## 三、交付产物

| 产物 | 路径 | 说明 |
|------|------|------|
| 账号 Handler | `code/backend/internal/handler/account/` | 5 个接口 Handler |
| 账号 Service | `code/backend/internal/service/account/` | 业务逻辑 |
| 账号 Repository | `code/backend/internal/repository/account/` | 数据访问 |
| 权限中间件 | `code/backend/internal/middleware/permission.go` | 权限校验 |
| RBAC Service | `code/backend/internal/service/rbac/` | 角色/权限/缓存逻辑 |
| RBAC Repository | `code/backend/internal/repository/rbac/` | roles/permissions 数据访问 |
| 匿名化定时任务 | `code/backend/internal/service/account/anonymize.go` | cron 定时执行 |
| 单元测试 | `code/backend/test/account_test.go` | 账号模块测试 |
| 单元测试 | `code/backend/test/rbac_test.go` | 权限引擎测试 |

---

## 四、技术要点

### 4.1 RBAC 权限校验流程

```
请求 → JWT验证 → Token黑名单 → Membership验证(层级回溯) → 权限点校验 → Handler
```

1. 从 JWT Claims 获取 `account_id`
2. 从 Header 提取 `X-Organization-ID` / `X-Team-ID` / `X-Group-ID`
3. 查询 Redis 缓存 → miss 时查 DB → 写入缓存
4. 校验角色是否拥有该权限点（含层级继承）
5. 通过 → 继续；不通过 → 返回 403

### 4.2 缓存设计

| 缓存 Key | TTL | 失效时机 |
|----------|-----|---------|
| `perm:account:{id}:org:{org_id}` | 24h | 角色变更、成员邀请/移除 |
| `perm:account:{id}:team:{team_id}` | 24h | 同上 |
| `perm:account:{id}:group:{group_id}` | 24h | 同上 |

### 4.3 注销流程

```
用户申请注销 → 状态 active → deactivating → member role 全部清空 → 30 天宽限期
  ├── 宽限期内恢复 → 状态恢复 active，role 需重新分配
  └── 宽限期到期 → 匿名化定时任务 → PII 替换 + 状态 deactivated
```

### 4.4 修改密码副作用

修改密码后，当前 Access Token 的 jti 写入 Redis 黑名单，强制用户重新登录。所有 sessions 中的 refresh_token 标记为 revoked。

---

## 五、接口清单（本里程碑）

| # | 方法 | 路径 | 认证 | 权限点 | HTTP |
|---|------|------|------|--------|------|
| 1 | GET | `/api/v1/account/profile` | ✅ | `account.profile.read` | 200 |
| 2 | PUT | `/api/v1/account/profile` | ✅ | `account.profile.update` | 200 |
| 3 | PUT | `/api/v1/account/password` | ✅ | `account.password.update` | 200 |
| 4 | POST | `/api/v1/account/deactivate` | ✅ | `account.deactivate` | 200 |
| 5 | POST | `/api/v1/account/undeactivate` | ✅ | `account.undeactivate` | 200 |

---

## 六、风险评估

| 风险 | 等级 | 应对措施 |
|------|------|---------|
| RBAC 层级回溯复杂度 | 高 | 提前设计 + 单测覆盖所有回溯路径（L0→L5） |
| Redis 缓存一致性 | 高 | 角色/成员变更时必须主动失效缓存 |
| 注销恢复后 role 丢失 | 中 | 文档明确说明：恢复后 role 需重新分配 |
| 匿名化任务执行冲突 | 低 | 分布式锁保证单实例执行 |

---

## 七、关联文档

- [里程碑总览](./项目里程碑)
- [M2 — 核心认证能力](./M2-核心认证能力.md)
- [接口总览 — 账号模块](../../05-接口文档/05.01-接口总览/接口总览)
- [中间件专项方案 — 权限引擎](../../03-架构与方案设计/04-链路实现/01-中间件链专项方案.md)
- [核心技术专项方案](../../03-架构与方案设计/04-链路实现/链路实现.md)

## 关联文档


> 以下为知识图谱自动推荐的交叉引用，建议人工审阅确认后保留。

- [架构评审记录](../../../03-架构与方案设计/审核记录/架构评审记录.md) — 共享术语：rbac、接口、权限、缓存、账号（置信度 0.75）
- [02-账号接口](../../../03-架构与方案设计/03-数据模型与契约/02-接口设计/02-账号接口.md) — 共享术语：接口、账号（置信度 0.75）
