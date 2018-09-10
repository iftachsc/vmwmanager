package main

import (
	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi/vim25/mo"
)

func makeEsxHost(apiHost mo.HostSystem) contracts.EsxHost {
	return contracts.EsxHost{apiHost.Summary.Host.Value,
		apiHost.Summary.Config.Name,
		apiHost.Summary.Hardware.Model}
}

func apiHostsToContractEsxHosts(apiHosts []mo.HostSystem) []contracts.EsxHost {

	hosts := []contracts.EsxHost{}

	for _, host := range apiHosts {
		hosts = append(hosts, makeEsxHost(host))
		//hss := object.NewHostStorageSystem(c.Client, *host.ConfigManager.StorageSystem)
		//hss.Scs
		//scsiLuns := GetScsiLuns(host.ConfigManager.StorageSystem)
	}

	return hosts
}
