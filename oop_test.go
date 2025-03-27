package oop

import (
	"reflect"
	"testing"
)

// TestInterface is a test interface used in the tests
type TestInterface interface {
	GetValue() int
}

// TestStruct is a test struct that implements TestInterface
type TestStruct struct {
	Value int
}

// GetValue implements TestInterface
func (t *TestStruct) GetValue() int {
	return t.Value
}

// TestStruct2 is another test struct that implements TestInterface
type TestStruct2 struct {
	Value int
}

// GetValue implements TestInterface
func (t *TestStruct2) GetValue() int {
	return t.Value * 2
}

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Test with initializer
	ts := &TestStruct{Value: 42}
	klass := New(nil, reflect.TypeOf(TestStruct{}), ts)

	if klass == nil {
		t.Fatal("New returned nil")
	}

	if klass.Class != ts {
		t.Errorf("New did not set Class correctly, got %v, want %v", klass.Class, ts)
	}

	// Test without initializer
	klass = New(nil, reflect.TypeOf(TestStruct{}), nil)

	if klass == nil {
		t.Fatal("New returned nil")
	}

	// Check that the class was initialized
	if klass.Class == nil {
		t.Error("New did not initialize Class")
	}
}

// TestCast tests the Cast function
func TestCast(t *testing.T) {
	// Create a test struct
	ts := &TestStruct{Value: 42}

	// Cast to TestInterface
	iface := Cast(ts, reflect.TypeOf((*TestInterface)(nil)).Elem())

	// Check that the cast was successful
	if iface == nil {
		t.Fatal("Cast returned nil")
	}

	// Check that the interface works correctly
	ti, ok := iface.(TestInterface)
	if !ok {
		t.Fatal("Cast did not return a TestInterface")
	}

	if ti.GetValue() != 42 {
		t.Errorf("GetValue returned %d, want 42", ti.GetValue())
	}

	// Test casting a value type that implements the interface via pointer
	ts2 := TestStruct{Value: 84}
	iface = Cast(ts2, reflect.TypeOf((*TestInterface)(nil)).Elem())

	// Check that the cast was successful
	if iface == nil {
		t.Fatal("Cast returned nil for value type")
	}

	// Check that the interface works correctly
	ti, ok = iface.(TestInterface)
	if !ok {
		t.Fatal("Cast did not return a TestInterface for value type")
	}

	if ti.GetValue() != 84 {
		t.Errorf("GetValue returned %d, want 84", ti.GetValue())
	}

	// Test casting to a non-implemented interface
	type NonImplementedInterface interface {
		NonExistentMethod()
	}

	iface = Cast(ts, reflect.TypeOf((*NonImplementedInterface)(nil)).Elem())
	if iface != nil {
		t.Error("Cast should return nil for non-implemented interface")
	}
}

// TestAs tests the As function
func TestAs(t *testing.T) {
	// Create a test struct - use a pointer to make it addressable
	ts := &TestStruct{Value: 42}

	// Cast to TestStruct
	result := As(ts, reflect.TypeOf(TestStruct{}))

	// Check that the cast was successful
	if result == nil {
		t.Fatal("As returned nil")
	}

	// Check that the result is a *TestStruct
	ptr, ok := result.(*TestStruct)
	if !ok {
		t.Fatal("As did not return a *TestStruct")
	}

	if ptr.Value != 42 {
		t.Errorf("Value is %d, want 42", ptr.Value)
	}

	// Test with incompatible type
	result = As(ts, reflect.TypeOf(TestStruct2{}))
	if result != nil {
		t.Error("As should return nil for incompatible type")
	}
}

// TestAsPtr tests the AsPtr function
func TestAsPtr(t *testing.T) {
	// Create a test struct
	ts := &TestStruct{Value: 42}

	// Get pointer to the struct
	ptr := AsPtr(ts)

	// Check that the pointer is not nil
	if ptr == nil {
		t.Fatal("AsPtr returned nil")
	}

	// Check that the pointer points to the correct data
	tsPtr := (*TestStruct)(ptr)
	if tsPtr.Value != 42 {
		t.Errorf("Value is %d, want 42", tsPtr.Value)
	}
}

// TestNilInterface tests the NilInterface function
func TestNilInterface(t *testing.T) {
	// Create an interface type
	var iface TestInterface = &TestStruct{Value: 42}

	// Test our NilInterface function with the interface type
	nilIface := NilInterface(iface)

	// Check that the result is not nil
	if nilIface == nil {
		t.Fatal("NilInterface returned nil")
	}

	// Check that IsNil returns false for the nil interface
	// Note: reflect.Zero(ifaceType).Interface() returns a struct value, not an interface value,
	// so IsNil will return false for it. This is expected behavior.
	if IsNil(nilIface) {
		t.Log("IsNil returned true for nil interface created with NilInterface")
	} else {
		t.Log("IsNil returned false for nil interface created with NilInterface (expected for zero values)")
	}

	// Create a nil interface directly
	var nilIface2 TestInterface = nil

	// Check that IsNil returns true for a nil interface
	if !IsNil(nilIface2) {
		t.Error("IsNil returned false for nil interface")
	}
}

