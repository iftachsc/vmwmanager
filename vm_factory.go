package main

import (
	"fmt"
	"reflect"

	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func makeVM(apiVM mo.VirtualMachine) contracts.VsphereVM {

	disks := []contracts.VsphereDisk{}
	nics := []types.VirtualEthernetCard{}

	var deviceList object.VirtualDeviceList = apiVM.Config.Hardware.Device

	scsiControllers := deviceList.SelectByType((*types.VirtualSCSIController)(nil))
	disksList := deviceList.SelectByType((*types.VirtualDisk)(nil))
	nicsList := deviceList.SelectByType((*types.VirtualEthernetCard)(nil))

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
			//tasks on VirtualDisk that depends on owning scsi controller here
			vd.SetScsiLocation(scsiController.BusNumber, *disk.UnitNumber)

			disks = append(disks, vd)
		}
	}

	for _, nicVd := range nicsList {

		nic := nicVd.(types.BaseVirtualEthernetCard).GetVirtualEthernetCard()
		fmt.Printf("nic type -->", reflect.TypeOf(nicVd).String())

		//vd := makeVsphereDisk(disk)
		//vd.SetScsiLocation(scsiController.BusNumber, *disk.UnitNumber)

		nics = append(nics, *nic)
	}

	return contracts.VsphereVM{apiVM.Reference().Value,
		apiVM.Config.Name,
		apiVM.Guest.HostName,
		apiVM.Guest.GuestFullName,
		disks,
		nics}
}
