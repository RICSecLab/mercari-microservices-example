package main_test

import (
  "testing"
  "context"
  "runtime"
  "strconv"
  "unicode/utf8"
  "net"
  "google.golang.org/grpc"
  "google.golang.org/grpc/credentials/insecure"
  "github.com/mercari/mercari-microservices-example/pkg/logger"
  db "github.com/mercari/mercari-microservices-example/platform/db/proto"
  app "github.com/mercari/mercari-microservices-example/services/customer/grpc"
  app_proto "github.com/mercari/mercari-microservices-example/services/customer/proto"
)

func toValidUTF8Char( v uint32 ) string {
  if v < 0x20 {
    v = 0x20
  }
  if v > 0x7F && v < 0xA0 {
    v = 0x20
  }
  c := rune( v );
  temp := make([]byte,4)
  size := utf8.EncodeRune( temp, c )
  return string(temp[:size])
}

func toValidUTF8String( r []byte, l int ) string {
  if len( r ) > l * 3 {
    return toValidUTF8String( r[:l*3], l )
  }
  if len( r ) < 1 {
    return ""
  }
  if len( r ) < 2 {
    v := uint32( r[ 0 ] )
    return toValidUTF8Char( v )
  }
  if len( r ) < 3 {
    v := ( uint32( r[ 0 ] ) << 8 ) | uint32( r[ 1 ] ) 
    return toValidUTF8Char( v )
  }
  if len( r ) < 4 {
    v := ( uint32( r[ 0 ] ) << 16 ) | ( uint32( r[ 1 ] ) << 8 ) | ( uint32( r[ 2 ] ) ) 
    return toValidUTF8Char( v )
  }
  v := ( uint32( r[ 0 ] ) << 16 ) | ( uint32( r[ 1 ] ) << 8 ) | ( uint32( r[ 2 ] ) ) 
  return toValidUTF8Char( v ) + toValidUTF8String( r[3:], l )
}

type FakeDBServiceServer struct {
  db.UnimplementedDBServiceServer
  CustomerId int
  ItemId int
  Name string
  Title string
  Price int64
}

func (this *FakeDBServiceServer) CreateCustomer( ctx context.Context, request *db.CreateCustomerRequest) (*db.CreateCustomerResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &db.CreateCustomerResponse{ Customer : &db.Customer{ Id: customer_id, Name: request.Name } }, nil
}
func (this *FakeDBServiceServer) GetCustomer( ctx context.Context, request *db.GetCustomerRequest) (*db.GetCustomerResponse, error) {
  return &db.GetCustomerResponse{ Customer : &db.Customer{ Id: request.Id, Name: this.Name } }, nil
}
func (this *FakeDBServiceServer) GetCustomerByName( ctx context.Context, request *db.GetCustomerByNameRequest) (*db.GetCustomerByNameResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &db.GetCustomerByNameResponse{ Customer : &db.Customer{ Id: customer_id, Name: request.Name } }, nil
}
func (this *FakeDBServiceServer) CreateItem( ctx context.Context, request *db.CreateItemRequest) (*db.CreateItemResponse, error) {
  item_id := strconv.Itoa( this.ItemId )
  return &db.CreateItemResponse{ Item : &db.Item{ Id: item_id, CustomerId: request.CustomerId, Title: request.Title, Price: request.Price } }, nil
}
func (this *FakeDBServiceServer) GetItem( ctx context.Context, request *db.GetItemRequest) (*db.GetItemResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &db.GetItemResponse{ Item : &db.Item{ Id: request.Id, CustomerId: customer_id, Title: this.Title, Price: this.Price } }, nil
}
func (this *FakeDBServiceServer) ListItems(context.Context, *db.ListItemsRequest) (*db.ListItemsResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  item_id := strconv.Itoa( this.ItemId )
  return &db.ListItemsResponse{ Items : []*db.Item{ &db.Item{ Id: item_id, CustomerId: customer_id, Title: this.Title, Price: this.Price } } }, nil
}

func TestCustomer(t *testing.T) {
  socket, e := net.Listen( "tcp", "localhost:0" )
  if e != nil {
    panic( "unable to listen" )
  }
  server := grpc.NewServer()
  db.RegisterDBServiceServer( server, &FakeDBServiceServer{
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
  fake_server := &FakeDBServiceServer{
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
      name := toValidUTF8String( name_, 30 )
      title := toValidUTF8String( title_, 30 )
      fake_server.CustomerId = customer_id
      fake_server.ItemId = customer_id
      fake_server.Name = name
      fake_server.Title = title
      fake_server.Price = price
      customer_id_in_str := strconv.Itoa( customer_id )
      //if utf8.ValidString( name ) {
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
      //}
    })
  }
}


