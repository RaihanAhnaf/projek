package controllers

import (
	"../helpers"
	. "../models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type ResultInfo struct {
	IsError bool
	Message string
	Data    interface{}
}
type IBaseController interface {
}
type BaseController struct {
	base       IBaseController
	Ctx        *orm.DataContext
	UploadPath string
	PdfPath    string
	FontPath   string
	LogoFile   string
	DocPath    string
	DbHost     string
	DbUsername string
	DbPassword string
	DbName     string
	BasePath   string
	AppMode    string
}

type Previlege struct {
	View         bool
	Create       bool
	Edit         bool
	Delete       bool
	Approve      bool
	Process      bool
	Menuid       string
	Menuname     string
	Username     string
	TopMenu      string
	Pdfname      string
	Rolename     interface{}
	Fullname     interface{}
	LocationId   int
	LocationName string
}

func (b *BaseController) IsAuthenticate(k *knot.WebContext) {
	if k.Session("userid") == nil {
		b.Redirect(k, "login", "default")
	}
	return
}

func (b *BaseController) LoadBase(k *knot.WebContext) []tk.M {
	k.Config.NoLog = true
	if k.Session("userid") == nil {
		http.Redirect(k.Writer, k.Request, "/login/default", http.StatusTemporaryRedirect)
		return nil
	}

	access := b.AccessMenu(k)
	return access
}

func (b *BaseController) Redirect(k *knot.WebContext, controller string, action string) {
	// tk.Println("Eaciit - Accounting -->> redirecting to " + controller + "/" + action)
	http.Redirect(k.Writer, k.Request, "/"+controller+"/"+action, http.StatusTemporaryRedirect)
}

func (b *BaseController) AccessMenu(k *knot.WebContext) []tk.M {
	url := k.Request.URL.String()
	if strings.Index(url, "?") > -1 {
		url = url[:strings.Index(url, "?")]
	}
	sessionRoles := k.Session("roles")
	access := []tk.M{}
	if sessionRoles != nil {
		accesMenu := sessionRoles.([]SysRolesModel)
		if len(accesMenu) > 0 {
			for _, o := range accesMenu[0].Menu {
				if o.Url == url {
					obj := tk.M{}
					obj.Set("View", o.View)
					obj.Set("Create", o.Create)
					obj.Set("Approve", o.Approve)
					obj.Set("Delete", o.Delete)
					obj.Set("Process", o.Process)
					obj.Set("Edit", o.Edit)
					obj.Set("Menuid", o.Menuid)
					obj.Set("Menuname", o.Menuname)
					obj.Set("Username", k.Session("username").(string))
					obj.Set("Rolename", accesMenu[0].Name)
					obj.Set("Fullname", k.Session("fullname").(string))
					obj.Set("LocationId", k.Session("locationid").(int))
					obj.Set("LocationName", k.Session("locationname").(string))
					access = append(access, obj)
					return access
				}

			}
		}
	}
	return access
}

func (b *BaseController) SetResultInfo(isError bool, msg string, data interface{}) ResultInfo {
	r := ResultInfo{}
	r.IsError = isError
	r.Message = msg
	r.Data = data
	return r
}
func (b *BaseController) SetResultFile(isError bool, msg string, file string) ResultInfo {
	r := ResultInfo{}
	r.IsError = isError
	r.Message = msg
	r.Data = file
	return r
}

func (b *BaseController) ErrorResultInfo(msg string, data interface{}) ResultInfo {
	r := ResultInfo{}
	r.IsError = true
	r.Message = msg
	r.Data = data
	return r
}
func (b *BaseController) ErrorMessageOnly(msg string) ResultInfo {
	r := ResultInfo{}
	r.IsError = true
	r.Message = msg
	r.Data = nil
	return r
}
func (b *BaseController) SetResultOK(data interface{}) *tk.Result {
	r := tk.NewResult()
	r.Data = data

	return r
}
func (b *BaseController) GetTopMenuName(menuname string) string {
	res := []tk.M{}

	ctx := b.Ctx.Connection

	csr, err := ctx.NewQuery().Select().From("Sysmenus").Where(db.Eq("Title", menuname)).Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	err = csr.Fetch(&res, 0, false)
	defer csr.Close()
	parentId := ""

	for _, val := range res {
		parentId = val.GetString("Parent")
	}

	csr, err = ctx.NewQuery().Select().From("Sysmenus").Where(db.Eq("_id", parentId)).Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	res2 := []tk.M{}
	err = csr.Fetch(&res2, 0, false)
	defer csr.Close()
	Title := ""
	for _, val := range res2 {
		Title = val.GetString("Title")
	}

	return Title
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}

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

