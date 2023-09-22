package store

import (
	v1 "github.com/superjcd/productservice/genproto/v1"
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
type AmzProductInactiveDetail struct {
	gorm.Model
	Asin         string `json:"asin" gorm:"column:asin;index:idx_inactive_product_details,priority:1"`
	Country      string `json:"country" gorm:"column:country;index:idx_inactive_product_details,priority:2"`
	Title        string `json:"title" gorm:"column:title"`
	BulletPoints string `json:"bullet_points" gorm:"column:bullet_points"`
	CreateDate   string `json:"create_date" gorm:"column:create_date;idx_inactive_product_details,priority:3"`
}

type AmzProductActiveDetail struct {
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

type ProductDetails struct {
	TotalCount int64                     `json:"totalCount"`
	Items      []*AmzProductActiveDetail `json:"items"`
}

func (p *ProductDetails) ConvertToListProductdetails(msg string, status v1.Status) v1.ListAmzProductDetailsResponse {
	details := make([]*v1.AmzProductActiveDetail, 0, 16)

	for _, item := range p.Items {
		details = append(details, &v1.AmzProductActiveDetail{
			Asin:            item.Asin,
			Country:         item.Country,
			Price:           item.Price,
			Currency:        item.Currency,
			Star:            item.Star,
			Ratings:         item.Ratings,
			Image:           item.Image,
			ParentAsin:      item.ParentAsin,
			CategoryInfo:    item.CategoryInfo,
			TopCategoryName: item.TopCategoryName,
			TopCategoryRank: item.TopCategoryRank,
			Color:           item.Color,
			Weight:          item.Weight,
			WeightUnit:      item.WeightUnit,
			Dimensions:      item.Dimensions,
			DimensionsUnit:  item.DimensionsUnit,
			CreateDate:      item.CreateDate,
		})
	}

	return v1.ListAmzProductDetailsResponse{Msg: msg, Status: status, Details: details}

}

type ProductHistoryInfoRecord struct {
	Datetime string `json:"datetime"`
	Value    string `json:"value"`
}

type ProductChange struct {
	gorm.Model
	Asin       string `json:"asin" gorm:"column:asin"`
	Country    string `json:"country" gorm:"column:country"`
	Field      string `json:"field" gorm:"column:field"`
	OldValue   string `json:"old_value" gorm:"column:old_value"`
	NewValue   string `json:"new_value" gorm:"column:new_value"`
	CreateDate string `json:"create_date" gorm:"column:create_date"`
}

type ProductChangeList struct {
	TotalCount int             `json:"totalCount"`
	Items      []ProductChange `json:"items"`
}

func (pcl *ProductChangeList) ConvertToListProductChangeResponse(msg string, status v1.Status) v1.ListProductChangesResponse {
	pcs := make([]*v1.ProductChange, 0, 8)

	for _, pc := range pcl.Items {
		pcs = append(pcs, &v1.ProductChange{
			Country:  pc.Country,
			Asin:     pc.Asin,
			Field:    pc.Field,
			OldValue: pc.OldValue,
			NewValue: pc.NewValue,
		})
	}
	return v1.ListProductChangesResponse{
		Msg:            msg,
		Status:         status,
		ProductChanges: pcs,
	}
}

// 创建一张和ProductChange相同结构的表
type InactiveProductChange struct {
	gorm.Model
	Asin       string `json:"asin" gorm:"column:asin"`
	Country    string `json:"country" gorm:"column:country"`
	Field      string `json:"field" gorm:"column:field"`
	OldValue   string `json:"old_value" gorm:"column:old_value"`
	NewValue   string `json:"new_value" gorm:"column:new_value"`
	CreateDate string `json:"create_date" gorm:"column:create_date"`
}

type InactiveProductChangeList struct {
	TotalCount int                     `json:"totalCount"`
	Items      []InactiveProductChange `json:"items"`
}

func (rcl *InactiveProductChangeList) ConvertToListInactiveProductChangeResponse(msg string, status v1.Status) v1.ListProductChangesResponse {
	pcs := make([]*v1.ProductChange, 0, 8)

	for _, pc := range rcl.Items {
		pcs = append(pcs, &v1.ProductChange{
			Country:  pc.Country,
			Asin:     pc.Asin,
			Field:    pc.Field,
			OldValue: pc.OldValue,
			NewValue: pc.NewValue,
		})
	}
	return v1.ListProductChangesResponse{
		Msg:            msg,
		Status:         status,
		ProductChanges: pcs,
	}
}

func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(Product{}, ProductChange{}, InactiveProductChange{}, AmzProductActiveDetail{}, AmzProductInactiveDetail{}); err != nil {
		return err
	}

	return nil
}
