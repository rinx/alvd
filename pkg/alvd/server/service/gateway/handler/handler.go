package handler

import (
	"context"
	"fmt"
	"io"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rinx/alvd/internal/errors"
	"github.com/rinx/alvd/internal/net/grpc/codes"
	"github.com/rinx/alvd/internal/net/grpc/status"
	"github.com/rinx/alvd/pkg/alvd/server/service/manager"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
)

const (
	defaultTimeout = 3 * time.Second
)

type server struct {
	manager    manager.Manager
	numReplica int
}

func New(man manager.Manager, replicas int) vald.Server {
	return &server{
		manager:    man,
		numReplica: replicas,
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
			return err
		}

		if res != nil && res.GetId() != "" {
			once.Do(func() {
				id = &payload.Object_ID{
					Id: res.GetId(),
				}
				cancel()
			})
		}

		return nil
	})
	if err != nil || id == nil || id.GetId() == "" {
		return nil, status.WrapWithNotFound(fmt.Sprintf("not found: %s", err), err)
	}

	return id, nil
}

func (s *server) Search(
	ctx context.Context,
	req *payload.Search_Request,
) (res *payload.Search_Response, err error) {
	cfg := req.GetConfig()

	timeout := getTimeout(cfg)
	num := int(cfg.GetNum())

	res = new(payload.Search_Response)
	res.Results = make([]*payload.Object_Distance, 0, s.manager.GetAgentCount()*num)
	dch := make(chan *payload.Object_Distance, cap(res.GetResults())/2)

	var maxDist uint32
	atomic.StoreUint32(&maxDist, math.Float32bits(math.MaxFloat32))

	ctx, cancel := context.WithTimeout(ctx, timeout)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer cancel()
		defer wg.Done()

		visited := sync.Map{}

		err = s.manager.Broadcast(ctx, func(ctx context.Context, client vald.Client) error {
			res, err := client.Search(ctx, req)
			if err != nil {
				return err
			}

			if res == nil || len(res.GetResults()) == 0 {
				return errors.New("not found")
			}

			for _, dist := range res.GetResults() {
				if dist == nil {
					continue
				}

				if dist.GetDistance() >= math.Float32frombits(atomic.LoadUint32(&maxDist)) {
					return nil
				}

				if _, already := visited.LoadOrStore(dist.GetId(), struct{}{}); !already {
					select {
					case <-ctx.Done():
						return nil
					case dch <- dist:
					}
				}
			}

			return nil
		})
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			close(dch)
			if num != 0 && len(res.GetResults()) > num {
				res.Results = res.Results[:num]
			}

			return res, nil
		case dist := <-dch:
			nres := len(res.GetResults())

			if nres >= num && dist.GetDistance() >= math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				continue
			}

			idx := -1

			for i := nres - 1; i >= 0; i-- {
				if res.GetResults()[i].GetDistance() <= dist.GetDistance() {
					idx = i
					break
				}
			}

			switch idx {
			case nres:
				res.Results = append(res.Results, dist)
			case -1:
				res.Results = append([]*payload.Object_Distance{dist}, res.Results...)
			default:
				res.Results = append(res.GetResults()[:idx+1], res.GetResults()[idx:]...)
				res.Results[idx+1] = dist
			}

			if num != 0 && nres+1 > num {
				res.Results = res.GetResults()[:num]
				nres--
			}

			if last := res.GetResults()[nres].GetDistance(); last < math.Float32frombits(atomic.LoadUint32(&maxDist)) {
				atomic.StoreUint32(&maxDist, math.Float32bits(last))
			}
		}
	}
}

func (s *server) SearchByID(
	ctx context.Context,
	req *payload.Search_IDRequest,
) (res *payload.Search_Response, err error) {
	vec, err := s.GetObject(ctx, &payload.Object_VectorRequest{
		Id: &payload.Object_ID{
			Id: req.GetId(),
		},
	})
	if err != nil {
		return nil, err
	}

	return s.Search(ctx, &payload.Search_Request{
		Vector: vec.GetVector(),
		Config: req.GetConfig(),
	})
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Search_StreamResponse{
							Payload: &payload.Search_StreamResponse_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Search_StreamResponse{
							Payload: &payload.Search_StreamResponse_Response{
								Response: res,
							},
						})
						mu.Unlock()
					}

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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Search_StreamResponse{
							Payload: &payload.Search_StreamResponse_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Search_StreamResponse{
							Payload: &payload.Search_StreamResponse_Response{
								Response: res,
							},
						})
						mu.Unlock()
					}

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
	mu := sync.Mutex{}
	ce = &payload.Object_Location{
		Uuid: req.GetVector().GetId(),
		Ips:  make([]string, 0, s.numReplica),
	}

	succeeded := uint32(0)

	err = s.manager.Range(ctx, 1, func(ctx context.Context, client vald.Client) error {
		if atomic.LoadUint32(&succeeded) >= uint32(s.numReplica) {
			return nil
		}

		loc, err := client.Insert(ctx, req)
		if err != nil {
			return err
		}

		atomic.AddUint32(&succeeded, 1)

		mu.Lock()
		defer mu.Unlock()

		ce.Ips = append(ce.GetIps(), loc.GetIps()...)
		ce.Name = loc.GetName()

		return nil
	})
	if err != nil && succeeded < uint32(s.numReplica) {
		return nil, err
	}

	return ce, nil
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Location{
								Location: res,
							},
						})
						mu.Unlock()
					}

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
	res, err = s.Remove(ctx, &payload.Remove_Request{
		Id: &payload.Object_ID{
			Id: req.GetVector().GetId(),
		},
		Config: &payload.Remove_Config{
			SkipStrictExistCheck: req.GetConfig().GetSkipStrictExistCheck(),
		},
	})
	if err != nil {
		return nil, err
	}

	res, err = s.Insert(ctx, &payload.Insert_Request{
		Vector: req.GetVector(),
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
			Filters:              req.GetConfig().GetFilters(),
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Location{
								Location: res,
							},
						})
						mu.Unlock()
					}

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

