import React from 'react';
import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const { useGetDatabaseSummaryQuery, tableProps } = vi.hoisted(() => ({
  useGetDatabaseSummaryQuery: vi.fn((_args: Record<string, unknown>) => ({
    data: { tables: [], totalTables: 40 },
    refetch: vi.fn(),
  })),
  tableProps: { current: null as any },
}));

vi.mock('@/rtk-query/system', () => ({ useGetDatabaseSummaryQuery }));
vi.mock('@/graphql/queries/ResetDatabaseQuery', () => ({ default: vi.fn() }));
vi.mock('@/store/slices/mesheryUi', () => ({ updateProgress: vi.fn() }));
vi.mock('@/assets/styles/general/tool.styles', () => ({
  ToolWrapper: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('../utils/hooks/useNotification', () => ({
  useNotification: () => ({ notify: vi.fn() }),
}));
// Run debounced handlers synchronously so the test doesn't depend on timers.
vi.mock('../utils/debounce', () => ({ default: (fn: any) => fn }));
vi.mock('../utils/custom-search', () => ({
  default: ({ onSearch }: any) => (
    <button data-testid="search-meshsync" onClick={() => onSearch('meshsync')}>
      search
    </button>
  ),
}));
vi.mock('@sistent/sistent', () => ({
  Button: ({ children, onClick }: any) => <button onClick={onClick}>{children}</button>,
  Typography: ({ children }: any) => <span>{children}</span>,
  useTheme: () => ({ palette: { error: { dark: '#000' } } }),
  PROMPT_VARIANTS: { DANGER: 'danger' },
  ResponsiveDataTable: (props: any) => {
    tableProps.current = props;
    return <div data-testid="table" />;
  },
}));

import DatabaseSummary from './DatabaseSummary';

const lastQueryArgs = () => useGetDatabaseSummaryQuery.mock.calls.at(-1)![0];

describe('DatabaseSummary', () => {
  beforeEach(() => {
    useGetDatabaseSummaryQuery.mockClear();
    tableProps.current = null;
  });

  it('goes back to the first page when searching from a later page', async () => {
    render(<DatabaseSummary promptRef={{ current: null }} />);

    act(() => tableProps.current.options.onChangePage(2));
    expect(lastQueryArgs()).toMatchObject({ page: 2, search: '' });

    await userEvent.click(screen.getByTestId('search-meshsync'));
    expect(lastQueryArgs()).toMatchObject({ page: 0, search: 'meshsync' });
  });

  it('only lets the name column be sorted, and sorts by the column that changed', () => {
    render(<DatabaseSummary promptRef={{ current: null }} />);

    const sortable = Object.fromEntries(
      tableProps.current.columns.map((c: any) => [c.name, c.options.sort]),
    );
    expect(sortable).toEqual({ name: true, count: false });

    act(() => tableProps.current.options.onColumnSortChange('name', 'desc'));
    expect(lastQueryArgs()).toMatchObject({ order: 'name desc' });
  });
});
