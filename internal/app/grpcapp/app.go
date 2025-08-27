package grpcapp

import (
	"auth-service/internal/grpc/authgrpc"
	"auth-service/internal/service"
	"context"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
	listener   net.Listener
}

func New(log *slog.Logger, addr string, authSvc *service.Auth) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authSvc) // <- передаём готовый экземпляр Auth

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       addr,
	}
}

// Запуск сервера и сохранение listener
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	a.listener = l

	log.Info("grpc server is running", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Shutdown реализует интерфейс Stoppable для graceful shutdown
func (a *App) Shutdown(ctx context.Context) error {
	const op = "grpcapp.Shutdown"
	a.log.With(slog.String("op", op)).Info("starting graceful shutdown", slog.String("port", a.port))

	done := make(chan struct{})
	go func() {
		a.gRPCServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		a.log.With(slog.String("op", op)).Warn("timeout exceeded, forcing stop")
		a.gRPCServer.Stop() // жёсткая остановка, если таймаут
	case <-done:
		a.log.With(slog.String("op", op)).Info("graceful shutdown complete")
	}
	return nil
}

// Остановка сервера
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.String("port", a.port))

	a.gRPCServer.GracefulStop()
}
