package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	config *Config
}

func NewGorm(config *Config) *DB {
	return &DB{
		config: config,
	}
}

func (d *DB) Connection() (*gorm.DB, error) {
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	d.config.Database.User,
	// 	d.config.Database.Pass,
	// 	d.config.Database.Host,
	// 	d.config.Database.Port,
	// 	d.config.Database.Name,
	// )

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: false})
	// if err != nil {
	// 	return nil, err
	// }

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		d.config.Database.Host,
		d.config.Database.User,
		d.config.Database.Pass,
		d.config.Database.Name,
		d.config.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
