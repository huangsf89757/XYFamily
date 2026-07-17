package init

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"xyfamily/pkg/logger"
)

func SeedSystemConfigs(ctx context.Context, pool *pgxpool.Pool) error {
	configs := []struct {
		key, value, desc string
	}{
		{"token.access_ttl", "1800", "Access Token TTL in seconds (30min)"},
		{"token.refresh_ttl", "604800", "Refresh Token TTL in seconds (7d)"},
		{"login.rate_limit.threshold", "5", "Login failure threshold"},
		{"login.rate_limit.window", "300", "Login rate limit window in seconds"},
		{"login.rate_limit.lock", "900", "Login lockout duration in seconds"},
		{"verification_code.ttl", "300", "Verification code TTL in seconds"},
		{"audit.retention_days", "365", "Audit log retention days"},
		{"account.deactivate_grace_days", "30", "Account deactivation grace period days"},
		{"login.auto_register", "false", "Auto register on login (P2 feature, default off)"},
	}
	for _, c := range configs {
		_, err := pool.Exec(ctx,
			"INSERT INTO system_configs (config_key, config_value, description) VALUES ($1, $2, $3) ON CONFLICT (config_key) DO NOTHING",
			c.key, c.value, c.desc)
		if err != nil {
			return fmt.Errorf("seed config %s: %w", c.key, err)
		}
	}
	logger.Get().Info("system_configs seeded", zap.Int("count", len(configs)))
	return nil
}

func SeedAll(ctx context.Context, pool *pgxpool.Pool) error {
	if err := SeedRoles(ctx, pool); err != nil {
		return err
	}
	if err := SeedPermissions(ctx, pool); err != nil {
		return err
	}
	if err := SeedRolePermissions(ctx, pool); err != nil {
		return err
	}
	if err := SeedSystemConfigs(ctx, pool); err != nil {
		return err
	}
	logger.Get().Info("all seed data initialized successfully")
	return nil
}
