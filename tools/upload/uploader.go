package upload

import (
	"mime/multipart"
	"os"
)

//上传文件的方法
type Uploader interface {
	Upload(file multipart.File, fileHeader *multipart.FileHeader) (string, string, error)
}

// 获取文件大小的接口
type Size interface {
	Size() int64
}

// 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}
