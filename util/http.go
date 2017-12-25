package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

func JSONReq3(method, urlStr string, res interface{}) (*http.Response, []byte, error) {
	return JSONReq6(method, urlStr, nil, nil, http.DefaultClient, res)
}

func JSONReq5(method, urlStr string, body io.Reader, header http.Header, res interface{}) (*http.Response, []byte, error) {
	return JSONReq6(method, urlStr, body, header, http.DefaultClient, res)
}

func JSONReq6(method, urlStr string, body io.Reader, header http.Header, c *http.Client, res interface{}) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "NewRequest")
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "RequestDo")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "ReadAll")
	}
	if res == nil {
		return resp, b, nil
	}
	err = json.Unmarshal(b, res)
	if err != nil {
		errMsg := fmt.Sprintf("json.Unmarshal error %v %s", err, b)
		glog.Errorf(errMsg)
		return nil, nil, errors.Wrap(err, errMsg)
	}
	return resp, b, nil
}
