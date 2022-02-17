package UI

import (
	"fmt"
	"github.com/andlabs/ui"
	"prescription-circulate/service"
)
var currentPreItems []service.PrescriptionItem
//var currentPreItem PrescriptionItem
var pageNumber int
var currentPage int
var currentUser User
//var countLabel *ui.Label
var mainWindow *ui.Window
//var count int

func (app *Application) insuranceMakeControlWritePage() ui.Control{
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	refreshButton := ui.NewButton("刷新")
    hbox.Append(refreshButton, false)
    submitButton := ui.NewButton("报销处方")
	hbox.Append(submitButton, false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("写入处方")
	group.SetMargined(true)
	vbox.Append(group, true)

    group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
    
    preIdLabel := ui.NewLabel("处方号: ")
	preStateLabel := ui.NewLabel("处方状态： ")
	certHashLabel := ui.NewLabel("签名公钥地址： ")
	signatureHashLabel := ui.NewLabel("处方存储地址： ")
	surgeryDateLabel := ui.NewLabel("就诊日期: ")
	patientNameLabel := ui.NewLabel("患者姓名: ")
	genderLabel := ui.NewLabel("患者性别: ")
	iDNumberLabel := ui.NewLabel("身份证号: ")
	phoneLabel := ui.NewLabel("患者电话: ")
	addressLabel := ui.NewLabel("现居地址: ")
	ageLabel := ui.NewLabel("患者年龄: ")
	diagnosisLabel := ui.NewLabel("诊断结果: ")
	recipeLabel := ui.NewLabel("治疗方案: ")
	hospitalNameLabel := ui.NewLabel("医院: ")
	physicianNameLabel := ui.NewLabel("主治医师: ")
	feeLabel := ui.NewLabel("费用: ")
	pharmacistNameLabel := ui.NewLabel("药房: ")
	pharmacyLabel := ui.NewLabel("药剂师: ")
	paymentLabel := ui.NewLabel("支付方式: ")
	insuranceLabel := ui.NewLabel("保险公司: ")
	policyNumberLabel := ui.NewLabel("保险单号: ")
	vbox.Append(preIdLabel, false)
	vbox.Append(preStateLabel, false)
	vbox.Append(certHashLabel, false)
	vbox.Append(signatureHashLabel, false)
	vbox.Append(surgeryDateLabel, false)
	vbox.Append(patientNameLabel, false)
	vbox.Append(genderLabel, false)
	vbox.Append(iDNumberLabel, false)
	vbox.Append(phoneLabel, false)
	vbox.Append(addressLabel, false)
	vbox.Append(ageLabel, false)
	vbox.Append(diagnosisLabel, false)
	vbox.Append(recipeLabel, false)
	vbox.Append(hospitalNameLabel, false)
	vbox.Append(physicianNameLabel, false)
	vbox.Append(feeLabel, false)
	vbox.Append(pharmacistNameLabel, false)
	vbox.Append(pharmacyLabel, false)
	vbox.Append(paymentLabel, false)
	vbox.Append(insuranceLabel, false)
	vbox.Append(policyNumberLabel, false)
		
	paymentEntry := ui.NewEntry()
	insuranceEntry := ui.NewEntry()
	policyNumberEntry := ui.NewEntry()
	entryForm.Append("支付方式",paymentEntry,false)
	entryForm.Append("保险公司", insuranceEntry, false)
	entryForm.Append("保险单号", policyNumberEntry, false)
    
	submitButton.OnClicked(func(*ui.Button){	
		if pageNumber==-1 || currentPage > pageNumber{
			return
		}
		var preItem service.PrescriptionItem
	    var exist bool=false
	    for _,preItem = range currentPreItems{
			if preItem.PreHash.PreState == "Conducted"{
				exist = true
				break
			}
		}
		if !exist{
			return
		}		
		paymentText := paymentEntry.Text()
		insuranceText := insuranceEntry.Text()
		policyNumberText := policyNumberEntry.Text()
		if paymentText=="" || insuranceText=="" || policyNumberText==""{
			ui.MsgBox(mainWindow,
				"关键信息不可为空",
				"请输入支付方式，保险单号和保险公司全名")
			return 
		}
		var pre service.Prescription
		pre = preItem.Pre
		pre.Payment = paymentText
		pre.PolicyNumber = policyNumberText
		pre.InsuranceCompany = insuranceText
		err:= app.Setup.ReimbursePrescription(pre,preItem)
		if err != nil{
			ui.MsgBox(mainWindow,
				"报销处方",
				err.Error())
			return
		} 
		if err == nil{
			ui.MsgBox(mainWindow,
				"处方报销成功",
			    "报销数据已上传区块链")
			return
		} 
	})

    refreshButton.OnClicked(func(*ui.Button) {
		if pageNumber!=-1 && currentPage <= pageNumber{
			var preItem service.PrescriptionItem
	        var exist bool=false
	        for _,preItem = range currentPreItems{
				if preItem.PreHash.PreState == "Conducted"{
				exist = true
				break
			    }
		    }
	        if exist {
				preIdLabel.SetText("处方号: "+preItem.PreHash.PreID)		
		        preStateLabel.SetText("处方状态： 已执行")
		        certHashLabel.SetText("签名公钥地址： "+preItem.PreHash.CertHash)
		        signatureHashLabel.SetText("处方存储地址： "+preItem.PreHash.SignatureHash)
		        surgeryDateLabel.SetText("就诊日期: "+preItem.Pre.SurgeryDate)
		        patientNameLabel.SetText("患者姓名: "+preItem.Pre.PatientName)
		        genderLabel.SetText("患者性别: "+preItem.Pre.Gender)
		        iDNumberLabel.SetText("身份证号: "+preItem.Pre.IDNumber)
		        phoneLabel.SetText("患者电话:"+preItem.Pre.Phone)
		        addressLabel.SetText("现居地址:"+preItem.Pre.Address)
		        ageLabel.SetText("患者年龄:"+preItem.Pre.Age)
		        diagnosisLabel.SetText("诊断结果:"+preItem.Pre.Diagnosis)
		        recipeLabel.SetText("治疗方案:"+preItem.Pre.Recipe)
		        hospitalNameLabel.SetText("医院:"+preItem.Pre.HospitalName)
		        physicianNameLabel.SetText("主治医师:"+preItem.Pre.PhysicianName)
		        feeLabel.SetText("费用:"+preItem.Pre.Fee)
		        pharmacistNameLabel.SetText("药房:"+preItem.Pre.PharmacistName)
		        pharmacyLabel.SetText("药剂师:"+preItem.Pre.PharmacyName)
		        paymentLabel.SetText("支付方式:"+preItem.Pre.Payment)
		        insuranceLabel.SetText("保险公司:"+preItem.Pre.InsuranceCompany)
		        policyNumberLabel.SetText("保险单号:"+preItem.Pre.PolicyNumber)
		    }
	    }
	})

	return vbox
}

func (app *Application) pharmacyMakeControlWritePage() ui.Control{
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
	refreshButton := ui.NewButton("刷新")
    hbox.Append(refreshButton, false)
    submitButton := ui.NewButton("执行处方")
	hbox.Append(submitButton, false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("写入处方")
	group.SetMargined(true)
	vbox.Append(group, true)

    group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
    
    preIdLabel := ui.NewLabel("处方号: ")
	preStateLabel := ui.NewLabel("处方状态： ")
	certHashLabel := ui.NewLabel("签名公钥地址： ")
	signatureHashLabel := ui.NewLabel("处方存储地址： ")
	surgeryDateLabel := ui.NewLabel("就诊日期: ")
	patientNameLabel := ui.NewLabel("患者姓名: ")
	genderLabel := ui.NewLabel("患者性别: ")
	iDNumberLabel := ui.NewLabel("身份证号: ")
	phoneLabel := ui.NewLabel("患者电话: ")
	addressLabel := ui.NewLabel("现居地址: ")
	ageLabel := ui.NewLabel("患者年龄: ")
	diagnosisLabel := ui.NewLabel("诊断结果: ")
	recipeLabel := ui.NewLabel("治疗方案: ")
	hospitalNameLabel := ui.NewLabel("医院: ")
	physicianNameLabel := ui.NewLabel("主治医师: ")
	feeLabel := ui.NewLabel("费用: ")
	pharmacistNameLabel := ui.NewLabel("药房: ")
	pharmacyLabel := ui.NewLabel("药剂师: ")
	paymentLabel := ui.NewLabel("支付方式: ")
	insuranceLabel := ui.NewLabel("保险公司: ")
	policyNumberLabel := ui.NewLabel("保险单号: ")
	vbox.Append(preIdLabel, false)
	vbox.Append(preStateLabel, false)
	vbox.Append(certHashLabel, false)
	vbox.Append(signatureHashLabel, false)
	vbox.Append(surgeryDateLabel, false)
	vbox.Append(patientNameLabel, false)
	vbox.Append(genderLabel, false)
	vbox.Append(iDNumberLabel, false)
	vbox.Append(phoneLabel, false)
	vbox.Append(addressLabel, false)
	vbox.Append(ageLabel, false)
	vbox.Append(diagnosisLabel, false)
	vbox.Append(recipeLabel, false)
	vbox.Append(hospitalNameLabel, false)
	vbox.Append(physicianNameLabel, false)
	vbox.Append(feeLabel, false)
	vbox.Append(pharmacistNameLabel, false)
	vbox.Append(pharmacyLabel, false)
	vbox.Append(paymentLabel, false)
	vbox.Append(insuranceLabel, false)
	vbox.Append(policyNumberLabel, false)
		
	pharmacistNameEntry := ui.NewEntry()
	pharmacyEntry := ui.NewEntry()
	entryForm.Append("药房",pharmacistNameEntry,false)
	entryForm.Append("药剂师", pharmacyEntry, false)
    
	submitButton.OnClicked(func(*ui.Button){
		if pageNumber==-1 || currentPage > pageNumber{
			return
		}
		var preItem service.PrescriptionItem
	    var exist bool=false
	    for _,preItem = range currentPreItems{
			if preItem.PreHash.PreState == "Generated"{
				exist = true
				break
			}
		}
		if !exist{
			return
		}		
		pharmacistNameText := pharmacistNameEntry.Text()
		pharmacyText := pharmacyEntry.Text()
		if pharmacistNameText=="" || pharmacyText==""{
			ui.MsgBox(mainWindow,
				"关键信息不可为空",
				"请输入经手药剂师姓名和药房全名")
			return 
		}
		var pre service.Prescription
		pre = preItem.Pre
		pre.PharmacistName = pharmacistNameText
		pre.PharmacyName = pharmacyText
		pre.Payment = "自费"
		err:= app.Setup.ConductPrescription(pre,preItem)
		if err != nil{
			ui.MsgBox(mainWindow,
				"执行处方失败",
				err.Error())
			return
		} 
		if err == nil {
			ui.MsgBox(mainWindow,
				"处方执行成功",
				"执行数据已上传区块链")
		}
	})

    refreshButton.OnClicked(func(*ui.Button) {
		// fmt.Println("pageNumber:",pageNumber)
		// fmt.Println("currentPage:",currentPage)
		if pageNumber!=-1 && currentPage <= pageNumber{
			var preItem service.PrescriptionItem
	        var exist bool=false
	        for _,preItem = range currentPreItems{
				if preItem.PreHash.PreState == "Generated"{
				exist = true
				break
			    }
		    }
			// fmt.Println(preItem)
	        if exist {
				preIdLabel.SetText("处方号: "+preItem.PreHash.PreID)		
		        preStateLabel.SetText("处方状态： 已开具")
		        certHashLabel.SetText("签名公钥地址： "+preItem.PreHash.CertHash)
		        signatureHashLabel.SetText("处方存储地址： "+preItem.PreHash.SignatureHash)
		        surgeryDateLabel.SetText("就诊日期: "+preItem.Pre.SurgeryDate)
		        patientNameLabel.SetText("患者姓名: "+preItem.Pre.PatientName)
		        genderLabel.SetText("患者性别: "+preItem.Pre.Gender)
		        iDNumberLabel.SetText("身份证号: "+preItem.Pre.IDNumber)
		        phoneLabel.SetText("患者电话:"+preItem.Pre.Phone)
		        addressLabel.SetText("现居地址:"+preItem.Pre.Address)
		        ageLabel.SetText("患者年龄:"+preItem.Pre.Age)
		        diagnosisLabel.SetText("诊断结果:"+preItem.Pre.Diagnosis)
		        recipeLabel.SetText("治疗方案:"+preItem.Pre.Recipe)
		        hospitalNameLabel.SetText("医院:"+preItem.Pre.HospitalName)
		        physicianNameLabel.SetText("主治医师:"+preItem.Pre.PhysicianName)
		        feeLabel.SetText("费用:"+preItem.Pre.Fee)
		        pharmacistNameLabel.SetText("药房:"+preItem.Pre.PharmacistName)
		        pharmacyLabel.SetText("药剂师:"+preItem.Pre.PharmacyName)
		        paymentLabel.SetText("支付方式:"+preItem.Pre.Payment)
		        insuranceLabel.SetText("保险公司:"+preItem.Pre.InsuranceCompany)
		        policyNumberLabel.SetText("保险单号:"+preItem.Pre.PolicyNumber)
		    }
	    }
	})

	return vbox

}

func (app *Application) doctorMakeControlWritePage() ui.Control{
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
    submitButton := ui.NewButton("开具处方")
	hbox.Append(submitButton, false)

	vbox.Append(ui.NewHorizontalSeparator(), false)

	group := ui.NewGroup("写入处方")
	group.SetMargined(true)
	vbox.Append(group, true)

    group.SetChild(ui.NewNonWrappingMultilineEntry())

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
    
	preIDEntry := ui.NewEntry()
	surgeryDateEntry := ui.NewEntry()
	patientNameEntry := ui.NewEntry()
	genderEntry := ui.NewEntry()
	iDNumberEntry := ui.NewEntry()
	phoneEntry := ui.NewEntry()
	addressEntry := ui.NewEntry()
	ageEntry := ui.NewEntry()
	diagnosisEntry := ui.NewMultilineEntry()
	recipeEntry := ui.NewMultilineEntry()
	hospitalNameEntry := ui.NewEntry()
    physicianNameEntry := ui.NewEntry()
	feeEntry := ui.NewEntry()

	entryForm.Append("处方号",preIDEntry,false)
	entryForm.Append("就诊日期",surgeryDateEntry,false)
	entryForm.Append("患者姓名",patientNameEntry,false)
	entryForm.Append("患者性别",genderEntry,false)
	entryForm.Append("身份证号",iDNumberEntry,false)
	entryForm.Append("患者电话",phoneEntry,false)
	entryForm.Append("现居地址",addressEntry,false)
	entryForm.Append("患者年龄",ageEntry,false)
	entryForm.Append("诊断结果",diagnosisEntry,false)
	entryForm.Append("治疗方案",recipeEntry,false)
	entryForm.Append("医院",hospitalNameEntry,false)
	entryForm.Append("主治医师",physicianNameEntry,false)
	entryForm.Append("费用",feeEntry,false)

	submitButton.OnClicked(func(*ui.Button) {
		preIDText := preIDEntry.Text()
		surgeryDateText := surgeryDateEntry.Text()
		patientNameText := patientNameEntry.Text()
		genderText := genderEntry.Text()
		iDNumberText := iDNumberEntry.Text()
		phoneText := phoneEntry.Text()
		addressText := addressEntry.Text()
		ageText := ageEntry.Text()
		diagnosisText := diagnosisEntry.Text()
		recipeText := recipeEntry.Text()
		hospitalNameText := hospitalNameEntry.Text()
		physicianNameText := physicianNameEntry.Text()
		feeText := feeEntry.Text()
        if preIDText=="" || surgeryDateText=="" || patientNameText=="" || genderText=="" || ageText=="" || diagnosisText=="" || recipeText ==""{
			ui.MsgBox(mainWindow,
				"关键信息不可为空",
				"请重新开具处方")
			return 
		}
		var pre service.Prescription
		pre.PreID = preIDText
		pre.PatientName = patientNameText
		pre.Gender = genderText
		pre.Age = ageText
		pre.IDNumber = iDNumberText
		pre.Address = addressText
		pre.Phone = phoneText
		pre.SurgeryDate = surgeryDateText
		pre.Diagnosis = diagnosisText
		pre.Recipe = recipeText
		pre.Fee = feeText
		pre.PhysicianName =physicianNameText
		pre.HospitalName = hospitalNameText
		pre.PharmacistName = "*"
		pre.PharmacyName = "*"
		pre.Payment = "*"
		pre.PolicyNumber = "*"
		pre.InsuranceCompany = "*"
		err := app.Setup.GeneratePrescription(pre)
		if err != nil{
			ui.MsgBox(mainWindow,
				"处方 ["+preIDText+"] 开具失败",
				err.Error())
		}
		if err ==nil{
			ui.MsgBox(mainWindow,
				"处方 ["+preIDText+"] 开具成功",
				"处方已上传区块链")
		}
	})
	
    return vbox
}

func (app *Application) makeControlQueryPage() ui.Control{	
	var identity string="identity"
	var certHash string="Hash"
	var err error
    switch currentUser.Identity{
		case "Doctor":
			identity = "医师"
			err = app.Setup.ChangeIdentity("Doctor")
		case "Pharmacy":
			identity = "药剂师"
		    err = app.Setup.ChangeIdentity("Pharmacy")
		case "Insurance":
			identity = "保险公司"
		    err = app.Setup.ChangeIdentity("Insurance")
		default:
			err = fmt.Errorf("default idendity")
	}
	certHash = app.Setup.GetCertAddress()

	// if currentUser.Identity == "Doctor"{
	// 	identity = "医师"
	// 	certHash = "Doctor_hash"
	//     err = app.Setup.ChangeIdentity("Doctor")
	// }
	// if currentUser.Identity == "Pharmacy"{
	// 	identity = "药剂师"
	// 	certHash = "Pharmacy_hash"
	// 	err = app.Setup.ChangeIdentity("Pharmacy")
	// }
	// if currentUser.Identity == "Insurance" {
	// 	identity = "保险公司"
	// 	certHash = "Insurance_hash"
	// 	err = app.Setup.ChangeIdentity("Insurance")
	// }
    vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	if err!=nil {
		ui.MsgBoxError(mainWindow,
			"身份传递异常",
			err.Error())
		return vbox
	}
	vbox.Append(ui.NewLabel("用户账号:  "+currentUser.LoginName), false)
	vbox.Append(ui.NewLabel("身份:            "+identity), false)
	vbox.Append(ui.NewLabel("公钥地址:  "+certHash), false)
    hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)
	vbox.Append(hbox, false)
    backButton := ui.NewButton("返回")
	hbox.Append(backButton, false)
    nextPageButton := ui.NewButton("翻页")
	hbox.Append(nextPageButton, false)
	queryButton := ui.NewButton("查询")
	hbox.Append(queryButton, false)
	group := ui.NewGroup("处方历史溯源")
	group.SetMargined(true)
	vbox.Append(group, true)

    group.SetChild(ui.NewNonWrappingMultilineEntry())
    entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	group.SetChild(entryForm)
	preIDEntry := ui.NewSearchEntry()
    entryForm.Append("处方号查询", preIDEntry, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)
	//currentPage = 0
    //-------------//
	preIdLabel := ui.NewLabel("处方号: ")
	preStateLabel := ui.NewLabel("处方状态： ")
	certHashLabel := ui.NewLabel("签名公钥地址： ")
	signatureHashLabel := ui.NewLabel("处方存储地址： ")
	surgeryDateLabel := ui.NewLabel("就诊日期: ")
	patientNameLabel := ui.NewLabel("患者姓名: ")
	genderLabel := ui.NewLabel("患者性别: ")
	iDNumberLabel := ui.NewLabel("身份证号: ")
	phoneLabel := ui.NewLabel("患者电话: ")
	addressLabel := ui.NewLabel("现居地址: ")
	ageLabel := ui.NewLabel("患者年龄: ")
	diagnosisLabel := ui.NewLabel("诊断结果: ")
	recipeLabel := ui.NewLabel("治疗方案: ")
	hospitalNameLabel := ui.NewLabel("医院: ")
	physicianNameLabel := ui.NewLabel("主治医师: ")
	feeLabel := ui.NewLabel("费用: ")
	pharmacistNameLabel := ui.NewLabel("药房: ")
	pharmacyLabel := ui.NewLabel("药剂师: ")
	paymentLabel := ui.NewLabel("支付方式: ")
	insuranceLabel := ui.NewLabel("保险公司: ")
	policyNumberLabel := ui.NewLabel("保险单号: ")
	//-------------//
    // if pageNumber==-1 || currentPage > pageNumber{
	// 	//vbox.Append(ui.NewLabel("此处方未被区块链确认"), false)
	// 	return vbox
	// }
	// //currentPage=0
	vbox.Append(preIdLabel, false)
	vbox.Append(preStateLabel, false)
	vbox.Append(certHashLabel, false)
	vbox.Append(signatureHashLabel, false)
	vbox.Append(surgeryDateLabel, false)
	vbox.Append(patientNameLabel, false)
	vbox.Append(genderLabel, false)
	vbox.Append(iDNumberLabel, false)
	vbox.Append(phoneLabel, false)
	vbox.Append(addressLabel, false)
	vbox.Append(ageLabel, false)
	vbox.Append(diagnosisLabel, false)
	vbox.Append(recipeLabel, false)
	vbox.Append(hospitalNameLabel, false)
	vbox.Append(physicianNameLabel, false)
	vbox.Append(feeLabel, false)
	vbox.Append(pharmacistNameLabel, false)
	vbox.Append(pharmacyLabel, false)
	vbox.Append(paymentLabel, false)
	vbox.Append(insuranceLabel, false)
	vbox.Append(policyNumberLabel, false)

    queryButton.OnClicked(func(*ui.Button) {
		if preIDEntry.Text() == ""{
			return
		}
        result,err := app.Setup.FindPrescriptionByPreID(preIDEntry.Text())
		if err != nil{
			ui.MsgBoxError(mainWindow,
				"此处方未被区块链确认",err.Error())
		    return
		}
		// if pageNumber==-1 || currentPage > pageNumber{
		// 	ui.MsgBoxError(mainWindow,
		// 		"此处方未被区块链确认",err)
		//     return 
		// }
		currentPreItems = result
		currentPage = 0
		pageNumber = len(currentPreItems)-1
		// mainWindow.Destroy()
		// ui.Main(app.generalUI)
		preIdLabel.SetText("处方号: "+currentPreItems[currentPage].PreHash.PreID)
		switch currentPreItems[currentPage].PreHash.PreState {
		case "Reimbursed":
			preStateLabel.SetText("处方状态： 已报销")
		case "Conducted":
			preStateLabel.SetText("处方状态： 已执行")
	    default:
			preStateLabel.SetText("处方状态： 已开具")
		}	
		certHashLabel.SetText("签名公钥地址： "+currentPreItems[currentPage].PreHash.CertHash)
		signatureHashLabel.SetText("处方存储地址： "+currentPreItems[currentPage].PreHash.SignatureHash)
		surgeryDateLabel.SetText("就诊日期: "+currentPreItems[currentPage].Pre.SurgeryDate)
		patientNameLabel.SetText("患者姓名: "+currentPreItems[currentPage].Pre.PatientName)
		genderLabel.SetText("患者性别: "+currentPreItems[currentPage].Pre.Gender)
		iDNumberLabel.SetText("身份证号: "+currentPreItems[currentPage].Pre.IDNumber)
		phoneLabel.SetText("患者电话:"+currentPreItems[currentPage].Pre.Phone)
		addressLabel.SetText("现居地址:"+currentPreItems[currentPage].Pre.Address)
		ageLabel.SetText("患者年龄:"+currentPreItems[currentPage].Pre.Age)
		diagnosisLabel.SetText("诊断结果:"+currentPreItems[currentPage].Pre.Diagnosis)
		recipeLabel.SetText("治疗方案:"+currentPreItems[currentPage].Pre.Recipe)
		hospitalNameLabel.SetText("医院:"+currentPreItems[currentPage].Pre.HospitalName)
		physicianNameLabel.SetText("主治医师:"+currentPreItems[currentPage].Pre.PhysicianName)
		feeLabel.SetText("费用:"+currentPreItems[currentPage].Pre.Fee)
		pharmacistNameLabel.SetText("药房:"+currentPreItems[currentPage].Pre.PharmacistName)
		pharmacyLabel.SetText("药剂师:"+currentPreItems[currentPage].Pre.PharmacyName)
		paymentLabel.SetText("支付方式:"+currentPreItems[currentPage].Pre.Payment)
		insuranceLabel.SetText("保险公司:"+currentPreItems[currentPage].Pre.InsuranceCompany)
		policyNumberLabel.SetText("保险单号:"+currentPreItems[currentPage].Pre.PolicyNumber)

	})
    
    nextPageButton.OnClicked(func(*ui.Button) {
		if pageNumber == -1{
			return
		}
		currentPage = currentPage+1
		if currentPage > pageNumber{
			currentPage =0
		}		//return vbox
			
		preIdLabel.SetText("处方号: "+currentPreItems[currentPage].PreHash.PreID)
		switch currentPreItems[currentPage].PreHash.PreState {
		case "Reimbursed":
			preStateLabel.SetText("处方状态： 已报销")
		case "Conducted":
			preStateLabel.SetText("处方状态： 已执行")
	    default:
			preStateLabel.SetText("处方状态： 已开具")
		}	
		certHashLabel.SetText("签名公钥地址： "+currentPreItems[currentPage].PreHash.CertHash)
		signatureHashLabel.SetText("处方存储地址： "+currentPreItems[currentPage].PreHash.SignatureHash)
		surgeryDateLabel.SetText("就诊日期: "+currentPreItems[currentPage].Pre.SurgeryDate)
		patientNameLabel.SetText("患者姓名: "+currentPreItems[currentPage].Pre.PatientName)
		genderLabel.SetText("患者性别: "+currentPreItems[currentPage].Pre.Gender)
		iDNumberLabel.SetText("身份证号: "+currentPreItems[currentPage].Pre.IDNumber)
		phoneLabel.SetText("患者电话:"+currentPreItems[currentPage].Pre.Phone)
		addressLabel.SetText("现居地址:"+currentPreItems[currentPage].Pre.Address)
		ageLabel.SetText("患者年龄:"+currentPreItems[currentPage].Pre.Age)
		diagnosisLabel.SetText("诊断结果:"+currentPreItems[currentPage].Pre.Diagnosis)
		recipeLabel.SetText("治疗方案:"+currentPreItems[currentPage].Pre.Recipe)
		hospitalNameLabel.SetText("医院:"+currentPreItems[currentPage].Pre.HospitalName)
		physicianNameLabel.SetText("主治医师:"+currentPreItems[currentPage].Pre.PhysicianName)
		feeLabel.SetText("费用:"+currentPreItems[currentPage].Pre.Fee)
		pharmacistNameLabel.SetText("药房:"+currentPreItems[currentPage].Pre.PharmacistName)
		pharmacyLabel.SetText("药剂师:"+currentPreItems[currentPage].Pre.PharmacyName)
		paymentLabel.SetText("支付方式:"+currentPreItems[currentPage].Pre.Payment)
		insuranceLabel.SetText("保险公司:"+currentPreItems[currentPage].Pre.InsuranceCompany)
		policyNumberLabel.SetText("保险单号:"+currentPreItems[currentPage].Pre.PolicyNumber)	

	})

	backButton.OnClicked(func(*ui.Button) {
		//fmt.Println("help")
		currentPage = 0
		pageNumber = -1
		currentPreItems = []service.PrescriptionItem{}
		mainWindow.Destroy()
		ui.Main(app.setupUI)	
	})
	return vbox
}

func (app *Application) generalUI(){
	mainWindow = ui.NewWindow("电子处方流转系统", 640, 480, true)
	mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	tab := ui.NewTab()
	mainWindow.SetChild(tab)
	mainWindow.SetMargined(true)

	tab.Append("处方查询", app.makeControlQueryPage())
	tab.SetMargined(0, true)
    
    switch currentUser.Identity{
	case "Doctor":
		tab.Append("处方开立", app.doctorMakeControlWritePage())
	    tab.SetMargined(1, true)
	case "Pharmacy":
		tab.Append("处方执行", app.pharmacyMakeControlWritePage())
		tab.SetMargined(1, true)
	case "Insurance":
		tab.Append("处方报销", app.insuranceMakeControlWritePage())
		tab.SetMargined(1, true)
	default:
		break
	}

	// if currentUser.Identity == "Doctor"{
	//     tab.Append("处方开立", app.doctorMakeControlWritePage())
	//     tab.SetMargined(1, true)
	// }

	// if currentUser.Identity == "Pharmacy"{
	//     tab.Append("处方执行", app.pharmacyMakeControlWritePage())
	// 	tab.SetMargined(1, true)
	// }

	// if currentUser.Identity == "Insurance"{
	//     tab.Append("处方报销", app.insuranceMakeControlWritePage())
	// 	tab.SetMargined(1, true)
	// }

	mainWindow.Show()
}

func (app *Application) setupUI() {
	mainWindow = ui.NewWindow("登录", 640, 480, true)
	mainWindow.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainWindow.Destroy()
		return true
	})

	vbContainer := ui.NewVerticalBox()
	vbContainer.SetPadded(true)

	inputGroup := ui.NewGroup("电子处方流转系统")
	inputGroup.SetMargined(true)

	vbInput := ui.NewVerticalBox()
	vbInput.SetPadded(true)

	inputForm := ui.NewForm()
	inputForm.SetPadded(true)

	IDNumber := ui.NewEntry()
	IDNumber.SetText("用户ID/证件号/邮箱")
	inputForm.Append("用户账号", IDNumber, false)

	Password := ui.NewPasswordEntry()
	//Password.SetText("请输入密码")
	inputForm.Append("密码", Password, false)

	logInButton := ui.NewButton("登录")

	vbInput.Append(inputForm, false)
	vbInput.Append(logInButton, false)

	inputGroup.SetChild(vbInput)

	vbContainer.Append(inputGroup, false)	

	mainWindow.SetChild(vbContainer)

	logInButton.OnClicked(func(*ui.Button) {
		// Update the UI directly as it is called from the main thread
		// messageLabel.SetText(IDNumber.Text())
		// fmt.Println(IDNumber.Text())
		// fmt.Println(Password.Text())
		flag := false
		loginName := IDNumber.Text()
	    password := Password.Text()
		for _, user:= range Users{
			if user.LoginName == loginName && user.Password == password {
				currentUser = user
				fmt.Println(currentUser)
				flag = true
				break
			}
		}
		if !flag{
			ui.MsgBox(mainWindow,
				"此身份尚未加入联盟链",
				"请先申请颁发数字证书")
		}
		if flag{
			mainWindow.Destroy()
			ui.Main(app.generalUI)
			//return true
		}
	})

	mainWindow.Show()

}

