package proxy

import (
	"fmt"
	"strings"
)

const HY2_SURGE_BASE_FORMATTER = "%s = hysteria2, %s, %d, password=%s"

type Hysteria2 struct {
	ClashProxy
	Password       string `json:"password"`
	Up             string `json:"up,omitempty"`
	Down           string `json:"down,omitempty"`
	SkipCertVerify bool   `json:"skip-cert-verify,omitempty"`
	SNI            string `json:"sni,omitempty"`
}

func (s Hysteria2) GetType() string {
	return s.Type
}

func (s Hysteria2) ExportSurgeProxy() (string, error) {
	configBuilder := strings.Builder{}
	fmt.Fprintf(&configBuilder, HY2_SURGE_BASE_FORMATTER, s.Name, s.Server, s.Port, s.Password)
	if len(s.Up) > 0 {
		fmt.Fprintf(&configBuilder, ", upload-bandwidth=%s", s.Up)
	}
	if len(s.Down) > 0 {
		fmt.Fprintf(&configBuilder, ", download-bandwidth=%s", s.Down)
	}
	if len(s.SNI) > 0 {
		fmt.Fprintf(&configBuilder, ", sni=%s", s.SNI)
	}
	if s.SkipCertVerify {
		configBuilder.WriteString(", skip-cert-verify=true")
	}
	return configBuilder.String(), nil
}
