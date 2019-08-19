package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type HistoryAsset struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId      `bson:"_id",json:"_id"`
	MonthYear     string             `bson:"MonthYear",json:"MonthYear"`
	DateCreated   time.Time          `bson:"DateCreated",json:"DateCreated"`
	ListAsset     []ListHistoryAsset `bson:"ListAsset",json:"ListAsset"`
}
type ListHistoryAsset struct {
	orm.ModelBase       `bson:"-" json:"-"`
	Id                  bson.ObjectId `bson:"_id" json:"_id"`
	Description         string        `bson:"Description" json:"Description"`
	Category            string        `bson:"Category" json:"Category"`
	Qty                 int           `bson:"Qty" json:"Qty"`
	Price               float64       `bson:"Price" json:"Price"`
	Total               float64       `bson:"Total" json:"Total"`
	PostingDate         time.Time     `bson:"PostingDate" json:"PostingDate"`
	DatePeriod          time.Time     `bson:"DatePeriod" json:"DatePeriod"`
	SumDepreciation     int           `bson:"SumDepreciation" json:"SumDepreciation"`
	MonthlyDepreciation float64       `bson:"MonthlyDepreciation" json:"MonthlyDepreciation"`
	User                string        `bson:"User" json:"User"`
}

func NewHistoryAsset() *HistoryAsset {
	m := new(HistoryAsset)
	m.ID = bson.NewObjectId()
	return m
}

func (u *HistoryAsset) TableName() string {
	return "HistoryAsset"
}

func (u *HistoryAsset) RecordID() interface{} {
	return u.ID

}
