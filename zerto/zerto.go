//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

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

type Zerto struct {
	url		string
	username	string
	password	string

	sessionToken	string
}

func (z *Zerto) makeRequest(reqType string, action string, p RequestParams) (*http.Response, error)  {
	_url := z.url + "/v1/" + action

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
	var netClient = http.Client{Transport: tr}

	body := p.body

	_url += "?" + p.params.Encode()

	req, err := http.NewRequest(reqType, _url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/JSON")

	if z.sessionToken == "" {
		req.SetBasicAuth(z.username, z.password)
	} else {
		req.Header.Set("X-Zerto-Session", z.sessionToken)
	}

	resp, err := netClient.Do(req)
	if err != nil {	log.Fatal(err); return nil, err }

	return resp, nil
}



func NewZerto(url string, username string, password string) *Zerto {
//	log.SetOutput(os.Stdout)
//	log.SetPrefix("Zerto Logger")

	return &Zerto {
		url: url,
		username: username,
		password: password,
		sessionToken: "",
	}
}
