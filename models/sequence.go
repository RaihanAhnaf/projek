package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequenceModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	Typejournal   string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequenceModel() *SequenceModel {
	m := new(SequenceModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequenceModel) RecordID() interface{} {
	return e.Id
}

func (m *SequenceModel) TableName() string {
	return "Sequence"
}
