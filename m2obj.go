package m2obj

import (
	"regexp"
	"strconv"
)

// Err Definition

type indexOverflowErr struct {
	Index int
}

func (e indexOverflowErr) Error() string {
	return "no such index[" + strconv.Itoa(e.Index) + "]"
}

type invalidKeyStrErr string

func (e invalidKeyStrErr) Error() string {
	return "invalid key string: " + string(e)
}

type unknownTypeErr string

func (e unknownTypeErr) Error() string {
	if string(e) == "" {
		return "unknown ObjectType"
	} else {
		return "the key {" + string(e) + "} has an unknown ObjectType"
	}
}

type invalidTypeErr string

func (e invalidTypeErr) Error() string {
	if string(e) == "" {
		return "invalid ObjectType"
	} else {
		return "the key {" + string(e) + "} has an invalid ObjectType"
	}
}

// Type Definition

type Object struct {
	val      interface{}
	parent   *Object
	onChange func() // used by fileSyncer
}

type Group map[string]interface{}
type Array []interface{}
type groupData map[string]*Object
type arrayData []*Object

// Method Definition

func (o *Object) callOnChange() {
	tObj := o
	for tObj != nil {
		if tObj.onChange != nil {
			tObj.onChange()
		}
		tObj = tObj.parent
	}
}

func (o *Object) Set(keyStr string, value interface{}) (err error) {
	defer func() { // recover any panic to error and return the error
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj := splitAndDig(o, keyStr, true, true)
	obj.SetVal(value)
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
	obj = splitAndDig(o, keyStr, false, true)
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
			parentObj = splitAndDig(o, parentKeyStr, false, true)
		} else {
			key = keyStr
			parentObj = o
		}
		switch parentObj.val.(type) {
		case *groupData:
			delete(*parentObj.val.(*groupData), key)
			o.callOnChange()
			return true
		default:
			return false
		}
	} else {
		return true // Not exists, regarded as remove successfully.
	}
}

// staticize without the wrapper, for different object type, it returns different type:
//     group: map[string]interface{}
//     array: []interface{}
//     value: interface{}
func (o *Object) staticize() interface{} {
	switch o.val.(type) {
	case *groupData: // Group
		m := make(map[string]interface{})
		for k, v := range *o.val.(*groupData) {
			if v == nil || v.Val() == nil {
				m[k] = nil
			} else {
				m[k] = v.staticize()
			}
		}
		return m
	case *arrayData: // Array
		m := make([]interface{}, len(*o.val.(*arrayData)))
		for i, v := range *o.val.(*arrayData) {
			if v == nil || v.Val() == nil {
				m[i] = nil
			} else {
				m[i] = v.staticize()
			}
		}
		return m
	default: // Value
		if o == nil || o.Val() == nil {
			return nil
		} else {
			return o.val
		}
	}
}

func (o *Object) Staticize() map[string]interface{} {
	switch o.val.(type) {
	case *groupData: // Group
		return o.staticize().(map[string]interface{})
	case *arrayData: // Array
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
	case *groupData: // Group
		newObj = newWithParent(groupData{}, o.parent)
		for k, obj := range *o.val.(*groupData) {
			_ = newObj.Set(k, obj.Clone())
		}
		return
	case *arrayData: // Array
		newObj = newWithParent(arrayData{}, o.parent)
		for _, obj := range *o.val.(*arrayData) {
			_ = newObj.ArrPush(obj.Clone())
		}
		return
	default: // Value
		newObj = newWithParent(o.val, o.parent)
		return
	}
}

func New(value interface{}) *Object {
	return newWithParent(value, nil)
}

func newWithParent(value interface{}, parent *Object) *Object {
	t := getDeepestValue(value)
	return &Object{
		val:    t,
		parent: parent,
	}
}
