package main

import (
	"fmt"
	"strconv"
	"time"
)

var MessageCreateAsset = "This topic will be created by bot %s"

func toStr(v interface{}) string {
	switch typedV := v.(type) {
	case string:
		return typedV
	case fmt.Stringer:
		return typedV.String()
	case int:
		return strconv.Itoa(typedV)
	case uint:
		return strconv.FormatUint(uint64(typedV), 10)
	case uint64:
		return strconv.FormatUint(typedV, 10)
	case int64:
		return strconv.FormatInt(typedV, 10)
	case float64:
		return strconv.FormatFloat(typedV, 'g', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(typedV), 'g', -1, 32)
	case int32:
		return strconv.FormatInt(int64(typedV), 10)
	case int16:
		return strconv.FormatInt(int64(typedV), 10)
	case int8:
		return strconv.FormatInt(int64(typedV), 10)
	case []byte:
		return string(typedV)
	case bool:
		return strconv.FormatBool(typedV)
	case time.Time:
		return typedV.Format(time.RFC3339)
	}
	return fmt.Sprintf("%v", v)
}
