package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type eviCC struct {
}

type pos struct {
	X int
	Y int
}

type Evidence struct {
	EvidenceID string `json:"evidenceID"` // 存证ID
	//EvidenceName string `json:"evidenceName"` //存证名称
	SourceUAV   string `json:"sourceUAV"`   //来源无人机
	CollectTime string `json:"collectTime"` //采集时间
	CollectPos  pos    `json:"collectPos"`  //采集位置
	HashAlg     string `json:"hashAlg"`     //哈希算法
	DataHash    string `json:"dataHash"`    //哈希摘要
}

func main() {
	if err := shim.Start(new(eviCC)); err != nil {
		fmt.Printf("Error starting testCC chaincode: %s", err)
	}
}

func (e *eviCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("accumulator state:"))
}

func (e *eviCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fc, args := stub.GetFunctionAndParameters()
	fmt.Println("invoking is running " + fc)
	if fc == "initiate" {
		return e.initSamples(stub)
	} else if fc == "save" {
		return e.saveEvidence(stub, args)
	} else if fc == "read" {
		return e.queryEvidence(stub, args)
	}
	return shim.Error("Unknown function invoking")
}

func (e *eviCC) initSamples(stub shim.ChaincodeStubInterface) peer.Response {
	exampleEvidences := []Evidence{
		{
			EvidenceID: "testEvidenceID1",
			//EvidenceName: "Evidence1",
			SourceUAV:   "UAV1",
			CollectTime: "202402260115",
			CollectPos:  pos{X: 150, Y: 150},
			HashAlg:     "SHA256",
			DataHash:    "lm12lmfo219jfmalwkm12f12",
		},
		{
			EvidenceID: "testEvidenceID2",
			//EvidenceName: "Evidence2",
			SourceUAV:   "UAV2",
			CollectTime: "202402260515",
			CollectPos:  pos{X: 200, Y: 200},
			HashAlg:     "SHA256",
			DataHash:    "faw351351gawfawjlng56afw",
		},
		{
			EvidenceID: "testEvidenceID3",
			//EvidenceName: "Evidence3",
			SourceUAV:   "UAV3",
			CollectTime: "202402260915",
			CollectPos:  pos{X: 300, Y: 300},
			HashAlg:     "SHA256",
			DataHash:    "fawnjk21nfoii2l921nflk12n",
		},
	}
	for i := 0; i < len(exampleEvidences); i++ {
		evidenceAsBytes, _ := json.Marshal(exampleEvidences[i])
		err := stub.PutState(exampleEvidences[i].EvidenceID, evidenceAsBytes)
		if err != nil {
			fmt.Println("Encounter an error when init Evidence"+strconv.Itoa(i), err)
			continue
		}
	}
	return shim.Success([]byte("init ledger success"))
}

func (e *eviCC) saveEvidence(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments.")
	}
	evidenceId := args[0]
	sourceUAV := args[1]
	collectTime := args[2]
	collectPosX := args[3]
	collectPosY := args[4]
	hashAlg := args[5]
	dataHash := args[6]
	posX, _ := strconv.Atoi(collectPosX)
	posY, _ := strconv.Atoi(collectPosY)

	var evidence = Evidence{
		EvidenceID:  evidenceId,
		SourceUAV:   sourceUAV,
		CollectTime: collectTime,
		CollectPos:  pos{X: posX, Y: posY},
		HashAlg:     hashAlg,
		DataHash:    dataHash,
	}

	evidenceAsBytes, _ := json.Marshal(evidence)

	err := stub.PutState(evidenceId, evidenceAsBytes)
	if err != nil {
		fmt.Println("Encounter an error when saving Evidence  ", evidenceId)
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveRequest", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *eviCC) queryEvidence(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting an evidence name")
	}
	key := args[0]
	result, err := stub.GetState(key) //负责查询账本，返回指定键的对应值
	if err != nil {
		fmt.Println("Encounter an error when querying by the evidence name ", key)
		return shim.Error(err.Error())
	}
	if result == nil {
		shim.Error("No information was found on the evidence name.")
	}
	return shim.Success(result)
}
