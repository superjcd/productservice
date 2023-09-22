package sql

import (
	"fmt"
	"sync"

	"github.com/superjcd/productservice/service/store"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Products() store.ProductStore {
	return &products{db: ds.db}
}

func (ds *datastore) ProductDetails() store.ProductDetailsStore {
	return &product_details{db: ds.db}
}

func (ds *datastore) ProductChanges() store.ProductChangeStore {
	return &product_changes{db: ds.db}
}

func (ds *datastore) InactiveProductChanges() store.InactiveProductChangeStore {
	return &inactive_product_changes{db: ds.db}
}

var (
	sqlFactory store.Factory
	once       sync.Once
)

func NewSqlStoreFactory(db *gorm.DB) (store.Factory, error) {
	if db == nil && sqlFactory == nil {
		return nil, fmt.Errorf("failed to get sql store fatory")
	}
	once.Do(func() {
		store.MigrateDatabase(db)
		sqlFactory = &datastore{db: db}
	})

	return sqlFactory, nil
}

func (ds *datastore) Close() error {
	db, _ := ds.db.DB()

	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
