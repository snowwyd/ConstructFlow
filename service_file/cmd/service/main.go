package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"service-file/internal/app"
	http "service-file/internal/controller/middleware"
	pb "service-file/internal/proto"
	"service-file/pkg/config"
	"service-file/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("config loaded successfully",
		slog.String("env", cfg.Env),
		slog.Any("http_server", cfg.HTTPServer),
		slog.Any("grpc_server", cfg.GRPCServer),
	)

	appInstance, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go func() {
		router := gin.Default()
		router.Use(gin.Recovery(), http.CORSMiddleware())

		setupRoutes(router, appInstance, cfg)

		port := cfg.HTTPServer.Address
		log.Info("starting HTTP server", slog.String("port", port))
		if err := router.Run(port); err != nil {
			log.Error("HTTP server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Запуск gRPC-сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("failed to listen for gRPC", slog.String("error", err.Error()))
		os.Exit(1)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor),
	)
	pb.RegisterFileServiceServer(grpcServer, appInstance.GRPCServer)

	log.Info("starting gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("gRPC server failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func setupLogger() *slog.Logger {
	opts := logger.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}

func setupRoutes(router *gin.Engine, appInstance *app.App, cfg *config.Config) {
	directoriesGroup := router.Group("/directories", http.AuthMiddleware(cfg))
	{
		directoriesGroup.POST("/create", appInstance.TreeHandler.CreateDirectory)
		directoriesGroup.DELETE("", appInstance.TreeHandler.DeleteDirectory)
		directoriesGroup.POST("", appInstance.TreeHandler.GetTree)
	}

	filesGroup := router.Group("/files", http.AuthMiddleware(cfg))
	{
		filesGroup.POST("/upload", appInstance.TreeHandler.UploadFile)
		filesGroup.DELETE("", appInstance.TreeHandler.DeleteFile)
		filesGroup.GET("/:file_id", appInstance.TreeHandler.GetFileInfo)
	}
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	godotenv.Load()
	expectedAPIKey := os.Getenv("GRPC_KEY")

	apiKeys := md.Get("x-api-key")
	if len(apiKeys) == 0 || apiKeys[0] != expectedAPIKey {
		return nil, status.Error(codes.Unauthenticated, "invalid API key")
	}

	return handler(ctx, req)
}
