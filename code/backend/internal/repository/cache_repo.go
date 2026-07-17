package repository

import (
	"context"
	"time"
	"fmt"

	"github.com/redis/go-redis/v9"

)

type CacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rc *RedisClient) *CacheRepository {
	return &CacheRepository{rdb: rc.Client}
}

func (r *CacheRepository) SetVerificationCode(ctx context.Context, target, codeType, code string, ttl int) error {
	key := fmt.Sprintf("vc:%s:%s", target, codeType)
	return r.rdb.Set(ctx, key, code, 0).Err()
}

func (r *CacheRepository) GetVerificationCode(ctx context.Context, target, codeType string) (string, error) {
	key := fmt.Sprintf("vc:%s:%s", target, codeType)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *CacheRepository) DelVerificationCode(ctx context.Context, target, codeType string) error {
	key := fmt.Sprintf("vc:%s:%s", target, codeType)
	return r.rdb.Del(ctx, key).Err()
}

func (r *CacheRepository) SetRateLimit(ctx context.Context, key string, threshold, window int) (bool, error) {
	cnt, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if cnt == 1 {
		r.rdb.Expire(ctx, key, 0)
	}
	if cnt > int64(threshold) {
		return false, nil
	}
	return true, nil
}

func (r *CacheRepository) IsRateLimited(ctx context.Context, key string) (bool, error) {
	val, err := r.rdb.Get(ctx, key).Int()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

func (r *CacheRepository) SetLockout(ctx context.Context, key string, ttl int) error {
	return r.rdb.Set(ctx, key, "1", 0).Err()
}

func (r *CacheRepository) BlacklistToken(ctx context.Context, jti string, ttl int) error {
	key := fmt.Sprintf("bl:%s", jti)
	return r.rdb.Set(ctx, key, "1", 0).Err()
}

func (r *CacheRepository) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	key := fmt.Sprintf("bl:%s", jti)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val != "", nil
}

func (r *CacheRepository) IncrResetCount(ctx context.Context, accountID string, ttl int) (int64, error) {
	key := fmt.Sprintf("reset:acct:%s:hourly", accountID)
	cnt, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if cnt == 1 {
		r.rdb.Expire(ctx, key, 0)
	}
	return cnt, nil
}
func (r *CacheRepository) PushAuditEvent(ctx context.Context, data map[string]interface{}) error {
	args := make([]interface{}, 0, len(data)*2)
	for k, v := range data { args = append(args, k, v) }
	return r.rdb.XAdd(ctx, &redis.XAddArgs{Stream: "audit_stream", MaxLen: 100000, Approx: true, Values: args}).Err()
}
func (r *CacheRepository) ReadAuditEvents(ctx context.Context, group, consumer string, count int64, block time.Duration) ([]redis.XStream, error) {
	return r.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{Group: group, Consumer: consumer, Streams: []string{"audit_stream", ">"}, Count: count, Block: block}).Result()
}
func (r *CacheRepository) AckAuditEvent(ctx context.Context, group string, ids ...string) error {
	return r.rdb.XAck(ctx, "audit_stream", group, ids...).Err()
}
func (r *CacheRepository) CreateAuditGroup(ctx context.Context, group string) error {
	return r.rdb.XGroupCreate(ctx, "audit_stream", group, "$").Err()
}
