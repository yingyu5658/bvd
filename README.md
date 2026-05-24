<img src="https://github.com/yingyu5658/bvd/raw/refs/heads/main/images/banner.png" alt="BVD" width="700">

**快速、高效、易用的B站视频下载工具。**

---

## Notice: De-GitHubbed

本仓库已停止在 GitHub 的开发，移至个人自建的代码庇护所。后续的所有提交与权威上游均以新地址为准。

- **新家地址:** https://git.verdant.ee/

### 如何关注 / 协同参与
1. 点击上方链接即可直接浏览极简的代码树结构。
2. 本项目拒绝使用任何基于网页端的 Pull Request。若您发现了 Bug 或有改进意愿，请直接在本地使用 `git format-patch` 生成纯文本补丁，并发送到 `im@verdant.ee`。

## 用法

```
NAME:
   bvd - 快速、高效、易用的下载B站视频 CLI 工具

USAGE:
   bvd [global options] command [command options]

COMMANDS:
   download  下载指定 BV 号的视频
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

示例：
```
➜  build git:(main) ✗ ./bvd download BV1ewwxesEu4
成功获取 1 个下载链接
开始下载第 1 个文件
保存到: downloads/gugugaga🐧🐧🐧.mp4
文件下载成功: downloads/gugugaga🐧🐧🐧.mp4
```
