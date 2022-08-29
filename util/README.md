# Package util

This package contains any functionality that can be called as common usage and etc.

## util.NullTime
Use this as wrapper for `sql.NullTime` with predefined `null` return value when marshalling JSON.
No need to validate null value from your DB scan result.

Usage in struct :

```go
type MyData struct {
	UpdatedAt util.NullTime `json:"updated_at"`
}
```

Assign a null value :

```go
data.UpdatedAt = util.NullTime{time.Time{}, true}
```

## util.NullString
Use this as wrapper for `sql.NullString` with predefined `null` return value when marshalling JSON.
No need to validate null value from your DB scan result.

Usage in struct :

```go
type MyData struct {
	Name util.NullString `json:"name"`
}
```

Assign a null value :

```go
data.Name = util.NullString{"", true}
```