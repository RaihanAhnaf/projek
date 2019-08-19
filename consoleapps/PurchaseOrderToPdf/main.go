package main

import (
	"eaciit/proactive-dev/helpers"
	. "eaciit/proactive-dev/models"
	db "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/mgo.v2/bson"
	// "os"
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
	// dataPO := bson.ObjectIdHex("59927528c692da2e40d1b449")
	PurchaseOrderPDF(dateStart, dateEnd)
}

func PurchaseOrderPDF(dateStart time.Time, dateEnd time.Time) interface{} {

	// Digunakan
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}

	Ctx := orm.New(conn)
	csr, e := Ctx.Connection.NewQuery().Select().From("PurchaseOrder").Where(db.Eq("_id", bson.ObjectIdHex("59bf662ac54475102c7a8162"))).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []PurchaseOrder{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	tk.Println(DATA)

	if DATA.Currency == "USD" {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceIDR = 0
			DATA.ListDetail[i].AmountIDR = 0
		}
		// discount = DATA.Discount
	} else {
		for i, _ := range DATA.ListDetail {
			DATA.ListDetail[i].PriceUSD = 0
			DATA.ListDetail[i].AmountUSD = 0
		}
	}

	config := helpers.ReadConfig()
	Img := config["imgpath"]

	// PurchaseOrderHeadTable := []string{"Item", "Qty", "Price USD", "Price IDR", "Amount USD", "Amount IDR"}

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
	pdf.CellFormat(0, 15, "PURCHASE ORDER", "", 0, "L", false, 0, "")

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

	pdf.GetY()
	pdf.SetY(y1 + 29)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, "NPWP : 70.842.637.4 - 609.000", "", 0, "L", false, 0, "")
	pdf.SetX(10)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+30, 286, pdf.GetY()+30) //garis horizontal2

	pdf.Line(pdf.GetX(), pdf.GetY()+9, pdf.GetX(), 101)         // garis vertikal 1 dari kiri/156 -> 136
	pdf.Line(pdf.GetX()+276, pdf.GetY()+9, pdf.GetX()+276, 101) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()+138, pdf.GetY()+9, pdf.GetX()+138, 101) // garis vertikal 2 tengah
	pdf.Line(pdf.GetX()+69, pdf.GetY()+30, pdf.GetX()+69, 91)   // garis vertikal 4 tenga2nya field to//101

	y2 := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(y2 + 10)
	pdf.MultiCell(138.5, 5, " TO:", "", "L", false)

	pdf.SetX(9)
	pdf.SetY(y2 + 15)
	pdf.MultiCell(138.5, 5, DATA.SupplierName, "", "L", false)

	pdf.SetY(y2 + 31)
	pdf.SetX(11)
	pdf.MultiCell(138.5, 5, " Payment :", "", "L", false)

	filCash := false
	fillInstallment := false
	if DATA.Payment == "CASH" {
		filCash = true
		fillInstallment = false
	} else {
		fillInstallment = true
		filCash = false
	}
	// tk.Println(filCash)
	pdf.SetY(y2 + 36)
	pdf.SetX(13)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filCash, 0, "")
	pdf.SetY(y2 + 36)
	pdf.SetX(20)
	pdf.MultiCell(138.5, 5, "Cash", "", "L", false)

	pdf.SetY(y2 + 43)
	pdf.SetX(13)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", fillInstallment, 0, "")
	pdf.SetY(y2 + 43)
	pdf.SetX(20)
	pdf.MultiCell(138.5, 5, "Instalment", "", "L", false)

	// p
	// pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)

	pdf.SetY(y2 + 50)
	pdf.SetX(12)
	pdf.MultiCell(50, 5, "Down Payment :", "", "L", false)

	pdf.SetY(y2 + 50)
	pdf.SetX(41)
	pdf.MultiCell(50, 5, DATA.DateStr, "", "L", false)

	//1 table 2 kolom

	pdf.SetY(y2 + 31)
	pdf.SetX(80)
	pdf.MultiCell(138.5, 5, " Type :", "", "L", false)

	filHardware := false
	filService := false
	filCunsumable := false
	filOthers := false
	if DATA.Type == "HARDWARE" {
		filOthers = false
		filHardware = true
		filCunsumable = false
		filService = false
	} else if DATA.Type == "SERVICE" {
		filService = true
		filHardware = false
		filCunsumable = false
		filOthers = false
	} else if DATA.Type == "CONSUMABLE" {
		filCunsumable = true
		filHardware = false
		filService = false
		filOthers = false
	} else {
		filOthers = true
		filHardware = false
		filService = false
		filCunsumable = false
	}

	pdf.SetY(y2 + 36)
	pdf.SetX(82)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filHardware, 0, "")
	pdf.SetY(y2 + 36)
	pdf.SetX(89)
	pdf.MultiCell(138.5, 5, "Hardware", "", "L", false)

	pdf.SetY(y2 + 43)
	pdf.SetX(82)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filService, 0, "")
	pdf.SetY(y2 + 43)
	pdf.SetX(89)
	pdf.MultiCell(138.5, 5, "Service", "", "L", false)

	pdf.SetY(y2 + 36)
	pdf.SetX(110)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filCunsumable, 0, "")
	pdf.SetY(y2 + 36)
	pdf.SetX(117)
	pdf.MultiCell(138.5, 5, "Consumable", "", "L", false)

	pdf.SetY(y2 + 43)
	pdf.SetX(110)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filOthers, 0, "")
	pdf.SetY(y2 + 43)
	pdf.SetX(117)
	pdf.MultiCell(138.5, 5, "Others", "", "L", false)

	pdf.SetY(y2 + 50)
	pdf.SetX(81)
	pdf.MultiCell(138.5, 5, "Currency :", "", "L", false)

	pdf.SetY(y2 + 50)
	pdf.SetX(100)
	pdf.MultiCell(80, 5, DATA.Currency, "", "L", false)

	//TABLE KANAN
	pdf.SetY(y2 + 10)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Order No :", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(168.5)
	pdf.MultiCell(138.5, 5, DATA.DocumentNumber, "", "L", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 15)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Date :", "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetY(y2 + 15)
	pdf.SetX(160.5)
	pdf.MultiCell(136.5, 5, DATA.DateStr, "", "L", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 31)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Place to delivery :", "", "LB", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 36)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " PT. Wiyasa Teknologi Nusantara", "", "LB", false)
	pdf.SetY(pdf.GetY())
	// pdf.Ln(2)

	pdf.SetY(y2 + 43)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Jl. Imam Bonjol No.120, DR. Soetomo - Tegalsari", "", "LB", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 50)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Surabaya, Indonesia", "", "LB", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 55)
	pdf.SetX(10)
	pdf.MultiCell(15, 10, "No.", "BRT", "C", false)
	pdf.SetY(pdf.GetY())

	// y3 := pdf.GetY()
	pdf.SetFont("Arial", "", 11)
	pdf.SetY(y2 + 55)
	pdf.SetX(25)
	pdf.MultiCell(100, 10, "Item & Specification", "BT", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 55)
	pdf.SetX(125)
	pdf.MultiCell(23, 10, "Qty", "LTB", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 55)
	pdf.SetX(148)
	pdf.MultiCell(50, 5, "Price / Unit", "TBR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 60)
	pdf.SetX(148)
	pdf.MultiCell(20, 5, "USD", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 60)
	pdf.SetX(168)
	pdf.MultiCell(30, 5, "IDR", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 55)
	pdf.SetX(198)
	pdf.MultiCell(50, 5, "Amount", "TB", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 60)
	pdf.SetX(198)
	pdf.MultiCell(20, 5, "USD", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 60)
	pdf.SetX(218)
	pdf.MultiCell(30, 5, "IDR", "B", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y2 + 55)
	pdf.SetX(248)
	pdf.MultiCell(38, 10, "Remark", "LTB", "C", false)
	pdf.SetY(pdf.GetY())
	pdf.SetFont("Arial", "", 10)
	pdf.SetY(y2 + 65)
	pdf.SetX(248)
	pdf.MultiCell(38, 5, DATA.Remark, "LRBT", "C", false)
	pdf.SetY(pdf.GetY())
	// ytable := y2 + 65
	pdf.SetY(y2 + 65)
	yline1 := y2 + 65
	// tk.Println("lalal", yline1)
	pdf.SetFont("Arial", "", 10)
	for i, list := range DATA.ListDetail {
		y4 := pdf.GetY()
		// pdf.SetY(y2 + 65)
		pdf.SetX(9)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(15, 5, numberstr, "", "C", false)

		pdf.SetY(y4)
		pdf.SetX(25)
		item := list.Item
		pdf.MultiCell(100, 5, item, "", "L", false)
		a0 := pdf.GetY()

		pdf.SetY(y4)
		pdf.SetX(125)
		qty := strconv.Itoa(list.Qty)
		pdf.MultiCell(23, 5, qty, "", "C", false)
		a1 := pdf.GetY()

		pdf.SetY(y4)
		pdf.SetX(148)
		priceusd := tk.Sprintf("%.2f", list.PriceUSD)
		if priceusd == "0.00" {
			priceusd = ""
		}
		pdf.MultiCell(20, 5, priceusd, "0", "R", false)
		a2 := pdf.GetY()

		pdf.SetY(y4)
		pdf.SetX(168)
		priceidr := tk.Sprintf("%.2f", list.PriceIDR)
		if priceidr == "0.00" {
			priceidr = ""
		}
		pdf.MultiCell(30, 5, priceidr, "0", "R", false)
		a3 := pdf.GetY()

		pdf.SetY(y4)
		pdf.SetX(198)
		amountusd := tk.Sprintf("%.2f", list.AmountUSD)
		if amountusd == "0.00" {
			amountusd = ""
		}
		pdf.MultiCell(20, 5, amountusd, "0", "R", false)
		a4 := pdf.GetY()
		// pdf.SetY(pdf.GetY())

		pdf.SetY(y4)
		pdf.SetX(218)
		amountidr := tk.Sprintf("%.2f", list.AmountIDR)
		if amountidr == "0.00" {
			amountidr = ""
		}
		pdf.MultiCell(30, 5, amountidr, "0", "R", false)
		a5 := pdf.GetY()
		// tk.Println(a5)
		// pdf.SetY(pdf.GetY())

		allA := []float64{a0, a1, a2, a3, a4, a5}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	y3 := pdf.GetY()
	// tk.Println(yline1, y3)

	// pdf.Line(pdf.GetX(), yline1, pdf.GetX(), y3)         //garis vertical No KIRI
	pdf.Line(pdf.GetX()+15, yline1, pdf.GetX()+15, y3)   //garis vertical No KANAN
	pdf.Line(pdf.GetX()+115, yline1, pdf.GetX()+115, y3) //garis itemspe Kanan
	pdf.Line(pdf.GetX()+138, yline1, pdf.GetX()+138, y3) //qty kanan
	pdf.Line(pdf.GetX()+158, yline1, pdf.GetX()+158, y3) // usd price kanan
	pdf.Line(pdf.GetX()+188, yline1, pdf.GetX()+188, y3) // idr price kanan
	pdf.Line(pdf.GetX()+208, yline1, pdf.GetX()+208, y3) // usd amount kanan
	pdf.Line(pdf.GetX()+238, yline1, pdf.GetX()+238, y3) // idr price kanan
	// pdf.Line(pdf.GetX()+280, yline1, pdf.GetX()+280, y3) // vertical remark kanan

	pdf.Line(pdf.GetX(), y3, 248, y3) //garis horizontal kiri No

	//Sub Total
	y3 = pdf.GetY()
	pdf.SetY(y3)
	pdf.SetX(148)
	pdf.MultiCell(20, 5, "Sub Total", "LBR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(168)
	pdf.MultiCell(30, 5, "", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(198)
	total := tk.Sprintf("%.2f", DATA.TotalUSD)
	if DATA.Currency == "IDR" {
		total = ""
	}
	pdf.MultiCell(20, 5, total, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(218)
	total = tk.Sprintf("%.2f", DATA.TotalIDR)
	if DATA.Currency == "USD" {
		total = ""
	}
	pdf.MultiCell(30, 5, total, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	//Diskon
	y3 = pdf.GetY()
	pdf.SetY(y3)
	pdf.SetX(148)
	pdf.MultiCell(20, 5, "Disc.", "LBR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(168)
	pdf.MultiCell(30, 5, "", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(198)
	diskon := tk.Sprintf("%.2f", DATA.Discount)
	if DATA.Currency == "IDR" {
		diskon = ""
	}
	pdf.MultiCell(20, 5, diskon, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(218)
	diskon = tk.Sprintf("%.2f", DATA.Discount)
	if DATA.Currency == "USD" {
		diskon = ""
	}
	pdf.MultiCell(30, 5, diskon, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	//Data VAT
	y3 = pdf.GetY()
	pdf.SetY(y3)
	pdf.SetX(148)
	pdf.MultiCell(20, 5, "VAT 10%", "LBR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(168)
	pdf.MultiCell(30, 5, "", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(198)
	vat := tk.Sprintf("%.2f", DATA.VAT)
	if DATA.Currency == "IDR" {
		vat = ""
	}
	pdf.MultiCell(20, 5, vat, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(218)
	vat = tk.Sprintf("%.2f", DATA.VAT)
	if DATA.Currency == "USD" {
		vat = ""
	}
	pdf.MultiCell(30, 5, vat, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	//GranTotal
	y3 = pdf.GetY()
	pdf.SetY(y3)
	pdf.SetX(148)
	pdf.MultiCell(20, 5, "Total", "LBR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(168)
	pdf.MultiCell(30, 5, "", "BR", "C", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(198)
	grandTotal := tk.Sprintf("%.2f", DATA.GrandTotalUSD)
	if DATA.Currency == "IDR" {
		grandTotal = ""
	}
	pdf.MultiCell(20, 5, grandTotal, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	pdf.SetY(y3)
	pdf.SetX(218)
	grandTotal = tk.Sprintf("%.2f", DATA.GrandTotalIDR)
	if DATA.Currency == "USD" {
		grandTotal = ""
	}
	pdf.MultiCell(30, 5, grandTotal, "BR", "R", false)
	pdf.SetY(pdf.GetY())

	// ysig := y3
	pdf.SetY(y3)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(15, 12, "PT. Wiyasa Teknologi Nusantara", "", 0, "L", false, 0, "")
	pdf.Ln(20)

	pdf.SetY(y3)
	pdf.SetX(89)
	pdf.CellFormat(15, 12, "Approved By,", "", 0, "L", false, 0, "")
	pdf.Ln(20)

	y3 = pdf.GetY()
	pdf.SetY(y3)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(15, 12, "Andik Suncahyo", "", 0, "L", false, 0, "")

	pdf.SetY(y3)
	pdf.SetX(89)
	pdf.CellFormat(15, 12, "Arif Darmawan", "", 0, "L", false, 0, "")

	// pdf.CellFormat(3, 3, "", "1", 0, "L", true, 0, "")
	// pdf.CellFormat(30, 10, "Document Number ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, DATA.DocumentNumber, "", 0, "L", false, 0, "")
	// pdf.CellFormat(20, 10, "Date", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(35, 10, DATA.DateStr, "", 0, "L", false, 0, "")
	// pdf.Ln(5)

	// pdf.GetY()
	// pdf.SetX(12)
	// pdf.CellFormat(30, 10, "Supplier Code ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, DATA.SupplierCode, "", 0, "L", false, 0, "")
	// pdf.CellFormat(20, 10, "Payment", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(35, 10, DATA.Payment, "", 0, "L", false, 0, "")
	// pdf.Ln(5)

	// pdf.GetY()
	// pdf.SetX(12)
	// pdf.CellFormat(30, 10, "Supplier Name ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, DATA.SupplierName, "", 0, "L", false, 0, "")
	// pdf.CellFormat(20, 10, "Type", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(35, 10, DATA.Type, "", 0, "L", false, 0, "")
	// pdf.Ln(5)

	// pdf.GetY()
	// pdf.SetX(12)
	// pdf.CellFormat(30, 10, "Account Code ", "", 0, "L", false, 0, "")
	// pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// pdf.CellFormat(80, 10, DATA.AccountCode, "", 0, "L", false, 0, "")

	// // vat := tk.Sprintf("%.2f", DATA.VAT)
	// // pdf.CellFormat(20, 10, "VAT ", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(35, 10, vat, "", 0, "L", false, 0, "")
	// pdf.Ln(10)

	// y3 := pdf.GetY()
	// widthHead := []float64{50.0, 15.0, 30.0, 30.0, 30.0, 30.0}
	// for i, head := range PurchaseOrderHeadTable {
	// 	pdf.SetY(y3)
	// 	x := 13.0
	// 	for y, z := range widthHead {
	// 		if i > y {
	// 			x += z
	// 		} else {
	// 			x += 0.0
	// 		}
	// 	}
	// 	pdf.SetX(x)
	// 	if i == 0 {
	// 		pdf.MultiCell(widthHead[i], 7, head, "LRT", "C", false)
	// 	} else {
	// 		pdf.MultiCell(widthHead[i], 7, head, "RT", "C", false)
	// 	}

	// }

	// y4 := pdf.GetY()

	// for _, list := range DATA.ListDetail {
	// 	y4 = pdf.GetY()
	// 	pdf.SetY(y4)
	// 	x := 13.0
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	item := list.Item
	// 	pdf.MultiCell(widthHead[0], 4, item, "LT", "C", false)
	// 	pdf.SetY(y4)

	// 	x += widthHead[0]
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	qty := strconv.Itoa(list.Qty)
	// 	pdf.MultiCell(widthHead[1], 4, qty, "T", "C", false)
	// 	pdf.SetY(y4)

	// 	x = 13.0 + widthHead[0] + widthHead[1]
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	priceusd := tk.Sprintf("%.2f", list.PriceUSD)
	// 	pdf.MultiCell(widthHead[2], 4, priceusd, "T", "C", false)
	// 	pdf.SetY(y4)

	// 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2]
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	priceidr := tk.Sprintf("%.2f", list.PriceIDR)
	// 	pdf.MultiCell(widthHead[3], 4, priceidr, "T", "C", false)
	// 	pdf.SetY(y4)

	// 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	amountusd := tk.Sprintf("%.2f", list.AmountUSD)
	// 	pdf.MultiCell(widthHead[3], 4, amountusd, "T", "C", false)
	// 	pdf.SetY(y4)

	// 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// 	pdf.SetLeftMargin(x)
	// 	pdf.SetX(x)
	// 	amountidr := tk.Sprintf("%.2f", list.AmountIDR)
	// 	pdf.MultiCell(widthHead[3], 4, amountidr, "TR", "C", false)
	// 	pdf.SetY(y4)
	// }

	// pdf.Ln(4)
	// y4 = pdf.GetY()
	// pdf.SetY(y4)
	// x := 13.0
	// pdf.SetLeftMargin(x)
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[0]+widthHead[1], 4, "Remark : ", "T", "L", false)
	// pdf.SetY(y4)
	// x = 33.0
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[5], 4, DATA.Remark, "", "R", false)

	// pdf.SetY(y4)
	// x = 13.0 + widthHead[0] + widthHead[1]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "Total :", "T", "R", false)
	// pdf.SetY(y4)
	// total := tk.Sprintf("%.2f", DATA.Total)
	// x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[5], 4, total, "T", "C", false)
	// pdf.SetY(y4)

	// pdf.Ln(4)
	// y4 = pdf.GetY()
	// pdf.SetY(y4)
	// x = 13.0 + widthHead[0] + widthHead[1]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "VAT :", "0", "R", false)
	// pdf.SetY(y4)
	// vat := tk.Sprintf("%.2f", DATA.VAT)
	// x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[5], 4, vat, "", "C", false)
	// pdf.SetY(y4)

	// pdf.Ln(4)
	// y4 = pdf.GetY()
	// pdf.SetY(y4)
	// x = 13.0 + widthHead[0] + widthHead[1]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "Diskon :", "0", "R", false)
	// pdf.SetY(y4)
	// diskon := tk.Sprintf("%.2f", DATA.Discount)
	// x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[5], 4, diskon, "", "C", false)
	// pdf.SetY(y4)
	// // pdf.SetY(y4)

	// pdf.Ln(4)
	// y4 = pdf.GetY()
	// pdf.SetY(y4)
	// x = 13.0 + widthHead[0] + widthHead[1]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "GrandTotal :", "0", "R", false)
	// pdf.SetY(y4)
	// grandTotal := tk.Sprintf("%.2f", DATA.GrandTotal)
	// x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// pdf.SetX(x)
	// pdf.MultiCell(widthHead[5], 4, grandTotal, "", "C", false)
	// pdf.SetY(y4)

	namepdf := "-purchaseorder.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	e = pdf.OutputFileAndClose(fileName)
	tk.Println(e)
	if e != nil {
		e.Error()
	}
	return ""
}
