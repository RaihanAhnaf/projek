package main

import (
	// xx "eaciit/proactive-dev/controllers"
	"./helpers"
	. "../../models"
	db "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
	// "os"
	// "math"
	"strconv"
	"time"
)

func PrepareConnection() (db.IConnection, error) {
	config := helpers.ReadConfig()
	ci := &db.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := db.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}
func main() {
	// conn, err := PrepareConnection()
	// if err != nil {
	// 	tk.Println(err)
	// }
	// err := os.Remove("income.pdf")
	// if err != nil {
	// 	tk.Println(err)
	// 	return
	// }
	dateStart := time.Date(2017, time.August, 1, 0, 0, 0, 0, time.UTC)
	dateEnd := time.Now()
	BalanceSheet(dateStart, dateEnd)

	// IncomeStatementPDF(dateStart, dateEnd)
	// GetCurrentEarning(dateStart, dateEnd)
}

type IncomeStatementModel struct {
	Id       bson.ObjectId
	Acc_Code int
	Acc_Name string
	Ending   float64
	Sales    string
}

type BalanceSheetModel struct {
	Id       bson.ObjectId
	Acc_Code int
	Acc_Name string
	Amount   float64
}

func GetDebet(DateStart time.Time, DateEnd time.Time) (string, []CoaModel) {

	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	Ctx := orm.New(conn)
	config := helpers.ReadConfig()
	Img := config["imgpath"]
	tk.Println(Img)
	csr, e := Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(db.Eq("category", "BALANCE SHEET"), db.Eq("debet_credit", "DEBET"), db.Ne("main_acc_code", 0))).Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	results := []CoaModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	//tk.Println(results)

	var pipes []tk.M
	StartTime, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)

	tk.Println(StartTime)
	tk.Println(EndTime)
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gt": StartTime, "$lt": EndTime}}})

	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	csr, e = Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	resultJournal := make([]tk.M, 0)
	e = csr.Fetch(&resultJournal, 0, false)
	if e != nil {
		return e.Error(), nil
	}

	//dataEarningg := getEarningg(DateStart, DateEnd)
	for _, each := range resultJournal {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("SumDebet")
				results[i].Credit = each.GetFloat64("SumCredit")
				results[i].Saldo = results[i].Debet - results[i].Credit
			}

			// if results[i].ACC_Code == 4400 {
			// 	results[i].Saldo = dataEarning.Data.(float64)
			// }
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		for _, z := range results {
			parent := z.Main_Acc_Code
			if parent == accCode {
				results[i].Debet = results[i].Debet + z.Debet
				results[i].Credit = results[i].Credit + z.Credit
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
		}
	}
	return "", results
}

func GetCredit(DateStart time.Time, DateEnd time.Time) (string, []CoaModel) {
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	Ctx := orm.New(conn)
	config := helpers.ReadConfig()
	Img := config["imgpath"]
	tk.Println(Img)
	csr, e := Ctx.Connection.NewQuery().Select().From("Coa").Where(db.And(db.Eq("category", "BALANCE SHEET"), db.Eq("debet_credit", "CREDIT"))).Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	results := []CoaModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error(), nil
	}

	var pipes []tk.M
	StartTime, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	EndTime, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	EndTime = EndTime.AddDate(0, 0, 1)

	tk.Println(StartTime)
	tk.Println(EndTime)
	pipes = append(pipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gt": StartTime, "$lt": EndTime}}})

	pipes = append(pipes, tk.M{"$unwind": "$ListDetail"})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$ListDetail.Acc_Code",
		"SumDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"SumCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})

	csr, e = Ctx.Connection.NewQuery().Command("pipe", pipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), nil
	}
	defer csr.Close()
	resultJournal := make([]tk.M, 0)
	e = csr.Fetch(&resultJournal, 0, false)
	if e != nil {
		return e.Error(), nil
	}
	// tk.Println(resultJournal)
	_, dataEarningg := GetCurrentEarning(DateStart, DateEnd)
	tk.Println(dataEarningg)
	for _, each := range resultJournal {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				results[i].Debet = each.GetFloat64("SumDebet")
				results[i].Credit = each.GetFloat64("SumCredit")
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
			// if results[i].ACC_Code == 4400 {
			// 	results[i].Saldo = dataEarning.Data.(float64)
			// }
		}
	}
	for i, _ := range results {
		accCode := results[i].ACC_Code
		for _, z := range results {
			parent := z.Main_Acc_Code
			if parent == accCode {
				results[i].Debet = results[i].Debet + z.Debet
				results[i].Credit = results[i].Credit + z.Credit
				results[i].Saldo = results[i].Debet - results[i].Credit
			}
		}

		if accCode == 4400 {
			results[i].Saldo = dataEarningg
		}
	}

	return "", results
}

