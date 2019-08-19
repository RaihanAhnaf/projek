package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TrackPurchaseInventoryModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId                   `bson:"_id",json:"_id"`
	DateStr        string                          `bson:"DateStr",json:"DateStr"`
	DateCreated    time.Time                       `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber string                          `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode   string                          `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName   string                          `bson:"SupplierName",json:"SupplierName"`
	SalesCode      string                          `bson:"SalesCode",json:"SalesCode"`
	DatePO         time.Time                       `bson:"DatePO",json:"DatePO"`
	DatePI         time.Time                       `bson:"DatePI",json:"DatePI"`
	DatePP         time.Time                       `bson:"DatePP",json:"DatePP"`
	Status         string                          `bson:"Status",json:"Status"`
	Remark         string                          `bson:"Remark",json:"Remark"`
	History        []HistoryTrackPurchaseInventory `bson:"History",json:"History"`
}
type HistoryTrackPurchaseInventory struct {
	Id               bson.ObjectId `bson:"_id",json:"_id"`
	DateStr          string        `bson:"DateStr",json:"DateStr"`
	DateCreated      time.Time     `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber   string        `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode     string        `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName     string        `bson:"SupplierName",json:"SupplierName"`
	SalesCode        string        `bson:"SalesCode",json:"SalesCode"`
	DatePO           time.Time     `bson:"DatePO",json:"DatePO"`
	DatePI           time.Time     `bson:"DatePI",json:"DatePI"`
	DatePP           time.Time     `bson:"DatePP",json:"DatePP"`
	DocumentNumberPP string        `bson:"DocumentNumberPP",json:"DocumentNumberPP"`
	Status           string        `bson:"Status",json:"Status"`
	Remark           string        `bson:"Remark",json:"Remark"`
}

func NewTrackPurchaseInventoryModel() *TrackPurchaseInventoryModel {
	m := new(TrackPurchaseInventoryModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TrackPurchaseInventoryModel) TableName() string {
	return "TrackingPurchaseInventory"
}

func (u *TrackPurchaseInventoryModel) RecordID() interface{} {
	return u.ID

}
