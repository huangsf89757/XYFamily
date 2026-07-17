package init

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"xyfamily/pkg/logger"
)

var rolePermMap = map[string][]string{
	"public": {"auth.register", "auth.login", "auth.reset_password"},
	"regular_member": {"account.profile.read", "account.profile.update", "account.password.update", "account.deactivate", "account.undeactivate"},
	"group_ordinary_admin": {"group.read", "group.member.list", "group.member.invite", "group.member.remove"},
	"group_core_admin": {"group.read", "group.update", "group.delete", "group.member.list", "group.member.invite", "group.member.remove", "group.role.assign", "group.role.downgrade"},
	"team_ordinary_admin": {"team.read", "team.member.list", "team.member.invite", "team.member.remove"},
	"team_core_admin": {"team.read", "team.update", "team.archive", "team.member.list", "team.member.invite", "team.member.remove", "team.role.assign", "team.role.downgrade", "team.group.create"},
	"organization_ordinary_admin": {"org.read", "org.member.list", "org.member.invite", "org.member.remove", "org.team.create"},
	"organization_core_admin": {"org.read", "org.update", "org.member.list", "org.member.invite", "org.member.remove", "org.role.assign", "org.role.downgrade", "org.team.create", "audit.login.read", "audit.operation.read", "audit.query"},
}

func SeedRolePermissions(ctx context.Context, pool *pgxpool.Pool) error {
	count := 0
	for roleKey, perms := range rolePermMap {
		for _, perm := range perms {
			_, err := pool.Exec(ctx,
				"INSERT INTO role_permissions (role_key, permission_key) VALUES ($1, $2) ON CONFLICT (role_key, permission_key) DO NOTHING",
				roleKey, perm)
			if err != nil {
				return fmt.Errorf("seed role_perm %s->%s: %w", roleKey, perm, err)
			}
			count++
		}
	}
	allPerms := []string{}
	for _, p := range permissionPoints {
		allPerms = append(allPerms, p.key)
	}
	for _, perm := range allPerms {
		_, err := pool.Exec(ctx,
			"INSERT INTO role_permissions (role_key, permission_key) VALUES ($1, $2) ON CONFLICT (role_key, permission_key) DO NOTHING",
			"super_admin", perm)
		if err != nil {
			return fmt.Errorf("seed super_admin perm %s: %w", perm, err)
		}
		count++
	}
	rm := []string{"account.profile.read", "account.profile.update", "account.password.update", "account.deactivate", "account.undeactivate"}
	for _, perm := range rm {
		for _, rk := range []string{"group_ordinary_admin", "group_core_admin", "team_ordinary_admin", "team_core_admin", "organization_ordinary_admin", "organization_core_admin"} {
			_, err := pool.Exec(ctx, "INSERT INTO role_permissions (role_key, permission_key) VALUES ($1, $2) ON CONFLICT DO NOTHING", rk, perm)
			if err != nil {
				return fmt.Errorf("seed inherit %s->%s: %w", rk, perm, err)
			}
			count++
		}
	}
	ginh := []string{"group.read"}
	for _, perm := range ginh {
		for _, rk := range []string{"team_ordinary_admin", "team_core_admin", "organization_ordinary_admin", "organization_core_admin"} {
			_, err := pool.Exec(ctx, "INSERT INTO role_permissions (role_key, permission_key) VALUES ($1, $2) ON CONFLICT DO NOTHING", rk, perm)
			if err != nil {
				return fmt.Errorf("seed inherit %s->%s: %w", rk, perm, err)
			}
			count++
		}
	}
	logger.Get().Info("role_permissions seeded", zap.Int("entries", count))
	return nil
}
