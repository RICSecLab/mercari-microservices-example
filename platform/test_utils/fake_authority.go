package test_utils

import (
  "context"
  "strconv"
  customer "github.com/mercari/mercari-microservices-example/services/customer/proto"
  authority "github.com/mercari/mercari-microservices-example/services/authority/proto"
)

type FakeAuthorityServiceServer struct {
  authority.UnimplementedAuthorityServiceServer
  CustomerId int
}

func (this *FakeAuthorityServiceServer) Signup( ctx context.Context, request *authority.SignupRequest) (*authority.SignupResponse, error) {
  customer_id := strconv.Itoa( this.CustomerId )
  return &authority.SignupResponse{ Customer : &customer.Customer{ Id: customer_id, Name: request.Name } }, nil
}
func (this *FakeAuthorityServiceServer) Signin( ctx context.Context, request *authority.SigninRequest) (*authority.SigninResponse, error) {
  token := GetAccessToken( this.CustomerId )
  return &authority.SigninResponse{ AccessToken : token }, nil
}

