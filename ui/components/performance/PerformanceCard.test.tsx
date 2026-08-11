import React from 'react';
import { render, screen } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const getUserByIdQuery = vi.fn();
const getProviderCapabilitiesQuery = vi.fn();

// Exercise the real useResourceOwner hook (single source of truth for the
// owner/cloud-link precedence) through mocked RTK queries, mirroring how the
// sibling FiltersCard test asserts avatar resolution and cloud-link gating.
vi.mock('@/rtk-query/user', () => ({
  useGetUserByIdQuery: (...args: unknown[]) => getUserByIdQuery(...args),
  useGetProviderCapabilitiesQuery: () => getProviderCapabilitiesQuery(),
}));

vi.mock('../../constants/endpoints', () => ({
  MESHERY_CLOUD_PROD: 'https://meshery.io',
}));

vi.mock('../../css/icons.styles', () => ({ iconMedium: {} }));

vi.mock('@meshery/schemas/permissions', () => ({
  Keys: new Proxy({}, { get: () => ({}) }),
}));

vi.mock('@/utils/hooks/useTestIDs', () => ({
  default: () => (id: string) => `test-${id}`,
}));

vi.mock('@/assets/icons', () => ({
  Delete: () => <svg data-testid="delete-icon" />,
}));

vi.mock('@sistent/sistent', () => ({
  CustomTooltip: ({ children, title }: any) => <div data-tip={title}>{children}</div>,
  Typography: ({ children }: any) => <span>{children}</span>,
  Avatar: ({ src }: any) => <div data-testid="avatar" data-src={src} />,
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  useTheme: () => ({
    palette: {
      mode: 'light',
      text: { disabled: 'disabled' },
    },
  }),
  Grid2: ({ children }: any) => <div>{children}</div>,
  IconButton: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  Link: ({ children }: any) => <a>{children}</a>,
  Table: ({ children }: any) => (
    <table>
      <tbody>{children}</tbody>
    </table>
  ),
  TableCell: ({ children }: any) => <td>{children}</td>,
  TableRow: ({ children }: any) => <tr>{children}</tr>,
  EditIcon: () => <svg data-testid="edit-icon" />,
}));

vi.mock('../general/FlipCard', () => ({
  default: ({ children }: any) => (
    <div data-testid="flip-card">
      {Array.isArray(children)
        ? children.map((c: any, i: number) => <div key={i}>{c}</div>)
        : children}
    </div>
  ),
}));

vi.mock('./PerformanceResults', () => ({
  default: () => <div data-testid="performance-results" />,
}));

vi.mock('./style', () => ({
  BottomPart: ({ children }: any) => <div>{children}</div>,
  CardButton: ({ children }: any) => <div>{children}</div>,
  ResultContainer: ({ children }: any) => <div>{children}</div>,
}));

import PerformanceCard from './PerformanceCard';

describe('PerformanceCard', () => {
  beforeEach(() => {
    getUserByIdQuery.mockReset();
    getUserByIdQuery.mockReturnValue({ data: { avatarUrl: 'https://a.io/u.png' } });
    getProviderCapabilitiesQuery.mockReset();
    getProviderCapabilitiesQuery.mockReturnValue({ data: { providerType: 'remote' } });
  });

  const baseProps = {
    profile: {
      id: 'p-1',
      owner: 'user-42',
      name: 'my perf profile',
      endpoints: ['http://x'],
      loadGenerators: ['fortio'],
      totalResults: 5,
    },
    handleDelete: vi.fn(),
    handleEdit: vi.fn(),
    handleRunTest: vi.fn(),
    handleProfile: vi.fn(),
    requestFullSize: vi.fn(),
    requestSizeRestore: vi.fn(),
  };

  it('resolves the owner avatar by id via useResourceOwner', () => {
    render(<PerformanceCard {...baseProps} />);
    expect(getUserByIdQuery).toHaveBeenCalledWith(
      'user-42',
      expect.objectContaining({ skip: false }),
    );
    expect(screen.getByTestId('avatar')).toHaveAttribute('data-src', 'https://a.io/u.png');
  });

  it('links the owner avatar on a remote provider', () => {
    render(<PerformanceCard {...baseProps} />);
    expect(screen.getByTestId('avatar').closest('a')).not.toBeNull();
  });

  it('does not link the owner avatar on the built-in local provider', () => {
    // The local provider's user has no Meshery Cloud profile page, so the link
    // would 404.
    getProviderCapabilitiesQuery.mockReturnValue({ data: { providerType: 'local' } });

    render(<PerformanceCard {...baseProps} />);

    expect(screen.getByTestId('avatar')).toBeInTheDocument();
    expect(screen.getByTestId('avatar').closest('a')).toBeNull();
  });
});
