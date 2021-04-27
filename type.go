package m2obj

type Formatter interface {
	Marshal(obj *Object) (data []byte, err error)
	Unmarshal(data []byte) (obj *Object, err error)
}
