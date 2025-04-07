package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"service-file/internal/domain/interfaces"

	filegrpc "service-file/internal/infrastructure/grpc"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *slog.Logger, gRPCUsecase interfaces.GRPCUsecase, port string) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)

	filegrpc.NewGRPCServer(gRPCServer, gRPCUsecase)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run запускает приложение
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("operation", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", a.port)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is running", slog.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("operation", op)).Info("grpc server is stopping")
	a.gRPCServer.GracefulStop()
	a.log.Info("grpc server stopped gracefully")
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	godotenv.Load()
	expectedAPIKey := os.Getenv("GRPC_KEY")

	apiKeys := md.Get("x-api-key")
	fmt.Println(apiKeys)
	if len(apiKeys) == 0 || apiKeys[0] != expectedAPIKey {
		return nil, status.Error(codes.Unauthenticated, "invalid API key")
	}

	return handler(ctx, req)
}
