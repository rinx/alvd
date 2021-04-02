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

type EgressFilter interface {
	Do([]*payload.Object_Distance) ([]*payload.Object_Distance, error)
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

func (f *filter) Do(origin []*payload.Object_Distance) (results []*payload.Object_Distance, err error) {
	state := lua.NewState()
	defer state.Close()

	libs.Preload(state)

	results = origin

	state.SetGlobal("results", luar.New(state, results))

	fn := state.NewFunctionFromProto(f.proto)
	state.Push(fn)
	err = state.PCall(0, lua.MultRet, nil)
	if err != nil {
		return origin, err
	}

	return results, nil
}
