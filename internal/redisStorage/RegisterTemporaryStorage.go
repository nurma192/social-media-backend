package redisStorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"social-media-back/models/request"
	"time"
)

func (s *RedisService) SaveRegisteredUserData(data *request.RegisterRequest) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}

	err = s.Client.Set(s.ctx, "register:"+data.Email, dataJson, 10*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("error saving data to redis: %v", err)
	}

	return nil
}

func (s *RedisService) GetRegisteredUserByEmail(email string) (*request.RegisterRequest, error) {
	dataJSON, err := s.Client.Get(s.ctx, "register:"+email).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	if len(dataJSON) == 0 {
		return nil, nil
	}
	var data request.RegisterRequest
	err = json.Unmarshal([]byte(dataJSON), &data)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %v", err)
	}
	return &data, nil
}

func (s *RedisService) DeleteRegisteredUserByEmail(email string) error {
	err := s.Client.Del(s.ctx, "register:"+email).Err()
	if err != nil {
		return fmt.Errorf("error deleting register: %v", err)
	}
	return nil
}
