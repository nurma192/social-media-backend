package awsStorage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"mime/multipart"
	"strings"
)

func (a *AWSService) UploadFile(file multipart.File, filename, contentType string) (string, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}

	_, err = a.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &a.Bucket,
		Key:         &filename,
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("Error when try to save to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", a.Bucket, filename)
	return fileURL, nil
}

func (a *AWSService) GetFile(filename string) (*s3.GetObjectOutput, error) {
	out, err := a.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &a.Bucket,
		Key:    &filename,
	})
	if err != nil {
		return nil, fmt.Errorf("Error when try to get file from  S3: %w", err)
	}
	return out, nil
}

func (a *AWSService) DeleteFile(fileUrl string) error {
	parts := strings.Split(fileUrl, "/")
	filename := parts[len(parts)-1]

	_, err := a.S3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &a.Bucket,
		Key:    &filename,
	})
	if err != nil {
		return fmt.Errorf("Error when try to delete file from S3: %w", err)
	}
	return nil
}
