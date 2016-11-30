package main

import (
	"net/http"
	"os"

	"path/filepath"

	"github.com/eaciit/config"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

func App() *knot.App {
	app := knot.NewApp("mfg")
	wd, _ := os.Getwd()
	//wd = filepath.Join(wd)
	app.ViewsPath = filepath.Join(wd, "views")
	toolkit.Println(app.ViewsPath)
	app.LayoutTemplate = "_layout.html"
	app.Static("static", filepath.Join(wd, "assets"))
	app.Register(&Dashboard{})
	app.Register(&Rest{})
	app.Register(&Login{})
	app.DefaultOutputType = knot.OutputHtml
	return app
}

func DB() *orm.DataContext {
	return orm.New(conn())
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
	//wd := config.GetDefault("workingpath", "").(string)

	otherRoutes := map[string]knot.FnContent{
		"/": func(r *knot.WebContext) interface{} {
			http.Redirect(r.Writer, r.Request, "/login/index", 301)
			return true
		},
	}

	app := App()
	knot.StartAppWithFn(app, toolkit.Sprintf("%s:%d", serveraddress, port), otherRoutes)
}

func conn() dbox.IConnection {
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
	conn.Connect()
	return conn
}
