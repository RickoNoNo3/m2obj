package m2obj

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	type TestData struct {
		err     error
		wantStr string
	}
	testData := []TestData{
		{
			indexOverflowErr{
				Index: 10,
			},
			"no such index[10]",
		},
		{
			invalidKeyStrErr("key"),
			"invalid key string: key",
		},
		{
			unknownTypeErr("key"),
			"the key {key} has an unknown ObjectType",
		},
		{
			invalidTypeErr("key"),
			"the key {key} has an invalid ObjectType",
		},
	}
	for _, data := range testData {
		assert.EqualError(t, data.err, data.wantStr)
	}
}

func TestNewAndStaticize(t *testing.T) {
	type TestData struct {
		obj     *Object
		wantMap map[string]interface{}
	}
	testData := []TestData{
		// New Value
		{
			New("abc"),
			map[string]interface{}{
				"val": "abc",
			},
		},
		{
			New("abc123123123"),
			map[string]interface{}{
				"val": "abc123123123",
			},
		},
		{
			New(3),
			map[string]interface{}{
				"val": 3,
			},
		},
		// New Group
		{
			New(groupData{}),
			map[string]interface{}{},
		},
		{
			New(groupData{
				"a": New(1),
				"b": New("2"),
				"c": New(true),
			}),
			map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
			},
		},
		{
			New(Group{
				"a": New(1),
				"b": New("2"),
				"c": New(true),
			}),
			map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
			},
		},
		{
			New(Group{
				"a": 1,
				"b": "2",
				"c": true,
			}),
			map[string]interface{}{
				"a": 1,
				"b": "2",
				"c": true,
			},
		},
		// New Array
		{
			New(arrayData{}),
			map[string]interface{}{
				"list": []interface{}{},
			},
		},
		{
			New(arrayData{
				0: New(1),
				1: New("2"),
				3: New(true),
			}),
			map[string]interface{}{
				"list": []interface{}{
					0: 1,
					1: "2",
					2: nil,
					3: true,
				},
			},
		},
		{
			New(Array{
				0: New(1),
				1: New("2"),
				3: New(true),
			}),
			map[string]interface{}{
				"list": []interface{}{
					0: 1,
					1: "2",
					2: nil,
					3: true,
				},
			},
		},
		{
			New(Array{
				0: 1,
				1: "2",
				3: true,
			}),
			map[string]interface{}{
				"list": []interface{}{
					0: 1,
					1: "2",
					2: nil,
					3: true,
				},
			},
		},
		// New From Map
		{
			wantMap: map[string]interface{}{
				"val": 35,
			},
		},
		{
			wantMap: map[string]interface{}{
				"list": []interface{}{},
			},
		},
		{
			wantMap: map[string]interface{}{
				"list": []interface{}{
					1: 1,
					2: "2",
					3: true,
				},
				"map": map[string]interface{}{
					"1": 1,
					"2": "2",
					"3": true,
				},
			},
		},
	}
	for i := len(testData) - 3; i < len(testData); i++ {
		testData[i].obj = NewFromMap(testData[i].wantMap)
	}
	for i, data := range testData {
		fmt.Printf("TestNewAndStaticize: testData[%d]\n", i)
		assert.Equal(t, data.wantMap, data.obj.Staticize())
	}
}

