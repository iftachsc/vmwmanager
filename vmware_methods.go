package main

import (
	"context"

	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/iftachsc/vmware"
)

type VmwareVcenter struct {
	Host     string
	User     string
	Password string
	Client   *govmomi.Client
}

type OpenNebola struct {
	Host     string
	User     string
	Password string
}

func (vc VmwareVcenter) GetUser() string {
	return vc.User
}
func (vc VmwareVcenter) GetPassword() string {
	return vc.Password
}

func (vc *VmwareVcenter) InitializeClient(ctx context.Context) error {
	if vc.Client != nil {
		return nil
	}

	println("Initiazlizing vim client for ", vc.Host)
	client, err := vmware.NewClient(ctx, vc.Host, vc.User, vc.Password)
	if err != nil {
		vc.Client = nil
		return err
	}

	vc.Client = client
	return nil
}

//GetVMs returns list of VsphereVM objects made out of mo.VirtualMachine api objects
func (vc *VmwareVcenter) GetVMs(ctx context.Context) (vms []contracts.VsphereVM, err error) {
	movms, error := vmware.GetVM(vc.Client, ctx)

	if err != nil {
		return nil, error
	}
	//translate mo.VirtualMachine to our VsphereVM
	for _, movm := range movms {
		vms = append(vms, makeVM(movm))
	}
	return vms, nil
}

//GetHost returns list of EsxHost objects made of mo.HostSystem api objects
func (vc *VmwareVcenter) GetHost(ctx context.Context) ([]contracts.EsxHost, error) {
	hostsystems, err := vmware.GetEsxHost(vc.Client, ctx)
	hosts := []contracts.EsxHost{}
	if err != nil {
		return nil, err
	}
	for _, hs := range hostsystems {
		hosts = append(hosts, makeEsxHost(hs))
	}
	return hosts, nil
}

//GetScsiLunDisks returns list unique list of contracts.Lun objects made of ScsiLun api objects
func (vc *VmwareVcenter) GetScsiLunDisks(ctx context.Context) (luns []types.ScsiLun, err error) {
	luns, error := vmware.GetScsiLunDisks(ctx, vc.Client)

	if err != nil {
		return nil, error
	}
	// for _, lun := range luns {
	// 	luns = append(luns, makeLun(hs))
	// }
	return luns, nil
}
