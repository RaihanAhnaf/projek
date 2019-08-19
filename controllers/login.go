package controllers

import (
	. "../helpers"
	. "../models"
	"fmt"
	"net/http"
	"time"

	db "github.com/eaciit/dbox"
	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	gomail "gopkg.in/gomail.v2"
)

type LoginController struct {
	*BaseController
}

func (c *LoginController) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	k.Config.IncludeFiles = []string{
		"_loader.html",
		"_head.html",
		"_loader2.html",
	}
	if k.Session("userid") == nil {
		k.Session("userid", nil)
	}

	return ""
}

func (c *LoginController) Do(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	formdata := struct {
		Username string
		Password string
	}{}

	message := ""

	err := k.GetPayload(&formdata)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	ctx := c.Ctx.Connection

	filter := []*db.Filter{}
	filter = append(filter, db.Eq("username", formdata.Username))

	csr, err := ctx.NewQuery().Select().From("SysUsers").Where(filter...).Cursor(nil)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	res := make([]SysUserModel, 0)
	resroles := make([]SysRolesModel, 0)
	resurl := []tk.M{}

	defer csr.Close()

	err = csr.Fetch(&res, 0, false)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	if len(res) > 0 {
		resUser := res[0]
		if GetMD5Hash(formdata.Password) == resUser.Password {
			if resUser.Enable == true {
				fil := []*db.Filter{}
				fil = append(fil, db.Eq("name", resUser.Roles))
				cursor, err := ctx.NewQuery().Select().From("SysRoles").Where(fil...).Cursor(nil)
				if err != nil {
					return c.SetResultInfo(true, err.Error(), nil)
				}

				err = cursor.Fetch(&resroles, 0, false)
				if err != nil {
					// tk.Println(err.Error())
					return c.SetResultInfo(true, err.Error(), nil)
				}

				defer cursor.Close()

				k.SetSession("userid", resUser.Id.Hex())
				k.SetSession("username", resUser.Username)
				k.SetSession("fullname", resUser.Fullname)
				k.SetSession("usermodel", resUser)
				k.SetSession("roles", resroles)
				k.SetSession("rolesid", resroles[0].Id.Hex())
				k.SetSession("stime", time.Now())
				k.SetSession("locationid", resUser.LocationID)
				k.SetSession("locationname", resUser.LocationName)

				cursor, err = ctx.NewQuery().Select().From("SysMenus").Where(db.Eq("Title", resroles[0].Landing)).Cursor(nil)
				if err != nil {
					// tk.Println(err.Error())
					return c.SetResultInfo(true, err.Error(), nil)
				}

				err = cursor.Fetch(&resurl, 0, false)

				defer cursor.Close()
				c.LogActivity("Login", "Log In", resUser.Username, k)

			} else {
				message = "Your account is disabled, please contact administrator to enable it."
			}
		} else {
			message = "Invalid Username or password!"
		}
	} else {
		message = "Invalid Username or password!"
		return c.SetResultInfo(true, message, nil)
	}
	return c.SetResultInfo(false, message, tk.M{}.Set("Roles", resurl))
}

func (c *LoginController) SessionCheckTimeOut(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	lastActive := k.Session("stime").(time.Time)
	duration := time.Since(lastActive)

	if duration.Minutes() > 14 {
		ret.IsError = true
		ret.Data = "/login/default"
	}
	return ret
}

func (c *LoginController) HeartBeat(k *knot.WebContext) interface{} {
	c.IsAuthenticate(k)
	k.Config.OutputType = knot.OutputJson
	ret := ResultInfo{}
	k.SetSession("stime", time.Now())
	return ret
}

func (c *LoginController) CP(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = "login/cp.html"
	k.Config.IncludeFiles = []string{
		"_loader2.html",
		"_head.html",
		"_loader.html",
	}
	token := k.Request.FormValue("tk")
	// fmt.Printf("token %v \n", token)
	ctx := c.Ctx.Connection
	checkToken := []tk.M{}
	csr, e := ctx.NewQuery().Where(db.Eq("tokenid", token)).From("SysUsers").Cursor(nil)

	e = csr.Fetch(&checkToken, 0, false)
	defer csr.Close()
	if e != nil {
		return c.SetResultInfo(true, e.Error(), nil)
	}

	if len(checkToken) == 0 {
		http.Redirect(k.Writer, k.Request, "/login/default", 302)
	}
	return ""
}

