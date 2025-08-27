package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Stoppable interface {
	Shutdown(ctx context.Context) error
	Stop()
}

func WaitForSignals(timeout time.Duration, services ...Stoppable) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, srv := range services {
		if err := srv.Shutdown(shutdownCtx); err != nil {
			srv.Stop()
		}
	}
}
