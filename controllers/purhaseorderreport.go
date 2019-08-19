package controllers

import (
	"os"
	"strconv"
	"strings"
	"time"

	"../helpers"
	. "../models"

	"github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

func (c *ReportController) GetDataPurchaseOrderReport(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   string
		DateEnd     string
		ReportType  string
		ReportBy    string
		ValueFilter string
		IsInventory bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	Results := []tk.M{}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd = dateEnd.AddDate(0, 0, 1)
	data, e := c.GeneratePurchaseOrderReport(dateStart, dateEnd, p.ReportBy, p.ValueFilter, strings.ToLower(p.ReportType), p.IsInventory)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	Results = data
	return c.SetResultInfo(false, "success", Results)
}
func (c *ReportController) GeneratePurchaseOrderReport(datestart time.Time, dateend time.Time, reportby string, valueCode string, repType string, isInventory bool) ([]tk.M, error) {
	pipes := []tk.M{}
	match := tk.M{}.Set("DatePosting", tk.M{}.Set("$gte", datestart).Set("$lte", dateend))
	if valueCode != "" {
		if reportby == "Sales" {
			match.Set("SalesCode", valueCode)
		} else if reportby == "Supplier" {
			match.Set("SupplierCode", valueCode)
		}
	}
	pipes = append(pipes, tk.M{}.Set("$match", match))

	var project tk.M
	var table string
	if !isInventory {
		table = "PurchaseOrder"
		project = tk.M{
			"Date":           "$DatePosting",
			"DocumentNumber": "$DocumentNumber",
			"SupplierCode":   "$SupplierCode",
			"SupplierName":   "$SupplierName",
			"Total":          "$GrandTotalIDR",
		}
		if repType == "detail" {
			pipes = append(pipes, tk.M{}.Set("$unwind", "$ListDetail"))
			project.Set("Item", "$ListDetail.Item")
			project.Set("Qty", "$ListDetail.Qty")
			project.Set("Price", "$ListDetail.PriceIDR")
			project.Set("VAT", "$VAT")
		}
	} else {
		project = tk.M{
			"Date":           "$DatePosting",
			"DocumentNumber": "$DocumentNumber",
			"SupplierCode":   "$SupplierCode",
			"SupplierName":   "$SupplierName",
			"SalesCode":      "$SalesCode",
			"SalesName":      "$SalesName",
			"Total":          "$GrandTotalIDR",
		}
		if repType == "detail" {
			pipes = append(pipes, tk.M{}.Set("$unwind", "$ListDetail"))
			project.Set("Item", "$ListDetail.Item")
			project.Set("Qty", "$ListDetail.Qty")
			project.Set("Price", "$ListDetail.PriceIDR")
			project.Set("VAT", "$VAT")
		}
		table = "PurchaseInventory"
	}

	pipes = append(pipes, tk.M{}.Set("$project", project))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{}.Set("Date", 1)))
	csr, e := c.Ctx.Connection.NewQuery().From(table).Command("pipe", pipes).Cursor(nil)
	//csr, e := c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(filter...).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return nil, e
	}
	results := []tk.M{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return nil, e
	}
	for i, each := range results {
		cutsName := strings.Replace(each.GetString("SupplierName"), each.GetString("SupplierCode")+"-", "", -1)
		results[i].Set("SupplierName", cutsName)
	}
	return results, nil

}
func (c *ReportController) ExportToPdfPurchaseOrderReport(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   string
		DateEnd     string
		ReportType  string
		ReportBy    string
		ValueFilter string
		IsInventory bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd = dateEnd.AddDate(0, 0, 1)
	data, e := c.GeneratePurchaseOrderReport(dateStart, dateEnd, p.ReportBy, p.ValueFilter, strings.ToLower(p.ReportType), p.IsInventory)
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
	if p.ReportType == "Summary" {
		pdf.CellFormat(0, 15, "SUMMARY PURCHASE ORDER REPORT", "", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(0, 15, "DETAIL PURCHASE ORDER REPORT", "", 0, "L", false, 0, "")
	}
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

	pdf.Line(12.0, pdf.GetY()+9, 282, pdf.GetY()+9)
	pdf.SetFont("Century_Gothic", "", 8)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	pdf.CellFormat(15, 10, "Date Periode  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(15, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	if p.ValueFilter != "" && p.ReportBy != "All" {
		pdf.Ln(5)
		if p.ReportBy == "Sales" {
			pdf.GetY()
			pdf.SetX(12)
			pdf.CellFormat(15, 10, "Sales Code ", "", 0, "L", false, 0, "")
			pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(80, 10, ": "+p.ValueFilter, "", 0, "L", false, 0, "")
		} else {
			pdf.GetY()
			pdf.SetX(12)
			pdf.CellFormat(15, 10, "Supplier Code  ", "", 0, "L", false, 0, "")
			pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
			pdf.CellFormat(80, 10, ": "+p.ValueFilter, "", 0, "L", false, 0, "")
		}
	}
	pdf.Ln(8)
	if p.ReportType == "Summary" {
		gridHead := []string{"No. ", "Date", "Document Number", "Supplier Name", "Sales Name", "Total"}
		widthHead := []float64{10, 30.0, 40.0, 80.0, 80.0, 30.0}
		if !p.IsInventory {
			gridHead[4] = ""
			widthHead[4] = 0
			widthHead[5] += 20
			widthHead[2] += 20
			widthHead[3] += 40
		}
		y0 := pdf.GetY()
		for i, head := range gridHead {
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

			if widthHead[i] != 0 {
				pdf.MultiCell(widthHead[i], 10, head, "1", "C", false)
			}

		}
		pdf.SetFont("Century_Gothic", "", 7)
		y0 = pdf.GetY()
		pdf.SetY(pdf.GetY())
		lastbigest := 0.0
		var length = len(data) + 1
		onepage := true
		for i, each := range data {
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 12.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "", "C", false)
			a0 := pdf.GetY()
			x += widthHead[0]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[1], 5, each.Get("Date").(time.Time).Format("02 January 2006"), "", "L", false)
			a1 := pdf.GetY()
			x += widthHead[1]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[2], 5, each.GetString("DocumentNumber"), "", "L", false)
			a2 := pdf.GetY()
			x += widthHead[2]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[3], 5, each.GetString("SupplierName"), "", "L", false)
			a3 := pdf.GetY()
			x += widthHead[3]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[4], 5, each.GetString("SalesName"), "", "L", false)
			a4 := pdf.GetY()
			x += widthHead[4]
			pdf.SetXY(x, y1)
			Total := tk.Sprintf("%.2f", each.GetFloat64("Total"))
			Total = c.ConvertToCurrency(Total)
			pdf.MultiCell(widthHead[5], 5, Total, "", "R", false)
			a5 := pdf.GetY()
			allA := []float64{a0, a1, a2, a3, a4, a5}
			var n, biggest float64
			for _, v := range allA {
				if v > n {
					n = v
					biggest = n
				}
			}
			pdf.SetY(biggest)
			if biggest >= 177.0 {
				onepage = false
				if y0 != 10.00125 {
					pdf.Line(12.0, y0, 12.0, biggest)
					pdf.Line(x+widthHead[5], y0, x+widthHead[5], biggest)
					pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
					pdf.Line(12.0, biggest+5, x+widthHead[5], biggest+5)
				}
				pdf.AddPage()
				y0 = pdf.GetY()
				if y0 == 10.00125 && i != length {
					pdf.Line(12.0, y0, 12.0, biggest+5) // vertical
					pdf.Line(12.0, y0, x+widthHead[5], y0)
					pdf.Line(x+widthHead[5], y0, x+widthHead[5], biggest+5)
					pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest+5)
					pdf.Line(12.0, biggest+5, x+widthHead[5], biggest+5)
					lastbigest = biggest + 5
				}
			}
		}
		y2 = pdf.GetY()
		if onepage {
			pdf.Line(12.0, y0, 12.0, y2) // vertical
			// pdf.Line(12.0, y0, 12.0+widthHead[5], y0)
			// pdf.Line(12.0+widthHead[5], y0, 12.0+widthHead[5], y2)
			pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y2)
			// pdf.Line(12.0, lastbigest+5, 12.0+widthHead[5], lastbigest+5)
		}
		pdf.LinearGradient(11.0, y2+0.2, 280, lastbigest, 255, 255, 255, 255, 255, 255, 12.0, y2, 12.0, lastbigest) //DELETE MORE LINE
		pdf.Line(12.0, y2, 282, y2)
		pdf.SetY(pdf.GetY())
	} else {
		gridHead := []string{"No. ", "Date", "Document Number", "Supplier Name", "Sales Name", "Item", "Amount", "Qty", "Total"}
		widthHead := []float64{10.0, 30.0, 30.0, 40.0, 45.0, 45.0, 30.0, 10.0, 30.0}
		if !p.IsInventory {
			gridHead[4] = ""
			widthHead[4] = 0
			widthHead[8] += 10
			widthHead[2] += 10
			widthHead[3] += 10
			widthHead[5] += 15
		}
		y0 := pdf.GetY()
		for i, head := range gridHead {
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
			if widthHead[i] != 0 {
				pdf.MultiCell(widthHead[i], 10, head, "1", "C", false)
			}
		}
		pdf.SetFont("Century_Gothic", "", 7)
		y0 = pdf.GetY()
		pdf.SetY(pdf.GetY())
		lastbigest := 0.0
		var length = len(data) + 1
		onepage := true
		for i, each := range data {
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 12.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "", "C", false)
			a0 := pdf.GetY()
			x += widthHead[0]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[1], 5, each.Get("Date").(time.Time).Format("02 January 2006"), "", "L", false)
			a1 := pdf.GetY()
			x += widthHead[1]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[2], 5, each.GetString("DocumentNumber"), "", "L", false)
			a2 := pdf.GetY()
			x += widthHead[2]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[3], 5, each.GetString("SupplierName"), "", "L", false)
			a3 := pdf.GetY()
			x += widthHead[3]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[4], 5, each.GetString("SalesName"), "", "L", false)
			a4 := pdf.GetY()
			x += widthHead[4]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[5], 5, each.GetString("Item"), "", "L", false)
			a5 := pdf.GetY()
			x += widthHead[5]
			pdf.SetXY(x, y1)
			Amount := tk.Sprintf("%.2f", each.GetFloat64("Price"))
			Amount = c.ConvertToCurrency(Amount)
			pdf.MultiCell(widthHead[6], 5, Amount, "", "R", false)
			a6 := pdf.GetY()
			x += widthHead[6]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[7], 5, strconv.Itoa(each.GetInt("Qty")), "", "C", false)
			a7 := pdf.GetY()
			x += widthHead[7]
			pdf.SetXY(x, y1)
			Total := tk.Sprintf("%.2f", each.GetFloat64("Total"))
			Total = c.ConvertToCurrency(Total)
			pdf.MultiCell(widthHead[8], 5, Total, "", "R", false)
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
			if biggest >= 177.0 {
				onepage = false
				if y0 != 10.00125 {
					pdf.Line(12.0, y0, 12.0, biggest)
					pdf.Line(x+widthHead[8], y0, x+widthHead[8], biggest)
					pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest)
					pdf.Line(12.0, biggest+5, x+widthHead[8], biggest+5)
				}
				pdf.AddPage()
				y0 = pdf.GetY()
				if y0 == 10.00125 && i != length {
					pdf.Line(12.0, y0, 12.0, biggest+5) // vertical
					pdf.Line(12.0, y0, x+widthHead[8], y0)
					pdf.Line(x+widthHead[8], y0, x+widthHead[8], biggest+5)
					pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest+5)
					pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest+5)
					pdf.Line(12.0, biggest+5, x+widthHead[8], biggest+5)
					lastbigest = biggest + 5
				}
			}
		}
		y2 = pdf.GetY()
		if onepage {
			pdf.Line(12.0, y0, 12.0, y2) // vertical
			// pdf.Line(12.0, y0, 12.0+widthHead[5], y0)
			// pdf.Line(12.0+widthHead[5], y0, 12.0+widthHead[5], y2)
			pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y2)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y2)
			// pdf.Line(12.0, lastbigest+5, 12.0+widthHead[5], lastbigest+5)
		}
		pdf.LinearGradient(11.0, y2+0.2, 280, lastbigest, 255, 255, 255, 255, 255, 255, 12.0, y2, 12.0, lastbigest) //DELETE MORE LINE
		pdf.Line(12.0, y2, 282, y2)
		pdf.SetY(pdf.GetY())
	}

	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-reportinvoice.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}

