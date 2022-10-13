package http

import (
	"net"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	authoritypb "github.com/mercari/mercari-microservices-example/services/authority/proto"
	catalogpb "github.com/mercari/mercari-microservices-example/services/catalog/proto"
)

func RunServer(ctx context.Context, port int, authorityAddr string, catalogAddr string, runningAt *string ) error {
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: false,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
	}

	authorityConn, err := grpc.DialContext(ctx, authorityAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial to authority grpc server: %w", err)
	}
	if err = authoritypb.RegisterAuthorityServiceHandlerClient(ctx, mux, authoritypb.NewAuthorityServiceClient(authorityConn)); err != nil {
		return fmt.Errorf("failed to create a authority grpc client: %w", err)
	}

	catalogConn, err := grpc.DialContext(ctx, catalogAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial to catalog grpc server: %w", err)
	}
	if err := catalogpb.RegisterCatalogServiceHandlerClient(ctx, mux, catalogpb.NewCatalogServiceClient(catalogConn)); err != nil {
		return fmt.Errorf("failed to create a catalog grpc client: %w", err)
	}
	


	errCh := make(chan error, 1)
	go func() {
		socket, e := net.Listen("tcp", fmt.Sprintf(":%d", port) )
		if e != nil {
			panic( e )
		}
		*runningAt = socket.Addr().String()
		errCh <- http.Serve( socket, mux )
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failed to serve http server: %w", err)
	case <-ctx.Done():

		if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("failed to close http server: %w", err)
		}

		return nil
	}
}
