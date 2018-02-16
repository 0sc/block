package main

import (
	"reflect"
	"testing"
)

func Test_newBlockchain(t *testing.T) {
	bc := newBlockchain()
	if l := len(bc.chain); l != 1 {
		t.Error("Expected block length to be 1 got ", l)
	}

	genesisBlock := bc.chain[0]

	if genesisBlock.Index != 1 {
		t.Error("Expected genesis block to have index 1 got ", genesisBlock.Index)
	}
}

func Test_blockchain_newBlock(t *testing.T) {
	bc := newBlockchain()

	b := bc.newBlock(500)

	if lastBlock := bc.chain[len(bc.chain)-1]; !reflect.DeepEqual(b, lastBlock) {
		t.Errorf("blockchain.newBlock() = %v, want %v", b, lastBlock)
	}
}

func Test_blockchain_createBlock(t *testing.T) {
	bc := blockchain{}
	b := bc.createBlock(1, "some-hash")

	if lastBlock := bc.chain[len(bc.chain)-1]; !reflect.DeepEqual(b, lastBlock) {
		t.Errorf("blockchain.createBlock() = %v, want %v", b, lastBlock)
	}
}

func Test_blockchain_newTransaction(t *testing.T) {
	bc := newBlockchain()
	if l := len(bc.currentTransactions); l != 0 {
		t.Errorf("expect len of current transactions to be 0 got %d", l)
	}

	blockIndex := bc.newTransaction(&transaction{})

	if blockIndex != 2 {
		t.Errorf("blockchain.newTransaction() = %v, want %v", blockIndex, 1)
	}

	if lastTransaction := bc.currentTransactions[len(bc.currentTransactions)-1]; !reflect.DeepEqual(&transaction{}, lastTransaction) {
		t.Errorf("blockchain.newTransaction() = %v, want %v", lastTransaction, &transaction{})
	}
}

func Test_blockchain_registerNode(t *testing.T) {
	bc := newBlockchain()
	node := "http://12345"
	bc.registerNode(node)

	if l := len(bc.nodes); l != 1 {
		t.Errorf("blockchain.registerNode expected len to be %v got %v", 1, l)
	}

	if n, ok := bc.nodes[node]; !ok || n != node {
		t.Errorf("blockchain.registerNode expected node to have %v with value %v but got %v", node, node, n)
	}
}
