package grpc

import (
	"fmt"
	"net"

	authgrpc "github.com/varkis-ms/service-auth/internal/grpc/auth"
	"github.com/varkis-ms/service-auth/internal/rpc/login"
	"github.com/varkis-ms/service-auth/internal/rpc/signup"
	"github.com/varkis-ms/service-auth/internal/rpc/validate"

	"github.com/varkis-ms/service-auth/pkg/logging"

	"google.golang.org/grpc"
)

type App struct {
	log        *logging.Logger
	gRPCServer *grpc.Server
	port       int64
}

// New creates new gRPC server app.
func New(
	log *logging.Logger,
	port int64,
	loginHandler *login.Handler,
	signupHandler *signup.Handler,
	validateHandler *validate.Handler,
) *App {
	//loggingOpts := []logging.Option{
	//	logging.WithLogOnEvents(
	//		//logging.StartCall, logging.FinishCall,
	//		logging.PayloadReceived, logging.PayloadSent,
	//	),
	//	// Add any other option (check functions starting with logging.With).
	//}

	//recoveryOpts := []recovery.Option{
	//	recovery.WithRecoveryHandler(func(p interface{}) (err error) {
	//		log.Error("Recovered from panic", log.WithField("panic", p))
	//
	//		return status.Errorf(codes.Internal, "internal error")
	//	}),
	//}

	gRPCServer := grpc.NewServer(
	//grpc.ChainUnaryInterceptor(
	//recovery.UnaryServerInterceptor(recoveryOpts...),
	//logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	//)
	)

	authgrpc.Register(gRPCServer, loginHandler, signupHandler, validateHandler)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", a.log.WithFields(
		map[string]interface{}{
			"addr": l.Addr().String(),
			"op":   op,
		}),
	)

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Shutdown stops gRPC server.
func (a *App) Shutdown() {
	const op = "grpcapp.Shutdown"

	a.log.Info("stopping gRPC server", a.log.WithFields(
		map[string]interface{}{
			"port": a.port,
			"op":   op,
		}),
	)

	a.gRPCServer.GracefulStop()
}
