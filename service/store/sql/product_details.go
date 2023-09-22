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

func (p *product_details) AppendActiveDetail(ctx context.Context, rq *v1.AppendAmzProductActiveDetailRequest) error {
	activeInfos := make([]store.AmzProductActiveDetail, 0, 16)

	for _, d := range rq.Details {
		activeInfo := store.AmzProductActiveDetail{
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

func (p *product_details) ListActiveDetails(ctx context.Context, rq *v1.ListAmzProductDetailsRequest) (*store.ProductDetails, error) {
	result := &store.ProductDetails{} // 店铺国家时间

	sql := fmt.Sprintf(`
		SELECT 
		  t1.*
		FROM
		(
			SELECT 
				*
			FROM 
				amz_product_active_details
			WHERE 
				create_date = '%s'
		)t1 LEFT JOIN 	(
			SELECT
				asin,
				country
			FROM 
				products
			WHERE shop = '%s'
				and country = '%s'
			) t2 ON t1.asin = t2.asin 
			   AND t1.country = t2.country`, rq.CreateDate, rq.Shop, rq.Country)

	d := p.db.Raw(sql).Offset(int(rq.Offset)).Limit(int(rq.Limit)).Scan(&result.Items).Offset(-1).Limit(-1).Count(&result.TotalCount)
	return result, d.Error
}

func (p *product_details) AppendInactiveDetail(ctx context.Context, rq *v1.AppendAmzProductInactiveDetailRequest) error {
	inactiveInfos := make([]store.AmzProductInactiveDetail, 0, 16)

	for _, d := range rq.Details {
		inactiveInfo := store.AmzProductInactiveDetail{
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
		return p.db.Where("create_date < ?", rq.MinCreateDate).Delete(&store.AmzProductActiveDetail{}).Error
	}

	return fmt.Errorf("min_create_date不能为空")
}

func (p *product_details) DeleteInactiveDetail(ctx context.Context, rq *v1.DeleteAmzProductInactiveDetailRequest) error {
	if rq.MinCreateDate != "" {
		return p.db.Unscoped().Where("create_date < ?", rq.MinCreateDate).Delete(&store.AmzProductInactiveDetail{}).Error
	}
	return fmt.Errorf("min_create_date不能为空")
}

func (p *product_details) GetProductHistoryInfo(ctx context.Context, rq *v1.GetProductHistoryInfoRequest) ([]store.ProductHistoryInfoRecord, error) {
	records := make([]store.ProductHistoryInfoRecord, 0, 16)

	sql := fmt.Sprintf(`
		SELECT 
			create_date as datetime,
			%s as value
		FROM 
		    amz_product_active_details
		WHERE 
			asin='%s'
			AND country='%s'
			AND create_date BETWEEN '%s' AND '%s'
	`, rq.Field, rq.Asin, rq.Country, rq.StartDate, rq.EndDate)

	d := p.db.Raw(sql).Scan(records)

	if d.Error != nil {
		return nil, d.Error
	}
	return records, nil
}
