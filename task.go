package main

import (
	"encoding/json"
	"es/model"
	"github.com/spf13/cast"
	"sync"
)

type Task struct {
	data   *model.Chat_data
	result chan<- *News
}

func (t *Task) do() {
	var (
		m map[string]any
	)

	_ = json.Unmarshal([]byte(*t.data.Metadata), &m)

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

	if t.data.URL != nil {
		news.Url = *t.data.URL
	}
	if v, ok := m["tolist"]; ok {
		news.Tolist = cast.ToStringSlice(v)
	}
	if v, ok := m[news.Msgtype]; ok {
		news.Detail = v
	}

	t.result <- &news
}

func InitTask(taskCh chan<- *Task, r chan *News, rows []*model.Chat_data) {
	for _, row := range rows {
		taskCh <- &Task{
			data:   row,
			result: r,
		}
	}
	close(taskCh)
}

func DistributeTask(taskChan <-chan *Task, wait *sync.WaitGroup, result chan *News) {
	for v := range taskChan {
		wait.Add(1)
		go ProcessTask(v, wait)
	}
	wait.Wait()
	close(result)
}

func ProcessTask(t *Task, wait *sync.WaitGroup) {
	t.do()
	wait.Done()
}

func ProcessResult(resultChan chan *News) []*News {
	arr := make([]*News, 0)

	for r := range resultChan {
		arr = append(arr, r)
	}

	return arr
}
