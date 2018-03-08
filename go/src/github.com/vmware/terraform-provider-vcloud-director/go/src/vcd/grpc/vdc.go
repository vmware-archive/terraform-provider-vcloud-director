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

type VdcProvider interface {
	Create(lc proto.CreateVdcInfo) (*proto.CreateVdcResult, error)

	Delete(lc proto.DeleteVdcInfo) (*proto.DeleteVdcResult, error)

	Read(lc proto.ReadVdcInfo) (*proto.ReadVdcResult, error)

	Update(lc proto.UpdateVdcInfo) (*proto.UpdateVdcResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type VdcProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.VdcServer
}

// GRPCClient is an implementation of KV that talks over RPC.
type VdcGRPCClient struct {
	client proto.VdcClient
}

// Here is the gRPC server that GRPCClient talks to.
type VdcGRPCServer struct {
	// This is the real implementation
	Impl proto.VdcServer
}

// RPCClient is an implementation of KV that talks over RPC.
type VdcRPCClient struct{ client *rpc.Client }

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

type VdcRPCServer struct {
	// This is the real implementation

	VdcImpl proto.VdcServer
}

func (m *VdcGRPCClient) Create(lc proto.CreateVdcInfo) (*proto.CreateVdcResult, error) {
	result, err := m.client.Create(context.Background(), &lc)
	return result, err
}
func (m *VdcGRPCClient) Delete(lc proto.DeleteVdcInfo) (*proto.DeleteVdcResult, error) {
	result, err := m.client.Delete(context.Background(), &lc)
	return result, err
}
func (m *VdcGRPCClient) Read(lc proto.ReadVdcInfo) (*proto.ReadVdcResult, error) {
	result, err := m.client.Read(context.Background(), &lc)
	return result, err
}

func (m *VdcGRPCClient) Update(lc proto.UpdateVdcInfo) (*proto.UpdateVdcResult, error) {
	result, err := m.client.Update(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED
func (m *VdcGRPCServer) Create(
	ctx context.Context,
	req *proto.CreateVdcInfo) (*proto.CreateVdcResult, error) {

	v, err := m.Impl.Create(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.CreateVdcResult{Created: v.Created}, err
}

func (m *VdcGRPCServer) Delete(
	ctx context.Context,
	req *proto.DeleteVdcInfo) (*proto.DeleteVdcResult, error) {
	return &proto.DeleteVdcResult{}, nil
}

func (m *VdcGRPCServer) Read(
	ctx context.Context,
	req *proto.ReadVdcInfo) (*proto.ReadVdcResult, error) {
	return &proto.ReadVdcResult{}, nil
}

func (m *VdcGRPCServer) Update(
	ctx context.Context,
	req *proto.UpdateVdcInfo) (*proto.UpdateVdcResult, error) {
	return &proto.UpdateVdcResult{}, nil
}

// DUMMY
func (*VdcProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &VdcRPCClient{client: c}, nil
}

func (p *VdcProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &VdcRPCServer{}, nil
}

func (p *VdcProviderPlugin) GRPCServer(s *grpc.Server) error {

	return nil
}

// ONLY GRPC CLIENT IS USE ON THIS SIDE
func (p *VdcProviderPlugin) GRPCClient(c *grpc.ClientConn) (interface{}, error) {
	logging.Plog("VdcProviderPlugin GRPCClient")
	return &VdcGRPCClient{client: proto.NewVdcClient(c)}, nil
}
