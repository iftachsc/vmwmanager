package main

import (
	"reflect"

	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func makeVM(apiVM mo.VirtualMachine) contracts.VsphereVM {

	disks := []contracts.VsphereDisk{}
	var deviceList object.VirtualDeviceList = apiVM.Config.Hardware.Device

	scsiControllers := deviceList.SelectByType((*types.VirtualSCSIController)(nil))
	disksList := deviceList.SelectByType((*types.VirtualDisk)(nil))

	for _, controller := range scsiControllers {

		scsiController := controller.(types.BaseVirtualSCSIController).GetVirtualSCSIController()

		//fmt.Printf(reflect.TypeOf(controller).String())
		for _, controllerDisk := range disksList.Select(func(disk types.BaseVirtualDevice) bool {
			return disk.GetVirtualDevice().ControllerKey == scsiController.Key
		}) {

			disk, ok := controllerDisk.(*types.VirtualDisk)
			if !ok {
				println("not ok:", reflect.TypeOf(disk).String())
			}

			vd := makeVsphereDisk(disk)
			vd.SetScsiLocation(scsiController.BusNumber, *disk.UnitNumber)

			disks = append(disks, vd)
		}
	}

	return contracts.VsphereVM{apiVM.Reference().Value,
		apiVM.Config.Name,
		apiVM.Guest.HostName,
		apiVM.Guest.GuestFullName,
		disks}
}

func apiVmsToContractVms(apiVms []mo.VirtualMachine) []contracts.VsphereVM {

	vms := []contracts.VsphereVM{}

	for _, vm := range apiVms {
		vms = append(vms, makeVM(vm))
	}

	return vms
}
