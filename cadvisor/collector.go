package cadvisor

import (
	"context"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/google/cadvisor/cache/memory"
	"github.com/google/cadvisor/container"
	info "github.com/google/cadvisor/info/v2"
	"github.com/google/cadvisor/manager"
	"github.com/google/cadvisor/utils/sysfs"

	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// Plugin specifc constants
const (
	PluginName    = "cadvisor"
	PluginVendor  = "grafanalabs"
	PluginVersion = 1
)

// kubernetes label constants,
// !Importing from kubernetes adds overhead
const (
	KubernetesPodNameLabel       = "io.kubernetes.pod.name"
	KubernetesPodNamespaceLabel  = "io.kubernetes.pod.namespace"
	KubernetesContainerNameLabel = "io.kubernetes.container.name"
)

var (
	storageDuration             = 1 * time.Minute                                                   // How long to keep data stored
	defaultHousekeepingInterval = 10 * time.Second                                                  // Interval for cadvisor to perform housekeeping, effects cpu usage
	maxHousekeepingInterval     = 60 * time.Second                                                  // Largest interval to allow between container housekeepings
	allowDynamicHousekeeping    = true                                                              // Whether to allow the housekeeping interval to be dynamic
	ignoreMetrics               = container.MetricSet{container.NetworkUdpUsageMetrics: struct{}{}} // Cadvisor will ignore udp metrics
	_                           = flag.CommandLine.Parse([]string{})                                // Removes noise output from glog imported by cAdvisor
)

// Collector contains the components to collect cadvisor metrics
type Collector struct {
	mng      manager.Manager
	manifest Manifest
	lock     *sync.Mutex
	interval time.Duration
}

func init() {
	// Override cAdvisor flag defaults.
	flagOverrides := map[string]string{
		// Override the default cAdvisor housekeeping interval.
		"housekeeping_interval": defaultHousekeepingInterval.String(),
		// Disable event storage by default.
		"event_storage_event_limit": "default=0",
		"event_storage_age_limit":   "default=0",
	}
	for name, defaultValue := range flagOverrides {
		if f := flag.Lookup(name); f != nil {
			f.DefValue = defaultValue
			f.Value.Set(defaultValue)
		} else {
			glog.Errorf("Expected cAdvisor flag %q not found", name)
		}
	}
}

// buildOrganizer waits on updates from the active task manifest on which metrics to collect
func (c *Collector) buildOrganizer(in chan []plugin.Metric) {
	for {
		newMetrics := <-in
		c.lock.Lock()
		c.interval = c.manifest.buildMetricsList(newMetrics)
		c.lock.Unlock()
	}
}

// StreamMetrics takes both an in and out channel of []plugin.Metric
//
// The mtxIn channel is used to set/update the metrics that Snap is
// currently requesting to be collected by the plugin.
//
// The mtxOut channel is used by the plugin to send the collected metrics
// to Snap.
func (c *Collector) StreamMetrics(ctx context.Context, mtxIn chan []plugin.Metric, mtxOut chan []plugin.Metric, chanErr chan string) error {
	go c.buildOrganizer(mtxIn)
	var contInfo [3]string
	var ok bool
	var err error
	c.mng, err = manager.New(memory.New(storageDuration, nil), sysfs.NewRealSysFs(), maxHousekeepingInterval, allowDynamicHousekeeping, ignoreMetrics, http.DefaultClient)
	if err != nil {
		log.Fatalf("Failed to create a Container Manager: %v", err)
		chanErr <- err.Error()
	}
	// Start the manager.
	if err := c.mng.Start(); err != nil {
		log.Fatalf("Failed to start container manager: %v", err)
		chanErr <- err.Error()
	}

	for {
		c.lock.Lock()
		containers, err := c.mng.GetContainerInfoV2("/", info.RequestOptions{Count: 1, Recursive: true, IdType: info.TypeName})
		if err != nil {
			log.Printf("unable to gather container metrics: %v", err)
		}
		metrics := []plugin.Metric{}
		for _, cont := range containers {
			if len(cont.Stats) < 1 {
				log.Printf("no container stats currently available")
				continue
			}
			if contInfo, ok = checkContainer(cont.Spec.Labels); !ok {
				continue
			}
			if cont.Spec.HasNetwork {
				for _, key := range c.manifest.tcpMetrics {
					m, ok := tcpMap[key]
					if !ok {
						log.Printf("metric: %v does not exist in the tcp metric map\n", key)
						continue
					}
					metrics = append(metrics, plugin.Metric{
						Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2]),
						Description: m.Description,
						Unit:        m.Unit,
						Data:        m.Data(cont.Stats[0]),
						Timestamp:   cont.Stats[0].Timestamp,
					})
				}
				for _, key := range c.manifest.tcp6Metrics {
					m, ok := tcp6Map[key]
					if !ok {
						log.Printf("metric: %v does not exist in the tcp6 metric map\n", key)
						continue
					}
					metrics = append(metrics, plugin.Metric{
						Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2]),
						Description: m.Description,
						Unit:        m.Unit,
						Data:        m.Data(cont.Stats[0]),
						Timestamp:   cont.Stats[0].Timestamp,
					})
				}
				for _, key := range c.manifest.ifaceMetrics {
					m, ok := ifaceMap[key]
					if !ok {
						log.Printf("metric: %v does not exist in the tcp6 metric map\n", key)
						continue
					}
					for _, iface := range cont.Stats[0].Network.Interfaces {
						metrics = append(metrics, plugin.Metric{
							Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2], iface.Name),
							Description: m.Description,
							Unit:        m.Unit,
							Data:        m.Data(iface),
							Timestamp:   cont.Stats[0].Timestamp,
						})
					}
				}
			}

			if cont.Spec.HasMemory {
				for _, key := range c.manifest.memMetrics {
					m, ok := memMap[key]
					if !ok {
						log.Printf("metric: %v does not exist in the mem metric map\n", key)
						continue
					}
					metrics = append(metrics, plugin.Metric{
						Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2]),
						Description: m.Description,
						Unit:        m.Unit,
						Data:        m.Data(cont.Stats[0]),
						Timestamp:   cont.Stats[0].Timestamp,
					})
				}
			}

			if cont.Spec.HasCpu {
				for _, key := range c.manifest.cpuMetrics {
					m, ok := cpuMap[key]
					if !ok {
						log.Printf("metric: %v does not exist in the cpu metric map\n", key)
						continue
					}
					metrics = append(metrics, plugin.Metric{
						Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2]),
						Description: m.Description,
						Unit:        m.Unit,
						Data:        m.Data(cont.Stats[0]),
						Timestamp:   cont.Stats[0].Timestamp,
					})
				}
			}

			if cont.Spec.HasFilesystem {
				for _, key := range c.manifest.fsMetrics {
					m, ok := fsMap[key]
					if !ok {
						log.Printf("metric: %v does not exist in the fs metric map\n", key)
						continue
					}
					metrics = append(metrics, plugin.Metric{
						Namespace:   m.Namespace(contInfo[0], contInfo[1], contInfo[2]),
						Description: m.Description,
						Unit:        m.Unit,
						Data:        m.Data(cont.Stats[0]),
						Timestamp:   cont.Stats[0].Timestamp,
					})
				}
			}
		}
		c.lock.Unlock()
		mtxOut <- metrics
		time.Sleep(c.interval)
	}
}

