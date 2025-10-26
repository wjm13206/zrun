// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// URL
const UpdateURL = "https://raw.githubusercontent.com/wjm13206/zrun/refs/heads/main/version.json"

// 警告信息
const versionIncompatibleWarning = "警告：检测到语法版本已升级可能不兼容此版本的脚本，请谨慎升级"

// 远程版本信息结构
type RemoteVersionInfo struct {
	LatestVersion       string `json:"version"`
	LatestSyntaxVersion string `json:"syntax_version"`
	DownloadURL         string `json:"download_url"`
}

// 获取版本信息
func fetchRemoteVersionInfo(url string) (*RemoteVersionInfo, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var versionInfo RemoteVersionInfo
	err = json.Unmarshal(body, &versionInfo)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

// 检查版本更新
func CheckSyntaxUpdates(currentVersion string, syntaxVersion string) {

	remoteInfo, err := fetchRemoteVersionInfo(UpdateURL)
	if err != nil {
		fmt.Printf("无法获取远程版本信息: %v\n", err)
		return
	}

	if currentVersion == remoteInfo.LatestVersion {
		fmt.Println("当前已是最新版本")
		return
	}

	fmt.Printf("发现新版本: %s\n", remoteInfo.LatestVersion)

	if syntaxVersion != remoteInfo.LatestSyntaxVersion {
		fmt.Println(versionIncompatibleWarning)
	}

	fmt.Printf("下载地址: %s\n", remoteInfo.DownloadURL)
}
