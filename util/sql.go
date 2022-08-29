package util

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"log"
	"reflect"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//NullString is a wrapper around sql.NullString
type NullString sql.NullString

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

func (x NullString) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}

	return x.String, nil
}

func (x *NullString) Scan(v interface{}) error {
	if v == nil {
		*x = NullString{"", false}
		return nil
	}
	var i sql.NullString
	if err := i.Scan(v); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(v) == nil {
		*x = NullString{i.String, false}
	} else {
		*x = NullString{i.String, true}
	}
	return nil
}

//NullInt is a wrapper around sql.NullInt
type NullInt16 sql.NullInt16

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullInt16
func (x *NullInt16) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Int16)
}

func (x NullInt16) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}

	return x.Int16, nil
}

func (x *NullInt16) Scan(v interface{}) error {
	if v == nil {
		*x = NullInt16{0, false}
		return nil
	}
	var i sql.NullInt16
	if err := i.Scan(v); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(v) == nil {
		*x = NullInt16{i.Int16, false}
	} else {
		*x = NullInt16{i.Int16, true}
	}
	return nil
}

type NullTime sql.NullTime

func (x *NullTime) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Time)
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

func (x *NullTime) Scan(v interface{}) error {
	if v == nil {
		*x = NullTime{time.Time{}, false}
		return nil
	}
	var i sql.NullTime
	if err := i.Scan(v); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(v) == nil {
		*x = NullTime{i.Time, false}
	} else {
		*x = NullTime{i.Time, true}
	}
	return nil
}

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func ComparePassword(pwd []byte, hash []byte) error {
	err := bcrypt.CompareHashAndPassword(hash, pwd)
	return err
}
