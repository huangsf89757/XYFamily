package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"xyfamily/pkg/config"
	"xyfamily/pkg/logger"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(cfg *config.Config) (*DB, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}
	poolCfg.MaxConns = cfg.Database.MaxConns
	poolCfg.MinConns = cfg.Database.MinConns

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	logger.Get().Info("database connected", zap.String("host", cfg.Database.Host))
	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	logger.Get().Info("redis connected", zap.String("addr", cfg.Redis.Addr()))
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
