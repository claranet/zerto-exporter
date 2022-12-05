//
// Zerto API Interface Wrapper
//
// Author: Martin Weber <martin.weber@de.clara.net>
//

package zerto

import "encoding/json"

type ZertoVm struct {
	ActualRPO		int
//	EnabledActions		*ZertoEnabledActions
	Entities								*ZertoVmEntities
	IOPs										int
	IsVmExists							bool
	JournalHardLimit				*ZertoVmLimit
	JournalUsedStorageMb		int
	JournalWarningThreshold	*ZertoVmLimit
//	LastTest
//	Link
	OrganizationName				string
	OutgoingBandWidthInMbps	float64
	Priority		int
//	ProtectedSite
	ProvisionedStorageInMB	int
//	RecoverySite
	SourceSite							string
	Status									int
	SubStatus								int
	TargetSite							string
	ThroughputInMB					float64
	UsedStorageInMB					int
	VmIdentifier						string
	VmName									string
	Volumes									[]*ZertoVmVolume
	VpgIdentifier						string
	VpgName									string
}

type ZertoVmVolume struct {
	VmVolumeIdentifier	string
}

type ZertoVmLimit struct {
	LimitType			int
	LimitValue		int
}

type ZertoVmEntities struct {
	Protected		int
	Recovery		int
	Source			int
	Target			int
}

func newZertoVm() *ZertoVm {
	vm := ZertoVm{}
	return &vm
}

func (z *Zerto) ListVms() []ZertoVm {
	resp, _ := z.makeRequest("GET", "/v1/vms", RequestParams{})
	data := json.NewDecoder(resp.Body)

	var d []ZertoVm
	data.Decode(&d)

	return d
}
