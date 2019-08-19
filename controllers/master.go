package controllers

import (
	"../helpers"
	. "../helpers"
	"../library/strformat"
	"../library/tealeg/xlsx"
	. "../models"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

type MasterController struct {
	*BaseController
}

func (c *MasterController) CoaDefault(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) AccNo(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) TypePurchase(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) TypeStock(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) Category(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) Sales(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) Departement(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) InvoiceCategory(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) LocationMaster(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *MasterController) GetDataCOA(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Start      time.Time
		End        time.Time
		Filter     bool
		TextSearch string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	// crs, e := c.Ctx.Find(new(CoaModel), tk.M{})
	// defer crs.Close()
	// if p.TextSearch != "" {
	filter := []*db.Filter{}
	filter = append(filter, db.Contains("account_name", p.TextSearch))
	filter = append(filter, db.Contains("debet_credit", p.TextSearch))
	filter = append(filter, db.Contains("category", p.TextSearch))
	mainFilter := new(db.Filter)
	mainFilter = db.Or(filter...)
	crs, e := c.Ctx.Connection.NewQuery().From("Coa").Where(mainFilter).Cursor(nil)
	defer crs.Close()
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	// }
	results := make([]CoaModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	if p.Filter && len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}
	StartTime, _ := time.Parse("2006-01-02", p.Start.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", p.End.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)
	var pipes []tk.M
	if p.Filter == true {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}, "Status": tk.M{"$ne": "draft"}}})
	}
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := make([]tk.M, 0)
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultJournal {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("SumDebet")
				results[i].Credit = each.GetFloat64("SumCredit")
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		for _, z := range results {
			parent := z.Main_Acc_Code
			if parent == accCode {
				results[i].Debet = results[i].Debet + z.Debet
				results[i].Credit = results[i].Credit + z.Credit
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
		}
	}
	Debet1100 := 0.0
	Credit1100 := 0.0

	for i, _ := range results {
		accCode := results[i].ACC_Code
		if accCode%1000 == 0 {
			lastCode := accCode + 1000
			for _, each := range results {
				if each.ACC_Code > accCode && each.ACC_Code < lastCode && each.Main_Acc_Code == 0 {
					results[i].Debet = results[i].Debet + each.Debet
					results[i].Credit = results[i].Credit + each.Credit
					results[i].Saldo = results[i].Debet - results[i].Credit
				}
			}
		}
		if results[i].Main_Acc_Code == 1100 {
			Debet1100 = Debet1100 + results[i].Debet
			Credit1100 = Credit1100 + results[i].Credit
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		if accCode == 1100 {
			results[i].Debet = Debet1100
			results[i].Credit = Credit1100
			results[i].Saldo = Debet1100 - Credit1100
		}
	}
	data := struct {
		Data  []CoaModel
		Total int
	}{
		Data:  results,
		Total: len(results),
	}

	return data
}

func (c *MasterController) UploadFiles(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filename := k.Request.FormValue("filename")
	result := tk.NewResult()
	pathToSave, _ := filepath.Abs("assets/docs/datamaster")
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
	// ==================================TEALEG====================================
	model := new(CoaModel)
	excelFileName := fileToProcess
	xlFile, er := xlsx.OpenFile(excelFileName)
	if er != nil {
		return er.Error()
	}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			model.ID = bson.NewObjectId()
			model.ACC_Code, _ = cells[0].Int()
			model.Account_Name, _ = cells[1].String()
			model.Debet_Credit, _ = cells[2].String()
			model.Category, _ = cells[3].String()
			main, _ := cells[4].String()
			if strings.Trim(main, " ") == "" {
				model.Main_Acc_Code = 0
			} else {
				model.Main_Acc_Code, _ = cells[4].Int()
			}
			csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("acc_code", model.ACC_Code)).Cursor(nil)
			if e != nil {
				return result.SetError(e)
			}
			csr.Close()
			if csr.Count() == 0 {
				c.Ctx.Save(model)
			}
		}
	}
	c.LogActivity("Master", "Upload Files", filename, k)
	return result
}
func (c *MasterController) InsertNewCoa(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		AccCode       int
		MainAccCode   int
		AccName       string
		DebetOrCredit string
		Category      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	mdl := new(CoaModel)
	mdl.ID = bson.NewObjectId()
	mdl.ACC_Code = p.AccCode
	mdl.Account_Name = p.AccName
	mdl.Debet_Credit = p.DebetOrCredit
	mdl.Category = p.Category
	mdl.Main_Acc_Code = p.MainAccCode
	var filter []*db.Filter
	filter = append(filter, db.Eq("acc_code", p.AccCode))
	filter = append(filter, db.Eq("account_name", p.AccName))
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Or(filter...)).Cursor(nil)
	if e != nil {
		return e.Error()
	}
	results := []tk.M{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	var sts bool
	var msg string
	csr.Close()
	// tk.Println(csr.Count())
	if csr.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Account has been saved"
	} else {
		sts = false
		msg = "Account is already exist"
	}
	mod := new(AccountBLModel)
	mod.ID = bson.NewObjectId()
	mod.ACC_Code = p.AccCode
	mod.Account_Name = p.AccName
	mod.Active = true
	if p.DebetOrCredit == "DEBET" {
		mod.Type = "AKTIVA"
	} else {
		mod.Type = "PASSIVA"
	}
	e = c.Ctx.Save(mod)
	c.LogActivity("Master", "Insert New COA", p.AccName, k)
	return tk.M{}.Set("status", sts).Set("Message", msg).Set("Data", results)
}
func (c *MasterController) GetDataCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Connection.NewQuery().From("Category").Cursor(nil)
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
func (c *MasterController) GetDataInvoiceCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	crs, e := c.Ctx.Connection.NewQuery().From("InvoiceCategory").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	results := make([]tk.M, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *MasterController) InsertNewInvoiceCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID   bson.ObjectId
		Code string
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	mdl := new(InvoiceCategoryModel)
	mdl.ID = p.ID
	if mdl.ID == "" {
		mdl.ID = bson.NewObjectId()
	}
	kode := ""
	numbStr := ""
	name := strings.ToUpper(strformat.Filter(p.Name, strformat.CharsetAlphaNumeric))[:3]
	numb := 0
	numb = c.GetLastNumberInvoiceCategory()
	if numb < 10 {
		numbStr = "000" + strconv.Itoa(numb)
	} else if numb <= 10 && numb < 100 {
		numbStr = "00" + strconv.Itoa(numb)
	} else if numb <= 100 && numb < 1000 {
		numbStr = "0" + strconv.Itoa(numb)
	} else {
		numbStr = strconv.Itoa(numb)
	}
	kode = "INVCAT/" + name + "/" + numbStr
	mdl.CategoryCode = kode
	mdl.CategoryName = strings.ToUpper(p.Name)
	var filter []*db.Filter
	filter = append(filter, db.Eq("CategoryCode", kode))
	filter = append(filter, db.Eq("CategoryName", strings.ToUpper(p.Name)))
	crs, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceCategory").Where(db.And(filter...)).Cursor(nil)
	if e != nil {
		return e.Error()
	}
	crs.Close()
	var sts bool
	var msg string
	if crs.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Category has been saved"
	} else {
		sts = false
		msg = "Category is already exist"
	}
	if p.ID == "" {
		c.LogActivity("Master", "Insert New Category", p.Code, k)
	} else {
		c.LogActivity("Master", "Update Category", p.Code, k)
	}
	return tk.M{}.Set("Status", sts).Set("Message", msg)
}
func (c *MasterController) InsertNewCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID   bson.ObjectId
		Code string
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}

	mdl := new(CategoryModel)
	mdl.ID = p.ID
	if mdl.ID == "" {
		mdl.ID = bson.NewObjectId()
	}
	kode := ""
	numbStr := ""
	name := strings.ToUpper(strformat.Filter(p.Name, strformat.CharsetAlphaNumeric))[:3]
	numb := 0
	numb = c.GetLastNumberCategory()
	if numb < 10 {
		numbStr = "000" + strconv.Itoa(numb)
	} else if numb <= 10 && numb < 100 {
		numbStr = "00" + strconv.Itoa(numb)
	} else if numb <= 100 && numb < 1000 {
		numbStr = "0" + strconv.Itoa(numb)
	} else {
		numbStr = strconv.Itoa(numb)
	}
	kode = "CAT/" + name + "/" + numbStr
	mdl.Code = kode
	mdl.Name = strings.ToUpper(p.Name)
	var filter []*db.Filter
	filter = append(filter, db.Eq("code", kode))
	filter = append(filter, db.Eq("name", strings.ToUpper(p.Name)))
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Category").Where(db.And(filter...)).Cursor(nil)
	if e != nil {
		return e.Error()
	}
	var sts bool
	var msg string
	csr.Close()
	if csr.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Category has been saved"
	} else {
		sts = false
		msg = "Category is already exist"
	}
	if p.ID == "" {
		c.LogActivity("Master", "Insert New Category", kode, k)
	} else {
		c.LogActivity("Master", "Update Category", p.Code, k)
	}

	return tk.M{}.Set("status", sts).Set("Message", msg)
}

