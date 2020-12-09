//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

type ZertoVra struct {
	DatastoreClusterIdentifier	string
	DatastoreClusterName		string
	DatastoreIdentifier		string
	DatastoreName			string
	HostIdentifier			string
	HostVersion			string
	IpAddress			string

	MemoryInGB			int
	NetworkIdentifier		string
	NetworkName			string
	Progress			int
	ProtectedCounters		*ZertoVraCounter
	RecoveryCounters		*ZertoVraCounter
	SelfProtectedVpgs		int
	Status				int
	VraGroup			string
//	VraIdentifier			string
	VraName				string
//	VraNetworkDataApi
	VraVersion			string
}

type ZertoVraCounter struct {
	Vms		int
	Volumes		int
	Vpgs		int
}

func (z *Zerto) ListVras() []ZertoVra {
	resp, _ := z.makeRequest("GET", "/vras", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d []ZertoVra
	data.Decode(&d)

	return d
}
