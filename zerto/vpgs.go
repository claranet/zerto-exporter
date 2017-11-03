//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

//
// /v1/vpgs
//
type ZertoVpg struct {
	ActualRPO		int
	BackupEnabled		bool
	ConfiguredRpoSeconds	int
	IOPs			int
	HistoryStatusApi	HistoryStatusApi
	OrganizationName	string
	Priority		int
	ProgressPercentage	int
	SourceSite		string
	TargetSite		string
	ThroughputInMB		float64
	UsedStorageInMB		int
	VpgName			string
	VmsCount		int
	Status			int
	SubStatus		int
}

type HistoryStatusApi struct {
	ActualHistoryInMinutes		int
	ConfiguredHistoryInMinutes	int
	EarliestCheckpoint		EarliestCheckpoint
}

type EarliestCheckpoint struct {
	CheckpointIdentifier		string
	Tag				string
	TimeStamp			string
	Vss				bool
}

func (z *Zerto) ListVpg() []ZertoVpg {
	resp, _ := z.makeRequest("GET", "/vpgs", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var c []ZertoVpg
	data.Decode(&c)
	for i:=0;i<len(c);i++ { if c[i].OrganizationName == "" { c[i].OrganizationName = "LOCAL" } }

	return c
}
