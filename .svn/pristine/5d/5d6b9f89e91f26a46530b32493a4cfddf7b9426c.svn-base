package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type CategoryModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	Code          string
	Name          string
}

func NewCategoryModel() *CategoryModel {
	m := new(CategoryModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *CategoryModel) TableName() string {
	return "Category"
}

func (u *CategoryModel) RecordID() interface{} {
	return u.ID

}
