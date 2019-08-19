package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type UnitModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	UnitName      string        `bson:"UnitName" json:"UnitName"`
}

func NewUnitModel() *UnitModel {
	m := new(UnitModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *UnitModel) TableName() string {
	return "UnitModel"
}

func (u *UnitModel) RecordID() interface{} {
	return u.ID

}
