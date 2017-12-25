package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"

	"housing"
)

var (
	dirname string
)

func init() {
	flag.StringVar(&dirname, "dirname", "", "directory containing 實價登錄 files")
}

func printAddress(fname string, rowID int, row []string) error {
	if row[2] == "" {
		glog.Errorf("empty address %s:%d %+v", fname, rowID, row)
		return nil
	}

	fmt.Printf("%s\n", row[2])
	return nil
}

func main() {
	flag.Parse()

	err := housing.ScanDir(dirname, printAddress)
	if err != nil {
		glog.Errorf("%+v", err)
	}
}
