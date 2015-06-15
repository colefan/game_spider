package page

import (
	"database/sql"

	"git.oschina.net/ciweilao/game_spider.git/common/logs"
	_ "github.com/go-sql-driver/mysql"
)

type PageOutSql struct {
	pageResult map[string]*PageItems
	db         *sql.DB
	audit      int
}

func NewPageOutSql() *PageOutSql {
	db, err := sql.Open("mysql", "colefan:123456@tcp(192.168.13.21:3306)/news_tation")
	if err != nil {
		logs.GetFirstLogger().Error("mysql open error :" + err.Error())
	}

	list := make(map[string]*PageItems)
	return &PageOutSql{pageResult: list, db: db, audit: -1}
}

func (this *PageOutSql) Release() {
	if this.db != nil {
		this.db.Close()
	}

}

func (this *PageOutSql) Process(itemList []*PageItems, taskname string) {
	for _, item := range itemList {
		if len(item.GetItem("news_content")) > 0 {
			v, ok := this.pageResult[item.GetKey()]
			if ok {
				v.AddItem("news_content", item.GetItem("news_content"))
				v.AddItem("news_src", item.GetItem("news_src"))
				//v.AddItem("time", item.GetItem("time"))
				//输出这个数据
				this.output(v)
				delete(this.pageResult, item.GetKey())
			} else {
				this.pageResult[item.GetKey()] = item
			}

		} else if len(item.GetItem("title")) > 0 {
			v, ok := this.pageResult[item.GetKey()]
			if ok {
				v.AddItem("title", item.GetItem("title"))
				v.AddItem("pic", item.GetItem("pic"))
				v.AddItem("info", item.GetItem("info"))
				v.AddItem("detailurl", item.GetItem("detailurl"))
				v.AddItem("time", item.GetItem("time"))
				this.output(v)
				delete(this.pageResult, item.GetKey())

			} else {
				this.pageResult[item.GetKey()] = item
			}
		}

	}

}

func (this *PageOutSql) getConfigData() {
	if this.audit != -1 {
		return
	}
	err := this.db.Ping()
	if err != nil {
		logs.GetFirstLogger().Error("ping mysql error :" + err.Error())
		return
	}

	rows, err := this.db.Query("select status from deploy")
	if err != nil {
		logs.GetFirstLogger().Error("select status from deploy error :" + err.Error())
		return
	}
	defer rows.Close()
	var status int
	if rows.Next() {
		err := rows.Scan(&status)
		if err != nil {
			logs.GetFirstLogger().Error("fetch status value error :" + err.Error())
		}
	}

}

func (this *PageOutSql) output(item *PageItems) {
	logs.GetFirstLogger().Info("NEWS BEGIN==============================")
	logs.GetFirstLogger().Info("url\t= " + item.GetKey())
	logs.GetFirstLogger().Info("url2\t= " + item.GetItem("detailurl"))
	logs.GetFirstLogger().Info("title\t= " + item.GetItem("title"))
	logs.GetFirstLogger().Info("pic\t= " + item.GetItem("pic"))
	logs.GetFirstLogger().Info("breif\t= " + item.GetItem("info"))
	logs.GetFirstLogger().Info("time\t= " + item.GetItem("time"))
	logs.GetFirstLogger().Info("src\t= " + item.GetItem("news_src"))
	logs.GetFirstLogger().Info("news\t= " + item.GetItem("news_content"))
	logs.GetFirstLogger().Info("NEWS END  ==============================")
	sqlStr := "insert into news(`origin_url`,`title`,`icon_url`,`brief`,`get_time`,`origin`,`content`,`type`,`status`) values ("
	sqlStr = sqlStr + "'" + item.GetKey() + "'," + "'" + item.GetItem("title") + "'," + "'" + item.GetItem("pic") + "'," + "'" + item.GetItem("info") + "'," + "'" + item.GetItem("time") + "'," + "'" + item.GetItem("news_src") + "'," + "'" + item.GetItem("news_content") + "'," + "0,0)"
	//	logs.GetFirstLogger().Info(sqlStr)

	if this.audit == -1 {
		this.getConfigData()
	}

	logs.GetFirstLogger().Info(sqlStr)
	rows, err := this.db.Query("select count(*) from news where title='" + item.GetItem("title") + "'")
	if err != nil {
		logs.GetFirstLogger().Error("select count(*) from news error : " + err.Error())
		return
	}
	defer rows.Close()
	var sameCount int = 0

	if rows.Next() {
		err := rows.Scan(&sameCount)
		if err != nil {
			logs.GetFirstLogger().Error("fetch status value error :" + err.Error())
			return
		}
	}

	if sameCount > 0 {
		return
	}

	_, err2 := this.db.Exec(sqlStr)
	if err2 != nil {
		logs.GetFirstLogger().Error("insert error : " + err2.Error())
		logs.GetFirstLogger().Error("sql = " + sqlStr)
	}

}
