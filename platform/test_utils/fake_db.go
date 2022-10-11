package test_utils

import (
  "context"
  "strconv"
  db "github.com/mercari/mercari-microservices-example/platform/db/proto"
)

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

