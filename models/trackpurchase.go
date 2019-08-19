package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TrackPurchaseModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId          `bson:"_id",json:"_id"`
	DateStr        string                 `bson:"DateStr",json:"DateStr"`
	DateCreated    time.Time              `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber string                 `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode   string                 `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName   string                 `bson:"SupplierName",json:"SupplierName"`
	DatePO         time.Time              `bson:"DatePO",json:"DatePO"`
	DatePI         time.Time              `bson:"DatePI",json:"DatePI"`
	DatePP         time.Time              `bson:"DatePP",json:"DatePP"`
	Status         string                 `bson:"Status",json:"Status"`
	Remark         string                 `bson:"Remark",json:"Remark"`
	History        []HistoryTrackPurchase `bson:"History",json:"History"`
}
type HistoryTrackPurchase struct {
	Id               bson.ObjectId `bson:"_id",json:"_id"`
	DateStr          string        `bson:"DateStr",json:"DateStr"`
	DateCreated      time.Time     `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber   string        `bson:"DocumentNumber",json:"DocumentNumber"`
	SupplierCode     string        `bson:"SupplierCode",json:"SupplierCode"`
	SupplierName     string        `bson:"SupplierName",json:"SupplierName"`
	DatePO           time.Time     `bson:"DatePO",json:"DatePO"`
	DatePI           time.Time     `bson:"DatePI",json:"DatePI"`
	DatePP           time.Time     `bson:"DatePP",json:"DatePP"`
	DocumentNumberPP string        `bson:"DocumentNumberPP",json:"DocumentNumberPP"`
	Status           string        `bson:"Status",json:"Status"`
	Remark           string        `bson:"Remark",json:"Remark"`
}

func NewTrackPurchaseModel() *TrackPurchaseModel {
	m := new(TrackPurchaseModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TrackPurchaseModel) TableName() string {
	return "TrackingPurchase"
}

func (u *TrackPurchaseModel) RecordID() interface{} {
	return u.ID

}
