package gorm_validator

import (
	validator_v10 "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func New(tx *gorm.DB) *validator_v10.Validate {
	validator := validator_v10.New()
	validator.RegisterValidation("uniqueness", validateUniqueness(tx))
	return validator
}

func validateUniqueness(tx *gorm.DB) validator_v10.Func {
	return func(fl validator_v10.FieldLevel) bool {
		field := fl.Field().Interface()
		fieldName := fl.FieldName()
		id := fl.Parent().FieldByName("ID").Interface()
		model := fl.Parent().Interface()

		var count int64
		result := tx.Model(&model).Where(fieldName+" = ?", field).Where("id != ?", id).Limit(1).Count(&count)

		if err := result.Error; err != nil {
			return false
		}

		return count == 0
	}
}
