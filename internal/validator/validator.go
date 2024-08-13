package validator

import (
	"log"
	"sync"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

var (
	instance *Validator
	once     sync.Once
)

// NewValidator initializes a validator as a singleton
func NewValidator() *Validator {
	once.Do(func() {
		instance = &Validator{
			validate: validator.New(validator.WithRequiredStructEnabled()),
		}
		instance.registerValidations()
	})

	return instance
}

// GetValidator returns an instance of the validator
func GetValidator() *Validator {
	return instance
}

// ValidateStruct validates any struct according to validation rules embedded in the struct
func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

// registerValidations registers custom written validations
func (v *Validator) registerValidations() {
	err := instance.validate.RegisterValidation("int_len", validateNumberOfDigits)
	if err != nil {
		log.Fatal("Failed to register validation 'int_len'", err)
	}
	err = instance.validate.RegisterValidation("int_max_len", validateNumberOfDigitsLte)
	if err != nil {
		log.Fatal("Failed to register validation 'int_max_len'", err)
	}
	err = instance.validate.RegisterValidation("int_min_len", validateNumberOfDigitsGte)
	if err != nil {
		log.Fatal("Failed to register validation 'int_min_len'", err)
	}
	err = instance.validate.RegisterValidation("str_date_gt", validateDateInFuture)
	if err != nil {
		log.Fatal("Failed to register validation 'str_date_gt'", err)
	}
	err = instance.validate.RegisterValidation("month_year_gt", validateMonthAndYearInFuture)
	if err != nil {
		log.Fatal("Failed to register validation 'month_year_gt'", err)
	}
}
