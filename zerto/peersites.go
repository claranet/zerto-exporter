//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

type ZertoPeersite struct {
	HostName		string
	IncomingThroughputInMb	float64
	Location		string
	OutgoingBandWidth	float64
	PairingStatus		int
	PeerSiteName		string
	Port			int
	ProvisionedStorage	int64
	SiteIdentifier		string
	SiteType		string
	UsedStorage		int64
	Version			string
}

func (z *Zerto) ListPeersites() []ZertoPeersite {
	resp, _ := z.makeRequest("GET", "/peersites", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d []ZertoPeersite
	data.Decode(&d)

	return d
}
