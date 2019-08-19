package controllers

import (
	"../helpers"
	. "../helpers"
	"../library/strformat"
	"../library/tealeg/xlsx"
	. "../models"

	"fmt"
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
)

func (c *TransactionController) GetAccount(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$sort": tk.M{"acc_code": 1}})
	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Coa").Select().Where(dbox.Ne("main_acc_code", 0)).Cursor(nil)
	defer crs.Close()
	results := make([]CoaModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
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
func (c *TransactionController) SaveJournal(k *knot.WebContext) interface{} {
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
	data := NewMainJournal()
	tk.UnjsonFromString(formData["data"][0], data)
	m := data.CreateDate.Month()
	y := data.CreateDate.Year()
	folder := fmt.Sprintf("%d%02d", y, m)
	codejurnal := fmt.Sprintf("%02d%d", m, y)

	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
	number := fmt.Sprintf("%04d", ids)

	if data.ID == "" {
		data.ID = tk.RandomString(32)
	}

	if data.IdJournal == "" {
		data.IdJournal = "JUR/" + codejurnal + "/" + number
	}
	data.DateStr = tk.Date2String(data.PostingDate, "dd MMM yyyy")

	//Create Directory
	baseImagePath := ReadConfig()["uploadpath"]
	pathfolder := filepath.Join(baseImagePath, folder)
	if _, err = os.Stat(pathfolder); os.IsNotExist(err) {
		os.MkdirAll(pathfolder, 0777)
	}

	// tk.Println(time.Now().UTC())
	for i, _ := range data.ListDetail {
		file, handler, err := k.Request.FormFile("fileUpload" + strconv.Itoa(i))
		if file != nil {
			// tk.Println(file)
			if err != nil {
				return c.ErrorResultInfo(err.Error(), nil)
			}

			defer file.Close()

			tk.RandomString(32)

			fileName := tk.RandomString(6) + filepath.Ext(handler.Filename)
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

			data.ListDetail[i].Attachment = fileName

		} else {
			data.ListDetail[i].Attachment = data.ListDetail[0].Attachment
		}
	}

	err = c.Ctx.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	c.LogActivity("Journal", "Insert Journal", data.IdJournal, k)

	return c.SetResultOK(nil)
}
func (c *TransactionController) GetDocumentNumber(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Type string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var pipes []tk.M
	now := time.Now()
	DateFilter, _ := time.Parse("2006-01-02", now.Format("2006-01-02"))
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": DateFilter}, "Journal_Type": tk.M{"$eq": p.Type}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"documentnumber": "$ListDetail.DocumentNumber",
		"date":           "$ListDetail.PostingDate",
	}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{}, "DocNo": tk.M{"$last": "$documentnumber"}}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Journal").Cursor(nil)
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
func (c *TransactionController) GetDataJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Type       string
		DateStart  string
		DateEnd    string
		Filter     bool
		TextSearch string
		IdJournal  string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	var pipes []tk.M
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)

	if p.Type == "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Journal_Type": tk.M{"$ne": p.Type}, "Status": tk.M{"$eq": "posting"}}})
	} else if p.Type == "All" {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Journal_Type": tk.M{"$ne": p.Type}, "Status": tk.M{"$eq": "posting"}}})

	} else {
		pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}, "Journal_Type": tk.M{"$eq": p.Type}, "Status": tk.M{"$eq": "posting"}}})

	}
	if p.IdJournal != "" {
		pipes = []tk.M{}
		pipes = append(pipes, tk.M{"$match": tk.M{"IdJournal": tk.M{"$eq": p.IdJournal}}})
	}
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
		"SalesName":      "$ListDetail.SalesName",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
	}})
	if p.Filter == true {
		if p.TextSearch != "" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Description": tk.M{"$regex": ".*" + p.TextSearch + ".*"}}})
		}

	}

	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Journal").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e = c.Ctx.Connection.NewQuery().From("GeneralLedger").Where(dbox.Eq("IdJournal", p.IdJournal)).Cursor(nil)
	defer csr.Close()
	general := make([]tk.M, 0)
	e = csr.Fetch(&general, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	Data := struct {
		Data1 []tk.M
		Data2 []tk.M
	}{
		Data1: results,
		Data2: general,
	}
	return c.SetResultInfo(false, "Success", Data)
}
func (c *TransactionController) GetDataDraftJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Type      string
		DateStart string
		DateEnd   string
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filters := []*dbox.Filter{}
	if p.Filter == true {
		filters = append(filters, dbox.Gte("PostingDate", dateStart))
		filters = append(filters, dbox.Lt("PostingDate", dateEnd))
	}
	filters = append(filters, dbox.Eq("Status", "draft"))
	if p.Type != "" {
		filters = append(filters, dbox.Eq("Journal_Type", p.Type))
	}

	csr, e := c.Ctx.Connection.NewQuery().From("Journal").Select().Where(dbox.And(filters...)).Cursor(nil)

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
func (c *TransactionController) PrintDraftJournalToPdf(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStr      string
		Journal_Type string
		Status       string
		User         string
		ListDetail   []Journal
		TotalDebet   string
		TotalCredit  string
		Balance      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	J_type := ""
	if p.Journal_Type == "CashIn" {
		J_type = "Cash In"
	} else if p.Journal_Type == "CashOut" {
		J_type = "Cash Out"
	} else if p.Journal_Type == "General" {
		J_type = "General Jurnal"
	}
	Status := ""
	if Status == "draft" {
		Status = "Draft"
	}
	if Status == "posting" {
		Status = "Posting"
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	JournalHeadTable := []string{"No", "Date", "Document Number", "Acc. Code", "Acc. Name", "Description", "Debet", "Credit"}
	pdf.AddPage()
	pdf.SetXY(10, 5)
	pdf.SetFont("Arial", "B", 20)
	// pdf.Image(c.LogoFile+"logo-sm.jpg", 8, 8, 0, 0, false, "", 0, "")
	// pdf.Image(c.LogoFile+"proactive.png", 157, 11, 38, 8, false, "", 0, "")
	pdf.Ln(3)
	pdf.GetY()
	pdf.SetX(10)
	pdf.CellFormat(0, 12, "Draft Journal", "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(10)
	pdf.CellFormat(0, 10, corp.Name, "", 0, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 15)
	pdf.Ln(5)
	pdf.GetY()
	pdf.SetX(10)
	pdf.WriteAligned(0, 5, "_______________________________________________________________", "L")
	pdf.SetFont("Arial", "", 7)
	pdf.Ln(5)
	pdf.GetY()
	pdf.SetX(10)
	pdf.CellFormat(30, 10, "Journal ", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, J_type, "", 0, "L", false, 0, "")
	pdf.CellFormat(20, 10, "Date", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 10, p.DateStr, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.GetY()
	pdf.SetX(10)
	pdf.CellFormat(30, 10, "Status ", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, "Draft", "", 0, "L", false, 0, "")
	pdf.CellFormat(20, 10, "User", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(35, 10, p.User, "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetX(10)
	y0 := pdf.GetY()
	widthHead := []float64{9.0, 20.0, 38.0, 20.0, 25.0, 35.0, 25.0, 25.0}
	for i, head := range JournalHeadTable {
		pdf.SetY(y0)
		x := 5.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 7, head, "LRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 7, head, "RT", "C", false)
		}

	}
	y1 := pdf.GetY()
	for _, list := range p.ListDetail {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 5.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		No := strconv.Itoa(list.No)
		pdf.MultiCell(widthHead[0], 4, No, "T", "C", false)
		a0 := pdf.GetY()
		pdf.SetY(y1)
		x += widthHead[0]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[1], 4, list.PostingDate.Local().Format("2006 Jan 02"), "T", "C", false)
		a1 := pdf.GetY()
		pdf.SetY(y1)
		x = 5.0 + widthHead[0] + widthHead[1]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[2], 4, list.DocumentNumber, "T", "C", false)
		a2 := pdf.GetY()
		pdf.SetY(y1)

		x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		acccode := strconv.Itoa(list.Acc_Code)
		pdf.MultiCell(widthHead[3], 4, acccode, "T", "C", false)
		a3 := pdf.GetY()
		pdf.SetY(y1)
		x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[4], 4, list.Acc_Name, "T", "C", false)
		a4 := pdf.GetY()
		pdf.SetY(y1)
		x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		debet := tk.Sprintf("%.2f", list.Debet)
		pdf.MultiCell(widthHead[5], 4, list.Description, "T", "C", false)
		a5 := pdf.GetY()
		pdf.SetY(y1)
		x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		Credit := tk.Sprintf("%.2f", list.Credit)
		pdf.MultiCell(widthHead[6], 4, debet, "T", "C", false)
		a6 := pdf.GetY()
		pdf.SetY(y1)
		x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[7], 4, Credit, "T", "C", false)
		a7 := pdf.GetY()
		pdf.SetY(y1)

		a8 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7, a8}

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
	pdf.SetY(y1)
	x := 5.0
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], 4, "Total", "T", "C", false)
	b1 := pdf.GetY()
	pdf.SetY(y1)
	x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[6], 4, p.TotalDebet, "T", "C", false)
	b2 := pdf.GetY()
	pdf.SetY(y1)
	x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[7], 4, p.TotalCredit, "T", "C", false)
	pdf.SetY(y1)

	allB := []float64{b1, b2}
	var n, biggest float64
	for _, v := range allB {
		if v > n {
			n = v
			biggest = n
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetY(y1)
	x = 5.0
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], 4, "Balance", "", "C", false)
	c1 := pdf.GetY()
	pdf.SetY(y1)
	x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[6], 4, "", "", "C", false)
	c2 := pdf.GetY()
	pdf.SetY(y1)
	x = 5.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
	pdf.SetLeftMargin(x)
	pdf.SetX(x)
	pdf.MultiCell(widthHead[7], 4, p.Balance, "", "C", false)
	pdf.SetY(y1)

	c3 := pdf.GetY()
	allC := []float64{c1, c2, c3}
	for _, v := range allC {
		if v > n {
			n = v
			biggest = n
		}
	}
	pdf.SetY(biggest)
	//end
	pdf.Ln(1)
	pdf.SetX(10)

	e = os.RemoveAll(c.PdfPath + "/journal")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/journal", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-journal.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + "-" + p.User + namepdf
	fileName := FixName

	location := c.PdfPath + "/journal"
	err := pdf.OutputFileAndClose(location + "/" + fileName)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", fileName)
}
func (c *TransactionController) PrintJournaltoPdf(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Type      string
		DateStart string
		DateEnd   string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filters := []*dbox.Filter{}
	filters = append(filters, dbox.Gte("PostingDate", dateStart))
	filters = append(filters, dbox.Lt("PostingDate", dateEnd))
	filters = append(filters, dbox.Eq("Status", "posting"))
	if p.Type != "" {
		filters = append(filters, dbox.Eq("Journal_Type", p.Type))
	}
	csr, e := c.Ctx.Connection.NewQuery().From("Journal").Where(dbox.And(filters...)).Cursor(nil)
	defer csr.Close()
	results := []MainJournal{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	JournalHeadTable := []string{"No", "Date", "Document Number", "Acc. Code", "Acc. Name", "Description", "Debet", "Credit"}
	for _, jr := range results {
		pdf.AddPage()
		pdf.SetXY(10, 5)
		pdf.SetFont("Arial", "", 12)
		// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
		pdf.Ln(2)

		yy := pdf.GetY()
		pdf.SetY(yy + 4)
		pdf.SetX(10)
		pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
		pdf.SetY(yy + 10)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
		pdf.SetX(152)
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 15, "POSTED JOURNAL", "", 0, "L", false, 0, "")

		pdf.SetY(yy + 17)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
		pdf.SetX(10)
		// pdf.Ln(1)

		pdf.SetY(yy + 23)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
		pdf.SetX(10)
		pdf.Line(pdf.GetX(), pdf.GetY()+9, 200, pdf.GetY()+9) //garis horizontal1

		pdf.SetFont("Arial", "B", 14)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)

		pdf.SetFont("Arial", "", 7)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Journal ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Journal_Type, "", 0, "L", false, 0, "")

		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Status ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Status, "", 0, "L", false, 0, "")

		pdf.Ln(10)

		y0 := pdf.GetY()
		widthHead := []float64{9.0, 20.0, 31.0, 20.0, 25.0, 35.0, 25.0, 25.0}
		for i, head := range JournalHeadTable {
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
		y1 := pdf.GetY()
		totalDebet := 0.0
		totalCredit := 0.0
		for _, list := range jr.ListDetail {
			totalDebet += list.Debet
			totalCredit += list.Credit
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 10.0
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			No := strconv.Itoa(list.No)
			pdf.MultiCell(widthHead[0], 4, No, "", "C", false)
			a0 := pdf.GetY()
			pdf.SetY(y1)
			x += widthHead[0]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[1], 4, list.PostingDate.Local().Format("2006 Jan 02"), "", "L", false)
			a1 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[2], 4, list.DocumentNumber, "", "L", false)
			a2 := pdf.GetY()
			pdf.SetY(y1)

			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			acccode := strconv.Itoa(list.Acc_Code)
			pdf.MultiCell(widthHead[3], 4, acccode, "", "L", false)
			a3 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[4], 4, list.Acc_Name, "", "L", false)
			a4 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			debet := tk.Sprintf("%.2f", list.Debet)
			debet = c.ConvertToCurrency(debet)
			pdf.MultiCell(widthHead[5], 4, list.Description, "", "L", false)
			a5 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			Credit := tk.Sprintf("%.2f", list.Credit)
			Credit = c.ConvertToCurrency(Credit)
			pdf.MultiCell(widthHead[6], 4, debet, "", "R", false)
			a6 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[7], 4, Credit, "", "R", false)
			a7 := pdf.GetY()
			pdf.SetY(y1)

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

		y3 := pdf.GetY()
		pdf.Line(pdf.GetX()-165, y0, pdf.GetX()-165, y3) //garis vertical 1
		pdf.Line(pdf.GetX()-156, y0, pdf.GetX()-156, y3) //garis vertical 2
		pdf.Line(pdf.GetX()-136, y0, pdf.GetX()-136, y3) //garis vertical 3
		pdf.Line(pdf.GetX()-105, y0, pdf.GetX()-105, y3) //garis vertical 4
		pdf.Line(pdf.GetX()-85, y0, pdf.GetX()-85, y3)   //garis vertical 5
		pdf.Line(pdf.GetX()-60, y0, pdf.GetX()-60, y3)   //garis vertical 6
		pdf.Line(pdf.GetX()-25, y0, pdf.GetX()-25, y3)   //garis vertical 7
		pdf.Line(pdf.GetX(), y0, pdf.GetX(), y3)         //garis vertical 8
		pdf.Line(pdf.GetX()+25, y0, pdf.GetX()+25, y3)   //garis vertical 9
		pdf.Line(pdf.GetX()-165, y3, 200, y3)            //garis horizontal

		Balance := totalDebet - totalCredit
		y3 = pdf.GetY()
		pdf.SetY(y3)
		pdf.SetX(10 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Total", "LBR", "C", false)
		pdf.SetY(pdf.GetY())

		x := 10.0

		pdf.SetX(x)

		b1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totDeb := tk.Sprintf("%.2f", totalDebet)
		totDeb = c.ConvertToCurrency(totDeb)
		pdf.MultiCell(widthHead[6], 4, totDeb, "RB", "R", false)
		b2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totCred := tk.Sprintf("%.2f", totalCredit)
		totCred = c.ConvertToCurrency(totCred)
		pdf.MultiCell(widthHead[7], 4, totCred, "RB", "R", false)
		pdf.SetY(y3)

		b3 := pdf.GetY()
		allB := []float64{b1, b2, b3}
		var n, biggest float64
		for _, v := range allB {
			if v > n {
				n = v
				biggest = n
			}
			pdf.SetY(biggest)
		}

		y3 = pdf.GetY()
		pdf.SetY(y3)
		x = 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Balance", "LBR", "C", false)

		c1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[6], 4, "", "BR", "R", false)
		c2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		fixBalance := tk.Sprintf("%.2f", Balance)
		pdf.MultiCell(widthHead[6], 4, fixBalance, "RB", "R", false)
		pdf.SetY(y1)

		c3 := pdf.GetY()
		allC := []float64{c1, c2, c3}
		for _, v := range allC {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		fullName := ""
		if jr.User == k.Session("username").(string) {
			fullName = k.Session("fullname").(string)
		} else {
			csr, e := c.Ctx.Connection.NewQuery().Select().From(NewSysUserModel().TableName()).Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			user := SysUserModel{}
			e = csr.Fetch(&user, 1, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			csr.Close()
			fullName = user.Fullname
		}
		pdf.Ln(5)
		pdf.SetX(5)
		pdf.CellFormat(50, 6, "User", "", 0, "C", false, 0, "")
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "Approve", "", 0, "C", false, 0, "")
		pdf.Ln(18)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(40, 6, fullName, "B", 0, "C", false, 0, "")
		pdf.Ln(2)
		pdf.GetY()
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "__________________________", "", 0, "C", false, 0, "")
		//end
		pdf.Ln(1)
		pdf.SetX(10)
	}
	e = os.RemoveAll(c.PdfPath + "/journal")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/journal", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-journal.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + "-" + "posting" + namepdf
	fileName := FixName
	location := c.PdfPath + "/journal"
	err := pdf.OutputFileAndClose(location + "/" + fileName)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", fileName)
}
func (c *TransactionController) SavePosting(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	// fmt.Println(payload.GetString("Id"))
	ids := payload.GetString("Id")
	roles := payload.GetString("Role")

	ctx := c.Ctx.Connection
	data := make([]MainJournal, 0)

	crsData, errData := ctx.NewQuery().From("Journal").Select().Where(dbox.Eq("_id", ids)).Cursor(nil)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	errData = crsData.Fetch(&data, 0, false)
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	time.Now().UTC()
	var m time.Month
	var y int

	ModelJournal := NewMainJournal()
	ModelGeneral := NewGeneralLedger()
	for _, dt := range data {
		NewListDetail := []Journal{}
		ModelJournal.ID = dt.ID
		ModelJournal.IdJournal = dt.IdJournal
		ModelJournal.CreateDate = dt.CreateDate

		if roles == "supervisor" {
			ModelJournal.PostingDate = dt.PostingDate
			ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else if roles == "administrator" {
			ModelJournal.PostingDate = dt.PostingDate
			ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else {
			ModelJournal.PostingDate = time.Now().UTC()
			ModelJournal.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
		}
		ModelJournal.User = dt.User
		ModelJournal.Journal_Type = dt.Journal_Type
		ModelJournal.Department = dt.Department
		ModelJournal.SalesCode = dt.SalesCode
		ModelJournal.SalesName = dt.SalesName
		ModelJournal.Status = "posting"
		numbertemporary := 0
		if ModelJournal.Department == "MULTIJOURNAL" {
			valDebet := 0.0
			valCredit := 0.0
			for i, arrList := range dt.ListDetail {
				headcode := arrList.DocumentNumber[:11]
				// if roles == "supervisor" {
				m = dt.PostingDate.Month()
				y = dt.PostingDate.Year()
				// } else {
				// 	m = time.Now().UTC().Month()
				// 	y = time.Now().UTC().Year()
				// }
				balance := valDebet - valCredit
				valDebet += arrList.Debet
				valCredit += arrList.Credit
				if i == 0 {
					idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
					number := fmt.Sprintf("%04d", idx)
					arrList.DocumentNumber = headcode + number
					numbertemporary = idx
				} else {
					// tk.Println(i, arrList.DocumentNumber, valDebet, valCredit)
					if balance == 0 {
						idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						arrList.DocumentNumber = headcode + number
						numbertemporary = idx
						// valDebet = 0
						// valCredit = 0
					} else {
						number := fmt.Sprintf("%04d", numbertemporary)
						arrList.DocumentNumber = headcode + number
					}
				}
				// tk.Println(arrList.DocumentNumber)
				arrList.PostingDate = ModelJournal.PostingDate
				arrList.DateStr = tk.Date2String(ModelJournal.PostingDate, "dd MMM yyyy")

				NewListDetail = append(NewListDetail, arrList)
			}
		} else {
			for i, arrList := range dt.ListDetail {
				headcode := arrList.DocumentNumber[:11]
				// tk.Println(i)
				if i%2 == 0 {
					// if roles == "supervisor" {
					m = dt.PostingDate.Month()
					y = dt.PostingDate.Year()
					// } else {
					// 	m = time.Now().UTC().Month()
					// 	y = time.Now().UTC().Year()
					// }
					idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
					number := fmt.Sprintf("%04d", idx)
					arrList.DocumentNumber = headcode + number
					numbertemporary = idx
				} else {
					number := fmt.Sprintf("%04d", numbertemporary)
					arrList.DocumentNumber = headcode + number
				}

				arrList.PostingDate = ModelJournal.PostingDate
				arrList.DateStr = tk.Date2String(ModelJournal.PostingDate, "dd MMM yyyy")

				NewListDetail = append(NewListDetail, arrList)
			}
		}
		ModelJournal.ListDetail = NewListDetail
		c.Ctx.Save(ModelJournal)

		//==================General Ledger Save=======================
		NewListDetailGeneral := []GeneralDetail{}
		ModelGeneral.ID = dt.ID
		ModelGeneral.IdJournal = dt.IdJournal
		ModelGeneral.CreateDate = dt.CreateDate
		if roles == "supervisor" {
			ModelGeneral.PostingDate = dt.PostingDate
			ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else if roles == "administrator" {
			ModelGeneral.PostingDate = dt.PostingDate
			ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else {
			ModelGeneral.PostingDate = time.Now().UTC()
			ModelGeneral.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
		}
		ModelGeneral.User = dt.User
		ModelGeneral.Journal_Type = dt.Journal_Type
		ModelGeneral.Status = "posting"
		ModelGeneral.Department = dt.Department
		ModelGeneral.SalesCode = dt.SalesCode
		ModelGeneral.SalesName = dt.SalesName
		numbertemporarytwo := 0
		if ModelJournal.Department == "MULTIJOURNAL" {
			valDebet := 0.0
			valCredit := 0.0
			for u, arrListGeneral := range dt.ListDetail {
				datalist := GeneralDetail{}
				headcode := arrListGeneral.DocumentNumber[:11]
				// if roles == "supervisor" {
				m = dt.PostingDate.Month()
				y = dt.PostingDate.Year()
				// } else {
				// 	m = time.Now().UTC().Month()
				// 	y = time.Now().UTC().Year()
				// }

				balance := valDebet - valCredit
				valDebet += arrListGeneral.Debet
				valCredit += arrListGeneral.Credit
				if u == 0 {
					idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
					number := fmt.Sprintf("%04d", idx)
					datalist.DocumentNumber = headcode + number
					numbertemporarytwo = idx
				} else {
					if balance == 0 {
						idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						datalist.DocumentNumber = headcode + number
						numbertemporarytwo = idx
						// valDebet = 0
						// valCredit = 0
					} else {
						number := fmt.Sprintf("%04d", numbertemporarytwo)
						datalist.DocumentNumber = headcode + number
					}
				}

				datalist.PostingDate = ModelGeneral.PostingDate
				datalist.DateStr = tk.Date2String(ModelGeneral.PostingDate, "dd MMM yyyy")

				datalist.Id = arrListGeneral.Id
				datalist.No = arrListGeneral.No
				datalist.Journal_Type = arrListGeneral.Journal_Type
				datalist.Acc_Code = arrListGeneral.Acc_Code
				datalist.Acc_Name = arrListGeneral.Acc_Name
				datalist.Debet = arrListGeneral.Debet
				datalist.Credit = arrListGeneral.Credit
				datalist.Description = arrListGeneral.Description
				datalist.Attachment = arrListGeneral.Attachment
				datalist.User = arrListGeneral.User
				datalist.Department = arrListGeneral.Department
				datalist.SalesCode = arrListGeneral.SalesCode
				datalist.SalesName = arrListGeneral.SalesName

				NewListDetailGeneral = append(NewListDetailGeneral, datalist)
			}
		} else {
			for u, arrListGeneral := range dt.ListDetail {
				datalist := GeneralDetail{}
				headcode := arrListGeneral.DocumentNumber[:11]

				if u%2 == 0 {
					// if roles == "supervisor" {
					m = dt.PostingDate.Month()
					y = dt.PostingDate.Year()
					// } else {
					// 	m = time.Now().UTC().Month()
					// 	y = time.Now().UTC().Year()
					// }
					idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
					number := fmt.Sprintf("%04d", idx)
					datalist.DocumentNumber = headcode + number
					numbertemporarytwo = idx
				} else {
					number := fmt.Sprintf("%04d", numbertemporarytwo)
					datalist.DocumentNumber = headcode + number
				}

				datalist.PostingDate = ModelGeneral.PostingDate
				datalist.DateStr = tk.Date2String(ModelGeneral.PostingDate, "dd MMM yyyy")

				datalist.Id = arrListGeneral.Id
				datalist.No = arrListGeneral.No
				datalist.Journal_Type = arrListGeneral.Journal_Type
				datalist.Acc_Code = arrListGeneral.Acc_Code
				datalist.Acc_Name = arrListGeneral.Acc_Name
				datalist.Debet = arrListGeneral.Debet
				datalist.Credit = arrListGeneral.Credit
				datalist.Description = arrListGeneral.Description
				datalist.Attachment = arrListGeneral.Attachment
				datalist.User = arrListGeneral.User
				datalist.Department = arrListGeneral.Department
				datalist.SalesCode = arrListGeneral.SalesCode
				datalist.SalesName = arrListGeneral.SalesName

				NewListDetailGeneral = append(NewListDetailGeneral, datalist)
			}
		}
		ModelGeneral.ListDetail = NewListDetailGeneral
		c.Ctx.Save(ModelGeneral)
		c.LogActivity("Journal", "Save Posting Journal", dt.IdJournal, k)
	}

	return c.SetResultOK(nil)
}
func (c *TransactionController) DeleteDraftJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	result := new(MainJournal)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("Journal", "Delete Journal", result.IdJournal, k)

	return c.SetResultInfo(false, "", nil)
}
func (c *TransactionController) SavePostingAndPrint(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	// fmt.Println(payload.GetString("Id"))
	ids := payload.GetString("Id")
	roles := payload.GetString("Role")

	ctx := c.Ctx.Connection
	data := make([]MainJournal, 0)

	crsData, errData := ctx.NewQuery().From("Journal").Select().Where(dbox.Eq("_id", ids)).Cursor(nil)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	errData = crsData.Fetch(&data, 0, false)
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	time.Now().UTC()
	var m time.Month
	var y int

	ModelJournal := NewMainJournal()
	ModelGeneral := NewGeneralLedger()
	dataPrint := []MainJournal{}
	for _, dt := range data {
		NewListDetail := []Journal{}
		ModelJournal.ID = dt.ID
		ModelJournal.IdJournal = dt.IdJournal
		ModelJournal.CreateDate = dt.CreateDate

		if roles == "supervisor" {
			ModelJournal.PostingDate = dt.PostingDate
			ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else if roles == "administrator" {
			ModelJournal.PostingDate = dt.PostingDate
			ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else {
			ModelJournal.PostingDate = time.Now().UTC()
			ModelJournal.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
		}
		ModelJournal.User = dt.User
		ModelJournal.Journal_Type = dt.Journal_Type
		ModelJournal.Status = "posting"
		ModelJournal.Department = dt.Department
		ModelJournal.SalesCode = dt.SalesCode
		ModelJournal.SalesName = dt.SalesName
		numbertemporary := 0
		for i, arrList := range dt.ListDetail {
			headcode := arrList.DocumentNumber[:11]
			// tk.Println(i)
			if i%2 == 0 {
				// if roles == "supervisor" {
				m = dt.PostingDate.Month()
				y = dt.PostingDate.Year()
				// } else {
				// 	m = time.Now().UTC().Month()
				// 	y = time.Now().UTC().Year()
				// }
				idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
				number := fmt.Sprintf("%04d", idx)
				arrList.DocumentNumber = headcode + number
				numbertemporary = idx
			} else {
				number := fmt.Sprintf("%04d", numbertemporary)
				arrList.DocumentNumber = headcode + number
			}

			arrList.PostingDate = ModelJournal.PostingDate
			arrList.DateStr = tk.Date2String(ModelJournal.PostingDate, "dd MMM yyyy")

			NewListDetail = append(NewListDetail, arrList)
		}
		ModelJournal.ListDetail = NewListDetail
		c.Ctx.Save(ModelJournal)
		postPrint := MainJournal{}
		postPrint.ID = ModelJournal.ID
		postPrint.IdJournal = ModelJournal.IdJournal
		postPrint.CreateDate = ModelJournal.CreateDate
		postPrint.PostingDate = ModelJournal.PostingDate
		postPrint.DateStr = ModelJournal.DateStr
		postPrint.User = ModelJournal.User
		postPrint.Journal_Type = ModelJournal.Journal_Type
		postPrint.Status = ModelJournal.Status
		postPrint.ListDetail = ModelJournal.ListDetail
		dataPrint = append(dataPrint, postPrint)
		//==================General Ledger Save=======================
		NewListDetailGeneral := []GeneralDetail{}
		ModelGeneral.ID = dt.ID
		ModelGeneral.IdJournal = dt.IdJournal
		ModelGeneral.CreateDate = dt.CreateDate
		if roles == "supervisor" {
			ModelGeneral.PostingDate = dt.PostingDate
			ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else if roles == "administrator" {
			ModelGeneral.PostingDate = dt.PostingDate
			ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
		} else {
			ModelGeneral.PostingDate = time.Now().UTC()
			ModelGeneral.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
		}
		ModelGeneral.User = dt.User
		ModelGeneral.Journal_Type = dt.Journal_Type
		ModelGeneral.Department = dt.Department
		ModelGeneral.SalesCode = dt.SalesCode
		ModelGeneral.SalesName = dt.SalesName
		ModelGeneral.Status = "posting"
		numbertemporarytwo := 0
		for u, arrListGeneral := range dt.ListDetail {
			datalist := GeneralDetail{}
			headcode := arrListGeneral.DocumentNumber[:11]

			if u%2 == 0 {
				// if roles == "supervisor" {
				m = dt.PostingDate.Month()
				y = dt.PostingDate.Year()
				// } else {
				// 	m = time.Now().UTC().Month()
				// 	y = time.Now().UTC().Year()
				// }
				idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
				number := fmt.Sprintf("%04d", idx)
				datalist.DocumentNumber = headcode + number
				numbertemporarytwo = idx
			} else {
				number := fmt.Sprintf("%04d", numbertemporarytwo)
				datalist.DocumentNumber = headcode + number
			}

			datalist.PostingDate = ModelGeneral.PostingDate
			datalist.DateStr = tk.Date2String(ModelGeneral.PostingDate, "dd MMM yyyy")

			datalist.Id = arrListGeneral.Id
			datalist.No = arrListGeneral.No
			datalist.Journal_Type = arrListGeneral.Journal_Type
			datalist.Acc_Code = arrListGeneral.Acc_Code
			datalist.Acc_Name = arrListGeneral.Acc_Name
			datalist.Debet = arrListGeneral.Debet
			datalist.Credit = arrListGeneral.Credit
			datalist.Description = arrListGeneral.Description
			datalist.Attachment = arrListGeneral.Attachment
			datalist.User = arrListGeneral.User
			datalist.Department = arrListGeneral.Department
			datalist.SalesCode = arrListGeneral.SalesCode
			datalist.SalesName = arrListGeneral.SalesName

			NewListDetailGeneral = append(NewListDetailGeneral, datalist)
		}
		ModelGeneral.ListDetail = NewListDetailGeneral
		c.Ctx.Save(ModelGeneral)
		c.LogActivity("Journal", "Save Posting Journal", dt.IdJournal, k)
	}

	export := c.ExportPdfAfterPosting(dataPrint, k)

	return c.SetResultInfo(false, "OK", export)
}
func (c *TransactionController) ExportPdfAfterPosting(results []MainJournal, k *knot.WebContext) string {
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return e.Error()
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	JournalHeadTable := []string{"No", "Date", "Document Number", "Acc. Code", "Acc. Name", "Description", "Debet", "Credit"}
	for _, jr := range results {
		pdf.AddPage()
		pdf.SetXY(10, 5)
		pdf.SetFont("Arial", "", 12)
		// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
		pdf.Ln(2)

		yy := pdf.GetY()
		pdf.SetY(yy + 4)
		pdf.SetX(10)
		pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
		pdf.SetY(yy + 10)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
		pdf.SetX(152)
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 15, "POSTED JOURNAL", "", 0, "L", false, 0, "")

		pdf.SetY(yy + 17)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
		pdf.SetX(10)

		pdf.SetY(yy + 23)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
		pdf.SetX(10)
		pdf.Line(pdf.GetX(), pdf.GetY()+9, 200, pdf.GetY()+9) //garis horizontal1

		pdf.SetFont("Arial", "B", 14)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)

		pdf.SetFont("Arial", "", 7)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Journal ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Journal_Type, "", 0, "L", false, 0, "")

		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Status ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Status, "", 0, "L", false, 0, "")

		pdf.Ln(10)

		y0 := pdf.GetY()
		widthHead := []float64{9.0, 20.0, 31.0, 20.0, 25.0, 35.0, 25.0, 25.0}
		for i, head := range JournalHeadTable {
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
		y1 := pdf.GetY()
		totalDebet := 0.0
		totalCredit := 0.0
		for _, list := range jr.ListDetail {
			totalDebet += list.Debet
			totalCredit += list.Credit
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 10.0
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			No := strconv.Itoa(list.No)
			pdf.MultiCell(widthHead[0], 4, No, "", "C", false)
			a0 := pdf.GetY()
			pdf.SetY(y1)
			x += widthHead[0]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[1], 4, list.PostingDate.Local().Format("2006 Jan 02"), "", "L", false)
			a1 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[2], 4, list.DocumentNumber, "", "L", false)
			a2 := pdf.GetY()
			pdf.SetY(y1)

			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			acccode := strconv.Itoa(list.Acc_Code)
			pdf.MultiCell(widthHead[3], 4, acccode, "", "L", false)
			a3 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[4], 4, list.Acc_Name, "", "L", false)
			a4 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			debet := tk.Sprintf("%.2f", list.Debet)
			debet = c.ConvertToCurrency(debet)
			pdf.MultiCell(widthHead[5], 4, list.Description, "", "L", false)
			a5 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			Credit := tk.Sprintf("%.2f", list.Credit)
			Credit = c.ConvertToCurrency(Credit)
			pdf.MultiCell(widthHead[6], 4, debet, "", "R", false)
			a6 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[7], 4, Credit, "", "R", false)
			a7 := pdf.GetY()
			pdf.SetY(y1)

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

		y3 := pdf.GetY()
		pdf.Line(pdf.GetX()-165, y0, pdf.GetX()-165, y3) //garis vertical 1
		pdf.Line(pdf.GetX()-156, y0, pdf.GetX()-156, y3) //garis vertical 2
		pdf.Line(pdf.GetX()-136, y0, pdf.GetX()-136, y3) //garis vertical 3
		pdf.Line(pdf.GetX()-105, y0, pdf.GetX()-105, y3) //garis vertical 4
		pdf.Line(pdf.GetX()-85, y0, pdf.GetX()-85, y3)   //garis vertical 5
		pdf.Line(pdf.GetX()-60, y0, pdf.GetX()-60, y3)   //garis vertical 6
		pdf.Line(pdf.GetX()-25, y0, pdf.GetX()-25, y3)   //garis vertical 7
		pdf.Line(pdf.GetX(), y0, pdf.GetX(), y3)         //garis vertical 8
		pdf.Line(pdf.GetX()+25, y0, pdf.GetX()+25, y3)   //garis vertical 8

		pdf.Line(pdf.GetX()-165, y3, 200, y3) //garis horizontal

		Balance := totalDebet - totalCredit
		y3 = pdf.GetY()
		pdf.SetY(y3)
		pdf.SetX(10 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Total", "LBR", "C", false)
		pdf.SetY(pdf.GetY())

		x := 10.0

		pdf.SetX(x)

		b1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totDeb := tk.Sprintf("%.2f", totalDebet)
		totDeb = c.ConvertToCurrency(totDeb)
		pdf.MultiCell(widthHead[6], 4, totDeb, "RB", "R", false)
		b2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totCred := tk.Sprintf("%.2f", totalCredit)
		totCred = c.ConvertToCurrency(totCred)
		pdf.MultiCell(widthHead[7], 4, totCred, "RB", "R", false)
		pdf.SetY(y3)

		b3 := pdf.GetY()
		allB := []float64{b1, b2, b3}
		var n, biggest float64
		for _, v := range allB {
			if v > n {
				n = v
				biggest = n
			}
			pdf.SetY(biggest)
		}

		y3 = pdf.GetY()
		pdf.SetY(y3)
		x = 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Balance", "LBR", "C", false)

		c1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[6], 4, "", "BR", "R", false)
		c2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		fixBalance := tk.Sprintf("%.2f", Balance)
		pdf.MultiCell(widthHead[6], 4, fixBalance, "RB", "R", false)
		pdf.SetY(y1)

		c3 := pdf.GetY()
		allC := []float64{c1, c2, c3}
		for _, v := range allC {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		fullName := ""
		if jr.User == k.Session("username").(string) {
			fullName = k.Session("fullname").(string)
		} else {
			csr, e := c.Ctx.Connection.NewQuery().Select().From(NewSysUserModel().TableName()).Cursor(nil)
			if e != nil {
				return e.Error()
			}
			user := SysUserModel{}
			e = csr.Fetch(&user, 1, false)
			if e != nil {
				return e.Error()
			}
			csr.Close()
			fullName = user.Fullname
		}
		pdf.Ln(5)
		pdf.SetX(5)
		pdf.CellFormat(50, 6, "User", "", 0, "C", false, 0, "")
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "Approve", "", 0, "C", false, 0, "")
		pdf.Ln(18)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(40, 6, fullName, "B", 0, "C", false, 0, "")
		pdf.Ln(2)
		pdf.GetY()
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "__________________________", "", 0, "C", false, 0, "")
		//end
		pdf.Ln(1)
		pdf.SetX(10)
	}
	e = os.RemoveAll(c.PdfPath + "/journal")
	if e != nil {
		return ""
	}
	e = os.MkdirAll(c.PdfPath+"/journal", 0777)
	if e != nil {
		return ""
	}

	namepdf := "-journal.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + "-" + "posting" + namepdf
	fileName := FixName
	location := c.PdfPath + "/journal"
	err := pdf.OutputFileAndClose(location + "/" + fileName)
	if err != nil {
		tk.Println(err.Error())
	}
	return fileName
}
func (c *TransactionController) PrintDraftJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		ID string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Journal").Where(dbox.Eq("_id", p.ID)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []MainJournal{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	JournalHeadTable := []string{"No", "Date", "Document Number", "Acc. Code", "Acc. Name", "Description", "Debet", "Credit"}
	for _, jr := range results {
		pdf.AddPage()
		pdf.SetXY(10, 5)

		pdf.SetFont("Arial", "", 12)
		// pdf.Image(c.LogoFile+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
		pdf.Ln(2)
		yy := pdf.GetY()
		pdf.SetY(yy + 4)
		pdf.SetX(10)
		pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
		pdf.SetY(yy + 10)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
		pdf.SetX(152)
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 15, "DRAFT JOURNAL", "", 0, "L", false, 0, "")

		pdf.SetY(yy + 17)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.City+"-"+corp.Country, "", 0, "L", false, 0, "")
		pdf.SetX(10)

		pdf.SetY(yy + 23)
		pdf.SetX(10)
		pdf.SetFont("Arial", "", 11)
		pdf.CellFormat(0, 12, corp.NoTelp, "", 0, "L", false, 0, "")
		pdf.SetX(10)
		pdf.Line(pdf.GetX(), pdf.GetY()+9, 200, pdf.GetY()+9) //garis horizontal1

		pdf.SetFont("Arial", "B", 14)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)

		pdf.SetFont("Arial", "", 7)
		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Journal ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Journal_Type, "", 0, "L", false, 0, "")

		pdf.Ln(5)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(10, 10, "Status ", "", 0, "L", false, 0, "")
		pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, jr.Status, "", 0, "L", false, 0, "")

		pdf.Ln(10)

		y0 := pdf.GetY()
		widthHead := []float64{9.0, 20.0, 31.0, 20.0, 25.0, 35.0, 25.0, 25.0}
		for i, head := range JournalHeadTable {
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
		y1 := pdf.GetY()
		totalDebet := 0.0
		totalCredit := 0.0
		for _, list := range jr.ListDetail {
			totalDebet += list.Debet
			totalCredit += list.Credit
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 10.0
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			No := strconv.Itoa(list.No)
			pdf.MultiCell(widthHead[0], 4, No, "", "C", false)

			a0 := pdf.GetY()
			pdf.SetY(y1)
			x += widthHead[0]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[1], 4, list.PostingDate.Local().Format("2006 Jan 02"), "", "L", false)
			a1 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[2], 4, list.DocumentNumber, "", "L", false)
			a2 := pdf.GetY()
			pdf.SetY(y1)

			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			acccode := strconv.Itoa(list.Acc_Code)
			pdf.MultiCell(widthHead[3], 4, acccode, "", "L", false)
			a3 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[4], 4, list.Acc_Name, "", "L", false)
			a4 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			debet := tk.Sprintf("%.2f", list.Debet)
			debet = c.ConvertToCurrency(debet)
			pdf.MultiCell(widthHead[5], 4, list.Description, "", "L", false)
			a5 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			Credit := tk.Sprintf("%.2f", list.Credit)
			Credit = c.ConvertToCurrency(Credit)
			pdf.MultiCell(widthHead[6], 4, debet, "", "R", false)
			a6 := pdf.GetY()
			pdf.SetY(y1)
			x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
			pdf.SetLeftMargin(x)
			pdf.SetX(x)
			pdf.MultiCell(widthHead[7], 4, Credit, "", "R", false)
			a7 := pdf.GetY()
			pdf.SetY(y1)

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

		y3 := pdf.GetY()
		pdf.Line(pdf.GetX()-165, y0, pdf.GetX()-165, y3) //garis vertical 1
		pdf.Line(pdf.GetX()-156, y0, pdf.GetX()-156, y3) //garis vertical 2
		pdf.Line(pdf.GetX()-136, y0, pdf.GetX()-136, y3) //garis vertical 3
		pdf.Line(pdf.GetX()-105, y0, pdf.GetX()-105, y3) //garis vertical 4
		pdf.Line(pdf.GetX()-85, y0, pdf.GetX()-85, y3)   //garis vertical 5
		pdf.Line(pdf.GetX()-60, y0, pdf.GetX()-60, y3)   //garis vertical 6
		pdf.Line(pdf.GetX()-25, y0, pdf.GetX()-25, y3)   //garis vertical 7
		pdf.Line(pdf.GetX(), y0, pdf.GetX(), y3)         //garis vertical 8
		pdf.Line(pdf.GetX()+25, y0, pdf.GetX()+25, y3)   //garis vertical 8
		pdf.Line(pdf.GetX()-165, y3, 200, y3)            //garis horizontal

		Balance := totalDebet - totalCredit

		y3 = pdf.GetY()
		pdf.SetY(y3)
		pdf.SetX(10 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Total", "LBR", "C", false)
		pdf.SetY(pdf.GetY())

		x := 10.0

		pdf.SetX(x)

		b1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totDeb := tk.Sprintf("%.2f", totalDebet)
		totDeb = c.ConvertToCurrency(totDeb)
		pdf.MultiCell(widthHead[6], 4, totDeb, "RB", "R", false)
		b2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		totCred := tk.Sprintf("%.2f", totalCredit)
		totCred = c.ConvertToCurrency(totCred)
		pdf.MultiCell(widthHead[7], 4, totCred, "RB", "R", false)
		pdf.SetY(y3)

		b3 := pdf.GetY()
		allB := []float64{b1, b2, b3}
		var n, biggest float64
		for _, v := range allB {
			if v > n {
				n = v
				biggest = n
			}
			pdf.SetY(biggest)
		}
		y3 = pdf.GetY()
		pdf.SetY(y3)
		x = 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4])
		pdf.MultiCell(widthHead[5], 4, "Balance", "LBR", "C", false)

		c1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[6], 4, "", "BR", "C", false)
		c2 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		fixBalance := tk.Sprintf("%.2f", Balance)
		pdf.MultiCell(widthHead[6], 4, fixBalance, "RB", "R", false)
		pdf.SetY(y3)

		c3 := pdf.GetY()
		allC := []float64{c1, c2, c3}
		for _, v := range allC {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
		fullName := ""
		if jr.User == k.Session("username").(string) {
			fullName = k.Session("fullname").(string)
		} else {
			csr, e := c.Ctx.Connection.NewQuery().Select().From(NewSysUserModel().TableName()).Cursor(nil)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			user := SysUserModel{}
			e = csr.Fetch(&user, 1, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			csr.Close()
			fullName = user.Fullname
		}
		yApp := pdf.GetY()
		pdf.SetY(yApp)
		pdf.Ln(5)
		pdf.SetX(5)
		pdf.CellFormat(50, 6, "User", "", 0, "C", false, 0, "")
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "Approve", "", 0, "C", false, 0, "")
		pdf.Ln(18)
		pdf.GetY()
		pdf.SetX(10)
		pdf.CellFormat(40, 6, fullName, "B", 0, "C", false, 0, "")
		pdf.Ln(2)
		pdf.GetY()
		pdf.SetX(156.7)
		pdf.CellFormat(50, 6, "__________________________", "", 0, "C", false, 0, "")
		//end
		pdf.Ln(1)
		pdf.SetX(10)
	}

	e = os.RemoveAll(c.PdfPath + "/listjournal")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/listjournal", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}

	namepdf := "-listjournal.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/listjournal"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	// tk.Println(e)
	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName
}
func (c *TransactionController) ExportToExcelALl(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Journal").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	defer csr.Close()
	results := []MainJournal{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	// var font *xlsx.Style
	file = xlsx.NewFile()
	sheet, e = file.AddSheet("Sheet1")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	font := xlsx.NewFont(11, "Calibri")
	style := xlsx.NewStyle()
	style.Font = *font
	row = sheet.AddRow()
	IDJournal := row.AddCell()
	IDJournal.Value = "Id_Journal"
	IDJournal.SetStyle(style)
	Date := row.AddCell()
	Date.Value = "Date"
	Date.SetStyle(style)
	Journal_Type := row.AddCell()
	Journal_Type.Value = "Journal_Type"
	Journal_Type.SetStyle(style)
	Desc := row.AddCell()
	Desc.Value = "Description"
	Desc.SetStyle(style)
	DocNum := row.AddCell()
	DocNum.Value = "DocumentNumber"
	DocNum.SetStyle(style)
	User := row.AddCell()
	User.Value = "User"
	User.SetStyle(style)
	Department := row.AddCell()
	Department.Value = "Department"
	Department.SetStyle(style)
	SalesName := row.AddCell()
	SalesName.Value = "Sales"
	SalesName.SetStyle(style)
	// ID := row.AddCell()
	// ID.Value = "Id"
	// ID.SetStyle(style)
	for _, each := range results {
		row = sheet.AddRow()
		IDJournal := row.AddCell()
		IDJournal.Value = each.IdJournal
		Date := row.AddCell()
		Date.Value = each.DateStr
		Journal_Type := row.AddCell()
		Journal_Type.Value = each.Journal_Type
		Desc := row.AddCell()
		Desc.Value = each.ListDetail[0].Description
		DocNum := row.AddCell()
		DocNum.Value = each.ListDetail[0].DocumentNumber
		User := row.AddCell()
		User.Value = each.User
		Department := row.AddCell()
		Department.Value = each.Department
		SalesName := row.AddCell()
		SalesName.Value = each.SalesName
		// ID := row.AddCell()
		// ID.Value = each.ID
		IDJournal.SetStyle(style)
		Date.SetStyle(style)
		Journal_Type.SetStyle(style)
		Desc.SetStyle(style)
		DocNum.SetStyle(style)
		User.SetStyle(style)
		Department.SetStyle(style)
		SalesName.SetStyle(style)
		// ID.SetStyle(style)
	}
	e = os.RemoveAll(c.PdfPath + "/journal")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/journal", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	FixName := time.Now().Format("2006-01-02T150405") + "MyXLSXFile.xlsx"
	fileName := FixName
	location := c.UploadPath + "/journal/"
	e = file.Save(location + fileName)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}
func (c *TransactionController) GetDataDepartment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	crs, e := c.Ctx.Connection.NewQuery().From("Department").Select().Cursor(nil)
	defer crs.Close()

	result := []DepartmentModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		// tk.Println(e.Error())
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "OK", result)
}
func (c *TransactionController) SaveNewDepartment(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DepartmentName string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}
	numbStr := ""
	kode := ""
	name := strings.ToUpper(strformat.Filter(p.DepartmentName, strformat.CharsetAlphaNumeric))[:3]
	numb := 0
	numb = c.GetLastNumberDepartment()
	if numb < 10 {
		numbStr = "000" + strconv.Itoa(numb)
	} else if numb <= 10 && numb < 100 {
		numbStr = "00" + strconv.Itoa(numb)
	} else if numb <= 100 && numb < 1000 {
		numbStr = "0" + strconv.Itoa(numb)
	} else {
		numbStr = strconv.Itoa(numb)
	}
	kode = "DEPT/" + name + "/" + numbStr
	model := NewDepartmentModel()
	model.DepartmentCode = kode
	model.DepartmentName = p.DepartmentName
	e = c.Ctx.Save(model)
	return c.SetResultInfo(false, "OK", nil)
}
func (c *TransactionController) GetLastNumberDepartment() int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceCustomer").Select().Where(dbox.Eq("collname", "department")).Cursor(nil)
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
		model.Collname = "department"
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
func (c *TransactionController) UploadFiles(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filename := k.Request.FormValue("filename")
	pathToSave, _ := filepath.Abs("assets/docs/journal/uploads")
	os.MkdirAll(pathToSave, 0777)
	e, _ := UploadHandlerCopy(k, "filedoc", pathToSave+"/")
	if e != nil {
		tk.Println("Error : " + e.Error())
		return c.SetResultInfo(true, e.Error(), nil)
	}

	fileToProcess := pathToSave + "/" + filename

	_, err := os.Stat(fileToProcess)
	if os.IsNotExist(err) {
		tk.Println(err.Error())
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// ==================================TEALEG====================================
	model := NewMainJournal()
	// gl := NewGeneralLedger()
	excelFileName := fileToProcess
	xlFile, er := xlsx.OpenFile(excelFileName)
	if er != nil {
		return c.SetResultInfo(true, er.Error(), nil)
	}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			model.ID = tk.RandomString(32)
			m := time.Now().UTC().Month()
			y := time.Now().UTC().Year()
			codejurnal := fmt.Sprintf("%02d%d", m, y)
			ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
			number := fmt.Sprintf("%04d", ids)
			model.IdJournal = "JUR/" + codejurnal + "/" + number
			postingDate, _ := cells[1].Float()
			model.PostingDate = xlsx.TimeFromExcelTime(postingDate, false)
			model.CreateDate = xlsx.TimeFromExcelTime(postingDate, false)
			model.DateStr = tk.Date2String(model.PostingDate, "dd MMM yyyy")
			model.User = k.Session("username").(string)
			if len(cells) > 12 {
				model.Department, _ = cells[12].String()
			}
			if len(cells) > 13 {
				model.SalesCode, _ = cells[13].String()
			}
			if len(cells) > 14 {
				model.SalesName, _ = cells[14].String()
			}
			typJournal, _ := cells[11].String()
			tk.Println("=>", typJournal)
			typJournalLow := strings.ToLower(typJournal)
			headcode := "GEM"
			if typJournalLow == "general journal" {
				model.Journal_Type = "General"
			} else if typJournalLow == "cash out" {
				model.Journal_Type = "CashOut"
				headcode = "BKK"
			} else if typJournalLow == "cash in" {
				model.Journal_Type = "CashIn"
				headcode = "BBM"
			}
			model.Status = "draft"
			listDetail := []Journal{}
			list := Journal{}
			list.Id = tk.RandomString(12)
			list.No = 1
			list.Journal_Type = model.Journal_Type
			list.PostingDate = model.PostingDate
			list.DateStr = model.DateStr
			// NYI: Change this to auto sequence

			// idx, _ := c.GetNextIdSeq("DocumentNumber", list.Journal_Type, int(m), y)
			// dn := fmt.Sprintf("%04d", idx)
			// list.DocumentNumber = headcode + "/" + tk.Date2String(model.PostingDate, "ddMMyy") + "/" + dn
			list.DocumentNumber = headcode + "/" + tk.Date2String(model.PostingDate, "ddMMyy") + "/temp-" + model.User + "/0001"
			// docCell, _ := cells[2].String()
			// docSplit := strings.Split(docCell, "/")
			// list.DocumentNumber, _ = cells[2].String()
			// if docSplit[0] == "MJO" {
			// 	list.DocumentNumber = "GEM" + "/" + docSplit[1] + "/" + docSplit[2]
			// }

			list.Acc_Code, _ = cells[3].Int()
			list.Acc_Name, _ = cells[4].String()
			list.Debet, _ = cells[5].Float()
			list.Description, _ = cells[10].String()
			list.Department = model.Department
			list.SalesCode = model.SalesCode
			list.SalesName = model.SalesName
			list.User = model.User
			listDetail = append(listDetail, list)

			list2 := Journal{}
			list2.Id = tk.RandomString(12)
			list2.No = 1
			list2.Journal_Type = model.Journal_Type
			list2.PostingDate = model.PostingDate
			list2.DateStr = model.DateStr
			list2.DocumentNumber = list.DocumentNumber
			list2.Acc_Code, _ = cells[7].Int()
			list2.Acc_Name, _ = cells[8].String()
			list2.Credit, _ = cells[9].Float()
			list2.Description, _ = cells[10].String()
			list2.Department = model.Department
			list2.SalesCode = model.SalesCode
			list2.SalesName = model.SalesName
			list2.User = model.User
			listDetail = append(listDetail, list2)
			model.ListDetail = listDetail
			no, _ := cells[0].Int()
			if no < 1734 {
				c.Ctx.Save(model)
			}
			// // GL
			// gl.ID = tk.RandomString(32)
			// gl.IdJournal = model.IdJournal
			// gl.PostingDate = model.PostingDate
			// gl.CreateDate = model.CreateDate
			// gl.DateStr = model.DateStr
			// gl.User = model.User
			// gl.Journal_Type = model.Journal_Type
			// gl.User = model.User
			// gl.Status = model.Status
			// glListDetail := []GeneralDetail{}

			// glList := GeneralDetail{}
			// glList.Id = tk.RandomString(32)
			// glList.No = list.No
			// glList.Journal_Type = model.Journal_Type
			// glList.PostingDate = list.PostingDate
			// glList.DateStr = list.DateStr
			// glList.DocumentNumber = list.DocumentNumber
			// glList.Acc_Code = list.Acc_Code
			// glList.Acc_Name = list.Acc_Name
			// glList.Debet = list.Debet
			// glList.Description = list.Description
			// glList.Attachment = list.Attachment
			// glList.User = list.User
			// glListDetail = append(glListDetail, glList)

			// glList2 := GeneralDetail{}
			// glList2.Id = tk.RandomString(32)
			// glList2.No = list2.No
			// glList2.Journal_Type = model.Journal_Type
			// glList2.PostingDate = list2.PostingDate
			// glList2.DateStr = list2.DateStr
			// glList2.DocumentNumber = list2.DocumentNumber
			// glList2.Acc_Code = list2.Acc_Code
			// glList2.Acc_Name = list2.Acc_Name
			// glList2.Credit = list2.Credit
			// glList2.Description = list2.Description
			// glList2.Attachment = list2.Attachment
			// glList2.User = list2.User
			// glListDetail = append(glListDetail, glList2)
			// gl.ListDetail = glListDetail
			// if no < 1734 {
			// 	c.Ctx.Save(gl)
			// }
		}
	}

	os.Remove(fileToProcess)

	return c.SetResultInfo(false, "Success", nil)
}

