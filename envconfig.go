package envconfig

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

const CONFIG_VALUES_DELIMITER = ","

func Parse(cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		val := ft.Tag.Get("default")

		envKey := ft.Tag.Get("env")
		if envKey == "" {
			envKey = strings.ToUpper(ft.Name)
		}
		if envValue := os.Getenv(envKey); envValue != "" {
			val = envValue
		}

		if val != "" {
			if fv.CanSet() {
				switch fv.Kind() {
				case reflect.Slice:
					vals := strings.Split(val, CONFIG_VALUES_DELIMITER)
					for _, v := range vals {
						fv.Set(reflect.Append(fv, reflect.ValueOf(v)))
					}
				case reflect.Bool:
					var b bool
					if val == "true" || val == "1" {
						b = true
					}
					fv.SetBool(b)
				case reflect.String:
					fv.SetString(val)
				case reflect.Int:
					valInt, err := strconv.Atoi(val)
					if err != nil {
						panic(err)
					}

					fv.SetInt(int64(valInt))
				}
			}
		}
	}
}

func Usage(cfg interface{}) {
	fmt.Printf("USAGE: [options] %s\n", path.Base(os.Args[0]))
	fmt.Printf("OPTIONS:\n")

	t := reflect.TypeOf(cfg).Elem()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)

		e := ft.Tag.Get("env")
		if e == "" {
			e = strings.ToUpper(ft.Name)
		}
		d := ft.Tag.Get("default")
		u := ft.Tag.Get("usage")

		fmt.Printf("\t%s [%v] %s\n", e, d, u)
	}

	fmt.Printf("EXAMPLE:\n\tVAR1=VALUE1 %s", path.Base(os.Args[0]))
}
