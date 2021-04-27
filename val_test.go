package m2obj

import (
	"github.com/stretchr/testify/assert"
	"math"
	"reflect"
	"testing"
)

// 老实人转换
func TestObject_Vals(t *testing.T) {
	obj := New(groupData{})
	assert.NoError(t, obj.Set("1", 1))
	assert.Equal(t, 1, obj.MustGet("1").ValInt())
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
	arr := obj.MustGet("arr").val.(*arrayData)
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
	assert.Equal(t, map[string]interface{}{
		"val": 1,
	}, obj.Staticize())
}

// 数字乱转
func TestObject_Vals2(t *testing.T) {
	obj := New(nil)
	obj.SetVal(1)
	assert.Equal(t, 1, obj.ValInt())
	assert.Equal(t, byte(1), obj.ValByte())
	assert.Equal(t, int8(1), obj.ValInt8())
	assert.Equal(t, int16(1), obj.ValInt16())
	assert.Equal(t, int32(1), obj.ValInt32())
	assert.Equal(t, int64(1), obj.ValInt64())
	assert.Equal(t, uint64(1), obj.ValUint())
	assert.Equal(t, float32(1), obj.ValFloat32())
	assert.Equal(t, float64(1), obj.ValFloat64())
	obj.SetVal(int8(1))
	assert.Equal(t, 1, obj.ValInt())
	assert.Equal(t, int8(1), obj.ValInt8())
	assert.Equal(t, int16(1), obj.ValInt16())
	assert.Equal(t, int32(1), obj.ValInt32())
	assert.Equal(t, int64(1), obj.ValInt64())
	assert.Equal(t, byte(1), obj.ValByte())
	assert.Equal(t, uint64(1), obj.ValUint())
	assert.Equal(t, float32(1), obj.ValFloat32())
	assert.Equal(t, float64(1), obj.ValFloat64())
	obj.SetVal(int64(1))
	assert.Equal(t, 1, obj.ValInt())
	assert.Equal(t, int8(1), obj.ValInt8())
	assert.Equal(t, int16(1), obj.ValInt16())
	assert.Equal(t, int32(1), obj.ValInt32())
	assert.Equal(t, int64(1), obj.ValInt64())
	assert.Equal(t, byte(1), obj.ValByte())
	assert.Equal(t, uint64(1), obj.ValUint())
	assert.Equal(t, float32(1), obj.ValFloat32())
	assert.Equal(t, float64(1), obj.ValFloat64())
	obj.SetVal(float32(1.8))
	assert.Equal(t, 1, obj.ValInt())
	assert.Equal(t, int8(1), obj.ValInt8())
	assert.Equal(t, int16(1), obj.ValInt16())
	assert.Equal(t, int32(1), obj.ValInt32())
	assert.Equal(t, int64(1), obj.ValInt64())
	assert.Equal(t, uint64(1), obj.ValUint())
	assert.Equal(t, byte(1), obj.ValByte())
	assert.True(t, func() bool {
		except := float32(1.8)
		actual := obj.ValFloat32()
		return math.Abs(float64(except-actual)) < 0.000001
	}())
	assert.True(t, func() bool {
		except := 1.8
		actual := obj.ValFloat64()
		return math.Abs(except-actual) < 0.000001
	}())
	obj.SetVal(1.8)
	assert.Equal(t, 1, obj.ValInt())
	assert.Equal(t, int8(1), obj.ValInt8())
	assert.Equal(t, int16(1), obj.ValInt16())
	assert.Equal(t, int32(1), obj.ValInt32())
	assert.Equal(t, int64(1), obj.ValInt64())
	assert.Equal(t, uint64(1), obj.ValUint())
	assert.Equal(t, byte(1), obj.ValByte())
	assert.True(t, func() bool {
		except := float32(1.8)
		actual := obj.ValFloat32()
		return math.Abs(float64(except-actual)) < 0.000001
	}())
	assert.True(t, func() bool {
		except := 1.8
		actual := obj.ValFloat64()
		return math.Abs(except-actual) < 0.000001
	}())
}

