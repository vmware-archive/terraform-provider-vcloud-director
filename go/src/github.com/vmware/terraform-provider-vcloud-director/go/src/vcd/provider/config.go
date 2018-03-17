/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package provider

import (
	"fmt"

	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/grpc"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
)

type Config struct {
	L proto.LoginCredentials
}

type VCDClient struct {
	*plugin.Client
	*plugin.GRPCClient
}

func (v VCDClient) getProvider() grpc.PyVcloudProvider {

	// Request the plugin
	raw, err := v.GRPCClient.Dispense("PY_PLUGIN")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	provider := raw.(grpc.PyVcloudProvider)
	return provider

}

func (c Config) CreateClient() (*VCDClient, error) {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: grpc.Handshake,
		Plugins:         grpc.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("PY_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	//defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("PY_PLUGIN")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	provider := raw.(grpc.PyVcloudProvider)

	result, err := provider.Login(c.L)

	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(result.Token))

	//

	rawDiskProvider, errDisk := rpcClient.Dispense("DISK_PLUGIN")
	if errDisk != nil {
		fmt.Println("Error:", errDisk.Error())

	}

	diskProvider := rawDiskProvider.(grpc.IndependentDiskProvider)

	//

	//

	rawOrgProvider, errOrg := rpcClient.Dispense("ORG_PLUGIN")
	if errOrg != nil {
		fmt.Println("Error:", errOrg.Error())

	}

	orgProvider := rawOrgProvider.(grpc.OrgProvider)

	rawUserProvider, errUser := rpcClient.Dispense("USER_PLUGIN")
	if errUser != nil {
		fmt.Println("Error:", errUser.Error())

	}
	userProvider := rawUserProvider.(grpc.UserProvider)

	rawVdcProvider, errVdc := rpcClient.Dispense("VDC_PLUGIN")
	if errUser != nil {
		fmt.Println("Error:", errVdc.Error())

	}
	vdcProvider := rawVdcProvider.(grpc.VdcProvider)

	rawVappVmProvider, errVappVm := rpcClient.Dispense("VAPP_VM_PLUGIN")
	if errVappVm != nil {
		fmt.Println("Error:", errVappVm.Error())

	}
	vappVmProvider := rawVappVmProvider.(grpc.VappVmProvider)

	vcdclient := &VCDClient{client, rpcClient.(*plugin.GRPCClient)}

	//save for later use
	providerGlobalRefPointer = &ProviderGlobalRef{
		pyVcloudProvider:        provider,
		independentDiskProvider: diskProvider,
		orgProvider:             orgProvider,
		userProvider:            userProvider,
		vdcProvider:             vdcProvider,
		vappVmProvider:          vappVmProvider,
	}
	return vcdclient, err

}

func getProvider(m interface{}) grpc.PyVcloudProvider {
	vcdClient := m.(*VCDClient)

	provider := vcdClient.getProvider()

	return provider
}
