package util

import "golang.org/x/crypto/bcrypt"

// BcryptEncrypt 加密字符串（密码）
func BcryptEncrypt(password string) (string, error) {
	// 生成盐值（cost=10，值越大加密越慢但越安全）
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// BcryptCompare 校验加密密码和明文密码是否匹配
func BcryptCompare(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil // 匹配返回true
}