package main

import (
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type Rest struct {
}

func (r *Rest) Tables(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	c := conn()
	defer c.Close()

	cs, _ := c.NewQuery().From("metadata").Order("_id").Cursor(nil)
	defer cs.Close()

	data := []*toolkit.M{}
	cs.Fetch(&data, 0, false)

	/*
		    names := []string{}
			for _, d := range data {
				names = append(names, d.GetString("_id"))
			}
	*/

	return toolkit.NewResult().SetData(data)
}
