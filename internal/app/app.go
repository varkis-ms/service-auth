package app

import (
	"os"
	"os/signal"
	grpcapp "service-auth/internal/app/grpc"
	"service-auth/internal/config"
	"service-auth/internal/rpc/login"
	"service-auth/internal/rpc/signup"
	"service-auth/internal/rpc/validate"
	"service-auth/internal/storage"
	"service-auth/pkg/database/postgresdb"
	"service-auth/pkg/logging"
	"syscall"
)

type App struct {
	GRPCServer *grpcapp.App
}

func Run(configPath string) {
	// Config
	logger := logging.GetLogger()
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.WithError(err).Fatal("no config")
	}

	// Repository
	logger.Info("Initializing postgres...")
	db, err := postgresdb.New(&cfg)
	if err != nil {
		logger.WithError(err).Fatal("app.Run - postgresdb.New")
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
