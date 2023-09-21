//go:build wireinject
// +build wireinject

package server

import (
	"github.com/google/wire"

	"github.com/superjcd/productservice/config"
	v1 "github.com/superjcd/productservice/genproto/v1"
	"github.com/superjcd/productservice/service"
)

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.ProductServiceServer, error) {

	wire.Build(
		service.NewClient,
		service.NewServer,
	)

	return &service.Server{}, nil

}
