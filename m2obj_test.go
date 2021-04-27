package m2obj

import (
	"fmt"
	"reflect"
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
			invalidKeyStrErr("keyStr"),
			"invalid key string: keyStr",
		},
		{
			unknownTypeErr("key"),
			"the key {key} has an unknown ObjectType",
		},
		{
			invalidTypeErr("key"),
			"the key {key} has an invalid ObjectType",
		},
		{
			unknownTypeErr(""),
			"unknown ObjectType",
		},
		{
			invalidTypeErr(""),
			"invalid ObjectType",
		},
	}
	for _, data := range testData {
		assert.EqualError(t, data.err, data.wantStr)
	}
}

func TestNewAndStaticize(t *testing.T) {
	assert.NotPanics(t, func() {
		m := New(nil).Staticize()
		assert.Equal(t, map[string]interface{}{
			"val": nil,
		}, m)
	})
	// -----------
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
		testData[i].obj = New(testData[i].wantMap)
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
				"map": New(map[string]interface{}{
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
	staticMap := map[string]interface{}{
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
	}
	obj2 := obj.Clone()
	t.Run("Test Pointer Address", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.NotEqual(t, reflect.ValueOf(obj2).Pointer(), reflect.ValueOf(obj).Pointer())
			assert.NotEqual(t, reflect.ValueOf(obj2.MustGet("arr")).Pointer(), reflect.ValueOf(obj.MustGet("arr")))
			assert.NotEqual(t, reflect.ValueOf(obj2.MustGet("arr.[2]")).Pointer(), reflect.ValueOf(obj.MustGet("arr.[2]")))
			assert.NotEqual(t, reflect.ValueOf(obj2.MustGet("arr.[2].1")).Pointer(), reflect.ValueOf(obj.MustGet("arr.[2].1")).Pointer())
			assert.NotEqual(t, reflect.ValueOf(obj2.MustGet("1")).Pointer(), reflect.ValueOf(obj.MustGet("1")).Pointer())
		})
	})
	t.Run("Test Edit Source In Arr", func(t *testing.T) {
		assert.NotPanics(t, func() {
			assert.NoError(t, obj.Set("arr.[2].1", "arr.[2].1 - test"))
			assert.Equal(t, "arr.[2].1 - test", obj.MustGet("arr.[2].1").ValStr())
			assert.Equal(t, 1, obj2.MustGet("arr.[2].1").ValInt())
			assert.NoError(t, obj.MustGet("arr").ArrPushAll(10, 11, 12))
			assert.Greater(t, len(*obj.MustGet("arr").val.(*arrayData)), 5)
			assert.Equal(t, len(*obj2.MustGet("arr").val.(*arrayData)), 3)
		})
		assert.Equal(t, staticMap, obj2.Staticize())
	})
	t.Run("Test Edit Source By Merge", func(t *testing.T) {
		assert.NoError(t, obj.GroupMerge(New(Group{
			"arr": New(Array{
				"world",
				0,
				Group{
					"1": nil,
					"2": nil,
					"3": nil,
				},
			}),
			"1": 2233,
			"2": "2233",
			"3": false,
		}), true))
		assert.Equal(t, staticMap, obj2.Staticize())
	})
	t.Run("Test Remove Source", func(t *testing.T) {
		obj.SetVal(nil)
		assert.Equal(t, map[string]interface{}{
			"val": nil,
		}, obj.Staticize())
		assert.Equal(t, staticMap, obj2.Staticize())
	})
}

func TestObject_Parent(t *testing.T) {
	obj := New(Group{
		"a": Group{
			"b": Group{
				"c": "c",
			},
			"d": New("d"),
		},
		"e": "e",
		"f": Array{
			Array{"g", "h"},
			Group{"i": "i"},
			map[string]interface{}{
				"j": "j",
				"k": []interface{}{"l", "m"},
			},
			"n",
		},
	})
	assert.NoError(
		t,
		obj.Set("o.p.q", Group{
			"r": map[string]interface{}{
				"s": map[string]interface{}{
					"t": "t",
				},
			},
			"u": Array{
				"v",
				Group{
					"w": "w",
				},
			},
		}),
	)
	assert.NoError(
		t,
		obj.Set(
			"o.p.q.u.[1].x.y.z",
			"z",
		),
	)
	assert.NotPanics(t, func() {
		a := obj.MustGet("a")
		b := a.MustGet("b")
		c := b.MustGet("c")
		d := a.MustGet("d")
		e := obj.MustGet("e")
		f := obj.MustGet("f")
		f0 := f.MustGet("[0]")
		g := f0.MustGet("[0]")
		h := f0.MustGet("[1]")
		f1 := f.MustGet("[1]")
		i := f1.MustGet("i")
		f2 := f.MustGet("[2]")
		j := f2.MustGet("j")
		k := f2.MustGet("k")
		l := k.MustGet("[0]")
		m := k.MustGet("[1]")
		n := f.MustGet("[3]")
		o := obj.MustGet("o")
		p := o.MustGet("p")
		q := p.MustGet("q")
		r := q.MustGet("r")
		s := r.MustGet("s")
		T := s.MustGet("t")
		u := q.MustGet("u")
		v := u.MustGet("[0]")
		u1 := u.MustGet("[1]")
		w := u1.MustGet("w")
		x := u1.MustGet("x")
		y := x.MustGet("y")
		z := y.MustGet("z")
		assert.Equal(t, obj, a.Parent())
		assert.Equal(t, a, b.Parent())
		assert.Equal(t, b, c.Parent())
		assert.Equal(t, a, d.Parent())
		assert.Equal(t, obj, e.Parent())
		assert.Equal(t, obj, f.Parent())
		assert.Equal(t, f, f0.Parent())
		assert.Equal(t, f0, g.Parent())
		assert.Equal(t, f0, h.Parent())
		assert.Equal(t, f, f1.Parent())
		assert.Equal(t, f1, i.Parent())
		assert.Equal(t, f, f2.Parent())
		assert.Equal(t, f2, j.Parent())
		assert.Equal(t, f2, k.Parent())
		assert.Equal(t, k, l.Parent())
		assert.Equal(t, k, m.Parent())
		assert.Equal(t, f, n.Parent())
		assert.Equal(t, obj, o.Parent())
		assert.Equal(t, o, p.Parent())
		assert.Equal(t, p, q.Parent())
		assert.Equal(t, q, r.Parent())
		assert.Equal(t, r, s.Parent())
		assert.Equal(t, s, T.Parent())
		assert.Equal(t, q, u.Parent())
		assert.Equal(t, u, v.Parent())
		assert.Equal(t, u, u1.Parent())
		assert.Equal(t, u1, w.Parent())
		assert.Equal(t, u1, x.Parent())
		assert.Equal(t, x, y.Parent())
		assert.Equal(t, y, z.Parent())
	})
}