func UiStart(app Application){
	// preHash := []PrescriptionHash{
	// 	{PreID:"KC0000001", PreState: "Genrated",   CertHash: "QmRYyR4m8CntJdvmMisuzXPCM8wqHDSTJSwcQaU3zSDaq7", SignatureHash: "QmZTGPUUTPNWNTfNryBExHfNm1732WKLjuQtqAeRfNbf1e"},
	// 	{PreID:"KC0000001", PreState: "Conducted",  CertHash: "QmSpQCBBAi9uji1s7GHDTjrAThrNXv62WcMrbbK3gvkGZM", SignatureHash: "QmPxd2PpwqUzpczMbvBx3EGbcQE386ax3hiHdiYqXjT8Fy"},
	// 	{PreID:"KC0000001", PreState: "Reimbursed", CertHash: "QmXM5i7mDHjGPjzmXcCUofSAJhiAZZFfTeK5bFqmePRbXF", SignatureHash: "QmPxd2PpwqUzpczMbvBx3EGbcQE386ax3hiHdiYqXjT8Fy"},
	//  }
	//  prescription := []Prescription{
	// 	{
	// 		PreID: "KC0000001",	
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
	// 	},
	// 	{
	// 		PreID: "KC0000001",	
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
	// 	},
	// 	{
	// 		PreID: "KC0000001",	
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
	// 	},
	// }
	
	
    // var preItem PrescriptionItem
	// preItem.Pre = prescription[0]
	// preItem.PreHash = preHash[0]
	// currentPreItems = append(currentPreItems,preItem)
	// preItem.Pre = prescription[1]
	// preItem.PreHash = preHash[1]
	// currentPreItems = append(currentPreItems,preItem)
	// preItem.Pre = prescription[2]
	// preItem.PreHash = preHash[2]
	// currentPreItems = append(currentPreItems,preItem)
	
	Init()
	pageNumber = -1
	currentPage = 0
	currentPreItems = []service.PrescriptionItem{}
	ui.Main(app.setupUI)
}
