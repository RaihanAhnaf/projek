package controllers

import (
	"github.com/eaciit/knot/knot.v1"
)

type BackOfficeController struct {
	*BaseController
}

func (c *BackOfficeController) Default(k *knot.WebContext) interface{} {
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
