// utils/json.go
package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ParseJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("request body kosong")
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(v)
}

func ParseAndValidate(r *http.Request, v interface{}) error {
	if err := ParseJSON(r, v); err != nil {
		return err
	}
	return validate.Struct(v)
}

func ValidateStruct(v interface{}) error {
	return validate.Struct(v)
}