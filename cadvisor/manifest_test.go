package cadvisor

import (
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

var (
	metricList1 = []plugin.Metric{
		plugin.Metric{
			Namespace: plugin.NewNamespace(PluginVendor, PluginName, "container", "*", "*", "*", "tcp", "ESTABLISHED"),
		},
	}
)

func TestBuildMetricsList(t *testing.T) {
	newManifest := Manifest{}
	newManifest.buildMetricsList(metricList1)
	if newManifest.tcpMetrics[0] != "ESTABLISHED" {
		t.Errorf("ESTABLISHED not added to TCP Metrics")
		t.Fail()
	}
}
