# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**XYFamily** — a multi-tenant account & permission platform (SaaS 多租户账号权限底座) in the **pre-implementation / documentation phase**. There is **no application code yet** — the repo currently holds product specs (PRD), a structured Wiki, and a plan to generate the architecture & dev-standard docs. Future iterations will scaffold a Go backend.

### Tech Stack (target, per PRD + plan)
- **Language / Framework**: Go (1.22+) + Gin
- **Database**: PostgreSQL 17
- **Cache / Auth state**: Redis (verification codes w/ TTL, login rate-limiter, token blacklist, distributed locks)
- **Auth**: JWT — Access Token 30min (stateless), Refresh Token 7 days (stored in `sessions`), HS256
- **Storage**: MinIO/S3 for avatar uploads
- **API**: JSON envelope, base path `/api/v1/`, CORS
- **Deployment target**: stateless multi-instance + LB + PG master/slave + Redis Sentinel

## Repository Structure

```
.
├── README.md                  # project stub
├── PRD-Demo.md                # consolidated product requirements (the source PRD)
├── LICENSE
├── .gitignore                 # (⚠️ currently a generic Xcode/Swift template — not Go-relevant; review before scaffolding code)
├── .zcode/plans/              # internal work plan
└── wiki/                      # structured knowledge base (10-level hierarchy)
    ├── README.md              # wiki navigation & naming/copyright conventions
    ├── 00-项目总览/              # project overview, milestones, risk ledger (partly empty)
    ├── 01-立项与项目管理前置/     # pre-project management (skeleton only)
    ├── 02-需求与产品设计/         # **product PRD — 39 docs, the canonical requirements source**
    ├── 03-技术架构与方案设计/     # architecture / ADR / DB design (✅ filled — 8 docs)
    ├── 04-开发规范与编码手册/     # dev standards (✅ filled — 11 docs)
    ├── 05-接口与模块落地文档/     # interface spec (✅ filled — 11 docs)
    ├── 06-测试与质量保障/         # (✅ filled — 3 docs)
    ├── 07-部署运维与应急故障/     # (✅ filled — 3 docs)
    ├── 08-合规安全与数据治理/     # (✅ filled — 3 docs)
    ├── 09-迭代复盘与业务沉淀/     # (⬜ skeleton — README + 模板占位)
    └── 10-团队资产与培训模板/     # (✅ filled — 3 docs)
```

## Canonical Sources (read these before implementing)

- **PRD (product truth)**: `wiki/02-需求与产品设计/02.02-产品PRD/PRD模块总览-V1.0.0.md` and its 39 section docs. Root consolidated version is `PRD-Demo.md`.
- **PRD review gaps**: `wiki/02-需求与产品设计/02.04-需求复盘/PRD审查报告-V1.0.0.md` — P0 gaps that the architecture layer must close (super-admin init, Session/JWT strategy, multi-tenant isolation, account identity, role binding).
- **Architecture plan**: `.zcode/plans/plan-sess_3defe108-*.md` — describes the 16 docs to author across 03 → 04 → 00 before code begins.

## Key Design Decisions (already fixed)

- **Hierarchy**: Organization → Team → Group → Member, fully data-isolated between organizations (application-layer + `org_id`, not PG RLS).
- **RBAC**: 6 层角色体系（含 Public 第 0 层），45 个权限点。Single-role-per-member (role stored in `*_members.role` column). Permission inheritance L5→L1.
- **Accounts**: UUID primary key `account_id` + global-unique indexes on phone/email/username.
- **JWT**: Access = stateless short; Refresh = stateful in `sessions` table + blacklist on logout/password change.
- **Soft delete** (`deleted_at`) on all entity tables; 30-day grace period for account deactivation.
- **Security constants**: bcrypt cost 12; verification code 6-digit / 5-min TTL; login rate limit 5 fails per IP in 5 min → 15 min lockout; audit logs retained 1 year, anonymized on deletion; DB pool min 10 / max 100.

## Wiki Conventions (follow when editing wiki/ docs)

- **Naming**: `文档名称-V1.0.0.md` (semantic version in filename).
- **Required frontmatter per doc**: 密级 / 版本 / 编写人 / 审核人 / 生效时间 / 废弃时间 / 关联标签 / 关联目录, plus a 变更记录 table.
- **`编写人` / `变更人`**: use the tool/agent name that performed the edit (e.g., `ClaudeCode`, `CatPaw`, `ZCode`), or the editor's real name for manual edits.
- **Document hierarchy** is the 10-level directory system documented in `wiki/README.md`.

## Development Workflow (once code is scaffolded)

- Target layout (per plan): `cmd/xyfamily/`, `internal/{handler,service,repository,model,middleware}/`, `pkg/`, `configs/`, `migrations/`.
- Build/test tooling to be defined in `wiki/04-04` (Make, go build/test, golangci-lint).
- No `go.mod` / code exists yet — first code task will be scaffolding the Go project per `03.01` and `04.01`.

## Important Caveats

- The `.gitignore` is an **Xcode/Swift template**, not appropriate for a Go project — replace it when scaffolding Go code.
- `README.md` is a 20-byte stub; update it once the project shape is settled.
- PRD-Demo.md at the root and the wiki/ PRD sections can drift — treat `wiki/02.02-产品PRD/` as authoritative per the plan.
