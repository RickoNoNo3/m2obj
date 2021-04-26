package m2obj

import "reflect"

func (o *Object) IsGroup() bool {
	if o.val == nil {
		return false
	}
	switch o.val.(type) {
	case *groupData:
		return true
	default:
		return false
	}
}

func (o *Object) IsArray() bool {
	if o.val == nil {
		return false
	}
	switch o.val.(type) {
	case *arrayData:
		return true
	default:
		return false
	}
}

func (o *Object) IsValue() bool {
	if o.val == nil {
		return false
	}
	return !o.IsGroup() && !o.IsArray()
}

func (o *Object) IsNil() bool {
	return o.val == nil
}

// Is returns if the type of the val of the object is just ty.
func (o *Object) Is(ty reflect.Type) bool {
	return reflect.TypeOf(o.val) == ty
}

// IsLike returns if the val of the object has the same type with the param v.
func (o *Object) IsLike(v interface{}) bool {
	tv := getDeepestValue(v)
	return reflect.TypeOf(o.val) == reflect.TypeOf(tv)
}

func (o *Object) SetVal(value interface{}) {
	o.setVal(value, true)
}

func (o *Object) setVal(value interface{}, needCallOnChange bool) {
	o.val = getDeepestValue(value)
	if needCallOnChange {
		o.callOnChange()
	}
}

func (o *Object) Val() interface{} {
	return o.val
}

func (o *Object) ValStr() string {
	return o.val.(string)
}

func (o *Object) ValBool() bool {
	return o.val.(bool)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt8() int8 {
	return o.val.(int8)
}

func (o *Object) ValInt16() int16 {
	return o.val.(int16)
}

func (o *Object) ValInt32() int32 {
	return o.val.(int32)
}

func (o *Object) ValInt64() int64 {
	return o.val.(int64)
}

func (o *Object) ValFloat32() float32 {
	return o.val.(float32)
}

func (o *Object) ValFloat64() float64 {
	return o.val.(float64)
}
