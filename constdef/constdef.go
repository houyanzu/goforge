package constdef

const PackageMainStr = "package main\n\n"

const ControllerImportStr = `
import (
	"github.com/gin-gonic/gin"
)

`

const ControllerFuncStr = `
// %s
func (co %s) %s(c *gin.Context) {
	//TODO: edit
}
`

const ControllerLoginFuncStr = `
// %s
func (co %s) %s(c *gin.Context, userID uint) {
	//TODO: edit
}
`
const RouterImportStr = `import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"github.com/gin-gonic/gin"
	"github.com/houyanzu/work-box/tool/middleware"
`

const RouterHeaderStr = `// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.` + "\n\n"

const RouterRegistControllerStr = "func RegisterController(controller interface{}) {\n\tcontrollers = append(controllers, controller)\n}\n\n"

const RouterOtherStr = `
func Register(router *gin.Engine) {
	for _, v := range controllers {
		AutoRegisterRoutes(router, v)
	}
}

// 自动注册路由
func AutoRegisterRoutes(router *gin.Engine, controller interface{}) {
	controllerType := reflect.TypeOf(controller)
	controllerValue := reflect.ValueOf(controller)

	// 获取基础路由前缀
	baseRoute, pkgPath := buildBaseRoute(controllerType)

	// 遍历控制器的所有方法并注册路由
	for i := 0; i < controllerType.NumMethod(); i++ {
		method := controllerType.Method(i)
		methodName := getControllerRouterName(method.Name)
		methodType := method.Type
		numIn := methodType.NumIn()
		if numIn < 2 {
			continue
		}

		firstParamType := getParamTypeName(methodType.In(1))
		if firstParamType != "Context" {
			continue
		}

		key := pkgPath+"."+controllerType.Name()+"."+method.Name
		httpMethod := MethodTags[key]
		// 注册方法为 Gin 的 Post 路由
		route := fmt.Sprintf("api/%s/%s", baseRoute, methodName)
		switch numIn {
		case 2:
			switch httpMethod {
			case "POST":
				router.POST(route, func(ctx *gin.Context) {
					method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx)})
				})
			case "GET":
				router.GET(route, func(ctx *gin.Context) {
					method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx)})
				})
			default:
				router.POST(route, func(ctx *gin.Context) {
					method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx)})
				})
			}

		case 3:
			secondParamType := getParamTypeName(methodType.In(2))
			if secondParamType == "uint" {
				switch httpMethod {
				case "POST":
					router.POST(route, middleware.Login(), func(ctx *gin.Context) {
						userID := middleware.GetUserId(ctx)
						method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx), reflect.ValueOf(userID)})
					})
				case "GET":
					router.GET(route, middleware.Login(), func(ctx *gin.Context) {
						userID := middleware.GetUserId(ctx)
						method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx), reflect.ValueOf(userID)})
					})
				default:
					router.POST(route, middleware.Login(), func(ctx *gin.Context) {
						userID := middleware.GetUserId(ctx)
						method.Func.Call([]reflect.Value{controllerValue, reflect.ValueOf(ctx), reflect.ValueOf(userID)})
					})
				}

			}
		}

	}
}


// 构建控制器的基础路由前缀
func buildBaseRoute(controllerType reflect.Type) (string, string) {
	// 解析包路径
	pkgPath := controllerType.PkgPath()
	pkgParts := strings.Split(pkgPath, "/")

	// 构建路由路径
	var routeBuilder strings.Builder
	flag := false

	for _, part := range pkgParts {
		if part == "controller" {
			break
		}
		if flag {
			routeBuilder.WriteString(part + "/")
		}
		if part == "home" {
			flag = true
		}
	}

	// 添加控制器名
	controllerName := getControllerRouterName(controllerType.Name())
	routeBuilder.WriteString(controllerName)

	return routeBuilder.String(), pkgPath
}

// 去掉字符串末尾的 "Controller" 并将首字母变为小写
func getControllerRouterName(input string) string {
	input = strings.TrimSuffix(input, "Controller")
	if input == "" {
		return input
	}

	runes := []rune(input)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func getParamTypeName(paramType reflect.Type) string {
	if paramType.Kind() == reflect.Ptr {
		return paramType.Elem().Name()
	} else {
		return paramType.Name()
	}
}

`

