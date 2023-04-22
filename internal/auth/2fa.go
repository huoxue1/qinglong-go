package auth

import (
	"github.com/pquerna/otp/totp"
)

// GenerateTOTP 生成 TOTP 密钥和二维码
func GenerateTOTP(user string, issuer string) (string, string, error) {
	// 生成 TOTP 密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: user,
	})
	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// VerifyTOTP 验证 TOTP 密码
func VerifyTOTP(secret string, code string) bool {
	return totp.Validate(code, secret)
}
