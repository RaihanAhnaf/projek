package controllers

import (
	"../helpers"
	"../library/tealeg/xlsx"
	. "../models"
	"os"
	"strconv"
	"strings"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

type FinancialController struct {
	*BaseController
}
type IncomePeriode struct {
	Id                 bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code           int
	Main_Code          int
	Acc_Name           string
	Amount_Begining    float64
	Amount_Debet       float64
	Amount_Credit      float64
	Amount_Ending      float64
	Amount_Transaction float64
	Amount             float64
	Periode            []tk.M
}
type IncomePeriodeActiva struct {
	Id                 bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code           int
	Main_Code          int
	Acc_Name           string
	Amount_Begining    float64
	Amount_Debet       float64
	Amount_Credit      float64
	Amount_Ending      float64
	Amount_Transaction float64
	Periode            []tk.M
}

type IncomePeriodePasiva struct {
	Id                 bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code           int
	Main_Code          int
	Acc_Name           string
	Amount_Begining    float64
	Amount_Debet       float64
	Amount_Credit      float64
	Amount_Ending      float64
	Amount_Transaction float64
	Periode            []tk.M
}

// type FieldActiva struct {
// 	Begining float64
// 	Credit   float64
// 	Debet    float64
// 	Ending   float64
// }

func (c *FinancialController) TrialBalance(k *knot.WebContext) interface{} {
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
func (c *FinancialController) BalanceSheet(k *knot.WebContext) interface{} {
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
func (c *FinancialController) IncomeStatement(k *knot.WebContext) interface{} {
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
func (c *FinancialController) AccountBL(Type string) []int {
	// accNIN := []int{4400, 4200}
	// // var NINacc []interface{} = accNIN
	// NINacc := make([]interface{}, len(accNIN))
	// for i, s := range accNIN {
	// 	NINacc[i] = s
	// }
	csr, e := c.Ctx.Connection.NewQuery().Select().From("AccountForBL").Where(db.And(db.Eq("type", Type), db.Eq("active", true))).Cursor(nil)
	if e != nil {
		tk.Println(e.Error())
		return nil
	}
	defer csr.Close()
	results := []AccountBLModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println(e.Error())
		return nil
	}
	Account := []int{}
	for _, each := range results {
		Account = append(Account, each.ACC_Code)
	}
	return Account
}
func (c *FinancialController) GetDataBalanceSheetCOA(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Category  string
		DebCred   string
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	accNIN := []int{4400, 4200}
	Account := c.AccountBL("AKTIVA")
	if p.DebCred == "CREDIT" {
		Account = c.AccountBL("PASSIVA")
	}
	var AccountInter []interface{}
	for _, each := range Account {
		AccountInter = append(AccountInter, each)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.In("acc_code", AccountInter...)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer csr.Close()
	results := []CoaModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	var pipes []tk.M
	// if p.Filter == true {
	StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
	// } else {
	// 	defaultTime := time.Now()
	// 	year := defaultTime.Year()
	// 	month := defaultTime.Month()
	// 	Start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 	Start = Start.AddDate(0, -1, 0)

	// 	End := defaultTime
	// 	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": Start, "$lt": End}}})

	// }
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"_id": tk.M{"$nin": accNIN}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resultJournal := make([]tk.M, 0)
	e = csr.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dataEarning := c.getEarning(p.DateStart, p.DateEnd)
	for _, each := range resultJournal {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("SumDebet")
				results[i].Credit = each.GetFloat64("SumCredit")
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
			if results[i].ACC_Code == 4400 {
				results[i].Saldo = dataEarning.Data.(float64)
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
	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)

	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Ne("main_acc_code", 0))).Cursor(nil)
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
		var pipes []tk.M
		StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTime = EndTime.AddDate(0, 0, -2)
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
		pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		pipes = append(pipes, tk.M{"$match": tk.M{"_id": tk.M{"$eq": 4200}}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		resultJournal2 := make([]tk.M, 0)
		e = csr.Fetch(&resultJournal2, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		for _, each := range resultJournal2 {
			for i, _ := range results {
				if each.GetInt("_id") == results[i].ACC_Code {
					results[i].Debet = each.GetFloat64("SumDebet")
					results[i].Credit = each.GetFloat64("SumCredit")
					results[i].Saldo = results[i].Debet - results[i].Credit
				}
			}
		}
	}
	for i, _ := range results {
		if results[i].ACC_Code > 3000 && results[i].ACC_Code < 4400 {
			results[i].Saldo = results[i].Saldo * -1
		}
	}
	newResults := []CoaModel{}
	for i, _ := range results {
		for j, _ := range ResultBegin {
			if results[i].ACC_Code == ResultBegin[j].ACC_Code {
				results[i].Saldo = results[i].Saldo + ResultBegin[j].Ending
			}
		}
		if results[i].Saldo != 0 {
			newResults = append(newResults, results[i])
		}
	}
	return c.SetResultInfo(true, "success", newResults)
}
func (c *FinancialController) GetDataActivaPeriode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		TextSearch string
		Filter     bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var filter []tk.M
	filter = append(filter, tk.M{"$match": tk.M{"category": tk.M{"$eq": "BALANCE SHEET"}}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":       "$_id",
		"acc_code":  "$acc_code",
		"acc_name":  "$account_name",
		"main_code": "$main_acc_code",
	}})
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Where(db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 1000), db.Lt("acc_code", 3000), db.Ne("main_acc_code", 0)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer crs.Close()
	resultCOA := []IncomePeriodeActiva{}
	e = crs.Fetch(&resultCOA, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	amount := []string{"$ListDetail.Debet", "$listDetail.Credit"}
	dateStart, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)

	periodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodeStart = periodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(periodeStart.Month()))
	year := strconv.Itoa(periodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	crs, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 1000), db.Lt("acc_code", 3000), db.Ne("main_acc_code", 0))).Cursor(nil)
	resultBegin := []CoaCloseModel{}
	e = crs.Fetch(&resultBegin, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	// tk.Println("Amount :", amount)
	// tk.Println("DateStart :", dateStart)
	// tk.Println("DateEnd :", End)
	type activaPeriodeAmount struct {
		Acc_Code  int
		Acc_Name  string
		Month     int
		Year      int
		Amount    float64
		SumDebet  float64
		SumCredit float64
	}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
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
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
		"Year":           tk.M{"$year": "$ListDetail.PostingDate"},
	}})
	// pipes = append(pipes, tk.M{"$group": tk.M{
	// 	"_id":       "$ListDetail.Acc_Code",
	// 	"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
	// 	"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	// }})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"acc_code":  "$ListDetail.Acc_Code",
		"acc_name":  "$ListDetail.Acc_Name",
		"month":     "$ListDetail.Month",
		"year":      "$ListDetail.Year",
		"amount":    "$Amount",
		"sumdebet":  "$Debet",
		"sumcredit": "$Credit",
	}})
	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := []activaPeriodeAmount{}
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultBeginNow := make([]tk.M, 0)
	e = crs.Fetch(&resultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(resultBegin) == 0 {
		for i, _ := range resultCOA {
			for _, res := range resultBeginNow {
				if resultCOA[i].Acc_Code == res.GetInt("_id") {
					Saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + Saldo
				}
			}
		}
	} else {
		for _, each := range resultBegin {
			for _, res := range resultBeginNow {
				if each.ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					each.Ending = each.Ending + saldo
				}
			}
		}
		for i, _ := range resultCOA {
			for _, each := range resultBegin {
				if each.ACC_Code == resultCOA[i].Acc_Code {
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + each.Ending
				}
			}
		}
	}

	var Tpipes []tk.M
	var startTrans time.Time
	var endTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		startTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		endTrans = endTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearTrans := tNow.Year()
		monthTrans := tNow.Month()
		startTrans = time.Date(yearTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans = tNow
	}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": startTrans, "$lte": endTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultTransaction := make([]tk.M, 0)
	e = crs.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultTransaction {
		for i, _ := range resultCOA {
			if each.GetInt("_id") == resultCOA[i].Acc_Code {

				resultCOA[i].Amount_Debet = each.GetFloat64("TransactionDeb")
				resultCOA[i].Amount_Credit = each.GetFloat64("TransactionCred")
				resultCOA[i].Amount_Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			}
		}
	}
	for i, _ := range resultCOA {
		resultCOA[i].Amount_Ending = resultCOA[i].Amount_Begining + resultCOA[i].Amount_Transaction
	}

	// tk.Println("ResultJournal =>", resultJournal)
	for i, _ := range resultJournal {
		if resultJournal[i].Acc_Code == 5110 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code == 5100 || resultJournal[i].Acc_Code == 5200 {

			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code > 7000 && resultJournal[i].Acc_Code < 8000 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
	}
	results := struct {
		DataAcc    []IncomePeriodeActiva
		DataAmount []activaPeriodeAmount
	}{
		DataAcc:    resultCOA,
		DataAmount: resultJournal,
	}
	return c.SetResultInfo(false, "Success", results)

	// ==========================================================================================================
}
func (c *FinancialController) GetDataActiva(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 1000), db.Lt("acc_code", 3000), db.Ne("main_acc_code", 0))).
		Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 1000), db.Lt("acc_code", 3000), db.Ne("main_acc_code", 0))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var begin []tk.M
	begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}}})
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
	ResultBeginNow := make([]tk.M, 0)
	e = csr.Fetch(&ResultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if len(ResultBegin) == 0 {
		for i, _ := range results {
			for _, res := range ResultBeginNow {
				if results[i].ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Begining = results[i].Begining + saldo

				}
			}
		}
	} else {
		if len(ResultBeginNow) == 0 {
			for i, _ := range results {
				for _, each := range ResultBegin {
					if each.ACC_Code == results[i].ACC_Code {
						results[i].Begining = results[i].Begining + each.Ending
					}
				}
			}
		} else {
			for _, each := range ResultBegin {
				for _, res := range ResultBeginNow {
					if each.ACC_Code == res.GetInt("_id") {
						saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
						each.Ending = each.Ending + saldo

					}
				}
			}
			for i, _ := range results {
				for _, each := range ResultBegin {
					if each.ACC_Code == results[i].ACC_Code {
						results[i].Begining = results[i].Begining + each.Ending
					}
				}
			}
		}
	}

	var Tpipes []tk.M
	var StartTrans time.Time
	var EndTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans, _ = time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTrans = EndTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearsTrans := tNow.Year()
		monthTrans := tNow.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans = tNow
	}

	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTrans, "$lt": EndTrans}}})
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
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("TransactionDeb")
				results[i].Credit = each.GetFloat64("TransactionCred")
				results[i].Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")

			}
		}
	}
	for i, _ := range results {
		results[i].Ending = results[i].Begining + results[i].Transaction
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) GetDataPasiva(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accNIN := []int{4400, 4200}
	// NINacc := make([]interface{}, len(accNIN))
	// for i, s := range accNIN {
	// 	NINacc[i] = s
	// }
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 3000), db.Lt("acc_code", 5000), db.Ne("main_acc_code", 0))).
		Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 3000), db.Lt("acc_code", 5000), db.Ne("main_acc_code", 0))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	var begin []tk.M
	begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}}})
	begin = append(begin, tk.M{"$unwind": "$ListDetail"})
	begin = append(begin, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	begin = append(begin, tk.M{"$match": tk.M{"_id": tk.M{"$nin": accNIN}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", begin).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	ResultBeginNow := make([]tk.M, 0)
	e = csr.Fetch(&ResultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// }
	if len(ResultBegin) == 0 {
		if len(ResultBegin) == 0 {
			var pipes []tk.M
			StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
			EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
			EndTime = EndTime.AddDate(0, 0, -2)
			pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
			pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":       "$ListDetail.Acc_Code",
				"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
				"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
			}})
			pipes = append(pipes, tk.M{"$match": tk.M{"_id": tk.M{"$eq": 4200}}})
			csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			defer csr.Close()
			resultJournal2 := make([]tk.M, 0)
			e = csr.Fetch(&resultJournal2, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			for i, _ := range results {
				for _, res := range resultJournal2 {
					if results[i].ACC_Code == res.GetInt("_id") {
						saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
						results[i].Transaction = results[i].Transaction + saldo
					}
				}
			}
		}
		for i, _ := range results {
			for _, res := range ResultBeginNow {
				if results[i].ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Begining = results[i].Begining + saldo
				}
			}
		}
	} else {
		if len(ResultBeginNow) == 0 {
			for i, _ := range results {
				for _, each := range ResultBegin {
					if each.ACC_Code == results[i].ACC_Code {
						results[i].Begining = results[i].Begining + each.Ending
					}
				}
			}
		} else {
			for _, each := range ResultBegin {
				for _, res := range ResultBeginNow {
					if each.ACC_Code == res.GetInt("_id") {
						saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
						each.Ending = each.Ending + saldo

					}
				}
			}
			for i, _ := range results {
				for _, each := range ResultBegin {
					if each.ACC_Code == results[i].ACC_Code {
						results[i].Begining = results[i].Begining + each.Ending
					}
				}
			}
		}
	}

	var Tpipes []tk.M
	var StartTrans time.Time
	var EndTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans, _ = time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTrans = EndTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearsTrans := tNow.Year()
		monthTrans := tNow.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans = tNow
	}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTrans, "$lt": EndTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$nin": accNIN}}})
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
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("TransactionDeb")
				results[i].Credit = each.GetFloat64("TransactionCred")
				results[i].Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
				results[i].Ending = results[i].Begining + results[i].Transaction
			}
		}
	}
	for i, _ := range results {
		if results[i].ACC_Code > 3000 && results[i].ACC_Code < 4400 {
			results[i].Transaction = results[i].Transaction * -1
		}
		results[i].Ending = results[i].Begining + results[i].Credit - results[i].Debet
	}

	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) GetDataPasivaPeriod(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		TextSearch string
		Filter     bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accNin := []int{4400, 4200}
	var filter []tk.M
	filter = append(filter, tk.M{"$match": tk.M{"category": tk.M{"$eq": "BALANCE SHEET"}}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":       "$_id",
		"acc_code":  "$acc_code",
		"acc_name":  "$account_name",
		"main_code": "$main_acc_code",
	}})
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Where(db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 3000), db.Lt("acc_code", 5000), db.Ne("main_acc_code", 0)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer crs.Close()
	resultCOA := []IncomePeriodePasiva{}
	e = crs.Fetch(&resultCOA, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	amount := []string{"$ListDetail.Debet", "$listDetail.Credit"}
	dateStart, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)

	periodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodeStart = periodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(periodeStart.Month()))
	year := strconv.Itoa(periodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	crs, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("category", "BALANCE SHEET"),
		db.Gte("acc_code", 3000), db.Lt("acc_code", 5000), db.Ne("main_acc_code", 0))).Cursor(nil)
	resultBegin := []CoaCloseModel{}
	e = crs.Fetch(&resultBegin, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	// tk.Println("Amount :", amount)
	// tk.Println("DateStart :", dateStart)
	// tk.Println("DateEnd :", End)
	type pasivaPeriodeAmount struct {
		Acc_Code  int
		Acc_Name  string
		Month     int
		Year      int
		Amount    float64
		SumDebet  float64
		SumCredit float64
	}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
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
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
		"Year":           tk.M{"$year": "$ListDetail.PostingDate"},
	}})
	// pipes = append(pipes, tk.M{"$group": tk.M{
	// 	"_id":       "$ListDetail.Acc_Code",
	// 	"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
	// 	"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	// }})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"acc_code":  "$ListDetail.Acc_Code",
		"acc_name":  "$ListDetail.Acc_Name",
		"month":     "$ListDetail.Month",
		"year":      "$ListDetail.Year",
		"amount":    "$Amount",
		"sumdebet":  "$Debet",
		"sumcredit": "$Credit",
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{
		"_id": tk.M{"$nin": accNin},
	}})
	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := []pasivaPeriodeAmount{}
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultBeginNow := make([]tk.M, 0)
	e = crs.Fetch(&resultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(resultBegin) == 0 {
		for i, _ := range resultCOA {
			for _, res := range resultBeginNow {
				if resultCOA[i].Acc_Code == res.GetInt("_id") {
					Saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + Saldo
				}
			}
		}
	} else {
		for _, each := range resultBegin {
			for _, res := range resultBeginNow {
				if each.ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					each.Ending = each.Ending + saldo
				}
			}
		}
		for i, _ := range resultCOA {
			for _, each := range resultBegin {
				if each.ACC_Code == resultCOA[i].Acc_Code {
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + each.Ending
				}
			}
		}
	}

	var Tpipes []tk.M
	var startTrans time.Time
	var endTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		startTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		endTrans = endTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearTrans := tNow.Year()
		monthTrans := tNow.Month()
		startTrans = time.Date(yearTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans = tNow
	}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": startTrans, "$lte": endTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultTransaction := make([]tk.M, 0)
	e = crs.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultTransaction {
		for i, _ := range resultCOA {
			if each.GetInt("_id") == resultCOA[i].Acc_Code {

				resultCOA[i].Amount_Debet = each.GetFloat64("TransactionDeb")
				resultCOA[i].Amount_Credit = each.GetFloat64("TransactionCred")
				resultCOA[i].Amount_Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			}
		}
	}
	for i, _ := range resultCOA {
		resultCOA[i].Amount_Ending = resultCOA[i].Amount_Begining + resultCOA[i].Amount_Transaction
	}

	// tk.Println("ResultJournal =>", resultJournal)
	for i, _ := range resultJournal {
		if resultJournal[i].Acc_Code == 5110 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code == 5100 || resultJournal[i].Acc_Code == 5200 {

			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code > 7000 && resultJournal[i].Acc_Code < 8000 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
	}
	results := struct {
		DataAcc    []IncomePeriodePasiva
		DataAmount []pasivaPeriodeAmount
	}{
		DataAcc:    resultCOA,
		DataAmount: resultJournal,
	}
	return c.SetResultInfo(false, "Success", results)

	// ==========================================================================================================
}
func (c *FinancialController) GetDataTrialIncome(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		TextSearch string
		Filter     bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(p.TextSearch)
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("category", "INCOME STATEMENT")).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)

	var begin []tk.M
	begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": PeriodeStart, "$lt": dateStart}}})
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
	ResultBeginNow := make([]tk.M, 0)
	e = csr.Fetch(&ResultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(ResultBeginNow) > 0 {
		for i, _ := range results {
			for _, res := range ResultBeginNow {
				if results[i].ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Begining = results[i].Begining + saldo
				}
			}
		}
	}

	var Tpipes []tk.M
	var StartTrans time.Time
	var EndTrans time.Time
	filter := []*db.Filter{}
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, TimeTransaction.Day(), 0, 0, 0, 0, time.UTC)
		EndTrans, _ = time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTrans = EndTrans.AddDate(0, 0, 1)
		if p.TextSearch != "" {
			filter = append(filter, db.Contains("Acc_Name", p.TextSearch))
		}
	} else {
		tNow := time.Now()
		yearsTrans := tNow.Year()
		monthTrans := tNow.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans = tNow
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa").Where(filter...).Cursor(nil)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTrans, "$lt": EndTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
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
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("TransactionDebet")
				results[i].Credit = each.GetFloat64("TransactionCredit")
				results[i].Transaction = each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
				results[i].Ending = results[i].Begining + results[i].Transaction
			}
		}
	}
	for i, _ := range results {
		results[i].Begining = 0.0
		if results[i].ACC_Code == 5110 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Credit - results[i].Debet
		}
		if results[i].ACC_Code == 5100 || results[i].ACC_Code == 5200 {
			// results[i].Debet = results[i].Debet * -1
			// results[i].Credit = results[i].Credit * -1
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Credit - results[i].Debet
		}
		if results[i].ACC_Code > 7000 && results[i].ACC_Code < 8000 {
			// results[i].Debet = results[i].Debet * -1
			// results[i].Credit = results[i].Credit * -1
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Credit - results[i].Debet
		}
	}
	// tk.Println("Data => ", results)
	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) GetDataTAX(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Or(db.Eq("account_name", "PAJAK"), db.Eq("acc_code", 6999))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	//Get Begining
	// var StartBegin time.Time
	// var EndBegin time.Time
	var Bpipes []tk.M

	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)

	Bpipes = append(Bpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": PeriodeStart, "$lt": dateStart}}})
	Bpipes = append(Bpipes, tk.M{"$unwind": "$ListDetail"})
	Bpipes = append(Bpipes, tk.M{"$group": tk.M{
		"_id":       tk.M{"Code": "$ListDetail.Acc_Code", "Name": "$ListDetail.Acc_Name"},
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Bpipes = append(Bpipes, tk.M{"$match": tk.M{"_id.Code": tk.M{"$eq": 6999}, "_id.Name": tk.M{"$eq": "PAJAK"}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Bpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resultBegining := make([]tk.M, 0)
	e = csr.Fetch(&resultBegining, 0, false)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(resultBegining) > 0 {
		for i, _ := range results {
			for _, res := range resultBegining {
				if results[i].ACC_Code == res.GetInt("_id.Code") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Begining = results[i].Begining + saldo

				}
			}
		}
	}

	var Tpipes []tk.M
	var StartTrans time.Time
	var EndTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, TimeTransaction.Day(), 0, 0, 0, 0, time.UTC)
		EndTrans, _ = time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTrans = EndTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearsTrans := tNow.Year()
		monthTrans := tNow.Month()
		StartTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		EndTrans = tNow
	}

	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTrans, "$lt": EndTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             tk.M{"Code": "$ListDetail.Acc_Code", "Name": "$ListDetail.Acc_Name"},
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id.Code": tk.M{"$ne": 6888}, "_id.Name": tk.M{"$eq": "PAJAK"}}})
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

	for _, each := range resultTransaction {
		for i, _ := range results {
			results[i].Debet = each.GetFloat64("TransactionDeb")
			results[i].Credit = each.GetFloat64("TransactionCred")
			results[i].Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) GetIncomeStatementForEarning(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	type dataTransaction struct {
		Id          bson.ObjectId
		AccountCode int
		Transaction float64
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)

	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	Penjualan := 0.0
	PotongandanRetur := 0.0
	HPP := 0.0
	for _, each := range res {
		if each.GetInt("_id") == 5110 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			Penjualan += trans
		}
		if each.GetInt("_id") == 5120 || each.GetInt("_id") == 5130 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			PotongandanRetur += trans
		}
		if each.GetInt("_id") == 5210 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			HPP += trans
		}
	}
	Penjualan = Penjualan * -1
	totalPenjualBersih := Penjualan - PotongandanRetur
	totalLabaKotor := totalPenjualBersih - HPP
	//BIAYA UMUM DAN ADMINISTRASI
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	LabaUsaha := totalLabaKotor - totOpex
	//PENDAPATAN/(BIAYA) DI LUAR USAHA
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	// BIAYA LAIN - LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	totalBiayaDanPendapatanDiluarUsaha := totLainLain - totBebanLainLain
	LabaBersihSebelumPajak := LabaUsaha + totalBiayaDanPendapatanDiluarUsaha
	//TAX
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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

	totAfterTax := LabaBersihSebelumPajak - totTax
	return c.SetResultInfo(false, "Success", totAfterTax)
}
func (c *FinancialController) getEarning(DateStart time.Time, DateEnd time.Time) ResultInfo {
	type dataTransaction struct {
		Id          bson.ObjectId
		AccountCode int
		Transaction float64
	}
	dateStart, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	Penjualan := 0.0
	PotongandanRetur := 0.0
	HPP := 0.0
	for _, each := range res {
		if each.GetInt("_id") == 5110 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			Penjualan += trans
		}
		// if each.GetInt("_id") == 5120 || each.GetInt("_id") == 5130 {
		// 	trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		// 	PotongandanRetur += trans
		// }
		if each.GetInt("_id") == 5210||each.GetInt("_id") == 5211||each.GetInt("_id") == 5212 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			HPP += trans
		}
	}
	Penjualan = Penjualan * -1
	totalPenjualBersih := Penjualan - PotongandanRetur
	totalLabaKotor := totalPenjualBersih - HPP
	//BIAYA UMUM DAN ADMINISTRASI
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	LabaUsaha := totalLabaKotor - totOpex
	//PENDAPATAN/(BIAYA) DI LUAR USAHA
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	// BIAYA LAIN - LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
	totalBiayaDanPendapatanDiluarUsaha := totLainLain - totBebanLainLain
	LabaBersihSebelumPajak := LabaUsaha + totalBiayaDanPendapatanDiluarUsaha
	//TAX
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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

	totAfterTax := LabaBersihSebelumPajak - totTax
	type modelTemp struct {
		amount float64
	}
	model := new(modelTemp)
	model.amount = totAfterTax
	return c.SetResultInfo(false, "success", totAfterTax)
}

