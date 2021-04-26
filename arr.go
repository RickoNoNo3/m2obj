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
		err = invalidKeyStrErr(keyStr)
		return
	} else { // matched
		index, err = strconv.Atoi(reg.FindStringSubmatch(key)[1])
		if err != nil { // the key can not trans to an Integer
			err = invalidKeyStrErr(keyStr)
			return
		} else { // gotten an integer as the index
			arr := *o.val.(*arrayData)
			if len(arr) <= index { // the index overflows from the arr
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

func (o *Object) ArrPush(value interface{}) (err error) {
	switch o.val.(type) {
	case *arrayData:
		*o.val.(*arrayData) = append(*o.val.(*arrayData), New(value))
		o.buildParentLink(o.parent)
		o.callOnChange()
		return nil
	default:
		return invalidTypeErr("")
	}
}

func (o *Object) ArrPop() (err error) {
	switch o.val.(type) {
	case *arrayData:
		*o.val.(*arrayData) = (*o.val.(*arrayData))[:len(*o.val.(*arrayData))-1]
		o.buildParentLink(o.parent)
		o.callOnChange()
		return nil
	default:
		return invalidTypeErr("")
	}
}

func (o *Object) ArrSet(index int, value interface{}) (err error) {
	switch o.val.(type) {
	case *arrayData:
		(*o.val.(*arrayData))[index] = New(value)
		o.buildParentLink(o.parent)
		o.callOnChange()
		return nil
	default:
		return invalidTypeErr("")
	}
}

// ArrGet is an alias of `*Object.Get("[index]")`
func (o *Object) ArrGet(index int) (obj *Object, err error) {
	switch o.val.(type) {
	case *arrayData:
		return (*o.val.(*arrayData))[index], nil
	default:
		return nil, invalidTypeErr("")
	}
}

func (o *Object) ArrInsert(index int, value interface{}) (err error) {
	switch o.val.(type) {
	case *arrayData:
		var (
			arr, arrBefore, arrAfter, arrRes arrayData
		)
		arr = *o.val.(*arrayData)
		// overflow
		if index < 0 || index >= len(arr) {
			return indexOverflowErr{index}
		}
		// before
		arrBefore = arrayData{}
		if index > 0 {
			arrBefore = append(arrBefore, arr[:index]...)
		}
		// after
		arrAfter = arrayData{}
		if index < len(arr)-1 {
			arrAfter = append(arrAfter, arr[index:]...)
		}
		// generate
		arrRes = append(arrBefore, New(value))
		arrRes = append(arrRes, arrAfter...)
		*o.val.(*arrayData) = arrRes
		o.buildParentLink(o.parent)
		o.callOnChange()
		return nil
	default:
		return invalidTypeErr("")
	}
}

func (o *Object) ArrRemove(index int) (err error) {
	switch o.val.(type) {
	case *arrayData:
		var (
			arr, arrBefore, arrAfter, arrRes arrayData
		)
		arr = *o.val.(*arrayData)
		// overflow
		if index < 0 || index >= len(arr) {
			return indexOverflowErr{index}
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
		o.buildParentLink(o.parent)
		o.callOnChange()
		return nil
	default:
		return invalidTypeErr("")
	}
}

func (o *Object) ArrShift() (err error) {
	return o.ArrRemove(0)
}

func (o *Object) ArrUnshift(value interface{}) (err error) {
	return o.ArrInsert(0, value)
}

func (o *Object) ArrForeach(do func(index int, obj *Object) error) (err error) {
	switch o.val.(type) {
	case *arrayData:
		for i, obj := range *o.val.(*arrayData) {
			if err = do(i, obj); err != nil {
				o.buildParentLink(o.parent)
				return
			}
		}
		o.buildParentLink(o.parent)
		return nil
	default:
		return invalidTypeErr("")
	}
}

// ArrPushArr Push all of elements in an ARRAY Object named o2 into the ARRAY Object
func (o *Object) ArrPushArr(o2 *Object) (err error) {
	switch o.val.(type) {
	case *arrayData:
		switch o2.val.(type) {
		case *arrayData: // Group
			newArr := o.Clone()
			err = o2.ArrForeach(func(index int, obj *Object) error {
				return newArr.ArrPush(obj)
			})
			if err == nil {
				o.SetVal(newArr)
				o.buildParentLink(o.parent)
			}
			return
		default:
			return invalidTypeErr("")
		}
	default:
		return invalidTypeErr("")
	}
}

func (o *Object) ArrPushAll(values ...interface{}) (err error) {
	switch o.val.(type) {
	case *arrayData:
		o2 := New(Array(values))
		err = o.ArrPushArr(o2)
		o.buildParentLink(o.parent)
	default:
		return invalidTypeErr("")
	}
	return
}

func (o *Object) ArrLen() (int, error) {
	switch o.val.(type) {
	case *arrayData:
		return len(*o.val.(*arrayData)), nil
	default:
		return -1, invalidTypeErr("")
	}
}
