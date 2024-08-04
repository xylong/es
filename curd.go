package main

import (
	"context"
	"es/dao"
	"fmt"
	"log"
)

func getAll(ctx context.Context, query *dao.Query) {
	data := query.Chat_data

	rows, err := query.WithContext(ctx).Chat_data.
		Select(data.ID, data.URL, data.Metadata).
		Find()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(len(rows))
}
