package helper

import (
	"GoChatServer/consts"
	"GoChatServer/dal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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

// GenerateStaticUrl 生成可访问的全路径url
func GenerateStaticUrl(filename string) string {
	if len(filename) == 0 {
		return ""
	}
	return Configs.Server.Host + path.Join(Configs.Server.StaticFileServerPath, filename)
}

// FormatFileMessageContent 对消息的内容格式化后返回，主要针对文件类的json中的附件信息添加域名前缀
func FormatFileMessageContent(t int32, content string) string {
	if t == consts.MessageTypeText {
		return content
	}
	fileData := types.MessageFileInfo{}
	err := json.Unmarshal([]byte(content), &fileData)
	if err != nil {
		Logger.Error("formatFileMessageContentUnmarshal:" + err.Error())
		return content
	}
	fileData.Filepath = GenerateStaticUrl(fileData.Filepath)

	fileJson, err := json.Marshal(&fileData)
	if err != nil {
		Logger.Error("formatFileMessageContentMarshal:" + err.Error())
		return content
	}
	return string(fileJson)
}

// UploadFileCheck 检测文件是否允许上传
func UploadFileCheck(file *multipart.FileHeader) error {
	//文件大小检测
	maxSizeMb := Configs.Server.MaxUploadFileSizeMb
	maxSizeByte := maxSizeMb * 1024 * 1024
	if file.Size > maxSizeByte {
		return fmt.Errorf("文件大小超过上限值：%d Mb", maxSizeMb)
	}

	//文件后缀格式检测
	ext := strings.ToLower(strings.Trim(filepath.Ext(file.Filename), "."))
	extIsOk := false
	for _, v := range Configs.Server.AllowUploadExtensions {
		if v == ext {
			extIsOk = true
		}
	}
	if !extIsOk {
		return fmt.Errorf("文件不允许上传，仅限上传的格式：%s", strings.Join(Configs.Server.AllowUploadExtensions, "、"))
	}
	return nil
}

// UploadFile 文件上传到本地服务器
func UploadFile(c *gin.Context, file *multipart.FileHeader, subject string) (filepath string, err error) {
	dateYmd := time.Now().Local().Format(consts.DateYMD)
	uuider := uuid.NewV4()
	filepath = path.Join(subject, dateYmd, uuider.String()+path.Ext(file.Filename))
	dst := path.Join(Configs.Server.UploadFilePath, filepath)

	// 上传文件至指定的完整文件路径
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		return "", err
	}
	return filepath, nil
}
