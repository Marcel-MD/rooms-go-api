package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Marcel-MD/rooms-go-api/rdb"
	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type ILoginLimiterService interface {
	IncrementAttempts(email string) error
}

type LoginLimiterService struct {
	rdb         *redis.Client
	ctx         context.Context
	maxAttempts int
	window      time.Duration
}

var (
	loginLimiterOnce    sync.Once
	loginLimiterService ILoginLimiterService
)

func GetLoginLimiterService() ILoginLimiterService {
	loginLimiterOnce.Do(func() {
		log.Info().Msg("Initializing loginLimiter service")

		attemptsStr := os.Getenv("LOGIN_ATTEMPTS")
		attempts, err := strconv.Atoi(attemptsStr)
		if err != nil {
			attempts = 5
		}

		windowStr := os.Getenv("LOGIN_WINDOW")
		window, err := time.ParseDuration(windowStr)
		if err != nil {
			window = 10 * time.Minute
		}

		rdb, ctx := rdb.GetRDB()

		loginLimiterService = &LoginLimiterService{
			rdb:         rdb,
			ctx:         ctx,
			maxAttempts: attempts,
			window:      window,
		}
	})
	return loginLimiterService
}

func (s *LoginLimiterService) IncrementAttempts(email string) error {
	now := time.Now().UnixNano()

	s.rdb.ZRemRangeByScore(s.ctx, email, "0", fmt.Sprint(now-(s.window.Nanoseconds()))).Result()

	attempts, err := s.rdb.ZRange(s.ctx, email, 0, -1).Result()
	if err != nil {
		return err
	}

	log.Info().Msg(fmt.Sprint(len(attempts)))

	if len(attempts) >= s.maxAttempts {
		return errors.New("too many attempts")
	}

	s.rdb.ZAddNX(s.ctx, email, redis.Z{Score: float64(now), Member: float64(now)})
	s.rdb.Expire(s.ctx, email, s.window)

	return nil
}
