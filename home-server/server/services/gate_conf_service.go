package services

import (
	"context"
	"homeserver/common"
	"homeserver/infrastructure"
	"homeserver/models/gateway"
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

func (s *GatewayConfService) ListAllConfs(ctx echo.Context) (*[]gateway.ConfigFileInfo, error) {
	config := common.GetAppConfig()
	client, err := infrastructure.GetSSHClient(context.Background(), config.GatewaySSHUser, config.GatewaySSHAddress, []byte(config.SSHPrivateKey))
	if err != nil {
		common.Log.Error().Err(err).Msg("Get ssh client error.")
		return nil, err
	}
	files, err := client.SFTPClient.ReadDir(config.GatewayConfigPath)
	if err != nil {
		common.Log.Error().Err(err).Str("path", config.GatewayConfigPath).Msg("read config path dir error.")
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
	return &fileList, nil
}

func (s *GatewayConfService) GetConfContent(ctx echo.Context, nameList []string) (map[string][]byte, error) {
	if len(nameList) == 0 {
		return make(map[string][]byte), nil
	}
	config := common.GetAppConfig()
	client, err := infrastructure.GetSSHClient(context.Background(), config.GatewaySSHUser, config.GatewaySSHAddress, []byte(config.SSHPrivateKey))
	if err != nil {
		common.Log.Error().Err(err).Msg("Get ssh client error.")
		return nil, err
	}
	resp := make(map[string][]byte)
	for _, name := range nameList {
		filePath := path.Join(config.GatewayConfigPath, name)
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
		resp[name] = content
		_ = f.Close()
	}
	return resp, nil
}
