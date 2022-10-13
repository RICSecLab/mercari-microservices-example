package test_utils

import (
  "context"
  "strconv"
  catalog "github.com/mercari/mercari-microservices-example/services/catalog/proto"
)

type FakeCatalogServiceServer struct {
  catalog.UnimplementedCatalogServiceServer
  CustomerId int
  ItemId int
  Title string
  Price int64
}

func (this *FakeCatalogServiceServer) CreateItem( ctx context.Context, request *catalog.CreateItemRequest) (*catalog.CreateItemResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  item_id := strconv.Itoa( this.ItemId )
  return &catalog.CreateItemResponse{ Item : &catalog.Item{ Id: item_id, CustomerId: customer_id, Title: request.Title, Price: request.Price } }, nil
}
func (this *FakeCatalogServiceServer) GetItem( ctx context.Context, request *catalog.GetItemRequest) (*catalog.GetItemResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &catalog.GetItemResponse{ Item : &catalog.Item{ Id: request.Id, CustomerId: customer_id, Title: this.Title, Price: this.Price } }, nil
}
func (this *FakeCatalogServiceServer) ListItems( ctx context.Context, request *catalog.ListItemsRequest) (*catalog.ListItemsResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  item_id := strconv.Itoa( this.ItemId )
  return &catalog.ListItemsResponse{ Items : []*catalog.Item{ &catalog.Item{ Id: item_id, CustomerId: customer_id, Title: this.Title, Price: this.Price } } }, nil
}

