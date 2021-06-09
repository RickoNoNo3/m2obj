|Language|
|:------:|
|[English](https://github.com/RickoNoNo3/m2obj)|
|[中文](https://github.com/RickoNoNo3/m2obj/README_CN.md)|

# M2Obj

一个类JSON的、动态的、可持久化的 Golang 【对象结构】，可用于管理配置项、缓存、模板引擎数据，也可以单纯用来存储动态JSON格式对象。

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
	obj := m2obj.New(m)
	_ = obj.Set("d.f.g", 4)
	m2 := obj.Staticize()
	fmt.Println(m2)
}
```

Then the `m2` will be:

```json
{
  "a": 1,
  "b": "2",
  "c": true,
  "d": {
    "e": "3",
    "f": {
      "g": 4
    }
  }
}
```

You can also use the `m2json.Formatter` to transform a JSON string, like:

```go
package main

import (
	"fmt"

	"github.com/rickonono3/m2obj/m2json"
)

var json = []byte(`{
	"a": 1,
	"b": "2",
	"c": true,
	"d": {
		"e": "3"
	}
}`)

func main() {
	formatter := m2json.Formatter{}
	obj, err := formatter.Unmarshal(json)
	if err == nil {
		_ = obj.Set("d.f.g", 4)
		json, _ = formatter.Marshal(obj)
		fmt.Println(string(json))
	}
}

```

Note that: As a character of `json` package of Go, the number variables are always parsed as `float64`. Strictly using the `ValXxx()` methods only, lets you can ignore this character because of the internal implementation of M2Obj. Or, such as to use `Val()`, you have to check it by yourself.

By the way, you can implement the `Formatter` interface easily by yourself, to support more customized functions.

### As a configuration manager

Easily Get/Set for configurations with any structure you like. There is a `FileSyncer` to sync between your config file and m2obj object.

The following example demonstrates filtering DEBUG output by changing the global DEBUG level:

`config/config.go`:

```go
package config

import (
	"github.com/rickonono3/m2obj"
	"github.com/rickonono3/m2obj/m2json"
)

const (
	LevelInfo = iota
	LevelWarn
	LevelError
)

var Config *m2obj.Object
var FileSyncer *m2obj.FileSyncer

func init() {
	Config = m2obj.New(m2obj.Group{
		"Debug": m2obj.Group{
			"IsDebugging": true,
			"Level":       LevelWarn,
		},
	})
	FileSyncer = m2obj.NewFileSyncer("./config.json", m2json.Formatter{})
	FileSyncer.BindObject(Config)
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

Make use of `Staticize()`, the Group object can be easily transformed to an `map[string]interface{}`. You can append global configurations to Go Template in one line. Of course, you can also perform more data operations on it.

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

> For the detailed documents and APIs of M2Obj, See [THIS](https://pkg.go.dev/github.com/rickonono3/m2obj).

### Data Types

| Type Name | Type Description | Note |
| --------- | ---------------- | ---- |
| `Object` | `type Object struct` | The base type of all object nodes. Always appears as `*Object` |
| `Group` | `map[string]interface{}`| Used like a JSON object |
| `Array` | `[]interface{}` | Used like a JSON array |
| `Formatter` | `type Formatter interface` | Converts the object from/to a given data format (like JSON, XML, etc.) |
| `FileSyncer` | `type FileSyncer struct` | Syncs between files and memory, uses Formatter |

### Special Definition

**Object Type**

- All elements have the same type: `*Object`.
- There are three Object Types: `Group`, `Array` and `Value`. They can only be differentiated by `IsGroup`, `IsArray` and `IsValue`.
- `Group` is a key-value map.
  - Definition: `map[string]interface{}`.
  - Like `{}` in JSON.
  - To create a Group Object, use `m2obj.New(m2obj.Group{"k1":v1,"k2":v2 ...})`.
- `Array` is an array(or, slice).
  - Definition: `[]interface{}`.
  - Like `[]` in JSON.
  - To create an Array Object, use `m2obj.New(m2obj.Array{v1,v2 ...})`.
- `Value` is any other type of value.
  - The inner val of a Value Object will never be `Object`/`*Object`, if the `New()` or `SetVal()` called with an Object param, it will be dismounting by a private method named `getDeepestValue`. It means, All the methods that have `interface{}` params can be called with a wrapped `Object` or just a value, they are all worked.

**Key String**

- To locale an object(element) simply. Used by `Get`/`Set`/`Has`/`Remove`.

- Named `keyStr` in the code.

- Example: `"A.B.[0].C"`

- Explain: the example means that ***Group `A` -> Array `B` -> Group `[0]` -> Any `C`***.

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

| Function | Note |
| -------- | ---- |
| `New` | Create an object. Use `New(Group{...})` / `New(Array{...})` to create multi-element objects |
| `NewFileSyncer` | Create a FileSyncer |

### Methods / Fields

`*Object` Base:

| Method / Field | Note |
| -------------- | ---- |
| `Set()` | Set the value of a child of the object assigned by the `keyStr`. If it exists, replace its value, or else, create it. |
| `SetIfHas()` | Set value when only the child exists. |
| `SetIfNotHas()` | Set value when only the child doesn't exist. |
| `Get()` | Get the value of a child of the object assigned by the `keyStr`. If it exists, returns `obj, nil`, or else, returns `nil, err`. |
| `MustGet()` | Similar to `Get`, but panic when the child doesn't exist. |
| `Has()` | Returns if the child assigned by the `keyStr` exists. |
| `Remove()` | Remove a child (and its children as well) assigned by the `keyStr`. If the removing is successful or the child doesn't exist at all, return `true`, or else, return `false`. |
| `SetVal()` | Set the inner value of an Object. |
| `Val()` | Get the inner value of an Object, as a type of `interface{}`. You can do your own operations on it, like `switch (type)` and `.(type)`, and even some `reflect` methods. |
| `ValStr()` | Get the inner value of an Object, and assert it is or transform it to a `string`. |
| `ValBool()` | Get the inner value of an Object, and assert it is or transform it to a `bool`. |
| `ValByte()` | Get the inner value of an Object, and assert it is or transform it to a `byte`. |
| `ValBytes()` | Get the inner value of an Object, and assert it is or transform it to a `[]byte`. |
| `ValRune()` | Get the inner value of an Object, and assert it is or transform it to a `rune`. |
| `ValRunes()` | Get the inner value of an Object, and assert it is or transform it to a `[]rune`. |
| `ValInt()` | Get the inner value of an Object, and assert it is or transform it to an `int`. |
| `ValInt8()` | Get the inner value of an Object, and assert it is or transform it to an `int8`. |
| `ValInt16()` | Get the inner value of an Object, and assert it is or transform it to an `int16`. |
| `ValInt32()` | Get the inner value of an Object, and assert it is or transform it to an `int32`. |
| `ValInt64()` | Get the inner value of an Object, and assert it is or transform it to an `int64`. |
| `ValUint()` | Get the inner value of an Object, and assert it is or transform it to an `uint64`. |
| `ValFloat32()` | Get the inner value of an Object, and assert it is or transform it to a `float32`. |
| `ValFloat64()` | Get the inner value of an Object, and assert it is or transform it to a `float64`. |
| `Staticize()` | Peel the object and all of its children to a `map[string]interface{}` |
| `Clone()` | Deep clone an object. |
| `Is()` | Use `reflect` to judge the type of an Object's value. |
| `IsLike()` | Use `reflect` to compare and judge if the type of the Object's value is same as a variable. |
| `IsNil()` | Judge if the Object's value is `nil`. |
| `IsGroup()` | Judge if the Object is a Group Object |
| `IsArray()` | Judge if the Object is a Array Object |
| `IsValue()` | Judge if the Object is a Value Object |
| `Parent()` | Get the parent Object of an Object, if the Object is root node, return `nil` |

`*Object` as a Group:

| Method / Field | Note |
| -------------- | ---- |
| `GroupMerge()` | Merge another Group Object to this Group. Enable the forced option to force replacement when the key already exists. |
| `GroupForeach()` | |

`*Object` as an Array:

| Method / Field | Note |
| -------------- | ---- |
| `ArrPush()` | |
| `ArrMerge()` | Push another Array Object to an Array Object |
| `ArrPushAll()` | Push all params(variable-length) to an Array Object |
| `ArrPop()` | |
| `ArrShift()` | |
| `ArrUnshift()` | |
| `ArrSet()` | |
| `ArrGet()` | An alias of `*Object.Get("[index]")` |
| `ArrInsert()` | |
| `ArrRemove()` | |
| `ArrForeach()` | |
| `ArrLen()` | |

`Formatter`:

| Method / Field | Note |
| -------------- | ---- |
| `Marshal()` | Transform an Object to bytes |
| `Unmarshal()` | Transform bytes to an Object |

`*FileSyncer`:

| Method / Field | Note |
| -------------- | ---- |
| `Load()` | Load from the file |
| `Save()` | Save to the file |
| `SetFilePath()` | |
| `GetFilePath()` | |
| `SetFormatter()` | |
| `BindObject()` | Bind a Group Object to start syncing |
| `GetBoundObject()` | |
| `HardLoad` | `bool`, appoint the behavior of `Load()`. If `true`, the loading will remove all the keys in the bound Object that are not found in the current file, or the keys will be kept. (Default: `false`) |
| `AutoSaveTime` | `int64`, the milliseconds interval to trigger `Save()`. If it is less than 0, auto saving is disabled. If it equals to 0, auto saving is triggered when the Object changed. If it is greater than 0, auto saving is triggered on each interval. (Default: 0) |
| `AutoLoadTime` | `int64`, the milliseconds interval to trigger `Load()`. If it is <= 0, auto loading is disabled. Or else, auto loading is triggered on each interval and **auto saving is disabled whether the `AutoSaveTime` is**. (Default: 0) |

# TODO

- [x] `IsGroup` / `IsArray` / `IsValue`
- [x] More `Arr` Methods
- [ ] More `Formatter`
- [ ] Performance Optimizations and Bench Tests.
- [x] Stronger type definition
