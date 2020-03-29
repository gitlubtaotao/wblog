package qiniu

import (
	"bytes"
	"context"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/gitlubtaotao/wblog/tools/upload"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/qiniu/x/rpc.v7"
	"mime/multipart"
	"time"
)

type IQiuNiuUploader interface {
	FilerUpload(file multipart.File, fileHeader *multipart.FileHeader) (url string, key string, err error)
	ByteUpload(bytes []byte, fileName string) (url string, key string, err error)
	PublicReadUrl(key string) string
	PrivateReadURL(key string) string
	FileInfo(key string) (storage.FileInfo, error)
	ChangeMimeType(key, mimeType string) error
	DeleteFile(key string) error
}

// 自定义返回值结构体
type MyPutRet struct {
	Key    string
	Hash   string
	Fsize  int
	Bucket string
	Name   string
}

//
type DirectUploader struct {
	Bucket        string
	FileServer    string
	Zone          *storage.Zone
	UseHTTPS      bool
	UseCdnDomains bool
	Expires       uint32
}

/*default config
default config
*/
func NewUploaderDefault() DirectUploader {
	zone := &storage.ZoneHuanan
	config := system.GetConfiguration()
	return NewUploader(config.QiniuBucket,
		config.QiniuFileServer, zone, false,
		false, 7200)
}

func NewUploader(bucket string, fileServer string, zone *storage.Zone, useHttps bool,
	useCdnDomains bool, expires uint32) DirectUploader {
	return DirectUploader{
		Bucket:        bucket,
		FileServer:    fileServer,
		Zone:          zone,
		UseHTTPS:      useHttps,
		UseCdnDomains: useCdnDomains,
		Expires:       expires,
	}
}

func (d *DirectUploader) Upload(file multipart.File, fileHeader *multipart.FileHeader) (string, string, error) {
	return d.FilerUpload(file, fileHeader)
}

/*
	FilerUpload: 文件直传方式
*/
func (d *DirectUploader) FilerUpload(file multipart.File, fileHeader *multipart.FileHeader) (url string, key string, err error) {
	var (
		ret  MyPutRet
		size int64
	)
	if statInterface, ok := file.(upload.Stat); ok {
		fileInfo, _ := statInterface.Stat()
		size = fileInfo.Size()
	}
	if sizeInterface, ok := file.(upload.Size); ok {
		size = sizeInterface.Size()
	}
	putPolicy := storage.PutPolicy{
		Scope:   d.Bucket,
		Expires: d.Expires,
	}
	mac := d.newMac()
	token := putPolicy.UploadToken(mac)
	cfg := d.storageConfig()
	uploader := storage.NewFormUploader(&cfg)
	putExtra := storage.PutExtra{}
	err = uploader.PutWithoutKey(context.Background(), &ret, token, file, size, &putExtra)
	if err != nil {
		return
	}
	
	return d.FileServer + ret.Key, ret.Key, nil
}

/*
上传bytes
*/
func (d *DirectUploader) ByteUpload(data []byte, fileName string) (url string, key string, err error) {
	var ret MyPutRet
	putPolicy := storage.PutPolicy{
		Scope: d.Bucket,
	}
	mac := d.newMac()
	upToken := putPolicy.UploadToken(mac)
	cfg := d.storageConfig()
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := storage.PutExtra{}
	dataLen := int64(len(data))
	err = formUploader.Put(context.Background(), &ret, upToken, fileName, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return
	}
	return d.FileServer + ret.Key, ret.Key, nil
}

func (d *DirectUploader) PublicReadUrl(key string) string {
	publicAccessURL := storage.MakePublicURL(d.FileServer, key)
	return publicAccessURL
}

/*
对于私有空间，首先需要按照公开空间的文件访问方式构建对应的公开空间访问链接，然后再对这个链接进行私有授权签名。
*/
func (d *DirectUploader) PrivateReadURL(key string) string {
	mac := d.newMac()
	deadline := time.Now().Add(time.Second * 3600).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, d.FileServer, key, deadline)
	return privateAccessURL
}

/*
获取文件信息
*/
func (d *DirectUploader) FileInfo(key string) (storage.FileInfo, error) {
	fileInfo, sErr := d.BucketManager().Stat(d.Bucket, key)
	return fileInfo, sErr
}

// ChangeMime
func (d *DirectUploader) ChangeMimeType(key string, mimeType string) error {
	return d.BucketManager().ChangeMime(d.Bucket, key, mimeType)
}

//删除文件
func (d *DirectUploader) DeleteFile(key string) error {
	return d.BucketManager().Delete(d.Bucket, key)
}

func (d *DirectUploader) BatchFileInfo(keys []string) []storage.BatchOpRet {
	statOps := make([]string, 0, len(keys))
	for _, key := range keys {
		statOps = append(statOps, storage.URIStat(d.Bucket, key))
	}
	rets, err := d.BucketManager().Batch(statOps)
	if err != nil {
		// 遇到错误
		if _, ok := err.(*rpc.ErrorInfo); ok {
			return rets
		} else {
			return nil
		}
	} else {
		return rets
	}
}

/*
 
 */
func (d *DirectUploader) BucketManager() *storage.BucketManager {
	mac := d.newMac()
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS:      d.UseHTTPS,
		Zone:          d.Zone,
		UseCdnDomains: d.UseCdnDomains,
	}
	return storage.NewBucketManager(mac, &cfg)
}

func (d *DirectUploader) newMac() (mac *qbox.Mac) {
	return qbox.NewMac(system.GetConfiguration().QiniuAccessKey, system.GetConfiguration().QiniuSecretKey)
}

func (d *DirectUploader) storageConfig() storage.Config {
	return storage.Config{
		Zone:          d.Zone,
		UseHTTPS:      d.UseHTTPS,
		UseCdnDomains: d.UseCdnDomains,
	}
}
