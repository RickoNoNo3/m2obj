package m2obj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArrs(t *testing.T) {
	obj := New(Group{
		"arr": Array{
			"hello",
			2,
			Group{
				"1": 1,
			},
		},
	})
	// errors
	t.Run("Errors Test", func(t *testing.T) {
		assert.Panics(t, func() {
			obj.ArrPush(100)
		})
		assert.Panics(t, func() {
			obj.ArrPop()
		})
		assert.Panics(t, func() {
			obj.ArrSet(0, 100)
		})
		assert.Panics(t, func() {
			obj.ArrGet(0)
		})
		assert.Panics(t, func() {
			obj.ArrInsert(0, 100)
		})
		assert.Panics(t, func() {
			obj.ArrRemove(0)
		})
		assert.Panics(t, func() {
			obj.ArrForeach(func(i int, obj *Object) error { return nil })
		})
		assert.Panics(t, func() {
			obj.ArrLen()
		})
	})
	var arr *Object
	assert.NotPanics(t, func() {
		arr = obj.MustGet("arr")
	})
	// push/len
	assert.NotPanics(t, func() {
		arr.ArrPush(100)
	})
	assert.Equal(t, 4, arr.ArrLen())
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
	assert.NotPanics(t, func() {
		arr.ArrPop()
	})
	assert.Equal(t, 3, arr.ArrLen())
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
	assert.NotPanics(t, func() {
		assert.Equal(t, map[string]interface{}{
			"1": 1,
		}, arr.MustGet("[2]").Staticize())
		arr.ArrSet(2, 100)
		assert.Equal(t, 100, arr.MustGet("[2]").ValInt())
		assert.Equal(t, 100, arr.ArrGet(2).ValInt())
	})
	// insert
	assert.Panics(t, func() {
		arr.ArrInsert(arr.ArrLen()+1, "awd")
	})
	assert.Panics(t, func() {
		arr.ArrInsert(-1, "awd")
	})
	assert.NotPanics(t, func() {
		arr.ArrInsert(arr.ArrLen(), "awd")
		arr.ArrInsert(0, "awd")
		arr.ArrInsert(2, "awd")
		assert.Equal(t, map[string]interface{}{
			"arr": []interface{}{
				"awd",
				"hello",
				"awd",
				2,
				100,
				"awd",
			},
		}, obj.Staticize())
	})
	// remove
	assert.Panics(t, func() {
		arr.ArrRemove(arr.ArrLen())
	})
	assert.Panics(t, func() {
		arr.ArrRemove(-1)
	})
	assert.NotPanics(t, func() {
		arr.ArrRemove(0)
		arr.ArrRemove(0)
		arr.ArrRemove(3)
		arr.ArrRemove(1)
		assert.Equal(t, map[string]interface{}{
			"arr": []interface{}{
				"awd",
				100,
			},
		}, obj.Staticize())
	})
	// foreach
	assert.NotPanics(t, func() {
		assert.NoError(t, arr.ArrForeach(func(index int, obj *Object) error {
			obj.SetVal(false)
			return nil
		}))
		assert.Equal(t, map[string]interface{}{
			"arr": []interface{}{
				false,
				false,
			},
		}, obj.Staticize())
	})
}

func TestArrs2(t *testing.T) {
	arr := New(Array{
		"1",
		2,
		true,
	})
	assert.NotPanics(t, func() {
		arr.ArrShift()
	})
	assert.Equal(t, New(Array{
		2,
		true,
	}).Staticize(), arr.Staticize())
	assert.NotPanics(t, func() {
		arr.ArrUnshift("1")
	})
	assert.Equal(t, New(Array{
		"1",
		2,
		true,
	}).Staticize(), arr.Staticize())
}

func TestObject_ArrPushArr(t *testing.T) {
	arr1 := New(Array{1, 2, 3})
	arr2 := New(Array{4, 5, "6"})
	obj := New(Group{"1": 1})
	assert.Panics(t, func() {
		arr1.ArrMerge(obj)
	})
	assert.Panics(t, func() {
		obj.ArrMerge(arr1)
	})
	assert.NotPanics(t, func() {
		arr1.ArrMerge(arr2)
	})
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, "6"}).Staticize(), arr1.Staticize())
	assert.NotPanics(t, func() {
		arr1.ArrMerge(New([]interface{}{7, 8, 9}))
	})
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, "6", 7, 8, 9}).Staticize(), arr1.Staticize())
}

func TestObject_ArrPushAll(t *testing.T) {
	arr1 := New(Array{1, 2, 3})
	arr2 := []int{4, 5, 6}
	obj := New(Group{"1": 1})
	// type error
	assert.Panics(t, func() {
		obj.ArrPushAll(4, 5, 6)
	})

	// ... add
	assert.NotPanics(t, func() {
		arr1.ArrPushAll(4, 5, 6)
	})
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, 6}).Staticize(), arr1.Staticize())

	// slice to []interface{} then ... add
	arr1 = New(Array{1, 2, 3})
	arrInterface2 := make([]interface{}, len(arr2))
	for i, v := range arr2 {
		arrInterface2[i] = v
	}
	assert.NotPanics(t, func() {
		arr1.ArrPushAll(arrInterface2...)
	})
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, 6}).Staticize(), arr1.Staticize())

	// slice for each add
	arr1 = New(Array{1, 2, 3})
	for _, v := range arr2 {
		assert.NotPanics(t, func() {
			arr1.ArrPush(v)
		})
	}
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, 6}).Staticize(), arr1.Staticize())
}

func TestObject_ArrForeach(t *testing.T) {
	var (
		arr1 *Object
		arr2 *Object
		arr3 []int
	)
	arr1 = New(Array{1, 2, 3})
	arr2 = New(Array{1, 2, "3"})
	arr3 = make([]int, 0)
	// no error and no panic
	assert.NotPanics(t, func() {
		assert.NoError(t, arr1.ArrForeach(func(index int, obj *Object) error {
			arr3 = append(arr3, obj.ValInt())
			return nil
		}))
	})
	assert.Equal(t, []int{1, 2, 3}, arr3)
	arr3 = make([]int, 0)
	// transform the panic to error
	assert.NotPanics(t, func() {
		assert.Error(t, arr2.ArrForeach(func(index int, obj *Object) error {
			switch obj.Val().(type) {
			case int:
				arr3 = append(arr3, obj.ValInt())
				return nil
			default:
				return invalidTypeErr("")
			}
		}))
	})
	assert.Equal(t, []int{1, 2}, arr3)
	arr3 = make([]int, 0)
	// throw the panic directly
	assert.Panics(t, func() {
		assert.NoError(t, arr2.ArrForeach(func(index int, obj *Object) error {
			arr3 = append(arr3, obj.ValInt())
			return nil
		}))
	})
	assert.Equal(t, []int{1, 2}, arr3)
	assert.Panics(t, func() {
		New(arr3).ArrForeach(func(index int, obj *Object) error {
			return nil
		})
	})
}
