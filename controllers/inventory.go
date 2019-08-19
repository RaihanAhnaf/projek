package controllers

import (
	. "../helpers"
	"../library/strformat"
	"../library/tealeg/xlsx"
	. "../models"
	"os"
	"path/filepath"
	"strconv"

	"github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"

	"strings"
	"time"

	db "github.com/eaciit/dbox"
)

func (c *MasterController) Inventory(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *MasterController) GetDataSupplier(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	crs, e := c.Ctx.Connection.NewQuery().From("Customer").Select().Where(dbox.Eq("Type", "SUPPLIER")).Cursor(nil)
	defer crs.Close()
	results := make([]CustomerModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		CreateResult(false, nil, e.Error())
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *MasterController) SaveInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Data InventoryModel
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	anName := strformat.Filter(p.Data.INVDesc, strformat.CharsetAlphaNumeric)
	if len(anName) < 3 {
		return c.SetResultInfo(true, "Inventory Name must have at least 3 characters", nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().From("Location").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]LocationModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	for _, fields := range results {
		model := NewInventoryModel()
		model.ID = p.Data.ID
		model.INVID = p.Data.INVID
		if model.ID == "" {
			model.ID = bson.NewObjectId()
			numbStr := ""
			name := strings.ToUpper(anName[:3])
			numb := c.GetLastNumberInventory(name)
			if numb < 10 {
				numbStr = "000" + strconv.Itoa(numb)
			} else if numb <= 10 && numb < 100 {
				numbStr = "00" + strconv.Itoa(numb)
			} else if numb <= 100 && numb < 1000 {
				numbStr = "0" + strconv.Itoa(numb)
			} else {
				numbStr = strconv.Itoa(numb)
			}
			model.INVID = name + "/" + numbStr
			model.LastDate = time.Now()
		}
		model.INVDesc = p.Data.INVDesc
		model.Unit = p.Data.Unit
		model.Type = p.Data.Type
		model.StoreLocation = fields.LocationID
		model.StoreLocationName = fields.LocationName

		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

		e = c.Ctx.Save(model)
		if e != nil {
			return c.SetResultInfo(true, e.Error(), nil)
		}

	}
	c.LogActivity("Inventory", "Save Inventory", p.Data.INVID, k)
	return c.SetResultInfo(false, "OK", "")
}

func (c *MasterController) GetAllDataInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var p = struct {
		LocationID int
	}{}
	k.GetPayload(&p)

	q := c.Ctx.Connection.NewQuery().From("Inventory")
	if p.LocationID != 0 {
		q = q.Where(db.Eq("StoreLocation", p.LocationID))
	}
	csr, e := q.Cursor(nil)
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
func (c *MasterController) GetAllDataInventoryCentral(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := c.Ctx.Connection.NewQuery().Select().From("Inventory").Where(db.Eq("StoreLocation", 1000)).Cursor(nil)
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
func (c *MasterController) GetDataInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		LocationId int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Inventory").Where(db.Eq("StoreLocation", p.LocationId)).Cursor(nil)
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
func (c *MasterController) DeleteInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	result := new(InventoryModel)
	e = c.Ctx.GetById(result, p.Id)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	e = c.Ctx.Delete(result)

	c.LogActivity("Inventory", "Delete Inventory", result.INVID, k)
	return c.SetResultInfo(false, "OK", nil)
}

func (c *MasterController) GetLastNumberInventory(name string) int {
	crs, e := c.Ctx.Connection.NewQuery().From("SequenceInventory").Select().Where(dbox.Eq("collname", name)).Cursor(nil)
	defer crs.Close()

	result := []SequenceInventoryModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
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
		model := NewSequenceInventoryModel()
		model.Collname = name
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

func (c *MasterController) GetDataUnit(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	crs, e := c.Ctx.Connection.NewQuery().From("UnitModel").Select().Cursor(nil)
	defer crs.Close()

	result := []UnitModel{}
	e = crs.Fetch(&result, 0, false)
	if e != nil {
		// tk.Println(e.Error())
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "OK", result)
}

func (c *MasterController) SaveNewUnit(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		UnitName string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(false, e.Error(), nil)
	}

	model := NewUnitModel()
	model.UnitName = p.UnitName
	e = c.Ctx.Save(model)
	c.LogActivity("Inventory", "Save New Unit", model.UnitName, k)
	return c.SetResultInfo(false, "OK", nil)
}

func (c *MasterController) UploadFilesInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filename := k.Request.FormValue("filename")
	result := tk.NewResult()
	pathToSave, _ := filepath.Abs("assets/docs/datamaster")
	os.MkdirAll(pathToSave, 0777)
	e, _ := UploadHandlerCopy(k, "filedoc", pathToSave+"/")
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	fileToProcess := pathToSave + "/" + filename

	_, err := os.Stat(fileToProcess)
	if os.IsNotExist(err) {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	// ==================================TEALEG====================================
	model := new(InventoryModel)
	excelFileName := fileToProcess
	xlFile, er := xlsx.OpenFile(excelFileName)
	if er != nil {
		return c.SetResultInfo(true, "Can't read your excel file. Please download available template, copy-paste your data there, and try import data again!", nil)
		// return c.SetResultInfo(true, er.Error(), nil)
	}

	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			model.ID = bson.NewObjectId()

			model.INVDesc, _ = cells[0].String()

			//Generate Code
			anName := strformat.Filter(model.INVDesc, strformat.CharsetAlphaNumeric)
			if len(anName) < 3 {
				return c.SetResultInfo(true, "Inventory Name must have at least 3 characters", nil)
			}
			numbStr := ""
			name := strings.ToUpper(anName[:3])
			numb := c.GetLastNumberInventory(name)
			if numb < 10 {
				numbStr = "000" + strconv.Itoa(numb)
			} else if numb <= 10 && numb < 100 {
				numbStr = "00" + strconv.Itoa(numb)
			} else if numb <= 100 && numb < 1000 {
				numbStr = "0" + strconv.Itoa(numb)
			} else {
				numbStr = strconv.Itoa(numb)
			}
			model.INVID = name + "/" + numbStr

			_, e0 := cells[1].String()
			if e0 != nil {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data UNIT in row : ", i ,", column : ",2,".\nMake sure the data is filled properly")
				return c.SetResultInfo(true, msg, nil)
			} else{
				model.Unit, _ = cells[1].String()
				if model.Unit == "" {
					msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data UNIT in row : ", i ,", column : ",2,".\nMake sure the data is not blank")
					return c.SetResultInfo(true, msg, nil)
				}
			}

			_, e1 := cells[2].Int()
			if e1 != nil {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data BEGINNING in row : ", i ,", column : ",3,".\nMake sure the data is an integers number")
				return c.SetResultInfo(true, msg, nil)
			} else{
				model.Beginning, _ = cells[2].Int()
			}

			_, e2 := cells[3].Float()
			if e2 != nil {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data PRICE in row : ", i ,", column : ",4,".\nMake sure the data is an integers or fraction number. \n Fraction separator is dot ('.') not comma (',')")
				return c.SetResultInfo(true, msg, nil)
			} else{
				model.UnitCost, _ = cells[3].Float()
			}
			
			stockunit := 0
			stockunit += model.Beginning
			model.Saldo = stockunit
			model.Total = float64(model.Beginning) * model.UnitCost

			
			model.Type, _ = cells[4].String()
			if model.Type == "" {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data TYPE STOCK row : ", i ,", column : ",5,".\nMake sure the data is not blank")
				return c.SetResultInfo(true, msg, nil)
			}
			

			_, e3 := cells[5].Int()
			if e3 != nil {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data LOCATION ID in row : ", i ,", column : ",6,".\nMake sure the data is an integers number")
				return c.SetResultInfo(true, msg, nil)
			} else{
				model.StoreLocation, _ = cells[5].Int()
			}

			_, e4 := cells[6].String()
			if e4 != nil {
				msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data LOCATION NAME in row : ", i ,", column : ",7,".\nMake sure the data is filled properly")
				return c.SetResultInfo(true, msg, nil)
			} else{
				model.StoreLocationName, _ = cells[6].String()
				if model.StoreLocationName == "" {
					msg := tk.Sprintf("%s%d%s%d%s", "Error retrive data LOCATION NAME in row : ", i ,", column : ",7,".\nMake sure the data is not blank")
					return c.SetResultInfo(true, msg, nil)
				}
			}
			
			
			model.LastDate = time.Now()
			csr, e := c.Ctx.Connection.NewQuery().Select().From("Inventory").Where(db.Eq("INVID", model.INVID)).Cursor(nil)
			if e != nil {
				// return result.SetError(e)
				return c.SetResultInfo(true, e.Error(), nil)
			}
			csr.Close()
			if csr.Count() == 0 {
				e = c.Ctx.Save(model)
				if e != nil {
					// return result.SetError(e)
					return c.SetResultInfo(true, e.Error(), nil)
				}
			}

			// //loginventory
			logitem := LogInventoryModel{}
			logitem.Id = bson.NewObjectId()
			logitem.CodeItem = model.INVID
			logitem.Item = model.INVDesc
			logitem.StorehouseId = model.StoreLocation
			logitem.StoreHouseName = model.StoreLocationName
			logitem.Date = model.LastDate
			logitem.Description = "Beginning"
			logitem.TypeTransaction = "BEGINNING"
			logitem.Price = model.UnitCost
			logitem.StockUnit = 0
			logitem.CountTransaction = model.Beginning
			logitem.Increment = model.Beginning

			mysaldo := logitem.StockUnit
			mysaldo += model.Beginning
			logitem.TotalSaldo = mysaldo

			e = c.Ctx.Save(&logitem)
			if e != nil {
				return c.SetResultInfo(true, e.Error(), nil)
			}
		}
	}
	c.LogActivity("Master", "Import Data Excel Inventory", filename, k)
	return result
}

