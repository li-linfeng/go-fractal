# Fractal Go 库

## 一、简介

`fractal` 是一个用于将数据资源转换为 `JSON` 格式的开源库。它允许你轻松地定义资源对象和转换器，以便在 API 响应中返回复杂的、嵌套的 JSON 数据。

## 二、 功能

- **资源定义**：定义您的数据资源对象，并指定如何将其转换为 JSON 格式。
- **转换器支持**：创建可自定义的转换器，以控制资源的转换过程。
- **嵌套包含**：支持嵌套的包含关系，使您能够在 API 响应中嵌套包含相关资源的数据。
- **项目和集合资源**：处理单个和基于列表的数据结构。

## 三、安装

要安装 `fractal` 库，请使用 `go get` 命令：

```go
go get github.com/li-linfeng/go-fractal
```

## 四、使用

### 4.1 基本用法

1. 定义您的转换器：

```go
type UserTransformer struct {
	fractal.Transformer
}

func (ut *UserTransformer) Transform() map[string]interface{} {
	data := ut.GetData().(User) // 假设要返回的数据模型是User
	return map[string]interface{}{
		"id":    data.ID,
		"name":  data.Name,
	}
}
```

2. 转换单个资源 `ItemResource`

```go
user := User{ID: 1, Name: "John Doe"}
resource := &fractal.ItemResource{}
resource.SetData(user)
resource.SetTransformer(&UserTransformer{})

//你可以获取转为json的数据
if resp, err := resource.ToJson(); err != nil {
    //错误处理
    return
}
fmt.Println(resp)
```

3. 转换集合 `CollectionResource`

```go
users := []User{
    {ID: 1, Name: "John Doe", },
    {ID: 2, Name: "Jane Smith"},
}

resource := &fractal.CollectionResource{}
resource.SetData(users)
resource.SetTransformer(&UserTransformer{})
if resp, err := resource.ToJson(); err != nil {
    //错误处理
    return
}
fmt.Println(resp)
```

### 4.2 资源类型

- `ItemResource` 单个资源
- `CollectionResource` 集合资源
- `PrimitiveResource` 基础类型资源

#### 4.2.1 资源接口 `ResourceInterface`

```go
type ResourceInterface interface {
    SetData(data interface{})//设置data
    SetIncludes(includes []string) //设置嵌套关系
    SetNestedInclude(map[string]interface{}) // 设置 格式化后的关联关系
    SetTransformer(transformer TransformerInterface) // 设置transformer
    TransformResource() error //执行转换
    GetTransformedData() interface{} //获取转换后的结果
    ToJson() (string, error) //获取转换后的json结果
}

```

#### 4.2.2 单个资源 `ItemResource`

```go
user := User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}
resource := &fractal.ItemResource{}
resource.SetData(user)
resource.SetTransformer(&UserTransformer{})

resource.ToJson() //{"id":1,"name":"John Doe", "email":"john.doe@example.com"}
```

> 当然你也可以不传值或者传空值， 此时`ToJson()`的结果为 `{}`

#### 4.2.3 集合资源 `CollectionResource`

```go
users := []User{
    {ID: 1, Name: "John Doe"},
    {ID: 2, Name: "Jane Smith"},
}
resource := &fractal.CollectionResource{}
resource.SetData(user)
resource.SetTransformer(&UserTransformer{})
resp, _ := resource.ToJson()

fmt.Println(resp) //[{"id":1,"name":"John Doe"},{"id":2,"name":"Jane Smith"}]
```

> 当然你也可以不传值或者传空值， 此时`ToJson()`的结果为 `[]`

#### 4.2.4 基础数据资源 `PrimitiveResource`

string, ~int, ~unit ,nil 等（除了 map 以及 struct 外）

```go
resource := &fractal.PrimitiveResource{}
resource.SetData("string") //"string"
resp, _ := resource.ToJson()
fmt.Println(resp)
// resource.SetData("nil") // null
// resource.SetData("12.3") // 12.3
// resource.SetData([]string{"a","b","c"}) // ["a","b","c"]
// resource.SetData([]int{1,2,3}) // [1,2,3]
```

#### 4.2.5 转换的结果

- 返回值为 `json`

  ```go
  if resp, err := resource.ToJson(); err != nil{
  	...
  }
  fmt.Println(resp)
  ```

- 返回值为 `map[string]interface{}`

  ```go
  if err:= resource.TransformResource; err !=nil{
  	...
  }
  fmt.Println(resource.GetTransformedData())
  ```

### 4.3 转换器 `Transformer`

