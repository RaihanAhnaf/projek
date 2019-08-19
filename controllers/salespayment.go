package controllers

import (
	"../helpers"
	. "../helpers"
	. "../models"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

func (c *TransactionController) SalesPayment(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *TransactionController) GetLastNumberSP(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// m := time.Now().UTC().Month()
	y := time.Now().UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceSP").Select().Where(dbox.And(dbox.Eq("collname", "salespayment"),
		dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequencePurchasePaymentModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequencePurchasePaymentModel()
		model.Collname = "salespayment"
		model.TypePo = "salespayment"
		model.Lastnumber = 0
		model.Month = 0
		model.Year = y
		// e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data
	}
	sec := result[0]
	sec.Lastnumber = sec.Lastnumber + 1
	data.Number = sec.Lastnumber
	data.Msg = "Success"

	return data
}
func (c *TransactionController) GetDataInvoiceForSP(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart string
		DateEnd   string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var pipes []tk.M
	// dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	// End, _ := time.Parse("2006-01-02", p.DateEnd)
	// dateEnd := End.AddDate(0, 0, 1)

	// pipes = append(pipes, tk.M{"$match": tk.M{"DateCreated": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Status": "POSTING"}})

	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("InvoiceNonInv").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransactionController) SaveSalesPayment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	e := k.Request.ParseMultipartForm(1000000000)
	if e != nil {
		c.ErrorResultInfo(e.Error(), nil)
	}
	payload := new(tk.M)
	_, formData, err := k.GetPayloadMultipart(payload)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := SalesPaymentModel{}
	tk.UnjsonFromString(formData["data"][0], &data)
	p := data
	lastnumber, _ := strconv.Atoi(formData["LastNumber"][0])
	m := data.DatePosting.UTC().Month()
	y := data.DatePosting.UTC().Year()
	folder := tk.Sprintf("%d%02d", y, m)
	//Create Directory
	baseImagePath := ReadConfig()["uploadpath"]
	pathfolder := filepath.Join(baseImagePath, folder)
	if _, err = os.Stat(pathfolder); os.IsNotExist(err) {
		os.MkdirAll(pathfolder, 0777)
	}
	file, handler, err := k.Request.FormFile("fileUpload")
	if file != nil {
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		defer file.Close()
		fileName := "sp-" + tk.RandomString(6) + filepath.Ext(handler.Filename)
		filePath := filepath.Join(pathfolder, fileName)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		p.Attachment = fileName
	}
	// p := t.Data
	// m := p.DatePosting.UTC().Month()
	// y := p.DatePosting.UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceSP").Select().Where(dbox.And(dbox.Eq("collname", "salespayment"),
		dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequencePurchasePaymentModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	secSP := 0
	if len(result) == 0 {
		mod := NewSequencePurchasePaymentModel()
		mod.Collname = "salespayment"
		mod.TypePo = "salespayment"
		mod.Lastnumber = 1
		mod.Month = 0
		mod.Year = y
		e = c.Ctx.Save(mod)
		secSP = 1
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		c.Ctx.Save(&sec)
		secSP = sec.Lastnumber
	}
	strNum := ""
	if secSP < 10 {
		strNum = "000"
	} else if secSP >= 10 && secSP < 100 {
		strNum = "00"
	} else if secSP >= 100 && secSP < 1000 {
		strNum = "0"
	}

	model := NewSalesPaymentModel()

	if p.ID == "" {
		p.ID = bson.NewObjectId()
		c.SaveLastNumber(lastnumber)
	}
	model.ID = p.ID
	model.DateStr = p.DateStr
	model.DatePosting = p.DatePosting
	dateFormat := time.Now().Format("02012006")
	model.DocumentNumber = "PP/" + dateFormat + "/" + strNum + strconv.Itoa(secSP)
	model.CustomerCode = p.CustomerCode
	model.CustomerName = p.CustomerName
	model.PaymentAccount = p.PaymentAccount
	model.PaymentName = p.PaymentName
	if p.Attachment != "" {
		model.Attachment = p.Attachment
	}
	model.ListDetail = p.ListDetail
	model.User = k.Session("username").(string)
	model.StoreLocationId = k.Session("locationid").(int)
	model.StoreLocationName = k.Session("locationname").(string)
	totalBalanceAll := 0.0
	if p.ListDetail[0].Id == "" {
		for key, _ := range p.ListDetail {
			// tk.Println("=======Receive=========")
			// tk.Println(p.ListDetail[key].Receive)
			// tk.Println("========AlreadyPaid========")
			// tk.Println(p.ListDetail[key].AlreadyPaid)
			idList := model.DocumentNumber + "/" + strconv.Itoa(key)
			p.ListDetail[key].Id = idList
			p.ListDetail[key].AlreadyPaid += p.ListDetail[key].Receive
			p.ListDetail[key].Balance = p.ListDetail[key].Amount - p.ListDetail[key].AlreadyPaid
			totalBalanceAll += (p.ListDetail[key].Amount - p.ListDetail[key].AlreadyPaid)
		}
	}
	model.BalanceAll = totalBalanceAll
	model.ListDetail = p.ListDetail
	e = c.Ctx.Save(model)
	if e != nil {
		return c.SetResultInfo(true, "ERROR", nil)
	}
	c.SaveSPtoJournal(p.ListDetail, k.Session("username").(string), p.PaymentAccount, p.PaymentName, p.DocumentNumber, p.DatePosting)
	c.SaveSPtoINV(p.ListDetail, k.Session("username").(string), p.DocumentNumber, p.DatePosting, model.DocumentNumber)

	if p.ID == "" {
		c.LogActivity("Sales Payment", "Insert Salespayment", p.DocumentNumber, k)
	} else {
		c.LogActivity("Sales Payment", "Update Salespayment", p.DocumentNumber, k)
	}

	return c.SetResultOK(nil)
}
func (c *TransactionController) SaveSPtoJournal(ListDetail []DetailSalesPayment, User string, accCode int, accName string, DocNumber string, DatePosting time.Time) interface{} {
	m := DatePosting.UTC().Month()
	y := DatePosting.UTC().Year()
	codejurnal := tk.Sprintf("%02d%d", m, y)
	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
	number := tk.Sprintf("%04d", ids)
	mdl := new(MainJournal)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = "JUR/" + codejurnal + "/" + number
	mdl.PostingDate = DatePosting
	mdl.CreateDate = time.Now()
	mdl.DateStr = DatePosting.Format("02 Jan 2006")
	mdl.Journal_Type = "CashIn"
	mdl.Department = "COMMERCE"
	mdl.Status = "posting"
	mdl.User = User
	dataListDetail := []Journal{}
	totalForCredit := 0.0
	idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "BBM/" + DatePosting.Format("020106") + "/"
	docNumberJournal := headDOC + numberDOC

	thisSalesCode := ""
	thisSalesName := ""
	for i, each := range ListDetail {
		if each.Pay == true {
			resultINV := c.GetInvoiceBasedOnDocNum(each.InvNumber)
			list := Journal{}
			list.Id = tk.RandomString(32)
			list.No = i + 1
			list.PostingDate = mdl.PostingDate
			list.DateStr = mdl.DateStr
			list.Journal_Type = mdl.Journal_Type
			list.Debet = each.Receive
			list.User = User
			list.DocumentNumber = docNumberJournal
			list.Acc_Code = accCode
			list.Description = each.InvNumber
			list.Acc_Name = accName
			list.Department = "COMMERCE"
			list.SalesCode = resultINV[0].SalesCode
			list.SalesName = resultINV[0].SalesName
			totalForCredit += each.Receive
			dataListDetail = append(dataListDetail, list)

			if i == 0 {
				thisSalesCode = resultINV[0].SalesCode
				thisSalesName = resultINV[0].SalesName
			}
		}
	}
	list2 := Journal{}
	list2.Id = tk.RandomString(32)
	list2.No = len(dataListDetail) + 1
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = docNumberJournal
	list2.Credit = totalForCredit
	list2.User = User
	list2.Acc_Code = 1210
	list2.Acc_Name = "PIUTANG DAGANG"
	list2.Description = DocNumber
	list2.Department = "COMMERCE"
	list2.SalesCode = thisSalesCode
	list2.SalesName = thisSalesName
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	c.SaveSPtoGeneralLedger(ListDetail, User, accCode, accName, DocNumber, mdl.IdJournal, docNumberJournal, DatePosting)
	return ""
}
func (c *TransactionController) SaveSPtoGeneralLedger(ListDetail []DetailSalesPayment, User string, accCode int, accName string, DocNumber string, idJournal string, docNumberJournal string, DatePosting time.Time) interface{} {
	m := DatePosting.UTC().Month()
	y := DatePosting.UTC().Year()
	mdl := new(MainGeneralLedger)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = idJournal
	mdl.PostingDate = DatePosting
	mdl.CreateDate = time.Now()
	mdl.DateStr = DatePosting.Format("02 Jan 2006")
	mdl.Journal_Type = "CashIn"
	mdl.Department = "COMMERCE"
	if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
		mdl.Status = "posting"
	} else {
		mdl.Status = "close"
	}
	mdl.User = User
	dataListDetail := []GeneralDetail{}
	totalForCredit := 0.0
	idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "BBM/" + DatePosting.Format("020106") + "/"
	docNumberGL := headDOC + numberDOC

	thisSalesCode := ""
	thisSalesName := ""
	for i, each := range ListDetail {
		if each.Pay == true {
			resultINV := c.GetInvoiceBasedOnDocNum(each.InvNumber)

			list := GeneralDetail{}
			list.Id = tk.RandomString(32)
			list.No = i + 1
			list.PostingDate = mdl.PostingDate
			list.DateStr = mdl.DateStr
			list.Journal_Type = mdl.Journal_Type
			list.Debet = each.Receive
			list.User = User
			list.DocumentNumber = docNumberGL
			list.Acc_Code = accCode
			list.Description = each.InvNumber
			list.Acc_Name = accName
			list.Department = "COMMERCE"
			list.SalesCode = resultINV[0].SalesCode
			list.SalesName = resultINV[0].SalesName
			totalForCredit += each.Receive
			dataListDetail = append(dataListDetail, list)

			if i == 0 {
				thisSalesCode = resultINV[0].SalesCode
				thisSalesName = resultINV[0].SalesName
			}
		}
	}
	list2 := GeneralDetail{}
	list2.Id = tk.RandomString(32)
	list2.No = len(dataListDetail) + 1
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = docNumberGL
	list2.Credit = totalForCredit
	list2.User = User
	list2.Acc_Code = 1210
	list2.Acc_Name = "PIUTANG DAGANG"
	list2.Description = DocNumber
	list2.Department = "COMMERCE"
	list2.SalesCode = thisSalesCode
	list2.SalesName = thisSalesName
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	return ""
}

