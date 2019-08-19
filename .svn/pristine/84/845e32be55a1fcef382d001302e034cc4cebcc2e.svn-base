package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SalesPaymentModel struct {
	orm.ModelBase     `bson:"-",json:"-"`
	ID                bson.ObjectId        `bson:"_id",json:"_id"`
	DateStr           string               `bson:"DateStr",json:"DateStr"`
	DatePosting       time.Time            `bson:"DatePosting",json:"DatePosting"`
	DocumentNumber    string               `bson:"DocumentNumber",json:"DocumentNumber"`
	CustomerCode      string               `bson:"CustomerCode",json:"CustomerCode"`
	CustomerName      string               `bson:"CustomerName",json:"CustomerName"`
	PaymentAccount    int                  `bson:"PaymentAccount",json:"PaymentAccount"`
	PaymentName       string               `bson:"PaymentName",json:"PaymentName"`
	ListDetail        []DetailSalesPayment `bson:"ListDetail",json:"ListDetail"`
	User              string               `bson:"User",json:"User"`
	BalanceAll        float64              `bson:"BalanceAll",json:"BalanceAll"`
	Attachment        string               `bson:"Attachment",json:"Attachment"`
	StoreLocationId   int                  `bson:"StoreLocationId",json:"StoreLocationId"`
	StoreLocationName string               `bson:"StoreLocationName",json:"StoreLocationName"`
}
type DetailSalesPayment struct {
	Id          string    `bson:"_id",json:"_id"`
	DatePayment time.Time `bson:"DatePayment",json:"DatePayment"`
	InvNumber   string    `bson:"InvNumber",json:"InvNumber"`
	Amount      float64   `bson:"Amount",json:"Amount"`
	AlreadyPaid float64   `bson:"AlreadyPaid",json:"AlreadyPaid"`
	Receive     float64   `bson:"Receive",json:"Receive"`
	Balance     float64   `bson:"Balance",json:"Balance"`
	Pay         bool      `bson:"Pay",json:"Pay"`
}

func NewSalesPaymentModel() *SalesPaymentModel {
	m := new(SalesPaymentModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *SalesPaymentModel) TableName() string {
	return "SalesPayment"
}

func (u *SalesPaymentModel) RecordID() interface{} {
	return u.ID

}
