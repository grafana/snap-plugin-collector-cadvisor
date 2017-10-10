package cadvisor

import (
	"log"
	"time"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// Manifest maintains a list of metrics that must be gathered
// by the colllector
type Manifest struct {
	tcpMetrics    []string
	tcp6Metrics   []string
	cpuMetrics    []string
	loadMetrics   []string
	ifaceMetrics  []string
	fsMetrics     []string
	diskIoMetrics []string
	memMetrics    []string
}

func (m *Manifest) buildMetricsList(metrics []plugin.Metric) time.Duration {
	m.tcpMetrics = []string{}
	m.tcp6Metrics = []string{}
	m.cpuMetrics = []string{}
	m.loadMetrics = []string{}
	m.ifaceMetrics = []string{}
	m.fsMetrics = []string{}
	m.memMetrics = []string{}
	m.diskIoMetrics = []string{}
	intervalVal, err := metrics[0].Config.GetInt("interval")
	var interval time.Duration
	if err != nil {
		interval = time.Second * 15
	} else {
		interval = time.Second * time.Duration(intervalVal)
	}
	for _, mtx := range metrics {
		if mtx.Namespace.Element(6).Value == "tcp" {
			m.tcpMetrics = append(m.tcpMetrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "tcp6" {
			m.tcp6Metrics = append(m.tcp6Metrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "cpu" {
			m.cpuMetrics = append(m.cpuMetrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "load" {
			m.loadMetrics = append(m.loadMetrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "iface" {
			m.ifaceMetrics = append(m.ifaceMetrics, mtx.Namespace.Element(8).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "fs" {
			m.fsMetrics = append(m.fsMetrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "mem" {
			m.memMetrics = append(m.memMetrics, mtx.Namespace.Element(7).Value)
			continue
		}
		if mtx.Namespace.Element(6).Value == "diskio" {
			m.diskIoMetrics = append(m.diskIoMetrics, mtx.Namespace.Element(8).Value)
			continue
		}
		log.Printf("metric %v not found but requested\n", mtx.Namespace.String())
	}
	return interval
}
