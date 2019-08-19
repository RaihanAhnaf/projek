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
	"strings"
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
	InvoicePDF(dateStart, dateEnd)
}

func InvoicePDF(dateStart time.Time, dateEnd time.Time) interface{} {

	// Digunakan
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}

	Ctx := orm.New(conn)
	csr, e := Ctx.Connection.NewQuery().Select().From("Invoice").Where(db.Eq("_id", bson.ObjectIdHex("59e4686692de4d4e9098ce3a"))).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []InvoiceModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	tk.Println(DATA)

	if DATA.Currency == "USD" {
		for i, _ := range DATA.ListItem {
			DATA.ListItem[i].PriceIDR = 0
			DATA.ListItem[i].AmountIDR = 0
		}
		// discount = DATA.Discount
	} else {
		for i, _ := range DATA.ListItem {
			DATA.ListItem[i].PriceUSD = 0
			DATA.ListItem[i].AmountUSD = 0
		}
	}
	csr, e = Ctx.Connection.NewQuery().Select().From("Customer").Where(db.Eq("Kode", DATA.CustomerCode)).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	resultsCustomer := []CustomerModel{}
	e = csr.Fetch(&resultsCustomer, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	cust := resultsCustomer[0]
	config := helpers.ReadConfig()
	Img := config["imgpath"]
	fonth := config["fonthpath"]
	// PurchaseOrderHeadTable := []string{"Item", "Qty", "Price USD", "Price IDR", "Amount USD", "Amount IDR"}

	pdf := gofpdf.New("P", "mm", "A4", fonth)
	pdf.AddFont("Century_Gothic", "", "Century_Gothic.json")
	pdf.AddFont("Century_Gothicb", "B", "Century_Gothicb.json")
	pdf.AddPage()
	pdf.SetXY(8, 5)
	x_defaulft := 28.0
	pdf.SetFont("Century_Gothic", "", 10)
	// pdf.Image(imageNameStr, x, y, w, h, flow, tp, link, linkStr)
	pdf.Image(Img+"eaciit-logo.png", 55, 10, 13, 14, false, "", 0, "")
	pdf.Ln(2)
	pdf.SetY(23)
	pdf.SetX(x_defaulft)
	pdf.CellFormat(0, 12, "PT. WIYASA TEKNOLOGI NUSANTARA", "", 0, "L", false, 0, "")
	pdf.SetY(pdf.GetY() + 15)
	pdf.SetX(x_defaulft)
	pdf.SetFont("Century_Gothicb", "B", 20)
	pdf.CellFormat(0, 12, "INVOICE", "", 0, "L", false, 0, "")
	pdf.SetFont("Century_Gothicb", "B", 15)
	pdf.Ln(6)
	pdf.SetX(x_defaulft)
	pdf.CellFormat(0, 12, DATA.DocumentNo, "", 0, "L", false, 0, "")
	pdf.SetY(pdf.GetY())
	pdf.Ln(5)
	pdf.SetFont("Century_Gothicb", "B", 20)
	pdf.SetX(x_defaulft)
	pdf.SetLineWidth(1)
	pdf.Line(28, pdf.GetY()+5, 182, pdf.GetY()+5)
	// pdf.CellFormat(0, 1, "_______________________________________", "", 0, "L", false, 0, "")
	pdf.SetLineWidth(0)
	pdf.SetFont("Century_Gothicb", "B", 11)
	pdf.Ln(5)
	y := pdf.GetY() + 5.0
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(80, 5, "Prepared", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(80, 5, "To", "", "L", false)
	pdf.SetFont("Century_Gothic", "", 10)
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(80, 5, "Andiek Suncahyo", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(80, 5, cust.Name, "", "L", false) // customer name
	// pdf.MultiCell(80, 5, "Eaciit Vyasa Pte.Ltd", "", "L", false) // customer name
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(80, 5, "Jl. Imam Bonjol No. 120", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(80, 5, cust.Address, "", "L", false) //address
	// pdf.MultiCell(80, 5, "One Raffles Place #41-01", "", "L", false) //address
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(80, 5, "Tegalsari - Surabaya", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(80, 5, cust.City, "", "L", false) //city
	// pdf.MultiCell(80, 5, "One Raffles Place", "", "L", false) //city
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(80, 5, "031-5676223", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(80, 5, cust.NoTelp, "", "L", false) //telp
	// pdf.MultiCell(80, 5, "Singapore 048616", "", "L", false) //telp
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(40, 5, "NPWP ", "", "L", false)
	pdf.SetXY(48, y)
	pdf.MultiCell(80, 5, ": 70.842.637.4-609.000", "", "L", false)
	pdf.SetXY(125, y)
	pdf.MultiCell(40, 5, "NPWP ", "", "L", false)
	pdf.SetXY(140, y)
	pdf.MultiCell(80, 5, ": "+cust.NPWP, "", "L", false)
	pdf.Ln(7)
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(40, 5, "Date ", "", "L", false)
	pdf.SetXY(48, y)
	pdf.MultiCell(40, 5, ": "+DATA.DateCreated.Format("Jan 2, 2006 "), "", "L", false)
	y = pdf.GetY()

	invHead := []string{"No", "Item", "Amount(IDR)"}
	widHead := []float64{10.0, 100.0, 44.0}
	yt1 := pdf.GetY()
	pdf.SetDrawColor(253, 77, 0)
	pdf.Line(28, yt1, 182, yt1) //horizontal
	for i, head := range invHead {
		pdf.SetY(y)
		x := x_defaulft
		for j, w := range widHead {
			if i > j {
				x += w
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)

		// if i == 2 {
		pdf.MultiCell(widHead[i], 5, head, "", "C", false)
		// } else {
		// 	pdf.MultiCell(widHead[i], 5, head, "RTB", "C", false)
		// }
	}
	yt2 := pdf.GetY()
	pdf.Line(28, yt2, 182, yt2) //horizontal
	for i, list := range DATA.ListItem {
		y4 := pdf.GetY()
		// pdf.SetY(y2 + 65)
		pdf.SetX(x_defaulft)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(10, 5, numberstr, "", "C", false)

		pdf.SetY(y4)
		pdf.SetX(38)
		item := list.Item
		pdf.MultiCell(100, 5, item, "", "L", false)
		a0 := pdf.GetY()
		pdf.SetY(y4)
		pdf.SetX(138)
		// ConvertToCurrency(amountidr)
		amountidr := tk.Sprintf("%.2f", list.AmountIDR)
		if DATA.Currency == "USD" {
			amount := list.AmountUSD * DATA.Rate
			amountidr = tk.Sprintf("%.2f", amount)
		}
		pdf.MultiCell(44, 5, ConvertToCurrency(amountidr), "0", "R", false)
		a2 := pdf.GetY()
		// tk.Println(a5)
		// pdf.SetY(pdf.GetY())

		allA := []float64{a0, a2}
		var n, biggest float64
		for _, v := range allA {
			if v > n {
				n = v
				biggest = n
			}
		}
		pdf.SetY(biggest)
	}
	pdf.Ln(10)
	yt3 := pdf.GetY()
	pdf.Line(28, yt3, 182, yt3)  //horizontal
	pdf.Line(38, yt1, 38, yt3)   //vertical
	pdf.Line(138, yt1, 138, yt3) //vertical
	y = pdf.GetY()
	pdf.SetFont("Century_Gothicb", "B", 10)
	pdf.SetY(y)
	pdf.SetX(100)
	pdf.MultiCell(38, 5, "Total", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(138)
	totalIDR := tk.Sprintf("%.2f", DATA.Total)
	vat := tk.Sprintf("%.2f", DATA.VAT)
	if DATA.Currency == "USD" {
		tot := DATA.Total * DATA.Rate
		totalIDR = tk.Sprintf("%.2f", tot)
		datVat := DATA.VAT * DATA.Rate
		vat = tk.Sprintf("%.2f", datVat)
	}
	pdf.MultiCell(44, 5, ConvertToCurrency(totalIDR), "", "R", false)
	y = pdf.GetY()
	pdf.SetY(y)
	pdf.SetX(100)
	pdf.MultiCell(38, 5, "PPN 10%", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(138)
	pdf.MultiCell(44, 5, ConvertToCurrency(vat), "", "R", false)
	y = pdf.GetY()
	pdf.SetY(y)
	pdf.SetX(100)
	pdf.MultiCell(38, 5, "Total Due (IDR)", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(138)
	pdf.MultiCell(44, 5, ConvertToCurrency(tk.Sprintf("%.2f", DATA.GrandTotalIDR)), "", "R", false)
	y = pdf.GetY()
	pdf.SetY(y)
	pdf.SetX(100)
	pdf.MultiCell(38, 5, "Total Due (USD)", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(138)
	pdf.MultiCell(44, 5, "$ "+ConvertToCurrency(tk.Sprintf("%.2f", DATA.GrandTotalUSD)), "", "R", false)
	y = pdf.GetY()
	pdf.Ln(5)
	pdf.SetFont("Century_Gothic", "", 10)
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(20, 5, "Say:", "", "L", false)
	pdf.SetDrawColor(253, 77, 0)
	pdf.SetXY(48, y)
	pdf.MultiCell(134, 5, "money to word", "1", "L", false)
	// pdf.MultiCell(134, 5, DATA.Description, "1", "L", false)
	pdf.Ln(5)
	pdf.SetX(x_defaulft)
	pdf.SetFont("Century_Gothicb", "B", 10)
	pdf.MultiCell(134, 5, "All payment can be made to", "", "L", false)
	pdf.GetY()
	pdf.SetFont("Century_Gothic", "", 10)
	pdf.SetX(x_defaulft)
	pdf.MultiCell(134, 5, "Bank Mandiri Branch Margorejo Surabaya", "", "L", false)
	pdf.GetY()
	pdf.SetX(x_defaulft)
	y = pdf.GetY()
	pdf.MultiCell(40, 5, "A/C No", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(48)
	pdf.MultiCell(40, 5, ": 1420013617674", "", "L", false)
	pdf.GetY()
	y = pdf.GetY()
	pdf.SetX(x_defaulft)
	pdf.MultiCell(40, 5, "Name", "", "L", false)
	pdf.SetY(y)
	pdf.SetX(48)
	pdf.MultiCell(100, 5, ": PT WIYASA TEKNOLOGI NUSANTARA", "", "L", false)
	y = pdf.GetY()
	pdf.SetXY(x_defaulft, y)
	pdf.MultiCell(20, 5, "Desc:", "", "L", false)
	pdf.SetDrawColor(253, 77, 0)
	pdf.SetXY(48, y)
	pdf.MultiCell(134, 5, DATA.Description, "1", "L", false)
	pdf.GetY()
	pdf.SetX(138)
	pdf.MultiCell(44, 5, "Director", "", "C", false)
	pdf.Ln(20)
	pdf.SetX(138)
	pdf.SetFont("Arial", "BU", 10)
	pdf.MultiCell(44, 5, "Andiek Suncahyo", "", "C", false)
	// pdf.CellFormat(134, 5, DATA.Description, "1", 0, "L", false, 0, "")
	// pdf.Line(28, pdf.GetY(), 182, pdf.GetY())
	// pdf.CellFormat(0, 15, "Jl. Imam Bonjol No. 120, DR.Soetomo - Tegalsari", "", 0, "L", false, 0, "")
	// pdf.SetX(180)
	// pdf.SetFont("Arial", "B", 14)
	// pdf.CellFormat(0, 15, "INVOICE", "", 0, "L", false, 0, "")

	// pdf.SetY(y1 + 17)
	// pdf.SetX(11)
	// pdf.SetFont("Arial", "", 11)
	// pdf.CellFormat(0, 12, "Surabaya-Indonesia", "", 0, "L", false, 0, "")
	// pdf.SetX(10)
	// // pdf.Ln(1)

	// pdf.SetY(y1 + 23)
	// pdf.SetX(11)
	// pdf.SetFont("Arial", "", 11)
	// pdf.CellFormat(0, 12, "Telp. 0351-5676223", "", 0, "L", false, 0, "")
	// pdf.SetX(10)

	// pdf.GetY()
	// pdf.SetY(y1 + 29)
	// pdf.SetX(11)
	// pdf.SetFont("Arial", "", 11)
	// pdf.CellFormat(0, 12, "NPWP : 70.842.637.4 - 609.000", "", 0, "L", false, 0, "")
	// pdf.SetX(10)

	// pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	// pdf.Line(pdf.GetX(), pdf.GetY()+30, 286, pdf.GetY()+30) //garis horizontal2

	// pdf.Line(pdf.GetX(), pdf.GetY()+9, pdf.GetX(), 101)         // garis vertikal 1 dari kiri/156 -> 136
	// pdf.Line(pdf.GetX()+276, pdf.GetY()+9, pdf.GetX()+276, 101) // garis vertikal 3 paling kanan
	// pdf.Line(pdf.GetX()+138, pdf.GetY()+9, pdf.GetX()+138, 101) // garis vertikal 2 tengah
	// // pdf.Line(pdf.GetX()+69, pdf.GetY()+30, pdf.GetX()+69, 91)   // garis vertikal 4 tenga2nya field to//101

	// y2 := pdf.GetY()
	// pdf.SetX(9)
	// pdf.SetY(y2 + 10)
	// pdf.MultiCell(138.5, 5, " TO:", "", "L", false)

	// pdf.SetX(9)
	// pdf.SetY(y2 + 15)
	// pdf.MultiCell(138.5, 5, DATA.CustomerName, "", "L", false)

	// pdf.SetY(y2 + 31)
	// pdf.SetX(11)
	// pdf.MultiCell(138.5, 5, " Account Number :", "", "L", false)

	// filSales := false
	// fillRevenue := false
	// if DATA.AccountName == "SALES" {
	// 	filSales = true
	// 	fillRevenue = false
	// } else {
	// 	fillRevenue = true
	// 	filSales = false
	// }
	// // tk.Println(filCash)
	// pdf.SetY(y2 + 36)
	// pdf.SetX(13)
	// pdf.CellFormat(4.5, 5, "", "1", 0, "L", filSales, 0, "")
	// pdf.SetY(y2 + 36)
	// pdf.SetX(20)
	// pdf.MultiCell(138.5, 5, "Sales", "", "L", false)

	// pdf.SetY(y2 + 43)
	// pdf.SetX(13)
	// pdf.CellFormat(4.5, 5, "", "1", 0, "L", fillRevenue, 0, "")
	// pdf.SetY(y2 + 43)
	// pdf.SetX(20)
	// pdf.MultiCell(138.5, 5, "Revenue", "", "L", false)

	// // p
	// // pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)

	// pdf.SetY(y2 + 50)
	// pdf.SetX(12)
	// pdf.MultiCell(50, 5, "Currency :", "", "L", false)

	// pdf.SetY(y2 + 50)
	// pdf.SetX(32)
	// pdf.MultiCell(50, 5, DATA.Currency, "", "L", false)

	// //1 table 2 kolom

	// // pdf.SetY(y2 + 31)
	// // pdf.SetX(80)
	// // pdf.MultiCell(138.5, 5, " Type :", "", "L", false)

	// // filHardware := false
	// // filService := false
	// // filCunsumable := false
	// // filOthers := false
	// // if DATA.Type == "HARDWARE" {
	// // 	filOthers = false
	// // 	filHardware = true
	// // 	filCunsumable = false
	// // 	filService = false
	// // } else if DATA.Type == "SERVICE" {
	// // 	filService = true
	// // 	filHardware = false
	// // 	filCunsumable = false
	// // 	filOthers = false
	// // } else if DATA.Type == "CONSUMABLE" {
	// // 	filCunsumable = true
	// // 	filHardware = false
	// // 	filService = false
	// // 	filOthers = false
	// // } else {
	// // 	filOthers = true
	// // 	filHardware = false
	// // 	filService = false
	// // 	filCunsumable = false
	// // }

	// // pdf.SetY(y2 + 36)
	// // pdf.SetX(82)
	// // pdf.CellFormat(4.5, 5, "", "1", 0, "L", filHardware, 0, "")
	// // pdf.SetY(y2 + 36)
	// // pdf.SetX(89)
	// // pdf.MultiCell(138.5, 5, "Hardware", "", "L", false)

	// // pdf.SetY(y2 + 43)
	// // pdf.SetX(82)
	// // pdf.CellFormat(4.5, 5, "", "1", 0, "L", filService, 0, "")
	// // pdf.SetY(y2 + 43)
	// // pdf.SetX(89)
	// // pdf.MultiCell(138.5, 5, "Service", "", "L", false)

	// // pdf.SetY(y2 + 36)
	// // pdf.SetX(110)
	// // pdf.CellFormat(4.5, 5, "", "1", 0, "L", filCunsumable, 0, "")
	// // pdf.SetY(y2 + 36)
	// // pdf.SetX(117)
	// // pdf.MultiCell(138.5, 5, "Consumable", "", "L", false)

	// // pdf.SetY(y2 + 43)
	// // pdf.SetX(110)
	// // pdf.CellFormat(4.5, 5, "", "1", 0, "L", filOthers, 0, "")
	// // pdf.SetY(y2 + 43)
	// // pdf.SetX(117)
	// // pdf.MultiCell(138.5, 5, "Others", "", "L", false)

	// // pdf.SetY(y2 + 50)
	// // pdf.SetX(81)
	// // pdf.MultiCell(138.5, 5, "Currency :", "", "L", false)

	// // pdf.SetY(y2 + 50)
	// // pdf.SetX(100)
	// // pdf.MultiCell(80, 5, DATA.Currency, "", "L", false)

	// //TABLE KANAN
	// pdf.SetY(y2 + 10)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " Order No :", "", "L", false)
	// pdf.SetY(y2 + 10)
	// pdf.SetX(168.5)
	// pdf.MultiCell(138.5, 5, DATA.DocumentNo, "", "L", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 15)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " Date :", "", "L", false)
	// pdf.SetY(pdf.GetY())
	// pdf.SetY(y2 + 15)
	// pdf.SetX(160.5)
	// date := DATA.DateCreated.Format("01-02-2006")
	// pdf.MultiCell(136.5, 5, date, "", "L", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 31)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " Place to delivery :", "", "LB", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 36)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " PT. Wiyasa Teknologi Nusantara", "", "LB", false)
	// pdf.SetY(pdf.GetY())
	// // pdf.Ln(2)

	// pdf.SetY(y2 + 43)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " Jl. Imam Bonjol No.120, DR. Soetomo - Tegalsari", "", "LB", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 50)
	// pdf.SetX(147.5)
	// pdf.MultiCell(138.5, 5, " Surabaya, Indonesia", "", "LB", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 55)
	// pdf.SetX(10)
	// pdf.MultiCell(15, 10, "No.", "RBT", "C", false)
	// pdf.SetY(pdf.GetY())

	// // y3 := pdf.GetY()
	// pdf.SetFont("Arial", "", 11)
	// pdf.SetY(y2 + 55)
	// pdf.SetX(25)
	// pdf.MultiCell(100, 10, "Item & Specification", "BT", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 55)
	// pdf.SetX(125)
	// pdf.MultiCell(23, 10, "Qty", "LTB", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 55)
	// pdf.SetX(148)
	// pdf.MultiCell(50, 5, "Price / Unit", "TBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 60)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "USD", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 60)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 5, "IDR", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 55)
	// pdf.SetX(198)
	// pdf.MultiCell(50, 5, "Amount", "TB", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 60)
	// pdf.SetX(198)
	// pdf.MultiCell(20, 5, "USD", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 60)
	// pdf.SetX(218)
	// pdf.MultiCell(30, 5, "IDR", "B", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y2 + 55)
	// pdf.SetX(248)
	// pdf.MultiCell(38, 10, "Description", "LBT", "C", false)
	// pdf.SetY(pdf.GetY())
	// pdf.SetFont("Arial", "", 10)
	// pdf.SetY(y2 + 65)
	// pdf.SetX(248)
	// pdf.MultiCell(38, 5, DATA.Description, "LRB", "C", false)
	// pdf.SetY(pdf.GetY())
	// // ytable := y2 + 65
	// pdf.SetY(y2 + 65)
	// yline1 := y2 + 65
	// // tk.Println("lalal", yline1)
	// pdf.SetFont("Arial", "", 10)
	// for i, list := range DATA.ListItem {
	// 	y4 := pdf.GetY()
	// 	// pdf.SetY(y2 + 65)
	// 	pdf.SetX(9)
	// 	number := i + 1
	// 	numberstr := strconv.Itoa(number)
	// 	pdf.MultiCell(15, 5, numberstr, "", "C", false)

	// 	pdf.SetY(y4)
	// 	pdf.SetX(25)
	// 	item := list.Item
	// 	pdf.MultiCell(100, 5, item, "", "L", false)
	// 	a0 := pdf.GetY()

	// 	pdf.SetY(y4)
	// 	pdf.SetX(125)
	// 	qty := strconv.Itoa(list.Qty)
	// 	pdf.MultiCell(23, 5, qty, "", "C", false)
	// 	a1 := pdf.GetY()

	// 	pdf.SetY(y4)
	// 	pdf.SetX(148)
	// 	priceusd := tk.Sprintf("%.2f", list.PriceUSD)
	// 	if priceusd == "0.00" {
	// 		priceusd = ""
	// 	}
	// 	pdf.MultiCell(20, 5, priceusd, "0", "R", false)
	// 	a2 := pdf.GetY()

	// 	pdf.SetY(y4)
	// 	pdf.SetX(168)
	// 	priceidr := tk.Sprintf("%.2f", list.PriceIDR)
	// 	if priceidr == "0.00" {
	// 		priceidr = ""
	// 	}
	// 	pdf.MultiCell(30, 5, priceidr, "0", "R", false)
	// 	a3 := pdf.GetY()

	// 	pdf.SetY(y4)
	// 	pdf.SetX(198)
	// 	amountusd := tk.Sprintf("%.2f", list.AmountUSD)
	// 	if amountusd == "0.00" {
	// 		amountusd = ""
	// 	}
	// 	pdf.MultiCell(20, 5, amountusd, "0", "R", false)
	// 	a4 := pdf.GetY()
	// 	// pdf.SetY(pdf.GetY())

	// 	pdf.SetY(y4)
	// 	pdf.SetX(218)
	// 	amountidr := tk.Sprintf("%.2f", list.AmountIDR)
	// 	if amountidr == "0.00" {
	// 		amountidr = ""
	// 	}
	// 	pdf.MultiCell(30, 5, amountidr, "0", "R", false)
	// 	a5 := pdf.GetY()
	// 	// tk.Println(a5)
	// 	// pdf.SetY(pdf.GetY())

	// 	allA := []float64{a0, a1, a2, a3, a4, a5}
	// 	var n, biggest float64
	// 	for _, v := range allA {
	// 		if v > n {
	// 			n = v
	// 			biggest = n
	// 		}
	// 	}
	// 	pdf.SetY(biggest)
	// }
	// y3 := pdf.GetY()
	// // tk.Println(yline1, y3)

	// pdf.Line(pdf.GetX(), yline1, pdf.GetX(), y3)         //garis vertical No KIRI
	// pdf.Line(pdf.GetX()+15, yline1, pdf.GetX()+15, y3)   //garis vertical No KANAN
	// pdf.Line(pdf.GetX()+115, yline1, pdf.GetX()+115, y3) //garis itemspe Kanan
	// pdf.Line(pdf.GetX()+138, yline1, pdf.GetX()+138, y3) //qty kanan
	// pdf.Line(pdf.GetX()+158, yline1, pdf.GetX()+158, y3) // usd price kanan
	// pdf.Line(pdf.GetX()+188, yline1, pdf.GetX()+188, y3) // idr price kanan
	// pdf.Line(pdf.GetX()+208, yline1, pdf.GetX()+208, y3) // usd amount kanan
	// pdf.Line(pdf.GetX()+238, yline1, pdf.GetX()+238, y3) // idr price kanan
	// // pdf.Line(pdf.GetX()+280, yline1, pdf.GetX()+280, y3) // vertical remark kanan

	// pdf.Line(pdf.GetX(), y3, 248, y3) //garis horizontal kiri No

	// //Sub Total
	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "Total", "LBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(198)
	// total := tk.Sprintf("%.2f", DATA.Total)
	// if DATA.Currency == "IDR" {
	// 	total = ""
	// }
	// pdf.MultiCell(20, 5, total, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// total = tk.Sprintf("%.2f", DATA.Total)
	// if DATA.Currency == "USD" {
	// 	total = ""
	// }
	// pdf.MultiCell(30, 5, total, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //VAT
	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "VAT 10%", "LBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(198)
	// vat := tk.Sprintf("%.2f", DATA.VAT)
	// if DATA.Currency == "IDR" {
	// 	vat = ""
	// }
	// pdf.MultiCell(20, 5, vat, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// vat = tk.Sprintf("%.2f", DATA.VAT)
	// if DATA.Currency == "USD" {
	// 	vat = ""
	// }
	// pdf.MultiCell(30, 5, vat, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //Data GrandTotalidr
	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "Grand Total", "LBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 10, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(198)
	// grandtotal1 := tk.Sprintf("%.2f", DATA.GrandTotalUSD)
	// pdf.MultiCell(20, 10, grandtotal1, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// grandtotal := tk.Sprintf("%.2f", DATA.GrandTotalIDR)
	// pdf.MultiCell(30, 10, grandtotal, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //grandtotalusd
	// // y3 = pdf.GetY()
	// // pdf.SetY(y3)
	// // pdf.SetX(148)
	// // pdf.MultiCell(20, 5, "Grand Total USD", "LBR", "C", false)
	// // pdf.SetY(pdf.GetY())

	// // pdf.SetY(y3)
	// // pdf.SetX(168)
	// // pdf.MultiCell(30, 5, "", "BR", "C", false)
	// // pdf.SetY(pdf.GetY())

	// // pdf.SetY(y3)
	// // pdf.SetX(218)
	// // grandtotal1 := tk.Sprintf("%.2f", DATA.GrandTotalUSD)
	// // pdf.MultiCell(30, 5, grandtotal1, "BR", "R", false)
	// // pdf.SetY(pdf.GetY())

	// // //GranTotal
	// // y3 = pdf.GetY()
	// // pdf.SetY(y3)
	// // pdf.SetX(148)
	// // pdf.MultiCell(20, 5, "Total", "LBR", "C", false)
	// // pdf.SetY(pdf.GetY())

	// // pdf.SetY(y3)
	// // pdf.SetX(168)
	// // pdf.MultiCell(30, 5, "", "BR", "C", false)
	// // pdf.SetY(pdf.GetY())

	// // pdf.SetY(y3)
	// // pdf.SetX(198)
	// // grandTotal := tk.Sprintf("%.2f", DATA.GrandTotalUSD)
	// // if DATA.Currency == "IDR" {
	// // 	grandTotal = ""
	// // }
	// // pdf.MultiCell(20, 5, grandTotal, "BR", "R", false)
	// // pdf.SetY(pdf.GetY())

	// // pdf.SetY(y3)
	// // pdf.SetX(218)
	// // grandTotal = tk.Sprintf("%.2f", DATA.GrandTotalIDR)
	// // if DATA.Currency == "USD" {
	// // 	grandTotal = ""
	// // }
	// // pdf.MultiCell(30, 5, grandTotal, "BR", "R", false)
	// // pdf.SetY(pdf.GetY())

	// // ysig := y3
	// pdf.SetY(y3)
	// pdf.SetX(11)
	// pdf.SetFont("Arial", "", 11)
	// pdf.CellFormat(15, 12, "PT. Wiyasa Teknologi Nusantara", "", 0, "L", false, 0, "")
	// pdf.Ln(20)

	// pdf.SetY(y3)
	// pdf.SetX(89)
	// pdf.CellFormat(15, 12, "Approved By,", "", 0, "L", false, 0, "")
	// pdf.Ln(20)

	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(11)
	// pdf.SetFont("Arial", "", 11)
	// pdf.CellFormat(15, 12, "Andiek Suncahyo", "", 0, "L", false, 0, "")

	// pdf.SetY(y3)
	// pdf.SetX(89)
	// pdf.CellFormat(15, 12, "Arief Darmawan", "", 0, "L", false, 0, "")

	// // pdf.CellFormat(3, 3, "", "1", 0, "L", true, 0, "")
	// // pdf.CellFormat(30, 10, "Document Number ", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(80, 10, DATA.DocumentNumber, "", 0, "L", false, 0, "")
	// // pdf.CellFormat(20, 10, "Date", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(35, 10, DATA.DateStr, "", 0, "L", false, 0, "")
	// // pdf.Ln(5)

	// // pdf.GetY()
	// // pdf.SetX(12)
	// // pdf.CellFormat(30, 10, "Supplier Code ", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(80, 10, DATA.SupplierCode, "", 0, "L", false, 0, "")
	// // pdf.CellFormat(20, 10, "Payment", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(35, 10, DATA.Payment, "", 0, "L", false, 0, "")
	// // pdf.Ln(5)

	// // pdf.GetY()
	// // pdf.SetX(12)
	// // pdf.CellFormat(30, 10, "Supplier Name ", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(80, 10, DATA.SupplierName, "", 0, "L", false, 0, "")
	// // pdf.CellFormat(20, 10, "Type", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(35, 10, DATA.Type, "", 0, "L", false, 0, "")
	// // pdf.Ln(5)

	// // pdf.GetY()
	// // pdf.SetX(12)
	// // pdf.CellFormat(30, 10, "Account Code ", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // pdf.CellFormat(80, 10, DATA.AccountCode, "", 0, "L", false, 0, "")

	// // // vat := tk.Sprintf("%.2f", DATA.VAT)
	// // // pdf.CellFormat(20, 10, "VAT ", "", 0, "L", false, 0, "")
	// // // pdf.CellFormat(5, 10, ":", "", 0, "L", false, 0, "")
	// // // pdf.CellFormat(35, 10, vat, "", 0, "L", false, 0, "")
	// // pdf.Ln(10)

	// // y3 := pdf.GetY()
	// // widthHead := []float64{50.0, 15.0, 30.0, 30.0, 30.0, 30.0}
	// // for i, head := range PurchaseOrderHeadTable {
	// // 	pdf.SetY(y3)
	// // 	x := 13.0
	// // 	for y, z := range widthHead {
	// // 		if i > y {
	// // 			x += z
	// // 		} else {
	// // 			x += 0.0
	// // 		}
	// // 	}
	// // 	pdf.SetX(x)
	// // 	if i == 0 {
	// // 		pdf.MultiCell(widthHead[i], 7, head, "LRT", "C", false)
	// // 	} else {
	// // 		pdf.MultiCell(widthHead[i], 7, head, "RT", "C", false)
	// // 	}

	// // }

	// // y4 := pdf.GetY()

	// // for _, list := range DATA.ListDetail {
	// // 	y4 = pdf.GetY()
	// // 	pdf.SetY(y4)
	// // 	x := 13.0
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	item := list.Item
	// // 	pdf.MultiCell(widthHead[0], 4, item, "LT", "C", false)
	// // 	pdf.SetY(y4)

	// // 	x += widthHead[0]
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	qty := strconv.Itoa(list.Qty)
	// // 	pdf.MultiCell(widthHead[1], 4, qty, "T", "C", false)
	// // 	pdf.SetY(y4)

	// // 	x = 13.0 + widthHead[0] + widthHead[1]
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	priceusd := tk.Sprintf("%.2f", list.PriceUSD)
	// // 	pdf.MultiCell(widthHead[2], 4, priceusd, "T", "C", false)
	// // 	pdf.SetY(y4)

	// // 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2]
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	priceidr := tk.Sprintf("%.2f", list.PriceIDR)
	// // 	pdf.MultiCell(widthHead[3], 4, priceidr, "T", "C", false)
	// // 	pdf.SetY(y4)

	// // 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	amountusd := tk.Sprintf("%.2f", list.AmountUSD)
	// // 	pdf.MultiCell(widthHead[3], 4, amountusd, "T", "C", false)
	// // 	pdf.SetY(y4)

	// // 	x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// // 	pdf.SetLeftMargin(x)
	// // 	pdf.SetX(x)
	// // 	amountidr := tk.Sprintf("%.2f", list.AmountIDR)
	// // 	pdf.MultiCell(widthHead[3], 4, amountidr, "TR", "C", false)
	// // 	pdf.SetY(y4)
	// // }

	// // pdf.Ln(4)
	// // y4 = pdf.GetY()
	// // pdf.SetY(y4)
	// // x := 13.0
	// // pdf.SetLeftMargin(x)
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[0]+widthHead[1], 4, "Remark : ", "T", "L", false)
	// // pdf.SetY(y4)
	// // x = 33.0
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[5], 4, DATA.Remark, "", "R", false)

	// // pdf.SetY(y4)
	// // x = 13.0 + widthHead[0] + widthHead[1]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "Total :", "T", "R", false)
	// // pdf.SetY(y4)
	// // total := tk.Sprintf("%.2f", DATA.Total)
	// // x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[5], 4, total, "T", "C", false)
	// // pdf.SetY(y4)

	// // pdf.Ln(4)
	// // y4 = pdf.GetY()
	// // pdf.SetY(y4)
	// // x = 13.0 + widthHead[0] + widthHead[1]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "VAT :", "0", "R", false)
	// // pdf.SetY(y4)
	// // vat := tk.Sprintf("%.2f", DATA.VAT)
	// // x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[5], 4, vat, "", "C", false)
	// // pdf.SetY(y4)

	// // pdf.Ln(4)
	// // y4 = pdf.GetY()
	// // pdf.SetY(y4)
	// // x = 13.0 + widthHead[0] + widthHead[1]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "Diskon :", "0", "R", false)
	// // pdf.SetY(y4)
	// // diskon := tk.Sprintf("%.2f", DATA.Discount)
	// // x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[5], 4, diskon, "", "C", false)
	// // pdf.SetY(y4)
	// // // pdf.SetY(y4)

	// // pdf.Ln(4)
	// // y4 = pdf.GetY()
	// // pdf.SetY(y4)
	// // x = 13.0 + widthHead[0] + widthHead[1]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[2]+widthHead[3]+widthHead[4], 4, "GrandTotal :", "0", "R", false)
	// // pdf.SetY(y4)
	// // grandTotal := tk.Sprintf("%.2f", DATA.GrandTotal)
	// // x = 13.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
	// // pdf.SetX(x)
	// // pdf.MultiCell(widthHead[5], 4, grandTotal, "", "C", false)
	// // pdf.SetY(y4)

	namepdf := "-invoice.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	e = pdf.OutputFileAndClose(fileName)
	tk.Println(e)
	if e != nil {
		e.Error()
	}
	return ""
}
func ConvertToCurrency(FloatString string) string {
	firstStr := strings.Split(FloatString, ".")[0]
	lenFirstStr := len(firstStr) % 3
	lenComa := 0
	if lenFirstStr == 0 {
		lenComa = len(firstStr)/3 - 1

	} else {
		lenComa = len(firstStr) / 3

	}
	upsideDown := ""
	x := 1
	y := 1
	for i := len(firstStr) - 1; i >= 0; i-- {
		upsideDown = upsideDown + string(firstStr[i])

		if x == 3 && y <= lenComa {
			upsideDown = upsideDown + ","
			x = 1
			y++
		} else {

			x++
		}
	}
	res := ""
	for i := len(upsideDown) - 1; i >= 0; i-- {
		res = res + string(upsideDown[i])
	}
	results := res + "." + strings.Split(FloatString, ".")[1]
	return results
}
