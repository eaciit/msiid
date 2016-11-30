package main

import (
	"eaciit/msiid/helper"
	"os"

	"path/filepath"

	"strings"

	"io/ioutil"

	"flag"

	"errors"

	"time"

	"github.com/eaciit/config"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
)

var (
	wdpath            string
	datapath, mappath string
	log               *toolkit.LogEngine

	fileflag = flag.String("file", "*", "-file=bgpf --will read all file that contains bgpf")
	skipflag = flag.Int("skip", 0, "-skip=100 skip first 100 records. If not provided will skip none")
	takeflag = flag.Int("take", 0, "-take=100 process first 100 records after skipped records. If not provided will read all data")
)

func main() {
	flag.Parse()

	wdpath, _ = os.Getwd()
	cfgpath := filepath.Join(wdpath, "..", "config", "app.json")
	config.SetConfigFile(cfgpath)

	log, _ = toolkit.NewLog(true, false, "", "", "")
	defer log.Close()

	datapath = translateFileName(config.Get("source").(string))
	mappath = translateFileName(config.Get("map").(string))
	log.Info("Reading data " + datapath)

	files, _ := ioutil.ReadDir(datapath)
	for _, file := range files {
		if *fileflag == "*" {
			processfile(file, *takeflag, *skipflag)
		} else {
			if strings.Contains(file.Name(), *fileflag) {
				processfile(file, *takeflag, *skipflag)
			}
		}
	}
}

func translateFileName(p string) string {
	return strings.Replace(p, "$wdroot", wdpath, 1)
}

func prepareConn() dbox.IConnection {
	ci := &dbox.ConnectionInfo{}
	ciconn := config.Get("connections").(map[string]interface{})["default"].(map[string]interface{})

	ci.Host = ciconn["host"].(string)
	ci.UserName = ciconn["user"].(string)
	ci.Password = ciconn["password"].(string)
	ci.Database = ciconn["dbname"].(string)

	conn, _ := dbox.NewConnection("mongo", ci)
	return conn
}

func processfile(f os.FileInfo, take, skip int) error {
	conn := prepareConn()
	conn.Connect()
	defer conn.Close()
	tablename := strings.Split(f.Name(), ".")[0]

	filenamepath := filepath.Join(datapath, f.Name())
	mapfilepath := filepath.Join(mappath, f.Name())

	//-- get map
	if !toolkit.IsFileExist(mapfilepath) {
		return errors.New("Map is not exist: " + mapfilepath)
	}

	t0 := time.Now()
	maps := []toolkit.M{}
	mapidx := 0
	mf := &helper.FlatFile{}
	mf.Name = mapfilepath
	mf.IterFn = func(t string) error {
		if strings.HasSuffix(t, ",") {
			t = t[:len(t)-1]
		}
		ts := strings.Split(t, " ")
		fieldname := strings.ToLower(strings.Trim(ts[0], " "))

		fieldtypes := strings.Split(ts[1], "(")
		fieldtype := strings.Replace(strings.ToLower(strings.Trim(fieldtypes[0], " ")), ",", "", -1)
		fieldsize := 0
		//log.Info(fieldtype)
		if fieldtype == "integer" || strings.HasSuffix(fieldtype, "int") {
			fieldsize = 4
		} else if fieldtype != "decimal" {
			fieldsize = toolkit.ToInt(fieldtypes[1][:len(fieldtypes[1])-1], toolkit.RoundingAuto)
		} else {
			lens := strings.Split(fieldtypes[1], ",")
			fieldsize = toolkit.ToInt(lens[0], toolkit.RoundingAuto)
		}
		field := toolkit.M{}.Set("name", fieldname).Set("fieldtype", fieldtype).Set("size", fieldsize)
		maps = append(maps, field)
		mapidx++
		return nil
	}
	if e := mf.Open(); e != nil {
		return toolkit.Errorf("Unable to open %s: %s", mapfilepath, e.Error())
	} else {
		defer mf.Close()
	}
	mf.Exec(0, 0)
	qmap := conn.NewQuery().SetConfig("multiexec", true).From("metadata").Save()
	qmap.Exec(toolkit.M{}.Set("data", toolkit.M{}.Set("_id", tablename).Set("model", maps)))
	qmap.Close()

	//-- get line count
	cmdout, err := toolkit.RunCommand("wc", "-l", filenamepath)
	if err != nil {
		return err
	}
	linecount := 0
	cmdouts := strings.Split(cmdout, " ")
	for _, o := range cmdouts {
		linecountTemp := toolkit.ToInt(o, toolkit.RoundingAuto)
		if linecountTemp != 0 {
			linecount = linecountTemp
		}
	}
	log.Info(toolkit.Sprintf("Processing %s - %d lines", f.Name(), linecount))

	q := conn.NewQuery().SetConfig("multiexec", true).From(tablename).Save()
	//defer q.Close()

	iread := 0
	stage := linecount / 100
	nexttarget := stage

	ff := &helper.FlatFile{}
	ff.Name = filenamepath
	ff.IterFn = func(t string) error {
		records := toolkit.M{}

		//log.Info("Read " + t)
		if len(t) == 0 {
			return nil
		}

		ts := strings.Split(t, "\t")
		for k, v := range ts {
			if k < len(maps) {
				fi := maps[k]
				ftype := fi.Get("fieldtype").(string)
				if ftype == "decimal" {
					records.Set(fi.Get("name").(string), toolkit.ToFloat64(v, 4, toolkit.RoundingAuto))
				} else if strings.Contains(ftype, "int") {
					records.Set(fi.Get("name").(string), toolkit.ToInt(v, toolkit.RoundingAuto))
				} else {
					records.Set(fi.Get("name").(string), v)
				}
			}
		}
		iread++

		records.Set("_id", iread)
		q.Exec(toolkit.M{}.Set("data", records))

		if iread == nexttarget {
			log.Info(toolkit.Sprintf("Processing %d lines (%2.0f%%) - %v",
				iread,
				toolkit.ToFloat64(iread*100, 0, toolkit.RoundingAuto)/toolkit.ToFloat64(linecount, 0, toolkit.RoundingAuto),
				time.Since(t0)))
			nexttarget += stage
		}
		return nil
	}
	ff.Open()
	defer ff.Close()
	return ff.Exec(take, skip)

	return nil
}
