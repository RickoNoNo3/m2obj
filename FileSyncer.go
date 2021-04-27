package m2obj

import (
	"io/ioutil"
	"sync"
	"time"
)

// FileSyncer
//
// !!! Only Bind to GROUP Object please !!!
type FileSyncer struct {
	// HardLoad
	// use SetVal(), instead of GroupMerge()
	HardLoad bool
	// AutoSaveTiming
	//  <0: Don't auto save
	//  =0: DEFAULT. Auto save when obj changed
	//  >0(ms): Auto save when timer triggered
	AutoSaveTiming int64
	// AutoLoadTiming
	//  <=0: DEFAULT. Don't auto load
	//  >0(ms): Auto load when timer triggered
	// while AutoLoadTiming > 0, the auto saving is disabled
	AutoLoadTiming int64
	// path to the file
	filePath string
	// an instance of a kind of data formatters
	formatter Formatter
	// bound object
	obj *Object

	// file mutex
	fileMutex sync.Mutex
	// object mutex
	objMutex sync.Mutex
}

func (fs *FileSyncer) GetFilePath() (filePath string) {
	fs.fileMutex.Lock()
	defer fs.fileMutex.Unlock()
	return fs.filePath
}

func (fs *FileSyncer) SetFilePath(filePath string) {
	fs.fileMutex.Lock()
	defer fs.fileMutex.Unlock()
	fs.filePath = filePath
	return
}

func (fs *FileSyncer) GetBoundObject() (obj *Object) {
	fs.objMutex.Lock()
	defer fs.objMutex.Unlock()
	return fs.obj
}

func (fs *FileSyncer) SetFormatter(formatter Formatter) {
	fs.fileMutex.Lock()
	defer fs.fileMutex.Unlock()
	fs.objMutex.Lock()
	defer fs.objMutex.Unlock()
	fs.formatter = formatter
	return
}

func (fs *FileSyncer) Save() (err error) {
	fs.objMutex.Lock()
	if fs.obj == nil {
		fs.objMutex.Unlock()
		return unknownTypeErr("")
	} else {
		fs.objMutex.Unlock()
	}
	var buf []byte
	fs.objMutex.Lock()
	buf, err = fs.formatter.Marshal(fs.obj)
	fs.objMutex.Unlock()
	if err == nil {
		fs.fileMutex.Lock()
		err = ioutil.WriteFile(fs.filePath, buf, 0644)
		fs.fileMutex.Unlock()
	}
	return
}

func (fs *FileSyncer) Load() (err error) {
	fs.objMutex.Lock()
	if fs.obj == nil {
		fs.objMutex.Unlock()
		return unknownTypeErr("")
	} else {
		fs.objMutex.Unlock()
	}
	var buf []byte
	fs.fileMutex.Lock()
	buf, err = ioutil.ReadFile(fs.filePath)
	fs.fileMutex.Unlock()
	if err == nil {
		var obj *Object
		obj, err = fs.formatter.Unmarshal(buf)
		if err == nil {
			if fs.HardLoad {
				fs.objMutex.Lock()
				fs.obj.setVal(obj, false)
				fs.objMutex.Unlock()
			} else {
				fs.objMutex.Lock()
				err = fs.obj.groupMerge(obj, true, false)
				fs.objMutex.Unlock()
			}
		}
	}
	return
}

// BindObject
//
// !!! Only Bind to GROUP Object !!!
func (fs *FileSyncer) BindObject(obj *Object) {
	if obj == nil || !obj.IsGroup() {
		panic(invalidTypeErr(""))
	}
	fs.objMutex.Lock()
	defer fs.objMutex.Unlock()
	if fs.obj != nil {
		fs.obj.onChange = nil
	}
	fs.obj = obj
	fs.obj.onChange = func() {
		if fs.AutoLoadTiming <= 0 && fs.AutoSaveTiming == 0 {
			_ = fs.Save()
		}
	}
}

func NewFileSyncer(filePath string, formatter Formatter) *FileSyncer {
	fs := &FileSyncer{
		filePath:  filePath,
		formatter: formatter,
		HardLoad:  false,
		obj:       nil,
	}
	go func() {
		for {
			func() {
				defer func() {
					_ = recover()
				}()
				autoLoadTiming := fs.AutoLoadTiming
				autoSaveTiming := fs.AutoSaveTiming
				if autoLoadTiming > 0 {
					time.Sleep(time.Duration(autoLoadTiming * int64(time.Millisecond)))
					_ = fs.Load()
				} else if autoSaveTiming > 0 {
					time.Sleep(time.Duration(autoSaveTiming * int64(time.Millisecond)))
					_ = fs.Save()
				} else {
					time.Sleep(time.Second)
				}
			}()
		}
	}()
	return fs
}
