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

// OrgProvider is the interface that we're exposing as a plugin.
type OrgProvider interface {
	//Create a org
	Create(lc proto.CreateOrgInfo) (*proto.CreateOrgResult, error)

	//Delete a org
	Delete(lc proto.DeleteOrgInfo) (*proto.DeleteOrgResult, error)

	//Update org
	Update(lc proto.UpdateOrgInfo) (*proto.UpdateOrgResult, error)

	//Read Org
	Read(lc proto.ReadOrgInfo) (*proto.ReadOrgResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type OrgProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.OrgServer
}

// GRPCClient is an implementation of KV that talks over RPC.
type OrgGRPCClient struct {
	client proto.OrgClient
}

// Here is the gRPC server that GRPCClient talks to.
type OrgGRPCServer struct {
	// This is the real implementation
	Impl proto.OrgServer
}

// RPCClient is an implementation of KV that talks over RPC.
type OrgRPCClient struct{ client *rpc.Client }

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

type OrgRPCServer struct {
	// This is the real implementation

	OrgImpl proto.OrgServer
}

//Create a org
func (m *OrgGRPCClient) Create(lc proto.CreateOrgInfo) (*proto.CreateOrgResult, error) {
	result, err := m.client.Create(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED as go side grpc implementation is for client only
func (m *OrgGRPCServer) Create(
	ctx context.Context,
	req *proto.CreateOrgInfo) (*proto.CreateOrgResult, error) {

	v, err := m.Impl.Create(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.CreateOrgResult{Created: v.Created}, err
}

//Read a org
func (m *OrgGRPCClient) Read(lc proto.ReadOrgInfo) (*proto.ReadOrgResult, error) {
	result, err := m.client.Read(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED as go side grpc implementation is for client only
func (m *OrgGRPCServer) Read(
	ctx context.Context,
	req *proto.ReadOrgInfo) (*proto.ReadOrgResult, error) {

	v, err := m.Impl.Read(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.ReadOrgResult{Present: v.Present}, err
}

//Update a org
func (m *OrgGRPCClient) Update(lc proto.UpdateOrgInfo) (*proto.UpdateOrgResult, error) {
	result, err := m.client.Update(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED as go side grpc implementation is for client only
func (m *OrgGRPCServer) Update(
	ctx context.Context,
	req *proto.UpdateOrgInfo) (*proto.UpdateOrgResult, error) {

	v, err := m.Impl.Update(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.UpdateOrgResult{Updated: v.Updated}, err
}

//Delete a org
func (m *OrgGRPCClient) Delete(lc proto.DeleteOrgInfo) (*proto.DeleteOrgResult, error) {
	result, err := m.client.Delete(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED as go side grpc implementation is for client only
func (m *OrgGRPCServer) Delete(
	ctx context.Context,
	req *proto.DeleteOrgInfo) (*proto.DeleteOrgResult, error) {

	v, err := m.Impl.Delete(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.DeleteOrgResult{Deleted: v.Deleted}, err
}

// DUMMY
func (*OrgProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &OrgRPCClient{client: c}, nil
}

func (p *OrgProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &OrgRPCServer{}, nil
}

func (p *OrgProviderPlugin) GRPCServer(s *grpc.Server) error {

	return nil
}

// ONLY GRPC CLIENT IS USE ON THIS SIDE
func (p *OrgProviderPlugin) GRPCClient(c *grpc.ClientConn) (interface{}, error) {
	logging.Plog("OrgProviderPlugin GRPCClient")
	return &OrgGRPCClient{client: proto.NewOrgClient(c)}, nil
}
