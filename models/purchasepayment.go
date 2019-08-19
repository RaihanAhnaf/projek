package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type PurchasePaymentModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId           `bson:"_id",json:"_id"`
	DateStr        string                  `bson:"DateStr",json:"DateStr"`
	DatePosting    time.Time               `bson:"DatePosting",json:"DatePosting"`
	DocumentNumber string                  `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode   string                  `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName   string                  `bson:"SupplierName",json:"SupplierName"`
	PaymentAccount int                     `bson:"PaymentAccount",json:"PaymentAccount"`
	PaymentName    string                  `bson:"PaymentName",json:"PaymentName"`
	ListDetail     []DetailPurchasePayment `bson:"ListDetail",json:"ListDetail"`
	User           string                  `bson:"User",json:"User"`
	BalanceAll     float64                 `bson:"BalanceAll",json:"BalanceAll"`
	DownPayment    int                     `bson:"DownPayment",json:"DownPayment"`
	Department     string                  `bson:"Department",json:"Department"`
	Attachment     string                  `bson:"Attachment",json:"Attachment"`
	IsInventory    bool                    `bson:"IsInventory",json:"IsInventory"`
}
type DetailPurchasePayment struct {
	Id          string    `bson:"_id",json:"_id"`
	DatePayment time.Time `bson:"DatePayment",json:"DatePayment"`
	PoNumber    string    `bson:"PoNumber",json:"PoNumber"`
	Amount      float64   `bson:"Amount",json:"Amount"`
	AlreadyPaid float64   `bson:"AlreadyPaid",json:"AlreadyPaid"`
	Payment     float64   `bson:"Payment",json:"Payment"`
	Balance     float64   `bson:"Balance",json:"Balance"`
	Pay         bool      `bson:"Pay",json:"Pay"`
}

func NewPurchasePaymentModel() *PurchasePaymentModel {
	m := new(PurchasePaymentModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *PurchasePaymentModel) TableName() string {
	return "PurchasePayment"
}

func (u *PurchasePaymentModel) RecordID() interface{} {
	return u.ID

}
