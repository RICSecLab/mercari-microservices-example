package main_test

import (
  "testing"
  "context"
  "runtime"
//  "strconv"
  "net"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "google.golang.org/grpc/metadata"
  "github.com/mercari/mercari-microservices-example/pkg/logger"
  item "github.com/mercari/mercari-microservices-example/services/item/proto"
  customer "github.com/mercari/mercari-microservices-example/services/customer/proto"
  test_utils "github.com/mercari/mercari-microservices-example/platform/test_utils"
  app "github.com/mercari/mercari-microservices-example/services/catalog/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/catalog/proto"
)

func TestCatalog(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  customer.RegisterCustomerServiceServer( server, &test_utils.FakeCustomerServiceServer{
    CustomerId : 1,
    Name : "hoge" } )
  item.RegisterItemServiceServer( server, &test_utils.FakeItemServiceServer{
    CustomerId : 1,
    ItemId : 2,
    Title : "fuga",
    Price : 1234 } )
  go func() {
    if e := server.Serve( socket ); e != nil {
      panic( e )
    }
    socket.Close()
  }()
  l, e2 := logger.New()
  if e2 != nil {
    panic( e2 )
  }
  clogger := l.WithName("catalog")
  runningAt := ""
  go func() {
    app.RunServer( context.Background(), 0, clogger, socket.Addr().String(), socket.Addr().String(), &runningAt )
  }()
  for runningAt == "" {
    runtime.Gosched()
  }
  {
    connection, e3 := grpc.Dial(
      runningAt,
      grpc.WithTransportCredentials( insecure.NewCredentials() ),
      grpc.WithBlock() )
    if e3 != nil {
      panic( e3 )
    }
    defer connection.Close()
    md := metadata.New(map[string]string{"authorization": "bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImFhN2M2Mjg3LWM0NWQtNDk2Ni04NGI0LWExNjMzZTRlM2E2NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRob3JpdHkiLCJzdWIiOiIxIn0.VwwG12uOOBpWRvrTyKwa8ZQXpF8-xCxIQ9loeXc27Q89qBitbVTDYDAntGxcYEBb4J36EjOMgqNXbXZPXH0ITrMxP6wU7ofGP3E9zIaN6xibefTlmNF-fFB5QNivk8mI2tyWahLDCGRbpy_dAxyvDzZpzqdbvDKdw02mXz6AWXZBqAMlqC9C3N28JglJ8B-udSGalPD5UHxUC-kBF9A0TZ7tW54-Jcjjh5-6fiX7AYtaQzV31AI82XNNA4V1pc0ZLpEVYHQmeGYONNRpuAs83LNLb8JKBd_-SpqAATX7XTAZnIUfoZWnV0xuePjwjtcMroWMvJeYs5Elo0KxDnLERw"})
    ctx := metadata.NewOutgoingContext( context.Background(), md )
    client := app_proto.NewCatalogServiceClient( connection )
    {
      response, e4 := client.CreateItem( ctx, &app_proto.CreateItemRequest{ Title : "fuga", Price : 1234  } )
      if e4 != nil {
        panic( e4 )
      }
      if response.Item.Id != "2" {
        panic( "unexpected id" )
      }
      if response.Item.CustomerId != "1" {
        panic( "unexpected customer id" )
      }
      if response.Item.Title != "fuga" {
        panic( "unexpected name" )
      }
      if response.Item.Price != 1234 {
        panic( "unexpected price" )
      }
    }
  }
}


