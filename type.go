package m2obj

// Formatter
//
// Converting between []byte and *Object
type Formatter interface {
	Marshal(obj *Object) (data []byte, err error)
	Unmarshal(data []byte) (obj *Object, err error)
}
