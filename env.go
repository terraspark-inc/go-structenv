package structenv

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	truthyValues = []string{"y", "yes", "1", "true", "t", "on"}
)

func LoadFromEnv(objs ...interface{}) error {
	for _, obj := range objs {
		ov := reflect.Indirect(reflect.ValueOf(obj))
		ot := ov.Type()

		for i := 0; i < ot.NumField(); i++ {
			fv := ov.Field(i)

			tagValue := ot.Field(i).Tag.Get("env")
			defaultTagValue := ot.Field(i).Tag.Get("default")

			if tagValue == "" {
				continue
			}

			envValStr := os.Getenv(tagValue)

			if (envValStr == "") && (defaultTagValue != "") {
				envValStr = defaultTagValue
			}

			// Skip if env unset or set to empty string and no default provided
			if (envValStr == "") && (defaultTagValue == "") {
				continue
			}

			switch fv.Kind() {
			case reflect.String:
				fv.SetString(envValStr)
			case reflect.Int:
				v, err := strconv.ParseInt(envValStr, 10, 64)
				if err != nil {
					return fmt.Errorf("unable to parse %s as int64: %v", tagValue, err)
				}
				fv.SetInt(v)
			case reflect.Float64:
				v, err := strconv.ParseFloat(envValStr, 64)
				if err != nil {
					return fmt.Errorf("unable to parse %s as float64: %v", tagValue, err)
				}
				fv.SetFloat(v)
			case reflect.Bool:
				envValStrLower := strings.ToLower(envValStr)
				var v bool
				for _, val := range truthyValues {
					if envValStrLower == val {
						v = true
					}
				}
				fv.SetBool(v)
			case reflect.ValueOf(time.Second).Kind():
				v, err := time.ParseDuration(envValStr)
				if err != nil {
					return fmt.Errorf("unable to parse %s as time.Duration: %v", tagValue, err)
				}
				fv.Set(reflect.ValueOf(v))
			default:
				return fmt.Errorf("unsupported field type %+v of key %v", fv.Type(), tagValue)
			}
		}
	}
	return nil
}
