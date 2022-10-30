package database

import (
	"KPool/shares"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/kaspanet/kaspad/util"
)


//WORK IN PROGRESS (Or keep in ram?)

type sharesByMiner map[util.Address][]*shares.Share

var minerKey = []byte("sharesByMiner")

type Database struct {
	store *badger.DB
}


func NewDatabase() (*Database, error) {
	store, err := badger.Open(badger.DefaultOptions("./database"))
	if err != nil {
		return nil, err
	}
	return &Database{
		store: store,
	}, nil
}

func (db *Database) PutShare(share *shares.Share) error {
	searilizedShare, err := share.Searilize()
	if err != nil {
		return err
	}
	shareKey := share.Key()
	err = db.putSharesByMiner([]byte(share.PayoutAddress.String()), shareKey, searilizedShare)
	if err != nil {
		return err
	}
	return nil
}

func (db *Database) GetAllShares() ([]*shares.Share, error) {
	return nil, nil
}

func (db *Database) GetSharesByMiner(kaspaAddress util.Address) (sharesByMiner, error) {
	return nil, nil
}

func (db *Database) putSharesByMiner(kaspaAddress []byte, key []byte, value []byte) error {
	txn := db.store.NewTransaction(true)
	txn.Set([]byte(key),[]byte(value))
	return nil
}

func (db *Database) putSharesByBlock(kaspaAddress []byte, key []byte, value []byte) error {
	txn := db.store.NewTransaction(true)
	txn.Set([]byte(key),[]byte(value))
	return nil
}