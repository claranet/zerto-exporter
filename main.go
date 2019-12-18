
//
// zerto-exporter
//
// Prometheus Exportewr for Zerto API
//
// Author: Martin Weber <martin.weber@de.clara.net>
// Company: Claranet GmbH
//

package main

import (
	"github.com/claranet/zerto-exporter/zerto"

	"os"
	"fmt"
	"flag"
	"net/http"
	"time"
	"regexp"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/log"
)

const AppVersion = "0.1.1"

var (
	namespace		= "zerto"
	zertoUrl		= flag.String("zerto.url", "", "Zerto URL to connect https://zvm.local.host:9669")
	zertoUser		= flag.String("zerto.username", "", "Zerto API User")
	zertoPassword		= flag.String("zerto.password", "", "Zerto API User Password")
	zertoMaxSessionAge	= flag.Int("zerto.session-age", 3600, "Zerto Session recreation Time")
	listenAddress		= flag.String("listen-address", ":9403", "The address to lisiten on for HTTP requests.")
	version         = flag.Bool("version", false, "Prints current version")
)

var (
	// Zerto API
	zertoApi		*zerto.Zerto
	// Current Session Age
	zertoSessionAge		int64		= 0
)

type Exporter struct {
	ZertoVpgActualRpo							*prometheus.GaugeVec
	ZertoVpgCountVm								*prometheus.GaugeVec
	ZertoVpgThroughputInMB				*prometheus.GaugeVec
	ZertoVpgUsedStorageInMB				*prometheus.GaugeVec
	ZertoVpgConfiguredRpoSeconds	*prometheus.GaugeVec
	ZertoVpgIops									*prometheus.GaugeVec
	ZertoVpgStatus								*prometheus.GaugeVec
	ZertoVpgSubStatus							*prometheus.GaugeVec

	ZertoLocalsiteReplicationToSelf	*prometheus.GaugeVec
	ZertoPeersitePairingStatus			*prometheus.GaugeVec

	ZertoVraStatus									*prometheus.GaugeVec
	ZertoVraProtectedVms						*prometheus.GaugeVec
	ZertoVraProtectedVolumes				*prometheus.GaugeVec
	ZertoVraProtectedVpgs						*prometheus.GaugeVec

	ZertoVmActualRpo								*prometheus.GaugeVec
	ZertoVmIops											*prometheus.GaugeVec
	ZertoVmJournalUsedStorageMb			*prometheus.GaugeVec
	ZertoVmProvisionedStorageInMB		*prometheus.GaugeVec
	ZertoVmOutgoingBandWidthInMbps	*prometheus.GaugeVec
	ZertoVmThroughputInMB						*prometheus.GaugeVec
	ZertoVmUsedStorageInMB					*prometheus.GaugeVec
	ZertoVmStatus										*prometheus.GaugeVec

	ZertoTaskStatus									*prometheus.GaugeVec
	ZertoTaskProgress								*prometheus.GaugeVec
	ZertoTaskStarted								*prometheus.GaugeVec
	ZertoTaskCompleted							*prometheus.GaugeVec

	ZertoAlertsCount								*prometheus.GaugeVec

	ZertoLicenseExpiryTime					*prometheus.GaugeVec
}

