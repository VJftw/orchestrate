package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type APIClient struct {
	BaseURI        string
	ResponseBody   io.ReadCloser
	ResponseStatus int
	HTTPClient     *http.Client
}

func (a *APIClient) Get(uri string) error {
	res, err := http.Get(strings.Join([]string{a.BaseURI, uri}, ""))
	if err != nil {
		return err
	}
	a.ResponseStatus = res.StatusCode

	// body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// json.Unmarshal(body, &a.ResponseJSON)
	return nil
}

func (a *APIClient) Post(uri string, body interface{}) error {
	a.HTTPClient = &http.Client{}

	urlStr := strings.Join([]string{a.BaseURI, uri}, "")
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", urlStr, bytes.NewReader(bodyJSON))

	res, err := a.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	a.ResponseBody = res.Body
	a.ResponseStatus = res.StatusCode

	return nil
}

func (a *APIClient) UnmarshalTo(v interface{}) error {
	resBody, err := ioutil.ReadAll(a.ResponseBody)
	if err != nil {
		return err
	}

	json.Unmarshal(resBody, v)
	return nil
}
