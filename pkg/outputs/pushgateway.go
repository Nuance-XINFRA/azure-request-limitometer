package outputs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	remaining = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "azurerm_api_resource_request_left_count",
		Help: "The number of requests left for the resource type.",
	})
)

// PushGatewayServer This struct contains the information necessary to connect to a PushGateway server
// such as host and port
type PushGatewayServer struct {
	Host string
	Port string
}

// GetPushGatewayConfig Generates a server config from environment variables
func GetPushGatewayConfig() PushGatewayServer {
	server := PushGatewayServer{
		Host: os.Getenv("PUSHGATEWAY_HOST"),
		Port: os.Getenv("PUSHGATEWAY_PORT"),
	}
	return server
}

// WriteOutputPushGateway pushes metrics to the pushgatewayw
func WriteOutputPushGateway(values map[string]int) {
	s := GetPushGatewayConfig()

	for k, v := range values {
		pusher := push.New(fmt.Sprintf("http://%s:%s", s.Host, s.Port), "limitometer")
		remaining.Set(float64(v))
		// Note that / cannot be used as part of a label value or the job name,
		// even if escaped as %2F. (The decoding happens before the path routing kicks in,
		//cf. the Go documentation of URL.Path.)
		pusher.Collector(remaining).Grouping("type", strings.Replace(k, "/", "\\", 1))
		if err := pusher.Push(); err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Successfully wrote to PushGateway")
}
