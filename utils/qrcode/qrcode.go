package qrcode

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	qrcode2 "github.com/makiuchi-d/gozxing/qrcode"
	"github.com/skip2/go-qrcode"
)

// QrcodeFromStr 将字符串转为二维码
func QrcodeFromStr(str string) (string, error) {
	var png []byte
	png, err := qrcode.Encode(str, qrcode.High, 512)
	if err != nil {
		return "", err
	}

	base64Str := base64.StdEncoding.EncodeToString(png)
	return base64Str, nil
}

// 解读二维码
func DecodeQrcodeFromStr(base64Str string) (string, error) {
	// 转Image
	png, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	// 转Image
	img, _, err := image.Decode(bytes.NewReader(png))
	if err != nil {
		return "", err
	}
	// 解码 使用zxing
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}
	qrReader := qrcode2.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	return result.GetText(), nil
}
