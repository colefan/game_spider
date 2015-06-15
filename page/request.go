package page

import (
	"net/http"
)

//请求页面
type Request struct {
	url      string
	respType string
	method   string
	header   http.Header
	urlTag   string
	postData string
}

func NewRequest(url string, respType string, method string) *Request {
	return &Request{url: url, respType: respType, method: method}
}

func (this *Request) GetUrl() string {
	return this.url
}

func (this *Request) GetUrlTag() string {
	return this.urlTag
}

func (this *Request) SetUrlTag(tag string) {
	this.urlTag = tag
}

func (this *Request) GetRespType() string {
	return this.respType
}

func (this *Request) GetMethod() string {
	return this.method
}

func (this *Request) SetPostData(postdata string) {
	this.postData = postdata
}

func (this *Request) GetPostData() string {
	return this.postData
}
