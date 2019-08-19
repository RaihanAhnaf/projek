package controllers

import (
	"../library/reflection"
	. "../models"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type DashboardController struct {
	*BaseController
}

func (c *DashboardController) Default(k *knot.WebContext) interface{} {
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
func (c *DashboardController) GetDataChartRevenue(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Year int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accRevenue := 5110
	acc7000 := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 7000 {
			acc7000 = append(acc7000, each.ACC_Code)
		}
	}
	dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accRevenue}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resRev := []tk.M{}
	e = csr.Fetch(&resRev, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	totRev := 0.0
	for i, each := range resRev {
		fixAmount := each.GetFloat64("Revenue") * -1
		resRev[i].Set("Revenue", fixAmount)
		totRev += fixAmount
	}
	// 7000
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc7000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res7000 := []tk.M{}
	e = csr.Fetch(&res7000, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for i, each := range res7000 {
		fixAmount := each.GetFloat64("Revenue") * -1
		res7000[i].Set("Revenue", fixAmount)
		totRev += fixAmount
		//
		for k, each2 := range resRev {
			if each.GetInt("_id") == each2.GetInt("_id") {
				amount := each2.GetFloat64("Revenue") + fixAmount
				resRev[k].Set("Revenue", amount)
			}
		}
	}
	prevStart := dateStart.AddDate(-1, 0, 0)
	prevEnd := dateStart
	dataPrevTotal := c.PrevRevenue(prevStart, prevEnd)
	data := struct {
		Chart     []tk.M
		Total     float64
		PrevTotal float64
	}{
		Chart:     resRev,
		Total:     totRev,
		PrevTotal: dataPrevTotal.Data.(float64),
	}
	return c.SetResultInfo(false, "SUCCESS", data)
}
func (c *DashboardController) GetDataChartExpenses(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Year int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":      "$Month",
		"Expenses": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resExp := []tk.M{}
	e = csr.Fetch(&resExp, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	totExp := 0.0
	for i, each := range resExp {
		fixAmount := each.GetFloat64("Expenses")
		resExp[i].Set("Expenses", fixAmount)
		totExp += fixAmount
	}
	prevStart := dateStart.AddDate(-1, 0, 0)
	prevEnd := dateStart
	dataPrevTotal := c.PrevExpenses(prevStart, prevEnd)
	data := struct {
		Chart     []tk.M
		Total     float64
		PrevTotal float64
	}{
		Chart:     resExp,
		Total:     totExp,
		PrevTotal: dataPrevTotal.Data.(float64),
	}
	return c.SetResultInfo(false, "SUCCESS", data)
}
func (c *DashboardController) PrevRevenue(prevStart time.Time, prevEnd time.Time) ResultInfo {

	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accRevenue := 5110
	acc7000 := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 7000 {
			acc7000 = append(acc7000, each.ACC_Code)
		}
	}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": prevStart, "$lt": prevEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accRevenue}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	Rev := []tk.M{}
	e = csr.Fetch(&Rev, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	totRev := 0.0
	for i, each := range Rev {
		fixAmount := each.GetFloat64("Revenue") * -1
		Rev[i].Set("Revenue", fixAmount)
		totRev += fixAmount
	}
	// 7000
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": prevStart, "$lt": prevEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc7000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res7000 := []tk.M{}
	e = csr.Fetch(&res7000, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for i, each := range res7000 {
		fixAmount := each.GetFloat64("Revenue") * -1
		res7000[i].Set("Revenue", fixAmount)
		totRev += fixAmount
	}
	return c.SetResultInfo(false, "Success", totRev)
}
func (c *DashboardController) PrevExpenses(prevStart time.Time, prevEnd time.Time) ResultInfo {
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": prevStart, "$lt": prevEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":      "",
		"Expenses": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	Rev := []tk.M{}
	e = csr.Fetch(&Rev, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	totExp := 0.0
	for i, each := range Rev {
		fixAmount := each.GetFloat64("Expenses")
		Rev[i].Set("Expenses", fixAmount)
		totExp += fixAmount
	}
	return c.SetResultInfo(false, "Success", totExp)
}
func (c *DashboardController) GetDataChartNetProfit(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Year int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	type NetProfitModel struct {
		Id     bson.ObjectId
		Month  int
		Amount float64
	}
	netprofitMap := map[int]*NetProfitModel{}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accPenjualan := 5110
	accPotongandanRetur := []int{5120, 5130}
	accHPP := 5210
	accOpex := []int{}
	acc7000 := []int{}
	acc8000 := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accOpex = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			acc7000 = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			acc8000 = append(accOpex, each.ACC_Code)
		}
	}
	dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accPenjualan}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPenjualan := []tk.M{}
	e = csr.Fetch(&resPenjualan, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPenjualan := 0.0
	for i, each := range resPenjualan {
		fixAmount := each.GetFloat64("Amounts") * -1
		resPenjualan[i].Set("Amounts", fixAmount)
		totalPenjualan += fixAmount
		net := NetProfitModel{}
		net.Id = bson.NewObjectId()
		net.Month = each.GetInt("_id")
		net.Amount = fixAmount
		netprofitMap[each.GetInt("_id")] = &net
	}
	// Potongan And retur
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accPotongandanRetur}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPotongandanRetur := []tk.M{}
	e = csr.Fetch(&resPotongandanRetur, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPotonganDanRetur := 0.0
	for i, each := range resPotongandanRetur {
		fixAmount := each.GetFloat64("Amounts")
		resPotongandanRetur[i].Set("Amounts", fixAmount)
		totalPotonganDanRetur += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}

	}
	TotalPenjualanBersih := totalPenjualan - totalPotonganDanRetur
	//HPP
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accHPP}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resHPP := []tk.M{}
	e = csr.Fetch(&resHPP, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalHPP := 0.0
	for i, each := range resHPP {
		fixAmount := each.GetFloat64("Amounts")
		resHPP[i].Set("Amounts", fixAmount)
		totalHPP += fixAmount
		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	LabaKotor := TotalPenjualanBersih - totalHPP
	// OPEX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accOpex}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resOpex := []tk.M{}
	e = csr.Fetch(&resOpex, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalOpex := 0.0
	for i, each := range resOpex {
		fixAmount := each.GetFloat64("Amounts")
		resOpex[i].Set("Amounts", fixAmount)
		totalOpex += fixAmount
		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	totalLabaUsaha := LabaKotor - totalOpex

	// Pendapatan lain2
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc7000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res7000 := []tk.M{}
	e = csr.Fetch(&res7000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalpendatanlain2 := 0.0
	for i, each := range res7000 {
		fixAmount := each.GetFloat64("Amounts") * -1
		res7000[i].Set("Amounts", fixAmount)
		totalpendatanlain2 += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount + fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 + fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	// Biaya Lain-lain
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc8000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res8000 := []tk.M{}
	e = csr.Fetch(&res8000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalbiayalain2 := 0.0
	for i, each := range res8000 {
		fixAmount := each.GetFloat64("Amounts")
		res8000[i].Set("Amounts", fixAmount)
		totalbiayalain2 += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	total70008000 := totalpendatanlain2 - totalbiayalain2
	netbeforetax := totalLabaUsaha + total70008000

	// TAX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": 6999}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTax := []tk.M{}
	e = csr.Fetch(&resTax, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalTax := 0.0
	for i, each := range resTax {
		fixAmount := each.GetFloat64("Amounts")
		resTax[i].Set("Amounts", fixAmount)
		totalTax += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	netProfit := netbeforetax - totalTax
	dataChart := []NetProfitModel{}
	for _, v := range netprofitMap {
		net := NetProfitModel{}
		net.Id = v.Id
		net.Month = v.Month
		net.Amount = v.Amount
		dataChart = append(dataChart, net)
	}

	prevStart := dateStart.AddDate(-1, 0, 0)
	prevEnd := dateStart
	prevNet := c.PrevNet(prevStart, prevEnd)
	data := struct {
		Total     float64
		Chart     []NetProfitModel
		PrevTotal float64
	}{
		Total:     netProfit,
		Chart:     dataChart,
		PrevTotal: prevNet.Data.(float64),
	}
	return c.SetResultInfo(false, "Success", data)
}
func (c *DashboardController) PrevNet(prevStart time.Time, prevEnd time.Time) ResultInfo {
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accPenjualan := 5110
	accPotongandanRetur := []int{5120, 5130}
	accHPP := 5210
	accOpex := []int{}
	acc7000 := []int{}
	acc8000 := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accOpex = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			acc7000 = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			acc8000 = append(accOpex, each.ACC_Code)
		}
	}
	dateStart := prevStart
	dateEnd := prevEnd
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accPenjualan}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPenjualan := []tk.M{}
	e = csr.Fetch(&resPenjualan, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPenjualan := 0.0
	for i, each := range resPenjualan {
		fixAmount := each.GetFloat64("Revenue") * -1
		resPenjualan[i].Set("Amounts", fixAmount)
		totalPenjualan += fixAmount
	}
	// Potongan And retur
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accPotongandanRetur}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPotongandanRetur := []tk.M{}
	e = csr.Fetch(&resPotongandanRetur, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPotonganDanRetur := 0.0
	for i, each := range resPotongandanRetur {
		fixAmount := each.GetFloat64("Amounts")
		resPotongandanRetur[i].Set("Amounts", fixAmount)
		totalPotonganDanRetur += fixAmount
	}
	TotalPenjualanBersih := totalPenjualan - totalPotonganDanRetur
	//HPP
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accHPP}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resHPP := []tk.M{}
	e = csr.Fetch(&resHPP, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalHPP := 0.0
	for i, each := range resHPP {
		fixAmount := each.GetFloat64("Amounts")
		resHPP[i].Set("Amounts", fixAmount)
		totalHPP += fixAmount
	}
	LabaKotor := TotalPenjualanBersih - totalHPP
	// OPEX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accOpex}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resOpex := []tk.M{}
	e = csr.Fetch(&resOpex, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalOpex := 0.0
	for i, each := range resOpex {
		fixAmount := each.GetFloat64("Revenue")
		resOpex[i].Set("Amounts", fixAmount)
		totalOpex += fixAmount
	}
	totalLabaUsaha := LabaKotor - totalOpex
	// Pendapatan lain2
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc7000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res7000 := []tk.M{}
	e = csr.Fetch(&res7000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalpendatanlain2 := 0.0
	for i, each := range res7000 {
		fixAmount := each.GetFloat64("Amounts") * -1
		res7000[i].Set("Amounts", fixAmount)
		totalpendatanlain2 += fixAmount

	}
	// Biaya Lain-lain
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc8000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res8000 := []tk.M{}
	e = csr.Fetch(&res8000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalbiayalain2 := 0.0
	for i, each := range res8000 {
		fixAmount := each.GetFloat64("Amounts")
		res8000[i].Set("Amounts", fixAmount)
		totalbiayalain2 += fixAmount
	}
	total70008000 := totalpendatanlain2 - totalbiayalain2
	netbeforetax := totalLabaUsaha + total70008000

	// TAX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": 6999}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTax := []tk.M{}
	e = csr.Fetch(&resTax, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalTax := 0.0
	for i, each := range resTax {
		fixAmount := each.GetFloat64("Amounts")
		resTax[i].Set("Amounts", fixAmount)
		totalTax += fixAmount
	}
	netProfit := netbeforetax - totalTax
	return c.SetResultInfo(false, "Success", netProfit)
}
func (c *DashboardController) GetDataChartRevExNet(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := struct {
		Year int
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	type ChartModel struct {
		Id     bson.ObjectId
		Month  int
		Amount float64
		Type   string
	}
	fixData := []ChartModel{}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accRevenue := []int{}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5110 {
			accRevenue = append(accRevenue, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			accRevenue = append(accRevenue, each.ACC_Code)
		}
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	//Revenue
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accRevenue}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resRevenue := []tk.M{}
	e = csr.Fetch(&resRevenue, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resRevenue {
		mod := ChartModel{}
		mod.Id = bson.NewObjectId()
		mod.Month = each.GetInt("_id")
		mod.Type = "Revenue"
		mod.Amount = each.GetFloat64("Amounts") * -1
		fixData = append(fixData, mod)
	}
	//Expenses
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resExpenses := []tk.M{}
	e = csr.Fetch(&resExpenses, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for _, each := range resExpenses {
		mod := ChartModel{}
		mod.Id = bson.NewObjectId()
		mod.Month = each.GetInt("_id")
		mod.Type = "Expenses"
		mod.Amount = each.GetFloat64("Amounts")
		fixData = append(fixData, mod)
	}
	netprofitFunc := c.GetNetProfit(dateStart, dateEnd)
	for _, each := range netprofitFunc.Data.([]tk.M) {
		mod := ChartModel{}
		mod.Id = bson.NewObjectId()
		mod.Month = each.GetInt("Month")
		mod.Type = "Net Profit"
		mod.Amount = each.GetFloat64("Amount")
		fixData = append(fixData, mod)
	}
	return c.SetResultInfo(false, "Success", fixData)
}
func (c *DashboardController) GetNetProfit(dateStart time.Time, dateEnd time.Time) ResultInfo {
	type NetProfitModel struct {
		Id     bson.ObjectId
		Month  int
		Amount float64
	}
	netprofitMap := map[int]*NetProfitModel{}
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accPenjualan := 5110
	accPotongandanRetur := []int{5120, 5130}
	accHPP := 5210
	accOpex := []int{}
	acc7000 := []int{}
	acc8000 := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accOpex = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			acc7000 = append(accOpex, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			acc8000 = append(accOpex, each.ACC_Code)
		}
	}
	// dateStart := time.Date(p.Year, time.January, 1, 0, 0, 0, 0, time.UTC)
	// dateEnd := dateStart.AddDate(1, 0, 0)
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accPenjualan}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPenjualan := []tk.M{}
	e = csr.Fetch(&resPenjualan, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPenjualan := 0.0
	for i, each := range resPenjualan {
		fixAmount := each.GetFloat64("Amounts") * -1
		resPenjualan[i].Set("Amounts", fixAmount)
		totalPenjualan += fixAmount
		net := NetProfitModel{}
		net.Id = bson.NewObjectId()
		net.Month = each.GetInt("_id")
		net.Amount = fixAmount
		netprofitMap[each.GetInt("_id")] = &net
	}
	// Potongan And retur
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accPotongandanRetur}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPotongandanRetur := []tk.M{}
	e = csr.Fetch(&resPotongandanRetur, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalPotonganDanRetur := 0.0
	for i, each := range resPotongandanRetur {
		fixAmount := each.GetFloat64("Amounts")
		resPotongandanRetur[i].Set("Amounts", fixAmount)
		totalPotonganDanRetur += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}

	}
	TotalPenjualanBersih := totalPenjualan - totalPotonganDanRetur
	//HPP
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": accHPP}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resHPP := []tk.M{}
	e = csr.Fetch(&resHPP, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalHPP := 0.0
	for i, each := range resHPP {
		fixAmount := each.GetFloat64("Amounts")
		resHPP[i].Set("Amounts", fixAmount)
		totalHPP += fixAmount
		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	LabaKotor := TotalPenjualanBersih - totalHPP
	// OPEX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accOpex}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resOpex := []tk.M{}
	e = csr.Fetch(&resOpex, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalOpex := 0.0
	for i, each := range resOpex {
		fixAmount := each.GetFloat64("Amounts")
		resOpex[i].Set("Amounts", fixAmount)
		totalOpex += fixAmount
		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	totalLabaUsaha := LabaKotor - totalOpex

	// Pendapatan lain2
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc7000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res7000 := []tk.M{}
	e = csr.Fetch(&res7000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalpendatanlain2 := 0.0
	for i, each := range res7000 {
		fixAmount := each.GetFloat64("Amounts") * -1
		res7000[i].Set("Amounts", fixAmount)
		totalpendatanlain2 += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount + fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 + fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	// Biaya Lain-lain
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": acc8000}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	res8000 := []tk.M{}
	e = csr.Fetch(&res8000, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalbiayalain2 := 0.0
	for i, each := range res8000 {
		fixAmount := each.GetFloat64("Amounts")
		res8000[i].Set("Amounts", fixAmount)
		totalbiayalain2 += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	total70008000 := totalpendatanlain2 - totalbiayalain2
	netbeforetax := totalLabaUsaha + total70008000

	// TAX
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$eq": 6999}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "$Month",
		"Amounts": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTax := []tk.M{}
	e = csr.Fetch(&resTax, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	totalTax := 0.0
	for i, each := range resTax {
		fixAmount := each.GetFloat64("Amounts")
		resTax[i].Set("Amounts", fixAmount)
		totalTax += fixAmount

		if netprofitMap[each.GetInt("_id")] != nil {
			netprofitMap[each.GetInt("_id")].Amount = netprofitMap[each.GetInt("_id")].Amount - fixAmount
		} else {
			net := NetProfitModel{}
			net.Id = bson.NewObjectId()
			net.Month = each.GetInt("_id")
			net.Amount = 0.0 - fixAmount
			netprofitMap[each.GetInt("_id")] = &net
		}
	}
	netProfit := netbeforetax - totalTax
	dataChart := []tk.M{}
	for _, v := range netprofitMap {
		net := tk.M{}
		net.Set("Id", v.Id)
		net.Set("Month", v.Month)
		net.Set("Amount", v.Amount)
		net.Set("TotalAmount", netProfit)
		dataChart = append(dataChart, net)
	}
	return c.SetResultInfo(false, "", dataChart)
}
func (c *DashboardController) MonthlyRevenue(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	now := time.Now()
	dateStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()
	prevStart := dateStart.AddDate(0, -1, 0)
	prevEnd := dateStart
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accRevenue := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5110 {
			accRevenue = append(accRevenue, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			accRevenue = append(accRevenue, each.ACC_Code)
		}
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accRevenue}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resRev := []tk.M{}
	e = csr.Fetch(&resRev, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": prevStart, "$lt": prevEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accRevenue}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPrevRev := []tk.M{}
	e = csr.Fetch(&resPrevRev, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	res := 0.0
	resPrev := 0.0
	if len(resRev) != 0 {
		res = resRev[0].GetFloat64("Revenue") * -1
	}
	if len(resPrevRev) != 0 {
		resPrev = resPrevRev[0].GetFloat64("Revenue") * -1
	}
	Data := []float64{res, resPrev}
	return c.SetResultInfo(false, "Success", Data)
}
func (c *DashboardController) MonthlyExpenses(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	now := time.Now()
	dateStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()
	prevStart := dateStart.AddDate(0, -1, 0)
	prevEnd := dateStart
	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resExp := []tk.M{}
	e = csr.Fetch(&resExp, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": prevStart, "$lt": prevEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":     "",
		"Revenue": tk.M{"$sum": "$Amount"},
	}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resPrevExp := []tk.M{}
	e = csr.Fetch(&resPrevExp, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	res := 0.0
	resPrev := 0.0
	if len(resExp) != 0 {
		res = resExp[0].GetFloat64("Revenue")
	}
	if len(resPrevExp) != 0 {
		resPrev = resPrevExp[0].GetFloat64("Revenue")
	}
	Data := []float64{res, resPrev}
	return c.SetResultInfo(false, "Success", Data)
}
func (c *DashboardController) GetDataTopFiveExpenses(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	now := time.Now()
	// now = now.AddDate(0, -5, 0)
	dateStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()

	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code == 6000 && each.ACC_Code != 6999 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 8000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	// tk.Println(dateStart, dateEnd)
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$Acc_Name",
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"Amount":    tk.M{"$sum": "$Amount"},
		// "Amount":    tk.M{"$subtract": amount},
	}})
	// pipes = append(pipes, tk.M{"$group": tk.M{
	// 	"_id":    "$_id",
	// 	"Amount": tk.M{"$subtract": amount},
	// }})
	pipes = append(pipes, tk.M{"$sort": tk.M{"Amount": -1}})
	pipes = append(pipes, tk.M{"$limit": 5})
	// tk.Println("Pipies => ", pipes)
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTop5 := []tk.M{}
	// tk.Println("resTop5 =>", resTop5)
	e = csr.Fetch(&resTop5, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", resTop5)
}
func (c *DashboardController) GetDataTopFiveRevenue(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	now := time.Now()
	// now = now.AddDate(0, -5, 0)
	dateStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()

	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	accExpenses := []int{}
	for _, each := range resuts {
		if each.ACC_Code == 5110 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
		if each.Main_Acc_Code == 7000 {
			accExpenses = append(accExpenses, each.ACC_Code)
		}
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	// tk.Println(dateStart, dateEnd)
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accExpenses}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$Acc_Name",
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"Amount":    tk.M{"$sum": "$Amount"},
		// "Amount":    tk.M{"$subtract": amount},
	}})
	// pipes = append(pipes, tk.M{"$group": tk.M{
	// 	"_id":    "$_id",
	// 	"Amount": tk.M{"$subtract": amount},
	// }})
	pipes = append(pipes, tk.M{"$sort": tk.M{"Amount": -1}})
	pipes = append(pipes, tk.M{"$limit": 5})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTop5 := []tk.M{}
	e = csr.Fetch(&resTop5, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", resTop5)
}

// o.Main_Acc_Code > 0 && o.ACC_Code < 2000 && o.ACC_Code > 1000;
func (c *DashboardController) GetDataCurrentAsset(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	now := time.Now()
	// now = now.AddDate(0, -5, 0)
	dateStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()

	csr, e := c.Ctx.Connection.NewQuery().Select().From("Coa").Cursor(nil)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resuts := []CoaModel{}
	e = csr.Fetch(&resuts, 0, false)
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}
	allAcc := []int{}
	for _, each := range resuts {
		if each.Main_Acc_Code > 0 && each.ACC_Code > 1100 && each.ACC_Code < 1200 {
			allAcc = append(allAcc, each.ACC_Code)
		}
	}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	// tk.Println(dateStart, dateEnd)
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": allAcc}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "",
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"SumSaldo":  tk.M{"$sum": "$Amount"},
		// "Amount":    tk.M{"$subtract": amount},
	}})
	// pipes = append(pipes, tk.M{"$group": tk.M{
	// 	"_id":    "$_id",
	// 	"Amount": tk.M{"$subtract": amount},
	// }})
	pipes = append(pipes, tk.M{"$sort": tk.M{"Amount": -1}})
	csr, e = c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	resTop5 := []tk.M{}
	e = csr.Fetch(&resTop5, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", resTop5)
}
func (c *DashboardController) GetDataBank(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	// 1120 : CashOut
	// 1110 : Kas
	// 1121 : CashIn
	accountBANK := []int{1120, 1121, 1110, 1111}
	amount := []string{"$ListDetail.Debet", "$ListDetail.Credit"}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"Status": tk.M{"$ne": "draft"}}})
	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"_id":            "$ListDetail._id",
		"No":             "$ListDetail.No",
		"Journal_Type":   "$ListDetail.Journal_Type",
		"PostingDate":    "$ListDetail.PostingDate",
		"DateStr":        "$ListDetail.DateStr",
		"DocumentNumber": "$ListDetail.DocumentNumber",
		"Acc_Code":       "$ListDetail.Acc_Code",
		"Acc_Name":       "$ListDetail.Acc_Name",
		"Debet":          "$ListDetail.Debet",
		"Credit":         "$ListDetail.Credit",
		"Description":    "$ListDetail.Description",
		"Attachment":     "$ListDetail.Attachment",
		"User":           "$ListDetail.User",
		"Amount":         tk.M{"$subtract": amount},
		"Month":          tk.M{"$month": "$ListDetail.PostingDate"},
	}})
	pipes = append(pipes, tk.M{"$match": tk.M{"Acc_Code": tk.M{"$in": accountBANK}}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$Acc_Code",
		"SumDebet":  tk.M{"$sum": "$Debet"},
		"SumCredit": tk.M{"$sum": "$Credit"},
		"Amount":    tk.M{"$sum": "$Amount"},
		// "Amount":    tk.M{"$subtract": amount},
	}})
	csr, e := c.Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	defer csr.Close()
	results := []tk.M{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *DashboardController) GetDataInvoiceDashboard(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	csr, e := c.Ctx.Connection.NewQuery().From("Invoice").Where(db.Ne("Status", "DRAFT")).Cursor(nil)
	defer csr.Close()
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	results := []InvoiceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for i, each := range results {
		if each.Status != "PAID" {
			results[i].Status = "Out Standing"
		}
	}
	return c.SetResultInfo(false, "Success", results)
}
func (c *DashboardController) GetDataPoDashboard(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	csr, e := c.Ctx.Connection.NewQuery().From("PurchaseOrder").Where(db.Ne("Status", "PO")).Cursor(nil)
	defer csr.Close()
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	results := []PurchaseOrder{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}

	// Now add records from PurchaseInventory
	csi, e := c.Ctx.Connection.NewQuery().From("PurchaseInventory").Where(db.Ne("Status", "PI")).Cursor(nil)
	defer csi.Close()
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	resulti := []PurchaseInventory{}
	e = csi.Fetch(&resulti, 0, false)
	if e != nil {
		c.SetResultInfo(true, e.Error(), nil)
	}
	for _, pi := range resulti {
		var porec = PurchaseOrder{}
		e = reflection.Map(&pi, &porec)
		porec.DatePosting = pi.DatePosting
		if e == nil {
			results = append(results, porec)
		}
	}

	for i, each := range results {
		if each.Status != "PAID" {
			results[i].Status = "Out Standing"
		}
	}
	return c.SetResultInfo(false, "Success", results)
}
