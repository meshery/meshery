package dataloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// AddonBatchLoader batches and caches addon lookups per cluster for a single request.
type AddonBatchLoader struct {
	ctx            context.Context
	providerGetter func(ctx context.Context) models.GenericPersisterGetter
	mu             sync.Mutex
	cache          map[string][]*model.AddonList
	queue          map[string][]chan addonBatchResult
	timer          *time.Timer
	batchWait      time.Duration
}

type addonBatchResult struct {
	data []*model.AddonList
	err  error
}

// newAddonBatchLoader constructs a new loader instance scoped to the request context.
func newAddonBatchLoader(ctx context.Context, providerGetter func(ctx context.Context) models.GenericPersisterGetter) *AddonBatchLoader {
	return &AddonBatchLoader{
		ctx:            ctx,
		providerGetter: providerGetter,
		cache:          make(map[string][]*model.AddonList),
		queue:          make(map[string][]chan addonBatchResult),
		batchWait:      2 * time.Millisecond,
	}
}

// Load loads addons for a single cluster ID, batching concurrent requests within the batch window.
func (l *AddonBatchLoader) Load(clusterID string, selectors []model.MeshType) ([]*model.AddonList, error) {
	if clusterID == "" {
		return nil, nil
	}

	l.mu.Lock()
	if cached, ok := l.cache[clusterID]; ok {
		l.mu.Unlock()
		return cached, nil
	}

	resultChan := make(chan addonBatchResult, 1)
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

// LoadMany loads addons for multiple cluster IDs, leveraging batching and in-memory cache.
func (l *AddonBatchLoader) LoadMany(clusterIDs []string, selectors []model.MeshType) (map[string][]*model.AddonList, error) {
	if len(clusterIDs) == 0 {
		return make(map[string][]*model.AddonList), nil
	}

	type indexedResult struct {
		cid  string
		data []*model.AddonList
		err  error
	}

	resultsChan := make(chan indexedResult, len(clusterIDs))
	var wg sync.WaitGroup

	for _, cid := range clusterIDs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			addons, err := l.Load(id, selectors)
			resultsChan <- indexedResult{cid: id, data: addons, err: err}
		}(cid)
	}

	wg.Wait()
	close(resultsChan)

	resultMap := make(map[string][]*model.AddonList)
	for res := range resultsChan {
		if res.err != nil {
			return nil, res.err
		}
		resultMap[res.cid] = res.data
	}

	return resultMap, nil
}

// dispatch processes all queued cluster IDs by executing a batched query for all mesh types.
func (l *AddonBatchLoader) dispatch() {
	l.mu.Lock()
	queue := l.queue
	l.queue = make(map[string][]chan addonBatchResult)
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
		var res addonBatchResult
		if err != nil {
			res = addonBatchResult{err: err}
		} else {
			data := resultsByCluster[cid]
			l.cache[cid] = data
			res = addonBatchResult{data: data}
		}
		for _, ch := range chans {
			ch <- res
		}
	}
}

// safeBatch wraps BatchGetAddonsState with panic recovery to protect the timer goroutine.
func (l *AddonBatchLoader) safeBatch(selectors []model.MeshType, provider models.GenericPersisterGetter, clusterIDs []string) (res map[string][]*model.AddonList, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("addon batch loader panic: %v", r)
		}
	}()
	return model.BatchGetAddonsState(l.ctx, selectors, provider, clusterIDs)
}
