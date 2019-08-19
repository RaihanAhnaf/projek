package models

import (
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"

	// "gopkg.in/mgo.v2/bson"
	"time"
)

type MainGeneralLedger struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string          `bson:"_id",json:"_id"`
	IdJournal     string          `bson:"IdJournal",json:"IdJournal"`
	PostingDate   time.Time       `bson:"PostingDate",json:"PostingDate"`
	CreateDate    time.Time       `bson:"CreateDate",json:"CreateDate"`
	DateStr       string          `bson:"DateStr",json:"DateStr"`
	User          string          `bson:"User",json:"User"`
	Journal_Type  string          `bson:"Journal_Type",json:"Journal_Type"`
	ListDetail    []GeneralDetail `bson:"ListDetail",json:"ListDetail"`
	Status        string          `bson:"Status",json:"Status"`
	Department    string          `bson:"Department",json:"Department"`
	SalesCode     string          `bson:"SalesCode",json:"SalesCode"`
	SalesName     string          `bson:"SalesName",json:"SalesName"`
}
type GeneralDetail struct {
	Id             string    `bson:"_id",json:"_id"`
	No             int       `bson:"No",json:"No"`
	Journal_Type   string    `bson:"Journal_Type",json:"Journal_Type"`
	PostingDate    time.Time `bson:"PostingDate",json:"PostingDate"`
	DateStr        string    `bson:"DateStr",json:"DateStr"`
	DocumentNumber string    `bson:"DocumentNumber",json:"DocumentNumber"`
	Acc_Code       int       `bson:"Acc_Code",json:"Acc_Code"`
	Acc_Name       string    `bson:"Acc_Name",json:"Acc_Name"`
	Debet          float64   `bson:"Debet",json:"Debet"`
	Credit         float64   `bson:"Credit",json:"Credit"`
	Description    string    `bson:"Description",json:"Description"`
	Attachment     string    `bson:"Attachment",json:"Attachment"`
	User           string    `bson:"User",json:"User"`
	Department     string    `bson:"Department",json:"Department"`
	SalesCode      string    `bson:"SalesCode",json:"SalesCode"`
	SalesName      string    `bson:"SalesName",json:"SalesName"`
}

func NewGeneralLedger() *MainGeneralLedger {
	m := new(MainGeneralLedger)
	m.ID = tk.RandomString(32)
	return m
}

func (u *MainGeneralLedger) TableName() string {
	return "GeneralLedger"
}

func (u *MainGeneralLedger) RecordID() interface{} {
	return u.ID

}
