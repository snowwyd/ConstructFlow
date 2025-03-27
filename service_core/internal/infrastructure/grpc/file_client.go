package grpc

import (
	"context"
	"log"
	"os"
	"service-core/internal/domain"

	pb "service-core/internal/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type FileGRPCClient struct {
	client pb.FileServiceClient
}

func NewFileGRPCClient(address string) (*FileGRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(addAPIKey)) // Используйте TLS в production
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServiceClient(conn)
	return &FileGRPCClient{client: client}, nil
}

func addAPIKey(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	godotenv.Load()
	apiKey := os.Getenv("GRPC_KEY")
	ctx = metadata.AppendToOutgoingContext(ctx, "x-api-key", apiKey)
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (c *FileGRPCClient) GetFileWithDirectory(ctx context.Context, fileID uint) (*domain.File, error) {
	req := &pb.GetFileRequest{
		FileId: uint32(fileID),
	}

	resp, err := c.client.GetFileByID(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("gRPC error: %v", st.Message())
		return nil, err
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
	req := &pb.UpdateFileStatusRequest{
		FileId: uint32(fileID),
		Status: fileStatus,
	}

	_, err := c.client.UpdateFileStatus(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		log.Printf("gRPC error: %v", st.Message())
		return err
	}

	return nil
}

func (c *FileGRPCClient) GetFilesInfo(ctx context.Context, fileIDs []uint32) (map[uint32]string, error) {
	req := &pb.GetFilesRequest{
		FileIds: fileIDs,
	}
	resp, err := c.client.GetFilesInfo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.FileNames, nil
}

func uintPtrOrNil(value uint32) *uint {
	if value == 0 {
		return nil
	}
	v := uint(value)
	return &v
}
