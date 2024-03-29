package cadvisor

import (
	"github.com/google/cadvisor/info/v1"
	info "github.com/google/cadvisor/info/v2"
	"github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

// Metric type to translate v2.ContainerInfo into a snap Metric
type Metric struct {
	Namespace   func(ns string, pn string, cn string) plugin.Namespace
	Description string
	Unit        string
	Tags        map[string]string
	Data        func(s *info.ContainerStats) interface{}
}

// IfaceMetric type to translate v1.InterfaceStats into a snap Metric
type IfaceMetric struct {
	Namespace   func(ns string, pn string, cn string, name string) plugin.Namespace
	Description string
	Unit        string
	Tags        map[string]string
	Data        func(s v1.InterfaceStats) interface{}
}

// DiskIoMetric type to translate v1.InterfaceStats into a snap Metric
type DiskIoMetric struct {
	Namespace   func(ns string, pn string, cn string, name string) plugin.Namespace
	Description string
	Unit        string
	Tags        map[string]string
	Data        func(s v1.PerDiskStats) interface{}
}

func containerNamespace(ns string, pn string, cn string) plugin.Namespace {
	return plugin.Namespace{
		plugin.NamespaceElement{
			Value: PluginVendor,
		},
		plugin.NamespaceElement{
			Value: PluginName,
		},
		plugin.NamespaceElement{
			Value: "container",
		},
		plugin.NamespaceElement{
			Name:  "namespace",
			Value: ns,
		},
		plugin.NamespaceElement{
			Name:  "pod_name",
			Value: pn,
		},
		plugin.NamespaceElement{
			Name:  "container_name",
			Value: cn,
		},
	}
}

var (
	cpuMap = map[string]Metric{
		"total": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("cpu", "total", "usage")
			},
			Unit:        "ns",
			Description: "total CPU usage",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Cpu.Usage.Total
			},
		},
		"user": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("cpu", "user", "usage")
			},
			Unit:        "ns",
			Description: "user CPU usage",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Cpu.Usage.User
			},
		},
		"system": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("cpu", "system", "usage")
			},
			Unit:        "ns",
			Description: "system CPU usage",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Cpu.Usage.System
			},
		},
		"load": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("cpu", "load")
			},
			Unit:        "load",
			Description: " Load is smoothed over the last 10 seconds. Instantaneous value can be read",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Cpu.LoadAverage
			},
		},
	}

	tcpMap = map[string]Metric{
		"ESTABLISHED": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "ESTABLISHED")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'ESTABLISHED'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.Established
			},
		},
		"SYN_SENT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "SYN_SENT")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'SYN_SENT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.SynSent
			},
		},
		"SYN_RECV": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "SYN_RECV")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'SYN_RECV'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.SynRecv
			},
		},
		"FIN_WAIT_1": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "FIN_WAIT_1")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'FIN_WAIT_1'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.FinWait1
			},
		},
		"FIN_WAIT_2": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "FIN_WAIT_2")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'FIN_WAIT_2'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.FinWait2
			},
		},
		"TIME_WAIT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "TIME_WAIT")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'TIME_WAIT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.TimeWait
			},
		},
		"CLOSE": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "CLOSE")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'CLOSE'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.Close
			},
		},
		"CLOSE_WAIT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "CLOSE_WAIT")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'CLOSE_WAIT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.CloseWait
			},
		},
		"LAST_ACK": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "LAST_ACK")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'LAST_ACK'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.LastAck
			},
		},
		"LISTEN": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "LISTEN")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'LISTEN'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.Listen
			},
		},
		"CLOSING": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp", "CLOSING")
			},
			Unit:        "event",
			Description: "Count of TCP connections in state 'CLOSING'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp.Closing
			},
		},
	}

	tcp6Map = map[string]Metric{
		"ESTABLISHED": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "ESTABLISHED")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'ESTABLISHED'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.Established
			},
		},
		"SYN_SENT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "SYN_SENT")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'SYN_SENT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.SynSent
			},
		},
		"SYN_RECV": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "SYN_RECV")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'SYN_RECV'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.SynRecv
			},
		},
		"FIN_WAIT_1": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "FIN_WAIT_1")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'FIN_WAIT_1'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.FinWait1
			},
		},
		"FIN_WAIT_2": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "FIN_WAIT_2")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'FIN_WAIT_2'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.FinWait2
			},
		},
		"TIME_WAIT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "TIME_WAIT")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'TIME_WAIT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.TimeWait
			},
		},
		"CLOSE": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "CLOSE")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'CLOSE'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.Close
			},
		},
		"CLOSE_WAIT": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "CLOSE_WAIT")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'CLOSE_WAIT'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.CloseWait
			},
		},
		"LAST_ACK": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "LAST_ACK")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'LAST_ACK'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.LastAck
			},
		},
		"LISTEN": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "LISTEN")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'LISTEN'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.Listen
			},
		},
		"CLOSING": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("tcp6", "CLOSING")
			},
			Unit:        "event",
			Description: "Count of TCP6 connections in state 'CLOSING'",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Network.Tcp6.Closing
			},
		},
	}

	memMap = map[string]Metric{
		"cache": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "cache")
			},
			Unit:        "B",
			Description: "Number of bytes of page cache memory.",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.Cache
			},
		},
		"usage": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "usage")
			},
			Unit:        "B",
			Description: "Current memory usage, this includes all memory regardless of when it was accessed.",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.Usage
			},
		},
		"rss": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "rss")
			},
			Unit:        "B",
			Description: "The amount of anonymous and swap cache memory (includes transparent hugepages)",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.RSS
			},
		},
		"swap": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "swap")
			},
			Unit:        "B",
			Description: "The amount of swap currently used by the processes in this cgroup",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.Swap
			},
		},
		"working_set": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "working_set")
			},
			Unit:        "B",
			Description: "The amount of working set memory, this includes recently accessed memory, dirty memory, and kernel memory.",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.WorkingSet
			},
		},
		"failcnt": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("mem", "failcnt")
			},
			Unit:        "B",
			Description: "",
			Data: func(s *info.ContainerStats) interface{} {
				return s.Memory.Failcnt
			},
		},
	}

	fsMap = map[string]Metric{
		"total_usage": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("fs", "total_usage")
			},
			Unit:        "B",
			Description: "Total Number of bytes consumed by container.",
			Data: func(s *info.ContainerStats) interface{} {
				return *s.Filesystem.TotalUsageBytes
			},
		},
		"base_usage": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("fs", "base_usage")
			},
			Unit:        "B",
			Description: "Total Number of bytes consumed by container.",
			Data: func(s *info.ContainerStats) interface{} {
				return *s.Filesystem.BaseUsageBytes
			},
		},
		"inode_usage": Metric{
			Namespace: func(ns string, pn string, cn string) plugin.Namespace {
				return containerNamespace(ns, pn, cn).AddStaticElements("fs", "inode_usage")
			},
			Unit:        "inodes",
			Description: "Number of inodes used within the container's root filesystem.",
			Data: func(s *info.ContainerStats) interface{} {
				return *s.Filesystem.InodeUsage
			},
		},
	}

	diskIoMap = map[string]DiskIoMetric{
		"read_bytes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("read_bytes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "B",
			Description: "Total Number of bytes read",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"reads": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("reads")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of reads completed",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"queued_reads": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("queued_reads")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total Number of reads queued",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"sector_reads": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("sector_reads")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of sector reads completed",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"merged_reads": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("merged_reads")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of reads merged",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"read_time": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("read_time")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "ns",
			Description: "Total number of reads completed",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Read"]
			},
		},
		"write_bytes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("write_bytes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "B",
			Description: "Total Number of bytes write",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
		"writes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("writes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of writes completed",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
		"queued_writes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("queued_writes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total Number of writes queued",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
		"sector_writes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("sector_writes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of sector writes completed",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
		"merged_writes": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("merged_writes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "event",
			Description: "Total number of writes merged",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
		"write_time": DiskIoMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("diskio").AddDynamicElement("device_name", "name of the disk").AddStaticElement("write_time")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "ns",
			Description: "Total amount of time spent writing",
			Data: func(s v1.PerDiskStats) interface{} {
				return s.Stats["Write"]
			},
		},
	}

	ifaceMap = map[string]IfaceMetric{
		"in_bytes": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("in_bytes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "B",
			Description: "Cumulative count of bytes received",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.RxBytes
			},
		},
		"in_packets": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("in_packets")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of packets received",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.RxPackets
			},
		},
		"in_errors": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("in_errors")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of errors received by the container",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.RxErrors
			},
		},
		"in_dropped": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("in_dropped")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of bytes received by the container",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.RxDropped
			},
		},
		"out_bytes": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("out_bytes")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "B",
			Description: "Cumulative count of bytes transmitted",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.TxBytes
			},
		},
		"out_packets": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("out_packets")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of packets transmitted",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.TxPackets
			},
		},
		"out_errors": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("out_errors")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of errors transmitted by the container",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.TxErrors
			},
		},
		"out_dropped": IfaceMetric{
			Namespace: func(ns string, pn string, cn string, name string) plugin.Namespace {
				metName := containerNamespace(ns, pn, cn).AddStaticElement("iface").AddDynamicElement("device_name", "name of the interface").AddStaticElement("out_dropped")
				if name != "*" {
					metName[7].Value = name
				}
				return metName
			},
			Unit:        "pckt",
			Description: "Cumulative count of bytes transmitted by the container",
			Data: func(s v1.InterfaceStats) interface{} {
				return s.TxDropped
			},
		},
	}
)
