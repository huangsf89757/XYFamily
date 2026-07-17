-- V001__init.sql
-- XYFamily 多租户账号权限底座 - 初始数据库迁移
-- 包含：13张表 + 枚举 + 索引

-- 启用 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ========================================
-- 1. accounts (账号主表)
-- ========================================
CREATE TABLE accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id      UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    phone_encrypted   BYTEA,
    phone_hash        VARCHAR(64) UNIQUE,
    email_encrypted   BYTEA,
    email_hash        VARCHAR(64) UNIQUE,
    username          VARCHAR(64) UNIQUE,
    password_hash     VARCHAR(255),
    nickname          VARCHAR(64),
    avatar            VARCHAR(512),
    status            VARCHAR(16) NOT NULL DEFAULT 'active',
    deactivated_at    TIMESTAMPTZ,
    previous_username  VARCHAR(64),
    username_changed_at TIMESTAMPTZ,
    last_login_at     TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_accounts_phone_hash    ON accounts(phone_hash)    WHERE deleted_at IS NULL;
CREATE INDEX idx_accounts_email_hash    ON accounts(email_hash)    WHERE deleted_at IS NULL;
CREATE INDEX idx_accounts_username      ON accounts(username)      WHERE deleted_at IS NULL;
CREATE INDEX idx_accounts_status        ON accounts(status);
CREATE INDEX idx_accounts_deleted       ON accounts(deleted_at);
CREATE INDEX idx_accounts_prev_username ON accounts(previous_username) WHERE previous_username IS NOT NULL;

-- ========================================
-- 2. sessions (Refresh Token 有状态存储)
-- ========================================
CREATE TABLE sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id        UUID NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    device            VARCHAR(128),
    client_ip         INET,
    user_agent        TEXT,
    expires_at        TIMESTAMPTZ NOT NULL,
    revoked_at        TIMESTAMPTZ,
    last_activity_at  TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_account   ON sessions(account_id);
CREATE INDEX idx_sessions_expires   ON sessions(expires_at);
CREATE INDEX idx_sessions_revoked   ON sessions(revoked_at) WHERE revoked_at IS NULL;

-- ========================================
-- 3. invitations (统一邀请流)
-- ========================================
CREATE TABLE invitations (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id              UUID NOT NULL,
    scope_type          VARCHAR(16) NOT NULL,
    scope_id            UUID NOT NULL,
    inviter_id          UUID NOT NULL,
    invitee_account_id  UUID,
    invitee_contact     VARCHAR(255),
    role                VARCHAR(32) NOT NULL,
    token               VARCHAR(255) NOT NULL UNIQUE,
    status              VARCHAR(16) NOT NULL DEFAULT 'pending',
    expired_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    accepted_at         TIMESTAMPTZ,
    deleted_at          TIMESTAMPTZ
);

