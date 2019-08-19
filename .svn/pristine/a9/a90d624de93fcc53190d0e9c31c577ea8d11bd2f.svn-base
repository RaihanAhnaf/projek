package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TrackSalesCreditMemo struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId                 `bson:"_id",json:"_id"`
	DateStr        string                        `bson:"DateStr",json:"DateStr"`
	DateCreated    time.Time                     `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber string                        `bson:"DocumentNumber",json:"DocumentNumber"`
	CustomerCode   string                        `bson:"CustomerCode",json:"CustomerCode"`
	CustomerName   string                        `bson:"CustomerName",json:"CustomerName"`
	DateSCM        time.Time                     `bson:"DateSCM",json:"DateSCM"`
	DateSP         time.Time                     `bson:"DateSP",json:"DateSP"`
	Status         string                        `bson:"Status",json:"Status"`
	Remark         string                        `bson:"Remark",json:"Remark"`
	History        []HistoryTrackSalesCreditMemo `bson:"History",json:"History"`
}
type HistoryTrackSalesCreditMemo struct {
	Id               bson.ObjectId `bson:"_id",json:"_id"`
	DateStr          string        `bson:"DateStr",json:"DateStr"`
	DateCreated      time.Time     `bson:"DateCreated",json:"DateCreated"`
	DocumentNumber   string        `bson:"DocumentNumber",json:"DocumentNumber"`
	CustomerCode     string        `bson:"CustomerCode",json:"CustomerCode"`
	CustomerName     string        `bson:"CustomerName",json:"CustomerName"`
	DateSCM          time.Time     `bson:"DateSCM",json:"DateSCM"`
	DateSP           time.Time     `bson:"DateSP",json:"DateSP"`
	DocumentNumberSP string        `bson:"DocumentNumberSP",json:"DocumentNumberSP"`
	Status           string        `bson:"Status",json:"Status"`
	Remark           string        `bson:"Remark",json:"Remark"`
}

func NewTrackSalesCreditMemo() *TrackSalesCreditMemo {
	m := new(TrackSalesCreditMemo)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TrackSalesCreditMemo) TableName() string {
	return "TrackSalesCreditMemo"
}

func (u *TrackSalesCreditMemo) RecordID() interface{} {
	return u.ID

}
