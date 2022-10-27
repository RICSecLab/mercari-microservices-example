package main_test

import (
  "fmt"
//  "bytes"
  "testing"
  "context"
  "runtime"
  "strconv"
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
    Price : 1234,
    Count : 3 } )
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
    md := metadata.New(map[string]string{ "authorization": "bearer "+test_utils.GetAccessToken( 1 ) })
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

func FuzzCatalog(f *testing.F) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  fake_customer_server := &test_utils.FakeCustomerServiceServer{
    CustomerId : 1,
    Name : "hoge" }
  customer.RegisterCustomerServiceServer( server, fake_customer_server )
  fake_item_server := &test_utils.FakeItemServiceServer{
    CustomerId : 1,
    ItemId : 2,
    Title : "fuga",
    Price : 1234,
    Count : 3 }
  item.RegisterItemServiceServer( server, fake_item_server )
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
    client := app_proto.NewCatalogServiceClient( connection )
    f.Fuzz( func( t *testing.T, customer_id int, item_id int, name_ []byte, title_ []byte, price int64, count uint64 ) {
      name := test_utils.ToValidUTF8StringBiased( name_, 1000 )
      title := test_utils.ToValidUTF8StringBiased( title_, 1000 )
      count = count % 10
      fake_customer_server.CustomerId = customer_id
      fake_customer_server.Name = name
      fake_item_server.CustomerId = customer_id
      fake_item_server.ItemId = item_id
      fake_item_server.Title = title
      fake_item_server.Price = price
      fake_item_server.Count = count
      customer_id_in_str := strconv.Itoa( customer_id )
      {
        md := metadata.New(map[string]string{ "authorization": "bearer "+test_utils.GetAccessToken( customer_id ) })
        ctx := metadata.NewOutgoingContext( context.Background(), md )
        item_id_in_str := strconv.Itoa( item_id )
        {
          response, e4 := client.CreateItem( ctx, &app_proto.CreateItemRequest{ Title : title, Price : price } )
          if e4 != nil {
            panic( e4 )
          }
          if response.Item.Id != item_id_in_str {
            panic( "unexpected item id" )
          }
          if response.Item.CustomerId != customer_id_in_str {
            panic( "unexpected customer id" )
          }
          /*if response.Item.Title != title {
            panic( "unexpected name" )
          }*/
          if response.Item.Price != price {
            panic( "unexpected price" )
          }
        }
        {
          response, e4 := client.GetItem( ctx, &app_proto.GetItemRequest{ Id: item_id_in_str } )
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
          response, e4 := client.ListItems( ctx, &app_proto.ListItemsRequest{ Id: customer_id_in_str } )
          if e4 != nil {
            panic( e4 )
          }
	  if uint64( len( response.Items ) ) != count {
            panic( fmt.Sprintf( "unexpected response size %d %d", uint64( len( response.Items ) ), count ) )
          }
	  if len( response.Items ) > 0 {
            if response.Items[ 0 ].Id != item_id_in_str {
              panic( "unexpected item id" )
            }
            if response.Items[ 0 ].Title != title {
              panic( "unexpected name" )
            }
            if response.Items[ 0 ].Price != price {
              panic( "unexpected price" )
            }
          }
        }
      }
    })
  }
}


