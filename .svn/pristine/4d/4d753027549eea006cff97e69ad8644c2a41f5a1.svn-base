package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SequenceSCMModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Collname      string
	TypePo        string
	Lastnumber    int
	Locked        bool
	Month         int
	Year          int
}

func NewSequenceSCMModel() *SequenceSCMModel {
	m := new(SequenceSCMModel)
	m.Id = bson.NewObjectId()
	return m
}
func (e *SequenceSCMModel) RecordID() interface{} {
	return e.Id
}

func (m *SequenceSCMModel) TableName() string {
	return "SequenceSCM"
}
