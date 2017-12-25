package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"housing/util/jinma"
)

var (
	jinmaToken string
	appID      string
)

func init() {
	flag.StringVar(&jinmaToken, "jinmaToken", "", "Jinma user token")
}

func handleMsg(msg jinma.Msg) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}
	fmt.Printf("%s\n", b)
	return nil
}

func scanPartition(partition int, fn func(jinma.Msg) error) error {
	esk := ""
	for {
		resp, err := jinma.MsgsByAppUser(appID, jinmaToken, partition, esk)
		if err != nil {
			return errors.Wrap(err, "jinma.MsgsByAppUser")
		}
		for _, msg := range resp.Msgs {
			if err := fn(msg); err != nil {
				return errors.Wrap(err, "handle msg function")
			}
		}

		esk = resp.LastEvaluatedKey
		if esk == "" {
			break
		}
	}
	return nil
}

func scanAllPartitions(fn func(jinma.Msg) error) error {
	for partition := 0; partition < 1536; partition++ {
		if err := scanPartition(partition, fn); err != nil {
			return errors.Wrap(err, "scanPartition")
		}
		glog.Infof("finished scanning partition %d", partition)
	}
	return nil
}

func main() {
	flag.Parse()

	// Get the the appID of our token.
	me, err := jinma.Me(jinmaToken)
	if err != nil {
		glog.Fatalf("%+v", err)
	}
	appID = me.App.ID

	// Get all messages.
	if err := scanAllPartitions(handleMsg); err != nil {
		glog.Fatalf("%+v", err)
	}
}
