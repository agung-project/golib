package http

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog"
)

type HandlerContextV2 struct {
	E       map[error]*ErrorResponse
	IsDebug bool
	Logger  zerolog.Logger
}

func NewContextHandlerV2(isDebug bool) HandlerContextV2 {
	var errMap = map[error]*ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		ErrInvalidHeader:          ErrInvalidHeader,
		ErrUnauthorized:           ErrUnauthorized,
		ErrInvalidHeaderSignature: ErrInvalidHeaderSignature,
		ErrInvalidHeaderTime:      ErrInvalidHeaderTime,
	}

	return HandlerContextV2{
		E:       errMap,
		IsDebug: isDebug,
	}
}

func (hctx HandlerContextV2) AddError(key error, value *ErrorResponse) {
	hctx.E[key] = value
}

func (hctx HandlerContextV2) AddErrorMap(errMap map[error]*ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k] = v
	}
}

type CustomWriterV2 struct {
	C HandlerContextV2
}

func (c *CustomWriterV2) Write(w http.ResponseWriter, data interface{}, statusCode int, pagination *Pagination, msg []string) {
	var resp ResponseV2
	resp.Success = true

	if msg == nil {
		msg = make([]string, 0)
	}

	resp.Message = msg

	resp.StatusCode = statusCode

	default_data := []interface{}{}
	if data == nil {
		resp.Data = default_data
	} else {
		resp.Data = data
	}

	if pagination != nil {
		data := SuccessResponseV2{}
		data.Pagination = *pagination
		data.Data = resp.Data
		resp.Data = data
	}

	if statusCode == 0 {
		resp.StatusCode = http.StatusOK
	}

	writeResponseV2(w, resp, http.StatusOK)
}

func (c *CustomWriterV2) WritePlain(w http.ResponseWriter, data interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	writeResponse(w, data, "application/json", statusCode)
}

// WriteError sending error response based on err type
func (c *CustomWriterV2) WriteError(w http.ResponseWriter, err error, msg []string) {
	var resp ResponseV2
	resp.Success = false
	resp.Message = msg
	resp.Data = []interface{}{}

	statusCode := http.StatusBadRequest

	var errorResponse = &ErrorResponse{}

	if len(c.C.E) > 0 {
		errorResponse = LookupError(c.C.E, err)
		if errorResponse == nil {
			errorResponse = ErrUnknown
			statusCode = http.StatusInternalServerError
		}
	} else {
		if !(errors.As(err, &errorResponse)) {
			errorResponse = ErrUnknown
			statusCode = http.StatusInternalServerError
		}
	}

	if len(msg) <= 0 {
		resp.Message = append(resp.Message, errorResponse.ResponseDesc)
	}
	resp.StatusCode = errorResponse.HttpStatus

	writeResponseV2(w, resp, statusCode)

}

func writeResponseV2(w http.ResponseWriter, response ResponseV2, statusCode int) {
	writeResponse(w, response, "application/json", statusCode)
}
