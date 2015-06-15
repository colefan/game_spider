package spider

import (
	"git.oschina.net/ciweilao/game_spider.git/common/logs"
	"git.oschina.net/ciweilao/game_spider.git/downloader"
	"git.oschina.net/ciweilao/game_spider.git/page"
	"git.oschina.net/ciweilao/game_spider.git/scheduler"
	"sync"
	"time"
)

//蜘蛛的启动程序
//蜘蛛根据入口链接，抓取相应的内容，进行分析，并发现同类型的链接交给调度器去进行调度和分析
type Spider struct {
	m_taskName      string
	m_pageProcesser page.PageProcesser
	m_downLoader    downloader.DownLoader
	m_scheduler     scheduler.Scheduler
	m_outputs       []page.PageOut
	m_exitWhenDone  bool

	//管理锁
	m_lock sync.Mutex
	//处理中的对象
	m_handlingMap map[string]string
}

func NewSpider(processerInst page.PageProcesser, taskName string) *Spider {
	spiderInst := &Spider{m_taskName: taskName, m_pageProcesser: processerInst}
	spiderInst.m_exitWhenDone = true

	if spiderInst.m_scheduler == nil {
		spiderInst.SetScheduler(scheduler.NewQueueScheduler(false))
	}

	if spiderInst.m_downLoader == nil {
		spiderInst.SetDownLoader(downloader.NewHttpDownLoader())
	}

	logs.GetFirstLogger().Info("*** start spider ***")
	spiderInst.m_outputs = make([]page.PageOut, 0)
	spiderInst.m_handlingMap = make(map[string]string)
	return spiderInst

}

func (this *Spider) TaskName() string {
	return this.m_taskName
}

func (this *Spider) SetScheduler(inst scheduler.Scheduler) {
	this.m_scheduler = inst
}

func (this *Spider) SetPageProcesser(inst page.PageProcesser) {
	this.m_pageProcesser = inst
}

func (this *Spider) SetDownLoader(inst downloader.DownLoader) {
	this.m_downLoader = inst
}

//通过url分析一个页面，获得PageItems
func (this *Spider) Get(url string, respType string) *page.PageItems {
	//TODO
	req := page.NewRequest(url, respType, "GET")
	return this.GetByRequest(req)

}

func (this *Spider) GetByRequest(req *page.Request) *page.PageItems {
	var reqs []*page.Request
	reqs = append(reqs, req)
	items := this.GetAllByRequest(reqs)
	if len(items) != 0 {
		return items[0]
	}
	return nil
}

func (this *Spider) GetAllByRequest(reqs []*page.Request) []*page.PageItems {
	for _, req := range reqs {
		this.AddRequest(req)
	}

	output := page.NewPageOutCollectPageItems()
	this.AddPageOut(output)
	this.Run()

	return output.GetCollected()
}

func (this *Spider) AddRequest(req *page.Request) *Spider {
	if req == nil {
		logs.GetFirstLogger().Error("request is nil")
		return this
	} else if req.GetUrl() == "" {
		logs.GetFirstLogger().Error("request is empty")
	}

	this.m_scheduler.Push(req)
	return this
}

func (this *Spider) AddPageOut(out page.PageOut) *Spider {
	this.m_outputs = append(this.m_outputs, out)
	return this
}

func (this *Spider) AddUrl(url string, respType string, urlType string) *Spider {
	req := page.NewRequest(url, respType, "GET")
	req.SetUrlTag(urlType)
	this.AddRequest(req)
	return this
}

func (this *Spider) waitForReqProcesser(url string, tag string) {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()
	this.m_handlingMap[url] = tag

}

func (this *Spider) finishForReqProcesser(url string) {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()
	delete(this.m_handlingMap, url)

}

func (this *Spider) countHandlingUrl() int {
	this.m_lock.Lock()
	defer this.m_lock.Unlock()
	return len(this.m_handlingMap)
}

func (this *Spider) Run() *Spider {
	for {
		var req *page.Request = this.m_scheduler.Poll()

		if req == nil {
			if this.countHandlingUrl() == 0 {
				break
			}

			time.Sleep(500 * time.Millisecond)
			//这里需要判断有没有没有处理完的请求，如果都处理完了可以退出，如果还没有处理完，那么继续等待下一个时间片段
			continue
		}

		//req.GetUrl
		this.waitForReqProcesser(req.GetUrl(), req.GetUrlTag())

		go func(*page.Request) {
			//deal the page,may get new pages
			this.pageProcess(req)
		}(req)

	}

	//运行结束
	return this

}

func (this *Spider) pageProcess(req *page.Request) {
	var p *page.Page
	//下载页面
	for i := 0; i < 3; i++ {
		p = this.m_downLoader.DownLoad(req)
		if p.IsSucc() {
			break
		}
		time.Sleep(time.Microsecond * 1000)
	}

	if !p.IsSucc() {
		this.finishForReqProcesser(req.GetUrl())
		return
	}

	//分析页面内容
	this.m_pageProcesser.Process(p)

	//获取新的链接
	if p.CountNewUrls() > 0 {
		newUrls := p.GetNewUrls()
		for tmpUrl, tmpUrlTag := range newUrls {
			this.AddUrl(tmpUrl, "html", tmpUrlTag)
		}
	}

	this.finishForReqProcesser(req.GetUrl())

	//输出
	for _, tmpOut := range this.m_outputs {
		tmpOut.Process(p.GetPageItemsList(), p.GetRequest().GetUrl())
	}

}
