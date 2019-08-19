package controllers

import (
	"time"

	. "../models"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	"gopkg.in/mgo.v2/bson"
)

type MenuSettingController struct {
	*BaseController
}

func (c *MenuSettingController) Default(k *knot.WebContext) interface{} {
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

func (c *MenuSettingController) GetMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson

	ret := ResultInfo{}
	payLoad := struct {
		Id string
	}{}

	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	var menuAccess []interface{}

	id := k.Session("rolesid").(string)
	result := new(SysRolesModel)
	_ = c.Ctx.GetById(result, bson.ObjectIdHex(id))
	for _, o := range result.Menu {
		if o.View {
			menuAccess = append(menuAccess, o.Menuid)
		}
	}

	var dbFilter []*db.Filter
	dbFilter = append(dbFilter, db.Eq("Enable", true))
	dbFilter = append(dbFilter, db.In("_id", menuAccess...))

	queryTotal := tk.M{}
	query := tk.M{}
	data := make([]TopMenuModel, 0)
	total := make([]TopMenuModel, 0)
	retModel := tk.M{}

	if len(dbFilter) > 0 {
		query.Set("where", db.And(dbFilter...))
		queryTotal.Set("where", db.And(dbFilter...))
	}

	crsData, errData := c.Ctx.Find(NewTopMenuModel(), query)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}

	errData = crsData.Fetch(&data, 0, false)

	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	} else {
		retModel.Set("Records", data)
	}

	crsTotal, errTotal := c.Ctx.Find(NewTopMenuModel(), queryTotal)
	defer crsTotal.Close()
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	}
	errTotal = crsTotal.Fetch(&total, 0, false)

	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	} else {
		retModel.Set("Count", len(total))
	}
	ret.Data = retModel

	return ret
}

func (c *MenuSettingController) GetAccessMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	url := k.Request.URL.String()
	accesMenu := k.Session("roles").([]SysRolesModel)

	access := []tk.M{}
	for _, o := range accesMenu[0].Menu {

		if o.Url == url {
			obj := tk.M{}
			obj.Set("view", o.View)
			obj.Set("create", o.Create)
			obj.Set("approve", o.Approve)
			obj.Set("delete", o.Delete)
			obj.Set("process", o.Process)
			obj.Set("edit", o.Edit)
			access = append(access, obj)
		}

	}

	return access
}

func (c *MenuSettingController) GetSelectMenu(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	payLoad := struct {
		Id string
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	var dbFilter []*db.Filter
	if payLoad.Id != "" {
		dbFilter = append(dbFilter, db.Eq("_id", payLoad.Id))
	}

	queryTotal := tk.M{}
	query := tk.M{}
	data := make([]TopMenuModel, 0)
	total := make([]TopMenuModel, 0)
	retModel := tk.M{}

	if len(dbFilter) > 0 {
		query.Set("where", db.And(dbFilter...))
		queryTotal.Set("where", db.And(dbFilter...))
	}

	crsData, errData := c.Ctx.Find(NewTopMenuModel(), query)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)

	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	} else {
		retModel.Set("Records", data)
	}
	crsTotal, errTotal := c.Ctx.Find(NewTopMenuModel(), queryTotal)
	defer crsTotal.Close()
	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	}
	errTotal = crsTotal.Fetch(&total, 0, false)

	if errTotal != nil {
		return c.SetResultInfo(true, errTotal.Error(), nil)
	} else {
		retModel.Set("Count", len(total))
	}
	ret.Data = retModel

	return ret
}

type submenu struct {
	Id   string
	Name string
	Menu map[string]string
}

type addMenu struct {
	Url      string
	Access   bool
	Approve  bool
	Create   bool
	Delete   bool
	Edit     bool
	Enable   bool
	Haschild bool
	Menuid   string
	Menuname string
	Process  bool
	View     bool
	Parent   string
	Checkall bool
}

