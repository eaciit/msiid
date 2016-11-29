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

type DataRequest struct {
	Table string
	Take  int
}

var e error

func (r *Rest) Data(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	res := toolkit.NewResult()

	datareq := new(DataRequest)
	datareq.Table = ctx.QueryDef("table", "")
	datareq.Take = toolkit.ToInt(ctx.QueryDef("take", "10"), toolkit.RoundingAuto)

	c := conn()
	defer c.Close()

	cs, ecs := c.NewQuery().From(datareq.Table).Take(datareq.Take).Select().Cursor(nil)
	if ecs != nil {
		return res.SetErrorTxt("Error preparing query: " + ecs.Error())
	}
	defer cs.Close()

	data := []*toolkit.M{}
	cs.Fetch(&data, 0, false)

	return res.SetData(data)
}
