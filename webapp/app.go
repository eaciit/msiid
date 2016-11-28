package main

import (
	"os"
	"path/filepath"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"

	"github.com/eaciit/config"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	"eaciit/colony-core"
)

var cmm = clncore.NewDataModelManager()

func init() {
	cmm.SetObj(new(clncore.DataModel))
	cmm.SetObj(new(clncore.DataConnection))
}

func App(wd string) *knot.App {
	app := knot.NewApp(config.Get("name").(string))
	if wd == "" {
		wd, _ = os.Getwd()
	}
	app.ViewsPath = filepath.Join(wd, "views")
	app.LayoutTemplate = "_layout.html"
	app.Static("static", filepath.Join(wd, "assets"))
	app.Static("views", filepath.Join(wd, "views"))
	app.Register(new(RestAPI))
	app.Register(new(Dashboard))
	app.Register(new(DataManager))
	app.Register(new(Orchestrator))
	app.Register(new(Studio))
	app.DefaultOutputType = knot.OutputHtml
	return app
}

var log *toolkit.LogEngine

func main() {
	log, _ := toolkit.NewLog(true, false, "", "", "")

	configpath, _ := os.Getwd()
	configpath = filepath.Join(configpath, "..", "config", "app.json")
	econfig := config.SetConfigFile(configpath)
	if econfig != nil {
		log.Error("Error loading config file " + econfig.Error())
	}

	port := int(config.GetDefault("port", 9100).(float64))
	serveraddress := config.GetDefault("server", "0.0.0.0").(string)
	wd := config.GetDefault("workingpath", "").(string)
	app := App(wd)
	knot.StartApp(app, toolkit.Sprintf("%s:%d", serveraddress, port))
}

func prepareConn() dbox.IConnection {
	ci := &dbox.ConnectionInfo{}
	ciconn := config.Get("connections").(map[string]interface{})["default"].(map[string]interface{})

	ci.Host = ciconn["host"].(string)
	ci.UserName = ciconn["user"].(string)
	ci.Password = ciconn["password"].(string)
	ci.Database = ciconn["dbname"].(string)

	conn, e := dbox.NewConnection("mongo", ci)
	if e != nil {
		panic(e.Error())
	}
	return conn
}
