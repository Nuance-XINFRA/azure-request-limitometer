package common

// ComputeInstanceMetadata denotes the compute specific instance metadata from Azure Instance
// Metadata Service
type ComputeInstanceMetadata struct {
	Location             string `json:"location,omitempty"`
	Name                 string `json:"name,omitempty"`
	Offer                string `json:"offer,omitempty"`
	OsType               string `json:"osType,omitempty"`
	PlacementGroupID     string `json:"placement_group_id,omitempty"`
	PlatformFaultDomain  string `json:"platformFaultDomain,omitempty"`
	PlatformUpdateDomain string `json:"platformUpdateDomain,omitempty"`
	Publisher            string `json:"publisher,omitempty"`
	ResourceGroupName    string `json:"resourceGroupName,omitempty"`
	Sku                  string `json:"sku,omitempty"`
	SubscriptionID       string `json:"subscriptionId,omitempty"`
	Tags                 string `json:"tags,omitempty"`
	Version              string `json:"version,omitempty"`
	VMID                 string `json:"vmId,omitempty"`
	VMScaleSetName       string `json:"vmScaleSetName,omitempty"`
	VMSize               string `json:"vmSize,omitempty"`
	Zone                 string `json:"zone,omitempty"`
}

// Config returns a config with the Azure resource group and Azure location to perform requests
type Config struct {
	VMName         string
	SubscriptionID string
	Location       string
	ResourceGroup  string
}
