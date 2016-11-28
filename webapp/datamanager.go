package main

import (
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

var tablename string

type DataManager struct {
}

func (d *DataManager) Index(ctx *knot.WebContext) interface{} {
	tablename = ctx.QueryDef("table", "")
	if tablename == "" {
		
	}
	ctx.Config.OutputType = knot.OutputTemplate
	return struct{}{}
}

func (d *DataManager) Populate(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson
	result := toolkit.NewResult()

	pcode := func(ctx *knot.WebContext, pm *PopulateModel) *toolkit.Result {
		//-- read grid info
		skip := pm.Skip
		take := pm.Take
		sortfields := pm.Sortings

		//-- get count
		count := 5

		//-- calc offset
		nextoffset := skip + take
		if nextoffset > count {
			nextoffset = count
		}

		//-- get data
		models := []toolkit.M{}
		for i := 0; i < count; i++ {
			models = append(models,
				toolkit.M{}.
					Set("_id", toolkit.Sprintf("DS%d", i+1)).
					Set("title", toolkit.Sprintf("Datasource %d", i+1)))
		}

		//-- return the data
		result.SetData(toolkit.M{}.Set("data", models[skip:nextoffset]).Set("count", count).Set("skip", skip).
			Set("take", take).
			Set("sortfield", sortfields))

		return result
	}

	return populateGrid(ctx, pcode)
}

type PopulateModel struct {
	PageSize, Skip, Take, Count int
	Sortings                    []string
	Search                      toolkit.M
}

func populateGrid(ctx *knot.WebContext,
	populateFn func(*knot.WebContext, *PopulateModel) *toolkit.Result) *toolkit.Result {
	ctx.Config.OutputType = knot.OutputJson

	//-- get grid config
	pm := new(PopulateModel)
	pm.Skip = toolkit.ToInt(ctx.QueryDef("skip", ""), toolkit.RoundingAuto)
	pm.Take = toolkit.ToInt(ctx.QueryDef("take", ""), toolkit.RoundingAuto)
	sortdirection := ctx.Query("sort[0][dir]")
	if sortdirection == "asc" {
		pm.Sortings = []string{ctx.Query("sort[0][field]")}
	} else {
		pm.Sortings = []string{"-", ctx.Query("sort[0][field]")}
	}

	return populateFn(ctx, pm)
}
