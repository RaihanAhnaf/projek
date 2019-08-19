package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type TypeStockModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	Code          string
	Name          string
}

func NewTypeStockModel() *TypeStockModel {
	m := new(TypeStockModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *TypeStockModel) TableName() string {
	return "TypeStock"
}

func (u *TypeStockModel) RecordID() interface{} {
	return u.ID

}
