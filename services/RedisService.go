package services

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func (s *AppService) setVerificationCode(email string, code string) error {
	err := s.RedisClient.Set(s.RedisCtx, "verify:"+email, code, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	fmt.Println("setVerificationCode, savedCode:", code)

	return nil
}

func (s *AppService) getVerificationCode(email string) (string, error) {
	storedCode, err := s.RedisClient.Get(s.RedisCtx, "verify:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("verification code not found or code TTL expired")
		}
		return "", err
	}
	if storedCode == "" {
		return "", fmt.Errorf("verification code not found or code TTL expired")
	}
	fmt.Println("getVerificationCode, storedCode:", storedCode)
	return storedCode, nil
}

func (s *AppService) checkVerificationCode(email string) (bool, error) {
	storedCode, err := s.RedisClient.Get(s.RedisCtx, "verify:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	if storedCode == "" {
		return false, nil
	}
	fmt.Println("checkVerificationCode, storedCode:", storedCode)
	return true, nil
}

func (s *AppService) deleteVerificationCode(email string) error {
	err := s.RedisClient.Del(s.RedisCtx, "verify:"+email).Err()
	if err != nil {
		return err
	}
	return nil
}
