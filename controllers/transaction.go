package controllers

import knot "github.com/eaciit/knot/knot.v1"

type TransactionController struct {
	*BaseController
}

func (c *TransactionController) Journal(k *knot.WebContext) interface{} {
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
func (c *TransactionController) UpdateJournal(k *knot.WebContext) interface{} {
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
func (c *TransactionController) Asset(k *knot.WebContext) interface{} {
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
func (c *TransactionController) Rental(k *knot.WebContext) interface{} {
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
func (c *TransactionController) Capital(k *knot.WebContext) interface{} {
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

func (c *TransactionController) FirstAmount(k *knot.WebContext) interface{} {
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
func (c *TransactionController) PurchaseOrder(k *knot.WebContext) interface{} {
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

func (c *TransactionController) PurchaseCreditMemo(k *knot.WebContext) interface{} {
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
func (c *TransactionController) Invoice(k *knot.WebContext) interface{} {
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
func (c *TransactionController) SalesCreditMemo(k *knot.WebContext) interface{} {
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
func (c *TransactionController) TransactionOrder(k *knot.WebContext) interface{} {
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
