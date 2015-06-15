package page

//页面分析器，用来具体分析一个页面的信息，将分析的结果按需输出，与具体的爬虫逻辑有关系
type PageProcesser interface {
	Process(p *Page)
}
