package main

// func makeVsphereVmNic(disk *types.VirtualEthernetCard) contracts.VsphereDisk {
// 	var vsphereDisk contracts.VsphereDisk

// 	switch backing := disk.GetVirtualDevice().Backing.(type) {
// 	case *types.VirtualDiskFlatVer1BackingInfo:
// 		println("not implemented")
// 	case *types.VirtualDiskFlatVer2BackingInfo:
// 		vsphereDisk = makeFlatVsphereDiskV2(backing)

// 	case *types.VirtualDiskRawDiskMappingVer1BackingInfo:
// 		vsphereDisk = makeRdmVsphereDiskV1(backing)

// 	case *types.VirtualDiskRawDiskVer2BackingInfo:
// 		vsphereDisk = makeRdmVsphereDiskV2(backing)
// 	default:
// 		println("cannot determine type")
// 	}

// 	vsphereDisk.SetCapacityMB(disk.CapacityInKB / 1024)

// 	return vsphereDisk
// }
