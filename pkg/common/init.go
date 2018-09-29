package common

// Conf Loaded Configuration from azure.json
var Conf Config

// Client Authorized Azure Client
var Client AzureClient

func init() {
	Conf = LoadConfig()
	Client = NewClient(Conf)
}
