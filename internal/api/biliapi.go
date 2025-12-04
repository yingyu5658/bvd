package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VideoBaseInfo struct {
	Data []struct {
		Cid        int    `json:"cid"`         // 每一个视频的 CID
		Part       string `json:"part"`        // 分 P 标题
		Page       int    `json:"page"`        // 分 P 编号
		FirstFrame string `json:"first_frame"` // 封面图
	} `json:"data"`
}

type BiliAPI struct {
	Client *http.Client
}

func NewBiliAPI() *BiliAPI {
	return &BiliAPI{
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

// func getInfo(bvid string) (VideoBaseInfo, error) {
// 	var APIUrl string = "https://api.bilibili.com/x/player/pagelist?bvid="
// 	APIUrl += bvid
// 	client := &http.Client{
// 		Timeout: 30 * time.Second,
// 	}

// 	// 发送 GET 请求
// 	resp, err := client.Get(APIUrl)
// 	if err != nil {
// 		return VideoBaseInfo{}, fmt.Errorf("发送请求失败: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return VideoBaseInfo{}, fmt.Errorf("读取响应失败: %w", err)
// 	}

// 	var result VideoBaseInfo
// 	if err := json.Unmarshal(body, &result); err != nil {
// 		return VideoBaseInfo{}, fmt.Errorf("JSON 解析失败： %w", err)
// 	}

// 	return VideoBaseInfo{}, err
// }

func (a *BiliAPI) GetVideoInfo(bvid string) (*VideoBaseInfo, error) {
	var APIUrl string = "https://api.bilibili.com/x/player/pagelist?bvid="
	APIUrl += bvid
	client := a.Client
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	// 发送 GET 请求
	resp, err := client.Get(APIUrl)
	if err != nil {
		return &VideoBaseInfo{}, fmt.Errorf("❌ 发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &VideoBaseInfo{}, fmt.Errorf("❌ 读取响应失败: %w", err)
	}

	if len(body) > 0 && body[0] == '<' {
		snippet := string(body)
		if len(snippet) > 512 {
			snippet = snippet[:512]
		}
		return nil, fmt.Errorf("❌ 服务器返回 HTML 而非 JSON: %s", snippet)
	}

	var result VideoBaseInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return &VideoBaseInfo{}, fmt.Errorf("❌ JSON 解析失败： %w", err)
	}

	return &result, nil

}
