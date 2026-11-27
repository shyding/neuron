package dataengine

import (
	"fmt"
	"sync"
	"time"
)

// Engine is the built-in data collection engine (replaces external neuron process)
type Engine struct {
	mu            sync.RWMutex
	nodes         map[string]*Node
	plugins       []PluginInfo
	subscriptions []map[string]interface{}
}

// Node represents a configured device node
type Node struct {
	Name    string                 `json:"name"`
	Plugin  string                 `json:"plugin"`
	Type    int                    `json:"type"` // 1=south, 2=north
	State   int                    `json:"running"` // 1=init, 2=ready, 3=running, 4=stopped
	Link    int                    `json:"link"`
	Groups  map[string]*Group      `json:"-"`
	Setting map[string]interface{} `json:"-"`
}

// Group represents a tag group within a node
type Group struct {
	Name     string                   `json:"name"`
	Interval int                      `json:"interval"`
	Tags     []map[string]interface{} `json:"tags"`
}

// PluginInfo describes a plugin
type PluginInfo struct {
	Name        string `json:"name"`
	Library     string `json:"library"`
	Description string `json:"description"`
	Type        int    `json:"type"` // 1=south, 2=north
	Kind        int    `json:"kind"` // 0=static, 1=system
	NodeType    int    `json:"node_type"`
	DisplayName string `json:"display_name"`
	Schema      string `json:"schema"`
}

// NewEngine creates a new data engine with pre-registered plugins
func NewEngine() *Engine {
	e := &Engine{
		nodes: make(map[string]*Node),
	}
	e.registerBuiltinPlugins()
	return e
}

