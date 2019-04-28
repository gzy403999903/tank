package util

import (
	"archive/zip"
	"fmt"
	"github.com/eyebluecn/tank/code/tool/result"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Zip(srcPath string, destPath string) {

	if PathExists(destPath) {
		panic(result.BadRequest("%s 已经存在了", destPath))
	}

	// 创建准备写入的文件
	fileWriter, err := os.Create(destPath)
	PanicError(err)
	defer func() {
		err := fileWriter.Close()
		PanicError(err)
	}()

	// 通过 fileWriter 来创建 zip.Write
	zipWriter := zip.NewWriter(fileWriter)
	defer func() {
		// 检测一下是否成功关闭
		if err := zipWriter.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zipWriter ，因为有可能会有很多个目录及文件，所以递归处理
	err = filepath.Walk(srcPath, func(path string, fileInfo os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		fmt.Println("遍历文件： " + path)

		// 通过文件信息，创建 zip 的文件信息
		fileHeader, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fileHeader.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// 目录加上/
		if fileInfo.IsDir() {
			fileHeader.Name += "/"
		}

		fmt.Println("头部情况： " + fileHeader.Name)

		// 写入文件信息，并返回一个 Write 结构
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 writer
		// 如目录，也没有数据需要写
		if !fileHeader.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fileToBeZip, err := os.Open(path)
		defer func() {
			err = fileToBeZip.Close()
			PanicError(err)
		}()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 writer
		_, err = io.Copy(writer, fileToBeZip)
		if err != nil {
			return
		}

		// 输出压缩的内容
		//fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
	PanicError(err)
}