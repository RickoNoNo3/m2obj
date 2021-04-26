package m2obj

type Formatter interface {
	Marshal(obj *Object) (data []byte, err error)
	UnMarshal(data []byte) (obj *Object, err error)
}
