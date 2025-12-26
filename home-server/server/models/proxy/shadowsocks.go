package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
)

const SS_SURGE_BASE_FORMATTER = "%s = ss, %s, %d, encrypt-method=%s, password=%s"

type ShadowSocks struct {
	ClashProxy
	Password          string         `json:"password"`
	UDP               bool           `json:"udp,omitempty"`
	UDPOverTCP        bool           `proxy:"udp-over-tcp,omitempty"`
	UDPOverTCPVersion int            `proxy:"udp-over-tcp-version,omitempty"`
	Cipher            string         `json:"cipher"`
	Plugin            string         `json:"plugin,omitempty"`
	PluginOpts        map[string]any `json:"plugin-opts,omitempty"`
}

func (s ShadowSocks) GetType() string {
	return s.Type
}

func (s ShadowSocks) ExportSurgeProxy() (string, error) {
	surgeProxyConfig := fmt.Sprintf(SS_SURGE_BASE_FORMATTER, s.Name, s.Server, s.Port, s.Cipher, s.Password)
	var appendConfig string
	var err error
	switch s.Plugin {
	case "shadow-tls":
		appendConfig, err = parseShadowTLSOpts(s.PluginOpts)
	case "obfs":
		appendConfig, err = parseObfsSOpts(s.PluginOpts)
	}
	if err != nil {
		return "", err
	}
	if len(appendConfig) > 0 {
		surgeProxyConfig += appendConfig
	}
	return surgeProxyConfig, nil
}

func parseShadowTLSOpts(opts map[string]any) (string, error) {
	if len(opts) == 0 {
		return "", errors.New("Empty shadowTLS options")
	}
	config := ", shadow-tls-sni=%s, shadow-tls-password=%s, shadow-tls-version=%.0f"
	host, hostOk := opts["host"]
	password, passWordOk := opts["password"]
	version, versionOk := opts["version"]
	if hostOk && passWordOk && versionOk {
		return fmt.Sprintf(config, host, password, version), nil
	}
	optStr, _ := json.Marshal(opts)
	return "", fmt.Errorf("Invalid shadowdTLS options, opts=%s", string(optStr))
}

func parseObfsSOpts(opts map[string]any) (string, error) {
	if len(opts) == 0 {
		return "", errors.New("Empty obfs options")
	}
	config := ", obfs=%s, obfs-host=%s"
	host, hostOk := opts["host"]
	mode, modeOk := opts["mode"]
	if hostOk && modeOk {
		return fmt.Sprintf(config, mode, host), nil
	}
	optStr, _ := json.Marshal(opts)
	return "", fmt.Errorf("Invalid obfs options, opts=%s", string(optStr))
}
