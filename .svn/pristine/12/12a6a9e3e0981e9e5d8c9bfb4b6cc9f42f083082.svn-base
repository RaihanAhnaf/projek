package models

import (
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "gopkg.in/mgo.v2/bson"
)

type CustomerModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             string  `bson:"_id",json:"_id"`
	Kode           string  `bson:"Kode",json:"Kode"`
	Name           string  `bson:"Name",json:"Name"`
	Address        string  `bson:"Address",json:"Address"`
	City           string  `bson:"City",json:"City"`
	NoTelp         string  `bson:"NoTelp",json:"NoTelp"`
	Owner          string  `bson:"Owner",json:"Owner"`
	Bank           string  `bson:"Bank",json:"Bank"`
	AccountNo      string  `bson:"AccountNo",json:"AccountNo"`
	NPWP           string  `bson:"NPWP",json:"NPWP"`
	TrxCode        int     `bson:"TrxCode",json:"TrxCode"`
	VATReg         string  `bson:"VATReg",json:"VATReg"`
	PaymentTerm    string  `bson:"PaymentTerm",json:"PaymentTerm"`
	Email          string  `bson:"Email",json:"Email"`
	Type           string  `bson:"Type",json:"Type"`
	SalesCode      string  `bson:"SalesCode",json:"SalesCode"`
	DepartmentCode string  `bson:"DepartementCode",json:"DepartementCode"`
	Limit          float64 `bson:"Limit",json:"Limit"`
}

func NewCustomerModel() *CustomerModel {
	m := new(CustomerModel)
	m.ID = tk.RandomString(32)
	return m

}
func (u *CustomerModel) TableName() string {
	return "Customer"
}

func (u *CustomerModel) RecordID() interface{} {
	return u.ID

}
