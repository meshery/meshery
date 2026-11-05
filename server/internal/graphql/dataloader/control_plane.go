package dataloader

import (
	"context"
	"sync"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// ControlPlaneBatchLoader batches and caches control plane lookups per cluster for a single request.
type ControlPlaneBatchLoader struct {
	ctx            context.Context
	providerGetter func(ctx context.Context) models.Provider
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
func newControlPlaneBatchLoader(ctx context.Context, providerGetter func(ctx context.Context) models.Provider) *ControlPlaneBatchLoader {
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
			l.dispatch(selectors)
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

func (l *ControlPlaneBatchLoader) dispatch(selectors []model.MeshType) {
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

	var provider models.Provider
	if l.providerGetter != nil {
		provider = l.providerGetter(l.ctx)
	}

	if len(selectors) == 0 {
		selectors = model.AllMeshType
	}

	resultsByCluster, err := model.BatchGetControlPlaneState(l.ctx, selectors, provider, clusterIDs)

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

