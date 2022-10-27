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
  customer "github.com/mercari/mercari-microservices-example/services/customer/proto"
  test_utils "github.com/mercari/mercari-microservices-example/platform/test_utils"
  app "github.com/mercari/mercari-microservices-example/services/authority/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/authority/proto"
)

func TestAuthority(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  customer.RegisterCustomerServiceServer( server, &test_utils.FakeCustomerServiceServer{
    CustomerId : 1,
    Name : "hoge" } )
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
  clogger := l.WithName("authority")
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
    client := app_proto.NewAuthorityServiceClient( connection )
    {
      response, e4 := client.Signup( context.Background(), &app_proto.SignupRequest{ Name: "hoge" } )
      if e4 != nil {
        panic( e4 )
      }
      if response.Customer.Id != "1" {
        panic( "unexpected id" )
      }
      if response.Customer.Name != "hoge" {
        panic( "unexpected name" )
      }
    }
    {
      response, e4 := client.Signin( context.Background(), &app_proto.SigninRequest{ Name: "hoge" } )
      if e4 != nil {
        panic( e4 )
      }
      if len( response.AccessToken ) == 0 {
        panic( "invalid access token" )
      }
    }
  }
}

func FuzzAuthority(f *testing.F) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  fake_server := &test_utils.FakeCustomerServiceServer{
    CustomerId : 1,
    Name : "hoge" }
  customer.RegisterCustomerServiceServer( server, fake_server )
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
    client := app_proto.NewAuthorityServiceClient( connection )
    f.Fuzz( func( t *testing.T, customer_id int, name_ []byte ) {
      name := test_utils.ToValidUTF8StringBiased( name_, 65536*65536 )
      fake_server.CustomerId = customer_id
      fake_server.Name = name
      customer_id_in_str := strconv.Itoa( customer_id )
      {
        response, e4 := client.Signup( context.Background(), &app_proto.SignupRequest{ Name: name } )
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
        response, e4 := client.Signin( context.Background(), &app_proto.SigninRequest{ Name: name } )
        if e4 != nil {
          panic( e4 )
        }
        if len( response.AccessToken ) == 0 {
          panic( "invalid access token" )
        }
      }
    })
  }
}

