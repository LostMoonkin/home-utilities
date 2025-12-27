package services

import (
	"bytes"
	"homeserver/common"
	"homeserver/infrastructure"
	"homeserver/models/proxy"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

type ProxySubService struct{}

func NewProxySubService() *ProxySubService {
	return &ProxySubService{}
}

func (s *ProxySubService) ClashToSurge(ctx echo.Context, typeFilters []string) (string, error) {
	config := common.GetAppConfig()
	rawData, err := infrastructure.HttpGet(config.ClashSubscriptionURLs[0], map[string]any{}, config.HttpProxy, 0)
	if err != nil {
		return "", err
	}
	// common.Log.Info().Bytes("Clash config", rawData).Msg("Get clash config success.")
	proxies, err := yaml.PathString("$.proxies[*]")
	if err != nil {
		return "", err
	}
	var clashNodes []yaml.RawMessage
	proxies.Read(bytes.NewReader(rawData), &clashNodes)
	if len(clashNodes) == 0 {
		common.Log.Warn().Msg("Empty proxies list in subscription, please check you config.")
		return "", nil
	}
	var validExporters []proxy.SurgeProxyExporter
	for _, nodeData := range clashNodes {
		jsonData, _ := nodeData.MarshalJSON()
		typeNode := gjson.GetBytes(jsonData, "type")
		if typeNode.Type != gjson.String {
			common.Log.Warn().Bytes("node", jsonData).Msg("Could not parse `type` in node, skip this node.")
			continue
		}
		if !proxy.SupportsType(typeNode.Str) {
			continue
		}
		surgeNode, err := proxy.NewSurgeProxyExporter(typeNode.Str, jsonData)
		if err != nil {
			common.Log.Warn().Bytes("nodeData", jsonData).Err(err).Msg("Create surge node error, skip this node.")
			continue
		}
		validExporters = append(validExporters, surgeNode)
	}
	if len(validExporters) == 0 {
		return "", nil
	}
	filterSet := make(map[string]struct{})
	for _, val := range typeFilters {
		filterSet[val] = struct{}{}
	}
	var surgeProxies strings.Builder
	for _, exporter := range validExporters {
		if len(filterSet) > 0 {
			if _, ok := filterSet[exporter.GetType()]; !ok {
				continue
			}
		}
		proxy, err := exporter.ExportSurgeProxy()
		if err != nil {
			common.Log.Warn().Any("proxy", proxy).Err(err).Msg("Export to surge proxy error, skip.")
			continue
		}
		if len(proxy) != 0 {
			surgeProxies.WriteString((proxy + "\n"))
		}
	}
	return surgeProxies.String(), nil
}
