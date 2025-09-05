// Package redis implements redis connection.
package redis

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
)

// Redis -.
type Redis struct {
    Client *redis.Client
}

// New -.
func New(connString string) (*Redis, error) {
    opts, err := redis.ParseURL(connString)
    if err != nil {
        return nil, fmt.Errorf("redis - New - redis.ParseURL: %w", err)
    }

    client := redis.NewClient(opts)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err = client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("redis - New - client.Ping: %w", err)
    }

    log.Println("Redis connection established successfully")

    return &Redis{Client: client}, nil
}

// Close -.
func (r *Redis) Close() error {
    if err := r.Client.Close(); err != nil {
        return fmt.Errorf("redis - Close - client.Close: %w", err)
    }
    log.Println("Redis connection closed successfully")

    return nil
}