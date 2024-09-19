package toolfunc

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// GetImportPkg 传入文件的路径，返回import可用的包名
func GetImportPkg(module, filePath string) string {
	// 将反斜杠替换为正斜杠
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	// 拼接 MODULE 和文件路径
	updatedPath := module + "/" + filePath

	// 找到最后一个 "/" 的位置
	lastSlashIndex := strings.LastIndex(updatedPath, "/")
	if lastSlashIndex != -1 {
		// 去掉最后一个 "/" 以及其后面的部分
		updatedPath = updatedPath[:lastSlashIndex]
	}

	return updatedPath
}

// GetModuleName 从 go.mod 文件中读取模块名称
func GetModuleName() (string, error) {
	// 打开 go.mod 文件
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	// 创建读取器
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read from go.mod")
	}

	// 读取第一行
	line := scanner.Text()
	if !strings.HasPrefix(line, "module ") {
		return "", fmt.Errorf("first line does not start with 'module '")
	}

	// 返回模块名称
	moduleName := strings.TrimPrefix(line, "module ")
	return moduleName, nil
}

// 检查结构体名是否以 "Controller" 结尾
func IsControllerType(name string) bool {
	return strings.HasSuffix(name, "Controller")
}

// GetFilenameByRoute 传入路由，返回文件路径和文件名
func GetFilenameByRoute(route string) (filePath, fileName string, err error) {
	if strings.HasPrefix(route, "/") || strings.HasSuffix(route, "/") {
		err = fmt.Errorf("controller route must not start or end with /")
		return
	}
	// 找到最后一个斜杠的位置
	lastSlashIndex := strings.LastIndex(route, "/")
	if lastSlashIndex == -1 {
		// 如果没有找到斜杠，则直接返回带 .go 扩展名的字符串
		return fmt.Sprintf("controller/%s.go", route), route, nil
	}

	// 将字符串分为两部分
	part1 := route[:lastSlashIndex]
	part2 := route[lastSlashIndex+1:]

	// 合成新的字符串
	modifiedPath := fmt.Sprintf("%s/controller/%s.go", part1, part2)

	return modifiedPath, part2, nil
}

// FileExists 检查文件是否存在
func FileExists(filename string) bool {
	// 获取文件信息
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// ValidateFileName 检查文件名是否合法
func ValidateFileName(s string) error {
	if strings.Contains(s, " ") {
		return errors.New("controller name can not contain spaces")
	}
	if strings.Contains(s, "_") {
		return errors.New("controller name can not contain _")
	}
	if strings.Contains(s, "-") {
		return errors.New("controller name can not contain -")
	}
	return nil
}

// CapitalizeFirstLetter 首字母大写
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s // 如果字符串为空，直接返回
	}

	// 将首字母大写，剩余部分保持不变
	return strings.ToUpper(string(s[0])) + s[1:]
}
