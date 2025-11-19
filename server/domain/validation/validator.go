package validation

import (
	"portfolio/domain"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type Validator struct {
	errors []*domain.DomainError
}

func NewValidator() *Validator {
	return &Validator{
		errors: make([]*domain.DomainError, 0),
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) Errors() []*domain.DomainError {
	return v.errors
}

func (v *Validator) FirstError() *domain.DomainError {
	if len(v.errors) == 0 {
		return nil
	}
	return v.errors[0]
}

func (v *Validator) Required(fieldName string, value interface{}) *Validator {
	if isEmpty(value) {
		v.errors = append(v.errors, domain.NewRequiredFieldError(fieldName))
	}
	return v
}

func (v *Validator) MinLength(fieldName string, value string, min int) *Validator {
	if utf8.RuneCountInString(value) < min {
		v.errors = append(v.errors, domain.NewValidationError(
			"Value is too short", fieldName,
			nil,
		))
	}
	return v
}

func (v *Validator) MaxLength(fieldName string, value string, max int) *Validator {
	if utf8.RuneCountInString(value) > max {
		v.errors = append(v.errors, domain.NewValidationError(
			"Value is too long", fieldName,
			nil,
		))
	}
	return v
}

func (v *Validator) Email(fieldName string, value string) *Validator {
	if value == "" {
		return v
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		v.errors = append(v.errors, domain.NewInvalidFormatError(fieldName, "valid email address"))
	}
	return v
}

func (v *Validator) URL(fieldName string, value string) *Validator {
	if value == "" {
		return v
	}

	urlRegex := regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
	if !urlRegex.MatchString(value) {
		v.errors = append(v.errors, domain.NewInvalidFormatError(fieldName, "valid URL"))
	}
	return v
}

func (v *Validator) Phone(fieldName string, value string) *Validator {
	if value == "" {
		return v
	}

	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-\(\)\.\+]{8,20}$`)
	if !phoneRegex.MatchString(value) {
		v.errors = append(v.errors, domain.NewInvalidFormatError(fieldName, "valid phone number"))
	}
	return v
}

func (v *Validator) DateNotFuture(fieldName string, value time.Time) *Validator {
	if value.After(time.Now()) {
		v.errors = append(v.errors, domain.NewValidationError(
			"Date cannot be in the future", fieldName,
			nil,
		))
	}
	return v
}

func (v *Validator) Custom(fieldName string, condition bool, message string) *Validator {
	if !condition {
		v.errors = append(v.errors, domain.NewValidationError(message, fieldName, nil))
	}
	return v
}

func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) == ""
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr:
		return v.IsNil()
	case reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func ValidateStruct(s interface{}) []*domain.DomainError {
	validator := NewValidator()

	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return []*domain.DomainError{
			domain.NewValidationError("Expected a struct", "", nil),
		}
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		rules := strings.Split(tag, ",")
		for _, rule := range rules {
			applyValidationRule(validator, field.Name, fieldValue.Interface(), strings.TrimSpace(rule))
		}
	}

	return validator.Errors()
}

func applyValidationRule(validator *Validator, fieldName string, value interface{}, rule string) {
	switch {
	case rule == "required":
		validator.Required(fieldName, value)
	case strings.HasPrefix(rule, "min="):
	case strings.HasPrefix(rule, "max="):
	case rule == "email":
		if str, ok := value.(string); ok {
			validator.Email(fieldName, str)
		}
	case rule == "url":
		if str, ok := value.(string); ok {
			validator.URL(fieldName, str)
		}
	case rule == "phone":
		if str, ok := value.(string); ok {
			validator.Phone(fieldName, str)
		}
	}
}
