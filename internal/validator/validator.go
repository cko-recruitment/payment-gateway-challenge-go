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

func NewValidator() *Validator {
	once.Do(func() {
		instance = &Validator{
			validate: validator.New(validator.WithRequiredStructEnabled()),
		}
		instance.registerValidations()
	})

	return instance
}

func GetValidator() *Validator {
	return instance
}

func (v *Validator) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}

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
