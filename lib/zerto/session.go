//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto
// import	"github.com/prometheus/log"

import (
 	log "github.com/sirupsen/logrus"
	"errors"
)


func (z *Zerto) IsSessionOpen() bool {
	return z.sessionToken != ""
}

func (z *Zerto) OpenSession() error {

	resp, err := z.makeRequest("POST", "/session/add", RequestParams{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	z.sessionToken = resp.Header.Get("x-zerto-session")
	log.Debug(resp.Header)
	log.Debug("Session Token: " + z.sessionToken)
	if ! z.IsSessionOpen() {
		return errors.New("Cannot open a Session")
	}
	return nil
}

func (z *Zerto) CloseSession() {
	z.makeRequest("DELETE", "/session", RequestParams{})
	z.sessionToken = ""
}
