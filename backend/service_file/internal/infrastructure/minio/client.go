package minio

import (
	"bytes"
	"context"
	"fmt"
	"service-file/internal/domain"
	"service-file/pkg/utils"

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
	const op = "infrastructure.minio.client.CreateBucket"

	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}

func (m *MinIOClient) UploadFile(ctx context.Context, bucketName string, objectName string, data []byte, contentType string) error {
	const op = "infrastructure.minio.client.UploadFile"

	_, err := m.client.PutObject(ctx, bucketName, objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)

	}
	return nil
}

func (m *MinIOClient) UploadNewVersion(ctx context.Context, bucket string, baseKey string, data []byte, contentType string, version int) (string, error) {
	const op = "infrastructure.minio.client.UploadNewVersion"

	base, ext := utils.ParseBaseName(baseKey)

	// Формируем новое имя с версией
	newKey := fmt.Sprintf("%s_v%d%s", base, version, ext)

	// Загружаем файл
	err := m.UploadFile(ctx, bucket, newKey, data, contentType)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return newKey, nil
}

func (m *MinIOClient) GetFile(ctx context.Context, bucket string, key string) (*minio.Object, error) {
	const op = "infrastructure.minio.client.GetFile"

	object, err := m.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем существование файла
	stat, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, fmt.Errorf("%s: %w", op, domain.ErrFileNotFound)
	}

	if stat.Size == 0 {
		object.Close()
		return nil, fmt.Errorf("%s: %w", op, domain.ErrEmptyFile)
	}

	return object, nil
}
