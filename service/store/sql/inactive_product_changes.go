package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/productservice/genproto/v1"
	"github.com/superjcd/productservice/service/store"
	"gorm.io/gorm"
)

type inactive_product_changes struct {
	db *gorm.DB
}

// AmzProdutActiveDetai
var _ store.InactiveProductChangeStore = (*inactive_product_changes)(nil)

func (pc *inactive_product_changes) Append(ctx context.Context, rq *v1.AppendProductChangesRequest) error {
	inactiveProductChanges := make([]store.InactiveProductChange, 0, 512)
	sql := fmt.Sprintf(`
	        SELECT
				t1.country country,
				t1.asin asin,
				'%s' as field,
				t1.price as old_value,
				t2.price as new_value,
				'%s' as create_date
			FROM
			(
			Select
				country,
				asin,
				%s
			FROM amz_product_inactive_details
			WHERE create_date = '%s'
				and %s != ''
			)t1 LEFT JOIN
			(
				Select
				country,
				asin,
				%s
			FROM amz_product_inactive_details
			WHERE create_date = '%s'
				and %s != ''
			)t2
			on t1.asin = t2.asin
			and t1.country = t2.country
			where t1.price != t2.price
		`, rq.Field, rq.NewDate, rq.Field, rq.OldDate, rq.Field, rq.Field, rq.NewDate, rq.Field)

	d := pc.db.Raw(sql).Scan(&inactiveProductChanges)

	if d.Error != nil {
		return d.Error
	}

	if len(inactiveProductChanges) > 0 {
		return pc.db.Create(&inactiveProductChanges).Error
	}

	return nil
}

func (pc *inactive_product_changes) List(ctx context.Context, rq *v1.ListProductChangesRequest) (*store.InactiveProductChangeList, error) {
	productChanges := make([]store.InactiveProductChange, 0, 32)
	sql := `
		  SELECT
		    t1.asin,
			t1.country,
			t1.field,
			t1.old_value,
			t1.new_value,
			t1.create_date
          FROM (
			SELECT 
				asin,
				country,
				field,
				old_value, 
				new_value,
				create_date
			FROM inactive_product_changes
			WHERE  country = '%s'		
				AND create_date = '%s'
				AND field = '%s'
			) t1 LEFT JOIN (
				SELECT 
					asin,
					country
				FROM 
					products
				WHERE shop = '%s' 
					AND country = '%s'
			) t2 on t1.country = t2.country 
				AND t1.asin = t2.asin		
	`
	sql = fmt.Sprintf(sql, rq.Country, rq.CreateDate, rq.Field, rq.Shop, rq.Country)

	d := pc.db.Raw(sql).Scan(&productChanges)

	if d.Error != nil {
		return nil, d.Error
	}

	return &store.InactiveProductChangeList{
		TotalCount: len(productChanges),
		Items:      productChanges,
	}, nil

}

func (pc *inactive_product_changes) Delete(ctx context.Context, rq *v1.DeleteProductChangesRequest) error {
	if rq.MinCreateDate != "" {
		return pc.db.Unscoped().Where("create_date < ?", rq.MinCreateDate).Delete(&store.ProductChange{}).Error
	}

	return fmt.Errorf("min_create_date should not be empty")
}
