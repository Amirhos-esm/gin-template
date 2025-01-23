package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetQuery assigns query values from the context to the provided pointers,
// based on the provided query parameter names and pointers.
func GetQuery(c *gin.Context, nameAndArgs ...any) error {
	if len(nameAndArgs)%2 != 0 {
		return errors.New("must pass pairs of query name and pointer")
	}

	for i := 0; i < len(nameAndArgs); i += 2 {
		fieldName, ok := nameAndArgs[i].(string)
		if !ok {
			return errors.New("query parameter name must be a string")
		}

		arg := nameAndArgs[i+1]
		val := reflect.ValueOf(arg)
		if val.Kind() != reflect.Ptr || val.IsNil() {
			return fmt.Errorf("must pass a pointer for %s", fieldName)
		}
		

		valElem := val.Elem()
		// Fetch the query value based on the parameter name
		queryValue, isOk := c.GetQuery(fieldName)
		if !isOk {
			return fmt.Errorf("%s(%s) is required", fieldName,valElem.Kind())
		}

		// Set the appropriate type based on the kind of the argument
		
		switch valElem.Kind() {
		case reflect.Bool:
			parsedBool, err := strconv.ParseBool(queryValue)
			if err != nil {
				return fmt.Errorf("%s must be a boolean", fieldName)
			}
			valElem.SetBool(parsedBool)
		case reflect.String:
			valElem.SetString(queryValue)
		case reflect.Int:
			parsedInt, err := strconv.Atoi(queryValue)
			if err != nil {
				return fmt.Errorf("%s must be an integer", fieldName)
			}
			valElem.SetInt(int64(parsedInt))
		default:
			return fmt.Errorf("%s has unsupported type", fieldName)
		}
	}
	return nil
}

