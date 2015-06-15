package page

import (
	"fmt"
)

type PageOutConsole struct {
}

func NewPageOutConsole() *PageOutConsole {
	return &PageOutConsole{}
}

func (this *PageOutConsole) Process(item []*PageItems, taskname string) {
	fmt.Println("------------------------------------------------------")
	fmt.Println("Crawled url:\t")

}

func (this *PageOutConsole) Release() {

}
