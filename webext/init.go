package webext

import (
	. "../controllers"
	. "../helpers"
	"os"

	_ "github.com/eaciit/dbox/dbc/mongo"
	knot "github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func init() {

	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	cfg := ReadConfig()
	ctx := orm.New(conn)
	baseCtrl := new(BaseController)
	baseCtrl.Ctx = ctx
	baseCtrl.LogoFile = cfg["imgpath"]
	baseCtrl.FontPath = cfg["fonthpath"]
	baseCtrl.UploadPath = cfg["uploadpath"]
	baseCtrl.PdfPath = cfg["pdfpath"]
	app := knot.NewApp("ecfinancial")
	app.ViewsPath = wd + "views/"

	listPath := make([]string, 0)
	listPath = append(listPath, cfg["uploadpath"])
	listPath = append(listPath, cfg["pdfpath"])
	listPath = append(listPath, cfg["txtpath"])
	listPath = append(listPath, cfg["imgpath"])
	listPath = append(listPath, cfg["htmlpath"])
	listPath = append(listPath, cfg["jsonpath"])
	mkdir(listPath)
	/**REGISTER ALL CONTROLLERS HERE**/
	app.Register(&LoginController{baseCtrl})
	app.Register(&DashboardController{baseCtrl})
	app.Register(&LogoutController{baseCtrl})
	app.Register(&MenuSettingController{baseCtrl})
	app.Register(&UserSettingController{baseCtrl})
	app.Register(&SysRolesController{baseCtrl})
	app.Register(&DatamasterController{baseCtrl})
	app.Register(&PoInvoiceSummaryController{baseCtrl})
	// app.Register(&RefMasterController{baseCtrl})
	// app.Register(&RefVendorsController{baseCtrl})
	// app.Register(&CountryMasterController{baseCtrl})
	// app.Register(&RefDocumentClauseController{baseCtrl})
	app.Register(&BalanceSheetSettingController{baseCtrl})
	app.Register(&BackOfficeController{baseCtrl})
	app.Register(&SearchController{baseCtrl})
	app.Register(&ActivityLogController{baseCtrl})
	app.Register(&MasterController{baseCtrl})
	app.Register(&FinancialController{baseCtrl})
	app.Register(&ReportController{baseCtrl})
	app.Register(&TransactionController{baseCtrl})
	app.Register(&UserProfileController{baseCtrl})
	// app.Register(&InvoiceController{baseCtrl})
	app.Register(&ClosingController{baseCtrl})
	app.Register(&TransferOrderController{baseCtrl})

	app.Static("res", wd+"assets")
	app.Static("files", wd+"upload")
	app.LayoutTemplate = "_layout.html"
	knot.RegisterApp(app)
	tk.Println("___INIT FINISH_____")
}
func mkdir(paths []string) {
	for _, path := range paths {
		if _, errs := os.Stat(path); errs != nil {
			os.Mkdir(path, os.ModePerm)
		}

	}
}

// func readConfig() map[string]string {
// 	ret := make(map[string]string)
// 	file, err := os.Open(wd + "conf/app.conf")
// 	if err == nil {
// 		defer file.Close()

// 		reader := bufio.NewReader(file)
// 		for {
// 			line, _, e := reader.ReadLine()
// 			if e != nil {
// 				break
// 			}

// 			sval := strings.Split(string(line), "=")
// 			ret[sval[0]] = sval[1]
// 		}
// 	} else {
// 		tk.Println(err.Error())
// 	}

// 	return ret
// }
