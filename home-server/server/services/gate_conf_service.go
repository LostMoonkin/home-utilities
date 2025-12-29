package services

import (
	"encoding/base64"
	"homeserver/common"
	"homeserver/context"
	"homeserver/infrastructure"
	"homeserver/models/gateway"
	"homeserver/models/response"
	"io"
	"os"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

type GatewayConfService struct{}

func NewGatewayConfService() *GatewayConfService {
	return &GatewayConfService{}
}

func (s *GatewayConfService) ListAllConfigs(ctx context.GContext) (*response.APIResponse[*[]gateway.ConfigFileInfo], error) {
	client, err := getSSHClient(ctx)
	if err != nil {
		return nil, err
	}
	files, err := client.SFTPClient.ReadDir(ctx.GetAppConfig().GatewayConfigPath)
	if err != nil {
		common.Log.Error().Err(err).Str("path", ctx.GetAppConfig().GatewayConfigPath).Msg("read config path dir error.")
		return nil, err
	}
	var fileList []gateway.ConfigFileInfo
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), ".conf") {
			fileList = append(fileList, gateway.ConfigFileInfo{
				Name:    file.Name(),
				ModTime: file.ModTime(),
				Size:    file.Size(),
			})
		}
	}
	return response.Success(&fileList), nil
}

func (s *GatewayConfService) GetConfContent(ctx context.GContext, nameList []string) (*response.APIResponse[map[string]string], error) {
	client, err := getSSHClient(ctx)
	if err != nil {
		common.Log.Error().Err(err).Msg("Get ssh client error.")
		return nil, err
	}
	resp := make(map[string]string)
	for _, name := range nameList {
		filePath := path.Join(ctx.GetAppConfig().GatewayConfigPath, name)
		f, err := client.SFTPClient.OpenFile(filePath, os.O_RDONLY)
		if err != nil {
			common.Log.Error().Err(err).Str("path", filePath).Msg("open file error.")
			return nil, err
		}
		content, err := io.ReadAll(f)
		if err != nil {
			common.Log.Error().Err(err).Str("path", filePath).Msg("read file content error.")
			_ = f.Close()
			return nil, err
		}
		resp[name] = base64.StdEncoding.EncodeToString(content)
		_ = f.Close()
	}
	return response.Success(resp), nil
}

func (s *GatewayConfService) UpdateConfig(ctx echo.Context, name, current, expected string) error {
	return nil
}

func getSSHClient(ctx context.GContext) (*infrastructure.SSHClientWrapper, error) {
	config := ctx.GetAppConfig()
	client, err := infrastructure.GetSSHClient(ctx.Request().Context(), config.GatewaySSHUser, config.GatewaySSHAddress, []byte(config.SSHPrivateKey))
	if err != nil {
		common.Log.Error().Err(err).Msg("Get ssh client error.")
		return nil, err
	}
	return client, nil
}
