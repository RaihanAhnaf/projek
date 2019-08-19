package controllers

import (
	. "../helpers"
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

func (c *TransactionController) GetDataSalesCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// var locid = k.Session("locationid").(int)
	p := struct {
		CustomerCode string
		TextSearch   string
		DateStart    string
		DateEnd      string
		Filter       bool
		Status       string
		LocationID   int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Gte("DateCreated", dateStart))
	filter = append(filter, dbox.Lt("DateCreated", dateEnd))
	filter = CreateLocationFilter(filter, "StoreLocationId", p.LocationID, false)

	if p.Filter == true {
		if p.TextSearch != "" {
			// filter = append(filter, dbox.Contains("DocumentNo", p.TextSearch))
			filter = append(filter, dbox.Or(dbox.Contains("StoreLocationName", p.TextSearch), dbox.Contains("DocumentNo", p.TextSearch), dbox.Contains("CustomerName", p.TextSearch), dbox.Contains("Status", p.TextSearch)))
		}
		if p.CustomerCode != "" {
			filter = append(filter, dbox.Eq("CustomerCode", p.CustomerCode))
		}
		if p.Status != "" {
			filter = append(filter, dbox.Eq("Status", p.Status))
		}
	}
	crs, e := c.Ctx.Connection.NewQuery().Select().From("SalesCreditMemo").Where(filter...).Cursor(nil)

	defer crs.Close()
	results := make([]SalesCreditMemo, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	if p.Filter && len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransactionController) GetLastNumberSCM(Date time.Time) int {
	m := Date.UTC().Month()
	y := Date.UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceSCM").Select().Where(dbox.And(dbox.Eq("collname", "salescreditmemo"),
		dbox.Eq("month", int(m)), dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequenceSCMModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		// return c.SetResultInfo(true, e.Error(), nil)
		tk.Println(e.Error())
	}
	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceSCMModel()
		model.Collname = "salescreditmemo"
		model.TypePo = "salescreditmemo"
		model.Lastnumber = 10
		model.Month = int(m)
		model.Year = y
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	}
	sec := result[0]
	sec.Lastnumber = sec.Lastnumber + 1
	data.Number = sec.Lastnumber
	data.Msg = "Success"

	return data.Number
}

func (c *TransactionController) DeleteSalesCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id         bson.ObjectId
		DocumentNo string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(SalesCreditMemo)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		c.WriteLog(e)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("SalesCreditMemo", "Delete SalesCreditMemo Draft", p.DocumentNo, k)
	return c.SetResultOK(nil)
}

