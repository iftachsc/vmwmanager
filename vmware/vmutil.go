package vmware

import (
	"context"
	"log"

	"github.com/iftachsc/contracts"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

type VsphereVM contracts.VsphereVM

func makeVM(apiVM mo.VirtualMachine) VsphereVM {
	return VsphereVM{apiVM.Summary.Vm.Value,
		apiVM.Summary.Config.Name,
		apiVM.Summary.Guest.HostName,
		apiVM.Summary.Config.GuestFullName}
}

func GetVM(c *govmomi.Client, ctx context.Context) ([]VsphereVM, error) {

	vms := []VsphereVM{}

	//defer c.Logout(ctx)

	// Create view of VirtualMachine objects
	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all machines
	// Reference: http://pubs.vmware.com/vsphere-60/topic/com.vmware.wssdk.apiref.doc/vim.VirtualMachine.html
	var virtualmachines []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &virtualmachines)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, vm := range virtualmachines {
		vms = append(vms, makeVM(vm))
	}

	return vms, nil
}

/*
//should go to scenario
func RegisterVM() {
	vm := new(mo.HostStorageSystem)

	canBePoweredOn := sUnmirroredVmdkVmNames.Contains(selectedMachineData.Name)
	//AddToLog(string.Format("VM {0} has a VMDK that is not mirrored.", selectedMachineData.Name), "WARNING");

	//newMachineName = PrefixTB.Text + selectedMachineData.Name;
	//newVmxPath = MyVMWare.GetVmxPath(selectedMachineData.VMXPath, PrefixTB.Text);

	//HostSystem esx = null;
	//var replicatedRdms = WizardData.RdmVolumeSnapshotMappingCmode.Keys.Where(x => x.rdmVmdk.vmName == newMachineName.Substring(prefix.Length));

}


func RegisterVM(c *govmomi.Client, ctx context.Context, canBePoweredOn bool, rdmDatastore Datastore) (VsphereVM, error) {
	newVmxPath := ""

	types.ScsiLun
	mo.  .ScsiLun targetScsiLun = null;
           // try
           // {

                //Round-Robin
                if (replicatedRdms.Count() > 0)
                {
                    //getting vmhost that sees the rdm datastore in round robing
                    var hostMount = sVMwareClient.rdmDatastore.Host[nextEsxToRegisterCounter % sVMwareClient.rdmDatastore.Host.Length];
                    esx = sVMwareClient.Esxs.Where(vmh => vmh.MoRef.Value == hostMount.Key.Value).First();
                }
                else
                {
                    esx = sVMwareClient.Esxs[nextEsxToRegisterCounter % sVMwareClient.Esxs.Count];
                }
                VirtualMachine newVm = RegisterOnEsx(conn, selectedMachineData, esx, workingFolder, newMachineName, selectedMachineData.VMXPath);
                nextEsxToRegisterCounter++;

                if (newVm != null)
                {
                    try
                    {
                        MyVMWare.SetVmNetwork_V3(sVMwareClient.VimClient, newVm.MoRef, selectedMachineData.Adapters);

                        //fixing flat vmdks with full url filepath

                        var diskDevices = newVm.Config.Hardware.Device.Where(dev => dev is VirtualDisk);
                        var scsiControllers = newVm.Config.Hardware.Device.Where(dev => dev is VirtualSCSIController);

                        foreach (VirtualDisk dev in diskDevices.Where(x => !(x.Backing is VirtualDiskRawDiskVer2BackingInfo && x.Backing is VirtualDiskRawDiskMappingVer1BackingInfo)))
                        {
                            var backing = (VirtualDeviceFileBackingInfo)dev.Backing;

                            var vmdkBackingFileName = ((VirtualDeviceFileBackingInfo)dev.Backing).FileName;
                            if (vmdkBackingFileName.StartsWith("[] /vmfs/volumes/")) //filename with uuid form
                            {
                                var originDatastoreByUuidFileName = datastoreByVmdkUuidFileName(vmdkBackingFileName);
                                if (originDatastoreByUuidFileName != null)
                                {
                                    MyVMWare.fixVmdkFileName(sVMwareClient.VimClient, newVm, dev, vmdkBackingFileName, originDatastoreByUuidFileName.Name);
                                    AddToLog(string.Format("Fixed vmdk file path on vm {0}. file path {1} fixed to {2}", newVm.Name, originDatastoreByUuidFileName.Name),"INFO", "Registering Virtual Mahcines");
                                }
                                else
                                    MyVMWare.removeDisk(sVMwareClient.VimClient, newVm, dev);
                                    AddToLog(string.Format("Removing previous RDM disk. label {0}, vm {1}.",dev.DeviceInfo.Label, newVm.Name),"INFO", "Registering Virtual Mahcines");
                            }
                            else if (selectedMachineData.VMDKs.Find(vmdk => vmdk.Name == dev.DeviceInfo.Label && vmdk.Type == "RawDiskMapping") != null)
                            {
                                MyVMWare.removeDisk(sVMwareClient.VimClient, newVm, dev);
                            }
                            if ((newVm.Name.Substring(prefix.Length) == "Ghibli" || (newVm.Name.Substring(prefix.Length) == "exchange_2007")) && dev.DeviceInfo.Label != "Hard disk 1")
                            {
                                MyVMWare.removeDisk(sVMwareClient.VimClient, newVm, dev);
                            }
                        }

                        //fix rdms
                        if (sVMwareClient.Esxs.Count == 0)
                        {
                            AddToLog("Retreiving Esx Hosts for cluster " + sVMwareClient.ConnectedCluster.Name);
                            sVMwareClient.VimClient = MyVMWare.ConnectVC(sVMwareClient.ServerName, sVMwareClient.UserName, sVMwareClient.Password);
                            sVMwareClient.Esxs = VimExtensionFunctions.GetHosts(sVMwareClient.VimClient, sVMwareClient.ConnectedCluster.MoRef);
                        }

                        //Dictionary<ScsiLun, LunRuneTime> scsiLunsWithRuntime = null;
                        List<ScsiLun> scsiLunsWithRuntime = null;
                        HostSystem runningEsx = sVMwareClient.Esxs.Find(x => x.MoRef.Value == newVm.Summary.Runtime.Host.Value);
                        var storagesystem = new HostStorageSystem(sVMwareClient.VimClient, runningEsx.ConfigManager.StorageSystem);
                        storagesystem.UpdateViewData();

                        if (MainWindow.NCConnections.Count > 0)
                        {

                            //if (replicatedRdms.Count() > 0) scsiLunsWithRuntime = getScsiLunsWithRuntime(newVm.Summary.Runtime.Host);
                            if (replicatedRdms.Count() > 0) scsiLunsWithRuntime = getEsxScsiDiskLuns(storagesystem);
                            foreach (RdmSnapMirrorMappingCmode rdmMap in replicatedRdms)
                            {
                                newVm.UpdateViewData();
                                scsiControllers = newVm.Config.Hardware.Device.Where(dev => dev is VirtualSCSIController);
                                VMDKData rdmVmdk = rdmMap.rdmVmdk;

                                try
                                {
                                    //sanId is iscsi node name of the vserver or fcp wwn
                                    targetScsiLun = getTargetScsiLun(scsiLunsWithRuntime, storagesystem, rdmMap.lunMap.lunId, rdmMap.lunMap.sanId);
                                }
                                catch
                                {
                                    AddToLog(string.Format("Could not find target scsi lun for VM {0} ({1}) on vmhost {2} (LunId {3})", newVm.Name, rdmMap.rdmVmdk.Name, runningEsx.Summary.Config.Name, rdmMap.lunMap.lunId), "ERROR");
                                }
                                //var targetScsiLun = scsiLunsWithRuntime.Where(x => x.Value.lunId == rdmMap.lunMap.lunId).FirstOrDefault();
                                AddToLog(string.Format("Adding RDM {0}: Disk Device {1}", rdmMap.rdmVmdk.Name, targetScsiLun.CanonicalName), "INFO", string.Format("Handling RDMs on VM {0}", newVm.Name), null,  string.Format("Adding RDM {0} with Disk Device {1}", rdmMap.rdmVmdk.Name, targetScsiLun.CanonicalName));
                                VirtualDiskRawDiskMappingVer1BackingInfo backing = new VirtualDiskRawDiskMappingVer1BackingInfo();
                                backing.CompatibilityMode = VirtualDiskCompatibilityMode.physicalMode.ToString();
                                backing.DiskMode = "persistent";
                                backing.DeviceName = targetScsiLun.DeviceName;
                                var originalRdmFilename = rdmMap.rdmVmdk.FilePath.Split('/').Last();

                                //taking the whole path of newVm vmx without the filename it self and replacing with poroiginal rdm file name
                                //to achive placing rdm file in same directory of new registered vm but with original name to keep consistency
                                // and avoid conflics of filenameseasily
                                backing.FileName = string.Format("[{0}] {1}/{2}", sVMwareClient.rdmDatastore.Name, sVMwareClient.WorkingFolderName, originalRdmFilename);
                                var diskSpec = new VirtualDeviceConfigSpec();
                                diskSpec.Device = new VirtualDisk();
                                diskSpec.Device.Backing = backing;

                                diskSpec.Device.Key = -202;
                                diskSpec.FileOperation = VirtualDeviceConfigSpecFileOperation.create;
                                diskSpec.Operation = VirtualDeviceConfigSpecOperation.add;
                                int[] arrScsiLocation = rdmMap.rdmVmdk.SCSILocation.Split(':').Select(x => int.Parse(x)).ToArray();
                                diskSpec.Device.UnitNumber = arrScsiLocation[1];
                                var spec = new VirtualMachineConfigSpec();
                                if (arrScsiLocation[0] == 0 || (arrScsiLocation[0] == 1 && scsiControllers.Count() > 1))
                                {
                                    diskSpec.Device.ControllerKey = scsiControllers.Where(x => (x as VirtualSCSIController).BusNumber == arrScsiLocation[0]).First().Key;
                                    spec.DeviceChange = new VirtualDeviceConfigSpec[1];
                                    spec.DeviceChange[0] = diskSpec;
                                }
                                else
                                {
                                    var scsiDev = new VirtualDeviceConfigSpec();
                                    scsiDev.Operation = VirtualDeviceConfigSpecOperation.add;
                                    scsiDev.Device = new ParaVirtualSCSIController();
                                    scsiDev.Device.Key = -222;
                                    (scsiDev.Device as ParaVirtualSCSIController).BusNumber = arrScsiLocation[0];
                                    (scsiDev.Device as ParaVirtualSCSIController).SharedBus = VirtualSCSISharing.physicalSharing;

                                    scsiDev.Device.UnitNumber = 0;
                                    spec.DeviceChange = new VirtualDeviceConfigSpec[2];
                                    diskSpec.Device.ControllerKey = -222;
                                    spec.DeviceChange[0] = scsiDev;
                                    spec.DeviceChange[1] = diskSpec;
                                }

                                sVMwareClient.VimClient.WaitForTask(newVm.ReconfigVM_Task(spec));
                            }
                        }

                        // Erez
                        //if (MainWindow.NAConnections.Count > 0)
                        //{
                        //    if (replicatedRdms.Count() > 0) scsiLunsWithRuntime = getScsiLunsWithRuntime(newVm.Summary.Runtime.Host);
                        //    foreach (RdmSnapMirrorMapping sm in replicatedRdms)
                        //    {
                        //        newVm.UpdateViewData();
                        //        scsiControllers = newVm.Config.Hardware.Device.Where(dev => dev is VirtualSCSIController);
                        //        VMDKData rdmVmdk = sm.OriginObject as VMDKData;
                        //        //var targetReplicatedRdm = replicatedRdms.Where(x => x.rdmVmdk.Name == dev.DeviceInfo.Label);
                        //        var targetScsiLun = scsiLunsWithRuntime.Where(x => x.Value.lunId == sm.lunMap.lunId).FirstOrDefault();
                        //        AddToLog(string.Format("Adding RDM {0}: DevName {1} , LunId {2}, Target: {3}", rdmVmdk.Name, targetScsiLun.Key.DeviceName, targetScsiLun.Value.lunId, targetScsiLun.Value.target), "INFO", string.Format("Handling RDMs on VM {0}", newVm.Name), string.Format("Adding RDM {0}: LunId {1}, Target: {2}", rdmVmdk.Name, targetScsiLun.Value.lunId, targetScsiLun.Value.target));
                        //        VirtualDiskRawDiskMappingVer1BackingInfo backing = new VirtualDiskRawDiskMappingVer1BackingInfo();
                        //        backing.CompatibilityMode = VirtualDiskCompatibilityMode.physicalMode.ToString();
                        //        backing.DiskMode = "persistent";
                        //        backing.DeviceName = targetScsiLun.Key.DeviceName;
                        //        var originalRdmFilename = rdmVmdk.FilePath.Split('/').Last();
                        //        //taking the whole path of newVm vmx without the filename it self and replacing with poroiginal rdm file name
                        //        //to achive placing rdm file in same directory of new registered vm but with original name to keep consistency
                        //        // and avoid conflics of filenameseasily
                        //        backing.FileName = string.Format("[{0}] {1}/{2}", sVMwareClient.rdmDatastore.Name, sVMwareClient.WorkingFolderName, originalRdmFilename);

                        //        var diskSpec = new VirtualDeviceConfigSpec();
                        //        diskSpec.Device = new VirtualDisk();
                        //        diskSpec.Device.Key = -202;

                        //        diskSpec.FileOperation = VirtualDeviceConfigSpecFileOperation.create;
                        //        diskSpec.Operation = VirtualDeviceConfigSpecOperation.add;
                        //        //diskSpec.FileOperation = VirtualDeviceConfigSpecFileOperation.replace;
                        //        diskSpec.Device.Backing = backing;

                        //        int[] arrScsiLocation = rdmVmdk.SCSILocation.Split(':').Select(x => int.Parse(x)).ToArray();
                        //        diskSpec.Device.UnitNumber = arrScsiLocation[1];

                        //        var spec = new VirtualMachineConfigSpec();
                        //        if (arrScsiLocation[0] == 0 || (arrScsiLocation[0] == 1 && scsiControllers.Count() > 1))
                        //        {
                        //            diskSpec.Device.ControllerKey = scsiControllers.Where(x => (x as VirtualSCSIController).BusNumber == arrScsiLocation[0]).First().Key;
                        //            spec.DeviceChange = new VirtualDeviceConfigSpec[1];
                        //            spec.DeviceChange[0] = diskSpec;
                        //        }
                        //        else
                        //        {
                        //            //diskSpec.Device.ControllerKey = scsiControllers.First().Key;
                        //            var scsiDev = new VirtualDeviceConfigSpec();
                        //            scsiDev.Operation = VirtualDeviceConfigSpecOperation.add;
                        //            scsiDev.Device = new ParaVirtualSCSIController();
                        //            scsiDev.Device.Key = -222;
                        //            (scsiDev.Device as ParaVirtualSCSIController).BusNumber = arrScsiLocation[0];
                        //            (scsiDev.Device as ParaVirtualSCSIController).SharedBus = VirtualSCSISharing.physicalSharing;

                        //            scsiDev.Device.UnitNumber = 0;
                        //            spec.DeviceChange = new VirtualDeviceConfigSpec[2];
                        //            diskSpec.Device.ControllerKey = -222;
                        //            spec.DeviceChange[0] = scsiDev;
                        //            spec.DeviceChange[1] = diskSpec;
                        //        }


                        //        sVMwareClient.VimClient.WaitForTask(newVm.ReconfigVM_Task(spec));
                        //    }
                        //}

                        if (selectedMachineData.PostPowerOnSleep.HasValue && canBePoweredOn)
                        {
                            AddToLog(string.Format("powering on machine {0}", newVm.Name));
                            newVm.PowerOnVM_Task(null);

                            if (selectedMachineData.PostPowerOnSleep.Value >= 0)
                            {
                                AddToLog("Sleeping for " + selectedMachineData.PostPowerOnSleep.Value + " Minutes");
                                Thread.Sleep((selectedMachineData.PostPowerOnSleep.Value * 60  + 10)* 1000);
                                try
                                {
                                    newVm.UpdateViewData("runtime.question");
                                    if(newVm.Runtime.Question != null && newVm.Runtime.Question.Id == "_vmx1")
                                        newVm.AnswerVM("_vmx1", "2");
                                }
                                catch (Exception ex)
                                {
                                    AddToLog("failed to answer power on question, error : " + ex.Message, "ERROR", null, ex);
                                }
                            }
                        }

                        return newVm;
                    }
                    catch (Exception ex)
                    {
                        AddToLog(string.Format("Failed post-register process [VM Network, Power], error : {0}", ex.Message));
                    }
                }
            }
            catch (Exception ex)
            {
                AddToLog(string.Format("Could not start Registeration on vc {0} due to the following error : {1}", conn.ServerName, ex.Message));
            }

            return null;
        }
}
*/
