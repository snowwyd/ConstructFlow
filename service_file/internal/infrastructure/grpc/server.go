package filegrpc

import (
	"context"
	"errors"
	"service-file/internal/domain"
	"service-file/internal/domain/interfaces"
	pb "service-file/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedFileServiceServer
	usecase interfaces.GRPCUsecase
}

func NewGRPCServer(gRPC *grpc.Server, usecase interfaces.GRPCUsecase) {
	pb.RegisterFileServiceServer(gRPC, &GRPCServer{usecase: usecase})
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

func (s *GRPCServer) UpdateFileStatus(ctx context.Context, req *pb.UpdateFileStatusRequest) (*emptypb.Empty, error) {
	err := s.usecase.UpdateFileStatus(ctx, uint(req.GetFileId()), req.GetStatus())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFileNotFound):
			return &emptypb.Empty{}, status.Error(codes.NotFound, "file not found")
		default:
			return &emptypb.Empty{}, status.Error(codes.Internal, "internal error")
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) GetFilesInfo(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	fileIDs := req.FileIds
	fileNames := make(map[uint32]string)

	files, err := s.usecase.GetFilesByID(ctx, fileIDs)
	if err != nil {
		// TODO: custom errors
		return nil, status.Errorf(codes.Internal, "failed to get files: %v", err)
	}

	for _, file := range files {
		fileNames[uint32(file.ID)] = file.Name
	}

	return &pb.GetFilesResponse{
		FileNames: fileNames,
	}, nil
}

func (s *GRPCServer) CheckWorkflow(ctx context.Context, req *pb.CheckWorkflowRequest) (*pb.CheckWorkflowResponse, error) {
	exists, err := s.usecase.CheckWorkflow(ctx, uint(req.WorkflowId))
	if err != nil {
		// TODO: custom errors
		return nil, status.Errorf(codes.Internal, "failed to check workflow id")
	}

	return &pb.CheckWorkflowResponse{Exists: exists}, nil
}

func (s *GRPCServer) DeleteUserRelations(ctx context.Context, req *pb.DeleteUserRelationsRequest) (*emptypb.Empty, error) {
	err := s.usecase.DeleteUserRelations(ctx, uint(req.GetUserId()))
	if err != nil {
		// TODO: custom errors
		return nil, status.Errorf(codes.Internal, "failed to delete user relations")
	}

	return &emptypb.Empty{}, nil
}
