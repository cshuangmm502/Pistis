package main

//
//import (
//	"bytes"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
//	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
//	"golang.org/x/crypto/ed25519"
//	"io/ioutil"
//	"log"
//	"os"
//	"path/filepath"
//	"strconv"
//	"sync"
//	"time"
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
//var Candidates []Candidate
//
//type Candidate struct {
//	Pk   ed25519.PublicKey
//	Sk   ed25519.PrivateKey
//	Rand []byte
//	Name string
//	Pi   []byte
//}
//
//type Entry struct {
//	NodeID       string
//	RandomNumber string
//	Proof        string
//	PublicKey    string
//}
//
//func init() {
//	USERS = make(map[string]pairKey)
//	for i := 0; i < 400; i++ {
//		pk, sk, err := ed25519.GenerateKey(nil)
//		if err != nil {
//			fmt.Println(err)
//			break
//		}
//		pairkey := pairKey{pk: pk, sk: sk}
//		USERS["user"+strconv.Itoa(i+1)] = pairkey
//	}
//
//	for i := 0; i < 400; i++ {
//		pk, sk, _ := ed25519.GenerateKey(nil)
//		r, pi := Evaluate(sk, pk, []byte(BLOCK6210770))
//		candidate := Candidate{
//			Pk:   pk,
//			Sk:   sk,
//			Rand: r,
//			Name: "user" + strconv.Itoa(i+1),
//			Pi:   pi,
//		}
//		Candidates = append(Candidates, candidate)
//	}
//}
//
//func main() {
//	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
//	log.Println("============ application-golang starts ============")
//
//	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
//	if err != nil {
//		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
//	}
//	wallet, err := gateway.NewFileSystemWallet("wallet")
//	if err != nil {
//		log.Fatalf("Failed to create wallet: %v", err)
//	}
//
//	if !wallet.Exists("appUser") {
//		err = populateWallet(wallet)
//		if err != nil {
//			log.Fatalf("Failed to populate wallet contents: %v", err)
//		}
//	}
//	ccpPath := filepath.Join(
//		"..",
//		"test-network",
//		"organizations",
//		"peerOrganizations",
//		"org1.example.com",
//		"connection-org1.yaml", //fabric network connection file
//	)
//
//	gw, err := gateway.Connect(
//		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
//		gateway.WithIdentity(wallet, "appUser"),
//	)
//	if err != nil {
//		log.Fatalf("Failed to connect to gateway: %v", err)
//	}
//	defer gw.Close()
//
//	network, err := gw.GetNetwork("mychannel")
//	if err != nil {
//		log.Fatalf("Failed to get network: %v", err)
//	}
//
//	contract := network.GetContract("elecc400")
//
//	// Create an event listener for 'startElection' event
//	eventID := "startElection"
//	reg, notifier, err := contract.RegisterEvent(eventID)
//	if err != nil {
//		fmt.Printf("Failed to register event: %v\n", err)
//		return
//	}
//	defer contract.Unregister(reg)
//	initLedger(contract)
//	// Listen for the event in a goroutine
//	//time.Sleep(5 * time.Second)
//	//imitatePub(contract, "1")
//	//time.Sleep(5 * time.Second)
//	counter := 1
//	go func() {
//		for {
//			select {
//			case event := <-notifier:
//				now := time.Now()
//				log.Printf("接收到链码事件: %v,开始竞选阶段\n\n", string(event.Payload))
//				//rand.Seed(int64(time.Now().Nanosecond()))
//				//randomNum := rand.Intn(11) + 9 // 生成10到20之间的随机数
//				//fmt.Println("随机数:", randomNum)
//				var wg sync.WaitGroup
//				for i := 0; i < 400; i++ {
//					wg.Add(1)
//					go func(i int) {
//						defer wg.Done()
//						campaign(contract, "1", "user"+strconv.Itoa(i+1))
//						//t := time.Since(now).Milliseconds()
//						//log.Println("总计时间开销：")
//						//log.Println(t)
//					}(i)
//				}
//				wg.Wait()
//				t1 := time.Since(now).Milliseconds()
//				log.Printf("投票阶段结束,耗时%v", t1) //t1
//				log.Println("进入排序验证阶段")
//				verifyandSorting()
//				sortTime := time.Since(now).Milliseconds() - t1
//				log.Printf("排序验证阶段结束,耗时%v", sortTime) //t1
//				log.Println("进入提交阶段")
//				submit(contract, "1")
//				subTime := time.Since(now).Milliseconds() - t1 - sortTime
//				log.Printf("提交阶段结束,耗时%v", subTime) //t1
//				totalTime := time.Since(now).Milliseconds()
//				log.Printf("本轮选举结束，总耗时%v", totalTime)
//				saveElectionTime(counter, t1, sortTime, subTime, totalTime)
//				counter++
//				time.Sleep(5 * time.Second)
//				//fmt.Printf("接收到链码事件: %v\n", event)
//				//case <-time.After(time.Second * 20):
//				//	log.Printf("不能根据指定的事件ID接收到链码事件(%s)\n", eventID)
//			case <-time.After(time.Second * 20):
//				log.Printf("20秒内未接收到链码事件(%s)，仍在监听...\n", eventID)
//			}
//		}
//	}()
//	log.Println("============ Continuous monitoring election events ============")
//	for {
//	}
//	log.Println("============ application-golang ends ============")
//}
//
//func saveElectionTime(count int, voteTime, sortTime, subTime, totalTime int64) {
//	// 打开文件，如果不存在则创建
//	file, err := os.OpenFile("election400_30.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//
//	// 使用fmt.Fprintf在新的一行中写入这些数
//	_, err = fmt.Fprintf(file, "%d %d %d %d %d\n", count, voteTime, sortTime, subTime, totalTime)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("Content written to file")
//}
//
//func populateWallet(wallet *gateway.Wallet) error {
//	log.Println("============ Populating wallet ============")
//	credPath := filepath.Join(
//		"..",
//		"test-network",
//		"organizations",
//		"peerOrganizations",
//		"org1.example.com",
//		"users",
//		"User1@org1.example.com",
//		"msp",
//	)
//
//	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
//	// read the certificate pem
//	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
//	if err != nil {
//		return err
//	}
//
//	keyDir := filepath.Join(credPath, "keystore")
//	// there's a single file in this dir containing the private key
//	files, err := ioutil.ReadDir(keyDir)
//	if err != nil {
//		return err
//	}
//	if len(files) != 1 {
//		return fmt.Errorf("keystore folder should have contain one file")
//	}
//	keyPath := filepath.Join(keyDir, files[0].Name())
//	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
//	if err != nil {
//		return err
//	}
//
//	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))
//
//	return wallet.Put("appUser", identity)
//}
//
//func initLedger(contract *gateway.Contract) {
//	log.Println("--> Simultaneously submit Transaction")
//	log.Println("--> Submit Transaction: InitLedger, function creates the initial set of assets on the ledger")
//	result, err := contract.SubmitTransaction("initLedger")
//	if err != nil {
//		log.Fatalf("Failed to Submit transaction: %v", err)
//	}
//	log.Println(string(result))
//}
//
//func startElection(contract *gateway.Contract, epoch string, targetBlock string) {
//	log.Println("--> Evaluate Transaction: startElection, this function triggers an election event")
//	result, err := contract.SubmitTransaction("startElection", epoch, targetBlock)
//	if err != nil {
//		log.Fatalf("Failed to evaluate transaction: %v", err)
//	}
//	log.Println(string(result))
//}
//
//func readElection(contract *gateway.Contract) {
//	log.Println("--> Evaluate Transaction: ReadElection, function returns epoch and target block of currentElection")
//	result, err := contract.EvaluateTransaction("readCurrentElection")
//	if err != nil {
//		log.Fatalf("Failed to evaluate transaction: %v\n", err)
//	}
//	log.Println(string(result))
//}
//
//func campaign(contract *gateway.Contract, epoch string, nodeId string) {
//	log.Println("--> Evaluate Transaction: campaign, this function will generate composite key and public random number on chain")
//	nowtime := time.Now()
//	sk1 := USERS[nodeId].sk
//	pk1 := USERS[nodeId].pk
//	r1, pi1 := Evaluate(sk1, pk1, []byte(BLOCK6210770)) //生成随机数，证明
//	endtime := time.Since(nowtime).Milliseconds()
//	log.Printf("节点%s的选票生成耗时为%d", nodeId, endtime)
//	result, err := contract.SubmitTransaction("campaign", epoch, nodeId, hex.EncodeToString(r1), hex.EncodeToString(pk1), hex.EncodeToString(pi1))
//	if err != nil {
//		log.Fatalf("Failed to evaluate transaction: %v", err)
//	}
//	log.Println(string(result))
//}
//
//func listenToEvents(contract *gateway.Contract) {
//	// Create an event listener for 'startElection' event
//	eventID := "startElection"
//	reg, notifier, err := contract.RegisterEvent(eventID)
//	if err != nil {
//		fmt.Printf("Failed to register event: %v\n", err)
//		return
//	}
//	defer contract.Unregister(reg)
//
//	// Listen for the event in a goroutine
//	go func() {
//		event, ok := <-notifier
//		if !ok {
//			return
//		}
//		fmt.Printf("接收到链码事件: %v\n", event)
//	}()
//
//}
//
//func test(contract *gateway.Contract) {
//	log.Println("--> Evaluate Transaction: ReadElection, function returns epoch and target block of currentElection")
//	result, err := contract.EvaluateTransaction("test")
//	if err != nil {
//		log.Fatalf("Failed to evaluate transaction: %v\n", err)
//	}
//	log.Println(string(result))
//}
//
////验证并排序-链下行为
//func verifyandSorting() {
//	log.Println("--> 开始验证选票合法性.")
//	startTime := time.Now()
//	verifyCandidate()
//	vrDua := time.Since(startTime).Milliseconds()
//	log.Printf("--> 验证选票合法性完成，耗时%dms.", vrDua) //t2
//	log.Println("--> 开始对选票进行排序.")
//	stStartTime := time.Now()
//	quickSort(Candidates, 0, len(Candidates)-1)
//	stDua := time.Since(stStartTime).Milliseconds()
//	log.Printf("--> 选票排序完成，耗时%dms.", stDua) //t3
//	//totalTime := time.Since(startTime).Milliseconds()
//	//log.Printf("链下验证与排序完成，耗时%dms.", totalTime)	//t2+t3
//
//}
//
//func submit(contract *gateway.Contract, epoch string) {
//	log.Println("--> 开始提交最终名单.")
//	sbstartTime := time.Now()
//	var electedList []Entry
//	for i := 0; i < 30; i++ {
//		entry := Entry{NodeID: Candidates[i].Name, RandomNumber: hex.EncodeToString(Candidates[i].Rand), Proof: hex.EncodeToString(Candidates[i].Pi),
//			PublicKey: hex.EncodeToString(Candidates[i].Pk)}
//		electedList = append(electedList, entry)
//	}
//	str, err := json.Marshal(electedList)
//	if err != nil {
//		log.Fatalf("Failed to Marshal: %v\n", err)
//	}
//	result, err := contract.SubmitTransaction("public", epoch, string(str))
//	if err != nil {
//		log.Fatalf("Failed to invoke transaction: %v", err)
//	}
//	sbEndTime := time.Since(sbstartTime).Milliseconds()
//	log.Printf("名单上链完成，耗时%dms.", sbEndTime) //t4
//	log.Println(string(result))
//}
//
//func verifyCandidate() {
//	log.Println("--> 开始验证选票合法性:")
//	nowtime := time.Now()
//	count := 0
//	//log.Println(string(result))
//	log.Println(Candidates[0])
//	for i := 0; i < 400; i++ {
//		res, err := Verify(Candidates[i].Pk, Candidates[i].Pi, []byte(BLOCK6210770), Candidates[i].Rand)
//		if err != nil {
//			log.Println("Failed to verify")
//		}
//		if res {
//			count++
//		}
//	}
//	endtime := time.Since(nowtime).Milliseconds()
//	log.Println(count)
//	log.Println(endtime)
//	//log.Println(string(result))
//}
//
//func verifyElecte(contract *gateway.Contract, epoch string) {
//	log.Println("--> Evaluate Transaction: Verify, function check the validity of elected nodes' pk, pi and random")
//	nowtime := time.Now()
//	count := 0
//	result, err := contract.EvaluateTransaction("verify", epoch)
//	if err != nil {
//		log.Fatalf("Failed to evaluate transaction: %v", err)
//	}
//	var candidatesList []Entry
//	err = json.Unmarshal(result, &candidatesList)
//	if err != nil {
//		log.Println("Failed to unmarshal candidatesList")
//	}
//	//log.Println(string(result))
//	log.Println(candidatesList[0])
//	for i := 0; i < 400; i++ {
//		pk, err1 := hex.DecodeString(candidatesList[i].PublicKey)
//		if err1 != nil {
//			log.Println("Failed to decode publickey")
//		}
//		random, err2 := hex.DecodeString(candidatesList[i].RandomNumber)
//		if err2 != nil {
//			log.Println("Failed to decode random")
//		}
//		pi, err3 := hex.DecodeString(candidatesList[i].Proof)
//		if err3 != nil {
//			log.Println("Failed to decode proof")
//		}
//		res, err4 := Verify(pk, pi, []byte(BLOCK6210770), random)
//		if err4 != nil {
//			log.Println("Failed to verify")
//		}
//		if res {
//			count++
//		}
//	}
//	endtime := time.Since(nowtime).Milliseconds()
//	log.Println(count)
//	log.Println(endtime)
//	//log.Println(string(result))
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
