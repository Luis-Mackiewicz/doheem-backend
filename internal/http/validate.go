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

	validate.RegisterValidation("cpf", validateCPF)
	validate.RegisterValidation("cnpj", validateCNPJ)
	validate.RegisterValidation("document", validateDocument)
	validate.RegisterValidation("phone_br", validatePhoneBR)
	validate.RegisterValidation("cep_br", validateCEPBR)
}

func allDigitsSame(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return false
		}
	}
	return true
}

func checkDigit(digits string, weights []int) int {
	sum := 0
	for i, w := range weights {
		sum += int(digits[i]-'0') * w
	}
	if rem := sum % 11; rem >= 2 {
		return 11 - rem
	}
	return 0
}

func validateCPF(fl validator.FieldLevel) bool {
	return validateCPFFunc(onlyDigits(fl.Field().String()))
}

func validateCNPJ(fl validator.FieldLevel) bool {
	return validateCNPJFunc(onlyDigits(fl.Field().String()))
}

func validateDocument(fl validator.FieldLevel) bool {
	d := onlyDigits(fl.Field().String())
	switch len(d) {
	case 11:
		return validateCPFFunc(d)
	case 14:
		return validateCNPJFunc(d)
	default:
		return false
	}
}

func validateCPFFunc(d string) bool {
	if len(d) != 11 || allDigitsSame(d) {
		return false
	}
	return checkDigit(d, []int{10, 9, 8, 7, 6, 5, 4, 3, 2}) == int(d[9]-'0') &&
		checkDigit(d, []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2}) == int(d[10]-'0')
}

func validateCNPJFunc(d string) bool {
	if len(d) != 14 || allDigitsSame(d) {
		return false
	}
	return checkDigit(d, []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}) == int(d[12]-'0') &&
		checkDigit(d, []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}) == int(d[13]-'0')
}

func validatePhoneBR(fl validator.FieldLevel) bool {
	d := onlyDigits(fl.Field().String())
	if len(d) != 10 && len(d) != 11 {
		return false
	}
	areaCode := int(d[0]-'0')*10 + int(d[1]-'0')
	return areaCode >= 11 && areaCode <= 99
}

func validateCEPBR(fl validator.FieldLevel) bool {
	d := onlyDigits(fl.Field().String())
	return len(d) == 8
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
