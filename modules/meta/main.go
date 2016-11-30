package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/eaciit/config"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
)

var (
	conn  dbox.IConnection
	metas = []*toolkit.M{}
)

func main() {
	wdpath, _ := os.Getwd()
	cfgpath := filepath.Join(wdpath, "../..", "config", "app.json")
	config.SetConfigFile(cfgpath)

	conn := prepareConn()
	defer conn.Close()

	toolkit.Println("Reading Meta")
	cmetas, _ := conn.NewQuery().From("metadata").Select().Cursor(nil)
	defer cmetas.Close()
	cmetas.Fetch(&metas, 0, false)

	ctmps, _ := conn.NewQuery().From("tmpmeta").Select().Cursor(nil)
	defer ctmps.Close()
	tmps := []*toolkit.M{}
	ctmps.Fetch(&tmps, 0, false)

	for _, v := range tmps {
		UpdateMeta(v)
	}

	toolkit.Println("Saving")
	qsave := conn.NewQuery().From("metadata").SetConfig("multiexec", true).Save()
	defer qsave.Close()
	for _, v := range metas {
		qsave.Exec(toolkit.M{}.Set("data", *v))
	}
}

func UpdateMeta(tmp *toolkit.M) {
	table := strings.ToLower(tmp.GetString("TableName"))
	field := strings.ToLower(tmp.GetString("ColumnName"))
	colindex := tmp.GetInt("ColumnID")
	description := tmp.GetString("description")
	if description == "" || description == "0" {
		description = field
	}
	toolkit.Printfn("%s %s %s\n", table, field, description)

	for _, v := range metas {
		if v.GetString("_id") == table {
			models := v.Get("model").([]interface{})
			for _, model := range models {
				m := model.(toolkit.M)
				if m.GetString("name") == field {
					m.Set("colindex", colindex)
					m.Set("description", description)
				}
				model = m
			}
			v.Set("model", models)
		}
	}
}

func prepareConn() dbox.IConnection {
	ci := &dbox.ConnectionInfo{}
	ciconn := config.Get("connections").(map[string]interface{})["default"].(map[string]interface{})

	ci.Host = ciconn["host"].(string)
	ci.UserName = ciconn["user"].(string)
	ci.Password = ciconn["password"].(string)
	ci.Database = ciconn["dbname"].(string)

	conn, _ := dbox.NewConnection("mongo", ci)
	conn.Connect()
	return conn
}
