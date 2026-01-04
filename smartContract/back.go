package main

import (
	r "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ed25519"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

const (
	BLOCK6210776 = "b6e1a46d3a20a98077c775c0256ec267d3e404fc82810b88dad300567dd8e16c"
	BLOCK6210770 = "b5cfc7e4aa55354ffa2ae5f0f8673e80ec2dc6ac2b4920be25839d5fd099319b"
	BLOCK6210775 = "274878b1565a77fe8559521513c122e8ef85e1832726ed03646ecd1dc809f3a2"
)

var clientPk *rsa.PublicKey
var clientSk *rsa.PrivateKey

type Candidate struct {
	Pk   ed25519.PublicKey
	Sk   ed25519.PrivateKey
	Rand []byte
	Name string
	Pi   []byte
}

func init() {
	privateKey, err := rsa.GenerateKey(r.Reader, 2048)
	if err != nil {
		panic(err)
	}
	clientSk = privateKey
	clientPk = &privateKey.PublicKey
}

func TriggerElection() time.Time {
	log.Printf("T1 is generated")
	start := time.Now()
	return start
}

func vrfOutputToInt(r []byte) *big.Int {
	h := sha256.Sum256(r) // 使用 SHA256，确保是 256-bit 均匀分布
	return new(big.Int).SetBytes(h[:])
}

func testVRFThreshold(probability float64, rounds int) []Candidate {
	var Committee []Candidate
	// 上限：2^256
	maxVal := new(big.Int).Lsh(big.NewInt(1), 256)

	// 计算 threshold = 2^256 * probability
	p := new(big.Float).SetFloat64(probability)
	tFloat := new(big.Float).Mul(new(big.Float).SetInt(maxVal), p)

	threshold := new(big.Int)
	tFloat.Int(threshold) // 向下取整

	selected := 0
	// 模拟多轮 VRF 计算
	for i := 0; i < rounds; i++ {
		// 模拟一个固定的 seed，例如区块哈希
		//hexSeed := BLOCK6210776
		//seedBytes, err := hex.DecodeString(hexSeed)
		//if err != nil {
		//	panic(err)
		//}
		// 生成密钥
		pk, sk, _ := ed25519.GenerateKey(nil)
		r, pi := Evaluate(sk, pk, []byte(BLOCK6210770))

		// 将 r 转为大整数
		val := vrfOutputToInt(r)

		// 如果 VRF 输出 < 阈值，视为中选
		if val.Cmp(threshold) == -1 {
			winner := Candidate{
				Pk:   pk,
				Sk:   sk,
				Rand: r,
				Name: "user" + strconv.Itoa(i+1),
				Pi:   pi,
			}
			selected++
			Committee = append(Committee, winner)
		}
	}

	actualProb := float64(selected) / float64(rounds)
	fmt.Printf("📊 Target: %.2f, Actual: %.4f (%d / %d)\n", probability, actualProb, selected, rounds)
	return Committee
}

func eli_Ver(candidate Candidate, probability float64) bool {
	// 上限：2^256
	maxVal := new(big.Int).Lsh(big.NewInt(1), 256)

	// 计算 threshold = 2^256 * probability
	p := new(big.Float).SetFloat64(probability)
	tFloat := new(big.Float).Mul(new(big.Float).SetInt(maxVal), p)

	threshold := new(big.Int)
	tFloat.Int(threshold) // 向下取整
	isValid, _ := Verify(candidate.Pk, candidate.Pi, []byte(BLOCK6210770), candidate.Rand)
	return isValid
}

func encrpt(m []byte) []byte {
	ciphertext, err := rsa.EncryptPKCS1v15(r.Reader, clientPk, m)
	if err != nil {
		panic(err)
	}
	fmt.Println("Encrypted (base64):", base64.StdEncoding.EncodeToString(ciphertext))
	return ciphertext
}

func decrpt(ciphertext []byte) []byte {
	plaintext, err := rsa.DecryptPKCS1v15(r.Reader, clientSk, ciphertext)
	if err != nil {
		panic(err)
	}
	fmt.Println("Decrypted:", string(plaintext))
	return plaintext
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//start := time.Now()
	//message := []byte("可读性：\n\nRSA公钥和私钥的原始数据通常是二进制格式，对于人类来说直接阅读和理解这些二进制数据是非常困难的。")
	//ciphertext := encrpt(message)
	//decrpt(ciphertext)
	//cost := time.Since(start).Microseconds()
	//fmt.Println(cost)
	//
	//probs := []float64{0.1}
	//start := time.Now()
	//for _, p := range probs {
	committee := testVRFThreshold(0.1, 1000)
	//}
	//cost := time.Since(start).Microseconds()
	//fmt.Println(cost)
	//fmt.Println(committee[0])
	fmt.Println(len(committee))
	res, _ := Verify(committee[0].Pk, committee[0].Pi, []byte(BLOCK6210770), committee[0].Rand)
	fmt.Println(res)
	//cost = time.Since(start).Microseconds()
	start := time.Now()
	for i := 0; i < 30; i++ {
		res, _ := Verify(committee[i].Pk, committee[i].Pi, []byte(BLOCK6210770), committee[i].Rand)
		log.Printf("Pk check completed, result is %f", res)
		res = eli_Ver(committee[i], 0.1)
		log.Printf("Random check completed, result is %f", res)
		log.Printf("The %d time check is completed", i+1)
	}
	cost := time.Since(start).Microseconds()
	fmt.Println(cost)
}
