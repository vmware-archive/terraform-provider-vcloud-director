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

type VappVmProvider interface {
	CreateFromVapp(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error)

	CreateFromCatalog(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error)

	Create(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error)

	Delete(lc proto.DeleteVappVmInfo) (*proto.DeleteVappVmResult, error)

	Read(lc proto.ReadVappVmInfo) (*proto.ReadVappVmResult, error)

	PowerOn(lc proto.PowerOnVappVmInfo) (*proto.PowerOnVappVmResult, error)

	PowerOff(lc proto.PowerOffVappVmInfo) (*proto.PowerOffVappVmResult, error)

	ModifyCPU(lc proto.ModifyVappVmCPUInfo) (*proto.ModifyVappVmCPUResult, error)

	ModifyMemory(lc proto.ModifyVappVmMemoryInfo) (*proto.ModifyVappVmMemoryResult, error)

	Update(lc proto.UpdateVappVmInfo) (*proto.UpdateVappVmResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type VappVmProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.VappVmServer
}

// GRPCClient is an implementation of KV that talks over RPC.
type VappVmGRPCClient struct {
	client proto.VappVmClient
	broker *plugin.GRPCBroker
}

// Here is the gRPC server that GRPCClient talks to.
type VappVmGRPCServer struct {
	// This is the real implementation
	Impl proto.VappVmServer
}

// RPCClient is an implementation of KV that talks over RPC.
type VappVmRPCClient struct{ client *rpc.Client }

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

type VappVmRPCServer struct {
	// This is the real implementation

	VappVmImpl proto.VappVmServer
}

func (m *VappVmGRPCClient) CreateFromVapp(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error) {
	result, err := m.client.CreateFromVapp(context.Background(), &lc)
	return result, err
}
func (m *VappVmGRPCClient) CreateFromCatalog(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error) {
	result, err := m.client.CreateFromCatalog(context.Background(), &lc)
	return result, err
}
func (m *VappVmGRPCClient) Create(lc proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error) {
	result, err := m.client.Create(context.Background(), &lc)
	return result, err
}
func (m *VappVmGRPCClient) Delete(lc proto.DeleteVappVmInfo) (*proto.DeleteVappVmResult, error) {
	result, err := m.client.Delete(context.Background(), &lc)
	return result, err
}
func (m *VappVmGRPCClient) Read(lc proto.ReadVappVmInfo) (*proto.ReadVappVmResult, error) {
	result, err := m.client.Read(context.Background(), &lc)
	return result, err
}

func (m *VappVmGRPCClient) Update(lc proto.UpdateVappVmInfo) (*proto.UpdateVappVmResult, error) {
	result, err := m.client.Update(context.Background(), &lc)
	return result, err
}

func (m *VappVmGRPCClient) PowerOn(lc proto.PowerOnVappVmInfo) (*proto.PowerOnVappVmResult, error) {
	result, err := m.client.PowerOn(context.Background(), &lc)
	return result, err
}

func (m *VappVmGRPCClient) PowerOff(lc proto.PowerOffVappVmInfo) (*proto.PowerOffVappVmResult, error) {
	result, err := m.client.PowerOff(context.Background(), &lc)
	return result, err
}

func (m *VappVmGRPCClient) ModifyCPU(lc proto.ModifyVappVmCPUInfo) (*proto.ModifyVappVmCPUResult, error) {
	result, err := m.client.ModifyCPU(context.Background(), &lc)
	return result, err
}

func (m *VappVmGRPCClient) ModifyMemory(lc proto.ModifyVappVmMemoryInfo) (*proto.ModifyVappVmMemoryResult, error) {
	result, err := m.client.ModifyMemory(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED
func (m *VappVmGRPCServer) Create(
	ctx context.Context,
	req *proto.CreateVappVmInfo) (*proto.CreateVappVmResult, error) {

	v, err := m.Impl.Create(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.CreateVappVmResult{Created: v.Created}, err
}

func (m *VappVmGRPCServer) Delete(
	ctx context.Context,
	req *proto.DeleteVappVmInfo) (*proto.DeleteVappVmResult, error) {
	return &proto.DeleteVappVmResult{}, nil
}

func (m *VappVmGRPCServer) Read(
	ctx context.Context,
	req *proto.ReadVappVmInfo) (*proto.ReadVappVmResult, error) {
	return &proto.ReadVappVmResult{}, nil
}

func (m *VappVmGRPCServer) Update(
	ctx context.Context,
	req *proto.UpdateVappVmInfo) (*proto.UpdateVappVmResult, error) {
	return &proto.UpdateVappVmResult{}, nil
}

func (m *VappVmGRPCServer) PowerOn(
	ctx context.Context,
	req *proto.PowerOnVappVmInfo) (*proto.PowerOnVappVmResult, error) {
	return &proto.PowerOnVappVmResult{}, nil
}

func (m *VappVmGRPCServer) PowerOff(
	ctx context.Context,
	req *proto.PowerOffVappVmInfo) (*proto.PowerOffVappVmResult, error) {
	return &proto.PowerOffVappVmResult{}, nil
}

func (m *VappVmGRPCServer) ModifyCPU(
	ctx context.Context,
	req *proto.ModifyVappVmCPUInfo) (*proto.ModifyVappVmCPUResult, error) {
	return &proto.ModifyVappVmCPUResult{}, nil
}

func (m *VappVmGRPCServer) ModifyMemory(
	ctx context.Context,
	req *proto.ModifyVappVmMemoryInfo) (*proto.ModifyVappVmMemoryResult, error) {
	return &proto.ModifyVappVmMemoryResult{}, nil
}

// DUMMY
func (*VappVmProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &VappVmRPCClient{client: c}, nil
}

func (p *VappVmProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &VappVmRPCServer{}, nil
}

func (p *VappVmProviderPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {

	return nil
}

// ONLY GRPC CLIENT IS USE ON THIS SIDE
func (p *VappVmProviderPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	logging.Plog("VappVmProviderPlugin GRPCClient")
	return &VappVmGRPCClient{
		client: proto.NewVappVmClient(c),
		broker: broker,
	}, nil
}
