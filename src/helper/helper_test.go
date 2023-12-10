package helper_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Cosmin2410/proxy-backend-test/src/helper"
)

func TestRandomAttribute_Object(t *testing.T) {
	input := []byte(`{"name":"test"}`)
	expected := map[string]interface{}{
		"name": "test",
		"foo":  "bar",
	}

	output, err := helper.AddRandomAttribute(input)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatal("Expected valid JSON, got:", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Error("Result and Expected not the same")
	}
}

func TestAddRandomAttribute_Array(t *testing.T) {
	input := []byte(`[{"name":"test1"},{"name":"test2"}]`)
	expected := []interface{}{
		map[string]interface{}{"name": "test1", "foo": "bar"},
		map[string]interface{}{"name": "test2", "foo": "bar"},
	}

	output, err := helper.AddRandomAttribute(input)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	var result []interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		t.Fatal("Expected valid JSON, got:", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Error("Result and Expected not the same")
	}
}

func TestAddRandomAttribute_InvalidJSON(t *testing.T) {
	input := []byte(`{"name":"test"`)

	_, err := helper.AddRandomAttribute(input)
	if err == nil {
		t.Fatal("Expected an error for invalid JSON, got none")
	}
}

func TestAddRandomAttribute_NotObjectOrArray(t *testing.T) {
	input := []byte(`"just a string"`)
	expected := input

	output, err := helper.AddRandomAttribute(input)
	if err != nil {
		t.Fatal("Expected no error, got:", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Error("Output and expected not the same")
	}
}
