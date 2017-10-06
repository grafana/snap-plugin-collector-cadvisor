# Snap collector plugin - cadvisor
This plugin collects metrics from cadvisor which gathers information on running processes and system utilization (CPU, memory, disks, network). It is designed with kubernetes in mind and only collects container metrics on containers with kubernetes labels

It's used in the [Snap framework](http://github.com:intelsdi-x/snap).

## Getting Started
### System Requirements 
* [golang 1.9+](https://golang.org/dl/) (needed only for building)

### Operating systems
All OSs currently supported by snap:
* Linux/amd64
* Darwin/amd64

### Installation
#### Download cadvisor plugin binary:
You can get the pre-built binaries for your OS and architecture under the plugin's [release](https://github.com/grafana/snap-plugin-collector-cadvisor/releases) page.  For Snap, check [here](https://github.com/intelsdi-x/snap/releases).


#### To build the plugin binary:
Fork https://github.com/grafana/snap-plugin-collector-cadvisor

Clone repo into `$GOPATH/src/github.com/grafana/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-cadvisor.git
```

Build the plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started)

Available configuration option:
* interval - this is a streaming plugin that requires a set interval for how often to forward metrics from cadvisors. This is a positive integer

## Documentation
There are a number of other resources you can review to learn to use this plugin:

* [gocadvisor](https://github.com/shirou/gocadvisor/) (go based implementation)
* [cadvisor](https://pythonhosted.org/cadvisor/) (python based implementation)
* [Snap cadvisor integration test](https://github.com/grafana/snap-plugin-collector-cadvisor/blob/master/cadvisor/cadvisor_integration_test.go)
* [Snap cadvisor unit test](https://github.com/grafana/snap-plugin-collector-cadvisor/blob/master/cadvisor/cadvisor_test.go)
* [Snap cadvisor examples](#examples)

### Collected metrics
List of metrics collected by this plugin can be found in [METRICS.md file](METRICS.md).
