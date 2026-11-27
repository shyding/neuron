package license

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Secret signing key for NeuronEX enterprise licenses and activation codes
const SigningSecret = "NEURONEX-ENTERPRISE-MASTER-SIGNING-KEY-2024"

// DefaultAllPlugins list of all 59 industrial protocols supported by NeuronEX
var DefaultAllPlugins = []string{
	"ABB COMLI",
	"Allen-Bradley 5000 EtherNet/IP",
	"Allen-Bradley ControlLogix 5500",
	"Allen-Bradley DF1",
	"Allen-Bradley MicroLogix 1400",
	"BACnet/IP",
	"Beckhoff ADS",
	"CODESYS V3 TCP",
	"DLT645-1997",
	"DLT645-2007",
	"DNP 3.0",
	"EtherNet/IP(CIP)",
	"Fanuc Focas Ethernet",
	"GE Historian",
	"GE SRTP",
	"HJ212-2017",
	"HOSTLINK CMODE",
	"HollySys Modbus RTU",
	"HollySys Modbus TCP",
	"IEC60870-5-101",
	"IEC60870-5-102",
	"IEC60870-5-103",
	"IEC60870-5-104",
	"IEC61850",
	"Inovance Modbus TCP",
	"KNXnet/IP",
	"KUKA ETHERNET KRL TCP",
	"MTConnect",
	"Mitsubishi 1E",
	"Mitsubishi 3E",
	"Mitsubishi 4E",
	"Mitsubishi CNC",
	"Mitsubishi FX",
	"Modbus ASCII",
	"Modbus RTU",
	"Modbus TCP",
	"Modbus TCP QH",
	"NON A11",
	"OPC UA",
	"OPC UA Server",
	"OPCAE",
	"OPCDA",
	"Omron EtherNet/IP",
	"Omron FINS TCP",
	"Omron FINS UDP",
	"Panasonic Mewtocol",
	"SECS GEM HSMS",
	"Siemens FetchWrite",
	"Siemens MPI",
	"Siemens S7 ISOTCP",
	"Siemens S7 ISOTCP for 300/400",
	"SparkPlugB",
	"WebSocket",
	"XINJE Modbus RTU",
	"Brother CNC",
	"HEIDENHAIN CNC",
	"KND CNC",
	"Mazak CNC",
	"SYNTEC CNC",
}

// LicenseInfo represents the full NeuronEX license structure returned to the frontend
type LicenseInfo struct {
	LicenseType     string   `json:"licenseType"`     // "COMMERCIAL" / "OFFICIAL"
	CustomerName    string   `json:"customerName"`    // Customer name
	CustomerContact string   `json:"customerContact"` // Email or contact info
	StartDate       string   `json:"startDate"`       // e.g. "2024-01-01 00:00:00"
	ExpiredDate     string   `json:"expiredDate"`     // e.g. "2999-01-01 00:00:00"
	Expired         bool     `json:"expired"`         // false
	TotalTagCount   int      `json:"totalTagCount"`   // 1000000 (unlimited points)
	UsedTagCount    int      `json:"usedTagCount"`    // currently used points
	EnabledPlugins  []string `json:"enabledPlugins"`  // all enabled industrial protocols
	HardwareToken   string   `json:"hardwareToken"`   // machine identification code
	EkuiperOn       bool     `json:"ekuiperOn"`       // true
	OfficialEmail   string   `json:"officialEmail"`   // support@emqx.io
	Original        bool     `json:"original"`        // false (activated official enterprise license)
	Signature       string   `json:"signature,omitempty"`
}

// Manager manages license loading, storage, validation, and activation
type Manager struct {
	mu           sync.RWMutex
	current      *LicenseInfo
	filePath     string
	usedTagCount int
}

// NewDefaultLicense returns a full permanent enterprise commercial license with no self-limitations
func NewDefaultLicense() *LicenseInfo {
	return &LicenseInfo{
		LicenseType:     "COMMERCIAL",
		CustomerName:    "EMQ Technologies Inc.",
		CustomerContact: "support@emqx.io",
		StartDate:       "2024-01-01 00:00:00",
		ExpiredDate:     "2999-01-01 00:00:00",
		Expired:         false,
		TotalTagCount:   1000000,
		UsedTagCount:    0,
		EnabledPlugins:  DefaultAllPlugins,
		HardwareToken:   "NEURONEX-ENTERPRISE-CORE",
		EkuiperOn:       true,
		OfficialEmail:   "support@emqx.io",
		Original:        false, // Crucial: false means enterprise activated, eliminating trial popups
	}
}

// NewManager creates a license manager
func NewManager(filePath string) *Manager {
	m := &Manager{
		filePath: filePath,
		current:  NewDefaultLicense(),
	}
	m.Load()
	return m
}

// GetLicense returns a copy of current license info
func (m *Manager) GetLicense() LicenseInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lic := *m.current
	lic.UsedTagCount = m.usedTagCount
	return lic
}

// SetUsedTagCount updates the dynamic tag count
func (m *Manager) SetUsedTagCount(count int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.usedTagCount = count
}

