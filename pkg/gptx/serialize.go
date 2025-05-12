package gptx

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// TODO: Find a library to handle any -> map[string]string conversion

// ToEnvMap converts a Config to a map of environment variables
// using reflection and struct tags.
func (c *Config) ToEnvMap() map[string]string {
	envMap := make(map[string]string)

	// Use reflection to iterate through struct fields
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip if field doesn't have env tag
		envKey, hasEnv := field.Tag.Lookup("env")
		if !hasEnv {
			continue
		}

		// Get default value (if any)
		defaultVal, hasDefault := field.Tag.Lookup("default")

		// Format value according to its type
		formatted, isDefault := formatValue(value, defaultVal, hasDefault, envKey)
		if formatted == "" || isDefault {
			continue
		}

		// Add to environment map
		envMap[EnvVar(envKey)] = formatted
	}

	return envMap
}

// formatValue converts a reflect.Value to a string representation
// and checks if it's a default value that should be skipped.
func formatValue(value reflect.Value, defaultVal string, hasDefault bool, envKey string) (string, bool) {
	var strValue string
	isDefaultValue := false

	switch value.Kind() {
	case reflect.String:
		strValue = value.String()
		isDefaultValue = hasDefault && strValue == defaultVal

		// Special case for system prompt
		if envKey == "INSTRUCTIONS" && strValue == fmt.Sprintf(SYS_PROMPT, AppName) {
			isDefaultValue = true
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal := value.Int()
		strValue = fmt.Sprintf("%d", intVal)

		if hasDefault {
			defInt, err := strconv.ParseInt(defaultVal, 10, 64)
			if err == nil && intVal == defInt {
				isDefaultValue = true
			}
		}

	case reflect.Ptr:
		if value.IsNil() {
			return "", true
		}

		// Handle pointer types
		switch value.Elem().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			strValue = fmt.Sprintf("%d", value.Elem().Int())
		case reflect.String:
			strValue = value.Elem().String()
		}

	case reflect.Slice:
		if value.Len() == 0 {
			return "", true
		}

		// Handle string slices
		if value.Type().Elem().Kind() == reflect.String {
			elements := make([]string, value.Len())
			for j := 0; j < value.Len(); j++ {
				elements[j] = strings.ReplaceAll(value.Index(j).String(), ",", "\\,")
			}
			strValue = strings.Join(elements, ",")
		}
	}

	return strValue, isDefaultValue
}
