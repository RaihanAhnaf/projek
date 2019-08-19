package controllers

import (
	"../helpers"
	. "../helpers"
	. "../models"
	"os"
	"strconv"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
)

type TransferOrderController struct {
	*BaseController
}

func (c *TransferOrderController) GetAllLocations(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	crs, e := c.Ctx.Connection.NewQuery().From("Location").Select().Cursor(nil)
	defer crs.Close()
	results := make([]LocationModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		return []LocationModel{}
	}
	return results
}
func (c *TransferOrderController) GetUserLocations(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)
	filter := []*db.Filter{}
	filter = CreateLocationFilter(filter, "LocationID", locid, false)

	q := c.Ctx.Connection.NewQuery().From("Location").Select()
	if len(filter) > 0 {
		q.Where(db.And(filter...))
	}
	crs, e := q.Cursor(nil)
	defer crs.Close()
	results := make([]LocationModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		return []LocationModel{}
	}
	return results
}
func (c *TransferOrderController) GetUserDestLocations(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)
	filter := []*db.Filter{}
	filter = CreateLocationFilter(filter, "LocationID", locid, true)

	q := c.Ctx.Connection.NewQuery().From("Location").Select()
	if len(filter) > 0 {
		q.Where(db.And(filter...))
	}
	crs, e := q.Cursor(nil)
	defer crs.Close()
	results := make([]LocationModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		return []LocationModel{}
	}
	return results
}

