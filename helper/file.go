package helper

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// NewFile 新建一个文件
func NewFile(filename string) bool {
	//检查文件是否存在
	if _, err := os.Stat(filename); errors.Is(err, fs.ErrNotExist) {
		//如果文件不存在，则创建文件所在的目录
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			fmt.Println("创建目录失败：", err)
			return false
		}

		//创建文件
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("创建文件失败：", err)
			return false
		}
		defer func() {
			_ = file.Close()
		}()
		return true
	} else if err != nil {
		fmt.Println("检查文件失败：", err)
		return false
	} else {
		return true
	}
}
