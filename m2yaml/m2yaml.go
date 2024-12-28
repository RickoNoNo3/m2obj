package m2yaml

import (
	"github.com/rickonono3/m2obj"
	"gopkg.in/yaml.v3"
)

type Formatter struct {
}

func (f Formatter) Marshal(obj *m2obj.Object) (data []byte, err error) {
	return yaml.Marshal(obj.Staticize())
}

func (f Formatter) Unmarshal(data []byte) (obj *m2obj.Object, err error) {
	m := make(map[string]interface{})
	err = yaml.Unmarshal(data, &m)
	if err == nil {
		obj = m2obj.New(m)
	}
	return
}
