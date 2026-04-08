package utils

import (
	"strconv"
	"unicode"
)

func ValidateUsername(username string) bool {
	// 用户名长度2-20位，只能包含汉字、字母、数字、下划线
	if len(username) < 2 || len(username) > 20 {
		return false
	}
	for _, ch := range username {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || unicode.Is(unicode.Han, ch)) {
			return false
		}
	}
	return true
}

func ValidatePhoneNumber(phone string) bool {
	// 简单的手机号验证，可以根据实际情况调整
	if len(phone) != 11 {
		return false
	}
	for _, ch := range phone {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return phone[0] == '1'
}

func IntToString(num int) string {
	if num < 10 {
		return string(rune(num + '0'))
	}
	return string(rune(num/10+'0')) + string(rune(num%10+'0'))
}

func StringToUint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	return uint(u64), err
}

func StringToInt(s string) (int, error) {
	i64, err := strconv.ParseInt(s, 10, 64)
	return int(i64), err
}
