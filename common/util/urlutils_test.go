package util

import (
	"testing"
)

func TestUrl(t *testing.T) {
	println("Test GetDomin : " + GetUrlDomain("www.dde.com/dsfd.h"))
}

func TestGetRealUrl(t *testing.T) {
	println("Test get Real url : " + GetRealUrl("http://www.google.com/zixun/dd.html", "/zixun/2.html"))
}
