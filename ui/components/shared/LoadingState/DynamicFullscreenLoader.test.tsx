import React from 'react';
import { render, screen } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { DynamicFullScreenLoader } from './DynamicFullscreenLoader';

describe('DynamicFullScreenLoader', () => {
  let show: ReturnType<typeof vi.fn>;
  let hide: ReturnType<typeof vi.fn>;
  let setMessage: ReturnType<typeof vi.fn>;
  let resetMessage: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    show = vi.fn();
    hide = vi.fn();
    setMessage = vi.fn();
    resetMessage = vi.fn();
    window.Loader = { show, hide, setMessage, resetMessage };
  });

  afterEach(() => {
    delete window.Loader;
  });

  it('renders nothing while loading', () => {
    const { container } = render(
      <DynamicFullScreenLoader isLoading={true}>
        <div data-testid="children">should not show</div>
      </DynamicFullScreenLoader>,
    );

    expect(container.textContent).toBe('');
    expect(screen.queryByTestId('children')).not.toBeInTheDocument();
  });

  it('renders children when not loading', () => {
    render(
      <DynamicFullScreenLoader isLoading={false}>
        <div data-testid="children">Visible</div>
      </DynamicFullScreenLoader>,
    );

    expect(screen.getByTestId('children')).toHaveTextContent('Visible');
  });

  it('calls show and setMessage when loading with a message', () => {
    render(
      <DynamicFullScreenLoader isLoading={true} message="Working...">
        <div>nope</div>
      </DynamicFullScreenLoader>,
    );

    expect(show).toHaveBeenCalled();
    expect(setMessage).toHaveBeenCalledWith('Working...');
  });

  it('calls show and resetMessage when loading without a message', () => {
    render(
      <DynamicFullScreenLoader isLoading={true}>
        <div>nope</div>
      </DynamicFullScreenLoader>,
    );

    expect(show).toHaveBeenCalled();
    expect(resetMessage).toHaveBeenCalled();
    expect(setMessage).not.toHaveBeenCalled();
  });

  it('calls hide and resetMessage when not loading', () => {
    render(
      <DynamicFullScreenLoader isLoading={false} message="x">
        <div>visible</div>
      </DynamicFullScreenLoader>,
    );

    expect(hide).toHaveBeenCalled();
    expect(resetMessage).toHaveBeenCalled();
  });

  it('handles missing window.Loader gracefully', () => {
    delete window.Loader;
    expect(() =>
      render(
        <DynamicFullScreenLoader isLoading={false}>
          <div data-testid="children">Visible</div>
        </DynamicFullScreenLoader>,
      ),
    ).not.toThrow();

    expect(screen.getByTestId('children')).toHaveTextContent('Visible');
  });
});
