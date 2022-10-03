package main_test

import (
  "testing"
  "context"
  "runtime"
  "unicode/utf8"
  "net"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "github.com/mercari/mercari-microservices-example/pkg/logger"
  db "github.com/mercari/mercari-microservices-example/platform/db/proto"
  app "github.com/mercari/mercari-microservices-example/services/customer/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/customer/proto"
)

type FakeDBServiceServer struct {
  db.UnimplementedDBServiceServer
}

func (FakeDBServiceServer) CreateCustomer( ctx context.Context, request *db.CreateCustomerRequest) (*db.CreateCustomerResponse, error) {
  return &db.CreateCustomerResponse{ Customer : &db.Customer{ Id: "1", Name: request.Name } }, nil
}
func (FakeDBServiceServer) GetCustomer( ctx context.Context, request *db.GetCustomerRequest) (*db.GetCustomerResponse, error) {
  return &db.GetCustomerResponse{ Customer : &db.Customer{ Id: request.Id, Name: "hoge" } }, nil
}
func (FakeDBServiceServer) GetCustomerByName( ctx context.Context, request *db.GetCustomerByNameRequest) (*db.GetCustomerByNameResponse, error) {
  return &db.GetCustomerByNameResponse{ Customer : &db.Customer{ Id: "1", Name: request.Name } }, nil
}
func (FakeDBServiceServer) CreateItem( ctx context.Context, request *db.CreateItemRequest) (*db.CreateItemResponse, error) {
  return &db.CreateItemResponse{ Item : &db.Item{ Id: "2", CustomerId: request.CustomerId, Title: request.Title, Price: request.Price } }, nil
}
func (FakeDBServiceServer) GetItem( ctx context.Context, request *db.GetItemRequest) (*db.GetItemResponse, error) {
  return &db.GetItemResponse{ Item : &db.Item{ Id: request.Id, CustomerId: "1", Title: "fuga", Price: 1234 } }, nil
}
func (FakeDBServiceServer) ListItems(context.Context, *db.ListItemsRequest) (*db.ListItemsResponse, error) {
  return &db.ListItemsResponse{ Items : []*db.Item{ &db.Item{ Id: "2", CustomerId: "1", Title: "fuga", Price: 1234 } } }, nil
}

func TestCustomer(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  db.RegisterDBServiceServer( server, &FakeDBServiceServer{} )
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
  clogger := l.WithName("customer")
  runningAt := ""
  go func() {
    app.RunServer( context.Background(), 0, clogger, socket.Addr().String(), &runningAt )
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
    client := app_proto.NewCustomerServiceClient( connection )
    {
      response, e4 := client.CreateCustomer( context.Background(), &app_proto.CreateCustomerRequest{ Name: "piyo" } )
      if e4 != nil {
        panic( e4 )
      }
      if response.Customer.Id != "1" {
        panic( "unexpected id" )
      }
      if response.Customer.Name != "piyo" {
        panic( "unexpected name" )
      }
    }
  }
}
func FuzzCreateCustomer(f *testing.F) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  db.RegisterDBServiceServer( server, &FakeDBServiceServer{} )
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
  clogger := l.WithName("customer")
  runningAt := ""
  go func() {
    app.RunServer( context.Background(), 0, clogger, socket.Addr().String(), &runningAt )
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
    client := app_proto.NewCustomerServiceClient( connection )
    f.Add( "" )
    f.Fuzz( func( t *testing.T, name string ) {
      if utf8.ValidString( name ) {
        response, e4 := client.CreateCustomer( context.Background(), &app_proto.CreateCustomerRequest{ Name: name } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Customer.Id != "1" {
          panic( "unexpected id" )
        }
        if response.Customer.Name != name {
          panic( "unexpected name" )
        }
      }
    })
  }
}


