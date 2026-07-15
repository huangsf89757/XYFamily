# iOS（Swift）代码规范

> XYFamily iOS 端开发规范：Swift + SwiftUI。

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
│   │   ├── Account.swift
│   │   ├── Organization.swift
│   │   └── Permission.swift
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

## 三、命名规范

### 3.1 基础命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 类型（类/结构体/枚举/协议） | PascalCase | `LoginView`, `UserRepository`, `AuthState` |
| 函数/方法 | camelCase | `fetchUserList`, `handleLogin` |
| 常量 | UPPER_SNAKE_CASE（全局） | `let API_BASE_URL = ""` |
| 变量 | camelCase | `var currentUser`, `var isLoggedIn` |
| 布尔值 | is/has/can/should 前缀 | `isAuthenticated`, `hasPermission` |
| 资源文件 | PascalCase | `UserAvatar.swift`, `OrganizationIcon.swift` |

### 3.2 SwiftUI 组件命名

```swift
// View 结构体：PascalCase
struct UserProfileView: View {
    @State private var isLoading = false
    @State private var errorMessage: String?

    var body: some View {
        // ...
    }
}

// 状态管理：@Observable 或 @StateObject
@Observable
class OrgContext {
    var currentOrg: String? = nil

    func switchOrg(to orgId: String) {
        currentOrg = orgId
    }
}
```

---

## 四、API 调用封装

### 4.1 URLSession 封装

```swift
actor NetworkClient {
    private let session: URLSession
    private let tokenProvider: TokenProvider
    private let baseURL: URL

    func request<T: Decodable>(
        endpoint: APIEndpoint,
        type: T.Type,
        method: HTTPMethod = .get,
        body: Encodable? = nil
    ) async throws -> T {
        var request = URLRequest(url: baseURL.appendingPathComponent(endpoint.path))
        request.httpMethod = method.rawValue
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")

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

### 4.2 统一错误类型

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

### 4.3 响应处理

```swift
func handleResponse<T: Decodable>(
    data: Data,
    response: HTTPURLResponse,
    type: T.Type
) throws -> T {
    let apiResponse = try JSONDecoder().decode(APIResponse<T>.self, from: data)
    guard apiResponse.code == 0, let data = apiResponse.data else {
        throw APIError.unknown(apiResponse.code, apiResponse.message)
    }
    return data
}
```

---

## 五、安全存储

### 5.1 Token 存储

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
    private let service = Bundle.main.bundleIdentifier ?? "com.xyfamily"

    func saveRefreshToken(_ token: String) throws {
        try KeychainWrapper.standard.set(token, forKey: "xyfamily_refresh_token", in: service)
    }

    func getRefreshToken() throws -> String? {
        return try KeychainWrapper.standard.get("xyfamily_refresh_token", in: service)
    }
}
```

---

## 六、权限管理

### 6.1 权限校验

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

## 七、Token 自动刷新

```swift
actor TokenRefreshManager {
    private var isRefreshing = false

    func refreshToken() async throws -> String {
        guard !isRefreshing else {
            // 等待其他请求的刷新结果
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

## 八、关联文档

- [前端通用代码规范](./04.02-分端代码规范/前端通用代码规范-V1.0.0.md)
- [接口总览](../05-接口与模块落地文档/05.01-接口总览/接口总览-V1.0.0.md)
- [Token 管理 PRD](../02-需求与产品设计/02.02-产品PRD/02-用户认证模块/03-Token管理/Token管理-V1.0.0.md)
