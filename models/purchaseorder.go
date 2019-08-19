package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type PurchaseOrder struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId    `bson:"_id",json:"_id"`
	DateStr          string           `bson:"DateStr",json:"DateStr"`
	DatePosting      time.Time        `bson:"DatePosting",json:"DatePosting"`
	DocumentNumber   string           `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode     string           `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName     string           `bson:"SupplierName",json:"SupplierName"`
	AccountCode      int              `bson:"AccountCode",json:"AccountCode"`
	Payment          string           `bson:"Payment",json:"Payment"`
	Type             string           `bson:"Type",json:"Type"`
	ListDetail       []DetailItems    `bson:"ListDetail",json:"ListDetail"`
	ListPayment      []DetailPayments `bson:"ListPayment",json:"ListPayment"`
	Status           string           `bson:"Status",json:"Status"`
	TotalIDR         float64          `bson:"TotalIDR",json:"TotalIDR"`
	TotalUSD         float64          `bson:"TotalUSD",json:"TotalUSD"`
	Discount         float64          `bson:"Discount",json:"Discount"`
	VAT              float64          `bson:"VAT",json:"VAT"`
	GrandTotalIDR    float64          `bson:"GrandTotalIDR",json:"GrandTotalIDR"`
	GrandTotalUSD    float64          `bson:"GrandTotalUSD",json:"GrandTotalUSD"`
	Remark           string           `bson:"Remark",json:"Remark"`
	Currency         string           `bson:"Currency",json:"Currency"`
	Rate             float64          `bson:"Rate",json:"Rate"`
	User             string           `bson:"User",json:"User"`
	AlreadyPaid      float64          `bson:"AlreadyPaid",json:"AlreadyPaid"`
	DownPayment      int              `bson:"DownPayment",json:"DownPayment"`
	Department       string           `bson:"Department",json:"Department"`
	LocationID       int              `bson:"LocationID",json:"LocationID"`
	LocationName     string           `bson:"LocationName",json:"LocationName"`
	DatePostingPI    time.Time        `bson:"DatePostingPI",json:"DatePostingPI"`
	DateStrPI        string           `bson:"DateStrPI",json:"DateStrPI"`
	DocumentNumberPI string           `bson:"DocumentNumberPI",json:"DocumentNumberPI"`
}

type DetailItems struct {
	Id        bson.ObjectId `bson:"_id",json:"_id"`
	Item      string        `bson:"Item",json:"Item"`
	Qty       int           `bson:"Qty",json:"Qty"`
	PriceUSD  float64       `bson:"PriceUSD",json:"PriceUSD"`
	PriceIDR  float64       `bson:"PriceIDR",json:"PriceIDR"`
	AmountUSD float64       `bson:"AmountUSD",json:"AmountUSD"`
	AmountIDR float64       `bson:"AmountIDR",json:"AmountIDR"`
}
type DetailPayments struct {
	Id              bson.ObjectId `bson:"_id",json:"_id"`
	DatePayment     time.Time     `bson:"DatePayment",json:"DatePayment"`
	DocumentPayment string        `bson:"DocumentPayment",json:"DocumentPayment"`
	PaymentAmount   float64       `bson:"PaymentAmount",json:"PaymentAmount"`
}

func NewPurchaseOrder() *PurchaseOrder {
	m := new(PurchaseOrder)
	m.ID = bson.NewObjectId()
	return m
}

func (u *PurchaseOrder) TableName() string {
	return "PurchaseOrder"
}

func (u *PurchaseOrder) RecordID() interface{} {
	return u.ID

}
