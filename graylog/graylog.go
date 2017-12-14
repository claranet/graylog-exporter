
package graylog

import (
//	"os"
	"log"
	"net/http"
	"net/url"
	"strings"
	"crypto/tls"
)

type RequestParams struct {
	body, header	string
	params		url.Values
}

type Graylog struct {
	url		string
	username	string
	password	string
}

func (g *Graylog) makeRequest(reqType string, action string) (*http.Response, error) {
	return g.makeRequestWithParams(reqType, action, RequestParams{})
}

func (g *Graylog) makeRequestWithParams(reqType string, action string, p RequestParams) (*http.Response, error)  {
	_url := strings.Trim(g.url, "/") + action

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
	var netClient = http.Client{Transport: tr}

	body := p.body
	_url += "?" + p.params.Encode()

	req, err := http.NewRequest(reqType, _url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/JSON")

	req.SetBasicAuth(g.username, g.password)

	resp, err := netClient.Do(req)
	if err != nil {	log.Fatal(err); return nil, err }

	return resp, nil
}



func NewGraylog(url string, username string, password string) *Graylog {
//	log.SetOutput(os.Stdout)
//	log.SetPrefix("Graylog Logger")

	return &Graylog {
		url: url,
		username: username,
		password: password,
	}
}
