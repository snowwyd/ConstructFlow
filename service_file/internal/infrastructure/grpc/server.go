package grpc

import (
	"context"
	"errors"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	pb "service-file/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedFileServiceServer
	usecase interfaces.FileTreeUsecase
}

func NewGRPCServer(usecase interfaces.FileTreeUsecase) *GRPCServer {
	return &GRPCServer{
		usecase: usecase,
	}
}

func (s *GRPCServer) GetFileByID(ctx context.Context, req *pb.GetFileRequest) (*pb.FileResponse, error) {
	file, err := s.usecase.GetFileByID(ctx, uint(req.GetFileId()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			return nil, status.Error(codes.NotFound, "file not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &pb.FileResponse{
		Id:          uint32(file.ID),
		DirectoryId: uint32(file.DirectoryID),
		Name:        file.Name,
		Status:      file.Status,
		Version:     int32(file.Version),
		Directory: &pb.DirectoryResponse{
			Id:           uint32(file.Directory.ID),
			ParentPathId: uint32PtrOrNil(file.Directory.ParentPathID),
			Name:         file.Directory.Name,
			Status:       file.Directory.Status,
			WorkflowId:   uint32(file.Directory.WorkflowID),
		},
	}, nil
}

func uint32PtrOrNil(value *uint) uint32 {
	if value == nil {
		return 0
	}
	return uint32(*value)
}
