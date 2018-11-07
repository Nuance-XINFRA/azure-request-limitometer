package common

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AzureClient This is an authorized client for Azure communication.
type AzureClient struct {
	compute.VirtualMachinesClient
	compute.DisksClient
	compute.VirtualMachineScaleSetsClient
}

// NewClient Initialized an authorized Azure client
func NewClient(config Config) (client AzureClient) {
	client = AzureClient{
		compute.NewVirtualMachinesClient(config.SubscriptionID),
		compute.NewDisksClient(config.SubscriptionID),
		compute.NewVirtualMachineScaleSetsClient(config.SubscriptionID),
	}

	// Authorizing with Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		client.VirtualMachinesClient.Authorizer = authorizer
		client.DisksClient.Authorizer = authorizer
		client.VirtualMachineScaleSetsClient.Authorizer = authorizer
	}

	client.VirtualMachinesClient.RetryAttempts = 1
	client.DisksClient.RetryAttempts = 1
	client.VirtualMachineScaleSetsClient.RetryAttempts = 1

	return
}

// GetVM Returns a VirtualMachine object.
func (c AzureClient) GetVM(vmname string) (vm compute.VirtualMachine) {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	vm, err := c.VirtualMachinesClient.Get(ctx, Conf.ResourceGroup, vmname, compute.InstanceView)
	if err != nil {
		log.Panicf("failed to get VM: %v", err)
	}

	return
}

// GetAllVM Returns a ListResultPage of all VMs in the ResourceGroup of the Config
func (c AzureClient) GetAllVM() (result compute.VirtualMachineListResultPage) {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	result, err := c.VirtualMachinesClient.List(ctx, Conf.ResourceGroup)
	if err != nil {
		log.Panicf("failed to get all VMs: %v", err)
	}

	return
}

// PutVM returns the Virtual Machine object
func (c AzureClient) PutVM(vmname string) (res autorest.Response) {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	node := c.GetVM(vmname)

	req, err := c.VirtualMachinesClient.CreateOrUpdatePreparer(ctx, Conf.ResourceGroup, vmname, node)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "CreateOrUpdate", nil, "Failure preparing request")
		log.Panic(err)
	}

	var result *http.Response
	result, err = autorest.SendWithSender(c.VirtualMachinesClient, req,
		azure.DoRetryWithRegistration(c.VirtualMachinesClient.Client))
	err = autorest.Respond(result, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated))
	if err != nil {
		log.Panic(err)
	}
	res.Response = result

	return
}
