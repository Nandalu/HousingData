package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"math/rand"
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"housing/transaction"
	"housing/util/jinma"
)

var (
	infile     string
	jinmaToken string
	randomSeed int64
)

func init() {
	flag.StringVar(&infile, "infile", "", "input file containing messages")
	flag.StringVar(&jinmaToken, "jinmaToken", "", "Jinma user token")
	flag.Int64Var(&randomSeed, "randomSeed", 0, "random seed")
}

func handleMsg(rowID int, msg jinma.Msg, tsct transaction.Transaction) error {
	// Use the transaction date as the sortkey.
	// To avoid collided sortkeys, randomly a time interval.
	skf64 := float64(tsct.A交易年月日)
	skf64 += float64(rand.Intn(24*60*60 - 1))
	skf64 += rand.Float64()

	updatedMsg, err := jinma.MsgUpdate(msg.ID, jinmaToken, nil, &skf64)
	if err != nil {
		return errors.Wrap(err, "jinma.MsgUpdate")
	}

	glog.Infof("%d %+v", rowID, updatedMsg)
	return nil
}

func scanMsgs(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return errors.Wrap(err, "os.Open")
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	i := -1
	for scanner.Scan() {
		i += 1
		msg := jinma.Msg{}
		if err := json.Unmarshal([]byte(scanner.Text()), &msg); err != nil {
			return errors.Wrap(err, "json.Unmarshal msg")
		}
		tsct := transaction.Transaction{}
		if err := json.Unmarshal([]byte(msg.Body), &tsct); err != nil {
			return errors.Wrap(err, "json.Unmarshal msg.Body")
		}
		if err := handleMsg(i, msg, tsct); err != nil {
			return errors.Wrap(err, "handleMsg")
		}
	}
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "scanner.Err")
	}
	return nil
}

func main() {
	flag.Parse()
	rand.Seed(randomSeed)

	if err := scanMsgs(infile); err != nil {
		glog.Fatalf("%+v", err)
	}
}
