package uaparser

import "strings"

func Parse(ua string) (browser, os string, isBot bool) {
	lowUA := strings.ToLower(ua)

	// Browser
	browser = "Other"
	if strings.Contains(lowUA, "chrome") {
		browser = "Chrome"
	}
	if strings.Contains(lowUA, "firefox") {
		browser = "Firefox"
	}
	if strings.Contains(lowUA, "safari") && !strings.Contains(lowUA, "chrome") {
		browser = "Safari"
	}

	// OS
	os = "Other"
	if strings.Contains(lowUA, "windows") {
		os = "Windows"
	}
	if strings.Contains(lowUA, "linux") {
		os = "Linux"
	}
	if strings.Contains(lowUA, "android") {
		os = "Android"
	}
	if strings.Contains(lowUA, "iphone") {
		os = "iOS"
	}

	// Bot
	isBot = strings.Contains(lowUA, "bot") || strings.Contains(lowUA, "spider") || strings.Contains(lowUA, "crawler")

	return
}
