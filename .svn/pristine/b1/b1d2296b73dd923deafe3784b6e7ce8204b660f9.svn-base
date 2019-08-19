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
	SalesPaymentPDF(dateStart, dateEnd)
}

func SalesPaymentPDF(dateStart time.Time, dateEnd time.Time) interface{} {

	// Digunakan
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}

	Ctx := orm.New(conn)
	csr, e := Ctx.Connection.NewQuery().Select().From("SalesPayment").Where(db.Eq("_id", bson.ObjectIdHex("59ba2b1e53cfe31bc8e81f70"))).Cursor(nil)
	if e != nil {
		tk.Println("query", e.Error())
	}
	defer csr.Close()
	results := []SalesPaymentModel{}
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		tk.Println("fetch", e.Error())
	}
	DATA := results[0]
	tk.Println(DATA)

	config := helpers.ReadConfig()
	Img := config["imgpath"]

	pdf := gofpdf.New("L", "mm", "A4", "")
	SalesPaymentHeadTable := []string{"No", "Date", "Invoice Number", "Amount", "Already Paid", "Receive", "Balance", "Pay"}
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
	pdf.CellFormat(0, 15, "SALES PAYMENT", "", 0, "L", false, 0, "")

	pdf.SetY(y1 + 17)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, "Surabaya-Indonesia", "", 0, "L", false, 0, "")
	pdf.SetX(10)
	pdf.Ln(1)

	pdf.SetY(y1 + 23)
	pdf.SetX(11)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 12, "Telp. 0351-5676223", "", 0, "L", false, 0, "")
	pdf.SetX(10)

	pdf.Line(pdf.GetX(), pdf.GetY()+9, 286, pdf.GetY()+9)   //garis horizontal1
	pdf.Line(pdf.GetX(), pdf.GetY()+25, 286, pdf.GetY()+25) //garis horizontal2

	pdf.Line(pdf.GetX(), pdf.GetY()+9, pdf.GetX(), 80)         // garis vertikal 1 dari kiri/156 -> 136
	pdf.Line(pdf.GetX()+276, pdf.GetY()+9, pdf.GetX()+276, 75) // garis vertikal 3 paling kanan
	pdf.Line(pdf.GetX()+138, pdf.GetY()+9, pdf.GetX()+138, 55) // garis vertikal 2 tengah
	// pdf.Line(pdf.GetX()+69, pdf.GetY()+30, pdf.GetX()+69, 91)   // garis vertikal 4 tenga2nya field to//101

	y2 := pdf.GetY()
	pdf.SetX(9)
	pdf.SetY(y2 + 10)
	pdf.MultiCell(12, 5, " To : ", "", "L", false)

	// pdf.SetX(12)
	pdf.SetY(y2 + 10)
	pdf.SetX(19)
	pdf.MultiCell(130, 5, DATA.CustomerName, "", "L", false)

	pdf.SetY(y2 + 26)
	pdf.SetX(10)
	pdf.MultiCell(138.5, 5, " Payment :", "", "L", false)

	filPtcMandiri := false
	filPettyCash := false
	filUsdMandiri := false
	filIdrMandiri := false
	if DATA.PaymentName == "PTC - MANDIRI" {
		filPtcMandiri = true
		filPettyCash = false
		filUsdMandiri = false
		filIdrMandiri = false
	} else if DATA.PaymentName == "PETTYCASH" {
		filPtcMandiri = false
		filPettyCash = true
		filUsdMandiri = false
		filIdrMandiri = false
	} else if DATA.PaymentName == "USD MANDIRI" {
		filPtcMandiri = false
		filPettyCash = false
		filUsdMandiri = true
		filIdrMandiri = false
	} else {
		filPtcMandiri = false
		filPettyCash = false
		filUsdMandiri = false
		filIdrMandiri = true
	}

	pdf.SetY(y2 + 31)
	pdf.SetX(13)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filPtcMandiri, 0, "")
	pdf.SetY(y2 + 31)
	pdf.SetX(20)
	pdf.MultiCell(138.5, 5, "PTC - MANDIRI", "", "L", false)
	// pdf.GetY()

	pdf.SetY(y2 + 31)
	pdf.SetX(60)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filPettyCash, 0, "")
	pdf.SetY(y2 + 31)
	pdf.SetX(67)
	pdf.MultiCell(138.5, 5, "PETTY CASH", "", "L", false)
	// pdf.GetY()

	pdf.SetY(y2 + 31)
	pdf.SetX(107)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filUsdMandiri, 0, "")
	pdf.SetY(y2 + 31)
	pdf.SetX(114)
	pdf.MultiCell(138.5, 5, "USD MANDIRI", "", "L", false)
	// pdf.GetY()

	pdf.SetY(y2 + 31)
	pdf.SetX(154)
	pdf.CellFormat(4.5, 5, "", "1", 0, "L", filIdrMandiri, 0, "")
	pdf.SetY(y2 + 31)
	pdf.SetX(161)
	pdf.MultiCell(138.5, 5, "IDR MANDIRI", "", "L", false)
	// pdf.GetY()

	// //TABLE KANAN
	pdf.SetY(y2 + 10)
	pdf.SetX(147.5)
	pdf.MultiCell(27, 5, " Payment No", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(171.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 10)
	pdf.SetX(173.5)
	pdf.MultiCell(90, 5, DATA.DocumentNumber, "", "L", false)
	pdf.SetY(pdf.GetY())
	// pdf.GetY()

	pdf.SetY(y2 + 15)
	pdf.SetX(147.5)
	pdf.MultiCell(138.5, 5, " Date", "", "L", false)
	pdf.SetY(pdf.GetY())
	pdf.SetY(y2 + 15)
	pdf.SetX(171.5)
	pdf.MultiCell(4, 5, ":", "", "L", false)
	pdf.SetY(y2 + 15)
	pdf.SetX(173.5)
	pdf.MultiCell(90, 5, DATA.DateStr, "", "L", false)
	pdf.SetY(pdf.GetY())
	// pdf.GetY()

	y0 := pdf.GetY()
	widthHead := []float64{10.0, 28.0, 34.0, 45.0, 45.0, 45.0, 45.0, 24.0}
	for i, head := range SalesPaymentHeadTable {
		pdf.SetY(y0 + 23)
		x := 10.0
		for y, z := range widthHead {
			if i > y {
				x += z
			} else {
				x += 0.0
			}
		}
		pdf.SetX(x)
		if i == 0 {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		} else {
			pdf.MultiCell(widthHead[i], 7, head, "BRT", "C", false)
		}

	}
	// pdf.Ln()
	y3 := pdf.GetY()
	for i, list := range DATA.ListDetail {

		y3 = pdf.GetY()
		pdf.SetY(y3)
		x := 10.0
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		number := i + 1
		numberstr := strconv.Itoa(number)
		pdf.MultiCell(widthHead[0], 5, numberstr, "L", "C", false)
		// pdf.MultiCell(w, h, txtStr, borderStr, alignStr, fill)
		a0 := pdf.GetY()
		pdf.SetY(y3)
		x += widthHead[0]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[1], 5, list.DatePayment.Local().Format("02 Jan 2006"), "L", "L", false)

		a1 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		pdf.MultiCell(widthHead[2], 5, list.InvNumber, "L", "L", false)

		a2 := pdf.GetY()
		pdf.SetY(y3)
		// tk.Println(list.Acc_Code)
		x = 10 + widthHead[0] + widthHead[1] + widthHead[2]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		amount := tk.Sprintf("%.2f", list.Amount)
		pdf.MultiCell(widthHead[3], 5, amount, "L", "R", false)

		a3 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		alreadypaid := tk.Sprintf("%.2f", list.AlreadyPaid)
		pdf.MultiCell(widthHead[4], 5, alreadypaid, "L", "L", false)

		a4 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		receive := tk.Sprintf("%.2f", list.Receive)
		pdf.MultiCell(widthHead[5], 5, receive, "L", "L", false)

		a5 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		balance := tk.Sprintf("%.2f", list.Balance)
		pdf.MultiCell(widthHead[6], 5, balance, "L", "R", false)

		a6 := pdf.GetY()
		pdf.SetY(y3)
		x = 10.0 + widthHead[0] + widthHead[1] + widthHead[2] + widthHead[3] + widthHead[4] + widthHead[5] + widthHead[6]
		pdf.SetLeftMargin(x)
		pdf.SetX(x)
		if list.Pay {
			tk.Sprintf("%.2f", list.Pay)
			pdf.MultiCell(widthHead[7], 5, "Paid", "LR", "L", false)
		}

		a7 := pdf.GetY()
		pdf.SetY(y3)

		allA := []float64{a0, a1, a2, a3, a4, a5, a6, a7}
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
	// x := 10.0
	y4 := pdf.GetY()
	// pdf.Line(pdf.GetX()-236, y3, pdf.GetX()-236, y4) //garis vertical No KIRI
	// pdf.Line(pdf.GetX()-226, y3, pdf.GetX()-226, y4) //garis vertical No KANAN
	// pdf.Line(pdf.GetX()-192, y3, pdf.GetX()-192, y4) //garis itemspe Kanan
	// pdf.Line(pdf.GetX()-152, y3, pdf.GetX()-152, y4) //qty kanan
	// pdf.Line(pdf.GetX()-124, y3, pdf.GetX()-124, y4) // usd price kanan
	// pdf.Line(pdf.GetX()-80, y3, pdf.GetX()-80, y4)   // idr price kanan
	// pdf.Line(pdf.GetX()-40, y3, pdf.GetX()-40, y4)   // usd amount kanan
	// pdf.Line(pdf.GetX(), y3, pdf.GetX(), y4)         // idr price kanan
	// pdf.Line(pdf.GetX()+40, y3, pdf.GetX()+40, y4)   // vertical remark kanan

	// pdf.Line(pdf.GetX(), y4, 100, y4) //garis horizontal kiri No
	pdf.Line(pdf.GetX()-252, y4, 286, y4) //garis horizontal2

	// //Sub Total
	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "Sub Total", "LBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(198)
	// pdf.MultiCell(20, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// total := tk.Sprintf("%.2f", DATA.Total)
	// pdf.MultiCell(30, 5, total, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //Diskon
	// y3 = pdf.GetY()
	// pdf.SetY(y3)
	// pdf.SetX(148)
	// pdf.MultiCell(20, 5, "Disc.", "LBR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(168)
	// pdf.MultiCell(30, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(198)
	// pdf.MultiCell(20, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// diskon := tk.Sprintf("%.2f", DATA.Discount)
	// pdf.MultiCell(30, 5, diskon, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //Data VAT
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
	// pdf.MultiCell(20, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// vat := tk.Sprintf("%.2f", DATA.VAT)
	// pdf.MultiCell(30, 5, vat, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

	// //GranTotal
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
	// pdf.MultiCell(20, 5, "", "BR", "C", false)
	// pdf.SetY(pdf.GetY())

	// pdf.SetY(y3)
	// pdf.SetX(218)
	// grandTotal := tk.Sprintf("%.2f", DATA.GrandTotal)
	// pdf.MultiCell(30, 5, grandTotal, "BR", "R", false)
	// pdf.SetY(pdf.GetY())

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
	// pdf.CellFormat(15, 12, "Andik Suncahyo", "", 0, "L", false, 0, "")

	// pdf.SetY(y3)
	// pdf.SetX(89)
	// pdf.CellFormat(15, 12, "Arif Darmawan", "", 0, "L", false, 0, "")

	namepdf := "-salespayment.pdf"
	FixName := time.Now().Format("2006-01-02T150405") + namepdf
	fileName := FixName
	e = pdf.OutputFileAndClose(fileName)
	tk.Println(e)
	if e != nil {
		e.Error()
	}
	return ""
}
