package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequencePOModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequencePOModel() *SequencePOModel {
	m := new(SequencePOModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequencePOModel) RecordID() interface{} {
	return e.Id
}

func (m *SequencePOModel) TableName() string {
	return "SequencePO"
}
