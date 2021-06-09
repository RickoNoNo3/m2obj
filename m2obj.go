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

// callOnChange
//
// Bubble the onChange event until the root element is reached
func (o *Object) callOnChange() {
	tObj := o
	for tObj != nil {
		if tObj.onChange != nil {
			tObj.onChange()
		}
		tObj = tObj.parent
	}
}

// Parent
//
// Returns the parent of the current object in the object tree, and returns nil when it is the root element.
func (o *Object) Parent() *Object {
	return o.parent
}

// buildParentLink
//
// Build or Rebuild the parent-child relationship for object tree
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

// Set
//
// Set the value for a elements located by the keyStr. The elements in the keyStr will be created automatically.
//
// An invalidKeyStrErr will be reported when the keyStr doesn't meet the agreed format.
//
// An invalidTypeErr will be reported when the actual object type does not match the tag format in the keyStr.
//
//
// Example:
//
//   obj := m2obj.New(m2obj.Group{})
//   obj.Set("a.b.c", 0)              // {a:{b:{c:0}}}
//   obj.Set("a.b.c.d", 2)            // invalidType (a.b.c isn't a group)
//   obj.Set("a.b.c", m2obj.Array{0}) // {a:{b:{c:[0]}}}
//   obj.Set("a.b.c.[0]", 1)          // {a:{b:{c:[1]}}}
//   obj.Set("a.b.c.[1]", 1)          // indexOverflow
//   obj.Set("a.b.c.d", 1)            // invalidType
//
// For more info about keyStr, see https://github.com/rickonono3/m2obj
func (o *Object) Set(keyStr string, value interface{}) (err error) {
	defer func() { // recover any panic to error and return the error
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj, _ := splitAndDig(o, keyStr, true)
	obj.setVal(value, true)
	o.buildParentLink(o.Parent())
	obj.callOnChange()
	return
}

// SetIfHas
//
// See Set and Has
func (o *Object) SetIfHas(keyStr string, value interface{}) (err error) {
	if o.Has(keyStr) {
		return o.Set(keyStr, value)
	}
	return nil
}

// SetIfNotHas
//
// See Set and Has
func (o *Object) SetIfNotHas(keyStr string, value interface{}) (err error) {
	if !o.Has(keyStr) {
		return o.Set(keyStr, value)
	}
	return nil
}

// Get
//
// Get an object located by the keyStr. Any non-existing or unexpected tag in the keyStr will cause error.
//
// For more info about keyStr, see https://github.com/rickonono3/m2obj
func (o *Object) Get(keyStr string) (obj *Object, err error) {
	defer func() {
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	obj, _ = splitAndDig(o, keyStr, false)
	return
}

// MustGet
//
// Like Get, but panic when error occurred
func (o *Object) MustGet(keyStr string) (obj *Object) {
	var err error
	if obj, err = o.Get(keyStr); err != nil {
		panic(err)
	}
	return
}

// Has
//
// Check if the element located by the keyStr exists.
// If the keyStr is invalid, this func will panic.
func (o *Object) Has(keyStr string) bool {
	_, err := o.Get(keyStr)
	return err == nil
}

// Remove
//
// Remove the element located by the keyStr.
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
			parentObj, _ = splitAndDig(o, parentKeyStr, false)
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

// staticize
//
// staticize without the wrapper and prepare to recursive.
//
// For different object type, it returns differently:
//     Group: map[string]interface{}
//     Array: []interface{}
//     Value: interface{}
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

// Staticize
//
// Converts an object to map[string]interface{} recursively.
//
// For different child-object type, converts it differently:
//     Group: map[string]interface{}
//     Array: []interface{}
//     Value: interface{}
//
// For different root-object type, returns a map[string]interface{} that is packaged in the outermost layer:
//     Group: <result of itself>
//     Array: {
//              "list": <result of itself>
//            }
//     Value: {
//              "val": <result of itself>
//            }
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

// Clone
//
// Deeply clone for an object.
//
// Note that if you maintain pointer elements yourself in some values, these elements cannot be deep copied.
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
			newObj.ArrPush(obj.Clone())
		}
	default: // Value
		newObj = New(o.val)
	}
	newObj.buildParentLink(nil)
	return
}

// New
//
// Create a new m2obj.Object with value.
//
// The value can be a leaked value or another Object or m2obj.Group{} or m2obj.Array{}. All of this type arguments will be automatically recognized, parsed and saved as the deepest value without any outer shell.
//
// BTW, in m2obj project, all arguments with the type interface{} follow the above principles. That is to say, you can pass various things directly to the arguments without worrying about parsing issues. They can be normal values or Objects that encapsulates a normal value.
func New(value interface{}) *Object {
	t := getDeepestValue(value)
	obj := &Object{
		val:    t,
		parent: nil,
	}
	obj.buildParentLink(nil)
	return obj
}
