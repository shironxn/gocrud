package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DB struct {
	cfg *Config
}

func NewGorm(cfg *Config) *DB {
	return &DB{
		cfg: cfg,
	}
}

func (d *DB) Connection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.cfg.Database.User,
		d.cfg.Database.Pass,
		d.cfg.Database.Host,
		d.cfg.Database.Port,
		d.cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: false})
	if err != nil {
		return nil, err
	}

	return db, nil
}
