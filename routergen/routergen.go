package routergen

import (
	"fmt"
	"github.com/houyanzu/goforge/constdef"
	"github.com/houyanzu/goforge/toolfunc"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

var imports []string
var inits []string
var importMap map[string]string
var httpMethods map[string]string

func Routergen(root string) {
	var err error

	// 打开或创建 register.go 文件
	file, err := os.Create(root + "register.go")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(constdef.RouterHeaderStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// 写入 package 声明
	_, err = file.WriteString(constdef.PackageMainStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	imports = make([]string, 0)
	inits = make([]string, 0)
	importMap = make(map[string]string)
	httpMethods = make(map[string]string)

	err = scanDirectories(root)
	if err != nil {
		fmt.Println("Error scanning directories:", err)
		return
	}

	// 写入导入语句
	importStr2 := constdef.RouterImportStr
	for _, v := range imports {
		importStr2 += "\t" + v + "\n"
	}
	importStr2 += ")\n\n"
	_, err = file.WriteString(importStr2)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// 写入函数定义
	_, err = file.WriteString("var controllers []interface{}\n\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	httpMethodsStr := "var MethodTags = map[string]string{\n%s\n}\n\n"
	httpMethodContent := ""
	for k, v := range httpMethods {
		httpMethodContent += fmt.Sprintf("\t\"" + k + "\": \"" + v + "\",\n")
	}
	httpMethodsStr = fmt.Sprintf(httpMethodsStr, httpMethodContent)
	_, err = file.WriteString(httpMethodsStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	_, err = file.WriteString(constdef.RouterRegistControllerStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	ss := ""
	for _, v := range inits {
		// 写入函数定义
		ss += v
	}
	ss = "func init() {" + ss + "\n}\n"
	_, err = file.WriteString(ss)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	otherStr := constdef.RouterOtherStr
	_, err = file.WriteString(otherStr)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

// 递归遍历目录，处理 controller 目录中的 Go 文件
func scanDirectories(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == "controller" {
			// 遇到 controller 目录，处理其中的 Go 文件
			return filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
					return processGoFile(path)
				}
				return nil
			})
		}
		return nil
	})
}

// 处理 Go 文件中的控制器
func processGoFile(filePath string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return err
	}

	module, err := toolfunc.GetModuleName()
	if err != nil {
		return err
	}

	pak := toolfunc.GetImportPkg(module, filePath)
	alias, ok := importMap[pak]
	if !ok {
		alias = fmt.Sprintf("controller%d", len(imports))
	}

	have := false
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			if toolfunc.IsControllerType(typeSpec.Name.Name) {
				controllerName := typeSpec.Name.Name
				fullControllerName := alias + "." + controllerName
				key := pak + "." + controllerName
				inspectControllerMethods(node, controllerName, key)
				// Write out code to register the controller
				inits = append(inits, fmt.Sprintf("\n\tRegisterController(%s{})", fullControllerName))

				have = true
			}
		}
	}
	if have {
		if !ok {
			importMap[pak] = alias
			imports = append(imports, alias+" \""+pak+"\"")
		}
	}

	return nil
}

func inspectControllerMethods(node *ast.File, controllerName, key string) {
	ast.Inspect(node, func(n ast.Node) bool {
		// 只处理函数声明
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}

		// 确保这是控制器的成员方法
		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			// 获取接收者类型（即控制器名）
			recvType := getReceiverType(funcDecl.Recv.List[0].Type)
			if recvType == controllerName {
				// 打印方法名
				hmk := key + "." + funcDecl.Name.Name

				// 打印注解（注释）
				if funcDecl.Doc != nil {
					firstLine := funcDecl.Doc.List[0].Text
					firstLineArr := strings.Split(firstLine, " ")
					if len(firstLineArr) > 1 {
						httpMethods[hmk] = firstLineArr[1]
					}
				}
			}
		}

		return true
	})
}

// 获取接收者类型的名称
func getReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr: // Handle pointer receiver like *Controller
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}