func (c *MasterController) DeleteInvoiceCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		ID   bson.ObjectId
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	result := new(InvoiceCategoryModel)
	e = c.Ctx.GetById(result, p.ID)
	if e != nil {
		c.WriteLog(e)
	}
	e = c.Ctx.Delete(result)

	c.LogActivity("Master", "Delete Category", result.CategoryCode, k)
	return ""
}
func (c *MasterController) GetLastNumberInvoiceCategory() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "invoicecategory")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "invoicecategory"
		model.Lastnumber = 1
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		e = c.Ctx.Save(&sec)
		data.Number = sec.Lastnumber
		data.Msg = "Success"
	}

	return data.Number
}
func (c *MasterController) GetLastNumberCategory() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "category")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "category"
		model.Lastnumber = 1
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		e = c.Ctx.Save(&sec)
		data.Number = sec.Lastnumber
		data.Msg = "Success"
	}

	return data.Number
}
func (c *MasterController) CancelCategory(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) Customer(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) GetDataCustomer(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	filter := []*dbox.Filter{}
	if payload.Has("TextSearchCode") && payload.GetString("TextSearchCode") != "" {
		filter = append(filter, dbox.Contains("Kode", payload.GetString("TextSearchCode")))
	}
	if payload.Has("TextSearchName") && payload.GetString("TextSearchName") != "" {
		filter = append(filter, dbox.Contains("Name", payload.GetString("TextSearchName")))
	}

	csr, e := c.Ctx.Connection.NewQuery().From("Customer").Where(filter...).Cursor(nil)
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
func (c *MasterController) GetCustomerBalance(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		CustomerCode string
	}{}
	k.GetPayload(&p)

	cx, e := c.Ctx.Connection.NewQuery().From("Customer").Where(db.Eq("Kode", p.CustomerCode)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer cx.Close()
	res := []*CustomerModel{}
	e = cx.Fetch(&res, 0, false)
	if e != nil || len(res) == 0 {
		return c.SetResultInfo(true, "", nil)
	}
	cust := res[0]

	//Invoice Inventory
	csr, e := c.Ctx.Connection.NewQuery().From("Invoice").Where(db.And(db.Eq("CustomerCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []*InvoiceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	ret := struct {
		TotalIDR       float64
		TotalUSD       float64
		PaidIDR        float64
		PaidUSD        float64
		BalanceIDR     float64
		BalanceUSD     float64
		NearestDueDate time.Time
	}{}

	ret.NearestDueDate = time.Now().AddDate(200, 0, 0)
	for _, inv := range results {
		ret.TotalIDR += inv.GrandTotalIDR
		ret.TotalUSD += inv.GrandTotalUSD
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.AlreadyPaid
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.AlreadyPaid
		}

		if inv.Status != "PAID" {
			pt := strings.Replace(cust.PaymentTerm, " ", "", -1)
			dur, e := time.ParseDuration(pt)
			if e != nil {
				due := inv.DateCreated.Add(dur)
				if due.Before(ret.NearestDueDate) && due.After(time.Now()) {
					ret.NearestDueDate = due
				}
			}
		}
	}

	//Get Data Invoice Non Inventory
	crs, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceNonInv").Where(db.And(db.Eq("CustomerCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsInvNon := make([]InvoiceNonInvModel, 0)
	e = crs.Fetch(&resultsInvNon, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, inv := range resultsInvNon {
		ret.TotalIDR += inv.GrandTotalIDR
		ret.TotalUSD += inv.GrandTotalUSD
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.AlreadyPaid
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.AlreadyPaid
		}

		if inv.Status != "PAID" {
			pt := strings.Replace(cust.PaymentTerm, " ", "", -1)
			dur, e := time.ParseDuration(pt)
			if e != nil {
				due := inv.DateCreated.Add(dur)
				if due.Before(ret.NearestDueDate) && due.After(time.Now()) {
					ret.NearestDueDate = due
				}
			}
		}
	}

	//Get Data Credit Memo
	crs, e = c.Ctx.Connection.NewQuery().Select().From("SalesCreditMemo").Where(db.And(db.Eq("CustomerCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsCM := make([]SalesCreditMemo, 0)
	e = crs.Fetch(&resultsCM, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, inv := range resultsCM {
		ret.TotalIDR += 0
		ret.TotalUSD += 0
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.GrandTotalIDR
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.GrandTotalUSD
		}

		if inv.Status != "PAID" {
			pt := strings.Replace(cust.PaymentTerm, " ", "", -1)
			dur, e := time.ParseDuration(pt)
			if e != nil {
				due := inv.DateCreated.Add(dur)
				if due.Before(ret.NearestDueDate) && due.After(time.Now()) {
					ret.NearestDueDate = due
				}
			}
		}
	}

	ret.BalanceIDR = ret.TotalIDR - ret.PaidIDR
	ret.BalanceUSD = ret.TotalUSD - ret.PaidUSD

	return c.SetResultInfo(false, "Success", ret)
}

func (c *MasterController) GetSupplierBalance(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		CustomerCode string
	}{}
	k.GetPayload(&p)

	ret := struct {
		TotalIDR       float64
		TotalUSD       float64
		PaidIDR        float64
		PaidUSD        float64
		BalanceIDR     float64
		BalanceUSD     float64
		NearestDueDate time.Time
	}{}

	cx, e := c.Ctx.Connection.NewQuery().From("Customer").Where(db.Eq("Kode", p.CustomerCode)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer cx.Close()
	res := []*CustomerModel{}
	e = cx.Fetch(&res, 0, false)
	if e != nil || len(res) == 0 {
		return c.SetResultInfo(true, "", nil)
	}
	cust := res[0]

	// Get from PO
	csr, e := c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(db.And(db.Eq("SupplierCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []*PurchaseOrder{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	ret.NearestDueDate = time.Now().AddDate(200, 0, 0)
	for _, inv := range results {
		ret.TotalIDR += inv.GrandTotalIDR
		ret.TotalUSD += inv.GrandTotalUSD
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.AlreadyPaid
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.AlreadyPaid
		}
		if inv.Status != "PAID" {
			pt := strings.Replace(cust.PaymentTerm, " ", "", -1)
			dur, e := time.ParseDuration(pt)
			if e != nil {
				due := inv.DatePosting.Add(dur)
				if due.Before(ret.NearestDueDate) && due.After(time.Now()) {
					ret.NearestDueDate = due
				}
			}
		}
	}

	// Get from PI
	csr, e = c.Ctx.Connection.NewQuery().From("PurchaseInventory").Where(db.And(db.Eq("SupplierCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res2 := []*PurchaseInventory{}
	e = csr.Fetch(&res2, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	for _, inv := range res2 {
		ret.TotalIDR += inv.GrandTotalIDR
		ret.TotalUSD += inv.GrandTotalUSD
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.AlreadyPaid
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.AlreadyPaid
		}
	}

	//Get Data Purchase Credit Memo
	csr, e = c.Ctx.Connection.NewQuery().From("PurchaseCreditMemo").Where(db.And(db.Eq("SupplierCode", p.CustomerCode), db.Ne("Status", "DRAFT"))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res3 := make([]PurchaseCreditMemo, 0)
	e = csr.Fetch(&res3, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	for _, inv := range res3 {
		// ret.TotalIDR += inv.GrandTotalIDR
		// ret.TotalUSD += inv.GrandTotalUSD
		ret.TotalIDR += 0
		ret.TotalUSD += 0
		if inv.Currency == "IDR" {
			ret.PaidIDR += inv.GrandTotalIDR
		}
		if inv.Currency == "USD" {
			ret.PaidUSD += inv.GrandTotalUSD
		}
	}

	ret.BalanceIDR = ret.TotalIDR - ret.PaidIDR
	ret.BalanceUSD = ret.TotalUSD - ret.PaidUSD

	return c.SetResultInfo(false, "Success", ret)
}

func (c *MasterController) GetDataDetailCoa(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   time.Time
		DateEnd     time.Time
		Accountcode string
		ParentCode  string
		Filter      bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	Code, _ := strconv.Atoi(p.Accountcode)
	ParentCode, _ := strconv.Atoi(p.ParentCode)
	var accCodeArray []int
	if ParentCode == 0 && Code%1000 == 0 {
		last := Code + 1000
		prn, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(db.Gte("acc_code", Code), db.Lt("acc_code", last))).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer prn.Close()
		res := []CoaModel{}
		e = prn.Fetch(&res, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		for _, each := range res {
			skip := false
			doc := each.ACC_Code
			for _, num := range accCodeArray {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				accCodeArray = append(accCodeArray, each.ACC_Code)
			}
		}
	} else if ParentCode == 0 && Code%1000 != 0 && Code != 1100 {
		prn, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("main_acc_code", Code)).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer prn.Close()
		res := []CoaModel{}
		e = prn.Fetch(&res, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		for _, each := range res {
			skip := false
			doc := each.ACC_Code
			for _, num := range accCodeArray {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				accCodeArray = append(accCodeArray, each.ACC_Code)
			}
		}
	} else if ParentCode == 0 && Code%1000 != 0 && Code == 1100 {
		prn, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Or(db.Eq("main_acc_code", 1100), db.Eq("main_acc_code", 1120))).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer prn.Close()
		res := []CoaModel{}
		e = prn.Fetch(&res, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		for _, each := range res {
			skip := false
			doc := each.ACC_Code
			for _, num := range accCodeArray {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				accCodeArray = append(accCodeArray, each.ACC_Code)
			}
		}
	}
	// tk.Println("=====================", accCodeArray)
	var pipes []tk.M
	if p.Filter == true {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	} else {
		// tk.Println("filter is not active")
		pipes = append(pipes, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
	}
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"IdJournal":      "$IdJournal",
		"MultiJournal":   "$Department",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Department":     "$ListDetail.Department",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
	}})
	if ParentCode == 0 && Code%1000 != 0 {
		pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accCodeArray}}})
	} else if ParentCode == 0 && Code%1000 == 0 {
		lastCode := Code + 1000
		pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$gte": Code, "$lt": lastCode}}})
	} else {
		pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": Code}})
	}

	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	DocNumber := []string{}
	DocSaldoAwal := []string{}
	for _, each := range results {
		if each.GetString("Attachment") != "BEGIN" && each.GetString("Attachment") != "INVOICE" && each.GetString("MultiJournal") != "MULTIJOURNAL" {
			skip := false
			doc := each.GetString("DocumentNumber")
			for _, num := range DocNumber {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				DocNumber = append(DocNumber, each.GetString("DocumentNumber"))
			}
		} else {
			skip := false
			doc := each.GetString("DocumentNumber")
			for _, num := range DocSaldoAwal {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				DocSaldoAwal = append(DocSaldoAwal, each.GetString("DocumentNumber"))
			}
		}
	}
	NewResult := []ledgerModel{}
	//data debet
	var pipess []tk.M
	if p.Filter == true {
		pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	} else {
		pipess = append(pipess, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
	}
	pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
	pipess = append(pipess, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"no":             "$ListDetail.No",
		"idjournal":      "$IdJournal",
		"journal_type":   "$ListDetail.Journal_Type",
		"postingdate":    "$ListDetail.PostingDate",
		"datestr":        "$ListDetail.DateStr",
		"documentnumber": "$ListDetail.DocumentNumber",
		"acc_code":       "$ListDetail.Acc_Code",
		"acc_name":       "$ListDetail.Acc_Name",
		"debet":          "$ListDetail.Debet",
		"credit":         "$ListDetail.Credit",
		"description":    "$ListDetail.Description",
		"department":     "$ListDetail.Department",
		"attachment":     "$ListDetail.Attachment",
		"user":           "$ListDetail.User",
	}})
	if ParentCode == 0 && Code%1000 != 0 {

		pipess = append(pipess, tk.M{"$match": tk.M{"acc_code": tk.M{"$in": accCodeArray}, "credit": tk.M{"$eq": 0.0}}})
	} else if ParentCode == 0 && Code%1000 == 0 {
		lastCode := Code + 1000

		pipess = append(pipess, tk.M{"$match": tk.M{"acc_code": tk.M{"$gte": Code, "$lt": lastCode}, "credit": tk.M{"$eq": 0.0}}})
	} else {
		pipess = append(pipess, tk.M{"$match": tk.M{"acc_code": tk.M{"$eq": Code}, "credit": tk.M{"$eq": 0.0}}})

	}
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultDebet := []ledgerModel{}
	e = csr.Fetch(&ResultDebet, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	for _, each := range ResultDebet {

		NewResult = append(NewResult, each)
	}
	//datacredit

	pipess = []tk.M{}
	if p.Filter == true {
		pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	} else {
		pipess = append(pipess, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
	}
	pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
	pipess = append(pipess, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"no":             "$ListDetail.No",
		"idjournal":      "$IdJournal",
		"journal_type":   "$ListDetail.Journal_Type",
		"postingdate":    "$ListDetail.PostingDate",
		"datestr":        "$ListDetail.DateStr",
		"documentnumber": "$ListDetail.DocumentNumber",
		"acc_code":       "$ListDetail.Acc_Code",
		"acc_name":       "$ListDetail.Acc_Name",
		"debet":          "$ListDetail.Debet",
		"credit":         "$ListDetail.Credit",
		"description":    "$ListDetail.Description",
		"department":     "$ListDetail.Department",
		"attachment":     "$ListDetail.Attachment",
		"user":           "$ListDetail.User",
	}})
	if ParentCode == 0 {
		pipess = append(pipess, tk.M{"$match": tk.M{"documentnumber": tk.M{"$in": DocNumber}}})
	} else {
		pipess = append(pipess, tk.M{"$match": tk.M{"documentnumber": tk.M{"$in": DocNumber}}})
	}
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultCredit := []ledgerModel{}
	e = csr.Fetch(&ResultCredit, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if ParentCode == 0 || ParentCode == 0 && Code%1000 == 0 {

		for _, e := range accCodeArray {
			newCode := e

			pipes = []tk.M{}
			if p.Filter == true {
				pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
			} else {

				pipes = append(pipes, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
			}
			pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
			pipes = append(pipes, tk.M{"$project": tk.M{
				"_id":            "$ListDetail._id",
				"No":             "$ListDetail.No",
				"IdJournal":      "$IdJournal",
				"MultiJournal":   "$Department",
				"Journal_Type":   "$ListDetail.Journal_Type",
				"PostingDate":    "$ListDetail.PostingDate",
				"DateStr":        "$ListDetail.DateStr",
				"DocumentNumber": "$ListDetail.DocumentNumber",
				"Acc_Code":       "$ListDetail.Acc_Code",
				"Acc_Name":       "$ListDetail.Acc_Name",
				"Debet":          "$ListDetail.Debet",
				"Credit":         "$ListDetail.Credit",
				"Description":    "$ListDetail.Description",
				"Department":     "$ListDetail.Department",
				"Attachment":     "$ListDetail.Attachment",
				"User":           "$ListDetail.User",
			}})
			pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": newCode}})

			csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()
			results := make([]tk.M, 0)
			e = csr.Fetch(&results, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			DocNumber = []string{}
			DocSaldoAwal := []string{}
			for _, each := range results {
				if each.GetString("Attachment") != "BEGIN" && each.GetString("Attachment") != "INVOICE" && each.GetString("MultiJournal") != "MULTIJOURNAL" {
					skip := false
					doc := each.GetString("DocumentNumber")
					for _, num := range DocNumber {
						if doc == num {
							skip = true
							break
						}
					}
					if !skip {
						DocNumber = append(DocNumber, each.GetString("DocumentNumber"))
					}
				} else {
					skip := false
					doc := each.GetString("DocumentNumber")
					for _, num := range DocSaldoAwal {
						if doc == num {
							skip = true
							break
						}
					}
					if !skip {
						DocSaldoAwal = append(DocSaldoAwal, each.GetString("DocumentNumber"))
					}
				}
			}
			pipess = []tk.M{}
			if p.Filter == true {
				pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
			} else {
				pipess = append(pipess, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
			}
			pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
			pipess = append(pipess, tk.M{"$project": tk.M{
				"_id":            "$ListDetail._id",
				"no":             "$ListDetail.No",
				"idjournal":      "$IdJournal",
				"journal_type":   "$ListDetail.Journal_Type",
				"postingdate":    "$ListDetail.PostingDate",
				"datestr":        "$ListDetail.DateStr",
				"documentnumber": "$ListDetail.DocumentNumber",
				"acc_code":       "$ListDetail.Acc_Code",
				"acc_name":       "$ListDetail.Acc_Name",
				"debet":          "$ListDetail.Debet",
				"credit":         "$ListDetail.Credit",
				"description":    "$ListDetail.Description",
				"department":     "$ListDetail.Department",
				"attachment":     "$ListDetail.Attachment",
				"user":           "$ListDetail.User",
			}})
			pipess = append(pipess, tk.M{"$match": tk.M{"documentnumber": tk.M{"$in": DocNumber}}})
			csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()
			ResultCredit := []ledgerModel{}
			e = csr.Fetch(&ResultCredit, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			pipess = []tk.M{}
			if p.Filter == true {
				pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
			} else {
				pipess = append(pipess, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
			}
			pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
			pipess = append(pipess, tk.M{"$project": tk.M{
				"_id":            "$ListDetail._id",
				"no":             "$ListDetail.No",
				"idjournal":      "$IdJournal",
				"journal_type":   "$ListDetail.Journal_Type",
				"postingdate":    "$ListDetail.PostingDate",
				"datestr":        "$ListDetail.DateStr",
				"documentnumber": "$ListDetail.DocumentNumber",
				"acc_code":       "$ListDetail.Acc_Code",
				"acc_name":       "$ListDetail.Acc_Name",
				"debet":          "$ListDetail.Debet",
				"credit":         "$ListDetail.Credit",
				"description":    "$ListDetail.Description",
				"department":     "$ListDetail.Department",
				"attachment":     "$ListDetail.Attachment",
				"user":           "$ListDetail.User",
			}})
			pipess = append(pipess, tk.M{"$match": tk.M{"documentnumber": tk.M{"$in": DocSaldoAwal}, "acc_code": tk.M{"$eq": newCode},
				"debet": tk.M{"$eq": 0.0}}})
			csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()
			ResultCreditSaldoAwal := []ledgerModel{}
			e = csr.Fetch(&ResultCreditSaldoAwal, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

			for _, each := range ResultCreditSaldoAwal {
				NewResult = append(NewResult, each)
			}
			// for _, each := range ResultCredit {

			// 	debetTemp := each.Debet

			// 	if each.Credit == 0 && each.Acc_Code != newCode {
			// 		each.Credit = debetTemp
			// 		each.Debet = 0.0

			// 		NewResult = append(NewResult, each)
			// 	}
			// }
			for _, each := range ResultCredit {
				if each.Debet == 0.0 && each.Acc_Code == Code && each.Attachment != "BEGIN" && each.Attachment != "INVOICE" && each.MultiJournal != "MULTIJOURNAL" {
					debetTemp := each.Credit
					for _, each2 := range ResultCredit {
						if each.DocumentNumber == each2.DocumentNumber && each2.Credit == 0.0 {
							each2.Debet = 0.0
							each2.Credit = debetTemp
							NewResult = append(NewResult, each2)
							break
						}
					}
					// each.Credit = debetTemp
					// each.Debet = 0.0

				}
			}

		}

	} else {
		pipess = []tk.M{}
		if p.Filter == true {
			pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
		} else {
			pipess = append(pipess, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
		}
		pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
		pipess = append(pipess, tk.M{"$project": tk.M{
			"_id":            "$ListDetail._id",
			"no":             "$ListDetail.No",
			"idjournal":      "$IdJournal",
			"journal_type":   "$ListDetail.Journal_Type",
			"postingdate":    "$ListDetail.PostingDate",
			"datestr":        "$ListDetail.DateStr",
			"documentnumber": "$ListDetail.DocumentNumber",
			"acc_code":       "$ListDetail.Acc_Code",
			"acc_name":       "$ListDetail.Acc_Name",
			"debet":          "$ListDetail.Debet",
			"credit":         "$ListDetail.Credit",
			"description":    "$ListDetail.Description",
			"department":     "$ListDetail.Department",
			"attachment":     "$ListDetail.Attachment",
			"user":           "$ListDetail.User",
		}})
		pipess = append(pipess, tk.M{"$match": tk.M{"documentnumber": tk.M{"$in": DocSaldoAwal}, "acc_code": tk.M{"$eq": Code},
			"debet": tk.M{"$eq": 0.0}}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		ResultCreditSaldoAwal := []ledgerModel{}
		e = csr.Fetch(&ResultCreditSaldoAwal, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		for _, each := range ResultCreditSaldoAwal {
			NewResult = append(NewResult, each)
		}
		for _, each := range ResultCredit {
			debetTemp := each.Debet

			if each.Credit == 0 && each.Acc_Code != Code {
				each.Credit = debetTemp
				each.Debet = 0.0

				NewResult = append(NewResult, each)
			}
		}
	}

	// BEGINING
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.In("acc_code", Code))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	Begining := 0.0

	data := struct {
		Data     []ledgerModel
		Begining float64
	}{
		Data:     NewResult,
		Begining: Begining,
	}
	return data

}
func (c *MasterController) GetEditCustomer(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(CustomerModel)
	e = c.Ctx.GetById(result, bson.ObjectIdHex(p.Id))
	if e != nil {
		c.WriteLog(e)
	}
	return result
}
func (c *MasterController) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id   string
		Kode string
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(CustomerModel)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		c.WriteLog(e)
	}

	e = c.Ctx.Delete(result)

	if result.Type == "CUSTOMER" {
		c.LogActivity("Master", "Delete Customer", result.Kode, k)
	} else {
		c.LogActivity("Master", "Delete Supplier", result.Kode, k)
	}

	return ""
}

func (c *MasterController) InsertNewCustomer(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Data CustomerModel
	}{}
	err := k.GetPayload(&p)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	numbStr := ""
	kode := ""
	name := strings.ToUpper(strformat.Filter(p.Data.Name, strformat.CharsetAlphaNumeric))[:3]
	numb := 0
	if p.Data.Type == "CUSTOMER" {
		numb, _ = c.GetNextIdSeqCustSupp("customer", name)
		if numb < 10 {
			numbStr = "000" + strconv.Itoa(numb)
		} else if numb <= 10 && numb < 100 {
			numbStr = "00" + strconv.Itoa(numb)
		} else if numb <= 100 && numb < 1000 {
			numbStr = "0" + strconv.Itoa(numb)
		} else {
			numbStr = strconv.Itoa(numb)
		}
		kode = "CUST/" + name + "/" + numbStr
	} else {
		numb, _ = c.GetNextIdSeqCustSupp("supplier", name)
		if numb < 10 {
			numbStr = "000" + strconv.Itoa(numb)
		} else if numb <= 10 && numb < 100 {
			numbStr = "00" + strconv.Itoa(numb)
		} else if numb <= 100 && numb < 1000 {
			numbStr = "0" + strconv.Itoa(numb)
		} else {
			numbStr = strconv.Itoa(numb)
		}
		kode = "SUPP/" + name + "/" + numbStr
	}
	if p.Data.Kode == "" {
		p.Data.Kode = kode
	}

	payload := new(CustomerModel)
	cek := p.Data.ID
	// tk.Println(p.Data.ID)
	if p.Data.ID == "" {
		p.Data.ID = tk.RandomString(32)
		LastNumber := numb
		CollName := strings.ToLower(p.Data.Type)
		c.SaveLastNumber(LastNumber, CollName)
	}

	payload.ID = p.Data.ID
	payload.Kode = p.Data.Kode
	payload.Name = p.Data.Name
	payload.Address = p.Data.Address
	payload.City = p.Data.City
	payload.NoTelp = p.Data.NoTelp
	payload.Owner = p.Data.Owner
	payload.Bank = p.Data.Bank
	payload.AccountNo = p.Data.AccountNo
	payload.NPWP = p.Data.NPWP
	payload.TrxCode = p.Data.TrxCode
	payload.VATReg = p.Data.VATReg
	payload.TrxCode = p.Data.TrxCode
	payload.Email = p.Data.Email
	payload.PaymentTerm = p.Data.PaymentTerm
	payload.Type = p.Data.Type
	payload.SalesCode = p.Data.SalesCode
	payload.DepartmentCode = p.Data.DepartmentCode
	payload.Limit = p.Data.Limit

	c.Ctx.Save(payload)
	if cek == "" {
		if p.Data.Type == "CUSTOMER" {
			c.LogActivity("Master", "Insert Customer", payload.Kode, k)
		} else {
			c.LogActivity("Master", "Insert Supplier", payload.Kode, k)
		}
	} else {
		if p.Data.Type == "CUSTOMER" {
			c.LogActivity("Master", "Update Customer", payload.Kode, k)
		} else {
			c.LogActivity("Master", "Update Supplier", payload.Kode, k)
		}
	}
	return c.SetResultOK(nil)
}

func (c *MasterController) SaveLastNumber(LastNumber int, CollName string) interface{} {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(db.Eq("collname", CollName)).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	sec := result[0]
	sec.Lastnumber = LastNumber
	c.Ctx.Save(&sec)

	return c.SetResultInfo(false, "Succes", sec)
}

func (c *MasterController) GetDataTypePurchase(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Connection.NewQuery().From("TypePurchase").Cursor(nil)
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

func (c *MasterController) GetLastNumberTypePurchase() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "typepurchase")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "typepurchase"
		model.Lastnumber = 1
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		e = c.Ctx.Save(&sec)
		data.Number = sec.Lastnumber
		data.Msg = "Success"
	}

	return data.Number
}

func (c *MasterController) InsertNewTypePurchase(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID   bson.ObjectId
		Code string
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	mdl := new(TypePurchaseModel)
	mdl.ID = p.ID
	mdl.Code = p.Code
	if mdl.ID == "" {
		mdl.ID = bson.NewObjectId()
		numbStr := ""
		name := strings.ToUpper(strformat.Filter(p.Name, strformat.CharsetAlphaNumeric))[:3]
		numb := 0
		numb = c.GetLastNumberTypePurchase()
		if numb < 10 {
			numbStr = "000" + strconv.Itoa(numb)
		} else if numb <= 10 && numb < 100 {
			numbStr = "00" + strconv.Itoa(numb)
		} else if numb <= 100 && numb < 1000 {
			numbStr = "0" + strconv.Itoa(numb)
		} else {
			numbStr = strconv.Itoa(numb)
		}
		mdl.Code = "TP/" + name + "/" + numbStr
	}

	mdl.Name = p.Name
	var filter []*db.Filter
	filter = append(filter, db.Eq("code", p.Code))
	filter = append(filter, db.Eq("name", p.Name))
	crs, e := c.Ctx.Connection.NewQuery().Select().From("TypePurchase").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}
	var sts bool
	var msg string
	crs.Close()
	if crs.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Type Purchase has been saved"
	} else {
		sts = false
		msg = "Type Purchase is already exist"
	}
	if p.ID == "" {
		c.LogActivity("Master", "Insert New TypePurchase", p.Code, k)
	} else {
		c.LogActivity("Master", "Update TypePurchase", p.Code, k)
	}
	return tk.M{}.Set("status", sts).Set("Message", msg)
}
func (c *MasterController) DeleteTypePurchase(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := NewTypePurchaseModel()
	e = c.Ctx.GetById(result, p.ID)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	code := result.Code

	e = c.Ctx.Delete(result)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	c.LogActivity("Master", "Delete TypePurchase", code, k)
	return c.SetResultInfo(false, "OK", nil)
}

// func (c *MasterController) GetLastNumberCategory(name string) int {
// 	crs, e := c.Ctx.Connection.NewQuery().From("Category").Select().Cursor(nil)
// 	defer crs.Close()

// 	result := []CategoryModel{}
// 	e = crs.Fetch(&result, 0, false)
// 	if e != nil {
// 		tk.Println(e.Error())
// 	}

// 	data := struct {
// 		Number int
// 		Msg string
// 	} {
// 		Number: 0,
// 		Msg: "",
// 	}
// 	if len(result) == 0 {
// 		model := NewCategoryModel()
// 		model.Code = "CAT/"+name+"/000"
// 		e = c.Ctx.Save(model)
// 		data.Number = 1
// 		return data.Number
// 	}
// 	sec := result[0]
// 	sec.Code = sec.Code + "1"
// 	data.Number = sec.Code
// 	data.Msg = "Success"
// 	return data.Number
// }

func (c *MasterController) GetLastNumberCustomer() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(db.Eq("collname", "customer")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "customer"
		model.Lastnumber = 0
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
func (c *MasterController) GetLastNumberSupplier() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(db.Eq("collname", "supplier")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
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
		model := NewSequenceCustomerModel()
		model.Collname = "supplier"
		model.Lastnumber = 0
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

func (c *MasterController) DeleteCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		ID   bson.ObjectId
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(CategoryModel)
	e = c.Ctx.GetById(result, p.ID)
	if e != nil {
		c.WriteLog(e)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("Master", "Delete Category", result.Code, k)
	return ""
}
func (c *MasterController) ExportPDFCoa(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Filter    bool
		DateStart time.Time
		DateEnd   time.Time
		Data      []CoaModel
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	DATA := p.Data
	pdf := gofpdf.New("P", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	pdf.AddPage()
	pdf.SetXY(10, 5)
	pdf.SetFont("Century_Gothic", "", 12)
	// pdf.Image(c.LogoFile+"eaciit-logo.png", 12, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(11.5)
	pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(11.5)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
	pdf.SetX(145)
	pdf.SetFont("Century_Gothicb", "B", 14)
	pdf.CellFormat(0, 15, "CHART OF ACCOUNT", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11.5)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
	pdf.SetX(11.5)
	// pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11.5)
	pdf.SetFont("Century_Gothic", "", 11)
	pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
	pdf.SetX(6)

	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)

	pdf.Line(pdf.GetX()+3, pdf.GetY()+9, 198, pdf.GetY()+9)
	pdf.SetFont("Century_Gothic", "", 6)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	if p.Filter == true {
		pdf.CellFormat(10, 10, "Date Periode : ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, p.DateStart.Format("02 January 2006")+" - "+p.DateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
		pdf.Ln(4)
	}
	pdf.GetY()
	pdf.SetX(12.0)
	pdf.CellFormat(10, 10, "Date Created : ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(8)
	coaHead := []string{"Account", "Account Name", "Debit", "Credit", "Saldo", "Debit/Credit", "Category"}
	widthHead := []float64{17.0, 59.0, 25.0, 25.0, 25.0, 22.0, 23.0}
	y0 := pdf.GetY()
	for i, head := range coaHead {
		pdf.SetY(y0)
		x := 6.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 5, head, "1", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 5, head, "1", "C", false)
		}

	}
	pdf.SetY(pdf.GetY())
	y0 = pdf.GetY()
	for _, each := range DATA {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 6.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.ACC_Code), "", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Account_Name, "", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		debit := tk.Sprintf("%.2f", each.Debet)
		debit = c.ConvertToCurrency(debit)
		if each.Debet < 0 {
			debit = "(" + tk.Sprintf("%.2f", each.Debet*-1) + ")"
			debit = c.ConvertToCurrency(debit)
		}
		pdf.MultiCell(widthHead[2], 5, debit, "", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		credit := tk.Sprintf("%.2f", each.Credit)
		credit = c.ConvertToCurrency(credit)
		if each.Credit < 0 {
			credit = "(" + tk.Sprintf("%.2f", each.Credit*-1) + ")"
			credit = c.ConvertToCurrency(credit)
		}
		pdf.MultiCell(widthHead[3], 5, credit, "", "R", false)
		a3 := pdf.GetY()
		x += widthHead[3]
		pdf.SetXY(x, y1)
		saldo := tk.Sprintf("%.2f", each.Saldo)
		saldo = c.ConvertToCurrency(saldo)
		if each.Saldo < 0 {
			saldo = "(" + tk.Sprintf("%.2f", each.Saldo*-1) + ")"
			saldo = c.ConvertToCurrency(saldo)
		}
		pdf.MultiCell(widthHead[4], 5, saldo, "", "R", false)
		a4 := pdf.GetY()
		x += widthHead[4]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[5], 5, each.Debet_Credit, "", "L", false)
		a5 := pdf.GetY()
		x += widthHead[5]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[6], 5, each.Category, "", "L", false)
		allA := []float64{a0, a1, a2, a3, a4, a5}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest >= 273 {
			pdf.Line(6.0, biggest, x+widthHead[6], biggest)
			pdf.Line(6.0, y0, 6.0, biggest)
			pdf.Line(6.0+widthHead[0], y0, 6.0+widthHead[0], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1], y0, 6.0+widthHead[0]+widthHead[1], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
			pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)

			pdf.AddPage()
			y0 = pdf.GetY()
			pdf.Line(6.0, pdf.GetY(), x+widthHead[6], pdf.GetY())
		}
	}
	pdf.Line(6.0, y0, 6.0, pdf.GetY())
	pdf.Line(6.0+widthHead[0], y0, 6.0+widthHead[0], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1], y0, 6.0+widthHead[0]+widthHead[1], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], pdf.GetY())
	pdf.Line(6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 6.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], pdf.GetY())

	pdf.Line(6.0, pdf.GetY(), +202.0, pdf.GetY())
	// pdf.Line(5.0, y0, 5.0, pdf.GetY())

	e = os.RemoveAll(c.PdfPath + "/master")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/master", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}

	namepdf := "-coa.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/master"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}

func (c *MasterController) GetDataSales(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	filter := []*dbox.Filter{}
	if payload.Has("TextSearchCode") && payload.GetString("TextSearchCode") != "" {
		filter = append(filter, dbox.Contains("SalesID", payload.GetString("TextSearchCode")))
	}
	if payload.Has("TextSearchName") && payload.GetString("TextSearchName") != "" {
		filter = append(filter, dbox.Contains("SalesName", payload.GetString("TextSearchName")))
	}

	csr, e := c.Ctx.Connection.NewQuery().From("Sales").Where(filter...).Cursor(nil)
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

func (c *MasterController) InsertNewSales(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Data SalesModel
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	model := NewSalesModel()
	model.ID = p.Data.ID
	model.SalesID = p.Data.SalesID
	if model.ID == "" {
		model.ID = bson.NewObjectId()
		numbStr := ""
		name := strings.ToUpper(strformat.Filter(p.Data.SalesName, strformat.CharsetAlphaNumeric))[:3]
		numb := 0
		numb = c.GetLastNumberSales()
		if numb < 10 {
			numbStr = "000" + strconv.Itoa(numb)
		} else if numb <= 10 && numb < 100 {
			numbStr = "00" + strconv.Itoa(numb)
		} else if numb <= 100 && numb < 1000 {
			numbStr = "0" + strconv.Itoa(numb)
		} else {
			numbStr = strconv.Itoa(numb)
		}
		model.SalesID = "SALES/" + name + "/" + numbStr
	}
	model.SalesName = p.Data.SalesName
	model.Phone = p.Data.Phone

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	e = c.Ctx.Save(model)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if model.ID == "" {
		c.LogActivity("Master", "Insert New Sales", model.SalesID, k)
	} else {
		c.LogActivity("Master", "Update Sales", model.SalesID, k)
	}
	return c.SetResultInfo(false, "OK", "")
}
func (c *MasterController) DeleteSales(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := new(SalesModel)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("Sales", "Delete Sales", result.SalesID, k)
	return c.SetResultInfo(false, "OK", nil)
}

func (c *MasterController) GetLastNumberSales() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "sales")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "sales"
		model.Lastnumber = 1
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		e = c.Ctx.Save(&sec)
		data.Number = sec.Lastnumber
		data.Msg = "Success"
	}

	return data.Number
}

func (c *MasterController) GetDataTypeStock(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Connection.NewQuery().From("TypeStock").Cursor(nil)
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
func (c *MasterController) GetLastNumberTypeStock() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "typestock")).Cursor(nil)
	defer crs.Close()

	result := []SequenceCustomerModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		tk.Println(e.Error())
		// return c.SetResultInfo(true, e.Error(), nil)
	}

	data := struct {
		Number int
		Msg    string
	}{
		Number: 0,
		Msg:    "",
	}
	if len(result) == 0 {
		model := NewSequenceCustomerModel()
		model.Collname = "typestock"
		model.Lastnumber = 1
		e = c.Ctx.Save(model)
		data.Number = 1
		data.Msg = "Success"
		return data.Number
	} else {
		sec := result[0]
		sec.Lastnumber = sec.Lastnumber + 1
		e = c.Ctx.Save(&sec)
		data.Number = sec.Lastnumber
		data.Msg = "Success"
	}

	return data.Number
}

