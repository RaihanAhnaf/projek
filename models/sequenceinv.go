package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequenceINVModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
	LocationID    int
}

func NewSequenceINVModel() *SequenceINVModel {
	m := new(SequenceINVModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequenceINVModel) RecordID() interface{} {
	return e.Id
}

func (m *SequenceINVModel) TableName() string {
	return "SequenceINV"
}
