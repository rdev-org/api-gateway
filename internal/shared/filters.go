package shared

import (
	"net/url"
	"strconv"
	"strings"
)

type Filter struct {
	Query    url.Values
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	Sort     string `json:"sort" query:"sort"`
	Search   string `json:"search" query:"search"`

	UserId string `json:"user_id"`
	Status string `json:"status" query:"status"`
}

func (f *Filter) SanitizeFilter(cols []string) Filter {
	search := f.Query.Get("search")
	pageSizeInt, _ := strconv.Atoi(f.Query.Get("page_size"))
	pageSize := pageSizeInt
	if pageSize <= 0 {
		pageSize = 10
	}

	pageInt, _ := strconv.Atoi(f.Query.Get("page"))
	page := (pageInt - 1) * pageSize

	if search != "" {
		cleanSpaces := strings.TrimSpace(search)
		f.Search = cleanSpaces
	}

	f.Sort = f.Query.Get("sort")
	return Filter{
		PageSize: pageSize,
		Page:     page,
		Sort:     f.ValidateSort(cols),
		Search:   f.Search,
		Status:   f.Query.Get("status"),
		UserId:   f.Query.Get("user_id"),
	}
}

func (f *Filter) ValidateSort(allowedSortColumns []string) string {
	finalSort := "id ASC"
	if f.Sort != "" {
		for _, col := range allowedSortColumns {
			if f.Sort == col || f.Sort == col+" ASC" || f.Sort == col+" DESC" {
				finalSort = f.Sort
				break
			}
		}
	}

	return finalSort
}
