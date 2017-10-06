package cadvisor

import (
	"fmt"
	"testing"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

func TestGetMetricTypes(t *testing.T) {
	c := NewCollector()
	metrics, err := c.GetMetricTypes(plugin.Config{})
	if err != nil {
		fmt.Println(err)
	}
	for _, m := range metrics {
		fmt.Println(m.Namespace.String())
	}
}
