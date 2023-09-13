package sql

import (
	"context"
	"fmt"

	v1 "github.com/HooYa-Bigdata/productservice/genproto/v1"
	"github.com/HooYa-Bigdata/productservice/service/store"
	"gorm.io/gorm"
)

type products struct {
	db *gorm.DB
}

var _ store.ProductStore = (*products)(nil)

func (p *products) Create(ctx context.Context, rq *v1.CreateProductRequest) error {

	product := store.Product{
		Sku:     rq.Sku,
		Shop:    rq.Shop,
		Asin:    rq.Asin,
		Country: rq.Country,
	}

	return p.db.Create(&product).Error // 我只存储了用户， 但没有处理和用户group有关的逻辑
}

func (p *products) List(ctx context.Context, rq *v1.ListProductRequest) (*store.ProductList, error) {
	result := &store.ProductList{}

	var where_clause string
	if rq.Shop == "" {
		where_clause = "1 = 1"
	} else {
		where_clause = fmt.Sprintf("shop like '%%%s%%'", rq.Shop)
	}

	d := p.db.Where(where_clause).
		Offset(int(rq.Offset)).
		Limit(int(rq.Limit)).
		Find(&result.Items).
		Offset(-1).
		Limit(-1).
		Count(&result.TotalCount)

	return result, d.Error
}

func (p *products) Update(ctx context.Context, rq *v1.UpdateProductRequest) error {
	product := store.Product{}
	if err := p.db.Where("sku = ? and asin = ?", rq.Sku, rq.Asin).First(&product).Error; err != nil {
		return err
	}

	product.Country = rq.Country
	product.Shop = rq.Shop

	return p.db.Save(&product).Error
}

func (p *products) Delete(ctx context.Context, rq *v1.DeleteProductRequest) error {
	return p.db.Where("asin = ? and sku = ? ", rq.Asin, rq.Sku).Delete(&store.Product{}).Error
}