`Transformer` 是执行数据转换的地方，写在一个统一的地方，可以很方便的让不同的接口使用相同的返回数据，代码复用。

```go
type TransformerInterface interface {
    //执行数据转换
    Transform() map[string]interface{}
    // 获取要转换的数据，即 resource的data
    GetData() interface{}
    // 注册include用到的function，如果没有注册则会无法加载include的数据
    SetAvailableIncludeFunctions()
}
//内置 Transformer ，
type Transformer struct {
	data                      interface{}
	AvailableIncludeFunctions map[string]func() ResourceInterface
}

使用时需要将它加入到你自定义的 transformer 中，像这样：

DemoTransformer struct {
    Transformer
}
...
```

如果你需要根据不同的条件， 比如在不同的路由下返回不同的数据结构， 我觉得比较好的做法是：
将判断条件添加在 DemoTransformer 中, 在需要的地方传入

```go

DemoTransformer struct {
    Transformer
		RouteName string
}
...

trans := &DemoTransformer{}
trans.RouteName = "xxxxx"
...

DemoResource.SetTransformer(trans)

```

### 4.4 嵌套包含 `Include`

如果您要包含相关数据，例如用户的帖子：

4.4.1 扩展您的转换器以处理包含：

```go
UserTransformer :
package transformers

type UserTransformer struct {
	fractal.Transformer
}

//建议使用初始化的方式获取transformer ,以避免忘记设置SetAvailableIncludeFunctions
func GetUserTransformer() *UserTransformer {
	trans := &UserTransformer{}
	trans.SetAvailableIncludeFunctions() // 如果想加载嵌套的包含数据， 一定要设置
	return trans
}

//注册包含的嵌套数据
func (ut *UserTransformer) SetAvailableIncludeFunctions() {
	ut.AvailableIncludeFunctions = map[string]func() fractal.ResourceInterface{
		"posts": ut.includePosts,
	}
}

func (ut *UserTransformer) includePosts() fractal.ResourceInterface {
	// 获取并返回相关帖子。这是一个简化的示例。
	posts := []Post{{Title: "帖子 1"}, {Title: "帖子 2"}}
	resource := &fractal.CollectionResource{}
	resource.SetSliceData(posts)
	resource.SetTransformer(&PostTransformer{}) // 假设您有一个 PostTransformer
	return resource
}
```

4.4.2 在转换数据时加载包含关系

在使用时可以加载各种关联关系，只需要向 resource 中设置就可以

```go
resource.SetIncludes([]string{"posts"})
//[{"id":1,"name":"John Doe", "posts":[{"title":"帖子 1"}，{"title":"帖子 2"}]},{"id":2,"name":"Jane Smith","posts":[{"title":"帖子 1"}，{"title":"帖子 2"}]}]
```

4.4.3 嵌套关系可以嵌套以上定义过的资源

```go
resource.SetIncludes([]string{"is_active"})
...
UserTransformer :

//嵌套 PrimitiveResource
func (ut *UserTransformer) SetAvailableIncludeFunctions() {
...
	"is_active" :ut.includeIsActive,
}
func (ut *UserTransformer) includeIsActive() fractal.ResourceInterface {
	activeStatus := &fractal.PrimitiveResource{}
	activeStatus.SetData(true)
	return activeStatus
}

//嵌套 ItemResource

func (ut *UserTransformer) SetAvailableIncludeFunctions() {
...
	"logo" :ut.includeLogo,
}

func (ut *UserTransformer) includeLogo() fractal.ResourceInterface {
	// 获取并返回相关帖子。这是一个简化的示例。
	logo := Logo{Url:"xxxx"}
	resource := &fractal.ItemResource{}
	resource.SetData(logo)
	resource.SetTransformer(&LogoTransformer{}) // 假设您有一个 PostTransformer
	return resource
}

//嵌套 CollectionResource
func (ut *UserTransformer) SetAvailableIncludeFunctions() {
...
	"books" :ut.includeBooks,
}

func (ut *UserTransformer) includeBooks() fractal.ResourceInterface {

	//如果存在返回logo
    users := []User{
        {ID: 1, Name: "John Doe", },
        {ID: 2, Name: "Jane Smith"},
    }
    resource := &fractal.ItemResource{}
    resource.SetData(logo)
    resource.SetTransformer(&BookTransformer{}) // 假设您有一个 PostTransformer
    return resource
}
```

## 五、 参考

本项目主要参考了 [league/fractal](https://fractal.thephpleague.com/transformers/) 适合 `php` 转 `go` 的小伙伴

```

```
