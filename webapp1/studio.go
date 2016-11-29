package main

import (
	"github.com/eaciit/knot/knot.v1"
)


type Studio struct {
}

func (d *Studio) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}
