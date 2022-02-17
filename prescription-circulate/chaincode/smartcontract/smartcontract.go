package smartcontract

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)
var INDEX string="PreID~PreState~CertHash"
type SmartContract struct {
	contractapi.Contract
}
//var INDEX string="index"
// InitLedger adds a base set of preHashs to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	//var INDEX string="PreID~PreState~CertHash"
	preHashs := []PrescriptionHash{
		{PreID:"KC0000001", PreState: "Genrated",   CertHash: "QmRYyR4m8CntJdvmMisuzXPCM8wqHDSTJSwcQaU3zSDaq7", SignatureHash: "QmZTGPUUTPNWNTfNryBExHfNm1732WKLjuQtqAeRfNbf1e"},
	    {PreID:"KC0000001", PreState: "Conducted",  CertHash: "QmSpQCBBAi9uji1s7GHDTjrAThrNXv62WcMrbbK3gvkGZM", SignatureHash: "QmPxd2PpwqUzpczMbvBx3EGbcQE386ax3hiHdiYqXjT8Fy"},
	    {PreID:"KC0000001", PreState: "Reimbursed", CertHash: "QmXM5i7mDHjGPjzmXcCUofSAJhiAZZFfTeK5bFqmePRbXF", SignatureHash: "QmPxd2PpwqUzpczMbvBx3EGbcQE386ax3hiHdiYqXjT8Fy"},
	}
    
	for _, preHash := range preHashs {
		preJson, err := json.Marshal(preHash)
		if err != nil {
			return err
		}
		//keys := [...]string{preHash.PreID, preHash.PreState, preHash.CertHash}
        index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{preHash.PreID,preHash.PreState,preHash.CertHash})
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(index, preJson)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// GetPreHashByCompositeKey returns preHash 
func (s *SmartContract) GetPreHashByCompositeKey(ctx contractapi.TransactionContextInterface, preID, preState, certHash string) (PrescriptionHash, error) {
	var preHash PrescriptionHash
    index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{preID, preState, certHash})
	if err != nil {
		return preHash, err
	}
    preJson, err := ctx.GetStub().GetState(index)
	if err != nil {
		return preHash, fmt.Errorf("failed to read from world state: %v", err)
	}
	if preJson == nil {
		return preHash, fmt.Errorf("the prescription [PreID: %s,PreState: %s,CertHash: %s ] does not exist", preID, preState, certHash)
	}
	err = json.Unmarshal(preJson, &preHash)
	if err != nil {
		return preHash, err
	}
	return preHash, nil
}

// GetPreHashByPreIDAndState returns all preHashes with the same ID and state found in world state
func (s *SmartContract) GetPreHashByPreIDAndState(ctx contractapi.TransactionContextInterface, preID, preState string) ([]PrescriptionHash, error) {
	//var INDEX string="PreID~PreState~CertHash"
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(INDEX, []string{preID, preState})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var preHashs []PrescriptionHash
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var preHash PrescriptionHash
		err = json.Unmarshal(queryResponse.Value, &preHash)
		if err != nil {
			return nil, err
		}
		preHashs = append(preHashs, preHash)
	}

	return preHashs, nil
}

// GetPreHashByPreID returns all preHashes with a same ID found in world state
func (s *SmartContract) GetPreHashByPreID(ctx contractapi.TransactionContextInterface, preID string) ([]PrescriptionHash, error) {
	//var INDEX string="PreID~PreState~CertHash"
	resultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(INDEX, []string{preID})
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var preHashs []PrescriptionHash
	for resultsIterator.HasNext(){
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var preHash PrescriptionHash
		err = json.Unmarshal(queryResponse.Value, &preHash)
		if err != nil {
			return nil, err
		}
		preHashs = append(preHashs, preHash)
	}

	return preHashs, nil
}

// CreatePreHash add a new preHash to world state
func (s *SmartContract) CreatePreHash(ctx contractapi.TransactionContextInterface, preID, preState, certHash, signatureHash string) error {
	exists, err := s.PreHashExists(ctx, preID,preState,certHash)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the prescription [PreID: %s,PreState: %s,CertHash: %s, SignatureHash: %s ] already exists", preID,preState,certHash,signatureHash)
	}

	preHash := PrescriptionHash{
		PreID:         preID,
		PreState:      preState,
		CertHash:      certHash,
		SignatureHash: signatureHash,
	}
    preJson, err := json.Marshal(preHash)
	if err != nil {
		return err
	}
	index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{preHash.PreID,preHash.PreState,preHash.CertHash})
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(index, preJson)
	return err
}

// DeletePreHash deletes an given preHash from the world state.
func (s *SmartContract) DeletePreHash(ctx contractapi.TransactionContextInterface, preID, preState, certHash string) error {
	exists, err := s.PreHashExists(ctx, preID, preState, certHash)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the prescription [PreID: %s,PreState: %s,CertHash: %s ] does not exist", preID, preState, certHash)
	}
    index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{preID, preState, certHash})
	if err != nil {
		return err
	}
	return ctx.GetStub().DelState(index)
}

func (s *SmartContract) PreHashExists(ctx contractapi.TransactionContextInterface, preID, preState, certHash string) (bool, error) {
	//var INDEX string="PreID~PreState~CertHash"
	index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{preID, preState, certHash})
	preJson, err := ctx.GetStub().GetState(index)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return preJson != nil, nil
}

// // GetAllAssets returns all assets found in world state
// func (s *SmartContract) GetAllPreHash(ctx contractapi.TransactionContextInterface) ([]PrescriptionHash, error) {
// 	// range query with empty string for startKey and endKey does an
// 	// open-ended query of all assets in the chaincode namespace.
// 	//var INDEX string="PreID~PreState~CertHash"
// 	index, err := ctx.GetStub().CreateCompositeKey(INDEX, []string{""})
	
// 	resultsIterator, err := ctx.GetStub().GetStateByRange(index, index)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resultsIterator.Close()

// 	var preHashs []PrescriptionHash
// 	for resultsIterator.HasNext() {
// 		queryResponse, err := resultsIterator.Next()
// 		if err != nil {
// 			return nil, err
// 		}

// 		var preHash PrescriptionHash
// 		err = json.Unmarshal(queryResponse.Value, &preHash)
// 		if err != nil {
// 			return nil, err
// 		}
// 		preHashs = append(preHashs, preHash)
// 	}

// 	return preHashs, nil
// }


