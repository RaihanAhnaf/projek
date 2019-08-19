package controllers

import (
	. "../helpers"
	. "../models"
	"strings"

	"gopkg.in/mgo.v2/bson"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type UserProfileController struct {
	*BaseController
}

func (c *UserProfileController) Profile(k *knot.WebContext) interface{} {
	access := c.LoadBase(k)
	if access == nil {
		return nil
	}

	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "_loader2.html", "_head.html"}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}
func (c *UserProfileController) GetProfile(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Username string
		skip     string
		take     string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	csr, e := c.Ctx.Find(new(SysUserModel), tk.M{}.Set("where", db.Eq("username", p.Username)).Set("skip", p.skip).Set("limit", p.take))
	defer csr.Close()
	results := make([]SysUserModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}
	// tk.Println(p.Username, results)
	return CreateResult(true, results, "success")
}
func (c *UserProfileController) SaveProfile(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id       bson.ObjectId `bson:"_id" , json:"Id"`
		Name     string
		Username string
		Password string
		Celluler string
		Email    string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return CreateResult(false, nil, e.Error())
	}

	users := []SysUserModel{}
	csr, e := c.Ctx.Connection.NewQuery().From("SysUsers").Where(
		db.And(
			db.Eq("username", p.Username),
			db.Ne("_id", p.Id),
		),
	).Cursor(nil)
	e = csr.Fetch(&users, 0, false)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(users) > 0 {
		return c.SetResultInfo(true, "", "Username already exists")
	}

	users = []SysUserModel{}
	csr2, e := c.Ctx.Connection.NewQuery().From("SysUsers").Where(db.Eq("_id", p.Id)).Cursor(nil)
	e = csr2.Fetch(&users, 0, false)
	defer csr2.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	if len(users) == 0 {
		return c.SetResultInfo(true, "", nil)
	}

	u := users[0]
	mdl := NewSysUserModel()
	mdl.Id = u.Id
	mdl.Username = p.Username
	mdl.Fullname = p.Name
	mdl.Enable = u.Enable
	mdl.Email = p.Email
	mdl.Roles = u.Roles
	mdl.CellularNo = p.Celluler
	if strings.Trim(p.Password, " \t") != "" {
		mdl.Password = GetMD5Hash(p.Password)
	} else {
		mdl.Password = u.Password
	}
	mdl.LastLogin = u.LastLogin
	mdl.LocationID = u.LocationID
	mdl.LocationName = u.LocationName
	mdl.Potition = u.Potition
	mdl.TokenID = u.TokenID

	e = c.Ctx.Save(mdl)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if p.Id != "" {
		c.LogActivity("UserProfile", "Update Data User Profile", p.Username, k)
	} else {
		c.LogActivity("UserProfile", "Save Data User Profile", p.Username, k)
	}

	return c.SetResultInfo(false, "Data has been saved", nil)
}
