package http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HandlerV2Option func(*HttpHandlerV2)

type HttpHandlerV2 struct {
	// H is handler, with return interface{} as data object, error for error type
	H func(w http.ResponseWriter, r *http.Request) HttpHandleResultV2
	CustomWriterV2
	IsDebug bool
}

func NewHttpHandlerV2(c HandlerContextV2, opts ...HandlerV2Option) func(handler func(w http.ResponseWriter, r *http.Request) HttpHandleResultV2) HttpHandlerV2 {
	return func(handler func(w http.ResponseWriter, r *http.Request) HttpHandleResultV2) HttpHandlerV2 {
		h := HttpHandlerV2{H: handler, CustomWriterV2: CustomWriterV2{C: c}, IsDebug: c.IsDebug}

		for _, opt := range opts {
			opt(&h)
		}

		return h
	}
}

func (h HttpHandlerV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result := h.H(w, r)

	if h.IsDebug {
		// Read the content
		var bodyBytes []byte
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body.Close() //  must close
		// Restore the io.ReadCloser to its original state
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		// Use the content
		bodyString := string(bodyBytes)
		log.Logger.Info().Msgf("[DEBUG] Request: %v", bodyString)
	}

	if result.Error != nil {
		log.Logger.Error().Err(result.Error).Msgf("Response: %+v", result.Data)
		h.WriteError(w, result.Error, result.Message)
		return
	}

	if result.IsPlainResponse {
		h.WritePlain(w, result.Data, result.StatusCode)
	} else {
		h.Write(w, result.Data, result.StatusCode, result.Pagination, result.Message)
	}
}
