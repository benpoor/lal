<p align="center">
<a title="logo" target="_blank" href="https://github.com/q191201771/lal">
<img alt="Wide" src="https://pengrl.com/images/other/lallogo.png">
</a>
<br>
<a title="TravisCI" target="_blank" href="https://www.travis-ci.org/q191201771/lal"><img src="https://www.travis-ci.org/q191201771/lal.svg?branch=master"></a>
<a title="codecov" target="_blank" href="https://codecov.io/gh/q191201771/lal"><img src="https://codecov.io/gh/q191201771/lal/branch/master/graph/badge.svg?style=flat-square"></a>
<a title="goreportcard" target="_blank" href="https://goreportcard.com/report/github.com/q191201771/lal"><img src="https://goreportcard.com/badge/github.com/q191201771/lal?style=flat-square"></a>
<br>
<a title="codeline" target="_blank" href="https://github.com/q191201771/lal"><img src="https://sloc.xyz/github/q191201771/lal/?category=code"></a>
<a title="license" target="_blank" href="https://github.com/q191201771/lal/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square"></a>
<a title="lastcommit" target="_blank" href="https://github.com/q191201771/lal/commits/master"><img src="https://img.shields.io/github/commit-activity/m/q191201771/lal.svg?style=flat-square"></a>
<br>
<a title="pr" target="_blank" href="https://github.com/q191201771/lal/pulls"><img src="https://img.shields.io/github/issues-pr-closed/q191201771/lal.svg?style=flat-square&color=FF9966"></a>
<a title="hits" target="_blank" href="https://github.com/q191201771/lal"><img src="https://hits.b3log.org/q191201771/lal.svg?style=flat-square"></a>
<a title="toplanguage" target="_blank" href="https://github.com/q191201771/lal"><img src="https://img.shields.io/github/languages/top/q191201771/lal.svg?style=flat-square"></a>
<br>
</p>

---

**`app/lalserver`服务器支持的协议：**

| - | sub rtmp | sub http(s)-flv | sub http-ts | sub hls | sub rtsp | relay push rtmp |
| - | - | - | - | - | - | - |
| pub rtmp        | ✔ | ✔ | ✔ | ✔ | X | ✔ |
| pub rtsp        | ✔ | ✔ | ✔ | ✔ | ✔ | ✔ |
| relay pull rtmp | ✔ | ✔ | ✔ | ✔ | X | . |

| 编码类型 | rtmp | rtsp | hls | http(s)-flv | http-ts |
| - | - | - | - | - | - |
| aac       | ✔ | ✔ | ✔ | ✔ | ✔ |
| avc/h264  | ✔ | ✔ | ✔ | ✔ | ✔ |
| hevc/h265 | ✔ | ✔ | X | ✔ | X |

表格含义见： [《流媒体传输连接类型之session client, server, pub, sub, push, pull》](https://pengrl.com/p/20080)

**`app/lalserver`功能特性：**

- [x] **全平台**。(依托Go语言)：支持`(linux/macOS/windows)`多平台开发、调试、运行。支持交叉编译。生成的可执行文件(无任何库依赖)可独立运行。(开放源码的同时)提供各平台可执行文件，可(免编译)直接运行
- [x] **高性能**。多核多线程扩展
- [x] **多种直播流封装协议**。支持RTMP/RTSP/HTTP-FLV/HTTP-TS/HLS，不同封装协议支持相互转换
- [x] **多种编码格式**。视频支持H264/AVC，H265/HEVC，音频支持AAC
- [x] **录制**。支持HLS录制(HLS直播与录制可同时开启)
- [x] **HTTPS**。支持HTTPS-FLV拉流
- [x] **RTSP支持interleaved mode模式**。同时支持RTP over UDP/TCP
- [x] **HTTP API接口**。用于获取服务信息，向服务发送命令。见[《lal流媒体服务器的HTTP API接口》](https://pengrl.com/p/20100)
- [x] **HTTP Notify事件回调**。见[《lal HTTP Notify(or Callback or Webhook)事件回调》](https://pengrl.com/p/20101)
- [x] **分布式集群**。
- [x] **静态pull回源**。通过配置文件配置回源地址
- [x] **静态push转推**。支持转推多个地址。通过配置文件配置转推地址
- [x] **CORS跨域**。支持HTTP-FLV，HTTP-TS，HLS跨域拉流
- [x] **HTTP文件服务器**。比如HLS切片文件可直接播放，不需要额外的HTTP文件服务器
- [x] **秒开播放**。GOP缓冲

除了lalserver，还提供一些基于lal开发的demo： [《lal/app/demo》](https://github.com/q191201771/lal/blob/master/app/demo/README.md)

<img alt="Wide" src="https://pengrl.com/images/other/lalmodule.jpg?date=0829">

发行版本日志：[《CHANGELOG.md》](https://github.com/q191201771/lal/blob/master/CHANGELOG.md)

### 编译，运行，体验功能

#### 编译

方式1，从源码自行编译

```shell
$export GO111MODULE=on && export GOPROXY=https://goproxy.cn,https://goproxy.io,direct
$make
```

方式2，直接下载编译好的二进制可执行文件

[点我打开《github lal最新release版本页面》](https://github.com/q191201771/lal/releases/latest)，下载对应平台编译好的二进制可执行文件的zip压缩包。

#### 运行

```shell
$./bin/lalserver -c conf/lalserver.conf.json
```

#### 体验功能

快速体验lalserver服务器见： [《常见推拉流客户端软件的使用方式》](https://pengrl.com/p/20051/)

lalserver详细配置见： [《配置注释文档》](https://github.com/q191201771/lal/blob/master/conf/lalserver.conf.json.brief)

### 源码框架

<br>

简单来说，源码在`pkg/`，`app/lalserver/`，`app/demo/`三个目录下。

- `pkg/`：存放各package包，供本repo的程序以及其他业务方使用
- `app/lalserver`：基于lal编写的一个通用流媒体服务器程序入口
- `app/demo/`：存放各种基于`lal/pkg`开发的小程序（小工具），一个子目录是一个程序，详情见各源码文件中头部的说明

目前唯一的第三方依赖（我自己写的Go基础库）： [github.com/q191201771/naza](https://github.com/q191201771/naza)

### 联系我

- 个人微信号： q191201771
- 个人QQ号： 191201771
- 微信群： 加我微信好友后，告诉我拉你进群
- QQ群： 1090510973

欢迎任何技术和非技术的交流。

目前lal正在收集新一轮需求中。

并且，lal也十分欢迎开源贡献者的加入。提PR前请先阅读：[《yoko版本PR规范》](https://pengrl.com/p/20070/)

### 性能测试，测试过的第三方客户端

见[《TEST.md》](https://github.com/q191201771/lal/blob/master/TEST.md)

### 项目star趋势图

觉得项目还不错，就点个star支持一下吧 :)

[![Stargazers over time](https://starchart.cc/q191201771/lal.svg)](https://starchart.cc/q191201771/lal)

