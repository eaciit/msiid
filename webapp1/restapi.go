package main

import (
	"eaciit/colony-core"
	"strings"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type RestAPI struct {
}

func (d *RestAPI) Index(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputTemplate
	return nil
}

func (d *RestAPI) Metadata(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson

	metadataRequestModel := struct {
		ModelName string
	}{}

	metadataResponseModel := struct {
		ModelName string
		Fields    []*clncore.DataField
	}{}

	result := toolkit.NewResult()
	if e := ctx.GetPayload(&metadataRequestModel); e != nil {
		modelnameQuery := ctx.Query("modelname")
		if modelnameQuery == "" {
			return result.SetErrorTxt("Unable to get metadata: " + e.Error())
		}
		metadataRequestModel.ModelName = modelnameQuery
	}

	modelnameQuery := strings.ToLower(metadataRequestModel.ModelName)
	if modelnameQuery == "connection" {
		if model := cmm.Get("clncore.DataConnection"); model != nil {
			metadataResponseModel.Fields = model.FieldArray()
		}
	} else if modelnameQuery == "datamodel" {
		if model := cmm.Get("clncore.DataModel"); model != nil {
			metadataResponseModel.Fields = model.FieldArray()
		}
	}

	metadataResponseModel.ModelName = metadataRequestModel.ModelName

	result.SetData(&metadataResponseModel)

	return result
}

func (d *RestAPI) Populate(ctx *knot.WebContext) interface{} {
	ctx.Config.OutputType = knot.OutputJson

	type PopulateModel struct {
		ID    string `bson:"_id" json:"_id"`
		Title string
		Email string
	}

	models := []*PopulateModel{}

	for i := 0; i <= 10; i++ {
		iStr := toolkit.ToString(i)
		models = append(models, &PopulateModel{
			"data " + iStr,
			"DataSet " + iStr,
			"email" + iStr + "@email.com",
		})
	}

	return toolkit.NewResult().SetData(models)
}
