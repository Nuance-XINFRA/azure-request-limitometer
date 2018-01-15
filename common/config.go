package common

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

type Config struct {
	EnvironmentName string `json:"cloud" env:"AZURE_ENVIRONMENT_NAME"`
	SubscriptionID  string `json:"subscriptionId" env:"AZURE_SUBSCRIPTION_ID"`
	TenantID        string `json:"tenantId" env:"AZURE_TENANT_ID"`
	ClientID        string `json:"aadClientId" env:"AZURE_CLIENT_ID"`
	ClientSecret    string `json:"aadClientSecret" env:"AZURE_CLIENT_SECRET"`
	ResourceGroup   string `json:"resourceGroup" env:"AZURE_RESOURCE_GROUP"`
}

// LoadConfig reads the necessary Azure-related config from
// /etc/kubernetes/azure.json and environment variables and stores
// them into `config`.
//
// It's a bit of a mess because of the reflection, but
// it's a nice pattern.
func LoadConfig() (config Config) {
	var err error

	// JSON
	configFile, err := os.Open("/etc/kubernetes/azure.json")
	if err == nil {
		configContent, err := ioutil.ReadAll(configFile)
		xk(err)
		json.Unmarshal(configContent, &config)
	} else {
		log.Println("/etc/kubernetes/azure.json not found: falling back to environment")
	}

	// ENVIRONMENT
	configValue := reflect.ValueOf(&config).Elem()
	configType := configValue.Type()

	// this loop is an awkward mix of sanity checks and environment loading
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)

		fieldJson, ok := field.Tag.Lookup("json")
		assert(ok, "struct Config should have `json` tag on all fields")

		fieldEnv, ok := field.Tag.Lookup("env")
		assert(ok, "struct Config should have `env` tag on all fields")

		if envval := os.Getenv(fieldEnv); envval != "" {
			configValue.FieldByName(field.Name).SetString(envval)
		}

		finalValue := configValue.Field(i).Interface().(string)

		xk(os.Setenv(fieldEnv, finalValue))

		if finalValue == "" {
			err = multierror.Append(
				err,
				fmt.Errorf(
					"either %v must be defined in %v, or %v must be set",
					fieldJson,
					"/etc/kubernetes/azure.json",
					fieldEnv,
				))
		}
	}
	xk(err)
	return
}
