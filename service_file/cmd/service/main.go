package main

import (
	"log/slog"
	"os"
	"os/signal"
	"service-file/internal/app"
	"service-file/pkg/config"
	"service-file/pkg/logger"
	"syscall"

	_ "service-file/docs"
)

//	@title			Constructflow File Microservice
//	@version		0.2
//	@description	API File Microservice for Constructflow

//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization

func main() {
	cfg := config.MustLoadEnv()

	log := setupLogger()

	log.Info("config loaded successfully",
		slog.String("env", cfg.Env),
		slog.Any("http_server", cfg.HTTPServer),
		slog.Any("grpc_server", cfg.GRPCServer),
	)

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go application.HTTPSrv.MustRun()
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.HTTPSrv.Stop()
	application.GRPCSrv.Stop()

	log.Info("application stopped")
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
