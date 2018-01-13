/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package test

import (
	"fmt"
	//	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"runtime/debug"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
)

func main() {
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
	res := proto.ReadVAppResult{Name: "somename"}
	logging.Plogf(" ok starting to client %#v", res)
	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {

		fmt.Println("==============\n\n\nError:", err.Error())
		os.Exit(1)
	}
	logging.Plog("call dispense")
	// Request the plugin
	raw, err := rpcClient.Dispense("PY_PLUGIN")
	if err != nil {
		fmt.Println("Error:", err.Error())

	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	provider := raw.(grpc.PyVcloudProvider)

	lc := proto.LoginCredentials{
		Username: "user1",
		Password: "Admin!23",
		Org:      "O1",
		Ip:       "10.112.83.27",

		//UseVcdCliProfile: true,

		AllowInsecureFlag: true,
	}
	_, err = provider.Login(lc)
	if err != nil {
		debug.PrintStack()

		fmt.Errorf("error in login %s", err)
	}

	vappInfo := proto.ReadVAppInfo{
		Name: "cento_vapp11_2",
		Vdc:  "OVD2",
	}

	resv, errv := provider.ReadVApp(vappInfo)
	
	if errv != nil {
		logging.PlogErrorf("error occ = %#v",errv)

	}else{
	logging.Plogf("%#v", *resv)
	}


	captureInfo:=proto.CaptureVAppInfo{
		CatalogName: "c1",
		ItemName: "item21",
		VappName:"cento_vapp11_2",
		VdcName:"OVD2",
	}

	resp, errp := provider.CaptureVapp(captureInfo)
	if errp != nil {
		
		log.Print(fmt.Errorf("__ERROR__ Creating catalog Item: [%#v]", errp))
		return 
	}
	logging.Plogf("%+v",resp)

}
