package controllers

import (
	. "../helpers"
	. "../models"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"

	// tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
	// "strconv"
	"time"
)

type PoInvoiceSummaryController struct {
	*BaseController
}

func (c *PoInvoiceSummaryController) Default(k *knot.WebContext) interface{} {
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

func (c *PoInvoiceSummaryController) GetDataListTrackPO(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart    time.Time
		DateEnd      time.Time
		Filter       bool
		TextSearch   string
		SupplierCode string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	filter := []*db.Filter{}
	filter = append(filter, db.Ne("Status", ""))
	if p.Filter == true {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
		if p.TextSearch != "" {
			filterContain := []*db.Filter{}
			filterContain = append(filterContain, db.Contains("Remark", p.TextSearch))
			filterContain = append(filterContain, db.Contains("SupplierName", p.TextSearch))
			filterContain = append(filterContain, db.Contains("Status", p.TextSearch))
			filterContain = append(filterContain, db.Contains("DocumentNumber", p.TextSearch))
			filter = append(filter, db.Or(filterContain...))
		}
		if p.SupplierCode != "" {
			filter = append(filter, db.Eq("SupplierCode", p.SupplierCode))
		}
	} else {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
	}
	csr, e := c.Ctx.Connection.NewQuery().From("TrackingPurchase").Where(filter...).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrackPurchaseModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "SUCCESS", results)
}
func (c *PoInvoiceSummaryController) GetDataListTrackPOInventory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart    time.Time
		DateEnd      time.Time
		Filter       bool
		TextSearch   string
		SupplierCode string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	filter := []*db.Filter{}
	filter = append(filter, db.Ne("Status", ""))
	if p.Filter == true {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
		if p.TextSearch != "" {
			filterContain := []*db.Filter{}
			filterContain = append(filterContain, db.Contains("Remark", p.TextSearch))
			filterContain = append(filterContain, db.Contains("SupplierName", p.TextSearch))
			filterContain = append(filterContain, db.Contains("Status", p.TextSearch))
			filterContain = append(filterContain, db.Contains("DocumentNumber", p.TextSearch))
			filter = append(filter, db.Or(filterContain...))
		}
		if p.SupplierCode != "" {
			filter = append(filter, db.Eq("SupplierCode", p.SupplierCode))
		}
	} else {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
	}
	csr, e := c.Ctx.Connection.NewQuery().From("TrackingPurchaseInventory").Where(filter...).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrackPurchaseInventoryModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "SUCCESS", results)
}
func (c *PoInvoiceSummaryController) GetDataListTrackINV(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var locid = k.Session("locationid").(int)
	p := struct {
		DateStart    time.Time
		DateEnd      time.Time
		Filter       bool
		TextSearch   string
		CustomerCode string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	dateStart, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	filter := []*db.Filter{}
	filter = append(filter, db.Ne("Status", ""))
	if p.Filter == true {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
		if p.TextSearch != "" {
			filterContain := []*db.Filter{}
			filterContain = append(filterContain, db.Contains("Remark", p.TextSearch))
			filterContain = append(filterContain, db.Contains("CustomerName", p.TextSearch))
			filterContain = append(filterContain, db.Contains("Status", p.TextSearch))
			filterContain = append(filterContain, db.Contains("DocumentNumber", p.TextSearch))
			filter = append(filter, db.Or(filterContain...))
		}
		if p.CustomerCode != "" {
			filter = append(filter, db.Eq("CustomerCode", p.CustomerCode))
		}
		filter = CreateLocationFilter(filter, "StoreLocationId", locid, false)
	} else {
		filter = append(filter, db.Gte("DateCreated", dateStart))
		filter = append(filter, db.Lt("DateCreated", dateEnd))
	}
	csr, e := c.Ctx.Connection.NewQuery().From("TrackingInvoice").Where(filter...).Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []TrackInvoiceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "SUCCESS", results)
}
