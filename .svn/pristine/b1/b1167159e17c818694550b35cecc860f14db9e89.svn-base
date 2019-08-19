package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type InvoiceModel struct {
	orm.ModelBase     `bson:"-",json:"-"`
	Id                bson.ObjectId  `bson:"_id" , json:"_id"`
	AccountCode       int            `bson:"AccountCode",json:"AccountCode"`
	AccountName       string         `bson:"AccountName",json:"AccountName"`
	DocumentNo        string         `bson:"DocumentNo",json:"DocumentNo"`
	CustomerCode      string         `bson:"CustomerCode",json:"CustomerCode"`
	PoNumber          string         `bson:"PoNumber",json:"PoNumber"`
	CustomerName      string         `bson:"CustomerName",json:"CustomerName"`
	DateStr           string         `bson:"DateStr",json:"DateStr"`
	DateCreated       time.Time      `bson:"DateCreated",json:"DateCreated"`
	ListItem          []Items        `bson:"ListItem",json:"ListItem"`
	ListPayment       []ItemPayments `bson:"ListPayment",json:"ListPayment"`
	Total             float64        `bson:"Total",json:"Total"`
	Discount          float64        `bson:"Discount",json:"Discount"`
	VAT               float64        `bson:"VAT",json:"VAT"`
	GrandTotalIDR     float64        `bson:"GrandTotalIDR",json:"GrandTotalIDR"`
	GrandTotalUSD     float64        `bson:"GrandTotalUSD",json:"GrandTotalUSD"`
	Rate              float64        `bson:"Rate",json:"Rate"`
	Description       string         `bson:"Description",json:"Description"`
	Status            string         `bson:"Status",json:"Status"`
	Currency          string         `bson:"Currency",json:"Currency"`
	User              string         `bson:"User",json:"User"`
	AlreadyPaid       float64        `bson:"AlreadyPaid",json:"AlreadyPaid"`
	SalesCode         string         `bson:"SalesCode",json:"SalesCode"`
	SalesName         string         `bson:"SalesName",json:"SalesName"`
	CreditMemo        bool           `bson:"CreditMemo",json:"CreditMemo"`
	INVCMI            bool           `bson:"INVCMI",json:"INVCMI"`
	StoreLocationId   int            `bson:"StoreLocationId",json:"StoreLocationId"`
	StoreLocationName string         `bson:"StoreLocationName",json:"StoreLocationName"`
}

type Items struct {
	ID        string  `bson:"ID",json:"ID"`
	CodeItem  string  `bson:"CodeItem",json:"CodeItem"`
	Item      string  `bson:"Item",json:"Item"`
	Qty       int     `bson:"Qty",json:"Qty"`
	PriceUSD  float64 `bson:"PriceUSD",json:"PriceUSD"`
	PriceIDR  float64 `bson:"PriceIDR",json:"PriceIDR"`
	AmountUSD float64 `bson:"AmountUSD",json:"AmountUSD"`
	AmountIDR float64 `bson:"AmountIDR",json:"AmountIDR"`
}

type ItemPayments struct {
	Id              bson.ObjectId `bson:"_id",json:"_id"`
	DatePayment     time.Time     `bson:"DatePayment",json:"DatePayment"`
	DocumentPayment string        `bson:"DocumentPayment",json:"DocumentPayment"`
	PaymentAmount   float64       `bson:"PaymentAmount",json:"PaymentAmount"`
}

func NewInvoiceModel() *InvoiceModel {
	m := new(InvoiceModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *InvoiceModel) RecordID() interface{} {
	return e.Id
}

func (m *InvoiceModel) TableName() string {
	return "Invoice"
}
