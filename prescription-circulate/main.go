package main

import (
	"prescription-circulate/UI"
	"prescription-circulate/service"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)
// var HospitalCredPath string = "./key/hospital"
// var PharmacyCredPath string = "./key/pharmacy"
// var InsuranceCredPath string =  "./key/insurance"


func main() {
	log.Println("============ Prescription Circulate starts ============")

	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}

	if !wallet.Exists("appUser") {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}

	ccpPath := filepath.Join(
		"..",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)
    
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()
	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
    
	var serviceSetup service.ServiceSetup
	contract := network.GetContract("basic")

	//-----------------------test--------------------------------//
	// credPath := filepath.Join(
	// 	"..",
	// 	"fabric-samples",
	// 	"test-network",
	// 	"organizations",
	// 	"peerOrganizations",
	// 	"org1.example.com",
	// 	"users",
	// 	"User1@org1.example.com",
	// 	"msp",
	// )
	
    err = serviceSetup.Init(contract)
	if err != nil {
		log.Fatalf("Failed to init the service: %v", err)
	}
    // insuranceCert := serviceSetup.GetCertAddress()

	// err = serviceSetup.Init(contract,PharmacyCredPath)
	// if err != nil {
	// 	log.Fatalf("Failed to init the service: %v", err)
	// }
    // pharmacyCert := serviceSetup.GetCertAddress()

	// err = serviceSetup.Init(contract,HospitalCredPath)
	// if err != nil {
	// 	log.Fatalf("Failed to init the service: %v", err)
	// }
    // hospitalCert := serviceSetup.GetCertAddress()

	log.Println("Succeed initilizing service")
	// log.Println("Hospital: ",hospitalCert)
	// log.Println("Pharmacy: ",pharmacyCert)
	// log.Println("Insurance: ",insuranceCert)
    
	err = serviceSetup.ChangeIdentity("Doctor")
	if err != nil {
		log.Fatalf("Failed to change identity: %v", err)
	}
	//log.Println("=====================Hospital===========================")
    
	// pre := service.Prescription{
	// 	PreID: "KC0000001",	
	// 	    PatientName: "法外狂徒张三",
	//         Gender:"男",
	//         Age:"22",
	//         IDNumber:"66566199810253162",
	//         Address:"上海",
	//         Phone:"15327681440",
	//         SurgeryDate:"2021年01月01日",
	//         Diagnosis:"流行性感冒，鼻窦炎",
	//         Recipe:"阿达帕林*1盒，1天三次。",
	//         Fee:"¥153.00",
	//         PhysicianName:"张懿虎",
	//         HospitalName:"丸摩堂",
	//         PharmacistName:"*",
	//         PharmacyName:"*",
	//         Payment:"*",
	//         PolicyNumber:"*",
	//         InsuranceCompany:"*",
	// }
	// pre2 := service.Prescription{
	// 	PreID: "KC0000001",	
	// 	    PatientName: "法外狂徒张三",
	//         Gender:"男",
	//         Age:"22",
	//         IDNumber:"66566199810253162",
	//         Address:"上海",
	//         Phone:"15327681440",
	//         SurgeryDate:"2021年01月01日",
	//         Diagnosis:"流行性感冒，鼻窦炎",
	//         Recipe:"阿达帕林*1盒，1天三次。",
	//         Fee:"¥153.00",
	//         PhysicianName:"张懿虎",
	//         HospitalName:"丸摩堂",
	//         PharmacistName:"崔彩礼",
	//         PharmacyName:"一点点",
	//         Payment:"自费",
	//         PolicyNumber:"*",
	//         InsuranceCompany:"*",
    // }
	// pre3 := service.Prescription{
	// 	PreID: "KC0000001",	
	// 	    PatientName: "法外狂徒张三",
	//         Gender:"男",
	//         Age:"22",
	//         IDNumber:"66566199810253162",
	//         Address:"上海",
	//         Phone:"15327681440",
	//         SurgeryDate:"2021年01月01日",
	//         Diagnosis:"流行性感冒，鼻窦炎",
	//         Recipe:"阿达帕林*1盒，1天三次。",
	//         Fee:"¥153.00",
	//         PhysicianName:"张懿虎",
	//         HospitalName:"丸摩堂",
	//         PharmacistName:"崔彩礼",
	//         PharmacyName:"一点点",
	//         Payment:"保险报销",
	//         PolicyNumber:"KFC000036",
	//         InsuranceCompany:"社保",
    // }

	// err = serviceSetup.GeneratePrescription(pre)
	// if err != nil {
	// 	log.Fatalf("Failed to generate the prescription: %v", err)
	// }
	// log.Println("Succeed genrating prescription")

    // log.Println("=====================Pharmacy===========================")
	// serviceSetup.ChangeIdentity("Pharmacy")
    // if err != nil {
	// 	log.Fatalf("Failed to change identity to Pharmacy: %v", err)
	// }
	// log.Println("prescription KC0000001 status:")
	// result,err := serviceSetup.FindPrescriptionByPreID("KC0000001")
	// if err != nil {
	// 	log.Fatalf("Failed to generate the prescription: %v", err)
	// }
	// log.Println(result)

	// err = serviceSetup.ConductPrescription(pre2,result[0])
	// if err != nil {
	// 	log.Fatalf("Failed to conduct the prescription: %v", err)
	// }
    // log.Println("Succeed conducting prescription")

	// log.Println("prescription KC0000001 status:")
	// result,err = serviceSetup.FindPrescriptionByPreID("KC0000001")
	// if err != nil {
	// 	log.Fatalf("Failed to generate the prescription: %v", err)
	// }
	// log.Println(result)
    
	// log.Println("=====================Insurance===========================")
	// serviceSetup.ChangeIdentity("Insurance")
    // if err != nil {
	// 	log.Fatalf("Failed to change identity to Insurance: %v", err)
	// }
	// err = serviceSetup.ReimbursePrescription(pre3,result[0])
	// if err != nil {
	// 	log.Fatalf("Failed to conduct the prescription: %v", err)
	// }
    // log.Println("Succeed reimbursing prescription")

	// log.Println("prescription KC0000001 status:")
	// result,err = serviceSetup.FindPrescriptionByPreID("KC0000001")
	// if err != nil {
	// 	log.Fatalf("Failed to reimburse the prescription: %v", err)
	// }
	// log.Println(result)
	app := UI.Application{
		Setup: &serviceSetup,
	}
	UI.UiStart(app)
	
}

func populateWallet(wallet *gateway.Wallet) error {
	log.Println("============ Populating wallet ============")
	credPath := filepath.Join(
		"..",
		"fabric-samples",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}