package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const AppVersion = "1.0.0"

type VideoBaseInfo struct {
	Data []struct {
		Cid        int    `json:"cid"`         // 每一个视频的 CID
		Part       string `json:"part"`        // 分 P 标题
		Page       int    `json:"page"`        // 分 P 编号
		FirstFrame string `json:"first_frame"` // 封面图
	} `json:"data"`
}

// getDownloadUrl函数得到的JSON的匹配结构体
type getDownloadUrlJson struct {
	Data struct {
		Durl []struct {
			Url string `json: "durl"`
		} `json:"durl"`
	} `json: "data"`
}

func main() {
	app := &cli.App{
		Name:  "bvd",
		Usage: "快速、高效、易用的下载B站视频 CLI 工具",
		Commands: []*cli.Command{
			{
				Name:      "download",
				Usage:     "下载指定 BV 号的视频",
				Action:    downloadAction,
				ArgsUsage: "<BVID> 欲下载视频的 BV 号",
			},
		},

		Action: func(c *cli.Context) error {
			args := c.Args()
			if args.Len() == 0 {
				cli.ShowAppHelp(c)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadAction(c *cli.Context) error {
	// 从命令行获取 BVID 参数
	bvid := c.Args().First()

	if bvid == "" {
		return fmt.Errorf("请提供 BV 号")
	}

	return startDownload(bvid)
}

func startDownload(bvid string) error {

	video, err := getVideoBaseInfo(bvid)
	if err != nil {
		return fmt.Errorf("获取视频信息失败：%w", err)
	}

	err = downloadVideo(video, bvid)
	if err != nil {
		return err
	}

	return nil
}

func getVideoBaseInfo(bvid string) (VideoBaseInfo, error) {
	var APIUrl string = "https://api.bilibili.com/x/player/pagelist?bvid="
	APIUrl += bvid
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 发送 GET 请求
	resp, err := client.Get(APIUrl)
	if err != nil {
		return VideoBaseInfo{}, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return VideoBaseInfo{}, fmt.Errorf("读取响应失败: %w", err)
	}

	var result VideoBaseInfo
	if err := json.Unmarshal(body, &result); err != nil {
		return VideoBaseInfo{}, fmt.Errorf("JSON 解析失败： %w", err)
	}

	return result, nil
}

func getDownloadUrl(video VideoBaseInfo, bvid string) ([]string, error) {

	cid := make([]int, len(video.Data))
	for i := 0; i < len(video.Data); i++ {
		cid[i] = video.Data[i].Cid
	}

	var downloadUrls = make([]string, len(video.Data))
	var result getDownloadUrlJson
	client := &http.Client{Timeout: 30 * time.Second}

	for j := 0; j < len(cid); j++ {
		apiUrl := fmt.Sprintf("https://api.bilibili.com/x/player/playurl?&cid=%d&bvid=%s&qn=80", cid[j], bvid)

		resp, err := client.Get(apiUrl)
		if err != nil {
			return nil, fmt.Errorf("CID %d 请求失败: %w", cid[j], err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("CID %d 读取响应失败: %w", cid[j], err)
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("CID %d JSON解析失败: %w", cid[j], err)
		}

		downloadUrls[j] = result.Data.Durl[0].Url
	}

	fmt.Printf("成功获取 %d 个下载链接（每个CID一个）\n", len(downloadUrls))
	fmt.Printf("下载链接：%s\n", downloadUrls)
	fmt.Printf("下载链接1：%s\n", downloadUrls[0])
	fmt.Printf("下载链接2：%s\n", downloadUrls[1])
	return downloadUrls, nil
}

func downloadVideo(video VideoBaseInfo, bvid string) error {
	fmt.Println("downloadVideo is running")
	var urls []string
	urls, err := getDownloadUrl(video, bvid)
	if err != nil {
		return fmt.Errorf("获取视频下载链接失败：%w", err)
	}

	// 创建下载目录
	downloadDir := "./downloads/"
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return fmt.Errorf("创建下载目录失败：%w", err)
	}

	fmt.Printf("成功获取 %d 个下载链接\n", len(urls))

	for i := 0; i < len(urls); i++ {
		fmt.Println("url[1]", urls[1])
		fmt.Printf("开始下载第 %d 个文件\n", i+1)
		fmt.Printf("Url: %s\n", urls[i])

		// 创建支持重定向的 HTTP 客户端
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Printf("重定向到: %s\n", req.URL)
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
			return fmt.Errorf("HTTP请求失败：%w", err)
		}
		defer resp.Body.Close()

		// 检查响应状态
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
		}

		// 构建文件名
		filename := filepath.Join(downloadDir, video.Data[i].Part+".mp4")
		fmt.Printf("保存到: %s\n", filename)

		// 创建文件
		out, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("创建文件失败：%w", err)
		}

		// 下载文件
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			out.Close()
			return fmt.Errorf("下载失败：%w", err)
		}

		// 关闭文件
		if err := out.Close(); err != nil {
			return fmt.Errorf("关闭文件失败：%w", err)
		}

		fmt.Printf("文件下载成功: %s\n", filename)
	}

	return nil
}