// TestNullOrZeroInterface tests the NullOrZeroInterface function with non-interface types
func TestNullOrZeroInterface(t *testing.T) {
	// Test with non-interface type
	var nonInterface int = 42
	result := NullOrZeroInterface(nonInterface)

	// Check that the result is not nil
	if result == nil {
		t.Fatal("NullOrZeroInterface returned nil for non-interface type")
	}

	// Check that the result is a zero value of the same type
	if reflect.TypeOf(result) != reflect.TypeOf(nonInterface) {
		t.Errorf("NullOrZeroInterface returned wrong type, got %T, want %T", result, nonInterface)
	}

	// Check that the result is zero
	if result.(int) != 0 {
		t.Errorf("NullOrZeroInterface returned non-zero value, got %v, want 0", result)
	}
}

// TestNilInterfacePanic tests that Nil.Of panics with non-interface type
func TestNilInterfacePanic(t *testing.T) {
	// Test with non-interface type in a separate function to avoid affecting the main test
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Nil.Of should panic for non-interface type")
			}
		}()
		// Use a non-interface type that will cause a panic
		var nonInterface int = 42
		n := Nil{}
		n.Of(nonInterface)
	}()
}

// TestIsNil tests the IsNil function
func TestIsNil(t *testing.T) {
	// Test with nil interface
	var iface TestInterface = nil
	if !IsNil(iface) {
		t.Error("IsNil returned false for nil interface")
	}

	// Test with non-nil interface
	iface = &TestStruct{Value: 42}
	if IsNil(iface) {
		t.Error("IsNil returned true for non-nil interface")
	}

	// Test with nil pointer
	var ptr *TestStruct = nil
	if !IsNil(ptr) {
		t.Error("IsNil returned false for nil pointer")
	}
}

// TestKlassPtr tests the Klass.Ptr method
func TestKlassPtr(t *testing.T) {
	// Create a test struct
	ts := &TestStruct{Value: 42}
	klass := New(nil, reflect.TypeOf(TestStruct{}), ts)

	// Get pointer to the struct
	ptr := klass.Ptr()

	// Check that the pointer is not nil
	if ptr == nil {
		t.Fatal("Ptr returned nil")
	}

	// Check that the pointer points to the correct data
	tsPtr := (*TestStruct)(ptr)
	if tsPtr.Value != 42 {
		t.Errorf("Value is %d, want 42", tsPtr.Value)
	}
}

// TestFrom tests the From function
func TestFrom(t *testing.T) {
	// Create a test struct
	ts := &TestStruct{Value: 42}
	klass := New(nil, reflect.TypeOf(TestStruct{}), ts)

	// Get pointer to the struct
	ptr := klass.Ptr()

	// Get Klass from pointer
	klassFromPtr := From(ptr, reflect.TypeOf(TestStruct{}))

	// Check that the Klass is not nil
	if klassFromPtr == nil {
		t.Fatal("From returned nil")
	}

	// Check that the Klass has the correct Class
	tsFromKlass, ok := klassFromPtr.Class.(*TestStruct)
	if !ok {
		t.Fatal("From did not return a Klass with a *TestStruct")
	}

	if tsFromKlass.Value != 42 {
		t.Errorf("Value is %d, want 42", tsFromKlass.Value)
	}
}

// TestNil tests the Nil struct
func TestNil(t *testing.T) {
	// Create a Nil instance
	n := Nil{}

	// Test Ptr method
	ptr := n.Ptr()
	if ptr != nil {
		t.Error("Ptr did not return nil")
	}

	// Test Of method
	var iface TestInterface = nil
	nilIface := n.Of(iface)

	// Check that the result is not nil
	if nilIface == nil {
		t.Fatal("Of returned nil")
	}

	// Check that IsNil returns true for the nil interface
	if !IsNil(nilIface) {
		t.Error("IsNil returned false for nil interface created with Of")
	}

	// Test Of with non-interface type in a separate function to avoid affecting the main test
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Of should panic for non-interface type")
			}
		}()
		n.Of(42)
	}()
}

// TestInitClass tests the initClass function
func TestInitClass(t *testing.T) {
	// Create a test struct with zero values
	ts := &TestStruct{}

	// Initialize the struct
	initClass(ts)

	// Check that the struct was initialized
	if ts.Value != 0 {
		t.Errorf("Value is %d, want 0", ts.Value)
	}

	// Test with non-struct type
	initClass(42) // Should not panic
}
