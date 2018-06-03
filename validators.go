package main

import (
	"log"
	"reflect"

	validator "gopkg.in/go-playground/validator.v8"
)

func allowedType(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {

	log.Println(field.Interface())

	if userType, ok := field.Interface().(string); ok {
		if userType == "user" || userType == "company" {
			return true
		}
	}

	return false
}
