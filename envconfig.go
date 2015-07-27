package envconfig

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

// char for split values
const ConfigValuesDelimiter = ","

// Parse fill config values from env-variables
func Parse(cfg interface{}) {
	v := reflect.ValueOf(cfg).Elem()
	t := reflect.TypeOf(cfg).Elem()

	for i := 0; i < t.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if fv.CanSet() == false {
			break
		}

		envKey := getEnvKey(ft)

		val := ft.Tag.Get("default")
		if envValue := os.Getenv(envKey); envValue != "" {
			val = envValue
		}

		if val != "" {
			switch fv.Kind() {
			case reflect.Slice:
				for _, v2 := range parseStringSlice(val) {
					fv.Set(reflect.Append(fv, reflect.ValueOf(v2)))
				}
			case reflect.Bool:
				fv.SetBool(parseBool(val))
			case reflect.String:
				fv.SetString(val)
			case reflect.Int:
				fv.SetInt(parseInt(val))
			}
		}
	}
}

func getEnvKey(ft reflect.StructField) string {
	envKey := ft.Tag.Get("env")
	if envKey == "" {
		envKey = strings.ToUpper(ft.Name)
	}

	return envKey
}

func parseBool(val string) bool {
	return val == "true" || val == "1"
}

func parseInt(val string) int64 {
	valInt, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}

	return int64(valInt)
}

func parseStringSlice(val string) []string {
	return strings.Split(val, ConfigValuesDelimiter)
}

// Usage print help information
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
