// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package server

import (
	"github.com/HooYa-Bigdata/productservice/config"
	"github.com/HooYa-Bigdata/productservice/genproto/v1"
	"github.com/HooYa-Bigdata/productservice/service"
)

// Injectors from wire.go:

// InitServer Inject service's component
func InitServer(conf *config.Config) (v1.ProductServiceServer, error) {
	productServiceClient, err := service.NewClient(conf)
	if err != nil {
		return nil, err
	}
	productServiceServer, err := service.NewServer(conf, productServiceClient)
	if err != nil {
		return nil, err
	}
	return productServiceServer, nil
}
