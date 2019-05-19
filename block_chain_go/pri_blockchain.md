<p style="background-color:#00b0f0; padding: 5px; margin: 3px auto; widows:3; orphans:3;"><span style="color:#ffffff;font-size: 24px; font-weight: bold;">任务：Go语言实现区块链（一）</span></p>


###  一、 任务描述
通过前面的Go语言的学习，对Go语言的有一定的了解，现在我们进行实战演练，通过本实验你将使用Go语言开发自己的区块链(或者说用go语言搭建区块链)、理解哈希函数是如何保持区块链的完整性、掌握如何用Go语言编程创造并添加新的块、实现多个节点通过竞争生成块、通过浏览器来查看整个链、了解所有其他关于区块链的基础知识。但是，本实验中将暂时不会涉及工作量证明算法（PoW）以及权益证明算法（PoS）这类的共识算法，同时为了让你更清楚得查看区块链以及块的添加，我们将网络交互的过程简化了。
###  二、 任务目标
掌握区块链数据结构模型，熟悉Go语言编程模式，实战演练检验前面知识的掌握。

###  三、 任务环境
一台Ubuntu 16.04
主机登录名：root 密码：Simplexue123
操作机：192.168.1.1

实验环境介绍：
一台Ubuntu 16.04主机，已经安装好Go语言环境。

###  四、任务实施
#### 1、实战演练
1.1 创建项目，已经设置了`GOPATH=/home/go`,所以我们就以GOPATH作为我们的项目目录。在此目录下，首先创建文件夹blockchain.创建三个文件分别为block.go、main.go、blockchain.go。block.go文件存放区块结构以及对应函数，main.go为主文件，blockchain.go存放链结构以及对应函数，如图一所示。
![](pic/1.png)
<center>图一 文件夹blockchain机构</center>

1.2 打开文件block.go，开始编写区块结构以及对应函数。导入依赖包。

```
// block
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)
```
1.3 创建区块数据模型，定义一个结构体block,存放区块数据。

```
type Block struct {
	Index     int64  
	Timestamp int64  
	Hash      string 
	Pre_hash  string 
	Data      string 
}
```

- Index 是这个块在整个链中的位置
- Timestamp 显而易见就是块生成时的时间戳
- Hash 是这个块通过 SHA256 算法生成的散列值
- Pre_hash 代表前一个块的 SHA256 散列值
- Data 代码此区块数据

我们使用散列算法（SHA256）来确定和维护链中块和块正确的顺序，确保每一个块的 Pre_Hash值等于前一个块中的 Hash 值，这样就以正确的块顺序构建出链,如图二所示。

![](pic/2.png)
<center>图二 区块链结构演示</center>

1.4 Hash值和新生成的块

我们为什么需要Hash？主要是两个原因：

1. 在节省空间的前提下去唯一标识数据。散列是用整个块的数据计算得出，在我们的例子中，将整个块的数据通过 SHA256 计算成一个定长不可伪造的字符串。
2. 维持链的完整性。通过存储前一个块的散列值，我们就能够确保每个块在链中的正确顺序。任何对数据的篡改都将改变散列值，同时也就破坏了链。

我们接着写一个函数calulate_hash，用来计算给定的数据的 SHA256 散列值：
```
func calulate_hash(b Block) string { //计算hash
	block_hash := string(b.Index) + string(b.Timestamp) + b.Pre_hash + b.Data
	hashinbyte := sha256.Sum256([]byte(block_hash))
	return hex.EncodeToString(hashinbyte[:])
}
```
需要用到string()将块数据转换为string类型，然后利用crypto/sha256包中的sha256加密函数进行Hash计算，需要注意的是，sha256.Sum256()只接受byte类型，最后转换成string类型。

