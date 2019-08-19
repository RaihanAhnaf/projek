package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type AssetDepreciationModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId `bson:"_id" json:"_id"`
	MonthYear     string        `bson:"MonthYear" json:"MonthYear"`
	IdChecbox     string        `bson:"IdChecbox" json:"IdChecbox"`
	Checked       bool          `bson:"Checked" json:"Checked"`
	Date          time.Time     `bson:"Date" json:"Date"`
	DateMonthYear string        `bson:"DateMonthYear" json:"DateMonthYear"`
	Amount        float64       `bson:"Amount" json:"Amount"`
}

func NewAssetDepreciationModel() *AssetDepreciationModel {
	m := new(AssetDepreciationModel)
	m.ID = bson.NewObjectId()
	return m

}
func (u *AssetDepreciationModel) TableName() string {
	return "AssetDepreciation"
}

func (u *AssetDepreciationModel) RecordID() interface{} {
	return u.ID

}
