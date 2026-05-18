package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func MustLoad(cfg interface{}) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		panic("config: must pass pointer to struct")
	}

	elem := v.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		fieldType := typ.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envKey := parts[0]
		required := false
		defaultVal := ""

		for _, part := range parts[1:] {
			if part == "required" {
				required = true
			} else if strings.HasPrefix(part, "default=") {
				defaultVal = strings.TrimPrefix(part, "default=")
			}
		}

		envVal := os.Getenv(envKey)
		if envVal == "" {
			if required {
				panic(fmt.Sprintf("config: required env variable %s is not set", envKey))
			}
			envVal = defaultVal
		}

		if envVal == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envVal)
		case reflect.Int, reflect.Int64:
			intVal, err := strconv.ParseInt(envVal, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("config: invalid int value for %s: %v", envKey, err))
			}
			field.SetInt(intVal)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(envVal)
			if err != nil {
				panic(fmt.Sprintf("config: invalid bool value for %s: %v", envKey, err))
			}
			field.SetBool(boolVal)
		}
	}
}
