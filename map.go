package typemeta

import (
	"reflect"
)

// Map is type meta for a map
type Map struct {
	Key  TypeMeta
	Elem TypeMeta

	typ reflect.Type
}

var _ TypeMeta = &Map{}

// Type returns the type of the struct
func (s *Map) Type() reflect.Type {
	return s.typ
}

// JSONNonNull returns false for a map
func (s *Map) JSONNonNull() bool {
	return false
}

// Kind returns
func (s *Map) Kind() reflect.Kind {
	return reflect.Map
}

// Primitive returns true for maps with primitive keys and elements
func (s *Map) Primitive() bool {
	return s.Key.Primitive() && s.Elem.Primitive()
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Map) Name() string {
	return s.typ.Name()
}

func (s *Map) String() string {
	return "map[" + s.Key.String() + "]" + s.Elem.String()
}