func (c *MasterController) UploadFilesInventoryOLD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	filename := k.Request.FormValue("filename")
	result := tk.NewResult()
	pathToSave, _ := filepath.Abs("assets/docs/datamaster")
	os.MkdirAll(pathToSave, 0777)
	e, _ := UploadHandlerCopy(k, "filedoc", pathToSave+"/")
	if e != nil {
		// tk.Println("Error : " + e.Error())
		return e.Error()
	}

	fileToProcess := pathToSave + "/" + filename

	_, err := os.Stat(fileToProcess)
	if os.IsNotExist(err) {
		// tk.Println(err.Error())
		return result.SetError(err)
	}
	// ==================================TEALEG====================================
	model := new(InventoryModel)
	excelFileName := fileToProcess
	xlFile, er := xlsx.OpenFile(excelFileName)
	if er != nil {
	 	return result.SetError(er)
	}
	for _, sheet := range xlFile.Sheets {
		for i, row := range sheet.Rows {
			if i == 0 {
				continue
			}
			cells := row.Cells
			model.ID = bson.NewObjectId()
			model.INVID, _ = cells[0].String()
			model.INVDesc, _ = cells[1].String()
			model.Unit, _ = cells[2].String()
			model.Beginning, _ = cells[3].Int()
			model.InInventory, _ = cells[4].Int()
			model.OutInventory, _ = cells[5].Int()
			model.Saldo, _ = cells[6].Int()
			model.UnitCost, _ = cells[7].Float()
			model.Total, _ = cells[8].Float()
			datePeriode, _ := cells[9].Float()
			model.Type, _ = cells[10].String()
			model.StoreLocation, _ = cells[11].Int()
			model.StoreLocationName = c.GetLocationName(model.StoreLocation)
			model.LastDate = xlsx.TimeFromExcelTime(datePeriode, false)
			csr, e := c.Ctx.Connection.NewQuery().Select().From("Inventory").Where(db.Eq("INVID", model.INVID)).Cursor(nil)
			if e != nil {
				return result.SetError(e)
			}
			csr.Close()
			if csr.Count() == 0 {
				e = c.Ctx.Save(model)
				if e != nil {
					return result.SetError(e)
				}
			}

			//loginventory
			// logitem := LogInventoryModel{}
			// logitem.Id = bson.NewObjectId()
			// logitem.CodeItem = model.INVID
			// logitem.Item = model.INVDesc
			// logitem.StorehouseId = model.StoreLocation
			// logitem.StoreHouseName = model.StoreLocationName
			// logitem.Date = time.Now()
			// logitem.Description = "Begining"
			// logitem.TypeTransaction = "Begining"
			// logitem.Price = model.UnitCost
			// logitem.StockUnit = saldoBeforeUpdate
			// logitem.CountTransaction = detail.Qty
			// logitem.Increment = 0
			// logitem.Decrement = detail.Qty
			// logitem.TotalSaldo = logitem.StockUnit - logitem.Decrement
			// e = c.Ctx.Save(&logitem)
			// if e != nil {
			// 	return c.SetResultInfo(true, e.Error(), nil)
			// }
		}
	}
	c.LogActivity("Master", "Upload Files", filename, k)
	return result
}

func (c *MasterController) GetLocationName(locid int) string {
	crs, e := c.Ctx.Connection.NewQuery().From("Location").Select().Where(dbox.Eq("LocationID", locid)).Cursor(nil)
	defer crs.Close()
	results := make([]LocationModel, 0)
	e = crs.Fetch(&results, 0, false)
	if e != nil {
		tk.Println(e.Error())
	}

	LocationName := ""
	for _, i := range results {
		LocationName = i.LocationName
	}
	return LocationName
}
