# Azure Request Limitometer

**NOTE: This software is provided AS IS.**

This program makes requests to the Azure Resource Manager API in order to get the headers of the 
remaining requests available to the API. This is important as if we reach the limit, Kubernetes is 
no longer functioning properly as it cannot make requests to Azure.

We authenticate with the `Managed Service Identity` of the VMs running in Azure.

We then output the result into InfluxDB with the following format inside InfluxDB.

```
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


## Building the project

The quickest way to build the project is building it with Docker by running the following command on your computer.

```
docker build -t azure-request-limitometer .
```

If you want to do it without Docker, you must first make sure you have `dep`.

```
go get -u github.com/golang/dep/cmd/dep
```

From there simply run the following command from your computer.

```
go build azure-request-limitometer/cmd/limitometer
```