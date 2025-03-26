package grpc

import (
	"context"
	"service-file/internal/domain/interfaces"
	pb "service-file/internal/proto"
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
		return nil, err
	}

	return &pb.FileResponse{
		Id:     uint32(fileInfo.ID),
		Name:   fileInfo.NameFile,
		Status: fileInfo.Status,
	}, nil
}
