package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"golang.org/x/crypto/ed25519"
	"strconv"
)

const EntryKeyPrefix = "ELECTION_ENTRY"
const CandidateListPrefix = "ELECTION_CANDIDATES"
const ElectedListPrefix = "ELECTION_ELECTED"
const BLOCK6210770 = "b5cfc7e4aa55354ffa2ae5f0f8673e80ec2dc6ac2b4920be25839d5fd099319b"

type vrfCC struct {
}

type Election struct {
	Epoch       int
	TargetBlock int
	//NodeEntries map[string]Entry
	//WinnerList  []string
	//BlockHash   []string
}

type Entry struct {
	NodeID       string
	RandomNumber string
	Proof        string
	PublicKey    string
}

type Candidate struct {
	Pk   ed25519.PublicKey
	sk   ed25519.PrivateKey
	Rand []byte
	Name string
	Pi   []byte
}

func main() {
	if err := shim.Start(new(vrfCC)); err != nil {
		fmt.Printf("Error starting miniCC chaincode: %s", err)
	}
}

func (t *vrfCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("accumulator state:"))
}

func (t *vrfCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fc, args := stub.GetFunctionAndParameters()
	fmt.Println("invoking is running " + fc)
	if fc == "initLedger" {
		return t.initLedger(stub)
	} else if fc == "startElection" {
		return t.startElection(stub, args)
	} else if fc == "campaign" {
		return t.campaign(stub, args)
	} else if fc == "readCurrentElection" {
		return t.readCurrentElection(stub)
	} else if fc == "public" {
		return t.publicCandidate(stub, args)
	} else if fc == "verify" {
		return t.verify(stub, args)
	}
	return shim.Error("Unknown function invoking")
}

// admin invoke this func to start new epoch committee selection
func (t *vrfCC) startElection(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	currentElectionAsBytes, err := stub.GetState("currentElection")
	if err != nil {
		return shim.Error(err.Error())
	}
	currentElection := &Election{}
	json.Unmarshal(currentElectionAsBytes, currentElection)
	currentEpoch := currentElection.Epoch
	targetBlock, _ := strconv.Atoi(args[0])
	var newElection = Election{
		Epoch:       currentEpoch + 1,
		TargetBlock: targetBlock,
	}
	bytes, err := json.Marshal(newElection)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState("currentElection", bytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.SetEvent("startElection", []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(args[0]))
}

// peer chaincode invoke -C myc -n mycc -c '{"Args":["campaignFor","selfPublickey"]}'
func (t *vrfCC) campaign(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments")
	}
	epoch := args[0]
	nodeID := args[1]
	rand := args[2]
	pi := args[3]
	pk := args[4]
	entry := Entry{NodeID: nodeID, RandomNumber: rand, Proof: pi, PublicKey: pk}
	bytes, _ := json.Marshal(entry)
	electionNodeKey, err := generateElectionNodeKey(stub, epoch, nodeID)
	if err != nil {
		return shim.Error("Failed to generate composite key")
	}
	err = stub.PutState(electionNodeKey, bytes)
	if err != nil {
		fmt.Println("Encounter an error when submit randNumber  ")
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(electionNodeKey))
}

func (t *vrfCC) compareWithCurrentRand() {

}

func (t *vrfCC) initLedger(stub shim.ChaincodeStubInterface) peer.Response {
	var election = Election{
		Epoch:       -1,
		TargetBlock: -1,
	}
	eleAsByte, err := json.Marshal(election)
	if err != nil {
		fmt.Println("Encounter an error when init Ledger")
		return shim.Error(fmt.Sprintf("Encounter an error when init Ledger"))
	}
	_ = stub.PutState("currentElection", eleAsByte)
	_ = stub.PutState("test", []byte("test"))
	return shim.Success([]byte("success"))
}

func (t *vrfCC) readCurrentElection(stub shim.ChaincodeStubInterface) peer.Response {
	electionJson, err := stub.GetState("currentElection")
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(electionJson)
}

