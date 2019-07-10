// rpc
package main

import (
	"encoding/json"
	"io"
	"net/http"
)

var blockchain *Blockchain

func run() {
	http.HandleFunc("/blockchain/get", blockchainGet)
	http.HandleFunc("/blockchain/write", blockchainWrite)
	http.ListenAndServe("127.0.0.1:8888", nil)
}
func blockchainGet(w http.ResponseWriter, r *http.Request) {
	bytes, error := json.Marshal(blockchain.blocks)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
	blockchain.Print()
}
func blockchainWrite(w http.ResponseWriter, r *http.Request) {
	data := r.URL.Query()
	blockdata := data.Get("data")
	blockchain.Senddata(blockdata)
	blockchainGet(w, r)
}
