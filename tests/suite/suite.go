package suite

import (
	"auth-service/config"
	"context"
	"net"
	"testing"
	"time"

	"github.com/ILmira-116/protos/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient auth.AuthClient // grpc клиент
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// Загружаем конфиг из ENV
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("cannot load config: %v", err)
	}

	// Родительский контекст с таймаутом (для всех RPC-вызовов в тестах)
	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*10)

	// Отмена контекста после завершения всех тестов
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	// Адрес gRPC-сервера (host:port)
	grpcAddress := net.JoinHostPort(cfg.GRPC.ServerHost, cfg.GRPC.ServerPort)

	// Создаём gRPC-клиента
	cc, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	// Закрываем соединение после завершения тестов
	t.Cleanup(func() {
		_ = cc.Close()
	})

	// gRPC-клиент  Auth-сервиса
	authClient := auth.NewAuthClient(cc)

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: authClient,
	}

}
