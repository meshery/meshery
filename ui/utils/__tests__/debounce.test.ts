import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import debounce from '../debounce';

describe('debounce', () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('delays invocation by the default timeout (500ms)', () => {
    const fn = vi.fn();
    const debounced = debounce(fn);

    debounced();
    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(499);
    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(1);
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it('delays invocation by the provided timeout', () => {
    const fn = vi.fn();
    const debounced = debounce(fn, 1000);

    debounced();
    vi.advanceTimersByTime(500);
    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(500);
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it('resets the timer when called repeatedly within the timeout window', () => {
    const fn = vi.fn();
    const debounced = debounce(fn, 300);

    debounced();
    vi.advanceTimersByTime(200);
    debounced();
    vi.advanceTimersByTime(200);
    debounced();
    vi.advanceTimersByTime(299);
    expect(fn).not.toHaveBeenCalled();

    vi.advanceTimersByTime(1);
    expect(fn).toHaveBeenCalledTimes(1);
  });

  it('forwards the most recent arguments to the underlying function', () => {
    const fn = vi.fn();
    const debounced = debounce(fn, 100);

    debounced('first', 1);
    debounced('second', 2);
    vi.advanceTimersByTime(100);

    expect(fn).toHaveBeenCalledTimes(1);
    expect(fn).toHaveBeenCalledWith('second', 2);
  });

  it('handles no arguments', () => {
    const fn = vi.fn();
    const debounced = debounce(fn, 50);

    debounced();
    vi.advanceTimersByTime(50);

    expect(fn).toHaveBeenCalledWith();
  });

  it('returns separate timer state for separate debounced functions', () => {
    const fn1 = vi.fn();
    const fn2 = vi.fn();
    const d1 = debounce(fn1, 100);
    const d2 = debounce(fn2, 200);

    d1('a');
    d2('b');

    vi.advanceTimersByTime(100);
    expect(fn1).toHaveBeenCalledWith('a');
    expect(fn2).not.toHaveBeenCalled();

    vi.advanceTimersByTime(100);
    expect(fn2).toHaveBeenCalledWith('b');
  });

  it('can be invoked again after the timer has fired', () => {
    const fn = vi.fn();
    const debounced = debounce(fn, 100);

    debounced('one');
    vi.advanceTimersByTime(100);
    expect(fn).toHaveBeenCalledTimes(1);

    debounced('two');
    vi.advanceTimersByTime(100);
    expect(fn).toHaveBeenCalledTimes(2);
    expect(fn).toHaveBeenLastCalledWith('two');
  });
});
