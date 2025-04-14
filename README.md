# GoForge 命令行工具介绍

`GoForge` 是一个用于快速构建 Go 项目的命令行工具，它可以帮助你轻松地生成控制器、方法、路由文件等。下面是如何使用 GoForge 工具的基本操作。

## 安装 GoForge

使用以下命令安装 GoForge 命令行工具：

```bash
go install github.com/houyanzu/goforge@latest
```

## 创建项目

初始化一个名为 `demo` 的新项目：

```bash
goforge init demo
```

进入项目目录：

```bash
cd demo
```

## 添加控制器

### 创建控制器

添加一个名为 `user/account` 的控制器：

```bash
goforge addController user/account
```

### 创建控制器并添加方法

创建一个控制器并同时为其添加多个方法：

```bash
goforge addController user/account --methods 'register login setPassword:login:POST getName:GET'
```

- 该命令创建 `user/account` 控制器，并添加四个方法：
    - `register`
    - `login`
    - `setPassword`（需要登录，且请求方法为 POST）
    - `getName`（请求方法为 GET）

如果不指定请求方法，默认为 POST。

### 为已存在的控制器添加方法

如果控制器已经存在，可以使用以下命令为其添加方法：

```bash
goforge addMethods user/account --methods 'register login setPassword:login:POST getName:GET'
```

## 生成路由文件

生成路由文件，自动根据已创建的控制器和方法生成路由：

```bash
goforge routergen
```

## 构建 API 可执行文件

生成指定 API 的可执行文件，以下命令将生成 `home` API 的可执行文件，并将其放置在 `bin` 目录下：

```bash
.\build.ps1 .\app\api\home\
```

此命令还会自动生成路由文件。

---

通过 GoForge，你可以轻松管理你的 Go 项目中的控制器、方法和路由文件，提高开发效率。