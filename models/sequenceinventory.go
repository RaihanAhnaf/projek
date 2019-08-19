package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequenceInventoryModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	Lastnumber    int
	Locked        bool
}

func NewSequenceInventoryModel() *SequenceInventoryModel {
	m := new(SequenceInventoryModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequenceInventoryModel) RecordID() interface{} {
	return e.Id
}

func (m *SequenceInventoryModel) TableName() string {
	return "SequenceInventory"
}
