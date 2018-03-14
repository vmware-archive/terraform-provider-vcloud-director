/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package grpc

import (
	"net/rpc"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"golang.org/x/net/context"

	"github.com/vmware/terraform-provider-vcloud-director/go/src/util/logging"
)

// PyVcloudProvider is the interface that we're exposing as a plugin.
type IndependentDiskProvider interface {
	Create(lc proto.CreateDiskInfo) (*proto.CreateDiskResult, error)

	Delete(lc proto.DeleteDiskInfo) (*proto.DeleteDiskResult, error)

	Read(lc proto.ReadDiskInfo) (*proto.ReadDiskResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type IndependentDiskProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.IndependentDiskServer
}

// GRPCClient is an implementation of KV that talks over RPC.
type DiskGRPCClient struct {
	client proto.IndependentDiskClient
	broker *plugin.GRPCBroker
}

// Here is the gRPC server that GRPCClient talks to.
type DiskGRPCServer struct {
	// This is the real implementation
	Impl proto.IndependentDiskServer
}

// RPCClient is an implementation of KV that talks over RPC.
type DiskRPCClient struct{ client *rpc.Client }

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

type DiskRPCServer struct {
	// This is the real implementation

	DiskImpl proto.IndependentDiskServer
}

func (m *DiskGRPCClient) Create(lc proto.CreateDiskInfo) (*proto.CreateDiskResult, error) {
	result, err := m.client.Create(context.Background(), &lc)
	return result, err
}
func (m *DiskGRPCClient) Delete(lc proto.DeleteDiskInfo) (*proto.DeleteDiskResult, error) {
	result, err := m.client.Delete(context.Background(), &lc)
	return result, err
}
func (m *DiskGRPCClient) Read(lc proto.ReadDiskInfo) (*proto.ReadDiskResult, error) {
	result, err := m.client.Read(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED
func (m *DiskGRPCServer) Create(
	ctx context.Context,
	req *proto.CreateDiskInfo) (*proto.CreateDiskResult, error) {

	v, err := m.Impl.Create(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.CreateDiskResult{DiskId: v.DiskId}, err
}

func (m *DiskGRPCServer) Delete(
	ctx context.Context,
	req *proto.DeleteDiskInfo) (*proto.DeleteDiskResult, error) {
	return &proto.DeleteDiskResult{}, nil
}

func (m *DiskGRPCServer) Read(
	ctx context.Context,
	req *proto.ReadDiskInfo) (*proto.ReadDiskResult, error) {
	return &proto.ReadDiskResult{}, nil
}

// DUMMY
func (*IndependentDiskProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &DiskRPCClient{client: c}, nil
}

func (p *IndependentDiskProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &DiskRPCServer{}, nil
}

func (p *IndependentDiskProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {

	return nil
}

// ONLY GRPC CLIENT IS USE ON THIS SIDE
func (p *IndependentDiskProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	logging.Plog("IndependentDiskProviderPlugin GRPCClient")
	return &DiskGRPCClient{
		client: proto.NewIndependentDiskClient(c),
		broker: broker,
	}, nil
}
