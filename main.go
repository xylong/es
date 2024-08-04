package main

import (
	"context"
	"encoding/json"
	"es/dao"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
)

func init() {
	ConnectDB()
}

type News struct {
	Msgid      string
	Action     string
	From       string
	FromUserid string
	FromName   string
	FromAvatar string
	Tolist     []string
	Roomid     string
	Msgtime    int64
	Msgtype    string
	Url        string
	Detail     any
}

func main() {
	//db := ConnectDB()
	//GenerateTableStruct(db)
	var err error

	dao.SetDefault(GetDB())
	rows := getAll(context.Background(), dao.Q)
	if len(rows) == 0 {
		return
	}

	bulk := es.Bulk()
	for _, row := range rows {
		var m map[string]any

		err = json.Unmarshal([]byte(*row.Metadata), &m)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		news := News{
			Msgid:      m["msgid"].(string),
			Action:     m["action"].(string),
			From:       m["from"].(string),
			FromUserid: "",
			FromName:   "",
			FromAvatar: "",
			//Tolist:     nil,
			Roomid:  cast.ToString(m["roomid"]),
			Msgtime: cast.ToInt64(m["msgtime"]),
			Msgtype: m["msgtype"].(string),
			//Url:     *row.URL,
			//Detail: cast.ToString(m["metadata"]),
		}

		if row.URL != nil {
			news.Url = *row.URL
		}
		//if v, ok := m["tolist"]; ok {
		//
		//}
		if v, ok := m[news.Msgtype]; ok {
			news.Detail = v
		}

		req := elastic.NewBulkIndexRequest().Index("news").
			Id(news.Msgid).Doc(news)
		bulk.Add(req)
	}

	rsp, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rsp)
	}

}