func (c *LoginController) ForgetPassword(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	// fmt.Println(payload.GetString("UserEmail"))
	useremail := payload.GetString("UserEmail")
	UrlString := payload.GetString("UrlString")
	ctx := c.Ctx.Connection

	checkEmail := []tk.M{}
	csr, errore := ctx.NewQuery().
		Where(db.Eq("email", useremail)).
		From("SysUsers").Cursor(nil)

	errore = csr.Fetch(&checkEmail, 0, false)
	defer csr.Close()

	if errore != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	if len(checkEmail) > 0 {
		token := tk.GenerateRandomString("", 30)
		dataresult := tk.M{
			"tokenid": token,
		}

		err = ctx.NewQuery().Update().From("SysUsers").Where(db.Eq("email", useremail)).Exec(tk.M{"data": dataresult})

		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}

		linktoken := fmt.Sprintf("%v/login/cp?tk=%v", UrlString, token)

		mailmsg := fmt.Sprintf("Hi, <br/><br/> We received a request to reset your password, <br/><br/>")
		mailmsg = fmt.Sprintf("%vFollow the link below to set a new password : <br/><br/> %v <br/><br/>", mailmsg, linktoken)
		mailmsg = fmt.Sprintf("%vIf you don't want to change your password, you can ignore this email <br/><br/> Thanks,</body></html>", mailmsg)

		m := gomail.NewMessage()

		m.SetHeader("From", "admin.support@eaciit.com")
		m.SetHeader("To", useremail)

		m.SetHeader("Subject", "[no-reply] Reset Password")
		m.SetBody("text/html", mailmsg)

		smtpPass := DecryptAes128("U2FsdGVkX19K+uD+i9wB6qtNB9cVN2zE7lX0LWRCPeYlASA/jlvwEwxqWxKRiKiB", AES128KEY)

		// d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
		d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", smtpPass)
		err = d.DialAndSend(m)

		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}
		// c.LogActivity("Login", "Forget Password", useremail, k)
	}
	return c.SetResultInfo(false, "Success", checkEmail)
}

func (c *LoginController) SaveNewPass(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := tk.M{}
	err := k.GetPayload(&payload)
	if err != nil {
		return c.SetResultInfo(true, err.Error(), nil)
	}

	useremail := payload.GetString("Email")
	newpass := payload.GetString("NewPass")
	urlToLogin := payload.GetString("UrlToLogin")

	ctx := c.Ctx.Connection

	checkEmail := []tk.M{}
	csr, errore := ctx.NewQuery().
		Where(db.Eq("email", useremail)).
		From("SysUsers").Cursor(nil)

	errore = csr.Fetch(&checkEmail, 0, false)
	defer csr.Close()

	if errore != nil {
		return c.SetResultInfo(true, errore.Error(), nil)
	}

	if len(checkEmail) > 0 {
		dataresult := tk.M{
			"tokenid":  "",
			"password": newpass,
		}

		err = ctx.NewQuery().Update().From("SysUsers").Where(db.Eq("email", useremail)).Exec(tk.M{"data": dataresult})

		if err != nil {
			return c.SetResultInfo(true, err.Error(), nil)
		}

		// send email

		mailmsg := fmt.Sprintf("")
		mailmsg = fmt.Sprintf("%v<br/>Congratulations, you successfully change Accounting account Password", mailmsg)
		mailmsg = fmt.Sprintf("%vEmail :  %v", mailmsg, useremail)
		// mailmsg = fmt.Sprintf("%v<br/><br/>Please Login to http://accgo.eaciit.com", mailmsg)
		mailmsg = fmt.Sprintf("%v<br/><br/>Please Login to : %v", mailmsg, urlToLogin)
		mailmsg = fmt.Sprintf("%v<br/><br/>2016 - 2020 EACIIT Pte Ltd<br/>Earnings & Cash Improvement Information Technologies", mailmsg)

		m := gomail.NewMessage()

		m.SetHeader("From", "admin.support@eaciit.com")
		m.SetHeader("To", useremail)

		m.SetHeader("Subject", "[no-reply] Account Accounting Information")
		m.SetBody("text/html", mailmsg)

		smtpPass := DecryptAes128("U2FsdGVkX19K+uD+i9wB6qtNB9cVN2zE7lX0LWRCPeYlASA/jlvwEwxqWxKRiKiB", AES128KEY)

		d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", smtpPass)
		// d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
		err = d.DialAndSend(m)
		//end send email
		return c.SetResultInfo(false, "Success", dataresult)

		c.LogActivity("Login", "Save New Password", useremail, k)

	}

	return c.SetResultInfo(true, "reset password error", nil)

}
