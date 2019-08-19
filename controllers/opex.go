package controllers

import (
	"../helpers"
	"../library/tealeg/xlsx"
	. "../models"
	"os"
	"strconv"
	"time"

	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

type OpexModel struct {
	Id         bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code   int
	Acc_Name   string
	Amount     float64
	Percentage float64
}
type OpexModelPeriode struct {
	Id       bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code int
	Acc_Name string
	Amount   float64
	Periode  []tk.M
}

func (c *ReportController) GetDataOpex(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart          string
		DateEnd            string
		Filter             bool
		TextSearch         string
		Department         []string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}

	parent := []int{8000, 6000, 9000}
	accCode := []int{5211, 5212}
	var filter []tk.M
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	match := []tk.M{}
	match = append(match, tk.M{"main_acc_code": tk.M{"$in": parent}})
	match = append(match, tk.M{"acc_code": tk.M{"$in": accCode}})
	filter = append(filter, tk.M{"$match": tk.M{"$or": match}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":      "$_id",
		"acc_code": "$acc_code",
		"acc_name": "$account_name",
	}})
	if p.Filter == true {
		if p.TextSearch != "" {
			filter = append(filter, tk.M{"$match": tk.M{"acc_name": tk.M{"$regex": ".*" + p.TextSearch + ".*"}}})
		}
	}
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Cursor(nil)
	defer crs.Close()
	resultsCOA := []OpexModel{}
	e = crs.Fetch(&resultsCOA, 0, false)
	if e != nil {
		c.SetResultInfo(false, e.Error(), nil)
	}
	if p.Filter && len(resultsCOA) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

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

	for _, sub := range resultJournal {
		for i, _ := range resultsCOA {
			if sub.GetInt("_id") == resultsCOA[i].Acc_Code {
				resultsCOA[i].Amount = sub.GetFloat64("SumDebet") - sub.GetFloat64("SumCredit")
			}
		}
	}

	return c.SetResultInfo(false, "Success", resultsCOA)
}

