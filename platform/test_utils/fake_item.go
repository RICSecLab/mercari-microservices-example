package test_utils

import (
  "context"
  "strconv"
  item "github.com/mercari/mercari-microservices-example/services/item/proto"
)

type FakeItemServiceServer struct {
  item.UnimplementedItemServiceServer
  CustomerId int
  ItemId int
  Title string
  Price int64
}

func (this *FakeItemServiceServer) CreateItem( ctx context.Context, request *item.CreateItemRequest) (*item.CreateItemResponse, error) {
  item_id := strconv.Itoa( this.ItemId )
  return &item.CreateItemResponse{ Item : &item.Item{ Id: item_id, CustomerId: request.CustomerId, Title: request.Title, Price: request.Price } }, nil
}
func (this *FakeItemServiceServer) GetItem( ctx context.Context, request *item.GetItemRequest) (*item.GetItemResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &item.GetItemResponse{ Item : &item.Item{ Id: request.Id, CustomerId: customer_id, Title: this.Title, Price: this.Price } }, nil
}
func (this *FakeItemServiceServer) ListItems(context.Context, *item.ListItemsRequest) (*item.ListItemsResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  item_id := strconv.Itoa( this.ItemId )
  return &item.ListItemsResponse{ Items : []*item.Item{ &item.Item{ Id: item_id, CustomerId: customer_id, Title: this.Title, Price: this.Price } } }, nil
}

