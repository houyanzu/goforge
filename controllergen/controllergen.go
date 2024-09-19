package controllergen

import (
	"fmt"
	"goforge/constdef"
	"goforge/toolfunc"
	"os"
	"strings"
)

type method struct {
	Name       string
	Login      bool
	HTTPMethod string
}

func OperateController(root, route, action, methods string) {
	if route == "" {
		fmt.Println("please enter controller route:")
		_, err := fmt.Scanf("%s", &route)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if route == "" {
		fmt.Println("Error: controller route is empty")
		return
	}

	filePath, fileName, err := toolfunc.GetFilenameByRoute(route)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = toolfunc.ValidateFileName(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	controllerStructName := toolfunc.CapitalizeFirstLetter(fileName) + "Controller"

	filePath = root + filePath
	switch action {
	case "addController":
		if toolfunc.FileExists(filePath) {
			fmt.Println("Error: controller route already exists")
			return
		}
		// 获取文件所在的目录
		dir := filePath[:len(filePath)-len(fileName+".go")]

		// 创建目录（如果不存在的话）
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating dir:", err)
			return
		}
	case "addMethods":
		if !toolfunc.FileExists(filePath) {
			fmt.Println("Error: controller route does not exist.")
			return
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if action == "addController" {
		// 写入 package 声明
		_, err = file.WriteString("package controller\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		_, err = file.WriteString(constdef.ControllerImportStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		structStr := "type " + controllerStructName + " struct {}\n\n"
		_, err = file.WriteString(structStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

	}

	methodList, err := getMethods(methods)
	if err != nil {
		fmt.Println("Error getMethods:", err)
		return
	}
	for _, v := range methodList {
		methodStr := fmt.Sprintf(constdef.ControllerFuncStr, v.HTTPMethod, controllerStructName, v.Name)
		if v.Login {
			methodStr = fmt.Sprintf(constdef.ControllerLoginFuncStr, v.HTTPMethod, controllerStructName, v.Name)
		}
		_, err = file.WriteString(methodStr)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func getMethods(methods string) (res []method, err error) {
	if methods == "" {
		return
	}

	methods = strings.Trim(methods, " ")
	methodsArr := strings.Split(methods, " ")
	length := len(methodsArr)
	res = make([]method, length)
	for k, mtd := range methodsArr {
		mtdDetail := strings.Split(mtd, ":")
		for i, v := range mtdDetail {
			if i == 0 {
				res[k].Name = toolfunc.CapitalizeFirstLetter(v)
			} else {
				switch strings.ToLower(v) {
				case "login":
					res[k].Login = true
				case "post":
					res[k].HTTPMethod = "POST"
				case "get":
					res[k].HTTPMethod = "GET"
				default:
					err = fmt.Errorf("invalid method: %s", v)
					return
				}
			}
			if res[k].HTTPMethod == "" {
				res[k].HTTPMethod = "POST"
			}
		}
	}
	return
}
