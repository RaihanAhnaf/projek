package models

import (
	"gopkg.in/mgo.v2/bson"
)

type TrialBalanceModel struct {
	ID            bson.ObjectId `json: "_id", bson: "_id"`
	ACC_Code      int
	Account_Name  string
	Debet_Credit  string
	Begining      float64
	Debet         float64
	Credit        float64
	Transaction   float64
	Ending        float64
	Category      string
	Main_Acc_Code int
}

func NewTrialBalanceModel() *TrialBalanceModel {
	m := new(TrialBalanceModel)
	m.ID = bson.NewObjectId()
	return m

}
