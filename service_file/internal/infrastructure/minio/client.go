package minio

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"service-file/internal/domain"
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

func (m *MinIOClient) UploadFile(ctx context.Context, bucketName string, objectName string, data []byte, contentType string) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	return err
}

func (m *MinIOClient) UploadNewVersion(ctx context.Context, bucket string, baseKey string, data []byte, contentType string) (string, error) {
	// Разделяем имя файла и расширение
	ext := filepath.Ext(baseKey)
	baseName := strings.TrimSuffix(baseKey, ext)

	// Формируем новое имя с версией
	newKey := fmt.Sprintf("%s_v%d%s", baseName, time.Now().Unix(), ext)

	// Загружаем файл
	err := m.UploadFile(ctx, bucket, newKey, data, contentType)
	if err != nil {
		return "", err
	}

	return newKey, nil
}

func (m *MinIOClient) GetFile(ctx context.Context, bucket string, key string) (*minio.Object, error) {
	object, err := m.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	// Проверяем существование файла
	stat, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, domain.ErrFileNotFound
	}

	if stat.Size == 0 {
		object.Close()
		return nil, fmt.Errorf("empty file in MinIO")
	}

	return object, nil
}
