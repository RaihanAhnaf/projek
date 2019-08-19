package controllers

import (
	"../helpers"
	. "../models"
	"os"
	"strconv"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
)

type ApModel struct {
	SupplierCode string
	SupplierName string
	Term         int
	DocNum       string
	PurchaseDate time.Time
	Due_Date     time.Time
	Amount       float64
	Total_AR     float64
	Age1         float64
	Age2         float64
	Age3         float64
	Age4         float64
}
type ApModelDetail struct {
	Supplier     string
	SupplierCode string
	Item         []ApModel
}

func (c *ReportController) GetDataAP(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateAr   string
		Supplier string
		Type     string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	filterS := db.Ne("Kode", "")
	if p.Supplier != "" {
		filterS = db.Eq("Kode", p.Supplier)
	}
	// getdata supplier
	csr, e := c.Ctx.Connection.NewQuery().Select().From(NewCustomerModel().TableName()).Where(filterS).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	suppliers := []CustomerModel{}
	e = csr.Fetch(&suppliers, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	datasuppliers := map[string]CustomerModel{}
	for _, each := range suppliers {
		datasuppliers[each.Kode] = each
	}
	// getdata
	dateAr, _ := time.Parse("02-Jan-2006", p.DateAr)
	filter := []*db.Filter{}
	filter = append(filter, db.Eq("Status", "PI"))
	// filter = append(filter, db.Ne("AlreadyPaid", 0.0))
	if p.Supplier != "" {
		filter = append(filter, db.Eq("SupplierCode", p.Supplier))
	}
	csr, e = c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(filter...).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []PurchaseOrder{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	data := []ApModel{}
	for _, each := range results {
		mod := ApModel{}
		mod.SupplierCode = each.SupplierCode
		mod.SupplierName = each.SupplierName
		mod.Term = datasuppliers[each.SupplierCode].TrxCode
		mod.DocNum = each.DocumentNumber
		mod.Total_AR = each.GrandTotalIDR - each.AlreadyPaid
		mod.Amount = each.GrandTotalIDR
		mod.PurchaseDate = each.DatePosting
		mod.Due_Date = each.DatePosting.AddDate(0, 0, mod.Term)
		date1 := mod.Due_Date.AddDate(0, 0, 30)
		date2 := mod.Due_Date.AddDate(0, 0, 60)
		date3 := mod.Due_Date.AddDate(0, 0, 90)
		if dateAr.After(mod.Due_Date) && dateAr.Before(date1) {
			mod.Age1 = mod.Total_AR
		}
		if dateAr.After(date1) && dateAr.Before(date2) {
			mod.Age2 = mod.Total_AR
		}
		if dateAr.After(date2) && dateAr.Before(date3) {
			mod.Age3 = mod.Total_AR
		}
		if dateAr.After(date3) {
			mod.Age4 = mod.Total_AR
		}
		if dateAr.After(mod.Due_Date) {
			data = append(data, mod)
		}
	}
	csr, e = c.Ctx.Connection.NewQuery().From("PurchaseInventory").Where(filter...).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resultsInv := []PurchaseInventory{}
	e = csr.Fetch(&resultsInv, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultsInv {
		mod := ApModel{}
		mod.SupplierCode = each.SupplierCode
		mod.SupplierName = each.SupplierName
		mod.Term = datasuppliers[each.SupplierCode].TrxCode
		mod.DocNum = each.DocumentNumber
		mod.Total_AR = each.GrandTotalIDR - each.AlreadyPaid
		mod.Amount = each.GrandTotalIDR
		mod.PurchaseDate = each.DatePosting
		mod.Due_Date = each.DatePosting.AddDate(0, 0, mod.Term)
		date1 := mod.Due_Date.AddDate(0, 0, 30)
		date2 := mod.Due_Date.AddDate(0, 0, 60)
		date3 := mod.Due_Date.AddDate(0, 0, 90)
		if dateAr.After(mod.Due_Date) && dateAr.Before(date1) {
			mod.Age1 = mod.Total_AR
		}
		if dateAr.After(date1) && dateAr.Before(date2) {
			mod.Age2 = mod.Total_AR
		}
		if dateAr.After(date2) && dateAr.Before(date3) {
			mod.Age3 = mod.Total_AR
		}
		if dateAr.After(date3) {
			mod.Age4 = mod.Total_AR
		}
		if dateAr.After(mod.Due_Date) {
			data = append(data, mod)
		}
	}
	//getdata paid
	filterPaid := []*db.Filter{}
	filterPaid = append(filterPaid, db.Eq("Status", "PAID"))
	dateARstart := dateAr.AddDate(-1, 0, 0)
	filterPaid = append(filterPaid, db.Gte("ListPayment.DatePayment", dateARstart))
	csr, e = c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(db.And(filterPaid...)).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resultsPaidNon := []PurchaseInventory{}
	e = csr.Fetch(&resultsPaidNon, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resultsPaidNon {
		mod := ApModel{}
		mod.SupplierCode = each.SupplierCode
		mod.SupplierName = each.SupplierName
		mod.Term = datasuppliers[each.SupplierCode].TrxCode
		mod.DocNum = each.DocumentNumber
		mod.Amount = each.GrandTotalIDR
		mod.Total_AR = each.GrandTotalIDR - each.AlreadyPaid
		mod.PurchaseDate = each.DatePosting
		mod.Due_Date = each.DatePosting.AddDate(0, 0, mod.Term)
		date1 := mod.Due_Date.AddDate(0, 0, 30)
		date2 := mod.Due_Date.AddDate(0, 0, 60)
		date3 := mod.Due_Date.AddDate(0, 0, 90)
		if dateAr.After(mod.Due_Date) && dateAr.Before(date1) {
			mod.Age1 = mod.Total_AR
		}
		if dateAr.After(date1) && dateAr.Before(date2) {
			mod.Age2 = mod.Total_AR
		}
		if dateAr.After(date2) && dateAr.Before(date3) {
			mod.Age3 = mod.Total_AR
		}
		if dateAr.After(date3) {
			mod.Age4 = mod.Total_AR
		}
		datepayment := each.ListPayment[len(each.ListPayment)-1].DatePayment
		if dateAr.After(mod.Due_Date) && datepayment.After(dateAr) {
			data = append(data, mod)
		}
	}
	return c.SetResultInfo(false, "success", data)
}
func (c *ReportController) ExportPdfAP(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateAr      string
		Type        string
		DataDetail  []ApModelDetail
		DataSummary []ApModel
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	NameFile := ""
	if p.Type == "Summary" {
		NameFile = c.BuildPdfAPSummary(p.DataSummary, p.DateAr)
	} else {
		NameFile = c.BuildPdfAPDetail(p.DataDetail, p.DateAr)
	}
	return c.SetResultFile(false, "Succcess", NameFile)
}
func (c *ReportController) PdfSummaryAP(Data []ApModel, Date string) string {
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return ""
	}
	//header
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
	pdf.SetX(187)
	pdf.SetFont("Century_Gothicb", "B", 16)
	pdf.CellFormat(0, 15, "Account Payable Aging", "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 17)
	pdf.SetX(210)
	pdf.CellFormat(0, 15, "Summary", "", 0, "L", false, 0, "")
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
	//first line
	date, _ := time.Parse("02-Jan-2006", Date)
	pdf.CellFormat(10, 10, "Date  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+date.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")

	//body
	pdf.Ln(8)
	head := []string{"Supplier", "Term", "Due Date", "Total A/R", "Aged 1 - 30", "Aged 31 - 60", "Aged 61 - 90", "Aged > 91"}
	widthHead := []float64{80.0, 10.0, 30.0, 30.0, 30.0, 30.0, 30.0, 30.0}
	y0 := pdf.GetY()
	for i, head := range head {
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
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
		}

	}
	y0 = pdf.GetY()
	pdf.SetY(pdf.GetY())
	totalAR := 0.0
	ta1 := 0.0
	ta2 := 0.0
	ta3 := 0.0
	ta4 := 0.0
	var length = len(Data) + 1
	lastbigest := 0.0
	yfirtTable := pdf.GetY()
	morePage := false
	for i, each := range Data {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 12.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, each.SupplierName, "", "L", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, strconv.Itoa(each.Term), "", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[2], 5, each.Due_Date.Format("02 Jan 2006"), "", "L", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		total := tk.Sprintf("%.2f", each.Total_AR)
		total = c.ConvertToCurrency(total)
		if each.Total_AR < 0 {
			total = "(" + tk.Sprintf("%.2f", each.Total_AR*-1) + ")"
			total = c.ConvertToCurrency(total)
		}
		pdf.MultiCell(widthHead[3], 5, total, "", "R", false)
		a3 := pdf.GetY()
		x += widthHead[3]
		pdf.SetXY(x, y1)
		aged1 := tk.Sprintf("%.2f", each.Age1)
		aged1 = c.ConvertToCurrency(aged1)
		if each.Age1 < 0 {
			aged1 = "(" + tk.Sprintf("%.2f", each.Age1*-1) + ")"
			aged1 = c.ConvertToCurrency(aged1)
		}
		pdf.MultiCell(widthHead[4], 5, aged1, "", "R", false)
		a4 := pdf.GetY()
		x += widthHead[4]
		pdf.SetXY(x, y1)
		aged2 := tk.Sprintf("%.2f", each.Age2)
		aged2 = c.ConvertToCurrency(aged2)
		if each.Age2 < 0 {
			aged2 = "(" + tk.Sprintf("%.2f", each.Age2*-1) + ")"
			aged2 = c.ConvertToCurrency(aged2)
		}
		pdf.MultiCell(widthHead[5], 5, aged2, "", "R", false)
		a5 := pdf.GetY()
		x += widthHead[5]
		pdf.SetXY(x, y1)
		aged3 := tk.Sprintf("%.2f", each.Age3)
		aged3 = c.ConvertToCurrency(aged3)
		if each.Age3 < 0 {
			aged3 = "(" + tk.Sprintf("%.2f", each.Age3*-1) + ")"
			aged3 = c.ConvertToCurrency(aged3)
		}
		pdf.MultiCell(widthHead[6], 5, aged3, "", "R", false)
		a6 := pdf.GetY()
		x += widthHead[6]
		pdf.SetXY(x, y1)
		aged4 := tk.Sprintf("%.2f", each.Age4)
		aged4 = c.ConvertToCurrency(aged4)
		if each.Age4 < 0 {
			aged4 = "(" + tk.Sprintf("%.2f", each.Age4*-1) + ")"
			aged4 = c.ConvertToCurrency(aged4)
		}
		pdf.MultiCell(widthHead[7], 5, aged4, "", "R", false)
		a7 := pdf.GetY()
		x += widthHead[7]
		pdf.SetXY(x, y1)
		totalAR = totalAR + each.Total_AR
		ta1 = ta1 + each.Age1
		ta2 = ta2 + each.Age2
		ta3 = ta3 + each.Age3
		ta4 = ta4 + each.Age4
		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7}
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
				pdf.Line(x, y0, x, biggest)
				pdf.Line(12.0, biggest, x, biggest)
			}
			pdf.AddPage()
			y0 = pdf.GetY()
			// tk.Println(y0, biggest, a7)
			if y0 == 10.00125 && i != length {
				pdf.Line(12.0, y0, 12.0, biggest)
				pdf.Line(x, y0, x, biggest)
				pdf.Line(12.0, y0, x, y0)
				pdf.Line(12.0, biggest, x, biggest)
				lastbigest = biggest + 5
			}
			morePage = true
		}
	}
	y4 := pdf.GetY()
	if !morePage {
		pdf.Line(12, yfirtTable, 12, y4)                                                                                                                                                                                                                     //vertical 1
		pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], yfirtTable, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y4) //vertical 4
	}
	pdf.SetXY(12.0, pdf.GetY())
	pdf.MultiCell(120.0, 5, "Total : ", "TLB", "R", false)
	value := tk.Sprintf("%.2f", totalAR)
	value = c.ConvertToCurrency(value)
	if totalAR < 0 {
		value = "(" + tk.Sprintf("%.2f", totalAR*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(120.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta1)
	value = c.ConvertToCurrency(value)
	if ta1 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta1*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(150.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta2)
	value = c.ConvertToCurrency(value)
	if ta2 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta2*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(180.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta3)
	value = c.ConvertToCurrency(value)
	if ta3 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta3*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(210.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta4)
	value = c.ConvertToCurrency(value)
	if ta4 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta4*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(240.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLBR", "R", false)
	y2 = pdf.GetY()
	pdf.LinearGradient(11.0, y2+0.2, 280, lastbigest, 255, 255, 255, 255, 255, 255, 12.0, y2, 12.0, lastbigest) //DELETE MORE LINE
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return ""
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return ""
	}
	namepdf := "-ArSummary.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return ""
	}
	return fileName
}
func (c *ReportController) PdfDetailAP(Data []ApModelDetail, Date string) string {
	// tk.Println(Data)
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return ""
	}
	//header
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
	pdf.SetX(185)
	pdf.SetFont("Century_Gothicb", "B", 16)
	pdf.CellFormat(0, 15, "Account Payable Aging", "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 17)
	pdf.SetX(210)
	pdf.CellFormat(0, 15, "Detail", "", 0, "L", false, 0, "")
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
	//first line
	date, _ := time.Parse("02-Jan-2006", Date)
	pdf.CellFormat(10, 10, "Date  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+date.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")

	//body
	pdf.Ln(8)
	head := []string{"Supplier", "No.Doc", "Invoice Date", "Due Date", "Aged 1 - 30", "Aged 31 - 60", "Aged 61 - 90", "Aged > 91"}
	widthHead := []float64{60.0, 30.0, 30.0, 30.0, 30.0, 30.0, 30.0, 30.0}
	y0 := pdf.GetY()
	for i, head := range head {
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
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
		}

	}
	y0 = pdf.GetY()
	pdf.SetY(pdf.GetY())
	ta1 := 0.0
	ta2 := 0.0
	ta3 := 0.0
	ta4 := 0.0
	for _, each := range Data {
		lenItem := len(each.Item) * 5
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 12.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], float64(lenItem), each.Supplier, "LB", "L", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		lastYitem := 0.0
		for _, item := range each.Item {
			x = 72.0
			y11 := pdf.GetY()
			pdf.SetXY(x, y11)
			pdf.MultiCell(widthHead[1], 5, item.DocNum, "LB", "L", false)
			a1 := pdf.GetY()
			x += widthHead[1]
			pdf.SetXY(x, y11)
			pdf.MultiCell(widthHead[2], 5, item.PurchaseDate.Format("02 Jan 2006"), "LB", "L", false)
			a2 := pdf.GetY()
			x += widthHead[2]
			pdf.SetXY(x, y11)
			pdf.MultiCell(widthHead[3], 5, item.Due_Date.Format("02 Jan 2006"), "LB", "L", false)
			a3 := pdf.GetY()
			x += widthHead[3]
			pdf.SetXY(x, y11)
			aged1 := tk.Sprintf("%.2f", item.Age1)
			aged1 = c.ConvertToCurrency(aged1)
			if item.Age1 < 0 {
				aged1 = "(" + tk.Sprintf("%.2f", item.Age1*-1) + ")"
				aged1 = c.ConvertToCurrency(aged1)
			}
			pdf.MultiCell(widthHead[4], 5, aged1, "LB", "R", false)
			a4 := pdf.GetY()
			x += widthHead[4]
			pdf.SetXY(x, y11)
			aged2 := tk.Sprintf("%.2f", item.Age2)
			aged2 = c.ConvertToCurrency(aged2)
			if item.Age2 < 0 {
				aged2 = "(" + tk.Sprintf("%.2f", item.Age2*-1) + ")"
				aged2 = c.ConvertToCurrency(aged2)
			}
			pdf.MultiCell(widthHead[5], 5, aged2, "LB", "R", false)
			a5 := pdf.GetY()
			x += widthHead[5]
			pdf.SetXY(x, y11)
			aged3 := tk.Sprintf("%.2f", item.Age3)
			aged3 = c.ConvertToCurrency(aged3)
			if item.Age3 < 0 {
				aged3 = "(" + tk.Sprintf("%.2f", item.Age3*-1) + ")"
				aged3 = c.ConvertToCurrency(aged3)
			}
			pdf.MultiCell(widthHead[6], 5, aged3, "LB", "R", false)
			a6 := pdf.GetY()
			x += widthHead[6]
			pdf.SetXY(x, y11)
			aged4 := tk.Sprintf("%.2f", item.Age4)
			aged4 = c.ConvertToCurrency(aged4)
			if item.Age4 < 0 {
				aged4 = "(" + tk.Sprintf("%.2f", item.Age4*-1) + ")"
				aged4 = c.ConvertToCurrency(aged4)
			}
			pdf.MultiCell(widthHead[7], 5, aged4, "LBR", "R", false)
			a7 := pdf.GetY()
			x += widthHead[7]
			pdf.SetXY(x, y11)
			allA := []float64{a1, a2, a3, a4, a5, a6, a7}
			var n, biggest float64
			for _, v := range allA {
				if v > n {
					n = v
					biggest = n
				}
			}
			pdf.SetY(biggest)
			lastYitem = biggest
			ta1 = ta1 + item.Age1
			ta2 = ta2 + item.Age2
			ta3 = ta3 + item.Age3
			ta4 = ta4 + item.Age4
		}
		allA := []float64{a0, lastYitem}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y4 := pdf.GetY()
	pdf.SetXY(12.0, pdf.GetY())
	pdf.MultiCell(150.0, 5, "Total : ", "TLB", "R", false)
	value := tk.Sprintf("%.2f", ta1)
	value = c.ConvertToCurrency(value)
	if ta1 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta1*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(150.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta2)
	value = c.ConvertToCurrency(value)
	if ta2 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta2*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(180.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta3)
	value = c.ConvertToCurrency(value)
	if ta3 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta3*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(210.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLB", "R", false)
	value = tk.Sprintf("%.2f", ta4)
	value = c.ConvertToCurrency(value)
	if ta4 < 0 {
		value = "(" + tk.Sprintf("%.2f", ta4*-1) + ")"
		value = c.ConvertToCurrency(value)
	}
	pdf.SetXY(240.0+12.0, y4)
	pdf.MultiCell(30.0, 5, value, "TLBR", "R", false)
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return ""
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return ""
	}
	namepdf := "-ArDetail.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return ""
	}
	return fileName
}

