package controllers

import (
	"../helpers"
	. "../models"
	"os"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

type ledgerModel struct {
	Id             bson.ObjectId `json: "_id", bson: "_id"`
	No             int
	IdJournal      string
	MultiJournal   string
	Journal_Type   string
	PostingDate    time.Time
	DateStr        string
	DocumentNumber string
	Acc_Code       int
	Acc_Name       string
	Debet          float64
	Credit         float64
	Saldo          float64
	Description    string
	Department     string
	SalesCode     string
	SalesName     string
	Attachment     string
	User           string
}

func NewledgerModell() *ledgerModel {
	m := new(ledgerModel)
	m.Id = bson.NewObjectId()
	return m

}

func (c *FinancialController) DataBegining(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   time.Time
		DateEnd     time.Time
		Accountcode string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	Code, _ := strconv.Atoi(p.Accountcode)
	//Get Begining
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))

	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)
	csrr, e := c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(dbox.And(dbox.Eq("monthyear", monthYear), dbox.Eq("acc_code", Code))).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csrr.Close()
	ResultBegin := []CoaCloseModel{}
	e = csrr.Fetch(&ResultBegin, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	Begining := 0.0
	if len(ResultBegin) == 0 {
		dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))

		var pipes []tk.M
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}, "Status": tk.M{"$ne": "draft"}}})
		pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
		pipes = append(pipes, tk.M{"$project": tk.M{
			"_id":            "$ListDetail._id",
			"No":             "$ListDetail.No",
			"IdJournal":      "$IdJournal",
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
		pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": Code}})
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
		for _, each := range results {
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
		}
		var pipess []tk.M
		pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}, "Status": tk.M{"$ne": "draft"}}})
		pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
		pipess = append(pipess, tk.M{"$project": tk.M{
			"_id":            "$ListDetail._id",
			"No":             "$ListDetail.No",
			"IdJournal":      "$IdJournal",
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
		pipess = append(pipess, tk.M{"$match": tk.M{"DocumentNumber": tk.M{"$in": DocNumber}}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipess).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		results = make([]tk.M, 0)
		e = csr.Fetch(&results, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		AccountCode := []int{}
		for _, each := range results {
			skip := false
			doc := each.GetInt("Acc_Code")
			for _, num := range AccountCode {
				if doc == num {
					skip = true
					break
				}
			}
			if !skip {
				AccountCode = append(AccountCode, each.GetInt("Acc_Code"))
			}
		}
		var begin []tk.M
		begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}, "Status": tk.M{"$ne": "draft"}}})
		begin = append(begin, tk.M{"$unwind": "$ListDetail"})
		begin = append(begin, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		begin = append(begin, tk.M{"$match": tk.M{"_id": tk.M{"$in": AccountCode}}})
		csr, e = c.Ctx.Connection.NewQuery().Command("pipe", begin).From("GeneralLedger").Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		ResultBeginAll := make([]tk.M, 0)
		e = csr.Fetch(&ResultBeginAll, 0, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		for _, each := range ResultBeginAll {
			Saldo := each.GetFloat64("SumDebet") - each.GetFloat64("SumCredit")
			Begining = Begining + Saldo
		}

	} else {
		for _, each := range ResultBegin {
			Begining = each.Beginning
		}
	}
	return ResultBegining(false, "Success", Begining)
}
func ResultBegining(isError bool, msg string, data float64) ResultInfo {
	r := ResultInfo{}
	r.IsError = isError
	r.Message = msg
	r.Data = data
	return r
}

func (c *ReportController) GetDataLedger(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   time.Time
		DateEnd     time.Time
		Accountcode string
		ParentCode  int
		Filter      bool
		TextSearch  string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}

	filter := []*db.Filter{}
	filter = append(filter, db.Contains("Acc_Name", p.TextSearch))
	filter = append(filter, db.Contains("Description", p.TextSearch))
	mainFilter := new(db.Filter)
	mainFilter = db.Or(filter...)
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(mainFilter).Cursor(nil)
	// defer crs.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	Code, _ := strconv.Atoi(p.Accountcode)
	ParentCode := p.ParentCode
	var accCodeArray []int
	if ParentCode == 0 {
		prn, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Where(dbox.Eq("main_acc_code", Code)).Cursor(nil)
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
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
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
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
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

	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
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
	pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
	pipess = append(pipess, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"no":             "$ListDetail.No",
		"idjournal":      "$IdJournal",
		"multiournal":    "$Department",
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
		"salescode":     "$ListDetail.SalesCode",
		"salesname":     "$ListDetail.SalesName",
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
	pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})

	pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
	pipess = append(pipess, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"no":             "$ListDetail.No",
		"idjournal":      "$IdJournal",
		"multiournal":    "$Department",
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
		"salescode":     "$ListDetail.SalesCode",
		"salesname":     "$ListDetail.SalesName",
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
			pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})

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
				"SalesCode":     "$ListDetail.SalesCode",
				"SalesName":     "$ListDetail.SalesName",
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
			pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})

			pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
			pipess = append(pipess, tk.M{"$project": tk.M{
				"_id":            "$ListDetail._id",
				"no":             "$ListDetail.No",
				"idjournal":      "$IdJournal",
				"multijournal":   "$Department",
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
				"salescode":     "$ListDetail.SalesCode",
				"salesname":     "$ListDetail.SalesName",
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
			pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
			pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
			pipess = append(pipess, tk.M{"$project": tk.M{
				"_id":            "$ListDetail._id",
				"no":             "$ListDetail.No",
				"idjournal":      "$IdJournal",
				"multijournal":   "$Department",
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
				"salescode":     "$ListDetail.SalesCode",
				"salesname":     "$ListDetail.SalesName",
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
			for _, each := range ResultCredit {

				debetTemp := each.Debet

				if each.Credit == 0 && each.Acc_Code != newCode {
					each.Credit = debetTemp
					each.Debet = 0.0

					NewResult = append(NewResult, each)
				}
			}

		}

	} else {
		pipess = []tk.M{}
		pipess = append(pipess, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})

		pipess = append(pipess, tk.M{"$unwind": "$ListDetail"})
		pipess = append(pipess, tk.M{"$project": tk.M{
			"_id":            "$ListDetail._id",
			"no":             "$ListDetail.No",
			"idjournal":      "$IdJournal",
			"multijournal":   "$Department",
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
			"salescode":     "$ListDetail.SalesCode",
			"salesname":     "$ListDetail.SalesName",
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
		// for _, each := range ResultCredit {
		// 	debetTemp := each.Debet

		// 	if each.Credit == 0 && each.Acc_Code != Code {
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

	// BEGINING
	PeriodeStart := time.Date(dateStart.Year(), dateStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	PeriodeStart = PeriodeStart.AddDate(0, -1, 0)
	month := strconv.Itoa(int(PeriodeStart.Month()))
	year := strconv.Itoa(PeriodeStart.Year())
	my := month + year
	monthYear, _ := strconv.Atoi(my)

	csr, e = c.Ctx.Connection.NewQuery().Select().From("Coa_Close").Where(dbox.And(dbox.Eq("monthyear", monthYear), dbox.In("acc_code", Code))).Cursor(nil)
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
	if len(ResultBegin) > 0 {
		for _, each := range ResultBegin {
			Begining = Begining + each.Ending
		}
		var begin []tk.M
		begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}}})
		begin = append(begin, tk.M{"$unwind": "$ListDetail"})
		begin = append(begin, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		begin = append(begin, tk.M{"$match": tk.M{"_id": tk.M{"$eq": Code}}})
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

		for _, each := range ResultBeginNow {
			Saldo := each.GetFloat64("SumDebet") - each.GetFloat64("SumCredit")
			Begining = Begining + Saldo
		}
	} else {
		var begin []tk.M
		begin = append(begin, tk.M{"$match": tk.M{"PostingDate": tk.M{"$lt": dateStart}}})
		begin = append(begin, tk.M{"$unwind": "$ListDetail"})
		begin = append(begin, tk.M{"$group": tk.M{
			"_id":       "$ListDetail.Acc_Code",
			"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
			"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
		}})
		begin = append(begin, tk.M{"$match": tk.M{"_id": tk.M{"$eq": Code}}})
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
		for _, each := range ResultBeginNow {
			Saldo := each.GetFloat64("SumDebet") - each.GetFloat64("SumCredit")
			Begining = Begining + Saldo
		}
	}
	data := struct {
		Data     []ledgerModel
		Begining float64
	}{
		Data:     NewResult,
		Begining: Begining,
	}
	return data

}
func (c *ReportController) ExportPdfLedger(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Begining  string
		AccName   string
		AccCode   string
		DateStart time.Time
		DateEnd   time.Time
		Data      []tk.M
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
	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
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
	pdf.SetX(180)
	pdf.SetFont("Century_Gothicb", "B", 16)
	pdf.CellFormat(0, 15, "LEDGER", "", 0, "L", false, 0, "")

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

	pdf.Line(pdf.GetX()+3, pdf.GetY()+9, 280, pdf.GetY()+9)
	pdf.SetFont("Century_Gothic", "", 6)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	// if p.Filter == true {
	pdf.CellFormat(10, 10, "Date Periode  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+p.DateStart.Format("02 January 2006")+" - "+p.DateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(4)
	// }
	pdf.GetY()
	pdf.SetX(12.0)
	pdf.CellFormat(10, 10, "Account Code  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+p.AccCode, "", 0, "L", false, 0, "")
	pdf.Ln(4)
	pdf.GetY()
	pdf.SetX(12.0)
	pdf.CellFormat(10, 10, "Account Name  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+p.AccName, "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothic", "", 9)
	pdf.SetX(200)
	pdf.CellFormat(10, 10, "Begining ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+p.Begining, "", 0, "L", false, 0, "")
	pdf.Ln(8)
	pdf.SetFont("Century_Gothic", "", 6)
	coaHead := []string{"No. ", "Date", "Document Number", "Account", "Account Name", "Department", "Sales", "Description", "Debit", "Credit", "Saldo"}
	widthHead := []float64{10, 20.0, 30.0, 15.0, 30.0, 25.0, 25.0, 35.0, 27.0, 27.0, 27.0}
	y0 := pdf.GetY()
	for i, head := range coaHead {
		pdf.SetY(y0)
		x := 12.0
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
	y0 = pdf.GetY()
	pdf.SetY(pdf.GetY())
	totalDebet := 0.0
	totalCredit := 0.0
	totalSaldo := 0.0
	lastbigest := 0.0
	var length = len(DATA) + 1
	for i, each := range DATA {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 12.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "LR", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.GetString("DateStr"), "R", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[2], 5, each.GetString("DocumentNumber"), "R", "L", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, strconv.Itoa(each.GetInt("Acc_Code")), "R", "L", false)
		a3 := pdf.GetY()
		x += widthHead[3]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[4], 5, each.GetString("Acc_Name"), "R", "L", false)
		a4 := pdf.GetY()
		x += widthHead[4]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[5], 5, each.GetString("Department"), "R", "L", false)
		a5 := pdf.GetY()
		x += widthHead[5]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[6], 5, each.GetString("SalesName"), "R", "L", false)
		a6 := pdf.GetY()
		x += widthHead[6]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[7], 5, each.GetString("Description"), "R", "L", false)
		a7 := pdf.GetY()
		x += widthHead[7]
		pdf.SetXY(x, y1)
		// pdf.MultiCell(widthHead[7], 5, each.GetString("IdJournal"), "", "L", false)
		totalDebet = totalDebet + each.GetFloat64("Debet")
		Debet := tk.Sprintf("%.2f", each.GetFloat64("Debet"))
		Debet = c.ConvertToCurrency(Debet)
		if each.GetFloat64("Debet") < 0 {
			Debet = "(" + tk.Sprintf("%.2f", each.GetFloat64("Debet")*-1) + ")"
			Debet = c.ConvertToCurrency(Debet)
		}
		pdf.MultiCell(widthHead[8], 5, Debet, "R", "R", false)
		a8 := pdf.GetY()
		x += widthHead[8]
		pdf.SetXY(x, y1)
		totalCredit = totalCredit + each.GetFloat64("Credit")
		Credit := tk.Sprintf("%.2f", each.GetFloat64("Credit"))
		Credit = c.ConvertToCurrency(Credit)
		if each.GetFloat64("Credit") < 0 {
			Credit = "(" + tk.Sprintf("%.2f", each.GetFloat64("Credit")*-1) + ")"
			Credit = c.ConvertToCurrency(Credit)
		}
		pdf.MultiCell(widthHead[9], 5, Credit, "R", "R", false)
		a9 := pdf.GetY()
		x += widthHead[9]
		pdf.SetXY(x, y1)
		totalSaldo = each.GetFloat64("Saldo")
		Saldo := tk.Sprintf("%.2f", each.GetFloat64("Saldo"))
		Saldo = c.ConvertToCurrency(Saldo)
		if each.GetFloat64("Saldo") < 0 {
			Saldo = "(" + tk.Sprintf("%.2f", each.GetFloat64("Saldo")*-1) + ")"
			Saldo = c.ConvertToCurrency(Saldo)
		}
		pdf.MultiCell(widthHead[10], 5, Saldo, "R", "R", false)
		// pdf.SetXY(x, y1)
		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		if biggest >= 177.0 {
			// pdf.Line(12.0, biggest, x+widthHead[9], biggest)
			if y0 != 10.00125 {
				pdf.Line(12.0, y0, 12.0, biggest)
				pdf.Line(x+widthHead[9], y0, x+widthHead[9], biggest)
				pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest) // vertical last
				pdf.Line(12.0, biggest, x+widthHead[9], biggest)
			}
			pdf.AddPage()
			y0 = pdf.GetY()
			if y0 == 10.00125 && i != length {
				pdf.Line(12.0, y0, 12.0, biggest+5) // vertical
				pdf.Line(12.0, y0, x+widthHead[9], y0)
				pdf.Line(x+widthHead[9], y0, x+widthHead[9], biggest+5)
				pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest+5) // vertical last
				pdf.Line(12.0, biggest+5, x+widthHead[9], biggest+5)
				lastbigest = biggest + 5
			}
		}
	}
	pdf.SetFont("Century_Gothic", "", 9)
	y2 = pdf.GetY()
	// pdf.LinearGradient(x, y, w, h, r1, g1, b1, r2, g2, b2, x1, y1, x2, y2)
	pdf.LinearGradient(11.0, y2+0.2, 280, lastbigest, 255, 255, 255, 255, 255, 255, 12.0, y2, 12.0, lastbigest) //DELETE MORE LINE
	pdf.SetY(pdf.GetY())
	// pdf.Line(12.0, pdf.GetY(), 282.0, pdf.GetY())
	pdf.SetX(12.0)
	pdf.MultiCell(190, 5, "Total", "1", "C", false)
	Debet := tk.Sprintf("%.2f", totalDebet)
	Debet = c.ConvertToCurrency(Debet)
	if totalDebet < 0 {
		Debet = "(" + tk.Sprintf("%.2f", totalDebet*-1) + ")"
		Debet = c.ConvertToCurrency(Debet)
	}
	pdf.SetY(y2)
	pdf.SetX(190.0 + 12.0)
	pdf.MultiCell(27.0, 5, Debet, "TRB", "R", false)

	Credit := tk.Sprintf("%.2f", totalCredit)
	Credit = c.ConvertToCurrency(Credit)
	if totalCredit < 0 {
		Credit = "(" + tk.Sprintf("%.2f", totalCredit*-1) + ")"
		Credit = c.ConvertToCurrency(Credit)
	}
	pdf.SetY(y2)
	pdf.SetX(217.0 + 12.0)
	pdf.MultiCell(27.0, 5, Credit, "TRB", "R", false)
	pdf.SetY(y2)
	pdf.SetX(244.0 + 12.0)
	Saldo := tk.Sprintf("%.2f", totalSaldo)
	Saldo = c.ConvertToCurrency(Saldo)
	if totalSaldo < 0 {
		Saldo = "(" + tk.Sprintf("%.2f", totalSaldo*-1) + ")"
		Saldo = c.ConvertToCurrency(Saldo)
	}
	pdf.MultiCell(27.0, 5, Saldo, "TRB", "R", false)

	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-Ledger.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}
