package controllers

import (
	. "../models"
	"fmt"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"time"
)

func (c *TransactionController) SaveJournalFirstAmount(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := new(tk.M)
	_, formData, err := k.GetPayloadMultipart(payload)
	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	m := time.Now().UTC().Month()
	y := time.Now().UTC().Year()

	codejurnal := fmt.Sprintf("%02d%d", m, y)

	data := NewMainJournal()
	datagm := NewGeneralLedger()
	tk.UnjsonFromString(formData["data"][0], data)
	tk.UnjsonFromString(formData["data"][0], datagm)

	ids, _ := c.GetNextIdSeq("Journal", "", int(m), y)
	number := fmt.Sprintf("%04d", ids)

	idx, _ := c.GetNextIdSeq("DocumentNumber", data.Journal_Type, int(m), y)
	numberx := fmt.Sprintf("%04d", idx)

	idxn, _ := c.GetNextIdSeq("DocumentNumberGL", datagm.Journal_Type, int(m), y)
	numberxn := fmt.Sprintf("%04d", idxn)

	if data.ID == "" {
		data.ID = tk.RandomString(32)
		datagm.ID = tk.RandomString(32)
	}

	if data.IdJournal == "" {
		data.IdJournal = "JUR/" + codejurnal + "/" + number
		datagm.IdJournal = data.IdJournal
	}

	for i, _ := range data.ListDetail {
		data.ListDetail[i].Attachment = "BEGIN"
		data.ListDetail[i].DocumentNumber = data.ListDetail[i].DocumentNumber[:11] + numberx
	}

	for n, _ := range datagm.ListDetail {
		datagm.ListDetail[n].Attachment = "BEGIN"
		datagm.ListDetail[n].DocumentNumber = datagm.ListDetail[n].DocumentNumber[:11] + numberxn
	}

	err = c.Ctx.Save(data)
	errn := c.Ctx.Save(datagm)

	if err != nil {
		return c.ErrorResultInfo(err.Error(), nil)
	}

	if errn != nil {
		return c.ErrorResultInfo(errn.Error(), nil)
	}

	c.LogActivity("First Amount", "Save First Amount", data.IdJournal, k)

	return c.SetResultOK(nil)
}
