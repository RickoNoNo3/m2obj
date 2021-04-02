# M2Obj
A JSON-like, dynamic, persistent OBJECT STRUCTURE for configurations, caches, Go Template data or just to store dynamic JSON objects in Go.



## Install

```shell
go get github.com/rickonono3/m2obj
```


## Usages & Examples

- As a map/JSON binder
- As a configuration manager
- As a Go Template data wrapper



### As a map/JSON binder

M2Obj can help you to operate a map/JSON or any other JSON-like data easily.

```go
package main

import (
	"fmt"

	"github.com/rickonono3/m2obj"
)

var m = map[string]interface{}{
	"a": 1,
	"b": "2",
	"c": true,
	"d": map[string]interface{}{
		"e": "3",
	},
}

func main() {
	obj := m2obj.NewFromMap(m)
	_ = obj.Set("d.f.g", 4)
	m2 := obj.Staticize()
	fmt.Println(m2)
}
```

Then the `m2` will be:

```go
var m2 = map[string]interface{}{
	"a": 1,
	"b": "2",
	"c": true,
	"d": map[string]interface{}{
		"e": "3",
		"f": map[string]interface{}{
			"g": 4,
		},
	},
}
```

You can also use the `JsonDataFormatter` to transform a JSON string, like:

```go
package main

import (
	"fmt"

	"github.com/rickonono3/m2obj/m2json"
)

var jsonStr = `{
  "a": 1,
  "b": "2",
  "c": true,
  "d": {
    "e": "3"
  }
}`

func main() {
	dataFormatter := m2json.New()
	obj, err := dataFormatter.UnMarshal(jsonStr)
	if err == nil {
		_ = obj.Set("d.f.g", 4)
		jsonStr2, _ := dataFormatter.Marshal(obj)
		fmt.Println(jsonStr2)
	}
}
```



### As a configuration manager

Easily Get/Set for configurations with any structure you like.

`config/config.go`:

```go
package config

import (
	"github.com/rickonono3/m2obj"
)

const (
	LevelInfo = iota
	LevelWarn
	LevelError
)

var Config *m2obj.Object

func init() {
	Config = m2obj.New(m2obj.Group{
		"Debug": m2obj.Group{
			"IsDebugging": true,
			"Level":       LevelWarn,
		},
	})
}
```
`main.go`:

```go
package main

import (
	"fmt"

	"m2obj_demo/config"
)

func debugPrint(str string, level int) {
	debug := config.Config.MustGet("Debug")
	if debug.MustGet("IsDebugging").ValBool() {
		if level >= debug.MustGet("Level").ValInt() {
			fmt.Println(">> " + str)
		}
	}
}

func main() {
	debugPrint("This is Info1", config.LevelInfo)
	debugPrint("This is Warn1", config.LevelWarn)
	debugPrint("This is Error1", config.LevelError)

	fmt.Println("----------")
	_ = config.Config.Set("Debug.Level", config.LevelError)

	debugPrint("This is Info2", config.LevelInfo)
	debugPrint("This is Warn2", config.LevelWarn)
	debugPrint("This is Error2", config.LevelError)
}
```

`stdout`:

```
>> This is Warn1
>> This is Error1
----------
>> This is Error2
```



### As a Go Template data wrapper

Make use of `Staticize()`, the Group object can be easily transformed to an `map[string]interface{}`. You can append global configurations to Go Template in one line.

`main.go`:

```go
package main

import (
	"html/template"
	"os"

	"github.com/rickonono3/m2obj"
)

var Config = m2obj.New(m2obj.Group{
	"cdn": "https://example.com",
})

func main() {
	t, err := template.ParseFiles("index.gohtml")
	if err == nil {
		_ = t.Execute(os.Stdout, m2obj.New(m2obj.Group{
			"config": Config,
			"title":  "M2Obj Demo",
			"body": m2obj.Group{
				"h1":   "M2Obj Demo for Go Template Data Wrapper",
				"text": "Enjoy!",
			},
		}).Staticize())
	}
}
```

`index.gohtml`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <title>{{.title}}</title>
  <script src="{{.config.cdn}}/index.js"></script>
</head>
<body>
<h1>{{.body.h1}}</h1>
<p>{{.body.text}}</p>
</body>
</html>
```

`stdout`:

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <title>M2Obj Demo</title>
  <script src="https://example.com/index.js"></script>
</head>
<body>
<h1>M2Obj Demo for Go Template Data Wrapper</h1>
<p>Enjoy!</p>
</body>
</html>
```



## Docs

