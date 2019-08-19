package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type LocationModel struct {
	orm.ModelBase   `bson:"-" json:"-"`
	ID              bson.ObjectId `bson:"_id" json:"_id"`
	Main_LocationID int           `bson:"Main_LocationID" json:"Main_LocationID"`
	LocationID      int           `bson:"LocationID" json:"LocationID"`
	LocationName    string        `bson:"LocationName" json:"LocationName"`
	Description     string        `bson:"Description" json:"Description"`
}

func NewLocationModel() *LocationModel {
	m := new(LocationModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *LocationModel) TableName() string {
	return "Location"
}

func (u *LocationModel) RecordID() interface{} {
	return u.ID

}
