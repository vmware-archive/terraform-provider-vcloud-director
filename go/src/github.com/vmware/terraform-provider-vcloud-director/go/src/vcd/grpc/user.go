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

type UserProvider interface {
	Create(lc proto.CreateUserInfo) (*proto.CreateUserResult, error)

	Delete(lc proto.DeleteUserInfo) (*proto.DeleteUserResult, error)

	Read(lc proto.ReadUserInfo) (*proto.ReadUserResult, error)

	Update(lc proto.UpdateUserInfo) (*proto.UpdateUserResult, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
// We also implement GRPCPlugin so that this plugin can be served over
// gRPC.
type UserProviderPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl proto.UserServer
}

// GRPCClient is an implementation of KV that talks over RPC.
type UserGRPCClient struct {
	client proto.UserClient
}

// Here is the gRPC server that GRPCClient talks to.
type UserGRPCServer struct {
	// This is the real implementation
	Impl proto.UserServer
}

// RPCClient is an implementation of KV that talks over RPC.
type UserRPCClient struct{ client *rpc.Client }

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc

type UserRPCServer struct {
	// This is the real implementation

	UserImpl proto.UserServer
}

func (m *UserGRPCClient) Create(lc proto.CreateUserInfo) (*proto.CreateUserResult, error) {
	result, err := m.client.Create(context.Background(), &lc)
	return result, err
}
func (m *UserGRPCClient) Delete(lc proto.DeleteUserInfo) (*proto.DeleteUserResult, error) {
	result, err := m.client.Delete(context.Background(), &lc)
	return result, err
}
func (m *UserGRPCClient) Read(lc proto.ReadUserInfo) (*proto.ReadUserResult, error) {
	result, err := m.client.Read(context.Background(), &lc)
	return result, err
}

func (m *UserGRPCClient) Update(lc proto.UpdateUserInfo) (*proto.UpdateUserResult, error) {
	result, err := m.client.Update(context.Background(), &lc)
	return result, err
}

//DUMMY IMPL NOT INVOKED
func (m *UserGRPCServer) Create(
	ctx context.Context,
	req *proto.CreateUserInfo) (*proto.CreateUserResult, error) {

	v, err := m.Impl.Create(ctx, req)
	logging.Plog("======CHECK THIS ===========")
	return &proto.CreateUserResult{Created: v.Created}, err
}

func (m *UserGRPCServer) Delete(
	ctx context.Context,
	req *proto.DeleteUserInfo) (*proto.DeleteUserResult, error) {
	return &proto.DeleteUserResult{}, nil
}

func (m *UserGRPCServer) Read(
	ctx context.Context,
	req *proto.ReadUserInfo) (*proto.ReadUserResult, error) {
	return &proto.ReadUserResult{}, nil
}

func (m *UserGRPCServer) Update(
	ctx context.Context,
	req *proto.UpdateUserInfo) (*proto.UpdateUserResult, error) {
	return &proto.UpdateUserResult{}, nil
}

// DUMMY
func (*UserProviderPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {

	return &UserRPCClient{client: c}, nil
}

func (p *UserProviderPlugin) Server(*plugin.MuxBroker) (interface{}, error) {

	return &UserRPCServer{}, nil
}

func (p *UserProviderPlugin) GRPCServer(s *grpc.Server) error {

	return nil
}

// ONLY GRPC CLIENT IS USE ON THIS SIDE
func (p *UserProviderPlugin) GRPCClient(c *grpc.ClientConn) (interface{}, error) {
	logging.Plog("UserProviderPlugin GRPCClient")
	return &UserGRPCClient{client: proto.NewUserClient(c)}, nil
}
