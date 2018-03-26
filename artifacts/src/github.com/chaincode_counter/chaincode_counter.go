package main

import (
    "encoding/json"
    _ "errors"
    "fmt"
    "strconv"
    "github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type CounterChaincode struct {
}

type Counter struct {
    Name string `json:"name"`
    Counts uint64 `json:"counts"`
}

const numOfCounters int = 3

func (cc *CounterChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()

	var err error

	err = stub.PutState(args[0], []byte(args[1]))

	if err != nil {
		fmt.Print(err.Error())
	}

    var counters [numOfCounters]Counter
    var countersBytes [numOfCounters][]byte
    counters[0] = Counter{Name: "Officer Worker", Counts: 0}
    counters[1] = Counter{Name: "Home Worker", Counts: 0}
    counters[2] = Counter{Name: "Student", Counts: 0}

    for i := 0; i < len(counters); i++ {
        countersBytes[i], _ = json.Marshal(counters[i])
        err = stub.PutState(strconv.Itoa(i), countersBytes[i])

		if err != nil {
			fmt.Print(err.Error())
		}
    }

    return shim.Success(nil)
}

func (cc *CounterChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	var result []byte
	var err error

    if function == "countUp" {
        result, err = cc.countUp(stub, args)
    } else if function == "query" {
		result, err = cc.getCounters(stub, args)
	}

	if err != nil {
		return shim.Error(err.Error())
	}

    return shim.Success(result)
}

func (cc *CounterChaincode) Query(stub shim.ChaincodeStubInterface) peer.Response {

	function, args := stub.GetFunctionAndParameters()

	var result []byte
	var err error

    if function == "refresh" {
        result, err = cc.getCounters(stub, args)

		if err != nil {
			return shim.Error(err.Error())
		} else {
			return shim.Success([]byte(result))
		}
    }
    
    return shim.Error(err.Error())
}

func (cc *CounterChaincode) countUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    counterId := args[0]
    counterJson, _ := stub.GetState(counterId)
    
    counter := Counter{}
    json.Unmarshal(counterJson, &counter)
    
    counter.Counts++;
    
    counterJson, _ = json.Marshal(counter)
    stub.PutState(counterId, counterJson)
    
    return nil, nil
}

func (cc *CounterChaincode) getCounters(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    var counters [numOfCounters]Counter
    var countersBytes [numOfCounters][]byte
    
    for i := 0; i < len(counters); i++ {
        countersBytes[i], _ = stub.GetState(strconv.Itoa(i))
        
        counters[i] = Counter{}
        json.Unmarshal(countersBytes[i], &counters[i])
    }
    
    return json.Marshal(counters)
}

func main() {
    err := shim.Start(new(CounterChaincode))
    
    if err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}
