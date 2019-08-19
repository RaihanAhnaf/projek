package controllers

import (
	. "../helpers"
	. "../models"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ActivityLogController struct {
	*BaseController
}

func (c *ActivityLogController) Default(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (c *ActivityLogController) GetDataActivityLog(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		DateStart time.Time
		DateEnd   time.Time
		Filter    bool
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"accesstime": tk.M{"$gte": StartTime, "$lt": EndTime}}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Storelog").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}

func (c *ActivityLogController) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id []bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		c.WriteLog(e)
		tk.Println(e.Error())
	}
	result := NewActivityLogModel()
	for i := 0; i < len(p.Id); i++ {
		e = c.Ctx.GetById(result, p.Id[i])
		if e != nil {
			c.WriteLog(e)
		}
		e = c.Ctx.Delete(result)
		c.LogActivity("Activity Log", "Delete Activity", p.Id[i].Hex(), k)
	}
	return c.SetResultInfo(false, "OK", nil)

}
func (c *ActivityLogController) GetDataActivityLogYearly(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		DateStart time.Time
		DateEnd   time.Time
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	StartTime, _ := time.Parse("2006-01-02", p.DateStart.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", p.DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"accesstime": tk.M{"$gte": StartTime, "$lt": EndTime}}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Storelog").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := make([]ActivityLog, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	newResults := []NewActivityLog{}
	for _, each := range results {
		res := NewActivityLog{}
		res.Id = bson.NewObjectId()
		res.Title = each.UserName + ": " + each.Activity
		res.Start = each.AccessTime
		res.End = each.AccessTime
		newResults = append(newResults, res)
	}
	return c.SetResultInfo(false, "Success", newResults)
}

type NewActivityLog struct {
	Id    bson.ObjectId `bson:"id" json:"id"`
	Title string        `bson:"title" json:"title"`
	Start time.Time     `bson:"start" json:"start"`
	End   time.Time     `bson:"end" json:"end"`
}
