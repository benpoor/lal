// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package base_test

import (
	"testing"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/naza/pkg/assert"
)

type in struct {
	rawURL      string
	defaultPort int
}

// TODO chef: 测试IPv6的case

func TestParseURL(t *testing.T) {
	// 非法url
	_, err := base.ParseURL("invalidurl", -1)
	assert.IsNotNil(t, err)

	golden := map[in]base.URLContext{
		// 常见url，url中无端口，另外设置默认端口
		in{rawURL: "rtmp://127.0.0.1/live/test110", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1",
			HostWithPort:        "127.0.0.1:1935",
			Host:                "127.0.0.1",
			Port:                1935,
			PathWithRawQuery:    "/live/test110",
			Path:                "/live/test110",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "test110",
			RawQuery:            "",
		},
		// 域名url
		in{rawURL: "rtmp://localhost/live/test110", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "localhost",
			HostWithPort:        "localhost:1935",
			Host:                "localhost",
			Port:                1935,
			PathWithRawQuery:    "/live/test110",
			Path:                "/live/test110",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "test110",
			RawQuery:            "",
		},
		// 带参数url
		in{rawURL: "rtmp://127.0.0.1/live/test110?a=1", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1",
			HostWithPort:        "127.0.0.1:1935",
			Host:                "127.0.0.1",
			Port:                1935,
			PathWithRawQuery:    "/live/test110?a=1",
			Path:                "/live/test110",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "test110",
			RawQuery:            "a=1",
		},
		// path多级
		in{rawURL: "rtmp://127.0.0.1:19350/a/b/test110", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "/a/b/test110",
			Path:                "/a/b/test110",
			PathWithoutLastItem: "a/b",
			LastItemOfPath:      "test110",
			RawQuery:            "",
		},
		// url中无端口，没有设置默认端口
		in{rawURL: "rtmp://127.0.0.1/live/test110?a=1", defaultPort: -1}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1",
			HostWithPort:        "127.0.0.1",
			Host:                "127.0.0.1",
			Port:                0,
			PathWithRawQuery:    "/live/test110?a=1",
			Path:                "/live/test110",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "test110",
			RawQuery:            "a=1",
		},
		// url中有端口，设置默认端口
		in{rawURL: "rtmp://127.0.0.1:19350/live/test110?a=1", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "/live/test110?a=1",
			Path:                "/live/test110",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "test110",
			RawQuery:            "a=1",
		},
		// 无path
		in{rawURL: "rtmp://127.0.0.1:19350", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "",
			Path:                "",
			PathWithoutLastItem: "",
			LastItemOfPath:      "",
			RawQuery:            "",
		},
		// 无path2
		in{rawURL: "rtmp://127.0.0.1:19350/", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "/",
			Path:                "/",
			PathWithoutLastItem: "",
			LastItemOfPath:      "",
			RawQuery:            "",
		},
		// path不完整
		in{rawURL: "rtmp://127.0.0.1:19350/live", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "/live",
			Path:                "/live",
			PathWithoutLastItem: "",
			LastItemOfPath:      "live",
			RawQuery:            "",
		},
		// path不完整2
		in{rawURL: "rtmp://127.0.0.1:19350/live/", defaultPort: 1935}: {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1:19350",
			HostWithPort:        "127.0.0.1:19350",
			Host:                "127.0.0.1",
			Port:                19350,
			PathWithRawQuery:    "/live/",
			Path:                "/live/",
			PathWithoutLastItem: "live",
			LastItemOfPath:      "",
			RawQuery:            "",
		},
	}

	for k, v := range golden {
		ctx, err := base.ParseURL(k.rawURL, k.defaultPort)
		assert.Equal(t, nil, err)
		assert.Equal(t, v, ctx, k.rawURL)
	}
}

func TestParseRTMPURL(t *testing.T) {
	golden := map[string]base.URLContext{
		// 其他测试见ParseURL
		"rtmp://127.0.0.1/test110": {
			Scheme:              "rtmp",
			StdHost:             "127.0.0.1",
			HostWithPort:        "127.0.0.1:1935",
			Host:                "127.0.0.1",
			Port:                1935,
			PathWithRawQuery:    "/test110",
			Path:                "/test110",
			PathWithoutLastItem: "test110",
			LastItemOfPath:      "",
			RawQuery:            "",
		},
	}
	for k, v := range golden {
		ctx, err := base.ParseRTMPURL(k)
		assert.Equal(t, nil, err)
		assert.Equal(t, v, ctx, k)
	}
}
