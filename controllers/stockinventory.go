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

type StockItem struct {
	ItemId            string
	Item              string
	StoreLocationId   int
	StoreLocationName string
	IsFirtStockFilled bool
	FirstStock        int
	PO                int
	CMV               int
	INV               int
	CMI               int
	TS                int
	TR                int
	TotalStock        int
}

func (c *ReportController) GetDataStockItem(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  string
		DateEnd    string
		ItemCode   string
		LocationID int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	filterinv := db.Ne("INVID", "")
	if p.ItemCode != "" {
		filterinv = db.Eq("INVID", p.ItemCode)
	}
	if p.LocationID != 0 {
		filterinv = db.And(filterinv, db.Eq("StoreLocation", p.LocationID))
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From(NewInventoryModel().TableName()).Where(filterinv).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	results := []InventoryModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	datas := []StockItem{}
	for _, each := range results {
		data := StockItem{}
		data.ItemId = each.INVID
		data.Item = each.INVDesc
		data.StoreLocationId = each.StoreLocation
		data.StoreLocationName = each.StoreLocationName
		datas = append(datas, data)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd = dateEnd.AddDate(0, 0, 1)
	filter := []*db.Filter{}
	filter = append(filter, db.Gte("date", dateStart))
	filter = append(filter, db.Lte("date", dateEnd))
	if p.ItemCode != "" {
		filter = append(filter, db.Eq("codeitem", p.ItemCode))
	}
	if p.LocationID != 0 {
		filter = append(filter, db.Eq("storehouseid", p.LocationID))
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From(NewLogInventoryModel().TableName()).Where(db.And(filter...)).Order("_id").Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultLogs := []LogInventoryModel{}
	e = csr.Fetch(&resultLogs, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	//get last stock
	pipes := []tk.M{}
	match := tk.M{}.Set("date", tk.M{}.Set("$lte", dateStart))
	if p.ItemCode != "" {
		match.Set("codeitem", tk.M{}.Set("$eq", p.ItemCode))
	}
	if p.LocationID != 0 {
		match.Set("storehouseid", tk.M{}.Set("$eq", p.LocationID))
	}
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{}.Set("_id", 1)))
	pipes = append(pipes, tk.M{}.Set("$group", tk.M{
		"_id": tk.M{
			"codeitem": "$codeitem",
			"location": "$storehouseid",
		},
		"stock": tk.M{}.Set("$last", "$totalsaldo"),
	}))

	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(NewLogInventoryModel().TableName()).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultlastLogs := []tk.M{}
	e = csr.Fetch(&resultlastLogs, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	// tk.Println(tk.JsonString(resultLogs))
	for i, data := range datas {
		for _, last := range resultlastLogs {
			group := last.Get("_id").(tk.M)
			if data.ItemId == group.GetString("codeitem") && data.StoreLocationId == group.GetInt("location") {
				datas[i].FirstStock = last.GetInt("stock")
				datas[i].TotalStock = datas[i].FirstStock
			}
		}
		for _, log := range resultLogs {
			if data.ItemId == log.CodeItem && data.StoreLocationId == log.StorehouseId {
				if !datas[i].IsFirtStockFilled {
					datas[i].FirstStock = log.StockUnit
				}
				// tk.Println(log.TypeTransaction)
				switch log.TypeTransaction {
				case "PO":
					datas[i].PO += log.CountTransaction
				case "CMV":
					datas[i].CMV += log.CountTransaction
				case "INV":
					datas[i].INV += log.CountTransaction
				case "CMI":
					datas[i].CMI += log.CountTransaction
				case "TS":
					datas[i].TS += log.CountTransaction
				case "TR":
					datas[i].TR += log.CountTransaction
				}
				datas[i].TotalStock = datas[i].FirstStock + datas[i].PO - datas[i].CMV - datas[i].INV + datas[i].CMI - datas[i].TS + datas[i].TR
				datas[i].IsFirtStockFilled = true
			}
		}
	}
	return c.SetResultInfo(false, "success", datas)
}
func (c *ReportController) ExportStockItemPdf(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart  string
		DateEnd    string
		ItemCode   string
		LocationID int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	filterinv := db.Ne("INVID", "")
	if p.ItemCode != "" {
		filterinv = db.Eq("INVID", p.ItemCode)
	}
	if p.LocationID != 0 {
		filterinv = db.And(filterinv, db.Eq("StoreLocation", p.LocationID))
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From(NewInventoryModel().TableName()).Where(filterinv).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}

	results := []InventoryModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	datas := []StockItem{}
	for _, each := range results {
		data := StockItem{}
		data.ItemId = each.INVID
		data.Item = each.INVDesc
		data.StoreLocationId = each.StoreLocation
		data.StoreLocationName = each.StoreLocationName
		datas = append(datas, data)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEndDb := dateEnd.AddDate(0, 0, 1)
	filter := []*db.Filter{}
	filter = append(filter, db.Gte("date", dateStart))
	filter = append(filter, db.Lte("date", dateEndDb))
	if p.ItemCode != "" {
		filter = append(filter, db.Eq("codeitem", p.ItemCode))
	}
	csr, e = c.Ctx.Connection.NewQuery().Select().From(NewLogInventoryModel().TableName()).Where(db.And(filter...)).Order("_id").Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultLogs := []LogInventoryModel{}
	e = csr.Fetch(&resultLogs, 0, false)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	//get last stock
	pipes := []tk.M{}
	match := tk.M{}.Set("date", tk.M{}.Set("$lte", dateStart))
	if p.ItemCode != "" {
		match.Set("codeitem", tk.M{}.Set("$eq", p.ItemCode))
	}
	if p.LocationID != 0 {
		match.Set("storehouseid", tk.M{}.Set("$eq", p.LocationID))
	}
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{}.Set("_id", 1)))
	pipes = append(pipes, tk.M{}.Set("$group", tk.M{
		"_id": tk.M{
			"codeitem": "$codeitem",
			"location": "$storehouseid",
		},
		"stock": tk.M{}.Set("$last", "$totalsaldo"),
	}))

	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From(NewLogInventoryModel().TableName()).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	resultlastLogs := []tk.M{}
	e = csr.Fetch(&resultlastLogs, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	for i, data := range datas {
		for _, last := range resultlastLogs {
			group := last.Get("_id").(tk.M)
			if data.ItemId == group.GetString("codeitem") && data.StoreLocationId == group.GetInt("location") {
				datas[i].FirstStock = last.GetInt("stock")
				datas[i].TotalStock = datas[i].FirstStock
			}
		}
		for _, log := range resultLogs {
			if data.ItemId == log.CodeItem && data.StoreLocationId == log.StorehouseId {
				if !datas[i].IsFirtStockFilled {
					datas[i].FirstStock = log.StockUnit
				}
				switch log.TypeTransaction {
				case "PO":
					datas[i].PO += log.CountTransaction
				case "CMV":
					datas[i].CMV += log.CountTransaction
				case "INV":
					datas[i].INV += log.CountTransaction
				case "CMI":
					datas[i].CMI += log.CountTransaction
				case "TS":
					datas[i].TS += log.CountTransaction
				case "TR":
					datas[i].TR += log.CountTransaction
				}
				datas[i].TotalStock = datas[i].FirstStock + datas[i].PO - datas[i].CMV - datas[i].INV + datas[i].CMI - datas[i].TS + datas[i].TR
				datas[i].IsFirtStockFilled = true
			}
		}
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
	pdf.SetX(200)
	pdf.SetFont("Century_Gothicb", "B", 16)
	pdf.CellFormat(0, 15, "Tracking Stock Item", "", 0, "L", false, 0, "")

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
	// pdf.Ln(5)
	// pdf.SetX(12)
	// pdf.CellFormat(15, 10, "Item Code  ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, ": "+p.ItemCode, "", 0, "L", false, 0, "")
	// pdf.Ln(5)
	// pdf.SetX(12)
	// pdf.CellFormat(15, 10, "Item Name  ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, ": "+p.ItemName, "", 0, "L", false, 0, "")
	pdf.Ln(8)
	gridHead := []string{"No. ", "Item Name", "Location", "Stock", "PO", "CMV", "INV", "CMI", "TS", "TR", "Total Stock"}
	widthHead := []float64{10, 50.0, 40.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 30.0}
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
		if i == 0 {
			pdf.MultiCell(widthHead[i], 10, head, "1", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 10, head, "1", "C", false)
		}

	}
	pdf.SetFont("Century_Gothic", "", 7)
	y0 = pdf.GetY()
	pdf.SetY(y0)
	lastlineYborder := 0.0
	lastlineXborder := 0.0
	for i, each := range datas {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 12.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 8, strconv.Itoa(i+1), "", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 8, each.Item, "", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[2], 8, each.StoreLocationName, "", "L", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 8, strconv.Itoa(each.FirstStock), "", "C", false)
		a3 := pdf.GetY()
		x += widthHead[3]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[4], 8, strconv.Itoa(each.PO), "", "C", false)
		a4 := pdf.GetY()
		x += widthHead[4]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[5], 8, strconv.Itoa(each.CMV), "", "C", false)
		a5 := pdf.GetY()
		x += widthHead[5]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[6], 8, strconv.Itoa(each.INV), "", "C", false)
		a6 := pdf.GetY()
		x += widthHead[6]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[7], 8, strconv.Itoa(each.CMI), "", "C", false)
		a7 := pdf.GetY()
		x += widthHead[7]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[8], 8, strconv.Itoa(each.TS), "", "C", false)
		a8 := pdf.GetY()
		x += widthHead[8]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[9], 8, strconv.Itoa(each.TR), "", "C", false)
		a9 := pdf.GetY()
		x += widthHead[9]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[10], 8, strconv.Itoa(each.TotalStock), "", "C", false)
		a10 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		x += widthHead[10]
		if biggest >= 183.0 {
			pdf.Line(12.0, biggest, x, biggest)

			pdf.Line(12.0, y0, 12.0, biggest)
			// pdf.Line(x+widthHead[9], y0, x+widthHead[9], biggest)
			pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], biggest)
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], biggest)                                                       // vertical last
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9], biggest)                             // vertical last
			pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9]+widthHead[10], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9]+widthHead[10], biggest) // vertical last
			//this for after first page
			pdf.AddPage()
			y0 = pdf.GetY()
			pdf.Line(12.0, y0, x, y0)
			// pdf.SetY(15.0)
		}
		//  else {
		// 	pdf.SetY(biggest)
		// }
		lastlineYborder = biggest
		lastlineXborder = x
	}
	pdf.Line(12.0, y0, 12.0, lastlineYborder)
	pdf.Line(12.0+widthHead[0], y0, 12.0+widthHead[0], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1], y0, 12.0+widthHead[0]+widthHead[1], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7], lastlineYborder)
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8], lastlineYborder)                                                       // vertical last
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9], lastlineYborder)                             // vertical last
	pdf.Line(12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9]+widthHead[10], y0, 12.0+widthHead[0]+widthHead[1]+widthHead[2]+widthHead[3]+widthHead[4]+widthHead[5]+widthHead[6]+widthHead[7]+widthHead[8]+widthHead[9]+widthHead[10], lastlineYborder) // vertical last

	pdf.Line(12.0, lastlineYborder, lastlineXborder, lastlineYborder)
	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-stockitem.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}
