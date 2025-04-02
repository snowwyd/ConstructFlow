package interfaces

import "context"

type FileStorage interface {
	UploadFile(ctx context.Context, bucketName string, objectName string, data []byte) error
	GetFileURL(ctx context.Context, bucketName string, objectName string) (string, error)
	UploadNewVersion(ctx context.Context, bucket string, baseKey string, data []byte) (string, error)
}
