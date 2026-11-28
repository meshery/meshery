import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

// Avoid loading the full app store (which pulls @meshery/schemas / the RTT
// Query API) into this unit test. ThemeTogglerCore wraps its inner component
// in ProviderStoreWrapper, which reads `store` from `@/store`.
vi.mock('@/store', () => ({
  store: {
    dispatch: vi.fn(),
    getState: vi.fn(() => ({})),
    subscribe: vi.fn(() => vi.fn()),
  },
}));

// Only `useMediaQuery` is used by hooks.tsx from @sistent/sistent.
vi.mock('@sistent/sistent', () => ({
  useMediaQuery: vi.fn(() => false),
}));

const mockUseGetUserPrefQuery = vi.fn();
const mockUpdateUserPrefWithContext = vi.fn();

vi.mock('@/rtk-query/user', () => ({
  useGetUserPrefQuery: (...args: any[]) => mockUseGetUserPrefQuery(...args),
  useUpdateUserPrefWithContextMutation: () => [mockUpdateUserPrefWithContext],
}));

import { ThemeTogglerCore } from './hooks';

const TestComponent = ({ mode, toggleTheme }: any) => (
  <div data-testid="theme-toggler">
    <span data-testid="mode">{mode}</span>
    <button data-testid="toggle" onClick={toggleTheme}>
      Toggle
    </button>
  </div>
);

const renderToggler = () => render(<ThemeTogglerCore Component={TestComponent} />);

beforeEach(() => {
  vi.clearAllMocks();
  mockUseGetUserPrefQuery.mockReturnValue({
    data: { remoteProviderPreferences: { theme: 'light' } },
    isSuccess: true,
    isError: false,
    isLoading: false,
    error: null,
  });
  mockUpdateUserPrefWithContext.mockImplementation(() => ({}));
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('ThemeTogglerCore.toggleTheme', () => {
  it('sends only the theme field when another preference update is in flight', () => {
    // Simulate an in-flight preference update whose optimistic value
    // (catalogContent: true) sits in the getUserPref cache but is unconfirmed.
    mockUseGetUserPrefQuery.mockReturnValue({
      data: {
        remoteProviderPreferences: { theme: 'light', catalogContent: true },
      },
      isSuccess: true,
      isError: false,
      isLoading: false,
      error: null,
    });
    const captured: any[] = [];
    mockUpdateUserPrefWithContext.mockImplementation((arg: any) => {
      captured.push(arg);
      return {};
    });

    renderToggler();

    fireEvent.click(screen.getByTestId('toggle'));

    expect(mockUpdateUserPrefWithContext).toHaveBeenCalledTimes(1);
    expect(captured[0]).toEqual({
      body: { remoteProviderPreferences: { theme: 'dark' } },
    });
    expect(captured[0].body).not.toHaveProperty('catalogContent');
    expect(JSON.stringify(captured[0].body)).not.toContain('catalogContent');
  });

  it('does not persist unconfirmed nested preferences on a theme toggle', () => {
    mockUseGetUserPrefQuery.mockReturnValue({
      data: {
        remoteProviderPreferences: {
          theme: 'light',
          usersExtensionPreferences: { catalogContent: true },
        },
      },
      isSuccess: true,
      isError: false,
      isLoading: false,
      error: null,
    });
    const captured: any[] = [];
    mockUpdateUserPrefWithContext.mockImplementation((arg: any) => {
      captured.push(arg);
      return {};
    });

    renderToggler();

    fireEvent.click(screen.getByTestId('toggle'));

    const body = captured[0].body;
    expect(body).toEqual({ remoteProviderPreferences: { theme: 'dark' } });
    expect(body.usersExtensionPreferences).toBeUndefined();
    expect(JSON.stringify(body)).not.toContain('catalogContent');
  });

  it('persists only the theme on a successful theme change', () => {
    const captured: any[] = [];
    mockUpdateUserPrefWithContext.mockImplementation((arg: any) => {
      captured.push(arg);
      return {};
    });

    renderToggler();

    fireEvent.click(screen.getByTestId('toggle'));

    expect(mockUpdateUserPrefWithContext).toHaveBeenCalledTimes(1);
    expect(captured[0]).toEqual({
      body: { remoteProviderPreferences: { theme: 'dark' } },
    });
    expect(screen.getByTestId('mode')).toHaveTextContent('dark');
  });

  it('toggles back from dark to light', () => {
    mockUseGetUserPrefQuery.mockReturnValue({
      data: { remoteProviderPreferences: { theme: 'dark' } },
      isSuccess: true,
      isError: false,
      isLoading: false,
      error: null,
    });
    const captured: any[] = [];
    mockUpdateUserPrefWithContext.mockImplementation((arg: any) => {
      captured.push(arg);
      return {};
    });

    renderToggler();

    expect(screen.getByTestId('mode')).toHaveTextContent('dark');
    fireEvent.click(screen.getByTestId('toggle'));

    expect(captured[0]).toEqual({
      body: { remoteProviderPreferences: { theme: 'light' } },
    });
    expect(screen.getByTestId('mode')).toHaveTextContent('light');
  });

  it('defaults to light theme when no preferences are cached', () => {
    mockUseGetUserPrefQuery.mockReturnValue({
      data: undefined,
      isSuccess: true,
      isError: false,
      isLoading: false,
      error: null,
    });
    const captured: any[] = [];
    mockUpdateUserPrefWithContext.mockImplementation((arg: any) => {
      captured.push(arg);
      return {};
    });

    renderToggler();

    expect(screen.getByTestId('mode')).toHaveTextContent('light');
    fireEvent.click(screen.getByTestId('toggle'));

    expect(captured[0]).toEqual({
      body: { remoteProviderPreferences: { theme: 'dark' } },
    });
  });
});
