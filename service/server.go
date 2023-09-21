package service

import (
	"context"

	"github.com/superjcd/productservice/config"
	v1 "github.com/superjcd/productservice/genproto/v1"
	"github.com/superjcd/productservice/pkg/database"
	"github.com/superjcd/productservice/service/store"
	"github.com/superjcd/productservice/service/store/sql"
	"gorm.io/gorm"
)

var _DB *gorm.DB

// Server Server struct
type Server struct {
	v1.UnimplementedProductServiceServer
	datastore store.Factory
	client    v1.ProductServiceClient
	conf      *config.Config
}

// NewServer New service grpc server
func NewServer(conf *config.Config, client v1.ProductServiceClient) (v1.ProductServiceServer, error) {
	_DB = database.MustPreParePostgresqlDb(&conf.Pg)
	factory, err := sql.NewSqlStoreFactory(_DB)
	if err != nil {
		return nil, err
	}

	server := &Server{
		client:    client,
		datastore: factory,
		conf:      conf,
	}

	return server, nil
}

func (s *Server) CreateProduct(ctx context.Context, rq *v1.CreateProductRequest) (*v1.CreateProductResponse, error) {
	if err := s.datastore.Products().Create(ctx, rq); err != nil {
		return &v1.CreateProductResponse{Msg: "创建失败", Status: v1.Status_failure}, err
	}
	return &v1.CreateProductResponse{
		Msg:    "创建成功",
		Status: v1.Status_success,
	}, nil
}

func (s *Server) ListProduct(ctx context.Context, rq *v1.ListProductRequest) (*v1.ListProductResponse, error) {
	groups, err := s.datastore.Products().List(ctx, rq)
	if err != nil {
		return &v1.ListProductResponse{Msg: "获取列表失败", Status: v1.Status_failure}, err
	}

	resp := groups.ConvertToListProductResponse("成功获取列表", v1.Status_success)

	return &resp, nil
}

func (s *Server) UpdateProduct(ctx context.Context, rq *v1.UpdateProductRequest) (*v1.UpdateProductResponse, error) {
	if err := s.datastore.Products().Update(ctx, rq); err != nil {
		return &v1.UpdateProductResponse{Msg: "失败", Status: v1.Status_success}, err
	}

	return &v1.UpdateProductResponse{
		Msg:    "更新成功",
		Status: v1.Status_success,
	}, nil

}

func (s *Server) DeleteProduct(ctx context.Context, rq *v1.DeleteProductRequest) (*v1.DeleteProductResponse, error) {
	if err := s.datastore.Products().Delete(ctx, rq); err != nil {
		return &v1.DeleteProductResponse{Msg: "删除失败", Status: v1.Status_failure}, err
	}

	return &v1.DeleteProductResponse{Msg: "删除成功", Status: v1.Status_success}, nil
}

func (s *Server) GetAmzProductLatestInfo(ctx context.Context, rq *v1.GetAmzProductLatestInfoRequest) (*v1.GetAmzProductLatestInfoResponse, error) {
	if info, err := s.datastore.ProductDetails().GetlatestInfo(ctx, rq); err != nil {
		return &v1.GetAmzProductLatestInfoResponse{Msg: "删除失败", Status: v1.Status_failure}, err
	} else {
		resp := info.ConvertToGetLatestInfoResponse("成功获取列表", v1.Status_success)

		return &resp, nil
	}
}

func (s *Server) AppendAmzProductActiveDetail(ctx context.Context, rq *v1.AppendAmzProductActiveDetailRequest) (*v1.AppendAmzProductActiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendActiveDetail(ctx, rq); err != nil {
		return &v1.AppendAmzProductActiveDetailResponse{Msg: "追加Active details 失败", Status: v1.Status_failure}, err
	}
	return &v1.AppendAmzProductActiveDetailResponse{Msg: "追加Active details成功", Status: v1.Status_success}, nil
}

func (s *Server) AppendAmzProductInactiveDetail(ctx context.Context, rq *v1.AppendAmzProductInactiveDetailRequest) (*v1.AppendAmzProductInactiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendInactiveDetail(ctx, rq); err != nil {
		return &v1.AppendAmzProductInactiveDetailResponse{Msg: "追加Inactive details 失败", Status: v1.Status_failure}, err
	}
	return &v1.AppendAmzProductInactiveDetailResponse{Msg: "追加Inactive details成功", Status: v1.Status_success}, nil
}
