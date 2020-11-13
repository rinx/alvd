package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/rinx/alvd/internal/net/grpc/status"
	"github.com/rinx/alvd/pkg/alvd/server/service/manager"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

type server struct {
	manager manager.Manager
}

func New(man manager.Manager) vald.Server {
	return &server{
		manager: man,
	}
}

func (s *server) Exists(
	ctx context.Context,
	meta *payload.Object_ID,
) (id *payload.Object_ID, err error) {
	ctx, cancel := context.WithCancel(ctx)

	var once sync.Once

	err = s.manager.Broadcast(ctx, func(ctx context.Context, client vald.Client) error {
		res, err := client.Exists(ctx, meta)
		if err != nil {
			return nil
		}

		if res != nil && res.Id != "" {
			once.Do(func() {
				id = &payload.Object_ID{
					Id: res.Id,
				}
				cancel()
			})
		}

		return nil
	})
	if err != nil || id == nil || id.Id == "" {
		return nil, status.WrapWithNotFound(fmt.Sprintf("not found: %s", err), err)
	}

	return id, nil
}

func (s *server) Search(
	ctx context.Context,
	req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	return res, nil
}

func (s *server) SearchByID(
	ctx context.Context,
	req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	return res, nil
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) error {
	return nil
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
	return nil
}

func (s *server) MultiSearch(
	ctx context.Context,
	reqs *payload.Search_MultiRequest,
) (res *payload.Search_Responses, errs error) {
	return res, errs
}

func (s *server) MultiSearchByID(
	ctx context.Context,
	reqs *payload.Search_MultiIDRequest,
) (res *payload.Search_Responses, errs error) {
	return res, errs
}

func (s *server) Insert(
	ctx context.Context,
	req *payload.Insert_Request,
) (ce *payload.Object_Location, err error) {
	return ce, nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) error {
	return nil
}

func (s *server) MultiInsert(
	ctx context.Context,
	reqs *payload.Insert_MultiRequest,
) (locs *payload.Object_Locations, err error) {
	return locs, nil
}

func (s *server) Update(
	ctx context.Context,
	req *payload.Update_Request,
) (res *payload.Object_Location, err error) {
	return res, nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) error {
	return nil
}

func (s *server) MultiUpdate(
	ctx context.Context,
	reqs *payload.Update_MultiRequest,
) (res *payload.Object_Locations, err error) {
	return res, nil
}

func (s *server) Upsert(
	ctx context.Context,
	req *payload.Upsert_Request,
) (loc *payload.Object_Location, err error) {
	return loc, nil
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) error {
	return nil
}

func (s *server) MultiUpsert(
	ctx context.Context,
	reqs *payload.Upsert_MultiRequest,
) (locs *payload.Object_Locations, err error) {
	return locs, nil
}

func (s *server) Remove(
	ctx context.Context,
	req *payload.Remove_Request,
) (locs *payload.Object_Location, err error) {
	return locs, nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) error {
	return nil
}

func (s *server) MultiRemove(
	ctx context.Context,
	reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, err error) {
	return locs, nil
}

func (s *server) GetObject(
	ctx context.Context,
	id *payload.Object_ID,
) (vec *payload.Object_Vector, err error) {
	return vec, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) error {
	return nil
}
