package redisStorage

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func (s *RedisService) SetVerificationCode(email string, code string) error {
	err := s.Client.Set(s.ctx, "verify:"+email, code, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	fmt.Println("setVerificationCode, savedCode:", code)

	return nil
}

func (s *RedisService) GetVerificationCode(email string) (string, error) {
	storedCode, err := s.Client.Get(s.ctx, "verify:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		return "", err
	}
	if storedCode == "" {
		return "", nil
	}
	fmt.Println("getVerificationCode, storedCode:", storedCode)
	return storedCode, nil
}

func (s *RedisService) CheckVerificationCode(email string) (bool, error) {
	storedCode, err := s.Client.Get(s.ctx, "verify:"+email).Result()
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

func (s *RedisService) DeleteVerificationCode(email string) error {
	err := s.Client.Del(s.ctx, "verify:"+email).Err()
	if err != nil {
		return err
	}
	return nil
}
