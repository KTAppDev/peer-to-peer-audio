package getbrowser

import (
	"strings"
)

func GetBrowser(userAgent string) string {
	if strings.Contains(userAgent, "Chrome") || strings.Contains(userAgent, "Firefox") {
		return "Chrome or Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		return "Safari"
	}
	return "Unknown"
}
