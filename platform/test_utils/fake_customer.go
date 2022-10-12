package test_utils

import (
  "context"
  "strconv"
  customer "github.com/mercari/mercari-microservices-example/services/customer/proto"
)

type FakeCustomerServiceServer struct {
  customer.UnimplementedCustomerServiceServer
  CustomerId int
  Name string
}

func (this *FakeCustomerServiceServer) CreateCustomer( ctx context.Context, request *customer.CreateCustomerRequest) (*customer.CreateCustomerResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &customer.CreateCustomerResponse{ Customer : &customer.Customer{ Id: customer_id, Name: request.Name } }, nil
}
func (this *FakeCustomerServiceServer) GetCustomer( ctx context.Context, request *customer.GetCustomerRequest) (*customer.GetCustomerResponse, error) {
  return &customer.GetCustomerResponse{ Customer : &customer.Customer{ Id: request.Id, Name: this.Name } }, nil
}
func (this *FakeCustomerServiceServer) GetCustomerByName( ctx context.Context, request *customer.GetCustomerByNameRequest) (*customer.GetCustomerByNameResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &customer.GetCustomerByNameResponse{ Customer : &customer.Customer{ Id: customer_id, Name: request.Name } }, nil
}

