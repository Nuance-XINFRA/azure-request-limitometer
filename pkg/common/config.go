package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const azureInstanceMetadataEndpoint = "http://169.254.169.254/metadata/instance"

// Queries the Azure Instance Metadata Service for the instance's compute metadata
func retrieveComputeInstanceMetadata() (metadata ComputeInstanceMetadata, err error) {
	c := &http.Client{}

	req, _ := http.NewRequest("GET", azureInstanceMetadataEndpoint+"/compute", nil)
	req.Header.Add("Metadata", "True")
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("api-version", "2017-12-01")
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		err = fmt.Errorf("sending Azure Instance Metadata Service request failed: %v", err)
	}
	defer resp.Body.Close()

	rawJSON, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		err = fmt.Errorf("reading response body failed: %v", err)
		return
	}
	if err := json.Unmarshal(rawJSON, &metadata); err != nil {
		err = fmt.Errorf("unmarshaling JSON response failed: %v", err)
	}

	return
}

// LoadConfig Returns a Config struct created from Environment Variables
func LoadConfig() (config Config) {
	m, err := retrieveComputeInstanceMetadata()
	if err != nil {
		err = fmt.Errorf("unable to load the config: %v", err)
	}

	config = Config{
		VMName:         m.Name,
		SubscriptionID: m.SubscriptionID,
		Location:       m.Location,
		ResourceGroup:  m.ResourceGroupName,
	}

	return
}
