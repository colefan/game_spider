package scheduler

import (
	"container/list"
	"git.oschina.net/ciweilao/game_spider.git/page"
	"sync"
)

//队列调度器
type QueueScheduler struct {
	rm     bool
	lock   *sync.Mutex
	queue  *list.List
	rmKeys map[string]string
}

func NewQueueScheduler(rmDuplicate bool) *QueueScheduler {
	queue := list.New()
	lock := new(sync.Mutex)
	keys := make(map[string]string)
	return &QueueScheduler{rm: rmDuplicate, lock: lock, queue: queue, rmKeys: keys}
}

func (this *QueueScheduler) Push(req *page.Request) {
	//TODO
	this.lock.Lock()
	defer this.lock.Unlock()
	_, ok := this.rmKeys[req.GetUrl()]
	if ok {
		return
	}
	this.rmKeys[req.GetUrl()] = req.GetUrlTag()
	this.queue.PushBack(req)
}

func (this *QueueScheduler) Poll() *page.Request {
	//TODO
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.queue.Len() <= 0 {
		return nil
	}
	e := this.queue.Front()
	req := e.Value.(*page.Request)
	this.queue.Remove(e)
	return req
}

func (this *QueueScheduler) Count() int {
	return this.queue.Len()
}
