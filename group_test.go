package m2obj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObject_GroupForeach(t *testing.T) {
	var (
		grp1 *Object
		grp2 *Object
		grp3 map[string]string
	)
	grp1 = New(Group{
		"1": "1",
		"2": "2",
		"3": "3",
	})
	grp2 = New(Group{
		"1": "1",
		"2": "2",
		"3": 3,
	})
	grp3 = make(map[string]string)
	// no error and no panic
	assert.NotPanics(t, func() {
		assert.NoError(t, grp1.GroupForeach(func(key string, obj *Object) error {
			grp3[key] = obj.ValStr()
			return nil
		}))
	})
	assert.Equal(t, map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
	}, grp3)
	grp3 = make(map[string]string)
	// transform the panic to error
	assert.NotPanics(t, func() {
		assert.Error(t, grp2.GroupForeach(func(key string, obj *Object) error {
			switch obj.Val().(type) {
			case string:
				grp3[key] = obj.ValStr()
				return nil
			default:
				return invalidTypeErr("")
			}
		}))
	})
	grp3 = make(map[string]string)
	// throw the panic directly
	assert.Panics(t, func() {
		assert.NoError(t, grp2.GroupForeach(func(key string, obj *Object) error {
			grp3[key] = obj.ValStr()
			return nil
		}))
	})
	assert.Error(t, New(grp3).GroupForeach(func(key string, obj *Object) error {
		return nil
	}))
}

func groupMergeInit() (grp1, grp2 *Object) {
	grp1 = New(Group{
		"a": "a1",
		"b": "b1",
		"c": true,
		"grp": Group{
			"aa": "aa1",
			"cc": "cc1",
			"grp_inner": Group{
				"aaa": "aaa1",
			},
		},
		"grp2": Group{
			"dd": "dd1",
		},
		"arr": Array{1, 2, 3},
		"dynamic": Group{
			"1": 1,
		},
	})
	grp2 = New(Group{
		"a": "a2",
		"b": "b2",
		"c": false,
		"d": "d2",
		"grp": Group{
			"aa": "aa2",
			"bb": "bb2",
			"grp_inner": Group{
				"bbb": "bbb2",
			},
		},
		"grp2":    nil,
		"arr":     Array{4, 5, 6},
		"dynamic": "dynamic",
	})
	return
}

func TestObject_GroupMerge(t *testing.T) {
	var grp1, grp2 *Object
	// error type
	grp1, grp2 = groupMergeInit()
	assert.Error(t, New(grp1).GroupMerge(New(nil), true))
	assert.Error(t, New(nil).GroupMerge(New(grp2), true))
	assert.Error(t, New(nil).GroupMerge(New(nil), true))
	assert.NoError(t, New(grp1).GroupMerge(New(grp2), true))
	// forced
	grp1, grp2 = groupMergeInit()
	assert.NoError(t, grp1.GroupMerge(grp2, true))
	assert.Equal(t, New(Group{
		"a": "a2",
		"b": "b2",
		"c": false,
		"d": "d2",
		"grp": Group{
			"aa": "aa2",
			"bb": "bb2",
			"cc": "cc1",
			"grp_inner": Group{
				"aaa": "aaa1",
				"bbb": "bbb2",
			},
		},
		"grp2":    nil,
		"arr":     Array{4, 5, 6},
		"dynamic": "dynamic",
	}).Staticize(), grp1.Staticize())
	// unforced
	grp1, grp2 = groupMergeInit()
	assert.NoError(t, grp1.GroupMerge(grp2, false))
	assert.Equal(t, New(Group{
		"a": "a1",
		"b": "b1",
		"c": true,
		"d": "d2",
		"grp": Group{
			"aa": "aa1",
			"bb": "bb2",
			"cc": "cc1",
			"grp_inner": Group{
				"aaa": "aaa1",
				"bbb": "bbb2",
			},
		},
		"grp2": Group{
			"dd": "dd1",
		},
		"arr": Array{1, 2, 3},
		"dynamic": Group{
			"1": 1,
		},
	}).Staticize(), grp1.Staticize())
}