1.5 生成块的函数Create_block。
```
func Create_block(preblock Block, data string) Block { //生成新的区块
	b := Block{}
	b.Index = preblock.Index + 1
	b.Pre_hash = preblock.Hash
	b.Data = data
	b.Timestamp = time.Now().Unix()
	b.Hash = calulate_hash(b)
	return b
}
```
其中，Index 是从给定的前一块的 Index 递增得出，时间戳是直接通过 time.Now() 函数来获得的，Hash 值通过前面的calulate_hash 函数计算得出，Pre_Hash则是给定的前一个块的 Hash 值。

1.6 有个计算Hash函数，创造块的函数，那么可以构造出基本的区块的，一般链的第一块叫做创世块，可以使用Create_block进行创造，不过在这里我们为了方便，我们多增加一个函数Genesis_block，创造创世块。
```
func Genesis_block() Block { //创世区块
	preblock := Block{}
	preblock.Timestamp = time.Now().Unix()
	preblock.Data = "Before Genesis Block"
	data := "Genesis Block"
	b := Create_block(preblock, data)
	return b
}
```
所有数据手动定义好，可以直接调用即可以生成创世块。

1.7 区块数据结构大致了解完了，那么接下来我们介绍链结构，将一个一个区块串成区块链，其实类似于链表的结构。打开文件blockchain.go，引用包设置,定义一个结构体存放Slice作为链表结构使用。
```
package main

import (
	"fmt"
	"log"
)

type Blockchain struct {
	blocks []*Block
}
```
1.8 首先放入区块链中需要确定这个区块是否合法，所以这个时候需要增加一个函数Islegal用来判断区块是否合法。
```
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
```
主要是判断Index是否按照顺序，Hash值是否正确。

1.9 一旦生成了一个区块，就需要将区块加入到区块链中，需要使用函数Append_chain。
```
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
```
核心在于使用Go自带append函数将符合要求的结构体区块加入到Slice链中。

1.10 创建区块需要传递数据进入区块，现在设置一个函数Senddata用来传递数据。
```
func (bc *Blockchain) Senddata(data string) {
	preblock := bc.blocks[len(bc.blocks)-1]
	newblock := Create_block(*preblock, data)
	bc.Append_chain(&newblock)
}
```

1.11 然后开始生成一个新的区块链，并且这个链需要初始化。利用函数Newblockchain实现。
```
func Newblockchain() *Blockchain {
	firstblock := Genesis_block()
	bc := Blockchain{}
	bc.Append_chain(&firstblock)
	return &bc
}
```
首先创造创世区块，然后初始化区块链，将创世区块加入到初始化好的区块链中，这样一个区块链就生成好了。

1.12 区块链生成了，我们可以设置一个输出函数Print，将整个区块链输出来看看。
```
func (bc *Blockchain) Print() {
	for _, block := range bc.blocks {
		fmt.Println("Index:", block.Index)
		fmt.Println("Timestamp:", block.Timestamp)
		fmt.Println("Data:", block.Data)
		fmt.Println("Hash:", block.Hash)
		fmt.Println("Pre_Hash:", block.Pre_hash)
	}
}
```
使用range遍历区块链，一个一个区块的输出即可。
至此区块链的主要函数都已经介绍完毕了。

1.13 下面来介绍主函数main。打开文件main.go。
```
// main
package main

func main() {
	blockchain := Newblockchain()
	data1 := "This is the first block's data"
	data2 := "This is the second block's data"
	blockchain.Senddata(data1)
	blockchain.Senddata(data2)
	blockchain.Print()
}

```
main函数机构非常简单，就是调用函数，设置了block 2的data为“This is the first block's data”，block 3 的data为“This is the second block's data”，加上创世区块，总共三个区块，这样一个简单的区块链实现就完成，下面我们来运行测试一下。
使用go命令进行编译运行`go run main.go`。

如果出现错误
```
./main.go:6:16: undefined: Newblockchain
```

请使用
```
go build .
./blockchain
```
运行结果如图三所示。
![](pic/3.png)
<center>图三 运行结果</center>