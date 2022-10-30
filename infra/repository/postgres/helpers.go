package repository

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Sort  string `json:"sort,omitempty" query:"sort"`
	Limit int    `json:"limit,omitempty" query:"limit"`
	Page  int    `json:"page,omitempty" query:"page"`

	TotalRows  int         `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}

	return p.Page
}

// GetLimit provides the maximum number of entries to return.
func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p.Limit
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}

	return p.Sort
}

// GetOffset instructs the server where to start returning rows within the query result.
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func paginate(model any, pagination *Pagination, dbConn *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var countAll int64
	dbConn.Model(model).
		Where("deleted_on is null").
		Count(&countAll)
	totalPages := int(math.Ceil(float64(countAll) / float64(pagination.Limit)))

	pagination.TotalRows = int(countAll)
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
