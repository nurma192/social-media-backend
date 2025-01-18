package awsStorage

import (
	"context"
	"log"

	//"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MyConfig struct {
	AWSRegion      string
	AWSAccessKeyID string
	AWSSecretKey   string
	AWSS3Bucket    string
}

type AWSService struct {
	S3Client *s3.Client
	Bucket   string
}

var awsServices *AWSService

func InitAWS(cfg *MyConfig) *AWSService { // Используем вашу пользовательскую структуру
	if awsServices != nil {
		return awsServices
	}

	// Создаем AWS конфигурацию из ключей
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKeyID,
			cfg.AWSSecretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Ошибка при загрузке AWS конфигурации: %v", err)
	}

	// Инициализация S3 клиента
	s3Client := s3.NewFromConfig(awsCfg)

	awsServices = &AWSService{
		S3Client: s3Client,
		Bucket:   cfg.AWSS3Bucket,
	}

	log.Println("AWS Services успешно инициализированы")
	return awsServices
}
