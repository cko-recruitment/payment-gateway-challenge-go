package validator

import (
	"log"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// validateNumberOfDigits checks of number of digits in an int are equal to param passed in validator
func validateNumberOfDigits(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n == param
}

// validateNumberOfDigitsGte checks of number of digits in an int are greater than or equal to param passed in validator
func validateNumberOfDigitsGte(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n >= param
}

// validateNumberOfDigitsLte checks of number of digits in an int are less than or equal to param passed in validator
func validateNumberOfDigitsLte(fl validator.FieldLevel) bool {
	n := getNumberOfDigits(fl.Field().Int())
	param := convertParamToInt(fl.Param())
	return n <= param
}

// validateDateInFuture checks if date of format mm/yyyy is in the future
func validateDateInFuture(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	date, err := time.Parse("01/2006", v)
	if err != nil {
		return false
	}
	return time.Now().Before(date)
}

// validateMonthAndYearInFuture checks if year and month combination are in the future
func validateMonthAndYearInFuture(fl validator.FieldLevel) bool {
	v := fl.Field().Int()
	m := fl.Parent().FieldByName("ExpiryMonth").Int()
	date := time.Date(int(v), time.Month(m), 1, 0, 0, 0, 0, time.UTC)
	return time.Now().Before(date)
}

// getNumberOfDigits returns the number of digits in an int
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

// convertParamToInt converts a string to an int
func convertParamToInt(s string) int {
	param, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("error occured while converting custom validation param to int", err)
	}
	return param
}
