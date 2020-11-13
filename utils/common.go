package utils

import (
	"reflect"
	"strconv"
	"strings"
)

func ConvertStructToMap(s interface{}, tagName string, allString bool) (result map[string]interface{}) {
	result = make(map[string]interface{})
	val := reflect.ValueOf(s).Elem()
	key := reflect.TypeOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		mapKey := key.Field(i).Name
		if tagName != "" {
			mapKey = key.Field(i).Tag.Get(tagName)
		}
		if mapKey == "-" {
			continue
		}
		switch val.Field(i).Kind() {
		case reflect.String:
			result[mapKey] = val.Field(i).String()
		case reflect.Int64:
			if allString {
				result[mapKey] = strconv.Itoa(int(val.Field(i).Int()))
			} else {
				result[mapKey] = val.Field(i).Int()
			}
		case reflect.Int:
			if allString {
				result[mapKey] = strconv.Itoa(int(val.Field(i).Int()))
			} else {
				result[mapKey] = int(val.Field(i).Int())
			}
		case reflect.Bool:
			result[mapKey] = val.Field(i).Bool()
		case reflect.Ptr:
			tmp := ConvertStructToMap(val.Field(i).Interface(), tagName, allString)
			for k, v := range tmp {
				result[k] = v
			}
		case reflect.Struct:
			tmp := ConvertStructToMap(val.Field(i).Addr().Interface(), tagName, allString)
			for k, v := range tmp {
				result[k] = v
			}
		}
	}
	return
}

func Concatenate(str []string) string {
	var res strings.Builder
	for i := 0; i < len(str); i++ {
		res.WriteString(str[i])
	}
	return res.String()
}
