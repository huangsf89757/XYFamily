# Docker Compose 生产配置

> XYFamily 生产级 Docker Compose 配置，含 PG 主从、Redis Sentinel、应用多实例。

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
| 关联标签 | 核心文档、系统基础 |
| 关联目录 | 07-部署运维与应急故障/07.01-Docker Compose生产配置 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建 | ClaudeCode |

---

## 一、目录结构

```
deploy/
  docker-compose.prod.yml   # 生产环境编排
  docker-compose.staging.yml # 预发布环境编排
  docker-compose.dev.yml    # 本地开发环境编排（已有，见研发工具规范）
  .env.template             # 环境变量模板
  pg/
    init.sql                # PG 初始化脚本
  redis/
    sentinel/
      sentinel1.conf
      sentinel2.conf
      sentinel3.conf
  nginx/
    nginx.conf              # 反向代理配置
```

---

## 二、生产环境 docker-compose.prod.yml

```yaml
version: "3.8"

services:
  # ─── 应用实例 ───
  xyfamily-1:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: xyfamily-1
    restart: unless-stopped
    ports:
      - "8001:8080"
    environment:
      - XY_ENV=production
      - XY_DB_HOST=pg-master
      - XY_DB_PORT=5432
      - XY_DB_NAME=xyfamily
      - XY_DB_USER=${XY_DB_USER}
      - XY_DB_PASSWORD=${XY_DB_PASSWORD}
      - XY_REDIS_MASTER=xyfamily
      - XY_REDIS_SENTINEL_ADDR=${XY_REDIS_SENTINEL_ADDRS}
      - XY_REDIS_PASSWORD=${XY_REDIS_PASSWORD}
      - XY_JWT_SECRET=${XY_JWT_SECRET}
      - XY_LOG_LEVEL=info
    depends_on:
      pg-master:
        condition: service_healthy
      redis-sentinel-1:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - backend

  xyfamily-2:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: xyfamily-2
    restart: unless-stopped
    ports:
      - "8002:8080"
    environment:
      - XY_ENV=production
      - XY_DB_HOST=pg-master
      - XY_DB_PORT=5432
      - XY_DB_NAME=xyfamily
      - XY_DB_USER=${XY_DB_USER}
      - XY_DB_PASSWORD=${XY_DB_PASSWORD}
      - XY_REDIS_MASTER=xyfamily
      - XY_REDIS_SENTINEL_ADDR=${XY_REDIS_SENTINEL_ADDRS}
      - XY_REDIS_PASSWORD=${XY_REDIS_PASSWORD}
      - XY_JWT_SECRET=${XY_JWT_SECRET}
      - XY_LOG_LEVEL=info
    depends_on:
      pg-master:
        condition: service_healthy
      redis-sentinel-1:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - backend

  xyfamily-3:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: xyfamily-3
    restart: unless-stopped
    ports:
      - "8003:8080"
    environment:
      - XY_ENV=production
      - XY_DB_HOST=pg-master
      - XY_DB_PORT=5432
      - XY_DB_NAME=xyfamily
      - XY_DB_USER=${XY_DB_USER}
      - XY_DB_PASSWORD=${XY_DB_PASSWORD}
      - XY_REDIS_MASTER=xyfamily
      - XY_REDIS_SENTINEL_ADDR=${XY_REDIS_SENTINEL_ADDRS}
      - XY_REDIS_PASSWORD=${XY_REDIS_PASSWORD}
      - XY_JWT_SECRET=${XY_JWT_SECRET}
      - XY_LOG_LEVEL=info
    depends_on:
      pg-master:
        condition: service_healthy
      redis-sentinel-1:
        condition: service_started
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:8080/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - backend

  # ─── Nginx 反向代理 ───
  nginx:
    image: nginx:alpine
    container_name: nginx-lb
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - xyfamily-1
      - xyfamily-2
      - xyfamily-3
    networks:
      - backend

  # ─── PostgreSQL 主库 ───
  pg-master:
    image: postgres:17
    container_name: pg-master
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_DB=xyfamily
      - POSTGRES_USER=${XY_DB_USER}
      - POSTGRES_PASSWORD=${XY_DB_PASSWORD}
      - POSTGRES_INITDB_ARGS=--data-checksums
    command:
      - "postgres"
      - "-c"
      - "wal_level=replica"
      - "-c"
      - "max_wal_senders=5"
      - "-c"
      - "wal_keep_size=1024"
    volumes:
      - pg-master-data:/var/lib/postgresql/data
      - ./pg/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${XY_DB_USER} -d xyfamily"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - backend

  # ─── PostgreSQL 从库 ───
  pg-slave-1:
    image: postgres:17
    container_name: pg-slave-1
    restart: unless-stopped
    ports:
      - "5433:5432"
    environment:
      - POSTGRES_USER=${XY_DB_USER}
      - POSTGRES_PASSWORD=${XY_DB_PASSWORD}
    depends_on:
      pg-master:
        condition: service_healthy
    volumes:
      - pg-slave-1-data:/var/lib/postgresql/data
    networks:
      - backend

  # ─── Redis Master ───
  redis-master:
    image: redis:7-alpine
    container_name: redis-master
    restart: unless-stopped
    command: redis-server --requirepass ${XY_REDIS_PASSWORD}
    volumes:
      - redis-master-data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "-a", "${XY_REDIS_PASSWORD}", "ping"]
      interval: 10s
      timeout: 3s
      retries: 3
    networks:
      - backend

  # ─── Redis Replica 1 ───
  redis-replica-1:
    image: redis:7-alpine
    container_name: redis-replica-1
    restart: unless-stopped
    command: redis-server --replicaof redis-master 6379 --requirepass ${XY_REDIS_PASSWORD}
    depends_on:
      redis-master:
        condition: service_healthy
    volumes:
      - redis-replica-1-data:/data
    networks:
      - backend

  # ─── Redis Replica 2 ───
  redis-replica-2:
    image: redis:7-alpine
    container_name: redis-replica-2
    restart: unless-stopped
    command: redis-server --replicaof redis-master 6379 --requirepass ${XY_REDIS_PASSWORD}
    depends_on:
      redis-master:
        condition: service_healthy
    volumes:
      - redis-replica-2-data:/data
    networks:
      - backend

  # ─── Redis Sentinel 1 ───
  redis-sentinel-1:
    image: redis:7-alpine
    container_name: redis-sentinel-1
    restart: unless-stopped
    volumes:
      - ./redis/sentinel/sentinel1.conf:/usr/local/etc/redis/sentinel.conf:ro
    command: redis-sentinel /usr/local/etc/redis/sentinel.conf
    depends_on:
      redis-master:
        condition: service_healthy
    networks:
      - backend

  # ─── Redis Sentinel 2 ───
  redis-sentinel-2:
    image: redis:7-alpine
    container_name: redis-sentinel-2
    restart: unless-stopped
    volumes:
      - ./redis/sentinel/sentinel2.conf:/usr/local/etc/redis/sentinel.conf:ro
    command: redis-sentinel /usr/local/etc/redis/sentinel.conf
    depends_on:
      redis-master:
        condition: service_healthy
    networks:
      - backend

  # ─── Redis Sentinel 3 ───
  redis-sentinel-3:
    image: redis:7-alpine
    container_name: redis-sentinel-3
    restart: unless-stopped
    volumes:
      - ./redis/sentinel/sentinel3.conf:/usr/local/etc/redis/sentinel.conf:ro
    command: redis-sentinel /usr/local/etc/redis/sentinel.conf
    depends_on:
      redis-master:
        condition: service_healthy
    networks:
      - backend

volumes:
  pg-master-data:
  pg-slave-1-data:
  redis-master-data:
  redis-replica-1-data:
  redis-replica-2-data:

networks:
  backend:
    driver: bridge
```