func NewExporter() *Exporter {
	defaultLabels := []string {"vpg", "org"}
	vmDefaultLabels := []string{"org", "vm", "vpg"}
	taskDefaultLabels := []string{"taskid", "type"}
	alertsDefaultLabels := []string{"level"}

	return &Exporter{
		ZertoVpgCountVm: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_vm_count",
			Help: "Count Virtual machines in VPG.",
		}, defaultLabels, ),
		ZertoVpgThroughputInMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_throughput",
			Help: "Throughput of this VPG in MB (megabytes)",
		}, defaultLabels, ),
		ZertoVpgUsedStorageInMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_used_storage",
			Help: "Used Storage in MB (megabytes)",
		}, defaultLabels, ),
		ZertoVpgActualRpo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_actual_rpo",
			Help: "Actual RPO of the VPG.",
		}, append(defaultLabels, "source", "target"), ),
		ZertoVpgConfiguredRpoSeconds: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_configured_rpo",
			Help: "Configured RPO of the VPG.",
		}, defaultLabels, ),
		ZertoVpgIops: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_iops",
			Help: "IOPs of the VPG.",
		}, defaultLabels, ),
		ZertoVpgStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_status",
			Help: "Status of VPG. 0=Initializing, 1=MeetingSLA, 2=NotMeetingSLA, 3=RpoNotMeetingSLA, 4=HistoryNotMeetingSLA, 5=FailingOver, 6=Moving, 7=Deleting, 8=Recovered",
		}, defaultLabels, ),
		ZertoVpgSubStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vpg_substatus",
			Help: "SubStatus of VPG.",
		}, defaultLabels, ),

		ZertoLocalsiteReplicationToSelf: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "localsite_replication_to_self",
			Help: "Is Replication to self enabled on localsite",
		}, []string{"location", "version"}, ),
		ZertoPeersitePairingStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "peersite_pairingstatus",
			Help: "Pairing status with remote site. 0=Paired, 1=Unpairing, 2=Unpaired",
		}, []string{"peername", "location", "version"}, ),

		ZertoVraStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vra_status",
			Help: "Status of Zerto VRA Appliance. 0=Installed, 1=UnsupportedEsxVersion, 2=NotInstalled, 3=Installing, 4=Removing, 5=InstallationError, 6=HostPasswordChanged, 7=UpdatingIpSettings, 8=DuringChangeHost, 9=HostInMaintenanceMode",
		}, []string{"vra_group", "vra_name"}, ),
		ZertoVraProtectedVms: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vra_protected_vms",
			Help: "Count Protected VMs of VRA",
		}, []string{"vra_group", "vra_name"}, ),
		ZertoVraProtectedVolumes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vra_protected_volumes",
			Help: "Count Protected volumes of VRA",
		}, []string{"vra_group", "vra_name"}, ),
		ZertoVraProtectedVpgs: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vra_protected_vpgs",
			Help: "Count Protected VPGs of VRA",
		}, []string{"vra_group", "vra_name"}, ),


		ZertoVmActualRpo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_actual_rpo",
			Help: "Acutal RPO of VM",
		}, vmDefaultLabels, ),
		ZertoVmIops: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_iops",
			Help: "IOPs of VM (Input/Outputs per second)",
		}, vmDefaultLabels, ),
		ZertoVmJournalUsedStorageMb: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_journal_used_storage",
			Help: "Used Storage in MB (megabyte)",
		}, vmDefaultLabels, ),
		ZertoVmProvisionedStorageInMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_provisioned_storage",
			Help: "Provisioned storage in MB (megabytes)",
		}, vmDefaultLabels, ),
		ZertoVmOutgoingBandWidthInMbps: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_outgoing_bandwith",
			Help: "Outgoing Bandwith in MB/s (megabytes per seconds)",
		}, vmDefaultLabels, ),
		ZertoVmThroughputInMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_throughput",
			Help: "Throughput in MB (megabyte)",
		}, vmDefaultLabels, ),
		ZertoVmUsedStorageInMB: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_used_storage",
			Help: "Used Storage in MB (megabyte)",
		}, vmDefaultLabels, ),
		ZertoVmStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "vm_status",
			Help: "Status of VM. 0=Initializing, 1=MeetingSLA, 2=NotMeetingSLA, 3=RpoNotMeetingSLA, 4=HistoryNotMeetingSLA, 5=FailingOver, 6=Moving, 7=Deleting, 8=Recovered",
		}, vmDefaultLabels, ),


		ZertoTaskStatus: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "task_status",
			Help: "Status of Task. 0=FirstUnusedValue, 1=InProgress, 2=WaitingForUserInput, 3=Paused, 4=Failed, 5=Stopped, 6=Completed, 7=Cancelling",
		}, taskDefaultLabels, ),
		ZertoTaskProgress: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "task_progress",
			Help: "Status progress of Task oin percent",
		}, taskDefaultLabels, ),
		ZertoTaskStarted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "task_started",
			Help: "Task start time",
		}, taskDefaultLabels, ),
		ZertoTaskCompleted: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "task_completed",
			Help: "Task end time",
		}, taskDefaultLabels, ),

		ZertoAlertsCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "alerts_count",
			Help: "Count Zerto Alerts, level: error|warning",
		}, alertsDefaultLabels, ),

		ZertoLicenseExpiryTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace, Name: "licencse_expireration_time",
			Help: "License Expiry Time",
		}, []string {}, ),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.ZertoVpgThroughputInMB.Describe(ch)
	e.ZertoVpgUsedStorageInMB.Describe(ch)
	e.ZertoVpgCountVm.Describe(ch)
	e.ZertoVpgActualRpo.Describe(ch)
	e.ZertoVpgConfiguredRpoSeconds.Describe(ch)
	e.ZertoVpgIops.Describe(ch)
	e.ZertoVpgStatus.Describe(ch)
	e.ZertoVpgSubStatus.Describe(ch)

	e.ZertoLocalsiteReplicationToSelf.Describe(ch)
	e.ZertoPeersitePairingStatus.Describe(ch)

	e.ZertoVraStatus.Describe(ch)
	e.ZertoVraProtectedVms.Describe(ch)
	e.ZertoVraProtectedVolumes.Describe(ch)
	e.ZertoVraProtectedVpgs.Describe(ch)

	e.ZertoVmActualRpo.Describe(ch)
	e.ZertoVmIops.Describe(ch)
	e.ZertoVmProvisionedStorageInMB.Describe(ch)
	e.ZertoVmJournalUsedStorageMb.Describe(ch)
	e.ZertoVmOutgoingBandWidthInMbps.Describe(ch)
	e.ZertoVmThroughputInMB.Describe(ch)
	e.ZertoVmUsedStorageInMB.Describe(ch)
	e.ZertoVmStatus.Describe(ch)

	e.ZertoTaskStatus.Describe(ch)
	e.ZertoTaskProgress.Describe(ch)
	e.ZertoTaskStarted.Describe(ch)
	e.ZertoTaskCompleted.Describe(ch)

	e.ZertoAlertsCount.Describe(ch)

	e.ZertoLicenseExpiryTime.Describe(ch)
}


