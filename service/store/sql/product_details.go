package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/productservice/genproto/v1"
	"github.com/superjcd/productservice/service/store"
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
	activeInfos := make([]store.AmzProdutActiveDetail, 0, 16)

	for _, d := range rq.Details {
		activeInfo := store.AmzProdutActiveDetail{
			Asin:            d.Asin,
			Country:         d.Country,
			Price:           d.Price,
			Currency:        d.Currency,
			Coupon:          d.Coupon,
			Star:            d.Star,
			Ratings:         d.Ratings,
			Image:           d.Image,
			ParentAsin:      d.ParentAsin,
			CategoryInfo:    d.CategoryInfo,
			TopCategoryName: d.TopCategoryName,
			TopCategoryRank: d.TopCategoryRank,
			Color:           d.Color,
			Weight:          d.Weight,
			WeightUnit:      d.WeightUnit,
			Dimensions:      d.Dimensions,
			DimensionsUnit:  d.DimensionsUnit,
			CreateDate:      d.CreateDate,
		}
		activeInfos = append(activeInfos, activeInfo)
	}

	return p.db.Create(&activeInfos).Error
}

func (p *product_details) AppendInactiveDetail(ctx context.Context, rq *v1.AppendAmzProductInactiveDetailRequest) error {
	inactiveInfos := make([]store.AmzProdutInactiveDetail, 0, 16)

	for _, d := range rq.Details {
		inactiveInfo := store.AmzProdutInactiveDetail{
			Asin:         d.Asin,
			Country:      d.Country,
			Title:        d.Title,
			BulletPoints: d.BulletPoints,
			CreateDate:   d.CreateDate,
		}

		inactiveInfos = append(inactiveInfos, inactiveInfo)
	}

	return p.db.Create(&inactiveInfos).Error
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
