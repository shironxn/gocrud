package config

import (
	"fmt"

	"gorm.io/driver/mysql"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.config.Database.User,
		d.config.Database.Pass,
		d.config.Database.Host,
		d.config.Database.Port,
		d.config.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: false})
	if err != nil {
		return nil, err
	}

	return db, nil
}
