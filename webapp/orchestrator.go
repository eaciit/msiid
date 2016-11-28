package main

import (
	"github.com/eaciit/knot/knot.v1"
)

type Orchestrator struct {
}

func (d *Orchestrator) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}