func (c *TransactionController) SaveSPtoINV(ListDetail []DetailSalesPayment, User string, DocNumber string, DatePost time.Time, docNoSP string) interface{} {
	for _, each := range ListDetail {
		result := ItemPaymentsNonInv{}
		if each.Pay == true {
			var pipes []tk.M
			pipes = append(pipes, tk.M{"$match": tk.M{"DocumentNo": each.InvNumber}})
			csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("InvoiceNonInv").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()

			DataArr := []InvoiceNonInvModel{}
			err := csr.Fetch(&DataArr, 0, false)
			if err != nil {
				return err.Error()
			}
			csr.Close()

			temp := []ItemPaymentsNonInv{}
			if len(DataArr) > 0 {
				if len(DataArr[0].ListPayment) > 0 {
					for _, ea := range DataArr[0].ListPayment {
						res := ItemPaymentsNonInv{}
						res.Id = ea.Id
						res.DatePayment = ea.DatePayment
						res.DocumentPayment = ea.DocumentPayment
						res.PaymentAmount = ea.PaymentAmount
						temp = append(temp, res)
					}
				}
				result.Id = bson.NewObjectId()
				result.DatePayment = DatePost
				result.DocumentPayment = each.InvNumber
				result.PaymentAmount = each.Receive
				temp = append(temp, result)
				DataArr[0].ListPayment = temp
				DataArr[0].AlreadyPaid = DataArr[0].AlreadyPaid + each.Receive
				if each.Balance == 0 {
					DataArr[0].Status = "PAID"
				}
				erro := c.Ctx.Save(&DataArr[0])
				if erro != nil {
					return c.SetResultInfo(true, "ERROR", nil)
				}
			}
			history := HistoryTrackInvoice{}
			history.Id = bson.NewObjectId()
			history.DocumentNumber = each.InvNumber
			history.DateCreated = time.Now()
			history.DateStr = DatePost.Format("02-Jan-2006")
			history.DateSP = DatePost
			history.Status = "SP Pending"
			if each.Balance == 0 {
				history.Status = "SP Paid"
			}
			history.DocumentNumberSP = docNoSP
			history.Remark = DataArr[0].Description
			history.CustomerCode = DataArr[0].CustomerCode
			history.CustomerName = DataArr[0].CustomerName

			crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNumber", each.InvNumber)).From("TrackingInvoice").Cursor(nil)
			defer crs.Close()

			resultINV := []TrackInvoiceModel{}
			e = crs.Fetch(&resultINV, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

			mod := resultINV[0]
			inv := NewTrackInvoiceModel()
			inv.ID = mod.ID
			inv.DocumentNumber = mod.DocumentNumber
			inv.DateCreated = mod.DateCreated
			inv.DateStr = mod.DateStr
			inv.DateINV = mod.DateINV
			inv.DateSP = DatePost
			inv.Status = "SP Pending"
			if each.Balance == 0 {
				inv.Status = "SP Paid"
			}
			inv.Remark = mod.Remark
			inv.CustomerCode = mod.CustomerCode
			inv.CustomerName = mod.CustomerName
			history.DateINV = mod.DateINV
			inv.History = mod.History
			inv.History = append(inv.History, history)
			e = c.Ctx.Save(inv)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
		}
	}

	return c.SetResultOK(nil)
}

func (c *TransactionController) GetAllDataSalesPayment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)
	p := struct {
		DateStart    time.Time
		DateEnd      time.Time
		CustomerCode string
		TextSearch   string
		Filter       bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Ne("Status", ""))
	filter = append(filter, dbox.Gte("DatePosting", dateStart))
	filter = append(filter, dbox.Lt("DatePosting", dateEnd))
	if p.Filter == true {
		if p.TextSearch != "" {
			// filter = append(filter, dbox.Contains("DocumentNumber", p.TextSearch))
			filter = append(filter, dbox.Or(dbox.Contains("StoreLocationName", p.TextSearch), dbox.Contains("DocumentNumber", p.TextSearch), dbox.Contains("CustomerName", p.TextSearch)))
		}
		if p.CustomerCode != "" {
			filter = append(filter, dbox.Eq("CustomerCode", p.CustomerCode))
		}
		filter = CreateLocationFilter(filter, "StoreLocationId", locid, false)
	}
	query := tk.M{}.Set("where", dbox.And(filter...))
	csr, e := c.Ctx.Find(new(SalesPaymentModel), query)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if p.Filter && len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}

	return c.SetResultInfo(false, "Success", results)
}

