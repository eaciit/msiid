package helper

import (
	"bufio"
	"errors"
	"os"

	"github.com/eaciit/toolkit"
)

type FlatFile struct {
	Name       string
	readerptr  *os.File
	scannerptr *bufio.Scanner

	IterFn func(txt string) error
}

func (f *FlatFile) Open() error {
	if f.Name == "" {
		return errors.New("File name can not be blank")
	}

	if f.readerptr == nil {
		var err error
		f.readerptr, err = os.Open(f.Name)
		if err != nil {
			return errors.New("Error open file " + f.Name + ": " + err.Error())
		}
	}

	if f.scannerptr == nil {
		f.scannerptr = bufio.NewScanner(f.readerptr)
	}

	return nil
}

func (f *FlatFile) Exec(take, skip int) error {
	if f.scannerptr == nil {
		return errors.New("Scanner is not ready. Please call Open() method")
	}

	i := 0
	p := 0
	for f.scannerptr.Scan() {
		i++
		if i > skip {
			p++
			scanout := f.scannerptr.Text()
			if f.IterFn != nil {
				eiter := f.IterFn(scanout)
				if eiter != nil {
					return toolkit.Errorf("Fail to scan line %d: %s", i, eiter.Error())
				}
			}
			if take > 0 && p == take {
				return nil
			}
		}
	}

	return nil
}

func (f *FlatFile) Close() {
	if f.readerptr != nil {
		f.readerptr.Close()
		f.readerptr = nil
	}

	f.scannerptr = nil
}
