package validator

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func validateNumberOfDigits(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n == param
}

func validateNumberOfDigitsGte(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n >= param
}

func validateNumberOfDigitsLte(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n <= param
}

func validateDateInFuture(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	date, err := time.Parse("01/2006", v)
	if err != nil {
		return false
	}
	return time.Now().Before(date)
}

func validateMonthAndYearInFuture(fl validator.FieldLevel) bool {
	v := fl.Field().Int()
	m := fl.Parent().FieldByName("ExpiryMonth").Int()
	date := time.Date(int(v), time.Month(m), 1, 0, 0, 0, 0, time.UTC)
	return time.Now().Before(date)
}

func getNumberOfDigits(v int64) int {
	if v < 0 {
		return 0
	}
	n := 0
	for ; v > 0; v /= 10 {
		n += 1
	}
	return n
}

func convertParamToInt(s string) int {
	param, err := strconv.Atoi(s)
	if err != nil {
		panic(err.Error())
	}
	return param
}
