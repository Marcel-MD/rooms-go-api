package rdb

import (
	"context"
	"os"
	"sync"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	once sync.Once
	rdb  *redis.Client
	ctx  context.Context
)

func GetRDB() (*redis.Client, context.Context) {
	once.Do(func() {
		rdb, ctx = initRDB()
	})
	return rdb, ctx
}

func initRDB() (*redis.Client, context.Context) {
	log.Info().Msg("Initializing redis")
	dsn := os.Getenv("REDIS_URL")

	opt, err := redis.ParseURL(dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse redis connection string")
	}

	opt.ReadTimeout = -1 // temporary fix until issue is resolved

	rdb := redis.NewClient(opt)
	ctx := context.Background()

	status := rdb.Ping(ctx)
	if status.Err() != nil {
		log.Fatal().Err(status.Err()).Msg("Failed to connect to redis")
	}

	return rdb, ctx
}
