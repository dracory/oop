# Go OOP Package

This package provides an implementation of Object-Oriented Programming (OOP) concepts in Go, focusing on creating a runtime type system that supports features typically found in languages with more traditional OOP support.

## Overview

Go is not a traditional object-oriented language. While it has structs and interfaces, it lacks features like inheritance, method overriding, and runtime type information that are common in languages like C++, Java, or C#. This package provides a framework to simulate these OOP features in Go.

## Installation

```bash
import "project/pkg/oop"
```

## Key Features

### 1. Runtime Type Information (RTTI)

The implementation provides runtime type information through the `TypeInfo` struct, which stores:
- Type name
- Type ID (unique identifier)

This allows for type checking and identification at runtime, similar to C++'s `typeid` or Java's reflection.

### 2. Virtual Method Tables (Vtables)

The implementation simulates vtables, which are used in languages like C++ to support dynamic dispatch:
- `VtableInfo` associates a type ID with a pointer to its vtable
- This enables method overriding and polymorphic behavior

### 3. Dynamic Casting

Several casting functions are provided:
- `Cast`: Converts an object to a different type, handling interface and type conversions
- `As`: Performs a dynamic cast and returns an optional pointer
- `AsPtr`: Returns a pointer to the object's data

These functions allow for safe type conversions, similar to C++'s `dynamic_cast` or C#'s `as` operator.

### 4. Class Metadata

The implementation maintains class metadata through:
- `ClassInfo`: Holds runtime information about a class
- `KlassHeader`: Contains metadata for a class instance
- `Klass`: Represents a class instance with metadata

This metadata enables features like inheritance, initialization/deinitialization, and type checking.

### 5. Interface Objects

The `IObject` struct represents an interface object with:
- A pointer to the underlying data
- A pointer to the vtable

This allows for interface-based programming with dynamic dispatch.

### 6. Nil Interface Handling

The implementation provides utilities for working with nil interfaces:
- `Nil` struct and its methods for creating nil interface values
- `NilInterface` function to create a nil instance of a specified interface type
- `IsNil` function to check if an interface instance is nil

## Example Usage

### User-Friendly API

The example below demonstrates how to use the user-friendly API provided by this package:

```go
package main

import (
    "fmt"
    "reflect"
    
    "project/pkg/oop"
)

// Define an interface
type IAnimal interface {
    MakeSound()
}

// Define implementing structs
type Dog struct {
    Name string
}

func (d *Dog) MakeSound() {
    fmt.Printf("%s: Woof!\n", d.Name)
}

type Cat struct {
    Name string
}

func (c *Cat) MakeSound() {
    fmt.Printf("%s: Meow!\n", c.Name)
}

func main() {
    // Create an object factory
    factory := oop.NewObjectFactory()
    
    // Create Dog and Cat instances using the user-friendly API
    dogObj := factory.CreateObject(&Dog{Name: "Buddy"})
    catObj := factory.CreateObject(&Cat{Name: "Whiskers"})
    
    // Cast the objects to the IAnimal interface using the simplified API
    animalDog := dogObj.As((*IAnimal)(nil))
    animalCat := catObj.As((*IAnimal)(nil))
    
    // Polymorphic method calls
    animalDog.(IAnimal).MakeSound() // Outputs: Buddy: Woof!
    animalCat.(IAnimal).MakeSound() // Outputs: Whiskers: Meow!
    
    // Clean up resources
    dogObj.Destroy()
    catObj.Destroy()
}
```

### Low-Level API

For comparison, here's how to use the low-level API directly:

```go
// Create Dog and Cat instances
dog := oop.New(allocator, reflect.TypeOf(Dog{}), &Dog{Name: "Rex"})
cat := oop.New(allocator, reflect.TypeOf(Cat{}), &Cat{Name: "Felix"})

// Cast to the IAnimal interface
animalDog := oop.Cast(dog.Class, reflect.TypeOf((*IAnimal)(nil)).Elem())
animalCat := oop.Cast(cat.Class, reflect.TypeOf((*IAnimal)(nil)).Elem())

// Polymorphic method calls
animalDog.(IAnimal).MakeSound() // Outputs: Rex: Woof!
animalCat.(IAnimal).MakeSound() // Outputs: Felix: Meow!

// Clean up resources
dog.Deinit()
cat.Deinit()
```

## Technical Details

### Memory Layout

The implementation carefully manages memory layout to support features like:
- Finding a class instance from an object pointer
- Accessing vtables for dynamic dispatch
- Maintaining class hierarchies

### Reflection Usage

The implementation leverages Go's reflection package to:
- Get type information at runtime
- Create and manipulate objects dynamically
- Perform type checking and conversions

### Unsafe Operations

The implementation uses the `unsafe` package for low-level memory operations:
- Converting between different pointer types
- Calculating memory offsets
- Accessing memory directly

## User-Friendly API

The package includes a user-friendly API layer that simplifies working with the OOP implementation:

### ObjectFactory

The `ObjectFactory` class provides a simplified way to create objects:

```go
// Create a factory
factory := oop.NewObjectFactory()

// Create an object (much simpler than using New directly)
dogObj := factory.CreateObject(&Dog{Name: "Buddy"})
```

Benefits:
- Hides the complexity of the underlying OOP implementation
- No need to manually specify types using reflection
- Provides a more intuitive object creation process

### ObjectWrapper

The `ObjectWrapper` class wraps a `Klass` object and provides simplified methods:

```go
// Cast to an interface (much simpler than using Cast directly)
animalDog := dogObj.As((*IAnimal)(nil))

// Clean up resources
dogObj.Destroy()
```

Benefits:
- Simplifies casting operations
- Provides a cleaner API for resource management
- Hides the internal details of the OOP implementation

## Benefits and Use Cases

This OOP implementation is useful for:

1. **Complex Object Hierarchies**: When you need to model complex inheritance relationships
2. **Plugin Systems**: For dynamically loading and using components with polymorphic behavior
3. **Language Interoperability**: When interfacing with C++ or other OOP languages
4. **Legacy Code Migration**: When porting code from OOP languages to Go

## Limitations

While this implementation provides many OOP features, it has some limitations:

1. **Performance Overhead**: The use of reflection and dynamic dispatch adds runtime overhead
2. **Type Safety**: Some operations use `unsafe` and may bypass Go's type system
3. **Complexity**: The implementation adds complexity compared to idiomatic Go code
4. **Maintenance**: Code using this framework may be harder to maintain than standard Go code

## Conclusion

This OOP package demonstrates how to extend Go with features from traditional object-oriented languages. The addition of a user-friendly API layer shows how complex functionality can be wrapped in a more intuitive interface. While this approach is not recommended for most Go projects (which should follow Go's idiomatic patterns), it provides valuable insights into language design and API development, and can be useful for specific use cases where traditional OOP features are required.
