package m2json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"m2obj"
)

type JsonDataFormatter struct {
	FilePath string
}

func (f JsonDataFormatter) Marshal(obj *m2obj.Object) (objStr string, err error) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(obj.Staticize())
	if err == nil {
		objStr = buf.String()
	}
	return
}

func (f JsonDataFormatter) UnMarshal(objStr string) (obj *m2obj.Object, err error) {
	decoder := json.NewDecoder(bytes.NewBufferString(objStr))
	m := make(map[string]interface{})
	err = decoder.Decode(&m)
	if err == nil {
		obj = m2obj.NewFromMap(m)
	}
	return
}

func (f JsonDataFormatter) SaveToFile(obj *m2obj.Object, filePath string) (err error) {
	var str string
	if str, err = f.Marshal(obj); err == nil {
		err = ioutil.WriteFile(filePath, []byte(str), 0644)
	}
	return
}

func (f JsonDataFormatter) LoadFromFile(filePath string) (obj *m2obj.Object, err error) {
	var buf []byte
	buf, err = ioutil.ReadFile(filePath)
	if err == nil {
		obj, err = f.UnMarshal(string(buf))
	}
	return
}

func New() JsonDataFormatter {
	return JsonDataFormatter{}
}
