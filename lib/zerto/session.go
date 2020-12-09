//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto
import	"github.com/prometheus/log"


func (z *Zerto) IsSessionOpen() bool {
	return z.sessionToken != ""
}

func (z *Zerto) OpenSession() {

	resp, err := z.makeRequest("POST", "/session/add", RequestParams{})
	if err != nil {
//		log.Fatal(err)
	}
	z.sessionToken = resp.Header.Get("x-zerto-session")
	log.Debug(resp.Header)
	log.Debug("Session Token: " + z.sessionToken)

}

func (z *Zerto) CloseSession() {
	z.makeRequest("DELETE", "/session", RequestParams{})
	z.sessionToken = ""
}