// registerBuiltinPlugins registers all 55+ southbound and northbound driver definitions
func (e *Engine) registerBuiltinPlugins() {
	// === Southbound Drivers (type=1) ===
	southDrivers := []struct{ name, display, desc string }{
		{"Modbus TCP", "Modbus TCP", "Modbus TCP protocol plugin"},
		{"Modbus RTU", "Modbus RTU", "Modbus RTU protocol plugin"},
		{"OPC UA", "OPC UA", "OPC UA protocol plugin"},
		{"OPC DA", "OPC DA", "OPC DA protocol plugin"},
		{"Siemens S7 ISOTCP", "Siemens S7 ISO TCP", "Siemens S7-200/300/400/1200/1500 ISO TCP plugin"},
		{"Mitsubishi 3E", "Mitsubishi MC 3E", "Mitsubishi MC Protocol 3E frame plugin"},
		{"Mitsubishi 1E", "Mitsubishi MC 1E", "Mitsubishi MC Protocol 1E frame plugin"},
		{"Mitsubishi FX", "Mitsubishi FX", "Mitsubishi FX series serial plugin"},
		{"Omron FINS TCP", "Omron FINS TCP", "Omron FINS TCP plugin"},
		{"Omron FINS UDP", "Omron FINS UDP", "Omron FINS UDP plugin"},
		{"Beckhoff ADS", "Beckhoff ADS", "Beckhoff ADS/AMS protocol plugin"},
		{"BACnet/IP", "BACnet/IP", "BACnet/IP protocol plugin"},
		{"KNXnet/IP", "KNXnet/IP", "KNX building automation plugin"},
		{"IEC 60870-5-104", "IEC 60870-5-104", "IEC 60870-5-104 telecontrol plugin"},
		{"IEC 61850", "IEC 61850", "IEC 61850 MMS plugin"},
		{"DL/T 645-2007", "DL/T 645-2007", "DL/T 645-2007 electric meter plugin"},
		{"DL/T 645-1997", "DL/T 645-1997", "DL/T 645-1997 electric meter plugin"},
		{"EtherNet/IP CIP", "EtherNet/IP CIP", "Allen-Bradley EtherNet/IP CIP plugin"},
		{"Allen-Bradley DF1", "Allen-Bradley DF1", "Allen-Bradley DF1 serial plugin"},
		{"Profinet IO", "PROFINET IO", "PROFINET IO device plugin"},
		{"Fanuc Focas", "FANUC Focas", "FANUC CNC Focas plugin"},
		{"Heidenhain", "Heidenhain", "Heidenhain CNC DNC plugin"},
		{"ABB COMLI", "ABB COMLI", "ABB COMLI serial plugin"},
		{"Delta DVP", "Delta DVP", "Delta DVP PLC plugin"},
		{"Keyence CIP", "Keyence CIP", "Keyence PLC CIP plugin"},
		{"Panasonic Mewtocol", "Panasonic Mewtocol", "Panasonic FP series Mewtocol plugin"},
		{"GE SRTP", "GE SRTP", "GE Fanuc SRTP plugin"},
		{"Schneider Modbus", "Schneider Modbus", "Schneider Modbus plugin"},
		{"Honeywell", "Honeywell", "Honeywell Experion plugin"},
		{"Yokogawa", "Yokogawa", "Yokogawa FA-M3 plugin"},
		{"HollySys LK", "HollySys LK", "HollySys LK PLC plugin"},
		{"HollySys MACS", "HollySys MACS", "HollySys MACS DCS plugin"},
		{"Inovance Modbus", "Inovance Modbus", "Inovance Modbus TCP plugin"},
		{"XINJE Modbus", "XINJE Modbus", "XINJE XC/XD/XL PLC Modbus plugin"},
		{"CODESYS V3", "CODESYS V3", "CODESYS V3 Ethernet plugin"},
		{"Non A11", "Non A11", "NON-A11 environment sensor plugin"},
		{"NOA", "NOA", "NOA cloud gateway plugin"},
		{"MTConnect", "MTConnect", "MTConnect CNC monitoring plugin"},
		{"Modbus TCP QH", "Modbus TCP QH", "Modbus TCP Quantum-H variant plugin"},
		{"Modbus RTU QH", "Modbus RTU QH", "Modbus RTU Quantum-H variant plugin"},
		{"Hostlink FINS", "Hostlink FINS", "Omron Hostlink FINS plugin"},
		{"Hostlink Cmode", "Hostlink Cmode", "Omron Hostlink C-mode plugin"},
		{"ABB COMLI TCP", "ABB COMLI TCP", "ABB COMLI TCP plugin"},
		{"Siemens MPI", "Siemens MPI", "Siemens S7 MPI/PPI plugin"},
		{"Siemens Fetch Write", "Siemens Fetch Write", "Siemens S7 Fetch/Write plugin"},
		{"Modbus ASCII", "Modbus ASCII", "Modbus ASCII serial plugin"},
		{"CANopen", "CANopen", "CANopen plugin"},
		{"EtherCAT", "EtherCAT", "EtherCAT plugin"},
		{"CC-Link", "CC-Link", "CC-Link plugin"},
		{"DNP3", "DNP3", "DNP3 (Distributed Network Protocol 3) plugin"},
		{"Profibus DP", "Profibus DP", "Profibus DP plugin"},
		{"HART-IP", "HART-IP", "HART-IP plugin"},
		{"File", "File", "File data source plugin"},
		{"Simulator", "Simulator", "Neuron data simulator plugin"},
		{"Custom", "Custom", "Custom protocol plugin"},
	}

	for _, d := range southDrivers {
		e.plugins = append(e.plugins, PluginInfo{
			Name:        d.name,
			Library:     fmt.Sprintf("libplugin-%s.so", sanitizeName(d.name)),
			Description: d.desc,
			Type:        1,
			Kind:        0,
			NodeType:    1,
			DisplayName: d.display,
		})
	}

	// === Northbound Drivers (type=2) ===
	northDrivers := []struct{ name, display, desc string }{
		{"MQTT", "MQTT", "MQTT northbound data forwarding plugin"},
		{"SparkplugB", "SparkPlugB", "SparkPlug B specification plugin"},
		{"Monitor", "Monitor", "NeuronEX data monitoring plugin"},
		{"eKuiper", "eKuiper", "eKuiper stream processing integration plugin"},
		{"WebSocket", "WebSocket", "WebSocket data push plugin"},
		{"Kafka", "Kafka", "Apache Kafka producer plugin"},
		{"Azure IoT", "Azure IoT", "Azure IoT Hub plugin"},
		{"AWS IoT", "AWS IoT Core", "AWS IoT Core plugin"},
	}

	for _, d := range northDrivers {
		e.plugins = append(e.plugins, PluginInfo{
			Name:        d.name,
			Library:     fmt.Sprintf("libplugin-%s.so", sanitizeName(d.name)),
			Description: d.desc,
			Type:        2,
			Kind:        0,
			NodeType:    2,
			DisplayName: d.display,
		})
	}
}

