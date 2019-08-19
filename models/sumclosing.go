package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type SumClosingModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id",json:"_id"`
	PeriodeStart  time.Time     `json: "periodestart", bson: "periodestart"`
	PeriodeEnd    time.Time     `json: "periodeend", bson: "periodeend"`
	Periode       string        `json: "periode", bson: "periode"`
	MonthYear     int           `json: "mothyear", bson: "mothyear"`
	Beginning     float64       `json: "beginning", bson: "beginning"`
	Transaction   float64       `json: "transaction ", bson: "transaction"`
	Ending        float64       `json: "ending", bson: "ending"`
}

func NewSumClosingModel() *SumClosingModel {
	m := new(SumClosingModel)
	m.Id = bson.NewObjectId()
	return m

}
func (u *SumClosingModel) TableName() string {
	return "SumClosingperMonth"
}

func (u *SumClosingModel) RecordID() interface{} {
	return u.Id

}
