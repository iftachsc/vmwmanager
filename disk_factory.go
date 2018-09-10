package main

import (
	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi/vim25/types"
)

func makeRdmMetadataV1(backing types.VirtualDiskRawDiskMappingVer1BackingInfo) *contracts.VsphereRdmDisk {
	return &contracts.VsphereRdmDisk{
		DescriptorFileName:  backing.FileName,
		DescriptorDatastore: backing.Datastore.String(),
		LunUuid:             backing.LunUuid,
		StorageServer:       "",
	}
}

func makeRdmVsphereDiskV1(backing *types.VirtualDiskRawDiskMappingVer1BackingInfo) *contracts.VsphereRdmDisk {
	return &contracts.VsphereRdmDisk{
		DescriptorFileName:  backing.FileName,
		DescriptorDatastore: backing.Datastore.String(), //TODO extract datastore name
		LunUuid:             backing.LunUuid,
		StorageServer:       backing.BackingObjectId + "," + backing.ContentId + "," + backing.DeviceName,
	}
}

func makeRdmVsphereDiskV2(backing *types.VirtualDiskRawDiskVer2BackingInfo) *contracts.VsphereRdmDisk {
	return &contracts.VsphereRdmDisk{
		DescriptorFileName:  backing.DescriptorFileName,
		DescriptorDatastore: backing.DescriptorFileName, //TODO extract datastore name
		LunUuid:             backing.Uuid,
		StorageServer:       "",
	}
}

func makeFlatVsphereDiskV2(backing *types.VirtualDiskFlatVer2BackingInfo) *contracts.VsphereFlatDisk {
	return &contracts.VsphereFlatDisk{
		FlatFileName: backing.FileName,
		Datastore:    backing.Datastore.String(), //TODO extract datastore name
		DiskMode:     backing.DiskMode,
	}
}

func makeVsphereDisk(disk *types.VirtualDisk) contracts.VsphereDisk {
	var vsphereDisk contracts.VsphereDisk

	switch backing := disk.GetVirtualDevice().Backing.(type) {
	case *types.VirtualDiskFlatVer1BackingInfo:
		println("not implemented")
	case *types.VirtualDiskFlatVer2BackingInfo:
		vsphereDisk = makeFlatVsphereDiskV2(backing)

	case *types.VirtualDiskRawDiskMappingVer1BackingInfo:
		vsphereDisk = makeRdmVsphereDiskV1(backing)

	case *types.VirtualDiskRawDiskVer2BackingInfo:
		vsphereDisk = makeRdmVsphereDiskV2(backing)
	default:
		println("cannot determine type")
	}

	vsphereDisk.SetCapacityMB(disk.CapacityInKB / 1024)

	return vsphereDisk
}
