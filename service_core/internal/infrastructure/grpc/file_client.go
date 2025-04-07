package grpc

import (
	"context"
	"fmt"
	"log"
	"service-core/internal/domain"
	"service-core/pkg/config"

	pb "service-core/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type FileGRPCClient struct {
	client pb.FileServiceClient
}

func NewFileGRPCClient(cfg *config.Config) (*FileGRPCClient, error) {
	apiKeyInterceptor := createAPIKeyInterceptor(cfg.SecretKeys.KeyGrpc)

	conn, err := grpc.Dial(cfg.GRPCAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(apiKeyInterceptor)) // Используйте TLS в production
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServiceClient(conn)
	return &FileGRPCClient{client: client}, nil
}

func createAPIKeyInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", apiKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *FileGRPCClient) GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error) {
	const op = "infrastructure.grpc.fileclient.GetFileWithDirectory"

	req := &pb.GetFileRequest{
		FileId: uint32(fileID),
	}

	resp, err := c.client.GetFileByID(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("gRPC error: %v", st.Message())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	file := &domain.File{
		ID:          uint(resp.Id),
		DirectoryID: uint(resp.DirectoryId),
		Name:        resp.Name,
		Status:      resp.Status,
		Version:     int(resp.Version),
		Directory: &domain.Directory{
			ID:           uint(resp.Directory.Id),
			ParentPathID: uintPtrOrNil(resp.Directory.ParentPathId),
			Name:         resp.Directory.Name,
			Status:       resp.Directory.Status,
			WorkflowID:   uint(resp.Directory.WorkflowId),
		},
	}
	return file, nil
}

func (c *FileGRPCClient) UpdateFileStatus(ctx context.Context, fileID uint, fileStatus string) error {
	const op = "infrastructure.grpc.fileclient.UpdateFileStatus"

	req := &pb.UpdateFileStatusRequest{
		FileId: uint32(fileID),
		Status: fileStatus,
	}

	_, err := c.client.UpdateFileStatus(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("gRPC error: %v", st.Message())
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *FileGRPCClient) GetFilesInfo(ctx context.Context, fileIDs []uint32) (map[uint32]string, error) {
	const op = "infrastructure.grpc.fileclient.GetFilesInfo"

	req := &pb.GetFilesRequest{
		FileIds: fileIDs,
	}
	resp, err := c.client.GetFilesInfo(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return resp.FileNames, nil
}

func (c *FileGRPCClient) CheckWorkflow(ctx context.Context, workflowID uint) (bool, error) {
	const op = "infrastructure.grpc.fileclient.CheckWorkflow"

	req := &pb.CheckWorkflowRequest{
		WorkflowId: uint32(workflowID),
	}

	resp, err := c.client.CheckWorkflow(ctx, req)
	if err != nil {
		return false, nil
	}

	return resp.Exists, nil
}

func uintPtrOrNil(value uint32) *uint {
	if value == 0 {
		return nil
	}
	v := uint(value)
	return &v
}
