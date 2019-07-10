package main

import (
	"bytes"
	"fmt"
	"math/big"
)

var b58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(input []byte) []byte {
	var result []byte //定义一个字节切片，返回值

	x := big.NewInt(0).SetBytes(input) //字节数组input转化为了大整数big.Int

	base := big.NewInt(int64(len(b58Alphabet))) //长度58的大整数
	zero := big.NewInt(0)                       //0的大整数

	mod := &big.Int{} //大整数的指针
	for x.Cmp(zero) != 0 {
		x.DivMod(x, base, mod) // 对x取余数
		result = append(result, b58Alphabet[mod.Int64()])
	}
	ReverseBytes(result)

	for _, b := range input {

		if b == 0x00 {
			result = append([]byte{b58Alphabet[0]}, result...)
		} else {
			break
		}
	}

	return result

}

func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}

func Base58Decode(input []byte) []byte {
	result := big.NewInt(0)
	zeroBytes := 0
	for _, b := range input {
		if b == '1' {
			zeroBytes++
		} else {
			break
		}
	}
	payload := input[zeroBytes:]

	for _, b := range payload {
		charIndex := bytes.IndexByte(b58Alphabet, b) //反推出余数

		result.Mul(result, big.NewInt(58)) //之前的结果乘以58

		result.Add(result, big.NewInt(int64(charIndex))) //加上这个余数

	}
	decoded := result.Bytes()
	decoded = append(bytes.Repeat([]byte{0x00}, zeroBytes), decoded...)
	return decoded
}
func main() {
	org := []byte("qwerty")
	fmt.Println(string(org))
	ReverseBytes(org)
	fmt.Println(string(org))
	fmt.Printf("%s\n", string(Base58Encode([]byte("hello jonson"))))
	fmt.Printf("%s", string(Base58Decode([]byte("2yGEbwRFyav6CimZ7"))))
}
