package handler

import (
	"context"
	"io"
	"strconv"
	"sync"

	"github.com/rinx/alvd/internal/errors"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/model"
	"github.com/rinx/alvd/pkg/vald/agent/ngt/service"
	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

type Server interface {
	core.AgentServer
	vald.Server
}

type server struct {
	name string
	ngt  service.NGT
}

func New(name string, ngt service.NGT) Server {
	return &server{
		name: name,
		ngt:  ngt,
	}
}

func (s *server) newLocations(uuids ...string) (locs *payload.Object_Locations) {
	if len(uuids) == 0 {
		return nil
	}
	locs = &payload.Object_Locations{
		Locations: make([]*payload.Object_Location, 0, len(uuids)),
	}
	for _, uuid := range uuids {
		locs.Locations = append(locs.Locations, &payload.Object_Location{
			Name: s.name,
			Uuid: uuid,
			Ips:  []string{s.name},
		})
	}
	return locs
}

func (s *server) newLocation(uuid string) *payload.Object_Location {
	locs := s.newLocations(uuid)
	if locs != nil && locs.Locations != nil && len(locs.Locations) > 0 {
		return locs.Locations[0]
	}
	return nil
}

func (s *server) Exists(ctx context.Context, uid *payload.Object_ID) (res *payload.Object_ID, err error) {
	uuid := uid.GetId()
	oid, ok := s.ngt.Exists(uuid)
	if !ok {
		return nil, err
	}
	return &payload.Object_ID{
		Id: strconv.Itoa(int(oid)),
	}, nil
}

func (s *server) Search(ctx context.Context, req *payload.Search_Request) (*payload.Search_Response, error) {
	return toSearchResponse(
		s.ngt.Search(
			req.GetVector(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func (s *server) SearchByID(ctx context.Context, req *payload.Search_IDRequest) (*payload.Search_Response, error) {
	return toSearchResponse(
		s.ngt.SearchByID(
			req.GetId(),
			req.GetConfig().GetNum(),
			req.GetConfig().GetEpsilon(),
			req.GetConfig().GetRadius()))
}

func toSearchResponse(dists []model.Distance, err error) (res *payload.Search_Response, rerr error) {
	res = new(payload.Search_Response)

	res.Results = make([]*payload.Object_Distance, 0, len(dists))
	for _, dist := range dists {
		res.Results = append(res.Results, &payload.Object_Distance{
			Id:       dist.ID,
			Distance: dist.Distance,
		})
	}

	return res, err
}

func (s *server) StreamSearch(stream vald.Search_StreamSearchServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.Search(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) StreamSearchByID(stream vald.Search_StreamSearchByIDServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.SearchByID(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) MultiSearch(ctx context.Context, reqs *payload.Search_MultiRequest) (res *payload.Search_Responses, errs error) {
	return nil, nil
}

func (s *server) MultiSearchByID(ctx context.Context, reqs *payload.Search_MultiIDRequest) (res *payload.Search_Responses, errs error) {
	return nil, nil
}

func (s *server) Insert(ctx context.Context, req *payload.Insert_Request) (res *payload.Object_Location, err error) {
	vec := req.GetVector()
	err = s.ngt.Insert(vec.GetId(), vec.GetVector())
	if err != nil {
		return nil, err
	}

	return s.newLocation(vec.GetId()), nil
}

func (s *server) StreamInsert(stream vald.Insert_StreamInsertServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.Insert(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) MultiInsert(ctx context.Context, reqs *payload.Insert_MultiRequest) (res *payload.Object_Locations, err error) {
	return nil, nil
}

func (s *server) Update(ctx context.Context, req *payload.Update_Request) (res *payload.Object_Location, err error) {
	vec := req.GetVector()
	err = s.ngt.Update(vec.GetId(), vec.GetVector())
	if err != nil {
		return nil, err
	}

	return s.newLocation(vec.GetId()), nil
}

func (s *server) StreamUpdate(stream vald.Update_StreamUpdateServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.Update(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) MultiUpdate(ctx context.Context, reqs *payload.Update_MultiRequest) (res *payload.Object_Locations, err error) {
	return nil, nil
}

func (s *server) Upsert(ctx context.Context, req *payload.Upsert_Request) (*payload.Object_Location, error) {
	_, exists := s.ngt.Exists(req.GetVector().GetId())
	if exists {
		return s.Update(ctx, &payload.Update_Request{
			Vector: req.GetVector(),
		})
	}

	return s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
	})
}

func (s *server) StreamUpsert(stream vald.Upsert_StreamUpsertServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.Upsert(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) MultiUpsert(ctx context.Context, reqs *payload.Upsert_MultiRequest) (res *payload.Object_Locations, err error) {
	return nil, nil
}

func (s *server) Remove(ctx context.Context, req *payload.Remove_Request) (res *payload.Object_Location, err error) {
	id := req.GetId()
	uuid := id.GetId()
	err = s.ngt.Delete(uuid)
	if err != nil {
		return nil, err
	}

	return s.newLocation(uuid), nil
}

func (s *server) StreamRemove(stream vald.Remove_StreamRemoveServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.Remove(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) MultiRemove(ctx context.Context, reqs *payload.Remove_MultiRequest) (res *payload.Object_Locations, err error) {
	return nil, nil
}

func (s *server) GetObject(ctx context.Context, id *payload.Object_ID) (res *payload.Object_Vector, err error) {
	uuid := id.GetId()
	vec, err := s.ngt.GetObject(uuid)
	if err != nil {
		return nil, err
	}

	return &payload.Object_Vector{
		Id:     uuid,
		Vector: vec,
	}, nil
}

func (s *server) StreamGetObject(stream vald.Object_StreamGetObjectServer) error {
	ctx := stream.Context()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	errs := make([]error, 0)
	emu := sync.Mutex{}

	close := func() (err error) {
		if len(errs) != 0 {
			for _, e := range errs {
				err = errors.Wrap(err, e.Error())
			}
		}

		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()

			return close()
		default:
			req, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					wg.Wait()

					return close()
				}

				return err
			}

			if req != nil {
				wg.Add(1)

				go func() {
					defer wg.Done()

					res, err := s.GetObject(ctx, req)
					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}

					mu.Lock()
					err = stream.Send(res)
					mu.Unlock()

					if err != nil {
						emu.Lock()
						defer emu.Unlock()

						errs = append(errs, err)

						return
					}
				}()
			}
		}
	}
}

func (s *server) CreateIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.CreateIndex(ctx, c.GetPoolSize())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) SaveIndex(ctx context.Context, _ *payload.Empty) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.SaveIndex(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) CreateAndSaveIndex(ctx context.Context, c *payload.Control_CreateIndexRequest) (res *payload.Empty, err error) {
	res = new(payload.Empty)
	err = s.ngt.CreateAndSaveIndex(ctx, c.GetPoolSize())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) IndexInfo(ctx context.Context, _ *payload.Empty) (res *payload.Info_Index_Count, err error) {
	return &payload.Info_Index_Count{
		Stored:      uint32(s.ngt.Len()),
		Uncommitted: uint32(s.ngt.InsertVCacheLen()),
		Indexing:    s.ngt.IsIndexing(),
	}, nil
}
