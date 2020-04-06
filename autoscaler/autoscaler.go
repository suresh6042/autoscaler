package autoscaler

import (
	"context"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"log"
)

type ScaleConf struct {
	VimClient *vim25.Client
	ScaleConfAttr
}
type ScaleConfAttr struct {
	Name   string `yaml:"name"`
	CpuAlarm string `yaml:"cpualarm"`
	MemoryAlarm string `yaml:"memoryalarm"`
	CpuToScale int32 `yaml:"cputoscale"`
	MemoryToScale int64 `yaml:"memorytoscale"`
	MemoryThreshold int64 `yaml:"memorythreshold"`
	CpuThreshold int32 `yaml:"cputhreshold"`
}
func Errorhandle(err error) {
	if err != nil {
		panic(err)
	}
}
func (scaleconf ScaleConf ) Scale(ctx context.Context) {
	finder:=find.NewFinder(scaleconf.VimClient)
	vm,err:=finder.VirtualMachine(ctx,scaleconf.Name)
	Errorhandle(err)
	var objv mo.VirtualMachine
	err=vm.Properties(ctx,vm.Reference(),[]string{"summary.config.numCpu","summary.config.memorySizeMB"},&objv)
	Errorhandle(err)
	currentCpu:=objv.Summary.Config.NumCpu
	currentMem:=objv.Summary.Config.MemorySizeMB
	var objm mo.ManagedEntity
	err=vm.Properties(ctx,vm.Reference(),[]string{"triggeredAlarmState"},&objm)
	Errorhandle(err)
	alarmlist:=objm.TriggeredAlarmState
	cpuflag:=0
	memflag:=0
	for _,alarm:=range alarmlist{
		if cpuflag=1;alarm.Alarm.Value == scaleconf.CpuAlarm && alarm.OverallStatus == "red"{
			desiredCpu:=currentCpu+scaleconf.CpuToScale
			if desiredCpu < scaleconf.CpuThreshold {
				spec:=types.VirtualMachineConfigSpec{NumCPUs:desiredCpu}
				task, err := vm.Reconfigure(ctx, spec)
				Errorhandle(err)
        		err = task.Wait(ctx)
        		Errorhandle(err)
        		log.Printf("CPU scaled up to %d cores for VM %s successfully \n",desiredCpu,scaleconf.Name)
			} else {
				log.Printf("CPU alarm found for %s but unable to scale up as the threshold reached\n",scaleconf.Name)
			}
		}else if memflag=1; alarm.Alarm.Value == scaleconf.MemoryAlarm && alarm.OverallStatus == "red"{
			desiredMem:=int64(currentMem)+scaleconf.MemoryToScale
			if desiredMem < scaleconf.MemoryThreshold{
				spec:=types.VirtualMachineConfigSpec{MemoryMB:desiredMem}
				task, err := vm.Reconfigure(ctx, spec)
				Errorhandle(err)
        		err = task.Wait(ctx)
        		Errorhandle(err)
        		log.Printf("Memory scaled up to %d MB for VM %s successfully \n",desiredMem,scaleconf.Name)
			} else {
				log.Printf("Memory alarm triggered for %s but unable to scale up as the threshold reached\n",scaleconf.Name)
			}
		}
	}
	if cpuflag == 0{
		log.Printf("No CPU Alarm triggered for %s \n" ,scaleconf.Name)
	}
	if memflag == 0{
		log.Printf("No Memory Alarm triggered for %s \n" ,scaleconf.Name)
	}
}
