package common

import (
	"os"
)

// Config Returns a config with the Azure resource group and Azure location to perform requests
type Config struct {
	SubscriptionID string `env:"AZURE_SUBSCRIPTION_ID"`
	Location       string `env:"AZURE_LOCATION"`
	ResourceGroup  string `env:"AZURE_RESOURCE_GROUP"`
}

// LoadConfig Returns a Config struct created from Environment Variables
func LoadConfig() (config Config) {
	config = Config{
		SubscriptionID: os.Getenv("AZURE_SUBSCRIPTION_ID"),
		Location:       os.Getenv("AZURE_LOCATION"),
		ResourceGroup:  os.Getenv("AZURE_RESOURCE_GROUP"),
	}
	return
}