func (c *MasterController) InsertNewTypeStock(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID   bson.ObjectId
		Code string
		Name string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	mdl := new(TypeStockModel)
	mdl.ID = p.ID
	mdl.Code = p.Code
	if mdl.ID == "" {
		mdl.ID = bson.NewObjectId()
		numbStr := ""
		name := strings.ToUpper(p.Name[:3])
		numb := 0
		numb = c.GetLastNumberTypeStock()
		if numb < 10 {
			numbStr = "000" + strconv.Itoa(numb)
		} else if numb <= 10 && numb < 100 {
			numbStr = "00" + strconv.Itoa(numb)
		} else if numb <= 100 && numb < 1000 {
			numbStr = "0" + strconv.Itoa(numb)
		} else {
			numbStr = strconv.Itoa(numb)
		}
		mdl.Code = "TS/" + name + "/" + numbStr
	}

	mdl.Name = p.Name
	var filter []*db.Filter
	filter = append(filter, db.Eq("code", p.Code))
	filter = append(filter, db.Eq("name", p.Name))
	crs, e := c.Ctx.Connection.NewQuery().Select().From("TypeStock").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}
	var sts bool
	var msg string
	crs.Close()
	if crs.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Type Stock has been saved"
	} else {
		sts = false
		msg = "Type Stock is already exist"
	}
	if p.ID == "" {
		c.LogActivity("Master", "Insert New TypeStock", p.Code, k)
	} else {
		c.LogActivity("Master", "Update TypeStock", p.Code, k)
	}
	return tk.M{}.Set("status", sts).Set("Message", msg)
}
func (c *MasterController) DeleteTypeStock(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := NewTypeStockModel()
	e = c.Ctx.GetById(result, p.ID)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	code := result.Code

	e = c.Ctx.Delete(result)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	c.LogActivity("Master", "Delete TypeStock", code, k)
	return c.SetResultInfo(false, "OK", nil)
}
func (c *MasterController) GetDataLocation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// If using filter
	/*
		var p LocationModel
		e := k.GetPayload(&p)
		if e != nil {
			return CreateResult(false, nil, e.Error())
		}
		filter := []*db.Filter{}
		filter = append(filter, db.Contains("Description", p.Description))
		filter = append(filter, db.Contains("LocationName", p.LocationName))
		mainFilter := new(db.Filter)
		mainFilter = db.Or(filter...)
		crs, e := c.Ctx.Connection.NewQuery().From("Location").Where(mainFilter).Cursor(nil)
		defer crs.Close()
		if e != nil {
			CreateResult(false, nil, e.Error())
		}
	*/
	csr, e := c.Ctx.Connection.NewQuery().From("Location").Cursor(nil)
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
func (c *MasterController) GetDataLocationByUser(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filter := []*db.Filter{}
	locationid := k.Session("locationid").(int)
	if locationid != 1000 {
		filter = append(filter, db.Eq("LocationID", locationid))
		filter = append(filter, db.Eq("Main_LocationID", locationid))
	} else {
		filter = append(filter, db.Ne("LocationID", 0))
	}
	csr, e := c.Ctx.Connection.NewQuery().From("Location").Where(db.Or(filter...)).Cursor(nil)
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
func (c *MasterController) GetUserLocation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	locationid := k.Session("locationid").(int)
	if locationid > 0 {
		filter := []*db.Filter{}
		filter = append(filter, db.Eq("LocationID", locationid))
		csr, e := c.Ctx.Connection.NewQuery().From("Location").Where(db.Or(filter...)).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		var results []LocationModel
		e = csr.Fetch(&results, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		if len(results) < 0 {
			return c.SetResultInfo(true, "Location not found", nil)
		}
		return c.SetResultInfo(false, "Success", results[0])
	}
	return c.SetResultInfo(true, "Session not found", nil)
}
func (c *MasterController) InsertNewLocation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	src := struct {
		Data struct {
			ID              bson.ObjectId
			LocationID      int
			LocationName    string
			Description     string
			Main_LocationID int
		}
	}{}

	e := k.GetPayload(&src)

	p := src.Data

	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	mdl := new(LocationModel)
	mdl.ID = p.ID
	mdl.LocationID = p.LocationID
	mdl.LocationName = p.LocationName
	mdl.Description = p.Description
	mdl.Main_LocationID = p.Main_LocationID

	if mdl.ID == "" {
		mdl.ID = bson.NewObjectId()
	}

	if len(strconv.Itoa(mdl.LocationID)) != 4 {
		return tk.M{}.Set("status", false).Set("Message", "LocationID must be 4 digit numeric")
	}

	var filter []*db.Filter
	filter = append(filter, db.Eq("LocationID", p.LocationID))
	filter = append(filter, db.Ne("_id", mdl.ID))
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Location").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}
	var sts bool
	var msg string
	crs.Close()
	if crs.Count() == 0 {
		e = c.Ctx.Save(mdl)
		sts = true
		msg = "Location has been saved"
		csr, e := c.Ctx.Connection.NewQuery().Select().From("Inventory").Where(db.Eq("StoreLocation", 1000)).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		resultinv := make([]InventoryModel, 0)
		e = csr.Fetch(&resultinv, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		for _, fields := range resultinv {
			model := NewInventoryModel()
			model.ID = bson.NewObjectId()
			model.INVID = fields.INVID
			model.INVDesc = fields.INVDesc
			model.Unit = fields.Unit
			model.Type = fields.Type
			model.StoreLocation = p.LocationID
			model.StoreLocationName = p.LocationName
			model.LastDate = time.Now()

			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

			e = c.Ctx.Save(model)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}

		}
	} else {
		sts = false
		msg = "Location is already exist"
	}
	if p.ID == "" {
		c.LogActivity("Master", "Insert New Location", strconv.Itoa(p.LocationID), k)
	} else {
		c.LogActivity("Master", "Update Location", strconv.Itoa(p.LocationID), k)
	}
	return tk.M{}.Set("status", sts).Set("Message", msg)
}

