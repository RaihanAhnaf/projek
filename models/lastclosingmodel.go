package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LastClosingModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id",json:"_id"`
	LastClosing   time.Time     `bson:"lastclosing",json:"lastclosing"`
	MonthYear     int           `bson:"monthyear",json:"monthyear"`
}

func NewLastClosingModel() *LastClosingModel {
	m := new(LastClosingModel)
	m.Id = bson.NewObjectId()
	return m

}
func (u *LastClosingModel) TableName() string {
	return "Last_Closing"
}

func (u *LastClosingModel) RecordID() interface{} {
	return u.Id

}
