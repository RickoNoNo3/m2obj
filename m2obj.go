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

func (o *Object) Parent() *Object {
	return o.parent
}

// 这里的调用可以优化，只有以下情况需要重新build
//   初始化一个obj（初始化）
//   从另一个obj中接入val，另一个obj有自身的parent关系（换子）
//   对任意一个obj进行了直接赋值，含[key]和[index]（换父）
func (o *Object) buildParentLink(parent *Object) {
	o.parent = parent
	switch o.val.(type) {
	case *groupData:
		grp := *o.val.(*groupData)
		for k := range grp {
			if grp[k] == nil {
				continue
			}
			grp[k].parent = o
			if grp[k].IsGroup() || grp[k].IsArray() {
				grp[k].buildParentLink(o)
			}
		}
	case *arrayData:
		arr := *o.val.(*arrayData)
		for i := range arr {
			if arr[i] == nil {
				continue
			}
			arr[i].parent = o
			if arr[i].IsGroup() || arr[i].IsArray() {
				arr[i].buildParentLink(o)
			}
		}
	}
}

func (o *Object) Set(keyStr string, value interface{}) (err error) {
	defer func() { // recover any panic to error and return the error
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj := splitAndDig(o, keyStr, true)
	obj.SetVal(value)
	o.buildParentLink(o.parent)
	o.callOnChange()
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
		case *groupData:
			delete(*parentObj.val.(*groupData), key)
			o.buildParentLink(o.parent)
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
		newObj = New(groupData{})
		for k, obj := range *o.val.(*groupData) {
			_ = newObj.Set(k, obj.Clone())
		}
	case *arrayData: // Array
		newObj = New(arrayData{})
		for _, obj := range *o.val.(*arrayData) {
			_ = newObj.ArrPush(obj.Clone())
		}
	default: // Value
		newObj = New(o.val)
	}
	newObj.buildParentLink(nil)
	return
}

func New(value interface{}) *Object {
	t := getDeepestValue(value)
	obj := &Object{
		val:    t,
		parent: nil,
	}
	obj.buildParentLink(nil)
	return obj
}
