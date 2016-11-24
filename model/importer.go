package model

import (
	"bufio"
	"eaciit/colony-core"
	"os"

	"strings"

	"github.com/eaciit/toolkit"
)

type Importer struct {
	Model clncore.DataModel
}

func (imp *Importer) Load(filename string, mapfilename string, overwrite bool) error {
	var err error
	var reader *os.File

	if filename == "" {
		filename = imp.Model.ID + ".txt"
	}

	if mapfilename == "" {
		mapfilename = imp.Model.ID + ".map"
	}

	//--read map file
	if mapreader, maperr := os.Open(mapfilename); maperr != nil {
		return maperr
	} else {
		defer mapreader.Close()
		mapscanner := bufio.NewScanner(mapreader)
		for mapscanner.Scan() {
			scannedtxt := mapscanner.Text()
			maptxts := strings.Split(scannedtxt, ",")
			for _, maptxt := range maptxts {
				toolkit.Println(maptxt)
			}
		}
	}

	//--read txt file
	reader, err = os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	txtscanner := bufio.NewScanner(reader)
	for txtscanner.Scan() {
		txt := txtscanner.Text()
		toolkit.Println(txt)
	}
	return nil
}
