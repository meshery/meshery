package models

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// listenerBufferSize is the per-subscriber channel buffer. It absorbs short
// bursts of events published faster than a subscriber drains them; once it is
// full, Publish drops the newest event for that subscriber rather than blocking
// the publisher (see Publish).
const listenerBufferSize = 256

// replayBufferSize bounds how many recent ID-bearing events are retained per
// user for SSE Last-Event-ID replay across reconnects.
const replayBufferSize = 256

// defaultHistoryIdleTTL is how long an inactive user's replay buffer is kept
// after its last publish or subscribe activity before lazy eviction reclaims it.
const defaultHistoryIdleTTL = 15 * time.Minute

// defaultHistorySweepInterval bounds how often lazy history pruning scans
// Broadcast.history so Publish/SubscribeFrom stay O(1) in the common path.
const defaultHistorySweepInterval = time.Minute

// ReplayEvent pairs an event payload with the monotonic per-user SSE cursor
// assigned when it was published, so repeated event IDs remain distinguishable.
type ReplayEvent struct {
	Cursor string
	Data   interface{}
}

type replayEntry struct {
	id     string
	cursor string
	data   interface{}
}

type replayBuffer struct {
	mu           sync.Mutex
	entries      [replayBufferSize]replayEntry
	head         int
	count        int
	nextSeq      uint64
	lastActivity time.Time
	evicted      bool
}

func (rb *replayBuffer) append(eventID string, data interface{}, now time.Time) string {
	rb.nextSeq++
	cursor := fmt.Sprintf("%s:%d", eventID, rb.nextSeq)
	rb.entries[rb.head] = replayEntry{id: eventID, cursor: cursor, data: data}
	rb.head = (rb.head + 1) % replayBufferSize
	if rb.count < replayBufferSize {
		rb.count++
	}
	rb.lastActivity = now
	return cursor
}

func isValidLastEventID(s string) bool {
	if s == "" {
		return false
	}
	if _, err := uuid.FromString(s); err == nil {
		return true
	}
	if idx := strings.LastIndexByte(s, ':'); idx > 0 && idx < len(s)-1 {
		if _, err := uuid.FromString(s[:idx]); err == nil {
			if seq, err := strconv.ParseUint(s[idx+1:], 10, 64); err == nil && seq > 0 {
				return true
			}
		}
	}
	return false
}

// after returns a chronological snapshot of all ReplayEvents recorded strictly
// after lastEventID. If lastEventID matches a cursor (<uuid>:<seq>) or bare
// UUID in the buffer, it replays entries after that exact position. If
// lastEventID is a valid UUID/cursor that has been evicted from the buffer, it
// replays the full retained buffer as best effort. Empty or malformed IDs
// return nil. Caller must hold rb.mu.
func (rb *replayBuffer) after(lastEventID string) []ReplayEvent {
	if rb.count == 0 || !isValidLastEventID(lastEventID) {
		return nil
	}
	start := (rb.head - rb.count + replayBufferSize) % replayBufferSize
	for i := rb.count - 1; i >= 0; i-- {
		idx := (start + i) % replayBufferSize
		entry := rb.entries[idx]
		if entry.cursor == lastEventID || entry.id == lastEventID {
			remaining := rb.count - 1 - i
			if remaining == 0 {
				return nil
			}
			out := make([]ReplayEvent, 0, remaining)
			for j := i + 1; j < rb.count; j++ {
				e := rb.entries[(start+j)%replayBufferSize]
				out = append(out, ReplayEvent{Cursor: e.cursor, Data: e.data})
			}
			return out
		}
	}
	out := make([]ReplayEvent, 0, rb.count)
	for j := 0; j < rb.count; j++ {
		e := rb.entries[(start+j)%replayBufferSize]
		out = append(out, ReplayEvent{Cursor: e.cursor, Data: e.data})
	}
	return out
}

type listener struct {
	ch         chan interface{}
	withCursor bool
}

type clients struct {
	listeners []listener
	mu        *sync.Mutex
	evicted   bool
}

type Broadcast struct {
	clients              sync.Map
	history              sync.Map // core.Uuid -> *replayBuffer
	lastHistorySweepNs   atomic.Int64
	historyIdleTTL       time.Duration
	historySweepInterval time.Duration
	Name                 string
}

// Subscribe registers a new listener for id and returns its channel plus an
// idempotent unsubscribe func. The clients entry is stored behind a pointer so
// Subscribe, Publish and unsubscribe all share one mutex and one listeners
// slice - no stale copy can outlive a concurrent mutation. The channel is
// buffered (listenerBufferSize); the caller must drain it or unsubscribe,
// otherwise a full buffer makes Publish drop the newest events for this listener.
// When all listeners for id unsubscribe, the entry is pruned from the clients
// map to prevent memory leaks across disconnected sessions.
func (c *Broadcast) Subscribe(id core.Uuid) (chan interface{}, func()) {
	return c.subscribe(id, false)
}

