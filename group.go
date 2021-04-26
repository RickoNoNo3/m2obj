package m2obj

func (o *Object) GroupForeach(do func(key string, obj *Object) error) (err error) {
	switch o.val.(type) {
	case *groupData:
		for k, obj := range *o.val.(*groupData) {
			if err = do(k, obj); err != nil {
				return
			}
		}
		return nil
	default:
		return invalidTypeErr("")
	}
}

// GroupMerge merges two different GROUP Object deeply.
// If forced is true, the already exist key in the group will be replaced;
func (o *Object) GroupMerge(o2 *Object, forced bool) (err error) {
	return o.groupMerge(o2, forced, true)
}

func (o *Object) groupMerge(o2 *Object, forced bool, needCallOnChange bool) (err error) {
	switch o.val.(type) {
	case *groupData: // Group
		switch o2.val.(type) {
		case *groupData: // Group
			newObj := o.Clone()
			err = o2.GroupForeach(func(key string, o2obj *Object) error {
				if newObj.Has(key) {
					o1obj := newObj.MustGet(key)
					// o1obj type check
					switch o1obj.val.(type) {
					case *groupData:
						// o2obj type check
						switch o2obj.val.(type) {
						case *groupData:
							// Merge two sub group
							return o1obj.GroupMerge(o2obj, forced)
						default:
							if forced {
								return newObj.Set(key, o2obj)
							}
						}
					default:
						if forced {
							return newObj.Set(key, o2obj)
						}
					}
				} else {
					return newObj.Set(key, o2obj)
				}
				return nil
			})
			if err == nil {
				o.setVal(newObj, needCallOnChange)
			}
			return
		default: // Array or Value
			return invalidTypeErr("")
		}
	default: // Array or Value
		return invalidTypeErr("")
	}
}
