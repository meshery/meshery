package dataloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// ControlPlaneBatchLoader batches and caches control plane lookups per cluster for a single request.
type ControlPlaneBatchLoader struct {
	ctx            context.Context
	providerGetter func(ctx context.Context) models.GenericPersisterGetter
	mu             sync.Mutex
	cache          map[string][]*model.ControlPlane
	queue          map[string][]chan batchResult
	timer          *time.Timer
	batchWait      time.Duration
}

type batchResult struct {
	data []*model.ControlPlane
	err  error
}

// newControlPlaneBatchLoader constructs a new loader instance scoped to the request context.
func newControlPlaneBatchLoader(ctx context.Context, providerGetter func(ctx context.Context) models.GenericPersisterGetter) *ControlPlaneBatchLoader {
	return &ControlPlaneBatchLoader{
		ctx:            ctx,
		providerGetter: providerGetter,
		cache:          make(map[string][]*model.ControlPlane),
		queue:          make(map[string][]chan batchResult),
		batchWait:      2 * time.Millisecond,
	}
}

// Load loads control planes for a single cluster ID, batching concurrent requests within the batch window.
func (l *ControlPlaneBatchLoader) Load(clusterID string, selectors []model.MeshType) ([]*model.ControlPlane, error) {
	if clusterID == "" {
		return nil, nil
	}

	l.mu.Lock()
	if cached, ok := l.cache[clusterID]; ok {
		l.mu.Unlock()
		return cached, nil
	}

	resultChan := make(chan batchResult, 1)
	l.queue[clusterID] = append(l.queue[clusterID], resultChan)

	if l.timer == nil {
		l.timer = time.AfterFunc(l.batchWait, func() {
			l.dispatch()
		})
	}
	l.mu.Unlock()

	select {
	case res := <-resultChan:
		return res.data, res.err
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	}
}

// LoadMany loads control planes for multiple cluster IDs, leveraging batching and in-memory cache.
func (l *ControlPlaneBatchLoader) LoadMany(clusterIDs []string, selectors []model.MeshType) (map[string][]*model.ControlPlane, error) {
	if len(clusterIDs) == 0 {
		return make(map[string][]*model.ControlPlane), nil
	}

	type indexedResult struct {
		cid  string
		data []*model.ControlPlane
		err  error
	}

	resultsChan := make(chan indexedResult, len(clusterIDs))
	var wg sync.WaitGroup

	for _, cid := range clusterIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			cps, err := l.Load(id, selectors)
			resultsChan <- indexedResult{cid: id, data: cps, err: err}
		}(cid)
	}

	wg.Wait()
	close(resultsChan)

	resultMap := make(map[string][]*model.ControlPlane)
	for res := range resultsChan {
		if res.err != nil {
			return nil, res.err
		}
		resultMap[res.cid] = res.data
	}

	return resultMap, nil
}

// dispatch processes all queued cluster IDs by executing a batched query for all mesh types.
func (l *ControlPlaneBatchLoader) dispatch() {
	l.mu.Lock()
	queue := l.queue
	l.queue = make(map[string][]chan batchResult)
	l.timer = nil
	l.mu.Unlock()

	if len(queue) == 0 {
		return
	}

	clusterIDs := make([]string, 0, len(queue))
	for cid := range queue {
		clusterIDs = append(clusterIDs, cid)
	}

	var provider models.GenericPersisterGetter
	if l.providerGetter != nil {
		provider = l.providerGetter(l.ctx)
	}

	// Always load every mesh type: the batch window merges Load calls, and the
	// timer closure only captures the first caller's selectors, so a per-call
	// selector set would silently serve the wrong results to later callers
	// (and the per-cluster cache is selector-blind too). The resolver already
	// filters by selector name.
	selectors := model.AllMeshType

	resultsByCluster, err := l.safeBatch(selectors, provider, clusterIDs)

	l.mu.Lock()
	defer l.mu.Unlock()

	for cid, chans := range queue {
		var res batchResult
		if err != nil {
			res = batchResult{err: err}
		} else {
			data := resultsByCluster[cid]
			l.cache[cid] = data
			res = batchResult{data: data}
		}
		for _, ch := range chans {
			ch <- res
		}
	}
}

// safeBatch wraps BatchGetControlPlaneState with panic recovery to protect the timer goroutine.
func (l *ControlPlaneBatchLoader) safeBatch(selectors []model.MeshType, provider models.GenericPersisterGetter, clusterIDs []string) (res map[string][]*model.ControlPlane, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("control plane batch loader panic: %v", r)
		}
	}()
	return model.BatchGetControlPlaneState(l.ctx, selectors, provider, clusterIDs)
}
