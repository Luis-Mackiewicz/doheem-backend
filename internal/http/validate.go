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

func validateCPF(fl validator.FieldLevel) bool {
	d := onlyDigits(fl.Field().String())
	if len(d) != 11 {
		return false
	}

	allSame := true
	for i := 1; i < 11; i++ {
		if d[i] != d[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(d[i]-'0') * (10 - i)
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if d1 != int(d[9]-'0') {
		return false
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(d[i]-'0') * (11 - i)
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return d2 == int(d[10]-'0')
}

func validateCNPJ(fl validator.FieldLevel) bool {
	d := onlyDigits(fl.Field().String())
	if len(d) != 14 {
		return false
	}

	allSame := true
	for i := 1; i < 14; i++ {
		if d[i] != d[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}

	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(d[i]-'0') * w1[i]
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if d1 != int(d[12]-'0') {
		return false
	}

	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(d[i]-'0') * w2[i]
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return d2 == int(d[13]-'0')
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
	if len(d) != 11 {
		return false
	}
	allSame := true
	for i := 1; i < 11; i++ {
		if d[i] != d[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(d[i]-'0') * (10 - i)
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if d1 != int(d[9]-'0') {
		return false
	}
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(d[i]-'0') * (11 - i)
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return d2 == int(d[10]-'0')
}

func validateCNPJFunc(d string) bool {
	if len(d) != 14 {
		return false
	}
	allSame := true
	for i := 1; i < 14; i++ {
		if d[i] != d[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	w1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	w2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(d[i]-'0') * w1[i]
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if d1 != int(d[12]-'0') {
		return false
	}
	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(d[i]-'0') * w2[i]
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return d2 == int(d[13]-'0')
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
