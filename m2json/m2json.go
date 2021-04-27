package m2json

import (
	"bytes"
	"encoding/json"
	"github.com/rickonono3/m2obj"
)

type Formatter struct {
}

func (f Formatter) Marshal(obj *m2obj.Object) (data []byte, err error) {
	buf := bytes.Buffer{}
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(obj.Staticize())
	if err == nil {
		data = buf.Bytes()
	}
	return
}

func (f Formatter) Unmarshal(data []byte) (obj *m2obj.Object, err error) {
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	m := make(map[string]interface{})
	err = decoder.Decode(&m)
	if err == nil {
		obj = m2obj.New(m)
	}
	return
}
