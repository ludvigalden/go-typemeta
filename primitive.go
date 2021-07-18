package typemeta

import (
	"reflect"
)

// Primitive is type meta for a primitive type
type Primitive struct {
	typ  reflect.Type
	enum *Enum
	name string
}

// SetName sets the name of the struct type meta
func (s *Primitive) SetName(name string) *Primitive {
	s.name = name
	return s
}

// SetEnum sets the enum of the primitive type
func (s *Primitive) SetEnum(enum *Enum) *Primitive {
	s.enum = enum
	return s
}

// Enum returns the enum of the primitive type.
func (s *Primitive) Enum() *Enum {
	return s.enum
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Primitive) Name() string {
	if s.name != "" {
		return s.name
	}
	return s.typ.Name()
}

// Type returns the type of the primitive type
func (s *Primitive) Type() reflect.Type {
	return s.typ
}

// Kind returns the reflect kind of the primitive type
func (s *Primitive) Kind() reflect.Kind {
	return s.typ.Kind()
}

// Primitive returns true for primitive types
func (s *Primitive) Primitive() bool {
	return true
}

// JSONNonNull returns true for primitive types
func (s *Primitive) JSONNonNull() bool {
	return true
}

// Copy returns a copy of the primitive type meta
func (s *Primitive) Copy() *Primitive {
	ns := *s
	return &ns
}

func (s *Primitive) String() string {
	return s.Type().String()
}
