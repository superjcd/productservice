package sql

import (
	"context"
	"fmt"

	v1 "github.com/HooYa-Bigdata/productservice/genproto/v1"
	"github.com/HooYa-Bigdata/productservice/service/store"
	"gorm.io/gorm"
)

type product_details struct {
	db *gorm.DB
}

var _ store.ProductDetailsStore = (*product_details)(nil)

func (p *product_details) GetlatestInfo(ctx context.Context, rq *v1.GetAmzProductLatestInfoRequest) (*store.ProductLatestInfo, error) {
	inactive_details := store.AmzProdutInactiveDetail{}
	active_details := store.AmzProdutActiveDetail{}

	q1 := p.db.Where("asin = ? and country = ?", rq.Asin, rq.Country).
		Last(&inactive_details)

	if q1.Error != nil {
		return nil, q1.Error
	}

	q2 := p.db.Where("asin = ? and country = ?", rq.Asin, rq.Country).
		Last(&active_details)

	if q2.Error != nil {
		return nil, q2.Error
	}

	return &store.ProductLatestInfo{
		InactiveDetails: &inactive_details,
		ActiveDetales:   &active_details,
	}, nil

}

func (p *product_details) AppendActiveDetail(ctx context.Context, rq *v1.AppendAmzProductActiveDetailRequest) error {
	active_info := store.AmzProdutActiveDetail{
		Asin:            rq.Details.Asin,
		Country:         rq.Details.Country,
		Price:           rq.Details.Price,
		Currency:        rq.Details.Currency,
		Coupon:          rq.Details.Coupon,
		Star:            rq.Details.Star,
		Ratings:         rq.Details.Ratings,
		Image:           rq.Details.Image,
		ParentAsin:      rq.Details.ParentAsin,
		CategoryInfo:    rq.Details.CategoryInfo,
		TopCategoryName: rq.Details.TopCategoryName,
		TopCategoryRank: rq.Details.TopCategoryRank,
		Color:           rq.Details.Color,
		Weight:          rq.Details.Weight,
		WeightUnit:      rq.Details.WeightUnit,
		Dimensions:      rq.Details.Dimensions,
		DimensionsUnit:  rq.Details.DimensionsUnit,
		CreateDate:      rq.Details.CreateDate,
	}

	return p.db.Create(&active_info).Error
}

func (p *product_details) AppendInactiveDetail(ctx context.Context, rq *v1.AppendAmzProductInactiveDetailRequest) error {
	inactive_info := store.AmzProdutInactiveDetail{
		Asin:         rq.Asin,
		Country:      rq.Country,
		Title:        rq.Title,
		BulletPoints: rq.BulletPoints,
		CreateDate:   rq.CreateDate,
	}

	return p.db.Create(&inactive_info).Error
}

func (p *product_details) DeleteActiveDetail(ctx context.Context, rq *v1.DeleteAmzProductActiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return p.db.Where("create_date < ?", rq.MinCreateDate).Delete(&store.AmzProdutActiveDetail{}).Error
	}

	return fmt.Errorf("min_create_date不能为空")
}

func (p *product_details) DeleteInactiveDetail(ctx context.Context, rq *v1.DeleteAmzProductInactiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return p.db.Where("create_date < ?", rq.MinCreateDate).Delete(&store.AmzProdutInactiveDetail{}).Error
	}
	return fmt.Errorf("min_create_date不能为空")
}
