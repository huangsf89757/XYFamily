# Git版本管控规范

> XYFamily Git 分支策略、提交信息规范、Code Review 规则、版本标签规范。

---

## 文档信息

| 项目 | 内容 |
|------|------|
| 文档密级 | 内部 |
| 文档版本 | V1.0.0 |
| 编写人 | ClaudeCode |
| 审核人 | - |
| 生效时间 | 2026-07-12 |
| 废弃时间 | - |
| 关联标签 | 开发规范、Git、核心文档 |
| 关联目录 | 02-标准与规范/04-开发规范 |

## 变更记录


| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-19 | 文档新编 | ClaudeCode |

---

## 一、分支策略

### 1.1 分支模型

```
main          ────────────────────────────────────────→ (生产)
  \
   develop     ──────────────────────────────────────→ (集成)
     \                                                \
      feature/auth ──────→                              \
      feature/rbac ────────→                             \
      hotfix/token ─────────────→                       \
```

| 分支 | 用途 | 保护规则 |
|------|------|----------|
| `main` | 生产分支，稳定代码 | 禁止直接推送，必须通过 PR 合并 |
| `develop` | 集成分支，功能合并 | 禁止直接推送，必须通过 PR 合并 |
| `feature/*` | 功能开发 | 从 develop 拉取，完成后合并到 develop |
| `hotfix/*` | 紧急修复 | 从 main 拉取，修复后合并到 main + develop |

### 1.2 分支命名

| 类型 | 格式 | 示例 |
|------|------|------|
| 功能开发 | `feature/模块名-功能名` | `feature/auth-login`、`feature/rbac-permission` |
| 紧急修复 | `hotfix/问题描述` | `hotfix/token-blacklist` |
| 代码重构 | `refactor/模块名` | `refactor/service-layer` |
| 文档更新 | `docs/模块名` | `docs/prd-update` |

---

## 二、提交信息规范

### 2.1 Conventional Commits

```
<type>(<scope>): <subject>

<body>

<footer>
```

### 2.2 Type 类型

| 类型 | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat(auth): add login API` |
| `fix` | 修复 | `fix(rbac): fix permission cache invalidation` |
| `refactor` | 重构 | `refactor(service): extract org service` |
| `perf` | 性能优化 | `perf(db): add index for organization_members` |
| `docs` | 文档 | `docs(prd): update PRD-V1.0.0` |
| `test` | 测试 | `test(auth): add unit tests for login handler` |
| `ci` | CI 配置 | `ci: update GitHub Actions` |
| `chore` | 其他 | `chore: update dependencies` |

### 2.3 Scope 范围

| Scope | 对应模块 |
|-------|---------|
| `auth` | 认证模块 |
| `account` | 账号模块 |
| `org` | 组织模块 |
| `team` | 团队模块 |
| `group` | 小组模块 |
| `rbac` | RBAC 权限模块 |
| `audit` | 审计模块 |
| `admin` | 超级管理员模块 |
| `db` | 数据库 |
| `infra` | 基础设施 |
| `middleware` | 中间件 |

### 2.4 示例

```
feat(auth): add login API with JWT support

- Add POST /api/v1/auth/login/password
- Add POST /api/v1/auth/login/sms-code
- Implement JWT token generation

Refs: PRD-AUTH-001
```

```
fix(rbac): fix permission cache invalidation on role change

Permission cache was not invalidated when role was changed.
Now invalidates all related cache entries.

Fixes: #123
```

---

## 三、Code Review 规则

### 3.1 审查触发条件

- 所有合并到 `main` 和 `develop` 的 PR 必须经过 Code Review
- 安全敏感代码（认证、权限、密码、Token）必须经过安全审查

### 3.2 审查清单

| 类别 | 检查项 |
|------|--------|
| **安全性** | 无硬编码密钥、输入已验证、SQL 注入防护、权限校验完整 |
| **正确性** | 业务逻辑正确、边界条件处理、错误处理完整 |
| **性能** | 无 N+1 查询、缓存策略正确、数据库索引合理 |
| **可读性** | 命名清晰、函数简短、无深嵌套、注释充分 |
| **测试** | 新增功能有测试、测试覆盖率 ≥ 80% |

### 3.3 审批规则

| 变更规模 | 审批要求 |
|----------|----------|
| 小改动（≤ 100 行） | 1 人审批 |
| 中等改动（100-500 行） | 2 人审批 |
| 大改动（> 500 行） | 3 人审批 + 安全审查 |
| 核心模块改动 | 核心模块负责人必须审批 |

---

## 四、版本标签规范

### 4.1 标签格式

```
v{主版本}.{次版本}.{补丁版本}
```

### 4.2 版本含义

| 版本段 | 触发条件 |
|--------|----------|
| 主版本 | API 不兼容变更 |
| 次版本 | 新功能，向后兼容 |
| 补丁版本 | 修复，向后兼容 |

### 4.3 标签标注时机

- 合并到 `main` 后打标签
- 使用 `git tag -a v1.0.0 -m "Release v1.0.0"`

### 4.4 关联 API 版本

- API 版本 `/api/v1/` 与代码版本同步
- API 版本变更时主版本递增

---

## 五、关联文档

- [接口文档规范](./02-接口文档规范.md)
