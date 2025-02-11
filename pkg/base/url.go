// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package base

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// 见单元测试

var ErrURL = errors.New("lal.url: fxxk")

const (
	DefaultRTMPPort = 1935
	DefaultHTTPPort = 80
	DefaultRTSPPort = 554
)

type URLPathContext struct {
	PathWithRawQuery    string
	Path                string
	PathWithoutLastItem string // 注意，没有前面的'/'，也没有后面的'/'
	LastItemOfPath      string // 注意，没有前面的'/'
	RawQuery            string
}

type URLContext struct {
	Scheme       string
	StdHost      string // host or host:port
	HostWithPort string
	Host         string
	Port         int

	//URLPathContext
	PathWithRawQuery    string
	Path                string
	PathWithoutLastItem string // 注意，没有前面的'/'，也没有后面的'/'
	LastItemOfPath      string // 注意，没有前面的'/'
	RawQuery            string
}

func ParseURLPath(path string) (ctx URLPathContext, err error) {
	stdURL, err := url.Parse(path)
	if err != nil {
		return ctx, err
	}
	return parseURLPath(stdURL)
}

func ParseURL(rawURL string, defaultPort int) (ctx URLContext, err error) {
	stdURL, err := url.Parse(rawURL)
	if err != nil {
		return ctx, err
	}
	if stdURL.Scheme == "" {
		return ctx, ErrURL
	}

	ctx.Scheme = stdURL.Scheme
	ctx.StdHost = stdURL.Host

	h, p, err := net.SplitHostPort(stdURL.Host)
	if err != nil {
		// url中端口不存r

		ctx.Host = stdURL.Host
		if defaultPort == -1 {
			ctx.HostWithPort = stdURL.Host
		} else {
			ctx.HostWithPort = net.JoinHostPort(stdURL.Host, fmt.Sprintf("%d", defaultPort))
			ctx.Port = defaultPort
		}
	} else {
		// 端口存在

		ctx.Port, err = strconv.Atoi(p)
		if err != nil {
			return ctx, err
		}
		ctx.Host = h
		ctx.HostWithPort = stdURL.Host

	}

	pathCtx, err := parseURLPath(stdURL)
	if err != nil {
		return ctx, err
	}
	ctx.PathWithRawQuery = pathCtx.PathWithRawQuery
	ctx.Path = pathCtx.Path
	ctx.PathWithoutLastItem = pathCtx.PathWithoutLastItem
	ctx.LastItemOfPath = pathCtx.LastItemOfPath
	ctx.RawQuery = pathCtx.RawQuery
	return ctx, nil
}

func ParseRTMPURL(rawURL string) (ctx URLContext, err error) {
	ctx, err = ParseURL(rawURL, DefaultRTMPPort)
	if err != nil {
		return
	}
	if ctx.Scheme != "rtmp" || ctx.Host == "" || ctx.Path == "" {
		return ctx, ErrURL
	}

	// 注意，使用ffmpeg推流时，会把`rtmp://127.0.0.1/test110`中的test110作为appName(streamName则为空)
	// 这种其实已不算十分合法的rtmp url了
	// 我们这里也处理一下，和ffmpeg保持一致
	if ctx.PathWithoutLastItem == "" && ctx.LastItemOfPath != "" {
		tmp := ctx.PathWithoutLastItem
		ctx.PathWithoutLastItem = ctx.LastItemOfPath
		ctx.LastItemOfPath = tmp
	}
	return
}

func ParseHTTPFLVURL(rawURL string) (ctx URLContext, err error) {
	ctx, err = ParseURL(rawURL, DefaultHTTPPort)
	if err != nil {
		return
	}
	if (ctx.Scheme != "http" && ctx.Scheme != "https") || ctx.Host == "" || ctx.Path == "" || !strings.HasSuffix(ctx.LastItemOfPath, ".flv") {
		return ctx, ErrURL
	}

	return
}

func ParseRTSPURL(rawURL string) (ctx URLContext, err error) {
	ctx, err = ParseURL(rawURL, DefaultRTSPPort)
	if err != nil {
		return
	}
	if ctx.Scheme != "rtsp" || ctx.Host == "" || ctx.Path == "" {
		return ctx, ErrURL
	}

	return
}

func parseURLPath(stdURL *url.URL) (ctx URLPathContext, err error) {
	ctx.Path = stdURL.Path

	index := strings.LastIndexByte(ctx.Path, '/')
	if index == -1 {
		ctx.PathWithoutLastItem = ""
		ctx.LastItemOfPath = ""
	} else if index == 0 {
		if ctx.Path == "/" {
			ctx.PathWithoutLastItem = ""
			ctx.LastItemOfPath = ""
		} else {
			ctx.PathWithoutLastItem = ""
			ctx.LastItemOfPath = ctx.Path[1:]
		}
	} else {
		ctx.PathWithoutLastItem = ctx.Path[1:index]
		ctx.LastItemOfPath = ctx.Path[index+1:]
	}

	ctx.RawQuery = stdURL.RawQuery

	if ctx.RawQuery == "" {
		ctx.PathWithRawQuery = ctx.Path
	} else {
		ctx.PathWithRawQuery = fmt.Sprintf("%s?%s", ctx.Path, ctx.RawQuery)
	}

	return ctx, nil
}