func GetCurrentEarning(DateStart time.Time, DateEnd time.Time) (string, float64) {
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	Ctx := orm.New(conn)

	dateStart, _ := time.Parse("2006-01-02", DateStart.Format("2006-01-02"))
	dateEnd, _ := time.Parse("2006-01-02", DateEnd.Format("2006-01-02"))
	dateEnd = dateEnd.AddDate(0, 0, 1)
	// Sales And Revenue
	// ninsales := []int{5400, 5300}
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 5000, "$lt": 6000}}})
	csr, e := Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), 0
	}
	defer csr.Close()
	res := make([]tk.M, 0)
	e = csr.Fetch(&res, 0, false)
	if e != nil {
		return e.Error(), 0
	}
	GrossProfit := 0.0
	for _, each := range res {
		if each.GetInt("_id") != 5300 || each.GetInt("_id") != 5400 {
			trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
			if each.GetInt("_id") == 5100 || each.GetInt("_id") == 5200 {
				trans = trans * -1
				GrossProfit = GrossProfit + trans
			}
			if each.GetInt("_id") == 5600 {
				GrossProfit = GrossProfit - trans
			}
		}
	}
	// tk.Println(GrossProfit)
	//OPERATING EXPENSE
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 6000, "$lt": 7000}}})
	csr, e = Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), 0
	}
	defer csr.Close()
	resoperex := make([]tk.M, 0)
	e = csr.Fetch(&resoperex, 0, false)
	if e != nil {
		return e.Error(), 0
	}
	totOpex := 0.0
	for _, each := range resoperex {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totOpex = totOpex + trans
	}
	// tk.Println(totOpex)
	netProfit := GrossProfit - totOpex
	//PENGHASILAN LAIN LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 7000, "$lt": 8000}}})
	csr, e = Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), 0
	}
	defer csr.Close()
	reslainlain := make([]tk.M, 0)
	e = csr.Fetch(&reslainlain, 0, false)
	if e != nil {
		return e.Error(), 0
	}
	// tk.Println(netProfit)
	totLainLain := 0.0
	for _, each := range reslainlain {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totLainLain = totLainLain + trans
	}
	// tk.Println(totLainLain)
	totLainLain = totLainLain * -1
	//BEBAN LAIN - LAIN
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id": tk.M{"$gte": 8000, "$lt": 9000}}})
	csr, e = Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), 0
	}
	defer csr.Close()
	resBebanlainlain := make([]tk.M, 0)
	e = csr.Fetch(&resBebanlainlain, 0, false)
	if e != nil {
		return e.Error(), 0
	}
	// tk.Println(resBebanlainlain)
	totBebanLainLain := 0.0
	for _, each := range resBebanlainlain {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totBebanLainLain = totBebanLainLain + trans
	}
	totAllRevnEx := totLainLain - totBebanLainLain
	totalEarningBefore := netProfit + totAllRevnEx
	//TAX
	Tpipes = []tk.M{}
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               tk.M{"Code": "$ListDetail.Acc_Code", "Name": "$ListDetail.Acc_Name"},
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"_id.Code": tk.M{"$ne": 6888}, "_id.Name": tk.M{"$eq": "TAX"}}})
	csr, e = Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error(), 0
	}
	defer csr.Close()
	resTax := make([]tk.M, 0)
	e = csr.Fetch(&resTax, 0, false)
	if e != nil {
		return e.Error(), 0
	}
	totTax := 0.0
	for _, each := range resTax {
		trans := each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
		totTax = totTax + trans
	}
	// tk.Println(GrossProfit, totOpex, netProfit, totLainLain, totBebanLainLain, totAllRevnEx, totalEarningBefore, totTax)
	totAfterTax := totalEarningBefore - totTax
	type modelTemp struct {
		amount float64
	}
	tk.Println(totAfterTax)
	// model := new(modelTemp)
	// model.amount = totAfterTax
	// return SetResultInfo(false, "success", totAfterTax)
	// tk.Println(totAfterTax)
	return "", totAfterTax

}

