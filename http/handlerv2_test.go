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

func TestSuccessHandlerV2(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandlerV2(false)
	newHandler := NewHttpHandlerV2(handlerCtx)

	message := make([]string, 0)
	message = append(message, "success msg")

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResultV2) {
		response.Data = "OK"
		response.Message = message
		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &ResponseV2{}
	_ = json.Unmarshal(body, respJson)

	assert.Equal(t, 200, resp.StatusCode, "Expect 200 status code")
	assert.Equal(t, 200, respJson.StatusCode, "Expect 200 status code in body")
	assert.Equal(t, true, respJson.Success, "Expected success True")
	assert.Equal(t, `OK`, fmt.Sprintf("%s", respJson.Data), "Expect OK data")
	assert.Equal(t, `[success msg]`, fmt.Sprintf("%s", respJson.Message), "Expect Correct Message")
}

func TestFailureHandlerV2(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandlerV2(false)
	newHandler := NewHttpHandlerV2(handlerCtx)

	var customErr = errors.New("custom error")
	var ErrResp = &ErrorResponse{
		Response: Response{
			ResponseDesc: "Custom Error",
		},
		HttpStatus: http.StatusUnprocessableEntity,
	}
	handlerCtx.AddError(customErr, ErrResp)

	message := make([]string, 0)
	message = append(message, "error msg 1")
	message = append(message, "error msg 2")

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResultV2) {
		response.Error = customErr
		response.Message = message
		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &ResponseV2{}
	_ = json.Unmarshal(body, respJson)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Expect 400 status code")
	assert.Equal(t, http.StatusUnprocessableEntity, respJson.StatusCode, "Expect 402 status code in body")
	assert.Equal(t, false, respJson.Success, "Expect Success False")
	assert.Equal(t, `[error msg 1 error msg 2]`, fmt.Sprintf("%s", respJson.Message), "Expect Correct message")
	assert.Equal(t, `[]`, fmt.Sprintf("%s", respJson.Data), "Expect EMpty data")
}

type Mydata struct {
	PageSize    int         `json:"page_size"`
	CurrentPage int         `json:"current_page"`
	TotalPage   int         `json:"total_page"`
	NextPage    int         `json:"next_page"`
	Data        interface{} `json:"data"`
}

func TestSuccessPaginationHandlerV2(t *testing.T) {
	// Start a local HTTP server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()

	handlerCtx := NewContextHandlerV2(false)
	newHandler := NewHttpHandlerV2(handlerCtx)

	testHandler := newHandler(func(w http.ResponseWriter, r *http.Request) (response HttpHandleResultV2) {
		response.Pagination = &Pagination{
			PageSize:    5,
			CurrentPage: 0,
			TotalPage:   10,
			NextPage:    1,
			TotalData:   50,
		}
		response.StatusCode = 200
		response.Data = "OK"
		return
	})

	testHandler.ServeHTTP(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	respJson := &ResponseV2{}
	_ = json.Unmarshal(body, respJson)

	my_data := Mydata{}
	dataBody, _ := json.Marshal(respJson.Data)
	_ = json.Unmarshal(dataBody, &my_data)

	assert.Equal(t, 200, resp.StatusCode, "Expect 200 status code")
	assert.Equal(t, 200, respJson.StatusCode, "Expect 200 status code in body")
	assert.Equal(t, true, respJson.Success, "Expected success True")
	assert.Equal(t, `[]`, fmt.Sprintf("%s", respJson.Message), "Expect Correct message")
	assert.Equal(t, 5, my_data.PageSize, "Expect 5 page size")
	assert.Equal(t, 0, my_data.CurrentPage, "Expect 0 current page")
	assert.Equal(t, 10, my_data.TotalPage, "Expect 10 total page")
	assert.Equal(t, 1, my_data.NextPage, "Expect 1 next page")
}
