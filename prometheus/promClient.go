package prometheus


import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)
type PrometheusRegistryStruct struct {
	Client *prom.Registry
}
var PrometheusRegistry *PrometheusRegistryStruct

func init() {
	PrometheusRegistry = &PrometheusRegistryStruct{}
	PrometheusRegistry.Client = prom.NewRegistry()
}


func (m *PrometheusRegistryStruct) MustRegister(cs ...prom.Collector)  {
	m.Client.MustRegister(cs...)
}

func GetPromHttp(port string) *http.Server {
	return &http.Server{Handler:promhttp.HandlerFor(PrometheusRegistry.Client, promhttp.HandlerOpts{}), Addr: port}
}

func (m *PrometheusRegistryStruct) UnRegister(cs prom.Collector)  {
	m.Client.Unregister(cs)
}