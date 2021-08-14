package helpers

import (
	"reflect"
	"testing"
)

func TestUnderscore(t *testing.T) {
	tests := []struct {
		value  string
		result string
	}{
		{"SimpleTestCamelCase", "simple_test_camel_case"},
		{"simpletestdowncase", "simpletestdowncase"},
		{"SIMPLETESTUPPERCASE", "simpletestuppercase"},
		{"TESTCombined", "test_combined"},
	}

	for _, test := range tests {
		result := Underscore(test.value)
		if !reflect.DeepEqual(test.result, result) {
			t.Errorf("received value %s is not equal to expected value: %s", result, test.result)
		}
	}
}

func TestUndescoreUppercased(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"SimpleTestCamelCase", "SIMPLE_TEST_CAMEL_CASE"},
		{"simpletestdowncase", "SIMPLETESTDOWNCASE"},
		{"SIMPLETESTUPPERCASE", "SIMPLETESTUPPERCASE"},
		{"TESTCombined", "TEST_COMBINED"},
	}

	for _, test := range tests {
		result := UndescoreUppercased(test.value)
		if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("received value %s is not equal to expected value: %s", result, test.expected)
		}
	}
}

func TestCamelCase(t *testing.T) {
	tests := []struct {
		value  string
		result string
	}{
		{"simple_test_camel_case", "SimpleTestCamelCase"},
		{"simpletestdowncase", "Simpletestdowncase"},
		{"test_combined", "TestCombined"},
	}

	for _, test := range tests {
		result := CamelCase(test.value)
		if !reflect.DeepEqual(test.result, result) {
			t.Errorf("received value %s is not equal to expected value: %s", result, test.result)
		}
	}
}