func (c *TransferOrderController) TransferShipment(k *knot.WebContext) interface{} {
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
func (c *TransferOrderController) TransferReceipt(k *knot.WebContext) interface{} {
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

func (c *TransferOrderController) GetShipmentList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var f = struct {
		DateStart string
		DateEnd   string
		Location  int
		Filter    bool
	}{}
	err := k.GetPayload(&f)
	filter := []*db.Filter{}
	if err == nil && f.Filter {
		dt, _ := time.Parse("2006-01-02", f.DateStart)
		dt = dt.Truncate(24 * time.Hour)
		filter = append(filter, db.Gte("DatePosting", dt))
		dt, _ = time.Parse("2006-01-02", f.DateEnd)
		dt = dt.Truncate(24 * time.Hour).Add(24 * time.Hour)
		filter = append(filter, db.Lt("DatePosting", dt))
	} else {
		filter = append(filter, db.Ne("_id", ""))
	}
	filter = CreateLocationFilter(filter, "StoreHouseFrom", f.Location, false)

	pipe := []tk.M{}

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseFrom", // field in the orders collection
			"foreignField": "LocationID",     // field in the items collection
			"as":           "FromDetail",
		},
	})

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseTo", // field in the orders collection
			"foreignField": "LocationID",   // field in the items collection
			"as":           "ToDetail",
		},
	})

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("TransferShipment").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	results := make([]TransferOrderModelTS, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransferOrderController) GetReceiptList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// var locid = k.Session("locationid").(int)

	var f = struct {
		DateStart string
		DateEnd   string
		Location  int
		Filter    bool
	}{}
	err := k.GetPayload(&f)
	filter := []*db.Filter{}
	if err == nil && f.Filter {
		dt, _ := time.Parse("2006-01-02", f.DateStart)
		dt = dt.Truncate(24 * time.Hour)
		filter = append(filter, db.Gte("DatePosting", dt))
		dt, _ = time.Parse("2006-01-02", f.DateEnd)
		dt = dt.Truncate(24 * time.Hour).Add(24 * time.Hour)
		filter = append(filter, db.Lt("DatePosting", dt))
	} else {
		filter = append(filter, db.Ne("_id", ""))
	}
	filter = CreateLocationFilter(filter, "StoreHouseTo", f.Location, false)

	pipe := []tk.M{}

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseFrom", // field in the orders collection
			"foreignField": "LocationID",     // field in the items collection
			"as":           "FromDetail",
		},
	})

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseTo", // field in the orders collection
			"foreignField": "LocationID",   // field in the items collection
			"as":           "ToDetail",
		},
	})

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From("TransferReceipt").Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	results := make([]TransferOrderModelTR, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransferOrderController) ExportToPdfReportTransferOrder(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	py := struct {
		DateStart     string `json: "DateStart", bson: "DateStart"`
		DateEnd       string `json: "DateEnd", bson: "DateEnd"`
		StoreHouse    string `json: "StoreHouse", bson: "StoreHouse"`
		StoreHouseVal int    `json: "StoreHouseVal", bson: "StoreHouseVal"`
		TransferType  string `json: "TransferType", bson: "TransferType"`
	}{}

	e := k.GetPayload(&py)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", py.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", py.DateEnd)
	//dateEndDb := dateEnd.AddDate(0, 0, 1)
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
	if py.TransferType == "Transfer Shipment" {
		pdf.CellFormat(0, 15, "TRANSFER SHIPMENT REPORT", "", 0, "L", false, 0, "")
	} else {
		pdf.CellFormat(0, 15, "TRANSFER RECEIPT REPORT", "", 0, "L", false, 0, "")
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
	pdf.CellFormat(0, 12, "Telp. 031-5676223", "", 0, "L", false, 0, "")
	pdf.SetX(6)

	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)

	pdf.Line(12.0, pdf.GetY()+9, 282, pdf.GetY()+9)
	pdf.SetFont("Century_Gothic", "", 8)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	pdf.CellFormat(25, 10, "Date Periode  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.SetX(200)
	pdf.CellFormat(25, 10, "Date Created  ", "", 0, "L", false, 0, "")
	pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, ": "+time.Now().Format("02 January 2006"), "", 0, "L", false, 0, "")
	// if p.ValueFilter != "" && p.ReportBy != "All" {
	pdf.Ln(5)
	if py.TransferType == "Transfer Shipment" {
		pdf.GetY()
		pdf.SetX(12)
		pdf.CellFormat(25, 10, "Store House From ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+py.StoreHouse, "", 0, "L", false, 0, "")
	} else {
		pdf.GetY()
		pdf.SetX(12)
		pdf.CellFormat(25, 10, "Store House To  ", "", 0, "L", false, 0, "")
		pdf.CellFormat(8, 10, "", "", 0, "L", false, 0, "")
		pdf.CellFormat(80, 10, ": "+py.StoreHouse, "", 0, "L", false, 0, "")
	}
	// }
	pdf.Ln(8)
	if py.TransferType == "Transfer Shipment" {
		data, e := c.GetExcelData(dateStart, dateEnd, py.StoreHouseVal, "TS")
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		gridHead := []string{"No. ", "Document Number Shipment", "Store House From", "Store House To", "Date", "Description"}
		widthHead := []float64{10, 50.0, 50.0, 50.0, 30.0, 80.0}
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
		pdf.SetFont("Century_Gothic", "", 8)
		y0 = pdf.GetY()
		pdf.SetY(pdf.GetY())
		lastbigest := 0.0
		var length = len(data) + 1
		onepage := true
		for i, each := range data {
			//Kolom Nomor
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 12.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "", "C", false)

			//Kolom Data
			a0 := pdf.GetY()
			x += widthHead[0]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[1], 5, each.GetString("DocumentNumberShipment"), "", "C", false)

			HouseLabel := ""
			locd := each.Get("FromDetail").([]interface{})
			HouseLabel = locd[0].(tk.M).GetString("LocationName")
			// for _, ea := range locd{
			// 	HouseLabel = ea.(tk.M).GetString("LocationName"))
			// }

			a1 := pdf.GetY()
			x += widthHead[1]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[2], 5, HouseLabel+" ("+each.GetString("StoreHouseFrom")+")", "", "C", false)

			locd = each.Get("ToDetail").([]interface{})
			HouseLabel = locd[0].(tk.M).GetString("LocationName")

			a2 := pdf.GetY()
			x += widthHead[2]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[3], 5, HouseLabel+" ("+each.GetString("StoreHouseTo")+")", "", "C", false)
			a3 := pdf.GetY()
			x += widthHead[3]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[4], 5, each.GetString("DateStr"), "", "C", false)
			a4 := pdf.GetY()
			x += widthHead[4]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[5], 5, each.GetString("Description"), "", "C", false)

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
		data, e := c.GetExcelData(dateStart, dateEnd, py.StoreHouseVal, "TR")
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		gridHead := []string{"No. ", "Document Number Receipt", "Store House From", "Store House To", "Date", "Description"}
		widthHead := []float64{10, 50.0, 50.0, 50.0, 30.0, 80.0}
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
		pdf.SetFont("Century_Gothic", "", 8)
		y0 = pdf.GetY()
		pdf.SetY(pdf.GetY())
		lastbigest := 0.0
		var length = len(data) + 1
		onepage := true
		for i, each := range data {
			//Kolom Nomor
			y1 = pdf.GetY()
			pdf.SetY(y1)
			x := 12.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(i+1), "", "C", false)

			//Kolom Data
			a0 := pdf.GetY()
			x += widthHead[0]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[1], 5, each.GetString("DocumentNumberReceipt"), "", "C", false)

			HouseLabel := ""
			locd := each.Get("FromDetail").([]interface{})
			HouseLabel = locd[0].(tk.M).GetString("LocationName")
			// for _, ea := range locd{
			// 	HouseLabel = ea.(tk.M).GetString("LocationName"))
			// }

			a1 := pdf.GetY()
			x += widthHead[1]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[2], 5, HouseLabel+" ("+each.GetString("StoreHouseFrom")+")", "", "C", false)

			locd = each.Get("ToDetail").([]interface{})
			HouseLabel = locd[0].(tk.M).GetString("LocationName")

			a2 := pdf.GetY()
			x += widthHead[2]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[3], 5, HouseLabel+" ("+each.GetString("StoreHouseTo")+")", "", "C", false)

			a3 := pdf.GetY()
			x += widthHead[3]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[4], 5, each.GetString("DateStr"), "", "C", false)
			a4 := pdf.GetY()
			x += widthHead[4]
			pdf.SetXY(x, y1)
			pdf.MultiCell(widthHead[5], 5, each.GetString("Description"), "", "C", false)

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
	}

	e = os.RemoveAll(c.PdfPath + "/report/pdf")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-reportTO.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf"
	e = pdf.OutputFileAndClose(location + "/" + fileName)
	if e != nil {
		return c.SetResultFile(true, e.Error(), "")
	}
	return c.SetResultFile(false, "success", fileName)
}

