package main

import (
	"context"
	"encoding/json"
	"es/dao"
	"es/model"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"time"
)

func init() {
	ConnectDB()
}

type News struct {
	Msgid      string
	Action     string
	From       string
	FromType   int    `json:"from_type"`
	FromUserid string `json:"from_userid"`
	FromName   string `json:"from_name"`
	FromAvatar string `json:"from_avatar"`
	FromMobile string `json:"from_mobile"`
	Tolist     []string
	Roomid     string
	Msgtime    int64
	Msgdate    string
	Msgtype    string
	Url        string
	Detail     any
}

func main() {
	ctx := context.Background()
	db := ConnectDB()
	//GenerateTableStruct(db)

	dao.SetDefault(db)
	rows := getAll(context.Background(), dao.Q)
	if len(rows) == 0 {
		return
	}

	t1 := time.Now()
	fmt.Println(t1.Format(time.DateTime))
	syncChatData(ctx, rows)
	t2 := time.Now()
	fmt.Println(t2.Format(time.DateTime))
	fmt.Println(t2.Sub(t1))
}

func syncChatData(ctx context.Context, rows []*model.Chat_data) {
	var err error

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
		if v, ok := m["tolist"]; ok {
			news.Tolist = cast.ToStringSlice(v)
		}
		if v, ok := m[news.Msgtype]; ok {
			news.Detail = v
		}

		req := elastic.NewBulkIndexRequest().Index("chatdata").
			Id(news.Msgid).Doc(news)
		bulk.Add(req)
	}

	rsp, err := bulk.Do(context.Background())
	if err != nil {
		fmt.Println("aaaaaa")
		fmt.Println(err)
	} else {
		fmt.Println("bbbbbbbb")
		fmt.Println(rsp)
	}
}
