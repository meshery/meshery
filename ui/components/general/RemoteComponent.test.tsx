import React from 'react';
import { act, fireEvent, render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import RemoteComponent from './RemoteComponent';

const { notify, routerPush, remoteState } = vi.hoisted(() => ({
  notify: vi.fn(),
  routerPush: vi.fn(),
  remoteState: { loading: false, err: null as Error | null, component: null as unknown },
}));

vi.mock('@paciolan/remote-component', () => ({
  createUseRemoteComponent: () => () => [
    remoteState.loading,
    remoteState.err,
    remoteState.component,
  ],
  getDependencies: () => ({}),
  createRequires: () => () => ({}),
}));

vi.mock('next/router', () => ({
  useRouter: () => ({ push: routerPush }),
}));

vi.mock('@/utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify }),
}));

// Pulls in RTK Query and the support webhook, neither of which this file tests.
vi.mock('../TroubleshootingComponent', () => ({
  default: () => null,
}));

const url = { url: '/api/provider/extension/plugin.js' };
const loadError = new Error('HTTP Error Response: 404 Not Found');

describe('RemoteComponent', () => {
  beforeEach(() => {
    notify.mockClear();
    routerPush.mockClear();
    remoteState.loading = false;
    remoteState.err = null;
    remoteState.component = () => <div>extension content</div>;
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('stays on the page after a load failure until the user acts', () => {
    vi.useFakeTimers();
    remoteState.err = loadError;

    render(<RemoteComponent url={url} />);
    act(() => {
      vi.advanceTimersByTime(60_000);
    });

    expect(routerPush).not.toHaveBeenCalled();
    expect(screen.getByRole('alert')).toHaveTextContent('404 Not Found');
  });

  it('goes to the dashboard when the user asks to', () => {
    remoteState.err = loadError;

    render(<RemoteComponent url={url} />);
    fireEvent.click(screen.getByRole('button', { name: 'Go to dashboard' }));

    expect(routerPush).toHaveBeenCalledTimes(1);
    expect(routerPush).toHaveBeenCalledWith('/');
  });

  it('reports a load failure once across re-renders', () => {
    remoteState.err = loadError;

    const { rerender } = render(<RemoteComponent url={url} />);
    rerender(<RemoteComponent url={url} />);

    expect(notify).toHaveBeenCalledTimes(1);
    expect(notify).toHaveBeenCalledWith({
      message: 'Failed to load extension: Error: HTTP Error Response: 404 Not Found',
      event_type: expect.objectContaining({ type: 'error' }),
    });
  });

  it('stays silent when a circular loader fails', () => {
    remoteState.err = loadError;

    const { container } = render(<RemoteComponent url={url} loaderType="circular" />);

    expect(container).toBeEmptyDOMElement();
    expect(notify).not.toHaveBeenCalled();
  });

  it('renders the extension without reporting anything when it loads', () => {
    render(<RemoteComponent url={url} />);

    expect(screen.getByText('extension content')).toBeInTheDocument();
    expect(screen.queryByRole('alert')).not.toBeInTheDocument();
    expect(notify).not.toHaveBeenCalled();
  });
});
