package typemeta

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestConvertPrimitiveValue(t *testing.T) {
	timeValue, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t.Run("converts to string", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e string
		}{{0, "0"}, {100.52, "100.52"}, {true, "true"}, {false, "false"}, {stringifiable{"Test"}, "Test"}, {timeValue, timeValue.Format(time.RFC3339)}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to string: " + err.Error())
			} else if conv.String() != cd.e {
				t.Error("expected \"" + cd.e + "\"  but received \"" + conv.String() + "\"")
			}
		}
	})
	t.Run("converts to int", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e int64
		}{{"0", 0}, {100.52, 100}, {"100", 100}, {100, 100}, {false, 0}, {true, 1}, {"", 0}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to int: " + err.Error())
			} else if conv.Int() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Int()) + "\"")
			}
		}
	})
	t.Run("converts to float", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e float64
		}{{"0", 0}, {100.52, 100.52}, {"100.52", 100.52}, {"1.337", 1.337}, {false, 0}, {true, 1}, {"", 0}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to float: " + err.Error())
			} else if conv.Float() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Int()) + "\"")
			}
		}
	})
	t.Run("converts to bool", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e bool
		}{{"0", false}, {0, false}, {100.52, true}, {100, true}, {-1, false}, {"true", true}, {"false", false}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to bool: " + err.Error())
			} else if conv.Bool() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Bool()) + "\"")
			}
		}
	})
	t.Run("converts to time", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e time.Time
		}{{timeValue.Format(time.RFC3339), timeValue}, {0, time.Time{}}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to time: " + err.Error())
			} else if !conv.Interface().(time.Time).Equal(cd.e) {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Interface()) + "\"")
			}
		}
	})
}

func TestConverteValue(t *testing.T) {
	timeValue, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t.Run("converts to string", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e string
		}{{0, "0"}, {100.52, "100.52"}, {true, "true"}, {false, "false"}, {stringifiable{"Test"}, "Test"}, {timeValue, timeValue.Format(time.RFC3339)}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to string: " + err.Error())
			} else if conv.String() != cd.e {
				t.Error("expected \"" + cd.e + "\"  but received \"" + conv.String() + "\"")
			}
		}
	})
	t.Run("converts to int", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e int64
		}{{"0", 0}, {100.52, 100}, {"100", 100}, {100, 100}, {false, 0}, {true, 1}, {"", 0}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to int: " + err.Error())
			} else if conv.Int() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Int()) + "\"")
			}
		}
	})
	t.Run("converts to float", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e float64
		}{{"0", 0}, {100.52, 100.52}, {"100.52", 100.52}, {"1.337", 1.337}, {false, 0}, {true, 1}, {"", 0}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to float: " + err.Error())
			} else if conv.Float() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Int()) + "\"")
			}
		}
	})
	t.Run("converts to bool", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e bool
		}{{"0", false}, {0, false}, {100.52, true}, {100, true}, {-1, false}, {"true", true}, {"false", false}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to bool: " + err.Error())
			} else if conv.Bool() != cd.e {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Bool()) + "\"")
			}
		}
	})
	t.Run("converts to time", func(t *testing.T) {
		cd := []struct {
			v interface{}
			e time.Time
		}{{timeValue.Format(time.RFC3339), timeValue}, {0, time.Time{}}}
		for _, cd := range cd {
			conv, err := ConvertValue(reflect.ValueOf(cd.v), cd.e)
			if err != nil {
				t.Error("failed converting to time: " + err.Error())
			} else if !conv.Interface().(time.Time).Equal(cd.e) {
				t.Error("expected \"" + fmt.Sprint(cd.e) + "\"  but received \"" + fmt.Sprint(conv.Interface()) + "\"")
			}
		}
	})
}

type stringifiable struct {
	Name string
}

func (s stringifiable) String() string {
	return s.Name
}
