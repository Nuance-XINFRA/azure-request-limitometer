package common

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/disk"
	"github.com/Azure/go-autorest/autorest/azure"
	aru "github.com/Azure/go-autorest/autorest/utils"
	"net/http"
	"time"
)

type AzureClient struct {
	compute.VirtualMachinesClient
	disk.DisksClient
}

func (c AzureClient) GetVM(vmname string) compute.VirtualMachine {
	vm, err := c.VirtualMachinesClient.Get(Conf.ResourceGroup, vmname, compute.InstanceView)
	if err != nil {
		panic(fmt.Errorf("failed to get VM; check HTTP_PROXY: %v", err))
	}
	return vm
}

func (c AzureClient) GetAllVM() compute.VirtualMachineListResult {
	list, err := c.VirtualMachinesClient.List(Conf.ResourceGroup)
	if err != nil {
		panic(fmt.Errorf("failed to get all VMs; check HTTP_PROXY: %v", err))
	}
	return list
}

func (c AzureClient) PutVM(vmname string) compute.VirtualMachine {
	vm := c.GetVM(vmname)
	vm.Resources = nil // similar to https://github.com/Azure/azure-sdk-for-java/issues/600

	cancel := make(chan struct{})
	vms, errs := c.VirtualMachinesClient.CreateOrUpdate(
		Conf.ResourceGroup,
		vmname,
		vm, // TODO repeat, maybe recycle
		cancel,
	)

	select {
	case vm := <-vms:
		return vm
	case err := <-errs:
		if err != nil {
			cancel <- struct{}{}
			panic(err)
		}
	}

	assert(false, "should be unreachable")
	return <-vms
}

func NewClient(config Config) (client AzureClient) {
	cloud, err := azure.EnvironmentFromName(config.EnvironmentName)
	xk(err)

	// assumes some env vars are set - TODO mention which
	authorizer, err := aru.GetAuthorizer(cloud)
	xk(err)

	client = AzureClient{
		compute.NewVirtualMachinesClient(config.SubscriptionID),
		disk.NewDisksClient(config.SubscriptionID),
	}

	// TODO do this better
	// client.Authorizer = authorizer
	client.VirtualMachinesClient.Authorizer = authorizer
	client.DisksClient.Authorizer = authorizer
	client.VirtualMachinesClient.RetryAttempts = 1
	client.DisksClient.RetryAttempts = 1

	httpClient, ok := client.VirtualMachinesClient.Sender.(*http.Client)
	assert(ok, "default autorest Sender should be *http.Client")
	httpClient.Timeout = 10 * time.Second

	return
}
