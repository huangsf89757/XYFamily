package init

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"xyfamily/pkg/logger"
)

type permDef struct {
	key, module, desc string
	isPublic         bool
}

var permissionPoints = []permDef{
	{"auth.register", "auth", "Register new account", true},
	{"auth.login", "auth", "Login to account", true},
	{"auth.reset_password", "auth", "Reset password via verification code", true},
	{"auth.logout", "auth", "Logout and invalidate tokens", false},
	{"auth.refresh", "auth", "Refresh access token", false},
	{"account.profile.read", "account", "Read own profile", false},
	{"account.profile.update", "account", "Update own profile", false},
	{"account.password.update", "account", "Change password", false},
	{"account.deactivate", "account", "Deactivate account", false},
	{"account.undeactivate", "account", "Restore deactivated account", false},
	{"account.third_party.bind", "account", "Bind third-party identity", false},
	{"account.third_party.unbind", "account", "Unbind third-party identity", false},
	{"org.create", "org", "Create organization", false},
	{"org.read", "org", "Read organization info", false},
	{"org.update", "org", "Update organization", false},
	{"org.member.invite", "org", "Invite member to organization", false},
	{"org.member.remove", "org", "Remove member from organization", false},
	{"org.member.list", "org", "List organization members", false},
	{"org.role.assign", "org", "Assign role in organization", false},
	{"org.role.downgrade", "org", "Downgrade role in organization", false},
	{"org.team.create", "org", "Create team in organization", false},
	{"team.read", "team", "Read team info", false},
	{"team.update", "team", "Update team", false},
	{"team.member.invite", "team", "Invite member to team", false},
	{"team.member.remove", "team", "Remove member from team", false},
	{"team.member.list", "team", "List team members", false},
	{"team.role.assign", "team", "Assign role in team", false},
	{"team.role.downgrade", "team", "Downgrade role in team", false},
	{"team.group.create", "team", "Create group in team", false},
	{"team.archive", "team", "Archive team", false},
	{"group.read", "group", "Read group info", false},
	{"group.update", "group", "Update group", false},
	{"group.delete", "group", "Delete group", false},
	{"group.member.invite", "group", "Invite member to group", false},
	{"group.member.remove", "group", "Remove member from group", false},
	{"group.member.list", "group", "List group members", false},
	{"group.role.assign", "group", "Assign role in group", false},
	{"group.role.downgrade", "group", "Downgrade role in group", false},
	{"audit.login.read", "audit", "Read login audit logs", false},
	{"audit.operation.read", "audit", "Read operation audit logs", false},
	{"audit.query", "audit", "Query audit logs", false},
	{"admin.config.read", "admin", "Read system configs", false},
	{"admin.config.update", "admin", "Update system configs", false},
	{"admin.force_downgrade", "admin", "Force downgrade any user role", false},
	{"admin.audit.global", "admin", "View all organization audit logs", false},
}

func SeedPermissions(ctx context.Context, pool *pgxpool.Pool) error {
	for _, p := range permissionPoints {
		_, err := pool.Exec(ctx,
			"INSERT INTO permission_points (permission_key, module, description, is_public) VALUES ($1, $2, $3, $4) ON CONFLICT (permission_key) DO NOTHING",
			p.key, p.module, p.desc, p.isPublic)
		if err != nil {
			return fmt.Errorf("seed permission %s: %w", p.key, err)
		}
	}
	logger.Get().Info("permissions seeded", zap.Int("count", len(permissionPoints)))
	return nil
}
