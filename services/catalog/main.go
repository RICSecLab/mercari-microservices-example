package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"

	"github.com/mercari/mercari-microservices-example/pkg/logger"
	"github.com/mercari/mercari-microservices-example/services/catalog/grpc"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	ctx, stop := signal.NotifyContext(ctx, unix.SIGTERM, unix.SIGINT)
	defer stop()

	l, err := logger.New()
	if err != nil {
		_, ferr := fmt.Fprintf(os.Stderr, "failed to create logger: %s", err)
		if ferr != nil {
			// Unhandleable, something went wrong...
			panic(fmt.Sprintf("failed to write log:`%s` original error is:`%s`", ferr, err))
		}
		return 1
	}
	clogger := l.WithName("catalog")
        runningAt := ""
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpc.RunServer(ctx, 5000, clogger.WithName("grpc"),"customer.customer.svc.cluster.local:5000","item.item.svc.cluster.local:5000",&runningAt)
	}()

	select {
	case err := <-errCh:
		fmt.Println(err.Error())
		return 1
	case <-ctx.Done():
		fmt.Println("shutting down...")
		return 0
	}
}
