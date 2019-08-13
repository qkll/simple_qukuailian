package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"math/big"
	"os"
	"strings"
)

func main() {
	//0 - 有一个私有的ECDSA键
	fmt.Println("0、Have a private ECDSA button")
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	fmt.Println("private key：" + Tool_DecimalByteSlice2HexString(private.D.Bytes()))

	//1 - 生成相应的公钥
	fmt.Println("1、Generate the corresponding public key")
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	fmt.Println("public key：" + Tool_DecimalByteSlice2HexString(pubKey))
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}

	//2 - 将公钥进行sha256
	fmt.Println("2、The public key makes sha256")
	publicSHA256 := sha256.Sum256(pubKey)
	fmt.Println("这是公钥hash：" + Tool_DecimalByteSlice2HexString(publicSHA256[:]))

	//crypto.RIPEMD160.HashFunc().New()
	RIPMD160Hasher := ripemd160.New()

	//3 - 在SHA-256的结果上执行RIPEMD-160哈希。
	fmt.Println("3、Perform a RIPEMD-160 hash on the SHA-256 result.")
	_, err = RIPMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		fmt.Errorf(err.Error())
		os.Exit(1)
	}

	//4 - 在RIPEMD-160散列前添加版本字节(主网络的0x00)
	fmt.Println("4、Add version byte before the RIPEMD-160 hash (0x00 of the primary network)")
	var by [1]byte
	by[0] = 0x00
	publicRIPEMD160 := RIPMD160Hasher.Sum(by[:])
	fmt.Println("RIPEMD-160 hash：" + Tool_DecimalByteSlice2HexString(publicRIPEMD160))

	//注意下面的步骤是Base58Check编码，它有多个库选项可用来实现它。

	//5、在扩展的RIPEMD-160结果上执行SHA-256散列。
	fmt.Println("5、Perform SHA-256 hash on extended RIPEMD-160 results。")
	rehash := sha256.Sum256(publicRIPEMD160)
	fmt.Println("first hash：" + Tool_DecimalByteSlice2HexString(rehash[:]))

	//6、对之前的SHA-256散列的结果执行SHA-256散列。
	fmt.Println("6、Perform SHA-256 hash on extended RIPEMD-160 results。")
	rerehash := sha256.Sum256(rehash[:])
	fmt.Println("second hash：" + Tool_DecimalByteSlice2HexString(rerehash[:]))

	//7、以第二个SHA-256散列的前4个字节为例。这是地址校验和。
	fmt.Println("7、Take the first 4 bytes of the second SHA-256 hash as an example. This is the address checksum。")
	sum := checksum(rerehash[:])
	fmt.Println("check address：" + Tool_DecimalByteSlice2HexString(sum))

	//8、将第7阶段的4个校验和字节添加到第4阶段扩展的RIPEMD-160散列的末尾。这是25字节的二进制比特币地址。
	fmt.Println("8、Add the 4 checksum bytes of phase 7 to the end of the RIPEMD-160 hash of the 4th stage extension. This is a 25 byte binary bitcoin address.")
	var b bytes.Buffer
	b.Write(publicRIPEMD160[:])
	b.Write(sum)
	result := b.Bytes()
	fmt.Println("Stitch check address：" + Tool_DecimalByteSlice2HexString(result))

	//9、使用Base58Check编码将一个字节字符串的结果转换为base58字符串。这是最常用的比特币地址格式。
	fmt.Println("9、Convert the result of a byte string to a base58 string using Base58Check encoding. This is the most commonly used bitcoin address format.")
	address := Encode(result)

	fmt.Println("Final address：" + address)
}

func checksum(payload []byte) []byte {
	addressChecksumLen := 4
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func Tool_DecimalByteSlice2HexString(DecimalSlice []byte) string {
	var sa = make([]string, 0)
	for _, v := range DecimalSlice {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	ss := strings.Join(sa, "")
	return ss
}

const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// EncodeBig encodes src, appending to dst. Be sure to use the returned
// new value of dst.
func EncodeBig(dst []byte, src *big.Int) []byte {
	start := len(dst)
	n := new(big.Int)
	n.Set(src)
	radix := big.NewInt(58)
	zero := big.NewInt(0)

	for n.Cmp(zero) > 0 {
		mod := new(big.Int)
		n.DivMod(n, radix, mod)
		dst = append(dst, alphabet[mod.Int64()])
	}

	for i, j := start, len(dst)-1; i < j; i, j = i+1, j-1 {
		dst[i], dst[j] = dst[j], dst[i]
	}
	return dst
}

func Encode(encoded []byte) string {
	//Perform SHA-256 twice
	hash := sha256.Sum256(encoded)
	hash = sha256.Sum256(hash[:])

	//First 4 bytes if this double-sha'd byte array is the checksum
	//Append this checksum to the input bytes
	encoded = append(encoded, hash[0:4]...)

	//Convert this checksum'd version to a big Int
	bigIntEncodedChecksum := new(big.Int).SetBytes(encoded)

	//Encode the big int checksum'd version into a Base58Checked string
	base58EncodedChecksum := EncodeBig(nil, bigIntEncodedChecksum)

	//Now for each zero byte we counted above we need to prepend a 1 to our
	//base58 encoded string. The rational behind this is that base58 removes 0's (0x00).
	//So bitcoin demands we add leading 0s back on as 1s.
	buffer := make([]byte, 0, len(base58EncodedChecksum))

	//base58 alone is not enough. We need to first count each of the zero bytes
	//which are at the beginning of the encodedCheckSum

	for _, v := range encoded {
		if v != 0 {
			break
		}
		buffer = append(buffer, '1')
	}
	buffer = append(buffer, base58EncodedChecksum...)
	return string(buffer)
}
