package projectgen

import (
	"fmt"
	"github.com/houyanzu/goforge/constdef"
	"os"
)

var dirs = []string{
	"/app/api/home/",
	"/bin/",
	"/databases/",
	"/lib/mydef/",
	"/lib/outputmsg/",
	"/myconfig/",
}

var files = []string{
	"/app/api/home/main.go",
	"/lib/mydef/mydef.go",
	"/lib/outputmsg/outputmsg.go",
	"/myconfig/myconfig.go",
	"/.gitignore",
	"/go.mod",
	"/build.ps1",
	"/build.sh",
}

func InitProject(name string) {
	for _, dir := range dirs {
		dir = "./" + name + dir
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating dir:", err, dir)
			return
		}
	}

	for k, f := range files {
		f = "./" + name + f
		file, err := os.Create(f)
		if err != nil {
			fmt.Println("Error creating file:", err, f)
			return
		}
		defer file.Close()

		switch k {
		case 0:
			_, err = file.WriteString(fmt.Sprintf(constdef.ApiMainContent, name, name))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 1:
			_, err = file.WriteString(constdef.MyDefContent)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 2:
			_, err = file.WriteString(constdef.OutputMsgContent)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 3:
			_, err = file.WriteString(constdef.MyConfigContent)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 4:
			_, err = file.WriteString(constdef.GitIgnoreContent)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 5:
			_, err = file.WriteString(fmt.Sprintf(constdef.GoModContent, name))
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 6:
			_, err = file.WriteString(constdef.Buildps1Content)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		case 7:
			_, err = file.WriteString(constdef.BuildshContent)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}
}
