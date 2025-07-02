package rexCommon

import "strings"

func IsMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	return isMobile
}

func IsWxClient(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isWx := false
	mobileKeywords := []string{"micromessenger", "MicroMessenger"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isWx = true
			break
		}
	}

	return isWx
}

func IsDouyinClient(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	is := false
	mobileKeywords := []string{"aweme"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			is = true
			break
		}
	}

	return is
}
