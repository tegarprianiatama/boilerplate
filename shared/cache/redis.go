package cache

import (
	"boilerplate/shared/config"
	"boilerplate/shared/logger"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const redisConnectTimeout = 5 * time.Second

// OpenRedis creates a Redis client and verifies the connection with a PING.
func OpenRedis() (*redis.Client, error) {
	cfg := config.GetConfig()
	log := logger.Log()

	addr := net.JoinHostPort(cfg.RedisHost, cfg.RedisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), redisConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", addr, err)
	}

	log.Info("Connected to Redis", zap.String("addr", addr), zap.Int("db", cfg.RedisDB))
	return client, nil
}
