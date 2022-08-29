package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPagination(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=2&page_size=10", nil)
	paginationFields := PaginationFields{
		QueryField:  "search",
		SortField:   "order",
		LimitField:  "page_size",
		OffsetField: "page",
	}

	page := GetPagination(req,10, paginationFields)
	assert.Equal(t, 10, page.Offset, "Offset should be 10")
	assert.Equal(t, 10, page.Limit, "Limit should be 10")

	paginationFields.OffsetField = "hal"
	paginationFields.LimitField = "batas"
	page = GetPagination(req,1, paginationFields)
	assert.Equal(t, 0, page.Offset, "Offset should be 0")
	assert.Equal(t, 1, page.Limit, "Limit should be 1")
}

func TestGetNextPagination(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?page=1&page_size=10", nil)
	paginationFields := PaginationFields{
		QueryField:  "search",
		SortField:   "order",
		LimitField:  "page_size",
		OffsetField: "page",
	}

	reqPage := GetPagination(req,10, paginationFields)
	var dataCount int64 = 90

	nextPage := GetNextPagination(reqPage, dataCount)
	assert.Equal(t, 2, nextPage.NextPage, "NextPage should be 2")
	assert.Equal(t, 1, nextPage.CurrentPage, "CurrentPage should be 1")
	assert.Equal(t, 9, nextPage.TotalPage, "TotalPage should be 9")
	assert.Equal(t, 10, nextPage.PageSize, "PageSize should be 10")

	req = httptest.NewRequest(http.MethodGet, "/?page=2&page_size=10", nil)
	reqPage = GetPagination(req,10, paginationFields)
	nextPage = GetNextPagination(reqPage, dataCount)
	assert.Equal(t, 3, nextPage.NextPage, "NextPage should be 3")
	assert.Equal(t, 2, nextPage.CurrentPage, "CurrentPage should be 2")
	assert.Equal(t, 9, nextPage.TotalPage, "TotalPage should be 9")
	assert.Equal(t, 10, nextPage.PageSize, "PageSize should be 10")

	dataCount = 33

	req = httptest.NewRequest(http.MethodGet, "/?page=3&page_size=5", nil)
	reqPage = GetPagination(req,10, paginationFields)
	nextPage = GetNextPagination(reqPage, dataCount)
	assert.Equal(t, 4, nextPage.NextPage, "NextPage should be 4")
	assert.Equal(t, 3, nextPage.CurrentPage, "CurrentPage should be 3")
	assert.Equal(t, 7, nextPage.TotalPage, "TotalPage should be 7")
}