# 安卓Kotlin代码规范

> XYFamily 安卓端（Kotlin + Jetpack Compose）**项目特定约定与选型**。通用 Kotlin/Android 编码规范请直接参考官方与业界标准，本文不再重复。

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

- **[Kotlin Coding Conventions](https://kotlinlang.org/docs/coding-conventions.html)** — Kotlin 官方编码规范
- **[Android Kotlin Style Guide](https://developer.android.com/kotlin/style-guide)** — Android 官方 Kotlin 风格指南
- **[Jetpack Compose Guidelines](https://developer.android.com/develop/ui/compose)** — Compose 官方文档
- **[Now in Android Style Guide](https://github.com/android/nowinandroid/blob/main/docs/ArchitectureLearningJourney.md)** — Google 官方 Android 架构指南

---

## 一、技术栈选型

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

## 三、API 调用（项目约定）

### 3.1 Token 自动注入拦截器

```kotlin
class AuthInterceptor : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = SecureStorage.getAccessToken()
            ?: return chain.proceed(chain.request())
        val newRequest = chain.request().newBuilder()
            .addHeader("Authorization", "Bearer $token")
            .build()
        return chain.proceed(newRequest)
    }
}
```

### 3.2 统一响应处理

```kotlin
sealed class ApiResult<T> {
    data class Success<T>(val data: T) : ApiResult<T>()
    data class Error(val code: Int, val message: String, val detail: String? = null) : ApiResult<Nothing>()
    data class Exception(val throwable: Throwable) : ApiResult<Nothing>()
}
```

---

## 四、安全存储（项目约定）

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

## 五、权限管理（项目角色）

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

---

## 六、错误处理（项目错误码映射）

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

## 七、关联文档

- [前端通用代码规范](./07-前端通用代码规范.md)
- [接口文档](../../../04-接口文档/接口文档.md)
