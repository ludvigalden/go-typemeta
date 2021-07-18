package typemeta

import (
	"reflect"
)

// Array is type meta for an array
type Array struct {
	Elem TypeMeta
	typ  reflect.Type
}

// Primitive returns true for slices with primitive elements
func (s *Array) Primitive() bool {
	return s.Elem.Primitive()
}

// Type returns
func (s *Array) Type() reflect.Type {
	return s.typ
}

// Kind returns
func (s *Array) Kind() reflect.Kind {
	return reflect.Array
}

// JSONNonNull returns
func (s *Array) JSONNonNull() bool {
	return false
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Array) Name() string {
	return s.typ.Name()
}

func (s *Array) String() string {
	return s.Type().String()
}
