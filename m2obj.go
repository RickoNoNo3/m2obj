package m2obj

import (
	"regexp"
	"strconv"
)

// Err Definition

type IndexOverflowErr struct {
	Index int
}

func (e IndexOverflowErr) Error() string {
	return "no such index[" + strconv.Itoa(e.Index) + "]"
}

type InvalidKeyStrErr string

func (e InvalidKeyStrErr) Error() string {
	return "invalid key string: " + string(e)
}

type UnknownTypeErr string

func (e UnknownTypeErr) Error() string {
	return "the key {" + string(e) + "} has an unknown ObjectType"
}

type InvalidTypeErr string

func (e InvalidTypeErr) Error() string {
	return "the key {" + string(e) + "} has an invalid ObjectType"
}

// Type Definition
type Object struct {
	val interface{}
}

type GroupData map[string]*Object
type ArrayData []*Object

type DataFormatter interface {
	Marshal(obj *Object) (objStr string, err error)
	UnMarshal(objStr string) (obj *Object, err error)
	SaveToFile(obj *Object) (err error)
	LoadFromFile() (obj *Object, err error)
}

// Method Definition

func (o *Object) Set(keyStr string, value interface{}) (err error) {
	defer func() { // recover any panic to error and return the error
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj := splitAndDig(o, keyStr, true)
	obj.val = getDeepestValue(value)
	return
}

func (o *Object) SetIfHas(keyStr string, value interface{}) (err error) {
	if o.Has(keyStr) {
		return o.Set(keyStr, value)
	}
	return nil
}

func (o *Object) SetIfNotHas(keyStr string, value interface{}) (err error) {
	if !o.Has(keyStr) {
		return o.Set(keyStr, value)
	}
	return nil
}

func (o *Object) Get(keyStr string) (obj *Object, err error) {
	defer func() {
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj = splitAndDig(o, keyStr, false)
	return
}

func (o *Object) MustGet(keyStr string) (obj *Object) {
	var err error
	if obj, err = o.Get(keyStr); err != nil {
		panic(err)
	}
	return
}

func (o *Object) Has(keyStr string) bool {
	_, err := o.Get(keyStr)
	return err == nil
}

func (o *Object) Remove(keyStr string) bool {
	if keyStr == "" {
		return false
	}
	if o.Has(keyStr) {
		var (
			key       string
			parentObj *Object
			reg       = regexp.MustCompile("^(.+)[.]([^.]+)$")
		)
		if reg.MatchString(keyStr) {
			submatch := reg.FindStringSubmatch(keyStr)
			key = submatch[2]
			parentKeyStr := submatch[1]
			parentObj = splitAndDig(o, parentKeyStr, false)
		} else {
			key = keyStr
			parentObj = o
		}
		switch parentObj.val.(type) {
		case *GroupData:
			delete(*parentObj.val.(*GroupData), key)
			return true
		default:
			return false
		}
	} else {
		return true // Not exists, regarded as remove successfully.
	}
}

func (o *Object) SetVal(value interface{}) {
	o.val = getDeepestValue(value)
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

// !!! DANGEROUS
//
// Returns a pointer to the Object's core array val to achieve more advanced operations on it.
//
// This func is the only one which has writable access to the val of an Object. So be careful.
func (o *Object) ValArr() *ArrayData {
	return o.val.(*ArrayData)
}

// staticize without the wrapper, for different object type, it returns different type:
//     Group: map[string]interface{}
//     Array: []interface{}
//     Value: interface{}
func (o *Object) staticize() interface{} {
	switch o.val.(type) {
	case *GroupData: // Group
		m := make(map[string]interface{})
		for k, v := range *o.val.(*GroupData) {
			if v == nil {
				m[k] = nil
			} else {
				m[k] = v.staticize()
			}
		}
		return m
	case *ArrayData: // Array
		m := make([]interface{}, len(*o.val.(*ArrayData)))
		for i, v := range *o.val.(*ArrayData) {
			if v == nil {
				m[i] = nil
			} else {
				m[i] = v.staticize()
			}
		}
		return m
	default: // Value
		if o == nil {
			return nil
		} else {
			return o.val
		}
	}
}

func (o *Object) Staticize() map[string]interface{} {
	switch o.val.(type) {
	case *GroupData: // Group
		return o.staticize().(map[string]interface{})
	case *ArrayData: // Array
		return map[string]interface{}{
			"list": o.staticize().([]interface{}),
		}
	default: // Value
		return map[string]interface{}{
			"val": o.staticize(),
		}
	}
}

func (o *Object) Clone() (newObj *Object) {
	switch o.val.(type) {
	case *GroupData: // Group
		newObj = New(GroupData{})
		for k, obj := range *o.val.(*GroupData) {
			_ = newObj.Set(k, obj.Clone())
		}
		return
	case *ArrayData: // Array
		newObj = New(ArrayData{})
		for _, obj := range *o.val.(*ArrayData) {
			_ = newObj.ArrPush(obj.val)
		}
		return
	default: // Value
		newObj = New(o.val)
		return
	}
}

// TODO: 在新建Group和Array时减少New出现的频率, 把*Object类型限定放宽为interface{}, 然后由代码自己处理成*Object
func New(value interface{}) *Object {
	t := getDeepestValue(value)
	return &Object{
		val: t,
	}
}

func NewFromMap(m map[string]interface{}) *Object {
	obj := New(GroupData{})
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}:
			_ = obj.Set(k, NewFromMap(v.(map[string]interface{})))
		case []interface{}:
			arr := New(ArrayData{})
			for _, v2 := range v.([]interface{}) {
				_ = arr.ArrPush(v2)
			}
			_ = obj.Set(k, arr)
		default:
			_ = obj.Set(k, v)
		}
	}
	return obj
}
