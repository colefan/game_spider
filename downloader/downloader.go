package downloader

import (
	"git.oschina.net/ciweilao/game_spider.git/page"
)

type DownLoader interface {
	DownLoad(req *page.Request) *page.Page
}