func (c *ReportController) ExportToPdfPurchaseOrderDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart   string
		DateEnd     string
		ReportType  string
		ReportBy    string
		ValueFilter string
		IsInventory bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd = dateEnd.AddDate(0, 0, 1)

	// === BEGIN FETCHING DATA ===
	pipes := []tk.M{}
	match := tk.M{}.Set("DatePosting", tk.M{}.Set("$gte", dateStart).Set("$lte", dateEnd))
	if p.ValueFilter != "" {
		if p.ReportBy == "Sales" {
			match.Set("SalesCode", p.ValueFilter)
		} else if p.ReportBy == "Supplier" {
			match.Set("SupplierCode", p.ValueFilter)
		}
	}
	pipes = append(pipes, tk.M{}.Set("$match", match))

	var table string
	if !p.IsInventory {
		table = "PurchaseOrder"
	} else {
		table = "PurchaseInventory"
	}

	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{}.Set("Date", 1)))
	csr, e := c.Ctx.Connection.NewQuery().From(table).Command("pipe", pipes).Cursor(nil)
	//csr, e := c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(filter...).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	results := []PurchaseInventory{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// === END FETCHING DATA ===

	pdf := gofpdf.New("L", "mm", "A5", c.FontPath)
	var namepdf string

	// BEGIN FETCHING DATA
	for _, DATA := range results {
		/*if ix > 0 {
			pdf.AddPage()
		}*/

		csr, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.SupplierCode)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		resultsSupp := []CustomerModel{}
		e = csr.Fetch(&resultsSupp, 0, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}
		supp := resultsSupp[0]
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
		if DATA.Currency == "USD" {
			for i, _ := range DATA.ListDetail {
				DATA.ListDetail[i].PriceIDR = 0
				DATA.ListDetail[i].AmountIDR = 0
			}
		} else {
			for i, _ := range DATA.ListDetail {
				DATA.ListDetail[i].PriceUSD = 0
				DATA.ListDetail[i].AmountUSD = 0
			}
		}

		if p.IsInventory {
			/* ======== BEGIN INVENTORY ========== */
			namepdf = "-purchaseorderinv.pdf"
			pdf.SetDrawColor(2, 2, 2)
			pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
			pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
			pdf.AddPage()
			x_defaulft := 10.0
			pdf.SetFont("Century_Gothicb", "B", 15)
			pdf.SetXY(82, 10)
			pdf.CellFormat(0, 12, "PURCHASE ORDER", "", 0, "L", false, 0, "")
			pdf.Ln(5)
			pdf.SetX(85)

			pdf.SetFont("Century_Gothicb", "B", 10)
			pdf.CellFormat(0, 12, DATA.DocumentNumber, "", 0, "L", false, 0, "")
			pdf.SetFont("Century_Gothic", "", 8)
			y0 := pdf.GetY() + 5
			//
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Kode, "", "L", false) // customer name
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "Date", "", "L", false) // date
			date := DATA.DatePosting.Format("January 02, 2006")
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Name, "", "L", false) // customer name
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
			dueDate := DATA.DatePosting.AddDate(0, 0, 30).Format("January 02, 2006")
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(20, 3, "Phone :", "", "L", false) // phone
			pdf.SetXY(30, y0)
			pdf.MultiCell(60, 3, supp.NoTelp, "", "L", false) //phone
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+DATA.DocumentNumber, "", "L", false) // DocumentNo
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Address, "", "L", false) // address
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
			if supp.Bank != "" && supp.AccountNo != "" {
				pdf.MultiCell(60, 3, supp.Bank+"-"+supp.AccountNo, "", "L", false) //rek bank
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
			for i, list := range DATA.ListDetail {
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
			pdf.MultiCell(40, 3, DATA.Remark, "", "L", false) // Remark
			pdf.Ln(5)
			//
			y0 = pdf.GetY()
			pdf.SetY(y0)
			yTotal := pdf.GetY()
			headBottom := []string{"Prepared by :", "Approved by", "Finance", ""}
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
			// xx += widthBottom[2]
			// pdf.SetXY(xx, yB)
			// pdf.MultiCell(widthBottom[3], 4, "(                          )", "", "C", false)
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
			total := tk.Sprintf("%.2f", DATA.TotalIDR)
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
			// y0 = pdf.GetY()
			pdf.SetFont("Century_Gothic", "", 8)
			// pdf.SetXY(30, y0)
			// pdf.MultiCell(150, 3, "Barang yang sudah dibeli tidak bisa di kembalikan / di tukar dengan uang", "", "L", false) // alert
			// pdf.Ln(2)
			y0 = pdf.GetY()
			pdf.SetXY(10, y0)
			pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
			pdf.SetXY(30, y0)
			datenow := time.Now().Format("January 02, 2006")
			pdf.MultiCell(150, 3, datenow, "", "L", false) // date print
			e = os.RemoveAll(c.PdfPath + "/purchaseorder")
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			e = os.MkdirAll(c.PdfPath+"/purchaseorder", 0777)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), "")
			}
			/* ======== END INVENTORY ========== */
		} else {
			/* ======== BEGIN NON-INVENTORY ========== */
			namepdf = "-purchaseorder.pdf"
			pdf.SetDrawColor(2, 2, 2)
			pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
			pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
			pdf.AddPage()
			x_defaulft := 10.0
			pdf.SetFont("Century_Gothicb", "B", 15)
			pdf.SetXY(79, 10)
			pdf.CellFormat(0, 12, "PURCHASE ORDER", "", 0, "L", false, 0, "")
			pdf.Ln(5)
			pdf.SetX(85)

			pdf.SetFont("Century_Gothicb", "B", 10)
			pdf.CellFormat(0, 12, DATA.DocumentNumber, "", 0, "L", false, 0, "")
			pdf.SetFont("Century_Gothic", "", 8)
			y0 := pdf.GetY() + 5
			//
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Kode, "", "L", false) // customer name
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "Date", "", "L", false) // date
			date := DATA.DatePosting.Format("January 02, 2006")
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Name, "", "L", false) // customer name
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
			dueDate := DATA.DatePosting.AddDate(0, 0, 30).Format("January 02, 2006")
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(20, 3, "Phone :", "", "L", false) // phone
			pdf.SetXY(30, y0)
			pdf.MultiCell(60, 3, supp.NoTelp, "", "L", false) //phone
			pdf.SetXY(140, y0)
			pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
			pdf.SetXY(160, y0)
			pdf.MultiCell(40, 3, ": "+DATA.DocumentNumber, "", "L", false) // DocumentNo
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(80, 3, supp.Address, "", "L", false) // address
			pdf.SetXY(140, y0)
			// pdf.MultiCell(20, 3, "Sales", "", "L", false) // sales
			// pdf.SetXY(160, y0)
			// pdf.MultiCell(40, 3, ": "+DATA.SalesName, "", "L", false) // sales
			pdf.Ln(5)
			//
			y0 = pdf.GetY()
			pdf.SetXY(x_defaulft, y0)
			pdf.MultiCell(20, 3, "Rek Bank :", "", "L", false) // rek bank
			pdf.SetXY(30, y0)
			if supp.Bank != "" && supp.AccountNo != "" {
				pdf.MultiCell(60, 3, supp.Bank+"-"+supp.AccountNo, "", "L", false) //rek bank
			} else {
				pdf.MultiCell(60, 3, "", "", "L", false) //rek bank
			}
			pdf.Ln(1)
			//
			y0 = pdf.GetY()
			invHead := []string{"No", "", "Item", "Qty", "Price", "Disc. Amount", "Amount"}
			widHead := []float64{10.0, 0.0, 80.0, 10.0, 30.0, 30.0, 30.0}
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
			for i, list := range DATA.ListDetail {
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
				pdf.MultiCell(widHead[1], 4, "", "", "L", false)
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
			pdf.MultiCell(40, 3, DATA.Remark, "", "L", false) // Remark
			pdf.Ln(5)
			//
			y0 = pdf.GetY()
			pdf.SetY(y0)
			yTotal := pdf.GetY()
			headBottom := []string{"Prepared by :", "Approved by", "Finance", ""}
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
			// xx += widthBottom[2]
			// pdf.SetXY(xx, yB)
			// pdf.MultiCell(widthBottom[3], 4, "(                          )", "", "C", false)
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
			total := tk.Sprintf("%.2f", DATA.TotalIDR)
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
			// y0 = pdf.GetY()
			pdf.SetFont("Century_Gothic", "", 8)
			// pdf.SetXY(30, y0)
			// pdf.MultiCell(150, 3, "Barang yang sudah dibeli tidak bisa di kembalikan / di tukar dengan uang", "", "L", false) // alert
			// pdf.Ln(2)
			y0 = pdf.GetY()
			pdf.SetXY(10, y0)
			pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
			pdf.SetXY(30, y0)
			datenow := time.Now().Format("January 02, 2006")
			pdf.MultiCell(150, 3, datenow, "", "L", false) // date print
			/* ======== END NON-INVENTORY ========== */
		}
	}
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName
}

