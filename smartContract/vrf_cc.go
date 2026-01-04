package main

//
//import (
//	"fmt"
//	"golang.org/x/crypto/ed25519"
//	"strconv"
//)
//
//const (
//	BLOCK6210776 = "b6e1a46d3a20a98077c775c0256ec267d3e404fc82810b88dad300567dd8e16c"
//	BLOCK6210770 = "b5cfc7e4aa55354ffa2ae5f0f8673e80ec2dc6ac2b4920be25839d5fd099319b"
//	BLOCK6210775 = "274878b1565a77fe8559521513c122e8ef85e1832726ed03646ecd1dc809f3a2"
//)
//
//var USERS map[string]pairKey
//
//type pairKey struct {
//	pk ed25519.PublicKey
//	sk ed25519.PrivateKey
//}
//
//func init() {
//	USERS = make(map[string]pairKey)
//	for i := 1; i <= 10; i++ {
//		pk, sk, err := ed25519.GenerateKey(nil)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//		pairkey := pairKey{pk: pk, sk: sk}
//		USERS["user"+strconv.Itoa(i)] = pairkey
//	}
//}
//
//func main() {
//	const message = "message"
//	//start := time.Now()
//	//for user := range USERS {
//	//	//fmt.Println(hex.EncodeToString(USERS[user].pk))
//	//	fmt.Println(USERS[user].pk)
//	//	hexcode := hex.EncodeToString(USERS[user].pk)
//	//	fmt.Println(hexcode)
//	//	fmt.Println(hex.DecodeString(hexcode))
//	//}
//	fmt.Println(USERS["user1"])
//	//end := time.Since(start)
//	//fmt.Println("10个委员会生成的时间：" + end.String())
//}
