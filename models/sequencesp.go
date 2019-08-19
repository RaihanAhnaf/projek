package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequencePurchasePaymentModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequencePurchasePaymentModel() *SequencePurchasePaymentModel {
	m := new(SequencePurchasePaymentModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequencePurchasePaymentModel) RecordID() interface{} {
	return e.Id
}

func (m *SequencePurchasePaymentModel) TableName() string {
	return "SequenceSP"
}
