package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type AccountBLModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId `json: "_id", bson: "_id"`
	ACC_Code      int
	Account_Name  string
	Type          string
	Active        bool
}

func NewAccountBLModel() *AccountBLModel {
	m := new(AccountBLModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *AccountBLModel) TableName() string {
	return "AccountForBL"
}

func (u *AccountBLModel) RecordID() interface{} {
	return u.ID

}
