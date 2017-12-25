package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"housing/transaction"
	"housing/util/jinma"
)

var (
	infile       string
	infileOffset int
	jinmaToken   string
	randomSeed   int64
)

func init() {
	flag.StringVar(&infile, "infile", "", "input file containing the parsed transactions")
	flag.IntVar(&infileOffset, "infileOffset", 0, "line offset from which we should read from infile")
	flag.StringVar(&jinmaToken, "jinmaToken", "", "Jinma user token")
	flag.Int64Var(&randomSeed, "randomSeed", 0, "random seed")
}

func create(inTs transaction.Transaction) (*jinma.Msg, error) {
	// Make a copy of the transaction and remove the unneeded fields.
	ts := inTs
	// These fields are unneeded because they are contained in the jinma.Msg itself.
	ts.A編號 = ""
	ts.Lat = 0
	ts.Lng = 0

	tsbody, err := json.Marshal(ts)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	// Use the transaction date as the sortkey.
	// To avoid collided sortkeys, randomly a time interval.
	skf64 := float64(ts.A交易年月日)
	skf64 += float64(rand.Intn(24*60*60 - 1))
	skf64 += rand.Float64()

	customID := inTs.A編號
	if customID == "" {
		return nil, errors.Wrap(err, fmt.Sprintf("empty customID for %+v", inTs))
	}

	msg, err := jinma.MsgCreate(jinmaToken, string(tsbody), ts.Lat, ts.Lng, &skf64, customID)
	if err != nil {
		return nil, errors.Wrap(err, "jinma.MsgCreate")
	}
	return msg, nil
}

func pubFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return errors.Wrap(err, "open")
	}
	defer f.Close()
	fbody, err := ioutil.ReadAll(f)
	if err != nil {
		return errors.Wrap(err, "readall")
	}
	lines := strings.Split(string(fbody), "\n")
	for i, line := range lines {
		// When resuming a previous run, update the offset here
		// to continue from where we left.
		if i < infileOffset {
			continue
		}

		if line == "" {
			continue
		}

		ts := transaction.Transaction{}
		if err := json.Unmarshal([]byte(line), &ts); err != nil {
			return errors.Wrap(err, "unmarshal line")
		}

		msg, err := create(ts)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("create error %d %+v", i, ts))
		}
		glog.Infof("created row: %d, msg.ID: %s, ts: %+v", i, msg.ID, ts)
	}
	return nil
}

func main() {
	flag.Parse()
	rand.Seed(randomSeed)

	if err := pubFile(infile); err != nil {
		glog.Errorf("%+v", err)
	}
}
