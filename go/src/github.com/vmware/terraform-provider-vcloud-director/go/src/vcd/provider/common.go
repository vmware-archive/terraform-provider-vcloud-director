package provider

import (
	"github.com/vmware/terraform-provider-vcloud-director/go/src/vcd/grpc"
)

//used to share the rpc client providers created after login
type ProviderGlobalRef struct {
	pyVcloudProvider        grpc.PyVcloudProvider
	independentDiskProvider grpc.IndependentDiskProvider
	orgProvider             grpc.OrgProvider
	userProvider            grpc.UserProvider
}

var providerGlobalRefPointer *ProviderGlobalRef