// func (c *TransactionController) UploadFilesDepartment(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson
// 	filename := k.Request.FormValue("filename")
// 	pathToSave, _ := filepath.Abs("assets/docs/journal")
// 	os.MkdirAll(pathToSave, 0777)
// 	e, _ := UploadHandlerCopy(k, "filedoc", pathToSave+"/")
// 	if e != nil {
// 		// tk.Println("Error : " + e.Error())
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}

// 	fileToProcess := pathToSave + "/" + filename

// 	_, err := os.Stat(fileToProcess)
// 	if os.IsNotExist(err) {
// 		// tk.Println(err.Error())
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	// ==================================TEALEG====================================
// 	excelFileName := fileToProcess
// 	xlFile, er := xlsx.OpenFile(excelFileName)
// 	if er != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	csr, e := c.Ctx.Connection.NewQuery().From("Journal").Select().Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csr.Close()
// 	dataJournal := []MainJournal{}
// 	e = csr.Fetch(&dataJournal, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	csrr, e := c.Ctx.Connection.NewQuery().From("GeneralLedger").Select().Cursor(nil)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	defer csrr.Close()
// 	dataGeneral := []MainGeneralLedger{}
// 	e = csrr.Fetch(&dataGeneral, 0, false)
// 	if e != nil {
// 		return c.SetResultInfo(true, e.Error(), nil)
// 	}
// 	for _, sheet := range xlFile.Sheets {
// 		for i, row := range sheet.Rows {
// 			if i == 0 {
// 				continue
// 			}
// 			cells := row.Cells
// 			depart, _ := cells[6].String()
// 			idJournal, _ := cells[0].String()
// 			for i, _ := range dataJournal {
// 				if dataJournal[i].IdJournal == idJournal {

