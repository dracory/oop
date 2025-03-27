package oop

import (
	"fmt"
	"reflect"
)

// ObjectFactory provides a user-friendly way to create and manage objects.
// It abstracts away the complexity of the underlying OOP implementation.
type ObjectFactory struct {
	allocator interface{}
}

// NewObjectFactory creates a new ObjectFactory.
func NewObjectFactory() *ObjectFactory {
	return &ObjectFactory{
		allocator: nil, // Using nil allocator for simplicity
	}
}

// CreateObject creates a new object of the specified type with the given initializer.
// It simplifies the object creation process by hiding the reflection details.
// Example: factory.CreateObject(&Dog{Name: "Buddy"})
func (f *ObjectFactory) CreateObject(initializer interface{}) *ObjectWrapper {
	if initializer == nil {
		return nil
	}

	// Get the type of the initializer
	objType := reflect.TypeOf(initializer)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}

	// Create a new object using the underlying OOP implementation
	klass := New(f.allocator, objType, initializer)

	// Wrap the object for easier use
	return &ObjectWrapper{
		klass: klass,
	}
}

// ObjectWrapper provides a user-friendly wrapper around a Klass object.
// It simplifies common operations like casting and type checking.
type ObjectWrapper struct {
	klass *Klass
}

// As casts the object to the specified interface type.
// It simplifies the casting process by hiding the reflection details.
// Example: dogObj.As((*IAnimal)(nil))
// Returns the casted object or an error if the cast is not possible.
func (o *ObjectWrapper) As(interfacePtr interface{}) (interface{}, error) {
	// Check if the input is a valid interface pointer
	if interfacePtr == nil {
		return nil, fmt.Errorf("interfacePtr cannot be nil")
	}

	// if interfaceTypeKind == reflect.Int ||
	// 	interfaceTypeKind == reflect.Int8 ||
	// 	interfaceTypeKind == reflect.Int16 ||
	// 	interfaceTypeKind == reflect.Int32 ||
	// 	interfaceTypeKind == reflect.Int64 {
	// 	return interfacePtr, nil
	// }

	if o.klass == nil || o.klass.Class == nil {
		return nil, fmt.Errorf("object is not initialized")
	}

	// Get the interface type
	interfaceType := reflect.TypeOf(interfacePtr)
	interfaceTypeKind := interfaceType.Kind()

	// Check if the type is a pointer
	if interfaceTypeKind != reflect.Ptr {
		return nil, fmt.Errorf("interfacePtr must be a pointer to an interface type")
	}

	// Check if the pointer points to an interface
	if interfaceType.Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("interfacePtr must be a pointer to an interface type")
	}

	interfaceType = interfaceType.Elem()

	// Cast the object to the interface type
	return Cast(o.klass.Class, interfaceType), nil
}

// Destroy deinitializes and destroys the object.
func (o *ObjectWrapper) Destroy() {
	if o.klass != nil {
		o.klass.Deinit()
		o.klass = nil
	}
}

// GetUnderlyingObject returns the underlying object.
func (o *ObjectWrapper) GetUnderlyingObject() interface{} {
	if o.klass == nil {
		return nil
	}
	return o.klass.Class
}
