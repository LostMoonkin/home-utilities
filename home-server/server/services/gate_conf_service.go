package services

import (
	"bytes"
	"encoding/base64"
	"homeserver/common"
	"homeserver/context"
	"homeserver/models/gateway"
	"homeserver/models/response"
	"io"
	"os"
	"path"
	"strings"
	"sync"
)

type GatewayConfService struct {
	fileLocks   sync.Map
	restartLock sync.Mutex
}

func NewGatewayConfService() *GatewayConfService {
	return &GatewayConfService{fileLocks: sync.Map{}, restartLock: sync.Mutex{}}
}

func (s *GatewayConfService) ListAll(ctx context.GContext) (*response.APIResponse, error) {
	client, err := ctx.GetGatewaySSHClient()
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
		if strings.HasSuffix(file.Name(), ".conf") {
			fileList = append(fileList, gateway.ConfigFileInfo{
				Name:    file.Name(),
				ModTime: file.ModTime(),
				Size:    file.Size(),
			})
		}
	}
	return response.Success(&fileList), nil
}

func (s *GatewayConfService) Get(ctx context.GContext, nameList []string) (*response.APIResponse, error) {
	client, err := ctx.GetGatewaySSHClient()
	if err != nil {
		return nil, err
	}
	resp := make(map[string]string)
	for _, name := range nameList {
		if !isValidConfigName(name) {
			common.Log.Warn().Str("filename", name).Msg("invalid filename.")
		}
		content, err := client.ReadFile(path.Join(ctx.GetAppConfig().GatewayConfigPath, name))
		if err != nil {
			return nil, err
		}
		resp[name] = base64.StdEncoding.EncodeToString(content)
	}
	return response.Success(resp), nil
}

func (s *GatewayConfService) Create(ctx context.GContext, name, content string) (*response.APIResponse, error) {
	if !isValidConfigName(name) {
		return response.Fail(-1, "invalid config name"), nil
	}
	// decode content
	rawContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		common.Log.Error().Err(err).Str("name", name).Str("content", content).Msg("base64 decode content error")
		return response.Fail(-1, "invalid config content"), nil
	}
	// lock current file
	rawLock, _ := s.fileLocks.LoadOrStore(name, &sync.Mutex{})
	lock := rawLock.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	client, err := ctx.GetGatewaySSHClient()
	if err != nil {
		return nil, err
	}
	filePath := path.Join(ctx.GetAppConfig().GatewayConfigPath, name)
	// check file exists
	ok, err := client.FileExists(filePath)
	if err != nil {
		return nil, err
	}
	if ok {
		return response.Fail(-1, "config file already exists."), nil
	}
	f, err := client.SFTPClient.Create(filePath)
	if err != nil {
		common.Log.Error().Err(err).Str("path", filePath).Msg("create file error.")
		return nil, err
	}
	defer f.Close()
	if _, err = io.Copy(f, bytes.NewReader(rawContent)); err != nil {
		common.Log.Error().Err(err).Str("path", filePath).Msg("write content error.")
		return nil, err
	}
	return response.Success(nil), nil
}

func (s *GatewayConfService) Update(ctx context.GContext, name, current, expected string) (*response.APIResponse, error) {
	if !isValidConfigName(name) {
		return response.Fail(-1, "invalid config name"), nil
	}
	client, err := ctx.GetGatewaySSHClient()
	if err != nil {
		return nil, err
	}
	// lock current file
	rawLock, _ := s.fileLocks.LoadOrStore(name, &sync.Mutex{})
	lock := rawLock.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	filePath := path.Join(ctx.GetAppConfig().GatewayConfigPath, name)
	rawContent, err := client.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	// check content
	content := base64.StdEncoding.EncodeToString(rawContent)
	if content != current {
		return response.Fail(-1, "config content has been modified."), nil
	}
	// update content
	expectedRawContent, err := base64.StdEncoding.DecodeString(expected)
	if err != nil {
		common.Log.Error().Err(err).Str("path", filePath).Str("expected", expected).Msg("base64 decode content error")
		return response.Fail(-1, "invalid expected content."), nil
	}
	f, err := client.SFTPClient.OpenFile(filePath, os.O_RDWR)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = io.Copy(f, bytes.NewReader(expectedRawContent)); err != nil {
		common.Log.Error().Err(err).Str("path", filePath).Msg("write content error.")
		return nil, err
	}
	return response.Success(nil), nil
}

func (s *GatewayConfService) Restart(ctx context.GContext) (*response.APIResponse, error) {
	if !s.restartLock.TryLock() {
		common.Log.Warn().Msg("restartLock is already locked")
		return response.Fail(-1, "Process is locked, please try later."), nil
	}
	defer s.restartLock.Unlock()
	client, err := ctx.GetGatewaySSHClient()
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		common.Log.Error().Err(err).Msg("new session error")
		return nil, err
	}
	defer session.Close()
	cmd := "/usr/local/bin/docker restart gateway"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		common.Log.Error().Err(err).Msg("docker restart command error")
		return nil, err
	}
	common.Log.Debug().Str("output", string(output)).Msg("docker restart command output")
	if !strings.Contains(string(output), "gateway") {
		common.Log.Warn().Bytes("output", output).Msg("docker restart command output does not contain 'gateway'")
		return response.Fail(-1, "restart gateway failed."), nil
	}
	return response.Success(nil), nil
}

func isValidConfigName(filename string) bool {
	// Filenames cannot be empty strings
	if filename == "" {
		return false
	}
	// The null byte is the only universally invalid character
	if strings.ContainsRune(filename, '\000') {
		return false
	}
	// The forward slash is the path separator and cannot be in a filename
	if strings.ContainsRune(filename, '/') {
		return false
	}
	// The file must has suffix `.conf`
	if !strings.HasSuffix(filename, ".conf") {
		return false
	}
	return true
}
