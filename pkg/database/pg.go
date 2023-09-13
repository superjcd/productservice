package database

import (
	"fmt"

	"github.com/HooYa-Bigdata/productservice/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustPreParePostgresqlDb(opts *config.Pg) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		opts.Host,
		opts.Username,
		opts.Password,
		opts.Database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	return db
}