// func (c *FinancialController) getEarning(DateStart time.Time, DateEnd time.Time) ResultInfo {
// 	type dataTransaction struct {
// 		Id          bson.ObjectId
// 		AccountCode int
// 		Transaction float64
// 	}
// 	dateStart, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
// 	dateEnd, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
// 	dateEnd = dateEnd.AddDate(0, 0, 1)

// 	var Tpipes []tk.M
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
// 	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
// 	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
// 		"_id":               "$ListDetail.Acc_Code",
// 		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
// 		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
// 	}})
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 5000, "$lt": 6000}}})
// 	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	res := make([]tk.M, 0)
// 	e = csr.Fetch(&res, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	GrossProfit := 0.0
// 	for _, each := range res {
// 		if each.GetInt("_id") != 5300 || each.GetInt("_id") != 5400 {
// 			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
// 			if each.GetInt("_id") == 5100 || each.GetInt("_id") == 5200 {
// 				trans = trans * -1
// 				GrossProfit = GrossProfit + trans
// 			}
// 			if each.GetInt("_id") == 5600 {
// 				GrossProfit = GrossProfit - trans
// 			}
// 		}
// 	}
// 	//OPERATING EXPENSE
// 	Tpipes = []tk.M{}
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
// 	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
// 	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
// 		"_id":               "$ListDetail.Acc_Code",
// 		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
// 		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
// 	}})
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 6000, "$lt": 7000}}})
// 	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	resoperex := make([]tk.M, 0)
// 	e = csr.Fetch(&resoperex, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	totOpex := 0.0
// 	for _, each := range resoperex {
// 		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
// 		totOpex = totOpex + trans
// 	}
// 	netProfit := GrossProfit - totOpex
// 	//PENGHASILAN LAIN LAIN
// 	Tpipes = []tk.M{}
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
// 	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
// 	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
// 		"_id":               "$ListDetail.Acc_Code",
// 		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
// 		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
// 	}})
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 7000, "$lt": 8000}}})
// 	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	reslainlain := make([]tk.M, 0)
// 	e = csr.Fetch(&reslainlain, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	totLainLain := 0.0
// 	for _, each := range reslainlain {
// 		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
// 		totLainLain = totLainLain + trans
// 	}
// 	totLainLain = totLainLain * -1
// 	//BEBAN LAIN - LAIN
// 	Tpipes = []tk.M{}
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
// 	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
// 	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
// 		"_id":               "$ListDetail.Acc_Code",
// 		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
// 		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
// 	}})
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 8000, "$lt": 9000}}})
// 	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	resBebanlainlain := make([]tk.M, 0)
// 	e = csr.Fetch(&resBebanlainlain, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	totBebanLainLain := 0.0
// 	for _, each := range resBebanlainlain {
// 		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
// 		totBebanLainLain = totBebanLainLain + trans
// 	}
// 	totAllRevnEx := totLainLain - totBebanLainLain
// 	totalEarningBefore := netProfit + totAllRevnEx
// 	//TAX
// 	Tpipes = []tk.M{}
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
// 	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
// 	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
// 		"_id":               tk.M{"Code": "$ListDetail.Acc_Code", "Name": "$ListDetail.Acc_Name"},
// 		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
// 		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
// 	}})
// 	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id.Code": tk.M{"$ne": 6888}, "_id.Name": tk.M{"$eq": "TAX"}}})
// 	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	resTax := make([]tk.M, 0)
// 	e = csr.Fetch(&resTax, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	totTax := 0.0
// 	for _, each := range resTax {
// 		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
// 		totTax = totTax + trans
// 	}

