package main

import (
	"context"
	"flag"
	"github.com/vmware/govmomi"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"time"
	"vmware/autoscaler"
)
func main(){
	err:=os.Setenv("VCS_ENDPOINT","https://dstadmin@vsphere.local:Adm!n4DST@vcenter.qanon-prod.smdibm.com/sdk")
	_,vcs_endpoint_set:=os.LookupEnv("VCS_ENDPOINT")
	if !vcs_endpoint_set{
		panic("Environment variable VCS_ENDPOINT not set")
	}
	config:=flag.String("config",`C:\Users\SureshR\go\src\vmware\config.yaml`,"Config file path")
	if *config == "" {
		panic("Config file path not defined")
	}
	interval:=flag.String("interval","30s","Interval to monitor VCS alarms")
	flag.Parse()
	duration,err:=time.ParseDuration(*interval)
	autoscaler.Errorhandle(err)
	vcs_raw_url:=os.Getenv("VCS_ENDPOINT")
	vcs_url,err:=url.Parse(vcs_raw_url)
	ctx, cancel:=context.WithCancel(context.Background())
	defer cancel()
	client, err := govmomi.NewClient(ctx, vcs_url, true)
	autoscaler.Errorhandle(err)
    f,err:=os.Open(*config)
    autoscaler.Errorhandle(err)
    defer f.Close()
    var scaleconfattrsl []autoscaler.ScaleConfAttr
    dec := yaml.NewDecoder(f)
    err = dec.Decode(&scaleconfattrsl)
    autoscaler.Errorhandle(err)
    for{
    	time.Sleep(duration)
    	for _,scaleconfattr:= range scaleconfattrsl{
    		scaleconf:= autoscaler.ScaleConf{
			VimClient:     client.Client,
			ScaleConfAttr: scaleconfattr,
		}
		go scaleconf.Scale(ctx)
	}
	}
}

