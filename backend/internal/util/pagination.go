package util

import (
	"github.com/shironxn/gocrud/internal/core/domain"
	"gorm.io/gorm"
)

type Pagination struct {
	validator *Validator
}

func NewPagination(validator *Validator) Pagination {
	return Pagination{
		validator: validator,
	}
}

func (p *Pagination) Paginate(metadata *domain.Metadata) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case metadata.Limit > 100:
			metadata.Limit = 100
		case metadata.Limit < 1:
			metadata.Limit = 10
		}

		metadata.TotalPages = (int(metadata.TotalRecords) + metadata.Limit - 1) / metadata.Limit
		if metadata.TotalPages < 1 {
			metadata.TotalPages = 1
		}

		switch {
		case metadata.Page > metadata.TotalPages:
			metadata.Page = metadata.TotalPages
		case metadata.Page < 1:
			metadata.Page = 1
		}

		if err := p.validator.Validate(metadata); err != nil {
			for _, error := range err.Error.([]domain.ValidationError) {
				switch error.Field {
				case "Sort":
					metadata.Sort = "created_at"
				case "Order":
					metadata.Order = "asc"

				}
			}
		}

		offset := (metadata.Page - 1) * metadata.Limit
		return db.
			Order(metadata.Sort + " " + metadata.Order).
			Offset(offset).
			Limit(metadata.Limit)
	}
}
