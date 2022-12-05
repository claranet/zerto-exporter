//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

//
// /v1/localsite
//
type ZertoLocalsite struct {
	ContactEmail			string
	ContactName			string
	ContactPhone			string
	IpAddress			string
	IsReplicationToSelfEnabled	bool
	Link				*ZertoLink
	Location			string
	SiteIdentifier			string
	SiteName			string
	SiteType			string
	UtcOffsetInMinutes		int
	Version				string
}

type ZertoLink struct {
	Href				string `json:"href"`
	Identifier			string `json:"identifier"`
	Rel				string `json:"rel"`
	Type				string `json:"type"`
}

//
// Action: /v1/localsite
//
// Fetch Information about local Zerto instance
//
func (z *Zerto) Localsite() ZertoLocalsite {
	resp, _ := z.makeRequest("GET", "/v1/localsite", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d ZertoLocalsite
	data.Decode(&d)

	return d
}