func (c *ReportController) GetDataDetailOpex(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart          string
		DateEnd            string
		AccCode            int
		Department         []string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
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
		"Amount":         tk.M{"$subtract": amount},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": p.AccCode}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := make([]tk.M, 0)
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", resultJournal)
}
func (c *ReportController) ExportDetailOpexToExcell(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart          string
		DateEnd            string
		Department         []string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	account := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5211 || each.ACC_Code == 5212 {
			account = append(account, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 || each.Main_Acc_Code == 8000 || each.Main_Acc_Code == 9000 {
			account = append(account, each.ACC_Code)
		}
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
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
		"Amount":         tk.M{"$subtract": amount},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": account}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	resultJournal := make([]tk.M, 0)
	e = crs.Fetch(&resultJournal, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
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
	// fill := xlsx.NewFill("solid", "00FF0000", "00FFFFFF")
	style := xlsx.NewStyle()
	headStyle := xlsx.NewStyle()
	// headStyle.Fill = *fill
	headStyle.Font = *fontHead
	style.Font = *font

	footerStyle2 := xlsx.NewStyle()
	rightalign := xlsx.Alignment{Horizontal: "right"}
	footerStyle2.Alignment = rightalign
	footerStyle2.Font = *font
	footerStyle2.ApplyAlignment = true

	row = sheet.AddRow()
	Date := row.AddCell()
	Date.Value = "Date"
	Date.SetStyle(headStyle)
	DocNum := row.AddCell()
	DocNum.Value = "Document Number"
	DocNum.SetStyle(headStyle)
	accCode := row.AddCell()
	accCode.Value = "Account Code"
	accCode.SetStyle(headStyle)
	accName := row.AddCell()
	accName.Value = "Account Name"
	accName.SetStyle(headStyle)
	Department := row.AddCell()
	Department.Value = "Department"
	Department.SetStyle(headStyle)
	SalesCode := row.AddCell()
	SalesCode.Value = "Sales"
	SalesCode.SetStyle(headStyle)
	Description := row.AddCell()
	Description.Value = "Description"
	Description.SetStyle(headStyle)
	Amount := row.AddCell()
	Amount.Value = "Amount"
	Amount.SetStyle(headStyle)
	total := 0.0
	for _, each := range resultJournal {
		row = sheet.AddRow()
		Date := row.AddCell()
		Date.Value = each.GetString("DateStr")
		Date.SetStyle(style)
		DocNum := row.AddCell()
		DocNum.Value = each.GetString("DocumentNumber")
		DocNum.SetStyle(style)
		accCode := row.AddCell()
		accCode.Value = each.GetString("Acc_Code")
		accCode.SetStyle(style)
		accName := row.AddCell()
		accName.Value = each.GetString("Acc_Name")
		accName.SetStyle(style)
		Department := row.AddCell()
		Department.Value = each.GetString("Department")
		Department.SetStyle(style)
		SalesCode := row.AddCell()
		SalesCode.Value = each.GetString("SalesName")
		SalesCode.SetStyle(style)
		Description := row.AddCell()
		Description.Value = each.GetString("Description")
		Description.SetStyle(style)
		Amount := row.AddCell()
		amount := tk.Sprintf("%.2f", each.GetFloat64("Amount"))
		if each.GetFloat64("Amount") < 0 {
			amount = "(" + tk.Sprintf("%.2f", each.GetFloat64("Amount")*-1) + ")"
		}
		Amount.Value = c.ConvertToCurrency(amount)
		Amount.SetStyle(footerStyle2)
		total = total + each.GetFloat64("Amount")
	}

	row = sheet.AddRow()
	// totalCell := row.AddCell()
	// totalCell.Value = "Total Operating Expenses"
	// totalCell.SetStyle(headStyle)
	Date = row.AddCell()
	Date.Value = "Total Operating Expenses"
	footerStyle1 := xlsx.NewStyle()
	centerHalign := xlsx.Alignment{Horizontal: "center"}
	footerStyle1.Alignment = centerHalign
	footerStyle1.ApplyAlignment = true
	footerStyle1.Font = *fontHead
	Date.SetStyle(footerStyle1)
	Date.Merge(5, 0)
	DocNum = row.AddCell()
	DocNum.Value = ""
	DocNum.SetStyle(style)
	accCode = row.AddCell()
	accCode.Value = ""
	accCode.SetStyle(style)
	accName = row.AddCell()
	accName.Value = ""
	accName.SetStyle(style)
	Department = row.AddCell()
	Department.Value = ""
	Department.SetStyle(style)
	SalesCode = row.AddCell()
	SalesCode.Value = ""
	SalesCode.SetStyle(style)
	Description = row.AddCell()
	Description.Value = ""
	Description.SetStyle(style)
	Amount = row.AddCell()
	totalStr := tk.Sprintf("%.2f", total)
	if total < 0 {
		totalStr = "(" + tk.Sprintf("%.2f", total*-1) + ")"
	}
	Amount.Value = c.ConvertToCurrency(totalStr)
	// totalCell.Merge(1, 0)

	Amount.SetStyle(footerStyle2)
	FixName := time.Now().Format("2006-01-02T150405") + "-opex.xlsx"
	fileName := FixName
	location := c.UploadPath + "/report/excel/"
	e = file.Save(location + fileName)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	return c.SetResultInfo(false, "success", fileName)
}
func (c *ReportController) ExportPdfOpex(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Department         []string
		DateStart          string
		DateEnd            string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	account := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5211 || each.ACC_Code == 5212 {
			account = append(account, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 || each.Main_Acc_Code == 8000 || each.Main_Acc_Code == 9000 {
			account = append(account, each.ACC_Code)
		}
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
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
		"Amount":         tk.M{"$subtract": amount},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": account}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer crs.Close()
	DATA := make([]tk.M, 0)
	e = crs.Fetch(&DATA, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
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
	pdf.CellFormat(0, 15, "OPERATING EXPENSES", "", 0, "L", false, 0, "")

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
	pdf.SetFont("Century_Gothic", "", 7)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	// if p.Filter == true {
	pdf.CellFormat(10, 10, "Date Periode  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	// pdf.Ln(4)
	// // }
	// pdf.GetY()
	// pdf.SetX(12.0)
	// pdf.CellFormat(10, 10, "Account Code  ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, ": "+p.AccCode, "", 0, "L", false, 0, "")
	// pdf.Ln(4)
	// pdf.GetY()
	// pdf.SetX(12.0)
	// pdf.CellFormat(10, 10, "Account Name  ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, ": "+p.AccName, "", 0, "L", false, 0, "")
	// pdf.SetFont("Century_Gothic", "", 9)
	// pdf.SetX(200)
	// pdf.CellFormat(10, 10, "Begining ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, ": "+p.Begining, "", 0, "L", false, 0, "")
	pdf.Ln(8)
	// pdf.SetFont("Century_Gothic", "", 6)
	coaHead := []string{"No. ", "Date", "Document Number", "Account", "Account Name", "Department", "Sales", "Description", "Amount"}
	widthHead := []float64{10, 20.0, 30.0, 15.0, 55.0, 25.0, 25.0, 60.0, 30.0}
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
	totalAmount := 0.0
	lastbigest := 0.0
	var length = len(DATA) + 1
	for i, each := range DATA {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 12.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.GetString("DateStr"), "", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[2], 5, each.GetString("DocumentNumber"), "", "L", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, strconv.Itoa(each.GetInt("Acc_Code")), "", "L", false)
		a3 := pdf.GetY()
		x += widthHead[3]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[4], 5, each.GetString("Acc_Name"), "", "L", false)
		a4 := pdf.GetY()
		x += widthHead[4]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[5], 5, each.GetString("Department"), "", "L", false)
		a5 := pdf.GetY()
		x += widthHead[5]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[6], 5, each.GetString("SalesName"), "", "L", false)
		a6 := pdf.GetY()
		x += widthHead[6]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[7], 5, each.GetString("Description"), "R", "L", false)
		a7 := pdf.GetY()
		x += widthHead[7]
		pdf.SetXY(x, y1)
		// pdf.MultiCell(widthHead[7], 5, each.GetString("IdJournal"), "", "L", false)
		totalAmount = totalAmount + each.GetFloat64("Amount")
		Amount := tk.Sprintf("%.2f", each.GetFloat64("Amount"))
		Amount = c.ConvertToCurrency(Amount)
		if each.GetFloat64("Amount") < 0 {
			Amount = "(" + tk.Sprintf("%.2f", each.GetFloat64("Amount")*-1) + ")"
			Amount = c.ConvertToCurrency(Amount)
		}
		pdf.MultiCell(widthHead[8], 5, Amount, "R", "R", false)
		a8 := pdf.GetY()
		// pdf.SetXY(x, y1)
		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7, a8}
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
				pdf.Line(x+widthHead[7], y0, x+widthHead[7], biggest)
				pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)
				// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest)
				// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest) // vertical last
				pdf.Line(12.0, biggest, x+widthHead[7], biggest)
			}
			pdf.AddPage()
			y0 = pdf.GetY()
			if y0 == 10.00125 && i != length {
				pdf.Line(12.0, y0, 12.0, biggest+5) // vertical
				pdf.Line(12.0, y0, x+widthHead[7], y0)
				pdf.Line(x+widthHead[7], y0, x+widthHead[7], biggest+5)
				pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest+5)
				pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest+5)
				// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest+5)
				// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest+5) // vertical last
				pdf.Line(12.0, biggest+5, x+widthHead[7], biggest+5)
				lastbigest = biggest + 5
			}
		}
		if pdf.PageNo() == 1 {
			pdf.Line(12.0, y0, 12.0, biggest)
			pdf.Line(x+widthHead[7], y0, x+widthHead[7], biggest)
			pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)
			// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest)
			// pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest) // vertical last
			// pdf.Line(12.0, biggest, x+widthHead[7], biggest)
		}
	}
	pdf.SetFont("Century_Gothic", "", 9)
	y2 = pdf.GetY()
	// pdf.LinearGradient(x, y, w, h, r1, g1, b1, r2, g2, b2, x1, y1, x2, y2)
	pdf.LinearGradient(11.0, y2+0.2, 280, lastbigest, 255, 255, 255, 255, 255, 255, 12.0, y2, 12.0, lastbigest) //DELETE MORE LINE
	pdf.SetY(pdf.GetY())
	// pdf.Line(12.0, pdf.GetY(), 282.0, pdf.GetY())
	pdf.SetX(12.0)
	pdf.MultiCell(240, 5, "Total", "1", "C", false)
	Amount := tk.Sprintf("%.2f", totalAmount)
	Amount = c.ConvertToCurrency(Amount)
	if totalAmount < 0 {
		Amount = "(" + tk.Sprintf("%.2f", totalAmount*-1) + ")"
		Amount = c.ConvertToCurrency(Amount)
	}
	pdf.SetY(y2)
	pdf.SetX(240.0 + 12.0)
	pdf.MultiCell(30.0, 5, Amount, "TRB", "R", false)

	// Credit := tk.Sprintf("%.2f", totalCredit)
	// Credit = c.ConvertToCurrency(Credit)
	// if totalCredit < 0 {
	// 	Credit = "(" + tk.Sprintf("%.2f", totalCredit*-1) + ")"
	// 	Credit = c.ConvertToCurrency(Credit)
	// }
	// pdf.SetY(y2)
	// pdf.SetX(210.0 + 12.0)
	// pdf.MultiCell(30.0, 5, Credit, "TRB", "R", false)
	// pdf.SetY(y2)
	// pdf.SetX(240.0 + 12.0)
	// Saldo := tk.Sprintf("%.2f", totalSaldo)
	// Saldo = c.ConvertToCurrency(Credit)
	// if totalSaldo < 0 {
	// 	Saldo = "(" + tk.Sprintf("%.2f", totalSaldo*-1) + ")"
	// 	Saldo = c.ConvertToCurrency(Saldo)
	// }
	// pdf.MultiCell(30.0, 5, Saldo, "TRB", "R", false)
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-Opex.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}
func (c *ReportController) GetDataTopTenOpexPie(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Department         []string
		DateStart          string
		DateEnd            string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	account := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5211 || each.ACC_Code == 5212 {
			account = append(account, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 || each.Main_Acc_Code == 8000 || each.Main_Acc_Code == 9000 {
			account = append(account, each.ACC_Code)
		}
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
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
		"Department":     "$ListDetail.Department",
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": account}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{
			"Acc_Name": "$Acc_Name",
			"Acc_Code": "$Acc_Code",
		},
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"Amount":    tk.M{"$sum": "$Amount"},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"Amount": -1}})
	pipes = append(pipes, tk.M{"$limit": 10})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	type PieChartModel struct {
		Account string
		Amount  float64
	}
	newResults := []PieChartModel{}
	for _, each := range results {
		accName := each["_id"].(tk.M).GetString("Acc_Name")
		accCode := each["_id"].(tk.M).GetString("Acc_Code")
		mdl := PieChartModel{}
		mdl.Account = accCode + "-" + accName
		mdl.Amount = each.GetFloat64("Amount")
		newResults = append(newResults, mdl)
	}
	return c.SetResultInfo(false, "success", newResults)
}
func (c *ReportController) GetDataTopTenOpexDetailGrid(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Department         []string
		DateStart          string
		DateEnd            string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	account := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5211 || each.ACC_Code == 5212 {
			account = append(account, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 || each.Main_Acc_Code == 8000 || each.Main_Acc_Code == 9000 {
			account = append(account, each.ACC_Code)
		}
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
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
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": account}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$Acc_Code",
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"Amount":    tk.M{"$sum": "$Amount"},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"Amount": -1}})
	pipes = append(pipes, tk.M{"$limit": 10})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	acc10 := []int{}
	for _, each := range results {
		acc10 = append(acc10, each.GetInt("_id"))
	}
	pipes = []tk.M{}
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
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc10}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	results = make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "success", results)
}
func (c *ReportController) GetDataChartOpex(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Year               int
		Department         []string
		DepartmentContains string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	type ChartModel struct {
		Id     bson.ObjectId
		Month  int
		Amount float64
		Type   string
	}
	fixData := []ChartModel{}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	account := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5211 || each.ACC_Code == 5212 {
			account = append(account, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 || each.Main_Acc_Code == 8000 || each.Main_Acc_Code == 9000 {
			account = append(account, each.ACC_Code)
		}
	}
	dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
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
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": account}}})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		// pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
		if p.Department[0] == "ALL" {
			// tk.Println("bener")
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$nin": p.Department}}})
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id": tk.M{
					"Month":      "$Month",
					"Department": "$Department",
				},
				"Expenses":  tk.M{"$sum": "$Amount"},
				"SumDebet":  tk.M{"$sum": "$Debet"},
				"SumCredit": tk.M{"$sum": "$Credit"},
			}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$in": p.Department}}})
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id": tk.M{
					"Month":      "$Month",
					"Department": "$Department",
				},
				"Expenses":  tk.M{"$sum": "$Amount"},
				"SumDebet":  tk.M{"$sum": "$Debet"},
				"SumCredit": tk.M{"$sum": "$Credit"},
			}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id": tk.M{
				"Month":      "$Month",
				"Department": "$Department",
			},
			"Expenses":  tk.M{"$sum": "$Amount"},
			"SumDebet":  tk.M{"$sum": "$Debet"},
			"SumCredit": tk.M{"$sum": "$Credit"},
		}})
	} else if len(p.Department) == 0 && p.DepartmentContains == "" {
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":      "$Month",
			"Expenses": tk.M{"$sum": "$Amount"},
		}})
	}
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resExp := []tk.M{}
	e = csr.Fetch(&resExp, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(resExp)
	if len(p.Department) == 0 && p.DepartmentContains == "" {
		fixData = []ChartModel{}
		for _, each := range resExp {
			mod := ChartModel{}
			mod.Id = bson.NewObjectId()
			mod.Month = each.GetInt("_id")
			mod.Type = "Opex"
			mod.Amount = each.GetFloat64("Expenses")
			fixData = append(fixData, mod)
		}
	}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			fixData = []ChartModel{}
			for _, each := range resExp {
				mod := ChartModel{}
				mod.Id = bson.NewObjectId()
				mod.Month = each["_id"].(tk.M).GetInt("Month")
				// if mod.Month == 9 {
				// 	tk.Println("===", each)
				// }
				mod.Type = each["_id"].(tk.M).GetString("Department")
				mod.Amount = each.GetFloat64("Expenses")
				fixData = append(fixData, mod)
			}
		} else {
			fixData = []ChartModel{}
			for _, each := range resExp {
				mod := ChartModel{}
				mod.Id = bson.NewObjectId()
				mod.Month = each["_id"].(tk.M).GetInt("Month")
				// mod.Type = p.Department
				mod.Type = each["_id"].(tk.M).GetString("Department")
				mod.Amount = each.GetFloat64("Expenses")
				fixData = append(fixData, mod)
			}
		}
	}
	if len(p.Department) == 0 && p.DepartmentContains != "" {
		fixData = []ChartModel{}
		for _, each := range resExp {
			mod := ChartModel{}
			mod.Id = bson.NewObjectId()
			mod.Month = each["_id"].(tk.M).GetInt("Month")
			mod.Type = each["_id"].(tk.M).GetString("Department")
			mod.Amount = each.GetFloat64("Expenses")
			fixData = append(fixData, mod)
		}
	}
	return c.SetResultInfo(false, "Success", fixData)
}
func (c *ReportController) GetDataOpexPeriode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart          string
		DateEnd            string
		Filter             bool
		TextSearch         string
		Department         []string
		DepartmentContains string
		SalesCode          []string
		SalesContains      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}

	parent := []int{8000, 6000, 9000}
	accCode := []int{5211, 5212}
	var filter []tk.M
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	match := []tk.M{}
	match = append(match, tk.M{"main_acc_code": tk.M{"$in": parent}})
	match = append(match, tk.M{"acc_code": tk.M{"$in": accCode}})
	filter = append(filter, tk.M{"$match": tk.M{"$or": match}})
	filter = append(filter, tk.M{"$project": tk.M{
		"_id":      "$_id",
		"acc_code": "$acc_code",
		"acc_name": "$account_name",
	}})
	if p.Filter == true {
		if p.TextSearch != "" {
			filter = append(filter, tk.M{"$match": tk.M{"acc_name": tk.M{"$regex": ".*" + p.TextSearch + ".*"}}})
		}
	}
	crs, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Command("pipe", filter).Cursor(nil)
	defer crs.Close()
	resultsCOA := []OpexModelPeriode{}
	e = crs.Fetch(&resultsCOA, 0, false)
	if e != nil {
		c.SetResultInfo(false, e.Error(), nil)
	}
	if p.Filter && len(resultsCOA) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	// allDept := []string{"All"}
	if len(p.Department) > 0 {
		if p.Department[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$nin": p.Department}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$in": p.Department}}})
		}
	} else if len(p.Department) == 0 && p.DepartmentContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.Department": tk.M{"$regex": ".*" + p.DepartmentContains + ".*"}}})
	}

	if len(p.SalesCode) > 0 {
		if p.SalesCode[0] == "ALL" {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$nin": p.SalesCode}}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$in": p.SalesCode}}})
		}
	} else if len(p.SalesCode) == 0 && p.SalesContains != "" {
		pipes = append(pipes, tk.M{"$match": tk.M{"ListDetail.SalesCode": tk.M{"$regex": ".*" + p.SalesContains + ".*"}}})
	}

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
		"SalesCode":      "$ListDetail.SalesCode",
		"SalesName":      "$ListDetail.SalesName",
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
			"Acc_Name": "$Acc_Name",
			"Month":    "$Month",
			"Year":     "$Year",
		},
		"Amount":    tk.M{"$sum": "$Amount"},
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
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
	// for _, sub := range resultJournal {
	// 	for i, _ := range resultsCOA {
	// 		if sub.GetInt("_id") == resultsCOA[i].Acc_Code {
	// 			resultsCOA[i].Amount = sub.GetFloat64("SumDebet") - sub.GetFloat64("SumCredit")
	// 		}
	// 	}
	// }
	results := struct {
		DataAcc    []OpexModelPeriode
		DataAmount []tk.M
	}{
		DataAcc:    resultsCOA,
		DataAmount: resultJournal,
	}
	return c.SetResultInfo(false, "Success", results)
}
