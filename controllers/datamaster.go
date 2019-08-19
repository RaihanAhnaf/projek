package controllers

import (
	. "../models"

	"github.com/eaciit/dbox"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type DatamasterController struct {
	*BaseController
}

func (c *DatamasterController) Default(k *knot.WebContext) interface{} {
	c.LoadBase(k)
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	return ""
}

func (d *DatamasterController) GetUsername(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	query := d.Ctx.Connection.NewQuery()
	result := []tk.M{}
	csr, e := query.
		Select("username").From("SysUsers").Order("username").Cursor(nil)
	e = csr.Fetch(&result, 0, false)

	if e != nil {
		return result
	}
	defer csr.Close()

	return result
}

func (d *DatamasterController) GetRolesRestricted(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	query := d.Ctx.Connection.NewQuery()
	result := []tk.M{}

	q := query.Select("name").From("SysRoles")

	var ur = r.Session("roles").([]SysRolesModel)
	var hasAllActive = false
	for _, role := range ur {
		if role.Name == "All Active" {
			hasAllActive = true
		}
	}
	if !hasAllActive {
		q = q.Where(dbox.Ne("name", "All Active"))
	}

	csr, e := q.Order("name").Cursor(nil)
	e = csr.Fetch(&result, 0, false)

	if e != nil {
		return result
	}
	defer csr.Close()

	return result
}

func (d *DatamasterController) GetRoles(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	query := d.Ctx.Connection.NewQuery()
	result := []tk.M{}
	csr, e := query.
		Select("name").From("SysRoles").Order("name").Cursor(nil)
	e = csr.Fetch(&result, 0, false)

	if e != nil {
		return result
	}
	defer csr.Close()

	return result
}

func (d *DatamasterController) Delete(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	p := struct {
		Id   bson.ObjectId
		Name string
	}{}
	e := r.GetPayload(&p)
	if e != nil {
		d.WriteLog(e)
	}

	result := new(SysUserModel)
	e = d.Ctx.GetById(result, p.Id)
	if e != nil {
		d.WriteLog(e)
	}

	e = d.Ctx.Delete(result)

	d.LogActivity("User", "Delete User", p.Name, r)
	return ""
}
