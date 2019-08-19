package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TransactionOrderModel struct {
	orm.ModelBase   `bson:"-",json:"-"`
	ID              bson.ObjectId `bson:"_id",json:"_id"`
	DateStr         string        `bson:"DateStr",json:"DateStr"`
	DatePosting     time.Time     `bson:"DatePosting",json:"DatePosting"`
	DocumentNumber  string        `bson:"DocumentNumber",json:"DocumentNumber"`
	StoreHouseFrom  string        `bson:"StoreHouseFrom",json:"StoreHouseFrom"`
	StoreHouseTo    string        `bson:"StoreHouseTo",json:"StoreHouseTo"`
	Description     string        `bson:"Description",json:"Description"`
	ListDetailOrder []DetailOrder `bson:"ListDetailOrder",json:"ListDetailOrder"`
}

type DetailOrder struct {
	Id        bson.ObjectId `bson:"_id",json:"_id"`
	CodeItem  string        `bson:"CodeItem",json:"CodeItem"`
	Item      string        `bson:"Item",json:"Item"`
	StockUnit int           `bson:"StockUnit",json:"StockUnit"`
	Qty       int           `bson:"Qty",json:"Qty"`
}

func NewTransactionOrder() *TransactionOrderModel {
	m := new(TransactionOrderModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TransactionOrderModel) TableName() string {
	return "TransactionOrder"
}

func (u *TransactionOrderModel) RecordID() interface{} {
	return u.ID

}
