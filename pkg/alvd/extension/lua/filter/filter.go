//   This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package filter

import (
	"bufio"
	"os"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	luar "layeh.com/gopher-luar"

	libs "github.com/vadv/gopher-lua-libs"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

type filter struct {
	filePath string

	proto *lua.FunctionProto
}

type RetryConfig struct {
	Enabled           bool
	MaxRetries        int
	NextNumMultiplier int
}

type EgressFilter interface {
	Do(origin []*payload.Object_Distance) (results []*payload.Object_Distance, retry *RetryConfig, err error)
}

func NewEgressFilter(filePath string) (EgressFilter, error) {
	proto, err := CompileLua(filePath)
	if err != nil {
		return nil, err
	}

	return &filter{
		filePath: filePath,
		proto:    proto,
	}, nil
}

func CompileLua(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}

	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

func (f *filter) Do(origin []*payload.Object_Distance) (results []*payload.Object_Distance, retry *RetryConfig, err error) {
	state := lua.NewState()
	defer state.Close()

	libs.Preload(state)

	results = origin
	retry = &RetryConfig{
		Enabled:           false,
		MaxRetries:        3,
		NextNumMultiplier: 2,
	}

	state.SetGlobal("results", luar.New(state, results))
	state.SetGlobal("retry", luar.New(state, retry))

	fn := state.NewFunctionFromProto(f.proto)
	state.Push(fn)
	err = state.PCall(0, lua.MultRet, nil)
	if err != nil {
		return origin, retry, err
	}

	return results, retry, nil
}
