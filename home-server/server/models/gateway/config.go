package gateway

import "time"

type ConfigFileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type UpdateConfigRequest struct {
	Name            string `json:"name"`
	CurrentContent  string `json:"currentContent"`
	ExpectedContent string `json:"expectedContent"`
}

type CreateConfigRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}
