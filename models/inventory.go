package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type InventoryModel struct {
	orm.ModelBase     `bson:"-" json:"-"`
	ID                bson.ObjectId          `bson:"_id" json:"_id"`
	INVID             string                 `bson:"INVID" json:"INVID"`
	INVDesc           string                 `bson:"INVDesc" json:"INVDesc"`
	Unit              string                 `bson:"Unit" json:"Unit"`
	Type              string                 `bson:"Type" json:"Type"`
	Beginning         int                    `bson:"Beginning" json:"Beginning"`
	InInventory       int                    `bson:"InInventory" json:"InInventory"`
	OutInventory      int                    `bson:"OutInventory" json:"OutInventory"`
	CMVInventory      int                    `bson:"CMVInventory" json:"CMVInventory"`
	CMInventory       int                    `bson:"CMInventory" json:"CMInventory"`
	Saldo             int                    `bson:"Saldo" json:"Saldo"`
	TSInventory       int                    `bson:"TSInventory" json:"TSInventory"`
	TRInventory       int                    `bson:"TRInventory" json:"TRInventory"`
	UnitCost          float64                `bson:"UnitCost" json:"UnitCost"`
	Total             float64                `bson:"Total" json:"Total"`
	LastDate          time.Time              `bson:"LastDate" json:"LastDate"`
	ListInventory     []ListHistoryInventory `bson:"ListInventory",json:"ListInventory"`
	StoreLocation     int                    `bson:"StoreLocation" json:"StoreLocation"`
	StoreLocationName string                 `bson:"StoreLocationName" json:"StoreLocationName"`
}

type ListHistoryInventory struct {
	orm.ModelBase     `bson:"-" json:"-"`
	Id                bson.ObjectId `bson:"_id" json:"_id"`
	INVID             string        `bson:"INVID" json:"INVID"`
	INVDesc           string        `bson:"INVDesc" json:"INVDesc"`
	Unit              string        `bson:"Unit" json:"Unit"`
	Type              string        `bson:"Type" json:"Type"`
	Beginning         int           `bson:"Beginning" json:"Beginning"`
	InInventory       int           `bson:"InInventory" json:"InInventory"`
	OutInventory      int           `bson:"OutInventory" json:"OutInventory"`
	CMVInventory      int           `bson:"CMVInventory" json:"CMVInventory"`
	CMInventory       int           `bson:"CMInventory" json:"CMInventory"`
	TSInventory       int           `bson:"TSInventory" json:"TSInventory"`
	TRInventory       int           `bson:"TRInventory" json:"TRInventory"`
	Saldo             int           `bson:"Saldo" json:"Saldo"`
	UnitCost          float64       `bson:"UnitCost" json:"UnitCost"`
	Total             float64       `bson:"Total" json:"Total"`
	LastDate          time.Time     `bson:"LastDate" json:"LastDate"`
	StoreLocation     int           `bson:"StoreLocation" json:"StoreLocation"`
	StoreLocationName string        `bson:"StoreLocationName" json:"StoreLocationName"`
}

func NewInventoryModel() *InventoryModel {
	m := new(InventoryModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *InventoryModel) TableName() string {
	return "Inventory"
}

func (u *InventoryModel) RecordID() interface{} {
	return u.ID

}