func (c *ReportController) BuildPdfAPDetail(Data []ApModelDetail, Date string) string {
	// ====== Preparing data ======
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return ""
	}
	head := []string{"Supplier", "No.Doc", "Invoice Date", "Due Date", "Amount", "Aged 1 - 30", "Aged 31 - 60", "Aged 61 - 90", "Aged > 91"}
	widthHead := []float64{55.0, 40.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0}
	dataAlign := []string{"L", "L", "L", "L", "R", "R", "R", "R", "R"}
	var ta0, ta1, ta2, ta3, ta4 float64

	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	pdf.AddPage()
	pdf.SetFont("Century_Gothic", "", 12)

	// ====== BEGIN Drawing ======
	{ // ====== Draw Report Heading ======
		pdf.SetXY(10, 5)
		pdf.Ln(2)
		y1 := pdf.GetY()
		pdf.SetY(y1 + 4)
		pdf.SetX(11.5)
		pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
		pdf.SetY(y1 + 10)
		pdf.SetX(11.5)
		pdf.SetFont("Century_Gothic", "", 11)
		pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
		pdf.SetX(187)
		pdf.SetFont("Century_Gothicb", "B", 16)
		pdf.CellFormat(0, 15, "Account Payable Aging", "", 0, "L", false, 0, "")
		pdf.SetY(y1 + 17)
		pdf.SetX(210)
		pdf.CellFormat(0, 15, "Detail", "", 0, "L", false, 0, "")
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
		//first line
		date, _ := time.Parse("02-Jan-2006", Date)
		pdf.CellFormat(10, 10, "Date  ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+date.Format("02 January 2006"), "", 0, "L", false, 0, "")
		pdf.SetX(200)
		pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	}
	{ // ====== Draw Table Headers ======
		pdf.Ln(8)
		y0 := pdf.GetY()
		x := 12.0
		for i, head := range head {
			pdf.SetXY(x, y0)
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
			x += widthHead[i]
		}
	}
	{ // ====== Draw Table Data =======
		pdf.SetY(pdf.GetY())
		for _, each := range Data {
			lineHeight := 5
			y := pdf.GetY()
			yd := y
			x := 12.0

			// Format data first
			totalItemHeight := 0
			dataItems := make([][8]string, len(each.Item))
			dataLineHeight := make([]float64, len(each.Item))
			for i, item := range each.Item {
				dataItems[i][0] = item.DocNum
				dataItems[i][1] = item.PurchaseDate.Format("02 Jan 2006")
				dataItems[i][2] = item.Due_Date.Format("02 Jan 2006")
				dataItems[i][3] = c.FormatCurrency(item.Amount)
				dataItems[i][4] = c.FormatCurrency(item.Age1)
				dataItems[i][5] = c.FormatCurrency(item.Age2)
				dataItems[i][6] = c.FormatCurrency(item.Age3)
				dataItems[i][7] = c.FormatCurrency(item.Age4)

				// calculate max line count
				maxLineCount := 1
				for z, data := range dataItems[i] {
					lineCount := len(pdf.SplitLines([]byte(data), widthHead[z+1]))
					if lineCount > maxLineCount {
						maxLineCount = lineCount
					}
				}
				totalItemHeight += maxLineCount * lineHeight
				dataLineHeight[i] = float64(maxLineCount * lineHeight)

				// calculate ta
				ta1 += item.Age1
				ta2 += item.Age2
				ta3 += item.Age3
				ta4 += item.Age4
				ta0 += item.Amount
			}

			// Draw Customer
			pdf.SetXY(x, y)
			pdf.MultiCell(widthHead[0], float64(totalItemHeight), each.Supplier, "LB", "LC", false)
			x += widthHead[0]
			xd := x

			// Draw Data
			for ix, item := range dataItems {
				for c, cell := range item {
					pdf.SetXY(x, yd)
					lineCount := float64(len(pdf.SplitLines([]byte(cell), widthHead[c+1])))
					pdf.MultiCell(widthHead[c+1], dataLineHeight[ix]/lineCount, cell, "LBR", dataAlign[c+1], false)
					x += widthHead[c+1]
				}
				yd = pdf.GetY()
				x = xd
			}
		}
	}
	{ // ====== Draw Report Footer =======
		y := pdf.GetY()
		x := 12.0
		pdf.SetXY(x, y)
		wid := 0.0
		for _, w := range widthHead[0:4] {
			wid += w
		}
		pdf.MultiCell(wid, 5, "Total : ", "LB", "R", false)
		x += wid

		value := c.FormatCurrency(ta0)
		pdf.SetXY(x, y)
		wid = widthHead[4]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta1)
		pdf.SetXY(x, y)
		wid = widthHead[5]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta2)
		pdf.SetXY(x, y)
		wid = widthHead[6]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta3)
		pdf.SetXY(x, y)
		wid = widthHead[7]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta4)
		pdf.SetXY(x, y)
		wid = widthHead[8]
		pdf.MultiCell(wid, 5, value, "LBR", "R", false)
		x += wid
	}
	// ====== END Drawing ======

	// ====== Outputs the file ======
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return ""
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return ""
	}
	namepdf := "-ApDetail.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return ""
	}
	return fileName
}

