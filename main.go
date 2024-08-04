package main

import (
	"context"
	"es/dao"
)

func init() {
	ConnectDB()
}

func main() {
	//db := ConnectDB()
	//GenerateTableStruct(db)

	dao.SetDefault(GetDB())
	getAll(context.Background(), dao.Q)
}
