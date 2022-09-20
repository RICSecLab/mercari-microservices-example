package main_test

import (
	"testing"
	"github.com/google/gofuzz"

	"github.com/mercari/mercari-microservices-example/services/authority/grpc"
)

func TestAnalyzer(t *testing.T) {
  f := fuzz.New()
  object := "";
  for i := 0; i < 100; i++ {
    f.Fuzz(&object)
    grpc.CreateAccessToken( object )
  }
}
