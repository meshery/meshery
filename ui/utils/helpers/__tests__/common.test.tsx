import { describe, expect, it, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import {
  closeButtonForSnackbarAction,
  errorHandlerGenerator,
  successHandlerGenerator,
} from '../common';
import { EVENT_TYPES } from '../../../lib/event-types';

describe('closeButtonForSnackbarAction', () => {
  it('returns a curried builder that renders a close icon button bound to the supplied key', () => {
    const closeSnackbar = vi.fn();
    const buildAction = closeButtonForSnackbarAction(closeSnackbar);
    const action = buildAction('snack-1');

    render(action as React.ReactElement);
    const btn = screen.getByRole('button', { name: /close/i });
    fireEvent.click(btn);

    expect(closeSnackbar).toHaveBeenCalledWith('snack-1');
  });

  it('produces independent handlers per snackbar key', () => {
    const closeSnackbar = vi.fn();
    const buildAction = closeButtonForSnackbarAction(closeSnackbar);

    const a = buildAction('a');
    const b = buildAction('b');

    const { unmount } = render(a as React.ReactElement);
    fireEvent.click(screen.getByRole('button'));
    unmount();

    render(b as React.ReactElement);
    fireEvent.click(screen.getByRole('button'));

    expect(closeSnackbar).toHaveBeenNthCalledWith(1, 'a');
    expect(closeSnackbar).toHaveBeenNthCalledWith(2, 'b');
  });

  it('safely handles non-function closeSnackbar on click', () => {
    const buildAction = closeButtonForSnackbarAction(undefined as any);
    const action = buildAction('snack-1');
    render(action as React.ReactElement);
    const btn = screen.getByRole('button', { name: /close/i });
    expect(() => fireEvent.click(btn)).not.toThrow();
  });
});

describe('successHandlerGenerator', () => {
  it('notifies with the message + stringified result when result is an object', () => {
    const notify = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved');
    handler({ name: 'demo' });
    expect(notify).toHaveBeenCalledWith({
      message: 'Saved',
      details: '{"name":"demo"}',
      event_type: EVENT_TYPES.SUCCESS,
    });
  });

  it('forwards the raw result to the optional callback before notifying', () => {
    const notify = vi.fn();
    const cb = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved', cb);
    const res = { name: 'demo' };
    handler(res);
    expect(cb).toHaveBeenCalledWith(res);
    expect(notify).toHaveBeenCalled();
  });

  it('emits the string as-is when result is already a string', () => {
    const notify = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved');
    handler('ok');
    expect(notify).toHaveBeenCalledWith({
      message: 'Saved',
      details: 'ok',
      event_type: EVENT_TYPES.SUCCESS,
    });
  });

  it('does nothing when the result is undefined', () => {
    const notify = vi.fn();
    const cb = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved', cb);
    handler(undefined);
    expect(notify).not.toHaveBeenCalled();
    expect(cb).not.toHaveBeenCalled();
  });

  it('safely handles circular object payloads without throwing', () => {
    const notify = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved');
    const circular: Record<string, unknown> = { name: 'demo' };
    circular.self = circular;
    expect(() => handler(circular)).not.toThrow();
    expect(notify).toHaveBeenCalledWith(
      expect.objectContaining({
        message: 'Saved',
        event_type: EVENT_TYPES.SUCCESS,
      }),
    );
  });

  it('handles non-function cb and notify safely', () => {
    const handlerWithBadCb = successHandlerGenerator(vi.fn(), 'Saved', 'not-a-fn' as any);
    expect(() => handlerWithBadCb('ok')).not.toThrow();

    const handlerWithBadNotify = successHandlerGenerator(null as any, 'Saved');
    expect(() => handlerWithBadNotify('ok')).not.toThrow();
  });

  it('handles circular object payloads whose toString also throws', () => {
    const notify = vi.fn();
    const handler = successHandlerGenerator(notify, 'Saved');
    const circularThrowing: Record<string, unknown> = {
      toString() {
        throw new Error('circular toString failed');
      },
    };
    circularThrowing.self = circularThrowing;

    expect(() => handler(circularThrowing)).not.toThrow();
    expect(notify).toHaveBeenCalledWith({
      message: 'Saved',
      details: '',
      event_type: EVENT_TYPES.SUCCESS,
    });
  });
});

describe('errorHandlerGenerator', () => {
  it('passes the raw error to the optional callback and notifies', () => {
    const notify = vi.fn();
    const cb = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed', cb);
    const err = new Error('boom');
    handler(err);
    expect(cb).toHaveBeenCalledWith(err);
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: 'Error: boom',
      event_type: EVENT_TYPES.ERROR,
    });
  });

  it('emits a string error verbatim', () => {
    const notify = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed');
    handler('network down');
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: 'network down',
      event_type: EVENT_TYPES.ERROR,
    });
  });

  it('stringifies non-Error objects via .toString()', () => {
    const notify = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed');
    handler({ toString: () => '[obj]' });
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: '[obj]',
      event_type: EVENT_TYPES.ERROR,
    });
  });

  it('runs the optional callback even on string errors', () => {
    const notify = vi.fn();
    const cb = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed', cb);
    handler('oops');
    expect(cb).toHaveBeenCalledWith('oops');
    expect(notify).toHaveBeenCalled();
  });

  it('safely handles null and undefined error values without throwing', () => {
    const notify = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed');

    expect(() => handler(null)).not.toThrow();
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: '',
      event_type: EVENT_TYPES.ERROR,
    });

    expect(() => handler(undefined)).not.toThrow();
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: '',
      event_type: EVENT_TYPES.ERROR,
    });
  });

  it('handles non-function cb and notify gracefully', () => {
    const handlerWithBadCb = errorHandlerGenerator(vi.fn(), 'Failed', 'not-a-fn' as any);
    expect(() => handlerWithBadCb('error')).not.toThrow();

    const handlerWithBadNotify = errorHandlerGenerator(null as any, 'Failed');
    expect(() => handlerWithBadNotify('error')).not.toThrow();
  });

  it('notifies safely even when error toString() throws', () => {
    const notify = vi.fn();
    const handler = errorHandlerGenerator(notify, 'Failed');
    const throwingErr = {
      message: 'inner error message',
      toString() {
        throw new Error('toString failed');
      },
    };

    expect(() => handler(throwingErr)).not.toThrow();
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: 'inner error message',
      event_type: EVENT_TYPES.ERROR,
    });

    const throwingNoMessage = {
      toString() {
        throw new Error('toString failed');
      },
    };
    expect(() => handler(throwingNoMessage)).not.toThrow();
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed',
      details: '',
      event_type: EVENT_TYPES.ERROR,
    });
  });
});
