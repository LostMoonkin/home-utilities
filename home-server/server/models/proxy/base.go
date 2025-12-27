package proxy

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type SurgeProxyExporter interface {
	ExportSurgeProxy() (string, error)
	GetType() string
}

type ClashProxy struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Server string `json:"server"`
	Port   int    `json:"port"`
}

var supportedProxies map[string]reflect.Type

func init() {
	supportedProxies = map[string]reflect.Type{
		"ss":        reflect.TypeFor[ShadowSocks](),
		"vmess":     reflect.TypeFor[Vmess](),
		"hysteria2": reflect.TypeFor[Hysteria2](),
	}
}

func SupportsType(proxyType string) bool {
	_, ok := supportedProxies[proxyType]
	return ok
}

func NewSurgeProxyExporter(proxyType string, rawData []byte) (SurgeProxyExporter, error) {
	rType, ok := supportedProxies[proxyType]
	if !ok {
		return nil, fmt.Errorf("Unsupported clash proxy type: {%s}", proxyType)
	}
	v := reflect.New(rType).Interface()
	if err := json.Unmarshal(rawData, v); err != nil {
		return nil, err
	}
	return v.(SurgeProxyExporter), nil
}
