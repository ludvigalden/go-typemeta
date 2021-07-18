package typemeta

import (
	"reflect"
)

// Slice is type meta for a slice
type Slice struct {
	Elem TypeMeta

	typ reflect.Type
}

// Primitive returns true for slices with primitive elements
func (s *Slice) Primitive() bool {
	return s.Elem.Primitive()
}

// Type returns the type of the struct
func (s *Slice) Type() reflect.Type {
	return s.typ
}

// Kind returns
func (s *Slice) Kind() reflect.Kind {
	return reflect.Slice
}

// JSONNonNull returns false for slices
func (s *Slice) JSONNonNull() bool {
	return false
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Slice) Name() string {
	return s.typ.Name()
}

// Copy returns a copy of the ptr type meta
func (s *Slice) Copy() *Slice {
	ns := *s
	return &ns
}

func (s *Slice) String() string {
	return s.Type().String()
}
