package main

import (
	"time"

	"github.com/eaciit/knot/knot.v1"
)

type Dashboard struct {
}

func (d *Dashboard) Index(ctx *knot.WebContext) interface{} {
	loginid := ctx.Session("loginid", "").(string)
	if loginid == "" {
		ctx.Config.OutputType = knot.OutputJson
		ctx.Server.Log().Warning("No user")
		//http.Redirect(ctx.Writer, ctx.Request, "/login/index", 301)
		//return toolkit.NewResult()
	}
	ctx.Server.Log().Info("User: " + loginid)
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

type timedata struct {
	TimeStamp                             time.Time
	Power, Uptime, Downtime, Speed, Count float32
}