CREATE INDEX idx_invitations_org       ON invitations(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invitations_scope     ON invitations(scope_type, scope_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invitations_invitee   ON invitations(invitee_account_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_invitations_contact   ON invitations(invitee_contact) WHERE deleted_at IS NULL;
CREATE INDEX idx_invitations_token     ON invitations(token);

-- ========================================
-- 4. organizations (组织)
-- ========================================
CREATE TABLE organizations (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name          VARCHAR(128) NOT NULL,
    description   VARCHAR(256),
    owner_id      UUID NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    disabled_at   TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_organizations_deleted ON organizations(deleted_at);
CREATE INDEX idx_organizations_owner   ON organizations(owner_id);

-- ========================================
-- 5. org_members (组织成员)
-- ========================================
CREATE TABLE org_members (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL,
    account_id    UUID NOT NULL,
    role          VARCHAR(32) NOT NULL,
    joined_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    UNIQUE (org_id, account_id)
);

CREATE INDEX idx_org_members_org       ON org_members(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_org_members_account   ON org_members(account_id) WHERE deleted_at IS NULL;

-- ========================================
-- 6. teams (团队)
-- ========================================
CREATE TABLE teams (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id       UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL,
    name          VARCHAR(128) NOT NULL,
    description   VARCHAR(256),
    owner_id      UUID NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    archived_at   TIMESTAMPTZ,
    deleted_at    TIMESTAMPTZ,
    UNIQUE (org_id, team_id)
);

CREATE INDEX idx_teams_org       ON teams(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_archived  ON teams(archived_at);
CREATE INDEX idx_teams_deleted   ON teams(deleted_at);

-- ========================================
-- 7. team_members (团队成员)
-- ========================================
CREATE TABLE team_members (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL,
    team_id       UUID NOT NULL,
    account_id    UUID NOT NULL,
    role          VARCHAR(32) NOT NULL,
    joined_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    UNIQUE (team_id, account_id)
);

CREATE INDEX idx_team_members_org       ON team_members(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_members_team      ON team_members(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_team_members_account   ON team_members(account_id) WHERE deleted_at IS NULL;

-- ========================================
-- 8. groups (小组)
-- ========================================
CREATE TABLE groups (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id      UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL,
    team_id       UUID NOT NULL,
    name          VARCHAR(64) NOT NULL,
    description   TEXT,
    owner_id      UUID NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_groups_org      ON groups(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_groups_team     ON groups(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_groups_deleted  ON groups(deleted_at);

-- ========================================
-- 9. group_members (小组成员)
-- ========================================
CREATE TABLE group_members (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id        UUID NOT NULL,
    team_id       UUID NOT NULL,
    group_id      UUID NOT NULL,
    account_id    UUID NOT NULL,
    role          VARCHAR(32) NOT NULL,
    joined_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ,
    UNIQUE (group_id, account_id)
);

CREATE INDEX idx_group_members_org       ON group_members(org_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_members_team      ON group_members(team_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_members_group     ON group_members(group_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_group_members_account   ON group_members(account_id) WHERE deleted_at IS NULL;

-- ========================================
-- 10. roles (角色)
-- ========================================
CREATE TABLE roles (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_key      VARCHAR(32) NOT NULL UNIQUE,
    name          VARCHAR(64) NOT NULL,
    level         SMALLINT NOT NULL,
    scope         VARCHAR(16) NOT NULL,
    description   TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ========================================
-- 11. permission_points (权限点)
-- ========================================
CREATE TABLE permission_points (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_key VARCHAR(64) NOT NULL UNIQUE,
    module        VARCHAR(16) NOT NULL,
    description   TEXT,
    is_public     BOOLEAN NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ========================================
-- 12. role_permissions (角色-权限映射)
-- ========================================
CREATE TABLE role_permissions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_key      VARCHAR(32) NOT NULL,
    permission_key VARCHAR(64) NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (role_key, permission_key)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role_key);

-- ========================================
-- 13. system_configs (系统配置)
-- ========================================
CREATE TABLE system_configs (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key    VARCHAR(64) NOT NULL UNIQUE,
    config_value  TEXT NOT NULL,
    description   TEXT,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_by    UUID
);

-- ========================================
-- 14. audit_logs (审计日志)
-- ========================================
CREATE TABLE audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id        UUID NOT NULL UNIQUE,
    account_id      UUID,
    org_id          UUID,
    action_domain   VARCHAR(16) NOT NULL,
    action_type     VARCHAR(48) NOT NULL,
    target_type     VARCHAR(32),
    target_id       UUID,
    result          VARCHAR(16),
    failure_reason  VARCHAR(64),
    login_method    VARCHAR(16),
    details         JSONB,
    trace_id        VARCHAR(64),
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_account   ON audit_logs(account_id, created_at DESC);
CREATE INDEX idx_audit_org       ON audit_logs(org_id, created_at DESC) WHERE org_id IS NOT NULL;
CREATE INDEX idx_audit_action    ON audit_logs(action_type, created_at DESC);
CREATE INDEX idx_audit_login_method ON audit_logs(login_method, created_at DESC) WHERE login_method IS NOT NULL;
CREATE INDEX idx_audit_created   ON audit_logs(created_at DESC);

-- ========================================
-- 初始化完成
-- ========================================
