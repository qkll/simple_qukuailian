package main

import (
	"fmt"
	"log"
)

type Blockchain struct {
	blocks []*Block
}

func Newblockchain() *Blockchain {
	firstblock := Genesis_block()
	bc := Blockchain{}
	bc.Append_chain(&firstblock)
	return &bc
}
func (bc *Blockchain) Senddata(data string) {
	preblock := bc.blocks[len(bc.blocks)-1]
	newblock := Create_block(*preblock, data)
	bc.Append_chain(&newblock)
}
func (bc *Blockchain) Append_chain(b *Block) { //区块加入区块链
	if len(bc.blocks) == 0 {
		bc.blocks = append(bc.blocks, b)
		return
	}
	if Islegal(*bc.blocks[len(bc.blocks)-1], *b) {
		bc.blocks = append(bc.blocks, b)
	} else {
		log.Fatal("invalid block!")
	}
}
func (bc *Blockchain) Print() {
	for _, block := range bc.blocks {
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Data:", block.Data)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("Pre_Hash:", block.Pre_hash)
	}
}
func Islegal(preblock Block, b Block) bool { //判断区块合法
	if b.Index-1 != preblock.Index {
		return false
	}
	if b.Pre_hash != preblock.Hash {
		return false
	}
	if calulate_hash(b) != b.Hash {
		return false
	}
	return true
}
