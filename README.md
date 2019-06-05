# Azure Request Limitometer

[![Build Status](https://travis-ci.org/Nuance-Mobility/azure-request-limitometer.svg?branch=master)](https://travis-ci.org/Nuance-Mobility/azure-request-limitometer)

**NOTE: This software is provided AS IS.**

This program makes requests to the Azure Resource Manager API in order to get the headers of the
remaining requests available to the API. This is important as if we reach the limit, Kubernetes is no longer functioning properly as it cannot make requests to Azure.

We authenticate with the `Managed Service Identity` of the VMs running in Azure.

There are curerntly two supported output logic: `influxdb` and `pushgateway`

The `influxDB` format is the following:

```bash
> select * from "Microsoft.Compute/HighCostGet3Min" limit 5
name: Microsoft.Compute/HighCostGet3Min
time                requestsRemaining
----                -----------------
1536942668898861850 257
1536942729747705521 257
1536942790585657263 258
1536942850472174265 257
1536942909647820539 258
```

The `PushGateway` format is the following:

`azurerm_api_resource_request_remaining_count`

| Element | Value |
| --- | --- |
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\HighCostGet30Min"} |646|
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\HighCostGet3Min"}|137|
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\LowCostGet30Min"}|31522|
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\LowCostGet3Min"}|3976|
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\PutVM30Min"}|3611|
| azurerm_api_resource_request_left_count{job="limitometer",type="Microsoft.Compute\PutVM3Min"}|730 |
| azurerm_api_resource_request_left_count{job="limitometer",type="SubIDReads"}|11694|

## Building the project

The quickest way to build the project is building it with Docker by running the following command on your computer.

```bash
docker build -t azure-request-limitometer .
```

If you want to do it without Docker, you must first make sure you have `dep`.

```bash
go get -u github.com/golang/dep/cmd/dep
```

From there simply run the following command from your computer.

```bash
go build azure-request-limitometer/cmd/limitometer
```
