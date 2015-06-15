package scheduler

import (
	"git.oschina.net/ciweilao/game_spider.git/page"
)

type Scheduler interface {
	Push(req *page.Request)
	Poll() *page.Request
	Count() int
}
