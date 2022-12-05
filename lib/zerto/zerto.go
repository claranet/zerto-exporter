//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import (
//	"os"
	// "github.com/prometheus/log"
	"net/http"
	"net/url"
	"strings"
	"crypto/tls"

	log "github.com/sirupsen/logrus"
)

type RequestParams struct {
	body, header	string
	params				url.Values
}

type Zerto struct {
	url		string
	username	string
	password	string

	sessionToken	string
	authProvider	string
}

func (z *Zerto) makeRequest(reqType string, action string, p RequestParams) (*http.Response, error)  {
	_url := z.url + "/" + strings.TrimLeft(action, "/")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, }
	var netClient = http.Client{Transport: tr}

	body := p.body
	log.Debug(body)

	if(len(p.params) > 0) {
		_url += "?" + p.params.Encode()
	}

	// log.Debug(_url)
	req, _ := http.NewRequest(reqType, _url, strings.NewReader(body))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	req.Header.Set("Content-Type", "application/json")

	if z.authProvider == "zerto" {
		if z.sessionToken == "" {
			req.SetBasicAuth(z.username, z.password)
		} else {
			req.Header.Set("x-zerto-session", z.sessionToken)
		}
	}

	if z.authProvider == "keycloak" {
		if z.sessionToken == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req.Header.Set("Authorization", "Bearer " + z.sessionToken)
		}
	}

	log.Debug("Reuest Headers ", req.Header)

	resp, _ := netClient.Do(req)
	// if err != nil {	log.Fatal(err); return nil, err }

	return resp, nil
}



func NewZerto(url string, username string, password string) *Zerto {
	return &Zerto {
		url: url,
		username: username,
		password: password,
		sessionToken: "",
		authProvider: "zerto",
	}
}