func getBlockHashByBlockNumber(stub shim.ChaincodeStub, blockNumber string) peer.Response {
	response := stub.InvokeChaincode("qscc", [][]byte{[]byte("GetBlockByNumber"), []byte("mychannel"), []byte(blockNumber)}, "")
	if response.Status != shim.OK {
		return shim.Error(fmt.Sprintf("Error querying for block number %s: %s", blockNumber, response.Message))
	}
	block := &common.Block{}
	proto.Unmarshal(response.Payload, block)
	blockHash := base64.StdEncoding.EncodeToString(block.Header.DataHash)
	return shim.Success([]byte(blockHash))
}

func generateElectionNodeKey(stub shim.ChaincodeStubInterface, epoch string, nodeID string) (string, error) {
	return stub.CreateCompositeKey(EntryKeyPrefix, []string{epoch, nodeID})
}

func generateCandidateListKey(stub shim.ChaincodeStubInterface, epoch string) (string, error) {
	return stub.CreateCompositeKey(CandidateListPrefix, []string{epoch})
}

func generateElectedListKey(stub shim.ChaincodeStubInterface, epoch string) (string, error) {
	return stub.CreateCompositeKey(ElectedListPrefix, []string{epoch})
}

func (t *vrfCC) test(stub shim.ChaincodeStubInterface) peer.Response {
	test, err := stub.GetState("test")
	if err != nil {
		return shim.Error(err.Error())
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(test)
}

func (t *vrfCC) publicCandidate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	epoch := args[0]
	electedAsBytes := []byte(args[1])
	electedKey, err := generateElectedListKey(stub, epoch)
	if err != nil {
		return shim.Error("Failed to generate composite key")
	}
	_ = stub.PutState(electedKey, electedAsBytes)
	return shim.Success([]byte("success to public elected node and full candidates in this epoch"))
}

//func (t *vrfCC) publicCandidate(stub shim.ChaincodeStubInterface, args []string) peer.Response {
//	var elected []Candidate //所有人
//	var entrylist []Entry   //最终30人
//	epoch := args[0]
//	candis := []byte(args[1])
//	_ = json.Unmarshal([]byte(args[1]), &elected)
//
//	for i := 0; i < 30; i++ {
//		entry := Entry{
//			NodeID:       elected[i].Name,
//			RandomNumber: hex.EncodeToString(elected[i].Rand),
//			Proof:        hex.EncodeToString(elected[i].Pi),
//			PublicKey:    hex.EncodeToString(elected[i].Pk),
//		}
//		entrylist = append(entrylist, entry)
//	}
//	electionElectedKey, err := generateElectedListKey(stub, epoch)
//	electionCandidateKey, err := generateCandidateListKey(stub, epoch)
//	if err != nil {
//		return shim.Error("Failed to generate composite key")
//	}
//	candisAsbytes, _ := json.Marshal(entrylist)
//	_ = stub.PutState(electionElectedKey, candisAsbytes)
//	_ = stub.PutState(electionCandidateKey, candis)
//	return shim.Success([]byte("success to public elected node and full candidates in this epoch"))
//}

func (t *vrfCC) verify(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	//count := 0
	epoch := args[0]
	//var candidatesList []Entry
	electionElectedKey, err := generateElectedListKey(stub, epoch)
	if err != nil {
		return shim.Error("Failed to generate composite key")
	}
	listAsbytes, err := stub.GetState(electionElectedKey)
	if err != nil {
		return shim.Error("Failed to get candidatesList")
	}
	//err = json.Unmarshal(listAsbytes, &candidatesList)
	//if err != nil {
	//	return shim.Error("Failed to unmarshal candidatesList")
	//}
	//for i := 0; i < 30; i++ {
	//	pk, err1 := hex.DecodeString(candidatesList[i].PublicKey)
	//	if err1 != nil {
	//		return shim.Error("Failed to decode publickey")
	//	}
	//	random, err2 := hex.DecodeString(candidatesList[i].RandomNumber)
	//	if err2 != nil {
	//		return shim.Error("Failed to decode random")
	//	}
	//	pi, err3 := hex.DecodeString(candidatesList[i].Proof)
	//	if err3 != nil {
	//		return shim.Error("Failed to decode proof")
	//	}
	//	res, err4 := Verify(pk, pi, []byte(BLOCK6210770), random)
	//	if err4 != nil {
	//		return shim.Error("Failed to verify")
	//	}
	//	return shim.Success([]byte("test1"))
	//	if res {
	//		count++
	//	}
	//}
	return shim.Success(listAsbytes)
}
