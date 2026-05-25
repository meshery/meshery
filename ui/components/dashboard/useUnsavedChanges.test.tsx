import React from 'react';
import { act, render } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import useUnsavedChanges from './useUnsavedChanges';

type RouteHandler = (url: string) => void;

const onMock = vi.fn<(event: string, handler: RouteHandler) => void>();
const offMock = vi.fn<(event: string, handler: RouteHandler) => void>();
const emitMock = vi.fn();
const pushMock = vi.fn();

const eventListeners: Record<string, EventListenerOrEventListenerObject[]> = {};

vi.mock('next/router', () => ({
  useRouter: () => ({
    events: { on: onMock, off: offMock, emit: emitMock },
    push: pushMock,
  }),
}));

vi.mock('@/utils/hooks', () => ({
  useEventListener: (event: string, handler: EventListenerOrEventListenerObject) => {
    eventListeners[event] = eventListeners[event] || [];
    eventListeners[event].push(handler);
  },
}));

interface HookResult {
  showModal: boolean;
  hasUnsavedChanges: boolean;
  confirmNavigation: () => void;
  cancelNavigation: () => void;
}

const HookHarness: React.FC<{
  isEditMode: boolean;
  dashboardLayout: unknown;
  savedLayout: unknown;
  onChange: (val: HookResult) => void;
}> = ({ isEditMode, dashboardLayout, savedLayout, onChange }) => {
  const result = useUnsavedChanges({ isEditMode, dashboardLayout, savedLayout });
  React.useEffect(() => {
    onChange(result);
  });
  return null;
};

describe('useUnsavedChanges', () => {
  beforeEach(() => {
    onMock.mockReset();
    offMock.mockReset();
    emitMock.mockReset();
    pushMock.mockReset();
    Object.keys(eventListeners).forEach((key) => delete eventListeners[key]);
  });

  it('reports unsaved changes when in edit mode and layouts differ', () => {
    let captured: HookResult | undefined;
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={(v) => (captured = v)}
      />,
    );
    expect(captured?.hasUnsavedChanges).toBe(true);
  });

  it('does not report unsaved changes outside of edit mode', () => {
    let captured: HookResult | undefined;
    render(
      <HookHarness
        isEditMode={false}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={(v) => (captured = v)}
      />,
    );
    expect(captured?.hasUnsavedChanges).toBe(false);
  });

  it('opens modal and throws when navigation starts with pending changes', () => {
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={() => {}}
      />,
    );
    expect(onMock).toHaveBeenCalledWith('routeChangeStart', expect.any(Function));
    const handler = onMock.mock.calls[0][1] as RouteHandler;
    expect(() => handler('/new')).toThrow(/Navigation blocked/);
    expect(emitMock).toHaveBeenCalledWith('routeChangeError');
  });

  it('confirmNavigation pushes the pending url and clears state', () => {
    let captured: HookResult | undefined;
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={(v) => (captured = v)}
      />,
    );
    const routeHandler = onMock.mock.calls[0][1] as RouteHandler;
    expect(() => routeHandler('/destination')).toThrow();

    act(() => {
      captured?.confirmNavigation();
    });
    expect(pushMock).toHaveBeenCalledWith('/destination');
  });

  it('cancelNavigation clears pending url without pushing', () => {
    let captured: HookResult | undefined;
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={(v) => (captured = v)}
      />,
    );
    const routeHandler = onMock.mock.calls[0][1] as RouteHandler;
    expect(() => routeHandler('/destination')).toThrow();

    act(() => {
      captured?.cancelNavigation();
    });
    expect(pushMock).not.toHaveBeenCalled();
  });

  it('prevents the browser unload prompt when there are unsaved changes', () => {
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 2 }}
        onChange={() => {}}
      />,
    );

    const handlers = eventListeners['beforeunload'];
    expect(handlers).toBeDefined();

    const event = {
      preventDefault: vi.fn(),
      returnValue: '',
    } as unknown as BeforeUnloadEvent;
    (handlers[0] as (e: BeforeUnloadEvent) => void)(event);

    expect(event.preventDefault).toHaveBeenCalled();
  });

  it('does not block navigation when layouts match', () => {
    render(
      <HookHarness
        isEditMode={true}
        dashboardLayout={{ a: 1 }}
        savedLayout={{ a: 1 }}
        onChange={() => {}}
      />,
    );
    const handler = onMock.mock.calls[0][1] as RouteHandler;
    expect(() => handler('/anywhere')).not.toThrow();
  });
});