// Load loads license from disk or keeps default
func (m *Manager) Load() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.filePath == "" {
		return
	}

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return
	}

	var lic LicenseInfo
	if err := json.Unmarshal(data, &lic); err == nil {
		// Ensure non-limiting flags
		lic.Original = false
		if lic.TotalTagCount == 0 {
			lic.TotalTagCount = 1000000
		}
		if len(lic.EnabledPlugins) == 0 {
			lic.EnabledPlugins = DefaultAllPlugins
		}
		m.current = &lic
	}
}

// Save writes current license to disk
func (m *Manager) Save() error {
	if m.filePath == "" {
		return nil
	}
	dir := filepath.Dir(m.filePath)
	_ = os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(m.current, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

// Reset resets license to permanent default
func (m *Manager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.current = NewDefaultLicense()
	return m.Save()
}

// ActivateLicense activates a license from raw license string or JSON
func (m *Manager) ActivateLicense(raw string) (*LicenseInfo, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("empty license data")
	}

	var lic LicenseInfo
	if err := json.Unmarshal([]byte(raw), &lic); err != nil {
		// Could be base64 encoded
		decoded, decErr := base32.StdEncoding.DecodeString(raw)
		if decErr != nil {
			return nil, fmt.Errorf("invalid license format: %v", err)
		}
		if err := json.Unmarshal(decoded, &lic); err != nil {
			return nil, fmt.Errorf("invalid decoded license content: %v", err)
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	lic.Original = false
	if lic.LicenseType == "" {
		lic.LicenseType = "COMMERCIAL"
	}
	if lic.TotalTagCount <= 0 {
		lic.TotalTagCount = 1000000
	}
	if len(lic.EnabledPlugins) == 0 {
		lic.EnabledPlugins = DefaultAllPlugins
	}
	if lic.CustomerName == "" {
		lic.CustomerName = "Enterprise Customer"
	}
	if lic.OfficialEmail == "" {
		lic.OfficialEmail = "support@emqx.io"
	}
	lic.EkuiperOn = true
	lic.Expired = false

	m.current = &lic
	_ = m.Save()
	return m.current, nil
}

// ActivateCode activates a license from an activation code string
// Format: NX-<encoded_payload>-<hmac_checksum>
func (m *Manager) ActivateCode(code string) (*LicenseInfo, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("empty activation code")
	}

	lic, err := ParseActivationCode(code)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	lic.Original = false
	m.current = lic
	_ = m.Save()
	return m.current, nil
}

// SignLicense computes HMAC signature for a license info
func SignLicense(customer string, tags int, expire string) string {
	msg := fmt.Sprintf("%s|%d|%s", customer, tags, expire)
	h := hmac.New(sha256.New, []byte(SigningSecret))
	h.Write([]byte(msg))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// GenerateActivationCode creates a formatted activation code string
// Format: NX-ACT-<CustomerHex>-<TagsHex>-<ExpireHex>-<ChecksumHex>
func GenerateActivationCode(customer string, tags int, expire string) string {
	if tags <= 0 {
		tags = 1000000
	}
	if expire == "" {
		expire = "2999-01-01"
	}
	if customer == "" {
		customer = "Enterprise Customer"
	}

	payload := fmt.Sprintf("%s|%d|%s", customer, tags, expire)
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(payload))

	sig := SignLicense(customer, tags, expire)
	code := fmt.Sprintf("NX-%s-%s", encoded, strings.ToUpper(sig[:8]))

	// Format into clean chunks if too long or return clean dashed code
	return code
}

// ParseActivationCode validates and decodes an activation code
func ParseActivationCode(code string) (*LicenseInfo, error) {
	code = strings.TrimSpace(code)
	if !strings.HasPrefix(code, "NX-") {
		return nil, errors.New("invalid activation code prefix (must start with NX-)")
	}

	parts := strings.Split(code, "-")
	if len(parts) < 3 {
		return nil, errors.New("invalid activation code format")
	}

	encodedPayload := parts[1]
	chksum := parts[2]

	data, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(encodedPayload)
	if err != nil {
		return nil, fmt.Errorf("invalid activation code payload: %v", err)
	}

	fields := strings.Split(string(data), "|")
	if len(fields) < 3 {
		return nil, errors.New("malformed payload in activation code")
	}

	customer := fields[0]
	var tags int
	fmt.Sscanf(fields[1], "%d", &tags)
	expire := fields[2]

	expectedSig := strings.ToUpper(SignLicense(customer, tags, expire)[:8])
	if strings.ToUpper(chksum) != expectedSig {
		return nil, errors.New("activation code checksum verification failed")
	}

	expireFull := expire
	if len(expire) == 10 {
		expireFull = expire + " 00:00:00"
	}

	return &LicenseInfo{
		LicenseType:     "COMMERCIAL",
		CustomerName:    customer,
		CustomerContact: "support@emqx.io",
		StartDate:       time.Now().Format("2006-01-02 15:04:05"),
		ExpiredDate:     expireFull,
		Expired:         false,
		TotalTagCount:   tags,
		UsedTagCount:    0,
		EnabledPlugins:  DefaultAllPlugins,
		HardwareToken:   "NEURONEX-LICENSED-" + expectedSig,
		EkuiperOn:       true,
		OfficialEmail:   "support@emqx.io",
		Original:        false,
	}, nil
}
