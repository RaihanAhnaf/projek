package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type DepartmentModel struct {
	orm.ModelBase  `bson:"-" json:"-"`
	ID             bson.ObjectId `bson:"_id" json:"_id"`
	DepartmentCode string
	DepartmentName string
}

func NewDepartmentModel() *DepartmentModel {
	m := new(DepartmentModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *DepartmentModel) TableName() string {
	return "Department"
}

func (u *DepartmentModel) RecordID() interface{} {
	return u.ID

}
