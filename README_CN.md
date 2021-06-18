|语言|
|:------:|
|[English](https://github.com/RickoNoNo3/m2obj)|
|[中文](https://github.com/RickoNoNo3/m2obj/blob/master/README_CN.md)|

# M2Obj

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/rickonono3/m2obj/workflow?label=built%20%26%20test&style=flat-square)

一个类JSON的、动态的、可持久化的 Golang 【对象结构】, 可用于管理配置项、缓存、模板引擎数据, 也可以单纯用来存储动态JSON格式对象.

A JSON-like, dynamic, persistent OBJECT STRUCTURE for configurations, caches, Go Template data or just to store dynamic JSON objects in Go.

## 为何使用 M2Obj

- Go 原生的处理树形动态数据的方式需要频繁地对每层数据进行类型断言, 代码过长. 而 M2Obj 只需在定位数据后使用已封装好的取值方法即可.
- Go 原生的对 JSON 结构数据的支持有所欠缺. M2Obj 进行了大量封装和改善, 并且专为 JSON / 类 JSON 格式做了适配性开发.
- M2Obj 拥有高效、并发的 Goroutine 文件同步器, 在内存数据和文件间进行自动同步
- Without M2Obj:
  ```go
  var M = map[string]interface{}{
    "info": map[string]interface{}{
      "name": "UKP",
      "schools": []interface{}{
        "Engineering|engineering@ukp.edu|135",
        map[string]interface{}{
          "name":         "Law",
          "email":        "law@ukp.edu",
          "studentCount": 300,
          "notice":       "",
        },
      },
    },
  }

  // set student count of law school to 295
  M["info"].(map[string]interface{})["schools"].([]interface{})[1].(map[string]interface{})["studentCount"] = 295
  // delete notice of law school
  delete(M["info"].(map[string]interface{})["schools"].([]interface{})[1].(map[string]interface{}), "notice")
  // println "UKP"
  fmt.Println(M["info"].(map[string]interface{})["name"].(string))
  // println data as map
  fmt.Println(M)
  ```
- With M2Obj:
  ```go
  var M2 = m2obj.New(M) // m2obj can be constructed in a lot of styles, see examples below.

  // set student count of law school to 295
  M2.Set("info.schools.[1].studentCount", 295)
  // delete notice of law school
  M2.Remove("info.schools.[1].notice")
  // println "UKP"
  fmt.Println(M2.MustGet("info.name").ValStr())
  // println data as map
  fmt.Println(M2.Staticize())
  ```

## 安装

```shell
go get github.com/rickonono3/m2obj
```

## 用法和示例

> 以下示例程序的完整代码见 [RickoNoNo3/m2obj_examples](https://github.com/RickoNoNo3/m2obj_examples)

- 作为 map/JSON 绑定器
- 作为 配置管理器
- 作为 Go Template 数据封装器

### 作为 map/JSON 绑定器

M2Obj可以让你轻松操作 map/JSON 或其他任何类JSON格式的数据.

```go
// initial map
var m = map[string]interface{}{
  "a": 1,
  "b": "2",
  "c": true,
  "d": map[string]interface{}{
    "e": "3",
  },
}

func main() {
  // new Object with param map[string]interface{}
  obj := m2obj.New(m)
  _ = obj.Set("d.f.g", 4)
  // staticize the object to map
  m2 := obj.Staticize()
  fmt.Println(m2)
}
```

运行程序, `m2`即为:

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

为了在map和JSON格式间转换, 引入`m2json.Formatter`:

```go
// initial json string
var jsonStr = `{
 "a": 1,
 "b": "2",
 "c": true,
 "d": {
   "e": "3"
 }
}`

func main() {
  // new Formatter
  formatter := m2json.Formatter{}
  // convert jsonStr([]byte) to a new Object
  obj, err := formatter.Unmarshal([]byte(jsonStr))
  if err == nil {
    _ = obj.Set("d.f.g", 4)
    // convert Object to jsonStr2([]byte)
    jsonStr2, _ := formatter.Marshal(obj)
    fmt.Println(string(jsonStr2))
  }
}
```

请注意: Go的`json`包有一个特性, JSON字符串中的数字类型总是被解析为`float64`而不管有没有小数点. 严格保证只使用`ValXxx()`系列方法可以规避此特性, 因为M2Obj做了内部实现. 否则, 比如使用`Val()`方法, 你必须手动处理它.

另外, 你可以轻松实现一个自己的`Formatter`接口来支持许多自定义功能或者序列化格式.

### 作为配置管理器

在结构化的配置项中轻松进行 Get/Set. 同时, 可以使用`FileSyncer`来在 M2Obj Object 和你的配置文件间进行同步.

下面的示例演示了通过更改全局 DEBUG 级别来过滤 DEBUG 输出:

`config.go`:

```go
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
  // FileSyncer
  FileSyncer = m2obj.NewFileSyncer("./config.json", m2json.Formatter{})
  FileSyncer.BindObject(Config)
  // FileSyncer.Load()

  // DEFAULT FILE_SYNCER OPTIONS:
  //   Auto Saving  : On bound object changes
  //   Auto Loading : Never
  //   Hard Load    : False
}
```

`main.go`:

```go
// Print str if IsDebugging && level >= DebugLevel
func debugPrint(str string, level int) {
  debug := Config.MustGet("Debug")
  if debug.MustGet("IsDebugging").ValBool() {
    if level >= debug.MustGet("Level").ValInt() {
      fmt.Println(">> " + str)
    }
  }
}

func main() {
  debugPrint("This is Info1", LevelInfo)   // filtered
  debugPrint("This is Warn1", LevelWarn)   // printed
  debugPrint("This is Error1", LevelError) // printed

  fmt.Println("----------")
  _ = Config.Set("Debug.Level", LevelError)

  debugPrint("This is Info2", LevelInfo)   // filtered
  debugPrint("This is Warn2", LevelWarn)   // filtered
  debugPrint("This is Error2", LevelError) // printed
}
```

`stdout`:

```
>> This is Warn1
>> This is Error1
----------
>> This is Error2
```

### 作为 Go Template 数据绑定器

只需一行, 就可以将全局配置附加到 Go Template 的数据中. 当然也可以在其上进行更多数据操作.

利用 `Staticize()`, 可以轻松地将 Group 对象转换为 `map[string]interface{}`.

`main.go`:

```go
// Config is the global configuration
var Config = m2obj.New(m2obj.Group{
  "cdn": "https://example.com",
})

func main() {
  // new template
  t, err := template.ParseFiles("index.gohtml")
  if err == nil {
    // define data of the executed template
    data := m2obj.New(m2obj.Group{
      "title": "M2Obj Examples",
      "body": m2obj.Group{
        "h1":   "M2Obj Examples for Go Template Data Wrapper",
        "text": "Enjoy!",
      },
    })
    // add the config object to data of the executed template.
    data.Set("config", Config)
    // staticize the data
    _ = t.Execute(os.Stdout, data.Staticize())
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

## 文档

> M2Obj 的详细文档和 API 参见 [这里](https:pkg.go.devgithub.comrickonono3m2obj).

### 数据类型

| 类型名 | 定义 | 备注 |
| --------- | ---------------- | ---- |
| `Object` | `type Object struct` | 所有对象节点的基本类型. 始终以指针形式出现 |
| `Group` | `map[string]interface{}`| 像 JSON 对象一样 |
| `Array` | `[]interface{}` | 像 JSON 数组一样 |
| `Formatter` | `type Formatter interface` | 将对象转换为给定的数据格式 (如 JSON、XML 等) |
| `FileSyncer` | `type FileSyncer struct` | 在文件和内存之间同步, 使用`Formatter` |

### 特别约定

**Object Type**

- 所有元素都拥有相同的类型: `*Object`.
- 有三种 Object Types: `Group`, `Array` 以及 `Value`. 他们只能被 `IsGroup`, `IsArray` 和 `IsValue` 三个方法来区分.
- `Group` 是一种键值对.
  - 定义: `map[string]interface{}`.
  - 如同 JSON 中的 `{}`.
  - 要创建 Group Object, 这样写: `m2obj.New(m2obj.Group{"k1":v1,"k2":v2 ...})`.
- `Array` 是一种切片(slice).
  - 定义: `[]interface{}`.
  - 如同 JSON 中的 `[]`.
  - 要创建 Array Object, 这样写: `m2obj.New(m2obj.Array{v1,v2 ...})`.
- `Value` 是任何其他值.
  - 一个 Value Object 内维护的实际值永远不可能是 `Object`/`*Object`, 如果调用 `New()` 或 `SetVal()` 时传入了一个 Object 参数, 他会被一个私有方法 `getDeepestValue` 自动拆解. 也就是说, 所有类型为 `interface{}` 的参数, 都可以往进传 Object 或者 裸值, 这不影响最后存储的结果.

**Key String**

- Key String 用来方便地定位元素. 在 `Get`/`Set`/`Has`/`Remove` 中都有使用.

- 在代码中名为 `keyStr`.

- 示例: `"A.B.[0].C"`

- 示例解释: 这个示例表示 ***Group `A` -> Array `B` -> Group `[0]` -> Any `C`***.

- 换句话说:

  1. 最后一段是不规定 Object Type 的, 如 `C`.
  2. 后跟 `[下标]` 的段必须是 Array Object, 并且 `下标` 必须合法, 如 `B.[0]`.
  3. 所有其他段必须是 Group Objects, 如 `A`.

- 示例`keyStr`实际上反映了如下结构:
  ```go
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

### 函数

| 函数名 | 说明 |
| -------- | ---- |
| `New` | 创建一个 Object. 也可用 `New(Group{...})` / `New(Array{...})` 创建多元集合形式的 Object |
| `NewFileSyncer` | 创建一个 FileSyncer |

### 方法 / 属性

`*Object` 基本:

| 方法 / 属性 | 说明 |
| -------------- | ---- |
| `Set()` | 设置由 `keyStr` 定位的 Object 的值. 如果此位置存在原有数据则替换, 否则在该位置创建新值. |
| `SetIfHas()` | 仅在定位位置存在时`Set`. |
| `SetIfNotHas()` | 仅在定位位置不存在时`Set`. |
| `Get()` | 获取由 `keyStr` 定位的 Object 的值. 存在返回 `obj, nil`, 不存在返回 `nil, err`. |
| `MustGet()` | 类似 `Get`, 但是在不存在时爆出 panic. 单返回值便于连写. |
| `Has()` | 检查 `keyStr` 是否存在. |
| `Remove()` | 删除由 `keyStr` 定位的项 (及其子项). 如果移除成功或者孩子根本不存在, 则返回“true”, 否则返回“false”. |
| `SetVal()` | 设置 Object 本身内部维护的值. |
| `Val()` | 获取 Object 的内部值, 作为 `interface{}` 类型. 你可以对它做你自己的操作, 比如`switch (type)`和`.(type)`, 甚至`reflect`包的操作. |
| `ValStr()` | 获取 Object 的内部值, 并断言或转换其为 `string`. |
| `ValBool()` | 获取 Object 的内部值, 并断言或转换其为 `bool`. |
| `ValByte()` | 获取 Object 的内部值, 并断言或转换其为 `byte`. |
| `ValBytes()` | 获取 Object 的内部值, 并断言或转换其为 `[]byte`. |
| `ValRune()` | 获取 Object 的内部值, 并断言或转换其为 `rune`. |
| `ValRunes()` | 获取 Object 的内部值, 并断言或转换其为 `[]rune`. |
| `ValInt()` | 获取 Object 的内部值, 并断言或转换其为 `int`. |
| `ValInt8()` | 获取 Object 的内部值, 并断言或转换其为 `int8`. |
| `ValInt16()` | 获取 Object 的内部值, 并断言或转换其为 `int16`. |
| `ValInt32()` | 获取 Object 的内部值, 并断言或转换其为 `int32`. |
| `ValInt64()` | 获取 Object 的内部值, 并断言或转换其为 `int64`. |
| `ValUint()` | 获取 Object 的内部值, 并断言或转换其为 `uint64`. |
| `ValFloat32()` | 获取 Object 的内部值, 并断言或转换其为 `float32`. |
| `ValFloat64()` | 获取 Object 的内部值, 并断言或转换其为 `float64`. |
| `Staticize()` | 静态化对象及其所有子对象到一个整体的 `map[string]interface{}` |
| `Clone()` | 深拷贝一个对象. |
| `Is()` | 使用`reflect`判断 Object 的内部值的类型. |
| `IsLike()` | 使用`reflect`判断 Object 的内部值的类型是否与某个给定变量相同. |
| `IsNil()` | 判断 Object 的内部值是否为 `nil`. |
| `IsGroup()` | 判断 Object 是否是一个 Group Object |
| `IsArray()` | 判断 Object 是否是一个 Array Object |
| `IsValue()` | 判断 Object 是否是一个 Value Object |
| `Parent()` | 获取 Object 的父 Object, 如果 Object 是根节点则返回`nil` |

`*Object` 作为 Group 时的特殊内容:

| 方法 / 属性 | 说明 |
| -------------- | ---- |
| `GroupMerge()` | 将另一个 Group Object 合并到该 Array Object. 启用 forced 选项来在 key 已经存在时强制替换 |
| `GroupForeach()` | |

`*Object` 作为 Array 时的特殊内容:

| 方法 / 属性 | 说明 |
| -------------- | ---- |
| `ArrPush()` | |
| `ArrMerge()` | 把另一个 Array Object 加到该 Array Object 后面 |
| `ArrPushAll()` | 将所有参数 (可变长度)加到该 Array Object 后面 |
| `ArrPop()` | |
| `ArrShift()` | |
| `ArrUnshift()` | |
| `ArrSet()` | |
| `ArrGet()` | `*Object.Get("[index]")` 的别名 |
| `ArrInsert()` | |
| `ArrRemove()` | |
| `ArrForeach()` | |
| `ArrLen()` | |

`Formatter`:

| 方法 / 属性 | 说明 |
| -------------- | ---- |
| `Marshal()` | 转换 Object 到 `[]byte` |
| `Unmarshal()` | 转换 `[]bytes` 到 Object |

`*FileSyncer`:

| 方法 / 属性 | 说明 |
| -------------- | ---- |
| `Load()` | 从文件加载 |
| `Save()` | 保存到文件 |
| `SetFilePath()` | |
| `GetFilePath()` | |
| `SetFormatter()` | |
| `BindObject()` | 绑定一个 Group Object 来开始同步 |
| `GetBoundObject()` | |
| `HardLoad` | `bool`, 指定 `Load()` 的行为. 如果为 `true` , 则在加载时清理加载源中有但绑定的 Object 中没有的所有键, 否则将保留这些键 (默认值: `false`) |
| `AutoSaveTime` | `int64`, 触发 `Save()` 的毫秒间隔. 如果 < 0, 则禁用自动保存. 如果 == 0, 则在对象更改时触发自动保存. 如果 > 0, 则在每个间隔触发自动保存 (默认值: 0) |
| `AutoLoadTime` | `int64`, 触发 `Load()` 的毫秒间隔. 如果 <= 0, 则禁用自动加载, 否则在每个间隔触发自动加载并且**屏蔽所有自动保存** (默认值: 0) |

# TODO

- [x] `IsGroup` / `IsArray` / `IsValue`
- [x] 更多 `Arr` 方法
- [ ] 更多 `Formatter`
- [ ] 性能优化和基准测试
- [x] 更强的类型定义
- [ ] 解决 Struct 类型的 field 赋值问题
