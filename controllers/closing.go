package controllers

import (
	. "../models"
	"strconv"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type ClosingController struct {
	*BaseController
}

func (c *ClosingController) Default(k *knot.WebContext) interface{} {
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
func (c *ClosingController) GetDataCOAtoClosing(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		Filter     bool
		TextSearch string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	filter := []*db.Filter{}
	filter = append(filter, db.Contains("account_name", p.TextSearch))
	mainFilter := new(db.Filter)
	mainFilter = db.Or(filter...)
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(mainFilter, db.Ne("main_acc_code", 0))).Cursor(nil)
	// defer crs.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	// csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Ne("main_acc_code", 0)).Cursor(nil)
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }
	defer csr.Close()
	results := []CoaCloseModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if p.Filter && len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}

	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.Eq("monthyear", monthYear)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(ResultBegin) == 0 {
		var begin []tk.M
		begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": PeriodeStart, "$lt": dateStart}, "Status": tk.M{"$eq": "posting"}}})
		begin = append(begin, tk.M{"$unwind": "$ListDetail"})
		begin = append(begin, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", begin).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		ResultBeginThisMonth := make([]tk.M, 0)
		e = csr.Fetch(&ResultBeginThisMonth, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		for i, _ := range results {
			for _, res := range ResultBeginThisMonth {
				if results[i].ACC_Code == res.GetInt("_id") && results[i].Main_Acc_Code != 0 {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Beginning = results[i].Beginning + saldo
				}
			}
		}
	} else {

		for i, _ := range results {
			for _, each := range ResultBegin {
				if each.ACC_Code == results[i].ACC_Code {
					results[i].Beginning = results[i].Beginning + each.Ending
				}
			}
		}

	}
	//Get Transaction
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resultTransaction := make([]tk.M, 0)
	e = csr.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dataEarning := c.getEarning(p.DateStart, p.DateEnd)
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code && results[i].Main_Acc_Code != 0 {
				results[i].Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			}
		}
		// if each.GetInt("_id") == 1210 {
		// 	tk.Println(each, each.GetFloat64("TransactionDeb"), each.GetFloat64("TransactionCred"))
		// }
	}
	for i, _ := range results {
		if results[i].ACC_Code == 4400 {
			results[i].Transaction = dataEarning.Data.(float64)
		}
		if results[i].ACC_Code > 3000 && results[i].ACC_Code < 4400 {
			results[i].Transaction = results[i].Transaction * -1
		}
		if results[i].ACC_Code == 5100 || results[i].ACC_Code == 5200 {
			results[i].Transaction = results[i].Transaction * -1
		}
		if results[i].ACC_Code > 7000 && results[i].ACC_Code < 8000 {
			results[i].Transaction = results[i].Transaction * -1
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		for _, z := range results {
			parent := z.Main_Acc_Code
			if parent == accCode && parent%1000 != 0 {
				if results[i].Beginning == 0.0 {
					results[i].Beginning = results[i].Beginning + z.Beginning
				}
				results[i].Transaction = results[i].Transaction + z.Transaction
			}
		}
		if accCode%1000 == 0 {
			results[i].Beginning = 0
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		if accCode%1000 == 0 {
			lastCode := accCode + 1000
			for _, each := range results {
				if each.ACC_Code > accCode && each.ACC_Code < lastCode && each.Main_Acc_Code != 0 {
					results[i].Beginning = results[i].Beginning + each.Beginning
					results[i].Transaction = results[i].Transaction + each.Transaction
				}
			}
		}
	}
	for i, _ := range results {
		if results[i].Category == "INCOME STATEMENT" {
			results[i].Beginning = 0.0
		}
		results[i].Ending = results[i].Beginning + results[i].Transaction
		if results[i].ACC_Code == 4200 {
			results[i].Ending = results[i].Beginning + results[i].Transaction + dataEarning.Data.(float64)
		}
		if results[i].ACC_Code == 4400 {
			results[i].Ending = 0.0
		}
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *ClosingController) SaveAndClosing(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart      time.Time
		DateEnd        time.Time
		Data           []CoaCloseModel
		Periode        string
		SumBegining    float64
		SumTransaction float64
		SumEndig       float64
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	year := strconv.Itoa(dateStart.Year())
	month := strconv.Itoa(int(dateStart.Month()))
	my := month + year
	monthyear, _ := strconv.Atoi(my)
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	csr, e := c.Ctx.Connection.NewQuery().From("GeneralLedger").Where(db.And(db.Gte("PostingDate", dateStart), db.Lt("PostingDate", dateEnd))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []MainGeneralLedger{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	ModelGeneral := NewGeneralLedger()
	for _, dt := range results {
		//==================General Ledger Save=======================

		ModelGeneral.ID = dt.ID
		ModelGeneral.IdJournal = dt.IdJournal
		ModelGeneral.CreateDate = dt.CreateDate
		ModelGeneral.PostingDate = dt.PostingDate
		ModelGeneral.DateStr = dt.DateStr
		ModelGeneral.User = dt.User
		ModelGeneral.Journal_Type = dt.Journal_Type
		ModelGeneral.Status = "close"
		ModelGeneral.Department = dt.Department
		ModelGeneral.ListDetail = dt.ListDetail
		c.Ctx.Save(ModelGeneral)
	}

	for _, each := range p.Data {
		coa := new(CoaCloseModel)
		coa.Id = bson.NewObjectId()
		coa.PeriodeStart = dateStart
		coa.PeriodeEnd = End
		coa.MonthYear = monthyear
		coa.ACC_Code = each.ACC_Code
		coa.Account_Name = each.Account_Name
		coa.Debet_Credit = each.Debet_Credit
		coa.Beginning = each.Beginning
		coa.Transaction = each.Transaction
		coa.Ending = each.Ending
		coa.Category = each.Category
		coa.Main_Acc_Code = each.Main_Acc_Code
		e = c.Ctx.Save(coa)
		if coa.ACC_Code == 4400 {
			c.SaveToJournal4400(coa.Transaction, coa.PeriodeEnd, k.Session("username").(string))
		}
	}
	last := new(LastClosingModel)
	last.Id = bson.NewObjectId()
	last.LastClosing = End
	last.MonthYear = monthyear
	e = c.Ctx.Save(last)
	cls := new(SumClosingModel)
	cls.Id = bson.NewObjectId()
	cls.PeriodeStart, _ = time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	cls.PeriodeEnd, _ = time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	cls.Periode = p.Periode
	cls.MonthYear = monthyear
	cls.Beginning = p.SumBegining
	cls.Transaction = p.SumTransaction
	cls.Ending = p.SumEndig
	e = c.Ctx.Save(cls)
	c.LogActivity("Closing", "Save And Closing", "Closing", k)
	return c.SetResultInfo(false, "success", nil)
}
func (c *ClosingController) SaveToJournal4400(Value float64, Date time.Time, User string) interface{} {
	m := Date.UTC().Month()
	y := Date.UTC().Year()
	codejurnal := tk.Sprintf("%02d%d", m, y)
	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
	number := tk.Sprintf("%04d", ids)
	mdl := new(MainJournal)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = "JUR/" + codejurnal + "/" + number
	mdl.PostingDate = Date
	mdl.CreateDate = Date
	mdl.DateStr = Date.Format("02 Jan 2006")
	mdl.Journal_Type = "General"
	mdl.Status = "posting"
	mdl.User = User
	dataListDetail := []Journal{}
	if Value >= 0 {
		list := Journal{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = Value
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + Date.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Acc_Code = 4400
		list.Description = "Current Earning Closing"
		list.Attachment = "ClOSING"
		list.Acc_Name = "CURRENT EARNING"
		dataListDetail = append(dataListDetail, list)

		list2 := Journal{}
		list2.Id = tk.RandomString(32)
		list2.No = 2
		list2.PostingDate = mdl.PostingDate
		list2.DateStr = mdl.DateStr
		list2.Journal_Type = mdl.Journal_Type
		list2.DocumentNumber = list.DocumentNumber
		list2.Credit = Value
		list2.User = User
		list2.Acc_Code = 4200
		list2.Acc_Name = "RETAINED EARNING"
		list2.Attachment = "CLOSING"
		list2.Description = "Current Earning Closing"
		dataListDetail = append(dataListDetail, list2)
		mdl.ListDetail = dataListDetail
	} else {
		list := Journal{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = Value * -1
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + Date.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Acc_Code = 4200
		list.Description = "Current Earning Closing"
		list.Attachment = "ClOSING"
		list.Acc_Name = "RETAINED EARNING"
		dataListDetail = append(dataListDetail, list)

		list2 := Journal{}
		list2.Id = tk.RandomString(32)
		list2.No = 2
		list2.PostingDate = mdl.PostingDate
		list2.DateStr = mdl.DateStr
		list2.Journal_Type = mdl.Journal_Type
		list2.DocumentNumber = list.DocumentNumber
		list2.Credit = Value * -1
		list2.User = User
		list2.Acc_Code = 4400
		list2.Acc_Name = "CURRENT EARNING"
		list2.Attachment = "CLOSING"
		list2.Description = "Current Earning Closing"
		dataListDetail = append(dataListDetail, list2)
		mdl.ListDetail = dataListDetail
	}
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	c.SavetoGeneralLedger4400(Value, mdl.IdJournal, "Current Earning Closing", Date, User)
	return ""
}
func (c *ClosingController) SavetoGeneralLedger4400(Value float64, idJournal string, Desc string, Date time.Time, User string) interface{} {
	m := Date.UTC().Month()
	y := Date.UTC().Year()
	mdl := new(MainGeneralLedger)
	mdl.ID = tk.RandomString(32)
	mdl.IdJournal = idJournal
	mdl.PostingDate = Date
	mdl.CreateDate = Date
	mdl.DateStr = Date.Format("02 Jan 2006")
	mdl.Journal_Type = "General"
	mdl.Status = "close"
	mdl.User = User
	dataListDetail := []GeneralDetail{}
	if Value >= 0 {
		list := GeneralDetail{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = Value
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + Date.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Acc_Code = 4400
		list.Description = "Current Earning Closing"
		list.Attachment = "CLOSING"
		list.Acc_Name = "CURRENT EARNING"
		dataListDetail = append(dataListDetail, list)

		list2 := GeneralDetail{}
		list2.Id = tk.RandomString(32)
		list2.No = 2
		list2.PostingDate = mdl.PostingDate
		list2.DateStr = mdl.DateStr
		list2.Journal_Type = mdl.Journal_Type
		list2.DocumentNumber = list.DocumentNumber
		list2.Credit = Value
		list2.User = User
		list2.Acc_Code = 4200
		list2.Acc_Name = "RETAINED EARNING"
		list2.Attachment = "CLOSING"
		list2.Description = "Current Earning Closing"
		dataListDetail = append(dataListDetail, list2)
		mdl.ListDetail = dataListDetail
	} else {
		list := GeneralDetail{}
		list.Id = tk.RandomString(32)
		list.No = 1
		list.PostingDate = mdl.PostingDate
		list.DateStr = mdl.DateStr
		list.Journal_Type = mdl.Journal_Type
		list.Debet = Value * -1
		list.User = User
		idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
		numberDOC := tk.Sprintf("%04d", idx)
		headDOC := "GEM/" + Date.Format("020106") + "/"
		list.DocumentNumber = headDOC + numberDOC
		list.Acc_Code = 4200
		list.Description = "RETAINED EARNING"
		list.Attachment = "CLOSING"
		list.Acc_Name = "CURRENT EARNING"
		dataListDetail = append(dataListDetail, list)

		list2 := GeneralDetail{}
		list2.Id = tk.RandomString(32)
		list2.No = 2
		list2.PostingDate = mdl.PostingDate
		list2.DateStr = mdl.DateStr
		list2.Journal_Type = mdl.Journal_Type
		list2.DocumentNumber = list.DocumentNumber
		list2.Credit = Value * -1
		list2.User = User
		list2.Acc_Code = 4400
		list2.Acc_Name = "CURRENT EARNING"
		list2.Attachment = "CLOSING"
		list2.Description = "Current Earning Closing"
		dataListDetail = append(dataListDetail, list2)
		mdl.ListDetail = dataListDetail
	}
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	return ""
}
func (c *ClosingController) GetLastDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$sort": tk.M{"lastclosing": 1}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":      "$monthyear",
		"lastdate": tk.M{"$last": "$lastclosing"},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"lastdate": -1}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Last_Closing").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	result := []tk.M{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	return c.SetResultInfo(false, "success", result)
}
func (c *ClosingController) getEarning(DateStart time.Time, DateEnd time.Time) ResultInfo {
	type dataTransaction struct {
		Id          bson.ObjectId
		AccountCode int
		Transaction float64
	}
	dateStart, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 5000, "$lt": 6000}}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res := make([]tk.M, 0)
	e = csr.Fetch(&res, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	GrossProfit := 0.0
	for _, each := range res {
		if each.GetInt("_id") != 5300 || each.GetInt("_id") != 5400 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			if each.GetInt("_id") == 5100 || each.GetInt("_id") == 5200 {
				trans = trans * -1
				GrossProfit = GrossProfit + trans
			}
			if each.GetInt("_id") == 5600 {
				GrossProfit = GrossProfit - trans
			}
		}
	}
	//OPERATING EXPENSE
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 6000, "$lt": 7000}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resoperex := make([]tk.M, 0)
	e = csr.Fetch(&resoperex, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totOpex := 0.0
	for _, each := range resoperex {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totOpex = totOpex + trans
	}
	netProfit := GrossProfit - totOpex
	//PENGHASILAN LAIN LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 7000, "$lt": 8000}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	reslainlain := make([]tk.M, 0)
	e = csr.Fetch(&reslainlain, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totLainLain := 0.0
	for _, each := range reslainlain {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totLainLain = totLainLain + trans
	}
	totLainLain = totLainLain * -1
	//BEBAN LAIN - LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 8000, "$lt": 9000}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resBebanlainlain := make([]tk.M, 0)
	e = csr.Fetch(&resBebanlainlain, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totBebanLainLain := 0.0
	for _, each := range resBebanlainlain {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totBebanLainLain = totBebanLainLain + trans
	}
	totAllRevnEx := totLainLain - totBebanLainLain
	totalEarningBefore := netProfit + totAllRevnEx
	//TAX
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               tk.M{"Code": "$ListDetail.Acc_Code", "Name": "$ListDetail.Acc_Name"},
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id.Code": tk.M{"$ne": 6888}, "_id.Name": tk.M{"$eq": "TAX"}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTax := make([]tk.M, 0)
	e = csr.Fetch(&resTax, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totTax := 0.0
	for _, each := range resTax {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totTax = totTax + trans
	}
	totAfterTax := totalEarningBefore - totTax
	type modelTemp struct {
		amount float64
	}
	model := new(modelTemp)
	model.amount = totAfterTax
	return c.SetResultInfo(false, "success", totAfterTax)
}
