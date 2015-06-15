package downloader

import (
	"git.oschina.net/ciweilao/game_spider.git/common/logs"
	"git.oschina.net/ciweilao/game_spider.git/page"
	"io/ioutil"
	"net/http"
	"strings"
)

type HttpDownLoader struct {
}

func NewHttpDownLoader() *HttpDownLoader {
	return &HttpDownLoader{}
}

func (this *HttpDownLoader) DownLoad(req *page.Request) *page.Page {
	var p = page.NewPage(req)
	var respType = req.GetRespType()
	switch respType {
	case "html":
		return this.downloadHtml(p, req)
	default:
		logs.GetFirstLogger().Error("error request type : " + respType)
	}

	return nil
}

func (this *HttpDownLoader) downloadHtml(p *page.Page, req *page.Request) *page.Page {
	p, destBody := this.downloadFile(p, req)
	if !p.IsSucc() {
		return p
	}
	p.SetBody(destBody)
	return p
}

//下载文件，并对字符编码做相应的处理
func (this *HttpDownLoader) downloadFile(p *page.Page, req *page.Request) (*page.Page, string) {
	var err error
	var httpResp *http.Response
	var urlStr string
	var method string
	urlStr = req.GetUrl()
	if len(urlStr) == 0 {
		logs.GetFirstLogger().Error("url is empty")
		p.SetStatus(true, "url is empty")
		return p, ""
	}

	method = req.GetMethod()

	if method == "POST" {
		httpResp, err = http.Post(req.GetUrl(), "application/x-www-form-urlencoded", strings.NewReader(req.GetPostData()))
	} else {
		httpResp, err = http.Get(req.GetUrl())
	}

	if err != nil {
		logs.GetFirstLogger().Error("http visit error :" + err.Error())
		p.SetStatus(true, err.Error())
	}
	p.SetHeader(httpResp.Header)
	p.SetCookies(httpResp.Cookies())
	body, _ := ioutil.ReadAll(httpResp.Body)
	bodyStr := string(body)
	defer httpResp.Body.Close()
	return p, bodyStr
}
