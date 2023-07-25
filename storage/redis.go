package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"todolist/config"
	"todolist/entity"
)

func NewRedisClient(cfg config.Config) redis.UniversalClient {
	if cfg.RedisDriver == `sentinel` {
		return redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:    cfg.RedisMasterName,
			SentinelAddrs: []string{cfg.RedisHost + `:` + cfg.RedisPort},
			Password:      strings.TrimSpace(cfg.RedisPassword),
			DB:            cfg.RedisUseDefaultDB,
		})
	}

	return redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%v:%v", cfg.RedisHost, cfg.RedisPort),
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
