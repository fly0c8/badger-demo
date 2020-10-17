package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	badger "github.com/dgraph-io/badger/v2"
	"os"
	"strings"
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

func createManyContracts(db *badger.DB)  {
		if err := db.Update(func(txn *badger.Txn) error {
			for i:=0; i<60000; i++ {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				c := Contract{
					Id:          fmt.Sprintf("Contract:%d", i),
					Contractor:  fmt.Sprintf("Contractor:%d", i),
					Identifiers: []string {
						"lpn1","lpn2","lpn3",
					},
				}
				enc.Encode(c)
				if err := txn.Set([]byte(c.Id), buf.Bytes()); err != nil {
					return err
				}
			}
			return nil
		}); err !=nil {
			log.Fatal(err)
		}
}

func main() {
	db, err := badger.Open(badger.DefaultOptions("./contracts"))
	if err != nil {
		log.Fatal(err)
	}
	//createContracts(db, contracts)
	//createManyContracts(db)

	//if c, err := readContract(db, "Contract:6666"); err == nil {
	//	fmt.Printf("Found contract: %+v\n", c)
	//}
	//
	//if c, err := readContract(db, "Contract:47111"); err == nil {
	//	fmt.Printf("Found contract: %+v\n", c)
	//}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter ContractID:")
		contractID, _ := reader.ReadString('\n')
		contractID = "Contract:"+strings.TrimSpace(contractID)
		fmt.Println("Searching for:"+contractID)
		if c, err := readContract(db, contractID); err != nil {
			fmt.Println("No contract found:"+err.Error())
		} else {
			fmt.Printf("Found contract: %+v\n", c)
		}
	}


}

func readContract(db *badger.DB, contractId string) (*Contract, error) {
	var contract Contract
	if err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(contractId))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return gob.NewDecoder(bytes.NewReader(val)).Decode(&contract)
		});
	}); err != nil {
		return nil, err
	} else {
		return &contract, nil
	}
}
func createContracts(db *badger.DB, contracts []Contract) {
	for _, c := range contracts {
		if err := db.Update(func(txn *badger.Txn) error {
			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)
			enc.Encode(c)
			if err := txn.Set([]byte(c.Id), buf.Bytes()); err != nil {
				return err
			}
			return nil
		}); err !=nil {
			log.Fatal(err)
		}
	}
}
