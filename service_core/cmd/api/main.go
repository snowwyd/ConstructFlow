package main

import (
	"log/slog"
	"os"
	"os/signal"
	"service-core/internal/app"
	"service-core/pkg/config"
	"service-core/pkg/logger"
	"syscall"

	_ "service-core/docs"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger()

	log.Info("config loaded successfully",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.Address),
		slog.Any("http_server", cfg.HTTPServer),
	)

	application, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to initialize application", slog.String("error", err.Error()))
		os.Exit(1)
	}

	go application.HTTPSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.HTTPSrv.Stop()

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
