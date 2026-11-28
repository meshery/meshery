import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
  subscribeToControllersStatus,
  watchControllersStatus,
} from '../controllersStatusSubscription';

// jsdom has no EventSource. This fake records every instance so a test can
// drive messages and errors the way the browser would.
class FakeEventSource {
  static CONNECTING = 0;
  static OPEN = 1;
  static CLOSED = 2;
  static instances: FakeEventSource[] = [];

  url: string;
  readyState = FakeEventSource.CONNECTING;
  onmessage: ((message: { data: string }) => void) | null = null;
  onerror: ((error: Event) => void) | null = null;
  close = vi.fn(() => {
    this.readyState = FakeEventSource.CLOSED;
  });

  constructor(url: string) {
    this.url = url;
    FakeEventSource.instances.push(this);
  }

  emit(data: string) {
    this.onmessage?.({ data });
  }

  // `closed` mirrors the browser: CONNECTING while it retries a transient drop,
  // CLOSED once it gives up (non-2xx status or wrong content type).
  fail(closed: boolean) {
    this.readyState = closed ? FakeEventSource.CLOSED : FakeEventSource.CONNECTING;
    this.onerror?.(new Event('error'));
  }
}

const latest = () => FakeEventSource.instances[FakeEventSource.instances.length - 1];
const BASE = '/api/system/controllers/status/subscribe';
const ITEM = { connectionId: 'c1', controller: 'MESHSYNC', status: 'CONNECTED', version: 'v1' };

describe('controllers status subscription', () => {
  beforeEach(() => {
    FakeEventSource.instances = [];
    vi.stubGlobal('EventSource', FakeEventSource);
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });

  describe('subscribeToControllersStatus', () => {
    it('opens the stream with one connectionIds param per non-empty id', () => {
      subscribeToControllersStatus(['c1', '', 'c 2'], vi.fn());
      expect(latest().url).toBe(`${BASE}?connectionIds=c1&connectionIds=c+2`);
    });

    it('omits the query string when there are no ids', () => {
      subscribeToControllersStatus([], vi.fn());
      expect(latest().url).toBe(BASE);
    });

    it('passes each status array to onNext and coerces non-arrays to []', () => {
      const onNext = vi.fn();
      subscribeToControllersStatus(['c1'], onNext);

      latest().emit(JSON.stringify([ITEM]));
      latest().emit(JSON.stringify({ not: 'an array' }));

      expect(onNext.mock.calls).toEqual([[[ITEM]], [[]]]);
    });

    it('logs and skips a frame that is not JSON', () => {
      const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
      const onNext = vi.fn();
      subscribeToControllersStatus(['c1'], onNext);

      latest().emit('not json');

      expect(onNext).not.toHaveBeenCalled();
      expect(errSpy).toHaveBeenCalled();
    });

    it('stays quiet while the browser retries a transient drop', () => {
      const onError = vi.fn();
      subscribeToControllersStatus(['c1'], vi.fn(), onError);

      latest().fail(false);
      vi.advanceTimersByTime(10_000);

      expect(onError).not.toHaveBeenCalled();
    });

    it('reports a permanent close once, after the throttle delay', () => {
      const onError = vi.fn();
      subscribeToControllersStatus(['c1'], vi.fn(), onError);

      latest().fail(true);
      latest().fail(true);
      vi.advanceTimersByTime(2999);
      expect(onError).not.toHaveBeenCalled();
      vi.advanceTimersByTime(1);

      expect(onError).toHaveBeenCalledTimes(1);
    });

    it('closes the stream and cancels a pending error report on dispose', () => {
      const onError = vi.fn();
      const subscription = subscribeToControllersStatus(['c1'], vi.fn(), onError);

      latest().fail(true);
      subscription.dispose();
      vi.advanceTimersByTime(10_000);

      expect(latest().close).toHaveBeenCalled();
      expect(onError).not.toHaveBeenCalled();
    });
  });

  describe('watchControllersStatus', () => {
    // A permanent close is reported after the 3s throttle; the reopen then
    // waits the backoff delay for that attempt.
    const NOTIFY_MS = 3000;
    const closeLatest = () => latest().fail(true);

    it('reopens the same stream after a permanent close', () => {
      const onNext = vi.fn();
      watchControllersStatus(['c1'], onNext);
      const first = latest();

      closeLatest();
      vi.advanceTimersByTime(NOTIFY_MS + 1000);

      expect(FakeEventSource.instances).toHaveLength(2);
      const second = latest();
      expect(second.url).toBe(first.url);
      expect(first.close).toHaveBeenCalled();

      second.emit(JSON.stringify([ITEM]));
      expect(onNext).toHaveBeenCalledWith([ITEM]);
    });

    it('does not reopen on a transient drop the browser is retrying', () => {
      watchControllersStatus(['c1'], vi.fn());

      latest().fail(false);
      vi.advanceTimersByTime(60_000);

      expect(FakeEventSource.instances).toHaveLength(1);
    });

    it('backs off 1s, 2s, 4s ... between reopens, capped at 30s', () => {
      watchControllersStatus(['c1'], vi.fn());
      const expected = [1000, 2000, 4000, 8000, 16_000, 30_000, 30_000];

      expected.forEach((delay, i) => {
        closeLatest();
        vi.advanceTimersByTime(NOTIFY_MS + delay - 1);
        expect(FakeEventSource.instances).toHaveLength(i + 1);
        vi.advanceTimersByTime(1);
        expect(FakeEventSource.instances).toHaveLength(i + 2);
      });
    });

    it('resets the backoff once the reopened stream delivers a frame', () => {
      watchControllersStatus(['c1'], vi.fn());
      for (const delay of [1000, 2000, 4000]) {
        closeLatest();
        vi.advanceTimersByTime(NOTIFY_MS + delay);
      }

      latest().emit(JSON.stringify([ITEM]));
      closeLatest();
      vi.advanceTimersByTime(NOTIFY_MS + 1000);

      expect(FakeEventSource.instances).toHaveLength(5);
    });

    it('stops reopening once disposed before the close is reported', () => {
      const subscription = watchControllersStatus(['c1'], vi.fn());

      closeLatest();
      subscription.dispose();
      vi.advanceTimersByTime(60_000);

      expect(FakeEventSource.instances).toHaveLength(1);
    });

    it('stops reopening once disposed during the backoff wait', () => {
      const subscription = watchControllersStatus(['c1'], vi.fn());

      closeLatest();
      vi.advanceTimersByTime(NOTIFY_MS + 500);
      subscription.dispose();
      vi.advanceTimersByTime(60_000);

      expect(FakeEventSource.instances).toHaveLength(1);
    });

    it('closes the reopened stream, not just the first one, on dispose', () => {
      const subscription = watchControllersStatus(['c1'], vi.fn());
      closeLatest();
      vi.advanceTimersByTime(NOTIFY_MS + 1000);
      const reopened = latest();

      subscription.dispose();

      expect(reopened.close).toHaveBeenCalled();
      reopened.fail(true);
      vi.advanceTimersByTime(60_000);
      expect(FakeEventSource.instances).toHaveLength(2);
    });
  });
});
