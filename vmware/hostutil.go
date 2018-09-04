package vmware

import (
	"context"
	"log"

	"github.com/iftachsc/contracts"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

type EsxHost contracts.EsxHost

func makeEsxHost(apiHost mo.HostSystem) EsxHost {
	return EsxHost{apiHost.Summary.Host.Value,
		apiHost.Summary.Config.Name,
		apiHost.Summary.Hardware.Model}
}

// func (vm *VM) AddDisk(item VM) []VMItem {
//     vm.Items = append(vm.Items, item)
//     return vm.Items
// }

func GetEsxHost(c *govmomi.Client, ctx context.Context) ([]EsxHost, error) {

	hosts := []EsxHost{}
	//defer m.Destroy(ctx)

	//defer c.Logout(ctx)

	// Create view of VirtualMachine objects
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var hostsystems []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hostsystems)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// Print summary per vm (see also: govc/vm/info.go)
	//for _, vm := range virtualmachines {
	//	fmt.Printf("%s: %s\n", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
	//}

	for _, host := range hostsystems {
		hosts = append(hosts, makeEsxHost(host))
	}

	return hosts, nil
}
