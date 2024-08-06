package main

import (
	"context"
	"es/dao"
	"es/model"
	"log"
)

func getAll(ctx context.Context, query *dao.Query) []*model.Chat_data {
	data := query.Chat_data

	rows, err := query.WithContext(ctx).Chat_data.
		Select(data.ID, data.URL, data.Metadata).
		Find()
	if err != nil {
		log.Fatalln(err)
	}

	return rows
}

func getData() []*model.Chat_data {
	db := ConnectDB()
	//GenerateTableStruct(db)
	dao.SetDefault(db)
	return getAll(context.Background(), dao.Q)
}
