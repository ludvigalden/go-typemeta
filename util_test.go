package typemeta

import (
	"fmt"
	"testing"
)

func TestUtil(t *testing.T) {
	type StructA struct {
		Name   string
		Active bool
	}
	t.Run("StructOf", func(t *testing.T) {
		s := GetStruct(StructA{})
		sptr := Get(&StructA{})
		sptrsl := Get([]*StructA{})
		sptrslptr := Get(&[]*StructA{})

		salts := []TypeMeta{sptr, sptrsl, sptrslptr}
		for _, salt := range salts {
			if sf := StructOf(salt); sf != s {
				t.Error("StructOf(" + salt.String() + ") should be " + s.String() + " but received " + fmt.Sprint(sf))
			}
		}
	})
}
