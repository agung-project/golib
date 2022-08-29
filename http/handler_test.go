package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessHandler(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandler(false)
	newHandler := NewHttpHandler(handlerCtx)

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResult) {
		response.Pagination = nil
		response.Data = "OK"
		response.Error = nil
		response.StatusCode = 200

		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &SuccessResponse{}
	_ = json.Unmarshal(body, respJson)

	assert.Equal(t, 200, resp.StatusCode, "Expect 200 status code")
	assert.Equal(t, `[OK]`, fmt.Sprintf("%s", respJson.Data), "Expect OK data")
}

func TestFailureHandler(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandler(false)
	newHandler := NewHttpHandler(handlerCtx)

	var customErr = errors.New("custom error")
	var ErrResp = &ErrorResponse{
		Response: Response{
			ResponseDesc: "Custom Error",
		},
		HttpStatus: http.StatusBadRequest,
	}
	handlerCtx.AddError(customErr, ErrResp)

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResult) {
		response.Error = customErr
		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &ErrorResponse{}
	_ = json.Unmarshal(body, respJson)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expect 400 status code")
	assert.Equal(t, `Custom Error`, fmt.Sprintf("%s", respJson.ResponseDesc), "Expect Custom Error")
}

func TestSuccessPaginationHandler(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandler(false)
	newHandler := NewHttpHandler(handlerCtx)

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResult) {
		response.Pagination = &Pagination{
			PageSize:    5,
			CurrentPage: 0,
			TotalPage:   10,
			NextPage:    1,
			TotalData:   50,
		}
		response.StatusCode = 200
		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &SuccessResponse{}
	_ = json.Unmarshal(body, respJson)

	assert.Equal(t, 200, resp.StatusCode, "Expect 200 status code")
	assert.Equal(t, 5, respJson.PageSize, "Expect 5 page size")
	assert.Equal(t, 0, respJson.CurrentPage, "Expect 0 current page")
	assert.Equal(t, 10, respJson.TotalPage, "Expect 10 total page")
	assert.Equal(t, 1, respJson.NextPage, "Expect 1 next page")
}
