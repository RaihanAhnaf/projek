package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type CoaCloseModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id",json:"_id"`
	PeriodeStart  time.Time     `json: "periodestart", bson: "periodestart"`
	PeriodeEnd    time.Time     `json: "periodeend", bson: "periodeend"`
	MonthYear     int           `json: "mothyear", bson: "mothyear"`
	ACC_Code      int           `json: "acc_code", bson: "acc_code"`
	Account_Name  string        `json: "account_name", bson: "account_name"`
	Debet_Credit  string        `json: "debet_credit", bson: "debet_credit"`
	Beginning     float64       `json: "beginning", bson: "beginning"`
	Transaction   float64       `json: "transaction ", bson: "transaction"`
	Ending        float64       `json: "ending", bson: "ending"`
	Category      string        `json: "category", bson: "category"`
	Main_Acc_Code int           `json: "main_acc_code", bson: "main_acc_code"`
}

func NewCoaCloseModel() *CoaCloseModel {
	m := new(CoaCloseModel)
	m.Id = bson.NewObjectId()
	return m

}
func (u *CoaCloseModel) TableName() string {
	return "Coa_Close"
}

func (u *CoaCloseModel) RecordID() interface{} {
	return u.Id

}
