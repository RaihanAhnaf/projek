package models

import (
	"github.com/eaciit/orm"
	// tk "github.com/eaciit/toolkit"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TransferOrderModelTS struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId           `bson:"_id",json:"_id"`
	DateStr                   string                  `bson:"DateStr",json:"DateStr"`
	DatePosting               time.Time               `bson:"DatePosting",json:"DatePosting"`
	DocumentNumberShipment    string                  `bson:"DocumentNumberShipment",json:"DocumentNumberShipment"`
	DocumentNumberReceipt     string                  `bson:"DocumentNumberReceipt",json:"DocumentNumberReceipt"`
	StoreHouseFrom            int                     `bson:"StoreHouseFrom",json:"StoreHouseFrom"`
	StoreHouseTo              int                     `bson:"StoreHouseTo",json:"StoreHouseTo"`
	Description               string                  `bson:"Description",json:"Description"`
	ListDetailTransfer 		  []DetailTransfer 		  `bson:"ListDetailTransferShipment",json:"ListDetailTransferShipment"`
	FromDetail 		  		  []DetailStoreHouse 	  `bson:"FromDetail",json:"FromDetail"`
	ToDetail 		  		  []DetailStoreHouse 	  `bson:"ToDetail",json:"ToDetail"`
	Status                    string                  `bson:"Status",json:"Status"`
}

type TransferOrderModelTR struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId           `bson:"_id",json:"_id"`
	DateStr                   string                  `bson:"DateStr",json:"DateStr"`
	DatePosting               time.Time               `bson:"DatePosting",json:"DatePosting"`
	DocumentNumberShipment    string                  `bson:"DocumentNumberShipment",json:"DocumentNumberShipment"`
	DocumentNumberReceipt     string                  `bson:"DocumentNumberReceipt",json:"DocumentNumberReceipt"`
	StoreHouseFrom            int                     `bson:"StoreHouseFrom",json:"StoreHouseFrom"`
	StoreHouseTo              int                     `bson:"StoreHouseTo",json:"StoreHouseTo"`
	Description               string                  `bson:"Description",json:"Description"`
	ListDetailTransfer 		  []DetailTransfer 		  `bson:"ListDetailTransferReceipt",json:"ListDetailTransferReceipt"`
	FromDetail 		  		  []DetailStoreHouse 	  `bson:"FromDetail",json:"FromDetail"`
	ToDetail 		  		  []DetailStoreHouse 	  `bson:"ToDetail",json:"ToDetail"`
	Status                    string                  `bson:"Status",json:"Status"`
}

type DetailTransfer struct {
	Id        bson.ObjectId `bson:"_id",json:"_id"`
	CodeItem  string        `bson:"CodeItem",json:"CodeItem"`
	Item      string        `bson:"Item",json:"Item"`
	StockUnit int           `bson:"StockUnit",json:"StockUnit"`
	Qty       int           `bson:"Qty",json:"Qty"`
}

type DetailStoreHouse struct {
	LocationID		int		`bson:"LocationID",json:"LocationID"`
	LocationName	string	`bson:"LocationName",json:"LocationName"`
	Description		string	`bson:"Description",json:"Description"`
}
