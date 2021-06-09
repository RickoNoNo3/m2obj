package m2obj

import (
	"regexp"
	"strconv"
)

// arrCheckIndexFormat
//
// Likes *Object.arrCheckIndexKey but only match the format, no verifying on integer transform, no index overflow checking.
func arrCheckIndexFormat(key string) bool {
	reg := regexp.MustCompile(`\[(\d+)]`)
	return reg.MatchString(key)
}

// arrCheckIndexKey
//
// To get an element by index of an Array Object, the keyStr Must be formatted as this:
//     xxx.ArrayName.[index].xxx
// It means that there must be an index statement quoted with '[' and ']' after an Array Object.
//
// This func checks off the rule above.
func (o *Object) arrCheckIndexKey(key, keyStr string) (index int, err error) {
	reg := regexp.MustCompile(`\[(\d+)]`)

	if !reg.MatchString(key) { // the key doesn't be matched as [number]
		err = invalidTypeErr(keyStr)
		return
	} else { // matched
		index, err = strconv.Atoi(reg.FindStringSubmatch(key)[1])
		if err != nil { // the key can not trans to an Integer
			err = invalidKeyStrErr(keyStr)
			return
		} else { // gotten an integer as the index
			arr := *o.val.(*arrayData)
			if index < 0 || len(arr) <= index { // the index overflows from the arr
				err = indexOverflowErr{
					Index: index,
				}
				return
			} else { // no error, check passed
				return index, nil
			}
		}
	}
}

// ArrPush **!!! ONLY FOR ARR OBJECT**
//
// Push a value (or an Object) back into the Array Object.
func (o *Object) ArrPush(value interface{}) {
	switch o.val.(type) {
	case *arrayData:
		*o.val.(*arrayData) = append(*o.val.(*arrayData), New(value))
		o.ArrGet(o.ArrLen() - 1).buildParentLink(o)
		o.callOnChange()
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrPop **!!! ONLY FOR ARR OBJECT**
//
// Pop back from the Array Object.
func (o *Object) ArrPop() (value *Object) {
	switch o.val.(type) {
	case *arrayData:
		value = o.ArrGet(o.ArrLen() - 1)
		*o.val.(*arrayData) = (*o.val.(*arrayData))[:len(*o.val.(*arrayData))-1]
		o.callOnChange()
	default:
		panic(invalidTypeErr(""))
	}
	return
}

// ArrSet **!!! ONLY FOR ARR OBJECT**
//
// Set the value of the element which indexed at the Array Object.
func (o *Object) ArrSet(index int, value interface{}) {
	switch o.val.(type) {
	case *arrayData:
		(*o.val.(*arrayData))[index] = New(value)
		o.ArrGet(index).buildParentLink(o)
		o.callOnChange()
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrGet **!!! ONLY FOR ARR OBJECT**
//
// An alias of `MustGet("[index]")`
func (o *Object) ArrGet(index int) (obj *Object) {
	switch o.val.(type) {
	case *arrayData:
		return (*o.val.(*arrayData))[index]
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrInsert **!!! ONLY FOR ARR OBJECT**
//
// Insert Into the index before the element which indexed at the Array Object.
//
// Specially, if `index == o.ArrLen()` , this is same as `o.Push(value)` but has a lower performance.
func (o *Object) ArrInsert(index int, value interface{}) {
	switch o.val.(type) {
	case *arrayData:
		var (
			arr, arrBefore, arrAfter, arrRes arrayData
		)
		arr = *o.val.(*arrayData)
		// overflow
		if index < 0 || index > len(arr) {
			panic(indexOverflowErr{index})
		}
		// before
		arrBefore = arrayData{}
		if index > 0 {
			arrBefore = append(arrBefore, arr[:index]...)
		}
		// after
		arrAfter = arrayData{}
		if index < len(arr) {
			arrAfter = append(arrAfter, arr[index:]...)
		}
		// generate
		arrRes = append(arrBefore, New(value))
		arrRes = append(arrRes, arrAfter...)
		*o.val.(*arrayData) = arrRes
		o.ArrGet(index).buildParentLink(o)
		o.callOnChange()
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrRemove **!!! ONLY FOR ARR OBJECT**
//
// Remove the element which indexed at the Array Object.
func (o *Object) ArrRemove(index int) {
	switch o.val.(type) {
	case *arrayData:
		var (
			arr, arrBefore, arrAfter, arrRes arrayData
		)
		arr = *o.val.(*arrayData)
		// overflow
		if index < 0 || index >= len(arr) {
			panic(indexOverflowErr{index})
		}
		// before
		arrBefore = arrayData{}
		if index > 0 {
			arrBefore = append(arrBefore, arr[:index]...)
		}
		// after
		arrAfter = arrayData{}
		if index < len(arr)-1 {
			arrAfter = append(arrAfter, arr[index+1:]...)
		}
		// generate
		arrRes = append(arrBefore, arrAfter...)
		*o.val.(*arrayData) = arrRes
		o.callOnChange()
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrShift **!!! ONLY FOR ARR OBJECT**
//
// An alias of `ArrRemove(index)`
func (o *Object) ArrShift() {
	o.ArrRemove(0)
}

// ArrUnshift **!!! ONLY FOR ARR OBJECT**
//
// An alias of `ArrInsert(0, value)`
func (o *Object) ArrUnshift(value interface{}) {
	o.ArrInsert(0, value)
}

// ArrForeach **!!! ONLY FOR ARR OBJECT**
//
// Loop for range `[0...o.ArrLen()-1]`, foreach calls `do`.
//
// Stops when do returns a non-nil err and return it.
func (o *Object) ArrForeach(do func(index int, obj *Object) error) (err error) {
	switch o.val.(type) {
	case *arrayData:
		for i, obj := range *o.val.(*arrayData) {
			if err = do(i, obj); err != nil {
				break
			}
		}
		o.buildParentLink(o.parent)
	default:
		panic(invalidTypeErr(""))
	}
	return
}

// ArrMerge **!!! ONLY FOR ARR OBJECT**
//
// Push all of elements from an Array Object o2 into the Array Object
func (o *Object) ArrMerge(o2 *Object) {
	switch o.val.(type) {
	case *arrayData:
		switch o2.val.(type) {
		case *arrayData: // Group
			newArr := o.Clone()
			err := o2.ArrForeach(func(index int, obj *Object) error {
				newArr.ArrPush(obj)
				return nil
			})
			if err == nil {
				o.SetVal(newArr)
				o.buildParentLink(o.parent)
			}
		default:
			panic(invalidTypeErr(""))
		}
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrPushAll **!!! ONLY FOR ARR OBJECT**
//
// Push all the elements in the parameter to the array object
func (o *Object) ArrPushAll(values ...interface{}) {
	switch o.val.(type) {
	case *arrayData:
		o2 := New(Array(values))
		o.ArrMerge(o2)
	default:
		panic(invalidTypeErr(""))
	}
}

// ArrLen **!!! ONLY FOR ARR OBJECT**
//
// Get the length of an array object
func (o *Object) ArrLen() int {
	switch o.val.(type) {
	case *arrayData:
		return len(*o.val.(*arrayData))
	default:
		panic(invalidTypeErr(""))
	}
}
