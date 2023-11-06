package captcha

import (
	captcha "github.com/mojocn/base64Captcha"
)

var driver = new(captcha.DriverString)

func init() {
	driver.Height = 80
	driver.Width = 240
	driver.NoiseCount = 0
	driver.ShowLineOptions = captcha.OptionShowSlimeLine | captcha.OptionShowHollowLine
	driver.Length = 4
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
}

// 生成验证码图片 base64
func GenerateCaptcha() (id, b64s string, err error) {
	// 生成base64图片
	c := captcha.NewCaptcha(driver, captcha.DefaultMemStore)
	return c.Generate()
}

// 验证验证码
func VerifyCaptcha(id, answer string) bool {
	c := captcha.NewCaptcha(driver, captcha.DefaultMemStore)
	return c.Verify(id, answer, true)
}
