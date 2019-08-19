package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type AssetModel struct {
	orm.ModelBase       `bson:"-" json:"-"`
	ID                  bson.ObjectId `bson:"_id" json:"_id"`
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

func NewAssetModel() *AssetModel {
	m := new(AssetModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *AssetModel) TableName() string {
	return "Asset"
}

func (u *AssetModel) RecordID() interface{} {
	return u.ID

}
