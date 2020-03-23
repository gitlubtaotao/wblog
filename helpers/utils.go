package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"github.com/gitlubtaotao/wblog/system"
	"io"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"
	
	"github.com/pkg/errors"
	"github.com/snluu/uuid"
)

// 计算字符串的md5值
func Md5(source string) string {
	md5h := md5.New()
	md5h.Write([]byte(source))
	return hex.EncodeToString(md5h.Sum(nil))
}

func Truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) > n {
		return string(runes[:n])
	}
	return s
}

func UUID() string {
	return uuid.Rand().Hex()
}

func GetCurrentTime() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(loc)
}

func SendToMail(to, subject, body, mailtype string) error {
	config := system.GetConfiguration()
	user := config.SmtpUsername
	password := config.SmtpPassword
	host := config.SmtpHost
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var contentType string
	if mailtype == "html" {
		contentType = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	return smtp.SendMail(host, auth, user, sendTo, msg)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Decrypt(ciphertext []byte, keystring string) ([]byte, error) {
	// Key
	key := []byte(keystring)
	
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		err = errors.New("Text is too short")
		return nil, nil
	}
	
	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]
	
	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]
	
	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)
	
	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)
	
	return ciphertext, nil
}

func Encrypt(plaintext []byte, keystring string) ([]byte, error) {
	
	// Key
	key := []byte(keystring)
	
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	
	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	
	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]
	
	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	
	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)
	
	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	
	return ciphertext, nil
}

//正则匹配邮箱
func MatchEmail(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

//正则匹配手机电话号码
func MatchTelephone(phone string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phone)
}
