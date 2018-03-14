/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package grpc

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "PyCloud",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"PY_PLUGIN": &PyVcloudProviderPlugin{},

	"DISK_PLUGIN":    &IndependentDiskProviderPlugin{},
	"ORG_PLUGIN":     &OrgProviderPlugin{},
	"USER_PLUGIN":    &UserProviderPlugin{},
	"VDC_PLUGIN":     &VdcProviderPlugin{},
	"VAPP_VM_PLUGIN": &VappVmProviderPlugin{},
}

// PyVcloudProvider is the interface that we're exposing as a plugin.
type PyVcloudProvider interface {
	Login(lc proto.LoginCredentials) (*proto.LoginResult, error)

	//Catalog
	ReadCatalog(name string) (*proto.ReadCatalogResult, error)

	CreateCatalog(c proto.Catalog) (*proto.CreateCatalogResult, error)

	DeleteCatalog(name string) (*proto.DeleteCatalogResult, error)

	//Catalog Item
	CatalogUploadMedia(c proto.CatalogUploadMediaInfo) (*proto.CatalogUploadMediaResult, error)

	CatalogUploadOva(c proto.CatalogUploadOvaInfo) (*proto.CatalogUploadOvaResult, error)

	OvaCheckResolved(c proto.CatalogCheckResolvedInfo) (*proto.CheckResolvedResult, error)

	DeleteCatalogItem(c proto.DeleteCatalogItemInfo) (*proto.DeleteCatalogItemResult, error)

	IsPresentCatalogItem(c proto.IsPresentCatalogItemInfo) (*proto.IsPresentCatalogItemResult, error)

	//capture VAPP

	CaptureVapp(c proto.CaptureVAppInfo) (*proto.CaptureVAppResult, error)

	//VAPP
	CreateVApp(c proto.CreateVAppInfo) (*proto.CreateVAppResult, error)

	DeleteVApp(c proto.DeleteVAppInfo) (*proto.DeleteVAppResult, error)

	ReadVApp(c proto.ReadVAppInfo) (*proto.ReadVAppResult, error)

	UpdateVApp(c proto.UpdateVAppInfo) (*proto.UpdateVAppResult, error)

	//Plugin Remote control
	StopPlugin(c proto.StopInfo) (*proto.StopResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type PyVcloudProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.PyVcloudProviderServer
}

func (*PyVcloudProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &RPCClient{client: c}, nil
}

func (p *PyVcloudProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &RPCServer{Impl: p.Impl}, nil
}

func (p *PyVcloudProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {

	proto.RegisterPyVcloudProviderServer(s, &GRPCServer{
		Impl:   p.Impl,
		broker: broker,
	})
	return nil
}

func (p *PyVcloudProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {

	return &GRPCClient{
		client: proto.NewPyVcloudProviderClient(c),
		broker: broker,
	}, nil
}
