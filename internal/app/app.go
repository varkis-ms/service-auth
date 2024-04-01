package app

import (
	"os"
	"os/signal"
	"syscall"

	grpcapp "github.com/varkis-ms/service-auth/internal/app/grpc"
	"github.com/varkis-ms/service-auth/internal/config"
	"github.com/varkis-ms/service-auth/internal/rpc/login"
	"github.com/varkis-ms/service-auth/internal/rpc/signup"
	"github.com/varkis-ms/service-auth/internal/rpc/validate"
	"github.com/varkis-ms/service-auth/internal/storage"
	"github.com/varkis-ms/service-auth/pkg/database/postgresdb"
	"github.com/varkis-ms/service-auth/pkg/logging"
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
