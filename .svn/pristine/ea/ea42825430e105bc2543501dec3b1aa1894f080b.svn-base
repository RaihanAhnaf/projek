package controllers

import (
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	"time"
)

func (c *ReportController) GetDataGLpettyCash(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	End, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd := End.AddDate(0, 0, 1)
	var pipes []tk.M
	if p.Filter == true {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Status": tk.M{"$eq": "posting"}}})
	} else {
		pipes = append(pipes, tk.M{"$match": tk.M{"Status": tk.M{"$eq": "posting"}}})
	}
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
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": 1110}})
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
	return c.SetResultInfo(false, "Success", results)
}
