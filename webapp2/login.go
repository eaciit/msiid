package main

import "github.com/eaciit/knot/knot.v1"
import "github.com/eaciit/toolkit"
import "strings"

type Login struct {
}

func (l *Login) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

type AuthRequest struct {
	UserId, Password string
}

func (l *Login) Auth(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	res := toolkit.NewResult()

	a := new(AuthRequest)
	a.UserId = ctx.FormDef("UserId", "")
	a.Password = ctx.FormDef("Password", "")

	if strings.ToLower(a.UserId) != "msi" && a.Password != "M$iEaciit" {
		return res.SetErrorTxt("Invalid credential is provided")
	}

	ctx.Server.Log().Info("User " + a.UserId + " Login")
	ctx.SetSession("loginid", a.UserId)
	return res
}

func (l *Login) Logout(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	res := toolkit.NewResult()

	ctx.SetSession("loginid", "")
	return res
}
