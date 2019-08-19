package models

import (
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type UploadStatus struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	Processname   string
	Start         time.Time
	End           time.Time
	Status        string
	Precentage    int
}
