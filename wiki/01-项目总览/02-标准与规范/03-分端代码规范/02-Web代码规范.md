# Web代码规范

> XYFamily Web 管理后台（React + TypeScript）**项目特定约定与选型**。通用 React/TypeScript 编码规范请直接参考官方与业界标准，本文不再重复。

---

## 文档信息

| 项目 | 内容 |
|------|------|
| 文档密级 | 内部 |
| 文档版本 | V1.1.0 |
| 编写人 | ClaudeCode |
| 审核人 | - |
| 生效时间 | 2026-07-12 |
| 废弃时间 | - |
| 关联标签 | 核心文档、规范标准 |
| 关联目录 | 01-项目总览/标准与规范/03-分端代码规范 |

## 变更记录


| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.1.0 | 2026-07-19 | 文档新编 | ClaudeCode |

---

## 参考规范（官方/业界标准）

- **[Airbnb React/JSX Style Guide](https://github.com/airbnb/javascript/tree/master/react)** — React 编码规范
- **[React 官方文档](https://react.dev/)** — Hooks、组件设计等最佳实践
- **[TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/)** — TS 语言规范
- **[Ant Design 文档](https://ant.design/)** — UI 组件使用指南
- **[Zustand 文档](https://docs.pmnd.rs/zustand)** — 状态管理
- **[React Router 文档](https://reactrouter.com/)** — 路由使用

---

## 一、技术栈选型

| 类别 | 选型 |
|------|------|
| 语言 | TypeScript 5.x |
| 框架 | React 18+ |
| UI 组件 | Ant Design 5.x |
| 状态管理 | Zustand |
| 路由 | React Router v6 |
| 网络 | Axios |
| 表单 | Formik / React Hook Form |
| 图表 | ECharts |
| 构建 | Vite |
| 适配 | 响应式布局，最小宽度 1200px |

---

## 二、项目目录结构

```
apps/web/
├── src/
│   ├── api/                      # API 封装
│   │   ├── request.ts            # Axios 封装 + 拦截器
│   │   ├── auth.ts               # 认证相关
│   │   ├── org.ts                # 组织相关
│   │   ├── team.ts               # 团队相关
│   │   ├── audit.ts              # 审计日志
│   │   └── admin.ts              # SuperAdmin 管理
│   ├── store/
│   │   ├── authStore.ts          # 认证状态
│   │   ├── orgStore.ts           # 组织上下文
│   │   ├── permissionStore.ts    # 权限状态（45 个权限点）
│   └── pages/
│   │   ├── login/                # 登录页
│   │   ├── dashboard/            # 仪表盘
│   │   ├── org/                  # 组织管理
│   │   ├── team/                 # 团队管理
│   │   ├── group/                # 小组管理
│   │   ├── audit/                # 审计日志
│   │   └── admin/                # 系统管理（SuperAdmin）
│   ├── components/               # 公共组件
│   ├── layouts/                  # 布局（侧边栏 + 顶栏）
│   ├── utils/
│   │   ├── permission.ts         # 权限校验（前端辅助）
│   │   └── download.ts           # 文件下载
│   └── App.tsx
```

---

## 三、权限管理（前端辅助）

### 3.1 权限点定义

本项目共 45 个权限点，与后端 `role_permissions` 一致，部分示例：

```typescript
enum Permission {
  AUTH_REGISTER = 'auth.register',
  AUTH_LOGIN = 'auth.login',
  AUTH_LOGOUT = 'auth.logout',
  AUTH_REFRESH = 'auth.refresh',
  AUTH_RESET_PASSWORD = 'auth.reset_password',
  ACCOUNT_PROFILE_READ = 'account.profile.read',
  ACCOUNT_PROFILE_UPDATE = 'account.profile.update',
  ACCOUNT_PASSWORD_UPDATE = 'account.password.update',
  ACCOUNT_DEACTIVATE = 'account.deactivate',
  ACCOUNT_UNDEACTIVATE = 'account.undeactivate',
  // ...（完整 45 个）
}
```

### 3.2 权限校验 Hook

```typescript
function useHasPermission(code: Permission): boolean {
  const { permissions, currentRole } = useAuthStore();
  return permissions.has(code);
}
```

### 3.3 权限路由

```typescript
function ProtectedRoute({ children, permissions }: { children: ReactNode; permissions?: Permission[] }) {
  const hasAll = !permissions || permissions.every(p => useHasPermission(p));
  if (!hasAll) return <AccessDenied />;
  return children;
}
```

---

## 四、多租户上下文

```typescript
// 组织切换
function OrgSwitcher() {
  const { user, currentOrgId, switchOrg } = useAuthStore();
  const orgs = useUserOrgs();

  return (
    <Select value={currentOrgId} onChange={switchOrg}>
      {orgs.map(org => (
        <Option key={org.id} value={org.id}>{org.name}</Option>
      ))}
    </Select>
  );
}
```

---

## 五、审计日志

```typescript
function AuditLogTable() {
  const [data, setData] = useState<OperationAudit[]>([]);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 20 });

  const fetchLogs = async () => {
    const response = await api.getOperationLogs({
      page: pagination.current,
      page_size: pagination.pageSize,
      org_id: getCurrentOrgId(),
    });
    setData(response.data.items);
  };

  return (
    <Table
      dataSource={data}
      columns={AUDIT_COLUMNS}
      pagination={{ total: response.data.total, ...pagination }}
      onChange={setPagination}
      rowKey="id"
    />
  );
}
```

---

## 六、错误处理（项目错误码映射）

```typescript
// Axios 响应拦截器
instance.interceptors.response.use(
  (response) => {
    const { code } = response.data;
    if (code !== 0) {
      message.error(getErrorMessage(code, response.data.message));
      if (code >= 101001 && code <= 101009) {
        handleTokenExpired();
      }
      return Promise.reject(response.data);
    }
    return response.data;
  },
  (error) => {
    if (error.response?.status === 401) {
      handleTokenExpired();
    } else if (error.response?.status === 500) {
      message.error('服务异常，请稍后重试');
    }
    return Promise.reject(error);
  }
);

// 本项目错误码映射
function getErrorMessage(code: number, fallback: string): string {
  const map: Record<number, string> = {
    101001: '登录已过期，请重新登录',
    101002: '登录已过期，请重新登录',
    101003: '登录已过期，请重新登录',
    104290: '操作过于频繁，请稍后再试',
    603001: '当前账号无权限执行此操作',
    603002: '当前账号无权限执行此操作',
    805000: '服务异常，请稍后重试',
  };
  return map[code] || fallback;
}
```

---

## 七、关联文档

- [前端通用代码规范](./07-前端通用代码规范.md)
- [接口文档](../../../04-接口文档/接口文档.md)
- [权限管理模块](../../../02-需求与产品设计/01-产品PRD/01-多租户底座/06-权限管理模块/权限管理模块.md)

## 关联文档


> 以下为知识图谱自动推荐的交叉引用，建议人工审阅确认后保留。

- [架构评审记录](../../../03-架构与方案设计/审核记录/架构评审记录.md) — 共享术语：多租户、审计、权限（置信度 0.75）
- [ADR架构决策记录](../../../03-架构与方案设计/01-基座/02-ADR架构决策记录.md) — 共享术语：多租户、审计、权限（置信度 0.75）
- [多租户底座](../../../02-需求与产品设计/01-产品PRD/01-多租户底座/多租户底座.md) — 共享术语：多租户、审计、权限（置信度 0.75）
- [超级管理员模块](../../../02-需求与产品设计/01-产品PRD/01-多租户底座/07-超级管理员模块/超级管理员模块.md) — 共享术语：多租户、审计、权限（置信度 0.75）
- [PRD审核记录](../../../02-需求与产品设计/01-产品PRD/审核记录/PRD审核记录.md) — 共享术语：多租户、审计、权限（置信度 0.75）
