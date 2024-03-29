package controllers

import (
	// . "eaciit/proactive-inv/helpers"
	. "../models"

	"os"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

func (c *TransactionController) GetDataPurchaseCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		SupplierCode string
		TextSearch   string
		DateStart    string
		DateEnd      string
		Filter       bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Eq("Status", "PAID"))
	if p.Filter == true {
		filter = append(filter, dbox.Gte("DatePosting", dateStart))
		filter = append(filter, dbox.Lt("DatePosting", dateEnd))
		if p.TextSearch != "" {
			// filter = append(filter, dbox.Contains("Remark", p.TextSearch))
			filter = append(filter, dbox.Or(dbox.Contains("DocumentNumber", p.TextSearch), dbox.Contains("SupplierName", p.TextSearch), dbox.Contains("Remark", p.TextSearch)))
		}
		if p.SupplierCode != "" {
			filter = append(filter, dbox.Eq("SupplierCode", p.SupplierCode))
		}
	} else {
		filter = append(filter, dbox.Gte("DatePosting", dateStart))
		filter = append(filter, dbox.Lt("DatePosting", dateEnd))
	}

	csr, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseCreditMemo").Where(filter...).Cursor(nil)

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

func (c *TransactionController) InsertDraftPurchaseCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	t := struct {
		Data       PurchaseCreditMemo
		LastNumber int
	}{}

	err := k.GetPayload(&t)
	if err != nil {
		// tk.Println(err.Error())
		return c.ErrorResultInfo(err.Error(), nil)
	}

	p := t.Data
	// tk.Println("test", p)
	m := time.Now().UTC().Month()
	y := time.Now().UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequencePO").Select().Where(dbox.And(dbox.Eq("collname", "purchasecreditmemo"),
		dbox.Eq("month", int(m)), dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequencePOCMModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	secPO := 0
	if len(result) == 0 {
		mod := NewSequencePOCMModel()
		mod.Collname = "purchasecreditmemo"
		mod.TypePo = "purchasecreditmemo"
		mod.Lastnumber = 1
		mod.Month = int(m)
		mod.Year = y
		e = c.Ctx.Save(mod)
		secPO = 1
	} else {
		sec := result[0]
		if p.ID == "" {
			secPO = sec.Lastnumber
			sec.Lastnumber = sec.Lastnumber + 1
			c.Ctx.Save(&sec)
			secPO = sec.Lastnumber
		}
	}
	strNum := ""
	if secPO < 10 {
		strNum = "000"
	} else if secPO >= 10 && secPO < 100 {
		strNum = "00"
	} else if secPO >= 100 && secPO < 1000 {
		strNum = "0"
	}
	model := NewPurchaseCreditMemo()
	model.ID = p.ID
	model.DocumentNumberPO = p.DocumentNumber
	model.DatePosting = p.DatePosting
	// if p.ID == "" {
	model.ID = bson.NewObjectId()
	dateFormat := model.DatePosting.Format("02012006")
	model.DocumentNumber = "CMV/" + dateFormat + "/" + strNum + strconv.Itoa(secPO)
	// }
	model.Status = "PAID"
	model.AccountCode = p.AccountCode
	model.DateStr = model.DatePosting.Format("02-Jan-2006")
	// model.DocumentNumber = "PO/" + dateFormat + "/" + strNum + strconv.Itoa(secPO)
	model.SupplierCode = p.SupplierCode
	model.SupplierName = p.SupplierName
	model.SalesCode = p.SalesCode
	model.SalesName = p.SalesName
	model.Payment = p.Payment
	// model.Type = p.Type
	model.TotalIDR = p.TotalIDR
	model.TotalUSD = p.TotalUSD
	model.Discount = p.Discount
	model.VAT = p.VAT
	model.GrandTotalIDR = p.GrandTotalIDR
	model.GrandTotalUSD = p.GrandTotalUSD
	model.Rate = p.Rate
	model.User = k.Session("username").(string)
	model.Currency = p.Currency
	model.DownPayment = p.DownPayment
	model.Department = p.Department
	// model.StatusPayment = p.StatusPayment
	if p.ListDetail[0].Id == "" {
		for key, _ := range p.ListDetail {
			p.ListDetail[key].Id = bson.NewObjectId()
		}
	}

	model.ListDetail = p.ListDetail
	model.Remark = p.Remark
	e = c.Ctx.Save(model)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	// change PO
	{
		csr, e := c.Ctx.Connection.NewQuery().Select().From(NewPurchaseInventory().TableName()).Where(dbox.Eq("DocumentNumber", model.DocumentNumberPO)).Cursor(nil)
		defer csr.Close()
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		POdata := PurchaseInventory{}
		e = csr.Fetch(&POdata, 1, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		POdata.CreditMemo = true
		e = c.Ctx.Save(&POdata)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
	}
	//
	//get data PP
	// {
	csr, e := c.Ctx.Connection.NewQuery().Select().From(NewPurchasePaymentModel().TableName()).Where(dbox.Eq("ListDetail.PoNumber", model.DocumentNumberPO)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	PPdata := PurchasePaymentModel{}
	e = csr.Fetch(&PPdata, 1, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// }
	history := HistoryTrackPurchase{}
	history.Id = bson.NewObjectId()
	history.DocumentNumber = model.DocumentNumber
	history.DateCreated = p.DatePosting
	history.DateStr = model.DateStr
	history.DatePO = model.DatePosting
	history.Status = "PP"
	history.Remark = model.Remark
	history.SupplierCode = model.SupplierCode
	history.SupplierName = model.SupplierName
	crs, e = c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNumber", model.DocumentNumber)).From("TrackingPurchase").Cursor(nil)
	defer crs.Close()
	if crs.Count() == 0 {
		po := NewTrackPurchaseModel()
		po.ID = bson.NewObjectId()
		po.DocumentNumber = model.DocumentNumber
		po.DateCreated = p.DatePosting
		po.DateStr = model.DateStr
		po.DatePO = model.DatePosting
		po.Status = "POCM"
		po.Remark = model.Remark
		po.SupplierCode = model.SupplierCode
		po.SupplierName = model.SupplierName
		po.History = append(po.History, history)
		c.Ctx.Save(po)
	} else {
		resultPO := []TrackPurchaseModel{}
		e = crs.Fetch(&resultPO, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		mod := resultPO[0]
		po := TrackPurchaseModel{}
		po.ID = mod.ID
		po.DocumentNumber = model.DocumentNumber
		po.DateCreated = p.DatePosting
		po.DateStr = model.DateStr
		po.DatePO = model.DatePosting
		po.Status = "POCM"
		po.Remark = model.Remark
		po.SupplierCode = model.SupplierCode
		po.SupplierName = model.SupplierName
		po.History = append(po.History, history)
		c.Ctx.Save(&po)
	}
	for key, _ := range p.ListDetail {
		crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.And(dbox.Eq("INVID", p.ListDetail[key].CodeItem), dbox.Eq("StoreLocation", 1000))).From("Inventory").Cursor(nil)
		defer crs.Close()

		resultListDetailInventory := []InventoryModel{}
		e = crs.Fetch(&resultListDetailInventory, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		// tk.Println("result", resultListDetailInventory, p.ListDetail[key].CodeItem)
		dresult := resultListDetailInventory[0]
		pi := InventoryModel{}

		pi.ID = dresult.ID
		pi.INVID = dresult.INVID
		pi.INVDesc = dresult.INVDesc
		pi.Unit = dresult.Unit
		pi.Type = dresult.Type
		pi.Beginning = dresult.Beginning
		pi.InInventory = dresult.InInventory
		pi.OutInventory = dresult.OutInventory
		pi.CMVInventory = dresult.CMVInventory + p.ListDetail[key].Qty
		pi.CMInventory = dresult.CMInventory
		pi.TSInventory = dresult.TSInventory
		pi.TRInventory = dresult.TRInventory
		pi.Saldo = ((pi.InInventory + pi.CMInventory + pi.TRInventory) - (pi.OutInventory + pi.CMVInventory + pi.TSInventory))
		pi.UnitCost = p.ListDetail[key].PriceIDR
		pi.Total = float64(pi.Saldo) * pi.UnitCost
		pi.LastDate = p.DatePosting
		pi.StoreLocation = dresult.StoreLocation
		pi.StoreLocationName = dresult.StoreLocationName

		dataListDetail := []ListHistoryInventory{}
		listpi := ListHistoryInventory{}
		listpi.Id = bson.NewObjectId()
		listpi.INVID = dresult.INVID
		listpi.INVDesc = dresult.INVDesc
		listpi.Unit = dresult.Unit
		listpi.Type = dresult.Type
		listpi.Beginning = dresult.Beginning
		listpi.InInventory = dresult.InInventory
		listpi.OutInventory = dresult.OutInventory
		listpi.CMVInventory = dresult.CMVInventory + p.ListDetail[key].Qty
		listpi.CMInventory = dresult.CMInventory
		listpi.TSInventory = dresult.TSInventory
		listpi.TRInventory = dresult.TRInventory
		listpi.Saldo = ((listpi.InInventory + listpi.CMInventory + listpi.TRInventory) - (listpi.OutInventory + listpi.CMVInventory + listpi.TSInventory))
		listpi.UnitCost = p.ListDetail[key].PriceIDR
		listpi.Total = float64(pi.Saldo) * pi.UnitCost
		listpi.LastDate = p.DatePosting
		listpi.StoreLocation = dresult.StoreLocation
		listpi.StoreLocationName = dresult.StoreLocationName
		dataListDetail = append(dataListDetail, listpi)
		pi.ListInventory = dataListDetail

		e = c.Ctx.Save(&pi)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		//loginventory
		logitem := LogInventoryModel{}
		logitem.Id = bson.NewObjectId()
		logitem.CodeItem = dresult.INVID
		logitem.Item = dresult.INVDesc
		logitem.StorehouseId = dresult.StoreLocation
		logitem.StoreHouseName = dresult.StoreLocationName
		logitem.Date = p.DatePosting
		logitem.Description = model.DocumentNumber
		logitem.Price = p.ListDetail[key].PriceIDR
		logitem.TypeTransaction = "CMV"
		logitem.StockUnit = dresult.Saldo
		logitem.CountTransaction = p.ListDetail[key].Qty
		logitem.Increment = 0
		logitem.Decrement = p.ListDetail[key].Qty
		logitem.TotalSaldo = logitem.StockUnit - logitem.Decrement
		e = c.Ctx.Save(&logitem)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
	}
	//jornal
	c.SavetoJournalFromCMV(3110, "HUTANG DAGANG", 1401, model.Department, model.Discount, model.TotalIDR, model.VAT, model.GrandTotalIDR, k.Session("username").(string), model.DocumentNumber, model.DatePosting, p.SalesCode, p.SalesName)
	if p.ID == "" {
		c.LogActivity("Purchase Credit Memo", "Insert Purchase Credit Memo", p.DocumentNumber, k)
	} else {
		c.LogActivity("Purchase Credit Memo", "Update Purchase Credit Memo", p.DocumentNumber, k)
	}

	return c.SetResultOK(nil)
}
func (c *TransactionController) SavetoJournalFromCMV(AccountDebet int, AccountNameDebet string, AccountCredit int, Department string, disCount float64, Amount float64, VAT float64, GrandTotal float64, User string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) interface{} {
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
	mdl.Journal_Type = "General"
	mdl.Status = "posting"
	mdl.User = User
	mdl.Department = Department
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	dataListDetail := []Journal{}
	list := Journal{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = Amount
	if disCount != 0 {
		list.Debet = GrandTotal - VAT
	}
	list.User = User
	list.Department = Department
	list.SalesCode = SalesCode
	list.SalesName = SalesName
	idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Acc_Code = AccountDebet
	var accs []interface{}
	accs = append(accs, AccountDebet)
	accs = append(accs, AccountCredit)
	csr, e := c.Ctx.Connection.NewQuery().From("Coa").Where(dbox.In("acc_code", accs...)).Cursor(nil)
	defer csr.Close()
	result := []CoaModel{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}
	AccountCreditName := "PERSEDIAAN BARANG DAGANG"
	for _, each := range result {
		if each.ACC_Code == AccountDebet {
			AccountNameDebet = each.Account_Name
		}
		if each.ACC_Code == AccountCredit {
			AccountCreditName = each.Account_Name
		}
	}
	list.Description = Desc
	list.Attachment = "INVOICE"

	list.Acc_Name = AccountNameDebet
	dataListDetail = append(dataListDetail, list)
	if VAT != 0 {
		vat := Journal{}
		vat.Id = tk.RandomString(32)
		vat.No = 2
		vat.PostingDate = mdl.PostingDate
		vat.DateStr = mdl.DateStr
		vat.Journal_Type = mdl.Journal_Type
		vat.DocumentNumber = list.DocumentNumber
		vat.Debet = VAT
		vat.User = User
		vat.Acc_Code = 1321
		vat.Acc_Name = "PPN MASUKAN"
		vat.Attachment = "INVOICE"
		vat.Description = Desc
		vat.Department = Department
		vat.SalesCode = SalesCode
		vat.SalesName = SalesName
		dataListDetail = append(dataListDetail, vat)
	}
	list2 := Journal{}
	list2.Id = tk.RandomString(32)
	list2.No = 2
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = list.DocumentNumber
	if VAT != 0 {
		list2.No = 3
		list2.Credit = GrandTotal
	}
	list2.Credit = GrandTotal
	list2.User = User
	list2.Acc_Code = AccountCredit
	list2.Acc_Name = AccountCreditName
	list2.Attachment = "INVOICE"
	list2.Description = Desc
	list2.Department = Department
	list2.SalesCode = SalesCode
	list2.SalesName = SalesName
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e = c.Ctx.Save(mdl)
	c.SavetoGeneralLedgerFromCMV(AccountDebet, AccountNameDebet, AccountCredit, Department, disCount, Amount, VAT, GrandTotal, User, mdl.IdJournal, Desc, DatePosting, SalesCode, SalesName)
	return ""
}
func (c *TransactionController) SavetoGeneralLedgerFromCMV(AccountDebet int, AccountNameDebet string, AccountCredit int, Department string, disCount float64, Amount float64, VAT float64, GrandTotal float64, User string, idJournal string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) interface{} {
	m := DatePosting.UTC().Month()
	y := DatePosting.UTC().Year()
	mdl := new(MainGeneralLedger)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = idJournal
	mdl.PostingDate = DatePosting
	mdl.CreateDate = time.Now()
	mdl.DateStr = DatePosting.Format("02 Jan 2006")
	mdl.Journal_Type = "General"
	mdl.Department = Department
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
		mdl.Status = "posting"
	} else {
		mdl.Status = "close"
	}
	mdl.User = User
	dataListDetail := []GeneralDetail{}
	list := GeneralDetail{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = Amount
	list.Debet = Amount
	if disCount != 0 {
		list.Debet = GrandTotal - VAT
	}
	list.User = User
	list.Department = Department
	list.SalesCode = SalesCode
	list.SalesName = SalesName
	idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Acc_Code = AccountDebet
	list.Description = Desc
	csr, e := c.Ctx.Connection.NewQuery().From("Coa").Where(dbox.Eq("acc_code", AccountCredit)).Cursor(nil)
	defer csr.Close()
	result := []CoaModel{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}
	list.Attachment = "INVOICE"
	list.Acc_Name = AccountNameDebet
	dataListDetail = append(dataListDetail, list)
	if VAT != 0 {
		vat := GeneralDetail{}
		vat.Id = tk.RandomString(32)
		vat.No = 2
		vat.PostingDate = mdl.PostingDate
		vat.DateStr = mdl.DateStr
		vat.Journal_Type = mdl.Journal_Type
		vat.DocumentNumber = list.DocumentNumber
		vat.Debet = VAT
		vat.User = User
		vat.Acc_Code = 1361
		vat.Acc_Name = "VALUE ADDED TAX (VAT) IN"
		vat.Attachment = "INVOICE"
		vat.Description = Desc
		vat.Department = Department
		vat.SalesCode = SalesCode
		vat.SalesName = SalesName
		dataListDetail = append(dataListDetail, vat)
	}
	list2 := GeneralDetail{}
	list2.Id = tk.RandomString(32)
	list2.No = 2
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = list.DocumentNumber
	list2.Credit = GrandTotal
	if VAT != 0 {
		list2.No = 3
		list2.Credit = GrandTotal
	}
	list2.Attachment = "INVOICE"
	list2.User = User
	list2.Acc_Code = AccountCredit
	list2.Acc_Name = result[0].Account_Name
	list2.Description = Desc
	list2.Department = Department
	list2.SalesCode = SalesCode
	list2.SalesName = SalesName
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e = c.Ctx.Save(mdl)
	return ""
}
func (c *TransactionController) DeleteDraftPurchaseCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id             bson.ObjectId
		DocumentNumber string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(PurchaseCreditMemo)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		c.WriteLog(e)
	}
	e = c.Ctx.Delete(result)
	crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNumber", p.DocumentNumber)).From("TrackingPurchase").Cursor(nil)
	defer crs.Close()
	resultPO := []TrackPurchaseModel{}
	e = crs.Fetch(&resultPO, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	mod := resultPO[0]
	e = c.Ctx.Delete(&mod)
	c.LogActivity("Purchase Credit Memo", "Delete purchase order draft", p.DocumentNumber, k)
	return c.SetResultInfo(false, "OK", nil)
}
func (c *TransactionController) GetDataPurchaseInvoiceInventoryForCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		TextSearch string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	crs, e := c.Ctx.Connection.NewQuery().From("PurchaseInventory").Where(dbox.And(dbox.Or(dbox.Eq("Status", "PAID"), dbox.Eq("Status", "PI")), dbox.Eq("CreditMemo", false), dbox.Eq("DocumentNumber", p.TextSearch))).Select().Cursor(nil)
	defer crs.Close()
	results := make([]PurchaseInventory, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *TransactionController) ExportToPdfPurchaseCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	csr, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseCreditMemo").Where(dbox.Eq("_id", p.Id)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []PurchaseCreditMemo{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	// tk.Println(DATA)

	if DATA.Currency == "USD" {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceIDR = 0
			DATA.ListDetail[i].AmountIDR = 0
		}
		// discount = DATA.Discount
	} else {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceUSD = 0
			DATA.ListDetail[i].AmountUSD = 0
		}
	}

	csr, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.SupplierCode)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	resultsSupp := []CustomerModel{}
	e = csr.Fetch(&resultsSupp, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	supp := resultsSupp[0]
	//user
	csr, e = c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(dbox.Eq("username", DATA.User)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	User := SysUserModel{}
	e = csr.Fetch(&User, 1, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	pdf := gofpdf.New("L", "mm", "A5", c.FontPath)
	pdf.SetDrawColor(2, 2, 2)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	pdf.AddPage()
	x_defaulft := 10.0
	pdf.SetFont("Century_Gothicb", "B", 15)
	pdf.SetXY(70, 10)
	pdf.CellFormat(0, 12, "PURCHASE CREDIT MEMO", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(85)

	pdf.SetFont("Century_Gothicb", "B", 10)
	pdf.CellFormat(0, 12, DATA.DocumentNumber, "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothic", "", 8)
	y0 := pdf.GetY() + 5
	//
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Kode, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Date", "", "L", false) // date
	date := DATA.DatePosting.Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Name, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
	dueDate := DATA.DatePosting.AddDate(0, 0, 30).Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(20, 3, "Phone :", "", "L", false) // phone
	pdf.SetXY(30, y0)
	pdf.MultiCell(60, 3, supp.NoTelp, "", "L", false) //phone
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+DATA.DocumentNumber, "", "L", false) // DocumentNo
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Address, "", "L", false) // address
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Sales", "", "L", false) // sales
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+DATA.SalesName, "", "L", false) // sales
	pdf.Ln(5)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(20, 3, "Rek Bank :", "", "L", false) // rek bank
	pdf.SetXY(30, y0)
	if supp.Bank != "" && supp.AccountNo != "" {
		pdf.MultiCell(60, 3, supp.Bank+"-"+supp.AccountNo, "", "L", false) //rek bank
	} else {
		pdf.MultiCell(60, 3, "", "", "L", false) //rek bank
	}
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	invHead := []string{"No", "", "Item", "Qty", "Price", "Disc. Amount", "Amount"}
	widHead := []float64{10.0, 0.0, 80.0, 10.0, 30.0, 30.0, 30.0}
	for i, head := range invHead {
		pdf.SetY(y0)
		x := x_defaulft
		for j, w := range widHead {
			if i > j {
				x += w
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)

		pdf.MultiCell(widHead[i], 4, head, "TB", "C", false)
	}
	// grid
	y0 = pdf.GetY()
	lastbigest := y0
	for i, list := range DATA.ListDetail {
		yg := pdf.GetY()
		x := x_defaulft
		pdf.SetX(x)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(widHead[0], 4, numberstr, "", "C", false)
		pdf.SetY(yg)
		x += widHead[0]
		pdf.SetX(x)
		a0 := pdf.GetY()
		pdf.MultiCell(widHead[1], 4, "", "", "L", false)
		pdf.SetY(yg)
		x += widHead[1]
		pdf.SetX(x)
		a1 := pdf.GetY()
		pdf.MultiCell(widHead[2], 4, list.Item, "", "L", false)
		pdf.SetY(yg)
		x += widHead[2]
		pdf.SetX(x)
		a2 := pdf.GetY()
		pdf.MultiCell(widHead[3], 4, strconv.Itoa(list.Qty), "", "L", false)
		pdf.SetY(yg)
		x += widHead[3]
		pdf.SetX(x)
		a3 := pdf.GetY()
		priceidr := tk.Sprintf("%.2f", list.PriceIDR)
		pdf.MultiCell(widHead[4], 4, c.ConvertToCurrency(priceidr), "", "R", false)
		pdf.SetY(yg)
		x += widHead[4]
		pdf.SetX(x)
		a4 := pdf.GetY()
		pdf.MultiCell(widHead[5], 4, "-", "", "C", false)
		pdf.SetY(yg)
		x += widHead[5]
		pdf.SetX(x)
		a5 := pdf.GetY()
		amount := tk.Sprintf("%.2f", list.AmountIDR)
		pdf.MultiCell(widHead[6], 4, c.ConvertToCurrency(amount), "", "R", false)
		a6 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3, a4, a5, a6}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		lastbigest = biggest
	}
	y0 = lastbigest
	// if y0 < 80 {
	pdf.Line(x_defaulft, 80, 200, 80)
	y0 = 80.0
	// } else {
	// 	pdf.Line(x_defaulft, y0, 200, y0)
	// }
	pdf.SetY(y0)
	pdf.Ln(2)
	y0 = pdf.GetY()
	pdf.SetY(y0)
	pdf.MultiCell(20, 3, "Remarks :", "", "L", false) // Remark
	pdf.SetXY(30, y0)
	pdf.MultiCell(40, 3, DATA.Remark, "", "L", false) // Remark
	pdf.Ln(5)
	//
	y0 = pdf.GetY()
	pdf.SetY(y0)
	yTotal := pdf.GetY()
	headBottom := []string{"Prepared by :", "Approved by", "Finance", "Received by :"}
	widthBottom := []float64{40.0, 30.0, 30.0, 30.0}
	for i, head := range headBottom {
		pdf.SetY(y0)
		x := x_defaulft
		for j, w := range widthBottom {
			if i > j {
				x += w
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthBottom[i], 4, head, "", "L", false)
		} else {
			pdf.MultiCell(widthBottom[i], 4, head, "", "C", false)
		}
	}
	pdf.Ln(15)
	yB := pdf.GetY()
	xx := x_defaulft
	pdf.SetX(xx)
	pdf.SetY(yB)
	pdf.MultiCell(widthBottom[0], 4, User.Fullname, "", "L", false)
	b0 := pdf.GetY()
	xx += widthBottom[0]
	pdf.SetXY(xx, yB)
	pdf.MultiCell(widthBottom[1], 4, "(                          )", "", "C", false)
	b1 := pdf.GetY()
	xx += widthBottom[1]
	pdf.SetXY(xx, yB)
	pdf.MultiCell(widthBottom[2], 4, "(                          )", "", "C", false)
	b2 := pdf.GetY()
	xx += widthBottom[2]
	pdf.SetXY(xx, yB)
	pdf.MultiCell(widthBottom[3], 4, "(                          )", "", "C", false)
	b3 := pdf.GetY()
	allB := []float64{b0, b1, b2, b3}
	var n, biggestB float64
	for _, v := range allB {
		if v > n {
			n = v
			biggestB = n
		}
	}
	lastY := biggestB
	// total etc
	xtotal := 150.0
	yTotal += 0.7
	pdf.SetXY(xtotal, yTotal)
	pdf.MultiCell(20, 3, "Total ", "", "L", false) // Total
	pdf.SetXY(170, yTotal)
	total := tk.Sprintf("%.2f", DATA.TotalIDR)
	pdf.MultiCell(30, 3, c.ConvertToCurrency(total), "", "R", false) // Total
	pdf.Ln(1)
	//
	yTotal = pdf.GetY()
	pdf.SetXY(xtotal, yTotal)
	pdf.MultiCell(20, 3, "Discount ", "", "L", false) // discount
	pdf.SetXY(170, yTotal)
	pdf.MultiCell(30, 3, "0.00", "", "R", false) // discount
	pdf.Ln(1)
	//
	yTotal = pdf.GetY()
	pdf.SetXY(xtotal, yTotal)
	pdf.MultiCell(20, 3, "VAT 10% ", "", "L", false) // vat
	pdf.SetXY(170, yTotal)
	vat := tk.Sprintf("%.2f", DATA.VAT)
	pdf.MultiCell(30, 3, c.ConvertToCurrency(vat), "", "R", false) // vat
	pdf.Ln(1)
	//grantototal
	yB += 0.7
	pdf.SetFont("Century_Gothicb", "B", 8)
	pdf.SetXY(xtotal, yB)
	pdf.MultiCell(20, 4, "Grand Total ", "TB", "L", false) // grandtotal
	pdf.SetXY(170, yB)
	grandTotal := tk.Sprintf("%.2f", DATA.GrandTotalIDR)
	pdf.MultiCell(30, 4, c.ConvertToCurrency(grandTotal), "TB", "R", false) // grandtotal
	pdf.Ln(1)
	// end bottom
	y0 = lastY
	pdf.SetY(y0)
	pdf.Ln(2)
	// y0 = pdf.GetY()
	pdf.SetFont("Century_Gothic", "", 8)
	// pdf.SetXY(30, y0)
	// pdf.MultiCell(150, 3, "Barang yang sudah dibeli tidak bisa di kembalikan / di tukar dengan uang", "", "L", false) // alert
	// pdf.Ln(2)
	y0 = pdf.GetY()
	pdf.SetXY(10, y0)
	pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
	pdf.SetXY(30, y0)
	datenow := time.Now().Format("January 02, 2006")
	pdf.MultiCell(150, 3, datenow, "", "L", false) // date print
	e = os.RemoveAll(c.PdfPath + "/purchasecreditmemo")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/purchasecreditmemo", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-purchasecreditmemo.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/purchasecreditmemo"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName
}

func (c *TransactionController) GetAutoCInvNumPCM(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := KendoDatasourceQueryFilter{}
	k.GetPayload(&p)

	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Or(dbox.Eq("Status", "PAID"), dbox.Eq("Status", "PI")))
	filter = append(filter, dbox.Eq("CreditMemo", false))
	filter = append(filter, p.ToDboxFilter())

	crs, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseInventory").Where(filter...).Cursor(nil)
	defer crs.Close()
	results := make([]struct {
		DocumentNumber string `bson:"DocumentNumber"`
	}, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		c.SetResultInfo(true, "Error", e.Error())
	}

	return KendoDatasourceResult{
		Data:  results,
		Count: len(results),
	}
}
