package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type blockchain struct {
	currentTransactions []*transaction
	chain               []*block
	nodes               map[string]string
}

type block struct {
	Index        int            `json:"index"`
	Timestamp    time.Time      `json:"timestamp"`
	Transactions []*transaction `json:"transactions"`
	Proof        int            `json:"proof"`
	PreviousHash string         `json:"previous_hash"`
}

type transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

// NewBlockchain sets up a new blockchain object complete with genesis block
func newBlockchain() *blockchain {
	bc := &blockchain{}
	bc.createGenesisBlock()
	bc.nodes = map[string]string{}
	return bc
}

// creates the genesis block
func (bc *blockchain) createGenesisBlock() {
	bc.createBlock(100, "1")
}

// newBlock creates a new block and adds it to the blockchain
func (bc *blockchain) newBlock(proof int) *block {
	return bc.createBlock(proof, bc.hash(bc.lastBlock()))
}

func (bc *blockchain) createBlock(proof int, previousHash string) *block {
	// build the block
	b := &block{
		Index:        len(bc.chain) + 1,
		Timestamp:    time.Now(),
		Transactions: bc.currentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}

	// Add block to the chain
	bc.chain = append(bc.chain, b)
	// Reset the current list of transactions
	bc.currentTransactions = []*transaction{}
	return b
}

// newTransaction adds a new transaction to the list of transactions
func (bc *blockchain) newTransaction(trans *transaction) int {
	bc.currentTransactions = append(bc.currentTransactions, trans)
	return bc.lastBlockIndex() + 1
}

// hash creates a SHA-256 a block
func (bc *blockchain) hash(b *block) string {
	blkString, err := json.Marshal(b)
	if err != nil {
		log.Println("Error marshalling block: ", err)
	}
	hsh := sha256.Sum256(blkString)
	return fmt.Sprintf("%x", hsh)
}

// Simple Proof of Work Algorithm:
// - Find a number p' such that hash(pp') contains leading 4 zeroes,
// - p is the previous proof, and p' is the new proof
func (bc *blockchain) proofOfWork(lastProof int) (proof int) {
	for {
		if validProof(lastProof, proof) {
			return
		}
		proof++
	}
}

// Add a new node to the list of nodes
func (bc *blockchain) registerNode(address string) {
	bc.nodes[address] = address
}

// Determine if a given blockchain is valid
func validChain(chain []*block) bool {
	for i := 1; i < len(chain); i++ {
		b := chain[i]
		if b.PreviousHash != bc.hash(b) {
			return false
		}

		if !validProof(chain[i-1].Proof, b.Proof) {
			return false
		}
	}
	return true
}

// This is our Consensus Algorithm, it resolves conflicts
// by replacing our chain with the longest one in the network.
func (bc *blockchain) resolveConflict() bool {
	maxLength := len(bc.chain)
	newChain := []*block{}

	type chain struct {
		Length int      `json:"length"`
		Chain  []*block `json:"chain"`
	}

	for _, node := range bc.nodes {
		req, err := requestChainFrom(fmt.Sprintf("%s/chain", node))

		if err != nil {
			log.Println("Error fetching chain from node: ", node, err)
			continue
		}

		c := chain{}
		err = json.NewDecoder(req.Body).Decode(&c)
		if err != nil {
			log.Println("Error decoding chain info from node: ", node, err)
			continue
		}

		// ordering throws off the hashing
		if c.Length > maxLength /*&& validChain(c.Chain) */ {
			maxLength = c.Length
			newChain = c.Chain
			log.Printf("This is now c: %+v", c)
		}
	}

	if newChain != nil {
		bc.chain = newChain
		return true
	}
	return false
}

func validProof(lastProof int, proof int) bool {
	guess := sha256.Sum256([]byte(fmt.Sprintf("%d%d", lastProof, proof)))
	guessHash := fmt.Sprintf("%x", guess)
	return guessHash[:4] == "0000"
}

// returns the last block in the chain
func (bc *blockchain) lastBlock() *block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *blockchain) lastBlockIndex() int {
	return bc.lastBlock().Index
}
