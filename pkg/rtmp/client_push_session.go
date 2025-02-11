// Copyright 2019, Chef.  All rights reserved.
// https://github.com/q191201771/lal
//
// Use of this source code is governed by a MIT-style license
// that can be found in the License file.
//
// Author: Chef (191201771@qq.com)

package rtmp

import "github.com/q191201771/lal/pkg/base"

type PushSession struct {
	IsFresh bool

	core *ClientSession
}

type PushSessionOption struct {
	ConnectTimeoutMS int
	PushTimeoutMS    int
	WriteAVTimeoutMS int
}

var defaultPushSessionOption = PushSessionOption{
	ConnectTimeoutMS: 0,
	PushTimeoutMS:    0,
	WriteAVTimeoutMS: 0,
}

type ModPushSessionOption func(option *PushSessionOption)

func NewPushSession(modOptions ...ModPushSessionOption) *PushSession {
	opt := defaultPushSessionOption
	for _, fn := range modOptions {
		fn(&opt)
	}
	return &PushSession{
		IsFresh: true,
		core: NewClientSession(CSTPushSession, func(option *ClientSessionOption) {
			option.ConnectTimeoutMS = opt.ConnectTimeoutMS
			option.DoTimeoutMS = opt.PushTimeoutMS
			option.WriteAVTimeoutMS = opt.WriteAVTimeoutMS
		}),
	}
}

// 建立rtmp publish连接
// 阻塞直到收到服务端返回的rtmp publish对应结果的信令，或发生错误
func (s *PushSession) Push(rawURL string) error {
	return s.core.DoWithTimeout(rawURL)
}

func (s *PushSession) AsyncWrite(msg []byte) error {
	return s.core.AsyncWrite(msg)
}

func (s *PushSession) Flush() error {
	return s.core.Flush()
}

func (s *PushSession) Dispose() {
	s.core.Dispose()
}

func (s *PushSession) GetStat() base.StatSession {
	return s.core.GetStat()
}

func (s *PushSession) UpdateStat(interval uint32) {
	s.core.UpdateStat(interval)
}

func (s *PushSession) IsAlive() (readAlive, writeAlive bool) {
	return s.core.IsAlive()
}

func (s *PushSession) AppName() string {
	return s.core.AppName()
}

func (s *PushSession) StreamName() string {
	return s.core.StreamName()
}

func (s *PushSession) RawQuery() string {
	return s.core.RawQuery()
}

func (s *PushSession) Done() <-chan error {
	return s.core.Done()
}

func (s *PushSession) UniqueKey() string {
	return s.core.UniqueKey
}
