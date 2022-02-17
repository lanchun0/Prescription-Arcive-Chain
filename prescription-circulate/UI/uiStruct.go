package UI

import (
	"prescription-circulate/service"
)

type Application struct {
	Setup *service.ServiceSetup
}

type User struct {
	LoginName	string
	Password	string
	Identity	string
}

// type PrescriptionItem struct{
// 	Pre Prescription
// 	PreHash PrescriptionHash
// }

// type PrescriptionHash struct{
// 	PreID           string `json:"PreID"`                    //处方id 
// 	PreState        string `json:"PreState"`                 //处方当前状态-----> Genrated -----> Conducted --->Reimbursed
// 	CertHash        string `json:"CertHash"`                 //证书存放的hash地址
// 	SignatureHash   string `json:"SignatureHash"`            //签名过后SignedInfo的hash地址
// }

// type Prescription struct {
// 	PreID           string `json:"PreID"`                    //处方id
// 	PatientName      string `json:"PatientName"`            //患者姓名
// 	Gender           string `json:"Gender"`                 //患者性别
// 	Age              string `json:"Age"`                    //患者年龄
// 	IDNumber         string `json:"IDNumber"`               //患者身份证号
// 	Address          string `json:"Address"`                //患者住址
// 	Phone            string `json:"Phone"`                  //患者电话
// 	SurgeryDate      string `json:"SurgeryDate"`            //治疗时间
// 	Diagnosis        string `json:"Diagnosis"`              //诊断结果
// 	Recipe           string `json:"Recipe"`                 //治疗方法
// 	Fee              string `json:"Fee"`                    //费用
// 	PhysicianName    string `json:"PhysicianName"`          //医师姓名
// 	HospitalName     string `json:"HospitalName"`           //医院名称
// 	PharmacistName   string `json:"PharmacistName"`         //药房名
// 	PharmacyName     string `json:"PharmacyName"`           //药剂师姓名
// 	Payment          string `json:"Payment"`                //治疗价格
// 	PolicyNumber     string `json:"PolicyNumber"`           //保险单号
// 	InsuranceCompany string `json:"InsuranceCompany"`       //保险公司名称
// }

var Users []User

func Init() {
	doctor := User{LoginName:"111", Password:"123456", Identity:"Doctor"}
	pharmacy := User{LoginName:"222", Password:"123456", Identity:"Pharmacy"}
	insurance := User{LoginName:"333", Password:"123456", Identity:"Insurance"}

	Users = append(Users, doctor)
	Users = append(Users, pharmacy)
	Users = append(Users, insurance)

}