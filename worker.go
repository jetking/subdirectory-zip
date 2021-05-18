package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"io"
	"io/fs"
	"os"
	"path"
	"time"
)

func do(source, desc string) {
	t1 := time.Now()
	dirFS := os.DirFS(source)
	entries, err := fs.ReadDir(dirFS, ".")
	if err != nil {
		color.Red("读取目录[%s]失败:%s", source, err.Error())
		return
	}
	dirTotal := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirTotal++
	}
	color.White("待处理目录总数:%d", dirTotal)

	idx := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		idx++
		//fmt.Println("处理小区:", entry.Name())
		files, err := fs.ReadDir(dirFS, entry.Name())
		if err != nil {
			color.Red("读取目录[%s]失败:%s", entry.Name(), err.Error())
			continue
		}
		entriesDone := make(chan struct{})
		progress := make(chan struct{})
		bar := progressbar.NewOptions(len(files),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionSetWidth(40),
			progressbar.OptionSetDescription(fmt.Sprintf("[%d/%d]%s\t\t", idx, dirTotal, entry.Name())),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
			progressbar.OptionOnCompletion(func() {
				entriesDone <- struct{}{}
			}),
		)
		go func() { // 更新进度条
			for {
				select {
				case <-time.NewTicker(10 * time.Second).C: // 10秒钟还干不完一个文件?
					color.Red("处理文件超时退出")
					return
				case <-progress:
					bar.Add(1)
				}
			}
		}()

		// 创建 zip 文件
		zipFile, err := os.Create(path.Join(desc, entry.Name()+".zip"))
		if err != nil {
			fmt.Println("创建压缩包文件失败:", err)
			continue
		}
		defer zipFile.Close()
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()

		// 把文件放入 zip
		for _, f := range files {
			dst, err := zipWriter.Create(f.Name())
			if err != nil {
				color.Red("在压缩包内创建文件[%s]失败:%s", f.Name(), err.Error())
				continue
			}
			src, err := fs.ReadFile(dirFS, path.Join(entry.Name(), f.Name()))
			if err != nil {
				color.Red("读取原始文件[%s]失败:%s", f.Name(), err.Error())
				continue
			}
			if _, err := io.Copy(dst, bytes.NewReader(src)); err != nil {
				color.Red("复制文件[%s]失败:%s", f.Name(), err.Error())
			}
			progress <- struct{}{}
		}
		<-entriesDone
		fmt.Println("")
	}
	t2 := time.Now()
	color.Green("\n操作完成 ^_^ \n文件保存在:%s\n用时:%v", desc, t2.Sub(t1))
	return
}