// 	totAfterTax := totalEarningBefore - totTax
// 	type modelTemp struct {
// 		amount float64
// 	}
// 	model := new(modelTemp)
// 	model.amount = totAfterTax
// 	return c.SetResultInfo(false, "success", totAfterTax)
// }
func (c *FinancialController) ExportToPdfIncome(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
	}{}
	type IncomeStatementModel struct {
		Id       bson.ObjectId
		Acc_Code int
		Acc_Name string
		Ending   float64
		Sales    string
	}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("category", "INCOME STATEMENT")).Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()
	resultTransaction := make([]tk.M, 0)
	e = csr.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return e.Error()
	}
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {

				results[i].Transaction = each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
				results[i].Ending = results[i].Begining + results[i].Transaction
			}
		}
	}
	for i, _ := range results {
		results[i].Begining = 0.0
		if results[i].ACC_Code == 5110 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
		if results[i].ACC_Code == 5100 || results[i].ACC_Code == 5200 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
		if results[i].ACC_Code > 7000 && results[i].ACC_Code < 8000 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
	}
	totSalesAndRevenue := 0.0
	for _, each := range results {
		if each.ACC_Code == 5110 {
			totSalesAndRevenue = totSalesAndRevenue + each.Ending
		}
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return e.Error()
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetXY(10, 5)
	pdf.SetFont("Arial", "", 12)
	// pdf.Image(c.LogoFile+"eaciit-logo.png", 12, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(11.5)
	pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(11.5)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
	pdf.SetX(145)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 15, "INCOME STATEMENT", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11.5)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
	pdf.SetX(11.5)
	// pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11.5)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
	pdf.SetX(6)

	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)

	pdf.Line(pdf.GetX()+3, pdf.GetY()+9, 198, pdf.GetY()+9)
	pdf.SetFont("Arial", "", 7)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	pdf.CellFormat(10, 10, "Date Periode : ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(4)
	pdf.GetY()
	pdf.SetX(12)
	pdf.CellFormat(10, 10, "Date Created : ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(8)
	IncomeHead := []string{"Account Code", "Account Name", "Amount", "Sales"}
	widthHead := []float64{20.0, 115.0, 30.0, 20.0}
	y0 := pdf.GetY()
	for i, head := range IncomeHead {
		pdf.SetY(y0)
		x := 13.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 5, head, "LRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 5, head, "RT", "C", false)
		}

	}
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "PENJUALAN", "LRT", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)
	//penjualan
	salesNrevenueData := []IncomeStatementModel{}
	//potongan and reture
	potonganDANretur := []IncomeStatementModel{}
	totpotonganDANretur := 0.0
	totSalespotonganDANretur := 0.0
	//hpp
	hargaPokok := []IncomeStatementModel{}
	totHargaPokok := 0.0
	totSalesHP := 0.0
	//opex
	operatingExpensesData := []IncomeStatementModel{}
	totOpEx := 0.0
	totSalesOpex := 0.0
	//other revenue
	otherRevenue := []IncomeStatementModel{}
	totOtherRev := 0.0
	totSalesOtherRev := 0.0
	//other expense
	otherExpense := []IncomeStatementModel{}
	totOtherExpense := 0.0
	totSalesOtherEx := 0.0
	//TAX
	Tax := []IncomeStatementModel{}
	totTax := 0.0
	totSalesTax := 0.0
	for _, each := range results {
		if each.ACC_Code == 5110 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
			}
			IncomeMod.Sales = Percentage
			salesNrevenueData = append(salesNrevenueData, IncomeMod)
		}
		if each.ACC_Code == 5120 || each.ACC_Code == 5130 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalespotonganDANretur = tot
			}
			IncomeMod.Sales = Percentage
			totpotonganDANretur += each.Ending
			potonganDANretur = append(potonganDANretur, IncomeMod)
		}
		if each.ACC_Code == 5210 ||each.ACC_Code == 5211||each.ACC_Code == 5212 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesHP = tot
			}
			totHargaPokok += each.Ending
			IncomeMod.Sales = Percentage
			hargaPokok = append(hargaPokok, IncomeMod)
		}
		if each.Main_Acc_Code == 6000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOpex += tot
			}
			totOpEx += each.Ending
			IncomeMod.Sales = Percentage
			operatingExpensesData = append(operatingExpensesData, IncomeMod)
		}
		if each.Main_Acc_Code == 7000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOtherRev += tot
			}
			totOtherRev += each.Ending
			IncomeMod.Sales = Percentage
			otherRevenue = append(otherRevenue, IncomeMod)
		}
		if each.Main_Acc_Code == 8000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOtherEx += tot
			}
			totOtherExpense += each.Ending
			IncomeMod.Sales = Percentage
			otherExpense = append(otherExpense, IncomeMod)
		}
		if each.Main_Acc_Code == 9000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesTax += tot
			}
			totTax += each.Ending
			IncomeMod.Sales = Percentage
			Tax = append(Tax, IncomeMod)
		}
	}

	y1 = pdf.GetY()
	for _, each := range salesNrevenueData {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
			// ending = "("
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL PENJUALAN", "1", "C", false)
	x2 := 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalSR := tk.Sprintf("%.2f", totSalesAndRevenue)
	totalSR = c.ConvertToCurrency(totalSR)
	if totSalesAndRevenue < 0 {
		totalSR = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totSalesAndRevenue*-1)) + ")"
		// "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		// totalSR = c.ConvertToCurrency(totalSR)
		// totalSR = "(" + totalSR + ")"
	}
	pdf.MultiCell(30, 5, totalSR, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	salesNrevenueSales := "100.00%"
	if totSalesAndRevenue == 0 {
		salesNrevenueSales = "0.00%"
	}
	pdf.MultiCell(20, 5, salesNrevenueSales, "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "POTONGAN DAN RETUR", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range potonganDANretur {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest > 270.0 {
			pdf.AddPage()
		}
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL POTONGAN DAN RETUR", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalpotonganFix := tk.Sprintf("%.2f", totpotonganDANretur)
	totalpotonganFix = c.ConvertToCurrency(totalpotonganFix)
	if totpotonganDANretur < 0 {
		totalpotonganFix = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totpotonganDANretur*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totalpotonganFix, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalespotonganDANretur, 'f', 2, 64)+"%", "RT", "R", false)

	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "PENJUALAN BERSIH", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totpenjualanbersih := totSalesAndRevenue - totpotonganDANretur
	// tk.Println("total", totpenjualanbersih, totSalesAndRevenue, totpotonganDANretur)
	totalpenjualanbersihFix := tk.Sprintf("%.2f", totpenjualanbersih)
	totalpenjualanbersihFix = c.ConvertToCurrency(totalpenjualanbersihFix)
	if totpenjualanbersih < 0 {
		totalpenjualanbersihFix = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totpenjualanbersih*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totalpenjualanbersihFix, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	salesPenjualanbersih := totpenjualanbersih / totSalesAndRevenue * 100
	pdf.MultiCell(20, 5, strconv.FormatFloat(salesPenjualanbersih, 'f', 2, 64)+"%", "RT", "R", false)
	//===========
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "HARGA POKOK PENJUALAN", "LRT", "L", false)
	pdf.GetY()
	y1 = pdf.GetY()
	for _, each := range hargaPokok {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
			// ending = "("
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL HARGA POKOK PENJUALAN", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totHargaPokokFix := tk.Sprintf("%.2f", totHargaPokok)
	totHargaPokokFix = c.ConvertToCurrency(totHargaPokokFix)
	if totHargaPokok < 0 {
		totHargaPokokFix = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totHargaPokok*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totHargaPokokFix, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesHP, 'f', 2, 64)+"%", "RT", "R", false)

	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "LABA KOTOR", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalLabaKotor := tk.Sprintf("%.2f", totpenjualanbersih-totHargaPokok)
	totalLabaKotor = c.ConvertToCurrency(totalLabaKotor)
	if totSalesAndRevenue < 0 {
		totalLabaKotor = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totHargaPokok*-1)) + ")"
		// "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		// totalSR = c.ConvertToCurrency(totalSR)
		// totalSR = "(" + totalSR + ")"
	}
	pdf.MultiCell(30, 5, totalLabaKotor, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	labakotorVAL := totpenjualanbersih - totHargaPokok
	saleslabakotor := labakotorVAL / totSalesAndRevenue * 100
	pdf.MultiCell(20, 5, strconv.FormatFloat(saleslabakotor, 'f', 2, 64)+"%", "RT", "R", false)
	//===========
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "BIAYA UMUM DAN ADMINISTRASI", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range operatingExpensesData {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest > 270.0 {
			pdf.AddPage()
		}
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL BIAYA UMUM DAN ADMINISTRASI", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalOpexFix := tk.Sprintf("%.2f", totOpEx)
	totalOpexFix = c.ConvertToCurrency(totalOpexFix)
	if totOpEx < 0 {
		totalOpexFix = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totOpEx*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totalOpexFix, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOpex, 'f', 2, 64)+"%", "RT", "R", false)
	// pdf.GetY()
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "LABA USAHA", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	netProfit := labakotorVAL - totOpEx
	netProfitFix := tk.Sprintf("%.2f", netProfit)
	netProfitFix = c.ConvertToCurrency(netProfitFix)
	if netProfit < 0 {
		netProfitFix = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", netProfit*-1)) + ")"
	}
	pdf.MultiCell(30, 5, netProfitFix, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	pdf.MultiCell(20, 5, strconv.FormatFloat(100-totSalesOpex, 'f', 2, 64)+"%", "RTB", "R", false)
	// pdf.AddPage()
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "PENDAPATAN/(BIAYA) DI LUAR USAHA", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range otherRevenue {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest > 270.0 {
			pdf.AddPage()
		}
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL PENDAPATAN/(BIAYA) DI LUAR USAHA", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalOR := tk.Sprintf("%.2f", totOtherRev)
	totalOR = c.ConvertToCurrency(totalOR)
	if totOtherRev < 0 {
		totalOR = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totOtherRev*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totalOR, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherRev, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "BIAYA DI LUAR USAHA", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range otherExpense {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest > 270.0 {
			pdf.AddPage()
		}
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL BIAYA DI LUAR USAHA", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totOE := tk.Sprintf("%.2f", totOtherExpense)
	totOE = c.ConvertToCurrency(totOE)
	if totOtherExpense < 0 {
		totOE = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totOtherExpense*-1)) + ")"
	}
	pdf.MultiCell(30, 5, totOE, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherEx, 'f', 2, 64)+"%", "RT", "R", false)
	// pdf.GetY()
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "TOTAL BIAYA DAN PENDAPATAN DI LUAR USAHA", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	allTotRevAndEx := totOtherRev - totOtherExpense
	allTotRevAndExFIX := tk.Sprintf("%.2f", allTotRevAndEx)
	allTotRevAndExFIX = c.ConvertToCurrency(allTotRevAndExFIX)
	if allTotRevAndEx < 0 {
		allTotRevAndExFIX = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", allTotRevAndEx*-1)) + ")"
	}
	pdf.MultiCell(30, 5, allTotRevAndExFIX, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherRev-totSalesOtherEx, 'f', 2, 64)+"%", "RTB", "R", false)
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "LABA BERSIH SEBELUM PAJAK", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	earningBeforeTAX := netProfit + allTotRevAndEx
	earningBeforeTAXFIX := tk.Sprintf("%.2f", earningBeforeTAX)
	earningBeforeTAXFIX = c.ConvertToCurrency(earningBeforeTAXFIX)
	if earningBeforeTAX < 0 {
		earningBeforeTAXFIX = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", earningBeforeTAX*-1)) + ")"
	}
	pdf.MultiCell(30, 5, earningBeforeTAXFIX, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	// saleseEarningBefore := totSalesOtherRev - totSalesOtherEx
	// salesNet := 100 - totSalesOpex
	salesearningBeforeTAX := (earningBeforeTAX / totSalesAndRevenue) * 100
	pdf.MultiCell(20, 5, strconv.FormatFloat(salesearningBeforeTAX, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "TAX", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range Tax {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		ending = c.ConvertToCurrency(ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}

		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest > 270.0 {
			pdf.AddPage()
		}
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "LABA BERSIH SETALAH PAJAK", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totEaningAfter := earningBeforeTAX - totTax
	earningAfter := tk.Sprintf("%.2f", totEaningAfter)
	earningAfter = c.ConvertToCurrency(earningAfter)
	if totEaningAfter < 0 {
		earningAfter = "(" + c.ConvertToCurrency(tk.Sprintf("%.2f", totEaningAfter*-1)) + ")"
	}
	pdf.MultiCell(30, 5, earningAfter, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totEaningAfter/totSalesAndRevenue*100, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.Ln(5)
	pdf.SetX(150)
	pdf.CellFormat(50, 6, "Approve", "", 0, "C", false, 0, "")
	pdf.Ln(18)
	pdf.GetY()
	pdf.Ln(2)
	pdf.GetY()
	pdf.SetX(150)
	pdf.CellFormat(50, 6, "__________________________", "", 0, "C", false, 0, "")

	e = os.RemoveAll(c.PdfPath + "/incomestatement")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/incomestatement", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-incomestatement.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/incomestatement"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	return fileName

}

func (c *FinancialController) GetDebet(DateStart time.Time, DateEnd time.Time) (string, []CoaModel) {

	type BalanceSheetModel struct {
		Id       bson.ObjectId
		Acc_Code int
		Acc_Name string
		Amount   float64
	}

	Account := c.AccountBL("AKTIVA")

	var AccountInter []interface{}
	for _, each := range Account {
		AccountInter = append(AccountInter, each)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.In("acc_code", AccountInter...)).Cursor(nil)

	if e != nil {

		return e.Error(), nil
	}
	defer csr.Close()
	results := []CoaModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {

		return e.Error(), nil
	}

	var pipes []tk.M
	var StartTime time.Time
	var EndTime time.Time
	// if Filter == true {
	StartTime, _ = time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	EndTimes, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	EndTime = EndTimes.AddDate(0, 0, 1)

	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
	// } else {
	// 	defaultTime := time.Now()
	// 	year := defaultTime.Year()
	// 	month := defaultTime.Month()
	// 	StartTime = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 	StartTime = StartTime.AddDate(0, -1, 0)

	// 	EndTime := defaultTime
	// 	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})

	// }
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		// return e.Error(), nil
		return e.Error(), nil
	}
	defer csr.Close()
	resultJournal := make([]tk.M, 0)
	e = csr.Fetch(&resultJournal, 0, false)
	if e != nil {
		// return e.Error(), nil
		return e.Error(), nil
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
	// newResults := []CoaModel{}
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
		// if results[i].Saldo != 0 {
		// 	newResults = append(newResults, results[i])
		// }
	}
	// dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	PeriodeStart := time.Date(StartTime.Year(), StartTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)

	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Ne("main_acc_code", 0))).Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	// for i, _ := range results {
	// 	if results[i].ACC_Code > 3000 && results[i].ACC_Code < 4400 {
	// 		results[i].Saldo = results[i].Saldo * -1
	// 	}
	// }
	newResults := []CoaModel{}
	for i, _ := range results {
		for j, _ := range ResultBegin {
			if results[i].ACC_Code == ResultBegin[j].ACC_Code {
				results[i].Saldo = results[i].Saldo + ResultBegin[j].Ending
			}
		}
		if results[i].Saldo != 0 {
			newResults = append(newResults, results[i])
		}
	}
	return "", newResults
}

func (c *FinancialController) GetCredit(DateStart time.Time, DateEnd time.Time) (string, []CoaModel) {

	type BalanceSheetModel struct {
		Id       bson.ObjectId
		Acc_Code int
		Acc_Name string
		Amount   float64
	}

	Account := c.AccountBL("PASSIVA")
	var AccountInter []interface{}
	for _, each := range Account {
		AccountInter = append(AccountInter, each)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.In("acc_code", AccountInter...)).Cursor(nil)

	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	results := []CoaModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	accNIN := []int{4400, 4200}
	var pipes []tk.M
	var StartTime time.Time
	var EndTime time.Time
	// if Filter == true {
	StartTime, _ = time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	EndTime, _ = time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)

	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
	// } else {
	// 	defaultTime := time.Now()
	// 	year := defaultTime.Year()
	// 	month := defaultTime.Month()
	// 	StartTime = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	// 	StartTime = StartTime.AddDate(0, -1, 0)
	// 	EndTime = defaultTime
	// 	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})

	// }
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"_id": tk.M{"$nin": accNIN}}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	resultJournal := make([]tk.M, 0)
	e = csr.Fetch(&resultJournal, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	dataEarning := c.getEarning(DateStart, DateEnd)
	for _, each := range resultJournal {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("SumDebet")
				results[i].Credit = each.GetFloat64("SumCredit")
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
			if results[i].ACC_Code == 4400 {
				results[i].Saldo = dataEarning.Data.(float64)
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
	PeriodeStart := time.Date(StartTime.Year(), StartTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)

	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Ne("main_acc_code", 0))).Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	if len(ResultBegin) == 0 {
		var pipes []tk.M
		// StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		// EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		EndTime = EndTime.AddDate(0, 0, -2)
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": StartTime, "$lt": EndTime}}})
		pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		pipes = append(pipes, tk.M{"$match": tk.M{"_id": tk.M{"$eq": 4200}}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return e.Error(), nil
		}
		defer csr.Close()
		resultJournal2 := make([]tk.M, 0)
		e = csr.Fetch(&resultJournal2, 0, false)
		if e != nil {
			return e.Error(), nil
		}
		for i, _ := range results {
			for _, res := range resultJournal2 {
				if results[i].ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					results[i].Saldo = results[i].Saldo + saldo
				}
			}
		}
	}
	for i, _ := range results {
		if results[i].ACC_Code > 3000 && results[i].ACC_Code < 4400 {
			results[i].Saldo = results[i].Saldo * -1
		}
	}
	newResults := []CoaModel{}
	for i, _ := range results {
		for j, _ := range ResultBegin {
			if results[i].ACC_Code == ResultBegin[j].ACC_Code {
				results[i].Saldo = results[i].Saldo + ResultBegin[j].Ending
			}
		}
		if results[i].Saldo != 0 {
			newResults = append(newResults, results[i])
		}
	}
	return "", newResults
}