func (c *TransferOrderController) ExportDetailReportAllTS(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// Set the base filename
	outputFilename := "transfershipmentdetail"

	// Get payload
	py := struct {
		DateStart     string `json: "DateStart", bson: "DateStart"`
		DateEnd       string `json: "DateEnd", bson: "DateEnd"`
		StoreHouse    string `json: "StoreHouse", bson: "StoreHouse"`
		StoreHouseVal int    `json: "StoreHouseVal", bson: "StoreHouseVal"`
		TransferType  string `json: "TransferType", bson: "TransferType"`
	}{}

	e := k.GetPayload(&py)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", py.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", py.DateEnd)

	mydata, e := c.GetExcelData(dateStart, dateEnd, py.StoreHouseVal, "TS")
	// tk.Println(tk.JsonStringIndent(mydata, "/n >>>>>>"))
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	//Begin Create PDF
	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	dataSet := []TransferShipmentModel{}
	for _, each := range mydata {
		// Retrieve data
		dataSet = []TransferShipmentModel{}
		// crs, e := c.Ctx.Connection.NewQuery().From("TransferShipment").Select().Where(db.Eq("_id", bson.ObjectIdHex("5c86148b3152bf4d4c82a65f"))).Cursor(nil)
		crs, e := c.Ctx.Connection.NewQuery().From("TransferShipment").Select().Where(db.Eq("_id", each.Get("_id"))).Cursor(nil)
		defer crs.Close()
		e = crs.Fetch(&dataSet, 0, false)
		if e != nil {
			return c.SetResultInfo(true, "error", e)
		}

		data := dataSet[0]
		// Retrieve user
		csr, e := c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(db.Eq("username", data.CreatedBy)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		User := SysUserModel{}
		e = csr.Fetch(&User, 1, false)
		if e != nil {
			// tk.Println("fetch", e.Error())
		}
		// Retrieve location from
		csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(db.Eq("LocationID", data.StoreHouseFrom)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		LocFrom := LocationModel{}
		e = csr.Fetch(&LocFrom, 1, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}
		// Retrieve location to
		csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(db.Eq("LocationID", data.StoreHouseTo)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		LocTo := LocationModel{}
		e = csr.Fetch(&LocTo, 1, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}

		// Create pdf object

		pdf.AddPage()

		// ==== BEGIN WRITE PDF ====
		// > Write Header
		var x, y float64
		pdf.SetFont("Century_Gothic", "", 12)
		pdf.CellFormat(0, 12, "TRANSFER SHIPMENT ORDER", "", 0, "C", false, 0, "")
		pdf.Ln(1)
		y = pdf.GetY() + 5
		pdf.SetY(y)
		pdf.SetFont("Century_Gothic", "", 8)
		pdf.CellFormat(0, 8, data.DocumentNumberShipment, "", 0, "C", false, 0, "")
		pdf.Ln(1)
		y = pdf.GetY() + 6
		pdf.SetY(y)
		x = pdf.GetX()
		pdf.CellFormat(0, 8, "FROM", "", 0, "L", false, 0, "")
		x += 16
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, ":   "+LocFrom.LocationName, "", 0, "L", false, 0, "")
		x = 200
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, "Date", "", 0, "L", false, 0, "")
		x += 22
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, ":   "+data.DateStr, "", 0, "L", false, 0, "")
		pdf.Ln(4)
		x = 200
		pdf.SetX(x)
		pdf.CellFormat(0, 8, "DELIVER TO", "", 0, "L", false, 0, "")
		x += 22
		pdf.SetX(x)
		pdf.CellFormat(0, 8, ":   "+LocTo.LocationName, "", 0, "L", false, 0, "")

		pdf.Ln(6)
		x = 226
		pdf.SetX(x)
		pdf.MultiCell(0, 3, LocTo.Description, "", "L", false)
		pdf.Ln(1)
		y = pdf.GetY() + 5
		pdf.SetY(y)
		x = pdf.GetX()

		// > Write Table
		tabHead := []string{"No.", "Code Item", "Item", "Qty"}
		tabAlign := []string{"L", "L", "L", "C"}
		tabWidth := []float64{15, 85, 135, 35}
		y0 := y
		sumQty := 0
		for ci, th := range tabHead {
			tw := tabWidth[ci]
			cellAlign := tabAlign[ci]
			pdf.SetXY(x, y0)
			pdf.SetFont("Century_Gothic", "", 8)

			// Draw Header
			pdf.MultiCell(tw, 6, th, "TB", cellAlign, false)

			// Draw Spacer
			pdf.MultiCell(tw, 2, "", "", "L", false)

			pdf.SetFont("Century_Gothic", "", 7)
			// Draw Data
			for _, d := range data.ListDetailTransferShipment {
				cellData := ""
				switch ci {
				case 0:
					cellData = strconv.Itoa(ci + 1)
				case 1:
					cellData = d.CodeItem
				case 2:
					cellData = d.Item
				case 3:
					cellData = strconv.Itoa(d.Qty)
					sumQty += d.Qty
				}
				pdf.SetX(x)
				pdf.MultiCell(tw, 5, cellData, "", cellAlign, false)
			}

			y = pdf.GetY()
			if y < 85 {
				y = 85
			}
			pdf.SetXY(x, y)
			// Draw Spacer & Bottom Line
			pdf.MultiCell(tw, 2, "", "B", "C", false)

			x += tw
		}
		pdf.Ln(1)
		y = pdf.GetY()
		pdf.SetY(y)

		// Draw Remarks and Total Qty
		pdf.SetFont("Century_Gothic", "", 8)
		pdf.Cell(tabWidth[0]+tabWidth[1]+tabWidth[2]-20, 7, "Remarks: "+data.Description)
		pdf.CellFormat(20, 7, "Total Quantity: ", "", 0, "R", false, 0, "")
		pdf.CellFormat(tabWidth[3], 7, strconv.Itoa(sumQty), "", 0, "C", false, 0, "")
		pdf.Ln(10)
		y = pdf.GetY() - 4
		x = pdf.GetX()

		// > Write Sign Forms
		sgnWidth := []float64{65, 65, 65, 65}
		sgnNames := []string{User.Fullname, "", "", ""}
		sgnTitle := []string{"Prepared by", "Approved by", "Delivered by", "Received by"}
		y0 = y
		for ci, sw := range sgnWidth {
			sn := sgnNames[ci]
			st := sgnTitle[ci]
			if sn == "" {
				sn = "                   "
			}
			sn = "(" + sn + ")"
			pdf.SetXY(x, y0)

			pdf.CellFormat(sw, 25, st, "", 0, "C", false, 0, "")
			pdf.SetXY(x, y0+25)
			pdf.CellFormat(sw, 7, sn, "", 0, "C", false, 0, "")

			x += sw
		}
		pdf.Ln(1)
		x = pdf.GetX()
		y = pdf.GetY() + 8
		pdf.SetXY(x, y)

		// > Write Print Time
		pdf.SetFont("Century_Gothic", "", 6)
		pdf.CellFormat(0, 7, "Print Date: "+time.Now().Format(" January 02, 2006"), "", 0, "L", false, 0, "")

		// ==== END WRITE PDF ======
	}

	// Delete temp file, writes the pdf, then output the filename
	fileName := time.Now().Format("2006-01-02T150405") + "-" + outputFilename + ".pdf"
	outputPath := c.PdfPath + "/report/pdf/"
	DeleteTemporaryFiles(outputPath, "5m", true)
	pdfError := pdf.OutputFileAndClose(outputPath + fileName)
	if pdfError != nil {
		return c.SetResultInfo(true, "error", pdfError)
	}
	return c.SetResultFile(false, "success", fileName)
}

