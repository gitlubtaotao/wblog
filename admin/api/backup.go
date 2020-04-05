package admin

import (
	"fmt"
	"github.com/gitlubtaotao/wblog/api"
	"github.com/gitlubtaotao/wblog/encrypt"
	"github.com/gitlubtaotao/wblog/tools/upload/qiniu"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"github.com/gitlubtaotao/wblog/helpers"
	"github.com/gitlubtaotao/wblog/system"
	"github.com/pkg/errors"
)

func BackupPost(c *gin.Context) {
	var (
		err error
		res = gin.H{}
	)
	defer api.WriteJSON(c, res)
	err = Backup()
	if err != nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

func RestorePost(c *gin.Context) {
	var (
		fileName  string
		fileUrl   string
		err       error
		res       = gin.H{}
		resp      *http.Response
		bodyBytes []byte
	)
	defer api.WriteJSON(c, res)
	fileName = c.PostForm("fileName")
	if fileName == "" {
		res["message"] = "fileName cannot be empty."
		return
	}
	upload := qiniu.NewUploaderDefault()
	fileUrl = upload.PublicReadUrl(fileName)
	resp, err = http.Get(fileUrl)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	defer resp.Body.Close()
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	
	bodyString, err := encrypt.DeCryptData(string(bodyBytes), false, "admin")
	bodyBytes = []byte(bodyString)
	if err != nil {
		res["message"] = err.Error()
		return
	}
	err = ioutil.WriteFile(fileName, bodyBytes, os.ModePerm)
	if err == nil {
		res["message"] = err.Error()
		return
	}
	res["succeed"] = true
}

//备份
func Backup() (err error) {
	var (
		u           *url.URL
		exist       bool
		bodyBytes   []byte
		encryptData []byte
	)
	u, err = url.Parse(system.GetConfiguration().DSN)
	if err != nil {
		seelog.Debug("parse dsn error:%v", err)
		return
	}
	exist, _ = helpers.PathExists(u.Path)
	if !exist {
		err = errors.New("database file doesn't exists.")
		seelog.Debug("database file doesn't exists.")
		return
	}
	seelog.Debug("start backup...")
	bodyBytes, err = ioutil.ReadFile(u.Path)
	if err != nil {
		seelog.Error(err)
		return
	}
	encryptData, err = encrypt.EncryptBytes(bodyBytes, "admin")
	if err != nil {
		seelog.Error(err)
		return
	}
	uploader := qiniu.NewUploaderDefault()
	_, _, err = uploader.ByteUpload(encryptData, fmt.Sprintf("wblog_%s.db", helpers.GetCurrentTime().Format("20060102150405")))
	if err != nil {
		seelog.Debugf("backup error:%v", err)
		return
	}
	seelog.Debug("backup succeefully.")
	return err
}
