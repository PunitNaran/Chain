package blockchain

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"log"
	"time"

	"localhost/chain/lib/transaction"
)

// Block has some inorma
type Block struct {
	Timestamp    int64                      // Current timestamp
	Hash         []byte                     // Hash of this block
	Transactions []*transaction.Transaction // Data inside of this block
	PrevHash     []byte                     // Last blocks hash
	Nonce        int
}

// HashTransaction
func (b *Block) HashTransaction() []byte {
	var txHashes [][]byte
	var txHash [64]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha512.Sum512(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// CreateBlock creates a new block
func CreateBlock(transactions []*transaction.Transaction, prevHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte{}, transactions, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()
	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

// Aum starting of the block
func Aum(coinbase *transaction.Transaction) *Block {
	return CreateBlock([]*transaction.Transaction{coinbase}, []byte{})
}

// Serialize encodes the data so it can be compatiably encodable for BadgerDB
func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	Handle(encoder.Encode(b))
	return res.Bytes()
}

// Deserialize decodes the data so it can be compatiably decodable for BadgerDB
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	Handle(decoder.Decode(&block))
	return &block
}

// Handle throws an error is something unexpected happeneds
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
