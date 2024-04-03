package app

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/varkis-ms/service-auth/internal/app/grpc"
	"github.com/varkis-ms/service-auth/internal/config"
	"github.com/varkis-ms/service-auth/internal/pkg/database/postgresdb"
	"github.com/varkis-ms/service-auth/internal/pkg/logger/handlers/slogpretty"
	"github.com/varkis-ms/service-auth/internal/pkg/logger/sl"
	"github.com/varkis-ms/service-auth/internal/rpc/login"
	"github.com/varkis-ms/service-auth/internal/rpc/signup"
	"github.com/varkis-ms/service-auth/internal/rpc/validate"
	"github.com/varkis-ms/service-auth/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type App struct {
	GRPCServer *grpcapp.App
}

func Run(configPath string) {
	// Config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic("config.LoadConfig failed" + err.Error())
	}

	logger := setupLogger(cfg.Env)

	// Repository
	logger.Info("Initializing postgres...")
	db, err := postgresdb.New(&cfg)
	if err != nil {
		logger.Error("postgresdb.New failed", sl.Err(err))
	}
	defer db.Close()
	repositories := storage.New(db)

	// Handlers
	loginHandler := login.New(repositories, cfg.SecretKey, logger)
	signupHandler := signup.New(repositories, logger)
	validateHandler := validate.New(repositories, cfg.SecretKey, logger)

	// gRPC server
	app := grpcapp.New(
		logger,
		cfg.Port,
		loginHandler,
		signupHandler,
		validateHandler,
	)

	go func() {
		app.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Shutdown()
	logger.Info("Gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
