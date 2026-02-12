package main

import (
	"fmt"
	"log"
	"os"

	"bvd/internal/api"
	downloader "bvd/internal/download"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bvd",
		Usage: "快速、高效、易用的下载B站视频 CLI 工具",
		Commands: []*cli.Command{
			{
				Name:  "download",
				Usage: "下载指定 BV 号的视频",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Aliases: []string{"f"},
						Usage:   "指定单个视频下载时的输出文件名（不含扩展名）",
					},
					&cli.StringFlag{
						Name:    "director",
						Aliases: []string{"d"},
						Usage:   "指定下载目标文件夹（不存在时会自动创建）",
					},
				},
				Action: func(c *cli.Context) error {
					bvid := c.Args().First()
					if bvid == "" {
						return fmt.Errorf("请提供视频 BV 号")
					}

					outputFile := c.String("file")
					outputDir := c.String("director")
					// fmt.Printf("Debug: outputFile=%s, outputDir=%s\n", outputFile, outputDir)

					biliAPI := api.NewBiliAPI()
					dl := downloader.NewDownloader()
					opts := &downloader.DownloadOptions{
						OutputFile: outputFile,
						OutputDir:  outputDir,
					}

					return dl.Start(bvid, biliAPI, opts)
				},
				ArgsUsage: "<BVID> 欲下载视频的 BV 号",
			},
		},
		Version: "1.0.1",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
