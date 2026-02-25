package structs

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// MarshalerMap implements a custom [MarshalToMap].
type MarshalerMap interface {
	MarshalMap() map[string]any
}

type options struct {
	key string
	fn  func(any) (any, bool)
}

func newOpt(key string, fn func(any) (any, bool)) options {
	return options{key, fn}
}

var opts = []options{
	newOpt("omitempty", func(v any) (any, bool) {
		if v == nil {
			return nil, false
		}
		refVal := reflect.ValueOf(v)
		if reflect.DeepEqual(v, reflect.Zero(refVal.Type()).Interface()) {
			return v, false
		}
		return v, true
	}),
	newOpt("string", func(v any) (any, bool) { return fmt.Sprintf("%v", v), true }),
}

func getElem(v reflect.Value) any {
	if conv, ok := v.Interface().(MarshalerMap); ok {
		return conv.MarshalMap()
	}
	switch v.Kind() {
	case reflect.Struct:
		return MarshalToMap(v.Interface())
	case reflect.Pointer:
		if v.IsNil() {
			return nil
		}
		return getElem(v.Elem())
	default:
		return v.Interface()
	}
}

// MarshalToMap transforms a struct into a map.
//
// If v is not a map, it returns nil.
//
// Implements [MarshalerMap] to have a custom behavior.
func MarshalToMap(v any) map[string]any {
	if conv, ok := v.(MarshalerMap); ok {
		return conv.MarshalMap()
	}
	ref := reflect.ValueOf(v)
	switch ref.Kind() {
	case reflect.Struct:
	case reflect.Pointer:
		return MarshalToMap(ref.Elem().Interface())
	default:
		return nil
	}
	refType := ref.Type()
	fields := ref.NumField()
	mp := make(map[string]any, fields)
	for i := range fields {
		field := ref.Field(i)
		fieldType := refType.Field(i)
		val := getElem(field)
		name := fieldType.Name
		data := strings.Split(refType.Field(i).Tag.Get("json"), ",")
		if len(data) > 0 {
			if len(data[0]) > 0 {
				name = data[0]
			}
			if len(data) > 1 {
				tagOpts := data[1:]
				ok := true
				i := 0
				for i < len(opts) && ok {
					opt := opts[i]
					if slices.Contains(tagOpts, opt.key) {
						val, ok = opt.fn(val)
					}
					i++
				}
				if !ok {
					continue
				}
			}
		}
		mp[name] = val
	}
	return mp
}
