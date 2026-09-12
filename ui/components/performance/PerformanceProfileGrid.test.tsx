import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  Grid2: ({ children, container, ...rest }: any) =>
    container ? (
      <div data-testid="grid-container">{children}</div>
    ) : (
      <div data-testid="grid-item">{children}</div>
    ),
  Pagination: ({ count, onChange }: any) => (
    <div data-testid="pagination" data-count={count}>
      <button type="button" onClick={() => onChange(null, 2)}>
        page-2
      </button>
    </div>
  ),
}));

vi.mock('./style', () => ({
  PaginationWrapper: ({ children }: any) => <div data-testid="pagination-wrapper">{children}</div>,
}));

vi.mock('./PerformanceCard', () => ({
  default: ({ profile, handleEdit, handleDelete, handleProfile, handleRunTest }: any) => (
    <div data-testid="perf-card" data-profile={profile.id}>
      <button type="button" onClick={handleEdit}>
        edit
      </button>
      <button type="button" onClick={handleDelete}>
        delete
      </button>
      <button type="button" onClick={handleProfile}>
        profile
      </button>
      <button type="button" onClick={handleRunTest}>
        run
      </button>
    </div>
  ),
}));

import PerformanceProfileGrid from './PerformanceProfileGrid';

describe('PerformanceProfileGrid', () => {
  const baseProps = {
    profiles: [
      { id: 'p1', name: 'profile-1' },
      { id: 'p2', name: 'profile-2' },
    ],
    deleteHandler: vi.fn(),
    setProfileForModal: vi.fn(),
    testHandler: vi.fn(),
    pages: 3,
    setPage: vi.fn(),
  };

  it('renders a PerformanceCard per profile and a pagination control', () => {
    render(<PerformanceProfileGrid {...baseProps} />);
    expect(screen.getAllByTestId('perf-card')).toHaveLength(2);
    expect(screen.getByTestId('pagination')).toHaveAttribute('data-count', '3');
  });

  it('renders no pagination when there are no profiles', () => {
    render(<PerformanceProfileGrid {...baseProps} profiles={[]} />);
    expect(screen.queryByTestId('pagination')).not.toBeInTheDocument();
  });

  it('forwards delete to deleteHandler with the profile id', () => {
    render(<PerformanceProfileGrid {...baseProps} />);
    const cards = screen.getAllByTestId('perf-card');
    cards[0]
      .querySelector('button:nth-of-type(2)')
      ?.dispatchEvent(new MouseEvent('click', { bubbles: true }));
    expect(baseProps.deleteHandler).toHaveBeenCalledWith('p1');
  });

  it('forwards edit/profile callbacks to setProfileForModal', () => {
    const setProfileForModal = vi.fn();
    render(
      <PerformanceProfileGrid
        {...baseProps}
        profiles={[{ id: 'p1', name: 'profile-1' } as any]}
        setProfileForModal={setProfileForModal}
      />,
    );
    const card = screen.getByTestId('perf-card');
    (card.querySelector('button:nth-of-type(1)') as HTMLElement).click();
    (card.querySelector('button:nth-of-type(3)') as HTMLElement).click();
    expect(setProfileForModal).toHaveBeenCalledTimes(2);
  });

  it('invokes setPage with zero-based page number when pagination changes', () => {
    const setPage = vi.fn();
    render(<PerformanceProfileGrid {...baseProps} setPage={setPage} />);
    fireEvent.click(screen.getByRole('button', { name: 'page-2' }));
    expect(setPage).toHaveBeenCalledWith(1);
  });

  it('wraps the run callback with `{ ...profile, runTest: true }`', () => {
    const testHandler = vi.fn();
    render(
      <PerformanceProfileGrid
        {...baseProps}
        profiles={[{ id: 'p1', name: 'a' } as any]}
        testHandler={testHandler}
      />,
    );

    (screen.getByText('run') as HTMLElement).click();
    expect(testHandler).toHaveBeenCalledWith({ id: 'p1', name: 'a', runTest: true });
  });
});
