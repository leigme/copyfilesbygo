package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

// Config 配置文件结构体
type Config struct {
	Version string
	Source  string
	Target  string
}

func main() {
	config := &Config{}
	if analysis(config) {
		listFile(config.Source, config.Source, config.Target)
	}
}

// 解析配置文件到结构体中
func analysis(config *Config) bool {
	// 打开文件
	file, err := os.Open("config.json")

	if nil != err {
		log.Fatal(err)
		return false
	}
	// 延迟关闭文件
	defer file.Close()

	fileInfo, err := file.Stat()

	if nil != err {
		log.Fatal(err)
		return false
	}
	// 获取文件长度
	fileSize := fileInfo.Size()

	buffer := make([]byte, fileSize)
	// 读取文件内容到内存
	bytesRead, err := file.Read(buffer)

	if nil != err {
		log.Fatal(err)
		return false
	}

	fmt.Println(string(buffer[:bytesRead]))
	// 解析文件
	err = json.Unmarshal(buffer, config)

	if nil != err {
		log.Fatal(err)
		return false
	}

	return true
}

// 读取源目录结构
func listFile(root, source, target string) {
	// 获取目录下所有文件及目录
	files, err := ioutil.ReadDir(source)

	if nil != err {
		log.Fatal(err)
		return
	}

	for _, file := range files {

		fileName := source + "/" + file.Name()

		if file.IsDir() {
			// 递归出路子目录
			listFile(root, fileName, target)
		} else {
			// 创建目标目录结构
			targetFile := strings.Replace(fileName, root, target, 1)
			// 复制文件
			copyFile(targetFile, fileName)
		}
	}
}

func copyFile(dstName, srcName string) (written int64, err error) {
	// 读取源文件
	src, err := os.Open(srcName)

	if nil != err {
		log.Fatal(err)
		return
	}
	// 获取目标文件的相对目录
	tar := strings.Replace(dstName, path.Base(srcName), "", 1)

	fmt.Println("======>" + tar)

	if nil != err {
		log.Fatal(err)
		return
	}
	// 创建目录
	err = os.MkdirAll(tar, os.ModePerm)

	if nil != err {
		fmt.Printf("mkdir failed![%v]\n", err)
	} else {
		fmt.Printf("mkdir success!\n")
	}
	// 关闭源文件
	defer src.Close()
	// 创建目标文件
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)

	if nil != err {
		log.Fatal(err)
		return
	}
	// 延迟关闭目标文件
	defer dst.Close()
	// 复制文件
	return io.Copy(dst, src)
}
