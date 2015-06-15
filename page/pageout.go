package page

type PageOut interface {
	Process(items []*PageItems, taskName string)
	Release()
}

type CollectPageOut interface {
	PageOut
	GetCollected() []*PageItems
}
