package interfaces

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type FileStorage interface {
	GetFile(ctx context.Context, bucket string, key string) (*minio.Object, error)

	UploadFile(ctx context.Context, bucketName string, objectName string, data []byte, contentType string) error
	UploadNewVersion(ctx context.Context, bucket string, baseKey string, data []byte, contentType string) (string, error)
}