func (c *Broadcast) subscribe(id core.Uuid, withCursor bool) (chan interface{}, func()) {
	ch := make(chan interface{}, listenerBufferSize)

	for {
		actual, ok := c.clients.Load(id)
		if !ok {
			actual, _ = c.clients.LoadOrStore(id, &clients{mu: &sync.Mutex{}})
		}
		cl := actual.(*clients)

		cl.mu.Lock()
		if cl.evicted {
			cl.mu.Unlock()
			continue
		}
		cl.listeners = append(cl.listeners, listener{ch: ch, withCursor: withCursor})
		cl.mu.Unlock()

		var once sync.Once
		unsubscribe := func() {
			once.Do(func() {
				cl.mu.Lock()
				defer cl.mu.Unlock()
				for i, l := range cl.listeners {
					if l.ch == ch {
						cl.listeners = append(cl.listeners[:i], cl.listeners[i+1:]...)
						close(ch)
						break
					}
				}
				if len(cl.listeners) == 0 {
					cl.evicted = true
					c.clients.CompareAndDelete(id, actual)
				}
			})
		}
		return ch, unsubscribe
	}
}

// SubscribeFrom atomically snapshots any buffered events strictly after
// lastEventID and registers a live listener channel for id. Live and replayed
// items delivered on the returned slice/channel are wrapped as ReplayEvent when
// they carry an event ID. Holding the per-user replay lock across both steps
// guarantees no ID-bearing event is missed or duplicated between the replay
// slice and the live channel.
func (c *Broadcast) SubscribeFrom(id core.Uuid, lastEventID string) ([]ReplayEvent, chan interface{}, func()) {
	if !isValidLastEventID(lastEventID) {
		ch, unsubscribe := c.subscribe(id, true)
		return nil, ch, unsubscribe
	}
	now := time.Now()
	c.maybePruneIdleHistory(now)
	rb := c.lockReplayBuffer(id, now)
	rb.lastActivity = now
	replay := rb.after(lastEventID)
	ch, unsubscribe := c.subscribe(id, true)
	rb.mu.Unlock()
	return replay, ch, unsubscribe
}

func (c *Broadcast) getHistoryIdleTTL() time.Duration {
	if c.historyIdleTTL > 0 {
		return c.historyIdleTTL
	}
	return defaultHistoryIdleTTL
}

func (c *Broadcast) getHistorySweepInterval() time.Duration {
	if c.historySweepInterval > 0 {
		return c.historySweepInterval
	}
	return defaultHistorySweepInterval
}

func (c *Broadcast) lockReplayBuffer(id core.Uuid, now time.Time) *replayBuffer {
	for {
		actual, ok := c.history.Load(id)
		if !ok {
			actual, _ = c.history.LoadOrStore(id, &replayBuffer{lastActivity: now})
		}
		rb := actual.(*replayBuffer)
		rb.mu.Lock()
		if rb.evicted {
			rb.mu.Unlock()
			continue
		}
		return rb
	}
}

func (c *Broadcast) maybePruneIdleHistory(now time.Time) {
	interval := c.getHistorySweepInterval()
	prev := c.lastHistorySweepNs.Load()
	nowNs := now.UnixNano()
	if prev != 0 && nowNs-prev < interval.Nanoseconds() {
		return
	}
	if !c.lastHistorySweepNs.CompareAndSwap(prev, nowNs) {
		return
	}
	ttl := c.getHistoryIdleTTL()
	c.history.Range(func(key, val interface{}) bool {
		if _, active := c.clients.Load(key); active {
			return true
		}
		rb := val.(*replayBuffer)
		rb.mu.Lock()
		if !rb.evicted && now.Sub(rb.lastActivity) > ttl {
			rb.evicted = true
			c.history.CompareAndDelete(key, val)
		}
		rb.mu.Unlock()
		return true
	})
}

func extractEventID(data interface{}) string {
	switch e := data.(type) {
	case *events.Event:
		if e != nil && e.ID != uuid.Nil {
			return e.ID.String()
		}
	case events.Event:
		if e.ID != uuid.Nil {
			return e.ID.String()
		}
	}
	return ""
}

// Publish records ID-bearing events in the per-user replay buffer and delivers
// data to every current listener for id. Sends are non-blocking and run under
// the same mutex that guards unsubscribe, which makes this safe on two fronts:
// a channel present in listeners under the lock is guaranteed still open (so
// the send can never panic on a closed channel), and a slow consumer whose
// buffer is full simply has the newest event dropped rather than blocking the
// publisher.
func (c *Broadcast) Publish(id core.Uuid, data interface{}) {
	if eventID := extractEventID(data); eventID != "" {
		now := time.Now()
		c.maybePruneIdleHistory(now)
		rb := c.lockReplayBuffer(id, now)
		cursor := rb.append(eventID, data, now)
		c.publishToClients(id, data, cursor)
		rb.mu.Unlock()
		return
	}
	c.publishToClients(id, data, "")
}

func (c *Broadcast) publishToClients(id core.Uuid, data interface{}, cursor string) {
	actual, ok := c.clients.Load(id)
	if !ok {
		return
	}
	cl, ok := actual.(*clients)
	if !ok {
		return
	}

	cl.mu.Lock()
	defer cl.mu.Unlock()
	if cl.evicted {
		return
	}
	for _, l := range cl.listeners {
		msg := data
		if l.withCursor && cursor != "" {
			msg = ReplayEvent{Cursor: cursor, Data: data}
		}
		select {
		case l.ch <- msg:
		default:
		}
	}
}

func NewBroadcaster(name string) *Broadcast {
	return &Broadcast{
		clients: sync.Map{},
		Name:    name,
	}
}