func (c *FinancialController) ExportToPDFBalanceSheet(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	type BalanceSheetModel struct {
		Id       bson.ObjectId
		Acc_Code int
		Acc_Name string
		Amount   float64
	}

	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnds, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	// dateEnd := dateEnds.AddDate(0, 0, 1)
	_, datadebet := c.GetDebet(p.DateStart, p.DateEnd)
	_, datacredit := c.GetCredit(p.DateStart, p.DateEnd)

	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetXY(8, 5)

	pdf.SetFont("Arial", "", 12)
	// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(11)
	pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
	pdf.SetX(180)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 15, "NERACA", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
	pdf.SetX(10)
	// pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
	pdf.SetX(10)

	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)
	pdf.SetFont("Arial", "", 7)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+20, 286, pdf.GetY()+20) //garis horizontal2

	yAsset := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(yAsset + 10)
	pdf.MultiCell(138.5, 5, "Date Periode : "+dateStart.Format("02 January 2006")+" - "+dateEnds.Format("02 January 2006"), "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetX(9)
	pdf.SetY(yAsset + 15)
	pdf.MultiCell(138.5, 5, "Date Created : "+time.Now().Format("02 January 2006"), "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(yAsset + 21)
	pdf.SetX(10)
	pdf.MultiCell(138.5, 5, "ASSETS", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(pdf.GetY())
	pdf.SetX(37.0)
	pdf.MultiCell(138.5, 7, "AKTIVA LANCAR", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(yAsset + 21)
	pdf.SetX(148)
	pdf.MultiCell(138.5, 5, "PASSIVA", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(pdf.GetY())
	pdf.SetX(174.0)
	pdf.MultiCell(138.5, 7, "KEWAJIBAN LANCAR", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	BelanceHead := []string{"ACC NAME", "AMOUNT"}
	widthHead := []float64{99.0, 30.0}
	y0 := pdf.GetY()
	pdf.SetFont("Arial", "B", 7)
	for i, head := range BelanceHead {
		if i == 1 {
			pdf.SetY(y0)
			x := 13.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "", "R", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "", "R", false)
			}
		} else {
			pdf.SetY(y0)
			x := 13.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "", "L", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "", "L", false)
			}
		}

	}

	for i, head := range BelanceHead {
		if i == 1 {
			pdf.SetY(y0)
			x := 150.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "0", "R", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "0", "R", false)
			}
		} else {
			pdf.SetY(y0)
			x := 150.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "0", "L", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "0", "L", false)
			}
		}

	}
	pdf.SetFont("Arial", "", 7)
	datadebcred := []BalanceSheetModel{}
	data1110 := []BalanceSheetModel{}
	totalBANK := 0.0
	data1000 := []BalanceSheetModel{}
	total1000 := 0.0
	//data100
	data2000 := []BalanceSheetModel{}
	total2000 := 0.0
	//other revenu
	data3000 := []BalanceSheetModel{}
	total3000 := 0.0

	//other expense
	data4000 := []BalanceSheetModel{}
	total4000 := 0.0

	for _, each := range datadebet {
		if each.Saldo != 0 {
			if each.ACC_Code == 1110 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				data1110 = append(data1110, BalanMod)
				total1000 += BalanMod.Amount
			}
			// if each.ACC_Code == 1120 || each.ACC_Code == 1121 {
			// 	totalBANK = totalBANK + each.Saldo
			// }
			if each.ACC_Code > 1110 && each.ACC_Code < 2000 {
				if strings.Contains(each.Account_Name, "BANK") {
					totalBANK = totalBANK + each.Saldo
				}
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				if BalanMod.Acc_Code == 1401 {
					BalanMod.Acc_Name = "PERSEDIAAN"
				}
				BalanMod.Amount = each.Saldo
				data1000 = append(data1000, BalanMod)
				total1000 += BalanMod.Amount
			}
			if each.ACC_Code > 2100 && each.ACC_Code < 3000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				data2000 = append(data2000, BalanMod)

			}
			if each.ACC_Code == 2100 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}
			if each.ACC_Code == 2200 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}
			if each.ACC_Code == 2210 {
				var BalanMod BalanceSheetModel

				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)

			}
			if each.ACC_Code == 2300 {
				var BalanMod BalanceSheetModel

				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}
			if each.ACC_Code == 2310 {
				var BalanMod BalanceSheetModel

				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)

			}
			if each.ACC_Code == 2400 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				/*Percentage := "0.00%"
				totOpEx += each.Ending
				IncomeMod.Sales = Percentage*/
				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
				// datadebcred = BalanMod.Amount
			}
			if each.ACC_Code == 2410 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
			}
			if each.ACC_Code == 2500 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				/*Percentage := "0.00%"
				totOpEx += each.Ending
				IncomeMod.Sales = Percentage*/
				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
				// datadebcred = BalanMod.Amount
			}
		}
	}
	// tk.Println(total1000)
	//Data Credit
	totalPajak := 0.0
	data3110 := []BalanceSheetModel{}
	for _, each := range datacredit {
		if each.Saldo != 0 {
			if each.ACC_Code >= 3110 && each.ACC_Code < 3120 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				// if each.Saldo != 0 {
				// 	BalanMod.Amount = each.Saldo * -1
				// }
				// totOtherRev += each.Ending
				// IncomeMod.Sales = Percentage
				data3110 = append(data3110, BalanMod)
				total3000 += BalanMod.Amount
			}
			if each.ACC_Code == 3120 || each.ACC_Code == 3121 {
				totalPajak = totalPajak + each.Saldo
			}
			if each.ACC_Code > 3121 && each.ACC_Code < 4000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				// if each.Saldo != 0 {
				// 	BalanMod.Amount = each.Saldo * -1
				// }
				// totOtherRev += each.Ending
				// IncomeMod.Sales = Percentage
				data3000 = append(data3000, BalanMod)
				total3000 += BalanMod.Amount
			}
			if each.ACC_Code > 4000 && each.ACC_Code < 5000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				if BalanMod.Acc_Code == 4200 {
					// tk.Println("==== 4200==", BalanMod.Amount)
				}
				// if each.Saldo != 0 && each.ACC_Code != 4400 {
				// 	BalanMod.Amount = each.Saldo * -1
				// }

				data4000 = append(data4000, BalanMod)
				total4000 += BalanMod.Amount
			}
			// if each.ACC_Code == 2200 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo

			// 	datadebcred = append(datadebcred, BalanMod)
			// 	// datadebcred = BalanMod.Amount
			// }
			// if each.ACC_Code == 2600 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo

			// 	datadebcred = append(datadebcred, BalanMod)
			// 	total2000 += BalanMod.Amount
			// 	// datadebcred = BalanMod.Amount
			// }
			// if each.ACC_Code == 2300 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo

			// 	datadebcred = append(datadebcred, BalanMod)
			// 	// datadebcred = BalanMod.Amount
			// }

			// if each.ACC_Code == 2400 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo

			// 	datadebcred = append(datadebcred, BalanMod)
			// 	// datadebcred = BalanMod.Amount
			// }
			// if each.ACC_Code == 2700 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo

			// 	datadebcred = append(datadebcred, BalanMod)
			// 	total2000 += BalanMod.Amount

			// }

			// if each.ACC_Code == 2800 {
			// 	var BalanMod BalanceSheetModel
			// 	// IncomeMod := new(IncomeStatementModel)
			// 	BalanMod.Id = each.ID
			// 	BalanMod.Acc_Code = each.ACC_Code
			// 	BalanMod.Acc_Name = each.Account_Name
			// 	BalanMod.Amount = each.Saldo
			// 	/*Percentage := "0.00%"
			// 	totOpEx += each.Ending
			// 	IncomeMod.Sales = Percentage*/
			// 	datadebcred = append(datadebcred, BalanMod)
			// 	total2000 += BalanMod.Amount
			// 	// datadebcred = BalanMod.Amount
			// }
		}
	}

	// yfirst := pdf.GetY()
	yfix := pdf.GetY()
	yA := yfix
	pdf.SetY(yA)
	for _, each := range data1110 {
		x := 13.0
		pdf.SetX(x)
		// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		// x += widthHead[0]
		pdf.SetXY(x, yA)
		pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yA)
		amount := tk.Sprintf("%.2f", each.Amount)
		amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		yA = pdf.GetY()
	}
	x := 13.0
	pdf.SetX(x)
	// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
	// x += widthHead[0]
	pdf.SetXY(x, yA)
	// pdf.MultiCell(widthHead[0], 5, "BANK", "0", "L", false)
	// x += widthHead[0]
	// pdf.SetXY(x, yA)
	// amount := tk.Sprintf("%.2f", totalBANK)
	// amount = c.ConvertToCurrency(amount)
	// if totalBANK < 0 {
	// 	amount = tk.Sprintf("%.2f", totalBANK*-1)
	// 	amount = c.ConvertToCurrency(amount)
	// 	amount = "(" + amount + ")"
	// }
	// pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
	yA = pdf.GetY()
	// yfix = pdf.GetY()
	// yA = yfix
	pdf.SetY(yA)
	for _, each := range data1000 {
		x := 13.0
		pdf.SetX(x)
		// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		// x += widthHead[0]
		pdf.SetXY(x, yA)
		pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yA)
		amount := tk.Sprintf("%.2f", each.Amount)
		amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		yA = pdf.GetY()
	}
	// pdf.SetY(yA)
	// pdf.Ln(2)
	yA = pdf.GetY()

	yB := yfix
	pdf.SetY(yB)
	for _, each := range data3110 {
		x := 150.0
		pdf.SetX(x)
		// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		// x += widthHead[0]
		pdf.SetXY(x, yB)
		pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yB)
		amount := tk.Sprintf("%.2f", each.Amount)
		amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		yB = pdf.GetY()
	}

	if totalPajak != 0 {
		x = 150.0
		pdf.SetXY(x, yB)
		pdf.MultiCell(widthHead[0], 5, "HUTANG PAJAK", "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yB)
		amountP := tk.Sprintf("%.2f", totalPajak)
		amountP = c.ConvertToCurrency(amountP)
		if totalBANK < 0 {
			amountP = tk.Sprintf("%.2f", totalPajak*-1)
			amountP = c.ConvertToCurrency(amountP)
			amountP = "(" + amountP + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amountP, "0", "R", false)
	}

	yB = pdf.GetY()
	pdf.SetY(yB)
	for _, each := range data3000 {
		x := 150.0
		pdf.SetX(x)
		// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		// x += widthHead[0]
		pdf.SetXY(x, yB)
		pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yB)
		amount := tk.Sprintf("%.2f", each.Amount)
		amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		yB = pdf.GetY()
	}
	// pdf.SetY(yB)
	pdf.Ln(2)
	yB = pdf.GetY()

	if yA > yB {
		yfix = yA
	} else {
		yfix = yB
	}

	// ycoba := pdf.GetY()

	// pdf.Ln(8)
	yC := yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "JUMLAH AKTIVA LANCAR", "0", "L", false)
	x2 := 112.0
	pdf.SetXY(x2, yC)
	totalALL := total1000
	totalCurAss := tk.Sprintf("%.2f", totalALL)
	totalCurAss = c.ConvertToCurrency(totalCurAss)
	if totalALL < 0 {
		totalCurAss := tk.Sprintf("%.2f", totalALL*-1)
		totalCurAss = c.ConvertToCurrency(totalCurAss)
		totalCurAss = "(" + totalCurAss + ")"
	}
	pdf.MultiCell(widthHead[1], 5, totalCurAss, "", "R", false)
	pdf.Ln(2)
	yC = pdf.GetY()

	yD := yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "JUMLAH KEWAJIBAN LANCAR", "0", "L", false)
	xtocur := 249.0
	pdf.SetXY(xtocur, yD)
	total3000 = total3000 + totalPajak
	totalCurrLiab := tk.Sprintf("%.2f", total3000)
	totalCurrLiab = c.ConvertToCurrency(totalCurrLiab)
	if total3000 < 0 {
		totalCurrLiab := tk.Sprintf("%.2f", total3000*-1)
		totalCurrLiab = c.ConvertToCurrency(totalCurrLiab)
		totalCurrLiab = "(" + totalCurrLiab + ")"
	}
	pdf.MultiCell(widthHead[1], 5, totalCurrLiab, "", "R", false)
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}

	yC = yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "AKTIVA TETAP", "0", "L", false)

	//data credit
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2100 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2200 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
			// x += widthHead[2]
			// pdf.SetXY(x, yC)
			// a3 := pdf.GetY()
			// allA := []float64{a0, a1, a2, a3}

			// var n, biggest float64
			// for _, v := range allA {
			// 	if v > n {
			// 		n = v
			// 		biggest = n
			// 	}
			// }
			// pdf.SetY(biggest)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2210 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2300 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
			x += widthHead[1]
			// pdf.SetXY(x, yC)
			// pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)

		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2310 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
			// pdf.SetXY(x, yC)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2400 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2410 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2500 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			// x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		}
	}
	pdf.Ln(2)
	yC = pdf.GetY()

	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "MODAL", "0", "L", false)

	// pdf.Ln(2)
	yD = pdf.GetY()
	pdf.SetY(yD)
	for _, each := range data4000 {
		x := 150.0
		pdf.SetX(x)
		// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		// x += widthHead[0]
		pdf.SetXY(x, yD)
		pdf.MultiCell(widthHead[0], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yD)
		amount := tk.Sprintf("%.2f", each.Amount)
		amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
		yD = pdf.GetY()
	}
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}

	yC = yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "JUMLAH AKTIVA TETAP", "0", "L", false)
	x2 = 112.0
	pdf.SetXY(x2, yC)
	totalFixAss := tk.Sprintf("%.2f", total2000)
	totalFixAss = c.ConvertToCurrency(totalFixAss)
	if total2000 < 0 {
		totalFixAss := tk.Sprintf("%.2f", total2000*-1)
		totalFixAss = c.ConvertToCurrency(totalFixAss)
		totalFixAss = "(" + totalFixAss + ")"
	}
	pdf.MultiCell(30, 5, totalFixAss, "", "R", false)
	pdf.Ln(2)
	yC = pdf.GetY()

	// pdf.Ln(2)
	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "JUMLAH MODAL USAHA", "0", "L", false)
	xToCapandEar := 249.0
	pdf.SetXY(xToCapandEar, yD)
	// pdf.SetY(yD)
	totalCapEar := tk.Sprintf("%.2f", total4000)
	totalCapEar = c.ConvertToCurrency(totalCapEar)
	if total4000 < 0 {
		totalCapEar := tk.Sprintf("%.2f", total4000*-1)
		totalCapEar = c.ConvertToCurrency(totalCapEar)
		totalCapEar = "(" + totalCapEar + ")"
	}
	pdf.MultiCell(widthHead[1], 5, totalCapEar, "", "R", false)
	// tk.Println(totalCapEar)
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}
	pdf.SetXY(37.0, yfix)
	pdf.MultiCell(92.5, 7, "AKTIVA LAIN - LAIN", "0", "L", false)
	totalAKTIVAlain := 0.0
	yE := pdf.GetY()
	pdf.SetY(yE)
	for _, each := range datadebet {
		if each.Saldo != 0 {
			if each.ACC_Code >= 2600 && each.ACC_Code < 3000 {
				pdf.SetY(yE)
				x := 13.0
				pdf.SetX(x)
				// pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
				// x += widthHead[0]
				pdf.SetXY(x, yE)
				pdf.MultiCell(widthHead[0], 5, each.Account_Name, "0", "L", false)
				x += widthHead[0]
				pdf.SetXY(x, yE)
				totalAKTIVAlain = totalAKTIVAlain + each.Saldo
				amount := tk.Sprintf("%.2f", each.Saldo)
				amount = c.ConvertToCurrency(amount)
				if each.Saldo < 0 {
					amount = tk.Sprintf("%.2f", each.Saldo*-1)
					amount = c.ConvertToCurrency(amount)
					amount = "(" + amount + ")"
				}
				pdf.MultiCell(widthHead[1], 5, amount, "0", "R", false)
				yE = pdf.GetY()
			}
		}
	}
	yfix = pdf.GetY()
	yE = yfix
	pdf.SetXY(37.0, yE)
	pdf.MultiCell(92.5, 7, "TOTAL ACTIVA", "0", "L", false)
	x2 = 112.0
	pdf.SetXY(x2, yE)
	totalAsset := total1000 + total2000 + totalAKTIVAlain
	totalAssetFix := tk.Sprintf("%.2f", totalAsset)
	totalAssetFix = c.ConvertToCurrency(totalAssetFix)
	if totalAsset < 0 {
		totalAssetFix := tk.Sprintf("%.2f", totalAsset*-1)
		totalAssetFix = c.ConvertToCurrency(totalAssetFix)
		totalAssetFix = "(" + totalAssetFix + ")"
	}
	pdf.MultiCell(widthHead[1], 5, totalAssetFix, "", "R", false)

	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "TOTAL PASSIVA", "", "L", false)

	xToPass := 249.0
	pdf.SetXY(xToPass, yD)
	totalPassiva := total3000 + total4000
	totalPassivaFix := tk.Sprintf("%.2f", totalPassiva)
	totalPassivaFix = c.ConvertToCurrency(totalPassivaFix)
	if totalPassiva < 0 {
		totalPassivaFix := tk.Sprintf("%.2f", totalPassiva*-1)
		totalPassivaFix = c.ConvertToCurrency(totalPassivaFix)
		totalPassivaFix = "(" + totalPassivaFix + ")"
	}
	pdf.MultiCell(30, 5, totalPassivaFix, "", "R", false)

	y3 := pdf.GetY()
	yline := yAsset + 9
	pdf.Line(pdf.GetX(), yline, pdf.GetX(), y3)            // vertikal kiri
	pdf.Line(pdf.GetX()+138, yline+11, pdf.GetX()+138, y3) // vertikal tengah
	pdf.Line(pdf.GetX()+276, yline, pdf.GetX()+276, y3)    // vertikal kanan
	pdf.Line(pdf.GetX(), y3, 286, y3)

	e = os.RemoveAll(c.PdfPath + "/balancesheet")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/balancesheet", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}

	namepdf := "-balancesheet.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/balancesheet"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		e.Error()
	}

	return fileName
}
func (c *FinancialController) GetDataTrialIncomePeriode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		TextSearch string
		Filter     bool
		Check      bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(p.TextSearch)
	// csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("category", "INCOME STATEMENT")).
	// 	Cursor(nil)
	// if e != nil {
	// 	return c.SetResultInfo(false, e.Error(), nil)
	// }
	// defer csr.Close()
	// results := []TrialBalanceModel{}
	// e = csr.Fetch(&results, 0, false)
	// if e != nil {
	// 	return c.SetResultInfo(false, e.Error(), nil)
	// }
	var filter []tk.M
	filter = append(filter, tk.M{"$match": tk.M{"category": tk.M{"$eq": "INCOME STATEMENT"}}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":       "$_id",
		"acc_code":  "$acc_code",
		"acc_name":  "$account_name",
		"main_code": "$main_acc_code",
	}})
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Cursor(nil)
	defer crs.Close()
	resultsCOA := []IncomePeriode{}
	e = crs.Fetch(&resultsCOA, 0, false)
	if e != nil {
		c.SetResultInfo(false, e.Error(), nil)
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	dateStart, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02T15:04:05-0700", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)

	type incomePeriodeAmount struct {
		Acc_Code  int
		Acc_Name  string
		Month     int
		Year      int
		Amount    float64
		SumDebet  float64
		SumCredit float64
	}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
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
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
		"Year":           tk.M{"$year": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		// "_id":       "$ListDetail.Acc_Code",
		"_id": tk.M{
			"Acc_Code": "$Acc_Code",
			// "Acc_Name": "$Acc_Name",
			"Month": "$Month",
			"Year":  "$Year",
		},
		"Amount":    tk.M{"$sum": "$Amount"},
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
	}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"acc_code":  "$ListDetail.Acc_Code",
		"acc_name":  "$ListDetail.Acc_Name",
		"month":     "$ListDetail.Month",
		"year":      "$ListDetail.Year",
		"amount":    "$Amount",
		"sumdebet":  "$Debet",
		"sumcredit": "$Credit",
	}})
	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := []incomePeriodeAmount{}
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// accNin := []int{4400, 4200}
	// tk.Printf("%+v\n", IncomePeriodeStatement{})
	// resultsCOA = append(resultsCOA, resultCOACheck)

	periodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodeStart = periodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(periodeStart.Month()))
	year := strconv.Itoa(periodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	crs, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("category", "INCOME STATEMENT"),
		db.Gte("acc_code", 51000), db.Lt("acc_code", 9000), db.Ne("main_acc_code", 0))).Cursor(nil)
	resultBegin := []CoaCloseModel{}
	e = crs.Fetch(&resultBegin, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()

	resultBeginNow := make([]tk.M, 0)
	e = crs.Fetch(&resultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(resultBegin) == 0 {
		for i, _ := range resultsCOA {
			for _, res := range resultBeginNow {
				if resultsCOA[i].Acc_Code == res.GetInt("_id") {
					Saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					resultsCOA[i].Amount_Begining = resultsCOA[i].Amount_Begining + Saldo
				}
			}
		}
	} else {
		for _, each := range resultBegin {
			for _, res := range resultBeginNow {
				if each.ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					each.Ending = each.Ending + saldo
				}
			}
		}
		for i, _ := range resultsCOA {
			for _, each := range resultBegin {
				if each.ACC_Code == resultsCOA[i].Acc_Code {
					resultsCOA[i].Amount_Begining = resultsCOA[i].Amount_Begining + each.Ending
				}
			}
		}
	}

	var Tpipes []tk.M
	var startTrans time.Time
	var endTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		startTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		endTrans = endTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearTrans := tNow.Year()
		monthTrans := tNow.Month()
		startTrans = time.Date(yearTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans = tNow
	}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": startTrans, "$lte": endTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	resultTransaction := make([]tk.M, 0)
	e = crs.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultTransaction {
		for i, _ := range resultsCOA {
			if each.GetInt("_id") == resultsCOA[i].Acc_Code {

				resultsCOA[i].Amount_Debet = each.GetFloat64("TransactionDeb")
				resultsCOA[i].Amount_Credit = each.GetFloat64("TransactionCred")
				resultsCOA[i].Amount_Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			}
		}
	}
	for i, _ := range resultsCOA {
		resultsCOA[i].Amount_Ending = resultsCOA[i].Amount_Begining + resultsCOA[i].Amount_Transaction
	}
	// tk.Println("periode income:", resultJournal)
	for i, _ := range resultJournal {
		if resultJournal[i].Acc_Code == 5110 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code == 5100 || resultJournal[i].Acc_Code == 5200 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
		if resultJournal[i].Acc_Code > 7000 && resultJournal[i].Acc_Code < 8000 {
			resultJournal[i].Amount = resultJournal[i].Amount * -1
		}
	}
	// for i, _ := range results {
	// 	results[i].Begining = 0.0
	// 	if results[i].ACC_Code == 5100 || results[i].ACC_Code == 5200 {
	// 		results[i].Transaction = results[i].Transaction * -1
	// 		results[i].Ending = results[i].Begining + results[i].Transaction
	// 	}
	// 	if results[i].ACC_Code > 7000 && results[i].ACC_Code < 8000 {
	// 		results[i].Transaction = results[i].Transaction * -1
	// 		results[i].Ending = results[i].Begining + results[i].Transaction
	// 	}
	// }
	results := struct {
		DataAcc    []IncomePeriode
		DataAmount []incomePeriodeAmount
	}{
		DataAcc:    resultsCOA,
		DataAmount: resultJournal,
	}
	// tk.Println("Data =>", results)
	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) GetDataTAXPeriod(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  time.Time
		DateEnd    time.Time
		TextSearch string
		Filter     bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var filter []tk.M
	filter = append(filter, tk.M{"$match": tk.M{"acc_code": tk.M{"$eq": 6999}}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":       "$_id",
		"acc_code":  "$acc_code",
		"acc_name":  "$account_name",
		"main_code": "$main_acc_code",
	}})
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Cursor(nil)
	defer crs.Close()
	resultCOA := []IncomePeriode{}
	e = crs.Fetch(&resultCOA, 0, false)
	if e != nil {
		c.SetResultInfo(false, e.Error(), nil)
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)

	periodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodeStart = periodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(periodeStart.Month()))
	year := strconv.Itoa(periodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	crs, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(db.And(db.Eq("monthyear", monthYear), db.Eq("acc_code", 6999))).Cursor(nil)
	resultBegin := []CoaCloseModel{}

	type incomePeriodeAmount struct {
		Acc_Code  int
		Acc_Name  string
		Month     int
		Year      int
		Amount    float64
		SumDebet  float64
		SumCredit float64
	}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
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
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
		"Year":           tk.M{"$year": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": 6999}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		// "_id":       "$ListDetail.Acc_Code",
		"_id": tk.M{
			"Acc_Code": "$Acc_Code",
			"Acc_Name": "$Acc_Name",
			"Month":    "$Month",
			"Year":     "$Year",
		},
		"Amount":    tk.M{"$sum": "$Amount"},
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
	}})
	pipes = append(pipes, tk.M{"$unwind": "$_id"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"acc_code":  "$_id.Acc_Code",
		"acc_name":  "$_id.Acc_Name",
		"month":     "$_id.Month",
		"year":      "$_id.Year",
		"amount":    "$Amount",
		"sumdebet":  "$Debet",
		"sumcredit": "$Credit",
	}})
	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := []incomePeriodeAmount{}
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultBeginNow := make([]tk.M, 0)
	e = crs.Fetch(&resultBeginNow, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(resultBegin) == 0 {
		for i, _ := range resultCOA {
			for _, res := range resultBeginNow {
				if resultCOA[i].Acc_Code == res.GetInt("_id") {
					Saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + Saldo
				}
			}
		}
	} else {
		for _, each := range resultBegin {
			for _, res := range resultBeginNow {
				if each.ACC_Code == res.GetInt("_id") {
					saldo := res.GetFloat64("SumDebet") - res.GetFloat64("SumCredit")
					each.Ending = each.Ending + saldo
				}
			}
		}
		for i, _ := range resultCOA {
			for _, each := range resultBegin {
				if each.ACC_Code == resultCOA[i].Acc_Code {
					resultCOA[i].Amount_Begining = resultCOA[i].Amount_Begining + each.Ending
				}
			}
		}
	}

	var Tpipes []tk.M
	var startTrans time.Time
	var endTrans time.Time
	if p.Filter == true {
		TimeTransaction, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
		yearsTrans := TimeTransaction.Year()
		monthTrans := TimeTransaction.Month()
		startTrans = time.Date(yearsTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
		endTrans = endTrans.AddDate(0, 0, 1)
	} else {
		tNow := time.Now()
		yearTrans := tNow.Year()
		monthTrans := tNow.Month()
		startTrans = time.Date(yearTrans, monthTrans, 1, 0, 0, 0, 0, time.UTC)
		endTrans = tNow
	}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": startTrans, "$lte": endTrans}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":             "$ListDetail.Acc_Code",
		"TransactionDeb":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCred": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	crs, e = c.Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultTransaction := make([]tk.M, 0)
	e = crs.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultTransaction {
		for i, _ := range resultCOA {
			if each.GetInt("_id") == resultCOA[i].Acc_Code {

				resultCOA[i].Amount_Debet = each.GetFloat64("TransactionDeb")
				resultCOA[i].Amount_Credit = each.GetFloat64("TransactionCred")
				resultCOA[i].Amount_Transaction = each.GetFloat64("TransactionDeb") - each.GetFloat64("TransactionCred")
			}
		}
	}
	for i, _ := range resultCOA {
		resultCOA[i].Amount_Ending = resultCOA[i].Amount_Begining + resultCOA[i].Amount_Transaction
	}
	results := struct {
		DataAcc    []IncomePeriode
		DataAmount []incomePeriodeAmount
	}{
		DataAcc:    resultCOA,
		DataAmount: resultJournal,
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *FinancialController) ExportToExcelPeriod(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Data      []tk.M
		DataTax   []tk.M
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	//split data incomeStatement
	DataSales := []tk.M{}
	DataReture := []tk.M{}
	DataHPP := []tk.M{}
	DataOpEx := []tk.M{}
	DataOtherRev := []tk.M{}
	DataOtherExp := []tk.M{}
	DataPajak := []tk.M{}
	for _, each := range p.Data {
		if each.GetInt("Acc_Code") == 5110 {
			DataSales = append(DataSales, each)
		}
		if each.GetInt("Acc_Code") == 5120 || each.GetInt("Acc_Code") == 5130 {
			DataReture = append(DataReture, each)
		}
		if each.GetInt("Acc_Code") == 5210 || each.GetInt("Main_Acc_Code") == 5200 {
			DataHPP = append(DataHPP, each)
		}
		if each.GetInt("Main_Acc_Code") == 6000 {
			DataOpEx = append(DataOpEx, each)
		}
		if each.GetInt("Main_Acc_Code") == 7000 {
			DataOtherRev = append(DataOtherRev, each)
		}
		if each.GetInt("Main_Acc_Code") == 8000 {
			DataOtherExp = append(DataOtherExp, each)
		}
		if each.GetInt("Acc_Code") == 6999 {
			DataPajak = append(DataPajak, each)
		}
	}

	type monthyearModel struct {
		M  int
		Y  int
		MY string
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	monthyear := []monthyearModel{}
	for d := dateStart; d.Before(dateEnd); d = d.AddDate(0, 1, 0) {
		each := monthyearModel{}
		each.M = int(d.Month())
		each.Y = d.Year()
		each.MY = d.Format("Jan2006")
		monthyear = append(monthyear, each)
	}
	mergeParent := len(monthyear)
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	file = xlsx.NewFile()
	sheet, e = file.AddSheet("Sheet1")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	font := xlsx.NewFont(11, "Calibri")
	fontHead := xlsx.NewFont(11, "Calibri")
	fontHead.Bold = true
	fillHead := xlsx.NewFill("solid", "FFEB00", "FFFFFFFF")
	border := xlsx.NewBorder("thin", "thin", "thin", "thin")
	headStyle := xlsx.NewStyle()
	headStyle.Font = *fontHead
	headStyle.Fill = *fillHead
	headStyle.ApplyFill = true
	headStyle.Border = *border
	headStyle.ApplyBorder = true
	style := xlsx.NewStyle()
	style.Font = *font

	footerStyle2 := xlsx.NewStyle()
	fillFooter := xlsx.NewFill("solid", "DBDBDB", "00000000")
	rightalign := xlsx.Alignment{Horizontal: "right"}
	footerStyle2.Fill = *fillFooter
	footerStyle2.ApplyFill = true
	footerStyle2.Alignment = rightalign
	footerStyle2.Font = *font
	footerStyle2.ApplyAlignment = true

	style2 := xlsx.NewStyle()
	style2.Font = *font
	style2.Alignment = rightalign
	style2.ApplyAlignment = true

	row = sheet.AddRow()
	accName := row.AddCell()
	accName.Value = "Account Name"
	accName.SetStyle(headStyle)
	for _, each := range monthyear {
		field := row.AddCell()
		date, _ := time.Parse("Jan2006", each.MY)
		field.Value = date.Format("Jan 2006")
		field.SetStyle(headStyle)
	}
	total := row.AddCell()
	total.Value = "Total"
	total.SetStyle(headStyle)

	// Sales and revenue session

	row = sheet.AddRow()
	parent := row.AddCell()
	parent.Value = "PENJUALAN"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataSales {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalSalesAndRevenue := row.AddCell()
	totalSalesAndRevenue.Value = "TOTAL PENJUALAN"
	totalSalesAndRevenue.SetStyle(footerStyle2)
	totalSalesAndRevenueALL := 0.0
	for _, each := range monthyear {
		totalSalesAndRevenue = row.AddCell()
		GrossFloat, Grosstring := c.calculateGross(each.MY, DataSales)
		totalSalesAndRevenueALL += GrossFloat
		totalSalesAndRevenue.Value = Grosstring
		totalSalesAndRevenue.SetStyle(footerStyle2)
	}
	totalSalesAndRevenue = row.AddCell()
	totalStr := tk.Sprintf("%.2f", totalSalesAndRevenueALL)
	totalSalesAndRevenue.Value = c.ConvertToCurrency(totalStr)
	if totalSalesAndRevenueALL < 0 {
		totalStr = tk.Sprintf("%.2f", totalSalesAndRevenueALL*-1)
		totalSalesAndRevenue.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalSalesAndRevenue.SetStyle(footerStyle2)

	//RETURE
	row = sheet.AddRow()
	parent = row.AddCell()
	parent.Value = "POTONGAN DAN RETUR"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataReture {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalRetur := row.AddCell()
	totalRetur.Value = "TOTAL POTONGAN DAN RETUR"
	totalRetur.SetStyle(footerStyle2)
	totalReturAll := 0.0
	for _, each := range monthyear {
		totalRetur = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPerParent(each.MY, DataReture)
		totalReturAll += totalFloat
		totalRetur.Value = totalString
		totalRetur.SetStyle(footerStyle2)
	}
	totalRetur = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalReturAll)
	totalRetur.Value = c.ConvertToCurrency(totalStr)
	if totalReturAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalReturAll*-1)
		totalRetur.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalRetur.SetStyle(footerStyle2)
	//
	row = sheet.AddRow()
	totalPenjualanBersih := row.AddCell()
	totalPenjualanBersih.Value = "TOTAL PENJUALAN BERSIH"
	totalPenjualanBersih.SetStyle(footerStyle2)
	totalPenjualanBersihAll := 0.0
	for _, each := range monthyear {
		totalPenjualanBersih = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPenjualanBersih(each.MY, DataSales, DataReture)
		totalPenjualanBersihAll += totalFloat
		totalPenjualanBersih.Value = totalString
		totalPenjualanBersih.SetStyle(footerStyle2)
	}
	totalPenjualanBersih = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalPenjualanBersihAll)
	totalPenjualanBersih.Value = c.ConvertToCurrency(totalStr)
	if totalPenjualanBersihAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalPenjualanBersihAll*-1)
		totalPenjualanBersih.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalPenjualanBersih.SetStyle(footerStyle2)
	//HPP
	row = sheet.AddRow()
	parent = row.AddCell()
	parent.Value = "HARGA POKOK"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataHPP {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalHPP := row.AddCell()
	totalHPP.Value = "TOTAL HARGA POKOK PENJUALAN"
	totalHPP.SetStyle(footerStyle2)
	totalHPPAll := 0.0
	for _, each := range monthyear {
		totalHPP = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPerParent(each.MY, DataHPP)
		totalHPPAll += totalFloat
		totalHPP.Value = totalString
		totalHPP.SetStyle(footerStyle2)
	}
	totalHPP = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalHPPAll)
	totalHPP.Value = c.ConvertToCurrency(totalStr)
	if totalHPPAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalHPPAll*-1)
		totalHPP.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalHPP.SetStyle(footerStyle2)
	//
	row = sheet.AddRow()
	totalLabaKotor := row.AddCell()
	totalLabaKotor.Value = "LABA KOTOR"
	totalLabaKotor.SetStyle(footerStyle2)
	totalLabaKotorAll := 0.0
	for _, each := range monthyear {
		totalLabaKotor = row.AddCell()
		totalfloatPenjualanBerih, _ := c.calculateSummaryPenjualanBersih(each.MY, DataSales, DataReture)
		totalFloatHpp, _ := c.calculateSummaryPerParent(each.MY, DataHPP)
		totalFloat := totalfloatPenjualanBerih - totalFloatHpp
		totalString := tk.Sprintf("%.2f", totalFloat)
		totalString = c.ConvertToCurrency(totalStr)
		if totalFloat < 0 {
			totalString = tk.Sprintf("%.2f", totalFloat*-1)
			totalString = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		totalLabaKotorAll += totalFloat
		totalLabaKotor.Value = totalString
		totalLabaKotor.SetStyle(footerStyle2)
	}
	totalLabaKotor = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalLabaKotorAll)
	totalLabaKotor.Value = c.ConvertToCurrency(totalStr)
	if totalLabaKotorAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalLabaKotorAll*-1)
		totalLabaKotor.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalLabaKotor.SetStyle(footerStyle2)
	// Operating expenses session
	row = sheet.AddRow()
	parent = row.AddCell()
	parent.Value = "BIAYA UMUM DAN ADMINISTRASI"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataOpEx {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalOpEx := row.AddCell()
	totalOpEx.Value = "TOTAL BIAYA UMUM DAN ADMINISTRASI"
	totalOpEx.SetStyle(footerStyle2)
	totalOpExAll := 0.0
	for _, each := range monthyear {
		totalOpEx = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPerParent(each.MY, DataOpEx)
		totalOpExAll += totalFloat
		totalOpEx.Value = totalString
		totalOpEx.SetStyle(footerStyle2)
	}
	totalOpEx = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalOpExAll)
	totalOpEx.Value = c.ConvertToCurrency(totalStr)
	if totalOpExAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalOpExAll*-1)
		totalOpEx.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalOpEx.SetStyle(footerStyle2)
	//net profit
	row = sheet.AddRow()
	totalNetProfit := row.AddCell()
	totalNetProfit.Value = "LABA USAHA"
	totalNetProfit.SetStyle(footerStyle2)
	totalNetProfitALL := 0.0
	for _, each := range monthyear {
		totalNetProfit = row.AddCell()
		// totalFloat, totalString := c.calculateNetProfit(each.MY, DataSales, DataOpEx)
		totalfloatPenjualanBerih, _ := c.calculateSummaryPenjualanBersih(each.MY, DataSales, DataReture)
		totalFloatHpp, _ := c.calculateSummaryPerParent(each.MY, DataHPP)
		floatlabakotor := totalfloatPenjualanBerih - totalFloatHpp
		totalFloatOpex, _ := c.calculateSummaryPerParent(each.MY, DataOpEx)
		totalFloat := floatlabakotor - totalFloatOpex
		totalString := tk.Sprintf("%.2f", totalFloat)
		totalString = c.ConvertToCurrency(totalString)
		if totalFloat < 0 {
			totalString = tk.Sprintf("%.2f", totalFloat*-1)
			totalString = "(" + c.ConvertToCurrency(totalString) + ")"
		}
		totalNetProfitALL += totalFloat
		totalNetProfit.Value = totalString
		totalNetProfit.SetStyle(footerStyle2)
	}
	totalNetProfit = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalNetProfitALL)
	totalNetProfit.Value = c.ConvertToCurrency(totalStr)
	if totalNetProfitALL < 0 {
		totalStr = tk.Sprintf("%.2f", totalNetProfitALL*-1)
		totalNetProfit.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalNetProfit.SetStyle(footerStyle2)

	// Penghasilan lain - lain
	row = sheet.AddRow()
	parent = row.AddCell()
	parent.Value = "PENGHASILAN LAIN-LAIN"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataOtherRev {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalOtherRev := row.AddCell()
	totalOtherRev.Value = "TOTAL PENDAPATAN DI LUAR USAHA"
	totalOtherRev.SetStyle(footerStyle2)
	totalOtherRevAll := 0.0
	for _, each := range monthyear {
		totalOtherRev = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPerParent(each.MY, DataOtherRev)
		totalOtherRevAll += totalFloat
		totalOtherRev.Value = totalString
		totalOtherRev.SetStyle(footerStyle2)
	}
	totalOtherRev = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalOtherRevAll)
	totalOtherRev.Value = c.ConvertToCurrency(totalStr)
	if totalOtherRevAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalOtherRevAll*-1)
		totalOtherRev.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalOtherRev.SetStyle(footerStyle2)

	// BEBAN LAIN-LAIN
	row = sheet.AddRow()
	parent = row.AddCell()
	parent.Value = "BEBAN LAIN-LAIN"
	parent.Merge(mergeParent+1, 0)
	for _, each := range DataOtherExp {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalOtherExp := row.AddCell()
	totalOtherExp.Value = "TOTAL BIAYA DI LUAR USAHA"
	totalOtherExp.SetStyle(footerStyle2)
	totalOtherExpAll := 0.0
	for _, each := range monthyear {
		totalOtherExp = row.AddCell()
		totalFloat, totalString := c.calculateSummaryPerParent(each.MY, DataOtherExp)
		totalOtherExpAll += totalFloat
		totalOtherExp.Value = totalString
		totalOtherExp.SetStyle(footerStyle2)
	}
	totalOtherExp = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalOtherExpAll)
	totalOtherExp.Value = c.ConvertToCurrency(totalStr)
	if totalOtherExpAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalOtherExpAll*-1)
		totalOtherExp.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalOtherExp.SetStyle(footerStyle2)

	//TOTAL OTHER REVENUE AND EXPENSE
	row = sheet.AddRow()
	totalOtherRevExp := row.AddCell()
	totalOtherRevExp.Value = "TOTAL BIAYA DAN PENDAPATAN DI LUAR USAHA"
	totalOtherRevExp.SetStyle(footerStyle2)
	totalOtherRevExpALL := 0.0
	for _, each := range monthyear {
		totalOtherRevExp = row.AddCell()
		totalFloat, totalString := c.calculateRevExp(each.MY, DataOtherRev, DataOtherExp)
		totalOtherRevExpALL += totalFloat
		totalOtherRevExp.Value = totalString
		totalOtherRevExp.SetStyle(footerStyle2)
	}
	totalOtherRevExp = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalOtherRevExpALL)
	totalOtherRevExp.Value = c.ConvertToCurrency(totalStr)
	if totalOtherRevExpALL < 0 {
		totalStr = tk.Sprintf("%.2f", totalOtherRevExpALL*-1)
		totalOtherRevExp.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalOtherRevExp.SetStyle(footerStyle2)
	//EARNING BEFORE TAX
	row = sheet.AddRow()
	totalEarningBeforeTax := row.AddCell()
	totalEarningBeforeTax.Value = "LABA BERSIH SEBELUM PAJAK"
	totalEarningBeforeTax.SetStyle(footerStyle2)
	totalEarningBeforeTaxALL := 0.0
	for _, each := range monthyear {
		totalEarningBeforeTax = row.AddCell()
		// totalFloat, totalString := c.calculateEarningBeforeTax(each.MY, DataSales, DataOpEx, DataOtherRev, DataOtherExp)
		//
		totalfloatPenjualanBerih, _ := c.calculateSummaryPenjualanBersih(each.MY, DataSales, DataReture)
		totalFloatHpp, _ := c.calculateSummaryPerParent(each.MY, DataHPP)
		floatlabakotor := totalfloatPenjualanBerih - totalFloatHpp
		totalFloatOpex, _ := c.calculateSummaryPerParent(each.MY, DataOpEx)
		totalFloatLabaUsaha := floatlabakotor - totalFloatOpex
		totalFloatRevEx, _ := c.calculateRevExp(each.MY, DataOtherRev, DataOtherExp)
		totalFloat := totalFloatLabaUsaha + totalFloatRevEx
		totalString := tk.Sprintf("%.2f", totalFloat)
		totalString = c.ConvertToCurrency(totalString)
		if totalFloat < 0 {
			totalString = tk.Sprintf("%.2f", totalFloat*-1)
			totalString = "(" + c.ConvertToCurrency(totalString) + ")"
		}
		totalEarningBeforeTaxALL += totalFloat
		totalEarningBeforeTax.Value = totalString
		totalEarningBeforeTax.SetStyle(footerStyle2)
	}
	totalEarningBeforeTax = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalEarningBeforeTaxALL)
	totalEarningBeforeTax.Value = c.ConvertToCurrency(totalStr)
	if totalEarningBeforeTaxALL < 0 {
		totalStr = tk.Sprintf("%.2f", totalEarningBeforeTaxALL*-1)
		totalEarningBeforeTax.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalEarningBeforeTax.SetStyle(footerStyle2)
	// TAX
	for _, each := range p.DataTax {
		row = sheet.AddRow()
		content := row.AddCell()
		content.Value = each.GetString("Acc_Name")
		content.SetStyle(style)
		// tk.Println()
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		for _, my := range monthyear {
			content := row.AddCell()
			amount := md[my.MY].(float64)
			totalStr := tk.Sprintf("%.2f", amount)
			content.Value = c.ConvertToCurrency(totalStr)
			if amount < 0 {
				totalStr = tk.Sprintf("%.2f", amount*-1)
				content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
			}
			content.SetStyle(style2)
		}
		content = row.AddCell()
		amount := each.GetFloat64("TotalAmount")
		totalStr := tk.Sprintf("%.2f", amount)
		content.Value = c.ConvertToCurrency(totalStr)
		if amount < 0 {
			totalStr = tk.Sprintf("%.2f", amount*-1)
			content.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
		}
		content.SetStyle(style2)
	}
	row = sheet.AddRow()
	totalEarningAfterTax := row.AddCell()
	totalEarningAfterTax.Value = "LABA BERSIH SETELAH PAJAK"
	totalEarningAfterTax.SetStyle(footerStyle2)
	totalEarningAfterTaxpAll := 0.0
	for _, each := range monthyear {
		totalEarningAfterTax = row.AddCell()
		// totalFloat, totalString := c.calculateEarningAfterTax(each.MY, DataSales, DataOpEx, DataOtherRev, DataOtherExp, p.DataTax)
		totalfloatPenjualanBerih, _ := c.calculateSummaryPenjualanBersih(each.MY, DataSales, DataReture)
		totalFloatHpp, _ := c.calculateSummaryPerParent(each.MY, DataHPP)
		floatlabakotor := totalfloatPenjualanBerih - totalFloatHpp
		totalFloatOpex, _ := c.calculateSummaryPerParent(each.MY, DataOpEx)
		totalFloatLabaUsaha := floatlabakotor - totalFloatOpex
		totalFloatRevEx, _ := c.calculateRevExp(each.MY, DataOtherRev, DataOtherExp)
		totalFloatBeforeTax := totalFloatLabaUsaha + totalFloatRevEx
		totalDataTax := 0.0
		for _, each2 := range p.DataTax {
			dataItem := each2.Get("DataItem")
			md, _ := dataItem.(map[string]interface{})
			totalDataTax += md[each.MY].(float64)
		}
		totalFloat := totalFloatBeforeTax - totalDataTax
		totalString := tk.Sprintf("%.2f", totalFloat)
		totalString = c.ConvertToCurrency(totalString)
		if totalFloat < 0 {
			totalString = tk.Sprintf("%.2f", totalFloat*-1)
			totalString = "(" + c.ConvertToCurrency(totalString) + ")"
		}
		totalEarningAfterTaxpAll += totalFloat
		totalEarningAfterTax.Value = totalString
		totalEarningAfterTax.SetStyle(footerStyle2)
	}
	totalEarningAfterTax = row.AddCell()
	totalStr = tk.Sprintf("%.2f", totalEarningAfterTaxpAll)
	totalEarningAfterTax.Value = c.ConvertToCurrency(totalStr)
	if totalEarningAfterTaxpAll < 0 {
		totalStr = tk.Sprintf("%.2f", totalEarningAfterTaxpAll*-1)
		totalEarningAfterTax.Value = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	totalEarningAfterTax.SetStyle(footerStyle2)

	e = os.RemoveAll(c.UploadPath + "/report/excel")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.UploadPath+"/report/excel", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	FixName := time.Now().Format("2006-01-02T150405") + "-income.xlsx"
	fileName := FixName
	location := c.UploadPath + "/report/excel/"
	e = file.Save(location + fileName)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	return c.SetResultInfo(false, "success", fileName)
}
func (c *FinancialController) calculateGross(nameField string, data []tk.M) (float64, string) {
	DataSum := 0.0
	DataMin := 0.0
	for _, each := range data {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		if each.GetInt("Acc_Code") == 5600 {
			DataMin += md[nameField].(float64)
		} else {
			DataSum += md[nameField].(float64)
		}
	}
	Result := DataSum - DataMin
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = "(" + tk.Sprintf("%.2f", Result*-1) + ")"
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
func (c *FinancialController) calculateSummaryPerParent(nameField string, data []tk.M) (float64, string) {
	DataSum := 0.0
	for _, each := range data {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		DataSum += md[nameField].(float64)
	}
	Result := DataSum
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = "(" + tk.Sprintf("%.2f", Result*-1) + ")"
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
func (c *FinancialController) calculateSummaryPenjualanBersih(nameField string, dataSales []tk.M, dataReture []tk.M) (float64, string) {
	DataSumSales := 0.0
	for _, each := range dataSales {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		DataSumSales += md[nameField].(float64)
	}
	DataSumRetur := 0.0
	for _, each := range dataReture {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		DataSumRetur += md[nameField].(float64)
	}
	Result := DataSumSales - DataSumRetur
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = "(" + tk.Sprintf("%.2f", Result*-1) + ")"
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}

func (c *FinancialController) calculateNetProfit(nameField string, dataSales []tk.M, dataOpEx []tk.M) (float64, string) {
	totalDataSales := 0.0
	for _, each := range dataSales {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataSales += md[nameField].(float64)
	}
	totalDataOpex := 0.0
	for _, each := range dataOpEx {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataOpex += md[nameField].(float64)
	}
	Result := totalDataSales - totalDataOpex
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = tk.Sprintf("%.2f", Result*-1)
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
func (c *FinancialController) calculateRevExp(nameField string, Revenue []tk.M, Expenses []tk.M) (float64, string) {
	totalDataRev := 0.0
	for _, each := range Revenue {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataRev += md[nameField].(float64)
	}
	totalDataExp := 0.0
	for _, each := range Expenses {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataExp += md[nameField].(float64)
	}
	Result := totalDataRev - totalDataExp
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = tk.Sprintf("%.2f", Result*-1)
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
func (c *FinancialController) calculateEarningBeforeTax(nameField string, Sales []tk.M, Opex []tk.M, Revenue []tk.M, Expenses []tk.M) (float64, string) {
	totalSales := 0.0
	for _, each := range Sales {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalSales += md[nameField].(float64)
	}
	totalOpex := 0.0
	for _, each := range Opex {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalOpex += md[nameField].(float64)
	}
	netProvit := totalSales - totalOpex
	totalDataRev := 0.0
	for _, each := range Revenue {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataRev += md[nameField].(float64)
	}
	totalDataExp := 0.0
	for _, each := range Expenses {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataExp += md[nameField].(float64)
	}
	RevExp := totalDataRev - totalDataExp
	Result := netProvit + RevExp
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = tk.Sprintf("%.2f", Result*-1)
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
func (c *FinancialController) calculateEarningAfterTax(nameField string, Sales []tk.M, Opex []tk.M, Revenue []tk.M, Expenses []tk.M, Tax []tk.M) (float64, string) {
	totalSales := 0.0
	for _, each := range Sales {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalSales += md[nameField].(float64)
	}
	totalOpex := 0.0
	for _, each := range Opex {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalOpex += md[nameField].(float64)
	}
	netProvit := totalSales - totalOpex
	totalDataRev := 0.0
	for _, each := range Revenue {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataRev += md[nameField].(float64)
	}
	totalDataExp := 0.0
	for _, each := range Expenses {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataExp += md[nameField].(float64)
	}
	RevExp := totalDataRev - totalDataExp
	Before := netProvit + RevExp
	totalDataTax := 0.0
	for _, each := range Tax {
		dataItem := each.Get("DataItem")
		md, _ := dataItem.(map[string]interface{})
		totalDataTax += md[nameField].(float64)
	}
	// tk.Println(Before, totalDataTax)
	Result := Before - totalDataTax
	totalStr := tk.Sprintf("%.2f", Result)
	totalStr = c.ConvertToCurrency(totalStr)
	if Result < 0 {
		totalStr = tk.Sprintf("%.2f", Result*-1)
		totalStr = "(" + c.ConvertToCurrency(totalStr) + ")"
	}
	return Result, totalStr
}
