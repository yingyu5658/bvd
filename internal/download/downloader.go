package downloader

import (
	"bvd/internal/api"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct{}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func requestSend(urls []string, video *api.VideoBaseInfo) error {
	downloadDir := "./downloads/"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("❌ 创建下载目录失败：%w", err)
	}

	for i := 0; i < len(urls); i++ {
		fmt.Printf("✅ 开始下载第 %d 个文件\n", i+1)

		// 创建支持重定向的 HTTP 客户端
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil
			},
		}

		req, err := http.NewRequest("GET", urls[i], nil)
		if err != nil {
			return err
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
		req.Header.Set("referer", "https://www.bilibili.com")

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("❌ HTTP请求失败：%w", err)
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("❌ 服务器返回错误状态码: %d", resp.StatusCode)
		}

		// 构建文件名
		filename := filepath.Join(downloadDir, video.Data[i].Part+".mp4")
		fmt.Printf("✅ 保存到: %s\n", filename)

		// 创建文件
		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("❌ 创建文件失败：%w", err)
		}

		// 下载文件
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			out.Close()
			return fmt.Errorf("❌ 下载失败：%w", err)
		}

		// 关闭文件
		if err := out.Close(); err != nil {
			return fmt.Errorf("❌ 关闭文件失败：%w", err)
		}

		fmt.Printf("✅ 文件下载成功: %s\n", filename)

	}

	return nil
}

func (d *Downloader) Start(bvid string, apiClient *api.BiliAPI) error {
	fmt.Println("✅ 获取到 BV 号:", bvid)

	videoInfo, err := apiClient.GetVideoInfo(bvid)
	if err != nil {
		return err
	}

	// 2. 提取 提取所有 CID
	cids := make([]int, len(videoInfo.Data))
	for i, page := range videoInfo.Data {
		cids[i] = page.Cid
	}

	fmt.Println("✅ 获取到 CID 列表", cids)

	// 3. 为每个CID获取下载链接
	downloadUrls := make([]string, len(cids))

	for j, cid := range cids {
		apiUrl := fmt.Sprintf("https://api.bilibili.com/x/player/playurl?&cid=%d&bvid=%s&qn=80", cid, bvid)

		resp, err := apiClient.Client.Get(apiUrl)
		if err != nil {
			return fmt.Errorf("❌ CID %d 请求 请求失败: %w", cid, err)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Errorf("❌ CID %d 读取响应失败: %w", cid, err)
		}

		var result struct {
			Data struct {
				Durl []struct {
					Url string `json:"url"`
				} `json:"durl"`
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("❌ CID %d JSON解析失败: %w", cid, err)
		}

		if len(result.Data.Durl) == 0 {
			return fmt.Errorf("❌ CID %d 没有 没有获取到下载链接", cid)
		}

		downloadUrls[j] = result.Data.Durl[0].Url
	}

	fmt.Println("✅ 获取到所有下载链接")

	requestSend(downloadUrls, videoInfo)

	return nil
}
