package common

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/network/mgmt/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/golang/glog"
)

// AzureClient This is an authorized client for Azure communication.
type AzureClient struct {
	compute.VirtualMachinesClient
	compute.DisksClient
	compute.VirtualMachineScaleSetsClient
	network.InterfacesClient
}

// NewClient Initialized an authorized Azure client
func NewClient(config Config) (client AzureClient) {
	client = AzureClient{
		compute.NewVirtualMachinesClientWithBaseURI(config.EnvironmentEndpoint, config.SubscriptionID),
		compute.NewDisksClientWithBaseURI(config.EnvironmentEndpoint, config.SubscriptionID),
		compute.NewVirtualMachineScaleSetsClientWithBaseURI(config.EnvironmentEndpoint, config.SubscriptionID),
		network.NewInterfacesClientWithBaseURI(config.EnvironmentEndpoint, config.SubscriptionID),
	}

	// Authorizing with Managed Service Identity
	err := os.Setenv("AZURE_ENVIRONMENT", config.AzureEnvironment)
	if err != nil {
		log.Panicf("Unable to set Azure environment: %v", err)
	}
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		client.VirtualMachinesClient.Authorizer = authorizer
		client.DisksClient.Authorizer = authorizer
		client.VirtualMachineScaleSetsClient.Authorizer = authorizer
		client.InterfacesClient.Authorizer = authorizer
	}

	client.VirtualMachinesClient.RetryAttempts = 1
	client.DisksClient.RetryAttempts = 1
	client.VirtualMachineScaleSetsClient.RetryAttempts = 1
	client.InterfacesClient.RetryAttempts = 1

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

// GetNic Returns a Nic object.
func (c AzureClient) GetNic(vmName string) network.Interface {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	nicName := c.getNicNameFromVMName(vmName)

	nic, err := c.InterfacesClient.Get(ctx, Conf.ResourceGroup, nicName, "")
	if err != nil {
		glog.Fatalf("failed to get Nic: %v", err)
	}
	return nic
}

// GetAllNics Returns a ListResultPage of all Interfaces in the ResourceGroup of the Config
func (c AzureClient) GetAllNics() network.InterfaceListResultPage {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	result, err := c.InterfacesClient.List(ctx, Conf.ResourceGroup)
	if err != nil {
		glog.Fatalf("failed to get all Interfaces; check HTTP_PROXY: %v", err)
	}

	return result
}

// PutNic returns the Interface object
func (c AzureClient) PutNic(vmName string) autorest.Response {
	ctx, cancel := context.WithTimeout(context.Background(), 6000*time.Second)
	defer cancel()

	nicName := c.getNicNameFromVMName(vmName)

	nic := c.GetNic(nicName)

	req, err := c.InterfacesClient.CreateOrUpdatePreparer(ctx, Conf.ResourceGroup, nicName, nic)
	if err != nil {
		err = autorest.NewErrorWithError(err, "compute.InterfacesClient", "CreateOrUpdatePreparer", nil, "Failure preparing request")
	}

	var resp *http.Response
	resp, err = autorest.SendWithSender(c.InterfacesClient, req,
		azure.DoRetryWithRegistration(c.InterfacesClient.Client))
	err = autorest.Respond(resp, azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated))
	if err != nil {
		glog.Fatal(err)
	}

	return autorest.Response{Response: resp}
}

func (c AzureClient) getNicNameFromVMName(vmName string) string {
	vm := c.GetVM(vmName)

	primaryNicID, err := getPrimaryInterfaceID(vm)
	if err != nil {
		glog.Fatalf("failed to getPrimaryInterfaceID from VM: %v", err)
	}

	nicName, err := getLastSegment(primaryNicID)
	if err != nil {
		glog.Fatalf("failed to nic name from nicID: %v", err)
	}

	return nicName

}

// This returns the full identifier of the primary NIC for the given VM.
func getPrimaryInterfaceID(machine compute.VirtualMachine) (string, error) {
	if len(*machine.NetworkProfile.NetworkInterfaces) == 1 {
		return *(*machine.NetworkProfile.NetworkInterfaces)[0].ID, nil
	}

	for _, ref := range *machine.NetworkProfile.NetworkInterfaces {
		if *ref.Primary {
			return *ref.ID, nil
		}
	}

	return "", fmt.Errorf("failed to find a primary nic for the vm. vmname=%q", *machine.Name)
}

// returns the deepest child's identifier from a full identifier string.
func getLastSegment(ID string) (string, error) {
	parts := strings.Split(ID, "/")
	name := parts[len(parts)-1]
	if len(name) == 0 {
		return "", fmt.Errorf("resource name was missing from identifier")
	}

	return name, nil
}
