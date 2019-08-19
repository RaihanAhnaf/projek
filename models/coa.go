package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type CoaModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId `json: "_id", bson: "_id"`
	ACC_Code      int
	Account_Name  string
	Debet_Credit  string
	Debet         float64
	Credit        float64
	Saldo         float64
	Category      string
	Main_Acc_Code int
}

func NewCoaModel() *CoaModel {
	m := new(CoaModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *CoaModel) TableName() string {
	return "Coa"
}

func (u *CoaModel) RecordID() interface{} {
	return u.ID

}
