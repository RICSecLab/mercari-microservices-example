package main_test

import (
  "testing"
  "context"
  "runtime"
  "strconv"
  "net"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "github.com/mercari/mercari-microservices-example/pkg/logger"
  db "github.com/mercari/mercari-microservices-example/platform/db/proto"
  test_utils "github.com/mercari/mercari-microservices-example/platform/test_utils"
  app "github.com/mercari/mercari-microservices-example/services/customer/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/customer/proto"
)

func TestCustomer(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  db.RegisterDBServiceServer( server, &test_utils.FakeDBServiceServer{
    CustomerId : 1,
    ItemId : 2,
    Name : "hoge",
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
  fake_server := &test_utils.FakeDBServiceServer{
    CustomerId : 1,
    ItemId : 2,
    Name : "hoge",
    Title : "fuga",
    Price : 1234 }
  db.RegisterDBServiceServer( server, fake_server )
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
    f.Fuzz( func( t *testing.T, customer_id int, item_id int, name_ []byte, title_ []byte, price int64 ) {
      name := test_utils.ToValidUTF8StringBiased( name_, 30 )
      title := test_utils.ToValidUTF8StringBiased( title_, 30 )
      fake_server.CustomerId = customer_id
      fake_server.ItemId = item_id
      fake_server.Name = name
      fake_server.Title = title
      fake_server.Price = price
      customer_id_in_str := strconv.Itoa( customer_id )
      {
        response, e4 := client.CreateCustomer( context.Background(), &app_proto.CreateCustomerRequest{ Name: name } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Customer.Id != customer_id_in_str {
          panic( "unexpected id" )
        }
        if response.Customer.Name != name {
          panic( "unexpected name" )
        }
      }
      {
        response, e4 := client.GetCustomer( context.Background(), &app_proto.GetCustomerRequest{ Id: customer_id_in_str } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Customer.Id != customer_id_in_str {
          panic( "unexpected id" )
        }
        if response.Customer.Name != name {
          panic( "unexpected name" )
        }
      }
      {
	response, e4 := client.GetCustomerByName( context.Background(), &app_proto.GetCustomerByNameRequest{ Name: name } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Customer.Id != customer_id_in_str {
          panic( "unexpected id" )
        }
        if response.Customer.Name != name {
          panic( "unexpected name" )
        }
      }
    })
  }
}


