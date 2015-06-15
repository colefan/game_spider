package page

import (
	"net/http"
)

//请求的结果页面
type Page struct {
	errMsg string
	isfail bool
	//页面类型：LIST OR CONTENT
	urlTag       string
	req          *Request
	bodyStr      string
	header       http.Header
	cookies      []*http.Cookie
	charset      string
	moreNewLinks map[string]string
	parserResult []*PageItems
	canbreak     bool
}

func NewPage(req *Request) *Page {
	links := make(map[string]string)
	itemList := make([]*PageItems, 0)
	return &Page{req: req, moreNewLinks: links, parserResult: itemList}
}

func (this *Page) GetErrMsg() string {
	return this.errMsg
}
func (this *Page) IsSucc() bool {
	return !this.isfail
}

func (this *Page) SetUrlTag(tag string) {
	this.urlTag = tag
}

func (this *Page) GetUrlTag() string {
	return this.urlTag
}

func (this *Page) AddNewUrl(url string, tag string) {
	for tmpUrl, _ := range this.moreNewLinks {
		if tmpUrl == url {
			return
		}
	}
	this.moreNewLinks[url] = tag
}

func (this *Page) GetNewUrls() map[string]string {
	return this.moreNewLinks
}

func (this *Page) CountNewUrls() int {
	return len(this.moreNewLinks)
}

func (this *Page) SetStatus(fail bool, errMsg string) {
	this.isfail = fail
	this.errMsg = errMsg
}

func (this *Page) SetBody(bodyStr string) {
	this.bodyStr = bodyStr
}

func (this *Page) GetBody() string {
	return this.bodyStr
}

func (this *Page) SetHeader(header http.Header) {
	this.header = header
}

func (this *Page) GetHeader() http.Header {
	return this.header
}

func (this *Page) SetCookies(cookies []*http.Cookie) {
	this.cookies = cookies
}

func (this *Page) GetCookies() []*http.Cookie {
	return this.cookies
}

func (this *Page) GetCharset() string {
	this.charset = this.header.Get("charset")
	return this.charset
}

func (this *Page) GetRequest() *Request {
	return this.req
}

func (this *Page) GetPageItemsList() []*PageItems {
	return this.parserResult
}

func (this *Page) AddPageItems(page *PageItems) {
	this.parserResult = append(this.parserResult, page)
}

func (this *Page) SetBreak(b bool) {
	this.canbreak = b
}

func (this *Page) IsBreak() bool {
	return this.canbreak
}
