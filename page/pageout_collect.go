package page

type PageOutCollectPageItems struct {
	collector []*PageItems
}

func NewPageOutCollectPageItems() *PageOutCollectPageItems {
	collector := make([]*PageItems, 0)
	return &PageOutCollectPageItems{collector: collector}
}

func (this *PageOutCollectPageItems) Process(itemList []*PageItems, taskname string) {
	for _, item := range itemList {
		this.collector = append(this.collector, item)
	}

}

func (this *PageOutCollectPageItems) Release() {

}

func (this *PageOutCollectPageItems) GetCollected() []*PageItems {
	return this.collector
}
