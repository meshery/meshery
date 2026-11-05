import React from 'react';
import { act, render, renderHook, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { DynamicComponentProvider, useDynamicComponent } from '../dynamicContext';

describe('DynamicComponentProvider / useDynamicComponent', () => {
  it('throws when the hook is used without a provider', () => {
    // Silence the React error logging during the intentional error.
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    expect(() => renderHook(() => useDynamicComponent())).toThrow(
      'useDynamicComponent must be used within a DynamicComponentProvider',
    );

    errorSpy.mockRestore();
  });

  it('exposes a setComponent function that updates DynamicComponent state', async () => {
    const SetterButton = () => {
      const { DynamicComponent, setComponent } = useDynamicComponent();
      return (
        <div>
          <button
            type="button"
            onClick={() => setComponent(() => () => <span data-testid="dyn">hello</span>)}
          >
            install
          </button>
          {DynamicComponent ? <DynamicComponent /> : <span data-testid="empty">empty</span>}
        </div>
      );
    };

    render(
      <DynamicComponentProvider>
        <SetterButton />
      </DynamicComponentProvider>,
    );

    expect(screen.getByTestId('empty')).toBeInTheDocument();
    await userEvent.click(screen.getByRole('button', { name: 'install' }));
    expect(screen.getByTestId('dyn')).toHaveTextContent('hello');
  });

  it('initial DynamicComponent is null', () => {
    let captured: ReturnType<typeof useDynamicComponent> | null = null;
    const Probe = () => {
      captured = useDynamicComponent();
      return null;
    };
    render(
      <DynamicComponentProvider>
        <Probe />
      </DynamicComponentProvider>,
    );
    expect(captured!.DynamicComponent).toBeNull();
    expect(typeof captured!.setComponent).toBe('function');
  });

  it('setComponent can also clear the dynamic component back to null', async () => {
    const Probe = () => {
      const { DynamicComponent, setComponent } = useDynamicComponent();
      return (
        <div>
          <button type="button" onClick={() => setComponent(() => () => <i>installed</i>)}>
            install
          </button>
          <button type="button" onClick={() => setComponent(null)}>
            clear
          </button>
          {DynamicComponent ? <DynamicComponent /> : <span data-testid="empty">empty</span>}
        </div>
      );
    };

    render(
      <DynamicComponentProvider>
        <Probe />
      </DynamicComponentProvider>,
    );

    await userEvent.click(screen.getByRole('button', { name: 'install' }));
    expect(screen.getByText('installed')).toBeInTheDocument();

    await act(async () => {
      await userEvent.click(screen.getByRole('button', { name: 'clear' }));
    });

    expect(screen.getByTestId('empty')).toBeInTheDocument();
  });
});
