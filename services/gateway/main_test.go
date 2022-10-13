package main_test

import (
  "net"
  "testing"
  "runtime"
  "context"
  "google.golang.org/grpc"
  authority "github.com/mercari/mercari-microservices-example/services/authority/proto"
  catalog "github.com/mercari/mercari-microservices-example/services/catalog/proto"
  test_utils "github.com/mercari/mercari-microservices-example/platform/test_utils"
  http "github.com/mercari/mercari-microservices-example/services/gateway/http"
)

func TestGateway(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  fake_authority_server := &test_utils.FakeAuthorityServiceServer{
    CustomerId : 1}
  authority.RegisterAuthorityServiceServer( server, fake_authority_server )
  fake_catalog_server := &test_utils.FakeCatalogServiceServer{
    CustomerId : 1,
    ItemId : 2,
    Title : "fuga",
    Price : 1234 }
  catalog.RegisterCatalogServiceServer( server, fake_catalog_server )
  ctx := context.Background()
  go func() {
    if e := server.Serve( socket ); e != nil {
      panic( e )
    }
    socket.Close()
  }()

  httpErrCh := make(chan error, 1)
  runningAt := ""
  go func() {
    httpErrCh <- http.RunServer(ctx, 0, socket.Addr().String(), socket.Addr().String(), &runningAt )
  }()
  for runningAt == "" {
    runtime.Gosched()
  }
}
