package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type APIClient struct {
	BaseURI        string
	T              *testing.T
	ResponseBody   io.ReadCloser
	ResponseStatus int
	HTTPClient     *http.Client
	BearerToken    string
}

func (a *APIClient) RequestNoBody(method string, uri string) {
	a.HTTPClient = &http.Client{}

	urlStr := strings.Join([]string{a.BaseURI, uri}, "")

	req, _ := http.NewRequest(method, urlStr, nil)
	a.setBearerToken(req)

	res, err := a.HTTPClient.Do(req)
	convey.So(err, convey.ShouldBeNil)
	if err != nil {
		a.T.Error(err)
		return
	}

	a.ResponseBody = res.Body
	a.ResponseStatus = res.StatusCode

	return
}

func (a APIClient) setBearerToken(req *http.Request) {
	if len(a.BearerToken) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", a.BearerToken))
	}
}

func (a *APIClient) RequestWithBody(method string, uri string, body interface{}) {
	a.HTTPClient = &http.Client{}

	urlStr := strings.Join([]string{a.BaseURI, uri}, "")
	bodyJSON, err := json.Marshal(body)
	convey.So(err, convey.ShouldBeNil)
	if err != nil {
		fmt.Println(err)
		a.T.Error(err)
		a.T.Fail()
		return
	}
	fmt.Println(string(bodyJSON))
	req, _ := http.NewRequest(method, urlStr, bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	// a.setBearerToken(req)

	res, err := a.HTTPClient.Do(req)
	convey.So(err, convey.ShouldBeNil)
	if err != nil {
		a.T.Error(err)
		a.T.Fail()
		return
	}

	a.ResponseBody = res.Body
	a.ResponseStatus = res.StatusCode

	return
}

func (a *APIClient) UnmarshalTo(v interface{}) {
	resBody, err := ioutil.ReadAll(a.ResponseBody)
	convey.So(err, convey.ShouldBeNil)
	if err != nil {
		a.T.Error(err)
		return
	}

	json.Unmarshal(resBody, v)
}
