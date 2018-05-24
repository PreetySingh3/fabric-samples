package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("example_cc0")

// SimpleAsset implements a simple chaincode to manage an asset
type SimpleAsset struct {
}

type MDRMStruct struct {
	Id              string `json:"Id"`
	MDRM            string `json:"MDRM"`
	ItemDescription string `json:"ItemDescription"`
	Amount          string `json:"Amount"`
}

//var credit map[string]interface{}

// Init is called during chaincode instantiation to initialize any
// data. Note that chaincode upgrade also calls this function to reset
// or to migrate data.
/*func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
    // Get the args from the transaction proposal
    args := stub.GetStringArgs()
    if len(args) != 2 {
            return shim.Error("Incorrect arguments. Expecting a key and a value")
    }
    // Set up any variables or assets here by calling stub.PutState()
    // We store the key and the value on the ledger
    err := stub.PutState(args[0], []byte(args[1]))
    if err != nil {
            return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
    }
    return shim.Success(nil)
}*/

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	logger.Info("########### example_cc0 Init ###########")

	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	logger.Info("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// Invoke is called per transaction on the chaincode. Each transaction is
// either a 'get' or a 'set' on the asset created by Init function. The Set
// method may create a new asset by specifying a new key-value pair.
func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "createAmountEditRequest" {
		result, err = createAmountEditRequest(stub, args)
	} else {
		result, err = getMDMR(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func createAmountEditRequest(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	MDRMStruct := &MDRMStruct{args[0], args[1], args[2], args[3]}
	MDRMJSONasBytes, err := json.Marshal(MDRMStruct)
	if err != nil {
		return "", fmt.Errorf("Failed to Marshal asset: %s", args[0])
	}

	err = stub.PutState(args[0], []byte(MDRMJSONasBytes))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

func getMDMR(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return args[1], nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
