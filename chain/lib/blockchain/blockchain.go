package blockchain

import (
	"fmt"
	"localhost/chain/lib/transaction"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath  = "./tmp/blocks"
	dbFile  = "./tmp/blocks/MANIFEST"
	aumData = "First Transaction from Aum"
)

// BlockChain A chain of Blocks
type BlockChain struct {
	LastHash []byte
	Database *badger.DB.
}

// BlockchainIterator iterate over block in the blockchain
type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

// DBexists checks if the database is present
func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// ContinueBlockChain continues from the last block in the chain
func ContinueBlockChain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			if len(val) > 0 {
				lastHash = val
				err = nil
			} else {
				err = fmt.Errorf("Aum block/Last hash not present")
			}
			return err
		})
		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}
	return &chain
}

// InitBlockChain Initalise the BlockChain
func InitBlockChain(address string) *BlockChain {
	var lastHash []byte
	// Open the database
	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}
	opts := badger.DefaultOptions(dbPath)
	// Store keys, metadata and values in default options
	db, err := badger.Open(opts)
	Handle(err)
	err = db.Update(func(txn *badger.Txn) error {
		cbtx := transaction.NewCoinbaseTX(address, aumData)
		aum := Aum(cbtx)
		err = txn.Set(aum.Hash, aum.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), aum.Hash)
		lastHash = aum.Hash
		return err
	})
	Handle(err)
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

// AddBlock adds a lock to the blockchain
func (chain *BlockChain) AddBlock(transactions []*transaction.Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			if len(val) > 0 {
				lastHash = val
				err = nil
			} else {
				err = fmt.Errorf("Aum block/Last hash not present")
			}
			return err
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})

	Handle(err)
}

// Iterator obtains blocks from top to bottom (newest to oldest)
func (chain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

// Next This fetches the next block from the DB
func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			if len(val) > 0 {
				encodedBlock = val
				err = nil
			} else {
				err = fmt.Errorf("Aum block/Last hash not present")
			}
			return err
		})
		block = Deserialize(encodedBlock)
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash

	return block
}
