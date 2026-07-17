# XYFamily

XYFamily 是一个多租户账号权限底座（SaaS），为上层业务提供统一的账号认证、RBAC 权限管理、多租户隔离能力。

## 技术栈

- **后端**: Go 1.22+ / Gin / PostgreSQL 17 / Redis 7
- **前端**: React + TypeScript（Web/H5）、原生（iOS/Android/鸿蒙）
- **部署**: Docker / K8s

## 项目结构

```
code/backend/          # Go 后端服务
code/web/              # Web 管理后台（React + Ant Design）
code/h5/               # H5 移动浏览器端
code/miniprogram/      # 微信小程序
code/ios/              # iOS 原生 App
code/android/          # Android 原生 App
code/harmony/          # 鸿蒙 HarmonyOS App
```

## 快速启动

### 1. 启动开发环境（PostgreSQL + Redis）

```bash
cd code/backend
docker compose -f docker-compose.dev.yml up -d
```

### 2. 配置

复制配置文件并设置敏感参数：

```bash
cp configs/config.yaml configs/config.local.yaml
export XY_DB_PASSWORD=xyfamily_dev
export XY_JWT_SECRET=your-jwt-secret
export XY_PII_ENCRYPTION_KEY=your-encryption-key
export XY_PII_INDEX_KEY=your-index-key
```

### 3. 运行数据库迁移

```bash
make migrate
```

### 4. 启动服务

```bash
make run
```

服务启动后访问 `http://localhost:8080/api/v1/healthz` 检查健康状态。

### 5. 构建与测试

```bash
make build    # 编译
make test     # 运行测试
make lint     # 代码检查
```

## API 接口

基础路径: `/api/v1`

| 接口 | 方法 | 说明 | 认证 |
|------|------|------|------|
| /healthz | GET | 健康检查 | 否 |
| /auth/verification-codes | POST | 发送验证码 | 否 |
| /auth/register | POST | 注册 | 否 |
| /auth/login | POST | 登录 | 否 |
| /auth/refresh | POST | 刷新 Token | 否 |
| /auth/logout | POST | 登出 | 是 |
| /auth/reset-password | POST | 重置密码 | 否 |

## 文档

详见 `wiki/` 目录下的结构化知识库。

## License

MIT
