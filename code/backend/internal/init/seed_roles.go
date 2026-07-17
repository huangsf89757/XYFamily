package init

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"xyfamily/pkg/logger"
)

func SeedRoles(ctx context.Context, pool *pgxpool.Pool) error {
	roles := []struct {
		roleKey, name, scope, desc string
		level                     int
	}{
		{"public", "Public", "global", "Public role, auto-assigned to every user", 0},
		{"regular_member", "RegularMember", "personal", "Regular member with basic operations", 1},
		{"group_ordinary_admin", "GroupOrdinaryAdmin", "group", "Group ordinary admin", 2},
		{"group_core_admin", "GroupCoreAdmin", "group", "Group core admin", 3},
		{"team_ordinary_admin", "TeamOrdinaryAdmin", "team", "Team ordinary admin", 4},
		{"team_core_admin", "TeamCoreAdmin", "team", "Team core admin", 5},
		{"organization_ordinary_admin", "OrganizationOrdinaryAdmin", "organization", "Organization ordinary admin", 6},
		{"organization_core_admin", "OrganizationCoreAdmin", "organization", "Organization core admin", 7},
		{"super_admin", "SuperAdmin", "global", "Super admin with all permissions", 8},
	}
	for _, r := range roles {
		_, err := pool.Exec(ctx,
			"INSERT INTO roles (role_key, name, level, scope, description) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (role_key) DO NOTHING",
			r.roleKey, r.name, r.level, r.scope, r.desc)
		if err != nil {
			return fmt.Errorf("seed role %s: %w", r.roleKey, err)
		}
	}
	logger.Get().Info("roles seeded", zap.Int("count", len(roles)))
	return nil
}
