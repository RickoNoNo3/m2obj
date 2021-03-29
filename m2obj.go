package m2obj

import (
	"regexp"
	"strconv"
	"strings"
)

type IndexOverflowErr struct {
	Index int
}

func (e IndexOverflowErr) Error() string {
	return "no such index[" + strconv.Itoa(e.Index) + "]"
}

type InvalidKeyStrErr string

func (e InvalidKeyStrErr) Error() string { return "invalid key string: " + string(e) }

type UnknownTypeErr string

func (e UnknownTypeErr) Error() string { return "the key {" + string(e) + "} has an unknown ObjectType" }

type ObjectType int

const Group = ObjectType(0)
const Array = ObjectType(1)
const Value = ObjectType(2)

type DataFormatter interface {
	Marshal(*Object) (string, error)
	UnMarshal(string) (*Object, error)
}

type ObjectTypeAssert interface {
	Assert(*Object) bool
}

type Object struct {
	val  interface{}
	Type ObjectType
}

type mapVal *map[string]*Object
type arrVal []*Object

type splitterDo func(key string)

// split the keyStr to slice keys, foreach key in keys, call the do(key).
func splitter(keyStr string, do splitterDo) {
	if keyStr = strings.TrimSpace(keyStr); keyStr == "" {
		return
	}
	keys := strings.Split(keyStr, ".")
	for _, key := range keys {
		do(key)
	}
}

func splitterDoDrag(current **Object, keyStr string, createLostGroup bool) splitterDo {
	return func(key string) {
		switch (*current).Type {
		// Group, goto next directly if exists, or create and goto a new Object.
		case Group:
			next, err := (*current).Get(key)
			// check Has
			if err != nil {
				// not has: create first and goto.
				if createLostGroup {
					mapObj := *(*current).val.(mapVal)
					mapObj[key] = NewEmptyGroup()
					*current = mapObj[key]
				} else {
					panic(InvalidKeyStrErr(keyStr))
				}
			} else {
				// has: directly goto.
				*current = next
			}
		// Array, check if next is an index flag and goto the element correspondingly.
		case Array:
			if index, err := (*current).arrCheckKeyStringNext(key, keyStr); err != nil {
				panic(err)
			} else {
				*current = (*current).val.(arrVal)[index]
			}
		case Value:
			panic(InvalidKeyStrErr(keyStr))
		default:
			panic(UnknownTypeErr(key))
		}
	}
}

func (o *Object) arrCheckKeyStringNext(key, keyStr string) (index int, err error) {
	reg := regexp.MustCompile(`\[(\d+)]`)
	// check [n] format
	if !reg.MatchString(key) {
		err = InvalidKeyStrErr(keyStr)
		return
	} else {
		index, err := strconv.Atoi(reg.FindStringSubmatch(key)[1])
		// check index atoi
		if err != nil {
			err = InvalidKeyStrErr(keyStr)
			return
		} else {
			arr := o.val.(arrVal)
			// check index overflow
			if len(arr) <= index {
				err = IndexOverflowErr{
					Index: index,
				}
				return
			} else {
				return index, nil
			}
		}
	}
}

func (o *Object) ArrPush(value interface{}) (err error) {

}

func (o *Object) ArrPop(value interface{}) (err error) {

}

func (o *Object) ArrSet(index int, value interface{}) (err error) {

}

func (o *Object) ArrGet(index int) (value interface{}, err error) {

}

func (o *Object) ArrInsert(index int, value interface{}) (err error) {

}

func (o *Object) ArrDelete(index int, value interface{}) (err error) {

}

func (o *Object) ArrForeach(do func(index int, obj *Object)) (err error) {

}

