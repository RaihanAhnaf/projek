package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequencePPModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequencePPModel() *SequencePPModel {
	m := new(SequencePPModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequencePPModel) RecordID() interface{} {
	return e.Id
}

func (m *SequencePPModel) TableName() string {
	return "SequencePP"
}