---

## 三、环境变量模板（.env.template）

```bash
# 数据库
XY_DB_USER=xyfamily
XY_DB_PASSWORD=CHANGE_ME_STRONG_PASSWORD

# Redis
XY_REDIS_PASSWORD=CHANGE_ME_STRONG_PASSWORD
XY_REDIS_SENTINEL_ADDRS=sentinel1:26379,sentinel2:26379,sentinel3:26379

# JWT
XY_JWT_SECRET=CHANGE_ME_64_CHAR_RANDOM_SECRET

# 应用
XY_ENV=production
XY_LOG_LEVEL=info
```

---

## 四、Nginx 反向代理配置

```nginx
# nginx.conf

worker_processes auto;

events {
    worker_connections 1024;
}

http {
    upstream xyfamily_backend {
        server xyfamily-1:8080;
        server xyfamily-2:8080;
        server xyfamily-3:8080;
    }

    server {
        listen 443 ssl;
        server_name api.xyfamily.example.com;

        ssl_certificate     /etc/nginx/ssl/server.crt;
        ssl_certificate_key /etc/nginx/ssl/server.key;
        ssl_protocols       TLSv1.2 TLSv1.3;
        ssl_ciphers         HIGH:!aNULL:!MD5;

        location / {
            proxy_pass http://xyfamily_backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # 健康检查绕过认证
            location /healthz {
                access_log off;
            }
        }
    }

    server {
        listen 80;
        server_name api.xyfamily.example.com;
        return 301 https://$server_name$request_uri;
    }
}
```

---

## 五、Redis Sentinel 配置

```conf
# sentinel1.conf

port 26379
sentinel monitor xyfamily redis-master 6379 2
sentinel auth-pass xyfamily ${XY_REDIS_PASSWORD}
sentinel down-after-milliseconds xyfamily 5000
sentinel failover-timeout xyfamily 30000
sentinel parallel-syncs xyfamily 1
```

sentinel2/sentinel3 配置相同。

---

## 六、Dockerfile

```dockerfile
# Dockerfile

FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o xyfamily ./cmd/xyfamily

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /build/xyfamily /usr/local/bin/xyfamily
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD wget -q --spider http://localhost:8080/healthz || exit 1
CMD ["xyfamily"]
```

---

## 七、启动/停止

```bash
# 生产环境启动
cp .env.template .env
# 编辑 .env 填入真实密钥
docker-compose -f docker-compose.prod.yml up -d

# 查看状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f xyfamily-1

# 停止
docker-compose -f docker-compose.prod.yml down
```

---

## 八、关联文档

- [容灾多活架构](../03-技术架构与方案设计/03.07-容灾多活架构/容灾多活架构-V1.0.0.md)
- [研发工具规范](../04-开发规范与编码手册/04.04-研发工具规范/研发工具规范-V1.0.0.md)
