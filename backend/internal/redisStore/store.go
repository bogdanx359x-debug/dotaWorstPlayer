package redisstore

import (
	"context"
	"dotaWorstPlayerChacker/internal/openDota"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const heroesHashKey = "dota:heroes"

type Store struct {
	rdb *redis.Client
}

func NewRedis() *Store {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	pass := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     pass,
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	return &Store{rdb: rdb}
}

func (s *Store) RefreshHeroes(ctx context.Context, heroes openDota.Heroes) error {
	fields := make([]interface{}, 0, len(heroes)*2)
	for _, h := range heroes {
		fields = append(fields, strconv.Itoa(h.HeroID), h.LocalizedName)
	}
	if err := s.rdb.Del(ctx, heroesHashKey).Err(); err != nil {
		return err
	}
	if len(fields) > 0 {
		if err := s.rdb.HSet(ctx, heroesHashKey, fields...).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) GetHeroName(ctx context.Context, id int) (string, bool, error) {
	name, err := s.rdb.HGet(ctx, heroesHashKey, strconv.Itoa(id)).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return name, true, nil
}

func (s *Store) RDB() *redis.Client { return s.rdb }
