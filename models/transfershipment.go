package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TransferShipmentModel struct {
	orm.ModelBase              `bson:"-",json:"-"`
	ID                         bson.ObjectId            `bson:"_id",json:"_id"`
	DateStr                    string                   `bson:"DateStr",json:"DateStr"`
	DatePosting                time.Time                `bson:"DatePosting",json:"DatePosting"`
	DocumentNumberShipment     string                   `bson:"DocumentNumberShipment",json:"DocumentNumberShipment"`
	DocumentNumberReceipt      string                   `bson:"DocumentNumberReceipt",json:"DocumentNumberReceipt"`
	StoreHouseFrom             int                      `bson:"StoreHouseFrom",json:"StoreHouseFrom"`
	StoreHouseTo               int                      `bson:"StoreHouseTo",json:"StoreHouseTo"`
	Description                string                   `bson:"Description",json:"Description"`
	ListDetailTransferShipment []DetailTransferShipment `bson:"ListDetailTransferShipment",json:"ListDetailTransferShipment"`
	Status                     string                   `bson:"Status",json:"Status"`
	CreatedBy                  string                   `bson:"CreatedBy",json:"CreatedBy"`
}

type DetailTransferShipment struct {
	Id        bson.ObjectId `bson:"_id",json:"_id"`
	CodeItem  string        `bson:"CodeItem",json:"CodeItem"`
	Item      string        `bson:"Item",json:"Item"`
	StockUnit int           `bson:"StockUnit",json:"StockUnit"`
	Qty       int           `bson:"Qty",json:"Qty"`
}

func NewTransferShipment() *TransferShipmentModel {
	m := new(TransferShipmentModel)
	m.ID = bson.NewObjectId()
	return m
}

func (u *TransferShipmentModel) TableName() string {
	return "TransferShipment"
}

func (u *TransferShipmentModel) RecordID() interface{} {
	return u.ID

}
