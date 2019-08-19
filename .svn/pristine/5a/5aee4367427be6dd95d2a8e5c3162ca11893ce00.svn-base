package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TrackInvoiceModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId         `bson:"_id",json:"_id"`
	DateStr        string                `bson:"DateStr",json:"DateStr"`
	DateCreated    time.Time             `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber string                `bson:"DocumentNumber",json:"DocumentNumber"`
	CustomerCode   string                `bson:"CustomerCode",json:"CustomerCode"`
	CustomerName   string                `bson:"CustomerName",json:"CustomerName"`
	DateINV        time.Time             `bson:"DateINV",json:"DateINV"`
	DateSP         time.Time             `bson:"DateSP",json:"DateSP"`
	Status         string                `bson:"Status",json:"Status"`
	Remark         string                `bson:"Remark",json:"Remark"`
	History        []HistoryTrackInvoice `bson:"History",json:"History"`
	IsInventory    bool
}
type HistoryTrackInvoice struct {
	Id               bson.ObjectId `bson:"_id",json:"_id"`
	DateStr          string        `bson:"DateStr",json:"DateStr"`
	DateCreated      time.Time     `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber   string        `bson:"DocumentNumber",json:"DocumentNumber"`
	CustomerCode     string        `bson:"CustomerCode",json:"CustomerCode"`
	CustomerName     string        `bson:"CustomerName",json:"CustomerName"`
	DateINV          time.Time     `bson:"DateINV",json:"DateINV"`
	DateSP           time.Time     `bson:"DateSP",json:"DateSP"`
	DocumentNumberSP string        `bson:"DocumentNumberSP",json:"DocumentNumberSP"`
	Status           string        `bson:"Status",json:"Status"`
	Remark           string        `bson:"Remark",json:"Remark"`
}

func NewTrackInvoiceModel() *TrackInvoiceModel {
	m := new(TrackInvoiceModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TrackInvoiceModel) TableName() string {
	return "TrackingInvoice"
}

func (u *TrackInvoiceModel) RecordID() interface{} {
	return u.ID

}
