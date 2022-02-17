package service

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"fmt"
	"log"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"encoding/json"
	shell "github.com/ipfs/go-ipfs-api"
)
var sh *shell.Shell
var DoctorCred_Path string="./key/hospital"
var PharmacyCred_Path string="./key/pharmacy"
var InsuranceCred_Path string= "./key/insurance"
var DoctorCert_Hash string
var PharmacyCert_Hash string
var InsuranceCert_Hash string

//reimburse the prescription
//Args: "pre" the conducted prescription content
//      "preItem" the original prescription
func (t *ServiceSetup) ReimbursePrescription(pre Prescription, preItem PrescriptionItem) error{
	contract := t.Contract
	originalPreHash := preItem.PreHash
	//var preHash PrescriptionHash
	preItem, err := t.FindPrescriptionWithIdentity(originalPreHash.PreID, "Conducted",originalPreHash.CertHash)
	if err != nil{
		return err
	}
	onChainPre := preItem.Pre
	onChainPre.Payment = pre.Payment
	onChainPre.PolicyNumber = pre.PolicyNumber
	onChainPre.InsuranceCompany = pre.InsuranceCompany

	signInfo,err := signPrescription(t.Key, onChainPre)
	if err != nil{
		return err
	}
	signHash := uploadSignature(signInfo)
	_, err = contract.SubmitTransaction("CreatePreHash",onChainPre.PreID,"Reimbursed",t.CertAddress,signHash)
	if err != nil{
		return err
	}
	return nil
}

//conduct the prescription
//Args: "pre" the conducted prescription content
//      "preItem" the original prescription
func (t *ServiceSetup) ConductPrescription(pre Prescription, preItem PrescriptionItem) error{
	contract := t.Contract
	originalPreHash := preItem.PreHash
	//var preHash PrescriptionHash
	preItem, err := t.FindPrescriptionWithIdentity(originalPreHash.PreID, "Generated",originalPreHash.CertHash)
	if err != nil{
		return err
	}
	// if preItem == nil{
	// 	return fmt.Errorf("Prescription :[PreID: %v] does not exist",originalPreHash.PreID)
	// }
	onChainPre := preItem.Pre

	onChainPre.PharmacistName = pre.PharmacistName
	onChainPre.PharmacyName = pre.PharmacyName
	//onChainPre.Payment = pre.Payment

	signInfo,err := signPrescription(t.Key, onChainPre)
	if err != nil{
		return err
	}
	signHash := uploadSignature(signInfo)
	_, err = contract.SubmitTransaction("CreatePreHash",onChainPre.PreID,"Conducted",t.CertAddress,signHash)
	if err != nil{
		return err
	}
	return nil
}

func (t *ServiceSetup) GeneratePrescription(pre Prescription) error{
	signInfo,err := signPrescription(t.Key, pre)
	if err != nil{
		return err
	}
	signHash := uploadSignature(signInfo)
	contract :=t.Contract
	_, err = contract.SubmitTransaction("CreatePreHash",pre.PreID,"Generated",t.CertAddress,signHash)
	if err != nil{
		return err
	}
	return nil
}

func (t *ServiceSetup) FindPrescriptionWithIdentity(preID, preState, certHash string) (PrescriptionItem,error){
	var preItem PrescriptionItem
	var preHash PrescriptionHash
	contract := t.Contract
	preHashJson, err := contract.EvaluateTransaction("GetPreHashByCompositeKey",preID, preState, certHash)
	if err != nil{
		return preItem, err
	}
	err = json.Unmarshal(preHashJson, &preHash)
	if err != nil{
		return preItem, err
	}
	pre, err := VerifyPreHash(preHash)
	if err != nil{
		return preItem, err
	}
	preItem.Pre = pre
	preItem.PreHash = preHash
	return preItem,nil
}

