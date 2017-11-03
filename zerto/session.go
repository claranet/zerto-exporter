//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

func (z *Zerto) IsSessionOpen() bool {
	return z.sessionToken != ""
}

func (z *Zerto) OpenSession() {

	resp, err := z.makeRequest("POST", "/session/add", RequestParams{body: `{"AuthenticationMethod": "1"}`})
	if err != nil {
//		log.Fatal(err)
	}
	z.sessionToken = resp.Header.Get("X-Zerto-Session")
//	log.Printf("%#v", resp.Body)

///	log.Print("Session Token: " + z.sessionToken)

}

func (z *Zerto) CloseSession() {
	z.makeRequest("DELETE", "/session", RequestParams{})
	z.sessionToken = ""
}