// 字符串乱转
func TestObject_Vals3(t *testing.T) {
	obj := New(nil)
	obj.SetVal("😘哟哟切克闹")
	assert.Equal(t, "😘哟哟切克闹", obj.ValStr())
	assert.Equal(t, []byte("😘哟哟切克闹"), obj.ValBytes())
	assert.Equal(t, []rune("😘哟哟切克闹"), obj.ValRunes())
	obj.SetVal([]byte("😘哟哟切克闹"))
	assert.Equal(t, "😘哟哟切克闹", obj.ValStr())
	assert.Equal(t, []byte("😘哟哟切克闹"), obj.ValBytes())
	assert.Equal(t, []rune("😘哟哟切克闹"), obj.ValRunes())
	obj.SetVal([]rune("😘哟哟切克闹"))
	assert.Equal(t, "😘哟哟切克闹", obj.ValStr())
	assert.Equal(t, []byte("😘哟哟切克闹"), obj.ValBytes())
	assert.Equal(t, []rune("😘哟哟切克闹"), obj.ValRunes())
	// ---- rune ----
	obj2 := New('❤')
	assert.Equal(t, []byte(string('❤')), obj2.ValBytes())
	assert.Equal(t, "❤", obj2.ValStr())
}

func TestObject_Is(t *testing.T) {
	type testType struct {
		A int
	}
	obj := New(Group{
		"haha": "haha",
		"ho": Group{
			"hoho": "hoho",
		},
		"arr": Array{
			1, "2", true,
			testType{1},
		},
	})
	// obj
	assert.True(t, obj.IsGroup())
	assert.False(t, obj.IsArray())
	assert.False(t, obj.IsValue())
	// obj.haha
	assert.True(t, obj.MustGet("haha").IsValue())
	assert.False(t, obj.MustGet("haha").IsGroup())
	assert.False(t, obj.MustGet("haha").IsArray())
	// obj.ho
	assert.True(t, obj.MustGet("ho").IsGroup())
	assert.False(t, obj.MustGet("ho").IsArray())
	assert.False(t, obj.MustGet("ho").IsValue())
	// obj.arr
	assert.True(t, obj.MustGet("arr").IsArray())
	assert.False(t, obj.MustGet("arr").IsGroup())
	assert.False(t, obj.MustGet("arr").IsValue())

	// --------------
	// Is / IsLike
	// --------------

	// nil
	assert.True(t, New(nil).Is(reflect.TypeOf(nil)))
	// Is 不会自动解壳，所以会是*Object而非里面的val的类型，所以这里是false
	assert.False(t, New(nil).Is(reflect.TypeOf(New(nil))))

	// obj
	assert.True(t, obj.Is(reflect.TypeOf(&groupData{})))
	assert.False(t, obj.Is(reflect.TypeOf(New(groupData{}))))
	assert.True(t, obj.IsLike(New(Group{})))
	assert.False(t, obj.IsLike(New(Array{})))
	assert.False(t, obj.IsLike(testType{}))

	// obj.arr.[3]
	assert.True(t, obj.MustGet("arr.[3]").Is(reflect.TypeOf(testType{})))
	assert.False(t, obj.MustGet("arr.[3]").Is(reflect.TypeOf(New(testType{}))))
	assert.True(t, obj.MustGet("arr.[3]").IsLike(testType{}))
	assert.True(t, obj.MustGet("arr.[3]").IsLike(New(testType{})))
	assert.False(t, obj.MustGet("arr.[3]").IsLike(1))
	assert.False(t, obj.MustGet("arr.[3]").IsLike(New(nil)))
}

func TestObject_Is2(t *testing.T) {
	nilObj := New(nil)
	assert.False(t, nilObj.IsGroup())
	assert.False(t, nilObj.IsArray())
	assert.False(t, nilObj.IsValue())
	assert.True(t, nilObj.IsNil())
}
