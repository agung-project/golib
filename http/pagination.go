package http

import (
	"net/http"
	"strconv"
	"strings"
)

func GetPagination(r *http.Request, defaultPageSize int, fields PaginationFields) RequestPagination {
	query := r.URL.Query().Get(fields.QueryField)
	sort := r.URL.Query().Get(fields.SortField)
	page, _ := strconv.Atoi(r.URL.Query().Get(fields.OffsetField))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get(fields.LimitField))

	sorts := []string{}
	if sort != "" {
		sorts = strings.Split(sort, ",")
	}

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	return RequestPagination{
		Query:  query,
		Sort:   sorts,
		Limit:  pageSize,
		Offset: (page - 1) * pageSize,
	}
}

func GetNextPagination(reqPage RequestPagination, dataCount int64) Pagination {
	count := int(dataCount)
	currentPage := (reqPage.Offset / reqPage.Limit) + 1
	pageSize := reqPage.Limit

	totalPage := count / pageSize
	if count%pageSize > 0 {
		totalPage++
	}

	if currentPage >= totalPage {
		currentPage = totalPage
	}

	nextPage := currentPage + 1
	if nextPage >= totalPage {
		nextPage = totalPage
	}

	return Pagination{
		PageSize:    pageSize,
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		NextPage:    nextPage,
		TotalData:   count,
	}
}
