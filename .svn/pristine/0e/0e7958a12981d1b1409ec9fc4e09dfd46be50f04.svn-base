package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type SysUserModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id"`
	Username      string
	Fullname      string
	Enable        bool
	Email         string
	CellularNo    string
	Roles         string
	Password      string
	Potition      string
	TokenID       string
	LastLogin     time.Time
	LocationID    int
	LocationName  string
}

func NewSysUserModel() *SysUserModel {
	m := new(SysUserModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *SysUserModel) RecordID() interface{} {
	return e.Id
}

func (m *SysUserModel) TableName() string {
	return "SysUsers"
}
