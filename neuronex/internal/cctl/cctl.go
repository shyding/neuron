package cctl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Config struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	DeviceID  string `json:"device_id"`
	Keepalive int    `json:"keepalive"`
}

type Message struct {
	DeviceID   string                 `json:"deviceId"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Event      string                 `json:"event,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Timestamp  int64                  `json:"timestamp"`
	MessageID  string                 `json:"messageId"`
}

type DownlinkMsg struct {
	MessageType string                 `json:"messageType"`
	DeviceID    string                 `json:"deviceId"`
	FunctionID  string                 `json:"functionId,omitempty"`
	Inputs      map[string]interface{} `json:"inputs,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	MessageID   string                 `json:"messageId"`
}

type Connector struct {
	cfg       Config
	conn      net.Conn
	mu        sync.Mutex
	running   bool
	stopChan  chan struct{}
	connected bool

	// State of pump-03
	runState string

	// State of Parent Gateway & Sub-devices
	gwUptime       int64
	pump1State     string
	pump2State     string
	valveOpening   float64
	targetPressure float64
}

func NewConnector(cfg Config) *Connector {
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port <= 0 {
		cfg.Port = 18884
	}
	if cfg.DeviceID == "" {
		cfg.DeviceID = "pump-03"
	}
	if cfg.Keepalive <= 0 {
		cfg.Keepalive = 15
	}
	return &Connector{
		cfg:            cfg,
		stopChan:       make(chan struct{}),
		runState:       "running",
		pump1State:     "running",
		pump2State:     "running",
		valveOpening:   75.0,
		targetPressure: 1.60,
		gwUptime:       time.Now().Unix(),
	}
}

func (c *Connector) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.connected = false

	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to CCTL at %s: %w", addr, err)
	}

	c.conn = conn
	c.connected = true

	go c.readLoop(conn)

	return nil
}

func (c *Connector) readLoop(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			c.mu.Lock()
			if c.conn == conn {
				c.connected = false
			}
			c.mu.Unlock()
			fmt.Printf("[NeuronEX CCTL] readLoop disconnected: %v\n", err)
			return
		}

		fmt.Printf("[NeuronEX CCTL] Downlink raw: %s", string(line))

		var down DownlinkMsg
		if err := json.Unmarshal(line, &down); err == nil {
			c.handleDownlink(down)
		} else {
			fmt.Printf("[NeuronEX CCTL] json.Unmarshal error: %v\n", err)
		}
	}
}

func (c *Connector) handleDownlink(down DownlinkMsg) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected || c.conn == nil {
		return
	}

	fmt.Printf("[NeuronEX CCTL] Handling messageType: %s, deviceId: %s, functionId: %s\n",
		down.MessageType, down.DeviceID, down.FunctionID)

	if down.MessageType == "INVOKE_FUNCTION" {
		output := map[string]interface{}{
			"success":   true,
			"deviceId":  down.DeviceID,
			"timestamp": time.Now().UnixMilli(),
		}

		switch down.DeviceID {
		case "gw-pump-station":
			if down.FunctionID == "syncTime" {
				output["time"] = time.Now().Format("2006-01-02 15:04:05")
				output["status"] = "synced"
			}
		case "sub-pump-01":
			if down.FunctionID == "startPump" {
				c.pump1State = "running"
				output["runState"] = "running"
			} else if down.FunctionID == "stopPump" {
				c.pump1State = "stopped"
				output["runState"] = "stopped"
			} else if down.FunctionID == "setPressureTarget" {
				if val, ok := down.Inputs["targetPressure"]; ok {
					if f, ok := val.(float64); ok {
						c.targetPressure = f
						output["targetPressure"] = f
					}
				}
			}
		case "sub-pump-02":
			if down.FunctionID == "startPump" {
				c.pump2State = "running"
				output["runState"] = "running"
			} else if down.FunctionID == "stopPump" {
				c.pump2State = "stopped"
				output["runState"] = "stopped"
			}
		case "sub-valve-01":
			if down.FunctionID == "setOpening" {
				if val, ok := down.Inputs["opening"]; ok {
					if f, ok := val.(float64); ok {
						c.valveOpening = f
						output["opening"] = f
					}
				}
			}
		case "pump-03", "":
			if down.FunctionID == "startPump" {
				c.runState = "running"
				output["runState"] = "running"
			} else if down.FunctionID == "stopPump" {
				c.runState = "stopped"
				output["runState"] = "stopped"
			}
		}

		reply := map[string]interface{}{
			"messageType": "INVOKE_FUNCTION_REPLY",
			"deviceId":    down.DeviceID,
			"messageId":   down.MessageID,
			"output":      true,
			"data":        output,
			"success":     true,
		}
		data, _ := json.Marshal(reply)
		data = append(data, '\n')
		_ = c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, err := c.conn.Write(data)
		fmt.Printf("[NeuronEX CCTL] Function reply written, err=%v: %s", err, string(data))
	} else if down.MessageType == "WRITE_PROPERTY" {
		reply := map[string]interface{}{
			"messageType": "WRITE_PROPERTY_REPLY",
			"deviceId":    down.DeviceID,
			"messageId":   down.MessageID,
			"properties":  down.Properties,
			"success":     true,
		}
		data, _ := json.Marshal(reply)
		data = append(data, '\n')
		_ = c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		_, _ = c.conn.Write(data)
		fmt.Printf("[NeuronEX CCTL] Write property reply written: %s", string(data))
	}
}

func (c *Connector) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.connected
}

func (c *Connector) SendTelemetry(deviceID string, properties map[string]interface{}) error {
	if deviceID == "" {
		deviceID = c.cfg.DeviceID
	}
	if properties == nil {
		properties = make(map[string]interface{})
	}

	now := time.Now().UnixNano() / int64(time.Millisecond)
	msg := Message{
		DeviceID:   deviceID,
		Properties: properties,
		Timestamp:  now,
		MessageID:  fmt.Sprintf("nx_%d_%s", now, deviceID),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	payload := append(data, '\n')

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected || c.conn == nil {
		addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err != nil {
			c.connected = false
			return fmt.Errorf("failed to reconnect to CCTL at %s: %w", addr, err)
		}
		c.conn = conn
		c.connected = true
		go c.readLoop(conn)
	}

	_ = c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if _, err := c.conn.Write(payload); err != nil {
		_ = c.conn.Close()
		c.conn = nil
		c.connected = false
		return fmt.Errorf("failed to send telemetry to CCTL: %w", err)
	}

	return nil
}

func (c *Connector) SendHeartbeat() error {
	// Jitter for realism
	rnd := rand.Float64()*0.05 - 0.025

	// 1. pump-03 (Windows 8088 dedicated)
	_ = c.SendTelemetry("pump-03", map[string]interface{}{
		"pressure":  1.45 + rnd,
		"flow":      28.5 + rnd*10,
		"vibration": 0.72 + rnd,
		"runState":  c.runState,
	})

	// 2. Parent Gateway: gw-pump-station
	_ = c.SendTelemetry("gw-pump-station", map[string]interface{}{
		"gatewayState": "online",
		"deviceCount":  3,
		"cpuUsage":     12.5 + rand.Float64()*3.0,
		"memUsage":     34.2 + rand.Float64()*1.0,
		"uptime":       time.Now().Unix() - c.gwUptime,
	})

	// 3. Sub-device 1: sub-pump-01
	_ = c.SendTelemetry("sub-pump-01", map[string]interface{}{
		"pressure":  1.65 + rnd,
		"flow":      45.2 + rnd*15,
		"vibration": 0.42 + rnd*0.5,
		"runState":  c.pump1State,
	})

	// 4. Sub-device 2: sub-pump-02
	_ = c.SendTelemetry("sub-pump-02", map[string]interface{}{
		"pressure":  1.58 + rnd,
		"flow":      43.8 + rnd*15,
		"vibration": 0.38 + rnd*0.5,
		"runState":  c.pump2State,
	})

	// 5. Sub-device 3: sub-valve-01
	_ = c.SendTelemetry("sub-valve-01", map[string]interface{}{
		"opening": c.valveOpening,
		"status":  1,
	})

	return nil
}

func (c *Connector) Start() error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = true
	c.mu.Unlock()

	_ = c.Connect()
	_ = c.SendHeartbeat()

	go c.keepaliveLoop()
	return nil
}

func (c *Connector) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return
	}
	c.running = false
	close(c.stopChan)
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.connected = false
	c.mu.Unlock()
}

func (c *Connector) keepaliveLoop() {
	ticker := time.NewTicker(time.Duration(c.cfg.Keepalive) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			_ = c.SendHeartbeat()
		}
	}
}
