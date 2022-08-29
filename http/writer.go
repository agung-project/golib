package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"

	"github.com/rs/zerolog"
)

type HandlerContext struct {
	E       map[error]*ErrorResponse
	IsDebug bool
	Logger  zerolog.Logger
}

func NewContextHandler(isDebug bool) HandlerContext {
	var errMap = map[error]*ErrorResponse{
		// register general error here, so if there are new general error you must add it here
		ErrInvalidHeader:          ErrInvalidHeader,
		ErrUnauthorized:           ErrUnauthorized,
		ErrInvalidHeaderSignature: ErrInvalidHeaderSignature,
		ErrInvalidHeaderTime:      ErrInvalidHeaderTime,
	}

	return HandlerContext{
		E:       errMap,
		IsDebug: isDebug,
	}
}

func (hctx HandlerContext) AddError(key error, value *ErrorResponse) {
	hctx.E[key] = value
}

func (hctx HandlerContext) AddErrorMap(errMap map[error]*ErrorResponse) {
	for k, v := range errMap {
		hctx.E[k] = v
	}
}

type CustomWriter struct {
	C HandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, statusCode int, pagination *Pagination) {
	var successResp SuccessResponse
	voData := reflect.ValueOf(data)
	arrayData := []interface{}{}

	if voData.Kind() != reflect.Slice {
		if voData.IsValid() {
			arrayData = []interface{}{data}
		}
		successResp.Data = arrayData
	} else {
		if voData.Len() != 0 {
			successResp.Data = data
		} else {
			successResp.Data = arrayData
		}
	}

	if pagination != nil {
		successResp.Pagination = *pagination
	}

	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	writeSuccessResponse(w, successResp, statusCode)
}

func (c *CustomWriter) WritePlain(w http.ResponseWriter, data interface{}, statusCode int) {
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	writeResponse(w, data, "application/json", statusCode)
}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, err error) {
	if len(c.C.E) > 0 {
		errorResponse := LookupError(c.C.E, err)
		if errorResponse == nil {
			errorResponse = ErrUnknown
		}

		writeErrorResponse(w, errorResponse)
	} else {
		var errorResponse = &ErrorResponse{}
		if errors.As(err, &errorResponse) {
			writeErrorResponse(w, errorResponse)
		} else {
			errorResponse = ErrUnknown
			writeErrorResponse(w, errorResponse)
		}
	}
}

func writeResponse(w http.ResponseWriter, response interface{}, contentType string, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", contentType)

	w.WriteHeader(httpStatus)
	w.Write(res)
}

func writeSuccessResponse(w http.ResponseWriter, response SuccessResponse, statusCode int) {
	writeResponse(w, response, "application/json", statusCode)
}

func writeErrorResponse(w http.ResponseWriter, errorResponse *ErrorResponse) {
	writeResponse(w, errorResponse, "application/json", errorResponse.HttpStatus)
}

// LookupError will get error message based on error type, with variables if you want give dynamic message error
func LookupError(lookup map[error]*ErrorResponse, err error) (res *ErrorResponse) {
	if msg, ok := lookup[err]; ok {
		res = msg
	}

	return
}
