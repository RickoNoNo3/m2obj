package m2obj

import "reflect"

// IsGroup
//
// If this Object is maintaining a(n) Group, return true, or else return false.
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

// IsArray
//
// If this Object is maintaining a(n) Array, return true, or else return false.
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

// IsValue
//
// If this Object is maintaining a(n) value neither Group nor Array, return true, or else return false.
func (o *Object) IsValue() bool {
	if o.val == nil {
		return false
	}
	return !o.IsGroup() && !o.IsArray()
}

// IsNil
//
// If this Object is not maintaining anything, return true, or else return false.
func (o *Object) IsNil() bool {
	return o.val == nil
}

// Is
//
// If this Object is maintaining anything with type ty, return true, or else return false.
func (o *Object) Is(ty reflect.Type) bool {
	return reflect.TypeOf(o.val) == ty
}

// IsLike
//
// If this Object is maintaining anything with the same type of v, return true, or else return false.
func (o *Object) IsLike(v interface{}) bool {
	tv := getDeepestValue(v)
	return reflect.TypeOf(o.val) == reflect.TypeOf(tv)
}

// SetVal
//
// Change the Object's val, no type stipulation to value, like New
func (o *Object) SetVal(value interface{}) {
	o.val = getDeepestValue(value)
	o.buildParentLink(o.Parent())
	o.callOnChange()
}

func (o *Object) setVal(value interface{}, needCallOnChange bool) {
	o.val = getDeepestValue(value)
	if needCallOnChange {
		o.callOnChange()
	}
}

// Val
//
// Get the inner value of an Object
func (o *Object) Val() interface{} {
	return o.val
}

// ValStr
//
// Get the inner value of an Object, and assert it is or transform it to a `string`.
func (o *Object) ValStr() string {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(""))
	return v.String()
}

// ValBool
//
// Get the inner value of an Object, and assert it is or transform it to a `bool`.
func (o *Object) ValBool() bool {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(true))
	return v.Bool()
}

// ValByte
//
// Get the inner value of an Object, and assert it is or transform it to a `byte`.
func (o *Object) ValByte() byte {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(byte(1)))
	return v.Interface().(byte)
}

// ValBytes
//
// Get the inner value of an Object, and assert it is or transform it to a `[]byte`.
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

// ValRune
//
// Get the inner value of an Object, and assert it is or transform it to an `rune`.
func (o *Object) ValRune() rune {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(' '))
	return v.Interface().(rune)
}

// ValRunes
//
// Get the inner value of an Object, and assert it is or transform it to an `[]rune`.
func (o *Object) ValRunes() []rune {
	v := reflect.ValueOf(o.val)
	// 针对[]byte转[]rune做特殊处理
	if v.Type() == reflect.TypeOf([]byte{}) {
		return New(o.ValStr()).ValRunes()
	}
	v = v.Convert(reflect.TypeOf([]rune{}))
	return v.Interface().([]rune)
}

// ValInt
//
// Get the inner value of an Object, and assert it is or transform it to an `int`.
func (o *Object) ValInt() int {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(1))
	return v.Interface().(int)
}

// ValInt8
//
// Get the inner value of an Object, and assert it is or transform it to an `int8`.
func (o *Object) ValInt8() int8 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int8(1)))
	return v.Interface().(int8)
}

// ValInt16
//
// Get the inner value of an Object, and assert it is or transform it to an `int16`.
func (o *Object) ValInt16() int16 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int16(1)))
	return v.Interface().(int16)
}

// ValInt32
//
// Get the inner value of an Object, and assert it is or transform it to an `int32`.
func (o *Object) ValInt32() int32 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int32(1)))
	return v.Interface().(int32)
}

// ValInt64
//
// Get the inner value of an Object, and assert it is or transform it to an `int64`.
func (o *Object) ValInt64() int64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(int64(1)))
	return v.Int()
}

// ValUint
//
// Get the inner value of an Object, and assert it is or transform it to an `uint64`.
func (o *Object) ValUint() uint64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(uint64(1)))
	return v.Uint()
}

// ValFloat32
//
// Get the inner value of an Object, and assert it is or transform it to a `float32`.
func (o *Object) ValFloat32() float32 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(float32(1)))
	return v.Interface().(float32)
}

// ValFloat64
//
// Get the inner value of an Object, and assert it is or transform it to a `float64`.
func (o *Object) ValFloat64() float64 {
	v := reflect.ValueOf(o.val)
	v = v.Convert(reflect.TypeOf(float64(1)))
	return v.Float()
}