func (t *ServiceSetup) FindPrescriptionByPreIDAndState(preID, preState string) ([]PrescriptionItem,error){
	var preItems []PrescriptionItem
	var preItem PrescriptionItem
	//var pres []Prescription
	var preHashs []PrescriptionHash
	contract := t.Contract
	preHashsJson,err := contract.EvaluateTransaction("GetPreHashByPreIDAndState",preID, preState)
	if err != nil{
		return preItems, err
	}
    err = json.Unmarshal(preHashsJson, &preHashs)
	if err != nil{
		return preItems, err
	}

    for _, preHash := range preHashs{
		pre, err := VerifyPreHash(preHash)
		if err != nil {
			return preItems, err
		}
		preItem.Pre = pre
		preItem.PreHash =preHash
		preItems = append(preItems,preItem)
	}

	return preItems, nil
}

func (t *ServiceSetup) FindPrescriptionByPreID(preID string) ([]PrescriptionItem,error) {
	var preItems []PrescriptionItem
	var preItem PrescriptionItem
	var preHashs []PrescriptionHash
	contract := t.Contract
	preHashsJson,err := contract.EvaluateTransaction("GetPreHashByPreID",preID)
	if err != nil{
		return preItems, err
	}
    err = json.Unmarshal(preHashsJson, &preHashs)

    for _, preHash := range preHashs{
		pre, err := VerifyPreHash(preHash)
		if err != nil {
			return preItems, err
		}
		preItem.Pre = pre
		preItem.PreHash =preHash
		preItems = append(preItems,preItem)
	}

	return preItems, nil
}

// Allow to change identity, but it verify the certAddr firstly
func (t *ServiceSetup) ChangeIdentity(identity string) error{
	var credPath,certAddr string
	switch identity{
		case "Doctor":
			credPath = DoctorCred_Path
			certAddr = DoctorCert_Hash  
        case "Pharmacy":
			credPath = PharmacyCred_Path
			certAddr = PharmacyCert_Hash
		case "Insurance":
			credPath = InsuranceCred_Path
			certAddr = InsuranceCert_Hash
		default:
			return fmt.Errorf("Failed to change identity: input should be :[Doctor, Pharmacy or Insucance]")
	}
	
	certB, keyB, err := getCertAndKey(credPath)
	if err != nil{
		return err
	}

	cert := downloadCert(certAddr)
    if !bytes.Equal(certB,cert){
		return fmt.Errorf("Failed to change identity: IPFS hash value is not correct")
	}
    t.CertAddress = certAddr
	t.Key = keyB
	return nil
}

func (t *ServiceSetup) GetCertAddress() string{
	return t.CertAddress
}

func (t *ServiceSetup) SetCred(credPath string) error{
	certB, keyB, err := getCertAndKey(credPath)
	if err != nil{
		return err
	}
    
	certAddr, err := uploadCert(certB)
	if err != nil{
		return err
	}
    //log.Println("Success upload the cert: ",certAddr)
	t.Key = keyB
	t.CertAddress = certAddr
	return nil
}

func (t *ServiceSetup) Init(contract *gateway.Contract) error{
	t.Contract = contract
	err := t.SetCred(InsuranceCred_Path)
	if err != nil{
		return err
	}
	InsuranceCert_Hash = t.GetCertAddress()
	
	err = t.SetCred(PharmacyCred_Path)
	if err != nil{
		return err
	}
	PharmacyCert_Hash = t.GetCertAddress()
    
	err = t.SetCred(DoctorCred_Path)
	if err != nil{
		return err
	}
	DoctorCert_Hash = t.GetCertAddress()

	return nil
}
//verify PreHash and transform it into Prescription
func VerifyPreHash(preHash PrescriptionHash) (Prescription, error){
	var pre Prescription
	cert := downloadCert(preHash.CertHash)
	signInfo := downloadSignature(preHash.SignatureHash)
	if !verifySignature(cert, signInfo){
		return pre,fmt.Errorf("Prescription :[PreID: %v, State: %v, CertAddr: %v] is not valid",preHash.PreID,preHash.PreState,preHash.CertHash)
	}
	pre = signInfo.Pre
	return pre, nil
}

//private method