func (c *TransactionController) SaveLastNumber(LastNumber int) interface{} {
	// m := time.Now().UTC().Month()
	y := time.Now().UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceSP").Select().Where(dbox.And(dbox.Eq("collname", "salespayment"),
		dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequencePurchasePaymentModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	sec := result[0]
	sec.Lastnumber = LastNumber
	c.Ctx.Save(&sec)

	return c.SetResultInfo(false, "Succes", sec)
}

func (c *TransactionController) ExportToPdfSalesPayment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Data SalesPaymentModel
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	// tk.Println(p.Data)

	// csr, e := c.Ctx.Connection.NewQuery().Select().From("SalesPayment").Where(dbox.Eq("_id", p.Data.ID)).Cursor(nil)
	// if e != nil {
	// 	tk.Println("query", e.Error())
	// }
	// defer csr.Close()
	results := p.Data
	// e = p.Data.Fetch(&results, 0, false)
	// if e != nil {
	// 	tk.Println("fetch", e.Error())
	// }
	DATA := results
	// tk.Println(DATA)

	// config := helpers.ReadConfig()
	// Img := config["imgpath"]
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	SalesPaymentHeadTable := []string{"No", "Date", "Invoice Number", "Amount", "Paid", "Receive", "Balance", "Status"}
	pdf.AddPage()
	pdf.SetXY(8, 5)

	pdf.SetFont("Century_Gothic", "", 12)
	// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(30)
	pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(30)
	pdf.SetFont("Century_Gothic", "", 12)
	pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
	pdf.SetX(180)
	pdf.SetFont("Century_Gothic", "", 16)
	pdf.CellFormat(0, 15, "SALES PAYMENT", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
	pdf.SetX(10)
	pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
	pdf.SetX(10)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+25, 286, pdf.GetY()+25) //garis horizontal2

	// pdf.Line(pdf.GetX(), pdf.GetY()+9, pdf.GetX(), 80)         // garis vertikal 1 dari kiri/156 -> 136
	// pdf.Line(pdf.GetX()+276, pdf.GetY()+9, pdf.GetX()+276, 75) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()+138, pdf.GetY()+9, pdf.GetX()+138, 55) // garis vertikal 2 tengah
	// pdf.Line(pdf.GetX()+69, pdf.GetY()+30, pdf.GetX()+69, 91)   // garis vertikal 4 tenga2nya field to//101
	y00 := pdf.GetY() + 9
	y2 := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(y2 + 10)
	pdf.MultiCell(12, 5, " To : ", "", "L", false)

	// pdf.SetX(12)
	pdf.SetY(y2 + 10)
	pdf.SetX(20)
	pdf.MultiCell(130, 5, DATA.CustomerName, "", "L", false)

	pdf.SetY(y2 + 26)
	pdf.SetX(10)
	pdf.MultiCell(138.5, 5, " Payment :", "", "L", false)
	// //TABLE KANAN
	pdf.SetY(y2 + 10)
	pdf.SetX(147.5)
	pdf.MultiCell(30, 5, " Payment No", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(175.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(178.5)
	pdf.MultiCell(90, 5, DATA.DocumentNumber, "", "L", false)
	pdf.SetY(pdf.GetY())
	// pdf.GetY()

	pdf.SetY(y2 + 15)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Date", "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetY(y2 + 15)
	pdf.SetX(175.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 15)
	pdf.SetX(178.5)
	pdf.MultiCell(90, 5, DATA.DateStr, "", "L", false)
	pdf.SetY(pdf.GetY())

	// csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(dbox.Eq("main_acc_code", 1100)).Cursor(nil)
	// defer csr.Close()
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }
	// coadata := []CoaModel{}
	// e = csr.Fetch(&coadata, 0, false)
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }
	// pdf.SetX(13.0)
	// x0 := 10.0
	// pdf.SetY(y2 + 31)
	// for i, each := range coadata {
	// 	if i == 0 {
	// 		pdf.SetX(13.0)
	// 	}
	// 	if i > 0 && i%4 == 0 {
	// 		pdf.Ln(8.0)
	// 		pdf.SetX(13.0)
	// 	}
	// 	filled := false
	// 	if DATA.PaymentAccount == each.ACC_Code {
	// 		filled = true
	// 	}
	// 	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filled, 0, "")
	// 	pdf.CellFormat(50.0, 5, each.Account_Name, "", 0, "L", false, 0, "")
	// }
	x0 := 10.0
	pdf.SetY(y2 + 31)
	pdf.SetX(13.0)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50.0, 5, DATA.PaymentName, "", 0, "L", false, 0, "")
	pdf.Ln(8.0)

	// pdf.GetY()

	y0 := pdf.GetY()
	widthHead := []float64{10.0, 28.0, 40.0, 45.0, 45.0, 45.0, 45.0, 18.0}
	for i, head := range SalesPaymentHeadTable {
		pdf.SetY(y0)
		x := 10.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		}

	}
	// pdf.Ln()
	pdf.SetFont("Century_Gothic", "", 9)
	y3 := pdf.GetY()
	for i, list := range DATA.ListDetail {

		y3 = pdf.GetY()
		pdf.SetY(y3)
		x := 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(widthHead[0], 5, numberstr, "", "C", false)
		// pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)
		a0 := pdf.GetY()
		pdf.SetY(y3)
		x += widthHead[0]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[1], 5, list.DatePayment.Local().Format("02 Jan 2006"), "", "L", false)

		a1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[2], 5, list.InvNumber, "", "L", false)

		a2 := pdf.GetY()
		pdf.SetY(y3)
		// tk.Println(list.Acc_Code)
		x = 10 + widthHead[0] + widthHead[1] + widthHead[2]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		amount := tk.Sprintf("%.2f", list.Amount)
		amount = c.ConvertToCurrency(amount)
		pdf.MultiCell(widthHead[3], 5, amount, "", "R", false)

		a3 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		alreadypaid := tk.Sprintf("%.2f", list.AlreadyPaid)
		alreadypaid = c.ConvertToCurrency(alreadypaid)
		pdf.MultiCell(widthHead[4], 5, alreadypaid, "", "R", false)

		a4 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		receive := tk.Sprintf("%.2f", list.Receive)
		receive = c.ConvertToCurrency(receive)
		pdf.MultiCell(widthHead[5], 5, receive, "", "R", false)

		a5 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		balance := tk.Sprintf("%.2f", list.Balance)
		balance = c.ConvertToCurrency(balance)
		pdf.MultiCell(widthHead[6], 5, balance, "", "R", false)

		a6 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		if list.Pay {
			tk.Sprintf("%.2f", list.Pay)
			pdf.MultiCell(widthHead[7], 5, "Paid", "", "C", false)
		}

		a7 := pdf.GetY()
		pdf.SetY(y3)

		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}

	y4 := pdf.GetY()
	pdf.Line(pdf.GetX()-258, y3, pdf.GetX()-258, y4) //vertical 1
	pdf.Line(pdf.GetX()-248, y3, pdf.GetX()-248, y4) //vertical 2
	pdf.Line(pdf.GetX()-220, y3, pdf.GetX()-220, y4) //vertical 3
	pdf.Line(pdf.GetX()-180, y3, pdf.GetX()-180, y4) //vertical 4
	pdf.Line(pdf.GetX()-135, y3, pdf.GetX()-135, y4) //vertical 5
	pdf.Line(pdf.GetX()-90, y3, pdf.GetX()-90, y4)   //vertical 6
	pdf.Line(pdf.GetX()-45, y3, pdf.GetX()-45, y4)   //vertical 7
	pdf.Line(pdf.GetX(), y3, pdf.GetX(), y4)         //vertical 8
	pdf.Line(pdf.GetX()+18, y3, pdf.GetX()+18, y4)   //vertical 9

	pdf.Line(x0, y00, x0, y4)                                                                                                                                                                                                                     // garis vertikal 1 dari kiri/156 -> 136
	pdf.Line(10.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y00, 10.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y4) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()-258, y4, 286, y4)

	e = os.RemoveAll(c.PdfPath + "/salespayment")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/salespayment", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}

	namepdf := "-salespayment.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/salespayment"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	// tk.Println(e)
	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName

}

