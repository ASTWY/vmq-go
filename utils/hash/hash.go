package hash

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// GetMD5Hash 返回字符串的MD5值
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetSha256Hash 返回字符串的Sha256值
func GetSha256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
