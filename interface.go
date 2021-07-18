package typemeta

import (
	"reflect"
)

// Interface is type meta for an interface
type Interface struct {
	typ reflect.Type
}

// Primitive returns true for interfaces
func (s *Interface) Primitive() bool {
	return true
}

// Type returns the reflect type of the interface or nil
func (s *Interface) Type() reflect.Type {
	return s.typ
}

// Kind returns `reflect.Interface` for interfaces
func (s *Interface) Kind() reflect.Kind {
	return reflect.Interface
}

// JSONNonNull returns false for interfaces
func (s *Interface) JSONNonNull() bool {
	return false
}

// Name returns the type meta's explicitly set name or the type's name within its package for a defined type. For other (non-defined) types it returns the empty string.
func (s *Interface) Name() string {
	if s.typ == nil {
		return "interface{}"
	}
	return s.typ.Name()
}
func (s *Interface) String() string {
	if s.typ == nil {
		return "interface{}"
	}
	return s.typ.String()
}
