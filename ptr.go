package typemeta

import (
	"reflect"
)

// Ptr is type meta for a pointer
type Ptr struct {
	Elem TypeMeta

	typ reflect.Type
}

// Layers returns the amount of layers of pointers the pointer is.
// If the element of the pointer is a non-pointer, 1 is returned.
func (s *Ptr) Layers() int {
	switch s.Elem.(type) {
	case *Ptr:
		return s.Elem.(*Ptr).Layers() + 1
	default:
		return 1
	}
}

// Type returns the type of the struct
func (s *Ptr) Type() reflect.Type {
	return s.typ
}

// Kind returns
func (s *Ptr) Kind() reflect.Kind {
	return reflect.Ptr
}

// Primitive returns
func (s *Ptr) Primitive() bool {
	return s.Elem.Primitive()
}

// JSONNonNull returns false for a pointer
func (s *Ptr) JSONNonNull() bool {
	return false
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Ptr) Name() string {
	return s.typ.Name()
}

// Copy returns a copy of the ptr type meta
func (s *Ptr) Copy() *Ptr {
	ns := *s
	return &ns
}

// String returns a string representation of the map type
func (s *Ptr) String() string {
	return s.Type().String()
}
