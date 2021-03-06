package controllers

import (
	. "../helpers"
	"../library/tealeg/xlsx"
	. "../models"
	"os"
	"path/filepath"
	"time"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

func (c *TransactionController) InsertNewAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := struct {
		Data string
	}{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	data := NewAssetModel()
	tk.UnjsonFromString(payload.Data, &data)
	err = c.Ctx.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Asset").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataAsset := []AssetModel{}
	e = csr.Fetch(&dataAsset, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	date := data.PostingDate
	monthYear := date.Format("012006")
	csr, e = c.Ctx.Connection.NewQuery().Select().From("HistoryAsset").Where(db.Eq("MonthYear", monthYear)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataHistory := []HistoryAsset{}
	e = csr.Fetch(&dataHistory, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if csr.Count() != 0 {
		dataHis := dataHistory[0]
		history := new(HistoryAsset)
		history.ID = dataHis.ID
		history.MonthYear = dataHis.MonthYear
		history.DateCreated = dataHis.DateCreated
		listAsset := []ListHistoryAsset{}
		listAsset = dataHis.ListAsset
		set := ListHistoryAsset{}
		set.Id = data.ID
		set.Description = data.Description
		set.Category = data.Category
		set.Qty = data.Qty
		set.Price = data.Price
		set.Total = data.Total
		set.DatePeriod = data.DatePeriod
		set.PostingDate = data.PostingDate
		set.SumDepreciation = data.SumDepreciation
		set.MonthlyDepreciation = data.MonthlyDepreciation
		set.User = data.User
		// tk.Println(set)
		listAsset = append(listAsset, set)
		history.ListAsset = listAsset
		e = c.Ctx.Save(history)
	}
	c.LogActivity("Asset", "Insert Asset", data.Description, k)

	return c.SetResultInfo(false, "Success", data)
}

func (c *TransactionController) GetDataAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateFilter time.Time
		TextSearch string
		Filter     bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().From("Asset").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if p.Filter == true {
		dateFilter, _ := time.Parse("2006-01-02", p.DateFilter.Format("2006-01-02"))
		date := time.Now()
		dateNow := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		filter := []*db.Filter{}
		if dateFilter.After(dateNow) == false {
			monthyear := dateFilter.Format("012006")
			pipes := []tk.M{}
			pipes = append(pipes, tk.M{"$match": tk.M{"MonthYear": tk.M{"$eq": monthyear}}})
			pipes = append(pipes, tk.M{"$unwind": "$ListAsset"})
			pipes = append(pipes, tk.M{"$project": tk.M{
				"ID":                  "$ListAsset._id",
				"Description":         "$ListAsset.Description",
				"Category":            "$ListAsset._id",
				"Qty":                 "$ListAsset.Qty",
				"Price":               "$ListAsset.Price",
				"Total":               "$ListAsset.Total",
				"PostingDate":         "$ListAsset.PostingDate",
				"DatePeriod":          "$ListAsset.DatePeriod",
				"SumDepreciation":     "$ListAsset.SumDepreciation",
				"MonthlyDepreciation": "$ListAsset.MonthlyDepreciation",
				"User":                "$ListAsset.User",
			}})
			pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lte": dateFilter}}})
			csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("HistoryAsset").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()
			results = make([]tk.M, 0)
			e = csr.Fetch(&results, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

			if len(results) == 0 {
				csr, e := c.Ctx.Connection.NewQuery().From("Asset").Where(db.Lte("PostingDate", dateFilter)).Cursor(nil)
				if e != nil {
					return c.SetResultInfo(true, e.Error(), nil)
				}
				defer csr.Close()
				results = make([]tk.M, 0)
				e = csr.Fetch(&results, 0, false)
				if e != nil {
					return c.SetResultInfo(true, e.Error(), nil)
				}
				// tk.Println(results)
			}
		}
		if p.TextSearch != "" {
			filter = append(filter, db.Contains("Description", p.TextSearch))
			query := tk.M{}.Set("where", db.Or(filter...))
			csr, e := c.Ctx.Find(new(AssetModel), query)

			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()

			results = make([]tk.M, 0)
			e = csr.Fetch(&results, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			if p.Filter && len(results) == 0 {
				return c.SetResultInfo(true, "Please refine your search", nil)
			}
		}
	}

	return c.SetResultInfo(false, "Success", results)
}
func (c *TransactionController) ImportExcelAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	filename := k.Request.FormValue("filename")
	pathToSave, _ := filepath.Abs("assets/docs/asset")
	os.MkdirAll(pathToSave, 0777)
	e, _ := UploadHandlerCopy(k, "filedoc", pathToSave+"/")
	if e != nil {
		// tk.Println("Error : " + e.Error())
		return e.Error()
	}

	fileToProcess := pathToSave + "/" + filename

	_, err := os.Stat(fileToProcess)
	if os.IsNotExist(err) {
		// tk.Println(err.Error())
		return err.Error()
	}
	model := new(AssetModel)
	excelFileName := fileToProcess
	xlFile, er := xlsx.OpenFile(excelFileName)
	if er != nil {
		return er.Error()
	}
	for _, sheet := range xlFile.Sheet {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			model.ID = bson.NewObjectId()
			model.Description, _ = cells[1].String()
			model.Category, _ = cells[2].String()
			model.Qty, _ = cells[3].Int()
			model.Price, _ = cells[4].Float()
			model.Total, _ = cells[5].Float()
			postingDate, _ := cells[6].Float()
			model.PostingDate = xlsx.TimeFromExcelTime(postingDate, false)
			datePeriode, _ := cells[7].Float()
			model.DatePeriod = xlsx.TimeFromExcelTime(datePeriode, false)
			model.SumDepreciation, _ = cells[8].Int()
			model.MonthlyDepreciation, _ = cells[9].Float()
			model.User = k.Session("username").(string)
			c.Ctx.Save(model)
			c.LogActivity("Asset", "Import Asset", model.ID.Hex(), k)
		}
	}
	return c.SetResultInfo(false, "Success", nil)
}
func (c *TransactionController) EditDataAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Data AssetModel
		ID   string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateNow := time.Now()
	dateForHistory := dateNow.AddDate(0, -1, 0)
	monthYear := dateForHistory.Format("012006")
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Asset").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataAsset := []AssetModel{}
	e = csr.Fetch(&dataAsset, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From("HistoryAsset").Where(db.Eq("MonthYear", monthYear)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataHistory := []HistoryAsset{}
	e = csr.Fetch(&dataHistory, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if csr.Count() == 0 {
		history := new(HistoryAsset)
		history.ID = bson.NewObjectId()
		history.MonthYear = monthYear
		history.DateCreated = dateForHistory
		listAsset := []ListHistoryAsset{}
		for _, each := range dataAsset {
			set := ListHistoryAsset{}
			set.Id = each.ID
			set.Description = each.Description
			set.Category = each.Category
			set.Qty = each.Qty
			set.Price = each.Price
			set.Total = each.Total
			set.DatePeriod = each.DatePeriod
			set.PostingDate = each.PostingDate
			set.SumDepreciation = each.SumDepreciation
			set.MonthlyDepreciation = each.MonthlyDepreciation
			set.User = each.User
			listAsset = append(listAsset, set)
		}
		history.ListAsset = listAsset
		e = c.Ctx.Save(history)
	}
	Data := p.Data
	asset := new(AssetModel)
	asset.ID = bson.ObjectIdHex(p.ID)
	asset.Description = Data.Description
	asset.Category = Data.Category
	asset.Qty = Data.Qty
	asset.Price = Data.Price
	asset.Total = Data.Total
	asset.DatePeriod = Data.DatePeriod
	asset.PostingDate = Data.PostingDate
	asset.SumDepreciation = Data.SumDepreciation
	asset.MonthlyDepreciation = Data.MonthlyDepreciation
	asset.User = Data.User
	e = c.Ctx.Save(asset)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	c.LogActivity("Asset", "Update Asset", asset.Description, k)

	return c.SetResultInfo(false, "Success", Data)
}
func (c *TransactionController) DeleteDataAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateNow := time.Now()
	dateForHistory := dateNow.AddDate(0, -1, 0)
	monthYear := dateForHistory.Format("012006")
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Asset").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataAsset := []AssetModel{}
	e = csr.Fetch(&dataAsset, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From("HistoryAsset").Where(db.Eq("MonthYear", monthYear)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	dataHistory := []HistoryAsset{}
	e = csr.Fetch(&dataHistory, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if csr.Count() == 0 {
		history := new(HistoryAsset)
		history.ID = bson.NewObjectId()
		history.MonthYear = monthYear
		history.DateCreated = dateForHistory
		listAsset := []ListHistoryAsset{}
		for _, each := range dataAsset {
			set := ListHistoryAsset{}
			set.Id = each.ID
			set.Description = each.Description
			set.Category = each.Category
			set.Qty = each.Qty
			set.Price = each.Price
			set.Total = each.Total
			set.DatePeriod = each.DatePeriod
			set.PostingDate = each.PostingDate
			set.SumDepreciation = each.SumDepreciation
			set.MonthlyDepreciation = each.MonthlyDepreciation
			set.User = each.User
			listAsset = append(listAsset, set)
		}
		history.ListAsset = listAsset
		e = c.Ctx.Save(history)
	}
	result := new(AssetModel)
	e = c.Ctx.GetById(result, bson.ObjectIdHex(p.ID))
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	e = c.Ctx.Delete(result)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	c.LogActivity("Asset", "Delete Asset", result.Description, k)

	return c.SetResultInfo(false, "Success", nil)
}
func (c *TransactionController) GetDataDepreciationAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct{}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("AssetDepreciation").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []AssetDepreciationModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *TransactionController) SaveDataDepreciationAndJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Data []AssetDepreciationModel
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(p.Data)
	mod := new(AssetDepreciationModel)
	for _, each := range p.Data {
		//save data asset deprexiation
		mod.ID = bson.NewObjectId()
		mod.MonthYear = each.MonthYear
		mod.Date = each.Date
		mod.Checked = true
		mod.DateMonthYear = each.DateMonthYear
		mod.IdChecbox = each.IdChecbox
		mod.Amount = each.Amount
		e = c.Ctx.Save(mod)
		//save journal
		desc := "Depreciation Equipment " + each.DateMonthYear
		c.SavetoJournalFromAsset(6126, "FINANCE & ACCOUNTING", each.Amount, k.Session("username").(string), desc, each.Date)
		c.LogActivity("Asset", "Save Data Depreciation And Journal", desc, k)
	}
	return c.SetResultInfo(false, "Success", nil)
}
func (c *TransactionController) SavetoJournalFromAsset(AccountDebet int, Department string, Amount float64, User string, Desc string, DatePosting time.Time) interface{} {
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
	dataListDetail := []Journal{}
	list := Journal{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = Amount
	list.User = User
	list.Department = Department
	idx, _ := c.GetNextIdSeq("DocumentNumber", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Acc_Code = AccountDebet
	list.Description = Desc
	list.Attachment = ""
	list.Acc_Name = "DEPRECIATION EQUIPMENT"
	dataListDetail = append(dataListDetail, list)

	list2 := Journal{}
	list2.Id = tk.RandomString(32)
	list2.No = 2
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = list.DocumentNumber
	list2.Credit = Amount
	list2.User = User
	list2.Acc_Code = 2800
	list2.Acc_Name = "ACCUMULATED DEPRECIATION OFFICE EQUIPMENT"
	list2.Attachment = ""
	list2.Description = Desc
	list2.Department = Department
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	c.SavetoGeneralLedgerFromAsset(AccountDebet, Department, Amount, User, mdl.IdJournal, Desc, DatePosting)
	return ""
}
func (c *TransactionController) SavetoGeneralLedgerFromAsset(AccountDebet int, Department string, Amount float64, User string, idJournal string, Desc string, DatePosting time.Time) interface{} {
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
	// if int(mdl.CreateDate.Month()) <= int(mdl.PostingDate.Month()) && mdl.CreateDate.Year() == mdl.PostingDate.Year() {
	mdl.Status = "posting"
	// } else {
	// 	mdl.Status = "close"
	// }
	mdl.User = User
	dataListDetail := []GeneralDetail{}
	list := GeneralDetail{}
	list.Id = tk.RandomString(32)
	list.No = 1
	list.PostingDate = mdl.PostingDate
	list.DateStr = mdl.DateStr
	list.Journal_Type = mdl.Journal_Type
	list.Debet = Amount
	list.User = User
	list.Department = Department
	idx, _ := c.GetNextIdSeq("DocumentNumberGL", mdl.Journal_Type, int(m), y)
	numberDOC := tk.Sprintf("%04d", idx)
	headDOC := "GEM/" + DatePosting.Format("020106") + "/"
	list.DocumentNumber = headDOC + numberDOC
	list.Acc_Code = AccountDebet
	list.Description = Desc
	list.Attachment = ""
	list.Acc_Name = "DEPRECIATION EQUIPMENT"
	dataListDetail = append(dataListDetail, list)

	list2 := GeneralDetail{}
	list2.Id = tk.RandomString(32)
	list2.No = 2
	list2.PostingDate = mdl.PostingDate
	list2.DateStr = mdl.DateStr
	list2.Journal_Type = mdl.Journal_Type
	list2.DocumentNumber = list.DocumentNumber
	list2.Credit = Amount
	list2.User = User
	list2.Acc_Code = 2800
	list2.Acc_Name = "ACCUMULATED DEPRECIATION OFFICE EQUIPMENT"
	list2.Attachment = ""
	list2.Description = Desc
	list2.Department = Department
	dataListDetail = append(dataListDetail, list2)
	mdl.ListDetail = dataListDetail
	e := c.Ctx.Save(mdl)
	if e != nil {
		tk.Println(e.Error())
	}
	return ""
}
