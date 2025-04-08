package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// // GetQuery assigns query values from the context to the provided pointers,
// // based on the provided query parameter names and pointers.
// func GetQuery(c *gin.Context, nameAndArgs ...any) error {
// 	if len(nameAndArgs)%2 != 0 {
// 		return errors.New("must pass pairs of query name and pointer")
// 	}

// 	for i := 0; i < len(nameAndArgs); i += 2 {
// 		fieldName, ok := nameAndArgs[i].(string)
// 		if !ok {
// 			return errors.New("query parameter name must be a string")
// 		}

// 		arg := nameAndArgs[i+1]
// 		val := reflect.ValueOf(arg)
// 		if val.Kind() != reflect.Ptr || val.IsNil() {
// 			return fmt.Errorf("must pass a pointer for %s", fieldName)
// 		}

// 		valElem := val.Elem()
// 		// Fetch the query value based on the parameter name
// 		queryValue, isOk := c.GetQuery(fieldName)
// 		if !isOk {
// 			return fmt.Errorf("%s(%s) is required", fieldName, valElem.Kind())
// 		}

// 		// Set the appropriate type based on the kind of the argument

// 		switch valElem.Kind() {
// 		case reflect.Bool:
// 			parsedBool, err := strconv.ParseBool(queryValue)
// 			if err != nil {
// 				return fmt.Errorf("%s must be a boolean", fieldName)
// 			}
// 			valElem.SetBool(parsedBool)
// 		case reflect.String:
// 			valElem.SetString(queryValue)
// 		case reflect.Int:
// 			parsedInt, err := strconv.Atoi(queryValue)
// 			if err != nil {
// 				return fmt.Errorf("%s must be an integer", fieldName)
// 			}
// 			valElem.SetInt(int64(parsedInt))
// 		default:
// 			return fmt.Errorf("%s has unsupported type", fieldName)
// 		}
// 	}
// 	return nil
// }

type ParamConstraint interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		bool |
		string | uuid.UUID
}

func ParseParam[T ParamConstraint](key, paramValue string, validator func(T) error) (T, error) {
	var output T
	var err error
	switch any(output).(type) {
	case bool:
		var parsed bool
		parsed, err = strconv.ParseBool(paramValue)
		output = any(parsed).(T)

	case string:
		output = any(paramValue).(T)

	case int:
		var parsed int64
		parsed, err = strconv.ParseInt(paramValue, 10, 0)
		output = any(parsed).(T)

	case int8:
		var parsed int64
		parsed, err = strconv.ParseInt(paramValue, 10, 8)
		output = any(parsed).(T)

	case int16:
		var parsed int64
		parsed, err = strconv.ParseInt(paramValue, 10, 16)
		output = any(parsed).(T)

	case int32:
		var parsed int64
		parsed, err = strconv.ParseInt(paramValue, 10, 32)
		output = any(parsed).(T)

	case int64:
		var parsed int64
		parsed, err = strconv.ParseInt(paramValue, 10, 64)
		output = any(parsed).(T)

	case uint:
		var parsed uint64
		parsed, err = strconv.ParseUint(paramValue, 10, 0)
		output = any(parsed).(T)

	case uint8:
		var parsed uint64
		parsed, err = strconv.ParseUint(paramValue, 10, 8)
		output = any(parsed).(T)

	case uint16:
		var parsed uint64
		parsed, err = strconv.ParseUint(paramValue, 10, 16)
		output = any(parsed).(T)

	case uint32:
		var parsed uint64
		parsed, err = strconv.ParseUint(paramValue, 10, 32)
		output = any(parsed).(T)

	case uint64:
		var parsed uint64
		parsed, err = strconv.ParseUint(paramValue, 10, 64)
		output = any(parsed).(T)
	case float32:
		var parsed float64
		parsed, err = strconv.ParseFloat(paramValue, 32)
		output = any(parsed).(T)
	case float64:
		var parsed float64
		parsed, err = strconv.ParseFloat(paramValue, 64)
		output = any(parsed).(T)
	case uuid.UUID:
		var parsed uuid.UUID
		parsed, err = uuid.Parse(paramValue)
		output = any(parsed).(T)

	}
	if err != nil {
		return output, fmt.Errorf("path parameter %s is invalid: %w", key, err)
	}
	// Validate the parsed value
	if validator != nil {
		if err := validator(output); err != nil {
			return output, err
		}
	}
	return output, nil
}

func GetPathParam[T ParamConstraint](c *gin.Context, key string, validator func(T) error) (T, error) {

	paramValue, ok := c.Params.Get(key)
	if !ok {
		var output T
		return output, fmt.Errorf("path parameter %s is required", key)
	}

	return ParseParam(key, paramValue, validator)

}
func GetQueryParam[T ParamConstraint](c *gin.Context, key string, validator func(T) error) (T, error) {

	paramValue, ok := c.GetQuery(key)
	if !ok {
		var output T
		return output, fmt.Errorf("Query parameter %s is required", key)
	}

	return ParseParam(key, paramValue, validator)

}

func SendError(c *gin.Context, err error, code int) {
	if err == nil {
		c.JSON(code, Response{})
		return
	}
	c.JSON(code, Response{
		Error: err.Error(),
	})
}
func SendJson(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Data: data,
	})
}
