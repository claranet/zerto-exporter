//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto
// import	"github.com/prometheus/log"

import (
 	log "github.com/sirupsen/logrus"
	"encoding/json"
	"errors"
)

func (z *Zerto) IsSessionOpen() bool {
	return z.sessionToken != ""
}

func (z *Zerto) OpenSession() error {

	resp, err := z.makeRequest("POST", "/v1/session/add", RequestParams{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	if resp.StatusCode >= 400 {
		z.authProvider = "keycloak"

		body := "grant_type=password&client_id=zerto-client&username=" + z.username + "&password=" + z.password
		log.Debug(body)
		resp, _ := z.makeRequest("POST", "/auth/realms/zerto/protocol/openid-connect/token", RequestParams{body: body})

		data := json.NewDecoder(resp.Body)
		var d map[string]string
		data.Decode(&d)

		z.sessionToken = d["access_token"]
	} else {
		z.sessionToken = resp.Header.Get("x-zerto-session")
	}

	log.Debug(resp.Header)
	log.Debug("Session Token: " + z.sessionToken)
	if ! z.IsSessionOpen() {
		return errors.New("Cannot open a Session")
	}
	return nil
}

func (z *Zerto) CloseSession() {
	z.makeRequest("DELETE", "/v1/session", RequestParams{})
	z.sessionToken = ""
}
