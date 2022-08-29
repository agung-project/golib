package http

import (
	"net/http"
)

type Pagination struct {
	PageSize    int `json:"page_size" mapstructure:"page_size"`
	CurrentPage int `json:"current_page" mapstructure:"current_page"`
	TotalPage   int `json:"total_page" mapstructure:"total_page"`
	NextPage    int `json:"next_page" mapstructure:"next_page"`
	TotalData   int `json:"total_data" mapstructure:"total_data"`
}

type PaginationFields struct {
	QueryField  string
	SortField   string
	LimitField  string
	OffsetField string
}

type RequestPagination struct {
	Query  string
	Sort   []string
	Limit  int
	Offset int
}

type HttpHandleResult struct {
	Data            interface{}
	StatusCode      int
	Pagination      *Pagination
	Error           error
	IsPlainResponse bool
}

type HttpHandleResultV2 struct {
	Data            interface{}
	StatusCode      int
	Pagination      *Pagination
	Error           error
	Message         []string
	IsPlainResponse bool
}

type Response struct {
	ResponseDesc string `json:"message" mapstructure:"message"`
}

type ResponseV2 struct {
	StatusCode int         `json:"status" mapstructure:"status"`
	Message    []string    `json:"message" mapstructure:"message"`
	Success    bool        `json:"success" mapstructure:"sucess"`
	Data       interface{} `json:"data" mapstructure:"data"`
}

type SuccessResponseV2 struct {
	Pagination
	Data interface{} `json:"data,omitempty" mapstructure:"data,omitempty"`
}

type SuccessResponse struct {
	Response
	Pagination
	Data interface{} `json:"data,omitempty" mapstructure:"data,omitempty"`
}

// error Response
type ErrorResponse struct {
	Response
	HttpStatus int `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return e.ResponseDesc
}

// ResponseDesc defines details data response
type ResponseDesc struct {
	EN string `json:"en" mapstructure:"en"`
}

var ErrUnknown = &ErrorResponse{
	Response: Response{
		ResponseDesc: "Unknown error",
	},
	HttpStatus: http.StatusInternalServerError,
}

var ErrUnauthorized = &ErrorResponse{
	Response: Response{
		ResponseDesc: "You are not authorized",
	},
	HttpStatus: http.StatusUnauthorized,
}

var ErrInvalidHeader = &ErrorResponse{
	Response: Response{
		ResponseDesc: "Invalid/incomplete header",
	},
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderSignature = &ErrorResponse{
	Response: Response{
		ResponseDesc: "Invalid header signature",
	},
	HttpStatus: http.StatusBadRequest,
}

var ErrInvalidHeaderTime = &ErrorResponse{
	Response: Response{
		ResponseDesc: "Request already expired",
	},
	HttpStatus: http.StatusBadRequest,
}

var ErrRequestEntityTooLarge = &ErrorResponse{
	Response: Response{
		ResponseDesc: "Request entity too large",
	},
	HttpStatus: http.StatusRequestEntityTooLarge,
}
