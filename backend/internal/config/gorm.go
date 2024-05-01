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
