package store

import (
	"context"

	v1 "github.com/HooYa-Bigdata/productservice/genproto/v1"
)

type Factory interface {
	Products() ProductStore
	ProductDetails() ProductDetailsStore
	Close() error
}

type ProductStore interface {
	Create(ctx context.Context, _ *v1.CreateProductRequest) error
	List(ctx context.Context, _ *v1.ListProductRequest) (*ProductList, error)
	Update(ctx context.Context, _ *v1.UpdateProductRequest) error
	Delete(ctx context.Context, _ *v1.DeleteProductRequest) error
}

type ProductDetailsStore interface {
	GetlatestInfo(ctx context.Context, _ *v1.GetAmzProductLatestInfoRequest) (*ProductLatestInfo, error)
	// GetHistoryInfo(ctx context.Context, _ *v1.GetAmzProductHistoryInfoRequest) (*ProductLatestInfo, error)
	AppendActiveDetail(ctx context.Context, _ *v1.AppendAmzProductActiveDetailRequest) error
	AppendInactiveDetail(ctx context.Context, _ *v1.AppendAmzProductInactiveDetailRequest) error
}
