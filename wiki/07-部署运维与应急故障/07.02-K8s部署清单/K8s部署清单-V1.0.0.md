# K8s 部署清单

> XYFamily Kubernetes 部署清单：Deployment、Service、Ingress、ConfigMap、Secret。

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
| 关联目录 | 07-部署运维与应急故障/07.02-K8s部署清单 |

## 变更记录

| 版本 | 日期 | 变更内容 | 变更人 |
|------|------|----------|--------|
| V1.0.0 | 2026-07-12 | 初始创建 | ClaudeCode |

---

## 一、目录结构

```
deploy/
  k8s/
    namespace.yaml          # Namespace
    secret.yaml             # Secret（数据库密码、JWT 密钥等）
    configmap.yaml          # ConfigMap（应用配置）
    deployment.yaml         # Deployment（3 副本）
    service.yaml            # Service（ClusterIP）
    ingress.yaml            # Ingress（HTTPS）
    hpa.yaml                # HPA（水平自动扩缩）
    pvc-pg.yaml             # PVC（PostgreSQL 持久化）
    statefulset-pg.yaml     # StatefulSet（PG 主从）
    pvc-redis.yaml          # PVC（Redis 持久化）
    sentinel.yaml           # Sentinel（Pod）
```

---

## 二、Namespace

```yaml
# namespace.yaml

apiVersion: v1
kind: Namespace
metadata:
  name: xyfamily
  labels:
    name: xyfamily
```

---

## 三、Secret（密钥管理）

```yaml
# secret.yaml
# ⚠️ 生产环境请使用外部 Secret Store（如 HashiCorp Vault / AWS Secrets Manager）

apiVersion: v1
kind: Secret
metadata:
  name: xyfamily-secrets
  namespace: xyfamily
type: Opaque
data:
  # base64 编码（生产环境从 Vault 注入）
  db-password: <base64:XY_DB_PASSWORD>
  redis-password: <base64:XY_REDIS_PASSWORD>
  jwt-secret: <base64:XY_JWT_SECRET>
```

---

## 四、ConfigMap

```yaml
# configmap.yaml

apiVersion: v1
kind: ConfigMap
metadata:
  name: xyfamily-config
  namespace: xyfamily
data:
  XY_ENV: "production"
  XY_DB_HOST: "pg-master"
  XY_DB_PORT: "5432"
  XY_DB_NAME: "xyfamily"
  XY_REDIS_MASTER: "xyfamily"
  XY_REDIS_SENTINEL_ADDR: "redis-sentinel-0:26379,redis-sentinel-1:26379,redis-sentinel-2:26379"
  XY_LOG_LEVEL: "info"
```

---

## 五、Deployment

```yaml
# deployment.yaml

apiVersion: apps/v1
kind: Deployment
metadata:
  name: xyfamily
  namespace: xyfamily
  labels:
    app: xyfamily
spec:
  replicas: 3
  selector:
    matchLabels:
      app: xyfamily
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    metadata:
      labels:
        app: xyfamily
    spec:
      containers:
        - name: xyfamily
          image: xyfamily:latest
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
              name: http
          envFrom:
            - configMapRef:
                name: xyfamily-config
          env:
            - name: XY_DB_USER
              value: "xyfamily"
            - name: XY_DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: xyfamily-secrets
                  key: db-password
            - name: XY_REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: xyfamily-secrets
                  key: redis-password
            - name: XY_JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: xyfamily-secrets
                  key: jwt-secret
          resources:
            requests:
              cpu: 250m
              memory: 256Mi
            limits:
              cpu: 1000m
              memory: 512Mi
          readinessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
            timeoutSeconds: 5
            failureThreshold: 3
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 30
            timeoutSeconds: 5
            failureThreshold: 3
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
```

**关键设计**：
- `maxUnavailable: 0` 保证滚动更新期间始终有 3 个实例运行
- readinessProbe 确保流量只打到就绪的 Pod
- livenessProbe 检测存活，故障时自动重启
- `terminationGracePeriodSeconds: 30` 优雅关闭已有连接

