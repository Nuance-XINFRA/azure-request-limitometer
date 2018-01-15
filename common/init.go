// wraps azure config/client boilerplate
// exports Conf and Client to be consumed
// by importer
package common

// for outside consumption
var Conf Config
var Client AzureClient

func init() {
	Conf = LoadConfig()
	Client = NewClient(Conf)
}
