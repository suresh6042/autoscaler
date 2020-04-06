package main
import (
        "context"
        "flag"
        "fmt"
        "github.com/vmware/govmomi"
        "github.com/vmware/govmomi/vim25/mo"
        "github.com/vmware/govmomi/vim25/types"
        "os"
        "net/url"
        "github.com/vmware/govmomi/find"
)

func errorhandle(err error){
                if err != nil {
                fmt.Println(err.Error())
                os.Exit(1)
        }
}
func main(){
        delta_cpu:=flag.Int("cpu",0,"Number of CPU Cores to be increased")
        delta_mem:=flag.Int("memory",0,"Amount of memory to be increased")
        vm_name:=flag.String("vm","","VM to be scaled up")
        flag.Parse()
        if *vm_name == ""{
                fmt.Println("Error: Flag -vm required")
                os.Exit(1)
        }
        _,vcs_endpoint_set:=os.LookupEnv("VCS_ENDPOINT")
        if !vcs_endpoint_set{
                fmt.Println("Error: Environment variable VCS_ENDPOINT not set")
                os.Exit(1)
        }
        vcs_raw_url:=os.Getenv("VCS_ENDPOINT")
        vcs_url,err:=url.Parse(vcs_raw_url)
        ctx, cancel:=context.WithCancel(context.Background())
        defer cancel()
        client, err := govmomi.NewClient(ctx, vcs_url, true)
        errorhandle(err)
        finder:=find.NewFinder(client.Client)
        vm,err:=finder.VirtualMachine(ctx,*vm_name)
        errorhandle(err)
        var obj mo.VirtualMachine
        err=vm.Properties(ctx,vm.Reference(),[]string{"summary.config.numCpu","summary.config.memorySizeMB"},&obj)
        errorhandle(err)
        current_cpu:=obj.Summary.Config.NumCpu
        current_mem:=obj.Summary.Config.MemorySizeMB
        desired_cpu:=current_cpu+int32(*delta_cpu)
        desired_mem:=int64(current_mem)+int64(*delta_mem)
        spec:=types.VirtualMachineConfigSpec{NumCPUs:desired_cpu,MemoryMB:desired_mem}
        task, err := vm.Reconfigure(ctx, spec)
        errorhandle(err)
        err = task.Wait(ctx)
        errorhandle(err)
        fmt.Println("Resources scaled up successfully")
        fmt.Printf("Current CPU: %d \n",desired_cpu)
        fmt.Printf("Current MemoryMB: %d \n",desired_mem)
}
