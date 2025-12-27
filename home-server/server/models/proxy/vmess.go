package proxy

import (
	"fmt"
	"strings"
)

const VMESS_SURGE_BASE_FORMATTER = "%s = vmess, %s, %d, username=%s, tls=%t"

type Vmess struct {
	ClashProxy
	UUID              string     `json:"uuid"`
	AlterID           int        `json:"alterId"`
	Cipher            string     `json:"cipher"`
	UDP               bool       `json:"udp,omitempty"`
	Network           string     `json:"network,omitempty"`
	TLS               bool       `json:"tls,omitempty"`
	SkipCertVerify    bool       `json:"skip-cert-verify,omitempty"`
	ClientFingerprint string     `json:"client-fingerprint,omitempty"`
	ServerName        string     `json:"servername,omitempty"`
	WSOpts            *WSOptions `json:"ws-opts,omitempty"`
}

type WSOptions struct {
	Path    string            `json:"path,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (s Vmess) GetType() string {
	return s.Type
}

func (s Vmess) ExportSurgeProxy() (string, error) {
	configBuilder := strings.Builder{}
	fmt.Fprintf(&configBuilder, VMESS_SURGE_BASE_FORMATTER, s.Name, s.Server, s.Port, s.UUID, s.TLS)
	if s.TLS && s.SkipCertVerify {
		configBuilder.WriteString(", skip-cert-verify=true")
	}
	if len(s.ServerName) > 0 {
		fmt.Fprintf(&configBuilder, ", sni=%s", s.ServerName)
	}
	if s.UDP {
		configBuilder.WriteString(", udp-relay=true")
	}
	if s.AlterID == 0 {
		configBuilder.WriteString(", vmess-aead=true")
	}
	if !strings.EqualFold("None", s.Cipher) {
		fmt.Fprintf(&configBuilder, ", encrypt-method=%s", s.Cipher)
	}
	if strings.EqualFold("ws", s.Network) {
		configBuilder.WriteString(", ws=true")
		if s.WSOpts != nil {
			headers := []string{}
			for header, value := range s.WSOpts.Headers {
				headers = append(headers, fmt.Sprintf("%s:%s", header, value))
			}
			if len(s.WSOpts.Path) > 0 {
				fmt.Fprintf(&configBuilder, ", ws-path=%s", s.WSOpts.Path)
			}
			if len(headers) > 0 {
				fmt.Fprintf(&configBuilder, ", ws-headers=%s", strings.Join(headers, "|"))
			}
		}
	}
	return configBuilder.String(), nil
}
