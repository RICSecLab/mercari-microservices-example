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
  app "github.com/mercari/mercari-microservices-example/services/item/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/item/proto"
)


func TestItem(t *testing.T) {
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
  clogger := l.WithName("item")
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
    client := app_proto.NewItemServiceClient( connection )
    {
      response, e4 := client.CreateItem( context.Background(), &app_proto.CreateItemRequest{ CustomerId : "1", Title : "fuga", Price : 1234  } )
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

func FuzzItem(f *testing.F) {
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
    client := app_proto.NewItemServiceClient( connection )
    f.Fuzz( func( t *testing.T, customer_id int, item_id int, name_ []byte, title_ []byte, price int64 ) {
      name := test_utils.ToValidUTF8StringBiased( name_, 30 )
      title := test_utils.ToValidUTF8StringBiased( title_, 30 )
      fake_server.CustomerId = customer_id
      fake_server.ItemId = item_id
      fake_server.Name = name
      fake_server.Title = title
      fake_server.Price = price
      customer_id_in_str := strconv.Itoa( customer_id )
      item_id_in_str := strconv.Itoa( item_id )
      {
        response, e4 := client.CreateItem( context.Background(), &app_proto.CreateItemRequest{ CustomerId : customer_id_in_str, Title : title, Price : price  } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Item.Id != item_id_in_str {
          panic( "unexpected item id" )
        }
        if response.Item.CustomerId != customer_id_in_str {
          panic( "unexpected customer id" )
        }
        if response.Item.Title != title {
          panic( "unexpected name" )
        }
        if response.Item.Price != price {
          panic( "unexpected price" )
        }
      }
      {
        response, e4 := client.GetItem( context.Background(), &app_proto.GetItemRequest{ Id: item_id_in_str } )
        if e4 != nil {
          panic( e4 )
        }
        if response.Item.Id != item_id_in_str {
          panic( "unexpected item id" )
        }
        if response.Item.CustomerId != customer_id_in_str {
          panic( "unexpected customer id" )
        }
        if response.Item.Title != title {
          panic( "unexpected name" )
        }
        if response.Item.Price != price {
          panic( "unexpected price" )
        }
      }
      {
	response, e4 := client.ListItems( context.Background(), &app_proto.ListItemsRequest{} )
        if e4 != nil {
          panic( e4 )
        }
        if len( response.Items ) != 1 {
          panic( "unexpected response size" )
        }
        if response.Items[ 0 ].Id != item_id_in_str {
          panic( "unexpected item id" )
        }
        if response.Items[ 0 ].CustomerId != customer_id_in_str {
          panic( "unexpected customer id" )
        }
        if response.Items[ 0 ].Title != title {
          panic( "unexpected name" )
        }
        if response.Items[ 0 ].Price != price {
          panic( "unexpected price" )
        }
      }
    })
  }
}


