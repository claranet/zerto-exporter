//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import (
 	log "github.com/sirupsen/logrus"
	"encoding/json"
	"net/url"
	"errors"
)

func (z *Zerto) IsSessionOpen() bool {
	return z.sessionToken != ""
}

func (z *Zerto) OpenSession() error {

	// Try to login to ZVM by zerto session manager
	// e.g. Windows based ZVM
	resp, err := z.makeRequest("POST", "/v1/session/add", RequestParams{})
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		// StatusCode is 404 if /v1/session/add does not exists
		// This happened on Linux based ZVM
		z.authProvider = "keycloak"

		body := url.Values{}
		body.Add("grant_type", "password")
		body.Add("client_id", "zerto-client")
		body.Add("username", z.username)
		body.Add("password", z.password)

		resp, _ := z.makeRequest("POST", "/auth/realms/zerto/protocol/openid-connect/token", RequestParams{body: body.Encode()})

		data := json.NewDecoder(resp.Body)
		var d map[string]string
		data.Decode(&d)

		log.Debug(d)
		z.sessionToken = d["access_token"]
	} else {
		z.authProvider = "zerto"
		log.Debug(resp.Header)
		z.sessionToken = resp.Header.Get("x-zerto-session")
	}

	log.Debug("Session Token: " + z.sessionToken)
	if ! z.IsSessionOpen() {
		return errors.New("Cannot open a Session")
	}

	return nil
}

func (z *Zerto) CloseSession() {
	if z.authProvider == "zerto" {
		z.makeRequest("DELETE", "/v1/session", RequestParams{})
	}
	z.sessionToken = ""
}
