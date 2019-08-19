package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type HeadGeneralLedgerModel struct {
	orm.ModelBase
	ID           bson.ObjectId `json: "_id", bson: "_id"`
	Acc_Code     int           `bson:"Acc_Code",json:"Acc_Code"`
	Bank_Account string        `bson:"Bank_Account",json:"Bank_Account"`
	Description  string        `bson:"Description",json:"Description"`
	Begining     float64       `bson:"Begining",json:"Begining"`
	Transaction  float64       `bson:"Transaction",json:"Transaction"`
	Ending       float64       `bson:"Ending",json:"Ending"`
}

func NewHeadGeneralLedgerModel() *HeadGeneralLedgerModel {
	m := new(HeadGeneralLedgerModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *HeadGeneralLedgerModel) TableName() string {
	return "HeadGenaralLedger"
}

func (u *HeadGeneralLedgerModel) RecordID() interface{} {
	return u.ID

}
