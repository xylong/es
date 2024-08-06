package main

import (
	"context"
	"encoding/json"
	"es/model"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/cast"
	"sync"
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
	var (
		wg          sync.WaitGroup
		dataCh      = make(chan *model.Chat_data, 100)
		resultCh    = make(chan *News, 1000)
		workerCount = 10
	)

	rows := getData()
	if len(rows) == 0 {
		return
	}

	t1 := time.Now()
	fmt.Println(t1.Format(time.DateTime))

	// 启动10个工作进程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			processData(context.Background(), dataCh, resultCh)
		}()
	}

	// 发送数据到dataCh
	for _, row := range rows {
		dataCh <- row
	}
	close(dataCh)

	// 等待所有结果处理完毕
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	bulk := es.Bulk()
	for ch := range resultCh {
		req := elastic.NewBulkIndexRequest().Index("chatdata").
			Id(ch.Msgid).Doc(ch)
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

	t2 := time.Now()
	fmt.Println(t2.Format(time.DateTime))
	fmt.Println(t2.Sub(t1))
}

func processData(ctx context.Context, dataCh <-chan *model.Chat_data, resultCh chan<- *News) {
	for {
		select {
		case data, ok := <-dataCh:
			if !ok {
				return
			}
			resultCh <- processItem(data)
		case <-ctx.Done():
			return
		}
	}
}

func processItem(data *model.Chat_data) *News {
	var (
		m map[string]any
	)

	_ = json.Unmarshal([]byte(*data.Metadata), &m)

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

	if data.URL != nil {
		news.Url = *data.URL
	}
	if v, ok := m["tolist"]; ok {
		news.Tolist = cast.ToStringSlice(v)
	}
	if v, ok := m[news.Msgtype]; ok {
		news.Detail = v
	}

	return &news
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
