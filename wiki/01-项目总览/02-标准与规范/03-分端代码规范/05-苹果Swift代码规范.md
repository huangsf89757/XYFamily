# iOS（Swift）代码规范

> XYFamily iOS 端（Swift + SwiftUI）**项目特定约定与选型**。通用 Swift/SwiftUI 编码规范请直接参考官方与业界标准，本文不再重复。

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
| V1.0.0 | 2026-07-12 | 初始创建 | ClaudeCode |
| V1.1.0 | 2026-07-16 | 移除通用规范，仅保留项目特有约定 | ClaudeCode |

---

## 参考规范（官方/业界标准）

- **[Swift API Design Guidelines](https://www.swift.org/documentation/api-design-guidelines/)** — Apple 官方 API 设计指南
- **[Swift Style Guide (Google)](https://google.github.io/swift/)** — Google Swift 编码规范
- **[SwiftUI Tutorials](https://developer.apple.com/tutorials/swiftui)** — Apple SwiftUI 官方教程
- **[Swift Concurrency](https://docs.swift.org/swift-book/documentation/the-swift-programming-language/concurrency/)** — Swift 并发编程指南

---

## 一、技术栈选型

| 类别 | 选型 |
|------|------|
| 语言 | Swift 5.9+ |
| UI 框架 | SwiftUI |
| 网络 | URLSession + 自定义封装 |
| 数据持久化 | Core Data（主）+ UserDefaults（轻量） |
| 安全存储 | Keychain（通过 KeychainAccess 封装） |
| 异步 | Swift Concurrency（async/await） |
| 状态管理 | @State / @ObservableObject / @EnvironmentObject |
| 图片 | SwiftUI Image + URLSession 加载 |
| 最低版本 | iOS 16.0 |
| 目标版本 | iOS 18.0 |

---

## 二、项目目录结构

```
Xyfamily/
├── Resources/
│   ├── Images.xcassets/
│   └── Info.plist
├── Common/
│   ├── Network/
│   │   ├── NetworkClient.swift           # URLSession 封装
│   │   ├── AuthInterceptor.swift          # Token 注入
│   │   ├── TokenRefreshManager.swift      # Token 自动刷新
│   │   └── APIError.swift                 # 统一错误类型
│   ├── DI/
│   └── Utils/
├── Data/
│   ├── Models/                            # 数据模型
│   └── Repositories/
├── Domain/
│   └── UseCases/
├── Presentation/
│   ├── Auth/                              # 认证模块
│   ├── Account/                           # 账号管理
│   ├── Organization/                      # 组织管理
│   ├── Team/                              # 团队管理
│   ├── Group/                             # 小组管理
│   └── Common/                            # 公共视图组件
└── Services/
    ├── SecureStorage.swift                # Keychain 封装
    └── PermissionService.swift            # 权限校验
```

---

## 三、API 调用（项目 Header 约定）

```swift
actor NetworkClient {
    func request<T: Decodable>(...) async throws -> T {
        var request = URLRequest(url: baseURL.appendingPathComponent(endpoint.path))

        // Token 注入
        if let token = tokenProvider.accessToken() {
            request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        }

        // 组织上下文
        if let orgId = tokenProvider.currentOrgId() {
            request.setValue(orgId, forHTTPHeaderField: "X-Organization-ID")
        }

        let (data, response) = try await session.data(for: request)
        return try handleResponse(data: data, response: response, type: type)
    }
}
```

---

## 四、安全存储（项目约定）

```swift
// Access Token：不持久化（内存）
actor TokenManager {
    private(set) var accessToken: String? = nil
    func setAccessToken(_ token: String) { accessToken = token }
    func getAccessToken() -> String? { accessToken }
    func clearAll() { accessToken = nil }
}

// Refresh Token：Keychain
struct KeychainStorage {
    func saveRefreshToken(_ token: String) throws {
        try KeychainWrapper.standard.set(token, forKey: "xyfamily_refresh_token", in: service)
    }
    func getRefreshToken() throws -> String? {
        return try KeychainWrapper.standard.get("xyfamily_refresh_token", in: service)
    }
}
```

---

## 五、Token 自动刷新

```swift
actor TokenRefreshManager {
    private var isRefreshing = false

    func refreshToken() async throws -> String {
        guard !isRefreshing else {
            return try await waitForRefresh()
        }
        isRefreshing = true
        defer { isRefreshing = false }

        let refreshToken = try KeychainStorage().getRefreshToken()
            ?? throw APIError.authExpired
        let newToken = try await NetworkClient().requestRefresh(refreshToken)
        TokenManager().setAccessToken(newToken)
        return newToken
    }
}
```

---

## 六、权限管理（项目角色）

```swift
enum UserRole: String {
    case superAdmin = "SuperAdmin"
    case orgCoreAdmin = "organization_core_admin"
    case teamCoreAdmin = "team_core_admin"
    case groupCoreAdmin = "group_core_admin"
    case regularMember = "RegularMember"
    case publicRole = "Public"
}

extension UserRole {
    func hasPermission(_ code: String) -> Bool {
        let permMap: [String: Set<String>] = [
            "SuperAdmin": Set(allPermissionCodes),
            "organization_core_admin": orgCoreAdminPerms,
            // ...
        ]
        return permMap[rawValue]?.contains(code) ?? false
    }
}
```

---

## 七、错误处理（项目错误码映射）

```swift
enum APIError: LocalizedError {
    case authExpired
    case rateLimited
    case forbidden
    case serverError
    case networkError(Error)
    case unknown(Int, String)

    var errorDescription: String? {
        switch self {
        case .authExpired: return "登录已过期，请重新登录"
        case .rateLimited: return "操作过于频繁，请稍后再试"
        case .forbidden: return "当前账号无权限执行此操作"
        case .serverError: return "服务异常，请稍后重试"
        case .networkError: return "网络连接异常"
        case .unknown(_, let msg): return msg
        }
    }
}
```

---

## 八、关联文档

- [前端通用代码规范](./07-前端通用代码规范.md)
- [接口总览](../../../04-接口文档/接口文档.md)

## 关联文档


> 以下为知识图谱自动推荐的交叉引用，建议人工审阅确认后保留。

- [04-安卓Kotlin代码规范](./04-安卓Kotlin代码规范.md) — 共享术语：token、安全、权限、角色（置信度 0.75）
