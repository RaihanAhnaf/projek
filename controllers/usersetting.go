package controllers

import (
	"strings"

	"../helpers"
	. "../models"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	"gopkg.in/mgo.v2/bson"
)

type UserSettingController struct {
	*BaseController
}

func (c *UserSettingController) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate

	k.Config.IncludeFiles = []string{}

	access := c.LoadBase(k)

	if access == nil {
		return nil
	}

	DataAccess := c.GetDataAccess(access)
	return DataAccess
}

func (d *UserSettingController) GetData(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}

	oo := struct {
		Id       string
		Username []interface{}
		Role     []interface{}
		Status   bool
		Take     int
		Skip     int
		Sort     []tk.M
	}{}

	err := r.GetPayload(&oo)
	if err != nil {
		return d.SetResultInfo(true, err.Error(), nil)
	}
	var dbFilter []*db.Filter
	if oo.Id != "" {
		dbFilter = append(dbFilter, db.Eq("_id", bson.ObjectIdHex(oo.Id)))
	} else {
		dbFilter = append(dbFilter, db.Eq("enable", oo.Status))
	}

	if len(oo.Username) != 0 {
		dbFilter = append(dbFilter, db.In("username", oo.Username...))
	}

	if len(oo.Role) != 0 {
		dbFilter = append(dbFilter, db.In("roles", oo.Role...))
	}

	sort := ""
	dir := ""
	if len(oo.Sort) > 0 {
		sort = strings.ToLower(oo.Sort[0].Get("field").(string))
		dir = oo.Sort[0].Get("dir").(string)
	}

	if sort == "" {
		sort = "username"
	}
	if dir != "" && dir != "asc" {
		sort = "-" + sort
	}

	queryTotal := tk.M{}
	query := tk.M{}
	data := make([]SysUserModel, 0)
	total := make([]SysUserModel, 0)
	retModel := tk.M{}
	query.Set("limit", oo.Take)
	query.Set("skip", oo.Skip)
	query.Set("order", []string{sort})

	var ur = r.Session("roles").([]SysRolesModel)
	var hasAllActive = false
	for _, role := range ur {
		if role.Name == "All Active" {
			hasAllActive = true
		}
	}
	if !hasAllActive {
		dbFilter = append(dbFilter, db.Ne("roles", "All Active"))
	}

	if len(dbFilter) > 0 {
		query.Set("where", db.And(dbFilter...))
		queryTotal.Set("where", db.And(dbFilter...))
	}

	crsData, errData := d.Ctx.Find(NewSysUserModel(), query)
	defer crsData.Close()
	if errData != nil {
		return d.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)

	if errData != nil {
		return d.SetResultInfo(true, errData.Error(), nil)
	} else {
		retModel.Set("Records", data)
	}
	crsTotal, errTotal := d.Ctx.Find(NewSysUserModel(), queryTotal)
	defer crsTotal.Close()
	if errTotal != nil {
		return d.SetResultInfo(true, errTotal.Error(), nil)
	}
	errTotal = crsTotal.Fetch(&total, 0, false)

	if errTotal != nil {
		return d.SetResultInfo(true, errTotal.Error(), nil)
	} else {
		retModel.Set("Count", len(total))
	}
	ret.Data = retModel

	return ret
}

func (c *UserSettingController) SaveData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id           string
		UserName     string
		FullName     string
		Enable       bool
		Email        string
		Role         string
		Password     string
		Potition     string
		LocationId   int
		LocationName string
	}{}
	e := k.GetPayload(&p)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	data := new(SysUserModel)
	if p.Id != "" {
		data.Id = bson.ObjectIdHex(p.Id)
	} else {
		data.Id = bson.NewObjectId()
	}

	data.Username = p.UserName
	data.Fullname = p.FullName
	data.Enable = p.Enable
	data.Email = p.Email
	data.Roles = p.Role
	data.Password = helpers.GetMD5Hash(p.Password)
	data.Potition = p.Potition
	data.LocationID = p.LocationId
	data.LocationName = p.LocationName
	err := c.Ctx.Save(data)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	c.LogActivity("UserSetting", "Save User", p.UserName, k)
	return c.SetResultInfo(false, "data has been saved", nil)
}

func (c *UserSettingController) UpdateData(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Id           string
		UserName     string
		FullName     string
		Enable       bool
		Email        string
		Role         string
		Password     string
		Potition     string
		LocationId   int
		LocationName string
	}{}
	e := k.GetPayload(&p)

	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	data := new(SysUserModel)
	if p.Id != "" {
		data.Id = bson.ObjectIdHex(p.Id)
	} else {
		data.Id = bson.NewObjectId()
	}

	data.Username = p.UserName
	data.Fullname = p.FullName
	data.Enable = p.Enable
	data.Email = p.Email
	data.Roles = p.Role
	data.Password = p.Password
	data.Potition = p.Potition
	data.LocationID = p.LocationId
	data.LocationName = p.LocationName
	err := c.Ctx.Save(data)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	c.LogActivity("UserSetting", "Update User", p.UserName, k)
	return c.SetResultInfo(false, "data has been saved", nil)
}

func (c *UserSettingController) ChangePassword(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	param := struct {
		OldPassword     string
		NewPassword     string
		ConfirmPassword string
	}{}

	err := k.GetPayload(&param)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	con, err := PrepareConnection()
	defer con.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	res := []SysUserModel{}

	wher := []*db.Filter{}
	onid := k.Session("userid")
	wher = append(wher, db.Eq("_id", bson.ObjectIdHex(onid.(string))))

	query, err := con.NewQuery().Select().From("SysUsers").Where(wher...).Cursor(nil)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	err = query.Fetch(&res, 0, false)
	defer query.Close()
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	for _, val := range res {
		if val.Password != helpers.GetMD5Hash(param.OldPassword) {
			return c.SetResultInfo(true, "Not Match Old Password", nil)
		} else if param.NewPassword != param.ConfirmPassword {
			return c.SetResultInfo(true, "Not Match New Password", nil)
		} else {
			val.Password = helpers.GetMD5Hash(param.NewPassword)
		}

		onhas := map[string]interface{}{"data": val}
		insert := con.NewQuery().From("SysUsers").SetConfig("multiexec", true).Save()
		err = insert.Exec(onhas)

		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}
		c.LogActivity("UserSetting", "Change Password", val.Username, k)
	}

	return c.SetResultInfo(false, "Data Save Successfully", nil)
}