func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	openZertoSession()

	vpgs := zertoApi.ListVpg()

	for i:=0; i<len(vpgs); i++ {
		g := e.ZertoVpgCountVm.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].VmsCount))
		g.Collect(ch)

		g = e.ZertoVpgThroughputInMB.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].ThroughputInMB))
		g.Collect(ch)

		g = e.ZertoVpgUsedStorageInMB.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].UsedStorageInMB))
		g.Collect(ch)

		g = e.ZertoVpgActualRpo.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName, vpgs[i].SourceSite, vpgs[i].TargetSite)
		g.Set(float64(vpgs[i].ActualRPO))
		g.Collect(ch)

		g = e.ZertoVpgConfiguredRpoSeconds.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].ConfiguredRpoSeconds))
		g.Collect(ch)

		g = e.ZertoVpgIops.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].IOPs))
		g.Collect(ch)

		g = e.ZertoVpgStatus.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].Status))
		g.Collect(ch)

		g = e.ZertoVpgSubStatus.WithLabelValues(vpgs[i].VpgName, vpgs[i].OrganizationName)
		g.Set(float64(vpgs[i].SubStatus))
		g.Collect(ch)
	}

	local := zertoApi.Localsite()
	g := e.ZertoLocalsiteReplicationToSelf.WithLabelValues(local.Location, local.Version)
	if local.IsReplicationToSelfEnabled { g.Set(1) } else { g.Set(0) }
	g.Collect(ch)

	license := zertoApi.LicenseInformations()
	t, _ := time.Parse("2006-01-02T15:04:05.000Z", license.Details.ExpiryTime)
	g = e.ZertoLicenseExpiryTime.WithLabelValues()
	g.Set(float64(t.Unix()))
	g.Collect(ch)

	peers := zertoApi.ListPeersites()
	for i:=0; i<len(peers);i++ {
		g = e.ZertoPeersitePairingStatus.WithLabelValues(peers[i].PeerSiteName, peers[i].Location, peers[i].Version)
		g.Set(float64(peers[i].PairingStatus))
		g.Collect(ch)
	}

	vras := zertoApi.ListVras()
	for i:=0;i<len(vras);i++ {
		g := e.ZertoVraStatus.WithLabelValues(vras[i].VraGroup, vras[i].VraName)
		g.Set(float64(vras[i].Status))
		g.Collect(ch)
		g = e.ZertoVraProtectedVms.WithLabelValues(vras[i].VraGroup, vras[i].VraName)
		g.Set(float64(vras[i].ProtectedCounters.Vms))
		g.Collect(ch)
		g = e.ZertoVraProtectedVolumes.WithLabelValues(vras[i].VraGroup, vras[i].VraName)
		g.Set(float64(vras[i].ProtectedCounters.Volumes))
		g.Collect(ch)
		g = e.ZertoVraProtectedVpgs.WithLabelValues(vras[i].VraGroup, vras[i].VraName)
		g.Set(float64(vras[i].ProtectedCounters.Vpgs))
		g.Collect(ch)
	}

	vms := zertoApi.ListVms()
	log.Debugf("Count VM's: %v", len(vms))
	for i:=0;i<len(vms);i++ {
		var (
			orgName = vms[i].OrganizationName
			vmName = vms[i].VmName
			vpgName = vms[i].VpgName
		)

		g := e.ZertoVmActualRpo.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].ActualRPO))
		g.Collect(ch)
		g = e.ZertoVmIops.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].IOPs))
		g.Collect(ch)
		g = e.ZertoVmProvisionedStorageInMB.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].ProvisionedStorageInMB))
		g.Collect(ch)
		g = e.ZertoVmJournalUsedStorageMb.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].JournalUsedStorageMb))
		g.Collect(ch)
		g = e.ZertoVmOutgoingBandWidthInMbps.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].OutgoingBandWidthInMbps))
		g.Collect(ch)
		g = e.ZertoVmThroughputInMB.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].ThroughputInMB))
		g.Collect(ch)
		g = e.ZertoVmUsedStorageInMB.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].UsedStorageInMB))
		g.Collect(ch)
		g = e.ZertoVmStatus.WithLabelValues(orgName, vmName, vpgName)
		g.Set(float64(vms[i].Status))
		g.Collect(ch)
	}

	tasks := zertoApi.ListTasks()
	re := regexp.MustCompile("[0-9]+")

	for i:=0;i<len(tasks);i++ {
		var (
			taskId string = tasks[i].TaskIdentifier
			taskType string = tasks[i].Type
		)

		startTime, _ := strconv.ParseFloat(re.FindString(tasks[i].Started), 64)
		completeTime, _ := strconv.ParseFloat(re.FindString(tasks[i].Completed), 64)

		g := e.ZertoTaskStatus.WithLabelValues(taskId, taskType)
		g.Set(float64(tasks[i].Status.State))
		g.Collect(ch)
		g = e.ZertoTaskProgress.WithLabelValues(taskId, taskType)
		g.Set(float64(tasks[i].Status.Progress))
		g.Collect(ch)
		g = e.ZertoTaskStarted.WithLabelValues(taskId, taskType)
		g.Set(startTime)
		g.Collect(ch)
		g = e.ZertoTaskCompleted.WithLabelValues(taskId, taskType)
		g.Set(completeTime)
		g.Collect(ch)
	}

	alerts := zertoApi.ListAlerts()

	{
		var cntErrors, cntWarnings, cntUnknown int = 0, 0, 0
		for i:=0;i<len(alerts);i++ {
			if alerts[i].Level == "Warning" {
				cntWarnings++
			} else if alerts[i].Level == "Error" {
				cntErrors++
			} else {
				cntUnknown++
			}
		}

		g := e.ZertoAlertsCount.WithLabelValues("warning")
		g.Set(float64(cntWarnings))
		g.Collect(ch)
		g = e.ZertoAlertsCount.WithLabelValues("error")
		g.Set(float64(cntErrors))
		g.Collect(ch)
		g = e.ZertoAlertsCount.WithLabelValues("unknown")
		g.Set(float64(cntUnknown))
		g.Collect(ch)
	}
}

func closeZertoSession() {
	log.Debug("Close Session")

	zertoApi.CloseSession()
	zertoSessionAge = 0
}

func openZertoSession() {
	var sessionAge int
	sessionAge = int(time.Now().Unix() - zertoSessionAge)
	log.Debugf("Session Age: %ds", sessionAge)

	if sessionAge > *zertoMaxSessionAge {
		log.Debug("Refresh Session")
		closeZertoSession()
	}

	if !zertoApi.IsSessionOpen() {
		log.Debug("Create new Session")
		zertoApi.OpenSession()
		zertoSessionAge = time.Now().Unix()
	}
}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	log.Debug("Create Zerto instance")
	zertoApi = zerto.NewZerto(*zertoUrl, *zertoUser, *zertoPassword)

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>NameNode Exporter</title></head>
		<body>
		<h1>NameNode Exporter</h1>
		<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>`))
	})

	log.Printf("Starting Server: %s", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
