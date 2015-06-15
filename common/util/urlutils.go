package util

import (
	"strings"
)

func GetUrlDomain(url string) string {
	if strings.Index(url, "http://") >= 0 {
		tmpUrl := string(url[7:len(url)])
		//println(tmpUrl)
		pos1 := strings.Index(tmpUrl, "/")
		if pos1 >= 0 {
			return "http://" + string(tmpUrl[0:pos1])
		} else {
			return url
		}
	} else {
		pos1 := strings.Index(url, "/")
		if pos1 >= 0 {
			return "http://" + string(url[0:pos1])
		} else {
			return "http://" + url
		}
	}
}

func IsRelativePath(path string) bool {
	if len(path) == 0 {
		return false
	}
	if strings.Index(path, "http:") >= 0 {
		return false
	} else if path[0] == '/' {
		return true
	} else if path[0] == '.' {
		return true
	} else {
		return false
	}
}

func GetRealUrl(baseUrl string, relativePath string) string {
	if len(relativePath) <= 0 {
		return ""
	}

	if relativePath[0] == '/' {
		return GetUrlDomain(baseUrl) + relativePath
	} else if relativePath[0] == '.' {
		pos1 := strings.Index(baseUrl, "http://")
		pos2 := strings.LastIndex(baseUrl, "/")
		if (pos2 > pos1) && pos2 > 0 {
			return string(baseUrl[0:pos2]) + "/" + relativePath
		} else {
			return baseUrl + "/" + relativePath
		}
	} else {
		return GetUrlDomain(baseUrl) + "/" + relativePath
	}

}
