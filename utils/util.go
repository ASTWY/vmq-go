package utils

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"time"

	"github.com/shopspring/decimal"
)

// float64ToSting float64转string
func Float64ToSting(f float64) string {
	// 如果是整数，直接返回
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	// 如果不是整数，当小数点后面有两位时，直接返回保留两位小数的字符串
	if f*100 == float64(int64(f*100)) {
		return fmt.Sprintf("%.2f", f)
	}
	// 如果不是整数，当小数点后面有一位时，直接返回保留一位小数的字符串
	if f*10 == float64(int64(f*10)) {
		return fmt.Sprintf("%.1f", f)
	}
	// 默认返回保留两位小数的字符串
	return fmt.Sprintf("%.2f", f)
}

// GetUnix13 返回13位时间戳 单位毫秒
func GetUnix13() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetUnix10 返回10位时间戳 单位秒
func GetUnix10() int64 {
	return time.Now().Unix()
}

// GetCWDir 返回当前工作目录
func GetCWDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

// DictionaryOrderSort 字典序排序 按照key的ASCII码从小到大排序（字典序）
func DictionaryOrderSort(m map[string]string) ([]string, map[string]string) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := make(map[string]string, len(m))
	for _, k := range keys {
		result[k] = m[k]
	}
	return keys, result
}

// SendEmailUsingSMTP 使用SMTP发送邮件
func SendEmailUsingSMTP(subject string, body string, smtpHost string, smtpPort int, smtpUser string, smtpPassword string, smtpSender string, to string, ssl bool) error {
	// 设置身份验证信息。
	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPassword,
		smtpHost,
	)
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	var err error
	if ssl {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpHost,
		}

		conn, errConn := tls.Dial("tcp", addr, tlsconfig)
		if errConn != nil {
			return errConn
		}

		c, errC := smtp.NewClient(conn, smtpHost)
		if errC != nil {
			return errC
		}
		err = c.Auth(auth)
		if err != nil {
			return err
		}
		err = c.Mail(smtpSender)
		if err != nil {
			return err
		}
		err = c.Rcpt(to)
		if err != nil {
			return err
		}

		wc, errW := c.Data()
		if errW != nil {
			return errW
		}
		defer wc.Close()
		_, err = fmt.Fprintf(wc, "Subject: "+subject+"\r\n\r\n"+body+"\r\n")
		if err != nil {
			return err
		}
		c.Quit()
	} else {
		err = smtp.SendMail(
			addr,
			auth,
			smtpSender,
			[]string{to},
			[]byte("Subject: "+subject+"\r\n\r\n"+body+"\r\n"),
		)

		if err != nil {
			return fmt.Errorf("发送邮件失败: %w", err)
		}
	}

	return nil
}

// 生成随机字符串 // 使用crypto/rand包
func RandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.Reader
	for i := 0; i < l; i++ {
		b := make([]byte, 1)
		r.Read(b)
		result = append(result, bytes[int(b[0])%len(bytes)])
	}
	return string(result)
}

// 生成订单号
func GenerateOrderNo() string {
	return fmt.Sprintf("%d%s", GetUnix13(), RandomString(7))
}

// 检查是否连续
func CeckConsecutive(arr []float64) (bool, []float64) {
	// sort.Float64s(arr)
	slices.Sort(arr)

	// 将arr 转为 []decimal
	var arrDecimal []decimal.Decimal
	for _, v := range arr {
		arrDecimal = append(arrDecimal, decimal.NewFromFloat(v))
	}

	var missing []float64
	switch len(arr) {
	case 0:
		return true, missing
	case 1:
		return true, missing
	}

	isConsecutive := true

	for i := 1; i < len(arr); i++ {
		diff := arrDecimal[i].Sub(arrDecimal[i-1]) // 相邻两个数的差
		// if res != 0.01 {
		if diff.Cmp(decimal.NewFromFloat(0.01)) != 0 {
			isConsecutive = false
			if diff.Cmp(decimal.NewFromFloat(0.01)) > 0 {
				for j := arr[i-1] + 0.01; j < arr[i]; j += 0.01 {
					missing = append(missing, j)
				}
			}
		}
	}
	slices.Sort(missing)
	return isConsecutive, missing
}

// 二分查找
func BinarySearch(arr []float64, target float64) bool {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) / 2
		if arr[mid] == target {
			return true
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return false
}

// 下载文件
func DownloadFile(url string, filePath string) error {
	// 创建目录
	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}

	// 创建一个 HTTP GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个新的文件保存下载的内容
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将下载的内容保存到文件中
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
