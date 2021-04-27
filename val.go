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
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(""))
	return v.String()
}

func (o *Object) ValBool() bool {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(true))
	return v.Bool()
}

func (o *Object) ValByte() byte {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(byte(1)))
	return v.Interface().(byte)
}

func (o *Object) ValUint() uint64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(uint64(1)))
	return v.Uint()
}

func (o *Object) ValBytes() []byte {
	v := reflect.ValueOf(o.val)
	// 针对rune转[]byte做特殊处理
	if v.Type() == reflect.TypeOf(' ') {
		return New(o.ValStr()).ValBytes()
	}
	// 针对[]rune转[]byte做特殊处理
	if v.Type() == reflect.TypeOf([]rune{}) {
		return New(o.ValStr()).ValBytes()
	}
	v = v.Convert(reflect.TypeOf([]byte{}))
	return v.Bytes()
}

func (o *Object) ValRune() rune {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(' '))
	return v.Interface().(rune)
}

func (o *Object) ValRunes() []rune {
	v := reflect.ValueOf(o.val)
	// 针对[]byte转[]rune做特殊处理
	if v.Type() == reflect.TypeOf([]byte{}) {
		return New(o.ValStr()).ValRunes()
	}
	v = v.Convert(reflect.TypeOf([]rune{}))
	return v.Interface().([]rune)
}

func (o *Object) ValInt() int {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(1))
	return v.Interface().(int)
}

func (o *Object) ValInt8() int8 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int8(1)))
	return v.Interface().(int8)
}

func (o *Object) ValInt16() int16 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int16(1)))
	return v.Interface().(int16)
}

func (o *Object) ValInt32() int32 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int32(1)))
	return v.Interface().(int32)
}

func (o *Object) ValInt64() int64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int64(1)))
	return v.Int()
}

func (o *Object) ValFloat32() float32 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(float32(1)))
	return v.Interface().(float32)
}

func (o *Object) ValFloat64() float64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(float64(1)))
	return v.Float()
}