func (o *Object) Set(keyStr string, value interface{}) (err error) {
	defer func() {
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	current := o
	splitter(keyStr, splitterDoDrag(&current, keyStr, true))
	tObj := NewWithData(value)
	o.val = tObj.val
	o.Type = tObj.Type
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

// TODO: SetIfHasAndIs

func (o *Object) Get(keyStr string) (obj *Object, err error) {
	defer func() {
		if pan := recover(); pan != nil {
			err = pan.(error)
		}
	}()
	current := o
	splitter(keyStr, splitterDoDrag(&current, keyStr, false))
	return current, nil
}

// TODO: GetIfIs

func (o *Object) MustGet(keyStr string) (obj *Object) {
	var err error
	if obj, err = o.Get(keyStr); err != nil {
		panic(InvalidKeyStrErr(keyStr))
	}
	return
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
// This function is the only one which has writable access to the val of an Object. So be careful.
func (o *Object) ValArr() *[]*Object {
	return (interface{}(&o.val)).(*[]*Object)
}

func (o *Object) ToMap() map[string]interface{} {
	// TODO: Staticize
	return
}

func (o *Object) ToString(formatter DataFormatter) (str string, err error) {
	str, err = formatter.Marshal(o)
	return
}

// TODO: 进行到这里
func (o *Object) ToFile(filePath string, formatter DataFormatter) (str string, err error) {
	str, err = formatter.Marshal(o)
	return
}

func (o *Object) ToString(formatter DataFormatter) (str string, err error) {
	str, err = formatter.Marshal(o)
	return
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

func (o *Object) ValInt() int {
	return o.val.(int)
}

// func (obj *Object) Has(key string) bool {
// 	if obj.Type != Group {
// 		panic(panicCallOnlyGroup)
// 	}
// 	return (*obj.List)[key] != nil
// }
//
// func (obj *Object) Del(key string) {
// 	if obj.Type != Group {
// 		panic(panicCallOnlyGroup)
// 	}
// 	if !obj.Has(key) {
// 		panic(panicNoKey + key)
// 	}
// 	delete(*obj.List, key)
// }
//
// func (obj *Object) DelIfExists(key string) {
// 	if obj.Type != Group {
// 		panic(panicCallOnlyGroup)
// 	}
// 	if !obj.Has(key) {
// 		panic(panicNoKey + key)
// 	}
// 	delete(*obj.List, key)
// }
//
// // 获取一个Object中一个key对应的子Object, 未找到key就会panic
// func (obj *Object) Get(key string) (val *Object) {
// 	if obj.Type != Group {
// 		panic(panicCallOnlyGroup)
// 	}
// 	if !obj.Has(key) {
// 		panic(panicNoKey + key)
// 	}
// 	return (*obj.List)[key]
// }
//
// // 给一个Object中的一个key设置对应的子Object, 可同时用于创建和更新
// func (obj *Object) SetObj(key string, val *Object) {
// 	if obj.Type != Group {
// 		panic(panicCallOnlyGroup)
// 	}
// 	(*obj.List)[key] = val
// }
//
// func (obj *Object) Set(key string, val interface{}) {
// 	var tObj *Object
// 	switch val.(type) {
// 	case Object:
// 		valObj := val.(Object)
// 		tObj = valObj.Clone()
// 	case *Object:
// 		tObj = val.(*Object).Clone()
// 	default:
// 		tObj = NewValue(obj)
// 	}
// 	obj.SetObj(key, NewValue(val))
// }
//
// func (obj *Object) SetIfNotExists(key string, val interface{}) {
// 	if !obj.Has(key) {
// 		obj.Set(key, val)
// 	}
// }
//
// func (obj *Object) SetIfNotExistsObj(key string, val interface{}) {
// 	if !obj.Has(key) {
// 		obj.Set(key, val)
// 	}
// }
//
// func (obj *Object) SetIfExists(key string, val *Object) {
// 	if obj.Has(key) {
// 		obj.Set(key, val)
// 	}
// }
//
// // 静态化. 将Object转为map[string]interface{}结构, 便于在其他地方使用(如构造json/模板传参).
// // 注意有子元素的会看作Group节点, 否则才看作Value节点, 因此有子元素的元素的Val不会被导出.
// func (obj *Object) Staticize() (res map[string]interface{}) {
// 	res = make(map[string]interface{})
// 	for k, v := range *obj.List {
// 		if len(*v.List) == 0 { // 不是Group
// 			res[k] = v.Val
// 		} else { // 是Group
// 			res[k] = v.Staticize()
// 		}
// 	}
// 	return
// }
//
// func (obj *Object) Clone() (newObj *Object) {
// 	if obj.Type == Group {
// 		// TODO: CloneGroup
// 		return
// 	} else {
// 		return NewValue(obj.Val)
// 	}
// }
//
// // Create a new Group Object and let Object.val = val
// //
// // is single-layer (no more map[string]interface{} as interface{} in the values). To create a new Group Object with more than one layer, use another inner NewGroup. See examples:
// func NewGroup(children map[string]interface{}) *Object {
// 	return &Object{
// 		Group,
// 		nil,
// 		&children,
// 	}
// }
//
// // Create a new Value Object and let Object.Val = val.
// //
// //   If val is an Object/*Object and is a Value, this function equals to *Object.Clone.
// //   It ensures that the Object.Val is always the real data directly instead of a nested Object.
// //
// //   If val is an Object/*Object and is a Group, it makes a panic.
// //   NEVER put an Object to the Object.Val!
// //
// //   If val is one of other types, create a new Value Object simply.
// //
// func NewValue(val interface{}) (newObj *Object) {
// 	switch val.(type) {
// 	case Object:
// 		valObj := val.(Object)
// 		if valObj.Type != Value {
// 			panic(panicArgOnlyValue)
// 		}
// 		newObj = valObj.Clone()
// 	case *Object:
// 		valObj := val.(*Object)
// 		if valObj.Type != Value {
// 			panic(panicArgOnlyValue)
// 		}
// 		newObj = valObj.Clone()
// 	default:
// 		newObj = &Object{
// 			Value,
// 			val,
// 			nil,
// 		}
// 	}
// 	return
// }
