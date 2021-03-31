package m2obj

import (
	"regexp"
	"strconv"
)

// like *Object.arrCheckIndexKey but only match the format, no verifying on integer transform, no index overflow checking.
func arrCheckIndexFormat(key string) bool {
	reg := regexp.MustCompile(`\[(\d+)]`)
	return reg.MatchString(key)
}

// To get an element by index of an Array Object, the keyStr Must be formatted as this:
//     xxx.ArrayName.[index].xxx
// It means that there must be an index statement quoted with '[' and ']' after an Array Object.
//
// This func checks off the rule above.
func (o *Object) arrCheckIndexKey(key, keyStr string) (index int, err error) {
	reg := regexp.MustCompile(`\[(\d+)]`)

	if !reg.MatchString(key) { // the key doesn't be matched as [number]
		err = InvalidKeyStrErr(keyStr)
		return
	} else { // matched
		index, err = strconv.Atoi(reg.FindStringSubmatch(key)[1])
		if err != nil { // the key can not trans to an Integer
			err = InvalidKeyStrErr(keyStr)
			return
		} else { // gotten an integer as the index
			arr := *o.val.(*ArrayData)
			if len(arr) <= index { // the index overflows from the arr
				err = IndexOverflowErr{
					Index: index,
				}
				return
			} else { // no error, check passed
				return index, nil
			}
		}
	}
}

func (o *Object) ArrPush(value interface{}) (err error) {
	switch o.val.(type) {
	case *ArrayData:
		*o.val.(*ArrayData) = append(*o.val.(*ArrayData), New(getDeepestValue(value)))
		return nil
	default:
		return InvalidTypeErr("")
	}
}

func (o *Object) ArrPop() (err error) {
	switch o.val.(type) {
	case *ArrayData:
		*o.val.(*ArrayData) = (*o.val.(*ArrayData))[:len(*o.val.(*ArrayData))-1]
		return nil
	default:
		return InvalidTypeErr("")
	}
}

func (o *Object) ArrSet(index int, value interface{}) (err error) {
	switch o.val.(type) {
	case *ArrayData:
		(*o.val.(*ArrayData))[index] = New(getDeepestValue(value))
		return nil
	default:
		return InvalidTypeErr("")
	}
}

// An alias of `*Object.Get("[index]")`
func (o *Object) ArrGet(index int) (obj *Object, err error) {
	switch o.val.(type) {
	case *ArrayData:
		return (*o.val.(*ArrayData))[index], nil
	default:
		return nil, InvalidTypeErr("")
	}
}

func (o *Object) ArrInsert(index int, value interface{}) (err error) {
	switch o.val.(type) {
	case *ArrayData:
		var (
			arr, arrBefore, arrAfter, arrRes ArrayData
		)
		arr = *o.val.(*ArrayData)
		// overflow
		if index < 0 || index >= len(arr) {
			return IndexOverflowErr{index}
		}
		// before
		arrBefore = ArrayData{}
		if index > 0 {
			arrBefore = append(arrBefore, arr[:index]...)
		}
		// after
		arrAfter = ArrayData{}
		if index < len(arr)-1 {
			arrAfter = append(arrAfter, arr[index:]...)
		}
		// generate
		arrRes = append(arrBefore, New(getDeepestValue(value)))
		arrRes = append(arrRes, arrAfter...)
		*o.val.(*ArrayData) = arrRes
		return nil
	default:
		return InvalidTypeErr("")
	}
}

func (o *Object) ArrRemove(index int) (err error) {
	switch o.val.(type) {
	case *ArrayData:
		var (
			arr, arrBefore, arrAfter, arrRes ArrayData
		)
		arr = *o.val.(*ArrayData)
		// overflow
		if index < 0 || index >= len(arr) {
			return IndexOverflowErr{index}
		}
		// before
		arrBefore = ArrayData{}
		if index > 0 {
			arrBefore = append(arrBefore, arr[:index]...)
		}
		// after
		arrAfter = ArrayData{}
		if index < len(arr)-1 {
			arrAfter = append(arrAfter, arr[index+1:]...)
		}
		// generate
		arrRes = append(arrBefore, arrAfter...)
		*o.val.(*ArrayData) = arrRes
		return nil
	default:
		return InvalidTypeErr("")
	}
}

func (o *Object) ArrForeach(do func(index int, obj *Object)) (err error) {
	switch o.val.(type) {
	case *ArrayData:
		for i, obj := range *o.ValArr() {
			do(i, obj)
		}
		return nil
	default:
		return InvalidTypeErr("")
	}
}

func (o *Object) ArrLen() (int, error) {
	switch o.val.(type) {
	case *ArrayData:
		return len(*o.val.(*ArrayData)), nil
	default:
		return -1, InvalidTypeErr("")
	}
}
