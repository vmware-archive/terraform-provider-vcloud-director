/*****************************************************************
* terraform-provider-vcloud-director
* Copyright (c) 2017 VMware, Inc. All Rights Reserved.
* SPDX-License-Identifier: BSD-2-Clause
******************************************************************/

package grpc

import (
	plugin "github.com/hashicorp/go-plugin"
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct {
	client proto.PyVcloudProviderClient
	broker *plugin.GRPCBroker
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl   proto.PyVcloudProviderServer
	broker *plugin.GRPCBroker
}

func (m *GRPCClient) Login(lc proto.LoginCredentials) (*proto.LoginResult, error) {
	result, err := m.client.Login(context.Background(), &lc)
	return result, err
}

func (m *GRPCServer) Login(
	ctx context.Context,
	req *proto.LoginCredentials) (*proto.LoginResult, error) {
	v, err := m.Impl.Login(ctx, req)
	return &proto.LoginResult{Token: v.Token}, err
}

func (m *GRPCClient) ReadCatalog(name string) (*proto.ReadCatalogResult, error) {
	result, err := m.client.ReadCatalog(context.Background(), &proto.Catalog{
		Name: name,
	})
	return result, err
}

func (m *GRPCServer) ReadCatalog(
	ctx context.Context,
	req *proto.Catalog) (*proto.ReadCatalogResult, error) {
	v, err := m.Impl.ReadCatalog(ctx, req)
	return &proto.ReadCatalogResult{Present: v.Present}, err
}

//catalog
func (m *GRPCClient) CreateCatalog(c proto.Catalog) (*proto.CreateCatalogResult, error) {
	result, err := m.client.CreateCatalog(context.Background(), &c)
	return result, err
}

func (m *GRPCServer) CreateCatalog(
	ctx context.Context,
	req *proto.Catalog) (*proto.CreateCatalogResult, error) {
	v, err := m.Impl.CreateCatalog(ctx, req)
	return &proto.CreateCatalogResult{Created: v.Created}, err

}

func (m *GRPCClient) DeleteCatalog(name string) (*proto.DeleteCatalogResult, error) {
	result, err := m.client.DeleteCatalog(context.Background(), &proto.Catalog{
		Name: name,
	})
	return result, err
}

func (m *GRPCServer) DeleteCatalog(
	ctx context.Context,
	req *proto.Catalog) (*proto.DeleteCatalogResult, error) {
	v, err := m.Impl.DeleteCatalog(ctx, req)
	return &proto.DeleteCatalogResult{Deleted: v.Deleted}, err

}

func (m *GRPCClient) UpdateCatalog(c proto.UpdateCatalogInfo) (*proto.UpdateCatalogResult, error) {
	result, err := m.client.UpdateCatalog(context.Background(), &c)
	return result, err
}

func (m *GRPCServer) UpdateCatalog(
	ctx context.Context,
	req *proto.UpdateCatalogInfo) (*proto.UpdateCatalogResult, error) {
	v, err := m.Impl.UpdateCatalog(ctx, req)
	return &proto.UpdateCatalogResult{Updated: v.Updated}, err

}

func (m *GRPCClient) ShareCatalog(c proto.ShareCatalogInfo) (*proto.ShareCatalogResult, error) {
	result, err := m.client.ShareCatalog(context.Background(), &c)
	return result, err
}

func (m *GRPCServer) ShareCatalog(
	ctx context.Context,
	req *proto.ShareCatalogInfo) (*proto.ShareCatalogResult, error) {
	v, err := m.Impl.ShareCatalog(ctx, req)
	return &proto.ShareCatalogResult{Success: v.Success}, err

}

// impl for CatalogUploadMedia

func (m *GRPCClient) CatalogUploadMedia(mediaInfo proto.CatalogUploadMediaInfo) (*proto.CatalogUploadMediaResult, error) {
	result, err := m.client.CatalogUploadMedia(context.Background(), &mediaInfo)
	return result, err
}

func (m *GRPCServer) CatalogUploadMedia(
	ctx context.Context,
	req *proto.CatalogUploadMediaInfo) (*proto.CatalogUploadMediaResult, error) {
	v, err := m.Impl.CatalogUploadMedia(ctx, req)
	return &proto.CatalogUploadMediaResult{Created: v.Created}, err

}

// impl for CatalogUploadOva

func (m *GRPCClient) CatalogUploadOva(ovaInfo proto.CatalogUploadOvaInfo) (*proto.CatalogUploadOvaResult, error) {
	result, err := m.client.CatalogUploadOva(context.Background(), &ovaInfo)
	return result, err
}

func (m *GRPCServer) CatalogUploadOva(
	ctx context.Context,
	req *proto.CatalogUploadOvaInfo) (*proto.CatalogUploadOvaResult, error) {
	v, err := m.Impl.CatalogUploadOva(ctx, req)
	return &proto.CatalogUploadOvaResult{Created: v.Created}, err

}

//impl to check resolved status

func (m *GRPCClient) OvaCheckResolved(ovaInfo proto.CatalogCheckResolvedInfo) (*proto.CheckResolvedResult, error) {
	result, err := m.client.OvaCheckResolved(context.Background(), &ovaInfo)
	return result, err
}

func (m *GRPCServer) OvaCheckResolved(
	ctx context.Context,
	req *proto.CatalogCheckResolvedInfo) (*proto.CheckResolvedResult, error) {
	v, err := m.Impl.OvaCheckResolved(ctx, req)
	return &proto.CheckResolvedResult{Resolved: v.Resolved}, err

}

// impl for DeleteCatalogItem

func (m *GRPCClient) DeleteCatalogItem(itemInfo proto.DeleteCatalogItemInfo) (*proto.DeleteCatalogItemResult, error) {
	result, err := m.client.DeleteCatalogItem(context.Background(), &itemInfo)
	return result, err
}

func (m *GRPCServer) DeleteCatalogItem(
	ctx context.Context,
	req *proto.DeleteCatalogItemInfo) (*proto.DeleteCatalogItemResult, error) {
	v, err := m.Impl.DeleteCatalogItem(ctx, req)
	return &proto.DeleteCatalogItemResult{Deleted: v.Deleted}, err

}

// impl for DeleteCatalogItem

func (m *GRPCClient) IsPresentCatalogItem(itemInfo proto.IsPresentCatalogItemInfo) (*proto.IsPresentCatalogItemResult, error) {
	result, err := m.client.IsPresentCatalogItem(context.Background(), &itemInfo)
	return result, err
}

func (m *GRPCServer) IsPresentCatalogItem(
	ctx context.Context,
	req *proto.IsPresentCatalogItemInfo) (*proto.IsPresentCatalogItemResult, error) {
	v, err := m.Impl.IsPresentCatalogItem(ctx, req)
	return &proto.IsPresentCatalogItemResult{Present: v.Present}, err

}

//Impl for Capture VAPP

// impl for DeleteCatalogItem

func (m *GRPCClient) CaptureVapp(itemInfo proto.CaptureVAppInfo) (*proto.CaptureVAppResult, error) {
	result, err := m.client.CaptureVapp(context.Background(), &itemInfo)
	return result, err
}

func (m *GRPCServer) CaptureVapp(
	ctx context.Context,
	req *proto.CaptureVAppInfo) (*proto.CaptureVAppResult, error) {
	v, err := m.Impl.CaptureVapp(ctx, req)
	return &proto.CaptureVAppResult{Captured: v.Captured}, err

}

// VApp

func (m *GRPCClient) CreateVApp(vappInfo proto.CreateVAppInfo) (*proto.CreateVAppResult, error) {
	result, err := m.client.CreateVApp(context.Background(), &vappInfo)
	return result, err
}

func (m *GRPCServer) CreateVApp(
	ctx context.Context,
	req *proto.CreateVAppInfo) (*proto.CreateVAppResult, error) {
	v, err := m.Impl.CreateVApp(ctx, req)
	return &proto.CreateVAppResult{Created: v.Created}, err

}

func (m *GRPCClient) DeleteVApp(vappInfo proto.DeleteVAppInfo) (*proto.DeleteVAppResult, error) {
	result, err := m.client.DeleteVApp(context.Background(), &vappInfo)
	return result, err
}

func (m *GRPCServer) DeleteVApp(
	ctx context.Context,
	req *proto.DeleteVAppInfo) (*proto.DeleteVAppResult, error) {
	v, err := m.Impl.DeleteVApp(ctx, req)
	return &proto.DeleteVAppResult{Deleted: v.Deleted}, err

}

func (m *GRPCClient) ReadVApp(itemInfo proto.ReadVAppInfo) (*proto.ReadVAppResult, error) {
	result, err := m.client.ReadVApp(context.Background(), &itemInfo)
	return result, err
}

func (m *GRPCServer) ReadVApp(
	ctx context.Context,
	req *proto.ReadVAppInfo) (*proto.ReadVAppResult, error) {
	v, err := m.Impl.ReadVApp(ctx, req)
	return &proto.ReadVAppResult{Present: v.Present}, err

}

func (m *GRPCClient) UpdateVApp(itemInfo proto.UpdateVAppInfo) (*proto.UpdateVAppResult, error) {
	result, err := m.client.UpdateVApp(context.Background(), &itemInfo)
	return result, err
}

func (m *GRPCServer) UpdateVApp(
	ctx context.Context,
	req *proto.UpdateVAppInfo) (*proto.UpdateVAppResult, error) {
	v, err := m.Impl.UpdateVApp(ctx, req)
	return &proto.UpdateVAppResult{Updated: v.Updated}, err

}

//IMPL for stop
func (m *GRPCClient) StopPlugin(stopInfo proto.StopInfo) (*proto.StopResult, error) {
	result, err := m.client.StopPlugin(context.Background(), &stopInfo)
	return result, err
}

func (m *GRPCServer) StopPlugin(
	ctx context.Context,
	req *proto.StopInfo) (*proto.StopResult, error) {
	v, err := m.Impl.StopPlugin(ctx, req)
	return &proto.StopResult{Stopped: v.Stopped}, err

}