// === Node Management ===

func (e *Engine) AddNode(body map[string]interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	name, _ := body["name"].(string)
	plugin, _ := body["plugin"].(string)
	nodeType := 1
	if t, ok := body["type"].(float64); ok {
		nodeType = int(t)
	}

	if name == "" {
		return fmt.Errorf("node name required")
	}
	if _, exists := e.nodes[name]; exists {
		return fmt.Errorf("node %s already exists", name)
	}

	e.nodes[name] = &Node{
		Name:   name,
		Plugin: plugin,
		Type:   nodeType,
		State:  1, // init
		Link:   0,
		Groups: make(map[string]*Group),
		Setting: make(map[string]interface{}),
	}
	return nil
}

func (e *Engine) DeleteNode(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.nodes, name)
}

func (e *Engine) GetNodes(nodeType string) []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var result []map[string]interface{}
	for _, n := range e.nodes {
		if nodeType != "" {
			t := 0
			if nodeType == "1" {
				t = 1
			} else if nodeType == "2" {
				t = 2
			}
			if t != 0 && n.Type != t {
				continue
			}
		}
		result = append(result, map[string]interface{}{
			"name":   n.Name,
			"plugin": n.Plugin,
			"type":   n.Type,
		})
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result
}

func (e *Engine) GetNodeStates() []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var states []map[string]interface{}
	for _, n := range e.nodes {
		subGroupCount := len(n.Groups)
		states = append(states, map[string]interface{}{
			"node":            n.Name,
			"running":         n.State,
			"link":            n.Link,
			"rtt":             0,
			"log_level":       "notice",
			"sub_group_count": subGroupCount,
		})
	}
	if states == nil {
		states = []map[string]interface{}{}
	}
	return states
}

func (e *Engine) GetNodeSetting(node string) map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if n, ok := e.nodes[node]; ok && n.Setting != nil {
		return n.Setting
	}
	return map[string]interface{}{}
}

func (e *Engine) SetNodeSetting(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, _ := body["node"].(string)
	if n, ok := e.nodes[node]; ok {
		for k, v := range body {
			if k != "node" {
				n.Setting[k] = v
			}
		}
	}
}

func (e *Engine) CtlNode(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, _ := body["node"].(string)
	cmd, _ := body["cmd"].(float64)
	if n, ok := e.nodes[node]; ok {
		switch int(cmd) {
		case 0: // start
			n.State = 3
		case 1: // stop
			n.State = 4
		}
	}
}

// === Group Management ===

func (e *Engine) AddGroup(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, _ := body["node"].(string)
	group, _ := body["group"].(string)
	interval := 1000
	if i, ok := body["interval"].(float64); ok {
		interval = int(i)
	}

	if n, ok := e.nodes[node]; ok {
		n.Groups[group] = &Group{
			Name:     group,
			Interval: interval,
			Tags:     []map[string]interface{}{},
		}
	}
}

func (e *Engine) DeleteGroup(node, group string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if n, ok := e.nodes[node]; ok {
		delete(n.Groups, group)
	}
}

func (e *Engine) GetGroups(node string) []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var result []map[string]interface{}
	if n, ok := e.nodes[node]; ok {
		for _, g := range n.Groups {
			result = append(result, map[string]interface{}{
				"name":      g.Name,
				"interval":  g.Interval,
				"tag_count": len(g.Tags),
			})
		}
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result
}