func checkContainer(labels map[string]string) ([3]string, bool) {
	var podName, nameSpace, containerName string
	var ok bool
	if podName, ok = labels[KubernetesPodNameLabel]; !ok {
		return [3]string{}, false
	}
	if nameSpace, ok = labels[KubernetesPodNamespaceLabel]; !ok {
		return [3]string{}, false
	}
	if containerName, ok = labels[KubernetesContainerNameLabel]; !ok {
		return [3]string{}, false
	}
	return [3]string{nameSpace, podName, containerName}, true
}

// GetMetricTypes will be called when your plugin is loaded in order to populate the metric catalog(where snaps stores all
// available metrics). Config info is passed in. This config information would come from global config snap settings.
// The metrics returned will be advertised to users who list all the metrics and will become targetable by tasks.
func (c Collector) GetMetricTypes(cfg plugin.Config) ([]plugin.Metric, error) {
	metrics := []plugin.Metric{}

	for _, m := range tcpMap {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}

	for _, m := range tcp6Map {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}

	for _, m := range fsMap {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}

	for _, m := range cpuMap {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}

	for _, m := range memMap {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}

	for _, m := range ifaceMap {
		metrics = append(metrics, plugin.Metric{
			Namespace:   m.Namespace("*", "*", "*", "*"),
			Description: m.Description,
			Unit:        m.Unit,
			Config:      cfg,
		})
	}
	return metrics, nil
}

// GetConfigPolicy returns the configPolicy for your plugin.
// A config policy is how users can provide configuration info to
// plugin. Here you define what sorts of config info your plugin
// needs and/or requires.
func (Collector) GetConfigPolicy() (plugin.ConfigPolicy, error) {
	policy := plugin.NewConfigPolicy()
	policy.AddNewIntRule([]string{PluginVendor, PluginName}, "interval", false, plugin.SetDefaultInt(15), plugin.SetMinInt(1))
	return *policy, nil
}

// NewCollector returns a new active cadvisor collector
func NewCollector() *Collector {
	return &Collector{
		lock:     &sync.Mutex{},
		manifest: Manifest{},
		interval: time.Second * 15,
	}
}