func (b *BaseController) GetDataAccess(access []tk.M) (DataAccess Previlege) {
	for _, o := range access {
		DataAccess.Create = o["Create"].(bool)
		DataAccess.View = o["View"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Process = o["Process"].(bool)
		DataAccess.Delete = o["Delete"].(bool)
		DataAccess.Edit = o["Edit"].(bool)
		DataAccess.Menuid = o["Menuid"].(string)
		DataAccess.Menuname = o["Menuname"].(string)
		DataAccess.Approve = o["Approve"].(bool)
		DataAccess.Username = o["Username"].(string)
		DataAccess.Rolename = o["Rolename"].(string)
		DataAccess.Fullname = o["Fullname"].(string)
		DataAccess.LocationId = o["LocationId"].(int)
		DataAccess.LocationName = o["LocationName"].(string)
	}

	DataAccess.TopMenu = b.GetTopMenuName(DataAccess.Menuname)

	return
}

// Fungsi ini tidak bisa concurent. Mohon diperbaiki jika usernya sudah banyak
func (b *BaseController) GetNextIdSeq(collName string, typeJournal string, month int, year int) (int, error) {
	mdl := []SequenceModel{}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"collname": collName, "typejournal": typeJournal, "month": month, "year": year}})
	csr, err := b.Ctx.Connection.NewQuery().Command("pipe", pipes).From("Sequence").Cursor(nil)
	// tk.Println("================")
	_, _ = json.Marshal(pipes)
	// tk.Println(string(i))
	// tk.Println(month)
	// tk.Println(year)
	if err != nil {
		return -9999, err
	}
	defer csr.Close()

	e := csr.Fetch(&mdl, 0, true)
	if e != nil {
		return 0, err
	}

	if len(mdl) == 0 {
		// tk.Println("========Kosong========")
		seq := NewSequenceModel()
		seq.Collname = collName
		seq.Typejournal = typeJournal
		seq.Lastnumber = 1
		seq.Month = month
		seq.Year = year
		e = b.Ctx.Save(seq)
		if e != nil {
			return 0, err
		}

		return 1, nil
	}

	seq := mdl[0]
	seq.Lastnumber = seq.Lastnumber + 1
	b.Ctx.Save(&seq)
	return seq.Lastnumber, nil
}

func (b *BaseController) GetNextIdSeqCustSupp(collName string, namecustomer string) (int, error) {
	mdl := []SequenceCustomerModel{}
	var pipes []tk.M
	pipes = append(pipes, tk.M{"$match": tk.M{"collname": collName, "namecust": namecustomer}})
	csr, err := b.Ctx.Connection.NewQuery().Command("pipe", pipes).From("SequenceCustomer").Cursor(nil)

	_, _ = json.Marshal(pipes)

	if err != nil {
		return -9999, err
	}
	defer csr.Close()

	e := csr.Fetch(&mdl, 0, true)
	if e != nil {
		return 0, err
	}

	if len(mdl) == 0 {
		// tk.Println("========Kosong========")
		seq := NewSequenceCustomerModel()
		seq.Collname = collName
		seq.NameCust = namecustomer
		seq.Lastnumber = 1
		e = b.Ctx.Save(seq)
		if e != nil {
			return 0, err
		}

		return 1, nil
	}

	seq := mdl[0]
	seq.Lastnumber = seq.Lastnumber + 1
	b.Ctx.Save(&seq)
	return seq.Lastnumber, nil
}
func (b *BaseController) LogActivity(pageName string, pageActivity string, desc string, k *knot.WebContext) {
	username := k.Session("username").(string)
	sUrl := k.Request.RequestURI
	// ipaddress := k.Request.RemoteAddr
	s := strings.Split(k.Request.RemoteAddr, ":")
	siP := ""
	for i := 0; i < len(s)-1; i++ {
		if i == 0 {
			siP = s[i]
			continue
		}
		siP = siP + ":" + s[i]
	}
	// tk.Println("ip address", siP)
	mdl := NewActivityLogModel()
	mdl.AccessTime = time.Now()
	mdl.AccessDate, _ = strconv.Atoi(time.Now().Format("20060102"))
	mdl.UserName = username
	mdl.IpAddress = siP
	mdl.PageName = pageName
	mdl.PageUrl = sUrl
	mdl.Activity = pageActivity
	mdl.Desc = desc
	b.Ctx.Save(mdl)
	return
}

// func (b *BaseController) InsertActivityLog(pageName string, pageActivity string, k *knot.WebContext) {
// 	username := k.Session("username").(string)
// 	sUrl := k.Request.URL.String()
// 	tk.Println("URL-BASE: ", sUrl)
// 	sUrl = sUrl[:strings.Index(sUrl, "/")]
// 	ipaddress := k.Request.RemoteAddr
// 	mdl := NewActivityLogModel()
// 	mdl.Username = username
// 	mdl.IpAddress = ipaddress[:strings.Index(ipaddress, ":")]
// 	mdl.PageName = pageName
// 	mdl.PageUrl = sUrl
// 	mdl.Activity = pageActivity
// 	mdl.AccessTime = time.Now()
// 	mdl.AccessDate, _ = strconv.Atoi(time.Now().Format("20060102"))
// 	b.Ctx.Save(mdl)
// 	return
// }
func (b *BaseController) ConvertToCurrency(FloatString string) string {
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
