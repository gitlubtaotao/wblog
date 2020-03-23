package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"github.com/gitlubtaotao/wblog/system"
)

func EnCryptData(originData string) (string, error) {
	secret := system.GetConfiguration().SessionSecret
	result, err := aesEncrypt([]byte(originData), []byte(secret))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), err
}

func DeCryptData(hash string) (string, error) {
	//解密base64字符串
	secret := system.GetConfiguration().SessionSecret
	pwdByte, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}
	//执行AES解密
	originData, err := aesDeCrypt(pwdByte, []byte(secret))
	return string(originData), err
}

//使用aes进行加密
func aesEncrypt(originData []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData := pKCS7Padding(originData, blockSize)
	//采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blocMode.CryptBlocks(crypted, origData)
	return crypted, nil
	
}

func pKCS7Padding(ciphered []byte, blockSize int) []byte {
	padding := blockSize - len(ciphered)%blockSize
	pretext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphered, pretext...)
}

func aesDeCrypt(canted []byte, key []byte) ([]byte, error) {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(canted))
	//这个函数也可以用来解密
	blockMode.CryptBlocks(origData, canted)
	//去除填充字符串
	origData, err = pKCS7UnPadding(origData)
	if err != nil {
		return nil, err
	}
	return origData, err
}

//填充的反向操作，删除填充字符串
func pKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}
