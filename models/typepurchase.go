package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type TypePurchaseModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	Code          string
	Name          string
}

func NewTypePurchaseModel() *TypePurchaseModel {
	m := new(TypePurchaseModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *TypePurchaseModel) TableName() string {
	return "TypePurchase"
}

func (u *TypePurchaseModel) RecordID() interface{} {
	return u.ID

}
