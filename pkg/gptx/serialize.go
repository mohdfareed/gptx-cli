package gptx

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ToEnvMap converts a Config to a map of environment variables
// using JSON as an intermediate representation.
func (c *Config) ToEnvMap() map[string]string {
	// First convert struct to a map using JSON marshal/unmarshal
	jsonData, err := json.Marshal(c)
	if err != nil {
		return map[string]string{} // Return empty map on error
	}

	// Unmarshal into a map
	var dataMap map[string]any
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return map[string]string{} // Return empty map on error
	}

	// Process each field in the struct directly
	result := make(map[string]string)
	for key, value := range dataMap {
		// Get the environment variable name
		envKey := EnvVar(c, key)

		// Convert the value to a string
		var strValue string
		switch v := value.(type) {
		case string:
			strValue = v
		case bool:
			strValue = strconv.FormatBool(v)
		case int:
			strValue = strconv.Itoa(v)
		case float64:
			strValue = strconv.FormatFloat(v, 'f', -1, 64)
		case []any:
			strSlice := make([]string, len(v))
			for i, item := range v {
				switch item := item.(type) {
				case string:
					strSlice[i] = item
				case bool:
					strSlice[i] = strconv.FormatBool(item)
				case int:
					strSlice[i] = strconv.Itoa(item)
				case float64:
					strSlice[i] = strconv.FormatFloat(item, 'f', -1, 64)
				default:
					strSlice[i] = fmt.Sprintf("%v", item) // Fallback to default string conversion
				}
			}
			strValue = strings.Join(strSlice, ",")
		}

		// Add to result map if not empty
		if strValue != "" {
			result[envKey] = strValue
		}
	}
	return result
}