func (c *TransactionController) ExportToPdfSalesPaymentListView(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id bson.ObjectId
		// Data SalesPaymentModel
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	// tk.Println(p.Data)

	csr, e := c.Ctx.Connection.NewQuery().Select().From("SalesPayment").Where(dbox.Eq("_id", p.Id)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []SalesPaymentModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	// tk.Println(DATA)

	// config := helpers.ReadConfig()
	// Img := config["imgpath"]
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	SalesPaymentHeadTable := []string{"No", "Date", "Invoice Number", "Amount", "Paid", "Receive", "Balance", "Status"}
	pdf.AddPage()
	pdf.SetXY(8, 5)

	pdf.SetFont("Century_Gothic", "", 12)
	// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(11)
	pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(11)
	pdf.SetFont("Century_Gothic", "", 12)
	pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
	pdf.SetX(180)
	pdf.SetFont("Century_Gothic", "", 16)
	pdf.CellFormat(0, 15, "SALES PAYMENT", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
	pdf.SetX(10)
	pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11)
	pdf.SetFont("Century_Gothic", "", 12)
	pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
	pdf.SetX(10)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+25, 286, pdf.GetY()+25) //garis horizontal2

	// pdf.Line(pdf.GetX(), pdf.GetY()+9, pdf.GetX(), 80)         // garis vertikal 1 dari kiri/156 -> 136
	// pdf.Line(pdf.GetX()+276, pdf.GetY()+9, pdf.GetX()+276, 75) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()+138, pdf.GetY()+9, pdf.GetX()+138, 55) // garis vertikal 2 tengah
	// pdf.Line(pdf.GetX()+69, pdf.GetY()+30, pdf.GetX()+69, 91)   // garis vertikal 4 tenga2nya field to//101
	y00 := pdf.GetY() + 9
	y2 := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(y2 + 10)
	pdf.MultiCell(12, 5, " To : ", "", "L", false)

	// pdf.SetX(12)
	pdf.SetY(y2 + 10)
	pdf.SetX(20)
	pdf.MultiCell(130, 5, DATA.CustomerName, "", "L", false)

	pdf.SetY(y2 + 26)
	pdf.SetX(10)
	pdf.MultiCell(138.5, 5, " Payment :", "", "L", false)
	// //TABLE KANAN
	pdf.SetY(y2 + 10)
	pdf.SetX(147.5)
	pdf.MultiCell(30, 5, " Payment No", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(175.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(178.5)
	pdf.MultiCell(90, 5, DATA.DocumentNumber, "", "L", false)
	pdf.SetY(pdf.GetY())
	// pdf.GetY()

	pdf.SetY(y2 + 15)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Date", "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetY(y2 + 15)
	pdf.SetX(175.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 15)
	pdf.SetX(178.5)
	pdf.MultiCell(90, 5, DATA.DateStr, "", "L", false)
	pdf.SetY(pdf.GetY())

	// csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa").Where(dbox.Eq("main_acc_code", 1100)).Cursor(nil)
	// defer csr.Close()
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }
	// coadata := []CoaModel{}
	// e = csr.Fetch(&coadata, 0, false)
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }
	// pdf.SetX(13.0)
	// x0 := 10.0
	// pdf.SetY(y2 + 31)
	// for i, each := range coadata {
	// 	if i == 0 {
	// 		pdf.SetX(13.0)
	// 	}
	// 	if i > 0 && i%4 == 0 {
	// 		pdf.Ln(8.0)
	// 		pdf.SetX(13.0)
	// 	}
	// 	filled := false
	// 	if DATA.PaymentAccount == each.ACC_Code {
	// 		filled = true
	// 	}
	// 	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filled, 0, "")
	// 	pdf.CellFormat(50.0, 5, each.Account_Name, "", 0, "L", false, 0, "")
	// }
	// pdf.Ln(8.0)
	x0 := 10.0
	pdf.SetY(y2 + 31)
	pdf.SetX(13.0)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50.0, 5, DATA.PaymentName, "", 0, "L", false, 0, "")
	pdf.Ln(8.0)

	// pdf.GetY()
	pdf.SetFont("Century_Gothic", "", 11)
	y0 := pdf.GetY()
	widthHead := []float64{10.0, 28.0, 40.0, 45.0, 45.0, 45.0, 45.0, 18.0}
	for i, head := range SalesPaymentHeadTable {
		pdf.SetY(y0)
		x := 10.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		}

	}
	// pdf.Ln()
	y3 := pdf.GetY()
	for i, list := range DATA.ListDetail {

		y3 = pdf.GetY()
		pdf.SetY(y3)
		x := 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(widthHead[0], 5, numberstr, "", "C", false)
		// pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)
		a0 := pdf.GetY()
		pdf.SetY(y3)
		x += widthHead[0]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[1], 5, list.DatePayment.Local().Format("02 Jan 2006"), "", "L", false)

		a1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[2], 5, list.InvNumber, "", "L", false)

		a2 := pdf.GetY()
		pdf.SetY(y3)
		// tk.Println(list.Acc_Code)
		x = 10 + widthHead[0] + widthHead[1] + widthHead[2]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		amount := tk.Sprintf("%.2f", list.Amount)
		amount = c.ConvertToCurrency(amount)
		pdf.MultiCell(widthHead[3], 5, amount, "", "R", false)

		a3 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		alreadypaid := tk.Sprintf("%.2f", list.AlreadyPaid)
		alreadypaid = c.ConvertToCurrency(alreadypaid)
		pdf.MultiCell(widthHead[4], 5, alreadypaid, "", "R", false)

		a4 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		// receive := tk.Sprintf("%.2f", list.Receive)
		// receive = c.ConvertToCurrency(receive)
		pdf.MultiCell(widthHead[5], 5, "0", "", "R", false)

		a5 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		balance := tk.Sprintf("%.2f", list.Balance)
		balance = c.ConvertToCurrency(balance)
		pdf.MultiCell(widthHead[6], 5, balance, "", "R", false)

		a6 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		if list.Pay {
			tk.Sprintf("%.2f", list.Pay)
			pdf.MultiCell(widthHead[7], 5, "Paid", "", "C", false)
		}

		a7 := pdf.GetY()
		pdf.SetY(y3)

		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}

	y4 := pdf.GetY()
	pdf.Line(pdf.GetX()-258, y3, pdf.GetX()-258, y4) //vertical 1
	pdf.Line(pdf.GetX()-248, y3, pdf.GetX()-248, y4) //vertical 2
	pdf.Line(pdf.GetX()-220, y3, pdf.GetX()-220, y4) //vertical 3
	pdf.Line(pdf.GetX()-180, y3, pdf.GetX()-180, y4) //vertical 4
	pdf.Line(pdf.GetX()-135, y3, pdf.GetX()-135, y4) //vertical 5
	pdf.Line(pdf.GetX()-90, y3, pdf.GetX()-90, y4)   //vertical 6
	pdf.Line(pdf.GetX()-45, y3, pdf.GetX()-45, y4)   //vertical 7
	pdf.Line(pdf.GetX(), y3, pdf.GetX(), y4)         //vertical 8
	pdf.Line(pdf.GetX()+18, y3, pdf.GetX()+18, y4)   //vertical 9

	pdf.Line(x0, y00, x0, y4)                                                                                                                                                                                                                     // garis vertikal 1 dari kiri/156 -> 136
	pdf.Line(10.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y00, 10.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y4) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()-258, y4, 286, y4)

	e = os.RemoveAll(c.PdfPath + "/salespayment")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/salespayment", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}

	namepdf := "-salespayment.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/salespayment"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	// tk.Println(e)
	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName

}

func (c *TransactionController) GetInvoiceBasedOnDocNum(InvNumber string) []InvoiceModel {
	typeInv := strings.Split(InvNumber, "/")[0]
	tableName := ""
	if typeInv == "INV" {
		tableName = "Invoice"
	} else {
		tableName = "InvoiceNonInv"
	}
	crs, _ := c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNo", InvNumber)).From(tableName).Cursor(nil)
	defer crs.Close()
	resultINV := []InvoiceModel{}
	_ = crs.Fetch(&resultINV, 0, false)
	return resultINV
}
