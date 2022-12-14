package grpc

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"google.golang.org/grpc"

	pkggrpc "github.com/mercari/mercari-microservices-example/pkg/grpc"
	"github.com/mercari/mercari-microservices-example/services/catalog/proto"
	customer "github.com/mercari/mercari-microservices-example/services/customer/proto"
	item "github.com/mercari/mercari-microservices-example/services/item/proto"
)

func RunServer(ctx context.Context, port int, logger logr.Logger, customerAddr string, itemAddr string, runningAt *string) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	cconn, err := grpc.DialContext(ctx, customerAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial customer grpc server: %w", err)
	}

	iconn, err := grpc.DialContext(ctx, itemAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial item grpc server: %w", err)
	}

	svc := &server{
		customerClient: customer.NewCustomerServiceClient(cconn),
		itemClient:     item.NewItemServiceClient(iconn),
	}

	return pkggrpc.NewServer(port, logger, func(s *grpc.Server) {
		proto.RegisterCatalogServiceServer(s, svc)
	}).Start(ctx,runningAt)
}
