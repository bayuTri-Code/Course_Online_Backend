package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
	}
}

func (r *RedisService) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)
	
	err := r.client.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	
	return nil
}

func (r *RedisService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)
	
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	
	return result > 0, nil
}

func (r *RedisService) DeleteTokenCache(ctx context.Context, token string) error {
	key := fmt.Sprintf("token:%s", token)
	
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete token cache: %w", err)
	}
	
	return nil
}

func (r *RedisService) DeleteRefreshToken(ctx context.Context, userID string) error {
	key := fmt.Sprintf("refresh_token:%s", userID)
	
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	
	return nil
}

func (r *RedisService) DeleteOTPSession(ctx context.Context, email string) error {
	otpKey := fmt.Sprintf("otp:%s", email)
	err := r.client.Del(ctx, otpKey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete OTP: %w", err)
	}
	
	attemptsKey := fmt.Sprintf("otp_attempts:%s", email)
	r.client.Del(ctx, attemptsKey)
	
	return nil
}

func (r *RedisService) DeleteAllUserTokens(ctx context.Context, firebaseUID string) error {
	pattern := "token:*"
	
	var cursor uint64
	var deletedCount int
	
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan tokens: %w", err)
		}
		
		for _, key := range keys {
			val, err := r.client.Get(ctx, key).Result()
			if err != nil {
				continue
			}
			
			if val == firebaseUID {
				token := key[6:]
				
				r.BlacklistToken(ctx, token, 24*time.Hour)
				
				r.client.Del(ctx, key)
				
				deletedCount++
			}
		}
		
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	
	return nil
}
