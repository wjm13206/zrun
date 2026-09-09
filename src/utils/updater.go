// Package utils 提供了 zrun 的工具函数。
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"zrun/src/version"
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

// 获取版本信息，带超时，可单测。
func fetchRemoteVersionInfo(url string) (*RemoteVersionInfo, error) {
	return fetchRemoteVersionInfoWithContext(context.Background(), url)
}

func fetchRemoteVersionInfoWithContext(ctx context.Context, url string) (*RemoteVersionInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	var versionInfo RemoteVersionInfo
	if err := json.Unmarshal(body, &versionInfo); err != nil {
		return nil, err
	}
	return &versionInfo, nil
}

// 返回：是否有更新、语法是否不兼容。
// 应用版本按 YY.大版本.小版本逐段数字比较，避免字符串比较误判（如 26.10.0 < 26.9.0）。
func CheckUpdate(currentVersion, syntaxVersion string, remote *RemoteVersionInfo) (hasUpdate bool, incompatible bool) {
	if remote == nil {
		return false, false
	}
	hasUpdate = version.CompareStrings(currentVersion, remote.LatestVersion) < 0
	incompatible = syntaxVersion != remote.LatestSyntaxVersion
	return hasUpdate, incompatible
}

// 检查版本更新
func CheckSyntaxUpdates(currentVersion string, syntaxVersion string) {
	remoteInfo, err := fetchRemoteVersionInfo(UpdateURL)
	if err != nil {
		fmt.Printf("无法获取远程版本信息: %v\n", err)
		return
	}

	hasUpdate, incompatible := CheckUpdate(currentVersion, syntaxVersion, remoteInfo)
	if !hasUpdate {
		fmt.Println("当前已是最新版本")
		return
	}

	fmt.Printf("发现新版本: %s\n", remoteInfo.LatestVersion)

	if incompatible {
		fmt.Println(versionIncompatibleWarning)
	}

	fmt.Printf("下载地址: %s\n", remoteInfo.DownloadURL)
}
