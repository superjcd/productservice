package sql

import (
	"context"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "github.com/superjcd/productservice/genproto/v1"
	"github.com/superjcd/productservice/service/store"
	"gorm.io/gorm"
)

var dbFile = "fake.db"

type FakeStoreTestSuite struct {
	suite.Suite
	Dbfile      string
	FakeFactory store.Factory
}

func (suite *FakeStoreTestSuite) SetupSuite() {
	file, err := os.Create(dbFile)
	assert.Nil(suite.T(), err)
	defer file.Close()

	suite.Dbfile = dbFile
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	assert.Nil(suite.T(), err)

	factory, err := NewSqlStoreFactory(db)
	assert.Nil(suite.T(), err)
	suite.FakeFactory = factory
}

func (suite *FakeStoreTestSuite) TearDownSuite() {
	var err error
	err = suite.FakeFactory.Close()
	assert.Nil(suite.T(), err)
	err = os.Remove(dbFile)
	assert.Nil(suite.T(), err)
}

// products
func (suite *FakeStoreTestSuite) TestCreateProduct() {
	products := &v1.CreateProductRequest{
		Products: []*v1.Product{
			{
				Sku:     "1001",
				Shop:    "apple",
				Asin:    "B1001",
				Country: "US",
			},
			{
				Sku:     "1010",
				Shop:    "pear",
				Asin:    "B11111",
				Country: "US",
			},
		},
	}

	err := suite.FakeFactory.Products().Create(context.Background(), products)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestListProducts() {
	rq := &v1.ListProductRequest{
		Shop:   "apple",
		Offset: 0,
		Limit:  10,
	}

	productList, err := suite.FakeFactory.Products().List(context.Background(), rq)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(productList.Items))
}

func (suite *FakeStoreTestSuite) TestUpdateProduct() {
	product := &v1.UpdateProductRequest{
		Sku:     "1001",
		Shop:    "elppa",
		Asin:    "B1001",
		Country: "US",
	}

	err := suite.FakeFactory.Products().Update(context.Background(), product)
	assert.Nil(suite.T(), err)

	// list with old name
	rq := &v1.ListProductRequest{
		Shop:   "apple",
		Offset: 0,
		Limit:  10,
	}
	productList, _ := suite.FakeFactory.Products().List(context.Background(), rq)
	assert.Equal(suite.T(), 0, len(productList.Items))

	// list with new name
	rq2 := &v1.ListProductRequest{
		Shop:   "elppa",
		Offset: 0,
		Limit:  10,
	}
	productList2, _ := suite.FakeFactory.Products().List(context.Background(), rq2)
	assert.Equal(suite.T(), 1, len(productList2.Items))
}

func (suite *FakeStoreTestSuite) TestZDeleteProduct() {
	rq := &v1.DeleteProductRequest{
		Sku:  "1001",
		Asin: "B1001",
	}

	err := suite.FakeFactory.Products().Delete(context.Background(), rq)
	assert.Nil(suite.T(), err)

}

func (suite *FakeStoreTestSuite) TestAppendeDetail() {
	// 导入active
	rq := &v1.AppendAmzProductActiveDetailRequest{
		Details: []*v1.AmzProductActiveDetail{
			{
				Asin:       "B1001",
				Country:    "US",
				Price:      "100",
				CreateDate: "2022-01-01",
			},
		},
	}

	err := suite.FakeFactory.ProductDetails().AppendActiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)

	// 导入inactive
	rq2 := &v1.AppendAmzProductInactiveDetailRequest{
		Details: []*v1.AmzProductInactivateDetail{
			{
				Asin:         "B1001",
				Country:      "US",
				Title:        "Iphone 15",
				BulletPoints: "1 good 2 cheap",
				CreateDate:   "2022-01-01",
			},
		},
	}

	err2 := suite.FakeFactory.ProductDetails().AppendInactiveDetail(context.Background(), rq2)
	assert.Nil(suite.T(), err2)

	// 导入更新的数据
	rq4 := &v1.AppendAmzProductActiveDetailRequest{
		Details: []*v1.AmzProductActiveDetail{
			{
				Asin:       "B1001",
				Country:    "US",
				Price:      "110",
				CreateDate: "2022-01-02",
			},
		},
	}

	err4 := suite.FakeFactory.ProductDetails().AppendActiveDetail(context.Background(), rq4)
	assert.Nil(suite.T(), err4)

}

func (suite *FakeStoreTestSuite) TestXListProductDetails() {
	rq := &v1.ListAmzProductDetailsRequest{
		Shop:       "elppa",
		Country:    "US",
		CreateDate: "2022-01-02",
		Offset:     0,
		Limit:      10,
	}

	_, err := suite.FakeFactory.ProductDetails().ListActiveDetails(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestYAppendProductChanges() {
	rq := &v1.AppendProductChangesRequest{
		OldDate: "2022-01-01",
		NewDate: "2022-01-02",
		Field:   "price",
	}

	err := suite.FakeFactory.ProductChanges().Append(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestYListProductChanges() {
	rq := &v1.ListProductChangesRequest{
		Shop:       "elppa",
		Country:    "US",
		CreateDate: "2022-01-02",
		Field:      "price",
		Offset:     0,
		Limit:      10,
	}
	list, err := suite.FakeFactory.ProductChanges().List(context.Background(), rq)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, list.TotalCount)

}

func (suite *FakeStoreTestSuite) TestZDeleteActiveDetail() {
	rq := &v1.DeleteAmzProductActiveDetailRequest{
		MinCreateDate: "2023-10-01",
	}
	err := suite.FakeFactory.ProductDetails().DeleteActiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestZDeleteInactiveDetail() {
	rq := &v1.DeleteAmzProductInactiveDetailRequest{
		MinCreateDate: "2023-10-01",
	}
	err := suite.FakeFactory.ProductDetails().DeleteInactiveDetail(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

// func (suite *FakeStoreTestSuite) TestGetProductHistoryInfo() {
// 	rq := &v1.GetProductHistoryInfoRequest{
// 		Asin:      "B1001",
// 		Country:   "US",
// 		Field:     "price",
// 		StartDate: "2020-01-01",
// 		EndDate:   "2023-01-01",
// 	}
// 	_, err := suite.FakeFactory.ProductDetails().GetProductHistoryInfo(context.Background(), rq)
// 	assert.Nil(suite.T(), err)
// }

func TestFakeStoreSuite(t *testing.T) {
	suite.Run(t, new(FakeStoreTestSuite))
}
