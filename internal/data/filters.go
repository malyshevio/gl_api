package data

import "gl_api.malyshev.io/internal/validator"

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
