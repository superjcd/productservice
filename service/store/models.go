package store

import (
	v1 "github.com/HooYa-Bigdata/productservice/genproto/v1"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Sku     string `json:"sku" gorm:"column:sku;index:idx_product,priority:2" validate:"required"`
	Shop    string `json:"shop" gorm:"column:shop" validate:"required"`
	Asin    string `json:"asin" gorm:"column:asin;index:idx_product,priority:1" validate:"required"`
	Country string `json:"country" gorm:"column:country" validate:"required"`
}

type ProductList struct {
	TotalCount int64      `json:"totalCount"`
	Items      []*Product `json:"items"`
}

func (pl *ProductList) ConvertToListProductResponse(msg string, status v1.Status) v1.ListProductResponse {
	products := make([]*v1.Product, 16)

	for _, product := range pl.Items {
		products = append(products, &v1.Product{
			Sku:     product.Sku,
			Shop:    product.Shop,
			Asin:    product.Asin,
			Country: product.Country,
		})
	}

	return v1.ListProductResponse{
		Msg:      msg,
		Status:   status,
		Products: products,
	}
}

// timestamp uint64,  可以通过create_at转化过来
type AmzProdutInactiveDetail struct {
	gorm.Model
	Asin         string `json:"asin" gorm:"column:asin;index:idx_inactive_product_details,priority:1"`
	Country      string `json:"country" gorm:"column:country;index:idx_inactive_product_details,priority:2"`
	Title        string `json:"title" gorm:"column:title"`
	BulletPoints string `json:"bullet_points" gorm:"column:bullet_points"`
	CreateDate   string `json:"create_date" gorm:"column:create_date;idx_inactive_product_details,priority:3"`
}

type AmzProdutActiveDetail struct {
	gorm.Model
	Asin            string `json:"asin" gorm:"column:asin;index:idx_product_details,priority:1"`
	Country         string `json:"country" gorm:"column:country;index:idx_product_details,priority:2"`
	Price           string `json:"price" gorm:"column:price"`
	Currency        string `json:"currency" gorm:"column:currency"`
	Coupon          string `json:"coupon" gorm:"column:coupon"`
	Star            string `json:"star" gorm:"column:star"`
	Ratings         uint32 `json:"ratings" gorm:"column:ratings"`
	Image           string `json:"image" gorm:"column:image"`
	ParentAsin      string `json:"parent_asin" gorm:"column:parent_asin"`
	CategoryInfo    string `json:"category_info" gorm:"column:category_info"`
	TopCategoryName string `json:"top_category_name" gorm:"column:top_category_name"`
	TopCategoryRank uint32 `json:"top_category_rank" gorm:"column:top_category_rank"`
	Color           string `json:"color" gorm:"column:color"`
	Weight          string `json:"weight" gorm:"column:weight"`
	WeightUnit      string `json:"weight_unit" gorm:"column:weight_unit"`
	Dimensions      string `json:"dimensions" gorm:"column:dimensions"`
	DimensionsUnit  string `json:"dimensions_unit" gorm:"column:dimensions_unit"`
	CreateDate      string `json:"create_date" gorm:"column:create_date;index:idx_product_details,priority:3"`
}

type ProductLatestInfo struct {
	InactiveDetails *AmzProdutInactiveDetail
	ActiveDetales   *AmzProdutActiveDetail
}

func (p *ProductLatestInfo) ConvertToGetLatestInfoResponse(msg string, status v1.Status) v1.GetAmzProductLatestInfoResponse {
	return v1.GetAmzProductLatestInfoResponse{
		Msg:    msg,
		Status: status,
		AmzProductIncativeDetails: &v1.AmzProductInactivateDetail{
			Asin:         p.InactiveDetails.Asin,
			Country:      p.InactiveDetails.Country,
			Title:        p.InactiveDetails.Title,
			BulletPoints: p.InactiveDetails.BulletPoints,
		},
		AmzProductAtiveDetails: &v1.AmzProductActiveDetail{
			Asin:            p.InactiveDetails.Asin,
			Country:         p.InactiveDetails.Country,
			Price:           p.ActiveDetales.Price,
			Currency:        p.ActiveDetales.Currency,
			Star:            p.ActiveDetales.Star,
			Ratings:         p.ActiveDetales.Ratings,
			Image:           p.ActiveDetales.Image,
			ParentAsin:      p.ActiveDetales.ParentAsin,
			CategoryInfo:    p.ActiveDetales.CategoryInfo,
			TopCategoryName: p.ActiveDetales.TopCategoryName,
			TopCategoryRank: p.ActiveDetales.TopCategoryRank,
			Color:           p.ActiveDetales.Color,
			Weight:          p.ActiveDetales.Weight,
			WeightUnit:      p.ActiveDetales.WeightUnit,
			Dimensions:      p.ActiveDetales.Dimensions,
			DimensionsUnit:  p.ActiveDetales.DimensionsUnit,
			CreateDate:      p.ActiveDetales.CreateDate,
		},
	}
}

type ProductHistoryInfos struct {
	TotalCount int64    `json:"totalCount"`
	Field      string   `json:"field"`
	Datetimes  []string `json:"datetimes"`
	Values     []string `json:"values"`
}

func (ph *ProductHistoryInfos) ConvertToListProductResponse(msg string, status v1.Status) v1.GetAmzProductHistoryInfoResponse {
	return v1.GetAmzProductHistoryInfoResponse{
		Msg:    msg,
		Status: status,
		Data: &v1.AmzProductHistoryInfoResponse{
			Datetimes: ph.Datetimes,
			Values:    ph.Values,
		},
	}
}

func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(Product{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(AmzProdutActiveDetail{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(AmzProdutInactiveDetail{}); err != nil {
		return err
	}
	return nil
}
