package store

import (
	"context"

	v1 "github.com/superjcd/productservice/genproto/v1"
)

type Factory interface {
	Products() ProductStore
	ProductDetails() ProductDetailsStore
	ProductChanges() ProductChangeStore
	InactiveProductChanges() InactiveProductChangeStore
	Close() error
}

type ProductStore interface {
	Create(ctx context.Context, _ *v1.CreateProductRequest) error
	List(ctx context.Context, _ *v1.ListProductRequest) (*ProductList, error)
	Update(ctx context.Context, _ *v1.UpdateProductRequest) error
	Delete(ctx context.Context, _ *v1.DeleteProductRequest) error
}

type ProductDetailsStore interface {
	AppendActiveDetail(ctx context.Context, _ *v1.AppendAmzProductActiveDetailRequest) error
	ListActiveDetails(ctx context.Context, _ *v1.ListAmzProductDetailsRequest) (*ProductDetails, error)
	DeleteActiveDetail(ctx context.Context, _ *v1.DeleteAmzProductActiveDetailRequest) error
	AppendInactiveDetail(ctx context.Context, _ *v1.AppendAmzProductInactiveDetailRequest) error
	DeleteInactiveDetail(ctx context.Context, _ *v1.DeleteAmzProductInactiveDetailRequest) error
	GetProductHistoryInfo(ctx context.Context, _ *v1.GetProductHistoryInfoRequest) ([]ProductHistoryInfoRecord, error)
}

type ProductChangeStore interface {
	Append(ctx context.Context, _ *v1.AppendProductChangesRequest) error
	List(ctx context.Context, _ *v1.ListProductChangesRequest) (*ProductChangeList, error)
	Delete(ctx context.Context, _ *v1.DeleteProductChangesRequest) error
}

type InactiveProductChangeStore interface {
	Append(ctx context.Context, _ *v1.AppendProductChangesRequest) error
	List(ctx context.Context, _ *v1.ListProductChangesRequest) (*InactiveProductChangeList, error)
	Delete(ctx context.Context, _ *v1.DeleteProductChangesRequest) error
}
