package main

//
//import (
//	"crypto/sha256"
//	"encoding/hex"
//	"fmt"
//	"math/big"
//	"math/rand"
//	"time"
//)
//
//const (
//	BLOCK6210776 = "b6e1a46d3a20a98077c775c0256ec267d3e404fc82810b88dad300567dd8e16c"
//)
//
//func vrfOutputToInt(r []byte) *big.Int {
//	h := sha256.Sum256(r) // 使用 SHA256，确保是 256-bit 均匀分布
//	return new(big.Int).SetBytes(h[:])
//}
//
//func testVRFThreshold(probability float64, rounds int) {
//	//var Committee []Candidate
//	// 上限：2^256
//	maxVal := new(big.Int).Lsh(big.NewInt(1), 256)
//
//	// 计算 threshold = 2^256 * probability
//	p := new(big.Float).SetFloat64(probability)
//	tFloat := new(big.Float).Mul(new(big.Float).SetInt(maxVal), p)
//
//	threshold := new(big.Int)
//	tFloat.Int(threshold) // 向下取整
//
//	selected := 0
//	// 模拟多个用户的 VRF 计算
//	for i := 0; i < rounds; i++ {
//		// 模拟一个变化的 seed，例如区块哈希
//		// 区块哈希（作为seed）
//		hexSeed := BLOCK6210776
//		seedBytes, err := hex.DecodeString(hexSeed)
//		if err != nil {
//			panic(err)
//		}
//		// 生成密钥
//		pk, sk := GenerateKey()
//
//		r, _ := Evaluate(sk[:], pk[:], seedBytes)
//
//		// 将 r 转为大整数
//		val := vrfOutputToInt(r)
//
//		// 如果 VRF 输出 < 阈值，视为中选
//		if val.Cmp(threshold) == -1 {
//			//winner := Candidate{
//			//	Pk:   pk[:],
//			//	Sk:   sk[:],
//			//	Rand: r,
//			//	Name: "user" + strconv.Itoa(i+1),
//			//	Pi:   pi,
//			//}
//			selected++
//			//Committee = append(Committee, winner)
//		}
//	}
//
//	actualProb := float64(selected) / float64(rounds)
//	fmt.Printf("📊 Target: %.2f, Actual: %.4f (%d / %d)\n", probability, actualProb, selected, rounds)
//	//return Committee
//}
//
//func main() {
//	rand.Seed(time.Now().UnixNano())
//
//	//probs := []float64{0.1}
//	start := time.Now()
//	//for _, p := range probs {
//	//	testVRFThreshold(p, 1000)
//	//}
//	GenerateKey()
//	cost := time.Since(start).Microseconds()
//	fmt.Println(cost)
//}
