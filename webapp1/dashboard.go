package main

import (
	"github.com/eaciit/knot/knot.v1"
)

type Dashboard struct {
}

func (d *Dashboard) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}