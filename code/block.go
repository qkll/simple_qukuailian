// block
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Block struct {
	Index     int64  `json:"Index"`
	Timestamp int64  `json:"Time"`
	Hash      string `json:"Hash"`
	Pre_hash  string `json:"Pre_hash"`
	Data      string `json:"Data"`
}

func calulate_hash(b Block) string { //计算hash
	block_hash := string(b.Index) + string(b.Timestamp) + b.Pre_hash + b.Data
	hashinbyte := sha256.Sum256([]byte(block_hash))
	return hex.EncodeToString(hashinbyte[:])
}
func Create_block(preblock Block, data string) Block { //生成新的区块
	b := Block{}
	b.Index = preblock.Index + 1
	b.Pre_hash = preblock.Hash
	b.Data = data
	b.Timestamp = time.Now().Unix()
	b.Hash = calulate_hash(b)
	return b
}
func Genesis_block() Block { //创世区块
	preblock := Block{}
	preblock.Timestamp = time.Now().Unix()
	preblock.Data = "Before Genesis Block"
	data := "Genesis Block"
	b := Create_block(preblock, data)
	return b
}
