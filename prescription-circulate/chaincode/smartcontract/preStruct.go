package smartcontract

type PrescriptionHash struct{
	PreID           string `json:"PreID"`                    //处方id 
	PreState        string `json:"PreState"`                 //处方当前状态-----> Genrated -----> Conducted --->Reimbursed
	CertHash        string `json:"CertHash"`                 //证书存放的hash地址
	SignatureHash   string `json:"SignatureHash"`            //签名过后SignedInfo的hash地址
}
