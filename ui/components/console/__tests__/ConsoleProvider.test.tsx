import React, { useEffect } from 'react';
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ConsoleProvider, { useConsole } from '../ConsoleProvider';
import MinimizedConsoles from '../MinimizedConsoles';
import type { OpenConsoleRequest } from '../ConsoleProvider';

/**
 * Counts how many times a console pane has been *mounted*, not rendered.
 *
 * A console is a live shell on the far end of a WebSocket. Every one of the
 * panel's states — docked, floating, minimized — must therefore be a style change
 * on one element tree, never a different tree: React remounting the pane to move
 * it would tear the socket down and kill the shell. This counter is what holds
 * that invariant honest, and it is the reason the shell does not simply swap
 * between a docked box and Sistent's `Panel`.
 */
let mounts = 0;

vi.mock('../ConsolePanel', () => ({
  default: () => {
    useEffect(() => {
      mounts += 1;
    }, []);
    return <div data-testid="console-pane" />;
  },
}));

const TARGET: OpenConsoleRequest = {
  connectionId: 'conn-1',
  kind: 'terminal',
  target: { resource: 'pod', name: 'nginx', namespace: 'default' },
};

const OpenButton = () => {
  const { openConsole } = useConsole();
  return (
    <button type="button" onClick={() => openConsole(TARGET)}>
      open
    </button>
  );
};

const renderPanel = () =>
  render(
    <ConsoleProvider>
      <OpenButton />
      {/* Mounted outside the provider in the app, as the Navigator is. */}
      <MinimizedConsoles isDrawerCollapsed={false} />
    </ConsoleProvider>,
  );

describe('ConsoleProvider', () => {
  beforeEach(() => {
    mounts = 0;
  });

  it('docks a newly opened console, and can detach, minimize and restore it without remounting', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.getByTestId('console-pane')).toBeInTheDocument();

    // Docked is the default, so the control on offer is the one that detaches.
    await user.click(screen.getByRole('button', { name: 'Detach consoles panel' }));
    await user.click(screen.getByRole('button', { name: 'Dock consoles panel' }));
    await user.click(screen.getByRole('button', { name: 'Detach consoles panel' }));

    // Minimizing hands the panel off to the Navigator rather than closing it.
    await user.click(screen.getByRole('button', { name: 'Minimize consoles panel' }));
    const restore = screen.getByRole('button', { name: 'Restore 1 console' });

    await user.click(restore);
    expect(screen.queryByRole('button', { name: 'Restore 1 console' })).not.toBeInTheDocument();

    expect(screen.getByTestId('console-pane')).toBeInTheDocument();
    expect(mounts).toBe(1);
  });

  it('keeps a minimized console mounted, so its shell survives', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    await user.click(screen.getByRole('button', { name: 'Minimize consoles panel' }));

    // Hidden from assistive tech and from view, but still in the tree.
    expect(screen.getByTestId('console-pane')).toBeInTheDocument();
    expect(mounts).toBe(1);
  });

  it('shows no restore control until a console is minimized', async () => {
    const user = userEvent.setup();
    renderPanel();

    expect(screen.queryByRole('button', { name: /Restore/ })).not.toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'open' }));
    expect(screen.queryByRole('button', { name: /Restore/ })).not.toBeInTheDocument();
  });

  it('closing the panel tears every console down', async () => {
    const user = userEvent.setup();
    renderPanel();

    await user.click(screen.getByRole('button', { name: 'open' }));
    await user.click(screen.getByRole('button', { name: 'Close all consoles' }));

    expect(screen.queryByTestId('console-pane')).not.toBeInTheDocument();
  });
});
