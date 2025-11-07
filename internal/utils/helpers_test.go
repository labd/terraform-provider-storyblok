package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestInterfacePointerToInt64_ValidInt64(t *testing.T) {
	var input interface{} = int64(42)
	result, err := InterfacePointerToInt64(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ValueInt64() != 42 {
		t.Errorf("expected 42, got %d", result.ValueInt64())
	}
}

func TestInterfacePointerToInt64_ValidInt(t *testing.T) {
	var input interface{} = int(42)
	result, err := InterfacePointerToInt64(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ValueInt64() != 42 {
		t.Errorf("expected 42, got %d", result.ValueInt64())
	}
}

func TestInterfacePointerToInt64_ValidFloat64(t *testing.T) {
	var input interface{} = float64(42.9)
	result, err := InterfacePointerToInt64(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ValueInt64() != 42 {
		t.Errorf("expected 42, got %d", result.ValueInt64())
	}
}

func TestInterfacePointerToInt64_ValidString(t *testing.T) {
	var input interface{} = "42"
	result, err := InterfacePointerToInt64(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ValueInt64() != 42 {
		t.Errorf("expected 42, got %d", result.ValueInt64())
	}
}

func TestInterfacePointerToInt64_InvalidString(t *testing.T) {
	var input interface{} = "invalid"
	_, err := InterfacePointerToInt64(&input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestInterfacePointerToInt64_UnsupportedType(t *testing.T) {
	var input interface{} = []int{1, 2, 3}
	_, err := InterfacePointerToInt64(&input)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestInterfacePointerToInt64_NilInput(t *testing.T) {
	var input *interface{} = nil
	result, err := InterfacePointerToInt64(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsNull() {
		t.Errorf("expected null result, got %v", result)
	}
}

func TestInt64ToStringInterfacePointer_ValidInt64(t *testing.T) {
	input := types.Int64Value(42)
	result := Int64ToStringInterfacePointer(input)
	if result == nil {
		t.Fatal("expected non-nil result, got nil")
	}
	if *result != "42" {
		t.Errorf("expected \"42\", got %v", *result)
	}
}

func TestInt64ToStringInterfacePointer_NullInput(t *testing.T) {
	input := types.Int64PointerValue(nil)
	result := Int64ToStringInterfacePointer(input)
	if result != nil {
		t.Errorf("expected nil result, got %v", *result)
	}
}
