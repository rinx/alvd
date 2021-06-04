//   This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lua

import (
	"io"
	"os"

	"github.com/yuin/gluamapper"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	libs "github.com/vadv/gopher-lua-libs"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

type LFunction = lua.LFunction

func MapConfig(filePath, varname string, st interface{}) error {
	state := lua.NewState()
	defer state.Close()

	libs.Preload(state)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = state.DoString(string(bytes))
	if err != nil {
		return err
	}

	table, ok := state.GetGlobal(varname).(*lua.LTable)
	if !ok {
		return nil
	}

	return gluamapper.Map(table, st)
}

type RetryConfig struct {
	Enabled           bool
	MaxRetries        int
	NextNumMultiplier int
}

type SearchResultInterceptor = func(
	*payload.Search_Config,
	[]*payload.Object_Distance,
) (
	[]*payload.Object_Distance,
	*RetryConfig,
	error,
)

func NewSearchResultInterceptorFn(sri *LFunction) SearchResultInterceptor {
	return func(
		cfg *payload.Search_Config,
		origin []*payload.Object_Distance,
	) (
		results []*payload.Object_Distance,
		retry *RetryConfig,
		err error,
	) {
		state := lua.NewState()
		defer state.Close()

		libs.Preload(state)

		results = origin
		retry = &RetryConfig{
			Enabled:           false,
			MaxRetries:        3,
			NextNumMultiplier: 2,
		}

		err = state.CallByParam(
			lua.P{
				Fn:      sri,
				NRet:    0,
				Protect: true,
			},
			luar.New(state, cfg),
			luar.New(state, results),
			luar.New(state, retry),
		)
		if err != nil {
			return origin, retry, err
		}

		return results, retry, nil
	}
}

type InsertDataInterceptor = func(*payload.Insert_Request) (
	*payload.Insert_Request,
	error,
)

func NewInsertDataInterceptorFn(idi *LFunction) InsertDataInterceptor {
	return func(origin *payload.Insert_Request) (
		req *payload.Insert_Request,
		err error,
	) {
		state := lua.NewState()
		defer state.Close()

		libs.Preload(state)

		req = origin

		err = state.CallByParam(
			lua.P{
				Fn:      idi,
				NRet:    0,
				Protect: true,
			},
			luar.New(state, req),
		)
		if err != nil {
			return origin, err
		}

		return req, nil
	}
}
