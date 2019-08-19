package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequencePOCMModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequencePOCMModel() *SequencePOCMModel {
	m := new(SequencePOCMModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequencePOCMModel) RecordID() interface{} {
	return e.Id
}

func (m *SequencePOCMModel) TableName() string {
	return "SequencePOCM"
}
