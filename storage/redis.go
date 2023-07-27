package storage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
	"todolist/config"
	"todolist/entity"
)

func NewRedisClient(cfg config.Config) (c redis.UniversalClient, err error) {
	switch cfg.RedisDriver {
	case "sentinel":
		c = SentinelRedisClient(cfg)
	case "cluster":
		c = ClusterRedisClient(cfg)
	default:
		c = RedisClient(cfg)
	}

	err = c.Ping(context.Background()).Err()

	return c, err
}

func SentinelRedisClient(cfg config.Config) *redis.ClusterClient {
	return redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:    cfg.RedisMasterName,
		SentinelAddrs: cfg.RedisSentinelAddr,
		Password:      cfg.RedisPassword,
		DB:            cfg.RedisUseDefaultDB,
	})
}

func ClusterRedisClient(cfg config.Config) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    cfg.RedisClusterPorts,
		Password: cfg.RedisPassword,
	})
}

func RedisClient(cfg config.Config) *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisUseDefaultDB,
		},
	)
}

func (u *UserStorage) SaveSessionRedis(ctx context.Context, session entity.Session) error {
	sessionJson, err := json.Marshal(session)
	if err != nil {
		return err
	}

	err = u.cache.Set(ctx, session.ID.String(), sessionJson, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStorage) SessionByIDRedis(ctx context.Context, id uuid.UUID) (entity.Session, error) {
	var session entity.Session

	res, err := u.cache.Get(ctx, id.String()).Bytes()
	if err != nil {
		return entity.Session{}, err
	}

	err = json.Unmarshal(res, &session)

	return session, nil
}