---

## 六、Service

```yaml
# service.yaml

apiVersion: v1
kind: Service
metadata:
  name: xyfamily-service
  namespace: xyfamily
spec:
  type: ClusterIP
  selector:
    app: xyfamily
  ports:
    - name: http
      port: 80
      targetPort: 8080
      protocol: TCP
```

---

## 七、Ingress（HTTPS）

```yaml
# ingress.yaml

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: xyfamily-ingress
  namespace: xyfamily
  annotations:
    kubernetes.io/ingress.class: "nginx"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "60"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "60"
spec:
  tls:
    - hosts:
        - api.xyfamily.example.com
      secretName: xyfamily-tls
  rules:
    - host: api.xyfamily.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: xyfamily-service
                port:
                  number: 80
```

**证书管理**：使用 cert-manager + Let's Encrypt 自动签发和续期。

---

## 八、HPA（水平自动扩缩）

```yaml
# hpa.yaml

apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: xyfamily-hpa
  namespace: xyfamily
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: xyfamily
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300  # 缩容前稳定 5 分钟
    scaleUp:
      stabilizationWindowSeconds: 60    # 扩容前稳定 1 分钟
```

---

## 九、PostgreSQL StatefulSet

```yaml
# statefulset-pg.yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: pg
  namespace: xyfamily
spec:
  serviceName: pg
  replicas: 2
  selector:
    matchLabels:
      app: pg
  template:
    metadata:
      labels:
        app: pg
    spec:
      containers:
        - name: postgres
          image: postgres:17
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: "xyfamily"
            - name: POSTGRES_USER
              value: "xyfamily"
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: xyfamily-secrets
                  key: db-password
            - name: PGDATA
              value: "/var/lib/postgresql/data/pgdata"
          args:
            - "postgres"
            - "-c"
            - "wal_level=replica"
            - "-c"
            - "max_wal_senders=5"
          volumeMounts:
            - name: pg-data
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              cpu: 500m
              memory: 512Mi
            limits:
              cpu: "2"
              memory: 2Gi
  volumeClaimTemplates:
    - metadata:
        name: pg-data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 50Gi
```

---

## 十、Redis Sentinel StatefulSet

```yaml
# sentinel.yaml

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: redis-sentinel
  namespace: xyfamily
spec:
  serviceName: redis-sentinel
  replicas: 3
  selector:
    matchLabels:
      app: redis-sentinel
  template:
    metadata:
      labels:
        app: redis-sentinel
    spec:
      containers:
        - name: sentinel
          image: redis:7-alpine
          ports:
            - containerPort: 26379
          command: ["redis-sentinel"]
          args: ["/usr/local/etc/redis/sentinel.conf"]
          env:
            - name: XY_REDIS_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: xyfamily-secrets
                  key: redis-password
```

---

## 十一、部署流程

```bash
# 1. 创建 Namespace
kubectl apply -f k8s/namespace.yaml

# 2. 创建 Secret（生产环境使用外部 Secret Store）
kubectl apply -f k8s/secret.yaml

# 3. 创建 ConfigMap
kubectl apply -f k8s/configmap.yaml

# 4. 部署 PostgreSQL（StatefulSet）
kubectl apply -f k8s/pvc-pg.yaml
kubectl apply -f k8s/statefulset-pg.yaml

# 5. 部署 Redis Sentinel
kubectl apply -f k8s/pvc-redis.yaml
kubectl apply -f k8s/sentinel.yaml

# 6. 部署应用
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
kubectl apply -f k8s/hpa.yaml

# 7. 查看状态
kubectl get pods -n xyfamily
kubectl get svc -n xyfamily
kubectl get ingress -n xyfamily
```

---

## 十二、关联文档

- [Docker Compose 生产配置](./07.01-Docker Compose生产配置/Docker%20Compose生产配置-V1.0.0.md)
- [容灾多活架构](../03-技术架构与方案设计/03.07-容灾多活架构/容灾多活架构-V1.0.0.md)
