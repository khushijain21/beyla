// Copyright Red Hat / IBM
// Copyright Grafana Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This implementation is a derivation of the code in
// https://github.com/netobserv/netobserv-ebpf-agent/tree/release-1.4

package beyla

import (
	"time"
)

type NetworkConfig struct {
	// Enable network observability.
	// Default value is false (disabled)
	Enable bool `yaml:"enable" env:"BEYLA_NETWORK_METRICS"`

	// AgentIP allows overriding the reported Agent IP address on each flow.
	AgentIP string `yaml:"agent_ip" env:"BEYLA_NETWORK_AGENT_IP"`

	// AgentIPIface specifies which interface should the agent pick the IP address from in order to
	// report it in the AgentIP field on each flow. Accepted values are: external (default), local,
	// or name:<interface name> (e.g. name:eth0).
	// If the AgentIP configuration property is set, this property has no effect.
	AgentIPIface string `yaml:"agent_ip_iface" env:"BEYLA_NETWORK_AGENT_IP_IFACE"`
	// AgentIPType specifies which type of IP address (IPv4 or IPv6 or any) should the agent report
	// in the AgentID field of each flow. Accepted values are: any (default), ipv4, ipv6.
	// If the AgentIP configuration property is set, this property has no effect.
	AgentIPType string `yaml:"agent_ip_type" env:"BEYLA_NETWORK_AGENT_IP_TYPE"`
	// Interfaces contains the interface names from where flows will be collected. If empty, the agent
	// will fetch all the interfaces in the system, excepting the ones listed in ExcludeInterfaces.
	// If an entry is enclosed by slashes (e.g. `/br-/`), it will match as regular expression,
	// otherwise it will be matched as a case-sensitive string.
	Interfaces []string `yaml:"interfaces" env:"BEYLA_NETWORK_INTERFACES" envSeparator:","`
	// ExcludeInterfaces contains the interface names that will be excluded from flow tracing. Default:
	// "lo" (loopback).
	// If an entry is enclosed by slashes (e.g. `/br-/`), it will match as regular expression,
	// otherwise it will be matched as a case-sensitive string.
	ExcludeInterfaces []string `yaml:"exclude_interfaces" env:"BEYLA_NETWORK_EXCLUDE_INTERFACES" envSeparator:","`
	// CacheMaxFlows specifies how many flows can be accumulated in the accounting cache before
	// being flushed for its later export. Default value is 5000.
	// Decrease it if you see the "received message larger than max" error in Beyla logs.
	CacheMaxFlows int `yaml:"cache_max_flows" env:"BEYLA_NETWORK_CACHE_MAX_FLOWS"`
	// CacheActiveTimeout specifies the maximum duration that flows are kept in the accounting
	// cache before being flushed for its later export.
	CacheActiveTimeout time.Duration `yaml:"cache_active_timeout" env:"BEYLA_NETWORK_CACHE_ACTIVE_TIMEOUT"`
	// Deduper specifies the deduper type. Accepted values are "none" (disabled) and "firstCome".
	// When enabled, it will detect duplicate flows (flows that have been detected e.g. through
	// both the physical and a virtual interface).
	// "firstCome" will forward only flows from the first interface the flows are received from.
	Deduper string `yaml:"deduper" env:"BEYLA_NETWORK_DEDUPER"`
	// DeduperFCExpiry specifies the expiry duration of the flows "firstCome" deduplicator. After
	// a flow hasn't been received for that expiry time, the deduplicator forgets it. That means
	// that a flow from a connection that has been inactive during that period could be forwarded
	// again from a different interface.
	// If the value is not set, it will default to 2 * CacheActiveTimeout
	DeduperFCExpiry time.Duration `yaml:"deduper_fc_expiry" env:"BEYLA_NETWORK_DEDUPER_FC_EXPIRY"`
	// DeduperJustMark will just mark duplicates (boolean field) instead of dropping them. Default: false.
	DeduperJustMark bool `yaml:"deduper_just_mark" env:"BEYLA_NETWORK_DEDUPER_JUST_MARK"`
	// Direction allows selecting which flows to trace according to its direction. Accepted values
	// are "ingress", "egress" or "both" (default).
	Direction string `yaml:"direction" env:"BEYLA_NETWORK_DIRECTION"`
	// Sampling holds the rate at which packets should be sampled and sent to the target collector.
	// E.g. if set to 100, one out of 100 packets, on average, will be sent to the target collector.
	Sampling int `yaml:"sampling" env:"BEYLA_NETWORK_SAMPLING"`
	// ListenInterfaces specifies the mechanism used by the agent to listen for added or removed
	// network interfaces. Accepted values are "watch" (default) or "poll".
	// If the value is "watch", interfaces are traced immediately after they are created. This is
	// the recommended setting for most configurations. "poll" value is a fallback mechanism that
	// periodically queries the current network interfaces (frequency specified by ListenPollPeriod).
	ListenInterfaces string `yaml:"listen_interfaces" env:"BEYLA_NETWORK_LISTEN_INTERFACES"`
	// ListenPollPeriod specifies the periodicity to query the network interfaces when the
	// ListenInterfaces value is set to "poll".
	ListenPollPeriod time.Duration `yaml:"listen_poll_period" env:"BEYLA_NETWORK_LISTEN_POLL_PERIOD"`
}

var defaultNetworkConfig = NetworkConfig{
	AgentIPIface:       "external",
	AgentIPType:        "any",
	ExcludeInterfaces:  []string{"lo"},
	CacheMaxFlows:      5000,
	CacheActiveTimeout: 5 * time.Second,
	Deduper:            "none",
	DeduperJustMark:    false,
	Direction:          "both",
	ListenInterfaces:   "watch",
	ListenPollPeriod:   10 * time.Second,
}
