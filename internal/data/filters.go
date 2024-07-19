package data

import (
	"math"
	"strings"

	"gl_api.malyshev.io/internal/validator"
)

// фильтры для фильмов
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "должно быть больше 0")
	v.Check(f.Page <= 10_000_000, "page", "не должно превышать 10 млн.")
	v.Check(f.PageSize > 0, "page_size", "должно быть больше 0")
	v.Check(f.PageSize <= 100, "page_size", "не должно превышать 10 млн.")

	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "неверное значение сортировки")
}

func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("неверный параметр сортировки: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, PageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     PageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(PageSize))),
		TotalRecords: totalRecords,
	}
}
