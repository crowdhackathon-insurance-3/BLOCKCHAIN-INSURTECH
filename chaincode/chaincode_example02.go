// package name must be main
package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"crypto/x509/pkix"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"crypto/x509"
	"github.com/golang/protobuf/proto"
	"encoding/pem"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/hyperledger/fabric/protos/msp"
	"encoding/asn1"
	"bytes"
)


const (
	InsuranceProductCompositeKey  = "insuranceProductCompositeKey"
	InsuranceEntryCompositeKey = "insuranceEntryCompositeKey"
)

// InsuranceProductType
type InsuranceProductType struct {
	Id  		    	string `json:"id"`
	Name        		string `json:"name"`
	StartDate  			string `json:"startDate"`
	EndDate  			string `json:"endDate"`
	MinAmount 			int64  `json:"minAmount"`
	MaxAmount 			int64  `json:"maxAmount"`
	Terms				string `json:"terms"`
	DRate	 			int64  `json:"dRate"`
	//
	ObjectType 			string `json:"docType"`
}

// InsuranceEntryType
type InsuranceEntryType struct {
	// in constructor
	Id     		     	string `json:"id"`
	Item	         	string `json:"item"`
	StartDate 		 	string `json:"startDate"`
	NumOfDays  		 	int64  `json:"numOfDays"`
	Amount           	int64  `json:"amount"`
	ProductId		 	string	`json:"productId"` // foreign key
	// calculated
	PremiumAmount    	int64  `json:"pAmount"`   // calculated (not in constructor)
	//
	ObjectType 			string `json:"docType"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init chaincode")

	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("Invoke function=%q\n", function)
	if function == "createInsuranceProduct" {
		// Create Insurance Product
		return t.createInsuranceProduct(stub, args)
	} else if function == "createInsuranceEntry" {
			// Create Insurance Entry
			return t.createInsuranceEntry(stub, args)
	} else if function == "retrieveInsuranceProduct" {
		// Retrieve Insurance Product
		return t.retrieveInsuranceProduct(stub, args)
	} else if function == "retrieveInsuranceEntry" {
		// Retrieve Insurance Entry
		return t.retrieveInsuranceEntry(stub, args)
	} else if function == "retrieveInsuranceProducts" {
		// Retrieve Insurance Products
		return t.retrieveInsuranceProducts(stub, args)
	} else if function == "retrieveInsuranceEntries" {
		// Retrieve Insurance Entries
		return t.retrieveInsuranceEntries(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"createInsuranceProduct\" \"createInsuranceEntry\" ")
}

func (t *SimpleChaincode) createInsuranceProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var insuranceProductKey string

	fmt.Printf("number of arguments received %d\n", len(args))

	if len(args) != 8 {
		msg:=fmt.Sprintf("Incorrect number of arguments. Expecting 8, received %d\n", len(args))
		fmt.Printf(msg)
		return shim.Error(msg)
	}
	fmt.Printf("- createInsuranceProduct(Id=%q,Name=%q,StartDate=%q,EndDate=%q,MinAmount=%q,MaxAmount=%q,terms=%q,dRate=%q)\n", args[0], args[1], args[2], args[3], args[4], args[5], args[6], args[7])

	insuranceProduct := InsuranceProductType{}
	insuranceProduct.ObjectType = "product"
	insuranceProduct.Id = args[0]
	insuranceProduct.Name = args[1]
	insuranceProduct.StartDate = args[2]
	insuranceProduct.EndDate = args[3]
	insuranceProduct.MinAmount, err = strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for MinAmount")
	}
	insuranceProduct.MaxAmount, err = strconv.ParseInt(args[5], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for MaxAmount")
	}
	insuranceProduct.Terms = args[6]
	insuranceProduct.DRate, err = strconv.ParseInt(args[7], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for DRate")
	}

	insuranceProductKey, _ = stub.CreateCompositeKey(InsuranceProductCompositeKey, []string{insuranceProduct.Id})

	insuranceProductJsonAsBytes, _ := json.Marshal(insuranceProduct)
	// Write the state to the ledger
	err = stub.PutState(insuranceProductKey, insuranceProductJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- created insuranceProduct=%s\n", insuranceProductJsonAsBytes)

	return shim.Success(insuranceProductJsonAsBytes)
}

func (t *SimpleChaincode) createInsuranceEntry(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var insuranceEntryKey string

	fmt.Printf("number of arguments received %d\n", len(args))

	if len(args) != 6 {
		msg:=fmt.Sprintf("Incorrect number of arguments. Expecting 6, received %d\n", len(args))
		fmt.Printf(msg)
		return shim.Error(msg)
	}
	fmt.Printf("- createInsuranceEntry(id=%q,item=%q,startDate=%q,numOfDays=%q,amount=%q,productId=%q)\n", args[0], args[1], args[2], args[3], args[4], args[5])

	insuranceEntry := InsuranceEntryType{}
	insuranceEntry.ObjectType = "entry"
	insuranceEntry.Id = args[0]
	insuranceEntry.Item = args[1]
	insuranceEntry.StartDate = args[2]
	insuranceEntry.NumOfDays, err = strconv.ParseInt(args[3], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for NumOfDays")
	}
	insuranceEntry.Amount, err = strconv.ParseInt(args[4], 10, 64)
	if err != nil {
		return shim.Error("Expecting integer value for Amount")
	}
	insuranceEntry.ProductId = args[5]

	// authorise
	// TODO authorise

	// Validate terms
	// TODO validate

	// TODO calculate Premium Amount based on InsuranceProduct dRate

	// calculate the premium amount
	insuranceEntry.PremiumAmount = int64( float64(insuranceEntry.Amount * insuranceEntry.NumOfDays) * 0.1 )

	insuranceEntryKey, _ = stub.CreateCompositeKey(InsuranceEntryCompositeKey, []string{insuranceEntry.Id})

	insuranceEntryJsonAsBytes, _ := json.Marshal(insuranceEntry)
	// Write the state to the ledger
	err = stub.PutState(insuranceEntryKey, insuranceEntryJsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- created insuranceEntry=%s\n", insuranceEntryJsonAsBytes)

	return shim.Success(insuranceEntryJsonAsBytes)
}

// Retrieve InsuranceProductType
func (t *SimpleChaincode) retrieveInsuranceProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var insuranceProductKey string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Printf("- retrieveInsuranceProduct(Id=%q)\n", args[0])
	var insuranceProductId = args[0]

	insuranceProductKey, _ = stub.CreateCompositeKey(InsuranceProductCompositeKey, []string{insuranceProductId})
	insuranceProductAsBytes, err := stub.GetState(insuranceProductKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- returned insurance product=%s\n", insuranceProductAsBytes)
	return shim.Success(insuranceProductAsBytes)
}

// Retrieve Insurance Entry
func (t *SimpleChaincode) retrieveInsuranceEntry(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var insuranceEntryKey string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	fmt.Printf("- retrieveInsuranceEntry(Id=%q)\n", args[0])
	var insuranceEntryId = args[0]

	insuranceEntryKey, _ = stub.CreateCompositeKey(InsuranceProductCompositeKey, []string{insuranceEntryId})
	insuranceEntryAsBytes, err := stub.GetState(insuranceEntryKey)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- returned insurance entry=%s\n", insuranceEntryAsBytes)
	return shim.Success(insuranceEntryAsBytes)
}

// Retrieve Insurance Products
func (t *SimpleChaincode) retrieveInsuranceProducts(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// verify number of parameters
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	fmt.Printf("- retrieveInsuranceProducts()\n")

	// perform query
	queryString :=  fmt.Sprintf(`
{
   "selector": {
	  "docType": "product",
      "_id": {
         "$gt": null
      }
   }
}`)
	fmt.Printf("- queryString=%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if resultsIterator!=nil { defer resultsIterator.Close() }
	if err != nil {
		return shim.Error(err.Error())
	}
	// process results
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		// for each user
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		//buffer.WriteString("{\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")
		//buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	//fmt.Printf("- return: %s\n", buffer.String())
	asBytes := buffer.Bytes()
	fmt.Printf("- return: %s\n", string(asBytes))
	return shim.Success(asBytes)
}

// Retrieve Insurance Entries
func (t *SimpleChaincode) retrieveInsuranceEntries(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// verify number of parameters
	if len(args) != 0 {
		return shim.Error("Incorrect number of arguments. Expecting 0")
	}
	fmt.Printf("- retrieveInsuranceEntries()\n")

	// perform query
	queryString :=  fmt.Sprintf(`
{
   "selector": {
	  "docType": "entry",
      "_id": {
         "$gt": null
      }
   }
}`)
	fmt.Printf("- queryString=%s\n", queryString)
	resultsIterator, err := stub.GetQueryResult(queryString)
	if resultsIterator!=nil { defer resultsIterator.Close() }
	if err != nil {
		return shim.Error(err.Error())
	}
	// process results
	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		// for each user
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		//buffer.WriteString("{\"Key\":")
		//buffer.WriteString("\"")
		//buffer.WriteString(queryResponse.Key)
		//buffer.WriteString("\"")
		//buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	//fmt.Printf("- return: %s\n", buffer.String())
	asBytes := buffer.Bytes()
	fmt.Printf("- return: %s\n", string(asBytes))
	return shim.Success(asBytes)
}

func isAdmin(stub shim.ChaincodeStubInterface) (bool, error) {
	userCN, err := getUserCN(stub)
	if err!=nil {
		return false, errors.New(fmt.Sprintf("Unable to get CN: %s", err.Error()))
	}
	if userCN=="admin" || userCN=="admin1" || userCN=="admin2" || userCN=="admin3" || userCN=="admin4" || userCN=="admin5" {
		return true, nil
	}
	return false, nil
}

func logAuthInfo(stub shim.ChaincodeStubInterface) {
	// logging
	user, err := cid.GetID(stub)
	if err!=nil {
		fmt.Printf("Failed to get user, err: %+v\n", err)
	}
	fmt.Printf("user=%s\n", user)
	userDN, err := getUserDN(stub)
	fmt.Printf("userDN=%s\n", userDN)
	mspid, err := cid.GetMSPID(stub)
	fmt.Printf("mspid=%s\n", mspid)
	val, ok, err := cid.GetAttributeValue(stub, "phone")
	if err != nil {
		// There was an error trying to retrieve the attribute
	}
	if !ok {
		// The client identity does not possess the attribute
		fmt.Printf("val phone not exists\n")
	} else {
		fmt.Printf("val phone=%s\n", val)
	}
}

func getUserCN(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := getCertificate(stub)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to get certificate: %s", err.Error()))
	}
	cn := getCN(&cert.Subject)
	return cn, nil
}
func getUserDN(stub shim.ChaincodeStubInterface) (string, error) {
	cert, err := getCertificate(stub)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Unable to get certificate: %s", err.Error()))
	}
	cn := getDN(&cert.Subject)
	return cn, nil
}

//
func getCertificate(stub shim.ChaincodeStubInterface) (*x509.Certificate, error) {
	serializedID, _ := stub.GetCreator()

	sId := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sId)
	if err != nil {
		errmsg:=fmt.Sprintf("Failed to deserialize identity, err %+v", err)
		return nil, errors.New(errmsg)
	}
	bl, _ := pem.Decode(sId.IdBytes)
	if bl == nil {
		return nil, errors.New("Failed to decode PEM structure")
	}
	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		errmsg:=fmt.Sprintf("ParseCertificate failed %s", err)
		return nil, errors.New(errmsg)
	}

	return cert, nil
}

var attributeTypeNames = map[string]string{
	"2.5.4.6":  "C",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.3":  "CN",
	"2.5.4.5":  "SERIALNUMBER",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.9":  "STREET",
	"2.5.4.17": "POSTALCODE",
}
func getDN(name *pkix.Name) string {

	r := name.ToRDNSequence()
	s := ""
	for i := 0; i < len(r); i++ {
		rdn := r[len(r)-1-i]
		if i > 0 {
			s += ","
		}
		for j, tv := range rdn {
			if j > 0 {
				s += "+"
			}
			typeString := tv.Type.String()
			typeName, ok := attributeTypeNames[typeString]
			if !ok {
				derBytes, err := asn1.Marshal(tv.Value)
				if err == nil {
					s += typeString + "=#" + hex.EncodeToString(derBytes)
					continue // No value escaping necessary.
				}
				typeName = typeString
			}
			valueString := fmt.Sprint(tv.Value)
			escaped := ""
			begin := 0
			for idx, c := range valueString {
				if (idx == 0 && (c == ' ' || c == '#')) ||
					(idx == len(valueString)-1 && c == ' ') {
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
					continue
				}
				switch c {
				case ',', '+', '"', '\\', '<', '>', ';':
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
				}
			}
			escaped += valueString[begin:]
			s += typeName + "=" + escaped
		}
	}
	return s
}
var cnAttributeTypeNames = map[string]string{
	"2.5.4.3":  "CN",
}
func getCN(name *pkix.Name) string {

	r := name.ToRDNSequence()
	s := ""
	for i := 0; i < len(r); i++ {
		rdn := r[len(r)-1-i]
		//if i > 0 {
		//	s += ","
		//}
		//		for j, tv := range rdn {
		for _, tv := range rdn {
			// if j > 0 {
			// 	s += "+"
			// }
			typeString := tv.Type.String()
			_, ok := cnAttributeTypeNames[typeString]
			if !ok {
				continue // No value escaping necessary.
			}
			valueString := fmt.Sprint(tv.Value)
			escaped := ""
			begin := 0
			for idx, c := range valueString {
				if (idx == 0 && (c == ' ' || c == '#')) ||
					(idx == len(valueString)-1 && c == ' ') {
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
					continue
				}
				switch c {
				case ',', '+', '"', '\\', '<', '>', ';':
					escaped += valueString[begin:idx]
					escaped += "\\" + string(c)
					begin = idx + 1
				}
			}
			escaped += valueString[begin:]
			//s += typeName + "=" + escaped
			s += escaped
		}
	}
	return s
}
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
