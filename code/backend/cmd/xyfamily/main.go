package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"xyfamily/internal/handler"
	initpkg "xyfamily/internal/init"
	"xyfamily/internal/middleware"
	"xyfamily/internal/repository"
	"xyfamily/internal/service"
	"xyfamily/pkg/config"
	"xyfamily/pkg/jwt"
	"xyfamily/pkg/logger"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil { fmt.Printf("load config failed: %v\n", err); os.Exit(1) }
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil { fmt.Printf("init logger failed: %v\n", err); os.Exit(1) }
	defer logger.Sync()
	// 安全加固（RISK-001）：JWT 密钥必须通过配置或环境变量 XYFAMILY_JWT_SECRET 注入，缺失则阻断启动
	if v := os.Getenv("XYFAMILY_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if cfg.JWT.Secret == "" {
		logger.Get().Fatal("jwt secret is empty; set configs jwt.secret or environment variable XYFAMILY_JWT_SECRET before starting")
	}
	logger.Get().Info("starting xyfamily backend", zap.String("mode", cfg.Server.Mode))
	db, err := repository.NewDB(cfg)
	if err != nil { logger.Get().Fatal("connect database failed", zap.Error(err)) }
	defer db.Close()
	redisClient, err := repository.NewRedis(cfg)
	if err != nil { logger.Get().Fatal("connect redis failed", zap.Error(err)) }
	defer redisClient.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := initpkg.SeedAll(ctx, db.Pool); err != nil { logger.Get().Error("seed data failed", zap.Error(err)) }
	jwtMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	accountRepo := repository.NewAccountRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	cacheRepo := repository.NewCacheRepository(redisClient)
	rbacRepo := repository.NewRBACRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	authService := service.NewAuthService(accountRepo, sessionRepo, cacheRepo, rbacRepo, auditRepo, jwtMgr, cfg)
	accountService := service.NewAccountService(accountRepo, sessionRepo, cacheRepo, auditRepo)
	orgService := service.NewOrgService(tenantRepo, cacheRepo, auditRepo)
	teamService := service.NewTeamService(tenantRepo, auditRepo)
	adminService := service.NewAdminService(accountRepo, rbacRepo, tenantRepo, auditRepo, cacheRepo)
	auditConsumer := service.NewAuditConsumer(db.Pool, cacheRepo)
	_ = middleware.NewAuditMiddleware(cacheRepo) // reserved for audit middleware
	authMW := middleware.NewAuthMiddleware(jwtMgr, cacheRepo)
	rateLimitMW := middleware.NewRateLimitMiddleware(cacheRepo, cfg.Security.RateLimitThreshold, cfg.Security.RateLimitWindow, cfg.Security.RateLimitLock)
	membershipMW := middleware.NewMembershipValidator(rbacRepo, cacheRepo)
	_ = middleware.NewPermissionChecker(rbacRepo, cacheRepo) // reserved for permission checks
	healthHandler := handler.NewHealthHandler(db, redisClient)
	authHandler := handler.NewAuthHandler(authService, rateLimitMW)
	accountHandler := handler.NewAccountHandler(accountService)
	orgHandler := handler.NewOrgHandler(orgService)
	teamHandler := handler.NewTeamHandler(teamService)
	adminHandler := handler.NewAdminHandler(adminService)
gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery())
	// CORS 安全加固：禁止 "*" 与凭证（AllowCredentials）同时开启，浏览器禁止该组合；含 "*" 时自动关闭凭证
	corsOrigins := cfg.CORS.Origins
	allowCreds := true
	for _, o := range corsOrigins {
		if o == "*" {
			allowCreds = false
			break
		}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Organization-ID", "X-Team-ID", "X-Group-ID"},
		ExposeHeaders: []string{"Content-Length"},
		AllowCredentials: allowCreds,
	}))
	api := r.Group("/api/v1")
	{
		api.GET("/healthz", healthHandler.Healthz)
		auth := api.Group("/auth")
		{
			auth.POST("/verification-codes", authHandler.SendCode)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/reset-password", authHandler.ResetPassword)
			authProtected := auth.Group("")
			authProtected.Use(authMW.RequireAuth())
			authProtected.POST("/logout", authHandler.Logout)
		}
		protected := api.Group("")
		protected.Use(authMW.RequireAuth())
		{
			account := protected.Group("/account")
			{
				account.GET("/profile", accountHandler.GetProfile)
				account.PUT("/profile", accountHandler.UpdateProfile)
				account.PUT("/password", accountHandler.ChangePassword)
				account.POST("/deactivate", accountHandler.Deactivate)
				account.POST("/undeactivate", accountHandler.Undeactivate)
			}
			orgs := protected.Group("/organizations")
			orgs.Use(membershipMW.ValidateScope())
			{
				orgs.POST("", orgHandler.Create)
				orgs.GET("/:organization_id", orgHandler.GetInfo)
				orgs.PUT("/:organization_id", orgHandler.Update)
				orgs.POST("/:organization_id/disable", orgHandler.Disable)
				orgs.POST("/:organization_id/enable", orgHandler.Enable)
				orgs.POST("/:organization_id/members/invitations", orgHandler.Invite)
				orgs.GET("/:organization_id/members", orgHandler.ListMembers)
				orgs.PUT("/:organization_id/members/:account_id/role", orgHandler.AssignRole)
				orgs.POST("/:organization_id/members/:account_id/downgrade", orgHandler.Downgrade)
				orgs.DELETE("/:organization_id/members/:account_id", orgHandler.RemoveMember)
				orgs.POST("/:organization_id/teams", teamHandler.Create)
			}
			teams := protected.Group("/teams")
			teams.Use(membershipMW.ValidateScope())
			{
				teams.GET("/:team_id", teamHandler.GetInfo)
				teams.PUT("/:team_id", teamHandler.Update)
				teams.POST("/:team_id/archive", teamHandler.Archive)
				teams.POST("/:team_id/groups", teamHandler.CreateGroup)
			}
			groups := protected.Group("/groups")
			groups.Use(membershipMW.ValidateScope())
			{
				groups.GET("/:group_id", teamHandler.GetGroup)
				groups.PUT("/:group_id", teamHandler.UpdateGroup)
				groups.DELETE("/:group_id", teamHandler.DeleteGroup)
			}
			audit := protected.Group("/audit-logs")
			audit.Use(membershipMW.ValidateScope())
			{
				audit.GET("", adminHandler.GlobalAuditList)
				audit.GET("/:id", adminHandler.AuditDetail)
			}
			admin := protected.Group("/admin")
			{
				admin.POST("/init", adminHandler.Init)
				admin.GET("/config", adminHandler.GetConfig)
				admin.PUT("/config", adminHandler.UpdateConfig)
				admin.POST("/force-downgrade", adminHandler.ForceDowngrade)
				admin.GET("/audit-logs", adminHandler.GlobalAuditList)
				admin.GET("/audit-logs/:id", adminHandler.AuditDetail)
			}
		}
	}
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Get().Info("server starting", zap.String("addr", addr))
	go auditConsumer.Start(context.Background())
	if err := r.Run(addr); err != nil { logger.Get().Fatal("server failed", zap.Error(err)) }
}
