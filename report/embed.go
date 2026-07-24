package report

import (
	"embed"
	"io/fs"
)

//go:embed ui/dist/*
var uiDist embed.FS

// GetUIFS 返回嵌入的前端文件系统
func GetUIFS() fs.FS {
	sub, err := fs.Sub(uiDist, "ui/dist")
	if err != nil {
		return uiDist
	}
	return sub
}
