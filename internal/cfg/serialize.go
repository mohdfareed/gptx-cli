package cfg

import (
	"encoding/json"
	"fmt"
)

// TODO: Refactor to delete this file.

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
		strValue := fmt.Sprintf("%v", value)

		// Add to result map if not empty
		if value != nil && strValue != "" {
			result[envKey] = strValue
		}
	}
	return result
}