func TestObject_SetGetRemove(t *testing.T) {
	obj := New(groupData{
		"a": New(Object{val: 3}),
		"b": New(arrayData{
			New(1),
			New("2"),
			New(groupData{
				"map": NewFromMap(map[string]interface{}{
					"a": "1",
					"b": 2,
					"c": []interface{}{1, 2, 3},
				}),
			}),
		}),
	})
	assert.Equal(t, map[string]interface{}{
		"a": 3,
		"b": []interface{}{
			1,
			"2",
			map[string]interface{}{
				"map": map[string]interface{}{
					"a": "1",
					"b": 2,
					"c": []interface{}{1, 2, 3},
				},
			},
		},
	}, obj.Staticize())
	assert.NoError(t, obj.Set("a", 10))
	assert.NoError(t, obj.Set("b.[1]", 15))
	assert.NoError(t, obj.SetIfHas("c", 20))
	assert.NoError(t, obj.SetIfNotHas("d", 25))
	assert.NoError(t, obj.Set("b.[2].map.b", 30))
	assert.NoError(t, obj.Set("b.[2].map.c.[2]", 35))
	assert.Equal(t, map[string]interface{}{
		"a": 10,
		"b": []interface{}{
			1,
			15,
			map[string]interface{}{
				"map": map[string]interface{}{
					"a": "1",
					"b": 30,
					"c": []interface{}{1, 2, 35},
				},
			},
		},
		"d": 25,
	}, obj.Staticize())
	assert.True(t, obj.Remove("a"))
	assert.True(t, obj.Remove("c"))
	assert.True(t, obj.Remove("d"))
	assert.False(t, obj.Remove("b.[0]"))
	assert.NoError(t, obj.Set("b.[2].map", nil))
	assert.NoError(t, obj.Set("b.[2].nil", nil))
	assert.NoError(t, obj.Set("b.[2].string", "string"))
	assert.NoError(t, obj.Set("b.[2].int", int64(2000)))
	assert.NoError(t, obj.Set("b.[2].group", New(groupData{
		"a": New(1),
		"b": New(New(2)),
		"c": New(New(New(3))),
	})))
	assert.NoError(t, obj.Set("b.[2].array", arrayData{
		New("哈哈"),
		New("吼吼"),
		New("嘿嘿"),
	}))
	assert.Equal(t, map[string]interface{}{
		"b": []interface{}{
			1,
			15,
			map[string]interface{}{
				"map":    nil,
				"nil":    nil,
				"string": "string",
				"int":    int64(2000),
				"group": map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": 3,
				},
				"array": []interface{}{
					"哈哈",
					"吼吼",
					"嘿嘿",
				},
			},
		},
	}, obj.Staticize())
	assert.Panics(t, func() {
		obj.MustGet("b.[99999]")
	})
	assert.Panics(t, func() {
		obj.MustGet("b.[]")
	})
	assert.Panics(t, func() {
		obj.MustGet("b.[hello]")
	})
	assert.Panics(t, func() {
		obj.MustGet("b.hello")
	})
	assert.NotPanics(t, func() {
		obj.MustGet("b.[0]")
	})
	assert.Equal(t, 1, obj.MustGet("b.[0]").ValInt())
	assert.Equal(t, 15, obj.MustGet("b.[1]").ValInt())
	assert.Equal(t, int64(2000), obj.MustGet("b.[2].int").ValInt64())
	assert.Equal(t, "哈哈", obj.MustGet("b.[2].array.[0]").ValStr())

	assert.NoError(t, obj.Set("aa.bb.cc.dd.ee.ff", 1))
	assert.Equal(t, 1, obj.MustGet("aa.bb.cc.dd.ee.ff").ValInt())
	assert.False(t, obj.Remove(""))
	assert.NoError(t, obj.MustGet("aa").Set("", 1))
	assert.Equal(t, 1, obj.MustGet("aa").ValInt())
}

func TestObject_SetGetRemove2(t *testing.T) {
	obj := New(groupData{})
	assert.NoError(t, obj.Set("aa.bb.cc.dd.ee.ff", 1))
	assert.NoError(t, obj.Set("aa...bb....cc..dd.ee...ff", 2))
	assert.Equal(t, 2, obj.MustGet("aa..bb..cc..dd..ee..ff").ValInt())
	assert.Error(t, obj.Set("aa.bb.cc.dd.ee.ff.gg", int8(10)))
	assert.NoError(t, obj.SetIfHas("aa.bb.cc.dd.ee.ff", int8(20)))
	assert.NoError(t, obj.SetIfNotHas("aa.bb.cc.dd.ee.ff", int8(30)))
	assert.Equal(t, map[string]interface{}{
		"ff": int8(20),
	}, obj.MustGet("aa.bb.cc.dd.ee").Staticize())
}

func TestObject_Vals(t *testing.T) {
	obj := New(groupData{})
	assert.NoError(t, obj.Set("1", int(1)))
	assert.Equal(t, int(1), obj.MustGet("1").ValInt())
	assert.NoError(t, obj.Set("8", int8(1)))
	assert.Equal(t, int8(1), obj.MustGet("8").ValInt8())
	assert.NoError(t, obj.Set("16", int16(1)))
	assert.Equal(t, int16(1), obj.MustGet("16").ValInt16())
	assert.NoError(t, obj.Set("32", int32(1)))
	assert.Equal(t, int32(1), obj.MustGet("32").ValInt32())
	assert.NoError(t, obj.Set("64", int64(1)))
	assert.Equal(t, int64(1), obj.MustGet("64").ValInt64())
	assert.NoError(t, obj.Set("f32", float32(1)))
	assert.Equal(t, float32(1), obj.MustGet("f32").ValFloat32())
	assert.NoError(t, obj.Set("f64", float64(1)))
	assert.Equal(t, float64(1), obj.MustGet("f64").ValFloat64())
	assert.NoError(t, obj.Set("bool", true))
	assert.Equal(t, true, obj.MustGet("bool").ValBool())
	assert.NoError(t, obj.Set("str", "str"))
	assert.Equal(t, "str", obj.MustGet("str").ValStr())
	assert.NoError(t, obj.Set("arr", arrayData{
		New(1),
		New("2"),
	}))
	arr := obj.MustGet("arr").valArr()
	assert.Equal(t, 1, (*arr)[0].ValInt())
	assert.Equal(t, "2", (*arr)[1].ValStr())
	type Tmp struct {
		A int
	}
	assert.NoError(t, obj.Set("custom", Tmp{
		A: 10,
	}))
	assert.Equal(t, 10, obj.MustGet("custom").Val().(Tmp).A)
	obj.SetVal(1)
	assert.Equal(t, 1, obj.ValInt())
	assert.Panics(t, func() {
		obj.ValBool()
	})
	assert.Panics(t, func() {
		obj.ValInt8()
	})
	assert.Panics(t, func() {
		obj.ValStr()
	})
	assert.Panics(t, func() {
		obj.valArr()
	})
	assert.Equal(t, map[string]interface{}{
		"val": 1,
	}, obj.Staticize())

}