> For the detailed documents and APIs of M2Obj, See [THIS](https://github.com/rickonono3/m2obj).

### Data Types

| Type Name       | Type Description               | Note                                                         |
| --------------- | ------------------------------ | ------------------------------------------------------------ |
| `Object`        | `type Object struct`           | The base type of all object nodes. Always appears as `*Object` |
| `Group`         | `map[string]interface{}`       | Used like a JavaScript object                                |
| `Array`         | `[]interface{}`                | Used like a JavaScript array                                 |
| `DataFormatter` | `type DataFormatter interface` | Transform the object to a given data formet (like JSON, XML, etc.) |



### Special Agreement

**Key String**

- To locale an object(element) simply. Used by `Get`/`Set`/`Has`/`Remove`.

- Called `keyStr` in the code.

- Example: `"A.B.[0].C"`

- Explain: the example means that ***Find an unconstrained object called `C` in the group object that is  found as the `[0]` element in an array object called `B` that is found in a group object called `A`***.

- In other words: 

  1. The last fragment is an unconstrained object which can be any type, like the `C`.
  2. The fragments followed by `[index]` must be array objects and the `index` must be valid, like the `B.[0]`.
  3. All of other fragments must be group objects, like the `A`.

- ```go
  var obj = m2obj.New(m2obj.Group{
    "A": m2obj.Group{
      "B": m2obj.Array{
        m2obj.Group{
          "C": "I am here!",
        },
      },
    },
  })
  ```



### Functions

| Function                                     | Note                                                         |
| -------------------------------------------- | ------------------------------------------------------------ |
| `New(interface{}) *Object`                   | Create an object. Use `New(Group{...})` / `New(Array{...})` to create a multielement object |
| `NewFromMap(map[string]interface{}) *Object` | Create a group object on behalf of a map                     |



### Methods

`*Object` Base:

| Method                                                      | Note                                                         |
| ----------------------------------------------------------- | ------------------------------------------------------------ |
| `Set(keyStr string, value interface{}) (err error)`         | Set the value of a child of the object assigned by the `keyStr`. If it exists, replace its value, or else, create it. |
| `SetIfHas(keyStr string, value interface{}) (err error)`    | Set value when only the child exists.                        |
| `SetIfNotHas(keyStr string, value interface{}) (err error)` | Set value when only the child doesn't exist.                 |
| `Get(keyStr string) (obj *Object, err error)`               | Get the value of a child of the object assigned by the `keyStr`. If it exists, returns `obj, nil`, or else, returns `nil, err`. |
| `MustGet(keyStr string) (obj *Object)`                      | Similar to `Get`, but panic when the child doesn't exist.          |
| `Has(keyStr string) bool`                                   | Returns if the child assigned by the `keyStr` exists.        |
| `Remove(keyStr string) bool`                                | Remove a child (and its children as well) assigned by the `keyStr`. If the removing is successful or the child doesn't exist at all, return `true`, or else, return `false`. |
| `SetVal(value interface{})`                                 | Set value of the object itself, not for its child.           |
| `Val() interface{}`                                         | Get value of the object itself, as a type of `interface{}`. You can do your own operations on it, like `switch (type)`. |
| `ValStr() string`                                           | Get value of the object itself, and assert that it is a `string`. |
| `ValBool() bool`                                            | Get value of the object itself, and assert that it is a `bool`. |
| `ValInt() int`                                              | Get value of the object itself, and assert that it is a `int`. |
| `ValInt8() int8`                                            | Get value of the object itself, and assert that it is a `int8`. |
| `ValInt16() int16`                                          | Get value of the object itself, and assert that it is a `int16`. |
| `ValInt32() int32`                                          | Get value of the object itself, and assert that it is a `int32`. |
| `ValInt64() int64`                                          | Get value of the object itself, and assert that it is a `int64`. |
| `ValFloat32() float32`                                      | Get value of the object itself, and assert that it is a `float32`. |
| `ValFloat64() float64`                                      | Get value of the object itself, and assert that it is a `float64`. |
| `Staticize() map[string]interface{}`                        | Peel the object and all of its children. If the object is a Group, the method returns a `map[string]interface{}` contains its children directly. If it is an Array, the method returns `map[string]interface{}{"list": []interface{}}` and its children will be push into the `list`. Or else, the method returns `map[string]interface{}{"val":interface{}}` and the `val` will be the value of the object. As for children, Group children to `map[string]interface{}` and Array children to `[]interface{}`. |
| `Clone() (newObj *Object)`                                  | Deep clone an object.                                        |
|                                                             |                                                              |

`*Object` as an Array:

| Method                                                    | Note                                 |
| --------------------------------------------------------- | ------------------------------------ |
| `ArrPush(value interface{}) (err error)`                  |                                      |
| `ArrPop() (err error)`                                    |                                      |
| `ArrSet(index int, value interface{}) (err error)`        |                                      |
| `ArrGet(index int) (obj *Object, err error)`              | An alias of `*Object.Get("[index]")` |
| `ArrInsert(index int, value interface{}) (err error)`     |                                      |
| `ArrRemove(index int) (err error)`                        |                                      |
| `ArrForeach(do func(index int, obj *Object)) (err error)` |                                      |
| `ArrLen() (int, error)`                                   |                                      |



# TODO

- [ ] `IsGroup` / `IsArray` / `IsValue`
- [ ] More `Arr*` Methods
- [ ] More `DataFormatter`
- [ ] Performance Optimizations and Bench Tests.
- [ ] Stronger type definition
