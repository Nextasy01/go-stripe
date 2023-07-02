package validators

import (
	"log"

	"github.com/go-playground/validator/v10"
)

func ValidateCurrency(price interface{}) error {
	validate := validator.New()

	validate.RegisterValidation("currency", Validate)

	err := validate.Struct(price)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Println(err)
		}
		return err
	}
	return nil
}

func Validate(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return currency == "usd" || currency == "eu"
}
