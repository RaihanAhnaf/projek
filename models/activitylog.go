package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ActivityLog struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            bson.ObjectId `bson:"_id",json:"_id"`
	AccessDate    int           `bson:"accessdate",json:"accessdate"`
	AccessTime    time.Time     `bson:"accesstime",json:"accesstime"`
	UserName      string        `bson:"username",json:"username"`
	IpAddress     string        `bson:"ipAddress",json:"ipAddress"`
	PageName      string        `bson:"pagename",json:"pagename"`
	PageUrl       string        `bson:"pageurl",json:"pageurl"`
	Activity      string        `bson:"activity",json:"activity"`
	Desc          string        `bson:"desc",json:"desc"`
}

func NewActivityLogModel() *ActivityLog {
	m := new(ActivityLog)
	m.Id = bson.NewObjectId()
	return m

}
func (u *ActivityLog) TableName() string {
	return "Storelog"
}

func (u *ActivityLog) RecordID() interface{} {
	return u.Id

}
