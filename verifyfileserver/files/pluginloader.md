# pluginloader:简化go语言plugin函数和对象调用
*2019-5-11 新增*

UnknownObject 类型增加了两个新方法：`Json` 、`CopyToStruct`，前一个导出JSON，后一个把结构体的可导出值复制到另一个相似的结构体中。

*2019-4-19 新修改*

- 修改 `pluginwrap` ，删除生成的文件中的 `InitxxxFuncs`，改为：`func GetxxxFuncs(p *pluginloader.PluginLoader) *xxxFuncs`， 函数都被包含在 `xxxFuncs`结构体中避免命名冲突。
- `pluginloader`新增用于调用未定义结构体的对象：UnknownObject。

*2019-4-18 新修改：pluginwrap几乎已经完美了！ 现在除了用户自定义类型，已经可以使用所有被导入库的类型。*

go语言的`plugin`调用还是很繁琐，而且功能限制也不少，我开发了的这个程序包含一个包`github.com/rocket049/pluginloader`和一个命令行工具`github.com/rocket049/pluginloader/cmd/pluginwrap`。

包`github.com/rocket049/pluginloader`引用到程序中，用于简化函数调用。

命令行工具`pluginwrap`用于从`plugin`的源代码生成可导出对象`struct`对应的`interface`和可导出`func`。

### 使用`pluginloader`
安装命令： `go get github.com/rocket049/pluginloader`

#### 可调用的函数：

```
type PluginLoader struct {
	P *plugin.Plugin
}

///Call return type must be: (res,error)
func (p *PluginLoader) Call(funcName string, p0 ...interface{}) (interface{}, error)

//CallValue Allow any number of return values,return type: []reflect.Value,error
func (p *PluginLoader) CallValue(funcName string, p0 ...interface{}) ([]reflect.Value, error)

//MakeFunc point a func ptr to plugin
func (s *PluginLoader) MakeFunc(fptr interface{}, name string) error 

//20190419 新增
//UnknownObject 成员'V' 必须是结构体指针的 Value: *struct{...}
type UnknownObject struct {
	V reflect.Value
}
//NewUnknownObject 参数'v' 必须是结构体指针的 Value: *struct{...}，否则返回 nil
func NewUnknownObject(v reflect.Value) *UnknownObject 

//Get 得到结构体成员的 Value
func (s *UnknownObject) Get(name string) reflect.Value

//Call 运行结构体的 method
func (s *UnknownObject) Call(fn string, args ...interface{}) []reflect.Value

//Json 把结构体编码为 JSON。 convert the struct to JSON. if error,return nil.
func (s *UnknownObject) Json() []byte 

//CopyToStruct 利用 reflect 技术把结构体的可 export 值复制到 v 中，v 必须是相似结构体的指针。 copy the exported value of a struct to v 
func (s *UnknownObject) CopyToStruct(v interface{}) error 
```
#### `pluginloader`导出`method`的具体说明

- Call，直接通过函数名调用对应的函数，但是返回值形式受到限制。
- CallValue，直接通过函数名调用对应的函数，返回值形式不受限制。
- MakeFunc，根据名字把预定义的函数类型变量指向插件中的函数。

#### `UnknownObject`说明

- `func NewUnknownObject(v reflect.Value) *UnknownObject`，参数'v' 必须是结构体指针的 Value: *struct{...}，否则返回 nil
- Method：`Get(name string) reflect.Value`，得到结构体成员的 Value
- Method：`Call(fn string, args ...interface{}) []reflect.Value`，运行结构体的 method

#### 调用示例
```
import "github.com/rocket049/pluginloader"

p, err := pluginloader.NewPluginLoader( "foo.so" )
if err != nil {
	panic(err)
}

res, err := p.Call("NameOfFunc", p0,p1,p3,...)
// ...

ret := p.CallValue("NameOfFunc", p0,p1,p3,...)
// ...

var Foo func(arg string)(string,error)
p.MakeFunc(&Foo,"Foo")
ret, err = Foo("something")
// ...

// 使用 UnknownObject. NewFoo return 'foo *Foo'
v, err := p.CallValue("NewFoo")
if err != nil {
	t.Fatal(err)
}
obj := NewUnknownObject(v[0])

id: = obj.Get("Id").Int()

err = obj.Call("Set", nil)
```
### 使用`pluginwrap`
安装命令： `go get github.com/rocket049/pluginloader/cmd/pluginwrap`

#### 用法
`pluginwrap path/to/plugin/foo`

生成的文件名字：
`fooWrap.go`

把这个文件加入你的工程就可以方便的调用`plugin`的导出函数了。

#### 功能
1. 生成导出对象的接口，以便用于类型断言。
2. 生成导出函数。

#### 限制
本程序基于标准包`plugin`和`reflect`实现，因为`go`语言的变量类型转换的使用有很多限制，所以本程序对导出函数的参数类型、返回值类型都有限制， ***导出参数、返回值的类型仅限于基本类型、标准库、第三方库中的类型，不能使用自定义类型。***

如果必须使用自定义类型，有两种办法：

1. 请使用 `pluginloader.Call` 或 `pluginloader.CallValue` 调用。
2. 把需要导出的复杂类型做成第三方库（`import "your/package"`），不要在`plugin`源代码中定义。

### 示例

#### 使用对象(`struct`)

```
	p, err := pluginloader.NewPluginLoader("foo.so")
	if err != nil {
		panic(err)
	}
	iface, err := p.Call("NewFoo",arg...)
	if err != nil {
		panic(err)
	}
	foo := iface.(IFoo)
	// use foo.Method
```

#### 使用`func`

```
	//2090419 new
	p, err := pluginloader.NewPluginLoader("foo.so")
	if err != nil {
		panic(err)
	}
	
	// MUST call GetxxxFuncs(p) before call funcs, xxx = plugin名字
	funcs := GetfooFuncs(p)
	
	// call methods in plugin foo
	funcs.Method(arg)
	
```