func TestObject_Clone(t *testing.T) {
	obj := New(groupData{
		"arr": New(arrayData{
			New("hello"),
			New(2),
			New(groupData{
				"1": New(1),
				"2": New("2"),
				"3": New(true),
			}),
		}),
		"1": New(1),
		"2": New("2"),
		"3": New(true),
	})
	obj2 := obj.Clone()
	obj.SetVal(nil)
	assert.Equal(t, map[string]interface{}{
		"val": nil,
	}, obj.Staticize())
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			"hello",
			2,
			map[string]interface{}{
				"1": 1,
				"2": "2",
				"3": true,
			},
		},
		"1": 1,
		"2": "2",
		"3": true,
	}, obj2.Staticize())
}

func TestArrs(t *testing.T) {
	obj := New(groupData{
		"arr": New(arrayData{
			New("hello"),
			New(2),
			New(groupData{
				"1": New(1),
			}),
		}),
	})
	// errors
	t.Run("Invalid Type Errors Test", func(t *testing.T) {
		var err error
		assert.Error(t, obj.ArrPush(100))
		assert.Error(t, obj.ArrPop())
		assert.Error(t, obj.ArrSet(0, 100))
		_, err = obj.ArrGet(0)
		assert.Error(t, err)
		assert.Error(t, obj.ArrInsert(0, 100))
		assert.Error(t, obj.ArrRemove(0))
		assert.Error(t, obj.ArrForeach(func(i int, obj *Object) {}))
		_, err = obj.ArrLen()
		assert.Error(t, err)
	})
	arr := obj.MustGet("arr")
	// push/len
	assert.NoError(t, arr.ArrPush(100))
	arrLen, err := arr.ArrLen()
	assert.NoError(t, err)
	assert.Equal(t, 4, arrLen)
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			"hello",
			2,
			map[string]interface{}{
				"1": 1,
			},
			100,
		},
	}, obj.Staticize())
	// pop
	assert.NoError(t, arr.ArrPop())
	arrLen, err = arr.ArrLen()
	assert.NoError(t, err)
	assert.Equal(t, 3, arrLen)
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			"hello",
			2,
			map[string]interface{}{
				"1": 1,
			},
		},
	}, obj.Staticize())
	// set/get
	assert.Equal(t, map[string]interface{}{
		"1": 1,
	}, arr.MustGet("[2]").Staticize())
	assert.NoError(t, arr.ArrSet(2, 100))
	assert.Equal(t, 100, arr.MustGet("[2]").ValInt())
	gotten, err := arr.ArrGet(2)
	assert.NoError(t, err)
	assert.Equal(t, 100, gotten.ValInt())
	// insert
	assert.Error(t, arr.ArrInsert(3, "awd"))
	assert.NoError(t, arr.ArrInsert(0, "awd"))
	assert.NoError(t, arr.ArrInsert(2, "awd"))
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			"awd",
			"hello",
			"awd",
			2,
			100,
		},
	}, obj.Staticize())
	// remove
	assert.Error(t, arr.ArrRemove(-1))
	assert.NoError(t, arr.ArrRemove(0))
	assert.NoError(t, arr.ArrRemove(0))
	assert.Error(t, arr.ArrRemove(3))
	assert.NoError(t, arr.ArrRemove(1))
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			"awd",
			100,
		},
	}, obj.Staticize())
	// foreach
	assert.NoError(t, arr.ArrForeach(func(index int, obj *Object) {
		obj.SetVal(false)
	}))
	assert.Equal(t, map[string]interface{}{
		"arr": []interface{}{
			false,
			false,
		},
	}, obj.Staticize())
}