func (s *server) MultiUpdate(
	ctx context.Context,
	reqs *payload.Update_MultiRequest,
) (res *payload.Object_Locations, err error) {
	return res, nil
}

func (s *server) Upsert(
	ctx context.Context,
	req *payload.Upsert_Request,
) (res *payload.Object_Location, err error) {
	_, err = s.Exists(ctx, &payload.Object_ID{
		Id: req.GetVector().GetId(),
	})
	if err != nil {
		res, err = s.Insert(ctx, &payload.Insert_Request{
			Vector: req.GetVector(),
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
			},
		})
	} else {
		res, err = s.Update(ctx, &payload.Update_Request{
			Vector: req.GetVector(),
			Config: &payload.Update_Config{
				SkipStrictExistCheck: true,
				Filters:              req.GetConfig().GetFilters(),
			},
		})
	}

	if err != nil {
		return nil, err
	}

	return res, nil
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Location{
								Location: res,
							},
						})
						mu.Unlock()
					}

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
	locs = &payload.Object_Location{
		Uuid: req.GetId().GetId(),
	}

	mu := sync.Mutex{}

	err = s.manager.Broadcast(ctx, func(ctx context.Context, client vald.Client) error {
		loc, err := client.Remove(ctx, req)
		if err != nil {
			return nil
		}

		mu.Lock()
		defer mu.Unlock()

		locs.Ips = append(locs.GetIps(), loc.GetIps()...)
		locs.Name = loc.GetName()

		return nil
	})
	if err != nil || locs == nil || locs.GetUuid() == "" || locs.GetName() == "" {
		return nil, status.WrapWithNotFound(fmt.Sprintf("not found: %s", err), err)
	}

	return locs, nil
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Object_StreamLocation{
							Payload: &payload.Object_StreamLocation_Location{
								Location: res,
							},
						})
						mu.Unlock()
					}

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

func (s *server) MultiRemove(
	ctx context.Context,
	reqs *payload.Remove_MultiRequest,
) (locs *payload.Object_Locations, err error) {
	return locs, nil
}

func (s *server) GetObject(
	ctx context.Context,
	req *payload.Object_VectorRequest,
) (vec *payload.Object_Vector, err error) {
	ctx, cancel := context.WithCancel(ctx)

	var once sync.Once

	err = s.manager.Broadcast(ctx, func(ctx context.Context, client vald.Client) error {
		res, err := client.GetObject(ctx, req)
		if err != nil {
			return err
		}

		if res != nil && res.GetId() != "" && res.GetVector() != nil {
			once.Do(func() {
				vec = &payload.Object_Vector{
					Id:     res.GetId(),
					Vector: res.GetVector(),
				}
				cancel()
			})
		}

		return nil
	})
	if err != nil || vec == nil || vec.GetId() == "" || vec.GetVector() == nil {
		return nil, status.WrapWithNotFound(fmt.Sprintf("not found: %s", err), err)
	}

	return vec, nil
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
						errs = append(errs, err)
						emu.Unlock()

						st, ok := status.FromError(err)

						if !ok {
							st = status.New(codes.Internal, errors.Wrap(err, "failed to parse grpc status from error").Error())
							err = errors.Wrap(st.Err(), err.Error())
						}

						mu.Lock()
						err = stream.Send(&payload.Object_StreamVector{
							Payload: &payload.Object_StreamVector_Status{
								Status: st.Proto(),
							},
						})
						mu.Unlock()
					} else {
						mu.Lock()
						err = stream.Send(&payload.Object_StreamVector{
							Payload: &payload.Object_StreamVector_Vector{
								Vector: res,
							},
						})
						mu.Unlock()
					}

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

func getTimeout(cfg *payload.Search_Config) time.Duration {
	if to := cfg.GetTimeout(); to != 0 {
		return time.Duration(to)
	} else {
		return defaultTimeout
	}
}
