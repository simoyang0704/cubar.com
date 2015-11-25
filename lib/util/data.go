package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

func DecodeJson(input io.ReadCloser, v interface{}) error {

	decoder := json.NewDecoder(input)
	decoder.UseNumber()
	return decoder.Decode(v)
}

func Md5Hash(str string) string {

	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func StructToJsonString(v interface{}) string {

	ret, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(ret)
}

func StringToStruct(str string, v interface{}) {

	decoder := json.NewDecoder(strings.NewReader(str))
	decoder.UseNumber()
	return decoder.Decode(v)
}

func StringToInt(str string) int {

	i, err := strconv.Atoi(strings.Trim(str, " "))
	if err != nil {
		return 0
	}

	return i
}

func IntToString(i int) string {

	return strconv.Itoa(i)
}

func StringToFloat64(str string) float64 {

	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		f = 0.0
	}

	return f
}

func Float64ToString(f float64) string {

	return strconv.FormatFloat(f, 'f', 10, 64)
}

func Int64ToString(i int64) string {

	return strconv.FormatInt(i, 10)
}

func StringToInt64(str string) int64 {

	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return int64(i)
}

func ParseBool(val interface{}) bool {

	var re bool
	switch val.(type) {
	case bool:
		re = val.(bool)
	case int64:
		re = val.(int64) > 0
	case int:
		re = val.(int) > 0
	case string:
		re = val.(string) != ""
	default:
		re = false
	}

	return re
}

func ParseFloat64(val interface{}) float64 {

	var re float64
	switch val.(type) {
	case int64:
		re = float64(val.(int64))
	case float64:
		re = val.(float64)
	case int:
		re = float64(val.(int))
	case string:
		re = StringToFloat64(val.(string))
	default:
		re = 0
	}

	return re
}

func ParseInt64(val interface{}) int64 {

	var re int64
	switch val.(type) {
	case int64:
		re = val.(int64)
	case float64:
		re = int64(val.(float64))
	case int:
		re = int64(val.(int))
	case string:
		re = StringToInt64(val.(string))
	default:
		re = 0
	}

	return re
}

func ParseString(val interface{}) string {

	var re string
	switch val.(type) {
	case string:
		re = val.(string)
	case int64:
		re = Int64ToString(val.(int64))
	case int:
		re = IntToString(val.(int))
	case float64:
		re = Float64ToString(val.(float64))
	default:
		re = ""
	}

	return re
}

func ParseIntString(val interface{}) string {

	switch val.(type) {
	case float64:
		return Int64ToString(int64(val.(float64)))
	default:
		return ParseString(val)
	}
}

func ToStringArray(val interface{}) []string {

	if val == nil || val == "" {
		return nil
	}

	switch val.(type) {
	case []interface{}:
		vals := val.([]interface{})
		ret := make([]string, len(vals))
		for index, v := range vals {
			ret[index] = ParseString(v)
		}

		return ret
	case []string:
		return val.([]string)
	}

	return nil
}

func ToInt64Array(val interface{}) []int64 {

	if val == nil || val == "" {
		return nil
	}

	switch val.(type) {
	case []interface{}:
		vals := val.([]interface{})
		ret := make([]int64, len(vals))
		for index, v := range vals {
			ret[index] = ParseInt64(v)
		}

		return ret
	case []float64:
		vals := val.([]float64)
		ret := make([]int64, len(vals))
		for index, v := range vals {
			ret[index] = int64(v)
		}

		return ret
	case []int64:
		return val.([]int64)
	}

	return nil
}