// 					dataJournal[i].Department = strings.ToUpper(depart)
// 					for j, _ := range dataJournal[i].ListDetail {
// 						dataJournal[i].ListDetail[j].Department = strings.ToUpper(depart)
// 					}
// 					e = c.Ctx.Save(&dataJournal[i])
// 					break
// 				}
// 			}
// 			for i, _ := range dataGeneral {
// 				if dataGeneral[i].IdJournal == idJournal {
// 					dataGeneral[i].Department = strings.ToUpper(depart)
// 					for j, _ := range dataGeneral[i].ListDetail {
// 						dataGeneral[i].ListDetail[j].Department = strings.ToUpper(depart)
// 					}
// 					e = c.Ctx.Save(&dataGeneral[i])
// 					break
// 				}
// 			}
// 			csr, e = c.Ctx.Connection.NewQuery().From("Department").Select().Where(dbox.Eq("departmentname", strings.ToUpper(depart))).Cursor(nil)
// 			if e != nil {
// 				return c.SetResultInfo(true, e.Error(), nil)
// 			}
// 			defer csr.Close()
// 			results := []DepartmentModel{}
// 			e = csr.Fetch(&results, 0, false)
// 			if e != nil {
// 				return c.SetResultInfo(true, e.Error(), nil)
// 			}
// 			if csr.Count() == 0 {
// 				numbStr := ""
// 				kode := ""
// 				name := strings.ToUpper(depart[:3])
// 				numb := 0
// 				numb = c.GetLastNumberDepartment()
// 				if numb < 10 {
// 					numbStr = "000" + strconv.Itoa(numb)
// 				} else if numb <= 10 && numb < 100 {
// 					numbStr = "00" + strconv.Itoa(numb)
// 				} else if numb <= 100 && numb < 1000 {
// 					numbStr = "0" + strconv.Itoa(numb)
// 				} else {
// 					numbStr = strconv.Itoa(numb)
// 				}
// 				kode = "DEPT/" + name + "/" + numbStr
// 				model := NewDepartmentModel()
// 				model.DepartmentCode = kode
// 				model.DepartmentName = depart
// 				e = c.Ctx.Save(model)
// 			}
// 		}
// 	}
// 	c.LogActivity("Journal", "Upload file department", filename, k)
// 	return c.SetResultInfo(false, "Success", nil)
// }
func (c *TransactionController) UnapplyJournal(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id   string
		Date string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	date, _ := time.Parse("2006-01-02", p.Date)
	csr, e := c.Ctx.Connection.NewQuery().From("Journal").Select().Where(dbox.Eq("IdJournal", p.Id)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultsJ := []MainJournal{}
	e = csr.Fetch(&resultsJ, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(resultsJ)
	// journal := resultsJ[0]
	m := date.Month()
	y := date.Year()
	codejurnal := fmt.Sprintf("%02d%d", m, y)
	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
	number := fmt.Sprintf("%04d", ids)
	idJournal := ""
	for _, journal := range resultsJ {
		journal.ID = tk.RandomString(32)
		idJournal = journal.ID
		journal.IdJournal = "JUR/" + codejurnal + "/" + number
		journal.CreateDate = time.Now()
		journal.PostingDate = date
		journal.Journal_Type = "General"
		journal.User = k.Session("username").(string)
		journal.DateStr = tk.Date2String(journal.PostingDate, "dd MMM yyyy")
		listdetail := []Journal{}
		totDebet := 0.0
		totCrdit := 0.0
		lastDocnum := ""
		for _, each := range journal.ListDetail {
			each.User = k.Session("username").(string)
			each.Journal_Type = "General"
			each.PostingDate = journal.PostingDate
			each.DateStr = journal.DateStr
			headcode := "GEM/"
			balance := totDebet - totCrdit
			dateDoc := tk.Date2String(each.PostingDate, "ddMMyy")
			if balance == 0 {
				// docNum = each.DocumentNumber
				idx, _ := c.GetNextIdSeq("DocumentNumber", "General", int(m), y)
				number := fmt.Sprintf("%04d", idx)
				each.DocumentNumber = headcode + dateDoc + "/" + number
				lastDocnum = each.DocumentNumber
			} else {
				each.DocumentNumber = lastDocnum
			}
			each.Id = tk.RandomString(10)
			unapplyDebet := each.Credit
			unapplyCredit := each.Debet
			each.Debet = unapplyDebet
			each.Credit = unapplyCredit
			if each.Attachment != "BEGIN" || each.Attachment != "INVOICE" || each.Attachment != "" {
				each.Attachment = ""
			}
			totDebet = totDebet + each.Debet
			totCrdit = totCrdit + each.Credit
			listdetail = append(listdetail, each)
		}
		journal.ListDetail = listdetail
		e = c.Ctx.Save(&journal)
	}
	csr, e = c.Ctx.Connection.NewQuery().From("GeneralLedger").Select().Where(dbox.Eq("IdJournal", p.Id)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultsG := []MainGeneralLedger{}
	e = csr.Fetch(&resultsG, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// general := resultsG[0]
	for _, general := range resultsG {
		general.ID = tk.RandomString(32)
		general.IdJournal = "JUR/" + codejurnal + "/" + number
		general.CreateDate = time.Now()
		general.PostingDate = date
		general.DateStr = tk.Date2String(general.PostingDate, "dd MMM yyyy")
		general.User = k.Session("username").(string)
		general.Journal_Type = "General"
		listdetailG := []GeneralDetail{}
		totDebet := 0.0
		totCrdit := 0.0
		lastDocnum := ""
		for _, each := range general.ListDetail {
			each.User = k.Session("username").(string)
			each.Journal_Type = "General"
			each.PostingDate = general.PostingDate
			each.DateStr = general.DateStr
			headcode := "GEM/"
			balance := totDebet - totCrdit
			dateDoc := tk.Date2String(each.PostingDate, "ddMMyy")
			if balance == 0 {
				// docNum = each.DocumentNumber
				idx, _ := c.GetNextIdSeq("DocumentNumberGL", "General", int(m), y)
				number := fmt.Sprintf("%04d", idx)
				each.DocumentNumber = headcode + dateDoc + "/" + number
				lastDocnum = each.DocumentNumber
			} else {
				each.DocumentNumber = lastDocnum
			}
			each.Id = tk.RandomString(10)

			unapplyDebet := each.Credit
			unapplyCredit := each.Debet

			// each.DocumentNumber = headcode + dateDoc + "/" + number
			each.Debet = unapplyDebet
			each.Credit = unapplyCredit
			if each.Attachment != "BEGIN" || each.Attachment != "INVOICE" || each.Attachment != "" {
				each.Attachment = ""
			}
			totDebet = totDebet + each.Debet
			totCrdit = totCrdit + each.Credit
			listdetailG = append(listdetailG, each)
		}
		general.ListDetail = listdetailG
		e = c.Ctx.Save(&general)
	}
	c.LogActivity("UnApply Journal", "Insert Journal", idJournal, k)
	return c.SetResultInfo(false, "Success", nil)
}

func (c *TransactionController) SavePostingEx(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	e := k.Request.ParseMultipartForm(1000000000)
	if e != nil {
		c.ErrorResultInfo(e.Error(), nil)
	}
	_, formData, err := k.GetPayloadMultipart(nil)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// Get destination data
	payload := struct {
		ID   string `json:"Id"`
		Role string `json:"Role"`
	}{}
	payload.ID = formData["Id"][0]
	payload.Role = formData["Role"][0]

	ids := payload.ID
	roles := payload.Role

	ctx := c.Ctx.Connection
	data := make([]MainJournal, 0)

	crsData, errData := ctx.NewQuery().From("Journal").Select().Where(dbox.Eq("_id", ids)).Cursor(nil)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	errData = crsData.Fetch(&data, 0, false)
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	singleData := MainJournal{}
	if crsData.Count() > 0 {
		singleData = data[0]
	}
	// Upload Attachment Part
	{
		//Create Directory
		m := singleData.CreateDate.Month()
		y := singleData.CreateDate.Year()
		folder := fmt.Sprintf("%d%02d", y, m)
		baseImagePath := ReadConfig()["uploadpath"]
		pathfolder := filepath.Join(baseImagePath, folder)
		if _, err = os.Stat(pathfolder); os.IsNotExist(err) {
			os.MkdirAll(pathfolder, 0777)
		}

		// move file to upload directory
		if len(formData["Details"]) > 0 {
			var details = formData["Details"][0]
			var detIDs = strings.Split(details, ",")
			for _, detID := range detIDs {
				file, handler, err := k.Request.FormFile("fileUpload_" + detID)
				if file != nil {
					// tk.Println(file)
					if err != nil {
						return c.ErrorResultInfo(err.Error(), nil)
					}

					defer file.Close()

					tk.RandomString(32)

					fileName := tk.RandomString(6) + filepath.Ext(handler.Filename)
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

					for i := range data {
						for w := range data[i].ListDetail {
							if data[i].ListDetail[w].Id == detID || data[i].ListDetail[w].Attachment == "" {
								data[i].ListDetail[w].Attachment = fileName
							}
						}
					}
				}
			}
		}

	}
	// Posting Part
	{
		time.Now().UTC()
		var m time.Month
		var y int

		ModelJournal := NewMainJournal()
		ModelGeneral := NewGeneralLedger()
		for _, dt := range data {
			NewListDetail := []Journal{}
			ModelJournal.ID = dt.ID
			ModelJournal.IdJournal = dt.IdJournal
			ModelJournal.CreateDate = dt.CreateDate

			if roles == "supervisor" {
				ModelJournal.PostingDate = dt.PostingDate
				ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
			} else if roles == "administrator" {
				ModelJournal.PostingDate = dt.PostingDate
				ModelJournal.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
			} else {
				ModelJournal.PostingDate = time.Now().UTC()
				ModelJournal.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
			}
			ModelJournal.User = dt.User
			ModelJournal.Journal_Type = dt.Journal_Type
			ModelJournal.Department = dt.Department
			ModelJournal.SalesCode = dt.SalesCode
			ModelJournal.SalesName = dt.SalesName
			ModelJournal.Status = "posting"
			numbertemporary := 0
			if ModelJournal.Department == "MULTIJOURNAL" {
				valDebet := 0.0
				valCredit := 0.0
				for i, arrList := range dt.ListDetail {
					headcode := arrList.DocumentNumber[:11]
					// if roles == "supervisor" {
					m = dt.PostingDate.Month()
					y = dt.PostingDate.Year()
					// } else {
					// 	m = time.Now().UTC().Month()
					// 	y = time.Now().UTC().Year()
					// }
					balance := valDebet - valCredit
					valDebet += arrList.Debet
					valCredit += arrList.Credit
					if i == 0 {
						idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						arrList.DocumentNumber = headcode + number
						numbertemporary = idx
					} else {
						// tk.Println(i, arrList.DocumentNumber, valDebet, valCredit)
						if balance == 0 {
							idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
							number := fmt.Sprintf("%04d", idx)
							arrList.DocumentNumber = headcode + number
							numbertemporary = idx
							// valDebet = 0
							// valCredit = 0
						} else {
							number := fmt.Sprintf("%04d", numbertemporary)
							arrList.DocumentNumber = headcode + number
						}
					}
					// tk.Println(arrList.DocumentNumber)
					arrList.PostingDate = ModelJournal.PostingDate
					arrList.DateStr = tk.Date2String(ModelJournal.PostingDate, "dd MMM yyyy")

					NewListDetail = append(NewListDetail, arrList)
				}
			} else {
				for i, arrList := range dt.ListDetail {
					headcode := arrList.DocumentNumber[:11]
					// tk.Println(i)
					if i%2 == 0 {
						// if roles == "supervisor" {
						m = dt.PostingDate.Month()
						y = dt.PostingDate.Year()
						// } else {
						// 	m = time.Now().UTC().Month()
						// 	y = time.Now().UTC().Year()
						// }
						idx, _ := c.GetNextIdSeq("DocumentNumber", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						arrList.DocumentNumber = headcode + number
						numbertemporary = idx
					} else {
						number := fmt.Sprintf("%04d", numbertemporary)
						arrList.DocumentNumber = headcode + number
					}

					arrList.PostingDate = ModelJournal.PostingDate
					arrList.DateStr = tk.Date2String(ModelJournal.PostingDate, "dd MMM yyyy")

					NewListDetail = append(NewListDetail, arrList)
				}
			}
			ModelJournal.ListDetail = NewListDetail
			c.Ctx.Save(ModelJournal)

			//==================General Ledger Save=======================
			NewListDetailGeneral := []GeneralDetail{}
			ModelGeneral.ID = dt.ID
			ModelGeneral.IdJournal = dt.IdJournal
			ModelGeneral.CreateDate = dt.CreateDate
			if roles == "supervisor" {
				ModelGeneral.PostingDate = dt.PostingDate
				ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
			} else if roles == "administrator" {
				ModelGeneral.PostingDate = dt.PostingDate
				ModelGeneral.DateStr = tk.Date2String(dt.PostingDate, "dd MMM yyyy")
			} else {
				ModelGeneral.PostingDate = time.Now().UTC()
				ModelGeneral.DateStr = tk.Date2String(time.Now().UTC(), "dd MMM yyyy")
			}
			ModelGeneral.User = dt.User
			ModelGeneral.Journal_Type = dt.Journal_Type
			ModelGeneral.Status = "posting"
			ModelGeneral.Department = dt.Department
			ModelGeneral.SalesCode = dt.SalesCode
			ModelGeneral.SalesName = dt.SalesName
			numbertemporarytwo := 0
			if ModelJournal.Department == "MULTIJOURNAL" {
				valDebet := 0.0
				valCredit := 0.0
				for u, arrListGeneral := range dt.ListDetail {
					datalist := GeneralDetail{}
					headcode := arrListGeneral.DocumentNumber[:11]
					// if roles == "supervisor" {
					m = dt.PostingDate.Month()
					y = dt.PostingDate.Year()
					// } else {
					// 	m = time.Now().UTC().Month()
					// 	y = time.Now().UTC().Year()
					// }

					balance := valDebet - valCredit
					valDebet += arrListGeneral.Debet
					valCredit += arrListGeneral.Credit
					if u == 0 {
						idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						datalist.DocumentNumber = headcode + number
						numbertemporarytwo = idx
					} else {
						if balance == 0 {
							idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
							number := fmt.Sprintf("%04d", idx)
							datalist.DocumentNumber = headcode + number
							numbertemporarytwo = idx
							// valDebet = 0
							// valCredit = 0
						} else {
							number := fmt.Sprintf("%04d", numbertemporarytwo)
							datalist.DocumentNumber = headcode + number
						}
					}

					datalist.PostingDate = ModelGeneral.PostingDate
					datalist.DateStr = tk.Date2String(ModelGeneral.PostingDate, "dd MMM yyyy")

					datalist.Id = arrListGeneral.Id
					datalist.No = arrListGeneral.No
					datalist.Journal_Type = arrListGeneral.Journal_Type
					datalist.Acc_Code = arrListGeneral.Acc_Code
					datalist.Acc_Name = arrListGeneral.Acc_Name
					datalist.Debet = arrListGeneral.Debet
					datalist.Credit = arrListGeneral.Credit
					datalist.Description = arrListGeneral.Description
					datalist.Attachment = arrListGeneral.Attachment
					datalist.User = arrListGeneral.User
					datalist.Department = arrListGeneral.Department
					datalist.SalesCode = arrListGeneral.SalesCode
					datalist.SalesName = arrListGeneral.SalesName

					NewListDetailGeneral = append(NewListDetailGeneral, datalist)
				}
			} else {
				for u, arrListGeneral := range dt.ListDetail {
					datalist := GeneralDetail{}
					headcode := arrListGeneral.DocumentNumber[:11]

					if u%2 == 0 {
						// if roles == "supervisor" {
						m = dt.PostingDate.Month()
						y = dt.PostingDate.Year()
						// } else {
						// 	m = time.Now().UTC().Month()
						// 	y = time.Now().UTC().Year()
						// }
						idx, _ := c.GetNextIdSeq("DocumentNumberGL", dt.Journal_Type, int(m), y)
						number := fmt.Sprintf("%04d", idx)
						datalist.DocumentNumber = headcode + number
						numbertemporarytwo = idx
					} else {
						number := fmt.Sprintf("%04d", numbertemporarytwo)
						datalist.DocumentNumber = headcode + number
					}

					datalist.PostingDate = ModelGeneral.PostingDate
					datalist.DateStr = tk.Date2String(ModelGeneral.PostingDate, "dd MMM yyyy")

					datalist.Id = arrListGeneral.Id
					datalist.No = arrListGeneral.No
					datalist.Journal_Type = arrListGeneral.Journal_Type
					datalist.Acc_Code = arrListGeneral.Acc_Code
					datalist.Acc_Name = arrListGeneral.Acc_Name
					datalist.Debet = arrListGeneral.Debet
					datalist.Credit = arrListGeneral.Credit
					datalist.Description = arrListGeneral.Description
					datalist.Attachment = arrListGeneral.Attachment
					datalist.User = arrListGeneral.User
					datalist.Department = arrListGeneral.Department
					datalist.SalesCode = arrListGeneral.SalesCode
					datalist.SalesName = arrListGeneral.SalesName

					NewListDetailGeneral = append(NewListDetailGeneral, datalist)
				}
			}
			ModelGeneral.ListDetail = NewListDetailGeneral
			c.Ctx.Save(ModelGeneral)
			c.LogActivity("Journal", "Save Posting Journal", dt.IdJournal, k)
		}
	}

	return c.SetResultOK(nil)
}
