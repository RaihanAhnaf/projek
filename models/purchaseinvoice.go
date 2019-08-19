package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type PurchaseInvoice struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId        `bson:"_id",json:"_id"`
	DateStr          string               `bson:"DateStr",json:"DateStr"`
	DatePosting      time.Time            `bson:"DatePosting",json:"DatePosting"`
	DocumentNumber   string               `bson:"DocumentNumber",json:"DocumentNumber"`
	PoNumber         string               `bson:"PoNumber",json:"PoNumber"`
	SupplierCode     string               `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName     string               `bson:"SupplierName",json:"SupplierName"`
	AccountCode      string               `bson:"AccountCode",json:"AccountCode"`
	Type             string               `bson:"Type",json:"Type"`
	ListDetail       []DetailItemsInvoice `bson:"ListDetail",json:"ListDetail"`
	Status           string               `bson:"Status",json:"Status"`
	Total            float64              `bson:"Total",json:"Total"`
	VAT              float64              `bson:"VAT",json:"VAT"`
	GrandTotal       float64              `bson:"GrandTotal",json:"GrandTotal"`
	Remark           string               `bson:"Remark",json:"Remark"`
	DownPayment      int                  `bson:"DownPayment",json:"DownPayment"`
	DatePostingPI    time.Time            `bson:"DatePostingPI",json:"DatePostingPI"`
	DateStrPI        string               `bson:"DateStrPI",json:"DateStrPI"`
	DocumentNumberPI string               `bson:"DocumentNumberPI",json:"DocumentNumberPI"`
}

type DetailItemsInvoice struct {
	Id        bson.ObjectId `bson:"_id",json:"_id"`
	Item      string        `bson:"Item",json:"Item"`
	Qty       int           `bson:"Qty",json:"Qty"`
	PriceUSD  float64       `bson:"PriceUSD",json:"PriceUSD"`
	PriceIDR  float64       `bson:"PriceIDR",json:"PriceIDR"`
	AmountUSD float64       `bson:"AmountUSD",json:"AmountUSD"`
	AmountIDR float64       `bson:"AmountIDR",json:"AmountIDR"`
}

func NewPurchaseInvoice() *PurchaseInvoice {
	m := new(PurchaseInvoice)
	m.ID = bson.NewObjectId()
	return m
}

func (u *PurchaseInvoice) TableName() string {
	return "PurchaseInvoice"
}

func (u *PurchaseInvoice) RecordID() interface{} {
	return u.ID

}