func (c *MenuSettingController) SaveMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id        string
		PageId    string
		Parent    string
		Title     string
		Url       string
		Icon      string
		IndexMenu int
		Enable    bool
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	mt := NewTopMenuModel()
	mt.Id = payLoad.Id
	mt.PageId = payLoad.PageId
	mt.Parent = payLoad.Parent
	mt.Title = payLoad.Title
	mt.Url = payLoad.Url
	mt.Icon = payLoad.Icon
	mt.IndexMenu = payLoad.IndexMenu
	mt.Enable = payLoad.Enable
	c.Ctx.Save(mt)

	// =============================== save/update menu on sysrole

	menu := Detailsmenu{
		Menuid:   mt.Id,
		Menuname: mt.Title,
		Access:   false,
		View:     false,
		Create:   false,
		Approve:  false,
		Delete:   false,
		Process:  false,
		Edit:     false,
		Parent:   mt.Parent,
		Haschild: false,
		Enable:   false,
		Url:      mt.Url,
		Checkall: false,
	}

	ress := new(SysRolesModel)
	query := tk.M{}
	csr, _ := c.Ctx.Find(ress, query)
	defer csr.Close()
	res := make([]SysRolesModel, 0)
	csr.Fetch(&res, 0, false)
	if len(res) > 0 {
		for _, v := range res {
			w := NewSysRolesModel()
			w.Id = v.Id
			w.Name = v.Name
			w.Menu = v.Menu
			w.Status = v.Status
			w.Landing = v.Landing

			savegaya := true
			for k2, v2 := range v.Menu {
				if v2.Menuid == mt.Id {
					w.Menu[k2].Menuname = mt.Title
					savegaya = false
				}
			}

			if savegaya {
				w.Menu = append(w.Menu, menu)
			}

			c.Ctx.Save(w)

		}
	}
	c.LogActivity("Menu Setting", "Save Menu", payLoad.Title, k)
	return c.SetResultInfo(false, "Menu has been successfully created.", nil)
}

func (c *MenuSettingController) DeleteMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id string
	}{}

	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	result := c.DeleteChildRecursive(payLoad.Id)

	//--------------- Delete SysRole -----------------
	data := make([]SysRolesModel, 0)
	crsData, errData := c.Ctx.Find(NewSysRolesModel(), nil)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	for _, dt := range data {
		NewMenu := []Detailsmenu{}
		ModelRole := NewSysRolesModel()
		ModelRole.Id = dt.Id
		ModelRole.Name = dt.Name
		for _, arrMenu := range dt.Menu {
			mVal, _ := tk.ToM(arrMenu)
			Menuid := mVal["Menuid"].(string)

			check := c.CheckAnyId(result, Menuid)

			if check {
				NewMenu = append(NewMenu, arrMenu)
			}
		}

		ModelRole.Menu = NewMenu
		ModelRole.Status = dt.Status
		ModelRole.Landing = dt.Landing
		c.Ctx.Save(ModelRole)
	}
	c.LogActivity("Menu Setting", "Delete Menu", payLoad.Id, k)
	return c.SetResultInfo(false, "Menu has been successfully created.", nil)
}

func (c *MenuSettingController) CheckAnyId(ids []string, id string) bool {
	for _, tempId := range ids {
		if tempId == id {
			return false
		}
	}

	return true
}

