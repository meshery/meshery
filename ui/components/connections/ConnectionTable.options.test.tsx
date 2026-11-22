import React from 'react';
import { renderHook } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('@sistent/sistent', () => ({
  useTheme: () => ({ palette: { common: { white: '#fff' }, error: { dark: '#900' } } }),
  Button: ({ children }: any) => <button type="button">{children}</button>,
  Grid2: ({ children }: any) => <div>{children}</div>,
  Table: ({ children }: any) => <table>{children}</table>,
  TableCell: ({ children }: any) => <td>{children}</td>,
  TableRow: ({ children }: any) => <tr>{children}</tr>,
  DeleteIcon: () => <svg />,
}));

vi.mock('./styles', () => ({
  ContentContainer: ({ children }: any) => <div>{children}</div>,
  InnerTableContainer: ({ children }: any) => <div>{children}</div>,
}));

vi.mock('../../css/icons.styles', () => ({ iconMedium: {} }));
vi.mock('./metadata', () => ({ default: () => null }));

import { useConnectionTableOptions } from './ConnectionTable.options';

const build = () => {
  const setPage = vi.fn();
  const setPageSizeAndResetPage = vi.fn();
  const { result } = renderHook(() =>
    useConnectionTableOptions({
      totalCount: 30,
      page: 2,
      pageSize: 10,
      setPage,
      setPageSizeAndResetPage,
      sortOrder: 'createdAt desc',
      setSortOrder: vi.fn(),
      rowsExpanded: [],
      setRowsExpanded: vi.fn(),
      columns: [{ name: 'name' }],
      filteredConnections: [],
      meshsyncControllerState: null,
      expansionFlags: { current: {} as any },
      handleDeleteConnections: vi.fn(),
    }),
  );
  return { options: result.current as any, setPage, setPageSizeAndResetPage };
};

describe('useConnectionTableOptions', () => {
  it('changeRowsPerPage sets the page size and resets the page in one call', () => {
    const { options, setPage, setPageSizeAndResetPage } = build();

    options.onTableChange('changeRowsPerPage', { rowsPerPage: 25, page: 2 });

    expect(setPageSizeAndResetPage).toHaveBeenCalledTimes(1);
    expect(setPageSizeAndResetPage).toHaveBeenCalledWith(25);
    // A separate setPage call would issue a second router.replace from stale
    // query state and drop the new page size.
    expect(setPage).not.toHaveBeenCalled();
  });

  it('changePage still updates only the page', () => {
    const { options, setPage, setPageSizeAndResetPage } = build();

    options.onTableChange('changePage', { page: 1 });

    expect(setPage).toHaveBeenCalledWith(1);
    expect(setPageSizeAndResetPage).not.toHaveBeenCalled();
  });
});
