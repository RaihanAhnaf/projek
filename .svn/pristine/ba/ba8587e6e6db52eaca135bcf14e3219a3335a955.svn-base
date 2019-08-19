package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type LogInventoryModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	Id               bson.ObjectId `bson:"_id",json:"_id"`
	CodeItem         string
	Item             string
	StorehouseId     int
	StoreHouseName   string
	Date             time.Time
	Description      string
	TypeTransaction  string
	Price            float64
	StockUnit        int
	CountTransaction int
	Increment        int
	Decrement        int
	TotalSaldo       int
}

func NewLogInventoryModel() *LogInventoryModel {
	m := new(LogInventoryModel)
	m.Id = bson.NewObjectId()
	return m
}

func (u *LogInventoryModel) TableName() string {
	return "LogInventory"
}

func (u *LogInventoryModel) RecordID() interface{} {
	return u.Id
}
