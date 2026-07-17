package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"xyfamily/internal/repository"
	"xyfamily/pkg/logger"
)

type AuditConsumer struct {
	pool      *pgxpool.Pool
	cacheRepo *repository.CacheRepository
	group     string
	consumer  string
}

func NewAuditConsumer(pool *pgxpool.Pool, cache *repository.CacheRepository) *AuditConsumer {
	return &AuditConsumer{pool: pool, cacheRepo: cache, group: "audit_group", consumer: "consumer-1"}
}

func (ac *AuditConsumer) Start(ctx context.Context) {
	_ = ac.cacheRepo.CreateAuditGroup(ctx, ac.group)
	if err := ac.cacheRepo.CreateAuditGroup(ctx, ac.group); err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			logger.Get().Error("create audit consumer group", zap.Error(err))
		}
	}
	logger.Get().Info("audit consumer started", zap.String("group", ac.group))
	go ac.consume(ctx)
}

func (ac *AuditConsumer) consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Get().Info("audit consumer stopped")
			return
		default:
			streams, err := ac.cacheRepo.ReadAuditEvents(ctx, ac.group, ac.consumer, 10, 5*time.Second)
			if err != nil && err != redis.Nil {
				logger.Get().Error("read audit events", zap.Error(err))
				time.Sleep(time.Second)
				continue
			}
			if err == redis.Nil { continue }
			for _, stream := range streams {
				for _, msg := range stream.Messages {
					ac.process(ctx, msg)
				}
			}
		}
	}
}

func (ac *AuditConsumer) process(ctx context.Context, msg redis.XMessage) {
	defer func() { _ = ac.cacheRepo.AckAuditEvent(ctx, ac.group, msg.ID) }()
	eventID, _ := msg.Values["event_id"].(string)
	if eventID == "" { return }
	accountIDStr, _ := msg.Values["account_id"].(string)
	orgIDStr, _ := msg.Values["org_id"].(string)
	var accountID, orgID *uuid.UUID
	if accountIDStr != "" { id, err := uuid.Parse(accountIDStr); if err == nil { accountID = &id } }
	if orgIDStr != "" { id, err := uuid.Parse(orgIDStr); if err == nil { orgID = &id } }
	details, _ := msg.Values["details"].(string)
	var detailsJSON []byte
	if details != "" { detailsJSON = []byte(details) }
	_, err := ac.pool.Exec(ctx,
		"INSERT INTO audit_logs (event_id, account_id, org_id, action_domain, action_type, target_type, target_id, result, failure_reason, login_method, details, trace_id, ip_address, user_agent) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) ON CONFLICT (event_id) DO NOTHING",
		eventID, accountID, orgID,
		msg.Values["action_domain"], msg.Values["action_type"],
		msg.Values["target_type"], msg.Values["target_id"],
		msg.Values["result"], msg.Values["failure_reason"], msg.Values["login_method"],
		detailsJSON, msg.Values["trace_id"], msg.Values["ip_address"], msg.Values["user_agent"],
	)
	if err != nil { logger.Get().Error("insert audit log", zap.Error(err), zap.String("event_id", eventID)) }
	var _ = json.Marshal // keep import
}
