package utils

import (
	"os"
	"strings"
)

// RemoveDomainError checks for the domain error
func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	return true
}

// EnforceHTTP enforces every url start with protocol pref http
func EnforceHTTP(url string) string {
	// make every url http
	if url[:4] != "http" {
		return "http://" + url
	}
	return url
}