func (c *MenuSettingController) DeleteChildRecursive(ids string) []string {

	var (
		query  []*db.Filter
		result []string
	)

	query = append(query, db.Eq("Parent", ids))

	crsData, errData := c.Ctx.Connection.NewQuery().Select("_Id").From(new(TopMenuModel).TableName()).Where(query...).Cursor(nil)
	defer crsData.Close()

	if errData != nil {
		panic(errData)
	}

	topMenus := []tk.M{}
	errData = crsData.Fetch(&topMenus, 0, false)

	if errData != nil {
		panic(errData)
	}

	if len(topMenus) > 0 {
		for _, topMenu := range topMenus {
			result = append(result, c.DeleteChildRecursive(topMenu.GetString("_id"))...)
		}
	}

	u := NewTopMenuModel()
	query = append(query[0:0], db.Eq("_id", ids))
	crsData, errData = c.Ctx.Connection.NewQuery().From(new(TopMenuModel).TableName()).Where(query...).Cursor(nil)

	if errData != nil {
		panic(errData)
	}

	errData = crsData.Fetch(&u, 1, true)
	errData = c.Ctx.Delete(u)

	if errData != nil {
		panic(errData)
	}

	result = append(result, ids)
	c.WriteLog(result)

	return result
}

func (c *MenuSettingController) UpdateMenuTop(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.OutputType = knot.OutputJson
	payLoad := struct {
		Id        string
		PageId    string
		Parent    string
		Title     string
		Url       string
		Icon      string
		IndexMenu int
		Enable    bool
	}{}
	err := k.GetPayload(&payLoad)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}
	mt := NewTopMenuModel()
	mt.Id = payLoad.Id
	mt.PageId = payLoad.PageId
	mt.Parent = payLoad.Parent
	mt.Title = payLoad.Title
	mt.Url = payLoad.Url
	mt.Icon = payLoad.Icon
	mt.IndexMenu = payLoad.IndexMenu
	mt.Enable = payLoad.Enable
	c.Ctx.Save(mt)
	//--------------- Delete Update -----------------
	data := make([]SysRolesModel, 0)
	crsData, errData := c.Ctx.Find(NewSysRolesModel(), nil)
	defer crsData.Close()
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	errData = crsData.Fetch(&data, 0, false)
	if errData != nil {
		return c.SetResultInfo(true, errData.Error(), nil)
	}
	for _, dt := range data {
		NewMenu := []Detailsmenu{}
		ModelRole := NewSysRolesModel()
		ModelRole.Id = dt.Id
		ModelRole.Name = dt.Name
		for _, arrMenu := range dt.Menu {
			mVal, _ := tk.ToM(arrMenu)
			Menuid := mVal["Menuid"].(string)
			if Menuid != payLoad.Id {
				NewMenu = append(NewMenu, arrMenu)
			} else {
				UpdateMenu := Detailsmenu{
					Menuid:   mVal["Menuid"].(string),
					Menuname: payLoad.Title,
					Access:   mVal["Access"].(bool),
					View:     mVal["View"].(bool),
					Create:   mVal["Create"].(bool),
					Approve:  mVal["Approve"].(bool),
					Delete:   mVal["Delete"].(bool),
					Process:  mVal["Process"].(bool),
					Edit:     mVal["Edit"].(bool),
					Parent:   payLoad.Parent,
					Haschild: mVal["Haschild"].(bool),
					Enable:   payLoad.Enable,
					Url:      payLoad.Url,
					Checkall: mVal["Checkall"].(bool),
				}
				NewMenu = append(NewMenu, UpdateMenu)
			}
		}
		ModelRole.Menu = NewMenu
		ModelRole.Status = dt.Status
		ModelRole.Landing = dt.Landing
		c.Ctx.Save(ModelRole)
	}

	data1 := make([]SysRolesModel, 0)
	crsData1, errData1 := c.Ctx.Find(NewSysRolesModel(), nil)
	defer crsData1.Close()
	if errData1 != nil {
		return c.SetResultInfo(true, errData1.Error(), nil)
	}
	errData1 = crsData1.Fetch(&data1, 0, false)
	if errData1 != nil {
		return c.SetResultInfo(true, errData1.Error(), nil)
	}

	k.SetSession("roles", data1)
	k.SetSession("rolesid", data1[0].Id.Hex())
	k.SetSession("stime", time.Now())

	c.LogActivity("Menu Setting", "Update Menu", payLoad.Title, k)
	return c.SetResultInfo(false, "Menu has been successfully update.", nil)
}
