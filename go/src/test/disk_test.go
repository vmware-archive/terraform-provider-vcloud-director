/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package test

import (
	"fmt"
	//	"io/ioutil"
	//"log"
	"testing"
	"os"
	"os/exec"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	//"runtime/debug"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
)
//go test disk_test.go -v run  TestDiskInterface
func TestDiskInterface(t *testing.T){
	// We don't want to see the plugin logs.
	//log.SetOutput(ioutil.Discard)
	//log.SetOutput(os.Stdout)
	logging.Init()

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: grpc.Handshake,
		Plugins:         grpc.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("PY_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	//defer client.Kill()
	
	logging.Plog(" ok starting to client  ")
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {

		fmt.Println("==============\n\n\nError:", err.Error())
		os.Exit(1)
	}
	logging.Plog("call dispense")
	// Request the plugin
	raw, err := rpcClient.Dispense("DISK_PLUGIN")
	if err != nil {
		fmt.Println("Error:", err.Error())

	}


	lraw, _ := rpcClient.Dispense("PY_PLUGIN")
	

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	lprovider := lraw.(grpc.PyVcloudProvider)

	lc := proto.LoginCredentials{
		Username: "user1",
		Password: "Admin!23",
		Org:      "O1",
		Ip:       "10.112.83.27",

		//UseVcdCliProfile: true,

		AllowInsecureFlag: true,
	}
	_, err = lprovider.Login(lc)


	
	provider := raw.(grpc.IndependentDiskProvider)
	/*
	res,err:=provider.Create(proto.CreateDiskInfo{Name:"one",Vdc:"OVD4",Size:"100"});

	if(err!=nil){
		//logging.Plog(fmt.Errorf("Error",err))
		t.Errorf("err %+v",err)
		return 
	}
	logging.Plogf("ok created id= %v",res.DiskId)
	*/

	disk,_:=provider.Read(proto.ReadDiskInfo{Name:"disk2",Vdc:"OVD4",DiskId:"-4"});
	logging.Plog(fmt.Sprintf(" %+v",disk.Present))


	//provider.Delete(proto.DeleteDiskInfo{Name:"one",Vdc:"OVD4",DiskId:res.DiskId});
	/*
	for {
		provider.Delete(proto.DeleteDiskInfo{Name:"one",Vdc:"OVD4"});
	}
	*/


	
}
