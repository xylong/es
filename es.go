package main

import (
	"github.com/olivere/elastic/v7"
	"log"
)

var (
	es *elastic.Client
)

func init() {
	var err error

	es, err = elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200/"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatal(err)
	}
}
