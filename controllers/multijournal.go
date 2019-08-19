package controllers

import (
	. "../helpers"
	. "../models"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

func (c *TransactionController) MultiJournal(k *knot.WebContext) interface{} {
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

// func (c *TransactionController) SaveMultiJournal(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	e := k.Request.ParseMultipartForm(1000000000)
// 	if e != nil {
// 		c.ErrorResultInfo(e.Error(), nil)
// 	}
// 	payload := new(tk.M)
// 	_, formData, err := k.GetPayloadMultipart(payload)
// 	if err != nil {
// 		return c.ErrorResultInfo(err.Error(), nil)
// 	}

// 	m := time.Now().UTC().Month()
// 	y := time.Now().UTC().Year()
// 	folder := fmt.Sprintf("%d%02d", y, m)
// 	codejurnal := fmt.Sprintf("%02d%d", m, y)

// 	data := NewMainJournal()
// 	datagm := NewGeneralLedger()
// 	tk.UnjsonFromString(formData["data"][0], data)
// 	tk.UnjsonFromString(formData["data"][0], datagm)

// 	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
// 	number := fmt.Sprintf("%04d", ids)

// 	idx, _ := c.GetNextIdSeq("DocumentNumber", data.Journal_Type, int(m), y)
// 	numberx := fmt.Sprintf("%04d", idx)

// 	idxn, _ := c.GetNextIdSeq("DocumentNumberGL", datagm.Journal_Type, int(m), y)
// 	numberxn := fmt.Sprintf("%04d", idxn)

// 	if data.ID == "" {
// 		data.ID = tk.RandomString(32)
// 		datagm.ID = tk.RandomString(32)
// 	}

// 	if data.IdJournal == "" {
// 		data.IdJournal = "JUR/" + codejurnal + "/" + number
// 		datagm.IdJournal = data.IdJournal
// 	}
// 	data.DateStr = tk.Date2String(data.PostingDate, "dd MMM yyyy")
// 	datagm.DateStr = tk.Date2String(datagm.PostingDate, "dd MMM yyyy")
// 	//Create Directory
// 	baseImagePath := ReadConfig()["uploadpath"]
// 	pathfolder := filepath.Join(baseImagePath, folder)
// 	if _, err = os.Stat(pathfolder); os.IsNotExist(err) {
// 		os.MkdirAll(pathfolder, 0777)
// 	}

// 	tk.Println(time.Now().UTC())
// 	for i, _ := range data.ListDetail {
// 		file, handler, err := k.Request.FormFile("fileUpload" + strconv.Itoa(i))
// 		if file != nil {
// 			tk.Println(file)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			defer file.Close()

// 			tk.RandomString(32)

// 			fileName := tk.RandomString(6) + filepath.Ext(handler.Filename)
// 			filePath := filepath.Join(pathfolder, fileName)
// 			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			defer f.Close()
// 			_, err = io.Copy(f, file)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			data.ListDetail[i].Attachment = fileName
// 			datagm.ListDetail[i].Attachment = fileName

// 		} else {
// 			data.ListDetail[i].Attachment = data.ListDetail[0].Attachment
// 		}
// 		data.ListDetail[i].DocumentNumber = data.ListDetail[i].DocumentNumber[:11] + numberx
// 		datagm.ListDetail[i].DocumentNumber = datagm.ListDetail[i].DocumentNumber[:11] + numberxn
// 	}

// 	err = c.Ctx.Save(data)
// 	if err != nil {
// 		return c.ErrorResultInfo(err.Error(), nil)
// 	}
// 	err = c.Ctx.Save(datagm)
// 	if err != nil {
// 		return c.ErrorResultInfo(err.Error(), nil)
// 	}

// 	c.LogActivity("Multi Journal", "Save Multi Journal", data.IdJournal, k)

// 	return c.SetResultOK(nil)
// }
// func (c *TransactionController) SaveMultiJournal(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	e := k.Request.ParseMultipartForm(1000000000)
// 	if e != nil {
// 		c.ErrorResultInfo(e.Error(), nil)
// 	}
// 	payload := new(tk.M)
// 	_, formData, err := k.GetPayloadMultipart(payload)
// 	if err != nil {
// 		return c.ErrorResultInfo(err.Error(), nil)
// 	}

// 	m := time.Now().UTC().Month()
// 	y := time.Now().UTC().Year()
// 	folder := fmt.Sprintf("%d%02d", y, m)
// 	codejurnal := fmt.Sprintf("%02d%d", m, y)

// 	data := NewMainJournal()
// 	tk.UnjsonFromString(formData["data"][0], data)

// 	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
// 	number := fmt.Sprintf("%04d", ids)

// 	if data.ID == "" {
// 		data.ID = tk.RandomString(32)
// 	}

// 	if data.IdJournal == "" {
// 		data.IdJournal = "JUR/" + codejurnal + "/" + number
// 	}
// 	data.DateStr = tk.Date2String(data.PostingDate, "dd MMM yyyy")

// 	//Create Directory
// 	baseImagePath := ReadConfig()["uploadpath"]
// 	pathfolder := filepath.Join(baseImagePath, folder)
// 	if _, err = os.Stat(pathfolder); os.IsNotExist(err) {
// 		os.MkdirAll(pathfolder, 0777)
// 	}

// 	tk.Println(time.Now().UTC())
// 	for i, _ := range data.ListDetail {
// 		file, handler, err := k.Request.FormFile("fileUpload" + strconv.Itoa(i))
// 		if file != nil {
// 			tk.Println(file)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			defer file.Close()

// 			tk.RandomString(32)

// 			fileName := tk.RandomString(6) + filepath.Ext(handler.Filename)
// 			filePath := filepath.Join(pathfolder, fileName)
// 			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			defer f.Close()
// 			_, err = io.Copy(f, file)
// 			if err != nil {
// 				return c.ErrorResultInfo(err.Error(), nil)
// 			}

// 			data.ListDetail[i].Attachment = fileName

// 		} else {
// 			data.ListDetail[i].Attachment = data.ListDetail[0].Attachment
// 		}
// 	}

// 	err = c.Ctx.Save(data)
// 	if err != nil {
// 		return c.ErrorResultInfo(err.Error(), nil)
// 	}

// 	c.LogActivity("Journal", "Insert Journal", data.IdJournal, k)

// 	return c.SetResultOK(nil)
// }
func (c *TransactionController) SaveMultiJournal(k *knot.WebContext) interface{} {
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

	m := time.Now().UTC().Month()
	y := time.Now().UTC().Year()
	folder := fmt.Sprintf("%d%02d", y, m)
	codejurnal := fmt.Sprintf("%02d%d", m, y)

	data := MainJournal{}
	tk.UnjsonFromString(formData["data"][0], &data)
	// tk.Println("asbyfbasbfaidvbidbcbviucb", data)
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
	fixData := []MainJournal{}
	// tk.Println(time.Now().UTC())
	valDebet := 0.0
	valCredit := 0.0
	splitListDetail := []Journal{}
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

		}
		// tk.Println("===============", splitListDetail)
		// tk.Println("adatsdas", data.ListDetail[i])
		//split data to 2 journal
		splitListDetail = append(splitListDetail, data.ListDetail[i])
		valDebet += data.ListDetail[i].Debet
		valCredit += data.ListDetail[i].Credit
		balance := valDebet - valCredit
		if i > 0 {
			if balance == 0 {
				valDebet = 0.0
				valCredit = 0.0
				newData := data
				newData.ID = tk.RandomString(32)
				newData.ListDetail = splitListDetail
				ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
				number := fmt.Sprintf("%04d", ids)
				newData.IdJournal = "JUR/" + codejurnal + "/" + number
				splitListDetail = []Journal{}
				fixData = append(fixData, newData)
				err = c.Ctx.Save(&newData)
				if err != nil {
					return c.ErrorResultInfo(err.Error(), nil)
				}
			}
		}
	}

	// err = c.Ctx.Save(data)
	// if err != nil {
	// 	return c.ErrorResultInfo(err.Error(), nil)
	// }
	c.LogActivity("Journal", "Insert Journal", data.IdJournal, k)

	return c.SetResultInfo(false, "Success", fixData)
}
