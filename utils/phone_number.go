package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	PhoneNumberMaskPrefix = 3
	PhoneNumberMaskSuffix = 2
	MinPhoneNumberLength  = PhoneNumberMaskPrefix + PhoneNumberMaskSuffix
)

func MaskPhoneNumber(number string) string {
	length := utf8.RuneCountInString(number)
	if length <= MinPhoneNumberLength {
		return strings.Repeat("*", length)
	}
	return fmt.Sprintf("%s%s%s",
		number[:PhoneNumberMaskPrefix],
		strings.Repeat("*", length-PhoneNumberMaskPrefix-PhoneNumberMaskSuffix),
		number[length-PhoneNumberMaskSuffix:length])
}
