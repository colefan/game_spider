package main

import (
	"git.oschina.net/ciweilao/game_spider.git/common/logs"
	"git.oschina.net/ciweilao/game_spider.git/common/util"
	"git.oschina.net/ciweilao/game_spider.git/page"
	"git.oschina.net/ciweilao/game_spider.git/spider"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type YouxiduoProcesser struct {
	exitDate string
}

func NewYouxiduoProcesser(date string) *YouxiduoProcesser {
	t := time.Now()
	todayTimeStr := t.Format("2006-01-02")
	if len(date) <= 1 {
		date = todayTimeStr
	}
	return &YouxiduoProcesser{exitDate: date}
}

func (this *YouxiduoProcesser) parseNewsLinkListInfo(content string, p *page.Page) *page.Page {
	//println("B LINK URLS")
	if p.IsBreak() {
		return p
	}
	reg, _ := regexp.Compile(`<a href(.)*<\/a>`)
	urlStr := reg.FindAllString(content, -1)
	for _, tmp := range urlStr {
		var pos1 int = strings.Index(tmp, "href=")
		var pos2 int = strings.Index(tmp, ">")
		if (pos2 - 1) > (pos1 + 6) {
			tmp = string(tmp[pos1+6 : pos2-1])
			if strings.Index(tmp, "http://") >= 0 {
				continue
			}
			tmp = util.GetRealUrl(p.GetRequest().GetUrl(), tmp)
			p.AddNewUrl(tmp, "list")
			//	println("list url = " + tmp)
		}
	}
	//println("E LINK URLS")
	return p
}

func (this *YouxiduoProcesser) parseNewsBreifInfo(content string, p *page.Page) *page.Page {
	logs.GetFirstLogger().Trace("B TEST LIST ITEMS")
	var pos1 int = strings.Index(content, "<li>")
	var pos2 int = strings.Index(content, "</li>")
	var count int = 1

	for pos1 >= 0 && pos2 >= 0 && (pos2 > pos1) {
		item := page.NewPageItems("")
		tmpStr := string(content[pos1 : pos2+5])
		content = string(content[pos2+5 : len(content)])

		pos1 = strings.Index(content, "<li>")
		pos2 = strings.Index(content, "</li>")
		logs.GetFirstLogger().Trace("B================>")
		reg, _ := regexp.Compile(`<span>(.)*[\d]{4}-[\d]{2}-[\d]{2}`)
		timeStr := reg.FindString(tmpStr)
		reg, _ = regexp.Compile(`[\d]{4}-[\d]{2}-[\d]{2}`)
		timeStr = reg.FindString(timeStr)
		if this.exitDate > timeStr {
			p.SetBreak(true)
			continue
		}
		item.AddItem("time", timeStr)

		reg, _ = regexp.Compile("title=\"(.)*\"")
		title := reg.FindString(tmpStr)
		title = string(title[strings.Index(title, "\"")+1 : len(title)])
		title = string(title[0:strings.Index(title, "\"")])
		logs.GetFirstLogger().Trace("title = " + title)
		//p.AddResultItem("title", title)
		item.AddItem("title", title)
		reg, _ = regexp.Compile("<img src=(.)*alt")
		pic := reg.FindString(tmpStr)
		pic = string(pic[strings.Index(pic, "\"")+1 : len(pic)])
		pic = string(pic[0:strings.Index(pic, "\"")])

		if util.IsRelativePath(pic) {
			pic = util.GetRealUrl(p.GetRequest().GetUrl(), pic)
		}
		logs.GetFirstLogger().Trace("pic = " + pic)
		//p.AddResultItem("pic", pic)
		item.AddItem("pic", pic)

		reg, _ = regexp.Compile("<p>(.)*</p>")
		info := reg.FindString(tmpStr)
		logs.GetFirstLogger().Trace("info = " + info)
		//p.AddResultItem("info", info)
		info = strings.Replace(info, "'", "\"", -1)
		info = strings.Replace(info, "&#39;", "\"", -1)

		item.AddItem("info", info)

		reg, _ = regexp.Compile("<span(.)*<a(.)*</span>")
		detailurl := reg.FindString(tmpStr)
		reg, _ = regexp.Compile("href(.)*\">")
		detailurl = reg.FindString(detailurl)
		detailurl = detailurl[strings.Index(detailurl, "\"")+1 : len(detailurl)]
		detailurl = detailurl[0:strings.Index(detailurl, "\"")]
		logs.GetFirstLogger().Trace("detailurl = " + detailurl)
		//p.AddResultItem("detailurl", detailurl)
		item.AddItem("detailurl", detailurl)
		//p.AddResultItem("key", detailurl)
		item.SetKey(detailurl)
		p.AddNewUrl(detailurl, "content")

		logs.GetFirstLogger().Trace("E================>")
		logs.GetFirstLogger().Tracef("count = %d", count)
		count = count + 1
		logs.GetFirstLogger().Warn(title)

		pos1 = strings.Index(content, "<li>")
		pos2 = strings.Index(content, "</li>")
		p.AddPageItems(item)
	}

	return p
}

func (this *YouxiduoProcesser) parseNewsDetail(content string, p *page.Page) *page.Page {
	logs.GetFirstLogger().Trace("B TEST ARTICLE")
	//println(content)
	//tile , 不用考虑，在前面已经获取过了
	item := page.NewPageItems(p.GetRequest().GetUrl())
	reg, _ := regexp.Compile(`<div><span><em(.)*<\/span></div>`)
	newssrc := reg.FindString(content)

	//news_src,新闻来源
	reg, _ = regexp.Compile(`<a(.)*<\/a>`)
	newssrc = reg.FindString(newssrc)
	newssrc = newssrc[strings.Index(newssrc, ">")+1 : len(newssrc)]
	if strings.Index(newssrc, "<") >= 0 {
		newssrc = newssrc[0:strings.Index(newssrc, "<")]
	}

	logs.GetFirstLogger().Trace("newssrc = " + newssrc)
	//p.AddResultItem("news_src", newssrc)
	item.AddItem("news_src", newssrc)
	//news_content,新闻内容
	reg, _ = regexp.Compile(`<div class=\"artCon\">(.|\s)*<\/div>(\s)*<div class=\"pagebreak\"`)
	news := reg.FindString(content)
	if len(news) > 0 {
		pbIndex := strings.Index(news, "<div class=\"pagebreak\"")
		if pbIndex > 0 {
			news = news[0:pbIndex]
		}

	}
	newsIndex1 := strings.Index(news, ">")
	newsIndex2 := strings.Index(news, "</div>")
	if newsIndex1 >= 0 && newsIndex2 >= 0 {
		news = news[newsIndex1+1 : newsIndex2]
	}

	//p.AddResultItem("news_content", news)
	news = strings.Replace(news, "'", "\"", -1)
	news = strings.Replace(news, "&#39;", "\"", -1)
	//	imgSrcIndex := strings.Index(news, "<img src=\"/")
	//	if imgSrcIndex >= 0 {
	//		news = strings.Replace(news, "<img src=\"/", "<img src=\""+util.GetUrlDomain(p.GetRequest().GetUrl())+"/", -1)
	//	}
	////////////////////
	imgSrcIndex := strings.Index(news, "<img ")
	if imgSrcIndex >= 0 {
		news = strings.Replace(news, "<img src=\"/", "<img src=\""+util.GetUrlDomain(p.GetRequest().GetUrl())+"/", -1)
		news = strings.Replace(news, "<img alt=\"[^\"]\" src=\"/", "<img src=\""+util.GetUrlDomain(p.GetRequest().GetUrl())+"/", -1)
		//println(news_content)

		//	println("===============")
		reg, _ = regexp.Compile(`<img[^>]*>`)
		imgList := reg.FindAllString(news, -1)
		for _, img := range imgList {
			//strings.Replace(news_content, img)
			//println("old img ==>" + img)
			newImg := img
			styleIndex := strings.Index(newImg, "style=\"")
			if styleIndex >= 0 {
				styleStr := newImg[styleIndex+len("style=\""):]
				endIndex := strings.Index(styleStr, "\"")
				if endIndex > 0 {
					styleStr = styleStr[0:endIndex]
				}
				newstyleStr := changeImgSize(styleStr)
				newImg = strings.Replace(img, styleStr, newstyleStr, -1)

			} else {
				//找width，找height
				reg2, _ := regexp.Compile(`width=\"[0-9]+\"`)
				tmpWidthStr := reg2.FindString(img)

				reg2, _ = regexp.Compile(`height=\"[0-9]+\"`)
				tmpHeightStr := reg2.FindString(img)
				//println("tmp height str = " + tmpHeightStr)
				var f float32 = 1.0
				if len(tmpWidthStr) > 0 {
					tmpStr1 := tmpWidthStr[strings.Index(tmpWidthStr, "\"")+1:]
					tmpStr1 = tmpStr1[0:strings.Index(tmpStr1, "\"")]

					tmpWidth, _ := strconv.Atoi(tmpStr1)
					if tmpWidth > 360 {
						f = float32(tmpWidth) / 360.0
						if len(tmpHeightStr) > 0 {
							tmpStr2 := tmpHeightStr[strings.Index(tmpHeightStr, "\"")+1:]
							tmpStr2 = tmpStr2[0:strings.Index(tmpStr2, "\"")]
							tmpHeight, _ := strconv.Atoi(tmpStr2)

							newImg = strings.Replace(img, tmpWidthStr, "width=\"360\"", -1)
							tmpHeight = int(float32(tmpHeight) / f)
							newImg = strings.Replace(newImg, tmpHeightStr, "height=\""+strconv.Itoa(tmpHeight)+"\"", -1)

						} else {
							newImg = strings.Replace(img, tmpWidthStr, "width=\"360\"", -1)
						}

					}
				}

			}

			//有没有STYLE,有style的处理style
			//有没有width
			//有没有height
			//println("new img ==>" + newImg)
			if img != newImg {
				news = strings.Replace(news, img, newImg, -1)
			}
		}
	}
	//////
	news = strings.Replace(news, "<a[^>]*>官方网站</a>", "", -1)

	logs.GetFirstLogger().Trace("news = " + news)
	//判断是否有视频在新闻中，如有则过滤到哦
	reg, _ = regexp.Compile(`<[^>]*shockwave-flash[^>]*>`)
	tmpN := reg.FindString(news)

	item.AddItem("news_content", news)
	//p.AddResultItem("key", p.GetRequest().GetUrl())
	if len(tmpN) <= 0 {
		p.AddPageItems(item)
	}

	logs.GetFirstLogger().Trace("E TEST ARTICLE")

	return p
}

func (this *YouxiduoProcesser) Process(p *page.Page) {
	if !p.IsSucc() {
		println(p.GetErrMsg())
		return
	}
	var body string = p.GetBody()
	var urlTag string = p.GetRequest().GetUrlTag()
	p.SetUrlTag(urlTag)
	//分析这个页面是LIST页面还是内容页面
	// <div class="infroList"><ul><li>...</div>===>LIST
	// <div class="pagebreak">...</div>===>LIST
	// CONTENT
	//<div class="article"

	if urlTag == "list" {
		//
		//1.寻找news-brief的content
		regList, err := regexp.Compile(`<div class=\"infroList\">(\s|.)*<\/ul>(\s|.)*<div class=\"pagebreak\">`)
		if err != nil {
			logs.GetFirstLogger().Error("分析页面出错，正则表达式错误了，url = " + p.GetRequest().GetUrl())
		}
		var infroList []string = regList.FindAllString(body, -1)

		if len(infroList) > 0 {
			this.parseNewsBreifInfo(infroList[0], p)
		} else {
			logs.GetFirstLogger().Info("No more list items")
		}
		//先寻找额外的LIST页面
		if !p.IsBreak() {
			regPageBreak, err := regexp.Compile(`<div class=\"pagebreak\">(\s|.)+<li class=\"lastPage\">`)
			if err != nil {
				logs.GetFirstLogger().Error("分析页面出错，翻页正则表达式错误，url = " + p.GetRequest().GetUrl())
			}
			var pageBreakList []string = regPageBreak.FindAllString(body, -1)
			if len(pageBreakList) > 0 {
				this.parseNewsLinkListInfo(pageBreakList[0], p)
			} else {
				logs.GetFirstLogger().Info("No more links")
			}

		}

	} else {
		//CONTENT
		this.parseNewsDetail(body, p)
	}

}

func changeImgSize(style string) string {
	heightIndex := strings.Index(style, "height:")
	widthIndex := strings.Index(style, "width:")
	if widthIndex >= 0 && heightIndex >= 0 {
		heightStr := style[heightIndex+len("height:"):]
		widthStr := style[widthIndex+len("width:"):]

		var f float32 = 1.0

		if strings.Index(widthStr, "px") > 0 {
			widthStr = widthStr[0:strings.Index(widthStr, "px")]
			widthStr = strings.Replace(widthStr, " ", "", -1)

			width, _ := strconv.Atoi(widthStr)

			if width <= 360 {
				return style
			}
			f = float32(width) / 360.0

			heightStr = heightStr[0:strings.Index(heightStr, "px")]
			heightStr = strings.Replace(heightStr, " ", "", -1)
			height, _ := strconv.Atoi(heightStr)
			height = int(float32(height) / f)

			return "width: 400px; height: " + strconv.Itoa(height) + "px"

		} else {
			return style
		}

	} else {
		return style
	}

}

func main() {
	logs.GetFirstLogger().SetLevel("warn")
	logs.GetFirstLogger().Info("crawl begin www.youxiduo.com")
	var sp *spider.Spider
	sp = spider.NewSpider(NewYouxiduoProcesser(""), "youxiduo")
	out := page.NewPageOutSql()
	sp.AddUrl("http://www.youxiduo.com/zixun/game/", "html", "list").AddPageOut(out).Run()
	//sp.AddUrl("http://www.youxiduo.com/zixun/game/108444.shtml", "html", "content").Run()
	out.Release()
	logs.GetFirstLogger().Info("crawl end www.youxiduo.com")

}
