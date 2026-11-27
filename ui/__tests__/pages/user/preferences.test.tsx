import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { Keys } from '@meshery/schemas/permissions';

const { useHasPermission, useGetUserPrefWithContextQuery } = vi.hoisted(() => ({
  useHasPermission: vi.fn(),
  useGetUserPrefWithContextQuery: vi.fn(),
}));

vi.mock('react-redux', () => ({
  useSelector: (selector: (state: { ui: { selectedK8sContext: string } }) => unknown) =>
    selector({ ui: { selectedK8sContext: 'ctx-1' } }),
}));

vi.mock('@sistent/sistent', () => ({ useHasPermission }));

vi.mock('@/rtk-query/user', () => ({
  useGetUserPrefWithContextQuery,
}));

vi.mock('../../../components/general/MesheryPage', () => ({
  MesheryPage: ({ title, children }: { title: string; children: React.ReactNode }) => (
    <div data-testid="meshery-page" data-title={title}>
      {children}
    </div>
  ),
}));

vi.mock('@/components/user-preferences', () => ({
  default: ({
    anonymousStats,
    perfResultStats,
  }: {
    anonymousStats: boolean;
    perfResultStats: boolean;
  }) => (
    <div
      data-testid="user-preferences"
      data-anon={String(anonymousStats)}
      data-perf={String(perfResultStats)}
    />
  ),
}));

vi.mock('@/components/general/error-404', () => ({
  default: ({ permissionKey }: { permissionKey: unknown }) => (
    <div data-testid="default-error" data-permission-key={JSON.stringify(permissionKey)} />
  ),
}));

import UserPref from '../../../pages/user/preferences';

describe('UserPref page access gate', () => {
  beforeEach(() => {
    useHasPermission.mockReset();
    useGetUserPrefWithContextQuery.mockReset();
  });

  it('skips the context query, ignores cached preferences, and renders DefaultError inside MesheryPage when denied', () => {
    useHasPermission.mockReturnValue(false);
    useGetUserPrefWithContextQuery.mockReturnValue({
      data: { anonymousUsageStats: true, anonymousPerfResults: true },
    });

    render(<UserPref />);

    expect(useHasPermission).toHaveBeenCalledWith(Keys.ExtensibilityViewMesheryUserPreferences);
    expect(useGetUserPrefWithContextQuery).toHaveBeenCalledWith('ctx-1', { skip: true });
    expect(screen.getByTestId('meshery-page')).toHaveAttribute('data-title', 'User Preferences');
    expect(screen.getByTestId('default-error')).toHaveAttribute(
      'data-permission-key',
      JSON.stringify(Keys.ExtensibilityViewMesheryUserPreferences),
    );
    expect(screen.queryByTestId('user-preferences')).not.toBeInTheDocument();
  });

  it('preserves the MesheryPage shell without rendering preferences while stats are loading', () => {
    useHasPermission.mockReturnValue(true);
    useGetUserPrefWithContextQuery.mockReturnValue({ data: undefined });

    render(<UserPref />);

    expect(useGetUserPrefWithContextQuery).toHaveBeenCalledWith('ctx-1', { skip: false });
    expect(screen.getByTestId('meshery-page')).toHaveAttribute('data-title', 'User Preferences');
    expect(screen.queryByTestId('user-preferences')).not.toBeInTheDocument();
    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
  });

  it('renders UserPreferences inside MesheryPage when permitted and stats are loaded', () => {
    useHasPermission.mockReturnValue(true);
    useGetUserPrefWithContextQuery.mockReturnValue({
      data: { anonymousUsageStats: true, anonymousPerfResults: false },
    });

    render(<UserPref />);

    expect(useGetUserPrefWithContextQuery).toHaveBeenCalledWith('ctx-1', { skip: false });
    expect(screen.getByTestId('meshery-page')).toHaveAttribute('data-title', 'User Preferences');
    expect(screen.getByTestId('user-preferences')).toHaveAttribute('data-anon', 'true');
    expect(screen.getByTestId('user-preferences')).toHaveAttribute('data-perf', 'false');
    expect(screen.queryByTestId('default-error')).not.toBeInTheDocument();
  });

  it('replaces UserPreferences with DefaultError and skips query if permission is revoked mid-session', () => {
    useHasPermission.mockReturnValue(true);
    useGetUserPrefWithContextQuery.mockReturnValue({
      data: { anonymousUsageStats: true, anonymousPerfResults: true },
    });

    const { rerender } = render(<UserPref />);
    expect(screen.getByTestId('user-preferences')).toBeInTheDocument();

    useHasPermission.mockReturnValue(false);
    rerender(<UserPref />);

    expect(useGetUserPrefWithContextQuery).toHaveBeenLastCalledWith('ctx-1', { skip: true });
    expect(screen.getByTestId('default-error')).toBeInTheDocument();
    expect(screen.queryByTestId('user-preferences')).not.toBeInTheDocument();
  });
});
