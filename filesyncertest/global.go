package filesyncertest

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/rickonono3/m2obj"
)

var obj *m2obj.Object

var filePath string

func initTestData(format string) {
	if runtime.GOOS == "windows" {
		filePath = filepath.Join(os.Getenv("USERPROFILE"), "test."+format)
	} else {
		filePath = "~/test." + format
	}
	obj = m2obj.New(m2obj.Group{
		"a": "a",
		"b": float64(2),
		"c": true,
		"d": m2obj.Group{
			"Sa": "Sa",
			"Sb": float64(233),
			"Sc": false,
		},
		"e": m2obj.Group{
			"ea": m2obj.Group{
				"eaa": m2obj.Group{
					"eaaa": m2obj.Array{float64(1), float64(2), "3"},
				},
			},
			"eb": "test",
		},
		"f": m2obj.Array{
			m2obj.Group{
				"f0": "f0",
			},
			m2obj.Group{
				"f1": "f1",
			},
			m2obj.Group{
				"f2": "f2",
			},
			"f3",
			float64(4),
		},
	})
	_ = os.Remove(filePath)
}

func allNumbersToFloat64(v interface{}) interface{} {
	switch v.(type) {
	case float64:
		return float64(v.(float64))
	case float32:
		return float64(v.(float32))
	case int:
		return float64(v.(int))
	case int8:
		return float64(v.(int8))
	case int16:
		return float64(v.(int16))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	case uint:
		return float64(v.(uint))
	case uint8:
		return float64(v.(uint8))
	case uint16:
		return float64(v.(uint16))
	case uint32:
		return float64(v.(uint32))
	case uint64:
		return float64(v.(uint64))
	}
	if m, ok := v.(map[string]interface{}); ok {
		return allNumbersToFloat64InMap(m)
	}
	if a, ok := v.([]interface{}); ok {
		return allNumbersToFloat64InArray(a)
	}
	return v
}

func allNumbersToFloat64InMap(m map[string]interface{}) map[string]interface{} {
	newM := make(map[string]interface{})
	for k, v := range m {
		newM[k] = allNumbersToFloat64(v)
	}
	return newM
}

func allNumbersToFloat64InArray(a []interface{}) []interface{} {
	newA := make([]interface{}, len(a))
	for i, v := range a {
		newA[i] = allNumbersToFloat64(v)
	}
	return newA
}
