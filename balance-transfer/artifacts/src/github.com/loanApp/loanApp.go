package main

import (
	"bytes"
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

type user struct {
	id        string `json:"id"`
	username  string `json:"username"`
	password  string `json:"password"`
	firstname string `json:"firstname"`
	lastname  string `json:"lastname"`
	role      string `json:"role"`
}

type loanApplication struct {
	UserId      string `json:"UserId"`
	Name        string `json:"Name"`
	SSN         string `json:"SSN"`
	LoanAmount  string `json:"LoanAmount"`
	Education   string `json:"Education"`
	Age         string `json:"Age"`
	Tenure      string `json:"Tenure"`
	Address     string `json:"Address"`
	BankId      string `json:"BankId"`
	Status      string `json:"Status"`
	CreditScore string `json:"CreditScore"`
}

type credit struct {
	UserId      string `json:"UserId"`
	CreditScore string `json:"CreditScore"`
}

type creditScoreStruct struct {
	CreditScore []credit
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
	fmt.Println("Args..............................", args)
	fmt.Println("fn..............................", fn)
	var result string
	var err error
	if fn == "createLoanRequest" {
		result, err = createLoanRequest(stub, args)
	} else if fn == "getLoanOfUser" { // assume 'get' even if fn is nil
		result, err = getLoanOfUser(stub, args)
	} else if fn == "queryLoanByBank" {
		return t.queryLoanByBank(stub, args)
	} else if fn == "setCreditScoreState" { // assume 'set' even if fn is nil

		result, err = setCreditScoreState(stub, args)
	} else {
		result, err = updateLoanStatus(stub, args)
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	// Return the result as success payload
	return shim.Success([]byte(result))
}

// Set credit score of bank
func setCreditScoreState(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	//var key, value string
	fmt.Println("testing.......................................", args)

	//credit := creditScoreStruct{}
	bytes1 := []byte(args[0])
	fmt.Println("bytes1................", bytes1)
	//json.Unmarshal(bytes1, &credit)
	var cs creditScoreStruct
	err := json.Unmarshal(bytes1, &cs)

	//	for k, v := range credit {
	//		fmt.Printf("key[%s] value[%s]\n", k, v)
	//	}

	fmt.Println("errr................", err)
	fmt.Println("credit................", cs)

	/*for i := 0; i < len(args); i++ {
		fmt.Println("for loop................", args[i])

		fmt.Printf("Id = %v, Name = %v", credit.UserId, credit.CreditScore)

	}*/

	//fmt.Println("err........................", err)
	//fmt.Println("UserId........................", credit)
	// String contains two JSON rows.
	/*text := "[{\"UserId\": \"L001\", \"CreditScore\": \"Go\"}]"
	// Get byte slice from string.
	bytes := []byte(text)

	// Unmarshal string into structs.
	var credit []creditScoreStruct
	json.Unmarshal(bytes, &credit)

	for l := range credit {
		fmt.Printf("Id = %v, Name = %v", credit[l].UserId, credit[l].CreditScore)
		fmt.Println()
	}*/

	/*for i := 0; i < len(args); i++ {
		fmt.Println("args[i] inside loop..............", args[i])
		key = args[i].UserId
		value = args[i].CreditScore
	}*/

	//loanAsBytes, _ := stub.GetState(key)
	//loanApplication := loanApplication{}

	/*json.Unmarshal(bytes1, &credit)
	loanApplication.CreditScore = value

	loanAsBytes, _ = json.Marshal(credit)
	stub.PutState(key, loanAsBytes)*/

	return args[0], nil
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func createLoanRequest(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 11 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	i1, err := strconv.Atoi(args[5])
	if err == nil {
		fmt.Println(i1)
	}

	if i1 < 18 || i1 > 60 {
		return "", fmt.Errorf("Applicant not eligible for loan due to age restrictions.")
	}

	loanApplication := &loanApplication{args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10]}
	LoanJSONasBytes, err := json.Marshal(loanApplication)
	if err != nil {
		return "", fmt.Errorf("Failed to Marshal asset: %s", args[0])
	}

	err = stub.PutState(args[0], []byte(LoanJSONasBytes))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}

// Get returns the value of the specified asset key
func getLoanOfUser(stub shim.ChaincodeStubInterface, args []string) (string, error) {
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
	return string(value), nil
}

// Set stores the asset (both key and value) on the ledger. If the key exists,
// it will override the value with the new one
func updateLoanStatus(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}
	var loanApplicationId = args[0]
	var loanApplicationStatus = args[1]

	loanAsBytes, _ := stub.GetState(loanApplicationId)
	loanApplication := loanApplication{}

	json.Unmarshal(loanAsBytes, &loanApplication)
	loanApplication.Status = loanApplicationStatus

	loanAsBytes, _ = json.Marshal(loanApplication)
	stub.PutState(loanApplicationId, loanAsBytes)

	return loanApplicationStatus, nil
}

// Query loan by bank
func (t *SimpleAsset) queryLoanByBank(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	BankId := args[0]

	queryString := fmt.Sprintf("{\"selector\": {\"BankId\": {\"$eq\": \"%s\" }}}", BankId)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
