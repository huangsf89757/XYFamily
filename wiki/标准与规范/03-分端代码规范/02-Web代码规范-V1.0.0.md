# Web 代码规范

> XYFamily Web 端开发规范：React + TypeScript，PC 管理后台。

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
| 关联标签 | 核心文档、规范标准 |
| 关联目录 | 04-开发规范与编码手册/04.02-分端代码规范 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建 | ClaudeCode |

---

## 一、技术栈

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

```typescript
// 45 个权限点枚举（与后端 role_permissions 一致）
enum Permission {
  // 认证
  AUTH_REGISTER = 'auth.register',
  AUTH_LOGIN = 'auth.login',
  AUTH_LOGOUT = 'auth.logout',
  AUTH_REFRESH = 'auth.refresh',
  AUTH_RESET_PASSWORD = 'auth.reset_password',
  // 账号
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

// 使用
function MemberList() {
  const canInvite = useHasPermission(Permission.ORG_MEMBER_INVITE);
  const canRemove = useHasPermission(Permission.ORG_MEMBER_REMOVE);

  return (
    <div>
      {canInvite && <Button onClick={handleInvite}>邀请成员</Button>}
      {canRemove && <Button onClick={handleRemove}>移除成员</Button>}
      <Table dataSource={members} />
    </div>
  );
}
```

---

## 四、布局与权限路由

```typescript
function ProtectedRoute({ children, permissions }: { children: ReactNode; permissions?: Permission[] }) {
  const hasAll = !permissions || permissions.every(p => useHasPermission(p));
  if (!hasAll) return <AccessDenied />;
  return children;
}

// 路由配置
const routes = [
  { path: '/org/:id', element: <ProtectedRoute permissions={[Permission.ORG_READ]}><OrgPage /></ProtectedRoute> },
  { path: '/audit/global', element: <ProtectedRoute permissions={[Permission.ADMIN_AUDIT_GLOBAL]}><GlobalAudit /></ProtectedRoute> },
];
```

---

## 五、多租户上下文

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

## 六、审计日志表格

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

## 七、错误处理

```typescript
// Axios 拦截器
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

## 八、响应式布局

```css
/* 最小宽度 1200px，侧边栏 200px */
.layout {
  min-width: 1200px;
}

.sidebar {
  width: 200px;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
}

.content {
  margin-left: 200px;
  padding: 24px;
}

/* 表格列宽自适应 */
.table-column {
  min-width: 120px;
  max-width: 300px;
}
```

---

## 九、关联文档

- [前端通用代码规范](./04.02-分端代码规范/前端通用代码规范-V1.0.0.md)
- [接口总览](../05-接口与模块落地文档/05.01-接口总览/接口总览-V1.0.0.md)
- [权限管理 PRD](../02-需求与产品设计/02.02-产品PRD/07-权限管理模块/权限管理模块-V1.0.0.md)
