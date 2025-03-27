package oop

import (
	"testing"
)

// TestAnimal is a test interface used in the tests
type TestAnimal interface {
	Sound() string
}

// TestDog is a test struct that implements TestAnimal
type TestDog struct {
	Name string
}

// Sound implements TestAnimal
func (d *TestDog) Sound() string {
	return d.Name + ": Woof!"
}

// TestCat is another test struct that implements TestAnimal
type TestCat struct {
	Name string
}

// Sound implements TestAnimal
func (c *TestCat) Sound() string {
	return c.Name + ": Meow!"
}

// TestNewObjectFactory tests the NewObjectFactory function
func TestNewObjectFactory(t *testing.T) {
	factory := NewObjectFactory()

	if factory == nil {
		t.Fatal("NewObjectFactory returned nil")
	}

	if factory.allocator != nil {
		t.Errorf("NewObjectFactory did not set allocator to nil, got %v", factory.allocator)
	}
}

// TestCreateObject tests the CreateObject method of ObjectFactory
func TestCreateObject(t *testing.T) {
	factory := NewObjectFactory()

	// Test with nil initializer
	obj := factory.CreateObject(nil)
	if obj != nil {
		t.Error("CreateObject should return nil for nil initializer")
	}

	// Test with valid initializer
	dog := &TestDog{Name: "Buddy"}
	obj = factory.CreateObject(dog)

	if obj == nil {
		t.Fatal("CreateObject returned nil for valid initializer")
	}

	if obj.klass == nil {
		t.Fatal("CreateObject did not set klass")
	}

	if obj.klass.Class != dog {
		t.Errorf("CreateObject did not set klass.Class correctly, got %v, want %v", obj.klass.Class, dog)
	}
}

// TestObjectWrapperAs tests the As method of ObjectWrapper
func TestObjectWrapperAs(t *testing.T) {
	factory := NewObjectFactory()

	// Create a test object
	dog := &TestDog{Name: "Buddy"}
	obj := factory.CreateObject(dog)

	// Cast to TestAnimal
	animal, err := obj.As((*TestAnimal)(nil))

	// Check that the cast was successful
	if err != nil {
		t.Fatalf("As returned error: %v", err)
	}
	if animal == nil {
		t.Fatal("As returned nil")
	}

	// Check that the interface works correctly
	a, ok := animal.(TestAnimal)
	if !ok {
		t.Fatal("As did not return a TestAnimal")
	}

	if a.Sound() != "Buddy: Woof!" {
		t.Errorf("Sound returned %q, want %q", a.Sound(), "Buddy: Woof!")
	}

	// Test with nil klass
	obj = &ObjectWrapper{klass: nil}
	animal, err = obj.As((*TestAnimal)(nil))
	if err == nil {
		t.Errorf("As should returned error for nil klass")
	}
	if animal != nil {
		t.Error("As should return nil for nil klass")
	}

	// Test with nil klass.Class
	obj = &ObjectWrapper{klass: &Klass{Class: nil}}
	animal, err = obj.As((*TestAnimal)(nil))
	if err == nil {
		t.Errorf("As should return error for nil klass.Class")
	}
	if animal != nil {
		t.Error("As should return nil for nil klass.Class")
	}

	// Test with non-pointer type (should return the type itself)
	var nonPointer int = 42
	nonPointerObj, err := obj.As(nonPointer)
	if err == nil {
		t.Error("As should return error for non-pointer type")
	}
	if nonPointerObj != nil {
		t.Error("As should return nil for non-pointer type")
	}
	// if err != nil {
	// 	t.Error("As should not return error for non-pointer type")
	// }
	// if nonPointerObj == nil {
	// 	t.Error("As should not return nil for non-pointer type")
	// }
	// if nonPointerObj != nonPointer {
	// 	t.Errorf("As returned %v, want %v", nonPointerObj, nonPointer)
	// }

	// Test with pointer to non-interface type (should return error)
	var nonInterface *int = new(int)
	*nonInterface = 42
	nonInterfaceObj, err := obj.As(nonInterface)
	if err == nil {
		t.Error("As should return error for pointer to non-interface type")
	}
	if nonInterfaceObj != nil {
		t.Error("As should not return nil for pointer to non-interface type")
	}

}

// TestObjectWrapperDestroy tests the Destroy method of ObjectWrapper
func TestObjectWrapperDestroy(t *testing.T) {
	factory := NewObjectFactory()

	// Create a test object
	dog := &TestDog{Name: "Buddy"}
	obj := factory.CreateObject(dog)

	// Destroy the object
	obj.Destroy()

	// Check that klass was set to nil
	if obj.klass != nil {
		t.Error("Destroy did not set klass to nil")
	}

	// Test with nil klass (should not panic)
	obj = &ObjectWrapper{klass: nil}
	obj.Destroy() // Should not panic
}

// TestGetUnderlyingObject tests the GetUnderlyingObject method of ObjectWrapper
func TestGetUnderlyingObject(t *testing.T) {
	factory := NewObjectFactory()

	// Create a test object
	dog := &TestDog{Name: "Buddy"}
	obj := factory.CreateObject(dog)

	// Get the underlying object
	underlying := obj.GetUnderlyingObject()

	// Check that the underlying object is correct
	if underlying != dog {
		t.Errorf("GetUnderlyingObject returned %v, want %v", underlying, dog)
	}

	// Test with nil klass
	obj = &ObjectWrapper{klass: nil}
	underlying = obj.GetUnderlyingObject()
	if underlying != nil {
		t.Error("GetUnderlyingObject should return nil for nil klass")
	}
}

// TestIntegration tests the integration of the user-friendly API
func TestIntegration(t *testing.T) {
	factory := NewObjectFactory()

	// Create Dog and Cat instances
	dogObj := factory.CreateObject(&TestDog{Name: "Buddy"})
	catObj := factory.CreateObject(&TestCat{Name: "Whiskers"})

	// Cast to TestAnimal
	dogAnimal, err := dogObj.As((*TestAnimal)(nil))
	if err != nil {
		t.Fatalf("As returned error for dog: %v", err)
	}
	catAnimal, err := catObj.As((*TestAnimal)(nil))
	if err != nil {
		t.Fatalf("As returned error for cat: %v", err)
	}

	// Check that the casts were successful
	if dogAnimal == nil {
		t.Fatal("As returned nil for dog")
	}
	if catAnimal == nil {
		t.Fatal("As returned nil for cat")
	}

	// Check that the interfaces work correctly
	dog, ok := dogAnimal.(TestAnimal)
	if !ok {
		t.Fatal("As did not return a TestAnimal for dog")
	}
	cat, ok := catAnimal.(TestAnimal)
	if !ok {
		t.Fatal("As did not return a TestAnimal for cat")
	}

	if dog.Sound() != "Buddy: Woof!" {
		t.Errorf("Dog sound returned %q, want %q", dog.Sound(), "Buddy: Woof!")
	}
	if cat.Sound() != "Whiskers: Meow!" {
		t.Errorf("Cat sound returned %q, want %q", cat.Sound(), "Whiskers: Meow!")
	}

	// Clean up
	dogObj.Destroy()
	catObj.Destroy()

	// Check that the objects were destroyed
	if dogObj.klass != nil {
		t.Error("Destroy did not set dog klass to nil")
	}
	if catObj.klass != nil {
		t.Error("Destroy did not set cat klass to nil")
	}
}
