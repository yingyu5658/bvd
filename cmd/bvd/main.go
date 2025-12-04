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
				Action: func(c *cli.Context) error {
					bvid := c.Args().First()
					if bvid == "" {
						return fmt.Errorf("请提供视频 BV 号")
					}

					biliAPI := api.NewBiliAPI()
					downloader := downloader.NewDownloader()

					return downloader.Start(bvid, biliAPI)
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
