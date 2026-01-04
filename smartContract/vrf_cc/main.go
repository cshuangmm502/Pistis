package main

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"golang.org/x/crypto/ed25519"
//	"strconv"
//	"time"
//)
//
//const (
//	BLOCK6210776 = "b6e1a46d3a20a98077c775c0256ec267d3e404fc82810b88dad300567dd8e16c"
//	BLOCK6210770 = "b5cfc7e4aa55354ffa2ae5f0f8673e80ec2dc6ac2b4920be25839d5fd099319b"
//	BLOCK6210775 = "274878b1565a77fe8559521513c122e8ef85e1832726ed03646ecd1dc809f3a2"
//)
//
//var Candidates []Candidate
//
//type Candidate struct {
//	Pk   ed25519.PublicKey
//	sk   ed25519.PrivateKey
//	Rand []byte
//	Name string
//	Pi   []byte
//}
//
////var USERS map[string]pairKey
//
//type pairKey struct {
//	pk ed25519.PublicKey
//	sk ed25519.PrivateKey
//}
//
//func init() {
//	for i := 0; i < 100; i++ {
//		pk, sk, _ := ed25519.GenerateKey(nil)
//		r, pi := Evaluate(sk, pk, []byte(BLOCK6210770))
//		candidate := Candidate{
//			Pk:   pk,
//			sk:   sk,
//			Rand: r,
//			Name: "user" + strconv.Itoa(i),
//			Pi:   pi,
//		}
//		Candidates = append(Candidates, candidate)
//	}
//}
//
////func init() {
////	USERS = make(map[string]pairKey)
////	for i := 1; i <= 100; i++ {
////		pk, sk, err := ed25519.GenerateKey(nil)
////		if err != nil {
////			fmt.Println(err)
////			break
////		}
////		pairkey := pairKey{pk: pk, sk: sk}
////		USERS["user"+strconv.Itoa(i)] = pairkey
////	}
////}
//
////func main() {
////	t := time.Now()
////	pk1, sk1, err := ed25519.GenerateKey(nil) //生成私钥和公钥
////	if err != nil {
////		fmt.Println("errors")
////	}
////	fmt.Println("pk1:")
////	fmt.Println(pk1)
////	fmt.Println("sk1:")
////	fmt.Println(sk1)
////	r1, pi1 := Evaluate(sk1, pk1, []byte(BLOCK6210770)) //生成随机数，证明
////	if err != nil {
////		fmt.Println("errors")
////	} else {
////		fmt.Println()
////		fmt.Println("Proof1:")
////		fmt.Println(pi1)
////		fmt.Println("random1:")
////		fmt.Println(r1)
////	}
////	pk2, sk2, err := ed25519.GenerateKey(nil) //生成私钥和公钥
////	if err != nil {
////		fmt.Println("errors")
////	}
////	fmt.Println("pk2:")
////	fmt.Println(pk2)
////	fmt.Println("sk2:")
////	fmt.Println(sk2)
////	r2, pi2 := Evaluate(sk2, pk2, []byte(BLOCK6210775)) //生成随机数，证明
////	if err != nil {
////		fmt.Println("errors")
////	} else {
////		fmt.Println()
////		fmt.Println("Proof2:")
////		fmt.Println(pi2)
////		fmt.Println("random2:")
////		fmt.Println(r2)
////	}
////	//pk3, sk3, err := ed25519.GenerateKey(nil)
////	//r3, _ := Evaluate(sk3, pk3, []byte(BLOCK6210776))
////	fmt.Println("")
////	fmt.Println(hex.EncodeToString(r1))
////	fmt.Println(hex.EncodeToString(r2))
////	fmt.Println(bytes.Compare(r1, r2))
////	t2 := time.Since(t)
////	fmt.Println(t2)
////	res, _ := Verify(pk1, pi1, []byte(BLOCK6210770), r1) //验证结果
////	fmt.Println("Result:")
////	fmt.Println(res)
////
////}
//
//func main() {
//	//fmt.Println(Candidates[0])
//	now := time.Now()
//	quickSort(Candidates, 0, len(Candidates)-1)
//	//for i := 0; i < 100; i++ {
//	//	fmt.Println(Candidates[i].rand)
//	//}
//	fmt.Println(Candidates[0])
//	fmt.Println("***************")
//	str, _ := json.Marshal(Candidates)
//	//fmt.Println(str)
//	var candis []Candidate
//	_ = json.Unmarshal(str, &candis)
//	fmt.Println(candis[0])
//	t := time.Since(now).Microseconds()
//	fmt.Println(t)
//}
//
//func quickSort(arr []Candidate, low, high int) {
//	if low < high {
//		pivotIndex := partition(arr, low, high)
//		quickSort(arr, low, pivotIndex-1)
//		quickSort(arr, pivotIndex+1, high)
//	}
//}
//
//func partition(arr []Candidate, low, high int) int {
//	pivot := arr[high].Rand
//	i := low - 1
//	for j := low; j < high; j++ {
//		if bytes.Compare(arr[j].Rand, pivot) < 0 {
//			i++
//			arr[i], arr[j] = arr[j], arr[i]
//		}
//	}
//	arr[i+1], arr[high] = arr[high], arr[i+1]
//	return i + 1
//}
