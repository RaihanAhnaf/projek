package controllers

import (
	. "../helpers"
	"../library/reflection"
	. "../models"
	"fmt"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
)

func (c *TransferOrderController) GetTransferReceipt(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)

	var f = struct {
		DateStart string
		DateEnd   string
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
	filter = CreateLocationFilter(filter, "StoreHouseTo", locid, false)

	crs, e := c.Ctx.Connection.NewQuery().From("TransferReceipt").Select().Where(db.And(filter...)).Cursor(nil)
	defer crs.Close()
	results := make([]TransferReceiptModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransferOrderController) GetTransferShipmentTS(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)

	var f = struct {
		DateStart string
		DateEnd   string
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
	filter = append(filter, db.Eq("Status", "TS"))
	filter = CreateLocationFilter(filter, "StoreHouseTo", locid, false)

	crs, e := c.Ctx.Connection.NewQuery().From("TransferShipment").Select().Where(filter...).Cursor(nil)
	defer crs.Close()
	results := make([]TransferShipmentModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *TransferOrderController) SaveTransferReceipt(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// Save Transaction Data
	var t = struct {
		Data TransferReceiptModel
	}{}
	err := k.GetPayload(&t)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	var data = &t.Data
	var newData = false
	if data.ID == "" {
		newData = true
		data.DocumentNumberReceipt = c.generateNewDocumentNumber(DocumentNumberFormatReceipt, strconv.Itoa(data.StoreHouseTo), "TransferReceipt")
		data.ID = bson.NewObjectId()
		for i := 0; i < len(data.ListDetailTransferReceipt); i++ {
			var detail = &data.ListDetailTransferReceipt[i]
			detail.Id = bson.NewObjectId()

			// Update Inventory
			var result = []InventoryModel{}
			f := []*db.Filter{}
			f = append(f, db.Eq("StoreLocation", data.StoreHouseTo))
			f = append(f, db.Eq("INVID", detail.CodeItem))
			crs, e := c.Ctx.Connection.NewQuery().From("Inventory").Select().Where(f...).Cursor(nil)
			defer crs.Close()
			e = crs.Fetch(&result, 0, false)
			modelLogInv := InventoryModel{}
			if e == nil && len(result) > 0 {
				// Update existing Inventory Record
				model := result[0]
				modelLogInv = model
				// Record Inventory History
				history := ListHistoryInventory{}
				reflection.Map(&model, &history)
				history.Id = bson.NewObjectId()

				// Update saldo and record TRInventory
				model.TRInventory += detail.Qty
				model.Saldo += detail.Qty
				model.LastDate = time.Now()
				model.Total = float64(model.Saldo) * model.UnitCost

				model.ListInventory = append(model.ListInventory, history)
				e = c.Ctx.Save(&model)
				if e != nil {
					fmt.Println(e.Error())
				}
			} else {
				// insert new inventory
				mdlStoreHouse := []LocationModel{}
				crs, e := c.Ctx.Connection.NewQuery().From("Location").Select().Where(db.Eq("LocationID", data.StoreHouseTo)).Cursor(nil)
				defer crs.Close()
				e = crs.Fetch(&mdlStoreHouse, 0, false)
				if e == nil && len(mdlStoreHouse) == 1 {
					var unitCost float64 = 0
					var invType string
					var invUnit string

					f := []*db.Filter{}
					f = append(f, db.Eq("StoreLocation", data.StoreHouseFrom))
					f = append(f, db.Eq("INVID", detail.CodeItem))
					cr2, e := c.Ctx.Connection.NewQuery().From("Inventory").Select().Where(f...).Cursor(nil)
					if e == nil {
						mdlInv := []InventoryModel{}
						defer cr2.Close()
						e = cr2.Fetch(&mdlInv, 0, false)
						if e == nil && len(mdlInv) > 0 {
							unitCost = mdlInv[0].UnitCost
							invType = mdlInv[0].Type
							invUnit = mdlInv[0].Unit
						}
					}

					model := InventoryModel{
						Saldo:             detail.Qty,
						TRInventory:       detail.Qty,
						ID:                bson.NewObjectId(),
						INVID:             detail.CodeItem,
						INVDesc:           detail.Item,
						StoreLocation:     mdlStoreHouse[0].LocationID,
						StoreLocationName: mdlStoreHouse[0].LocationName,
						LastDate:          time.Now(),
						UnitCost:          unitCost,
						Type:              invType,
						Unit:              invUnit,
					}
					model.Total = float64(model.Saldo) * model.UnitCost

					modelLogInv = model
					modelLogInv.Saldo = 0
					e = c.Ctx.Save(&model)
					if e != nil {
						fmt.Println(e.Error())
					}
				}

			}
			//loginventory
			logitem := LogInventoryModel{}
			logitem.Id = bson.NewObjectId()
			logitem.CodeItem = modelLogInv.INVID
			logitem.Item = modelLogInv.INVDesc
			logitem.StorehouseId = modelLogInv.StoreLocation
			logitem.StoreHouseName = modelLogInv.StoreLocationName
			logitem.Date = data.DatePosting
			logitem.Description = data.DocumentNumberReceipt
			logitem.TypeTransaction = "TR"
			logitem.Price = modelLogInv.UnitCost
			logitem.StockUnit = modelLogInv.Saldo
			logitem.CountTransaction = detail.Qty
			logitem.Increment = detail.Qty
			logitem.Decrement = 0
			logitem.TotalSaldo = logitem.StockUnit + logitem.Increment
			e = c.Ctx.Save(&logitem)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
		}

	}
	data.CreatedBy = k.Session("username").(string)

	// Update TS
	crs, e := c.Ctx.Connection.NewQuery().From("TransferShipment").Select().Where(db.Eq("DocumentNumberShipment", data.DocumentNumberShipment)).Cursor(nil)
	defer crs.Close()
	result := []TransferShipmentModel{}
	e = crs.Fetch(&result, 0, false)
	if e == nil && len(result) > 0 {
		model := result[0]
		model.DocumentNumberReceipt = data.DocumentNumberReceipt
		model.Status = "TR"
		e = c.Ctx.Save(&model)
	}

	// Save TR
	err = c.Ctx.Save(data)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	// Save Log
	if newData {
		c.LogActivity("Transaction Order", "Create Transfer Receipt", t.Data.DocumentNumberShipment, k)
	} else {
		c.LogActivity("Transaction Order", "Update Transfer Receipt", t.Data.DocumentNumberShipment, k)
	}

	// Update Saldo Inventory

	return c.SetResultOK(nil)
}

func (c *TransferOrderController) DeleteTransferReceipt(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// Delete Transaction Data
	var t = struct {
		ID bson.ObjectId
	}{}
	err := k.GetPayload(&t)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	var mdl = NewTransferReceipt()
	e := c.Ctx.GetById(mdl, t.ID)
	if e != nil {
		c.WriteLog(e)
	}
	var docNum = mdl.DocumentNumberReceipt
	c.Ctx.Delete(mdl)

	// Save Log
	c.LogActivity("Transaction Order", "Delete Transfer Shipment", docNum, k)

	// Update Saldo Inventory

	return c.SetResultInfo(false, "OK", nil)
}

func (c *TransferOrderController) ExportPdfPerDataTR(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// Set the base filename
	outputFilename := "transferreceiptprint"
	// Get payload
	p := struct {
		Id bson.ObjectId
	}{}
	k.GetPayload(&p)
	// Retrieve data
	dataSet := []TransferReceiptModel{}
	crs, e := c.Ctx.Connection.NewQuery().From("TransferReceipt").Select().Where(db.Eq("_id", p.Id)).Cursor(nil)
	defer crs.Close()
	e = crs.Fetch(&dataSet, 0, false)
	if e != nil {
		return c.SetResultInfo(true, "error", e)
	}
	data := dataSet[0]
	// Retrieve user
	csr, e := c.Ctx.Connection.NewQuery().Select().From("SysUsers").Where(dbox.Eq("username", data.CreatedBy)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	User := SysUserModel{}
	e = csr.Fetch(&User, 1, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	// Retrieve location from
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(dbox.Eq("LocationID", data.StoreHouseFrom)).Cursor(nil)
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
	csr, e = c.Ctx.Connection.NewQuery().Select().From("Location").Where(dbox.Eq("LocationID", data.StoreHouseTo)).Cursor(nil)
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
	pdf := gofpdf.New("L", "mm", "A5", c.FontPath)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
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
	x = 130
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
	x = 130
	pdf.SetX(x)
	pdf.CellFormat(0, 8, "DELIVER TO", "", 0, "L", false, 0, "")
	x += 22
	pdf.SetX(x)
	pdf.CellFormat(0, 8, ":   "+LocTo.LocationName, "", 0, "L", false, 0, "")

	pdf.Ln(6)
	x = 155
	pdf.SetX(x)
	pdf.MultiCell(0, 3, LocTo.Description, "", "L", false)
	pdf.Ln(1)
	y = pdf.GetY() + 5
	pdf.SetY(y)
	x = pdf.GetX()

	// > Write Table
	tabHead := []string{"No.", "Code Item", "Item", "Qty"}
	tabAlign := []string{"L", "L", "L", "C"}
	tabWidth := []float64{10, 40, 125, 15}
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
	pdf.Ln(1)
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
