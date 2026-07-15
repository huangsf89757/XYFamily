# 安卓（Kotlin）代码规范

> XYFamily 安卓端开发规范：Kotlin + Jetpack Compose。

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
| 语言 | Kotlin |
| UI 框架 | Jetpack Compose |
| 网络 | OkHttp + Retrofit |
| 数据持久化 | Room（SQLite） |
| 安全存储 | Android Keystore + EncryptedSharedPreferences |
| 依赖注入 | Hilt（Dagger） |
| 异步 | Kotlin Coroutines + Flow |
| 图片 | Coil |
| 状态管理 | Compose State + ViewModel |
| 最低版本 | Android API 24（7.0） |
| 目标版本 | Android API 35（15） |

---

## 二、项目目录结构

```
app/
  src/main/java/com/xyfamily/
    common/
      network/          # API 封装（Retrofit + OkHttp）
        RetrofitClient.kt
        AuthInterceptor.kt   # Token 注入 + 刷新
        ResponseHandler.kt   # 统一响应处理
      di/               # Hilt 依赖注入
      util/             # 工具类
    data/
      model/            # 数据模型
      repository/       # 数据仓库
    domain/
      usecase/          # 用例
    presentation/
      auth/             # 认证模块（登录/注册）
      account/          # 账号管理
      org/              # 组织管理
      team/             # 团队管理
      group/            # 小组管理
      rbac/             # 权限相关
      common/           # 公共组件
```

---

## 三、命名规范

### 3.1 基础命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 包名 | 全小写，`com.xyfamily.{module}` | `com.xyfamily.presentation.auth` |
| 类名 | PascalCase | `LoginViewModel`, `OrganizationRepository` |
| 函数/方法 | camelCase | `fetchUserList`, `handleLoginSuccess` |
| 常量 | UPPER_SNAKE_CASE | `MAX_LOGIN_ATTEMPTS`, `TOKEN_REFRESH_URL` |
| 变量 | camelCase | `currentUser`, `isLoggedIn` |
| 布尔值 | is/has/can 前缀 | `isLoading`, `hasPermission`, `canEdit` |
| 资源 ID | snake_case | `ic_user_avatar`, `tv_login_title` |

### 3.2 Compose 组件命名

```kotlin
// Composable 函数：PascalCase
@Composable
fun UserProfileCard(user: User) { ... }

// 状态变量
var isLoading by remember { mutableStateOf(false) }
var errorMessage by remember { mutableStateOf<String?>(null) }
```

---

## 四、API 调用封装

### 4.1 Retrofit 配置

```kotlin
@Module
@InstallIn(SingletonComponent::class)
object NetworkModule {

    @Provides
    @Singleton
    fun provideRetrofit(): Retrofit {
        val okHttpClient = OkHttpClient.Builder()
            .addInterceptor(AuthInterceptor())    // Token 注入
            .addInterceptor(TokenRefreshInterceptor()) // Token 自动刷新
            .build()

        return Retrofit.Builder()
            .baseUrl(BuildConfig.API_BASE_URL)
            .client(okHttpClient)
            .addConverterFactory(MoshiConverterFactory.create())
            .build()
    }
}
```

### 4.2 Token 注入拦截器

```kotlin
class AuthInterceptor : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = SecureStorage.getAccessToken()
            ?: return chain.proceed(chain.request())  // 无 Token，直接发请求

        val newRequest = chain.request().newBuilder()
            .addHeader("Authorization", "Bearer $token")
            .build()
        return chain.proceed(newRequest)
    }
}
```

### 4.3 统一响应处理

```kotlin
sealed class ApiResult<T> {
    data class Success<T>(val data: T) : ApiResult<T>()
    data class Error(val code: Int, val message: String, val detail: String? = null) : ApiResult<Nothing>()
    data class Exception(val throwable: Throwable) : ApiResult<Nothing>()
}

suspend fun <T> safeApiCall(apiCall: suspend ()): ApiResult<T> {
    return try {
        val response = apiCall()
        if (response.isSuccessful && response.body()?.code == 0) {
            ApiResult.Success(response.body()!!.data)
        } else {
            val errorBody = response.errorBody()?.string()
            val apiResponse = response.body()
            ApiResult.Error(
                code = apiResponse?.code ?: -1,
                message = apiResponse?.message ?: errorBody ?: "请求失败",
                detail = apiResponse?.detail
            )
        }
    } catch (e: Exception) {
        ApiResult.Exception(e)
    }
}
```

---

## 五、安全存储

### 5.1 Token 存储

```kotlin
// Access Token：不持久化（内存中）
class TokenManager {
    private var accessToken: String? = null

    fun setAccessToken(token: String) { accessToken = token }
    fun getAccessToken(): String? = accessToken
    fun clearTokens() { accessToken = null; securePrefs.edit().remove(KEY_REFRESH_TOKEN).apply() }
}

// Refresh Token：Keystore + EncryptedSharedPreferences
class SecureStorage {
    fun saveRefreshToken(token: String) {
        encryptedPrefs.edit().putString(KEY_REFRESH_TOKEN, token).apply()
    }
    fun getRefreshToken(): String? = encryptedPrefs.getString(KEY_REFRESH_TOKEN, null)
}
```

---

## 六、权限管理

### 6.1 权限点校验

```kotlin
// 前端权限校验（辅助，最终以后端为准）
class PermissionChecker(private val userRole: UserRole) {
    fun hasPermission(permCode: String): Boolean {
        return when (userRole) {
            UserRole.SUPER_ADMIN -> true
            UserRole.ORG_CORE_ADMIN -> ORG_CORE_ADMIN_PERMS.contains(permCode)
            UserRole.TEAM_CORE_ADMIN -> TEAM_CORE_ADMIN_PERMS.contains(permCode)
            UserRole.GROUP_CORE_ADMIN -> GROUP_CORE_ADMIN_PERMS.contains(permCode)
            UserRole.REGULAR_MEMBER -> REGULAR_MEMBER_PERMS.contains(permCode)
            UserRole.PUBLIC -> PUBLIC_PERMS.contains(permCode)
        }
    }
}
```

### 6.2 组织上下文

```kotlin
@HiltViewModel
class OrgContextViewModel @Inject constructor() : ViewModel() {
    private val _currentOrg = MutableStateFlow<String?>(null)
    val currentOrg: StateFlow<String?> = _currentOrg

    fun switchOrg(orgId: String) { _currentOrg.value = orgId }
}
```

---

## 七、错误处理与提示

### 7.1 统一错误处理

```kotlin
fun showError(context: Context, result: ApiResult<Nothing>) {
    when (val r = result) {
        is ApiResult.Error -> {
            val message = when (r.code) {
                101001..101009 -> "登录已过期，请重新登录"
                104290 -> "操作过于频繁，请稍后再试"
                114290 -> "验证码发送频率过高"
                603001, 603002 -> "当前账号无权限执行此操作"
                805000 -> "服务异常，请稍后重试"
                else -> r.message
            }
            Toast.makeText(context, message, Toast.LENGTH_SHORT).show()
        }
        is ApiResult.Exception -> {
            Toast.makeText(context, "网络连接异常", Toast.LENGTH_SHORT).show()
        }
        else -> {}
    }
}
```

---

## 八、关联文档

- [前端通用代码规范](./04.02-分端代码规范/前端通用代码规范-V1.0.0.md)
- [接口总览](../05-接口与模块落地文档/05.01-接口总览/接口总览-V1.0.0.md)
- [Token 管理 PRD](../02-需求与产品设计/02.02-产品PRD/02-用户认证模块/03-Token管理/Token管理-V1.0.0.md)