func downloadSignature(SigHash string) SignedInfo{
	sh = shell.NewShell("localhost:5001")
	read, err := sh.Cat(SigHash)
	if err != nil {
		log.Fatalf("Failed to download from IPFS: %v", err)
	}
	sigJson, err := ioutil.ReadAll(read)
    var signature SignedInfo
	err = json.Unmarshal(sigJson, &signature)
	if err != nil {
		log.Fatalf("Failed to Unmarshal the signature: %v", err)
	}
	return signature

}

func uploadSignature(signature SignedInfo) string{
	sigJson, err := json.Marshal(&signature)
	if err != nil {
		log.Fatalf("Failed to Marshal the signature: %v", err)
	}
	sh = shell.NewShell("localhost:5001")
	hash, err := sh.Add(bytes.NewBufferString(string(sigJson)))
	if err != nil {
		log.Fatalf("Failed to upload to IPFS: %v", err)
	}
	return hash

}

func signPrescription(key []byte, pre Prescription)(SignedInfo, error){
	var signature SignedInfo
	signature.Pre = pre
	preJson, err := json.Marshal(&pre)
	if err != nil{
		return signature, err
	}
	preHashValue := sha256.Sum256(preJson)

	privKey, _ := x509.ParsePKCS8PrivateKey(key)
	ecdsaKey := privKey.(*ecdsa.PrivateKey)
	r, s, _ := ecdsa.Sign(rand.Reader, ecdsaKey, preHashValue[:])

	signature.SignR, err=r.MarshalText()
	if err!=nil{
		return signature, err
	}

	signature.SignS, err=s.MarshalText()
	if err!=nil{
		return signature, err
	}
	return signature, err

}

func verifySignature(cert []byte, signature SignedInfo) bool{
	preJson, err := json.Marshal(&signature.Pre)
	if err != nil{
		log.Fatalf("Failed to Marshal the prescription: %v", err)
	}
	preHashValue := sha256.Sum256(preJson)
    var r,s big.Int							
	r.UnmarshalText(signature.SignR)
	s.UnmarshalText(signature.SignS)
	
	pubInterface,_ := x509.ParsePKIXPublicKey(cert)
	certKey := pubInterface.(*ecdsa.PublicKey)
	
	return ecdsa.Verify(certKey,preHashValue[:],&r,&s)

}

func uploadCert(cert []byte) (string, error){
	var hash string
	certJson, err := json.Marshal(&cert)
	if err != nil {
		return hash, fmt.Errorf("Failed to Marshal the certificate: %v", err)
	}
	sh := shell.NewShell("localhost:5001")
	hash, err = sh.Add(bytes.NewBufferString(string(certJson)))
	if err != nil {
		return hash, fmt.Errorf("Failed to upload to IPFS: %v", err)
	}
	return hash, nil
}

func downloadCert(certHash string) []byte{
	sh = shell.NewShell("localhost:5001")
	read, err := sh.Cat(certHash)
	if err != nil {
		log.Fatalf("Failed to download from IPFS: %v", err)
	}
	certJson, err := ioutil.ReadAll(read)
    var cert []byte
	err = json.Unmarshal(certJson, &cert)
	if err != nil {
		log.Fatalf("Failed to Unmarshal the signature: %v", err)
	}
	return cert
}

func getCertAndKey(credPath string) (certB []byte,keyB []byte,err error){
	var cert []byte
	var key []byte
	//-----------------get the key-----------//
	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err = ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return cert,key,err
	}

    keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return cert,key,err
	}
	if len(files) != 1 {
		return cert,key,err
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err = ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return cert,key,err
	}

    blkCert, _ := pem.Decode(cert)
	certKey, _ := x509.ParseCertificate(blkCert.Bytes)
	pubkey := certKey.PublicKey.(*ecdsa.PublicKey)
	certB,err = x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return certB , key, err
	}

	blkPriv, _ := pem.Decode(key)
	privKey, _ := x509.ParsePKCS8PrivateKey(blkPriv.Bytes)
	ecdsaKey := privKey.(*ecdsa.PrivateKey)
	keyB, err = x509.MarshalPKCS8PrivateKey(ecdsaKey)
	if err != nil {
		return certB , keyB, err
	}
    return certB , keyB, err
}