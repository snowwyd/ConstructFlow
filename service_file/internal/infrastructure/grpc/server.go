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
	fileInfo, err := s.usecase.GetFileByID(ctx, uint(req.GetFileId()))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			return nil, status.Error(codes.NotFound, "file not found")
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &pb.FileResponse{
		Id:          uint64(fileInfo.ID),
		Name:        fileInfo.NameFile,
		Status:      fileInfo.Status,
		DirectoryId: uint64(fileInfo.DirectoryID),
	}, nil
}
