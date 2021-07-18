package typemeta

import "reflect"

// TypeMeta is the interface for type meta
type TypeMeta interface {
	// String returns a string representation of the type.
	String() string
	// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
	Name() string
	// Type returns the reflect type of the type meta.
	Type() reflect.Type
	// Type returns the reflect kind of the type meta.
	Kind() reflect.Kind
	// JSONNonNull returns whether it's not possible for the value to be defined as `null` in JSON-format.
	JSONNonNull() bool
	// Primitive returns whether the type meta is for a primitive type.
	Primitive() bool
}
