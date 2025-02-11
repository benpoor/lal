// Copyright 2020, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package rtsp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/q191201771/lal/pkg/rtprtcp"
	"github.com/q191201771/naza/pkg/nazalog"

	"github.com/q191201771/naza/pkg/nazanet"
)

// TODO chef
// - 所有client session需要defer dispose？
// - lalserver接入pullrtsp，通过HTTP API的形式
// - 超时
// - 日志
// - stat
// - pub和sub存在一些重复代码
// - sub缺少主动发送sr
// - queue的策略，当一条流没数据之后
// - 用context重写其它pull session

var ErrRTSP = errors.New("lal.rtsp: fxxk")

const (
	MethodOptions  = "OPTIONS"
	MethodAnnounce = "ANNOUNCE"
	MethodDescribe = "DESCRIBE"
	MethodSetup    = "SETUP"
	MethodRecord   = "RECORD"
	MethodPlay     = "PLAY"
	MethodTeardown = "TEARDOWN"
)

const (
	HeaderFieldCSeq      = "CSeq"
	HeaderFieldTransport = "Transport"
	HeaderFieldSession   = "Session"
)

const (
	TransportFieldClientPort  = "client_port"
	TransportFieldServerPort  = "server_port"
	TransportFieldInterleaved = "interleaved"
)

const (
	Interleaved = uint8(0x24)
)

var (
	// TODO chef: 参考协议标准，不要使用固定值
	sessionID = "191201771"

	minServerPort = uint16(8000)
	maxServerPort = uint16(16000)

	unpackerItemMaxSize = 1024

	serverCommandSessionReadBufSize = 256
)

var availUDPConnPool *nazanet.AvailUDPConnPool

// 传入远端IP，RTPPort，RTCPPort，创建两个对应的RTP和RTCP的UDP连接对象，以及对应的本端端口
func initConnWithClientPort(rHost string, rRTPPort, rRTCPPort uint16) (rtpConn, rtcpConn *nazanet.UDPConnection, lRTPPort, lRTCPPort uint16, err error) {
	// NOTICE
	// 处理Pub时，
	// 一路流的rtp端口和rtcp端口必须不同。
	// 我尝试给ffmpeg返回rtp和rtcp同一个端口，结果ffmpeg依然使用rtp+1作为rtcp的端口。
	// 又尝试给ffmpeg返回rtp:a和rtcp:a+2的端口，结果ffmpeg依然使用a和a+1端口。
	// 也即是说，ffmpeg默认认为rtcp的端口是rtp的端口+1。而不管SETUP RESPONSE的rtcp端口是多少。
	// 我目前在Acquire2这个函数里做了保证，绑定两个可用且连续的端口。

	var rtpc, rtcpc *net.UDPConn
	rtpc, lRTPPort, rtcpc, lRTCPPort, err = availUDPConnPool.Acquire2()
	if err != nil {
		return
	}
	nazalog.Debugf("acquire udp conn. rtp port=%d, rtcp port=%d", lRTPPort, lRTCPPort)

	rtpConn, err = nazanet.NewUDPConnection(func(option *nazanet.UDPConnectionOption) {
		option.Conn = rtpc
		option.RAddr = net.JoinHostPort(rHost, fmt.Sprintf("%d", rRTPPort))
		option.MaxReadPacketSize = rtprtcp.MaxRTPRTCPPacketSize
	})
	if err != nil {
		return
	}
	rtcpConn, err = nazanet.NewUDPConnection(func(option *nazanet.UDPConnectionOption) {
		option.Conn = rtcpc
		option.RAddr = net.JoinHostPort(rHost, fmt.Sprintf("%d", rRTCPPort))
		option.MaxReadPacketSize = rtprtcp.MaxRTPRTCPPacketSize
	})
	return
}

// 从setup消息的header中解析rtp rtcp channel
func parseRTPRTCPChannel(setupTransport string) (rtp, rtcp uint16, err error) {
	return parseTransport(setupTransport, TransportFieldInterleaved)
}

// 从setup消息的header中解析rtp rtcp 端口
func parseClientPort(setupTransport string) (rtp, rtcp uint16, err error) {
	return parseTransport(setupTransport, TransportFieldClientPort)
}

func parseServerPort(setupTransport string) (rtp, rtcp uint16, err error) {
	return parseTransport(setupTransport, TransportFieldServerPort)
}

func parseTransport(setupTransport string, key string) (first, second uint16, err error) {
	var clientPort string
	items := strings.Split(setupTransport, ";")
	for _, item := range items {
		if strings.HasPrefix(item, key) {
			kv := strings.Split(item, "=")
			if len(kv) != 2 {
				continue
			}
			clientPort = kv[1]
		}
	}
	items = strings.Split(clientPort, "-")
	if len(items) != 2 {
		return 0, 0, ErrRTSP
	}
	iFirst, err := strconv.Atoi(items[0])
	if err != nil {
		return 0, 0, err
	}
	iSecond, err := strconv.Atoi(items[1])
	if err != nil {
		return 0, 0, err
	}
	return uint16(iFirst), uint16(iSecond), err
}

func init() {
	availUDPConnPool = nazanet.NewAvailUDPConnPool(minServerPort, maxServerPort)
}

