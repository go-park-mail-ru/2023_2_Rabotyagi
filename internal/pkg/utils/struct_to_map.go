package utils

import "reflect"

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	typeOf := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typeOf.Field(i).Name
		fieldName = CamelToSnakeCase(fieldName)

		if field.IsZero() && field.Type().Kind() != reflect.Bool {
			continue
		}

		result[fieldName] = field.Interface()
	}

	return result
}
