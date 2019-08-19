package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SalesModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	SalesID       string        `bson:"SalesID" json:"SalesID"`
	SalesName     string        `bson:"SalesName" json:"SalesName"`
	Phone         string        `bson:"Phone" json:"Phone"`
}

func NewSalesModel() *SalesModel {
	m := new(SalesModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *SalesModel) TableName() string {
	return "Sales"
}

func (u *SalesModel) RecordID() interface{} {
	return u.ID

}
