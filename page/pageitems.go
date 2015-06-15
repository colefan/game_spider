package page

type PageItems struct {
	items map[string]string
	skip  bool
	key   string
}

func NewPageItems(keyUrl string) *PageItems {
	list := make(map[string]string)
	return &PageItems{items: list, key: keyUrl}
}

func (this *PageItems) SetSkip(b bool) {
	this.skip = b
}

func (this *PageItems) Skipped() bool {
	return this.skip
}

func (this *PageItems) AddItem(key string, value string) {
	_, ok := this.items[key]
	if !ok {
		this.items[key] = value
	}
}

func (this *PageItems) GetItem(key string) string {
	v, ok := this.items[key]
	if !ok {
		return ""
	} else {
		return v
	}
}

func (this *PageItems) SetKey(k string) {
	this.key = k
}

func (this *PageItems) GetKey() string {
	return this.key
}