// ---------------------------------------------------------------------------------------------------------------------
// PUB
// ffmpeg -re -stream_loop -1 -i /Volumes/Data/tmp/wontcry.flv -acodec copy -vcodec copy -f rtsp rtsp://localhost:5544/live/test110

// read http request. method=OPTIONS, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:1 User-Agent:Lavf57.83.100], body= - server.go:95
// read http request. method=ANNOUNCE, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:2 Content-Length:490 Content-Type:application/sdp User-Agent:Lavf57.83.100], body=v=0
// o=- 0 0 IN IP4 127.0.0.1
// s=No Name
// c=IN IP4 127.0.0.1
// t=0 0
// a=tool:libavformat 57.83.100
// m=video 0 RTP/AVP 96
// a=rtpmap:96 H264/90000
// a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAFqyyAUBf8uAiAAADAAIAAAMAPB4sXJA=,aOvDyyLA; profile-level-id=640016
// a=control:streamid=0
// m=audio 0 RTP/AVP 97
// b=AS:128
// a=rtpmap:97 MPEG4-GENERIC/44100/2
// a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=121056E500
// a=control:streamid=1
// - server.go:95
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=0, headers=map[CSeq:3 Transport:RTP/AVP/UDP;unicast;client_port=32182-32183;mode=record User-Agent:Lavf57.83.100], body= - server.go:95
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=1, headers=map[CSeq:4 Session:191201771 Transport:RTP/AVP/UDP;unicast;client_port=32184-32185;mode=record User-Agent:Lavf57.83.100], body= - server.go:95
// read http request. method=RECORD, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:5 Range:npt=0.000- Session:191201771 User-Agent:Lavf57.83.100], body= - server.go:95
// read http request. method=TEARDOWN, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:6 Session:191201771 User-Agent:Lavf57.83.100], body= - server.go:95

// read udp packet failed. err=read udp [::]:8002: use of closed network connection - server_pub_session.go:199
// read udp packet failed. err=read udp [::]:8003: use of closed network connection - server_pub_session.go:199
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// PUB(rtp over tcp)
// ffmpeg -re -stream_loop -1 -i /Volumes/Data/tmp/wontcry.flv -acodec copy -vcodec copy -rtsp_transport tcp -f rtsp rtsp://localhost:5544/live/test110
//
// read http request. method=OPTIONS, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:1 User-Agent:Lavf57.83.100], body= - server.go:137
// read http request. method=ANNOUNCE, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:2 Content-Length:478 Content-Type:application/sdp User-Agent:Lavf57.83.100], body=v=0
// o=- 0 0 IN IP6 ::1
// s=No Name
// c=IN IP6 ::1
// t=0 0
// a=tool:libavformat 57.83.100
// m=video 0 RTP/AVP 96
// a=rtpmap:96 H264/90000
// a=fmtp:96 packetization-mode=1; sprop-parameter-sets=Z2QAFqyyAUBf8uAiAAADAAIAAAMAPB4sXJA=,aOvDyyLA; profile-level-id=640016
// a=control:streamid=0
// m=audio 0 RTP/AVP 97
// b=AS:128
// a=rtpmap:97 MPEG4-GENERIC/44100/2
// a=fmtp:97 profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3; config=121056E500
// a=control:streamid=1
// - server.go:137
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=0, headers=map[CSeq:3 Transport:RTP/AVP/TCP;unicast;interleaved=0-1;mode=record User-Agent:Lavf57.83.100], body= - server.go:137
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=1, headers=map[CSeq:4 Session:191201771 Transport:RTP/AVP/TCP;unicast;interleaved=2-3;mode=record User-Agent:Lavf57.83.100], body= - server.go:137
// read http request. method=RECORD, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:5 Range:npt=0.000- Session:191201771 User-Agent:Lavf57.83.100], body= - server.go:137
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// SUB
//
// read http request. method=OPTIONS, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:1 User-Agent:Lavf57.83.100], body= - server.go:108
// read http request. method=DESCRIBE, uri=rtsp://localhost:5544/live/test110, headers=map[Accept:application/sdp CSeq:2 User-Agent:Lavf57.83.100], body= - server.go:108
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=0, headers=map[CSeq:3 Transport:RTP/AVP/UDP;unicast;client_port=15690-15691 User-Agent:Lavf57.83.100], body= - server.go:108
// read http request. method=SETUP, uri=rtsp://localhost:5544/live/test110/streamid=1, headers=map[CSeq:4 Session:191201771 Transport:RTP/AVP/UDP;unicast;client_port=15692-15693 User-Agent:Lavf57.83.100], body= - server.go:108
// read http request. method=PLAY, uri=rtsp://localhost:5544/live/test110, headers=map[CSeq:5 Range:npt=0.000- Session:191201771 User-Agent:Lavf57.83.100], body= - server.go:108
// ---------------------------------------------------------------------------------------------------------------------

// ---------------------------------------------------------------------------------------------------------------------
// SUB(rtp over tcp)
//
// ---------------------------------------------------------------------------------------------------------------------

// 8000 video rtp
// 8001 video rtcp
// 8002 audio rtp
// 8003 audio rtcp
