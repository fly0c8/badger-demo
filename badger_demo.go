package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	badger "github.com/dgraph-io/badger/v2"

)

type Contract struct {
	Id          string
	Contractor  string
	Identifiers []string
}

var (
	contracts = []Contract{
		{
			Id:         "s1",
			Contractor: "skidata",
			Identifiers: []string{
				"lpn1", "lpn2", "lpn3",
			},
		},
		{
			Id:         "s2",
			Contractor: "skidata",
			Identifiers: []string{
				"lpn4", "lpn5", "lpn6",
			},
		},
	}
)

func main() {
	db, err := badger.Open(badger.DefaultOptions("contracts.db"))
	if err != nil {
		log.Fatal(err)
	}
	createContracts(db)
}



func createContracts(db *badger.DB) {
	for _, c := range contracts {
		if err := db.Update(func(txn *badger.Txn) error {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(c)
			if err := txn.Set([]byte(c.Id), buf.Bytes()); err != nil {
				return err
			}
		}); err !=nil {
			log.Fatal(err)
		}
	}
}