func (c *ReportController) BuildPdfAPSummary(Data []ApModel, Date string) string {
	// ====== Preparing data ======
	corp, e := helpers.GetDataCorporateJson()
	if e != nil {
		return ""
	}
	head := []string{"Customer", "Term", "Due Date", "Amount", "Total AR", "Aged 1 - 30", "Aged 31 - 60", "Aged 61 - 90", "Aged > 91"}
	widthHead := []float64{55.0, 40.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0, 25.0}
	dataAlign := []string{"L", "L", "L", "R", "R", "R", "R", "R", "R"}
	var ta0, ta1, ta2, ta3, ta4, tt float64

	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	pdf.AddPage()
	pdf.SetFont("Century_Gothic", "", 12)

	// ====== BEGIN Drawing ======
	{ // ====== Draw Report Heading ======
		pdf.SetXY(10, 5)
		pdf.Ln(2)
		y1 := pdf.GetY()
		pdf.SetY(y1 + 4)
		pdf.SetX(11.5)
		pdf.CellFormat(0, 12, corp.Name, "", 0, "L", false, 0, "")
		pdf.SetY(y1 + 10)
		pdf.SetX(11.5)
		pdf.SetFont("Century_Gothic", "", 11)
		pdf.CellFormat(0, 15, corp.Address, "", 0, "L", false, 0, "")
		pdf.SetX(187)
		pdf.SetFont("Century_Gothicb", "B", 16)
		pdf.CellFormat(0, 15, "Account Payable Aging", "", 0, "L", false, 0, "")
		pdf.SetY(y1 + 17)
		pdf.SetX(210)
		pdf.CellFormat(0, 15, "Summary", "", 0, "L", false, 0, "")
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
		//first line
		date, _ := time.Parse("02-Jan-2006", Date)
		pdf.CellFormat(10, 10, "Date  ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+date.Format("02 January 2006"), "", 0, "L", false, 0, "")
		pdf.SetX(200)
		pdf.CellFormat(10, 10, "Date Created  ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	}
	{ // ====== Draw Table Headers ======
		pdf.Ln(8)
		y0 := pdf.GetY()
		x := 12.0
		for i, head := range head {
			pdf.SetXY(x, y0)
			pdf.MultiCell(widthHead[i], 8, head, "1", "C", false)
			x += widthHead[i]
		}
	}
	{ // ====== Draw Table Data =======
		pdf.SetY(pdf.GetY())
		for _, each := range Data {
			lineHeight := 5
			y := pdf.GetY()
			yd := y
			x := 12.0

			// Format data first
			totalItemHeight := 0
			dataItems := make([]string, len(head))
			dataItems[0] = each.SupplierName
			dataItems[1] = each.DocNum
			dataItems[2] = each.Due_Date.Format("02 Jan 2006")
			dataItems[3] = c.FormatCurrency(each.Amount)
			dataItems[4] = c.FormatCurrency(each.Total_AR)
			dataItems[5] = c.FormatCurrency(each.Age1)
			dataItems[6] = c.FormatCurrency(each.Age2)
			dataItems[7] = c.FormatCurrency(each.Age3)
			dataItems[8] = c.FormatCurrency(each.Age4)

			// calculate max line count
			maxLineCount := 1
			for z, data := range dataItems {
				lineCount := len(pdf.SplitLines([]byte(data), widthHead[z]))
				if lineCount > maxLineCount {
					maxLineCount = lineCount
				}
			}
			totalItemHeight += maxLineCount * lineHeight

			// calculate ta
			ta1 += each.Age1
			ta2 += each.Age2
			ta3 += each.Age3
			ta4 += each.Age4
			ta0 += each.Amount
			tt += each.Total_AR

			// Draw Data
			pdf.SetXY(x, y)
			for c, cell := range dataItems {
				pdf.SetXY(x, yd)
				lineCount := float64(len(pdf.SplitLines([]byte(cell), widthHead[c])))
				pdf.MultiCell(widthHead[c], float64(totalItemHeight)/lineCount, cell, "LBR", dataAlign[c], false)
				x += widthHead[c]
			}
		}
	}
	{ // ====== Draw Report Footer =======
		y := pdf.GetY()
		x := 12.0
		pdf.SetXY(x, y)
		wid := 0.0
		for _, w := range widthHead[0:3] {
			wid += w
		}
		pdf.MultiCell(wid, 5, "Total : ", "LB", "R", false)
		x += wid

		value := c.FormatCurrency(ta0)
		pdf.SetXY(x, y)
		wid = widthHead[4]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(tt)
		pdf.SetXY(x, y)
		wid = widthHead[4]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta1)
		pdf.SetXY(x, y)
		wid = widthHead[5]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta2)
		pdf.SetXY(x, y)
		wid = widthHead[6]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta3)
		pdf.SetXY(x, y)
		wid = widthHead[7]
		pdf.MultiCell(wid, 5, value, "LB", "R", false)
		x += wid

		value = c.FormatCurrency(ta4)
		pdf.SetXY(x, y)
		wid = widthHead[8]
		pdf.MultiCell(wid, 5, value, "LBR", "R", false)
		x += wid
	}
	// ====== END Drawing ======

	// ====== Outputs the file ======
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return ""
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return ""
	}
	namepdf := "-ApDetail.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return ""
	}
	return fileName
}
