package common

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

// AzureClient This is an authorized client for Azure communication.
type AzureClient struct {
	compute.VirtualMachinesClient
	compute.DisksClient
}

// NewClient Initialized an authorized Azure client
func NewClient(config Config) (client AzureClient) {
	client = AzureClient{
		compute.NewVirtualMachinesClient(config.SubscriptionID),
		compute.NewDisksClient(config.SubscriptionID),
	}

	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		glog.Fatalf("Unable to authorize the client. Reporting Error: %s", err)
	}

	client.VirtualMachinesClient.Authorizer = authorizer
	client.DisksClient.Authorizer = authorizer

	client.VirtualMachinesClient.RetryAttempts = 1
	client.DisksClient.RetryAttempts = 1

	return
}

// GetVM Returns a VirtualMachine object.
func (c AzureClient) GetVM(vmname string) compute.VirtualMachine {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	vm, err := c.VirtualMachinesClient.Get(ctx, Conf.ResourceGroup, vmname, compute.InstanceView)
	if err != nil {
		glog.Fatalf("failed to get VM; check HTTP_PROXY: %v", err)
	}
	return vm
}

// GetAllVM Returns a ListResultPage of all VMs in the ResourceGroup of the Config
func (c AzureClient) GetAllVM() compute.VirtualMachineListResultPage {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	result, err := c.VirtualMachinesClient.List(ctx, Conf.ResourceGroup)
	if err != nil {
		glog.Fatalf("failed to get all VMs; check HTTP_PROXY: %v", err)
	}

	return result
}

// PutVM returns the Virtual Machine object
func (c AzureClient) PutVM(vmname string) autorest.Response {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	node := c.GetVM(vmname)

	req, err := c.VirtualMachinesClient.CreateOrUpdatePreparer(ctx, Conf.ResourceGroup, vmname, node)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.VirtualMachinesClient", "CreateOrUpdate", nil, "Failure preparing request")
	}

	var resp *http.Response
	resp, err = autorest.SendWithSender(c.VirtualMachinesClient, req,
		azure.DoRetryWithRegistration(c.VirtualMachinesClient.Client))
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated))
	if err != nil {
		glog.Fatal(err)
	}

	return autorest.Response{resp}
}
