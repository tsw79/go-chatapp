package util

import (
	"fmt"
	"os"
	"strconv"
)

type EnvironmentUtil struct {
	val interface{}
}

func Environment() *EnvironmentUtil {
	return &EnvironmentUtil{}
}

// Helper func
// Reads an environment or return a default (string) value
// 	Note: Instead of passing a default type as int, bool, etc.which  will add
// 		some overhead processing time for type assertions and type switches.
// 		Therefore, by design, all `defaultVal` paramaters are of string type,
// 		with its proper data type value inside a string.
func (e *EnvironmentUtil) GetEnv(key string, defaultVal string) *EnvironmentUtil {
	if value, exists := os.LookupEnv(key); exists {
		e.val = value
	} else {
		e.val = defaultVal
	}
	return e
}

func (e *EnvironmentUtil) AsString() string {
	return fmt.Sprintf("%v", e.val)
}

func (e *EnvironmentUtil) AsInt() int {
	// Use Atoi (ASCII to integer) function to convert value
	result, _ := strconv.Atoi(e.AsString())
	return result

	// val := e.val.(int)
	// intVal, _ := val.(int)
	// return intVal
}

func (e *EnvironmentUtil) AsBool() bool {
	result, _ := strconv.ParseBool(e.AsString())
	return result
}

// // Helper func
// // Reads an environment or return a default (string) value
// func (e *EnvironmentUtil) GetEnv(key string, defaultVal string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultVal
// }

// // Helper func
// // Reads an environment variable into integer or return a default value
// func (e *EnvironmentUtil) GetEnvAsInt(key string, defaultVal int) int {
// 	valueStr := e.GetEnv(key, "")
// 	// Use Atoi (ASCII to integer) function to convert value
// 	if value, err := strconv.Atoi(valueStr); err == nil {
// 		return value
// 	}
// 	return defaultVal
// }

// // Helper to read an environment variable into a bool or return default value
// func getEnvAsBool(name string, defaultVal bool) bool {
// 	valStr := getEnv(name, "")
// 	if val, err := strconv.ParseBool(valStr); err == nil {
// return val
// 	}

// 	return defaultVal
// }

// // Helper to read an environment variable into a string slice or return default value
// func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
// 	valStr := getEnv(name, "")

// 	if valStr == "" {
// return defaultVal
// 	}

// 	val := strings.Split(valStr, sep)

// 	return val
// }
