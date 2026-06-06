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
});
