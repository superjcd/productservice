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
	products, err := s.datastore.Products().List(ctx, rq)
	if err != nil {
		return &v1.ListProductResponse{Msg: "获取列表失败", Status: v1.Status_failure}, err
	}

	resp := products.ConvertToListProductResponse("成功获取列表", v1.Status_success)

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

func (s *Server) AppendAmzProductActiveDetail(ctx context.Context, rq *v1.AppendAmzProductActiveDetailRequest) (*v1.AppendAmzProductActiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendActiveDetail(ctx, rq); err != nil {
		return &v1.AppendAmzProductActiveDetailResponse{Msg: "追加Active details 失败", Status: v1.Status_failure}, err
	}
	return &v1.AppendAmzProductActiveDetailResponse{Msg: "追加Active details成功", Status: v1.Status_success}, nil
}

func (s *Server) ListProductAmzDetals(ctx context.Context, rq *v1.ListAmzProductDetailsRequest) (*v1.ListAmzProductDetailsResponse, error) {
	details, err := s.datastore.ProductDetails().ListActiveDetails(ctx, rq)

	if err != nil {
		return &v1.ListAmzProductDetailsResponse{Msg: "get product details failed", Status: v1.Status_failure}, err
	}

	resp := details.ConvertToListProductdetails("get product details successfully", v1.Status_success)

	return &resp, nil
}

func (s *Server) AppendAmzProductInactiveDetail(ctx context.Context, rq *v1.AppendAmzProductInactiveDetailRequest) (*v1.AppendAmzProductInactiveDetailResponse, error) {
	if err := s.datastore.ProductDetails().AppendInactiveDetail(ctx, rq); err != nil {
		return &v1.AppendAmzProductInactiveDetailResponse{Msg: "append Inactive details failed", Status: v1.Status_failure}, err
	}
	return &v1.AppendAmzProductInactiveDetailResponse{Msg: "append Inactive details success", Status: v1.Status_success}, nil
}

func (s *Server) AppendProductChanges(ctx context.Context, rq *v1.AppendProductChangesRequest) (*v1.AppendProductChangesResponse, error) {
	if err := s.datastore.ProductChanges().Append(ctx, rq); err != nil {
		return &v1.AppendProductChangesResponse{Msg: "append product changes failed", Status: v1.Status_failure}, err
	}
	return &v1.AppendProductChangesResponse{Msg: "append product changes success", Status: v1.Status_success}, nil
}

func (s *Server) ListProductChanges(ctx context.Context, rq *v1.ListProductChangesRequest) (*v1.ListProductChangesResponse, error) {
	productChanges, err := s.datastore.ProductChanges().List(ctx, rq)
	if err != nil {
		return &v1.ListProductChangesResponse{Msg: "get product changes list failed", Status: v1.Status_failure}, err
	}

	resp := productChanges.ConvertToListProductChangeResponse("get product changes list success", v1.Status_success)

	return &resp, nil
}

func (s *Server) DeleteProductChanges(ctx context.Context, rq *v1.DeleteProductChangesRequest) (*v1.DeleteProductChangesResponse, error) {
	if err := s.datastore.ProductChanges().Delete(ctx, rq); err != nil {
		return &v1.DeleteProductChangesResponse{Msg: "delete product changes failed", Status: v1.Status_failure}, err
	}
	return &v1.DeleteProductChangesResponse{Msg: "delete product changes success", Status: v1.Status_success}, nil
}

func (s *Server) AppendInactiveProductChanges(ctx context.Context, rq *v1.AppendProductChangesRequest) (*v1.AppendProductChangesResponse, error) {
	if err := s.datastore.InactiveProductChanges().Append(ctx, rq); err != nil {
		return &v1.AppendProductChangesResponse{Msg: "append product changes failed", Status: v1.Status_failure}, err
	}
	return &v1.AppendProductChangesResponse{Msg: "append product changes success", Status: v1.Status_success}, nil
}

func (s *Server) ListInactiveProductChanges(ctx context.Context, rq *v1.ListProductChangesRequest) (*v1.ListProductChangesResponse, error) {
	productChanges, err := s.datastore.InactiveProductChanges().List(ctx, rq)
	if err != nil {
		return &v1.ListProductChangesResponse{Msg: "get product changes list failed", Status: v1.Status_failure}, err
	}

	resp := productChanges.ConvertToListInactiveProductChangeResponse("get product changes list success", v1.Status_success)

	return &resp, nil
}

func (s *Server) DeleteInactiveProductChanges(ctx context.Context, rq *v1.DeleteProductChangesRequest) (*v1.DeleteProductChangesResponse, error) {
	if err := s.datastore.InactiveProductChanges().Delete(ctx, rq); err != nil {
		return &v1.DeleteProductChangesResponse{Msg: "delete product changes failed", Status: v1.Status_failure}, err
	}
	return &v1.DeleteProductChangesResponse{Msg: "delete product changes success", Status: v1.Status_success}, nil
}

// GetProductHistoryInfo
func (s *Server) GetProductHistoryInfo(ctx context.Context, rq *v1.GetProductHistoryInfoRequest) (*v1.GetProductHistoryInfoResponse, error) {
	records := make([]*v1.HistoryInfoRecord, 0, 16)

	items, err := s.datastore.ProductDetails().GetProductHistoryInfo(ctx, rq)

	if err != nil {
		return &v1.GetProductHistoryInfoResponse{Msg: "get product history info failed", Status: v1.Status_failure}, err
	}

	for _, item := range items {
		records = append(records, &v1.HistoryInfoRecord{Datetime: item.Datetime, Value: item.Value})
	}

	return &v1.GetProductHistoryInfoResponse{Msg: "get product history info successfully", Status: v1.Status_success, Records: records}, nil
}
