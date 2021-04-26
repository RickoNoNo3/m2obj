package m2obj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		assert.Error(t, obj.ArrForeach(func(i int, obj *Object) error { return nil }))
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
}

func TestArrs2(t *testing.T) {
	arr := New(Array{
		"1",
		2,
		true,
	})
	assert.NoError(t, arr.ArrShift())
	assert.Equal(t, New(Array{
		2,
		true,
	}).Staticize(), arr.Staticize())
	assert.NoError(t, arr.ArrUnshift("1"))
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
	assert.Error(t, arr1.ArrPushArr(obj))
	assert.Error(t, obj.ArrPushArr(arr1))
	assert.NoError(t, arr1.ArrPushArr(arr2))
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, "6"}).Staticize(), arr1.Staticize())
	assert.NoError(t, arr1.ArrPushArr(New([]interface{}{7, 8, 9})))
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, "6", 7, 8, 9}).Staticize(), arr1.Staticize())
}

func TestObject_ArrPushAll(t *testing.T) {
	arr1 := New(Array{1, 2, 3})
	arr2 := []int{4, 5, 6}
	obj := New(Group{"1": 1})
	// type error
	assert.Error(t, obj.ArrPushAll(4, 5, 6))

	// ... add
	assert.NoError(t, arr1.ArrPushAll(4, 5, 6))
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, 6}).Staticize(), arr1.Staticize())

	// slice to []interface{} then ... add
	arr1 = New(Array{1, 2, 3})
	arrInterface2 := make([]interface{}, len(arr2))
	for i, v := range arr2 {
		arrInterface2[i] = v
	}
	assert.NoError(t, arr1.ArrPushAll(arrInterface2...))
	assert.Equal(t, New(Array{1, 2, 3, 4, 5, 6}).Staticize(), arr1.Staticize())

	// slice for each add
	arr1 = New(Array{1, 2, 3})
	for _, v := range arr2 {
		assert.NoError(t, arr1.ArrPush(v))
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
	assert.Error(t, New(arr3).ArrForeach(func(index int, obj *Object) error {
		return nil
	}))
}