const ApiMainContent = `package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/houyanzu/work-box/database"
	"github.com/houyanzu/work-box/lib/boxlog"
	"github.com/houyanzu/work-box/tool/cache"
	_ "%s/lib/outputmsg"
	"%s/myconfig"
)

func main() {
	var configName string
	var port string
	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&configName, "config", "", "Config json file name")
	flag.Parse()
	err := myconfig.ParseConfig(configName)
	if err != nil {
		panic(err)
	}

	err = database.InitMysql()
	if err != nil {
		panic(err)
	}
	err = cache.InitCache()
	if err != nil {
		panic(err)
	}

	err = boxlog.Init("home_api", true)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	Register(router)
	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}
}

`

const GitIgnoreContent = `go.sum
.idea
.vscode
*.json
*.log
bin
`

const MyDefContent = `// 在这里定义一些通用的自定义常量、变量、方法等

package mydef
`

const OutputMsgContent = `// 只修改msgMap，不要添加其他代码

package outputmsg

import (
	"github.com/houyanzu/work-box/lib/output"
)

func init() {
	output.InitMsgMap(msgMap)
}

var msgMap = map[output.ErrorCode]map[string]string{
	0: {
		"zh": "ok",
		"en": "ok",
		"tw": "ok",
	},
	1: {
		"zh": "请重试",
		"en": "Please try again",
		"tw": "請重試",
	},
	3: {
		"zh": "登陆过期了，需要重新登录哟",
		"en": "Login Has Expired. Please Login Again",
		"tw": "登錄過期了，需要重新登錄喲",
	},
	6: {
		"zh": "错误",
		"en": "Error",
		"tw": "Error",
	},
}
`

const MyConfigContent = `package myconfig

import (
	"bytes"
	"encoding/json"
	"github.com/houyanzu/work-box/config"
	"io/ioutil"
	"os"
)

type MyConfig struct {
	config.Config
	Common commonConfig ` + "`json:\"common\"`" + `
}

type commonConfig struct {
}

var myConfig *MyConfig

func ParseConfig(configName string) (err error) {
	dat, err := ioutil.ReadFile(configName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, &myConfig)
	if err != nil {
		return err
	}

	config.ParseConfig(&myConfig.Config)

	return nil
}

func GetConfig() *MyConfig {
	return myConfig
}

func CreateConfigFile() {
	var conf MyConfig
	js, err := json.Marshal(conf)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	_ = json.Indent(&out, js, "", "  ")
	f, _ := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE, 0777)
	out.WriteTo(f)
}
`

const GoModContent = `module %s

go 1.22.3

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/houyanzu/work-box v1.7.3
)
`
const Buildps1Content = `param(
    [Parameter(Mandatory=$false)]
    [string]$path,

    [Parameter(Mandatory=$false)]
    [string]$goos
)

if ([string]::IsNullOrEmpty($path)) {
    Write-Host "Error: The 'path' parameter is not provided or is an empty string." -ForegroundColor Red
    exit 1
}

if (-not [string]::IsNullOrEmpty($goos)) {
    $env:GOOS = $goos
}

goforge routergen
go build -o ./bin/ $path

if (-not [string]::IsNullOrEmpty($goos)) {
    $env:GOOS = 'windows'
}
`

const BuildshContent = `#!/bin/sh

path=""
goos=""

original_goos=$GOOS

while [ "$1" != "" ]; do
    case $1 in
        --path )          shift
                          path=$1
                          ;;
        --goos )          shift
                          goos=$1
                          ;;
        * )               echo "Invalid parameter: $1"
                          exit 1
    esac
    shift
done

if [ -z "$path" ]; then
    echo "Error: The 'path' parameter is not provided or is an empty string."
    exit 1
fi

if [ -n "$goos" ]; then
    export GOOS=$goos
fi

goforge routergen
if [ $? -ne 0 ]; then
    echo "Error: goforge routergen command failed."
    exit 1
fi

go build -o ./bin/ "$path"
if [ $? -ne 0 ]; then
    echo "Error: go build command failed."
    exit 1
fi

if [ -n "$goos" ]; then
    export GOOS=$original_goos
fi
`