func (c *TransferOrderController) ExportDetailReportAllTR(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// Set the base filename
	outputFilename := "transferreceiptdetail"

	// Get payload
	py := struct {
		DateStart     string `json: "DateStart", bson: "DateStart"`
		DateEnd       string `json: "DateEnd", bson: "DateEnd"`
		StoreHouse    string `json: "StoreHouse", bson: "StoreHouse"`
		StoreHouseVal int    `json: "StoreHouseVal", bson: "StoreHouseVal"`
		TransferType  string `json: "TransferType", bson: "TransferType"`
	}{}

	e := k.GetPayload(&py)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", py.DateStart)
	dateEnd, _ := time.Parse("2006-01-02", py.DateEnd)

	mydata, e := c.GetExcelData(dateStart, dateEnd, py.StoreHouseVal, "TR")
	// tk.Println(tk.JsonStringIndent(mydata, "/n >>>>>>"))
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	//Begin Create PDF
	pdf := gofpdf.New("L", "mm", "A4", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	dataSet := []TransferReceiptModel{}
	for _, each := range mydata {
		// Retrieve data
		dataSet = []TransferReceiptModel{}
		// crs, e := c.Ctx.Connection.NewQuery().From("TransferShipment").Select().Where(db.Eq("_id", bson.ObjectIdHex("5c86148b3152bf4d4c82a65f"))).Cursor(nil)
		crs, e := c.Ctx.Connection.NewQuery().From("TransferReceipt").Select().Where(db.Eq("_id", each.Get("_id"))).Cursor(nil)
		defer crs.Close()
		e = crs.Fetch(&dataSet, 0, false)
		if e != nil {
			return c.SetResultInfo(true, "error", e)
		}

		data := dataSet[0]
		// Retrieve user

		csr, e := c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(db.Eq("username", data.CreatedBy)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		User := SysUserModel{}
		e = csr.Fetch(&User, 1, false)
		if e != nil {
			// tk.Println("fetch", e.Error())
		}
		// Retrieve location from
		csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(db.Eq("LocationID", data.StoreHouseFrom)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		LocFrom := LocationModel{}
		e = csr.Fetch(&LocFrom, 1, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}
		// Retrieve location to
		csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(db.Eq("LocationID", data.StoreHouseTo)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		LocTo := LocationModel{}
		e = csr.Fetch(&LocTo, 1, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}

		// Create pdf object

		pdf.AddPage()

		// ==== BEGIN WRITE PDF ====
		// > Write Header
		var x, y float64
		pdf.SetFont("Century_Gothic", "", 12)
		pdf.CellFormat(0, 12, "TRANSFER RECEIPT ORDER", "", 0, "C", false, 0, "")
		pdf.Ln(1)
		y = pdf.GetY() + 5
		pdf.SetY(y)
		pdf.SetFont("Century_Gothic", "", 8)
		pdf.CellFormat(0, 8, data.DocumentNumberReceipt, "", 0, "C", false, 0, "")
		pdf.Ln(1)
		y = pdf.GetY() + 6
		pdf.SetY(y)
		x = pdf.GetX()
		pdf.CellFormat(0, 8, "FROM", "", 0, "L", false, 0, "")
		x += 16
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, ":   "+LocFrom.LocationName, "", 0, "L", false, 0, "")
		x = 200
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, "Date", "", 0, "L", false, 0, "")
		x += 22
		pdf.SetXY(x, y)
		pdf.CellFormat(0, 8, ":   "+data.DateStr, "", 0, "L", false, 0, "")
		pdf.Ln(4)
		x = pdf.GetX()
		pdf.CellFormat(0, 8, "TS No.", "", 0, "L", false, 0, "")
		x += 16
		pdf.SetX(x)
		pdf.CellFormat(0, 8, ":   "+data.DocumentNumberShipment, "", 0, "L", false, 0, "")
		x = 200
		pdf.SetX(x)
		pdf.CellFormat(0, 8, "DELIVER TO", "", 0, "L", false, 0, "")
		x += 22
		pdf.SetX(x)
		pdf.CellFormat(0, 8, ":   "+LocTo.LocationName, "", 0, "L", false, 0, "")

		pdf.Ln(6)
		x = 226
		pdf.SetX(x)
		pdf.MultiCell(0, 3, LocTo.Description, "", "L", false)
		pdf.Ln(1)
		y = pdf.GetY() + 5
		pdf.SetY(y)
		x = pdf.GetX()

		// > Write Table
		tabHead := []string{"No.", "Code Item", "Item", "Qty"}
		tabAlign := []string{"L", "L", "L", "C"}
		tabWidth := []float64{15, 85, 135, 35}
		y0 := y
		sumQty := 0
		for ci, th := range tabHead {
			tw := tabWidth[ci]
			cellAlign := tabAlign[ci]
			pdf.SetXY(x, y0)
			pdf.SetFont("Century_Gothic", "", 8)

			// Draw Header
			pdf.MultiCell(tw, 6, th, "TB", cellAlign, false)

			// Draw Spacer
			pdf.MultiCell(tw, 2, "", "", "L", false)

			pdf.SetFont("Century_Gothic", "", 7)
			// Draw Data
			for _, d := range data.ListDetailTransferReceipt {
				cellData := ""
				switch ci {
				case 0:
					cellData = strconv.Itoa(ci + 1)
				case 1:
					cellData = d.CodeItem
				case 2:
					cellData = d.Item
				case 3:
					cellData = strconv.Itoa(d.Qty)
					sumQty += d.Qty
				}
				pdf.SetX(x)
				pdf.MultiCell(tw, 5, cellData, "", cellAlign, false)
			}

			y = pdf.GetY()
			if y < 85 {
				y = 85
			}
			pdf.SetXY(x, y)
			// Draw Spacer & Bottom Line
			pdf.MultiCell(tw, 2, "", "B", "C", false)

			x += tw
		}
		pdf.Ln(1)
		y = pdf.GetY()
		pdf.SetY(y)

		// Draw Remarks and Total Qty
		pdf.SetFont("Century_Gothic", "", 8)
		pdf.Cell(tabWidth[0]+tabWidth[1]+tabWidth[2]-20, 7, "Remarks: "+data.Description)
		pdf.CellFormat(20, 7, "Total Quantity: ", "", 0, "R", false, 0, "")
		pdf.CellFormat(tabWidth[3], 7, strconv.Itoa(sumQty), "", 0, "C", false, 0, "")
		pdf.Ln(10)
		y = pdf.GetY() - 4
		x = pdf.GetX()

		// > Write Sign Forms
		sgnWidth := []float64{40, 40, 40, 40}
		sgnNames := []string{User.Fullname, "", "", ""}
		sgnTitle := []string{"Prepared by", "Approved by", "Delivered by", "Received by"}
		y0 = y
		for ci, sw := range sgnWidth {
			sn := sgnNames[ci]
			st := sgnTitle[ci]
			if sn == "" {
				sn = "                   "
			}
			sn = "(" + sn + ")"
			pdf.SetXY(x, y0)

			pdf.CellFormat(sw, 25, st, "", 0, "C", false, 0, "")
			pdf.SetXY(x, y0+25)
			pdf.CellFormat(sw, 7, sn, "", 0, "C", false, 0, "")

			x += sw
		}
		pdf.Ln(1)
		x = pdf.GetX()
		y = pdf.GetY() + 8
		pdf.SetXY(x, y)

		// > Write Print Time
		pdf.SetFont("Century_Gothic", "", 6)
		pdf.CellFormat(0, 7, "Print Date: "+time.Now().Format(" January 02, 2006"), "", 0, "L", false, 0, "")

		// ==== END WRITE PDF ======
	}

	// Delete temp file, writes the pdf, then output the filename
	fileName := time.Now().Format("2006-01-02T150405") + "-" + outputFilename + ".pdf"
	outputPath := c.PdfPath + "/report/pdf/"
	DeleteTemporaryFiles(outputPath, "5m", true)
	pdfError := pdf.OutputFileAndClose(outputPath + fileName)
	if pdfError != nil {
		return c.SetResultInfo(true, "error", pdfError)
	}
	return c.SetResultFile(false, "success", fileName)
}

func (c *TransferOrderController) GetExcelData(datestart time.Time, dateend time.Time, location int, typeData string) ([]tk.M, error) {
	tableName := ""

	filter := []*db.Filter{}
	dt := datestart.Truncate(24 * time.Hour)
	filter = append(filter, db.Gte("DatePosting", dt))
	dt = dateend.Truncate(24 * time.Hour).Add(24 * time.Hour)
	filter = append(filter, db.Lt("DatePosting", dt))

	if typeData == "TS" {
		filter = CreateLocationFilter(filter, "StoreHouseFrom", location, false)
		tableName = "TransferShipment"
	} else {
		filter = CreateLocationFilter(filter, "StoreHouseTo", location, false)
		tableName = "TransferReceipt"
	}

	pipe := []tk.M{}

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseFrom", // field in the orders collection
			"foreignField": "LocationID",     // field in the items collection
			"as":           "FromDetail",
		},
	})

	pipe = append(pipe, tk.M{
		"$lookup": tk.M{
			"from":         "Location",
			"localField":   "StoreHouseTo", // field in the orders collection
			"foreignField": "LocationID",   // field in the items collection
			"as":           "ToDetail",
		},
	})

	crs, e := c.Ctx.Connection.NewQuery().Command("pipe", pipe).From(tableName).Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()

	results := []tk.M{}
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}

	return results, e
}
