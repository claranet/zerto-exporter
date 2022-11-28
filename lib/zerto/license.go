//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//
package zerto

import "encoding/json"
//import "github.com/prometheus/log"

//
// /v1/license
//
type ZertoLicense struct {
	Details						ZertoLicenseDetails
	Usage							ZertoLicenseUsage
}

type ZertoLicenseDetails struct {
	ExpiryTime				string
	LicenseKey				string
	LicenseType				string
	MaxVms						int
}

type ZertoLicenseUsage struct {
	SitesUsage				*ZertoLicenseSiteUsage
	TotalVmsCount			int
}

type ZertoLicenseSiteUsage struct {
	ProtectedVmsCount	float64
	SiteIdentifier		string
	SiteName					string
}
//
// Action: /v1/license
//
// Fetch Information about local Zerto instance
//
func (z *Zerto) LicenseInformations() ZertoLicense {
	resp, _ := z.makeRequest("GET", "/license", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d ZertoLicense
	data.Decode(&d)

	//log.Debug(d)

	return d
}
