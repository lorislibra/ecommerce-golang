package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/donnjedarko/paninaro/infrastructures/db"
	"github.com/donnjedarko/paninaro/src/domains"
	"github.com/go-redis/redis/v8"
)

type userTokenRepository struct {
	redis db.RedisInstance
}

func NewRefreshTokenRepository(redis db.RedisInstance) domains.RefreshTokenRepository {
	return &userTokenRepository{
		redis: redis,
	}
}

func (u *userTokenRepository) Client() *redis.Client {
	return u.redis.Client()
}

func (u *userTokenRepository) Exist(ctx context.Context, uid string, tokenId string) (bool, error) {
	key := fmt.Sprintf("%s:%s", uid, tokenId)

	res := u.Client().Get(ctx, key)
	if _, err := res.Result(); err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (u *userTokenRepository) SaveAndDelete(ctx context.Context, uid string, tokenId string, refreshToken string, expire time.Duration, oldTokenId string) error {
	key := fmt.Sprintf("%s:%s", uid, tokenId)
	keyOld := fmt.Sprintf("%s:%s", uid, oldTokenId)

	p := u.Client().Pipeline()
	p.Set(ctx, key, refreshToken, expire)
	p.Del(ctx, keyOld)

	if _, err := p.Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (u *userTokenRepository) Save(ctx context.Context, uid string, tokenId string, refreshToken string, expire time.Duration) error {
	key := fmt.Sprintf("%s:%s", uid, tokenId)

	res := u.Client().Set(ctx, key, refreshToken, expire)
	if _, err := res.Result(); err != nil {
		return err
	}

	return nil
}

func (u *userTokenRepository) Delete(ctx context.Context, uid string, tokenId string) (bool, error) {
	key := fmt.Sprintf("%s:%s", uid, tokenId)

	res := u.Client().Del(ctx, key)
	n, err := res.Result()
	if err != nil {
		return false, err
	}

	if n > 0 {
		return true, nil
	}

	return false, nil
}

func (u *userTokenRepository) DeleteAll(ctx context.Context, uid string) (bool, error) {
	key := fmt.Sprintf("%s:*", uid)

	client := u.Client()
	iter := client.Scan(ctx, 0, key, 10).Iterator()

	var deletedCount int64

	for iter.Next(ctx) {
		deleted, err := client.Del(ctx, iter.Val()).Result()
		if err != nil {
			return false, err
		}
		deletedCount += deleted
	}

	if err := iter.Err(); err != nil {
		return false, err
	}

	if deletedCount > 0 {
		return true, nil
	}

	return false, nil
}
