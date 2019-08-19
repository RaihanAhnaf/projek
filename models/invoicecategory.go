package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type InvoiceCategoryModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	CategoryCode  string
	CategoryName  string
}

func NewInvoiceCategoryModel() *InvoiceCategoryModel {
	m := new(InvoiceCategoryModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *InvoiceCategoryModel) TableName() string {
	return "InvoiceCategory"
}

func (u *InvoiceCategoryModel) RecordID() interface{} {
	return u.ID

}
