import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import React from 'react';

// ---------------------------------------------------------------------------
// `ability.tsx` is a React component module that orchestrates organization
// loading, ability fetching, and slow-session UX. We isolate its rendering
// from its heavy dependency graph by mocking the hooks/components it imports.
// ---------------------------------------------------------------------------

const updateMock = vi.fn();
vi.mock('../../utils/can', () => ({
  ability: { update: updateMock },
}));

const useGetUserKeysMock = vi.fn();
vi.mock('../userKeys', () => ({
  useGetUserKeysQuery: (...args: unknown[]) => useGetUserKeysMock(...args),
}));

const useGetProviderCapabilitiesMock = vi.fn();
const useGetSelectedOrganizationMock = vi.fn();
const useUpdateSelectedOrganizationMutationMock = vi.fn();
vi.mock('../user', () => ({
  useGetProviderCapabilitiesQuery: (...args: unknown[]) => useGetProviderCapabilitiesMock(...args),
  useGetSelectedOrganization: (...args: unknown[]) => useGetSelectedOrganizationMock(...args),
  useUpdateSelectedOrganizationMutation: (...args: unknown[]) =>
    useUpdateSelectedOrganizationMutationMock(...args),
}));

vi.mock('@/components/general/ErrorPage', () => ({
  default: ({ message }: { message?: string }) =>
    React.createElement('div', { 'data-testid': 'error-page' }, message ?? 'error'),
}));
vi.mock('@/components/general/error-404', () => ({
  default: () => React.createElement('div', { 'data-testid': 'error-404' }, '404'),
}));
vi.mock('@/components/shared/LoadingState/DynamicFullscreenLoader', () => ({
  DynamicFullScreenLoader: ({
    isLoading,
    children,
    message,
  }: {
    isLoading: boolean;
    children: React.ReactNode;
    message?: string;
  }) =>
    React.createElement(
      'div',
      { 'data-testid': 'loader', 'data-loading': isLoading, 'data-message': message ?? '' },
      children,
    ),
}));

import { render, screen } from '@testing-library/react';

describe('ability — useGetUserAbilities', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('maps key list to {action, subject} pairs and surfaces other RTK-Query state', async () => {
    useGetUserKeysMock.mockReturnValue({
      data: {
        keys: [
          { id: 'create', function: 'My Org' },
          { id: 'delete', function: 'TEAM' },
        ],
      },
      isLoading: false,
      isError: false,
    });
    const { useGetUserAbilities } = await import('../ability');
    const { renderHook } = await import('@testing-library/react');
    const { result } = renderHook(() => useGetUserAbilities({ id: 'org-1' }, false));
    expect(useGetUserKeysMock).toHaveBeenCalledWith({ orgId: 'org-1' }, { skip: false });
    expect(result.current.isLoading).toBe(false);
    expect(result.current.isError).toBe(false);
    expect(result.current.abilities).toEqual([
      { action: 'create', subject: 'my org' },
      { action: 'delete', subject: 'team' },
    ]);
  });

  it('returns empty abilities when data is missing', async () => {
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: true });
    const { useGetUserAbilities } = await import('../ability');
    const { renderHook } = await import('@testing-library/react');
    const { result } = renderHook(() => useGetUserAbilities({ id: 'x' }, true));
    expect(result.current.abilities).toEqual([]);
  });
});

describe('ability — useGetCurrentAbilities', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('skips the fetch when org is null', async () => {
    useGetUserKeysMock.mockReturnValue({ data: { keys: [] }, isLoading: false });
    const { useGetCurrentAbilities } = await import('../ability');
    const { renderHook } = await import('@testing-library/react');
    renderHook(() => useGetCurrentAbilities(null));
    expect(useGetUserKeysMock).toHaveBeenCalledWith({ orgId: undefined }, { skip: true });
  });

  it('calls ability.update with the mapped abilities when data is loaded', async () => {
    useGetUserKeysMock.mockReturnValue({
      data: { keys: [{ id: 'read', function: 'Connections' }] },
      isLoading: false,
    });
    const { useGetCurrentAbilities } = await import('../ability');
    const { renderHook } = await import('@testing-library/react');
    renderHook(() => useGetCurrentAbilities({ id: 'org-7' }));
    expect(updateMock).toHaveBeenCalledWith([{ action: 'read', subject: 'connections' }]);
  });
});