func (c *MasterController) GetDetailBillingCustomer(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		CustomerCode string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	var filter []*db.Filter
	filter = append(filter, db.Eq("CustomerCode", p.CustomerCode))

	var filterKhusus []*db.Filter
	filterKhusus = append(filterKhusus, db.Eq("CustomerCode", p.CustomerCode))
	filterKhusus = append(filterKhusus, db.Ne("Status", "DRAFT"))

	resultAll := []tk.M{}

	//Get Data Invoice Non Inventory
	crs, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceNonInv").Where(db.And(filterKhusus...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsInvNon := make([]InvoiceNonInvModel, 0)
	e = crs.Fetch(&resultsInvNon, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsInvNon {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Invoice (Non Inventory)")
		temp.Set("ID", e.Id)
		temp.Set("CustomerCode", e.CustomerCode)
		temp.Set("DocumentNumber", e.DocumentNo)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", e.DateCreated)
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		resultAll = append(resultAll, temp)
	}

	//Get Data Invoice Inventory
	crs, e = c.Ctx.Connection.NewQuery().Select().From("Invoice").Where(db.And(filterKhusus...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsInv := make([]InvoiceModel, 0)
	e = crs.Fetch(&resultsInv, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsInv {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Invoice")
		temp.Set("ID", e.Id)
		temp.Set("CustomerCode", e.CustomerCode)
		temp.Set("DocumentNumber", e.DocumentNo)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", e.DateCreated)
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		resultAll = append(resultAll, temp)
	}

	//Get Data Sales Payment
	crs, e = c.Ctx.Connection.NewQuery().Select().From("SalesPayment").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsPym := make([]SalesPaymentModel, 0)
	e = crs.Fetch(&resultsPym, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPym {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Sales Payment")
		temp.Set("ID", e.ID)
		temp.Set("CustomerCode", e.CustomerCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", e.DatePosting)
		totalPayment := 0.0
		for _, list := range e.ListDetail {
			totalPayment += list.Receive
		}
		temp.Set("Paid", totalPayment)
		temp.Set("Total", 0)
		resultAll = append(resultAll, temp)
	}

	//Get Data Credit Memo
	crs, e = c.Ctx.Connection.NewQuery().Select().From("SalesCreditMemo").Where(db.And(filterKhusus...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsCM := make([]SalesCreditMemo, 0)
	e = crs.Fetch(&resultsCM, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsCM {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Sales Credit Memo")
		temp.Set("ID", e.Id)
		temp.Set("CustomerCode", e.CustomerCode)
		temp.Set("DocumentNumber", e.DocumentNo)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", e.DateCreated)
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		resultAll = append(resultAll, temp)
	}

	return c.SetResultOK(resultAll)
}

func (c *MasterController) GetDetailBillingSupplierNonInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		SupplierCode string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	var filter []*db.Filter
	filter = append(filter, dbox.Ne("Status", "DRAFT"))
	filter = append(filter, db.Eq("SupplierCode", p.SupplierCode))

	resultAll := []tk.M{}

	//Get Data Purchase Order
	crs, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseOrder").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsPOrd := make([]PurchaseOrder, 0)
	e = crs.Fetch(&resultsPOrd, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPOrd {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Order (Non Inventory)")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		// datecreate, _ := time.Parse(time.RFC3339, e.Get("DatePosting"))
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		temp.Set("Remark", e.Remark)
		resultAll = append(resultAll, temp)
	}

	//Get Data Purchase Invoice
	var filter2 []*db.Filter
	filter2 = append(filter2, dbox.Ne("Status", "PO"))
	filter2 = append(filter2, dbox.Ne("Status", "DRAFT"))
	filter2 = append(filter2, dbox.Ne("DocumentNumberPI", ""))
	filter2 = append(filter2, db.Eq("SupplierCode", p.SupplierCode))

	crs, e = c.Ctx.Connection.NewQuery().Select().From("PurchaseOrder").Where(db.And(filter2...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsPInv := make([]PurchaseOrder, 0)
	e = crs.Fetch(&resultsPInv, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPInv {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Invoice (Non Inventory)")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		temp.Set("Remark", e.Remark)
		resultAll = append(resultAll, temp)
	}

	//Get Data Purchase Payment
	var filter3 []*db.Filter
	filter3 = append(filter3, db.Eq("IsInventory", false))
	filter3 = append(filter3, db.Eq("SupplierCode", p.SupplierCode))
	crs, e = c.Ctx.Connection.NewQuery().Select().From("PurchasePayment").Where(db.And(filter3...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		return e.Error()
	}

	resultsPPym := make([]PurchasePaymentModel, 0)
	e = crs.Fetch(&resultsPPym, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPPym {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Payment (Non Inventory)")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		totalPayment := 0.0
		for _, list := range e.ListDetail {
			totalPayment += list.Payment
		}
		temp.Set("Paid", totalPayment)
		temp.Set("Total", 0)
		temp.Set("Remark", "")
		resultAll = append(resultAll, temp)
	}

	//get data transaction inventory
	dataInventory := c.GetDetailBillingSupplierInventory(p.SupplierCode)
	for _, inv := range dataInventory {
		resultAll = append(resultAll, inv)
	}

	return c.SetResultOK(resultAll)
}

func (c *MasterController) GetDetailBillingSupplierInventory(SupplierCode string) []tk.M {
	// k.Config.OutputType = knot.OutputJson

	p := struct {
		SupplierCode string
	}{}

	// e := k.GetPayload(&p)
	// if e != nil {
	// 	return c.SetResultInfo(true, e.Error(), nil)
	// }

	p.SupplierCode = SupplierCode

	var filter []*db.Filter
	filter = append(filter, dbox.Ne("Status", "DRAFT"))
	filter = append(filter, db.Eq("SupplierCode", p.SupplierCode))
	resultAll := []tk.M{}

	//Get Data Purchase Order
	crs, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseInventory").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		tk.Println(e.Error())
	}

	resultsPOrd := make([]PurchaseInventory, 0)
	e = crs.Fetch(&resultsPOrd, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPOrd {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Order")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		temp.Set("Remark", e.Remark)
		resultAll = append(resultAll, temp)
	}

	//Get Data Purchase Invoice
	var filter2 []*db.Filter
	filter2 = append(filter2, dbox.Ne("Status", "PO"))
	filter2 = append(filter2, dbox.Ne("Status", "DRAFT"))
	filter2 = append(filter2, dbox.Ne("DateStrPI", ""))
	filter2 = append(filter2, db.Eq("SupplierCode", p.SupplierCode))

	crs, e = c.Ctx.Connection.NewQuery().Select().From("PurchaseInventory").Where(db.And(filter2...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		tk.Println(e.Error())
	}

	resultsPInv := make([]PurchaseInventory, 0)
	e = crs.Fetch(&resultsPInv, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPInv {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Invoice")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		temp.Set("Remark", e.Remark)
		resultAll = append(resultAll, temp)
	}

	//Get Data Purchase Payment
	var filter3 []*db.Filter
	// filter3 = append(filter3, db.Eq("Status", "PI"))
	filter3 = append(filter3, db.Eq("IsInventory", true))
	filter3 = append(filter3, db.Eq("SupplierCode", p.SupplierCode))
	crs, e = c.Ctx.Connection.NewQuery().Select().From("PurchasePayment").Where(db.And(filter3...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		tk.Println(e.Error())
	}

	resultsPPym := make([]PurchasePaymentModel, 0)
	e = crs.Fetch(&resultsPPym, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPPym {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Payment")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		totalPayment := 0.0
		for _, list := range e.ListDetail {
			totalPayment += list.Payment
		}
		temp.Set("Paid", totalPayment)
		temp.Set("Total", 0)
		temp.Set("Remark", "")
		resultAll = append(resultAll, temp)
	}

	//Get Data Purchase Credit Memo
	filter4 := []*dbox.Filter{}
	filter4 = append(filter4, dbox.Eq("Status", "PAID"))
	filter4 = append(filter4, dbox.Eq("SupplierCode", p.SupplierCode))
	crs, e = c.Ctx.Connection.NewQuery().Select().From("PurchaseCreditMemo").Where(db.And(filter4...)).Cursor(nil)
	defer crs.Close()
	if e != nil {
		tk.Println(e.Error())
	}

	resultsPCM := make([]PurchaseCreditMemo, 0)
	e = crs.Fetch(&resultsPCM, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	for _, e := range resultsPCM {
		m := tk.M{}
		temp := tk.M{}
		_ = tk.StructToM(e, &m)
		temp.Set("Type", "Purchase Credit Memo")
		temp.Set("ID", e.ID)
		temp.Set("SupplierCode", e.SupplierCode)
		temp.Set("DocumentNumber", e.DocumentNumber)
		temp.Set("DateStr", e.DateStr)
		temp.Set("DateCreated", (e.DatePosting).Format(time.RFC3339))
		temp.Set("Paid", 0)
		temp.Set("Total", e.GrandTotalIDR)
		temp.Set("Remark", e.Remark)
		resultAll = append(resultAll, temp)
	}
	// return c.SetResultOK(resultAll)
	return resultAll
}

func (c *MasterController) DeleteLocation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := new(LocationModel)
	e = c.Ctx.GetById(result, p.ID)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	ID := result.LocationID

	e = c.Ctx.Delete(result)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	c.LogActivity("Master", "Delete Location", strconv.Itoa(ID), k)
	return c.SetResultInfo(false, "OK", nil)
}

// func (c *MasterController) GetLastNumberSequence(colname string) (int, error) {
// 	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(db.Eq("collname", colname)).Cursor(nil)
// 	defer crs.Close()

// 	result := []SequenceCustomerModel{}
// 	e = crs.Fetch(&result, 0, false)
// 	if e != nil {
// 		return 0, e
// 	}

// 	data := struct {
// 		Number int
// 		Msg    string
// 	}{
// 		Number: 0,
// 		Msg:    "",
// 	}
// 	if len(result) == 0 {
// 		model := NewSequenceCustomerModel()
// 		model.Collname = colname
// 		model.Lastnumber = 0
// 		e = c.Ctx.Save(model)
// 		data.Number = 1
// 		data.Msg = "Success"
// 		return data.Number, nil
// 	}
// 	sec := result[0]
// 	sec.Lastnumber = sec.Lastnumber + 1
// 	data.Number = sec.Lastnumber
// 	data.Msg = "Success"

// 	return data.Number, nil
// }
func (c *MasterController) UpdateCoa(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id      bson.ObjectId
		AccName string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	result := new(CoaModel)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	result.Account_Name = p.AccName
	e = c.Ctx.Save(result)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// save in BL account
	csr, e := c.Ctx.Connection.NewQuery().Select().From(NewAccountBLModel().TableName()).Where(db.Eq("acc_code", result.ACC_Code)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	blAcc := AccountBLModel{}
	e = csr.Fetch(&blAcc, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if csr.Count() > 0 {
		blAcc.Account_Name = p.AccName
		e = c.Ctx.Save(&blAcc)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
	}
	c.LogActivity("Master", "Update Coa Name", p.Id.Hex(), k)
	return c.SetResultInfo(false, "Success", nil)
}