// === Tag Management ===

func (e *Engine) AddTags(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, _ := body["node"].(string)
	group, _ := body["group"].(string)
	tags, _ := body["tags"].([]interface{})

	if n, ok := e.nodes[node]; ok {
		if g, ok := n.Groups[group]; ok {
			for _, t := range tags {
				if tag, ok := t.(map[string]interface{}); ok {
					g.Tags = append(g.Tags, tag)
				}
			}
		}
	}
}

func (e *Engine) DeleteTags(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	node, _ := body["node"].(string)
	group, _ := body["group"].(string)
	tagNames, _ := body["tags"].([]interface{})

	if n, ok := e.nodes[node]; ok {
		if g, ok := n.Groups[group]; ok {
			nameSet := make(map[string]bool)
			for _, t := range tagNames {
				if s, ok := t.(string); ok {
					nameSet[s] = true
				}
			}
			var remaining []map[string]interface{}
			for _, tag := range g.Tags {
				name, _ := tag["name"].(string)
				if !nameSet[name] {
					remaining = append(remaining, tag)
				}
			}
			g.Tags = remaining
		}
	}
}

func (e *Engine) GetTags(node, group string) []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if n, ok := e.nodes[node]; ok {
		if g, ok := n.Groups[group]; ok {
			if g.Tags != nil {
				return g.Tags
			}
		}
	}
	return []map[string]interface{}{}
}

// === Subscription Management ===

func (e *Engine) AddSubscription(body map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.subscriptions = append(e.subscriptions, body)
}

func (e *Engine) GetSubscriptions(app string) []map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if app == "" {
		if e.subscriptions == nil {
			return []map[string]interface{}{}
		}
		return e.subscriptions
	}

	var result []map[string]interface{}
	for _, s := range e.subscriptions {
		if a, ok := s["app"].(string); ok && a == app {
			result = append(result, s)
		}
	}
	if result == nil {
		result = []map[string]interface{}{}
	}
	return result
}

// === Plugin queries ===

func (e *Engine) GetPlugins(pluginType string) []map[string]interface{} {
	var result []map[string]interface{}
	t := 0
	if pluginType == "1" {
		t = 1
	} else if pluginType == "2" {
		t = 2
	}

	for _, p := range e.plugins {
		if t != 0 && p.Type != t {
			continue
		}
		result = append(result, map[string]interface{}{
			"name":         p.Name,
			"library":      p.Library,
			"description":  p.Description,
			"type":         p.Type,
			"kind":         p.Kind,
			"node_type":    p.NodeType,
			"display_name": p.DisplayName,
		})
	}
	return result
}

func (e *Engine) GetGlobalDrivers() []map[string]interface{} {
	return e.GetPlugins("1")
}

func (e *Engine) GetPluginSchema(pluginName string) map[string]interface{} {
	// Return a generic schema for any plugin
	return map[string]interface{}{
		"tag_type": map[string]interface{}{
			"1": "attribute",
			"2": "static",
			"3": "timestamp",
		},
		"group_interval": map[string]interface{}{
			"min":     100,
			"max":     60000,
			"default": 1000,
		},
		"setting": map[string]interface{}{
			"host": map[string]interface{}{
				"type":        "string",
				"description": "Device IP address",
				"default":     "127.0.0.1",
				"required":    true,
			},
			"port": map[string]interface{}{
				"type":        "int",
				"description": "Device port",
				"default":     502,
				"required":    true,
			},
		},
	}
}

// === Counts ===

func (e *Engine) CountTags() int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	count := 0
	for _, n := range e.nodes {
		for _, g := range n.Groups {
			count += len(g.Tags)
		}
	}
	return count
}

func (e *Engine) CountNodes() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.nodes)
}

// helper
func sanitizeName(name string) string {
	result := make([]byte, 0, len(name))
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' {
			result = append(result, byte(c))
		} else if c == ' ' || c == '/' {
			result = append(result, '-')
		}
	}
	return string(result)
}

// unused but required by import
var _ = time.Now