func BalanceSheet(dateStart time.Time, dateEnd time.Time) interface{} {
	err, datadebet := GetDebet(dateStart, dateEnd)
	if err != "" {
		tk.Println(err)
	}
	err, datacredit := GetCredit(dateStart, dateEnd)
	if err != "" {
		tk.Println(err)
	}
	// err, datacurrentearning := GetCurrentEarning(dateStart, dateEnd)
	// if err != "" {
	// 	tk.Println(err)
	// }

	// datadebet := GetDebet(dateStart, dateEnd)
	// datacredit := GetCredit(dateStart, dateEnd)

	config := helpers.ReadConfig()
	Img := config["imgpath"]
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetXY(8, 5)

	pdf.SetFont("Arial", "", 12)
	pdf.Image(Img+"eaciit-logo.png", 10, 10, 17, 17, false, "", 0, "")
	pdf.Ln(2)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(30)
	pdf.CellFormat(0, 12, "PT. WIYASA TEKNOLOGI NUSANTARA", "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(30)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 15, "Jl. Imam Bonjol No. 120, DR.Soetomo - Tegalsari", "", 0, "L", false, 0, "")
	pdf.SetX(180)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 15, "BALANCE SHEET", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, "Surabaya-Indonesia", "", 0, "L", false, 0, "")
	pdf.SetX(10)
	// pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, "Telp. 0351-5676223", "", 0, "L", false, 0, "")
	pdf.SetX(10)

	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)
	pdf.SetFont("Arial", "", 7)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+20, 286, pdf.GetY()+20) //garis horizontal2

	yAsset := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(yAsset + 10)
	pdf.MultiCell(138.5, 5, "Date Periode : "+dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetX(9)
	pdf.SetY(yAsset + 15)
	pdf.MultiCell(138.5, 5, "Date Created : "+time.Now().Format("02 January 2006"), "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(yAsset + 21)
	pdf.SetX(10)
	pdf.MultiCell(138.5, 5, "ASSETS", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(pdf.GetY())
	pdf.SetX(37.0)
	pdf.MultiCell(138.5, 7, "CURRENT ASSETS", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(yAsset + 21)
	pdf.SetX(148)
	pdf.MultiCell(138.5, 5, "PASSIVA", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	pdf.SetY(pdf.GetY())
	pdf.SetX(174.0)
	pdf.MultiCell(138.5, 7, "CURRENT LIABILITIES", "0", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)

	BelanceHead := []string{"ACC CODE", "ACC NAME", "Amount"}
	widthHead := []float64{24.0, 75.0, 30.0}
	y0 := pdf.GetY()

	for i, head := range BelanceHead {
		if i == 2 {
			pdf.SetY(y0)
			x := 13.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "", "R", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "", "R", false)
			}
		} else {
			pdf.SetY(y0)
			x := 13.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "", "L", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "", "L", false)
			}
		}

	}

	for i, head := range BelanceHead {
		if i == 2 {
			pdf.SetY(y0)
			x := 150.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "0", "R", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "0", "R", false)
			}
		} else {
			pdf.SetY(y0)
			x := 150.0
			for y, z := range widthHead {
				if i > y {
					x += z
				} else {
					x += 0.0
				}
			}
			pdf.SetX(x)
			if i == 0 {
				pdf.MultiCell(widthHead[i], 5, head, "0", "L", false)
			} else {
				pdf.MultiCell(widthHead[i], 5, head, "0", "L", false)
			}
		}

	}

	datadebcred := []BalanceSheetModel{}

	data1000 := []BalanceSheetModel{}
	total1000 := 0.0
	//data100
	data2000 := []BalanceSheetModel{}
	total2000 := 0.0
	//other revenu
	data3000 := []BalanceSheetModel{}
	total3000 := 0.0

	//other expense
	data4000 := []BalanceSheetModel{}
	total4000 := 0.0

	for _, each := range datadebet {
		if each.Saldo != 0 {
			if each.ACC_Code >= 1111 && each.ACC_Code < 2000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				data1000 = append(data1000, BalanMod)
				total1000 += BalanMod.Amount
			}
			if each.ACC_Code > 2100 && each.ACC_Code < 3000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				data2000 = append(data2000, BalanMod)

			}
			if each.ACC_Code == 2200 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}
			if each.ACC_Code == 2600 {
				var BalanMod BalanceSheetModel

				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)

			}
			if each.ACC_Code == 2300 {
				var BalanMod BalanceSheetModel

				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}

			if each.ACC_Code == 2400 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				/*Percentage := "0.00%"
				totOpEx += each.Ending
				IncomeMod.Sales = Percentage*/
				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
				// datadebcred = BalanMod.Amount
			}
			if each.ACC_Code == 2800 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)

			}
		}
	}

	//Data Credit
	for _, each := range datacredit {
		if each.Saldo != 0 {
			if each.ACC_Code > 3000 && each.ACC_Code < 4000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				if each.Saldo != 0 {
					BalanMod.Amount = each.Saldo * -1
				}
				// totOtherRev += each.Ending
				// IncomeMod.Sales = Percentage
				data3000 = append(data3000, BalanMod)
				total3000 += BalanMod.Amount
			}
			if each.ACC_Code > 4000 && each.ACC_Code < 5000 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				if each.Saldo != 0 && each.ACC_Code != 4400 {
					BalanMod.Amount = each.Saldo * -1
				}

				data4000 = append(data4000, BalanMod)
				total4000 += BalanMod.Amount
			}
			if each.ACC_Code == 2200 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				// datadebcred = BalanMod.Amount
			}
			if each.ACC_Code == 2600 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
				// datadebcred = BalanMod.Amount
			}
			if each.ACC_Code == 2300 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				// datadebcred = BalanMod.Amount
			}

			if each.ACC_Code == 2400 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				// datadebcred = BalanMod.Amount
			}
			if each.ACC_Code == 2700 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo

				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount

			}

			if each.ACC_Code == 2800 {
				var BalanMod BalanceSheetModel
				// IncomeMod := new(IncomeStatementModel)
				BalanMod.Id = each.ID
				BalanMod.Acc_Code = each.ACC_Code
				BalanMod.Acc_Name = each.Account_Name
				BalanMod.Amount = each.Saldo
				/*Percentage := "0.00%"
				totOpEx += each.Ending
				IncomeMod.Sales = Percentage*/
				datadebcred = append(datadebcred, BalanMod)
				total2000 += BalanMod.Amount
				// datadebcred = BalanMod.Amount
			}
		}
	}

	// yfirst := pdf.GetY()
	yfix := pdf.GetY()
	yA := yfix
	pdf.SetY(yA)
	for _, each := range data1000 {
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yA)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[1]
		pdf.SetXY(x, yA)
		amount := tk.Sprintf("%.2f", each.Amount)
		// amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			// amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		yA = pdf.GetY()
	}
	// pdf.SetY(yA)
	pdf.Ln(2)
	yA = pdf.GetY()

	yB := yfix
	pdf.SetY(yB)
	for _, each := range data3000 {
		x := 150.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yB)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[1]
		pdf.SetXY(x, yB)
		amount := tk.Sprintf("%.2f", each.Amount)
		// amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			// amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		yB = pdf.GetY()
	}
	// pdf.SetY(yB)
	pdf.Ln(2)
	yB = pdf.GetY()

	if yA > yB {
		yfix = yA
	} else {
		yfix = yB
	}

	// ycoba := pdf.GetY()

	// pdf.Ln(8)
	yC := yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "TOTAL CURRENT ASSETS", "0", "L", false)
	x2 := 112.0
	pdf.SetXY(x2, yC)
	totalCurAss := tk.Sprintf("%.2f", total1000)
	// totalCurAss = c.ConvertToCurrency(totalCurAss)
	if total1000 < 0 {
		totalCurAss = "(" + tk.Sprintf("%.2f", total1000) + ")"
	}
	pdf.MultiCell(widthHead[2], 5, totalCurAss, "", "R", false)
	pdf.Ln(2)
	yC = pdf.GetY()

	yD := yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "TOTAL CURRENT LIABILITIES", "0", "L", false)
	xtocur := 249.0
	pdf.SetXY(xtocur, yD)
	totalCurrLiab := tk.Sprintf("%.2f", total3000)
	// totalCurrLiab = c.ConvertToCurrency(totalCurrLiab)
	if total3000 < 0 {
		totalCurrLiab = "(" + tk.Sprintf("%2.f", total3000*-1) + ")"
	}
	pdf.MultiCell(widthHead[2], 5, totalCurrLiab, "", "R", false)
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}

	yC = yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "FIX ASSETS", "0", "L", false)

	//data credit

	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2200 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
			// x += widthHead[2]
			// pdf.SetXY(x, yC)
			// a3 := pdf.GetY()
			// allA := []float64{a0, a1, a2, a3}

			// var n, biggest float64
			// for _, v := range allA {
			// 	if v > n {
			// 		n = v
			// 		biggest = n
			// 	}
			// }
			// pdf.SetY(biggest)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2600 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2300 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
			x += widthHead[2]
			// pdf.SetXY(x, yC)
			// pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)

		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2700 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
			// pdf.SetXY(x, yC)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2400 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		}
	}
	yC = pdf.GetY()
	for _, each := range datadebcred {
		if each.Acc_Code == 2800 {
			pdf.SetY(yC)
			x := 13.0
			pdf.SetX(x)
			pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
			x += widthHead[0]
			pdf.SetXY(x, yC)
			pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
			x += widthHead[1]
			pdf.SetXY(x, yC)
			amount := tk.Sprintf("%.2f", each.Amount)
			// amount = c.ConvertToCurrency(amount)
			if each.Amount < 0 {
				amount = tk.Sprintf("%.2f", each.Amount*-1)
				// amount = c.ConvertToCurrency(amount)
				amount = "(" + amount + ")"
			}
			pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		}
	}
	pdf.Ln(2)
	yC = pdf.GetY()

	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "CAPITAL AND EARNING", "0", "L", false)

	// pdf.Ln(2)
	yD = pdf.GetY()
	pdf.SetY(yD)
	for _, each := range data4000 {
		x := 150.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "0", "L", false)
		x += widthHead[0]
		pdf.SetXY(x, yD)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "0", "L", false)
		x += widthHead[1]
		pdf.SetXY(x, yD)
		amount := tk.Sprintf("%.2f", each.Amount)
		// amount = c.ConvertToCurrency(amount)
		if each.Amount < 0 {
			amount = tk.Sprintf("%.2f", each.Amount*-1)
			// amount = c.ConvertToCurrency(amount)
			amount = "(" + amount + ")"
		}
		pdf.MultiCell(widthHead[2], 5, amount, "0", "R", false)
		yD = pdf.GetY()
	}
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}

	yC = yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "TOTAL FIX ASSETS", "0", "L", false)
	x2 = 112.0
	pdf.SetXY(x2, yC)
	totalFixAss := tk.Sprintf("%.2f", total2000)
	// totalFixAss = c.ConvertToCurrency(totalFixAss)
	if total2000 < 0 {
		totalFixAss = "(" + tk.Sprintf("%.2f", total2000) + ")"
	}
	pdf.MultiCell(30, 5, totalFixAss, "", "R", false)
	pdf.Ln(2)
	yC = pdf.GetY()

	// pdf.Ln(2)

	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "TOTAL CAPITAL AND EARNING", "0", "L", false)
	xToCapandEar := 249.0
	pdf.SetXY(xToCapandEar, yD)
	// pdf.SetY(yD)
	totalCapEar := tk.Sprintf("%.2f", total4000)
	// totalCapEar = c.ConvertToCurrency(totalCapEar)
	if total4000 < 0 {
		totalCapEar = "(" + tk.Sprintf("%.2f", total4000*-1) + ")"
	}
	pdf.MultiCell(widthHead[2], 5, totalCapEar, "", "R", false)
	// tk.Println(totalCapEar)
	pdf.Ln(2)
	yD = pdf.GetY()

	if yC > yD {
		yfix = yC
	} else {
		yfix = yD
	}

	yC = yfix
	pdf.SetXY(37.0, yC)
	pdf.MultiCell(92.5, 7, "TOTAL ASSETS", "0", "L", false)
	x2 = 112.0
	pdf.SetXY(x2, yC)
	totalAsset := total1000 + total2000
	totalAssetFix := tk.Sprintf("%.2f", totalAsset)
	// totalAssetFix = c.ConvertToCurrency(totalAssetFix)
	if totalAsset < 0 {
		totalAssetFix = "(" + tk.Sprintf("%.2f", totalAsset) + ")"
	}
	pdf.MultiCell(widthHead[2], 5, totalAssetFix, "", "R", false)

	yD = yfix
	pdf.SetXY(174.0, yD)
	pdf.MultiCell(92.5, 7, "TOTAL PASSIVA", "", "L", false)

	xToPass := 249.0
	pdf.SetXY(xToPass, yD)
	totalPassiva := total3000 + total4000
	totalPassivaFix := tk.Sprintf("%.2f", totalPassiva)
	// totalPassivaFix = c.ConvertToCurrency(totalPassivaFix)
	if totalPassiva < 0 {
		totalPassivaFix = "(" + tk.Sprintf("%.2f", totalPassiva) + ")"
	}
	pdf.MultiCell(30, 5, totalPassivaFix, "", "R", false)

	y3 := pdf.GetY()
	yline := yAsset + 9
	pdf.Line(pdf.GetX(), yline, pdf.GetX(), y3)
	pdf.Line(pdf.GetX()+138, yline+11, pdf.GetX()+138, y3)
	pdf.Line(pdf.GetX()+276, yline, pdf.GetX()+276, y3)
	pdf.Line(pdf.GetX(), y3, 286, y3)

	namepdf := "-balancesheet.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	e := pdf.OutputFileAndClose(fileName)
	tk.Println(e)
	if e != nil {
		e.Error()
	}
	// tk.Println(GetDebet(dateStart, dateEnd))
	//tk.Println(GetCredit(dateStart, dateEnd))
	return ""
}

