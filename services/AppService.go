package services

import (
	"social-media-back/internal/awsStorage"
	"social-media-back/internal/mail"
	"social-media-back/internal/redisStorage"
	"social-media-back/internal/storage/DatabaseService"
	"social-media-back/internal/token"
)

type AppService struct {
	DBService    *DatabaseService.DBService
	JWTService   *token.JWTService
	RedisService *redisStorage.RedisService
	AWSService   *awsStorage.AWSService
	EmailService *mail.EmailService
}

func NewAppService(dbService *DatabaseService.DBService, jwtService *token.JWTService, redisService *redisStorage.RedisService, awsService *awsStorage.AWSService, emailService *mail.EmailService) *AppService {
	return &AppService{
		DBService:    dbService,
		JWTService:   jwtService,
		RedisService: redisService,
		AWSService:   awsService,
		EmailService: emailService,
	}
}
