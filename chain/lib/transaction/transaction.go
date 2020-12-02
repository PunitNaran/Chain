package transaction

import (
	"bytes"
	"crypto/sha512"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction combination of inputs and outputs
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// TXOutput transaction output
type TXOutput struct {
	Value        int    // your coins
	ScriptPubKey string // the user defined wallet
}

// IsCoinbase checks whether the transaction is coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [64]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha512.Sum512(encoded.Bytes())
	tx.ID = hash[:]
}

// TXInput transaction input
type TXInput struct {
	Txid      []byte // stores ID of input -> output transaction
	Vout      int    // index of the output in the transaction
	ScriptSig string // data to be used by the ouputs SrciptPubKey
}

// CanUnlockOutputWith checks whether the address initiated the transaction
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith checks if the output can be unlocked with the provided data
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX creates a `coinbase` transaction - doesnt require any outputs
// to create inputs.
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	// creates an empty input and store arbitrary data intead of scriptsig
	txin := TXInput{[]byte{}, -1, data}
	// subsidy is the reward amount
	txout := TXOutput{subsidy, to}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}
