package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type reCC struct {
}

type errRecord struct {
	RecordId  string
	SourceUAV string
	Content   string
}

func main() {
	if err := shim.Start(new(reCC)); err != nil {
		fmt.Printf("Error starting testCC chaincode: %s", err)
	}
}

func (e *reCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("Init ledger successful!"))
}

func (e *reCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fc, args := stub.GetFunctionAndParameters()
	fmt.Println("invoking is running " + fc)
	if fc == "save" {
		return e.saveRecord(stub, args)
	} else if fc == "read" {
		return e.readRecord(stub, args)
	}
	return shim.Error("Unknown function invoking")
}

func (e *reCC) saveRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}
	recordId := args[0]
	sourceUAV := args[1]
	content := args[2]
	var record = errRecord{
		RecordId:  recordId,
		SourceUAV: sourceUAV,
		Content:   content,
	}

	recordAsBytes, _ := json.Marshal(record)
	err := stub.PutState(recordId, recordAsBytes)
	if err != nil {
		fmt.Println("Encounter an error when saving record  ", recordId)
		return shim.Error(err.Error())
	}

	err = stub.SetEvent("saveRequest", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (e *reCC) readRecord(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting a recordId.")
	}
	key := args[0]
	result, err := stub.GetState(key) //负责查询账本，返回指定键的对应值
	if err != nil {
		fmt.Println("Encounter an error when querying by the recordId ", key)
		return shim.Error(err.Error())
	}
	if result == nil {
		shim.Error("No information was found by the recordId.")
	}
	return shim.Success(result)
}
