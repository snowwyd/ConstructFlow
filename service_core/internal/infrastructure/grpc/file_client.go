package grpc

import (
	"context"
	"log"
	"service-core/internal/domain"

	pb "service-core/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type FileGRPCClient struct {
	client pb.FileServiceClient
}

func NewFileGRPCClient(address string) (*FileGRPCClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Используйте TLS в production
	if err != nil {
		return nil, err
	}
	client := pb.NewFileServiceClient(conn)
	return &FileGRPCClient{client: client}, nil
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
		DirectoryID: uint(resp.DirectoryId),
		Name:        resp.Name,
		Status:      resp.Status,
		Version:     int(resp.Version),
		Directory: &domain.Directory{
			ParentPathID: uintPtrOrNil(resp.Directory.ParentPathId),
			Name:         resp.Directory.Name,
			Status:       resp.Directory.Status,
			WorkflowID:   uint(resp.Directory.WorkflowId),
		},
	}
	file.ID = uint(resp.Id)
	file.Directory.ID = uint(resp.Directory.Id)
	return file, nil
}

func uintPtrOrNil(value uint32) *uint {
	if value == 0 {
		return nil
	}
	v := uint(value)
	return &v
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