describe('ability — LoadSessionGuard component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // Provide sensible defaults
    useUpdateSelectedOrganizationMutationMock.mockReturnValue([vi.fn(), {}]);
    useGetProviderCapabilitiesMock.mockReturnValue({ data: { providerType: 'local' } });
    useGetUserKeysMock.mockReturnValue({ data: { keys: [] }, isLoading: false });
    vi.useFakeTimers({ shouldAdvanceTime: false });
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.restoreAllMocks();
  });

  it('renders children inside the loader while abilities are still loading', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-1' },
      didFallback: false,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: true });

    const { LoadSessionGuard } = await import('../ability');
    render(
      React.createElement(
        LoadSessionGuard,
        null,
        React.createElement('div', { 'data-testid': 'child' }, 'inside'),
      ),
    );
    expect(screen.getByTestId('loader').getAttribute('data-loading')).toBe('true');
    expect(screen.getByTestId('child')).toBeInTheDocument();
  });

  it('renders 404 + error message when org fetch errored', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: null,
      didFallback: false,
      isLoading: false,
      error: { message: 'boom' },
    });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(screen.getByTestId('error-404')).toBeInTheDocument();
    expect(screen.getByTestId('error-page')).toHaveTextContent(
      'Error occurred while fetching your current organization',
    );
  });

  it('renders 404 + error message when there is no organization and no fallback', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: null,
      didFallback: false,
      isLoading: false,
      error: null,
    });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(screen.getByTestId('error-404')).toBeInTheDocument();
    expect(screen.getByTestId('error-page')).toBeInTheDocument();
  });

  it('renders 404 + custom message when abilities loading errored with a message', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-9' },
      didFallback: false,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: { message: 'permission denied' },
    });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(screen.getByTestId('error-404')).toBeInTheDocument();
    expect(screen.getByTestId('error-page')).toHaveTextContent('permission denied');
  });

  it('renders 404 + default message when abilities loading errored without a message', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-9' },
      didFallback: false,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({
      data: undefined,
      isLoading: false,
      error: {},
    });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(screen.getByTestId('error-page')).toHaveTextContent(
      'An error occurred while fetching your organization permissions',
    );
  });

  it('calls updateSelectedOrganization when there is a fallback selection and no error', async () => {
    const updatePrefs = vi.fn();
    useUpdateSelectedOrganizationMutationMock.mockReturnValue([updatePrefs, {}]);
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'fallback-id' },
      didFallback: true,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: { keys: [] }, isLoading: false });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(updatePrefs).toHaveBeenCalledWith('fallback-id');
  });

  it('renders the slow-session message when loading exceeds the warning threshold', async () => {
    const { act } = await import('@testing-library/react');
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-1' },
      didFallback: false,
      isLoading: true,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: true });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    // Initially no slow message
    expect(screen.getByTestId('loader').getAttribute('data-message')).toBe('');

    // Advance fake timers past the 10s warning threshold, then flush React.
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10001);
    });
    expect(screen.getByTestId('loader').getAttribute('data-message')).toContain(
      'Still initializing',
    );
    expect(screen.getByTestId('loader').getAttribute('data-message')).toContain(
      'local provider session',
    );
  });

  it('formats pending session bootstrap message for organization-only', async () => {
    const { act } = await import('@testing-library/react');
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: null,
      didFallback: true,
      isLoading: true,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: false });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10001);
    });
    expect(screen.getByTestId('loader').getAttribute('data-message')).toContain(
      'your organization',
    );
  });

  it('formats pending bootstrap message for organization+permissions when both loading', async () => {
    const { act } = await import('@testing-library/react');
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-1' },
      didFallback: false,
      isLoading: true,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: true });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10001);
    });
    const message = screen.getByTestId('loader').getAttribute('data-message') || '';
    expect(message).toContain('your organization');
    expect(message).toContain('your organization permissions');
    expect(message).toContain(' and ');
  });

  it('uses generic provider session phrasing when providerType is not local', async () => {
    const { act } = await import('@testing-library/react');
    useGetProviderCapabilitiesMock.mockReturnValue({ data: { providerType: 'remote' } });
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'org-1' },
      didFallback: false,
      isLoading: true,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: undefined, isLoading: true });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10001);
    });
    expect(screen.getByTestId('loader').getAttribute('data-message')).toContain('your session');
    expect(screen.getByTestId('loader').getAttribute('data-message')).not.toContain(
      'local provider',
    );
  });

  it('does not call updateSelectedOrganization when there is no fallback selection', async () => {
    const updatePrefs = vi.fn();
    useUpdateSelectedOrganizationMutationMock.mockReturnValue([updatePrefs, {}]);
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'real' },
      didFallback: false,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: { keys: [] }, isLoading: false });

    const { LoadSessionGuard } = await import('../ability');
    render(React.createElement(LoadSessionGuard, null, React.createElement('div', null, 'x')));
    expect(updatePrefs).not.toHaveBeenCalled();
  });

  it('renders children when fully loaded and no errors', async () => {
    useGetSelectedOrganizationMock.mockReturnValue({
      selectedOrganization: { id: 'real' },
      didFallback: false,
      isLoading: false,
      error: null,
    });
    useGetUserKeysMock.mockReturnValue({ data: { keys: [] }, isLoading: false });

    const { LoadSessionGuard } = await import('../ability');
    render(
      React.createElement(
        LoadSessionGuard,
        null,
        React.createElement('div', { 'data-testid': 'app-child' }, 'inside'),
      ),
    );
    expect(screen.getByTestId('app-child')).toBeInTheDocument();
    expect(screen.getByTestId('loader').getAttribute('data-loading')).toBe('false');
  });
});
