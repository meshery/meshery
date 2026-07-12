import React, { useEffect } from 'react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SessionsProvider, { useSessions } from '../SessionsProvider';
import MinimizedSessions from '../MinimizedSessions';
import type { OpenSessionRequest } from '../SessionsProvider';

/**
 * Counts how many times a session pane has been *mounted*, not rendered.
 *
 * A session is a live shell on the far end of a WebSocket. Every one of the
 * panel's states — docked, floating, minimized — must therefore be a style change
 * on one element tree, never a different tree: React remounting the pane to move
 * it would tear the socket down and kill the shell. This counter is what holds
 * that invariant honest, and it is the reason the shell does not simply swap
 * between a docked box and Sistent's `Panel`.
 */
let mounts = 0;

vi.mock('../SessionPanel', () => ({
  default: () => {
    useEffect(() => {
      mounts += 1;
    }, []);
    return <div data-testid="session-pane" />;
  },
}));

const TARGET: OpenSessionRequest = {
  connectionId: 'conn-1',
  kind: 'terminal',
  target: { resource: 'pod', name: 'nginx', namespace: 'default' },
};

const OpenButton = () => {
  const { openSession } = useSessions();
  return (
    <button type="button" onClick={() => openSession(TARGET)}>
      open
    </button>
  );
};

const renderPanel = () =>
  render(
    <SessionsProvider>
      <OpenButton />
      {/* Mounted outside the provider in the app, as the Navigator is. */}
      <MinimizedSessions isDrawerCollapsed={false} />
    </SessionsProvider>,
  );

describe('SessionsProvider', () => {
  beforeEach(() => {
    mounts = 0;
  });

  it('docks a newly opened session, and can detach, minimize and restore it without remounting', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.getByTestId('session-pane')).toBeInTheDocument();

    // Docked is the default, so the control on offer is the one that detaches.
    await user.click(screen.getByRole('button', { name: 'Detach sessions panel' }));
    await user.click(screen.getByRole('button', { name: 'Dock sessions panel' }));
    await user.click(screen.getByRole('button', { name: 'Detach sessions panel' }));

    // Minimizing hands the panel off to the Navigator rather than closing it.
    await user.click(screen.getByRole('button', { name: 'Minimize sessions panel' }));
    const restore = screen.getByRole('button', { name: 'Restore 1 session' });

    await user.click(restore);
    expect(screen.queryByRole('button', { name: 'Restore 1 session' })).not.toBeInTheDocument();

    expect(screen.getByTestId('session-pane')).toBeInTheDocument();
    expect(mounts).toBe(1);
  });

  it('keeps a minimized session mounted, so its shell survives', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    await user.click(screen.getByRole('button', { name: 'Minimize sessions panel' }));

    // Hidden from assistive tech and from view, but still in the tree.
    expect(screen.getByTestId('session-pane')).toBeInTheDocument();
    expect(mounts).toBe(1);
  });

  it('shows no restore control until a session is minimized', async () => {
    const user = userEvent.setup();
    renderPanel();

    expect(screen.queryByRole('button', { name: /Restore/ })).not.toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.queryByRole('button', { name: /Restore/ })).not.toBeInTheDocument();
  });

  it('closing the panel tears every session down', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    await user.click(screen.getByRole('button', { name: 'Close all sessions' }));

    expect(screen.queryByTestId('session-pane')).not.toBeInTheDocument();
  });
});
