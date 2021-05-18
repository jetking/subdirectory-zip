package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func main() {
	color.Magenta("该工具用于将指定目录的各子目录分别打包成 zip 文件")
	source := ""
	dest := ""
	args := os.Args // ./exe [source] [dest]
	if len(args) != 3 {
		color.White("请输入来源目录位置:")
		fmt.Scanln(&source)
		color.White("请输入保存压缩包的目录位置:")
		fmt.Scanln(&dest)
	} else {
		source = args[1]
		dest = args[2]
	}
	fmt.Printf("来源目录:\t%s\n目标目录:\t%s\n", source, dest)
	do(source, dest)
}
