package storage

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
	"todolist/entity"
)

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:16379",
			Password: "", // no password set
			DB:       0,  // use default DB
		},
	)
	return rdb
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
