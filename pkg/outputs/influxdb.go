package outputs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/influxdata/influxdb/client/v2"
)

// InfluxDBServer This struct contains the information necessary to connect to a InfluxDB server
// such as host, port and database
type InfluxDBServer struct {
	Host     string
	Port     string
	Database string
}

// GetInfluxdbConfig Generates a server config from environment variables
func GetInfluxdbConfig() InfluxDBServer {
	server := InfluxDBServer{
		Host:     os.Getenv("INFLUXDB_HOST"),
		Port:     os.Getenv("INFLUXDB_PORT"),
		Database: os.Getenv("INFLUXDB_DATABASE"),
	}
	return server
}

// WriteOutputInflux Creates a Batch of points given a map
func WriteOutputInflux(values map[string]int, fieldName string) {
	s := GetInfluxdbConfig()

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:%s", s.Host, s.Port),
	})
	if err != nil {
		glog.Fatal(err)
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  s.Database,
		Precision: "s",
	})
	if err != nil {
		glog.Fatal(err)
	}

	for k, v := range values {
		tags := make(map[string]string)
		fields := map[string]interface{}{
			fieldName: v,
		}

		pt, err := client.NewPoint(k, tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	if err := c.Write(bp); err != nil {
		glog.Fatal(err)
	}

	glog.Info("Successfully wrote to InfluxDB")
}
