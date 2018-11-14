# CloudGo-IO

[CloudGo-IO](https://github.com/siskonemilia/CloudGo-IO) 是一个在 [CloudGo](https://github.com/siskonemilia/CloudGo) 的基础上开发的基于 Gin 框架开发的简单的 WEB 服务程序，支持静态文件服务、JavaScript 请求响应、模板输出、表单处理等功能。具体来讲，我们实现了一个拥有美观界面的支持用户注册和信息查看的网页。用户可以在其上登录他们的信息，然后通过相应的页面查看（表格太丑，做了些样式）。支持多个用户信息的存储和访问。支持用户查重、信息查重、信息格式查错。


![Preview](assets/preview.png)

## 使用 CloudGo-IO

配置好 Golang 环境的前提下，运行以下命令安装并使用 CloudGo-IO：

```bash
go install github.com/siskonemilia/CloudGo-IO
CloudGo-IO [-h/--hostname hostname] [-p/--port port]
```

然后你就可以访问 http://hostname:port，比如默认的 [http://localhost:8000](http://localhost:8000)，来使用功能了。一开始你会看到的是一个注册界面，在注册完成后你会自动跳转到信息页面。在知道其他用户的用户名的前提下，你也可以通过 http://hostname:port/detail?username=\[username\] 来查看其他用户的信息。

## Gin 框架基本用法

首先引用 Gin，然后我们就可以开始我们的开发之路了：

```golang
import(
  "github.com/gin-gonic/gin"
)
```

Gin 框架非常类似于我之前在 Node.js 开发时使用过的服务端框架：Express.js。他们都是实例化一个路由器（Router），然后对这个 Router 添加各种 Route 规则（也就是把特定路径的请求导向特定的处理函数）。如下例：

```golang
router := gin.Default()
router.GET("/getPath", func(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{
    "key": "value",
  })
})
router.Run("127.0.0.1:8000")
```

该例中，我们实例化了一个 Gin 的默认路由器，并且指定对于 /getPath 路径发起的 GET 方法将会由后面的处理函数接管。这个处理函数接受一个 *gin.Context 对象（其中包含了 Request 和 Response），并且完成所有的输入和输出的处理。在这里，它没有理会客户端发来的内容，而是简单地发挥了一个状态码为 200 的 JSON 响应。在最后，我们将服务器设置在了 127.0.0.1:8000 上，运行之后，客户端就可以访问 [http://127.0.0.1:8000/getPath](http://127.0.0.1:8000/getPath) 来查看结果了。

类似于上面的例子，Gin 还支持 POST、DELETE 等所有的 HTTP 请求，提供了多种预制的返回格式（HTML、XML等），使用起来非常方便，并且性能比语法相似的马天尼框架好大约 10 倍。

## 静态文件服务

在 Gin 中，静态文件服务的搭建是简单的。我们只需要为路由器指定一个 Static Route 即可：

```golang
router.Static("/public", "public")
// Now all visit towards /public will be served with
// the target file in ./public.
```

## JavaScript 请求响应 & 表单处理

以 POST 方法的 Ajax 为例。当客户端试图通过 Ajax 的 POST 方法向服务端发起请求时，它实际上是向服务端发送了一个 POST 方法的表单数据。因而，我们实际上是处理一个表单数据，gin 也为我们提供了完备的支持：

```golang
// This type is used to bind data sent to
// the server side.
type registerJSON struct {
	Username string `form:"username" binding:"required"`
	Stuid    string `form:"stuid" binding:"required"`
	Tel      string `form:"tel" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

// Create an instance of that type to bind data
var user registerJSON
// Try to bind the data sent from client side
if err := c.Bind(&user); err != nil {
  c.JSON(http.StatusBadRequest, gin.H{
    "code":    "BAD_REQUEST",
    "message": "Something wrong with the server",
  })
  fmt.Println(err.Error())
  return
}
```

在这段代码中，我们首先定义了一个用于绑定数据的类型，并且为其指定了字段。所谓的「绑定数据」指的是我们将会把 request.body 中的数据按字段填入到这个类型的实例中。在上面的示例代码中，我们为其分配了表单中的字段，从而使它可以去解析一个 POST 表单。进而进行进一步的逻辑处理。通过 Bind 方法，我们就能将请求中的数据悉数传入到给定的绑定实例上，而如果有任何问题发生，一个 err 将会被返回给你。

## 模板输出

Gin 支持众多的模板，但是简单起见，我们这里就使用最简单的 HTML 模板来做示例。那么，我们首先说说如何创建一个 HTML 模板吧。这种模板和普通的 HTML 区别不大，我们这里需要使用到的就是它的「变量」特性。在 HTML 模板中，所有的变量被标识为 `{.VariableName}` 的格式，这些被称为变量的字符串可以在渲染阶段被服务端轻易而统一地替换为有意义的字段。最后，我们把 `.html` 的后缀改为 `.tmpl` 即可。

为了加载一个 HTML 模板，我们需要使用 `LoadHTMLGlob` 或是 `LoadHTMLFiles` 方法：

```golang
router.LoadHTMLBlob("views/*")
// Load all files in folder views
```

之后，我们就可以在需要的地方这样子使用模板渲染一个 HTML 返回给用户了：

```golang
c.HTML(http.StatusOK, "ATemplate.tmpl", gin.H{
  "variable1": value1,
  "variable2": value2,
})
```

这样，用户收到的将会是一个变量被替换为对应值的 HTML。

## 「未实现」与 404

对于 `/unknown` 的「未实现」报错，我们也可以简单地通过最基本的语法实现：

```golang
router.GET("/unknown", func(c *gin.Context) {
  c.JSON(http.StatusNotImplemented, gin.H{
    "code":    "NOT_IMPLEMENTED",
    "message": "This page has not been implemented.",
  })
})
```

而对于 404 问题，Gin 提供了非常贴心的 NoRoute 方法，他可以捕获所有未能通过 Route 方法指定处理函数的请求，并且分配处理函数：

```golang
// PAGE_NOT_FOUND page for all paths without routing
router.NoRoute(func(c *gin.Context) {
  c.JSON(http.StatusNotFound, gin.H{
    "code":    "PAGE_NOT_FOUND",
    "message": "Target page not found.",
  })
})
```

这样用户就知道「服务器没有挂掉，只是我的网址写错了」啦。