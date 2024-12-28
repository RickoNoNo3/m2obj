package filesyncertest

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/rickonono3/m2obj"
	"github.com/rickonono3/m2obj/m2json"
	"github.com/stretchr/testify/assert"
)

// 需要注意：处理JSON数据时应当将所有数字视为float64

func TestFileSyncer_m2json_Save(t *testing.T) {
	initTestData("json")
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	assert.Error(t, fs.Save())
	fs.BindObject(obj)
	assert.NoError(t, fs.Save())
	fileBytes, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)
	fileObj, err := formatter.Unmarshal(fileBytes)
	assert.NoError(t, err)
	assert.Equal(t, allNumbersToFloat64(allNumbersToFloat64(obj.Staticize())), allNumbersToFloat64(allNumbersToFloat64(fileObj.Staticize())))
}

func TestFileSyncer_m2json_Load(t *testing.T) {
	initTestData("json")
	t.Run("TestFileSyncer_m2json_Save", TestFileSyncer_m2json_Save)
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	assert.Error(t, fs.Load())
	fs.BindObject(m2obj.New(m2obj.Group{}))
	assert.NoError(t, fs.Load())
	assert.Equal(t, allNumbersToFloat64(obj.Staticize()), allNumbersToFloat64(fs.GetBoundObject().Staticize()))
}

func TestFileSyncer_m2json_AutoSave(t *testing.T) {
	waitStartTime := 2 * time.Second
	waitStopTime := 1 * time.Second
	initTestData("json")
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	checkObj := func(t *testing.T) {
		fileBytes, err := ioutil.ReadFile(filePath)
		assert.NoError(t, err)
		fileObj, err := formatter.Unmarshal(fileBytes)
		assert.NoError(t, err)
		assert.Equal(t, allNumbersToFloat64(obj.Staticize()), allNumbersToFloat64(fileObj.Staticize()))
	}
	// init save
	fs.AutoSaveTiming = 500
	time.Sleep(waitStartTime)
	fs.BindObject(obj)
	time.Sleep(waitStartTime)
	fs.AutoSaveTiming = -1
	time.Sleep(waitStopTime)
	t.Run("init save", checkObj)
	// first save
	obj.Remove("e.ea")
	fs.AutoSaveTiming = 500
	time.Sleep(waitStartTime)
	fs.AutoSaveTiming = -1
	time.Sleep(waitStopTime)
	t.Run("first save", checkObj)
	// second save
	obj.SetVal(m2obj.Group{
		"secondView": true,
	})
	fs.AutoSaveTiming = 500
	time.Sleep(waitStartTime)
	fs.AutoSaveTiming = -1
	time.Sleep(waitStopTime)
	t.Run("second save", checkObj)
}

func TestFileSyncer_m2json_AutoLoad(t *testing.T) {
	waitStartTime := 2 * time.Second
	waitStopTime := 1 * time.Second
	initTestData("json")
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	fs.AutoSaveTiming = -1
	obj2 := obj.Clone()
	checkObjEqual := func(t *testing.T) {
		assert.Equal(t, allNumbersToFloat64(obj.Staticize()), allNumbersToFloat64(obj2.Staticize()))
	}
	checkObjNotEqual := func(t *testing.T) {
		assert.NotEqual(t, allNumbersToFloat64(obj.Staticize()), allNumbersToFloat64(obj2.Staticize()))
	}
	// first load
	fs.BindObject(obj)
	assert.NoError(t, fs.Save())
	obj.Remove("a")
	obj.Remove("b")
	obj.Remove("c")
	obj.Remove("d")
	obj.Remove("e")
	t.Run("first load before", checkObjNotEqual)
	fs.AutoLoadTiming = 500
	time.Sleep(waitStartTime)
	fs.AutoLoadTiming = -1
	time.Sleep(waitStopTime)
	t.Run("first load after", checkObjEqual)
	// second load
	obj.SetVal(m2obj.Group{
		"secondView": true,
	})
	assert.NoError(t, fs.Save())
	obj.SetVal(obj2.Clone())
	t.Run("second load before", checkObjEqual)
	fs.AutoLoadTiming = 500
	time.Sleep(waitStartTime)
	t.Run("second load after", checkObjNotEqual)
	t.Run("second load after content check", func(t *testing.T) {
		expect := obj2.Clone()
		assert.NoError(t, expect.GroupMerge(m2obj.New(m2obj.Group{
			"secondView": true,
		}), true))
		assert.Equal(t, allNumbersToFloat64(expect.Staticize()), allNumbersToFloat64(obj.Staticize()))
	})
	// HardLoad
	fs.HardLoad = true
	time.Sleep(waitStartTime)
	t.Run("second load after content check with HardLoad", func(t *testing.T) {
		expect := m2obj.New(m2obj.Group{
			"secondView": true,
		})
		assert.Equal(t, allNumbersToFloat64(expect.Staticize()), allNumbersToFloat64(obj.Staticize()))
	})
}

func TestFileSyncer_m2json_BindObject(t *testing.T) {
	initTestData("json")
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	obj2 := m2obj.New(m2obj.Group{
		"test": "test",
	})
	obj3 := m2obj.New("hello world!")
	// bind to obj
	assert.NotPanics(t, func() {
		fs.BindObject(obj)
	})
	assert.Equal(t, obj, fs.GetBoundObject())
	// bind to obj2
	assert.NotPanics(t, func() {
		fs.BindObject(obj2)
	})
	assert.Equal(t, obj2, fs.GetBoundObject())
	// bind to obj3
	assert.Panics(t, func() {
		fs.BindObject(obj3)
	})
	assert.Equal(t, obj2, fs.GetBoundObject())
}

func TestFileSyncer_m2json_GetSet(t *testing.T) {
	initTestData("json")
	fs := m2obj.NewFileSyncer(filePath, m2json.Formatter{})
	fs.BindObject(obj)
	assert.NotPanics(t, func() {
		assert.Equal(t, filePath, fs.GetFilePath())
		fs.SetFilePath(filePath + ".bak")
		assert.Equal(t, filePath+".bak", fs.GetFilePath())
		fs.SetFormatter(m2json.Formatter{})
	})
}

func TestFileSyncer_m2json_AutoLoadOnChange(t *testing.T) {
	initTestData("json")
	formatter := m2json.Formatter{}
	fs := m2obj.NewFileSyncer(filePath, formatter)
	getActualFileObj := func() (fileObj *m2obj.Object) {
		fileBytes, err := ioutil.ReadFile(filePath)
		assert.NoError(t, err)
		fileObj, err = formatter.Unmarshal(fileBytes)
		assert.NoError(t, err)
		return
	}
	checkObjsEqual := func() {
		// ignore number type difference - all to allNumbersToFloat64(float64
		assert.Equal(
			t,
			allNumbersToFloat64(obj.Staticize()),
			allNumbersToFloat64(getActualFileObj().Staticize()),
		)
	}

	fs.BindObject(obj)
	assert.NoError(t, fs.Save())
	checkObjsEqual()
	obj.SetVal(m2obj.Group{
		"test": "test",
	})
	checkObjsEqual()
	assert.NoError(t, obj.Set("ha.hi.hu.he.ho.arr", m2obj.Array{
		float64(1), float64(2), float64(3),
	}))
	checkObjsEqual()
	assert.NoError(t, obj.Set("t", "233"))
	checkObjsEqual()
	assert.True(t, obj.Remove("t"))
	assert.True(t, obj.Remove("t"))
	checkObjsEqual()
}