func (c *TransactionController) InsertSalesCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	t := struct {
		Data SalesCreditMemo
	}{}

	err := k.GetPayload(&t)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	p := t.Data
	model := NewSalesCreditMemo()
	model.DateStr = p.DateStr
	newDate, _ := time.Parse("02-Jan-2006", model.DateStr)
	model.DateCreated = newDate
	LastNumber := c.GetLastNumberSCM(model.DateCreated)
	// tk.Println(LastNumber)
	// if p.Id == "" {
	model.Id = bson.NewObjectId()
	c.SaveLastNumberSalesCreditMemo(LastNumber, model.DateCreated)
	p.DocumentNo = c.SetDocumentNumberSalesCreditMemo(LastNumber, model.DateCreated)
	// } else {
	// 	model.Id = p.Id
	// }
	model.DocumentNo = p.DocumentNo
	model.DocumentNoInvoice = p.DocumentNoInvoice
	model.CustomerCode = p.CustomerCode
	model.CustomerName = p.CustomerName
	model.ListItem = p.ListItem
	// model.PoNumber = p.PoNumber
	model.Status = p.Status
	model.Total = p.Total
	// model.AccountCode = p.AccountCode
	// model.AccountName = p.AccountName
	model.Currency = p.Currency
	model.VAT = p.VAT
	model.User = k.Session("username").(string)
	model.GrandTotalIDR = p.GrandTotalIDR
	model.GrandTotalUSD = p.GrandTotalUSD
	model.Rate = p.Rate
	model.StoreLocationId = p.StoreLocationId
	model.StoreLocationName = p.StoreLocationName
	if p.ListItem[0].ID == "" {
		for key, _ := range p.ListItem {

			idList := strconv.Itoa(key) + model.DocumentNo
			p.ListItem[key].ID = idList

		}
	}

	model.ListItem = p.ListItem
	model.Description = p.Description
	model.SalesCode = p.SalesCode
	model.SalesName = p.SalesName
	c.Ctx.Save(model)
	// change INV
	{
		csr, e := c.Ctx.Connection.NewQuery().Select().From(NewInvoiceModel().TableName()).Where(dbox.Eq("DocumentNo", model.DocumentNoInvoice)).Cursor(nil)
		defer csr.Close()
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		INVdata := InvoiceModel{}
		e = csr.Fetch(&INVdata, 1, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		INVdata.CreditMemo = true
		e = c.Ctx.Save(&INVdata)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
	}
	if p.Status == "POSTING" {
		history := HistoryTrackSalesCreditMemo{}
		history.Id = bson.NewObjectId()
		history.DocumentNumber = model.DocumentNo
		history.DateCreated = p.DateCreated
		history.DateStr = p.DateCreated.Format("2006-01-02")
		history.DateSCM = model.DateCreated
		history.Status = "SalesCreditMemo"
		history.Remark = model.Description
		history.CustomerCode = model.CustomerCode
		history.CustomerName = model.CustomerName

		crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNumber", model.DocumentNo)).From("TrackingSalesCreditMemo").Cursor(nil)
		defer crs.Close()

		if crs.Count() == 0 {
			SCM := NewTrackSalesCreditMemo()
			SCM.ID = bson.NewObjectId()
			SCM.DocumentNumber = model.DocumentNo
			SCM.DateCreated = p.DateCreated
			SCM.DateStr = p.DateCreated.Format("2006-01-02")
			SCM.DateSCM = model.DateCreated
			SCM.Status = "SalesCreditMemo"
			SCM.Remark = model.Description
			SCM.CustomerCode = model.CustomerCode
			SCM.CustomerName = model.CustomerName
			SCM.History = append(SCM.History, history)
			c.Ctx.Save(SCM)
		} else {
			resultSCM := []TrackSalesCreditMemo{}
			e = crs.Fetch(&resultSCM, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			mod := resultSCM[0]
			SCM := TrackSalesCreditMemo{}
			SCM.ID = mod.ID
			SCM.DocumentNumber = model.DocumentNo
			SCM.DateCreated = p.DateCreated
			SCM.DateStr = p.DateCreated.Format("2006-01-02")
			SCM.DateSCM = model.DateCreated
			SCM.Status = "SalesCreditMemo"
			SCM.Remark = model.Description
			SCM.CustomerCode = model.CustomerCode
			SCM.CustomerName = model.CustomerName
			SCM.History = append(SCM.History, history)
			c.Ctx.Save(&SCM)
		}
		accountJournal := tk.M{}.Set("debet", 5130).Set("credit", 1210)
		c.SavetoJournalFromSalesCreditMemo(accountJournal, p.Total, p.VAT, p.GrandTotalIDR, k.Session("username").(string), p.Description, model.DateCreated, p.SalesCode, p.SalesName)
		c.LogActivity("SalesCreditMemo", "Posting SalesCreditMemo", p.DocumentNo, k)
		accountJournals := []int{1401, 5210}
		dataSecondJournal := []tk.M{}
		for key, _ := range p.ListItem {
			crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.And(dbox.Eq("INVID", p.ListItem[key].CodeItem), dbox.Eq("StoreLocation", model.StoreLocationId))).From("Inventory").Cursor(nil)
			defer crs.Close()

			resultListDetailInventory := []InventoryModel{}
			e = crs.Fetch(&resultListDetailInventory, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

			dresult := resultListDetailInventory[0]
			amount := dresult.UnitCost * float64(p.ListItem[key].Qty)
			vat := 0.0
			if model.VAT > 0 {
				vat = (10 * amount) / 100
			}
			grandtotal := amount + vat
			dataSecondJournal = append(dataSecondJournal, tk.M{}.Set("amount", amount).Set("vat", vat).Set("grandtotal", grandtotal))
			pi := InventoryModel{}

			pi.ID = dresult.ID
			pi.INVID = dresult.INVID
			pi.INVDesc = dresult.INVDesc
			pi.Unit = dresult.Unit
			pi.Type = dresult.Type
			pi.Beginning = dresult.Beginning
			pi.InInventory = dresult.InInventory
			pi.OutInventory = dresult.OutInventory
			pi.CMVInventory = dresult.CMVInventory
			pi.CMInventory = dresult.CMInventory + p.ListItem[key].Qty
			pi.TSInventory = dresult.TSInventory
			pi.TRInventory = dresult.TRInventory
			pi.Saldo = ((pi.InInventory + pi.CMInventory + pi.TRInventory) - (pi.OutInventory + pi.CMVInventory + pi.TSInventory))
			pi.UnitCost = dresult.UnitCost
			pi.Total = float64(pi.Saldo) * pi.UnitCost
			pi.LastDate = p.DateCreated
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
			listpi.CMVInventory = dresult.CMVInventory
			listpi.CMInventory = dresult.CMInventory + p.ListItem[key].Qty
			listpi.TSInventory = dresult.TSInventory
			listpi.TRInventory = dresult.TRInventory
			listpi.Saldo = ((listpi.InInventory + listpi.CMInventory + listpi.TRInventory) - (listpi.OutInventory + listpi.CMVInventory + listpi.TSInventory))
			listpi.UnitCost = dresult.UnitCost
			listpi.Total = float64(pi.Saldo) * pi.UnitCost
			listpi.LastDate = p.DateCreated
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
			logitem.Date = p.DateCreated
			logitem.Description = model.DocumentNo
			logitem.TypeTransaction = "CMI"
			logitem.Price = p.ListItem[key].PriceIDR
			logitem.StockUnit = dresult.Saldo
			logitem.CountTransaction = p.ListItem[key].Qty
			logitem.Increment = p.ListItem[key].Qty
			logitem.Decrement = 0
			logitem.TotalSaldo = logitem.StockUnit + logitem.Decrement
			e = c.Ctx.Save(&logitem)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
		}
		// c.SavetoJournalFromSalesCreditMemo(p.AccountCode, p.AccountName, p.Total, p.VAT, p.GrandTotalIDR, k.Session("username").(string), p.Description, model.DateCreated)
		// c.LogActivity("SalesCreditMemo", "Posting SalesCreditMemo", p.DocumentNo, k)
		e = c.SaveMultiJournalFromSalesCreditMemo(accountJournals, dataSecondJournal, k.Session("username").(string), p.Description, model.DateCreated, p.SalesCode, p.SalesName)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
	} else {
		if p.Id == "" {
			c.LogActivity("SalesCreditMemo", "Insert SalesCreditMemo", p.DocumentNo, k)
		} else {
			c.LogActivity("SalesCreditMemo", "Update SalesCreditMemo", p.DocumentNo, k)
		}
	}

	return c.SetResultOK(nil)
}

func (c *TransactionController) SaveLastNumberSalesCreditMemo(LastNumber int, Date time.Time) interface{} {
	m := Date.UTC().Month()
	y := Date.UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceSCM").Select().Where(dbox.And(dbox.Eq("collname", "salescreditmemo"),
		dbox.Eq("month", int(m)), dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequenceSCMModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	sec := result[0]
	sec.Lastnumber = LastNumber
	c.Ctx.Save(&sec)

	return c.SetResultInfo(false, "Success", sec)
}

func (c *TransactionController) SetDocumentNumberSalesCreditMemo(LastNumber int, Date time.Time) string {
	strNum := ""
	if LastNumber < 10 {
		strNum = "000"
	} else if LastNumber >= 10 && LastNumber < 100 {
		strNum = "00"
	} else if LastNumber >= 100 && LastNumber < 1000 {
		strNum = "0"
	}

	dateFormat := Date.Format("02012006")
	return "SCM/" + dateFormat + "/" + strNum + strconv.Itoa(LastNumber)
}

func (c *TransactionController) SavetoJournalFromSalesCreditMemo(accounts tk.M, Amount float64, VAT float64, GrandTotal float64, User string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) interface{} {
	accountCodes := []int{accounts.GetInt("debet"), accounts.GetInt("credit")}
	var filteraccount []interface{}
	for _, each := range accountCodes {
		filteraccount = append(filteraccount, each)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(dbox.In("acc_code", filteraccount...)).Cursor(nil)
	defer csr.Close()
	result := []CoaModel{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}
	accs := map[int]string{}
	for _, acc := range result {
		accs[acc.ACC_Code] = acc.Account_Name
	}
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
	mdl.Department = "COMMERCE"
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	mdl.User = User
	dataListDetail := []Journal{}
	list := Journal{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = GrandTotal
	list.User = User
	idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Description = Desc
	list.Acc_Code = accounts.GetInt("debet")
	list.Acc_Name = accs[list.Acc_Code]
	list.Department = "COMMERCE"
	list.SalesCode = SalesCode
	list.SalesName = SalesName
	list.Attachment = "SalesCreditMemo"
	dataListDetail = append(dataListDetail, list)

	if VAT != 0 {
		vat := Journal{}
		vat.Id = tk.RandomString(32)
		vat.No = 2
		vat.PostingDate = mdl.PostingDate
		vat.DateStr = mdl.DateStr
		vat.Journal_Type = mdl.Journal_Type
		vat.DocumentNumber = list.DocumentNumber
		vat.Credit = VAT
		vat.User = User
		vat.Acc_Code = 3122
		vat.Acc_Name = "PPN KELUARAN"
		vat.Attachment = "SalesCreditMemo"
		vat.Description = Desc
		vat.Department = "COMMERCE"
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
	list2.Credit = Amount
	if VAT != 0 {
		list2.No = 3
		list2.Credit = Amount
	}
	list2.User = User
	list2.Acc_Code = accounts.GetInt("credit")
	list2.Acc_Name = accs[list2.Acc_Code]
	list2.Description = Desc
	list2.Department = "COMMERCE"
	list2.SalesCode = SalesCode
	list2.SalesName = SalesName
	list2.Attachment = "SalesCreditMemo"
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e = c.Ctx.Save(mdl)
	c.SavetoGeneralLedgerFromSalesCreditMemo(accounts, accs, Amount, VAT, GrandTotal, User, mdl.IdJournal, Desc, DatePosting, SalesCode, SalesName)
	if e != nil {
		tk.Println(e.Error())
	}
	return ""
}
func (c *TransactionController) SavetoGeneralLedgerFromSalesCreditMemo(accounts tk.M, accs map[int]string, Amount float64, VAT float64, GrandTotal float64, User string, idJournal string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) interface{} {
	m := DatePosting.UTC().Month()
	y := DatePosting.UTC().Year()
	mdl := new(MainGeneralLedger)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = idJournal
	mdl.PostingDate = DatePosting
	mdl.CreateDate = time.Now()
	mdl.DateStr = DatePosting.Format("02 Jan 2006")
	mdl.Journal_Type = "General"
	// if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
	mdl.Status = "posting"
	// } else {
	// 	mdl.Status = "close"
	// }
	mdl.User = User
	mdl.Department = "COMMERCE"
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	dataListDetail := []GeneralDetail{}
	list := GeneralDetail{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = GrandTotal
	list.User = User
	idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Acc_Code = accounts.GetInt("debet")
	list.Acc_Name = accs[list.Acc_Code]
	list.Description = Desc
	list.Attachment = "SalesCreditMemo"
	list.Department = "COMMERCE"
	list.SalesCode = SalesCode
	list.SalesName = SalesName
	dataListDetail = append(dataListDetail, list)

	if VAT != 0 {
		vat := GeneralDetail{}
		vat.Id = tk.RandomString(32)
		vat.No = 2
		vat.PostingDate = mdl.PostingDate
		vat.DateStr = mdl.DateStr
		vat.Journal_Type = mdl.Journal_Type
		vat.DocumentNumber = list.DocumentNumber
		vat.Credit = VAT
		vat.User = User
		vat.Acc_Code = 3122
		vat.Acc_Name = "PPN KELUARAN"
		vat.Description = Desc
		vat.Attachment = "SalesCreditMemo"
		vat.Department = "COMMERCE"
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
		list2.Credit = Amount
	}
	list2.User = User
	list2.Acc_Code = accounts.GetInt("credit")
	list2.Acc_Name = accs[list2.Acc_Code]
	list2.Description = Desc
	list2.Department = "COMMERCE"
	list2.SalesCode = SalesCode
	list2.SalesName = SalesName
	list2.Attachment = "SalesCreditMemo"
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	return ""
}

func (c *TransactionController) SaveMultiJournalFromSalesCreditMemo(accounts []int, DataSecondJournal []tk.M, User string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) error {
	//journal reture
	// tk.Println("journal reture")
	var filteraccount []interface{}
	for _, each := range accounts {
		filteraccount = append(filteraccount, each)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(dbox.In("acc_code", filteraccount...)).Cursor(nil)
	defer csr.Close()
	result := []CoaModel{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		return e
	}
	accs := map[int]string{}
	for _, acc := range result {
		accs[acc.ACC_Code] = acc.Account_Name
	}
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
	mdl.Department = "COMMERCE"
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	mdl.User = User
	dataListDetail := []Journal{}
	for _, each := range DataSecondJournal {
		list := Journal{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = each.GetFloat64("grandtotal")
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + DatePosting.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Description = Desc
		list.Acc_Code = 1401
		list.Acc_Name = accs[list.Acc_Code]
		list.Department = "COMMERCE"
		list.SalesCode = SalesCode
		list.SalesName = SalesName
		list.Attachment = "SalesCreditMemo"
		dataListDetail = append(dataListDetail, list)
		if each.GetFloat64("vat") != 0 {
			vat := Journal{}
			vat.Id = tk.RandomString(32)
			vat.No = 2
			vat.PostingDate = mdl.PostingDate
			vat.DateStr = mdl.DateStr
			vat.Journal_Type = mdl.Journal_Type
			vat.DocumentNumber = list.DocumentNumber
			vat.Credit = each.GetFloat64("vat")
			vat.User = User
			vat.Acc_Code = 3122
			vat.Acc_Name = "PPN KELUARAN"
			vat.Attachment = "SalesCreditMemo"
			vat.Description = Desc
			vat.Department = "COMMERCE"
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
		list2.Credit = each.GetFloat64("grandtotal")
		if each.GetFloat64("vat") != 0 {
			list2.No = 3
			list2.Credit = each.GetFloat64("amount")
		}
		list2.User = User
		list2.Acc_Code = 5210
		list2.Acc_Name = accs[list2.Acc_Code]
		list2.Description = Desc
		list2.Department = "COMMERCE"
		list2.SalesCode = SalesCode
		list2.SalesName = SalesName
		list2.Attachment = "SalesCreditMemo"
		dataListDetail = append(dataListDetail, list2)
	}

	mdl.ListDetail = dataListDetail
	e = c.Ctx.Save(mdl)
	if e != nil {
		return e
	}
	// c.SavetoGeneralLedgerFromInvoice(accounts, accs, Amount, VAT, GrandTotal, User, mdl.IdJournal, Desc, DatePosting)
	e = c.SavetoMultiGeneralLedgerFromSalesCreditMemo(accs, DataSecondJournal, User, mdl.IdJournal, Desc, DatePosting, SalesCode, SalesName)
	if e != nil {
		return e
	}
	return nil
}
func (c *TransactionController) SavetoMultiGeneralLedgerFromSalesCreditMemo(accs map[int]string, DataSecondJournal []tk.M, User string, idJournal string, Desc string, DatePosting time.Time, SalesCode string, SalesName string) error {
	m := DatePosting.UTC().Month()
	y := DatePosting.UTC().Year()
	mdl := new(MainGeneralLedger)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = idJournal
	mdl.PostingDate = DatePosting
	mdl.CreateDate = time.Now()
	mdl.DateStr = DatePosting.Format("02 Jan 2006")
	mdl.Journal_Type = "General"
	// if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
	mdl.Status = "posting"
	// } else {
	// 	mdl.Status = "close"
	// }
	mdl.User = User
	mdl.Department = "COMMERCE"
	mdl.SalesCode = SalesCode
	mdl.SalesName = SalesName
	dataListDetail := []GeneralDetail{}
	for _, each := range DataSecondJournal {
		list := GeneralDetail{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = each.GetFloat64("grandtotal")
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + DatePosting.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Acc_Code = 1401
		list.Acc_Name = accs[list.Acc_Code]
		list.Description = Desc
		list.Attachment = "SalesCreditMemo"
		list.Department = "COMMERCE"
		list.SalesCode = SalesCode
		list.SalesName = SalesName
		dataListDetail = append(dataListDetail, list)

		if each.GetFloat64("vat") != 0 {
			vat := GeneralDetail{}
			vat.Id = tk.RandomString(32)
			vat.No = 2
			vat.PostingDate = mdl.PostingDate
			vat.DateStr = mdl.DateStr
			vat.Journal_Type = mdl.Journal_Type
			vat.DocumentNumber = list.DocumentNumber
			vat.Credit = each.GetFloat64("vat")
			vat.User = User
			vat.Acc_Code = 3122
			vat.Acc_Name = "PPN KELUARAN"
			vat.Description = Desc
			vat.Attachment = "SalesCreditMemo"
			vat.Department = "COMMERCE"
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
		list2.Credit = each.GetFloat64("grandtotal")
		if each.GetFloat64("vat") != 0 {
			list2.No = 3
			list2.Credit = each.GetFloat64("amount")
		}
		list2.User = User
		list2.Acc_Code = 5210
		list2.Acc_Name = accs[list2.Acc_Code]
		list2.Description = Desc
		list2.Department = "COMMERCE"
		list2.SalesCode = SalesCode
		list2.SalesName = SalesName
		list2.Attachment = "SalesCreditMemo"
		dataListDetail = append(dataListDetail, list2)
	}

	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		return e
	}
	return nil
}

// func (c *TransactionController) SavetoJournalFromSalesCreditMemo(AccountKredit int, AccountName string, Amount float64, VAT float64, GrandTotal float64, User string, Desc string, DatePosting time.Time) interface{} {
// 	m := DatePosting.UTC().Month()
// 	y := DatePosting.UTC().Year()
// 	codejurnal := tk.Sprintf("%02d%d", m, y)
// 	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
// 	number := tk.Sprintf("%04d", ids)
// 	mdl := new(MainJournal)
// 	mdl.ID = tk.RandomString(32)
// 	mdl.IdJournal = "JUR/" + codejurnal + "/" + number
// 	mdl.PostingDate = DatePosting
// 	mdl.CreateDate = time.Now()
// 	mdl.DateStr = DatePosting.Format("02 Jan 2006")
// 	mdl.Journal_Type = "General"
// 	mdl.Status = "posting"
// 	mdl.Department = "COMMERCE"
// 	mdl.User = User
// 	dataListDetail := []Journal{}
// 	list := Journal{}
// 	list.Id = tk.RandomString(32)
// 	list.No = 1
// 	list.PostingDate = mdl.PostingDate
// 	list.DateStr = mdl.DateStr
// 	list.Journal_Type = mdl.Journal_Type
// 	list.Debet = GrandTotal
// 	list.User = User
// 	idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
// 	numberDOC := tk.Sprintf("%04d", idx)
// 	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
// 	list.DocumentNumber = headDOC + numberDOC
// 	list.Acc_Code = 5210
// 	list.Description = Desc
// 	list.Acc_Name = "SALES RETUR"
// 	list.Department = "COMMERCE"
// 	list.Attachment = "SalesCreditMemo"
// 	dataListDetail = append(dataListDetail, list)

// 	if VAT != 0 {
// 		vat := Journal{}
// 		vat.Id = tk.RandomString(32)
// 		vat.No = 2
// 		vat.PostingDate = mdl.PostingDate
// 		vat.DateStr = mdl.DateStr
// 		vat.Journal_Type = mdl.Journal_Type
// 		vat.DocumentNumber = list.DocumentNumber
// 		vat.Credit = VAT
// 		vat.User = User
// 		vat.Acc_Code = 3221
// 		vat.Acc_Name = "VALUE ADDED TAX (VAT) - OUT"
// 		vat.Attachment = "SalesCreditMemo"
// 		vat.Description = Desc
// 		vat.Department = "COMMERCE"
// 		dataListDetail = append(dataListDetail, vat)
// 	}

// 	list2 := Journal{}
// 	list2.Id = tk.RandomString(32)
// 	list2.No = 2
// 	list2.PostingDate = mdl.PostingDate
// 	list2.DateStr = mdl.DateStr
// 	list2.Journal_Type = mdl.Journal_Type
// 	list2.DocumentNumber = list.DocumentNumber
// 	list2.Credit = Amount
// 	if VAT != 0 {
// 		list2.No = 3
// 		list2.Credit = Amount
// 	}
// 	list2.User = User
// 	list2.Acc_Code = 1210
// 	list2.Acc_Name = "ACCOUNT RECEIVABLE"
// 	list2.Description = Desc
// 	list2.Department = "COMMERCE"
// 	list2.Attachment = "SalesCreditMemo"
// 	dataListDetail = append(dataListDetail, list2)

// 	list3 := Journal{}
// 	list3.Id = tk.RandomString(32)
// 	list3.No = 3
// 	list3.PostingDate = mdl.PostingDate
// 	list3.DateStr = mdl.DateStr
// 	list3.Journal_Type = mdl.Journal_Type
// 	list3.DocumentNumber = list.DocumentNumber
// 	list3.Debet = Amount
// 	if VAT != 0 {
// 		list3.No = 4
// 		list3.Debet = Amount
// 	}
// 	list3.User = User
// 	list3.Acc_Code = 1240
// 	list3.Acc_Name = "MERCHANDISE INVENTORY"
// 	list3.Description = Desc
// 	list3.Department = "COMMERCE"
// 	list3.Attachment = "SalesCreditMemo"
// 	dataListDetail = append(dataListDetail, list3)

// 	list4 := Journal{}
// 	list4.Id = tk.RandomString(32)
// 	list4.No = 4
// 	list4.PostingDate = mdl.PostingDate
// 	list4.DateStr = mdl.DateStr
// 	list4.Journal_Type = mdl.Journal_Type
// 	list4.DocumentNumber = list.DocumentNumber
// 	list4.Credit = Amount
// 	if VAT != 0 {
// 		list4.No = 5
// 		list4.Credit = Amount
// 	}
// 	list4.User = User
// 	list4.Acc_Code = 5400
// 	list4.Acc_Name = "COST OF GOODS SOLD"
// 	list4.Description = Desc
// 	list4.Department = "COMMERCE"
// 	list4.Attachment = "SalesCreditMemo"
// 	dataListDetail = append(dataListDetail, list4)

// 	mdl.ListDetail = dataListDetail

// 	e := c.Ctx.Save(mdl)
// 	c.SavetoGeneralLedgerFromSalesCreditMemo(1210, "ACCOUNT RECEIVABLE", Amount, VAT, GrandTotal, User, mdl.IdJournal, Desc, DatePosting)
// 	if e != nil {
// 		tk.Println(e.Error())
// 	}
// 	return ""
// }

// func (c *TransactionController) SavetoGeneralLedgerFromSalesCreditMemo(AccountKredit int, AccountName string, Amount float64, VAT float64, GrandTotal float64, User string, idJournal string, Desc string, DatePosting time.Time) interface{} {
// 	m := DatePosting.UTC().Month()
// 	y := DatePosting.UTC().Year()
// 	mdl := new(MainGeneralLedger)
// 	mdl.ID = tk.RandomString(32)
// 	mdl.IdJournal = idJournal
// 	mdl.PostingDate = DatePosting
// 	mdl.CreateDate = time.Now()
// 	mdl.DateStr = DatePosting.Format("02 Jan 2006")
// 	mdl.Journal_Type = "General"
// 	if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
// 		mdl.Status = "posting"
// 	} else {
// 		mdl.Status = "close"
// 	}
// 	mdl.User = User
// 	mdl.Department = "COMMERCE"
// 	dataListDetail := []GeneralDetail{}
// 	list := GeneralDetail{}
// 	list.Id = tk.RandomString(32)
// 	list.No = 1
// 	list.PostingDate = mdl.PostingDate
// 	list.DateStr = mdl.DateStr
// 	list.Journal_Type = mdl.Journal_Type
// 	list.Debet = GrandTotal
// 	list.User = User
// 	idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
// 	numberDOC := tk.Sprintf("%04d", idx)
// 	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
// 	list.DocumentNumber = headDOC + numberDOC
// 	list.Acc_Code = 5210
// 	list.Description = Desc
// 	list.Acc_Name = "SALES RETUR"
// 	list.Attachment = "SalesCreditMemo"
// 	list.Department = "COMMERCE"
// 	dataListDetail = append(dataListDetail, list)

// 	if VAT != 0 {
// 		vat := GeneralDetail{}
// 		vat.Id = tk.RandomString(32)
// 		vat.No = 2
// 		vat.PostingDate = mdl.PostingDate
// 		vat.DateStr = mdl.DateStr
// 		vat.Journal_Type = mdl.Journal_Type
// 		vat.DocumentNumber = list.DocumentNumber
// 		vat.Credit = VAT
// 		vat.User = User
// 		vat.Acc_Code = 3221
// 		vat.Acc_Name = "VALUE ADDED TAX (VAT) - OUT"
// 		vat.Description = Desc
// 		vat.Attachment = "SalesCreditMemo"
// 		vat.Department = "COMMERCE"
// 		dataListDetail = append(dataListDetail, vat)
// 	}

// 	list2 := GeneralDetail{}
// 	list2.Id = tk.RandomString(32)
// 	list2.No = 2
// 	list2.PostingDate = mdl.PostingDate
// 	list2.DateStr = mdl.DateStr
// 	list2.Journal_Type = mdl.Journal_Type
// 	list2.DocumentNumber = list.DocumentNumber
// 	list2.Credit = GrandTotal
// 	if VAT != 0 {
// 		list2.No = 3
// 		list2.Credit = Amount
// 	}
// 	list2.User = User
// 	list2.Acc_Code = AccountKredit
// 	list2.Acc_Name = AccountName
// 	list2.Description = Desc
// 	list2.Department = "COMMERCE"
// 	list2.Attachment = "SalesCreditMemo"
// 	dataListDetail = append(dataListDetail, list2)
// 	mdl.ListDetail = dataListDetail
// 	e := c.Ctx.Save(mdl)
// 	if e != nil {
// 		tk.Println(e.Error())
// 	}
// 	return ""
// }
func (c *TransactionController) ExportToPdfListSalesCreditMemo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id             bson.ObjectId
		WordGrandtotal string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("SalesCreditMemo").Where(dbox.Eq("_id", p.Id)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []SalesCreditMemo{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	// tk.Println(DATA)

	if DATA.Currency == "USD" {
		for i, _ := range DATA.ListItem {
			DATA.ListItem[i].PriceIDR = 0
			DATA.ListItem[i].AmountIDR = 0
		}
		// discount = DATA.Discount
	} else {
		for i, _ := range DATA.ListItem {
			DATA.ListItem[i].PriceUSD = 0
			DATA.ListItem[i].AmountUSD = 0
		}
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.CustomerCode)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	resultsCustomer := []CustomerModel{}
	e = csr.Fetch(&resultsCustomer, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	cust := resultsCustomer[0]
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
	pdf.SetXY(78, 10)
	pdf.CellFormat(0, 12, "SALES CREDIT MEMO", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(85)

	pdf.SetFont("Century_Gothicb", "B", 10)
	pdf.CellFormat(0, 12, DATA.DocumentNo, "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothic", "", 8)
	y0 := pdf.GetY() + 5
	//
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Kode, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Date", "", "L", false) // date
	date := DATA.DateCreated.Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Name, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
	dueDate := DATA.DateCreated.AddDate(0, 0, 30).Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(20, 3, "Phone :", "", "L", false) // phone
	pdf.SetXY(30, y0)
	pdf.MultiCell(60, 3, cust.NoTelp, "", "L", false) //phone
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+DATA.DocumentNo, "", "L", false) // DocumentNo
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Address, "", "L", false) // address
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
	if cust.Bank != "" && cust.AccountNo != "" {
		pdf.MultiCell(60, 3, cust.Bank+"-"+cust.AccountNo, "", "L", false) //rek bank
	} else {
		pdf.MultiCell(60, 3, "", "", "L", false) //rek bank
	}
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	invHead := []string{"No", "Code Item", "Item", "Qty", "Price", "Disc. Amount", "Amount"}
	widHead := []float64{10.0, 30.0, 50.0, 10.0, 30.0, 30.0, 30.0}
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
	for i, list := range DATA.ListItem {
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
		pdf.MultiCell(widHead[1], 4, list.CodeItem, "", "L", false)
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
	pdf.MultiCell(40, 3, DATA.Description, "", "L", false) // Remark
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
	total := tk.Sprintf("%.2f", DATA.Total)
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
	y0 = pdf.GetY()
	pdf.SetFont("Century_Gothic", "", 8)
	pdf.SetXY(30, y0)
	pdf.MultiCell(150, 3, "Barang yang sudah dibeli tidak bisa di kembalikan / di tukar dengan uang", "", "L", false) // alert
	pdf.Ln(2)
	y0 = pdf.GetY()
	pdf.SetXY(10, y0)
	pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
	pdf.SetXY(30, y0)
	datenow := time.Now().Format("January 02, 2006")
	pdf.MultiCell(150, 3, datenow, "", "L", false) // date print

	e = os.RemoveAll(c.PdfPath + "/salescreditmemo")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/salescreditmemo", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-SalesCreditMemo.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/salescreditmemo"
	// tk.Println(location)
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName
}

func (c *TransactionController) GetDataInvoiceSearch(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		TextSearch string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	filter := []*dbox.Filter{}

	// if p.TextSearch != "" {
	filter = append(filter, dbox.Eq("DocumentNo", p.TextSearch))
	filter = append(filter, dbox.Or(dbox.Eq("Status", "PAID"), dbox.Eq("Status", "POSTING")))
	filter = append(filter, dbox.Eq("CreditMemo", false))
	// }

	crs, e := c.Ctx.Connection.NewQuery().Select().From("Invoice").Where(dbox.And(filter...)).Cursor(nil)

	defer crs.Close()
	results := make([]InvoiceModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	if len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransactionController) GetAutoCInvNumSCM(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := KendoDatasourceQueryFilter{}
	k.GetPayload(&p)

	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Or(dbox.Eq("Status", "PAID"), dbox.Eq("Status", "POSTING")))
	filter = append(filter, dbox.Eq("CreditMemo", false))
	filter = append(filter, p.ToDboxFilter())

	crs, e := c.Ctx.Connection.NewQuery().Select().From("Invoice").Where(filter...).Cursor(nil)
	defer crs.Close()
	results := make([]struct {
		DocumentNo string `bson:"DocumentNo"`
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