func (c *TransactionController) ExportToPdfPurchaseOrderInventoryDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	csr, e := c.Ctx.Connection.NewQuery().Select().From("PurchaseInventory").Where(dbox.Eq("_id", p.Id)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []PurchaseInventory{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]

	if DATA.Currency == "USD" {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceIDR = 0
			DATA.ListDetail[i].AmountIDR = 0
		}
		// discount = DATA.Discount
	} else {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceUSD = 0
			DATA.ListDetail[i].AmountUSD = 0
		}
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.SupplierCode)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	resultsSupp := []CustomerModel{}
	e = csr.Fetch(&resultsSupp, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	supp := resultsSupp[0]
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
	pdf.SetXY(82, 10)
	pdf.CellFormat(0, 12, "PURCHASE ORDER", "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(85)

	pdf.SetFont("Century_Gothicb", "B", 10)
	pdf.CellFormat(0, 12, DATA.DocumentNumber, "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothic", "", 8)
	y0 := pdf.GetY() + 5
	//
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Kode, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Date", "", "L", false) // date
	date := DATA.DatePosting.Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Name, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
	dueDate := DATA.DatePosting.AddDate(0, 0, 30).Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(20, 3, "Phone :", "", "L", false) // phone
	pdf.SetXY(30, y0)
	pdf.MultiCell(60, 3, supp.NoTelp, "", "L", false) //phone
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+DATA.DocumentNumber, "", "L", false) // DocumentNo
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, supp.Address, "", "L", false) // address
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
	if supp.Bank != "" && supp.AccountNo != "" {
		pdf.MultiCell(60, 3, supp.Bank+"-"+supp.AccountNo, "", "L", false) //rek bank
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
	for i, list := range DATA.ListDetail {
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
	pdf.MultiCell(40, 3, DATA.Remark, "", "L", false) // Remark
	pdf.Ln(5)
	//
	y0 = pdf.GetY()
	pdf.SetY(y0)
	yTotal := pdf.GetY()
	headBottom := []string{"Prepared by :", "Approved by", "Finance", ""}
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
	// xx += widthBottom[2]
	// pdf.SetXY(xx, yB)
	// pdf.MultiCell(widthBottom[3], 4, "(                          )", "", "C", false)
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
	total := tk.Sprintf("%.2f", DATA.TotalIDR)
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
	// y0 = pdf.GetY()
	pdf.SetFont("Century_Gothic", "", 8)
	// pdf.SetXY(30, y0)
	// pdf.MultiCell(150, 3, "Barang yang sudah dibeli tidak bisa di kembalikan / di tukar dengan uang", "", "L", false) // alert
	// pdf.Ln(2)
	y0 = pdf.GetY()
	pdf.SetXY(10, y0)
	pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
	pdf.SetXY(30, y0)
	datenow := time.Now().Format("January 02, 2006")
	pdf.MultiCell(150, 3, datenow, "", "L", false) // date print
	e = os.RemoveAll(c.PdfPath + "/purchaseorder")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/purchaseorder", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-purchaseorder.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/purchaseorder"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		// e.Error()
		return c.SetResultInfo(true, e.Error(), "")
	}
	return fileName
}
