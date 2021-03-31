package m2obj

import (
	"strings"
)

// split splits the keyStr to keys.
func split(keyStr string) (keys []string) {
	keys = make([]string, 0)
	if keyStr = strings.TrimSpace(keyStr); keyStr == "" {
		return
	}
	return strings.Split(keyStr, ".")
}

// splitAndDig digs into current object in-depth assigned by `keyStr`, until it gets the last element and returns it.
//
// Set `createLost` to true if you want to create lost keys in the `keyStr`.
// All of lost middle keys of `keyStr` will be checked:
//   If it is marked as an Array (There is an `[index]` key behind it), send panic always.
//   Else, create as a Group Object.
// The last key of `keyStr` just lost will be created as an empty Value Object. You can do something by yourself.
//
// The func panic at:
//     1. the key is not found and `createLost` is false
//     2. the key is middle of `keyStr` and has an object type neither *GroupData nor *ArrayData
//     3. the key behind an Array Object key doesn't satisfy the rule with ArrayName.[index]
func splitAndDig(current *Object, keyStr string, createLost bool) *Object {
	tObj := current
	keys := split(keyStr)
	for i, key := range keys {
		if key == "" {
			continue
		}
		// Once the code runs here, the tObj means the parent of the param key.
		// After this switch, the tObj will be the object self assigned by the param key.
		switch tObj.val.(type) {
		case *GroupData:
			if next, ok := (*tObj.val.(*GroupData))[key]; ok { // the key exists
				tObj = next
			} else if createLost { // not exists but can be created
				mapObj := *tObj.val.(*GroupData)
				if i != len(keys)-1 { // is a middle key
					if arrCheckIndexFormat(keys[i+1]) { // is Array
						panic(InvalidKeyStrErr(keyStr))
					} else { // is Group
						mapObj[key] = New(GroupData{})
					}
				} else { // is the last key
					mapObj[key] = New(nil)
				}
				tObj = mapObj[key]
			} else { // not found and panic
				panic(InvalidKeyStrErr(keyStr))
			}
		case *ArrayData:
			if index, err := tObj.arrCheckIndexKey(key, keyStr); err == nil {
				tObj = (*tObj.val.(*ArrayData))[index]
			} else {
				panic(err)
			}
		default:
			panic(InvalidTypeErr(key))
		}
	}
	return tObj
}

func getDeepestValue(v interface{}) interface{} {
	tv := v
	for {
		switch tv.(type) {
		case Object:
			tv = tv.(Object).val
		case *Object:
			tv = tv.(*Object).val
		case GroupData:
			group := tv.(GroupData)
			return &group
		case *GroupData:
			return tv.(*GroupData)
		case ArrayData:
			arr := tv.(ArrayData)
			return &arr
		case *ArrayData:
			return tv.(*ArrayData)
		default:
			return tv
		}
	}
}
