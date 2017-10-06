package main

import (
	"github.com/grafana/snap-plugin-collector-cadvisor/cadvisor"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func main() {
	plugin.StartStreamCollector(cadvisor.NewCollector(), cadvisor.PluginName, cadvisor.PluginVersion, plugin.Exclusive(true))
}
