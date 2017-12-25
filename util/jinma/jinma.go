package jinma

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"housing/util"
)

const (
	host = "http://www.jinma.io"
)

type App struct {
	ID           string
	Name         string
	Icon         string
	MarketingURI string
}

type User struct {
	ID           string
	Name         string
	Picture      string
	ThirdPartyID string
	Privacy      string
	Language     string
}

type Msg struct {
	ID       string
	User     User
	Time     float64
	Body     string
	Lat      float64
	Lng      float64
	SKF64    float64
	Hashtags []string
	App      App
	CustomID string
}

type MsgsByAppUserResp struct {
	Msgs             []Msg
	LastEvaluatedKey string
}

type MeResp struct {
	User User
	App  App
}

func Me(token string) (*MeResp, error) {
	vals := url.Values{
		"Token": {token},
	}
	urlStr := host + "/Me?" + vals.Encode()
	resp := MeResp{}
	httpResp, respBody, err := util.JSONReq3("POST", urlStr, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "util.JSONReq3")
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %d %s", httpResp.StatusCode, respBody)
	}
	return &resp, nil
}

func MsgCreate(token, body string, lat, lng float64, skf64 *float64, customID string) (*Msg, error) {
	vals := url.Values{
		"Lat":   {strconv.FormatFloat(lat, 'f', -1, 64)},
		"Lng":   {strconv.FormatFloat(lng, 'f', -1, 64)},
		"Body":  {body},
		"Token": {token},
	}
	if skf64 != nil {
		vals.Set("SKF64", strconv.FormatFloat(*skf64, 'f', -1, 64))
	}
	if customID != "" {
		vals.Set("CustomID", customID)
	}
	urlStr := host + "/MsgCreate?" + vals.Encode()
	resp := Msg{}
	httpResp, respBody, err := util.JSONReq3("POST", urlStr, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "util.JSONReq3")
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %d %s", httpResp.StatusCode, respBody)
	}
	return &resp, nil
}

func MsgUpdate(id, token string, body []byte, skf64 *float64) (*Msg, error) {
	vals := url.Values{
		"MsgID": {id},
		"Token": {token},
	}
	if len(body) > 0 {
		vals.Set("Body", string(body))
	}
	if skf64 != nil {
		vals.Set("SKF64", strconv.FormatFloat(*skf64, 'f', -1, 64))
	}
	urlStr := host + "/MsgUpdate?" + vals.Encode()
	resp := Msg{}
	httpResp, respBody, err := util.JSONReq3("POST", urlStr, &resp)
	if err != nil {
		glog.Errorf("%+v", err)
		return nil, errors.Wrap(err, "util.JSONReq3")
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %d %s", httpResp.StatusCode, respBody)
	}
	return &resp, nil
}

func MsgsByAppUser(appID, token string, partition int, esk string) (*MsgsByAppUserResp, error) {
	vals := url.Values{
		"AppID": {appID},
		"Token": {token},
		"I":     {fmt.Sprintf("%d", partition)},
	}
	if esk != "" {
		vals.Set("ESK", esk)
	}
	urlStr := host + "/MsgsByAppUser?" + vals.Encode()
	resp := MsgsByAppUserResp{}
	httpResp, respBody, err := util.JSONReq3("GET", urlStr, &resp)
	if err != nil {
		return nil, errors.Wrap(err, "util.JSONReq3")
	}
	if httpResp.StatusCode != 200 {
		return nil, fmt.Errorf("request error: %d %s", httpResp.StatusCode, respBody)
	}
	return &resp, nil
}
