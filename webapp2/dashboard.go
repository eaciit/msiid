package main

import (
	//"eaciit/mfg"
	"time"

	//"github.com/eaciit/crowd.dev"

	"github.com/eaciit/knot/knot.v1"
)

type Dashboard struct {
}

func (d *Dashboard) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

type timedata struct {
	TimeStamp                             time.Time
	Power, Uptime, Downtime, Speed, Count float32
}
