package controllers

import (
	"../helpers"
	. "../helpers"
	"../library/strformat"
	. "../models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

func (c *TransactionController) GetDataInvoiceNonInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		CustomerCode string
		TextSearch   string
		DateStart    string
		DateEnd      string
		Filter       bool
		Status       string
		LocationID   int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Gte("DateCreated", dateStart))
	filter = append(filter, dbox.Lt("DateCreated", dateEnd))
	filter = CreateLocationFilter(filter, "StoreLocationId", p.LocationID, false)

	if p.Filter == true {
		if p.TextSearch != "" {
			// filter = append(filter, dbox.Contains("DocumentNo", p.TextSearch))
			filter = append(filter, dbox.Or(dbox.Contains("StoreLocationName", p.TextSearch), dbox.Contains("DocumentNo", p.TextSearch), dbox.Contains("CustomerName", p.TextSearch), dbox.Contains("Status", p.TextSearch)))
		}
		if p.CustomerCode != "" {
			filter = append(filter, dbox.Eq("CustomerCode", p.CustomerCode))
		}
		if p.Status != "" {
			filter = append(filter, dbox.Eq("Status", p.Status))
		}
	}
	crs, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceNonInv").Where(filter...).Cursor(nil)

	defer crs.Close()
	results := make([]InvoiceNonInvModel, 0)
	tk.Println("Check =>", results)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	if p.Filter && len(results) == 0 {
		return c.SetResultInfo(true, "Please refine your search", nil)
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransactionController) InsertInvoiceNonInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	t := struct {
		Data InvoiceNonInvModel
	}{}

	err := k.GetPayload(&t)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}
	p := t.Data
	// ==== Save Invoice BEGIN ====
	model := t.Data
	model.DateStr = p.DateStr
	newDate, _ := time.Parse("02-Jan-2006", model.DateStr)
	model.DateCreated = newDate.Add(time.Hour*time.Duration(time.Now().Hour()) +
		time.Minute*time.Duration(time.Now().Minute()) +
		time.Second*time.Duration(time.Now().Second()))
	LastNumber := c.GetLastNumberNoninv(model.DateCreated, p.StoreLocationId)
	if p.Id == "" {
		model.Id = bson.NewObjectId()
		err = c.SaveLastNumberInvoiceNon(LastNumber, model.DateCreated, p.StoreLocationId)
		if err != nil {
			return c.ErrorResultInfo(err.Error(), nil)
		}
		model.DocumentNo = c.SetDocumentNumberInvoice(LastNumber, model.DateCreated, p.StoreLocationId, "Non Inventory")
	} else {
		model.Id = p.Id
	}
	model.User = k.Session("username").(string)
	if p.ListItem[0].ID == "" {
		for key, _ := range p.ListItem {

			idList := strconv.Itoa(key) + model.DocumentNo
			p.ListItem[key].ID = idList

		}
	}
	model.ListItem = p.ListItem
	// tk.Println("Model =>", model)
	c.Ctx.Save(&model)
	// ==== Save Invoice END ====

	// ==== POSTING BEGIN ====
	if p.Status == "POSTING" {
		// ==== Save History BEGIN ====
		history := HistoryTrackInvoice{}
		history.Id = bson.NewObjectId()
		history.DocumentNumber = model.DocumentNo
		history.DateCreated = p.DateCreated
		history.DateStr = p.DateCreated.Format("2006-01-02")
		history.DateINV = model.DateCreated
		history.Status = "INVOICE"
		history.Remark = model.Description
		history.CustomerCode = model.CustomerCode
		history.CustomerName = model.CustomerName

		crs, e := c.Ctx.Connection.NewQuery().Select().Where(dbox.Eq("DocumentNumber", model.DocumentNo)).From("TrackingInvoice").Cursor(nil)
		defer crs.Close()

		if crs.Count() == 0 {
			inv := NewTrackInvoiceModel()
			inv.ID = bson.NewObjectId()
			inv.DocumentNumber = model.DocumentNo
			inv.DateCreated = p.DateCreated
			inv.DateStr = p.DateCreated.Format("2006-01-02")
			inv.DateINV = model.DateCreated
			inv.Status = "INVOICE"
			inv.Remark = model.Description
			inv.CustomerCode = model.CustomerCode
			inv.CustomerName = model.CustomerName
			inv.History = append(inv.History, history)
			c.Ctx.Save(inv)
		} else {
			resultINV := []TrackInvoiceModel{}
			e = crs.Fetch(&resultINV, 0, false)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
			mod := resultINV[0]
			inv := TrackInvoiceModel{}
			inv.ID = mod.ID
			inv.DocumentNumber = model.DocumentNo
			inv.DateCreated = p.DateCreated
			inv.DateStr = p.DateCreated.Format("2006-01-02")
			inv.DateINV = model.DateCreated
			inv.Status = "INVOICE"
			inv.Remark = model.Description
			inv.CustomerCode = model.CustomerCode
			inv.CustomerName = model.CustomerName
			inv.History = append(inv.History, history)
			c.Ctx.Save(&inv)
		}
		// ==== Save History END ====

		// ==== Save Journal BEGIN ====
		// accountJournal := tk.M{}.Set("debet", 1210).Set("credit", p.AccountCode)
		accountJournal := tk.M{}.Set("debet", 1210).Set("credit", 5110)
		c.SavetoJournalFromInvoice(accountJournal, p.Total, p.VAT, p.GrandTotalIDR, k.Session("username").(string), p.Description, model.DateCreated, p.SalesCode, p.SalesName)
		// ==== Save Journal END ====

		c.LogActivity("Invoice", "Posting Invoice", p.DocumentNo, k)
	} else {
		if p.Id == "" {
			c.LogActivity("Invoice", "Insert Invoice", p.DocumentNo, k)
		} else {
			c.LogActivity("Invoice", "Update Invoice", p.DocumentNo, k)
		}
	}
	// ==== POSTING END ====

	return c.SetResultOK(nil)
}
func (c *TransactionController) ExportToPdfListInvoiceNonInv(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id             bson.ObjectId
		WordGrandtotal string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceNonInv").Where(dbox.Eq("_id", p.Id)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []InvoiceNonInvModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}

	pdf := gofpdf.New("L", "mm", "A5", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	for _, DATA := range results {
		if DATA.Currency == "USD" {
			for i, _ := range DATA.ListItem {
				DATA.ListItem[i].PriceIDR = 0
				DATA.ListItem[i].AmountIDR = 0
			}
			// discount = DATA.Discount
		} else {
			for i, _ := range DATA.ListItem {
				DATA.ListItem[i].PriceUSD = 0
				DATA.ListItem[i].AmountUSD = 0
			}
		}
		csr, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.CustomerCode)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer csr.Close()
		resultsCustomer := []CustomerModel{}
		e = csr.Fetch(&resultsCustomer, 0, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}
		cust := resultsCustomer[0]
		// PurchaseOrderHeadTable := []string{"Item", "Qty", "Price USD", "Price IDR", "Amount USD", "Amount IDR"}
		corp, e := helpers.GetDataCorporateJson()
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		//user
		csr, e := c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(dbox.Eq("username", DATA.User)).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		defer csr.Close()
		User := SysUserModel{}
		e = csr.Fetch(&User, 1, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		c.addPageInvoiceNINV(pdf, &corp, &DATA, &cust, p.WordGrandtotal, &User)
		// Print only once
		break
	}

	e = os.RemoveAll(c.PdfPath + "/invoice")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/invoice", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-invoice.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/invoice"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		tk.Println(e.Error())
	}
	return fileName
}
func (c *TransactionController) DetailReportPdfNonInv(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		CustomerCode string
		TextSearch   string
		DateStart    string
		DateEnd      string
		Filter       bool
		Status       string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	dateStart, _ := time.Parse("2006-01-02", p.DateStart)
	End, _ := time.Parse("2006-01-02", p.DateEnd)
	dateEnd := End.AddDate(0, 0, 1)
	filter := []*dbox.Filter{}

	if p.Filter == true {
		filter = append(filter, dbox.Gte("DateCreated", dateStart))
		filter = append(filter, dbox.Lt("DateCreated", dateEnd))
		if p.TextSearch != "" {
			filter = append(filter, dbox.Contains("DocumentNo", p.TextSearch))
		}
		if p.CustomerCode != "" {
			filter = append(filter, dbox.Eq("CustomerCode", p.CustomerCode))
		}
		if p.Status != "" {
			filter = append(filter, dbox.Eq("Status", p.Status))
		}
	}
	crs, e := c.Ctx.Connection.NewQuery().Select().From("InvoiceNonInv").Where(filter...).Cursor(nil)

	defer crs.Close()
	results := make([]InvoiceNonInvModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}

	pdf := gofpdf.New("L", "mm", "A5", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	for _, DATA := range results {
		if DATA.Currency == "USD" {
			for i, _ := range DATA.ListItem {
				DATA.ListItem[i].PriceIDR = 0
				DATA.ListItem[i].AmountIDR = 0
			}
			// discount = DATA.Discount
		} else {
			for i, _ := range DATA.ListItem {
				DATA.ListItem[i].PriceUSD = 0
				DATA.ListItem[i].AmountUSD = 0
			}
		}
		crs, e = c.Ctx.Connection.NewQuery().Select().From("Customer").Where(dbox.Eq("Kode", DATA.CustomerCode)).Cursor(nil)
		if e != nil {
			tk.Println("query", e.Error())
		}
		defer crs.Close()
		resultsCustomer := []CustomerModel{}
		e = crs.Fetch(&resultsCustomer, 0, false)
		if e != nil {
			tk.Println("fetch", e.Error())
		}
		cust := resultsCustomer[0]
		// PurchaseOrderHeadTable := []string{"Item", "Qty", "Price USD", "Price IDR", "Amount USD", "Amount IDR"}
		corp, e := helpers.GetDataCorporateJson()
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		//user
		csr, e := c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(dbox.Eq("username", DATA.User)).Cursor(nil)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		csr.Close()
		User := SysUserModel{}
		e = csr.Fetch(&User, 1, false)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}
		numEN := strformat.NumeralCreateEnglish()
		wGT := strformat.Capitalize(numEN.ConvertCurrency(DATA.GrandTotalIDR))
		c.addPageInvoiceNINV(pdf, &corp, &DATA, &cust, wGT, &User)
	}

	e = os.RemoveAll(c.PdfPath + "/report/pdf/")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	e = os.MkdirAll(c.PdfPath+"/report/pdf/", 0777)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), "")
	}
	namepdf := "-invoice.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	location := c.PdfPath + "/report/pdf/"
	e = pdf.OutputFileAndClose(location + "/" + fileName)

	if e != nil {
		e.Error()
	}
	return c.SetResultInfo(false, "", fileName)
}
func (c *TransactionController) addPageInvoiceNINV(pdf *gofpdf.Fpdf, corp *CorporateJsonModel, DATA *InvoiceNonInvModel, cust *CustomerModel, WordGrandtotal string, User *SysUserModel) {
	pdf.AddPage()
	x_defaulft := 10.0
	pdf.SetFont("Century_Gothicb", "B", 15)
	pdf.SetXY(x_defaulft, 10)
	pdf.CellFormat(185, 12, "SALES INVOICE", "", 0, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetX(85)

	pdf.SetFont("Century_Gothicb", "B", 10)
	doc := strings.Split(DATA.DocumentNo, "/")
	docnum := doc[0] + "/" + doc[2] + "/" + doc[3]
	pdf.CellFormat(0, 12, docnum, "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothic", "", 8)
	y0 := pdf.GetY() + 5
	//
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Kode, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Date", "", "L", false) // date
	date := DATA.DateCreated.Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+date, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Name, "", "L", false) // customer name
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "Due Date", "", "L", false) // date
	dueDate := DATA.DateCreated.AddDate(0, 0, cust.TrxCode).Format("January 02, 2006")
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+dueDate, "", "L", false) // date
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, "Phone : "+cust.NoTelp, "", "L", false) // phone
	pdf.SetXY(30, y0)
	pdf.MultiCell(0, 3, "", "", "L", false) //phone
	pdf.SetXY(140, y0)
	pdf.MultiCell(20, 3, "DOC No.", "", "L", false) // DocumentNo
	pdf.SetXY(160, y0)
	pdf.MultiCell(40, 3, ": "+docnum, "", "L", false) // DocumentNo
	pdf.Ln(1)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(80, 3, cust.Address, "", "L", false) // address
	pdf.SetXY(140, y0)
	// pdf.MultiCell(20, 3, "Sales", "", "L", false) // sales
	// pdf.SetXY(160, y0)
	// pdf.MultiCell(40, 3, ": "+DATA.SalesName, "", "L", false) // sales
	y0 = pdf.GetY()
	pdf.SetY(y0)
	pdf.Ln(9)
	//
	y0 = pdf.GetY()
	pdf.SetXY(x_defaulft, y0)
	pdf.MultiCell(20, 3, "Rek Bank :", "", "L", false) // rek bank
	pdf.SetXY(30, y0)
	if cust.Bank != "" && cust.AccountNo != "" {
		pdf.MultiCell(60, 3, cust.Bank+"-"+cust.AccountNo, "", "L", false) //rek bank
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
	for i, list := range DATA.ListItem {
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
	pdf.MultiCell(100, 3, DATA.Description, "", "L", false) // Remark
	pdf.Ln(3)
	//
	y0 = pdf.GetY()
	pdf.SetY(y0)
	yTotal := pdf.GetY()
	headBottom := []string{"Prepared by :", "Approved by", "Finance", "Received by :"}
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
	pdf.Ln(10)
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
	xx += widthBottom[2]
	pdf.SetXY(xx, yB)
	pdf.MultiCell(widthBottom[3], 4, "(                          )", "", "C", false)
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
	total := tk.Sprintf("%.2f", DATA.Total)
	pdf.MultiCell(30, 3, c.ConvertToCurrency(total), "", "R", false) // Total
	pdf.Ln(1)

	yTotal = pdf.GetY()
	pdf.SetXY(xtotal, yTotal)
	pdf.MultiCell(20, 3, "Discount ", "", "L", false) // discount
	pdf.SetXY(170, yTotal)
	valuediscount := DATA.Discount / 100 * DATA.Total
	discount := tk.Sprintf("%.2f", valuediscount)
	pdf.MultiCell(30, 3, c.ConvertToCurrency(discount), "", "R", false) // discount
	pdf.Ln(1)

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
	pdf.MultiCell(20, 3, "Payment  :", "", "L", false)
	pdf.SetXY(30, y0)
	pdf.MultiCell(150, 3, corp.BankName+" : "+corp.AccNo+" - "+corp.AccName, "", "L", false) // rekening
	pdf.Ln(2)
	y0 = pdf.GetY()
	pdf.SetXY(10, y0)
	pdf.MultiCell(20, 3, "Print Date :", "", "L", false)
	pdf.SetXY(30, y0)
	datenow := time.Now().Format("January 02, 2006")
	pdf.MultiCell(150, 3, datenow, "", "L", false) // date print

}

func (c *TransactionController) GetLastNumberNoninv(Date time.Time, LocID int) int {
	m := Date.UTC().Month()
	y := Date.UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceINV").Select().Where(dbox.And(dbox.Eq("collname", "invoice"), dbox.Eq("typepo", "invoicenoninv"),
		dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequenceINVModel{}
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
		model := NewSequenceINVModel()
		model.Collname = "invoice"
		model.TypePo = "invoicenoninv"
		model.Lastnumber = 1
		model.Month = int(m)
		model.Year = y
		model.LocationID = LocID
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

func (c *TransactionController) SaveLastNumberInvoiceNon(LastNumber int, Date time.Time, LocID int) error {
	// m := Date.UTC().Month()
	y := Date.UTC().Year()
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceINV").Select().Where(dbox.And(dbox.Eq("collname", "invoice"), dbox.Eq("typepo", "invoicenoninv"),
		dbox.Eq("year", y))).Cursor(nil)

	defer crs.Close()
	result := []SequenceINVModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		return e
	}

	sec := result[0]
	sec.Lastnumber = LastNumber
	c.Ctx.Save(&sec)

	return nil
}

func (c *TransactionController) DeleteInvoiceNon(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id         bson.ObjectId
		DocumentNo string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
	}

	result := new(InvoiceNonInvModel)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		c.WriteLog(e)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("Invoice", "Delete Invoice Non Inventory Draft", p.DocumentNo, k)
	return c.SetResultOK(nil)
}
