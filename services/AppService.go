package services

import (
	"social-media-back/internal/awsStorage"
	"social-media-back/internal/mail"
	"social-media-back/internal/redisStorage"
	"social-media-back/internal/storage"
	"social-media-back/internal/token"
)

type AppService struct {
	DBService    *storage.DBService
	JWTService   *token.JWTService
	RedisService *redisStorage.RedisService
	AWSService   *awsStorage.AWSService
	EmailService *mail.EmailService
}

func NewAppService(dbService *storage.DBService, jwtService *token.JWTService, redisService *redisStorage.RedisService, awsService *awsStorage.AWSService, emailService *mail.EmailService) *AppService {
	return &AppService{
		DBService:    dbService,
		JWTService:   jwtService,
		RedisService: redisService,
		AWSService:   awsService,
		EmailService: emailService,
	}
}
