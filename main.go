package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	bc             *blockchain
	nodeIdentifier = "Soemting----1"
)

func mine(w http.ResponseWriter, r *http.Request) {
	lastProof := bc.lastBlock().Index

	// We run the proof of work algorithm to get the next proof...
	proof := bc.proofOfWork(lastProof)

	// We must receive a reward for finding the proof.
	// The sender is "0" to signify that this node has mined a new coin.
	transc := transaction{
		Sender:    "0",
		Recipient: nodeIdentifier,
		Amount:    1,
	}
	bc.newTransaction(&transc)

	// Forge the new Block by adding it to the chain
	block := bc.newBlock(proof)

	resp := struct {
		Message      string         `json:"message"`
		Index        int            `json:"index"`
		Transactions []*transaction `json:"transactions"`
		Proof        int            `json:"proof"`
		PreviousHash string         `json:"previous_hash"`
	}{
		Message:      "New Block Forged",
		Index:        block.Index,
		Transactions: block.Transactions,
		Proof:        block.Proof,
		PreviousHash: block.PreviousHash,
	}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding mining payload: ", err)
	}
}

func fullChain(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Length int      `json:"length"`
		Chain  []*block `json:"chain"`
	}{
		Length: len(bc.chain),
		Chain:  bc.chain,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding fullChain payload: ", err)
	}
}

func newTransaction(w http.ResponseWriter, r *http.Request) {
	transc := &transaction{}
	err := json.NewDecoder(r.Body).Decode(transc)

	if err != nil {
		log.Println("Error decoding transaction payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	} else if transc.Amount == 0 || transc.Recipient == "" || transc.Sender == "" {
		http.Error(w, "Missing values", http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	index := bc.newTransaction(transc)
	resp := []byte(fmt.Sprintf(`{"message": "Transaction will be added to block %d"}`, index))
	w.Write(resp)
}

func registerNode(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Nodes []string
	}{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("Error decoding register payload: ", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	for _, node := range payload.Nodes {
		bc.registerNode(node)
	}

	w.Header().Set("Content-Type", "application/json")
	keys := make([]string, len(bc.nodes))
	for k := range bc.nodes {
		keys = append(keys, k)
	}

	response := struct {
		Message    string   `json:"message"`
		TotalNodes []string `json:"total_nodes"`
	}{
		Message:    "New nodes have been added",
		TotalNodes: keys,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("Error encoding register response: ", err)
	}
}

func consensus(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Message string   `json:"message"`
		Chain   []*block `json:"chain"`
	}{}

	if bc.resolveConflict() {
		resp.Message = "Our chain was replaced"
	} else {
		resp.Message = "Our chain is authoritative"
	}

	resp.Chain = bc.chain
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Error encoding consensus response: ", err)
	}
}

func main() {
	bc = newBlockchain()

	http.HandleFunc("/mine", mine)
	http.HandleFunc("/transactions/new", newTransaction)
	http.HandleFunc("/chain", fullChain)
	http.HandleFunc("/nodes/register", registerNode)
	http.HandleFunc("/nodes/resolve", consensus)

	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		panic(err)
	}
}
