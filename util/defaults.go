package util

import (
	"reflect"
	"strconv"
	"strings"
)

// SetDefaults will read the default tag and set it to struct field. It is also possible
// to use namespaces for embedded or custom types. So for example:
//
// type A struct {
//	  B string `default[C]:"A Default" default[E]:"E Default"`
// }
//
// type B struct {
// 	  C A  `default.ns:"C"`
// }
//
// type D struct {
// 	  E A  `default.ns:"E"`
// }
//
// f := new(D); SetDefaults(f); fmt.Println(f.E) // will print 'E Default'
func SetDefaults(c interface{}) {
	parseValue(reflect.ValueOf(c), "")
}

func parseValue(v reflect.Value, ns string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i, l := 0, v.NumField(); i < l; i++ {
		value := v.Field(i)
		field := v.Type().Field(i)
		// embedded types
		if reflect.Struct == value.Kind() {
			parseValue(value, field.Tag.Get("default.ns"))
			continue
		}
		// types with namespaces
		if "" != ns {
			if defaults, ok := field.Tag.Lookup("default[" + ns + "]"); ok {
				setValue(value, defaults)
				continue
			}
		}
		// no namespace and fallback
		if defaults, ok := field.Tag.Lookup("default"); ok {
			setValue(value, defaults)
		}
	}
}

func setValue(v reflect.Value, d string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(d)
	case reflect.Array:
		switch v.Type().Elem().Kind() {
		case reflect.Int:
			list := strings.Split(d, ",")
			for i, l := 0, len(list); i < l; i++ {
				if val, err := strconv.Atoi(list[i]); err == nil {
					v.Index(i).SetInt(int64(val))
				}
			}
		}
	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		case reflect.String:
			list := strings.Split(d, ",")
			slice := reflect.MakeSlice(v.Type(), len(list), len(list))
			for i, l := 0, len(list); i < l; i++ {
				slice.Index(i).SetString(list[i])
			}
			v.Set(slice)
		}

	}
}
