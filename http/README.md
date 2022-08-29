# Package http

This package contains http handler related to help you write handler easier and have a standardized response.

## How to create custom handler
To create custom handler using this package, you need to declare the handler using this function signature:

```go
func(w http.ResponseWriter, r *http.Request) (result HttpHandleResult)
```

This handler requires you to return the HttpHandleResult struct with fields : `data`, `pagination`, `status code`, `error` and `is plain response`.
* `data` is the struct for the response data 
* `pagination` is the struct of Pagination that consist of `page size`, `current page`, `total page` and `next page`
* `status code` will be override if set, default is 200. 
* `error` is the error (if any)
* `isPlainResponse` is option to return custom response instead Default Success Response

Example:
```go
func HelloHandler(w http.ResponseWriter, r *http.Request) (result golib.HttpHandlerResult) {
    data, err = someProcess("john", "doe")
    
    pagination = somePagingProcess(data)
    
    statusCode = 202
    
    return golib.HttpHandlerResult {
        Data: data
        StatusCode: statusCode
        Pagination: pagination
        Error: err
        IsPlainResponse: false
    }
}

func main() {
	handlerCtx := phttp.NewContextHandler(false) //IsDebug if true will print Request body

	// add custom error
	var ErrCustom *ErrorResponse = &ErrorResponse{
		Response: Response{
			ResponseDesc: "Custom message",
		},
		HttpStatus: http.StatusInternalServerError,
	}
	handlerCtx.AddError(errors.New("custom error"), ErrCustom)

    // newHandler is a function that will create function for create new custom handler with injected handler context
    newHandler := phttp.NewHttpHandler(handlerCtx)
    
    // helloHandler is the handler
	helloHandler := newHandler(HelloHandler)

	router := chi.NewRouter()
	router.Get("/hello", helloHandler.ServeHTTP)

	http.ListenAndServe(":5678", router)
}
```