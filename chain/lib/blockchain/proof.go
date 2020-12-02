package blockchain

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The first few bytes must contain 0s

// This should slowly increase in having miners and the computation time
const Difficulty = 12

// ProofOfWork
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// IntToHex converts an int64 to a byte array
func InttoHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// NewProof produces a pointer to a proof of work
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(512-Difficulty))
	return &ProofOfWork{b, target}
}

// InitData grabs previous block hash and block data and combines it
// create a cohisive set of bytes
func (pow *ProofOfWork) InitData(nonce int) []byte {
	return bytes.Join([][]byte{
		InttoHex(pow.Block.Timestamp),
		pow.Block.PrevHash,
		pow.Block.HashTransaction(),
		InttoHex(int64(nonce)),
		InttoHex(int64(Difficulty)),
	},
		[]byte{},
	)
}

// Run prepares the data, hashes it, convests it into a big
// integer compares it with the target, this repeats recurrsively
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashVal big.Int
	var hash [64]byte
	nonce := 0
	fmt.Printf("Mining a new block")
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha512.Sum512(data)
		// fmt.Printf("\r%x", hash)
		hashVal.SetBytes(hash[:])
		if hashVal.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Println()
	return nonce, hash[:]
}

// Validate verifies that the proof of work is correct
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha512.Sum512(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
