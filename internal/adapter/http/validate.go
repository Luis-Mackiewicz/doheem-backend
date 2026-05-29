package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return fld.Name
		}
		return name
	})
}

type validationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
}

func decodeAndValidate(r *http.Request, req any) []validationError {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return []validationError{{Field: "body", Tag: "invalid_json"}}
	}
	if err := validate.Struct(req); err != nil {
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			out := make([]validationError, len(errs))
			for i, fe := range errs {
				out[i] = validationError{
					Field: fe.Field(),
					Tag:   fe.Tag(),
				}
			}
			return out
		}
		return []validationError{{Field: "body", Tag: "invalid"}}
	}
	return nil
}
