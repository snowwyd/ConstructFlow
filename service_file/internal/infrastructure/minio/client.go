package minio

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	cfg    *Config
}

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

func NewMinIOClient(cfg *Config) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOClient{
		client: client,
		cfg:    cfg,
	}, nil
}

func (m *MinIOClient) CreateBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

func (m *MinIOClient) UploadFile(ctx context.Context, bucketName string, objectName string, data []byte) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return err
}

func (m *MinIOClient) UploadNewVersion(ctx context.Context, bucket string, baseKey string, data []byte) (string, error) {
	// Разделяем имя файла и расширение
	ext := filepath.Ext(baseKey)
	baseName := strings.TrimSuffix(baseKey, ext)

	// Формируем новое имя с версией
	newKey := fmt.Sprintf("%s_v%d%s", baseName, time.Now().Unix(), ext)

	// Загружаем файл
	err := m.UploadFile(ctx, bucket, newKey, data)
	if err != nil {
		return "", err
	}

	return newKey, nil
}

func (m *MinIOClient) GetFileURL(ctx context.Context, bucketName string, objectName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment")

	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName,
		time.Second*24*60*60, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
