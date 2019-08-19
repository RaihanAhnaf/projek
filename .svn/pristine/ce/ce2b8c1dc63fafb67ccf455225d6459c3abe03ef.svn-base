package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequenceCustomerModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	NameCust      string
	Lastnumber    int
	Locked        bool
}

func NewSequenceCustomerModel() *SequenceCustomerModel {
	m := new(SequenceCustomerModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequenceCustomerModel) RecordID() interface{} {
	return e.Id
}

func (m *SequenceCustomerModel) TableName() string {
	return "SequenceCustomer"
}
