package db

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// ReadJSON reads and unmarshals JSON data from a file into the provided object.
func Read(filename string, v interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// WriteJSON marshals the given object and writes it to a file.
func Write(filename string, v interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

// VerifyJSON compares the JSON content of a file with the expected data.
// If they don't match, it returns an error explaining the difference.
func Check(filename string, expected interface{}) (bool, error) {
	var actual map[string]interface{}
	if err := Read(filename, &actual); err != nil {
		return false, fmt.Errorf("failed to read JSON file: %w", err)
	}

	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return false, fmt.Errorf("failed to marshal expected data: %w", err)
	}

	var expectedMap map[string]interface{}
	if err := json.Unmarshal(expectedBytes, &expectedMap); err != nil {
		return false, fmt.Errorf("failed to unmarshal expected data: %w", err)
	}

	return compareJSON(actual, expectedMap)
}

// compareJSON compares two JSON objects and returns an error if they don't match.
func compareJSON(actual, expected map[string]interface{}) (bool, error) {
	for key, expectedValue := range expected {
		actualValue, exists := actual[key]
		if !exists {
			return false, fmt.Errorf("missing key '%s' in actual JSON", key)
		}

		if !compareValues(actualValue, expectedValue) {
			return false, fmt.Errorf("value mismatch for key '%s': expected '%v', but got '%v'", key, reflect.TypeOf(expectedValue), reflect.TypeOf(actualValue))

		}
	}

	for key := range actual {
		if _, exists := expected[key]; !exists {
			return false, fmt.Errorf("unexpected key '%s' found in actual JSON", key)
		}
	}

	return true, nil
}

// compareValues performs a deep comparison of two values.
func compareValues(actual, expected interface{}) bool {
	return reflect.TypeOf(actual) == reflect.TypeOf(expected)
}