func IncomeStatementPDF(dateStart time.Time, dateEnd time.Time) interface{} {
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	Ctx := orm.New(conn)
	config := helpers.ReadConfig()
	Img := config["imgpath"]
	csr, e := Ctx.Connection.NewQuery().Select().From("Coa").Where(db.Eq("category", "INCOME STATEMENT")).Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()
	results := []TrialBalanceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	var Tpipes []tk.M
	Tpipes = append(Tpipes, tk.M{"$match": tk.M{"PostingDate": tk.M{"$gte": dateStart, "$lt": dateEnd}}})
	Tpipes = append(Tpipes, tk.M{"$unwind": "$ListDetail"})
	Tpipes = append(Tpipes, tk.M{"$group": tk.M{
		"_id":               "$ListDetail.Acc_Code",
		"TransactionDebet":  tk.M{"$sum": "$ListDetail.Debet"},
		"TransactionCredit": tk.M{"$sum": "$ListDetail.Credit"},
	}})
	csr, e = Ctx.Connection.NewQuery().Command("pipe", Tpipes).From("GeneralLedger").Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()
	resultTransaction := make([]tk.M, 0)
	e = csr.Fetch(&resultTransaction, 0, false)
	if e != nil {
		return e.Error()
	}
	for _, each := range resultTransaction {
		for i, _ := range results {
			if each.GetInt("_id") == results[i].ACC_Code {
				// if results[i].Debet_Credit == "DEBET" {
				// 	results[i].Transaction = each.GetFloat64("TransactionDebet")
				// } else {
				// 	results[i].Transaction = each.GetFloat64("TransactionCredit")
				// }
				results[i].Transaction = each.GetFloat64("TransactionDebet") - each.GetFloat64("TransactionCredit")
				results[i].Ending = results[i].Begining + results[i].Transaction
			}
		}
	}
	for i, _ := range results {
		results[i].Begining = 0.0
		if results[i].ACC_Code == 5100 || results[i].ACC_Code == 5200 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
		if results[i].ACC_Code > 7000 && results[i].ACC_Code < 8000 {
			results[i].Transaction = results[i].Transaction * -1
			results[i].Ending = results[i].Begining + results[i].Transaction
		}
	}
	totSalesAndRevenue := 0.0
	for _, each := range results {
		if each.Main_Acc_Code == 5000 {
			totSalesAndRevenue = totSalesAndRevenue + each.Ending
		}
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetXY(10, 5)
	pdf.SetFont("Arial", "B", 20)
	// pdf.Image(imageNameStr, x, y, w, h, flow, tp, link, linkStr)
	pdf.Image(Img+"logoeaciit2.png", 10, 10, 21, 17, false, "", 0, "")
	pdf.Ln(3)
	y1 := pdf.GetY()
	pdf.SetY(y1 + 4)
	pdf.SetX(32)
	pdf.CellFormat(0, 12, "INCOME STATEMENT REPORT", "", 0, "L", false, 0, "")
	pdf.SetY(y1 + 10)
	pdf.SetX(32)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 12, "PT. Wiyasa Tekhnologi Nusantara", "", 0, "L", false, 0, "")
	pdf.SetX(10)
	y2 := pdf.GetY()
	pdf.SetY(y2 + 2)
	pdf.Line(pdf.GetX()+3, pdf.GetY()+9, 198, pdf.GetY()+9)
	pdf.SetFont("Arial", "", 7)
	pdf.Ln(8)
	pdf.GetY()
	pdf.SetX(12)
	pdf.CellFormat(10, 10, "Periode", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 10, dateStart.Format("02 January 2006")+" - "+dateEnd.Format("02 January 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(8)
	IncomeHead := []string{"Account Code", "Account Name", "Amount", "Sales"}
	widthHead := []float64{20.0, 115.0, 30.0, 20.0}
	y0 := pdf.GetY()
	for i, head := range IncomeHead {
		pdf.SetY(y0)
		x := 13.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 5, head, "LRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 5, head, "RT", "C", false)
		}

	}
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "SALES AND REVENUE", "LRT", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetX(13.0)
	salesNrevenueData := []IncomeStatementModel{}
	//opex
	operatingExpensesData := []IncomeStatementModel{}
	totOpEx := 0.0
	totSalesOpex := 0.0
	//other revenue
	otherRevenue := []IncomeStatementModel{}
	totOtherRev := 0.0
	totSalesOtherRev := 0.0
	//other expense
	otherExpense := []IncomeStatementModel{}
	totOtherExpense := 0.0
	totSalesOtherEx := 0.0
	//TAX
	Tax := []IncomeStatementModel{}
	totTax := 0.0
	totSalesTax := 0.0
	for _, each := range results {
		if each.Main_Acc_Code == 5000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
			}
			IncomeMod.Sales = Percentage
			salesNrevenueData = append(salesNrevenueData, IncomeMod)
		}
		if each.Main_Acc_Code == 6000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOpex += tot
			}
			totOpEx += each.Ending
			IncomeMod.Sales = Percentage
			operatingExpensesData = append(operatingExpensesData, IncomeMod)
		}
		if each.Main_Acc_Code == 7000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOtherRev += tot
			}
			totOtherRev += each.Ending
			IncomeMod.Sales = Percentage
			otherRevenue = append(otherRevenue, IncomeMod)
		}
		if each.Main_Acc_Code == 8000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesOtherEx += tot
			}
			totOtherExpense += each.Ending
			IncomeMod.Sales = Percentage
			otherExpense = append(otherExpense, IncomeMod)
		}
		if each.Main_Acc_Code == 9000 {
			var IncomeMod IncomeStatementModel
			// IncomeMod := new(IncomeStatementModel)
			IncomeMod.Id = each.ID
			IncomeMod.Acc_Code = each.ACC_Code
			IncomeMod.Acc_Name = each.Account_Name
			IncomeMod.Ending = each.Ending
			Percentage := "0.00%"
			if totSalesAndRevenue != 0 {
				tot := each.Ending / totSalesAndRevenue * 100
				Percentage = strconv.FormatFloat(tot, 'f', 2, 64) + "%"
				totSalesTax += tot
			}
			totTax += each.Ending
			IncomeMod.Sales = Percentage
			Tax = append(Tax, IncomeMod)
		}
	}
	tk.Println(Tax)
	y1 = pdf.GetY()
	for _, each := range salesNrevenueData {
		y1 = pdf.GetY()
		pdf.SetY(y1)

		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}
		// tk.Println(allA)
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL SALES AND REVENUE", "1", "C", false)
	x2 := 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalSR := tk.Sprintf("%.2f", totSalesAndRevenue)
	if totSalesAndRevenue < 0 {
		totalSR = "(" + tk.Sprintf("%.2f", totSalesAndRevenue) + ")"
	}
	pdf.MultiCell(30, 5, totalSR, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	salesNrevenueSales := "100.00%"
	if totSalesAndRevenue == 0 {
		salesNrevenueSales = "0.00%"
	}
	pdf.MultiCell(20, 5, salesNrevenueSales, "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "OPERATING EXPENSES", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range operatingExpensesData {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}
		// tk.Println(allA)
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL OPERATING EXPENSES", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalOpexFix := tk.Sprintf("%.2f", totOpEx)
	if totOpEx < 0 {
		totalOpexFix = "(" + tk.Sprintf("%.2f", totOpEx) + ")"
	}
	pdf.MultiCell(30, 5, totalOpexFix, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOpex, 'f', 2, 64)+"%", "RT", "R", false)
	// pdf.GetY()
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "NET PROFIT", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	netProfit := totSalesAndRevenue - totOpEx
	netProfitFix := tk.Sprintf("%.2f", netProfit)
	if netProfit < 0 {
		netProfitFix = "(" + tk.Sprintf("%.2f", netProfit) + ")"
	}
	pdf.MultiCell(30, 5, netProfitFix, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	pdf.MultiCell(20, 5, strconv.FormatFloat(100-totSalesOpex, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.AddPage()
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "OTHER REVENUE", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range otherRevenue {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}
		// tk.Println(allA)
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL OTHER REVENUE", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totalOR := tk.Sprintf("%.2f", totOtherRev)
	if totOtherRev < 0 {
		totalOR = "(" + tk.Sprintf("%.2f", totOtherRev) + ")"
	}
	pdf.MultiCell(30, 5, totalOR, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherRev, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "OTHER EXPENSES", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range otherExpense {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}
		// tk.Println(allA)
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "TOTAL OTHER EXPENSES", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totOE := tk.Sprintf("%.2f", totOtherExpense)
	if totOtherExpense < 0 {
		totOE = "(" + tk.Sprintf("%.2f", totOtherExpense) + ")"
	}
	pdf.MultiCell(30, 5, totOE, "RT", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherEx, 'f', 2, 64)+"%", "RT", "R", false)
	// pdf.GetY()
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "TOTAL OTHER REVENUE AND EXPENSES", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	allTotRevAndEx := totOtherRev - totOtherExpense
	allTotRevAndExFIX := tk.Sprintf("%.2f", allTotRevAndEx)
	if allTotRevAndEx < 0 {
		allTotRevAndExFIX = "(" + tk.Sprintf("%.2f", allTotRevAndEx) + ")"
	}
	pdf.MultiCell(30, 5, allTotRevAndExFIX, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totSalesOtherRev-totSalesOtherEx, 'f', 2, 64)+"%", "RTB", "R", false)
	y2 = pdf.GetY()
	pdf.SetXY(13.0, y2)
	pdf.MultiCell(135, 5, "EARNING BEFORE TAX", "LRB", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y2)
	earningBeforeTAX := netProfit + allTotRevAndEx
	earningBeforeTAXFIX := tk.Sprintf("%.2f", earningBeforeTAX)
	if earningBeforeTAX < 0 {
		earningBeforeTAXFIX = "(" + tk.Sprintf("%.2f", earningBeforeTAX) + ")"
	}
	pdf.MultiCell(30, 5, earningBeforeTAXFIX, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y2)
	saleseEarningBefore := totSalesOtherRev - totSalesOtherEx
	salesNet := 100 - totSalesOpex
	pdf.MultiCell(20, 5, strconv.FormatFloat(salesNet+saleseEarningBefore, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.GetY()
	pdf.SetX(13.0)
	pdf.MultiCell(185, 7, "TAX", "LRT", "L", false)
	y1 = pdf.GetY()
	for _, each := range Tax {
		y1 = pdf.GetY()
		pdf.SetY(y1)
		x := 13.0
		pdf.SetX(x)
		pdf.MultiCell(widthHead[0], 5, strconv.Itoa(each.Acc_Code), "LRT", "C", false)
		a0 := pdf.GetY()
		x += widthHead[0]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[1], 5, each.Acc_Name, "LRT", "L", false)
		a1 := pdf.GetY()
		x += widthHead[1]
		pdf.SetXY(x, y1)
		ending := tk.Sprintf("%.2f", each.Ending)
		if each.Ending < 0 {
			ending = "(" + tk.Sprintf("%.2f", each.Ending) + ")"
		}
		pdf.MultiCell(widthHead[2], 5, ending, "TR", "R", false)
		a2 := pdf.GetY()
		x += widthHead[2]
		pdf.SetXY(x, y1)
		pdf.MultiCell(widthHead[3], 5, each.Sales, "TR", "R", false)
		a3 := pdf.GetY()
		allA := []float64{a0, a1, a2, a3}
		// tk.Println(allA)
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y1 = pdf.GetY()
	pdf.SetXY(13.0, pdf.GetY())
	pdf.MultiCell(135, 5, "EARNING AFTER TAX", "1", "C", false)
	x2 = 135.0 + 13.0
	pdf.SetXY(x2, y1)
	totEaningAfter := earningBeforeTAX - totTax
	earningAfter := tk.Sprintf("%.2f", totEaningAfter)
	if totEaningAfter < 0 {
		earningAfter = "(" + tk.Sprintf("%.2f", totEaningAfter) + ")"
	}
	pdf.MultiCell(30, 5, earningAfter, "RTB", "R", false)
	x2 = 165.0 + 13.0
	pdf.SetXY(x2, y1)
	pdf.MultiCell(20, 5, strconv.FormatFloat(totEaningAfter/totSalesAndRevenue*100, 'f', 2, 64)+"%", "RTB", "R", false)
	pdf.Ln(5)
	pdf.SetX(150)
	pdf.CellFormat(50, 6, "Approve", "", 0, "C", false, 0, "")
	pdf.Ln(18)
	pdf.GetY()
	pdf.Ln(2)
	pdf.GetY()
	pdf.SetX(150)
	pdf.CellFormat(50, 6, "__________________________", "", 0, "C", false, 0, "")
	namepdf := "-incomestatement.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	e = pdf.OutputFileAndClose(fileName)
	tk.Println(fileName)
	return ""
}
